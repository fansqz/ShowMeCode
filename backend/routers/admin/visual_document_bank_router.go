package admin

import (
	"github.com/fansqz/fancode-backend/controller/admin"
	"github.com/gin-gonic/gin"
)

func SetupVisualDocumentBankRoutes(r *gin.Engine, v admin.VisualDocumentBankManageController) {
	// 知识库相关路由
	document := r.Group("/manage/visual/document/bank")
	{
		document.POST("", v.InsertVisualDocumentBank)
		document.PUT("", v.UpdateVisualDocumentBank)
		document.DELETE("/:id", v.DeleteVisualDocumentBank)
		document.GET("/all", v.GetAllVisualDocumentBank)
		document.GET("/:id", v.GetVisualDocumentBankByID)
	}
}
