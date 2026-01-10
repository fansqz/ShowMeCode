package po

import (
	"gorm.io/gorm"
)

type SysRole struct {
	gorm.Model
	Name        string     `gorm:"column:name" json:"name"`               // 角色名称
	Description string     `gorm:"column:description" json:"description"` // 描述
	Apis        []*SysApi  `gorm:"many2many:role_apis;" json:"apis"`      // 角色关联的api
	Menus       []*SysMenu `gorm:"many2many:role_menus" json:"menus"`     // 角色关联的菜单
}
