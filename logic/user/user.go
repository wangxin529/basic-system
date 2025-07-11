package user

import (
	"elevate-hub/common/set"
	"elevate-hub/conversion"
	meta2 "elevate-hub/db/meta"
	"elevate-hub/db/models"
	"elevate-hub/meta"
	"elevate-hub/middleware/jwt"
	"elevate-hub/svc"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type UserLogic struct {
	svCtx *svc.ServiceContext
}

func NewUserLogic(ctx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		svCtx: ctx,
	}
}

func (l *UserLogic) AddUser(c *gin.Context) {
	var userDto AddUserReq
	err := c.ShouldBindJSON(&userDto)
	if err != nil {
		log.Printf("add user failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "用户新增失败")
		return
	}
	if userDto.Password != userDto.PasswordAgain {
		meta.ErrHandleWithMsg(c, "两次密码不同请确认")
		return
	}

	err = l.svCtx.Mysql.UserModel.Create(models.User{
		Username: userDto.Username,
		Password: jwt.MD5Password(userDto.Password),
		NickName: userDto.Nickname,
		Sex:      userDto.Sex,
		Email:    userDto.Email,
		Phone:    userDto.Phone,
	})
	if err != nil {
		log.Printf("add user in db create failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "用户新增失败请联系，系统管理员")
		return
	}
	meta.SuccessHandle(c, nil)
}

func (l *UserLogic) UpdateUser(c *gin.Context) {
	var userDto UpdateUserReq
	err := c.ShouldBindJSON(&userDto)
	if err != nil {
		meta.ErrHandle(c, err)
		return
	}
	err = l.updateUser(userDto)
	if err != nil {
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, nil)
}

func (l *UserLogic) updateUser(userDto UpdateUserReq) error {
	err := l.svCtx.Mysql.UserModel.Update(meta2.UpdateOption{
		ID: userDto.ID,
		Data: map[string]interface{}{
			"nickname":       userDto.Nickname,
			"sex":            userDto.Sex,
			"email":          userDto.Email,
			"phone":          userDto.Phone,
			"status":         userDto.Status,
			"supper_manager": userDto.IsSupper,
		},
	})
	if err != nil {
		return errors.Wrap(err, "用户数据修改失败")
	}
	var oldRoles []int64
	_, err = l.svCtx.Mysql.AuthRelationModel.ListAggregate(meta2.ListOption{
		DisableCount: true,
		Condition: map[string]interface{}{
			"user_id": userDto.ID,
		},
		Select: "role_id",
	}, &oldRoles)

	if err != nil {
		return errors.Wrap(err, "角色数据修改失败")
	}

	var (
		newRole = set.Array2Set[int64](userDto.Roles)
		oldRole = set.Array2Set[int64](oldRoles)
		addRole []int64
		delRole []int64
	)

	// 新增的
	for _, role := range userDto.Roles {
		if !oldRole.Contains(role) {
			addRole = append(addRole, role)
		}
	}

	// 删除的
	for _, role := range oldRoles {
		if !newRole.Contains(role) {
			delRole = append(delRole, role)
		}
	}

	err = l.svCtx.Mysql.AuthRelationModel.Delete(meta2.DeleteOption{
		Condition: map[string]interface{}{
			"user_id": userDto.ID,
		},
		InCondition: map[string]interface{}{
			"role_id": delRole,
		},
	})
	if err != nil {
		return errors.Wrap(err, "角色删除失败")
	}
	var addAuthRole = make([]*models.AuthRelation, 0, len(addRole))
	for _, role := range addRole {
		addAuthRole = append(addAuthRole, &models.AuthRelation{
			UserID: userDto.ID,
			RoleID: role,
		})
	}
	err = l.svCtx.Mysql.AuthRelationModel.Creates(addAuthRole)
	if err != nil {
		return errors.Wrap(err, "角色新转增失败")
	}
	return nil
}

func (l *UserLogic) DelUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		meta.ErrHandleWithMsg(c, "指定需要删除的用户")
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		meta.ErrHandleWithMsg(c, "删除用户信息错误")
		return
	}
	err = l.svCtx.Mysql.UserModel.Delete(meta2.DeleteOption{
		ID: idInt,
	})
	if err != nil {
		log.Printf("delete user failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}

	meta.SuccessHandle(c, nil)
	return

}

func (l *UserLogic) ListUser(c *gin.Context) {
	var req ListUserReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.Printf("list user failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "查询参数失败")
		return
	}
	res, count, err := l.svCtx.Mysql.UserModel.List(meta2.ListOption{
		Page: &req.Page,
	})

	if err != nil {
		log.Printf("list user failed:%v", err.Error())
		meta.ErrHandle(c, err)
		return
	}

	meta.SuccessHandleAndTotal(c, conversion.SliceUserAggregate2DTO(res), count)
	return
}
func (l *UserLogic) UserInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		meta.ErrHandleWithMsg(c, "参数错误")
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		meta.ErrHandleWithMsg(c, "参数错误")
		return
	}
	res, err := l.svCtx.Mysql.UserModel.GetByID(idInt)
	if err != nil {
		log.Printf("get user failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, conversion.UserAggregateDO2DTO(res))
	return
}

func (l *UserLogic) CurrentUser(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		meta.ErrHandleWithHttpCodeAndMsg(c, http.StatusUnauthorized, "请先登录")
		return
	}
	number, ok := userId.(json.Number)
	if !ok {
		meta.ErrHandleWithHttpCodeAndMsg(c, http.StatusUnauthorized, "请先登录")
		return
	}

	id, err := number.Int64()
	if err != nil {
		meta.ErrHandleWithMsg(c, "请重新登录")
		return
	}
	user, err := l.svCtx.Mysql.UserModel.Get(meta2.GetOption{
		ID: id,
	})
	if err != nil {
		meta.ErrHandleWithMsg(c, "当前用户查询失败，请联系管理员")
		return
	}
	meta.SuccessHandle(c, conversion.UserDO2DTO(user))
	return
}
