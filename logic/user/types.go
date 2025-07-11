package user

import "elevate-hub/db/meta"

type AddUserReq struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Nickname      string `json:"nickname"`
	PasswordAgain string `json:"password_again"`
	Sex           int    `json:"sex"`
	//IsSupper      int    `json:"is_supper"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ListUserReq struct {
	meta.Page
}

type UpdateUserReq struct {
	ID       int64   `json:"id"`
	Nickname string  `json:"nickname"`
	Sex      int     `json:"sex"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	IsSupper int     `json:"is_supper"`
	Status   int     `json:"status"`
	Roles    []int64 `json:"roles"`
}
