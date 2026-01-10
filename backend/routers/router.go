// Package routers
// @Author: fzw
// @Create: 2023/7/14
// @Description: 路由相关
package routers

import (
	conf "github.com/fansqz/fancode-backend/common/config"
	c "github.com/fansqz/fancode-backend/controller"
	"github.com/fansqz/fancode-backend/controller/admin"
	"github.com/fansqz/fancode-backend/controller/user"
	"github.com/fansqz/fancode-backend/interceptor"
	adminRouter "github.com/fansqz/fancode-backend/routers/admin"
	userRouter "github.com/fansqz/fancode-backend/routers/user"
	"github.com/gin-gonic/gin"
)

// SetupRouter
//
//	@Description: 启动路由
func SetupRouter(
	authController c.AuthController,
	accountController c.AccountController,
	commonController c.CommonController,
	userSavedCodeHandler *user.UserSavedCodeHandler,
	sysApiController admin.SysApiController,
	sysMenuController admin.SysMenuController,
	sysRoleController admin.SysRoleController,
	sysUserController admin.SysUserController,
	visualDocumentManageController admin.VisualDocumentManageController,
	visualDocumentBankManageController admin.VisualDocumentBankManageController,
	debugController user.DebugController,
	visualController user.VisualController,
	visualDocumentController user.VisualDocumentController,
	visualDocumentBankController user.VisualDocumentBankController,
	config *conf.AppConfig,
	panicInterceptor *interceptor.RecoverPanicInterceptor,
	corsInterceptor *interceptor.CorsInterceptor,
	requestInterceptor *interceptor.RequestInterceptor,
	loggerInterceptor *interceptor.LoggerInterceptor,
) *gin.Engine {
	if config.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(loggerInterceptor.LoggerInterceptor())
	// 拦截panic
	r.Use(panicInterceptor.RecoverPanic())
	// 允许跨域
	r.Use(corsInterceptor.Cors())
	// 拦截非法用户
	r.Use(interceptor.VisitorUIDInterceptor())
	r.Use(requestInterceptor.TokenAuthorize())
	// 限流器
	r.Use(interceptor.RateLimitInterceptor())

	//设置静态文件位置
	r.Static("/static", "/")

	//ping
	r.GET("/ping", c.Ping)

	SetupAuthRoutes(r, authController)
	SetupAccountRoutes(r, accountController)
	SetupCommonRoutes(r, commonController)
	userRouter.SetupUserSavedCodeRoutes(r, userSavedCodeHandler)
	adminRouter.SetupSysApiRoutes(r, sysApiController)
	adminRouter.SetupSysMenuRoutes(r, sysMenuController)
	adminRouter.SetupSysRoleRoutes(r, sysRoleController)
	adminRouter.SetupSysUserRoutes(r, sysUserController)
	adminRouter.SetupVisualDocumentRoutes(r, visualDocumentManageController)
	adminRouter.SetupVisualDocumentBankRoutes(r, visualDocumentBankManageController)
	userRouter.SetupDebugRoutes(r, debugController)
	userRouter.SetupVisualRoutes(r, visualController)
	userRouter.SetupVisualDocumentRoutes(r, visualDocumentController)
	userRouter.SetupVisualDocumentBankRoutes(r, visualDocumentBankController)
	return r
}
