package api

import "elevate-hub/db/meta"

type ListAPIReq struct {
	meta.Page
	APIType *int64 `form:"api_type"`
}

type AddAPIReq struct {
	Name     string `json:"name"`
	key      string `json:"key"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Status   int    `json:"status"`
	Describe string `json:"describe"`
	APIType  int    `json:"api_type"`
}

type UpdateAPIReq struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Status   int    `json:"status"`
	Describe string `json:"describe"`
}
