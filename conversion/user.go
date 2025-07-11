package conversion

import (
	"elevate-hub/db/models"
	"elevate-hub/types/dto"
)

func UserAggregateDO2DTO(in *models.UserAggregate) *dto.User {
	return &dto.User{
		ID:       in.ID,
		Username: in.Username,
		NickName: in.NickName,
		IsSupper: in.SupperManager,
		Phone:    in.Phone,
		Avatar:   in.Avatar,
		Status:   in.Status,
		Sex:      in.Sex,
		Email:    in.Email,
		Dept:     in.DepartmentName,
		DeptID:   in.DeptId,
		Post:     in.Post,
	}
}

func SliceUserAggregate2DTO(ins []*models.UserAggregate) []*dto.User {
	var out = make([]*dto.User, 0, len(ins))

	for _, user := range ins {
		out = append(out, UserAggregateDO2DTO(user))
	}
	return out
}

func UserDO2DTO(in *models.User) *dto.User {
	return &dto.User{
		ID:       in.ID,
		Username: in.Username,
		NickName: in.NickName,
		IsSupper: in.SupperManager,
		Phone:    in.Phone,
		Avatar:   in.Avatar,
		Status:   in.Status,
		Sex:      in.Sex,
		Email:    in.Email,
		Post:     in.Post,
		DeptID:   in.DeptId,
	}
}

func SliceUser2DTO(ins []*models.User) []*dto.User {
	var out = make([]*dto.User, 0, len(ins))

	for _, user := range ins {
		out = append(out, UserDO2DTO(user))
	}
	return out
}
