package models

import (
	"elevate-hub/db/mysql"
	"encoding/json"
	"gorm.io/gorm"
)

const (
	APPTableName = "app"
)

type APP struct {
	BaseModel
	ModelTime
	Name       string `json:"name" gorm:"column:name"`
	Key        string `json:"app_key" gorm:"column:app_key"`
	Entry      []byte `json:"entry" gorm:"column:entry"` // 前端入口
	Secret     string `json:"secret" gorm:"column:secret"`
	Status     int    `json:"status" gorm:"column:status"`
	SignMethod string `json:"sign_method" gorm:"column:sign_method"` // 加密算法     md5
	Proxy      []byte `json:"proxy" gorm:"column:proxy"`
}

type Proxy struct {
	Prefix string `json:"prefix"`
	Addr   string `json:"addr"`
}

type Entry struct {
	Name       string `json:"name"`
	Entry      string `json:"entry"`
	Container  string `json:"container"`
	ActiveRule string `json:"active_rule"`
}

func (a APP) GetProxy() (outs []Proxy) {
	err := json.Unmarshal(a.Proxy, &outs)
	if err != nil {
		return []Proxy{}
	}
	return outs
}

func (a APP) SetProxy(ins []Proxy) APP {
	marshal, _ := json.Marshal(ins)
	a.Proxy = marshal
	return a
}

func (a APP) GetEntry() (out Entry) {
	err := json.Unmarshal(a.Entry, &out)
	if err != nil {
		return Entry{}
	}
	return out
}

func (a APP) SetEntry(ins []Entry) APP {
	marshal, _ := json.Marshal(ins)
	a.Entry = marshal
	return a
}

type APPModel struct {
	mysql.Table[APP]
}

func NewAPPModel(db *gorm.DB) *APPModel {
	return &APPModel{
		mysql.Table[APP]{
			TableName: APPTableName,
			DB:        db,
		},
	}
}
