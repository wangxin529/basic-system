package meta

type Page struct {
	Page     int `json:"current" form:"current"`
	PageSize int `json:"pageSize" form:"pageSize"`
	//Total    int `json:"total"`
}

var (
	EmptyListOption = ListOption{
		Page:         nil,
		DisableCount: true,
	}
	EmptyCountOption = CountOption{}
	EmptyPage        = &Page{
		Page: -1,
	}
)

type GetOption struct {
	ID            int64
	Condition     interface{} // 查询条件 传入model 数据 str, map, model
	InCondition   map[string]interface{}
	LikeCondition map[string]interface{}
	IsLast        bool
}
type Between struct {
	Field  string
	Value1 interface{}
	Value2 interface{}
}
type Order struct {
	Filed string `json:"filed"`
	Sort  string `json:"sort,default=asc"`
}
type CustomCondition struct {
	SQL  string
	Args []interface{}
}
type ListOption struct {
	DisableCount    bool // 查询总量
	Page            *Page
	Condition       map[string]interface{}
	LikeCondition   map[string]interface{}
	InCondition     map[string]interface{}
	CustomCondition *CustomCondition
	Select          string
	IDs             []int64
	Between         *Between
	Order           *Order
	Join            string
}

type DeleteOption struct {
	ID          int64
	IDs         []int64
	Condition   interface{}
	InCondition map[string]interface{}
}

type UpdateOption struct {
	ID        int64
	Condition map[string]interface{}
	Data      interface{}
}

type CreateOption[T any] struct {
	Data T
}

type CountOption struct {
	IsDelete bool
}
