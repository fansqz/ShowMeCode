package interceptor

import (
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/gin-gonic/gin"
)

type LoggerInterceptor struct {
}

func NewLoggerInterceptor() *LoggerInterceptor {
	return &LoggerInterceptor{}
}

// LoggerInterceptor 日志拦截，设置有些必要的日志信息
func (*LoggerInterceptor) LoggerInterceptor() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 随机生成log_id
		logID := logger.GenerateLogID()
		context.Set(logger.LOG_ID_KEY, logID)
	}
}
