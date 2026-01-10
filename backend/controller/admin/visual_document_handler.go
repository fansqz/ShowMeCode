package admin

import (
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/controller/utils"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/visual_document_service"

	"github.com/gin-gonic/gin"
)

type VisualDocumentManageController interface {
	// GetVisualDocumentDirectory 获取可视化文档目录
	GetVisualDocumentDirectory(ctx *gin.Context)
	// GetVisualDocumentByID 根据id获取可视化文档
	GetVisualDocumentByID(ctx *gin.Context)

	// InsertVisualDocument 添加可视化文档
	InsertVisualDocument(ctx *gin.Context)

	// UpdateVisualDocument 更新可视化文档
	UpdateVisualDocument(ctx *gin.Context)
	// UpdateVisualDocumentDirectory 更新可视化文档的目录结构
	UpdateVisualDocumentDirectory(ctx *gin.Context)
	// DeleteVisualDocumentByID 删除可视化文档
	DeleteVisualDocumentByID(ctx *gin.Context)
}

func NewVisualDocumentManageController(vd visual_document_service.VisualDocumentService) VisualDocumentManageController {
	return &visualDocumentManageController{
		visualDocumentService: vd,
	}
}

type visualDocumentManageController struct {
	visualDocumentService visual_document_service.VisualDocumentService
}

func (v *visualDocumentManageController) GetVisualDocumentDirectory(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "bankID", 0)
	answer, err := v.visualDocumentService.GetVisualDocumentDirectory(ctx, uint(id), nil)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(answer)
}

func (v *visualDocumentManageController) GetVisualDocumentByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	document, err := v.visualDocumentService.GetVisualDocumentByID(ctx, uint(id), nil)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(document)
}

func (v *visualDocumentManageController) InsertVisualDocument(ctx *gin.Context) {
	result := r.NewResult(ctx)
	document := po.VisualDocument{}
	if err := ctx.BindJSON(&document); err != nil {
		result.Error(err)
		return
	}
	if document.BankID == 0 || document.Title == "" {
		result.Error(e.ErrBadRequest)
	}
	id, err := v.visualDocumentService.InsertVisualDocument(ctx, &document)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(id)
}

func (v *visualDocumentManageController) UpdateVisualDocument(ctx *gin.Context) {
	result := r.NewResult(ctx)
	document := dto.VisualDocumentDto{}
	if err := ctx.BindJSON(&document); err != nil {
		result.Error(err)
		return
	}
	if err := v.visualDocumentService.UpdateVisualDocument(ctx, &document); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("更新成功")
}

func (v *visualDocumentManageController) UpdateVisualDocumentDirectory(ctx *gin.Context) {
	result := r.NewResult(ctx)
	req := dto.UpdateVisualDocumentReq{}
	if err := ctx.BindJSON(&req); err != nil {
		result.Error(err)
		return
	}
	if err := v.visualDocumentService.UpdateVisualDocumentDirectory(ctx, &req); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("更新成功")
}

func (v *visualDocumentManageController) DeleteVisualDocumentByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	if err := v.visualDocumentService.DeleteVisualDocumentByID(ctx, uint(id)); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("删除成功")
}
