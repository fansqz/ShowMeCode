package config

import "gopkg.in/ini.v1"

// COSConfig @Description:oss的配置
type COSConfig struct {
	AppID             string `ini:"appID"`
	Region            string `ini:"region"`
	SecretID          string `ini:"secretID"`
	SecretKey         string `ini:"secretKey"`
	ProblemBucketName string `ini:"problemBucketName"`
	ImageBucketName   string `ini:"imageBucketName"`
}

func NewCOSConfig(cfg *ini.File) *COSConfig {
	cosConfig := &COSConfig{}
	cfg.Section("cos").MapTo(cosConfig)
	return cosConfig
}
