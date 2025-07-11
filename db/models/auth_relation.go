package models

import (
	"elevate-hub/db/meta"
	"elevate-hub/db/mysql"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

const (
	AuthRelationTableName = "auth_relation"
)

// user-role  AppID-ApiID   roleId - menuID-permission-json   menuID - ApiID
type AuthRelation struct {
	RoleID           int64  `json:"role_id" gorm:"column:role_id"`
	AppID            int64  `json:"app_id" gorm:"column:app_id"`
	MenuID           int64  `json:"menu_id" gorm:"column:menu_id"`
	ApiID            int64  `json:"api_id" gorm:"column:api_id"`
	UserID           int64  `json:"user_id" gorm:"column:user_id"`
	ButtonPermission []byte `json:"button_permission" gorm:"column:button_permission"`
}

func (a *AuthRelation) SetButtons(btns []string) *AuthRelation {
	marshal, _ := json.Marshal(btns)
	a.ButtonPermission = marshal
	return a
}

func (a *AuthRelation) GetButtons() (btns []string) {
	err := json.Unmarshal(a.ButtonPermission, &btns)
	if err != nil {
		return []string{}
	}
	return btns
}

type AuthRelationModel struct {
	mysql.Table[AuthRelation]
}

func NewAuthRelationAuthRelationModel(db *gorm.DB) *AuthRelationModel {
	return &AuthRelationModel{
		mysql.Table[AuthRelation]{
			TableName: AuthRelationTableName,
			DB:        db,
			Unscoped:  true,
		},
	}
}
func (m *AuthRelationModel) ListAPI(appID, menuID int64) ([]*API, error) {
	var outs []*API
	_, err := m.ListAggregate(meta.ListOption{
		DisableCount: true,
		Condition: map[string]interface{}{
			"app_id":  appID,
			"role_id": 0,
			"menu_id": menuID,
			"user_id": 0,
		},
		Join:   fmt.Sprintf("left join %s as api on %s.api_id =api.id", APITableName, m.TableName),
		Select: fmt.Sprintf("%s.*", APITableName),
	}, &outs)
	return outs, err
}

func (m *AuthRelationModel) DeleteAPI(appID, menuID int64, apis []int64) error {
	if len(apis) == 0 {
		return nil
	}
	err := m.Delete(meta.DeleteOption{
		Condition: map[string]interface{}{
			"app_id":  appID,
			"menu_id": menuID,
		},
		InCondition: map[string]interface{}{
			"api_id": apis,
		},
	})
	return err
}

func (m *AuthRelationModel) AddAPI(appID, menuID int64, apiIds []int64) error {
	if len(apiIds) == 0 {
		return nil
	}
	var auths = make([]*AuthRelation, 0, len(apiIds))
	for _, apiId := range apiIds {
		auths = append(auths, &AuthRelation{

			AppID:  appID,
			MenuID: menuID,
			ApiID:  apiId,
		})
	}

	err := m.Creates(auths)
	return err
}

func (m *AuthRelationModel) GetUserRoles(userId int64) ([]int64, error) {
	var roleIDs []int64
	_, err := m.ListAggregate(meta.ListOption{
		Condition: map[string]interface{}{
			"user_id": userId,
			"app_id":  0,
			"menu_id": 0,
			"api_id":  0,
		},
		CustomCondition: &meta.CustomCondition{
			SQL: "role_id != 0",
		},
		Select:       "role_id",
		DisableCount: true,
	}, &roleIDs)
	return roleIDs, err
}
