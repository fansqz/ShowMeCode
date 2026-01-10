package interceptor

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewCorsInterceptor,
	NewRequestInterceptor,
	NewRecoverPanicInterceptor,
	NewLoggerInterceptor,
)
