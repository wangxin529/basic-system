package models

type DepartmentAggregate struct {
	Department
	LeaderUser     string `json:"leader_user"`
	LeaderNickname string `json:"leader_nickname"`
	LeaderEmail    string `json:"leader_email"`
}

type UserAggregate struct {
	User
	DepartmentName string `json:"department_name"`
}

type APIAggregate struct {
	API
	CreatorUserName string `json:"creator_username"`
}
