package dto

import (
	"github.com/fansqz/fancode-backend/models/po"
	"time"
)

type ActivityItem struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// AccountInfo
// 和userInfo类似，但是比userInfo的数据多一些
type AccountInfo struct {
	Avatar       string `json:"avatar"`
	LoginName    string `json:"loginName"`
	UserName     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Introduction string `json:"introduction"`
	Sex          int    `json:"sex"`
	BirthDay     string `json:"birthDay"`
	CodingAge    int    `json:"codingAge"`
	// 用户权限和菜单权限
	Roles []uint   `json:"roles"`
	Menus []string `json:"menus"`
}

func NewAccountInfo(user *po.SysUser) *AccountInfo {
	userInfo := &AccountInfo{
		Avatar:       user.Avatar,
		LoginName:    user.LoginName,
		UserName:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		Introduction: user.Introduction,
		BirthDay:     user.BirthDay.Format("2006-01-02"),
		Sex:          user.Sex,
		CodingAge:    time.Now().Year() - user.CreatedAt.Year(),
	}
	userInfo.Roles = make([]uint, len(user.Roles))
	for i := 0; i < len(user.Roles); i++ {
		userInfo.Roles[i] = user.Roles[i].ID
		for j := 0; j < len(user.Roles[i].Menus); j++ {
			userInfo.Menus = append(userInfo.Menus, user.Roles[i].Menus[j].Code)
		}
	}
	return userInfo
}
