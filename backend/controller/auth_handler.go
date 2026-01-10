package controller

import (
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/models/po"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/common_service"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	// Login 用户登录
	Login(ctx *gin.Context)
	// SendAuthCode 发送验证码
	SendAuthCode(ctx *gin.Context)
	// UserRegister 用户注册
	UserRegister(ctx *gin.Context)
}

type authController struct {
	authService common_service.AuthService
}

func NewAuthController(authService common_service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (u *authController) SendAuthCode(ctx *gin.Context) {
	result := r.NewResult(ctx)
	email := ctx.PostForm("email")
	kind := ctx.PostForm("type")
	// 生成code
	if err := u.authService.SendAuthCode(ctx, email, kind); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("验证码发送成功")
}

func (u *authController) UserRegister(ctx *gin.Context) {
	result := r.NewResult(ctx)
	user := &po.SysUser{
		Email:    ctx.PostForm("email"),
		Username: ctx.PostForm("username"),
		Password: ctx.PostForm("password"),
	}
	code := ctx.PostForm("code")
	if err := u.authService.UserRegister(ctx, user, code); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("注册成功")
}

func (u *authController) Login(ctx *gin.Context) {
	result := r.NewResult(ctx)
	//获取并检验用户参数
	kind := ctx.PostForm("type")
	if kind == "password" {
		u.passwordLogin(ctx)
	} else if kind == "email" {
		u.emailLogin(ctx)
	} else {
		result.Error(e.ErrLoginType)
	}
}

// passportLogin 密码登陆
func (a *authController) passwordLogin(ctx *gin.Context) {
	result := r.NewResult(ctx)
	account := ctx.PostForm("account")
	password := ctx.PostForm("password")
	if account == "" || password == "" {
		result.Error(e.ErrBadRequest)
		return
	}
	token, err := a.authService.PasswordLogin(ctx, account, password)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(token)
}

// emailLogin 邮箱登陆
func (a *authController) emailLogin(ctx *gin.Context) {
	result := r.NewResult(ctx)
	email := ctx.PostForm("email")
	code := ctx.PostForm("code")
	if email == "" || code == "" {
		result.Error(e.ErrBadRequest)
		return
	}
	// 登录
	token, err := a.authService.EmailLogin(ctx, email, code)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(token)
}
