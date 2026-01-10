package dao

import (
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type VisualDocumentBankDao interface {
	// InsertVisualDocumentBank 添加题库
	InsertVisualDocumentBank(db *gorm.DB, bank *po.VisualDocumentBank) error
	// GetVisualDocumentBankByID 根据题库id获取题库
	GetVisualDocumentBankByID(db *gorm.DB, bankID uint) (*po.VisualDocumentBank, error)
	// UpdateVisualDocumentBank 更新题库
	UpdateVisualDocumentBank(db *gorm.DB, bank *po.VisualDocumentBank) error
	// DeleteVisualDocumentBankByID 删除题库
	DeleteVisualDocumentBankByID(db *gorm.DB, id uint) error
	// GetAllVisualDocumentBank 获取所有的题目数据
	GetAllVisualDocumentBank(db *gorm.DB) ([]*po.VisualDocumentBank, error)
}

type visualDocumentBankDao struct {
}

func NewVisualDocumentBankDao() VisualDocumentBankDao {
	return &visualDocumentBankDao{}
}

func (v *visualDocumentBankDao) InsertVisualDocumentBank(db *gorm.DB, bank *po.VisualDocumentBank) error {
	return db.Create(bank).Error
}

func (v *visualDocumentBankDao) GetVisualDocumentBankByID(db *gorm.DB, bankID uint) (*po.VisualDocumentBank, error) {
	bank := &po.VisualDocumentBank{}
	err := db.First(&bank, bankID).Error
	return bank, err
}

func (v *visualDocumentBankDao) UpdateVisualDocumentBank(db *gorm.DB, bank *po.VisualDocumentBank) error {
	return db.Model(bank).Updates(bank).Error
}

func (v *visualDocumentBankDao) DeleteVisualDocumentBankByID(db *gorm.DB, id uint) error {
	return db.Delete(&po.VisualDocumentBank{}, id).Error
}

func (v *visualDocumentBankDao) GetAllVisualDocumentBank(db *gorm.DB) ([]*po.VisualDocumentBank, error) {
	var banks []*po.VisualDocumentBank
	err := db.Find(&banks).Error
	return banks, err
}
