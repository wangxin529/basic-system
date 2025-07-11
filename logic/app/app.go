package app

import (
	"elevate-hub/common/utils"
	"elevate-hub/conversion"
	meta2 "elevate-hub/db/meta"
	"elevate-hub/db/models"
	"elevate-hub/meta"
	"elevate-hub/svc"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type AppLogic struct {
	svCtx *svc.ServiceContext
}

func NewAppLogic(ctx *svc.ServiceContext) *AppLogic {
	return &AppLogic{
		svCtx: ctx,
	}
}

func (l *AppLogic) AddApp(c *gin.Context) {
	var req AddAppReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("add App failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "应用新增失败")
		return
	}

	err = l.svCtx.Mysql.APPModel.Create(models.APP{
		Name:       req.Name,
		Key:        req.Key,
		SignMethod: req.SignMethod,
		Status:     req.Status,
		Secret:     utils.MD5SecureRandomString(10),
	})
	if err != nil {
		log.Printf("add App in db create failed:%v", err)
		meta.ErrHandleWithMsg(c, "应用新增失败请联系，系统管理员")
		return
	}
	meta.SuccessHandle(c, nil)
}

func (l *AppLogic) UpdateApp(c *gin.Context) {
	//var req
	//c.ShouldBindJSON(&)
	//l.svCtx.Mysql.APPModel.Update(meta2.UpdateOption{
	//	ID:
	//})
}

func (l *AppLogic) DelApp(c *gin.Context) {
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
	err = l.svCtx.Mysql.APPModel.Delete(meta2.DeleteOption{
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

func (l *AppLogic) ListApp(c *gin.Context) {

	var req ListAppReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.Printf("list App failed:%v", err)
		meta.ErrHandleWithMsg(c, "查询参数失败")
		return
	}
	res, count, err := l.svCtx.Mysql.APPModel.List(meta2.ListOption{
		Page: &req.Page,
	})

	if err != nil {
		log.Printf("list App failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}

	meta.SuccessHandleAndTotal(c, conversion.SliceAPPDO2DTO(res), count)
	return

}
func (l *AppLogic) App(c *gin.Context) {
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
	res, err := l.svCtx.Mysql.APPModel.Get(meta2.GetOption{
		ID: idInt,
	})
	if err != nil {
		log.Printf("get menu failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, conversion.APPDO2DTO(res))
	return
}

func (l *AppLogic) UpdateProxy(c *gin.Context) {
	var req UpdateProxyReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("update proxy failed:%v", err.Error())
		meta.ErrHandle(c, err)
		return
	}
	err = l.svCtx.Mysql.APPModel.Update(meta2.UpdateOption{
		ID: req.ID,
		Data: map[string]interface{}{
			"proxy": models.APP{}.SetProxy(req.Proxy).Proxy,
		},
	})
	if err != nil {
		log.Printf("update proxy failed:%v", err.Error())
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, nil)
	return
}

func (l *AppLogic) ListProxy(c *gin.Context) {
	id, err := meta.Param2Int("id", c)
	if err != nil {
		meta.ErrHandle(c, err)
		return
	}
	res, err := l.svCtx.Mysql.APPModel.Get(meta2.GetOption{
		ID: id,
	})
	if err != nil {
		log.Printf("get menu failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, res.GetProxy())
}
func (l *AppLogic) EntryList(c *gin.Context) {
	res, _, err := l.svCtx.Mysql.APPModel.List(meta2.ListOption{
		DisableCount: true,
		Select:       "entry",
		CustomCondition: &meta2.CustomCondition{
			SQL: "entry is not null",
		}},
	)
	if err != nil {
		log.Printf("get entry failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	var entrys = make([]models.Entry, 0)
	for _, item := range res {
		entrys = append(entrys, item.GetEntry())
	}
	meta.SuccessHandle(c, entrys)
	return
}

func (l *AppLogic) UpdateEntry(c *gin.Context) {
	var req UpdateAppEntry
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("update entry failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	marshal, _ := json.Marshal(req)
	err = l.svCtx.Mysql.APPModel.Update(meta2.UpdateOption{
		ID: req.ID,
		Data: map[string]interface{}{
			"entry": marshal,
		},
	})
	if err != nil {
		log.Printf("update entry failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, nil)
	return
}
