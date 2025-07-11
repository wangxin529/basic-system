package router

import (
	"elevate-hub/logic/user"
	"elevate-hub/svc"

	"github.com/gin-gonic/gin"
)

func registerUserRouter(v1 *gin.RouterGroup, svCtx *svc.ServiceContext) {
	logic := user.NewUserLogic(svCtx)
	dept := v1.Group("/user")
	{
		dept.GET("/:id", logic.UserInfo)
		dept.GET("/current", logic.CurrentUser)
		dept.POST("", logic.AddUser)
		dept.PUT("", logic.UpdateUser)
		dept.DELETE("/:id", logic.DelUser)
		dept.GET("", logic.ListUser)
	}
}
