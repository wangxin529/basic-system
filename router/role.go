package router

import (
	"elevate-hub/logic/role"
	"elevate-hub/svc"

	"github.com/gin-gonic/gin"
)

func registerRoleRouter(v1 *gin.RouterGroup, svCtx *svc.ServiceContext) {
	logic := role.NewRoleLogic(svCtx)
	dept := v1.Group("/role")
	{
		dept.GET("/:id", logic.Role)
		dept.POST("", logic.AddRole)
		dept.PUT("", logic.UpdateRole)
		dept.DELETE("/:id", logic.DelRole)
		dept.GET("", logic.ListRole)
	}
}
