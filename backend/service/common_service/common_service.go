package common_service

import (
	"context"
	conf "github.com/fansqz/fancode-backend/common/config"
)

type CommonService interface {
	// GetURL 根据路径获取url
	GetURL(ctx context.Context, path string) (string, error)
}

func NewCommonService(config *conf.AppConfig) CommonService {
	return &commonService{
		config: config,
	}
}

type commonService struct {
	config *conf.AppConfig
}

func (c *commonService) GetURL(ctx context.Context, urlPath string) (string, error) {
	return c.config.URLPrefix + urlPath, nil
}
