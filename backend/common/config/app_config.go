package config

// AppConfig
// @Description:应用配置
type AppConfig struct {
	Release         bool   `ini:"release"` //是否是上线模式
	Port            string `ini:"port"`    //端口
	URLPrefix       string `ini:"urlPrefix"`
	DefaultPassword string `ini:"defaultPassword"`
	DebuggerImage   string `ini:"debuggerImage"` // 调试器镜像名称
	*MySqlConfig
	*RedisConfig
	*EmailConfig
	*ReleasePathConfig
	*COSConfig
	*FilePathConfig
	*LoggerConfig
	*AIConfig
}

type ReleasePathConfig struct {
	StartWith []string
}
