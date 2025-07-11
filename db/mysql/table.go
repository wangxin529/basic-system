package mysql

import (
	"elevate-hub/db/meta"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type TableOperate[T any] interface {
	Create(data T) error
	Delete(option meta.DeleteOption) error
	Update(option meta.UpdateOption) error
	List(option meta.ListOption) (res []*T, count int64, err error)
	ListAggregate(option meta.ListOption, out interface{}) (count int64, err error)
	Get(option meta.GetOption) (res *T, err error)
	Count(option meta.CountOption) (int64, error)
	Copy(db *gorm.DB) TableOperate[T]
	//Transaction(f func(tx *gorm.DB) error) error
}

//func NewTable[T any](db *gorm.DB, tableName string) *TableOperate[T] {
//	return Table[T]{
//		DB:        db,
//		TableName: tableName,
//	}
//}

type Table[T any] struct {
	DB        *gorm.DB
	TableName string
	Unscoped  bool // 是否根据具delete_at字段
}

func (t *Table[T]) Copy(db *gorm.DB) TableOperate[T] {
	return &Table[T]{
		DB:        db,
		TableName: t.TableName,
		Unscoped:  t.Unscoped,
	}
}
func (t *Table[T]) Create(data T) error {
	return t.DB.Table(t.TableName).Create(&data).Error

}
func (t *Table[T]) Creates(data []*T) error {
	return t.DB.Table(t.TableName).Create(&data).Error

}
func (t *Table[T]) Delete(option meta.DeleteOption) error {
	tx := t.DB.Table(t.TableName)
	if option.ID != 0 {
		tx = tx.Where("id = ?", option.ID)
	}

	if option.IDs != nil {
		tx = tx.Where("id in (?)", option.IDs)
	}
	if option.Condition != nil {
		tx = tx.Where(option.Condition)
	}
	if option.InCondition != nil {
		for key, val := range option.InCondition {
			tx = tx.Where(key+" in ?", val)
		}
	}
	if t.Unscoped {

		return tx.Delete(t.TableName).Error
	}
	return tx.Update("deleted_at", time.Now()).Error
}
func (t *Table[T]) Update(option meta.UpdateOption) error {
	tx := t.DB.Table(t.TableName)
	if option.ID != 0 {
		tx = tx.Where("id = ?", option.ID)
	}
	if option.Condition != nil {
		for key, val := range option.Condition {
			tx = tx.Where(key+" = ? ", val)
		}
	}

	return tx.Updates(&option.Data).Error
}
func (t *Table[T]) Get(option meta.GetOption) (res *T, err error) {
	tx := t.DB.Table(t.TableName)
	if option.ID != 0 {
		tx = tx.Where("id = ?", option.ID)
	}
	if option.Condition != nil {
		tx = tx.Where(option.Condition)
	}
	if option.InCondition != nil {
		for key, val := range option.InCondition {
			tx = tx.Where(key+" in ?", val)
		}
	}
	if option.LikeCondition != nil {
		for key, val := range option.LikeCondition {
			tx = tx.Where(key+" like ?", fmt.Sprintf("%%%v%%", val))
		}
	}

	if option.IsLast {
		tx = tx.Last(&res)
	} else {
		tx = tx.First(&res)
	}
	err = tx.Error
	return
}

func (t *Table[T]) listToSql(option meta.ListOption) (tx *gorm.DB, count int64, err error) {
	page := option.Page

	tx = t.DB.Table(t.TableName)
	if option.Condition != nil && len(option.Condition) != 0 {
		tx = tx.Where(option.Condition)
	}

	if len(option.IDs) != 0 {
		tx = tx.Where("id in (?)", option.IDs)
	}

	if option.InCondition != nil {
		for key, val := range option.InCondition {
			tx = tx.Where(key+" in ?", val)
		}
	}

	if option.LikeCondition != nil && len(option.LikeCondition) != 0 {
		for key, value := range option.LikeCondition {
			tx = tx.Where(fmt.Sprintf("%s like ?", key), fmt.Sprintf("%%%v%%", value))
		}
	}
	if option.CustomCondition != nil {
		tx = tx.Where(option.CustomCondition.SQL, option.CustomCondition.Args...)
	}

	if option.Between != nil {
		tx = tx.Where(option.Between.Field+" between ? and ?", option.Between.Value1, option.Between.Value2)
	}
	if !option.DisableCount && !t.Unscoped {
		err = tx.Where(fmt.Sprintf("`%s`.deleted_at is null", t.TableName)).Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
	}
	if option.Join != "" {
		tx = tx.Joins(option.Join)
	}
	if option.Order != nil {
		tx = tx.Order(option.Order.Filed + " " + option.Order.Sort)
	}
	if option.Page != nil && option.Page.Page > -1 && option.Page.PageSize > 0 {
		tx = tx.Offset((page.Page - 1) * page.PageSize).Limit(page.PageSize)
	}
	if option.Select != "" {
		tx = tx.Select(option.Select)
	}
	if t.Unscoped {
		tx = tx.Unscoped()
	}
	return tx, count, nil
}
func (t *Table[T]) List(option meta.ListOption) (res []*T, count int64, err error) {
	tx, count, err := t.listToSql(option)
	if err != nil {
		return nil, 0, err
	}
	err = tx.Find(&res).Error
	if err == gorm.ErrRecordNotFound {
		return []*T{}, 0, nil
	}
	return
}
func (t *Table[T]) ListAggregate(option meta.ListOption, out interface{}) (count int64, err error) {
	tx, count, err := t.listToSql(option)
	if err != nil {
		return 0, err
	}
	err = tx.Find(out).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return
}

func (t *Table[T]) Count(option meta.CountOption) (count int64, err error) {
	tx := t.DB.Table(t.TableName)
	if !option.IsDelete {
		tx = tx.Where("deleted_at is null")
	}
	err = tx.Count(&count).Error
	return
}

func (t *Table[T]) listAllTable(option meta.ListOption) (res []*T, count int64, err error) {

	return res, 0, nil
}

//
//func Transaction(f func(tx *gorm.DB) error) error {
//
//	return mysql.DB.Transaction(f)
//}
