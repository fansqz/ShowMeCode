package user

import (
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/models/dto"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/visual_debug_servcie"

	"github.com/gin-gonic/gin"
)

// VisualController
// 可视化相关，依赖于debugger
type VisualController interface {
	// StructVisual 结构体导向可视化数据结构
	StructVisual(ctx *gin.Context)
	// ArrayVisual 数组可视化数据结构
	ArrayVisual(ctx *gin.Context)
	// Array2DVisual 二维数组可视化数据结构
	Array2DVisual(ctx *gin.Context)
	// GetVisualDescription 获取用户可视化描述解析结果
	GetVisualDescription(ctx *gin.Context)
}

type visualController struct {
	visualService visual_debug_servcie.VisualService
}

func NewVisualController(visualService visual_debug_servcie.VisualService) VisualController {
	return &visualController{
		visualService: visualService,
	}
}

// StructVisual 结构体导向可视化数据结构
func (v *visualController) StructVisual(ctx *gin.Context) {
	result := r.NewResult(ctx)
	req := dto.StructVisualRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	data, err := v.visualService.StructVisual(ctx, &req)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(data)
}

// ArrayVisual 变量为导向可视化数据结构
func (v *visualController) ArrayVisual(ctx *gin.Context) {
	result := r.NewResult(ctx)
	req := dto.ArrayVisualRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	data, err := v.visualService.ArrayVisual(ctx, &req)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(data)
}

func (v *visualController) Array2DVisual(ctx *gin.Context) {
	result := r.NewResult(ctx)
	req := dto.Array2DVisualRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	data, err := v.visualService.Array2DVisual(ctx, &req)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(data)
}

// GetVisualDescription 获取用户代码的分析结果
func (v *visualController) GetVisualDescription(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := ctx.Param("id")
	visualDescription, err := v.visualService.GetVisualDescription(ctx, id)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(visualDescription)
}
