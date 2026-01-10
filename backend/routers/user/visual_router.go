package user

import (
	"github.com/fansqz/fancode-backend/controller/user"

	"github.com/gin-gonic/gin"
)

func SetupVisualRoutes(r *gin.Engine, visualController user.VisualController) {
	//用户相关
	visual := r.Group("/visual/")
	{
		visual.POST("/debug/struct", visualController.StructVisual)
		visual.POST("/debug/array", visualController.ArrayVisual)
		visual.POST("/debug/array2d", visualController.Array2DVisual)
		visual.GET("/debug/description/:id", visualController.GetVisualDescription)
	}
}
