package dto

import (
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
)

// UserSavedCodeDtoForCreate 创建用户保存代码的DTO
type UserSavedCodeDtoForCreate struct {
	DocumentID *uint  `json:"documentID" binding:"omitempty"` // 文档ID，可选
	Language   string `json:"language" binding:"required"`    // 编程语言
	Code       string `json:"code" binding:"required"`        // 代码内容
	Remark     string `json:"remark"`                         // 备注
}

// UserSavedCodeDtoForUpdate 更新用户保存代码的DTO
type UserSavedCodeDtoForUpdate struct {
	ID       uint   `json:"id" binding:"required"`       // 代码ID
	Language string `json:"language" binding:"required"` // 编程语言
	Code     string `json:"code" binding:"required"`     // 代码内容
	Remark   string `json:"remark"`                      // 备注
}

// UserSavedCodeDtoForList 用户保存代码列表的DTO
type UserSavedCodeDtoForList struct {
	ID         uint       `json:"id"`
	DocumentID *uint      `json:"documentID"` // 文档ID，可选
	Language   string     `json:"language"`   // 编程语言
	Remark     string     `json:"remark"`     // 备注
	CreatedAt  utils.Time `json:"createdAt"`  // 创建时间
	UpdatedAt  utils.Time `json:"updatedAt"`  // 更新时间
}

// UserSavedCodeDtoForDetail 用户保存代码详情的DTO
type UserSavedCodeDtoForDetail struct {
	ID         uint       `json:"id"`
	DocumentID *uint      `json:"documentID"` // 文档ID，可选
	Language   string     `json:"language"`   // 编程语言
	Code       string     `json:"code"`       // 代码内容
	Remark     string     `json:"remark"`     // 备注
	CreatedAt  utils.Time `json:"createdAt"`  // 创建时间
	UpdatedAt  utils.Time `json:"updatedAt"`  // 更新时间
}

// UserSavedCodeDtoForQuery 查询用户保存代码的DTO
type UserSavedCodeDtoForQuery struct {
	DocumentID *uint  `json:"documentID" form:"documentID"` // 文档ID，可选
	Language   string `json:"language" form:"language"`     // 编程语言，可选
	Page       int    `json:"page" form:"page"`             // 页码
	PageSize   int    `json:"pageSize" form:"pageSize"`     // 每页大小
}

// NewUserSavedCodeDtoForList 创建列表DTO
func NewUserSavedCodeDtoForList(savedCode *po.UserSavedCode) *UserSavedCodeDtoForList {
	return &UserSavedCodeDtoForList{
		ID:         savedCode.ID,
		DocumentID: savedCode.DocumentID,
		Language:   savedCode.Language,
		Remark:     savedCode.Remark,
		CreatedAt:  utils.Time(savedCode.CreatedAt),
		UpdatedAt:  utils.Time(savedCode.UpdatedAt),
	}
}

// NewUserSavedCodeDtoForDetail 创建详情DTO
func NewUserSavedCodeDtoForDetail(savedCode *po.UserSavedCode) *UserSavedCodeDtoForDetail {
	return &UserSavedCodeDtoForDetail{
		ID:         savedCode.ID,
		DocumentID: savedCode.DocumentID,
		Language:   savedCode.Language,
		Code:       savedCode.Code,
		Remark:     savedCode.Remark,
		CreatedAt:  utils.Time(savedCode.CreatedAt),
		UpdatedAt:  utils.Time(savedCode.UpdatedAt),
	}
}
