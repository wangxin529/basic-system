package models

import (
	"elevate-hub/db/mysql"
	"encoding/json"
	"gorm.io/gorm"
)

type MenuConfigType = string

const (
	MenuTableName                     = "menu"
	ComponentPathType  MenuConfigType = "component_path"
	ChildrenRouterType MenuConfigType = "children_router"
)

type MenuButton struct {
	Name          string  `json:"name"`
	PermissionKey string  `json:"permission_key"`
	API           []int64 `json:"api"`
}
type MenuConfig struct {
	Type          string `json:"type"`
	APP           string `json:"app"`
	ComponentPath string `json:"component_path"`
	//PermissionKey string `json:"permission_key"`
}

type Menu struct {
	BaseModel
	ModelTime
	Name       string  `json:"name" gorm:"column:name"`
	Path       string  `json:"path" gorm:"column:path"`
	Status     int     `json:"status" gorm:"column:status"` // 0 可用  1不可用
	Parent     int64   `json:"parent" gorm:"column:parent"` //父级菜单
	Sort       float64 `json:"sort" gorm:"sort"`
	MenuConfig []byte  `json:"menu_config" gorm:"column:menu_config"`
	MenuType   int     `json:"menu_type" gorm:"column:menu_type"`
	Buttons    []byte  `json:"buttons" gorm:"column:buttons"`
}

func (m Menu) SetButtons(btns []*MenuButton) Menu {
	marshal, _ := json.Marshal(btns)
	m.Buttons = marshal
	return m
}

func (m Menu) GetButtons() []*MenuButton {
	var btns []*MenuButton
	err := json.Unmarshal(m.Buttons, &btns)
	if err != nil {
		btns = []*MenuButton{}
	}
	return btns
}

func (m Menu) SetMenuConfig(mc *MenuConfig) Menu {
	marshal, _ := json.Marshal(mc)
	m.MenuConfig = marshal
	return m
}

func (m Menu) GetMenuConfig() MenuConfig {
	var mc MenuConfig
	err := json.Unmarshal(m.MenuConfig, &mc)
	if err != nil {
		mc = MenuConfig{}
	}
	return mc
}

type MenuModel struct {
	mysql.Table[Menu]
}

func NewMenuModel(db *gorm.DB) *MenuModel {
	return &MenuModel{
		mysql.Table[Menu]{
			TableName: MenuTableName,
			DB:        db,
		},
	}
}
