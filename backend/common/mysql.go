// Package db
// @Author: fzw
// @Create: 2023/7/4
// @Description: 数据库开启关闭等
package common

import (
	"fmt"
	"github.com/fansqz/fancode-backend/common/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitMysql
//
//	@Description: 初始化mysql
//	@param cfg
//	@return error
func InitMysql(cfg *config.MySqlConfig) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	var err error
	Mysql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Errorf("[InitMysql] init mysql error, err = %v", err)
		panic(err)
	}
}
