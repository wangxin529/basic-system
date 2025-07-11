package conversion

import (
	"elevate-hub/db/models"
	"elevate-hub/types/dto"
)

func DepartmentDO2DTO(in *models.Department) *dto.Department {
	return &dto.Department{
		ID:     in.ID,
		Name:   in.Name,
		Parent: in.Parent,
		Leader: "",
	}
}

func SliceDepartmentDO2DTO(ins []*models.Department) []*dto.Department {
	var out = make([]*dto.Department, 0, len(ins))

	for _, in := range ins {
		out = append(out, DepartmentDO2DTO(in))
	}
	return out
}

func DepartmentAggregateDO2DTO(in *models.DepartmentAggregate) *dto.Department {
	return &dto.Department{
		ID:     in.ID,
		Name:   in.Name,
		Parent: in.Parent,
		Leader: in.LeaderNickname,
	}
}

func SliceDepartmentAggregateDO2DTO(ins []*models.DepartmentAggregate) []*dto.Department {
	var out = make([]*dto.Department, 0, len(ins))

	for _, in := range ins {
		out = append(out, DepartmentAggregateDO2DTO(in))
	}
	return out
}

func ProxyDO2DTO(in models.Proxy) dto.Proxy {
	return dto.Proxy{
		Prefix: in.Prefix,
		Addr:   in.Addr,
	}

}
func SliceProxyDO2DTO(ins []models.Proxy) []dto.Proxy {
	var out = make([]dto.Proxy, 0, len(ins))
	for _, in := range ins {
		out = append(out, ProxyDO2DTO(in))
	}
	return out
}

func APPDO2DTO(in *models.APP) *dto.APP {
	return &dto.APP{
		ID:         in.ID,
		Name:       in.Name,
		Key:        in.Key,
		Status:     in.Status,
		SignMethod: in.SignMethod,
		Entry:      in.GetEntry(),
		Proxy:      SliceProxyDO2DTO(in.GetProxy()),
	}
}

func SliceAPPDO2DTO(ins []*models.APP) []*dto.APP {
	var out = make([]*dto.APP, 0, len(ins))

	for _, in := range ins {
		out = append(out, APPDO2DTO(in))
	}
	return out
}
