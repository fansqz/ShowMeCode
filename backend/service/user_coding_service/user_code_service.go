package user_coding_service

import (
	"context"
	"github.com/fansqz/fancode-backend/common"
	conf "github.com/fansqz/fancode-backend/common/config"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/constants"
	"github.com/fansqz/fancode-backend/dao"
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/**
 * 放一些公用的方法
 */

const (
	AcmCCodeFilePath    = "/conf/acm_template/c"
	AcmGoCodeFilePath   = "/conf/acm_template/go"
	AcmJavaCodeFilePath = "/conf/acm_template/java"
	VisualDocumentPath  = "/conf/document/visual_document.md"
)

type UserCodeService interface {
	// SaveUserCode 保存用户代码
	SaveUserCode(ctx context.Context, userCode *po.UserCode) error
	// GetUserCode 读取用户代码
	GetUserCode(ctx context.Context, problemID uint, language constants.LanguageType) (string, error)
	// GetUserCodeByProblemID 根据题目id获取用户代码，无语言类型
	GetUserCodeByProblemID(ctx context.Context, problemID uint) (*po.UserCode, error)
	// GetProblemTemplateCode 获取题目的模板代码
	GetProblemTemplateCode(ctx context.Context, language string) (string, error)
}

type userCodeService struct {
	codeDao dao.UserCodeDao
}

func NewUserCodeService(config *conf.AppConfig, userCodeDao dao.UserCodeDao) UserCodeService {
	return &userCodeService{
		codeDao: userCodeDao,
	}
}

// SaveUserCode 保存用户代码
func (u *userCodeService) SaveUserCode(ctx context.Context, userCode *po.UserCode) error {
	userID := utils.GetUserIDWithCtx(ctx)
	userCode.UserID = userID
	exist, err := u.codeDao.CheckUserCode(common.Mysql, userID, userCode.ProblemID, constants.LanguageType(userCode.Language))
	if err != nil {
		logger.WithCtx(ctx).Errorf("[SaveUserCode] CheckUserCode fail, err = %v", err)
		return e.ErrUnknown
	}
	// 添加
	if !exist {
		if err = u.codeDao.InsertUserCode(common.Mysql, userCode); err != nil {
			logger.WithCtx(ctx).Errorf("[SaveUserCode] InsertUserCode fail, err = %v", err)
			return err
		}
		return nil
	}
	// 保存
	var code *po.UserCode
	code, err = u.codeDao.GetUserCode(common.Mysql, userID, userCode.ProblemID, constants.LanguageType(userCode.Language))
	if err != nil {
		logger.WithCtx(ctx).Errorf("[SaveUserCode] GetUserCode fail, err = %v", err)
		return e.ErrUnknown
	}
	code.Code = userCode.Code
	if err = u.codeDao.UpdateUserCode(common.Mysql, code); err != nil {
		logger.WithCtx(ctx).Errorf("[SaveUserCode] GetUserCode fail, err = %v", err)
		return e.ErrUnknown
	}
	return nil
}

// GetUserCode 读取用户代码
func (u *userCodeService) GetUserCode(ctx context.Context, problemId uint, language constants.LanguageType) (string, error) {
	userID := utils.GetUserIDWithCtx(ctx)
	exist, err := u.codeDao.CheckUserCode(common.Mysql, userID, problemId, language)
	if err != nil {
		log.Println(err)
		return "", e.ErrUnknown
	}
	// 如果用户代码不存在，那么读取模板
	if !exist {
		// 读取acm模板
		code, err := getAcmCodeTemplate(language)
		if err != nil {
			return "", e.ErrProblemGetFailed
		}
		return code, nil
	}
	var code *po.UserCode
	code, err = u.codeDao.GetUserCode(common.Mysql, userID, problemId, language)
	if err != nil {
		log.Println(err)
		return "", e.ErrUnknown
	}
	return code.Code, nil
}

// GetUserCodeByProblemID 根据题目id获取用户代码，无语言类型
func (u *userCodeService) GetUserCodeByProblemID(ctx context.Context, problemId uint) (*po.UserCode, error) {
	userID := utils.GetUserIDWithCtx(ctx)
	codeList, err := u.codeDao.GetUserCodeListByProblemID(common.Mysql, userID, problemId)
	if err != nil {
		log.Println(err)
		return nil, e.ErrUnknown
	}

	if len(codeList) != 0 {
		return codeList[0], nil
	}
	// 如果用户代码不存在，那么读取模板
	language := constants.LanguageGo
	// 读取acm模板
	code, err := getAcmCodeTemplate(language)
	if err != nil {
		logger.WithCtx(ctx).Printf("get template fail, err = %v", err)
		return nil, e.ErrProblemGetFailed
	}
	return &po.UserCode{
		Code:     code,
		Language: string(language),
	}, nil
}

func (u *userCodeService) GetProblemTemplateCode(ctx context.Context, language string) (string, error) {
	// 读取acm模板
	code, err := getAcmCodeTemplate(constants.LanguageType(language))
	if err != nil {
		return "", e.ErrProblemGetFailed
	}
	return code, nil
}

func getAcmCodeTemplate(language constants.LanguageType) (string, error) {
	// 规定配置读取位置在执行文件所在目录下的/conf目录下
	// 读取当前位置
	dir, _ := os.Getwd()
	dir = strings.ReplaceAll(dir, "\\", "/")
	// 通过环境变量控制读取哪个配置文件
	var filePath string
	switch language {
	case constants.LanguageC:
		filePath = filepath.Join(dir, AcmCCodeFilePath)
	case constants.LanguageGo:
		filePath = filepath.Join(dir, AcmGoCodeFilePath)
	case constants.LanguageJava:
		filePath = filepath.Join(dir, AcmJavaCodeFilePath)
	}
	code, err := os.ReadFile(filePath)
	return string(code), err
}
