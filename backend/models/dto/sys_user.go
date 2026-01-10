package dto

import (
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
)

// SysUserDtoForList 获取用户列表
type SysUserDtoForList struct {
	ID        uint       `json:"id"`
	LoginName string     `json:"loginName"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	UpdateAt  utils.Time `json:"updatedAt"`
	Roles     []string   `json:"roles"`
}

func NewSysUserDtoForList(user *po.SysUser) *SysUserDtoForList {
	response := &SysUserDtoForList{
		ID:        user.ID,
		LoginName: user.LoginName,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		UpdateAt:  utils.Time(user.UpdatedAt),
	}
	if user.Roles != nil {
		response.Roles = make([]string, len(user.Roles))
		for i, role := range user.Roles {
			response.Roles[i] = role.Name
		}
	}
	return response
}
