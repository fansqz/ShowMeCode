package routers

import (
	"github.com/fansqz/fancode-backend/controller"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine, authController controller.AuthController) {
	//用户相关
	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.UserRegister)
		auth.POST("/code/send", authController.SendAuthCode)
	}
}
