package user

import (
	"github.com/fansqz/fancode-backend/controller/utils"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/visual_document_service"

	"github.com/gin-gonic/gin"
)

type VisualDocumentController interface {
	// GetVisualDocumentDirectory 获取可视化文档目录
	GetVisualDocumentDirectory(ctx *gin.Context)
	// GetVisualDocumentByID 根据id获取可视化文档
	GetVisualDocumentByID(ctx *gin.Context)
}

func NewVisualDocumentController(vd visual_document_service.VisualDocumentService) VisualDocumentController {
	return &visualDocumentController{
		visualDocumentService: vd,
	}
}

type visualDocumentController struct {
	visualDocumentService visual_document_service.VisualDocumentService
}

func (v *visualDocumentController) GetVisualDocumentDirectory(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "bankID", 0)
	enable := true
	answer, err := v.visualDocumentService.GetVisualDocumentDirectory(ctx, uint(id), &enable)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(answer)
}

func (v *visualDocumentController) GetVisualDocumentByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	enable := true
	document, err := v.visualDocumentService.GetVisualDocumentByID(ctx, uint(id), &enable)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(document)
}
