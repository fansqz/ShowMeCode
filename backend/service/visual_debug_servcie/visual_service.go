package visual_debug_servcie

import (
	"context"

	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/service/visual_debug_servcie/debug_core"
)

type VisualService interface {
	// GetVisualDescription 获取用户代码的分析结果
	GetVisualDescription(ctx context.Context, debugID string) (*dto.VisualDescription, error)
	// StructVisual 结构体导向可视化数据结构
	StructVisual(ctx context.Context, request *dto.StructVisualRequest) (*debug_core.StructVisualData, error)
	// ArrayVisual 数组可视化请求
	ArrayVisual(ctx context.Context, request *dto.ArrayVisualRequest) (*debug_core.ArrayVisualData, error)
	// Array2DVisual 二维数组可视化请求
	Array2DVisual(ctx context.Context, request *dto.Array2DVisualRequest) (*debug_core.Array2DVisualData, error)
}

type visualService struct {
}

func NewVisualService() VisualService {
	return &visualService{}
}

// GetVisualDescription 获取用户代码的分析结果
func (v *visualService) GetVisualDescription(ctx context.Context, debugID string) (*dto.VisualDescription, error) {
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, debugID)
	if !ok {
		return nil, e.ErrDebuggerIsClosed
	}
	visualDescription, err := debugContext.VisualDescriptionAnalyzer.GetVisualDescription(ctx)
	if err != nil {
		return nil, err
	}
	return visualDescription, nil
}

// StructVisual 结构体导向可视化数据结构
func (v *visualService) StructVisual(ctx context.Context, request *dto.StructVisualRequest) (*debug_core.StructVisualData, error) {
	debugContext, _ := DebugSessionManage.GetDebugSession(ctx, request.DebugID)
	data, err := debugContext.Debugger.StructVisual(ctx, &request.Query)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[StructVisual] get struct visual error, err = %v", err)
		return nil, e.ErrMysql
	}
	return data, nil
}

// ArrayVisual 数组可视化
func (v *visualService) ArrayVisual(ctx context.Context, request *dto.ArrayVisualRequest) (*debug_core.ArrayVisualData, error) {
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, request.DebugID)
	if !ok {
		return nil, e.ErrDebuggerIsClosed
	}
	data, err := debugContext.Debugger.ArrayVisual(ctx, &request.Query)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[StructVisual] get variable visual error, err = %v", err)
		return nil, err
	}
	return data, nil
}

// Array2DVisual 二维数组可视化
func (v *visualService) Array2DVisual(ctx context.Context, request *dto.Array2DVisualRequest) (*debug_core.Array2DVisualData, error) {
	debugContext, ok := DebugSessionManage.GetDebugSession(ctx, request.DebugID)
	if !ok {
		return nil, e.ErrDebuggerIsClosed
	}
	data, err := debugContext.Debugger.Array2DVisual(ctx, &request.Query)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[StructVisual] get variable visual error, err = %v", err)
		return nil, err
	}
	return data, nil
}
