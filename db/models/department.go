package models

import (
	"elevate-hub/common/utils"
	"elevate-hub/db/meta"
	"elevate-hub/db/mysql"
	"fmt"
	"gorm.io/gorm"
)

const (
	DepartmentTableName = "department"
)

type Department struct {
	BaseModel
	ModelTime
	Name   string `json:"name" gorm:"column:name"`
	Code   string `json:"code" gorm:"column:code"`
	Leader int64  `json:"leader" gorm:"column:leader"` // userID
	Parent int64  `json:"parent" gorm:"column:parent"` //父级部门ID
}

func (d Department) GenerateCode() Department {
	d.Code = utils.MD5(fmt.Sprintf("%s_%d", d.Name, d.Parent))
	return d
}

type DepartmentModel struct {
	mysql.Table[Department]
}

func NewDepartmentModel(db *gorm.DB) *DepartmentModel {
	return &DepartmentModel{
		mysql.Table[Department]{
			TableName: DepartmentTableName,
			DB:        db,
		},
	}
}

func (m *DepartmentModel) GetByID(id int64) (out *DepartmentAggregate, err error) {
	sql := `select d.*, u.id as leader_user, u.nickname as leader_nickname, u.email as leader_emial from  
        	(select * from department  where id = ? and deleted_at IS NULL) as d left join user as u on d.leader = u.id `
	err = m.DB.Raw(sql, id).First(&out).Error
	return
}

func (m *DepartmentModel) List(option meta.ListOption) (out []*DepartmentAggregate, count int64, err error) {
	dSql := "from department  where deleted_at IS NULL"
	err = m.DB.Raw(fmt.Sprintf("select count(*) %s", dSql)).First(&count).Error
	if err != nil {
		return
	}
	if option.Page != nil && option.Page.PageSize > 0 && option.Page.Page > -1 {
		page := option.Page
		dSql += fmt.Sprintf(" limit %d,%d ", (page.Page-1)*page.PageSize, page.PageSize)
	}
	sql := fmt.Sprintf(` select d.*, u.id as leader_user, u.nickname as leader_nickname, u.email as leader_emial from  
        	(select * %s) as d left join user as u on d.leader = u.id `, dSql)
	err = m.DB.Raw(sql).Find(&out).Error
	return
}
