package ai_provider

import (
	"context"
	"github.com/fansqz/fancode-backend/common/config"
)

// AIProvider 只定义最基础的AI调用方法
// prompt为业务层拼接好的字符串，返回AI原始响应内容
// 具体实现只负责与AI服务商交互
// 如需支持多家AI，只需实现本接口
type AIProvider interface {
	Chat(ctx context.Context, prompt string) (string, error)
}

// 根据ai配置生成不同的AIProvider
func NewAIProvider(aiConfig *config.AIConfig) AIProvider {
	switch aiConfig.Provider {
	case "openai":
		return NewOpenAIProvider(aiConfig.ApiKey, aiConfig.ApiBase, aiConfig.Model, aiConfig.Timeout)
	case "volcengine":
		return NewVolcengineProvider(aiConfig.ApiKey, aiConfig.ApiBase, aiConfig.Model, aiConfig.Timeout)
	default:
		// 默认使用OpenAI
		return NewOpenAIProvider(aiConfig.ApiKey, aiConfig.ApiBase, aiConfig.Model, aiConfig.Timeout)
	}
}
