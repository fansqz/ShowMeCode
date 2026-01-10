package po

import (
	"gorm.io/gorm"
)

// UserCode
// 通过userCode可以找到用户在某个题目某个语言中保存的代码
type UserCode struct {
	gorm.Model
	ProblemID uint   `gorm:"column:problem_id" json:"problemID"`
	UserID    uint   `gorm:"column:user_id" json:"userID"`
	Language  string `gorm:"column:language" json:"language"`
	Code      string `gorm:"column:code" json:"code"`
}
