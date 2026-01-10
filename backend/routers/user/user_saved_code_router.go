package user

import (
	c "github.com/fansqz/fancode-backend/controller/user"
	"github.com/gin-gonic/gin"
)

// SetupUserSavedCodeRoutes 设置用户保存代码相关路由
func SetupUserSavedCodeRoutes(r *gin.Engine, handler *c.UserSavedCodeHandler) {
	// 用户保存代码相关路由组
	savedCodeGroup := r.Group("/user/savedCode")
	{
		// 创建用户保存的代码
		savedCodeGroup.POST("", handler.CreateUserSavedCode)
		// 更新用户保存的代码
		savedCodeGroup.PUT("", handler.UpdateUserSavedCode)
		// 获取用户保存的代码列表
		savedCodeGroup.GET("list", handler.GetUserSavedCodeList)
		// 根据ID获取用户保存的代码详情
		savedCodeGroup.GET("/:id", handler.GetUserSavedCodeByID)
		// 删除用户保存的代码
		savedCodeGroup.DELETE("/:id", handler.DeleteUserSavedCode)
	}
}
