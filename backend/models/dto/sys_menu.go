package dto

import (
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
)

type SysMenuTreeDto struct {
	ID           uint              `json:"id"`
	ParentMenuID uint              `json:"parentMenuID"`
	Code         string            `json:"code"`        // 请求路径
	Name         string            `json:"name"`        // 请求名称
	Description  string            `json:"description"` // 描述
	UpdatedAt    utils.Time        `json:"updatedAt"`   // 更新时间
	Children     []*SysMenuTreeDto `json:"children"`    //子menu
}

func NewSysMenuTreeDto(sysMenu *po.SysMenu) *SysMenuTreeDto {
	return &SysMenuTreeDto{
		ID:           sysMenu.ID,
		ParentMenuID: sysMenu.ParentMenuID,
		Code:         sysMenu.Code,
		Name:         sysMenu.Name,
		Description:  sysMenu.Description,
		UpdatedAt:    utils.Time(sysMenu.UpdatedAt),
	}
}
