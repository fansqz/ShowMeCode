package gosync

import (
	"context"
	"fmt"
)

// Go 封装的go协程工具，会兜住panic，但是目前只能传递ctx
func Go(ctx context.Context, task func(ctx context.Context)) {
	go func(ctx context.Context, f func(ctx context.Context)) { // 匿名函数的参数为业务逻辑函数
		defer func() {
			// 在每个协程内部接收该协程自身抛出来的 panic
			if err := recover(); err != nil {
				fmt.Println("defer", err)
			}
		}()

		f(ctx) // 业务函数调用执行

	}(ctx, task) // 将当前的业务函数名传递给协程
}
