package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// GetPwd 密码加密
func GetPwd(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return hash, err
}

// ComparePwd
//
//	@Description: 检验密码是否是否正确
//	@param pwd1   原密码
//	@param pwd2   加密以后密码
//	@return bool  是否正确
func ComparePwd(pwd1 string, pwd2 string) bool {
	// Returns true on success, pwd1 is for the database.
	err := bcrypt.CompareHashAndPassword([]byte(pwd1), []byte(pwd2))
	if err != nil {
		return false
	} else {
		return true
	}
}
