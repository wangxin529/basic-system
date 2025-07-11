package dto

import "elevate-hub/db/meta"

type ListMenuReq struct {
	meta.Page
}

type AddMenuReq struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Parent   int64  `json:"parent"`
	MenuType int    `json:"menu_type"`
}

type Button struct {
	Name          string  `json:"name"`
	PermissionKey string  `json:"permission_key"`
	Api           []int64 `json:"api"`
}

type MenuConfig struct {
	Type string `json:"type"`
	APP  string `json:"app"`
	//ChildrenRouter string `json:"children_router"`
	ComponentPath string `json:"component_path"`
	//PermissionKey string `json:"permission_key"`
}

type UpdateMenuReq struct {
	ID         int64       `json:"id"`
	Name       string      `json:"title"`
	Path       string      `json:"path"`
	Parent     int64       `json:"parent"`
	Status     int64       `json:"status"`
	MenuType   int         `json:"menu_type"`
	Buttons    []*Button   `json:"buttons"`
	MenuConfig *MenuConfig `json:"menu_config"`
}
