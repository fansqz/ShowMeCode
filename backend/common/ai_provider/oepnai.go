package ai_provider

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

// OpenAIProvider OpenAI服务商实现
// 只负责与OpenAI API交互
type OpenAIProvider struct {
	ApiKey  string
	ApiBase string
	Model   string
	Timeout int
}

func NewOpenAIProvider(apiKey, apiBase, model string, timeout int) *OpenAIProvider {
	return &OpenAIProvider{
		ApiKey:  apiKey,
		ApiBase: apiBase,
		Model:   model,
		Timeout: timeout,
	}
}

// Chat 实现AIProvider接口
func (p *OpenAIProvider) Chat(ctx context.Context, prompt string) (string, error) {
	body := fmt.Sprintf(`{"model": "%s", "messages": [{"role": "user", "content": "%s"}]}`,
		p.Model, prompt)
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "POST", p.ApiBase+"/chat/completions", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+p.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}
