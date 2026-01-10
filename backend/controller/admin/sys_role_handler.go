package admin

import (
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/controller/utils"
	"github.com/fansqz/fancode-backend/models/po"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/system_service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SysRoleController 角色管理相关功能
type SysRoleController interface {
	// GetRoleByID 根据id获取role
	GetRoleByID(ctx *gin.Context)
	// InsertSysRole 添加角色
	InsertSysRole(ctx *gin.Context)
	// DeleteSysRole 删除角色
	DeleteSysRole(ctx *gin.Context)
	// GetSysRoleList 获取角色列表
	GetSysRoleList(ctx *gin.Context)
	// UpdateSysRole 更新角色
	UpdateSysRole(ctx *gin.Context)
	// UpdateRoleApis 更新角色api
	UpdateRoleApis(ctx *gin.Context)
	// UpdateRoleMenus 更新角色menu
	UpdateRoleMenus(ctx *gin.Context)
	// GetApiIDsByRoleID
	GetApiIDsByRoleID(ctx *gin.Context)
	// GetMenuIDsByRoleID 通过角色id获取该角色拥有的menuID
	GetMenuIDsByRoleID(ctx *gin.Context)
}

type sysRoleController struct {
	sysRoleService system_service.SysRoleService
}

func NewSysRoleController(roleService system_service.SysRoleService) SysRoleController {
	return &sysRoleController{
		sysRoleService: roleService,
	}
}

func (s *sysRoleController) GetRoleByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	roleID := utils.GetIntParamOrDefault(ctx, "id", 0)
	user, err2 := s.sysRoleService.GetRoleByID(ctx, uint(roleID))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(user)
}

func (s *sysRoleController) InsertSysRole(ctx *gin.Context) {
	result := r.NewResult(ctx)
	sysRole := &po.SysRole{}
	sysRole.Name = ctx.PostForm("name")
	sysRole.Description = ctx.PostForm("description")
	//插入
	pID, err := s.sysRoleService.InsertSysRole(ctx, sysRole)
	if err != nil {
		result.Error(err)
		return
	}
	result.Success("角色添加成功", pID)
}

func (s *sysRoleController) UpdateSysRole(ctx *gin.Context) {
	result := r.NewResult(ctx)
	sysRoleIDString := ctx.PostForm("id")
	sysRoleID := utils.AtoiOrDefault(sysRoleIDString, 0)
	sysRole := &po.SysRole{
		Model: gorm.Model{
			ID: uint(sysRoleID),
		},
		Name:        ctx.PostForm("name"),
		Description: ctx.PostForm("description"),
	}
	if err2 := s.sysRoleService.UpdateSysRole(ctx, sysRole); err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData("角色修改成功")
}

func (s *sysRoleController) DeleteSysRole(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	if err := s.sysRoleService.DeleteSysRole(ctx, uint(id)); err != nil {
		result.Error(err)
		return
	}
	result.SuccessData("删除成功")
}

// GetSysRoleList 读取一个列表的角色
func (s *sysRoleController) GetSysRoleList(ctx *gin.Context) {
	result := r.NewResult(ctx)
	pageQuery, err := utils.GetPageQueryByQuery(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	role := &po.SysRole{
		Name:        ctx.Query("name"),
		Description: ctx.Query("description"),
	}
	pageQuery.Query = role
	pageInfo, err := s.sysRoleService.GetSysRoleList(ctx, pageQuery)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(pageInfo)
}

type insertApisToRoleRequest struct {
	RoleID uint   `json:"roleID"`
	ApiIDs []uint `json:"apiIDs"`
}

func (s *sysRoleController) UpdateRoleApis(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var json insertApisToRoleRequest
	if err := ctx.BindJSON(&json); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err2 := s.sysRoleService.UpdateRoleApis(ctx, json.RoleID, json.ApiIDs); err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("添加成功")
}

type insertMenusToRoleRequest struct {
	RoleID  uint   `json:"roleID"`
	MenuIDs []uint `json:"menuIDs"`
}

func (s *sysRoleController) UpdateRoleMenus(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var json insertMenusToRoleRequest
	if err := ctx.BindJSON(&json); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err2 := s.sysRoleService.UpdateRoleMenus(ctx, json.RoleID, json.MenuIDs); err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("添加成功")
}

func (s *sysRoleController) GetApiIDsByRoleID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	apiIDs, err2 := s.sysRoleService.GetApiIDsByRoleID(ctx, uint(id))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(apiIDs)
}

func (s *sysRoleController) GetMenuIDsByRoleID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	menuIDs, err2 := s.sysRoleService.GetMenuIDsByRoleID(ctx, uint(id))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(menuIDs)
}
