package dao

import (
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type UserSavedCodeDao interface {
	// CreateUserSavedCode 创建用户保存的代码
	CreateUserSavedCode(db *gorm.DB, savedCode *po.UserSavedCode) error
	// UpdateUserSavedCode 更新用户保存的代码
	UpdateUserSavedCode(db *gorm.DB, savedCode *po.UserSavedCode) error
	// DeleteUserSavedCode 删除用户保存的代码
	DeleteUserSavedCode(db *gorm.DB, id uint, userID uint) error
	// GetUserSavedCodeByID 根据ID获取用户保存的代码
	GetUserSavedCodeByID(db *gorm.DB, id uint, userID uint) (*po.UserSavedCode, error)
	// GetUserSavedCodeList 获取用户保存的代码列表
	GetUserSavedCodeList(db *gorm.DB, userID uint, documentID *uint, language string, page, pageSize int) ([]*po.UserSavedCode, int64, error)
	// CheckUserSavedCodeExists 检查用户保存的代码是否存在
	CheckUserSavedCodeExists(db *gorm.DB, id uint, userID uint) (bool, error)
}

type userSavedCodeDao struct {
}

func NewUserSavedCodeDao() UserSavedCodeDao {
	return &userSavedCodeDao{}
}

// CreateUserSavedCode 创建用户保存的代码
func (u *userSavedCodeDao) CreateUserSavedCode(db *gorm.DB, savedCode *po.UserSavedCode) error {
	return db.Create(savedCode).Error
}

// UpdateUserSavedCode 更新用户保存的代码
func (u *userSavedCodeDao) UpdateUserSavedCode(db *gorm.DB, savedCode *po.UserSavedCode) error {
	return db.Model(&po.UserSavedCode{}).Where("id = ? AND user_id = ?", savedCode.ID, savedCode.UserID).Updates(map[string]interface{}{
		"language": savedCode.Language,
		"code":     savedCode.Code,
		"remark":   savedCode.Remark,
	}).Error
}

// DeleteUserSavedCode 删除用户保存的代码
func (u *userSavedCodeDao) DeleteUserSavedCode(db *gorm.DB, id uint, userID uint) error {
	return db.Where("id = ? AND user_id = ?", id, userID).Delete(&po.UserSavedCode{}).Error
}

// GetUserSavedCodeByID 根据ID获取用户保存的代码
func (u *userSavedCodeDao) GetUserSavedCodeByID(db *gorm.DB, id uint, userID uint) (*po.UserSavedCode, error) {
	var savedCode po.UserSavedCode
	err := db.Where("id = ? AND user_id = ?", id, userID).First(&savedCode).Error
	if err != nil {
		return nil, err
	}
	return &savedCode, nil
}

// GetUserSavedCodeList 获取用户保存的代码列表
func (u *userSavedCodeDao) GetUserSavedCodeList(db *gorm.DB, userID uint, documentID *uint, language string, page, pageSize int) ([]*po.UserSavedCode, int64, error) {
	var savedCodes []*po.UserSavedCode
	var total int64

	query := db.Model(&po.UserSavedCode{}).Where("user_id = ?", userID)

	// 如果指定了文档ID，则按文档ID过滤
	if documentID != nil {
		query = query.Where("document_id = ?", *documentID)
	} else {
		// 如果没有指定文档ID，则只查询没有文档ID的记录（非知识库页面保存的代码）
		query = query.Where("document_id IS NULL")
	}

	// 如果指定了编程语言，则按语言过滤
	if language != "" {
		query = query.Where("language = ?", language)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&savedCodes).Error
	if err != nil {
		return nil, 0, err
	}

	return savedCodes, total, nil
}

// CheckUserSavedCodeExists 检查用户保存的代码是否存在
func (u *userSavedCodeDao) CheckUserSavedCodeExists(db *gorm.DB, id uint, userID uint) (bool, error) {
	var count int64
	err := db.Model(&po.UserSavedCode{}).Where("id = ? AND user_id = ?", id, userID).Count(&count).Error
	return count > 0, err
}
