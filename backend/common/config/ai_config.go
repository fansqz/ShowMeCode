package config

// AIConfig
// @Description: AI服务相关配置
type AIConfig struct {
	Provider string `ini:"provider"` // ai服务商，如openai、tongyi等
	ApiKey   string `ini:"api_key"`
	ApiBase  string `ini:"api_base"`
	Model    string `ini:"model"`
	Timeout  int    `ini:"timeout"`
	// Access Key ID
	AccessKeyID string `ini:"access_key_id"`
	// Access Key Secret
	AccessKeySecret string `ini:"access_key_secret"`
}
