package routers

import (
	"github.com/fansqz/fancode-backend/controller"
	"github.com/gin-gonic/gin"
)

func SetupAccountRoutes(r *gin.Engine, accountController controller.AccountController) {
	account := r.Group("/account")
	{
		account.GET("/info", accountController.GetAccountInfo)
		account.PUT("", accountController.UpdateAccountInfo)
		account.POST("/password/reset", accountController.ResetPassword)
		account.POST("/password", accountController.ChangePassword)
	}
	avatar := r.Group("/avatar")
	{
		avatar.POST("", accountController.UploadAvatar)
	}
}
