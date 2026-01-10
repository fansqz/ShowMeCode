// Package initialize
// @Author: fzw
// @Create: 2023/7/14
// @Description: 初始化时读取配置文件相关工具
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/ini.v1"
)

// InitSetting
//
//	@Description: 初始化配置
//	@param test_file 配置文件路径
//	@return error
func InitSetting() *AppConfig {
	// 规定配置读取位置在执行文件所在目录下的/conf目录下
	// 读取当前位置
	path, _ := os.Getwd()
	path = strings.ReplaceAll(path, "\\", "/")

	// 通过环境变量控制读取哪个配置文件
	path = filepath.Join(path, "conf", "config_local.ini")
	if os.Getenv("env") != "" {
		path = path + fmt.Sprintf("/conf/config_%s.ini", os.Getenv("env"))
	}
	config, err := InitSettingWithPath(path)
	if err != nil {
		panic(err)
	}
	return config
}

func InitSettingWithPath(path string) (*AppConfig, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	config := new(AppConfig)
	cfg.MapTo(config)

	config.MySqlConfig = NewMySqlConfig(cfg)
	config.RedisConfig = NewRedisConfig(cfg)
	config.EmailConfig = NewEmailConfig(cfg)
	config.COSConfig = NewCOSConfig(cfg)
	config.FilePathConfig = NewFilePathConfig(cfg)
	config.LoggerConfig = NewLoggerConfig(cfg)
	config.AIConfig = NewAIConfig(cfg)
	return config, nil
}

// NewAIConfig 从ini文件加载AI配置
func NewAIConfig(cfg *ini.File) *AIConfig {
	aiSection := cfg.Section("ai")
	return &AIConfig{
		Provider: aiSection.Key("provider").MustString("openai"),
		ApiKey:   aiSection.Key("api_key").String(),
		ApiBase:  aiSection.Key("api_base").String(),
		Model:    aiSection.Key("model").MustString("gpt-3.5-turbo"),
		Timeout:  aiSection.Key("timeout").MustInt(30),
	}
}
