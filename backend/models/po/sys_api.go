package po

import (
	"gorm.io/gorm"
)

type SysApi struct {
	gorm.Model
	ParentApiID uint   `gorm:"column:parent_api_id" json:"parentApiID"` // 父api的id
	Path        string `gorm:"column:path" json:"path"`                 // 请求路径
	Method      string `gorm:"column:method" json:"method"`             // 请求方法
	Name        string `gorm:"column:name" json:"name"`                 // 请求名称
	Description string `gorm:"column:description" json:"description"`   // 描述
}
