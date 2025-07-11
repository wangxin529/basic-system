package conversion

import (
	"elevate-hub/db/models"
	"elevate-hub/types/dto"
)

func RoleDO2DTO(in *models.Role, ar []*models.AuthRelation) *dto.Role {
	role := &dto.Role{
		ID:     in.ID,
		Name:   in.Name,
		Status: in.Status,
	}
	role.OperatePermissions = make([]dto.OperatePermissions, 0, len(ar))
	if len(ar) == 0 {
		return role
	}
	for _, v := range ar {
		role.OperatePermissions = append(role.OperatePermissions, dto.OperatePermissions{
			MenuID:            v.MenuID,
			ButtonPermissions: v.GetButtons(),
		})
	}
	return role
}

func SliceRoleDO2DTO(ins []*models.Role, arMap map[int64][]*models.AuthRelation) []*dto.Role {
	var out = make([]*dto.Role, 0, len(ins))
	if arMap == nil {
		arMap = make(map[int64][]*models.AuthRelation)
	}
	for _, in := range ins {
		ar := arMap[in.ID]
		out = append(out, RoleDO2DTO(in, ar))
	}
	return out
}
