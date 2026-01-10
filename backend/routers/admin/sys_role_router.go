package admin

import (
	"github.com/fansqz/fancode-backend/controller/admin"
	"github.com/gin-gonic/gin"
)

func SetupSysRoleRoutes(r *gin.Engine, roleController admin.SysRoleController) {
	//题目相关路由
	role := r.Group("/manage/role")
	{
		role.GET("/:id", roleController.GetRoleByID)
		role.POST("", roleController.InsertSysRole)
		role.PUT("", roleController.UpdateSysRole)
		role.DELETE("/:id", roleController.DeleteSysRole)
		role.GET("/list", roleController.GetSysRoleList)
		role.GET("/api/:id", roleController.GetApiIDsByRoleID)
		role.GET("/menu/:id", roleController.GetMenuIDsByRoleID)
		role.PUT("/api", roleController.UpdateRoleApis)
		role.PUT("/menu", roleController.UpdateRoleMenus)
	}
}
