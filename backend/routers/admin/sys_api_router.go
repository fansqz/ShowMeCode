package admin

import (
	"github.com/fansqz/fancode-backend/controller/admin"
	"github.com/gin-gonic/gin"
)

func SetupSysApiRoutes(r *gin.Engine, apiController admin.SysApiController) {
	//题目相关路由
	api := r.Group("/manage/api")
	{
		api.GET("/:id", apiController.GetApiByID)
		api.POST("", apiController.InsertApi)
		api.PUT("", apiController.UpdateApi)
		api.DELETE("/:id", apiController.DeleteApiByID)
		api.GET("/tree", apiController.GetApiTree)
	}
}
