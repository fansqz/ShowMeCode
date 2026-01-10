package visual_document_service

import (
	"context"
	"errors"
	"github.com/fansqz/fancode-backend/common"
	conf "github.com/fansqz/fancode-backend/common/config"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/dao"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
	"gorm.io/gorm"
)

// VisualDocumentBankService 可视化题库管理
type VisualDocumentBankService interface {
	// InsertVisualDocumentBank 添加可视化题库
	InsertVisualDocumentBank(ctx context.Context, visualDocumentBank *po.VisualDocumentBank) (uint, error)
	// UpdateVisualDocumentBank 更新可视化题库
	UpdateVisualDocumentBank(ctx context.Context, visualDocumentBank *po.VisualDocumentBank) error
	// DeleteVisualDocumentBank 删除题库
	DeleteVisualDocumentBank(ctx context.Context, id uint) error
	// GetAllVisualDocumentBank 获取所有可视化题库列表
	GetAllVisualDocumentBank(ctx context.Context, enable *bool) ([]*dto.VisualDocumentBankDto, error)
	// GetVisualDocumentBankByID 获取可视化题库信息
	GetVisualDocumentBankByID(ctx context.Context, id uint) (*dto.VisualDocumentBankDto, error)
}

type visualDocumentBankService struct {
	config                *conf.AppConfig
	visualDocumentBankDao dao.VisualDocumentBankDao
	visualDocumentDao     dao.VisualDocumentDao
	sysUserDao            dao.SysUserDao
}

func NewVisualDocumentBankService(config *conf.AppConfig, bankDao dao.VisualDocumentBankDao, documentDao dao.VisualDocumentDao,
	sysUserDao dao.SysUserDao) VisualDocumentBankService {
	return &visualDocumentBankService{
		config:                config,
		visualDocumentBankDao: bankDao,
		visualDocumentDao:     documentDao,
		sysUserDao:            sysUserDao,
	}
}

func (v *visualDocumentBankService) InsertVisualDocumentBank(ctx context.Context, visualDocumentBank *po.VisualDocumentBank) (uint, error) {
	userID := utils.GetUserIDWithCtx(ctx)
	visualDocumentBank.CreatorID = userID
	if err := v.visualDocumentBankDao.InsertVisualDocumentBank(common.Mysql, visualDocumentBank); err != nil {
		logger.WithCtx(ctx).Errorf("[InsertVisualDocumentBank] insert visual document bank error, err = %v", err)
		return 0, e.ErrMysql
	}
	return visualDocumentBank.ID, nil
}

func (v *visualDocumentBankService) UpdateVisualDocumentBank(ctx context.Context, visualDocumentBank *po.VisualDocumentBank) error {
	visualDocumentBank.CreatorID = 0
	err := v.visualDocumentBankDao.UpdateVisualDocumentBank(common.Mysql, visualDocumentBank)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateVisualDocumentBank] update visual document bank error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}

func (v *visualDocumentBankService) DeleteVisualDocumentBank(ctx context.Context, id uint) error {
	var err error
	if err = v.visualDocumentBankDao.DeleteVisualDocumentBankByID(common.Mysql, id); err != nil {
		logger.WithCtx(ctx).Errorf("[DeleteVisualDocumentBank] delete visual document bank error, err = %v", err)
		return e.ErrMysql
	}
	if err = v.visualDocumentDao.DeleteVisualDocumentByBankID(common.Mysql, id); err != nil {
		logger.WithCtx(ctx).Errorf("[DeleteVisualDocumentBank] DeleteVisualDocumentByBankID error, err = %v", err)
		return e.ErrMysql
	}

	return nil
}

func (v *visualDocumentBankService) GetAllVisualDocumentBank(ctx context.Context, enable *bool) ([]*dto.VisualDocumentBankDto, error) {
	banks, err := v.visualDocumentBankDao.GetAllVisualDocumentBank(common.Mysql)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetAllvisualDocumentBank] get all problem bank error, err = %v", err)
		return nil, e.ErrMysql
	}
	// 根据enable过滤
	if enable != nil {
		newDocumentList := make([]*po.VisualDocumentBank, 0, len(banks))
		for _, d := range banks {
			if d.Enable == *enable {
				newDocumentList = append(newDocumentList, d)
			}
		}
		banks = newDocumentList
	}

	// 设置更多信息
	answer := []*dto.VisualDocumentBankDto{}
	for i, _ := range banks {
		answer = append(answer, dto.NewVisualDocumentBankDto(banks[i]))
		answer[i].CreatorName, err = v.sysUserDao.GetUserNameByID(common.Mysql, answer[i].CreatorID)
	}
	return answer, nil
}

func (v *visualDocumentBankService) GetVisualDocumentBankByID(ctx context.Context, id uint) (*dto.VisualDocumentBankDto, error) {
	bank, err := v.visualDocumentBankDao.GetVisualDocumentBankByID(common.Mysql, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.ErrProblemNotExist
	}
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetvisualDocumentBankByID] get problem bank error, err = %v", err)
		return nil, e.ErrMysql
	}
	answer := dto.NewVisualDocumentBankDto(bank)
	answer.CreatorName, err = v.sysUserDao.GetUserNameByID(common.Mysql, answer.CreatorID)
	return answer, nil
}
