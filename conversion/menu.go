package conversion

import (
	"elevate-hub/db/models"
	"elevate-hub/types/dto"
	"fmt"
)

func ButtonDTO2DO(button *dto.Button) *models.MenuButton {
	return &models.MenuButton{
		Name:          button.Name,
		PermissionKey: button.PermissionKey,
		API:           button.Api,
	}
}

func SliceButtonDTO2DO(ins []*dto.Button) []*models.MenuButton {
	var out = make([]*models.MenuButton, 0, len(ins))

	for _, in := range ins {
		out = append(out, ButtonDTO2DO(in))
	}
	return out
}

func MenuConfigDTO2DO(in *dto.MenuConfig) *models.MenuConfig {
	return &models.MenuConfig{
		Type:          in.Type,
		APP:           in.APP,
		ComponentPath: in.ComponentPath,
		//PermissionKey: in.PermissionKey,
	}
}

func MenuConfigDO2DTO(in models.MenuConfig) *dto.MenuConfig {

	return &dto.MenuConfig{
		Type: in.Type,
		APP:  in.APP,
		//ChildrenRouter: in.ChildrenRouter,
		ComponentPath: in.ComponentPath,
		//PermissionKey: in.PermissionKey,
	}
}

func ButtonDO2DTO(button *models.MenuButton) *dto.Button {
	return &dto.Button{
		Name:          button.Name,
		PermissionKey: button.PermissionKey,
		Api:           button.API,
	}
}

func SliceButtonDO2DTO(ins []*models.MenuButton) []*dto.Button {
	var out = make([]*dto.Button, 0, len(ins))

	for _, in := range ins {
		out = append(out, ButtonDO2DTO(in))
	}
	return out
}

func MenuDO2DTO(in *models.Menu) *dto.Menu {
	return &dto.Menu{
		ID:         in.ID,
		Name:       in.Name,
		Value:      in.ID,
		Path:       in.Path,
		Status:     in.Status,
		Parent:     in.Parent,
		Buttons:    SliceButtonDO2DTO(in.GetButtons()),
		MenuType:   in.MenuType,
		MenuConfig: MenuConfigDO2DTO(in.GetMenuConfig()),
	}
}

func SliceMenuDO2DTO(ins []*models.Menu) []*dto.Menu {
	var out = make([]*dto.Menu, 0, len(ins))

	for _, in := range ins {
		out = append(out, MenuDO2DTO(in))
	}
	return out
}

func Menu2Child(Menus []*dto.Menu, root *dto.Menu) *dto.Menu {
	// 使用map存储部门，提高查找效率
	MenuMap := make(map[int64]*dto.Menu)
	for _, dept := range Menus {
		MenuMap[dept.ID] = dept
	}
	if root == nil {
		root = &dto.Menu{
			ID: 0,
		}
	}

	for _, dept := range Menus {
		if dept.Parent == root.ID {
			root.Children = append(root.Children, dept)
			dept.Children = []*dto.Menu{}
			buildMenuChildren(MenuMap, dept)
		}
	}
	return root
}

func buildMenuChildren(MenuMap map[int64]*dto.Menu, parent *dto.Menu) {
	for _, dept := range MenuMap {
		if dept.Parent == parent.ID {
			parent.Children = append(parent.Children, dept)
			dept.Children = []*dto.Menu{}
			buildMenuChildren(MenuMap, dept)
		}
	}
}

func MenuDO2DTOBrief(in *models.Menu) *dto.MenuBrief {
	return &dto.MenuBrief{
		ID:         menu_key_prefix(in.ID),
		Name:       in.Name,
		Path:       in.Path,
		Parent:     menu_key_prefix(in.Parent),
		MenuType:   in.MenuType,
		MenuConfig: MenuConfigDO2DTO(in.GetMenuConfig()),
	}
}

func SliceMenuDO2DTOBrief(ins []*models.Menu) []*dto.MenuBrief {

	var out = make([]*dto.MenuBrief, 0, len(ins))
	for _, in := range ins {
		out = append(out, MenuDO2DTOBrief(in))
	}
	return out
}

func menu_key_prefix(i int64) string {
	return fmt.Sprintf("menu_%d", i)
}
