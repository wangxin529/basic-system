package models

import (
	"elevate-hub/common/set"
	"elevate-hub/db/meta"
	"elevate-hub/db/mysql"
	"gorm.io/gorm"
)

const (
	UserTableName = "user"
)

type User struct {
	BaseModel
	ModelTime
	Username      string `json:"username" gorm:"column:username;comment:用户名"`
	Password      string `json:"-" gorm:"column:password;comment:密码"`
	NickName      string `json:"nickname" gorm:"column:nickname;comment:昵称"`
	SupperManager int    `json:"supper_manager" gorm:"column:supper_manager"` // 0 不是 1 是
	Phone         string `json:"phone" gorm:"column:phone;comment:手机号"`
	Avatar        string `json:"avatar" gorm:"column:avatar;comment:头像"`
	Sex           int    `json:"sex" gorm:"column:sex;comment:性别"` // 1 男 2 女
	Email         string `json:"email" gorm:"column:email;comment:邮箱"`
	DeptId        int64  `json:"dept_id" gorm:"column:dept_id;comment:部门"`
	Post          string `json:"post" gorm:"column:post;comment:岗位"`
	Status        int    `json:"status" gorm:"column:status"` // 0 可用  1不可用
}

type UserModel struct {
	mysql.Table[User]
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		mysql.Table[User]{
			TableName: UserTableName,
			DB:        db,
		},
	}
}

func (u UserModel) List(option meta.ListOption) (outs []*UserAggregate, count int64, err error) {
	res, count, err := u.Table.List(option)
	if err != nil {
		return nil, 0, err
	}
	var ids = set.NewSet[int64]()
	for _, item := range res {
		ids.Add(item.ID)
	}
	sql := `select u.*, d.name as department_name from
	(select * from user where id in (?) and deleted_at is null) as u
		left join department as d on u.dept_id = d.id
		order by  u.id asc
	`
	err = u.Table.DB.Raw(sql, ids.ToArr()).Find(&outs).Error

	return
}

func (u UserModel) GetByID(id int64) (out *UserAggregate, err error) {

	sql := `select u.*, d.name as department_name from
	(select * from user where id in (?) and deleted_at is null) as u
		left join department as d on u.dept_id = d.id
	`
	err = u.Table.DB.Raw(sql, id).Find(&out).Error

	return
}
