package dao

import (
	"errors"
	"github.com/fansqz/fancode-backend/constants"
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type UserCodeDao interface {
	CheckUserCode(db *gorm.DB, userId uint, problemId uint, languageType constants.LanguageType) (bool, error)
	GetUserCode(db *gorm.DB, userId uint, problemId uint, languageType constants.LanguageType) (*po.UserCode, error)
	InsertUserCode(db *gorm.DB, code *po.UserCode) error
	UpdateUserCode(db *gorm.DB, code *po.UserCode) error
	GetUserCodeListByProblemID(db *gorm.DB, userId uint, problemId uint) ([]*po.UserCode, error)
}

type userCodeDao struct {
}

func NewUserCodeDao() UserCodeDao {
	return &userCodeDao{}
}

func (u *userCodeDao) CheckUserCode(db *gorm.DB, userId uint, problemId uint,
	language constants.LanguageType) (bool, error) {
	userCode := po.UserCode{}
	err := db.Model(&po.UserCode{}).Where("user_id = ? and problem_id = ? and language = ?",
		userId, problemId, language).
		First(&userCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (u *userCodeDao) GetUserCode(db *gorm.DB, userId uint, problemId uint,
	language constants.LanguageType) (*po.UserCode, error) {
	userCode := po.UserCode{}
	err := db.Model(&po.UserCode{}).Where("user_id = ? and problem_id = ? and language = ?",
		userId, problemId, language).
		First(&userCode).Error
	return &userCode, err
}

func (u *userCodeDao) InsertUserCode(db *gorm.DB, code *po.UserCode) error {
	return db.Create(code).Error
}

func (u *userCodeDao) UpdateUserCode(db *gorm.DB, code *po.UserCode) error {
	return db.Model(&po.UserCode{}).Where("id = ?", code.ID).Updates(map[string]interface{}{
		"code": code.Code,
	}).Error
}

func (u *userCodeDao) GetUserCodeListByProblemID(db *gorm.DB, userId uint, problemId uint) ([]*po.UserCode, error) {
	var answer []*po.UserCode
	err := db.Where("user_id = ? and problem_id = ?", userId, problemId).Order("updated_at desc").Find(&answer).Error
	return answer, err
}
