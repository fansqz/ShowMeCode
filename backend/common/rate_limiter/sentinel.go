package rate_limiter

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
)

// InitSentinel 配置限流器
func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		panic(err)
	}

	// 设置限流规则，暂时所有请求配个30qps的限流
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "default",
			Threshold:              50,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})
	if err != nil {
		panic(err)
	}
}
