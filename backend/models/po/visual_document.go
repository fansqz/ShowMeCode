package po

import "gorm.io/gorm"

// VisualDocument 一篇可视化的文章
type VisualDocument struct {
	gorm.Model
	BankID uint `gorm:"column:bank_id" json:"bankID"`
	// ParentID 父节点id
	ParentID uint `gorm:"column:parent_id" json:"parentID"`
	// Order 文档的排列顺序
	Order     uint   `gorm:"order" json:"order"`
	Title     string `gorm:"column:title" json:"title"`
	Content   string `gorm:"column:content" json:"content"`
	CreatorID uint   `gorm:"column:creator_id" json:"creatorID"`
	// Enable 是否启用
	Enable bool `gorm:"column:enable" json:"enable"`
}

// VisualDocumentCode 存储可视化文章的代码
type VisualDocumentCode struct {
	gorm.Model
	DocumentID uint   `gorm:"column:document_id" json:"documentID"`
	Code       string `gorm:"column:code" json:"code"`
	Language   string `gorm:"column:language" json:"language"`
	// Breakpoints 初始断点
	Breakpoints string `gorm:"column:breakpoints" json:"breakpoints"`
}
