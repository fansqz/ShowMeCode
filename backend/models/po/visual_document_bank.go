package po

import "gorm.io/gorm"

// VisualDocumentBank 知识库
type VisualDocumentBank struct {
	gorm.Model
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	CreatorID   uint   `gorm:"column:creator_id" json:"creatorID"`
	// Enable 是否启用
	Enable bool `gorm:"column:enable" json:"enable"`
}
