package user

import (
	"github.com/fansqz/fancode-backend/controller/user"
	"github.com/gin-gonic/gin"
)

func SetupVisualDocumentBankRoutes(r *gin.Engine, v user.VisualDocumentBankController) {
	// 知识库相关路由
	bank := r.Group("/learn/visual/document/bank")
	{
		bank.GET("/all", v.GetAllVisualDocumentBank)
		bank.GET("/:id", v.GetVisualDocumentBankByID)
	}
}
