package user_saved_code_service

import (
	"context"

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

type UserSavedCodeService interface {
	// CreateUserSavedCode 创建用户保存的代码
	CreateUserSavedCode(ctx context.Context, req *dto.UserSavedCodeDtoForCreate) (*dto.UserSavedCodeDtoForDetail, error)
	// UpdateUserSavedCode 更新用户保存的代码
	UpdateUserSavedCode(ctx context.Context, req *dto.UserSavedCodeDtoForUpdate) (*dto.UserSavedCodeDtoForDetail, error)
	// DeleteUserSavedCode 删除用户保存的代码
	DeleteUserSavedCode(ctx context.Context, id uint) error
	// GetUserSavedCodeByID 根据ID获取用户保存的代码
	GetUserSavedCodeByID(ctx context.Context, id uint) (*dto.UserSavedCodeDtoForDetail, error)
	// GetUserSavedCodeList 获取用户保存的代码列表
	GetUserSavedCodeList(ctx context.Context, req *dto.UserSavedCodeDtoForQuery) (*dto.PageInfo, error)
}

type userSavedCodeService struct {
	savedCodeDao dao.UserSavedCodeDao
}

func NewUserSavedCodeService(config *conf.AppConfig, savedCodeDao dao.UserSavedCodeDao) UserSavedCodeService {
	return &userSavedCodeService{
		savedCodeDao: savedCodeDao,
	}
}

// CreateUserSavedCode 创建用户保存的代码
func (u *userSavedCodeService) CreateUserSavedCode(ctx context.Context, req *dto.UserSavedCodeDtoForCreate) (*dto.UserSavedCodeDtoForDetail, error) {
	userID := utils.GetUserIDWithCtx(ctx)

	// 创建新的保存代码记录
	savedCode := &po.UserSavedCode{
		UserID:     userID,
		DocumentID: req.DocumentID,
		Language:   req.Language,
		Code:       req.Code,
		Remark:     req.Remark,
	}

	if err := u.savedCodeDao.CreateUserSavedCode(common.Mysql, savedCode); err != nil {
		logger.WithCtx(ctx).Errorf("[CreateUserSavedCode] CreateUserSavedCode fail, err = %v", err)
		return nil, e.ErrUnknown
	}

	return dto.NewUserSavedCodeDtoForDetail(savedCode), nil
}

// UpdateUserSavedCode 更新用户保存的代码
func (u *userSavedCodeService) UpdateUserSavedCode(ctx context.Context, req *dto.UserSavedCodeDtoForUpdate) (*dto.UserSavedCodeDtoForDetail, error) {
	userID := utils.GetUserIDWithCtx(ctx)

	// 检查代码是否存在
	exists, err := u.savedCodeDao.CheckUserSavedCodeExists(common.Mysql, req.ID, userID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateUserSavedCode] CheckUserSavedCodeExists fail, err = %v", err)
		return nil, e.ErrUnknown
	}
	if !exists {
		return nil, e.NewRecordNotFoundErr("saved code not found")
	}

	// 更新保存代码记录
	savedCode := &po.UserSavedCode{
		Model:    gorm.Model{ID: req.ID},
		UserID:   userID,
		Language: req.Language,
		Code:     req.Code,
		Remark:   req.Remark,
	}

	if err = u.savedCodeDao.UpdateUserSavedCode(common.Mysql, savedCode); err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateUserSavedCode] UpdateUserSavedCode fail, err = %v", err)
		return nil, e.ErrUnknown
	}

	// 获取更新后的数据
	updatedCode, err := u.savedCodeDao.GetUserSavedCodeByID(common.Mysql, req.ID, userID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateUserSavedCode] GetUserSavedCodeByID fail, err = %v", err)
		return nil, e.ErrUnknown
	}

	return dto.NewUserSavedCodeDtoForDetail(updatedCode), nil
}

// DeleteUserSavedCode 删除用户保存的代码
func (u *userSavedCodeService) DeleteUserSavedCode(ctx context.Context, id uint) error {
	userID := utils.GetUserIDWithCtx(ctx)

	// 检查代码是否存在
	exists, err := u.savedCodeDao.CheckUserSavedCodeExists(common.Mysql, id, userID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[DeleteUserSavedCode] CheckUserSavedCodeExists fail, err = %v", err)
		return e.ErrUnknown
	}
	if !exists {
		return e.NewRecordNotFoundErr("saved code not found")
	}

	if err = u.savedCodeDao.DeleteUserSavedCode(common.Mysql, id, userID); err != nil {
		logger.WithCtx(ctx).Errorf("[DeleteUserSavedCode] DeleteUserSavedCode fail, err = %v", err)
		return e.ErrUnknown
	}

	return nil
}

// GetUserSavedCodeByID 根据ID获取用户保存的代码
func (u *userSavedCodeService) GetUserSavedCodeByID(ctx context.Context, id uint) (*dto.UserSavedCodeDtoForDetail, error) {
	userID := utils.GetUserIDWithCtx(ctx)

	savedCode, err := u.savedCodeDao.GetUserSavedCodeByID(common.Mysql, id, userID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetUserSavedCodeByID] GetUserSavedCodeByID fail, err = %v", err)
		return nil, e.NewRecordNotFoundErr("saved code not found")
	}

	return dto.NewUserSavedCodeDtoForDetail(savedCode), nil
}

// GetUserSavedCodeList 获取用户保存的代码列表
func (u *userSavedCodeService) GetUserSavedCodeList(ctx context.Context, req *dto.UserSavedCodeDtoForQuery) (*dto.PageInfo, error) {
	userID := utils.GetUserIDWithCtx(ctx)

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	savedCodes, total, err := u.savedCodeDao.GetUserSavedCodeList(common.Mysql, userID, req.DocumentID, req.Language, req.Page, req.PageSize)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetUserSavedCodeList] GetUserSavedCodeList fail, err = %v", err)
		return nil, e.ErrUnknown
	}

	// 转换为DTO
	result := make([]*dto.UserSavedCodeDtoForList, 0, len(savedCodes))
	for _, savedCode := range savedCodes {
		result = append(result, dto.NewUserSavedCodeDtoForList(savedCode))
	}

	return &dto.PageInfo{
		Total: total,
		Size:  int64(len(result)),
		List:  result,
	}, nil
}
