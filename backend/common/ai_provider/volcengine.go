package ai_provider

import (
	"context"
	"fmt"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

// VolcengineProvider 火山引擎AI服务商实现
// 使用官方Go SDK与火山引擎API交互
type VolcengineProvider struct {
	client  *arkruntime.Client
	model   string
	timeout int
}

// NewVolcengineProvider 创建火山引擎AI服务商实例
func NewVolcengineProvider(apiKey, apiBase, model string, timeout int) *VolcengineProvider {
	// 创建客户端配置
	client := arkruntime.NewClientWithApiKey(
		apiKey,
		arkruntime.WithBaseUrl(apiBase),
	)

	return &VolcengineProvider{
		client:  client,
		model:   model,
		timeout: timeout,
	}
}

// Chat 实现AIProvider接口
func (p *VolcengineProvider) Chat(ctx context.Context, prompt string) (string, error) {
	req := model.CreateChatCompletionRequest{
		Model: p.model,
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(prompt),
				},
			},
		},
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return "", err
	}
	return *resp.Choices[0].Message.Content.StringValue, nil
}
