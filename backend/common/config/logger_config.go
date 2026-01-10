package config

import "gopkg.in/ini.v1"

type LoggerConfig struct {
	// Logger 日志收集类型
	Type string `ini:"type"`
	Host string `ini:"host"`
	Port string `ini:"port"`
}

func NewLoggerConfig(cfg *ini.File) *LoggerConfig {
	config := &LoggerConfig{}
	cfg.Section("logger").MapTo(config)
	return config
}
