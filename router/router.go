package router

import (
	"elevate-hub/middleware"
	"elevate-hub/middleware/application"
	"elevate-hub/middleware/jwt"
	"elevate-hub/middleware/proxy"
	"elevate-hub/svc"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func Router(g *gin.Engine, ctx *svc.ServiceContext) {
	authMiddleware, err := middleware.AuthInit(ctx.Config, ctx)
	if err != nil {
		log.Fatalf("JWT Init Error, %s", err.Error())
	}
	// 注册登录路由
	noCheckRouter(g, ctx, authMiddleware)
	checkRouter(g, ctx, authMiddleware)

}

func noCheckRouter(g *gin.Engine, ctx *svc.ServiceContext, jwtMiddleware *jwt.GinJWTMiddleware) {
	g.POST("/basic/api/v1/login", jwtMiddleware.LoginHandler)
}

func checkRouter(g *gin.Engine, ctx *svc.ServiceContext, jwtMiddleware *jwt.GinJWTMiddleware) {
	// 注册路由
	am := application.NewApplicationMiddleware(ctx)
	permission := middleware.NewPermission(ctx)
	fmt.Println(permission)
	pm := proxy.NewProxy(ctx)
	g.Use(am.Handle()).Use(jwtMiddleware.MiddlewareFunc()) //.Use(middleware.AuthCheckRole())
	g.NoRoute(pm.ServeHTTP())
	cg := g.Group("/basic/api/v1")
	{
		cg.POST("/logout", jwtMiddleware.LogOut)
		registerDepartmentRouter(cg, ctx)
		registerUserRouter(cg, ctx)
		registerApiRouter(cg, ctx)
		registerAppRouter(cg, ctx)
		registerMenuRouter(cg, ctx)
		registerRoleRouter(cg, ctx)
		registerAuthRouter(cg, ctx)
	}
}
