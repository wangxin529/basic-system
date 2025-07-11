package router

import (
	"elevate-hub/logic/menu"
	"elevate-hub/svc"

	"github.com/gin-gonic/gin"
)

func registerMenuRouter(v1 *gin.RouterGroup, svCtx *svc.ServiceContext) {
	logic := menu.NewMenuLogic(svCtx)
	dept := v1.Group("/menu")
	{
		dept.GET("/:id", logic.Menu)
		dept.GET("/current", logic.CurrentUserMenu)
		dept.POST("", logic.AddMenu)
		dept.PUT("", logic.UpdateMenu)
		dept.DELETE("/:id", logic.DelMenu)
		dept.GET("", logic.ListMenu)
		dept.GET("/tree", logic.MenuTree)
	}
}
