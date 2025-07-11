package menu

import (
	"elevate-hub/common/set"
	"elevate-hub/common/utils"
	"elevate-hub/conversion"
	meta2 "elevate-hub/db/meta"
	"elevate-hub/db/models"
	"elevate-hub/meta"
	"elevate-hub/svc"
	"elevate-hub/types/dto"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuLogic struct {
	svCtx *svc.ServiceContext
}

func NewMenuLogic(ctx *svc.ServiceContext) *MenuLogic {
	return &MenuLogic{
		svCtx: ctx,
	}
}

func (l *MenuLogic) AddMenu(c *gin.Context) {
	var req dto.AddMenuReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("add menu failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "菜单新增失败")
		return
	}

	err = l.svCtx.Mysql.MenuModel.Create(models.Menu{
		Name:     req.Name,
		Parent:   req.Parent,
		Path:     req.Path,
		MenuType: req.MenuType,
	})
	if err != nil {
		log.Printf("add menu in db create failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "菜单新增失败请联系，系统管理员")
		return
	}
	meta.SuccessHandle(c, nil)
}

func (l *MenuLogic) UpdateMenu(c *gin.Context) {
	var req dto.UpdateMenuReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("update menu failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "菜单参数错误")
		return
	}

	menuDO := models.Menu{}.
		SetButtons(conversion.SliceButtonDTO2DO(req.Buttons))
	if req.MenuConfig != nil {
		menuDO = menuDO.SetMenuConfig(conversion.MenuConfigDTO2DO(req.MenuConfig))
	}
	err = l.svCtx.Mysql.MenuModel.Update(meta2.UpdateOption{
		ID: req.ID,
		Data: map[string]interface{}{
			"name":        req.Name,
			"path":        req.Path,
			"parent":      req.Parent,
			"status":      req.Status,
			"buttons":     menuDO.Buttons,
			"menu_type":   req.MenuType,
			"menu_config": menuDO.MenuConfig,
		},
	})
	if err != nil {
		log.Printf("update menu failed:%v", err)
		meta.ErrHandleWithMsg(c, "菜单修改失败")
		return
	}
	meta.SuccessHandle(c, nil)

}

func (l *MenuLogic) DelMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		meta.ErrHandleWithMsg(c, "指定需要删除的菜单")
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		meta.ErrHandleWithMsg(c, "删除菜单信息错误")
		return
	}
	err = l.svCtx.Mysql.MenuModel.Delete(meta2.DeleteOption{
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
func (l *MenuLogic) ListMenu(c *gin.Context) {
	var req dto.ListMenuReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.Printf("list menu failed:%v", err.Error())
		meta.ErrHandleWithMsg(c, "查询参数失败")
		return
	}

	res, _, err := l.svCtx.Mysql.MenuModel.List(meta2.EmptyListOption)

	if err != nil {
		log.Printf("list role failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	menus := conversion.SliceMenuDO2DTO(res)
	meta.SuccessHandle(c, menus)
	return
}

func (l *MenuLogic) MenuTree(c *gin.Context) {
	var req MenuTreeReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		log.Printf("list role failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	option := meta2.EmptyListOption
	if req.MenuType != 0 {
		option.AddToCondition("menu_type", req.MenuType)
	}
	res, _, err := l.svCtx.Mysql.MenuModel.List(option)
	if err != nil {
		log.Printf("list role failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	menus := conversion.SliceMenuDO2DTO(res)
	menuTree := conversion.Menu2Child(menus, nil)
	meta.SuccessHandle(c, menuTree.Children)
	return
}
func (l *MenuLogic) Menu(c *gin.Context) {
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
	res, err := l.svCtx.Mysql.MenuModel.Get(meta2.GetOption{
		ID: idInt,
	})
	if err != nil {
		log.Printf("get menu failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	meta.SuccessHandle(c, conversion.MenuDO2DTO(res))
	return
}

func (l *MenuLogic) CurrentUserMenu(c *gin.Context) {
	var (
		btns        = []string{}
		menuIds     []int64
		inCondition = map[string]interface{}{}
		ars         []*models.AuthRelation
	)
	roles, isAdmin, err := utils.GetRoles(c)
	if err != nil {
		log.Printf("get menu failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	if isAdmin {
		goto adminMenu
	}
	ars, _, err = l.svCtx.Mysql.AuthRelationModel.List(meta2.ListOption{
		DisableCount: true,
		InCondition: map[string]interface{}{
			"role_id": roles,
		},
		Condition: map[string]interface{}{
			"app_id":  0,
			"api_id":  0,
			"user_id": 0,
		},
	})
	if err != nil {
		log.Printf("get role failed:%v", err)
		meta.ErrHandle(c, err)
		return
	}
	menuIds = make([]int64, 0, len(ars))
	for _, role := range ars {
		btns = append(btns, role.GetButtons()...)
		menuIds = append(menuIds, role.MenuID)
	}

	inCondition["id"] = menuIds
adminMenu:
	menus, _, err := l.svCtx.Mysql.MenuModel.List(meta2.ListOption{
		DisableCount: true,
		Condition: map[string]interface{}{
			"status": 0,
		},
		InCondition: inCondition,
	})
	if err != nil {
		meta.ErrHandle(c, err)
		return
	}

	if len(menus) == 0 {
		meta.SuccessHandle(c, nil)
		return
	}
	var (
		appKeys = set.NewSet[string]()
		//route   []interface{}
	)
	for _, menu := range menus {
		menuConfig := menu.GetMenuConfig()
		if menuConfig.Type == models.ChildrenRouterType {
			//route = append(route, map[string]interface{}{
			//	"path":     menu.Path,
			//	"microApp": menuConfig.APP,
			//})
			appKeys.Add(menuConfig.APP)
		}
	}
	apps, _, err := l.svCtx.Mysql.APPModel.List(meta2.ListOption{
		DisableCount: true,
		InCondition: map[string]interface{}{
			"app_key": appKeys.ToArr(),
		},
	})
	if err != nil {
		return
	}
	var entrys = make([]interface{}, 0, len(apps))
	for _, app := range apps {
		entrys = append(entrys, app.GetEntry())
	}
	meta.SuccessHandle(c, map[string]any{
		"menu": conversion.Menu2Child(conversion.SliceMenuDO2DTO(menus), nil),
		"app":  entrys,
		//"app_route": route,
		"buttons": btns,
	})

}
