package router

import (
	"elevate-hub/logic/api"
	"elevate-hub/svc"
	"github.com/gin-gonic/gin"
)

func registerApiRouter(v1 *gin.RouterGroup, svCtx *svc.ServiceContext) {
	apiLogic := api.NewAPILogic(svCtx)
	dept := v1.Group("/api")
	{
		dept.GET("/:id", apiLogic.API)
		dept.POST("", apiLogic.AddAPI)
		dept.PUT("", apiLogic.UpdateAPI)
		dept.DELETE("/:id", apiLogic.DelAPI)
		dept.GET("", apiLogic.ListAPI)
	}
}
