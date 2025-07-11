package models

import (
	"elevate-hub/db/meta"
	"elevate-hub/db/mysql"
	"fmt"
	"gorm.io/gorm"
)

const (
	APITableName = "api"
)

type API struct {
	BaseModel
	ModelTime
	Name string `json:"name" gorm:"column:name"`
	//Key      string `json:"key" gorm:"column:key"`
	Status   int    `json:"status" gorm:"column:status"` // 0 可用  1不可用
	Method   string `json:"method" gorm:"column:method"`
	Path     string `json:"path" gorm:"column:path"`
	Creator  int64  `json:"creator" gorm:"column:creator"` // 创建人
	Describe string `json:"describe" gorm:"column:describe"`
	APIType  int    `json:"api_type" gorm:"column:api_type"` // 1 app   2    menu
}

type APIModel struct {
	mysql.Table[API]
}

func NewAPIModel(db *gorm.DB) *APIModel {
	return &APIModel{
		mysql.Table[API]{
			TableName: APITableName,
			DB:        db,
		},
	}
}

func (a *APIModel) List(option meta.ListOption) (outs []*APIAggregate, count int64, err error) {
	option.Join = fmt.Sprintf("left join %s on %s.id = %s.creator", UserTableName, UserTableName, a.TableName)
	option.Select = fmt.Sprintf("%s.*,%s.nickname as creator_username", a.TableName, UserTableName)
	count, err = a.ListAggregate(option, &outs)
	return

}
