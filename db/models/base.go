package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID int64 `json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
}

type ModelTime struct {
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
