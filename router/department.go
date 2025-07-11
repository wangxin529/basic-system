package router

import (
	"elevate-hub/logic/department"
	"elevate-hub/svc"
	"github.com/gin-gonic/gin"
)

func registerDepartmentRouter(v1 *gin.RouterGroup, svCtx *svc.ServiceContext) {
	deptLogic := department.NewDepartmentLogic(svCtx)
	dept := v1.Group("/department")
	{
		dept.GET("/:id", deptLogic.Department)
		dept.POST("", deptLogic.AddDepartment)
		dept.PUT("", deptLogic.UpdateDepartment)
		dept.DELETE("/:id", deptLogic.DelDepartment)
		dept.GET("", deptLogic.ListDepartment)
	}
}
