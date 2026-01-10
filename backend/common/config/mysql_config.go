package config

import "gopkg.in/ini.v1"

// MySqlConfig
// @Description: mysql相关配置
type MySqlConfig struct {
	User     string `ini:"user"`     //用户名
	Password string `ini:"password"` //密码
	DB       string `ini:"db"`       //要操作的数据库
	Host     string `ini:"host"`     //host
	Port     string `ini:"port"`     //端口
}

func NewMySqlConfig(cfg *ini.File) *MySqlConfig {
	mysqlConfig := &MySqlConfig{}
	cfg.Section("mysql").MapTo(mysqlConfig)
	return mysqlConfig
}
