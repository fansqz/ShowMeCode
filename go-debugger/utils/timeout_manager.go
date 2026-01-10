package utils

import (
	"context"
	"github.com/fansqz/go-debugger/utils/gosync"
	"github.com/sirupsen/logrus"
	"time"
)

// TimeoutManager 一个计时器
// 如果在timeout时间内没有执行reset命令，就会执行fun函数
// duration至少10s
type TimeoutManager struct {
	timer          *time.Timer
	timeout        time.Duration
	resetChannel   chan bool
	chancelChannel chan bool
	fun            func()
}

// NewTimeoutManager 创建一个新的计时器实例
func NewTimeoutManager() *TimeoutManager {
	return &TimeoutManager{}
}

// Start 开始计时
// 在timeout时间内没有执行reset命令，就会执行fun函数
func (t *TimeoutManager) Start(ctx context.Context, timeout time.Duration, option func()) {
	t.timer = time.NewTimer(timeout)
	t.timeout = timeout
	t.fun = option
	t.resetChannel = make(chan bool)
	t.chancelChannel = make(chan bool)
	gosync.Go(ctx, func(ctx context.Context) {
		for {
			select {
			case <-t.timer.C:
				logrus.Infof("[TimeoutManager] Timer expired, performing action")
				// Timer到期，执行命令
				t.fun()
				return
			case <-t.resetChannel:
				logrus.Infof("[TimeoutManager] reset")
				if !t.timer.Stop() {
					<-t.timer.C
				}
				t.timer.Reset(t.timeout)
			case <-t.chancelChannel:
				logrus.Infof("[TimeoutManager] chancel")
				if !t.timer.Stop() {
					<-t.timer.C // 确保Timer停止并清空通道
				}
				return
			}
		}
	})
}

// Reset 重置计时器
func (t *TimeoutManager) Reset() {
	t.resetChannel <- true
}

// Chancel 取消即使
func (t *TimeoutManager) Chancel() {
	t.chancelChannel <- true
}
