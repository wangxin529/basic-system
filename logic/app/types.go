package app

import (
	"elevate-hub/db/meta"
	"elevate-hub/db/models"
)

type ListAppReq struct {
	meta.Page
}

type AddAppReq struct {
	Name       string `json:"name"`
	Key        string `json:"key"`
	SignMethod string `json:"sign_method"`
	Status     int    `json:"status"`
}

type UpdateAppReq struct {
	ID     int64
	Status int `json:"status"`
}

type UpdateProxyReq struct {
	ID    int64          `json:"id"`
	Proxy []models.Proxy `json:"proxy"`
}

type UpdateAppEntry struct {
	ID int64 `json:"id"`
	models.Entry
}
