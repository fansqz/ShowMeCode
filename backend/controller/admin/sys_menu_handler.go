package admin

import (
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/controller/utils"
	"github.com/fansqz/fancode-backend/models/po"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/system_service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SysMenuController interface {
	// GetMenuCount 获取menu数目
	GetMenuCount(ctx *gin.Context)
	// DeleteMenuByID 删除menu
	DeleteMenuByID(ctx *gin.Context)
	// UpdateMenu 更新menu
	UpdateMenu(ctx *gin.Context)
	// GetMenuByID 根据id获取menu
	GetMenuByID(ctx *gin.Context)
	// GetMenuTree 获取menu树
	GetMenuTree(ctx *gin.Context)
	// InsertMenu 添加menu
	InsertMenu(ctx *gin.Context)
}

type sysMenuController struct {
	sysMenuService system_service.SysMenuService
}

func NewSysMenuController(menuService system_service.SysMenuService) SysMenuController {
	return &sysMenuController{
		sysMenuService: menuService,
	}
}

func (s *sysMenuController) GetMenuCount(ctx *gin.Context) {
	result := r.NewResult(ctx)
	count, err := s.sysMenuService.GetMenuCount(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(count)
}

func (s *sysMenuController) DeleteMenuByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	if err := s.sysMenuService.DeleteMenuByID(ctx, uint(id)); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("删除成功")
}

func (s *sysMenuController) UpdateMenu(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.AtoiOrDefault(ctx.PostForm("id"), 0)
	parentIDStr := ctx.PostForm("parentMenuID")
	parentID, err := strconv.Atoi(parentIDStr)
	if err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	menu := &po.SysMenu{
		Model: gorm.Model{
			ID: uint(id),
		},
		ParentMenuID: uint(parentID),
		Code:         ctx.PostForm("code"),
		Name:         ctx.PostForm("name"),
		Description:  ctx.PostForm("description"),
	}
	if err2 := s.sysMenuService.UpdateMenu(ctx, menu); err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("修改成功")
}

func (s *sysMenuController) GetMenuByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	menu, err := s.sysMenuService.GetMenuByID(ctx, uint(id))
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(menu)
}

func (s *sysMenuController) GetMenuTree(ctx *gin.Context) {
	result := r.NewResult(ctx)
	menuTree, err := s.sysMenuService.GetMenuTree(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(menuTree)
}

func (s *sysMenuController) InsertMenu(ctx *gin.Context) {
	result := r.NewResult(ctx)
	parentID := utils.AtoiOrDefault(ctx.PostForm("parentMenuID"), 0)
	menu := &po.SysMenu{
		ParentMenuID: uint(parentID),
		Code:         ctx.PostForm("code"),
		Name:         ctx.PostForm("name"),
		Description:  ctx.PostForm("description"),
	}
	id, err := s.sysMenuService.InsertMenu(ctx, menu)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(id)
}
