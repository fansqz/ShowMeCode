package dto

import (
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
)

// SysRoleDtoForList 获取角色列表
type SysRoleDtoForList struct {
	ID          uint       `json:"id"`
	UpdatedAt   utils.Time `json:"updatedAt"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}

func NewSysRoleDtoForList(role *po.SysRole) *SysRoleDtoForList {
	response := &SysRoleDtoForList{
		ID:          role.ID,
		UpdatedAt:   utils.Time(role.UpdatedAt),
		Name:        role.Name,
		Description: role.Description,
	}
	return response
}

// SimpleRoleDto 简单的角色，用于获取只有id和名称的角色列表
type SimpleRoleDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewSimpleRoleDto(role *po.SysRole) *SimpleRoleDto {
	response := &SimpleRoleDto{
		ID:   role.ID,
		Name: role.Name,
	}
	return response
}
