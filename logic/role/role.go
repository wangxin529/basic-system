package role

import (
	"elevate-hub/conversion"
	meta2 "elevate-hub/db/meta"
	"elevate-hub/db/models"
	"elevate-hub/meta"
	"elevate-hub/svc"
	"gorm.io/gorm"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleLogic struct {
	svCtx *svc.ServiceContext
}

func NewRoleLogic(ctx *svc.ServiceContext) *RoleLogic {
	return &RoleLogic{
		svCtx: ctx,
	}
}

func (l *RoleLogic) AddRole(c *gin.Context) {
	var role AddRoleReq
	err := c.ShouldBindJSON(&role)
	if err != nil {
		log.Printf("add role failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "角色新增失败")
		return
	}

	err = l.svCtx.Mysql.RoleModel.Create(models.Role{
		Name:   role.Name,
		Status: role.Status,
	})
	if err != nil {
		log.Printf("add role in db create failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "角色新增失败请联系，系统管理员")
		return
	}
	meta.SuccessHandle(c, nil)
}

func (l *RoleLogic) UpdateRole(c *gin.Context) {
	var req UpdateRoleReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		meta.ErrHandleWithMsg(c, "角色修改参数错误")
		return
	}

	err = l.svCtx.Mysql.RoleModel.Update(meta2.UpdateOption{
		ID:   req.ID,
		Data: map[string]interface{}{"name": req.Name, "status": req.Status},
	})
	if err != nil {
		log.Printf("update role failed:%v", err)
		meta.ErrHandleWithMsg(c, "角色修改失败，联系管理源")
		return
	}

	// 修改权限

	auths, _, err := l.svCtx.Mysql.AuthRelationModel.List(meta2.ListOption{
		DisableCount: true,
		Page:         meta2.EmptyPage,
		Condition: map[string]interface{}{
			"role_id": req.ID,
			"app_id":  0,
			"user_id": 0,
			"api_id":  0,
		},
	})
	var (
		news    []*models.AuthRelation
		updates []*models.AuthRelation
		authMap = make(map[int64]*models.AuthRelation)
	)

	for _, auth := range auths {
		authMap[auth.MenuID] = auth
	}

	for key, val := range req.OperatePermission {
		if op, ok := authMap[key]; ok {

			updates = append(updates, op.SetButtons(val))
		} else {
			op = &models.AuthRelation{MenuID: key, RoleID: req.ID}
			news = append(news, op.SetButtons(val))
		}
	}

	if len(news) > 0 {
		// 新增
		err = l.svCtx.Mysql.AuthRelationModel.Creates(news)
		if err != nil {
			log.Printf("update role failed:%v", err)
			meta.ErrHandleWithMsg(c, err.Error())
			return
		}
	}

	if len(updates) > 0 {
		for _, update := range updates {
			err = l.svCtx.Mysql.DB.Transaction(func(tx *gorm.DB) error {
				err = l.svCtx.Mysql.AuthRelationModel.Copy(tx).Update(meta2.UpdateOption{
					Condition: map[string]interface{}{
						"role_id": update.RoleID,
						"menu_id": update.MenuID,
					},
					Data: map[string]interface{}{
						"button_permission": update.ButtonPermission,
					},
				})
				if err != nil {
					log.Printf("update role failed:%v", err)
					return err
				}
				return nil
			})

			if err != nil {
				log.Printf("update role failed:%v", err)
				meta.ErrHandleWithMsg(c, err.Error())
				return
			}
		}
	}
	meta.SuccessHandle(c, nil)
	return
}

func (l *RoleLogic) DelRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		meta.ErrHandleWithMsg(c, "指定需要删除的角色")
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		meta.ErrHandleWithMsg(c, "删除角色信息错误")
		return
	}
	err = l.svCtx.Mysql.RoleModel.Delete(meta2.DeleteOption{
		ID: idInt,
	})
	if err != nil {
		log.Printf("delete role failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}

	meta.SuccessHandle(c, nil)
	return
}

func (l *RoleLogic) ListRole(c *gin.Context) {
	var req ListRoleReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.Printf("list role failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "查询参数失败")
		return
	}
	res, count, err := l.svCtx.Mysql.RoleModel.List(meta2.ListOption{
		Page: &req.Page,
	})

	if err != nil {
		log.Printf("list role failed:%v", err.Error())
		meta.ErrHandle(c, err)
		return
	}

	meta.SuccessHandleAndTotal(c, conversion.SliceRoleDO2DTO(res, nil), count)
	return
}
func (l *RoleLogic) Role(c *gin.Context) {
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
	ars, _, err := l.svCtx.Mysql.AuthRelationModel.List(meta2.ListOption{
		Page: meta2.EmptyPage,
		Condition: map[string]interface{}{
			"role_id": id,
			"app_id":  0,
			"user_id": 0,
			"api_id":  0,
		},
	})
	if err != nil {
		meta.ErrHandle(c, err)
		return
	}
	res, err := l.svCtx.Mysql.RoleModel.Get(meta2.GetOption{
		ID: idInt,
	})
	if err != nil {
		log.Printf("get role failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}

	meta.SuccessHandle(c, conversion.RoleDO2DTO(res, ars))
	return
}
