package admin

import (
	"github.com/fansqz/fancode-backend/controller/utils"
	"github.com/fansqz/fancode-backend/models/po"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/system_service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SysApiController interface {
	// GetApiCount 获取api数目
	GetApiCount(ctx *gin.Context)
	// DeleteApiByID 删除api
	DeleteApiByID(ctx *gin.Context)
	// UpdateApi 更新api
	UpdateApi(ctx *gin.Context)
	// GetApiByID 根据id获取api
	GetApiByID(ctx *gin.Context)
	// GetApiTree 获取api树
	GetApiTree(ctx *gin.Context)
	// InsertApi 添加api
	InsertApi(ctx *gin.Context)
}

type sysApiController struct {
	sysApiService system_service.SysApiService
}

func NewSysApiController(apiService system_service.SysApiService) SysApiController {
	return &sysApiController{
		sysApiService: apiService,
	}
}

func (s *sysApiController) GetApiCount(ctx *gin.Context) {
	result := r.NewResult(ctx)
	count, err := s.sysApiService.GetApiCount(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(count)
}

func (s *sysApiController) DeleteApiByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	if err := s.sysApiService.DeleteApiByID(ctx, uint(id)); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("删除成功")
}

func (s *sysApiController) UpdateApi(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.AtoiOrDefault(ctx.PostForm("id"), 0)
	parentID := utils.AtoiOrDefault(ctx.PostForm("parentApiID"), 0)
	api := &po.SysApi{
		Model: gorm.Model{
			ID: uint(id),
		},
		ParentApiID: uint(parentID),
		Path:        ctx.PostForm("path"),
		Method:      ctx.PostForm("method"),
		Name:        ctx.PostForm("name"),
		Description: ctx.PostForm("description"),
	}
	if err := s.sysApiService.UpdateApi(ctx, api); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("修改成功")
}

func (s *sysApiController) GetApiByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	api, err := s.sysApiService.GetApiByID(ctx, uint(id))
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(api)
}

func (s *sysApiController) GetApiTree(ctx *gin.Context) {
	result := r.NewResult(ctx)
	apiTree, err := s.sysApiService.GetApiTree(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(apiTree)
}

func (s *sysApiController) InsertApi(ctx *gin.Context) {
	result := r.NewResult(ctx)
	parentID := utils.AtoiOrDefault(ctx.PostForm("parentApiID"), 0)
	api := &po.SysApi{
		ParentApiID: uint(parentID),
		Path:        ctx.PostForm("path"),
		Method:      ctx.PostForm("method"),
		Name:        ctx.PostForm("name"),
		Description: ctx.PostForm("description"),
	}
	id, err := s.sysApiService.InsertApi(ctx, api)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(id)
}
