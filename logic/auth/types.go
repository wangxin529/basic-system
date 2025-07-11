package auth

import "elevate-hub/db/meta"

type ListAppReq struct {
	meta.Page
}

type ListAuthReq struct {
	MenuID int64 `json:"menuId" form:"menuId"`
	APPID  int64 `json:"appId" form:"appId"`
}
type AddAuthReq struct {
	MenuID int64   `json:"menuId"`
	APPID  int64   `json:"appId"`
	AppIds []int64 `json:"apiIds"`
}
