package department

import (
	"elevate-hub/conversion"
	meta2 "elevate-hub/db/meta"
	"elevate-hub/db/models"
	"elevate-hub/meta"
	"elevate-hub/svc"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type DepartmentLogic struct {
	svCtx *svc.ServiceContext
}

func NewDepartmentLogic(ctx *svc.ServiceContext) *DepartmentLogic {
	return &DepartmentLogic{
		svCtx: ctx,
	}
}

func (l *DepartmentLogic) AddDepartment(c *gin.Context) {
	var req AddDepartmentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("add Department failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "部门新增失败")
		return
	}
	dept := models.Department{
		Name:   req.Name,
		Parent: req.Parent,
		Leader: req.Leader,
	}.GenerateCode()

	_, err = l.svCtx.Mysql.DepartmentModel.Get(meta2.GetOption{
		Condition: map[string]interface{}{
			"code": dept.Code,
		},
	})

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		meta.ErrHandleWithMsg(c, "当前菜单已存在,请勿重复添加")
		return
	}
	err = l.svCtx.Mysql.DepartmentModel.Create(dept)
	if err != nil {
		log.Printf("add Department in db create failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "部门新增失败请联系，系统管理员")
		return
	}
	meta.SuccessHandle(c, nil)
}

func (l *DepartmentLogic) UpdateDepartment(c *gin.Context) {
	var req UpdateDepartmentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("update Department failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "部门修改失败，部门信息错误")
		return
	}

	if req.ID == 0 {
		meta.ErrHandleWithMsg(c, "未知的部门信息")
		return
	}
	var dept = make(map[string]interface{})
	// 添加待改的参数
	if req.Name != "" {
		dept["name"] = req.Name
	}
	dept["parent"] = req.Parent
	if req.Leader != 0 {
		dept["leader"] = req.Leader
	}
	err = l.svCtx.Mysql.DepartmentModel.Update(meta2.UpdateOption{
		ID:   req.ID,
		Data: dept,
	})
	if err != nil {
		log.Printf("update Department failed:%v", err)
		meta.ErrHandleWithMsg(c, "部门数据修改失败")
		return
	}
	meta.SuccessHandle(c, nil)
}

func (l *DepartmentLogic) DelDepartment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		meta.ErrHandleWithMsg(c, "指定需要删除的数据")
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		meta.ErrHandleWithMsg(c, "参数错误")
		return
	}
	err = l.svCtx.Mysql.DepartmentModel.Delete(meta2.DeleteOption{
		ID: idInt,
	})
	if err != nil {
		log.Printf("delete deparment failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}

	meta.SuccessHandle(c, nil)
	return
}

func (l *DepartmentLogic) ListDepartment(c *gin.Context) {

	res, _, err := l.svCtx.Mysql.DepartmentModel.List(meta2.EmptyListOption)
	if err != nil {
		log.Printf("list department failed:%v", err.Error())
		meta.ErrHandle(c, err)
		return
	}
	depts := conversion.SliceDepartmentAggregateDO2DTO(res)
	dept := conversion.Department2Child(depts, nil)
	meta.SuccessHandle(c, dept.Children)
	return

}
func (l *DepartmentLogic) Department(c *gin.Context) {
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
	res, err := l.svCtx.Mysql.DepartmentModel.GetByID(idInt)
	if err != nil {
		log.Printf("get menu failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, conversion.DepartmentAggregateDO2DTO(res))
	return
}
