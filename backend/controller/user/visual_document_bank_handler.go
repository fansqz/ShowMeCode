package user

import (
	"github.com/fansqz/fancode-backend/controller/utils"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/visual_document_service"

	"github.com/gin-gonic/gin"
)

// VisualDocumentBankController
// @Description: 知识库相关功能
type VisualDocumentBankController interface {
	// GetAllVisualDocumentBank 获取所有知识库
	GetAllVisualDocumentBank(ctx *gin.Context)
	// GetVisualDocumentBankByID 读取知识库信息
	GetVisualDocumentBankByID(ctx *gin.Context)
}

type visualDocumentBankController struct {
	visualDocumentBankService visual_document_service.VisualDocumentBankService
}

func NewVisualDocumentBankController(bankService visual_document_service.VisualDocumentBankService) VisualDocumentBankController {
	return &visualDocumentBankController{
		visualDocumentBankService: bankService,
	}
}

func (v *visualDocumentBankController) GetAllVisualDocumentBank(ctx *gin.Context) {
	result := r.NewResult(ctx)
	enable := true
	banks, err := v.visualDocumentBankService.GetAllVisualDocumentBank(ctx, &enable)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(banks)
}

func (v *visualDocumentBankController) GetVisualDocumentBankByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	bank, err := v.visualDocumentBankService.GetVisualDocumentBankByID(ctx, uint(id))
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(bank)
}
