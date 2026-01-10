package admin

import (
	"github.com/fansqz/fancode-backend/controller/utils"
	"github.com/fansqz/fancode-backend/models/po"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/visual_document_service"

	"github.com/gin-gonic/gin"
)

// VisualDocumentBankManageController
// @Description: 知识库管理相关功能
type VisualDocumentBankManageController interface {
	// InsertVisualDocumentBank 添加知识库
	InsertVisualDocumentBank(ctx *gin.Context)
	// UpdateVisualDocumentBank 更新知识库
	UpdateVisualDocumentBank(ctx *gin.Context)
	// DeleteVisualDocumentBank 删除知识库
	DeleteVisualDocumentBank(ctx *gin.Context)
	// GetVisualDocumentBankByID 读取知识库信息
	GetVisualDocumentBankByID(ctx *gin.Context)
	// GetAllVisualDocumentBank 获取所有知识库
	GetAllVisualDocumentBank(ctx *gin.Context)
}

type visualDocumentBankManageController struct {
	visualDocumentBankService visual_document_service.VisualDocumentBankService
}

func NewVisualDocumentBankManageController(bankService visual_document_service.VisualDocumentBankService) VisualDocumentBankManageController {
	return &visualDocumentBankManageController{
		visualDocumentBankService: bankService,
	}
}

func (v *visualDocumentBankManageController) InsertVisualDocumentBank(ctx *gin.Context) {
	result := r.NewResult(ctx)
	bank := po.VisualDocumentBank{}
	if err := ctx.BindJSON(&bank); err != nil {
		result.Error(err)
		return
	}
	pID, err := v.visualDocumentBankService.InsertVisualDocumentBank(ctx, &bank)
	if err != nil {
		result.Error(err)
		return
	}
	result.Success("知识库添加成功", pID)
}

func (v *visualDocumentBankManageController) UpdateVisualDocumentBank(ctx *gin.Context) {
	result := r.NewResult(ctx)
	bank := po.VisualDocumentBank{}
	if err := ctx.BindJSON(&bank); err != nil {
		result.Error(err)
		return
	}
	if err := v.visualDocumentBankService.UpdateVisualDocumentBank(ctx, &bank); err != nil {
		result.Error(err)
		return
	}
	result.SuccessData("知识库修改成功")
}

func (v *visualDocumentBankManageController) DeleteVisualDocumentBank(ctx *gin.Context) {
	result := r.NewResult(ctx)
	// 读取id
	bankID := uint(utils.GetIntParamOrDefault(ctx, "id", 0))
	// 删除知识库
	if err := v.visualDocumentBankService.DeleteVisualDocumentBank(ctx, bankID); err != nil {
		result.Error(err)
		return
	}
	result.SuccessData("知识库删除成功")
}

func (v *visualDocumentBankManageController) GetVisualDocumentBankByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	bank, err := v.visualDocumentBankService.GetVisualDocumentBankByID(ctx, uint(id))
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(bank)
}

func (v *visualDocumentBankManageController) GetAllVisualDocumentBank(ctx *gin.Context) {
	result := r.NewResult(ctx)
	banks, err := v.visualDocumentBankService.GetAllVisualDocumentBank(ctx, nil)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(banks)
}
