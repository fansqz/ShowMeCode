package interceptor

import (
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/models/vo"
	"github.com/gin-gonic/gin"
)

type RecoverPanicInterceptor struct {
}

func NewRecoverPanicInterceptor() *RecoverPanicInterceptor {
	return &RecoverPanicInterceptor{}
}

func (i *RecoverPanicInterceptor) RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			result := vo.NewResult(c)
			if err := recover(); err != nil {
				// 记录 panic 的错误信息
				logger.WithCtx(c).Errorf("Recovered from panic: %v", err)
				// 给客户端一个友好的错误提示
				result.SimpleErrorMessage("系统错误")
			}
		}()

		// 继续执行下一个中间件或者处理函数
		c.Next()
	}
}
