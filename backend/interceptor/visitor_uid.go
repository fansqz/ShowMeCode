package interceptor

import (
	"github.com/fansqz/fancode-backend/constants"
	"github.com/fansqz/fancode-backend/utils"
	"github.com/gin-gonic/gin"
)

// VisitorUIDInterceptor 前端header中获取访客id
func VisitorUIDInterceptor() gin.HandlerFunc {
	return func(context *gin.Context) {
		visitorUID := context.Request.Header.Get(constants.VisitorUIDHeader)
		if visitorUID != "" {
			context.Set(utils.CtxVisitorUID, visitorUID)
		}
		context.Next()
	}
}
