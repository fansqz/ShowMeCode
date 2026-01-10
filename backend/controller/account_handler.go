package controller

import (
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/models/po"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/common_service"
	"time"

	"github.com/gin-gonic/gin"
)

// AccountController 关于一些账号信息的handler
type AccountController interface {
	// UploadAvatar 上传头像
	UploadAvatar(ctx *gin.Context)
	// GetAccountInfo 获取账号信息
	GetAccountInfo(ctx *gin.Context)
	// UpdateAccountInfo 更新账号信息
	UpdateAccountInfo(ctx *gin.Context)
	// ChangePassword 修改密码
	ChangePassword(ctx *gin.Context)
	// ResetPassword 重置密码
	ResetPassword(ctx *gin.Context)
}

type accountController struct {
	accountService common_service.AccountService
}

func NewAccountController(accountService common_service.AccountService) AccountController {
	return &accountController{
		accountService: accountService,
	}
}

func (a *accountController) UploadAvatar(ctx *gin.Context) {
	result := r.NewResult(ctx)
	file, err := ctx.FormFile("avatar")
	if err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if file.Size > 2<<20 {
		result.SimpleErrorMessage("文件大小不能超过2m")
		return
	}
	path, err2 := a.accountService.UploadAvatar(ctx, file)
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(path)
}

func (a *accountController) GetAccountInfo(ctx *gin.Context) {
	result := r.NewResult(ctx)
	accountInfo, err := a.accountService.GetAccountInfo(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(accountInfo)
}

func (a *accountController) UpdateAccountInfo(ctx *gin.Context) {
	result := r.NewResult(ctx)
	user := &po.SysUser{}
	user.Avatar = ctx.PostForm("avatar")
	user.Username = ctx.PostForm("username")
	user.Introduction = ctx.PostForm("introduction")
	sex := ctx.PostForm("sex")
	if sex == "2" {
		user.Sex = 2
	} else if sex == "1" {
		user.Sex = 1
	}
	birthDay := ctx.PostForm("birthDay")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", birthDay, time.Local)
	if err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	user.BirthDay = t
	if err2 := a.accountService.UpdateAccountInfo(ctx, user); err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("提交成功，重新登录可更新数据")
}

func (a *accountController) ChangePassword(ctx *gin.Context) {
	result := r.NewResult(ctx)
	oldPassword := ctx.PostForm("oldPassword")
	newPassword := ctx.PostForm("newPassword")
	if oldPassword == "" {
		result.Error(e.ErrBadRequest)
		return
	}
	if newPassword == "" {
		result.Error(e.ErrBadRequest)
	}
	if err := a.accountService.ChangePassword(ctx, oldPassword, newPassword); err != nil {
		result.Error(err)
	}
	result.SuccessMessage("修改成功，请重新登录")
}

func (a *accountController) ResetPassword(ctx *gin.Context) {
	result := r.NewResult(ctx)
	err := a.accountService.ResetPassword(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("重置成功，请留意邮箱")
}
