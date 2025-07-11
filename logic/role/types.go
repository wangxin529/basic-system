package role

import "elevate-hub/db/meta"

type ListRoleReq struct {
	meta.Page
}

type AddRoleReq struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type UpdateRoleReq struct {
	ID                int64              `json:"id"`
	Name              string             `json:"name"`
	Status            int                `json:"status"`
	OperatePermission map[int64][]string `json:"operate_permission"`
}
