package interceptor

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	e "github.com/fansqz/fancode-backend/common/error"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/gin-gonic/gin"
)

// RateLimitInterceptor 限流器
func RateLimitInterceptor() gin.HandlerFunc {
	return func(context *gin.Context) {
		result := r.NewResult(context)
		// 埋点（流控规则方式）
		entry, b := sentinel.Entry("default", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			result.Error(e.ErrTooManyRequests)
			return
		} else {
			entry.Exit()
		}
		context.Next()
	}
}
