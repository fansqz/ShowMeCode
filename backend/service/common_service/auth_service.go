package common_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/fansqz/fancode-backend/common"
	conf "github.com/fansqz/fancode-backend/common/config"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/constants"
	"github.com/fansqz/fancode-backend/dao"
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
	"gorm.io/gorm"
	"net/mail"
	"strings"
	"time"
)

const (
	RegisterEmailProKey = "emailcode-register-"
	LoginEmailProKey    = "emailcode-login-"
)

type AuthService interface {

	// PasswordLogin 密码登录 account可能是邮箱可能是用户id
	PasswordLogin(ctx context.Context, account string, password string) (string, error)
	// EmailLogin 邮箱验证登录
	EmailLogin(ctx context.Context, email string, code string) (string, error)
	// SendAuthCode 获取邮件的验证码
	SendAuthCode(ctx context.Context, email string, kind string) error
	// UserRegister 用户注册
	UserRegister(ctx context.Context, user *po.SysUser, code string) error
}

type authService struct {
	config     *conf.AppConfig
	sysUserDao dao.SysUserDao
	sysMenuDao dao.SysMenuDao
	sysRoleDao dao.SysRoleDao
}

func NewAuthService(config *conf.AppConfig, userDao dao.SysUserDao, menuDao dao.SysMenuDao, roleDao dao.SysRoleDao) AuthService {
	return &authService{
		config:     config,
		sysUserDao: userDao,
		sysMenuDao: menuDao,
		sysRoleDao: roleDao,
	}
}

func (u *authService) PasswordLogin(ctx context.Context, account string, password string) (string, error) {
	var user *po.SysUser
	var err error
	if verifyEmailFormat(account) {
		email := strings.ToLower(account)
		user, err = u.sysUserDao.GetUserByEmail(common.Mysql, email)
	} else {
		user, err = u.sysUserDao.GetUserByLoginName(common.Mysql, account)
	}
	if user == nil || err == gorm.ErrRecordNotFound {
		logger.WithCtx(ctx).Infof("[PasswordLogin] user not exist")
		return "", e.ErrUserNotExist
	}
	if err != nil {
		logger.WithCtx(ctx).Errorf("[PasswordLogin] get user error, err = %v", err)
		return "", e.ErrUserUnknownError
	}
	// 比较密码
	if !utils.ComparePwd(user.Password, password) {
		return "", e.ErrUserNameOrPasswordWrong
	}
	token, err := utils.GenerateToken(utils.Claims{
		ID: user.ID,
	})
	if err != nil {
		logger.WithCtx(ctx).Errorf("[PasswordLogin] generate token error, err = %v", err)
		return "", e.ErrUnknown
	}
	return token, nil
}

func (u *authService) EmailLogin(ctx context.Context, email string, code string) (string, error) {
	email = strings.ToLower(email)
	// 获取用户
	user, err := u.sysUserDao.GetUserByEmail(common.Mysql, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.WithCtx(ctx).Infof("[EmailLogin] user not exist")
		return "", e.ErrEmailNotRegister
	}
	if err != nil {
		logger.WithCtx(ctx).Errorf("[EmailLogin] get user error, err = %v", err)
		return "", e.ErrUnknown
	}
	// 检测验证码
	key := LoginEmailProKey + email
	result, err := common.Redis.Get(key).Result()
	if err != nil {
		logger.WithCtx(ctx).Errorf("[EmailLogin] get code from redis fail, err = %v", err)
		return "", e.ErrUnknown
	}
	if result != code {
		return "", e.ErrLoginCodeWrong
	}
	token, err := utils.GenerateToken(utils.Claims{
		ID: user.ID,
	})
	if err != nil {
		logger.WithCtx(ctx).Errorf("[EmailLogin] generate token error, err = %v", err)
		return "", e.ErrUserUnknownError
	}
	return token, nil
}

// verifyEmailFormat email verify
func verifyEmailFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (u *authService) SendAuthCode(ctx context.Context, email string, kind string) error {
	if !verifyEmailFormat(email) {
		// 邮箱格式有问题
		return e.ErrEmailFormatWrong
	}
	email = strings.ToLower(email)
	if kind == "register" {
		return u.sendRegisterCode(ctx, email)
	} else if kind == "login" {
		return u.sendLoginCode(ctx, email)
	}
	return e.ErrBadRequest
}

// sendLoginEmail 发送登陆验证码
func (u *authService) sendLoginCode(ctx context.Context, email string) error {
	body := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ShowMeCode验证码</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f9f9f9;
            margin: 0;
            padding: 0;
        }
       .email-container {
            background-color: #fff;
            max-width: 550px;
            margin: 30px auto;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }
       .email-content {
            padding: 30px;
        }
       .email-content p {
            font-size: 16px;
            line-height: 1.6;
            color: #333;
            margin-bottom: 15px;
        }
       .code-block {
            background-color: #e7f3ff;
            border-left: 6px solid #007BFF;
            padding: 20px;
            text-align: center;
            font-size: 32px;
            font-weight: bold;
            color: #007BFF;
            margin: 20px 0;
        }
       .email-footer {
            background-color: #f2f2f2;
            text-align: center;
            padding: 20px;
            font-size: 14px;
            color: #777;
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="email-content">
            <p>尊敬的用户，您好！</p>
            <p>您正在尝试登录 ShowMeCode，为了确保您的账户安全，请使用以下验证码完成登录。请不要将验证码泄露给他人。</p>
            <div class="code-block">
                %s
            </div>
            <p>验证码的有效期为 10 分钟。如果这不是您本人的操作，请忽略此邮件。</p>
        </div>
        <div class="email-footer">
            <p>感谢您使用 ShowMeCode！</p>
        </div>
    </div>
</body>
</html>`
	// 校验email是否存在
	exist, err := u.sysUserDao.CheckEmailIsExist(common.Mysql, email)
	if err != nil {
		logger.WithCtx(ctx).Infof("[sendLoginCode] check email fail, err = %v", err)
		return err
	}
	// 邮箱不存在
	if !exist {
		return e.ErrEmailNotRegister
	}
	// 发送code
	code := utils.GetRandomNumber(6)
	message := common.EmailMessage{
		To:      []string{email},
		Subject: "ShowMeCode登录验证码",
		Body:    fmt.Sprintf(body, code),
	}
	// 直接返回成功
	if err = common.SendMail(u.config.EmailConfig, message); err != nil {
		logger.WithCtx(ctx).Errorf("[SendAuthCode] send mail error, err = %v", err)
	}
	key := LoginEmailProKey + email
	_, err = common.Redis.Set(key, code, 10*time.Minute).Result()
	if err != nil {
		logger.WithCtx(ctx).Errorf("[SendAuthCode] set redis error, err = %v", err)
		return e.ErrUnknown
	}
	return nil
}

// sendRegisterEmail 发送注册邮箱
func (u *authService) sendRegisterCode(ctx context.Context, email string) error {
	body := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ShowMeCode验证码</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f9f9f9;
            margin: 0;
            padding: 0;
        }
       .email-container {
            background-color: #fff;
            max-width: 550px;
            margin: 30px auto;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }
       .email-content {
            padding: 30px;
        }
       .email-content p {
            font-size: 16px;
            line-height: 1.6;
            color: #333;
            margin-bottom: 15px;
        }
       .code-block {
            background-color: #e7f3ff;
            border-left: 6px solid #007BFF;
            padding: 20px;
            text-align: center;
            font-size: 32px;
            font-weight: bold;
            color: #007BFF;
            margin: 20px 0;
        }
       .email-footer {
            background-color: #f2f2f2;
            text-align: center;
            padding: 20px;
            font-size: 14px;
            color: #777;
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="email-content">
            <p>尊敬的用户，您好！</p>
            <p>您正在尝试注册 ShowMeCode，为了确保您的账户安全，请使用以下验证码完成注册。请不要将验证码泄露给他人。</p>
            <div class="code-block">
                %s
            </div>
            <p>验证码的有效期为 10 分钟。如果这不是您本人的操作，请忽略此邮件。</p>
        </div>
        <div class="email-footer">
            <p>感谢您使用 ShowMeCode！</p>
        </div>
    </div>
</body>
</html>`
	// 校验email是否存在
	exist, err := u.sysUserDao.CheckEmailIsExist(common.Mysql, email)
	if err != nil {
		logger.WithCtx(ctx).Infof("[sendLoginCode] check email fail, err = %v", err)
		return err
	}
	if exist {
		return e.ErrUserEmailIsExist
	}
	// 发送code
	code := utils.GetRandomNumber(6)
	message := common.EmailMessage{
		To:      []string{email},
		Subject: "ShowMeCode注册验证码",
		Body:    fmt.Sprintf(body, code),
	}
	if err = common.SendMail(u.config.EmailConfig, message); err != nil {
		logger.WithCtx(ctx).Errorf("[SendAuthCode] send mail error, err = %v", err)
	}
	key := RegisterEmailProKey + email
	if _, err = common.Redis.Set(key, code, 10*time.Minute).Result(); err != nil {
		logger.WithCtx(ctx).Errorf("[SendAuthCode] set redis error, err = %v", err)
		return e.ErrUnknown
	}
	return nil
}

func (u *authService) UserRegister(ctx context.Context, user *po.SysUser, code string) error {
	// email统一设置为小写
	user.Email = strings.ToLower(user.Email)
	// 检测是否已注册过
	exist, err := u.sysUserDao.CheckEmailIsExist(common.Mysql, user.Email)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UserRegister] check email fail, err = %v", err)
		return err
	}
	if exist {
		// 用户已经注册过
		return e.ErrUserEmailIsExist
	}
	// 检测code
	result := common.Redis.Get(RegisterEmailProKey + user.Email)
	if result.Err() != nil {
		logger.WithCtx(ctx).Errorf("[SendAuthCode] get redis error, err = %v", err)
		return e.ErrUnknown
	}
	if result.Val() != code {
		return e.NewCustomMsg("验证码错误")
	}
	// 设置用户名
	if user.Username == "" {
		user.Username = "fancoder"
		return nil
	}
	// 生成用户名称，唯一名称
	loginName := utils.GetRandomNumber(11)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UserRegister] new user name error, err = %v", err)
		return e.ErrUnknown
	}
	loginName = loginName + utils.GetRandomNumber(3)
	for i := 0; i < 5; i++ {
		var exist bool
		exist, err = u.sysUserDao.CheckLoginName(common.Mysql, user.LoginName)
		if err != nil {
			logger.WithCtx(ctx).Errorf("[UserRegister] check user name error, err = %v", err)
			return e.ErrUnknown
		}
		if exist {
			loginName = loginName + utils.GetRandomNumber(1)
		} else {
			break
		}
	}
	user.LoginName = loginName
	if len(user.Password) < 6 {
		return e.ErrUserPasswordNotEnoughAccuracy
	}
	//进行注册操作
	newPassword, err := utils.GetPwd(user.Password)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UserRegister] get password error, err = %v", err)
		return e.ErrPasswordEncodeFailed
	}
	user.Password = string(newPassword)

	// 设置默认值
	user.BirthDay = time.Now()
	user.Sex = 1
	err = common.Mysql.Transaction(func(tx *gorm.DB) error {
		if err2 := u.sysUserDao.InsertUser(tx, user); err2 != nil {
			return err2
		}
		return u.sysUserDao.InsertRolesToUser(tx, user.ID, []uint{constants.UserID})
	})
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UserRegister] inser user error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}
