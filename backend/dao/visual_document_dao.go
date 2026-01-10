package dao

import (
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type VisualDocumentDao interface {
	GetAllSimpleDocument(db *gorm.DB, bankID uint) ([]*po.VisualDocument, error)
	GetSimpleDocumentByParentID(db *gorm.DB, parentID uint) ([]*po.VisualDocument, error)
	GetVisualDocumentByID(db *gorm.DB, id uint) (*po.VisualDocument, error)
	GetCodeListByDocumentID(db *gorm.DB, documentID uint) ([]*po.VisualDocumentCode, error)

	InsertVisualDocument(db *gorm.DB, code *po.VisualDocument) error
	InsertVisualDocumentCode(db *gorm.DB, code *po.VisualDocumentCode) error

	UpdateVisualDocumentTitle(db *gorm.DB, id uint, title string) error
	UpdateVisualDocumentEnable(db *gorm.DB, id uint, enable bool) error
	UpdateVisualDocument(db *gorm.DB, code *po.VisualDocument) error
	UpdateVisualDocumentParentAndOrder(db *gorm.DB, id uint, parentID uint, order uint) error
	UpdateVisualDocumentCode(db *gorm.DB, code *po.VisualDocumentCode) error

	DeleteVisualDocumentByID(db *gorm.DB, documentID uint) error
	DeleteVisualDocumentCodeByDocumentID(db *gorm.DB, documentID uint) error
	DeleteVisualDocumentCodeByID(db *gorm.DB, id uint) error
	DeleteVisualDocumentByBankID(db *gorm.DB, bankID uint) error
}

type VisualDocumentCodeDao interface {
	GetLanguageListByDocumentID(db *gorm.DB, documentID uint) ([]string, error)
}

func NewVisualDocumentDao() VisualDocumentDao {
	return &visualDocumentDao{}
}

type visualDocumentDao struct {
}

func (v *visualDocumentDao) GetAllSimpleDocument(db *gorm.DB, bankID uint) ([]*po.VisualDocument, error) {
	documents := []*po.VisualDocument{}
	err := db.Select("id", "parent_id", "title", "creator_id", "order", "enable").Where("bank_id = ?", bankID).Find(&documents).Error
	return documents, err
}

func (v *visualDocumentDao) GetSimpleDocumentByParentID(db *gorm.DB, parentID uint) ([]*po.VisualDocument, error) {
	documents := []*po.VisualDocument{}
	err := db.Select("id", "parent_id", "title", "creator_id", "order", "enable").Where("parent_id = ?", parentID).Find(&documents).Error
	return documents, err
}

func (v *visualDocumentDao) GetVisualDocumentByID(db *gorm.DB, id uint) (*po.VisualDocument, error) {
	document := &po.VisualDocument{}
	err := db.First(&document, id).Error
	return document, err
}

func (v *visualDocumentDao) GetCodeListByDocumentID(db *gorm.DB, documentID uint) ([]*po.VisualDocumentCode, error) {
	var codeList []*po.VisualDocumentCode
	err := db.Where("document_id = ?", documentID).Find(&codeList).Error
	return codeList, err
}

func (v *visualDocumentDao) InsertVisualDocument(db *gorm.DB, document *po.VisualDocument) error {
	return db.Create(document).Error
}

func (v *visualDocumentDao) InsertVisualDocumentCode(db *gorm.DB, code *po.VisualDocumentCode) error {
	return db.Create(code).Error
}

func (v *visualDocumentDao) UpdateVisualDocumentTitle(db *gorm.DB, id uint, title string) error {
	return db.Model(&po.VisualDocument{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title": title,
	}).Error
}

func (v *visualDocumentDao) UpdateVisualDocumentEnable(db *gorm.DB, id uint, enable bool) error {
	return db.Model(&po.VisualDocument{}).Where("id = ?", id).Updates(map[string]interface{}{
		"enable": enable,
	}).Error
}

func (v *visualDocumentDao) UpdateVisualDocument(db *gorm.DB, document *po.VisualDocument) error {
	return db.Model(&po.VisualDocument{}).Where("id = ?", document.ID).Updates(map[string]interface{}{
		"parent_id": document.ParentID,
		"title":     document.Title,
		"content":   document.Content,
		"enable":    document.Enable,
	}).Error
}

func (v *visualDocumentDao) UpdateVisualDocumentParentAndOrder(db *gorm.DB, id uint, parentID uint, order uint) error {
	return db.Model(&po.VisualDocument{}).Where("id = ?", id).Updates(map[string]interface{}{
		"order":     order,
		"parent_id": parentID,
	}).Error
}

func (v *visualDocumentDao) UpdateVisualDocumentCode(db *gorm.DB, code *po.VisualDocumentCode) error {
	return db.Model(&po.VisualDocumentCode{}).Where("id = ?", code.ID).Updates(map[string]interface{}{
		"code":       code.Code,
		"breakpoint": code.Breakpoints,
	}).Error
}

func (v *visualDocumentDao) DeleteVisualDocumentByID(db *gorm.DB, documentID uint) error {
	return db.Delete(&po.VisualDocument{}, documentID).Error
}

func (v *visualDocumentDao) DeleteVisualDocumentByBankID(db *gorm.DB, bankID uint) error {
	return db.Where("bank_id = ?", bankID).Delete(&po.VisualDocument{}).Error
}

func (v *visualDocumentDao) DeleteVisualDocumentCodeByDocumentID(db *gorm.DB, documentID uint) error {
	return db.Where("document_id = ?", documentID).Delete(&po.VisualDocumentCode{}).Error
}

func (v *visualDocumentDao) DeleteVisualDocumentCodeByID(db *gorm.DB, id uint) error {
	return db.Delete(&po.VisualDocumentCode{}, id).Error
}
