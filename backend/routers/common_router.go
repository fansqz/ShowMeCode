package routers

import (
	"github.com/fansqz/fancode-backend/controller"
	"github.com/gin-gonic/gin"
)

func SetupCommonRoutes(r *gin.Engine, commonController controller.CommonController) {
	//用户相关
	common := r.Group("/common")
	{
		common.GET("/getURL", commonController.GetURL)
	}
}
