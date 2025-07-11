package conversion

import (
	"elevate-hub/db/models"
	"elevate-hub/types/dto"
	"time"
)

func APIDO2DTO(in *models.API) *dto.API {
	return &dto.API{
		ID:       in.ID,
		Name:     in.Name,
		Method:   in.Method,
		Status:   in.Status,
		Path:     in.Path,
		APIType:  in.APIType,
		CreateAt: in.CreatedAt.Format(time.DateTime),
		Describe: in.Describe,
	}
}

func SliceAPIDO2DTO(ins []*models.API) []*dto.API {
	var out = make([]*dto.API, 0, len(ins))

	for _, in := range ins {
		out = append(out, APIDO2DTO(in))
	}
	return out
}

func APIAggregateDO2DTO(in *models.APIAggregate) *dto.API {
	return &dto.API{
		ID:       in.ID,
		Name:     in.Name,
		Method:   in.Method,
		Status:   in.Status,
		Path:     in.Path,
		CreateAt: in.CreatedAt.Format(time.DateTime),
		Describe: in.Describe,
		//Key:      in.Key,
		Creator: in.CreatorUserName,
		APIType: in.APIType,
	}
}

func SliceAPIAggregateDO2DTO(ins []*models.APIAggregate) []*dto.API {
	var out = make([]*dto.API, 0, len(ins))

	for _, in := range ins {
		out = append(out, APIAggregateDO2DTO(in))
	}
	return out
}
