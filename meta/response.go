package meta

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"` // 0 正常 -1 异常
	Message string `json:"message"`
}
type ErrorResponse struct {
	Response
}
type SuccessResponse struct {
	Response
	Data interface{} `json:"data"`
}

type SuccessResponseAndTotal struct {
	Response
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}

func ErrHandle(c *gin.Context, err error) {
	c.JSON(200, ErrorResponse{
		Response{
			Code:    -1,
			Message: err.Error(),
		}})

}
func ErrHandleWithMsg(c *gin.Context, err string) {
	c.JSON(200, ErrorResponse{
		Response{
			Code:    -1,
			Message: err,
		}})

}
func ErrHandleWithHttpCodeAndMsg(c *gin.Context, httpCode int, err string) {
	c.JSON(200, ErrorResponse{
		Response{
			Code:    -1,
			Message: err,
		}})

}

func SuccessHandle(c *gin.Context, data interface{}) {
	c.JSON(200, SuccessResponse{
		Response: Response{
			Code:    0,
			Message: "success",
		},
		Data: data,
	})

}

func SuccessHandleAndTotal(c *gin.Context, data interface{}, total int64) {
	c.JSON(200, SuccessResponseAndTotal{
		Response: Response{
			Code:    0,
			Message: "success",
		},
		Data:  data,
		Total: total,
	})

}
