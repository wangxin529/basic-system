package models

import (
	"elevate-hub/db/mysql"
	"gorm.io/gorm"
)

const (
	RoleTableName = "role"
)

type Role struct {
	BaseModel
	ModelTime
	Name   string `json:"name" gorm:"column:name"`
	Status int    `json:"status" gorm:"column:status"` // 0 可用 1 不可用
}

type RoleModel struct {
	mysql.Table[Role]
}

func NewRoleModel(db *gorm.DB) *RoleModel {
	return &RoleModel{
		mysql.Table[Role]{
			TableName: RoleTableName,
			DB:        db,
		},
	}
}
