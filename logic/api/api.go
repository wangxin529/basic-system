package api

import (
	"elevate-hub/common/utils"
	"elevate-hub/conversion"
	meta2 "elevate-hub/db/meta"
	"elevate-hub/db/models"
	"elevate-hub/meta"
	"elevate-hub/svc"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type APILogic struct {
	svCtx *svc.ServiceContext
}

func NewAPILogic(ctx *svc.ServiceContext) *APILogic {
	return &APILogic{
		svCtx: ctx,
	}
}

func (l *APILogic) AddAPI(c *gin.Context) {
	var req AddAPIReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("add API failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "API新增失败")
		return
	}
	id, _ := utils.GetUserID(c)
	err = l.svCtx.Mysql.APIModel.Create(models.API{
		Name:    req.Name,
		Path:    req.Path,
		Method:  req.Method,
		Creator: id,
		//Key:      req.key,
		Status:   req.Status,
		Describe: req.Describe,
		APIType:  req.APIType,
	})
	if err != nil {
		log.Printf("add API in db create failed:%v", err)
		meta.ErrHandleWithMsg(c, "API新增失败请联系，系统管理员")
		return
	}
	meta.SuccessHandle(c, nil)
}

func (l *APILogic) UpdateAPI(c *gin.Context) {
	var req UpdateAPIReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("update API failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "API修改失败")
		return
	}

	err = l.svCtx.Mysql.APIModel.Update(meta2.UpdateOption{
		ID: req.ID,
		Data: map[string]interface{}{
			"name":     req.Name,
			"path":     req.Path,
			"method":   req.Method,
			"describe": req.Describe,
			"status":   req.Status,
		},
	})
	if err != nil {
		log.Printf("update API failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "API修改失败")
		return
	}
	meta.SuccessHandle(c, nil)
}

func (l *APILogic) DelAPI(c *gin.Context) {
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
	err = l.svCtx.Mysql.APIModel.Delete(meta2.DeleteOption{
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

func (l *APILogic) ListAPI(c *gin.Context) {

	var req ListAPIReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.Printf("list API failed:%v", err)
		meta.ErrHandleWithMsg(c, "查询参数失败")
		return
	}
	options := meta2.ListOption{
		Page: &req.Page,
	}
	if req.APIType != nil {
		options.AddToCondition("api_type", req.APIType)
	}

	res, count, err := l.svCtx.Mysql.APIModel.List(options)

	if err != nil {
		log.Printf("list API failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandleAndTotal(c, conversion.SliceAPIAggregateDO2DTO(res), count)
	return

}
func (l *APILogic) API(c *gin.Context) {
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
	res, err := l.svCtx.Mysql.APIModel.Get(meta2.GetOption{
		ID: idInt,
	})
	if err != nil {
		log.Printf("get api failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, conversion.APIDO2DTO(res))
	return
}
