package router

import (
	"elevate-hub/logic/app"
	"elevate-hub/svc"
	"github.com/gin-gonic/gin"
)

func registerAppRouter(v1 *gin.RouterGroup, svCtx *svc.ServiceContext) {
	appLogic := app.NewAppLogic(svCtx)
	appRouter := v1.Group("/app")
	{
		appRouter.GET("/entry", appLogic.EntryList)
		appRouter.PUT("/entry", appLogic.UpdateEntry)
		appRouter.GET("/:id", appLogic.App)
		appRouter.POST("", appLogic.AddApp)
		appRouter.PUT("", appLogic.UpdateApp)
		appRouter.DELETE("/:id", appLogic.DelApp)
		appRouter.GET("", appLogic.ListApp)
		appRouter.PUT("/proxy", appLogic.UpdateProxy)
		appRouter.GET("/proxy/:id", appLogic.ListProxy)
	}
}
