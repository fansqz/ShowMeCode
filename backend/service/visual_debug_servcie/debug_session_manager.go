package visual_debug_servcie

import (
	"context"
	"errors"
	"fmt"
	"github.com/fansqz/fancode-backend/common/config"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/service/visual_debug_servcie/ai_analyze_core"
	"github.com/fansqz/fancode-backend/service/visual_debug_servcie/debug_core"
	"github.com/fansqz/fancode-backend/utils"
	"sync"
)

type debugSessionManage struct {
	// debugID - debugSession
	debugContextMap sync.Map
}

// DebugSession 调试上下文对象，用于存储用户的一次调试的信息
type DebugSession struct {
	// 用于停止循环处理service返回的event
	DestroyEventChan chan struct{}
	// DtoEventChan 将event返回给用户的channel
	DtoEventChan chan interface{}
	// Debugger 用户的调试器
	Debugger debug_core.Debugger
	// 可视化描述分析器
	VisualDescriptionAnalyzer *ai_analyze_core.VisualDescriptionAnalyzer
}

var DebugSessionManage = &debugSessionManage{}

// CreateDebugSession 创建调试上下文
func (d *debugSessionManage) CreateDebugSession(ctx context.Context, aiConfig *config.AIConfig, key string) error {
	_, ok := d.debugContextMap.Load(key)
	if ok {
		d.DestroyDebugSession(ctx, key)
	}
	de := debug_core.NewDebugger()
	va := ai_analyze_core.NewVisualDescriptionAnalyzer(aiConfig)
	d.debugContextMap.Store(key, &DebugSession{
		DestroyEventChan:          make(chan struct{}, 2),
		DtoEventChan:              make(chan interface{}, 10),
		Debugger:                  de,
		VisualDescriptionAnalyzer: va,
	})
	return nil
}

// GetDebugSession 根据key获取一个debugger
func (d *debugSessionManage) GetDebugSession(ctx context.Context, id string) (*DebugSession, bool) {
	debugContext, ok := d.debugContextMap.Load(id)
	if !ok {
		return nil, ok
	}
	return debugContext.(*DebugSession), ok
}

// DestroyDebugSession 销毁一个debugger
func (d *debugSessionManage) DestroyDebugSession(ctx context.Context, id string) {
	dcontext, ok := d.debugContextMap.Load(id)
	debugContext, _ := dcontext.(*DebugSession)
	if !ok {
		return
	}
	if err := debugContext.Debugger.Terminate(ctx); err != nil {
		logger.WithCtx(ctx).Errorf("[DestroyDebugSession] debug terminate fail, err = %s", err)
	}

	close(debugContext.DestroyEventChan)
	close(debugContext.DtoEventChan)
	d.debugContextMap.Delete(id)
}

// SendEvent 向用户发送调试事件
func (d *debugSessionManage) SendEvent(ctx context.Context, id string, event interface{}) {
	session, ok := d.GetDebugSession(ctx, id)
	if !ok {
		logger.WithCtx(ctx).Warnf("[SendEvent], key not exist, key = %v", id)
		return
	}
	session.DtoEventChan <- event
}

func (d *debugSessionManage) SendDestroyEvent(ctx context.Context, id string) {
	session, ok := d.GetDebugSession(ctx, id)
	if !ok {
		logger.WithCtx(ctx).Warnf("[SendDestroyEvent], key not exist, key = %v", id)
		return
	}
	session.DestroyEventChan <- struct{}{}
}

type userDebugManage struct {
	// userID - debugID
	userDebugMap sync.Map
}

var UserDebugManage = &userDebugManage{}

func (u *userDebugManage) GetUserDebugID(ctx context.Context) (string, error) {
	userID := gerVisitorID(ctx)
	if userID == "" {
		logger.WithCtx(ctx).Infof("[userDebugManage] get debug id fail, user id not found")
		return "", errors.New("get user id fail")
	}
	value, ok := u.userDebugMap.Load(userID)
	if ok {
		return value.(string), nil
	}
	return "", nil
}

func (u *userDebugManage) StoreUserDebugID(ctx context.Context, debugID string) error {
	userID := gerVisitorID(ctx)
	if userID == "" {
		logger.WithCtx(ctx).Infof("[userDebugManage] get debug id fail, user id not found")
		return errors.New("get user id fail")
	}
	u.userDebugMap.Store(userID, debugID)
	return nil
}

func gerVisitorID(ctx context.Context) string {
	var visitorID string
	userID := utils.GetUserIDWithCtx(ctx)
	if userID == 0 {
		// 读取访客id
		visitorID = utils.GetVisitorIDWithCtx(ctx)
		if visitorID != "" {
			visitorID = "visitor:" + visitorID
		}
	} else {
		visitorID = fmt.Sprintf("%d", userID)
	}
	return visitorID
}
