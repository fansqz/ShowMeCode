package po

import (
	"gorm.io/gorm"
)

// UserSavedCode 用户保存的代码
type UserSavedCode struct {
	gorm.Model
	UserID     uint   `gorm:"column:user_id;not null" json:"userID"`            // 用户ID
	DocumentID *uint  `gorm:"column:document_id" json:"documentID"`             // 文档ID，可选
	Language   string `gorm:"column:language;not null;size:50" json:"language"` // 编程语言
	Code       string `gorm:"column:code;type:longtext" json:"code"`            // 代码内容
	Remark     string `gorm:"column:remark;type:text" json:"remark"`            // 备注
}

// TableName 指定表名
func (UserSavedCode) TableName() string {
	return "user_saved_codes"
}
