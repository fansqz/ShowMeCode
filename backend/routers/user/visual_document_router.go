package user

import (
	"github.com/fansqz/fancode-backend/controller/user"
	"github.com/gin-gonic/gin"
)

func SetupVisualDocumentRoutes(r *gin.Engine, visualDocumentController user.VisualDocumentController) {
	//题目相关路由
	document := r.Group("/learn/visual/document")
	{
		document.GET("/directory/:bankID", visualDocumentController.GetVisualDocumentDirectory)
		document.GET("/:id", visualDocumentController.GetVisualDocumentByID)
	}
}
