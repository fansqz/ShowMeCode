package admin

import (
	"github.com/fansqz/fancode-backend/controller/admin"
	"github.com/gin-gonic/gin"
)

func SetupVisualDocumentRoutes(r *gin.Engine, visualDocumentController admin.VisualDocumentManageController) {
	//题目相关路由
	document := r.Group("/manage/visual/document")
	{
		document.GET("/directory/:bankID", visualDocumentController.GetVisualDocumentDirectory)
		document.POST("/directory", visualDocumentController.UpdateVisualDocumentDirectory)
		document.GET("/:id", visualDocumentController.GetVisualDocumentByID)
		document.POST("", visualDocumentController.InsertVisualDocument)
		document.PUT("", visualDocumentController.UpdateVisualDocument)
		document.DELETE("/:id", visualDocumentController.DeleteVisualDocumentByID)
	}
}
