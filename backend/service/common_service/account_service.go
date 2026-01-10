package common_service

import (
	"context"
	"github.com/fansqz/fancode-backend/common"
	conf "github.com/fansqz/fancode-backend/common/config"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/file_store"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/dao"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
	"mime/multipart"
	"path"
)

const (
	// UserAvatarPath cos中，用户图片存储的位置
	UserAvatarPath = "/avatar/user"
)

type AccountService interface {
	// UploadAvatar 上传头像
	UploadAvatar(ctx context.Context, file *multipart.FileHeader) (string, error)
	// GetAccountInfo 获取账号信息
	GetAccountInfo(ctx context.Context) (*dto.AccountInfo, error)
	// UpdateAccountInfo 更新账号信息
	UpdateAccountInfo(ctx context.Context, user *po.SysUser) error
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, oldPassword, newPassword string) error
	// ResetPassword 重置密码
	ResetPassword(ctx context.Context) error
}

func NewAccountService(config *conf.AppConfig, userDao dao.SysUserDao, sysRoleDao dao.SysRoleDao) AccountService {
	return &accountService{
		config:     config,
		sysUserDao: userDao,
		sysRoleDao: sysRoleDao,
	}
}

type accountService struct {
	config     *conf.AppConfig
	sysUserDao dao.SysUserDao
	sysRoleDao dao.SysRoleDao
}

func (a *accountService) UploadAvatar(ctx context.Context, file *multipart.FileHeader) (string, error) {
	cos := file_store.NewImageCOS(a.config.COSConfig)
	fileName := file.Filename
	fileName = utils.GetUUID() + "." + path.Base(fileName)
	file2, err := file.Open()
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UploadAvatar] request error, err = %v", err)
		return "", e.ErrBadRequest
	}
	err = cos.SaveFile(UserAvatarPath+"/"+fileName, file2)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UploadAvatar] save test_file error, err = %v", err)
		return "", e.ErrServer
	}
	return path.Join(UserAvatarPath, fileName), nil
}

func (a *accountService) GetAccountInfo(ctx context.Context) (*dto.AccountInfo, error) {
	userID := utils.GetUserIDWithCtx(ctx)
	user, err := a.sysUserDao.GetUserByID(common.Mysql, userID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetAccountInfo] get user error, err = %v", err)
		return nil, e.ErrMysql
	}
	// 读取菜单
	for i := 0; i < len(user.Roles); i++ {
		user.Roles[i].Menus, err = a.sysRoleDao.GetMenusByRoleID(common.Mysql, user.Roles[i].ID)
		if err != nil {
			logger.WithCtx(ctx).Errorf("[PasswordLogin] get menus error, err = %v", err)
			return nil, e.ErrUserUnknownError
		}
	}
	return dto.NewAccountInfo(user), nil
}

func (a *accountService) UpdateAccountInfo(ctx context.Context, user *po.SysUser) error {
	user.ID = utils.GetUserIDWithCtx(ctx)
	// 不能更新账号名称和密码
	user.LoginName = ""
	user.Password = ""
	err := a.sysUserDao.UpdateUser(common.Mysql, user)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetAccountInfo] update user error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}

func (a *accountService) ChangePassword(ctx context.Context, oldPassword, newPassword string) error {
	userID := utils.GetUserIDWithCtx(ctx)
	//检验用户名
	user, err := a.sysUserDao.GetUserByID(common.Mysql, userID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[ChangePassword] get user error, err = %v", err)
		return e.ErrMysql
	}
	if user == nil || user.LoginName == "" {
		logger.WithCtx(ctx).Errorf("[ChangePassword] user not exist")
		return e.ErrUserNotExist
	}
	//检验旧密码
	if !utils.ComparePwd(oldPassword, user.Password) {
		return e.ErrUserNameOrPasswordWrong
	}
	password, err := utils.GetPwd(newPassword)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[ChangePassword] get password error, err = %v", err)
		return e.ErrPasswordEncodeFailed
	}
	user.Password = string(password)
	err = a.sysUserDao.UpdateUser(common.Mysql, user)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[ChangePassword] update user error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}

func (a *accountService) ResetPassword(ctx context.Context) error {
	userID := utils.GetUserIDWithCtx(ctx)
	password := utils.GetRandomPassword(11)
	password2, err := utils.GetPwd(password)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[ResetPassword] get password, err = %v", err)
		return e.ErrUserUnknownError
	}
	// 更新密码
	tx := common.Mysql.Begin()
	user := &po.SysUser{}
	user.ID = userID
	user.Password = string(password2)
	err = a.sysUserDao.UpdateUser(tx, user)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[ResetPassword] update user error, err = %v", err)
		tx.Rollback()
		return e.ErrMysql
	}
	// 发送密码
	message := common.EmailMessage{
		To:      []string{user.Email},
		Subject: "fancode-重置密码",
		Body:    "新密码：" + password,
	}
	err = common.SendMail(a.config.EmailConfig, message)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[ResetPassword] send mail error, err = %v", err)
		tx.Rollback()
		return e.ErrUserUnknownError
	}
	return nil
}
