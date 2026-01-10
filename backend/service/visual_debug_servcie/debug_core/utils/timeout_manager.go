package utils

import (
	"context"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/utils/gosync"
	"time"
)

// TimeoutManager 一个计时器
// 如果在timeout时间内没有执行reset命令，就会执行fun函数
// duration至少10s
type TimeoutManager struct {
	timer          *time.Timer
	timeout        time.Duration
	resetChannel   chan struct{}
	chancelChannel chan struct{}
	fun            func()
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewTimeoutManager 创建一个新的计时器实例
func NewTimeoutManager() *TimeoutManager {
	return &TimeoutManager{
		resetChannel:   make(chan struct{}, 1), // 使用无缓冲通道并设置缓冲大小为1
		chancelChannel: make(chan struct{}, 1),
	}
}

// Start 开始计时
// 在timeout时间内没有执行reset命令，就会执行fun函数
func (t *TimeoutManager) Start(ctx context.Context, timeout time.Duration, option func()) {
	t.timer = time.NewTimer(timeout)
	t.timeout = timeout
	t.fun = option
	t.ctx, t.cancel = context.WithCancel(ctx)

	gosync.Go(t.ctx, func(ctx context.Context) {
		for {
			select {
			case <-t.timer.C:
				logger.WithCtx(ctx).Infof("[TimeoutManager] Timer expired, performing action")
				// Timer到期，执行命令
				t.cleanup()
				t.fun()
				return
			case <-t.resetChannel:
				logger.WithCtx(ctx).Infof("[TimeoutManager] reset")
				if !t.timer.Stop() {
					select {
					case <-t.timer.C:
					default:
					}
				}
				t.timer.Reset(t.timeout)
			case <-t.chancelChannel:
				logger.WithCtx(ctx).Infof("[TimeoutManager] cancel")
				t.cleanup()
				return
			case <-t.ctx.Done():
				logger.WithCtx(ctx).Infof("[TimeoutManager] context cancelled")
				t.cleanup()
				return
			}
		}
	})
}

// Reset 重置计时器
func (t *TimeoutManager) Reset() {
	select {
	case t.resetChannel <- struct{}{}:
	default:
		logger.WithCtx(t.ctx).Warnf("[TimeoutManager] Reset channel is full, ignoring reset")
	}
}

// Cancel 取消计时
func (t *TimeoutManager) Cancel() {
	select {
	case t.chancelChannel <- struct{}{}:
	default:
		logger.WithCtx(t.ctx).Warnf("[TimeoutManager] Cancel channel is full, ignoring cancel")
	}
}

// cleanup 清理资源
func (t *TimeoutManager) cleanup() {
	if t.cancel != nil {
		t.cancel()
	}
	if t.timer != nil {
		if !t.timer.Stop() {
			select {
			case <-t.timer.C:
			default:
			}
		}
	}
	t.resetChannel = nil
	t.chancelChannel = nil
}
