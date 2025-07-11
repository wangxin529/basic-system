package dto

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	NickName string `json:"nickname"`
	IsSupper int    `json:"is_supper"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Sex      int    `json:"sex"`
	Email    string `json:"email"`
	Dept     string `json:"dept"`
	DeptID   int64  `json:"dept_id"`
	Post     string `json:"post"`
	Status   int    `json:"status"`
}

type Role struct {
	ID                 int64                `json:"id"`
	Name               string               `json:"name"`
	Status             int                  `json:"status"`
	OperatePermissions []OperatePermissions `json:"operate_permission"`
}
type OperatePermissions struct {
	MenuID            int64    `json:"menu_id"`
	ButtonPermissions []string `json:"button_permission"`
}

type Menu struct {
	ID         int64       `json:"key"`   //`json:"id"`
	Name       string      `json:"title"` //`json:"name"`
	Value      int64       `json:"value"`
	Path       string      `json:"path"`
	Status     int         `json:"status"`
	Parent     int64       `json:"parent"`
	MenuType   int         `json:"menu_type"`
	Buttons    []*Button   `json:"buttons"`
	Children   []*Menu     `json:"children"`
	MenuConfig *MenuConfig `json:"menu_config"`
}

type Department struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	Parent   int64         `json:"parent"`
	Leader   string        `json:"leader"`
	Children []*Department `json:"children"`
}
type Proxy struct {
	Prefix string `json:"prefix"`
	Addr   string `json:"addr"`
}

type APP struct {
	ID         int64       `json:"id"`
	Name       string      `json:"name" `
	Key        string      `json:"key" `
	Status     int         `json:"status"`
	SignMethod string      `json:"sign_method"` // 加密算法     md5
	Entry      interface{} `json:"entry"`
	Proxy      []Proxy     `json:"proxy"`
}

type API struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	//Key      string `json:"key"`
	MenuID   int64  `json:"menu_id"`
	Status   int    `json:"status"`
	APIType  int    `json:"api_type"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Creator  string `json:"creator"`
	Describe string `json:"describe"`
	CreateAt string `json:"create_at"`
}

type MenuBrief struct {
	ID         string      `json:"id"`
	Name       string      `json:"title"`
	Path       string      `json:"path"`
	Parent     string      `json:"parent"`
	MenuType   int         `json:"menu_type"`
	MenuConfig *MenuConfig `json:"menu_config"`
}
