package po

import "gorm.io/gorm"

type SysMenu struct {
	gorm.Model
	ParentMenuID uint   `gorm:"column:parent_Menu_id" json:"parentMenuID"` // 父menu的id
	Name         string `gorm:"column:name" json:"name"`                   // 角色名称
	Code         string `gorm:"column:code" json:"code"`                   // 权限值
	Description  string `gorm:"column:description" json:"description"`     // 描述
}
