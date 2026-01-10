package admin

import (
	"github.com/fansqz/fancode-backend/controller/admin"
	"github.com/gin-gonic/gin"
)

func SetupSysMenuRoutes(r *gin.Engine, menuController admin.SysMenuController) {
	//题目相关路由
	menu := r.Group("/manage/menu")
	{
		menu.GET("/:id", menuController.GetMenuByID)
		menu.POST("", menuController.InsertMenu)
		menu.PUT("", menuController.UpdateMenu)
		menu.DELETE("/:id", menuController.DeleteMenuByID)
		menu.GET("/tree", menuController.GetMenuTree)
	}
}
