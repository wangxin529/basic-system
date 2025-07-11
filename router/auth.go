package router

import (
	"elevate-hub/logic/auth"
	"elevate-hub/svc"
	"github.com/gin-gonic/gin"
)

func registerAuthRouter(v1 *gin.RouterGroup, svCtx *svc.ServiceContext) {
	authLogic := auth.NewAuthLogic(svCtx)
	appRouter := v1.Group("/auth")
	{
		appRouter.POST("/api", authLogic.AddAPI)
		appRouter.GET("/api", authLogic.ListAPI)
	}
}
