package middleware

import (
	"elevate-hub/common/set"
	"elevate-hub/db/meta"
	"elevate-hub/db/models"
	"elevate-hub/svc"
)

type Permission struct {
	svCtx        *svc.ServiceContext
	rolesAPI     map[int64]set.Set[string] // 角色
	roleButtons  map[int64]set.Set[string]
	roleMenus    map[int64]set.Set[int64]
	needValidUrl set.Set[string] // 所有需要验证的URL
}

func NewPermission(svCtx *svc.ServiceContext) *Permission {
	p := &Permission{
		svCtx:        svCtx,
		rolesAPI:     make(map[int64]set.Set[string]),
		roleButtons:  make(map[int64]set.Set[string]),
		roleMenus:    make(map[int64]set.Set[int64]),
		needValidUrl: set.Set[string]{},
	}
	p.reload()
	return p
}

func (p *Permission) reload() {
	auths, _, err := p.svCtx.Mysql.AuthRelationModel.List(meta.ListOption{
		DisableCount: true,
		CustomCondition: &meta.CustomCondition{
			SQL:  "role_id != 0 and user_id = 0 and menu_id != 0 and api_id = 0 ",
			Args: []any{},
		},
	})
	if err != nil {
		return
	}

	menus, _, err := p.svCtx.Mysql.MenuModel.List(meta.ListOption{
		DisableCount: true,
	})
	if err != nil {
		return
	}
	apis, _, err := p.svCtx.Mysql.APIModel.List(meta.ListOption{
		Condition: map[string]interface{}{
			"api_type": 2,
		},
		DisableCount: true,
	})
	var apiDict = make(map[int64]*models.APIAggregate)
	for _, api := range apis {
		p.needValidUrl.Add(api.Path)
		apiDict[api.ID] = api
	}

	var menuDict = make(map[int64]*models.Menu)

	for _, menu := range menus {
		menuDict[menu.ID] = menu
	}

	for _, auth := range auths {
		if _, ok := p.roleButtons[auth.RoleID]; !ok {
			p.roleButtons[auth.RoleID] = set.NewSet[string]()
		}
		// 角色按钮

		if _, ok := p.roleMenus[auth.RoleID]; !ok {
			p.roleMenus[auth.RoleID] = set.NewSet[int64]()
		}
		p.roleButtons[auth.RoleID].Add(auth.GetButtons()...)
		// 角色菜单
		if _, ok := p.roleMenus[auth.RoleID]; !ok {
			p.roleMenus[auth.RoleID] = set.NewSet[int64]()
		}
		p.roleMenus[auth.RoleID].Add(auth.MenuID)
		if _, ok := menuDict[auth.MenuID]; ok {
			for _, btn := range menuDict[auth.MenuID].GetButtons() {
				for _, apiID := range btn.API {
					if api, okAPI := apiDict[apiID]; okAPI {
						if _, ok := p.rolesAPI[auth.RoleID]; !ok {
							p.rolesAPI[auth.RoleID] = set.NewSet[string]()
						}
						p.rolesAPI[auth.RoleID].Add(api.Path)
					}
				}

			}
		}
	}

}
