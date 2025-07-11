package department

import "elevate-hub/db/meta"

type ListDepartmentReq struct {
	meta.Page
}

type AddDepartmentReq struct {
	Name   string `json:"name"`
	Leader int64  `json:"leader"`
	Parent int64  `json:"parent"`
}

type UpdateDepartmentReq struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Leader int64  `json:"leader"`
	Parent int64  `json:"parent"`
}
