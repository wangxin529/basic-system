package auth

import (
	"elevate-hub/conversion"
	"elevate-hub/meta"
	"elevate-hub/svc"
	"github.com/gin-gonic/gin"
	"log"
)

type AuthLogic struct {
	svCtx *svc.ServiceContext
}

func NewAuthLogic(ctx *svc.ServiceContext) *AuthLogic {
	return &AuthLogic{
		svCtx: ctx,
	}
}

func (l *AuthLogic) ListAPI(c *gin.Context) {
	var req ListAuthReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		meta.ErrHandle(c, err)
		return
	}
	apis, err := l.svCtx.Mysql.AuthRelationModel.ListAPI(req.APPID, req.MenuID)
	if err != nil {
		log.Printf("get menu failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, conversion.SliceAPIDO2DTO(apis))
	return
}

func (l *AuthLogic) AddAPI(c *gin.Context) {
	var req AddAuthReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		meta.ErrHandle(c, err)
		return
	}
	err = l.addAPI(req)
	if err != nil {
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, nil)
	return
}
func (l *AuthLogic) addAPI(req AddAuthReq) error {
	apis, err := l.svCtx.Mysql.AuthRelationModel.ListAPI(req.APPID, req.MenuID)
	if err != nil {
		return err
	}
	// 区分应该保存或者删除的API
	var (
		newApiMap  = make(map[int64]struct{})
		oldApiMap  = make(map[int64]struct{})
		deleteAPIs []int64
		AddAPIs    []int64
	)
	for _, api := range req.AppIds {
		newApiMap[api] = struct{}{}
	}
	// 匹配删除的apiIds
	for _, api := range apis {
		oldApiMap[api.ID] = struct{}{}
		if _, ok := newApiMap[api.ID]; !ok {
			deleteAPIs = append(deleteAPIs, api.ID)
		}
	}
	// 匹配新增的apiIds
	for _, api := range req.AppIds {
		if _, ok := oldApiMap[api]; !ok {
			AddAPIs = append(AddAPIs, api)
		}
	}
	// 修改数据库
	//todo: 添加事务
	if len(deleteAPIs) > 0 {
		err = l.svCtx.Mysql.AuthRelationModel.DeleteAPI(req.APPID, req.MenuID, deleteAPIs)
	}
	if len(AddAPIs) > 0 {
		err = l.svCtx.Mysql.AuthRelationModel.AddAPI(req.APPID, req.MenuID, AddAPIs)
	}
	return err
}
