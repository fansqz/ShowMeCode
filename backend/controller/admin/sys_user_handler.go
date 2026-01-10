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

type SysUserController interface {
	// GetUserByID 根据id获取user
	GetUserByID(ctx *gin.Context)
	// InsertSysUser 添加用户
	InsertSysUser(ctx *gin.Context)
	// UpdateSysUser 更新用户，但是不更新密码
	UpdateSysUser(ctx *gin.Context)
	// DeleteSysUser 删除用户
	DeleteSysUser(ctx *gin.Context)
	// GetSysUserList 获取用户列表
	GetSysUserList(ctx *gin.Context)
	// UpdateUserRoles 更新角色roleIDs
	UpdateUserRoles(ctx *gin.Context)
	// GetRoleIDsByUserID 通过用户id获取所有角色id
	GetRoleIDsByUserID(ctx *gin.Context)
	// GetAllSimpleRole 获取简单角色列表
	GetAllSimpleRole(ctx *gin.Context)
}

type sysUserController struct {
	sysUserService system_service.SysUserService
}

func NewSysUserController(userService system_service.SysUserService) SysUserController {
	return &sysUserController{
		sysUserService: userService,
	}
}

func (s *sysUserController) GetUserByID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	userID := utils.GetIntParamOrDefault(ctx, "id", 0)
	user, err2 := s.sysUserService.GetUserByID(ctx, uint(userID))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(user)
}

func (s *sysUserController) InsertSysUser(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id, err := s.sysUserService.InsertSysUser(ctx, &po.SysUser{
		LoginName: ctx.PostForm("loginName"),
		Username:  ctx.PostForm("username"),
		Password:  ctx.PostForm("password"),
		Email:     ctx.PostForm("email"),
		Phone:     ctx.PostForm("phone"),
	})
	if err != nil {
		result.Error(err)
		return
	}
	result.Success("添加成功", id)
}

func (s *sysUserController) UpdateSysUser(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.AtoiOrDefault(ctx.PostForm("id"), 0)
	user := &po.SysUser{
		Model: gorm.Model{
			ID: uint(id),
		},
		LoginName: ctx.PostForm("loginName"),
		Username:  ctx.PostForm("username"),
		Password:  ctx.PostForm("password"),
		Email:     ctx.PostForm("email"),
		Phone:     ctx.PostForm("phone"),
	}
	if err2 := s.sysUserService.UpdateSysUser(ctx, user); err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("更新成功")
}

func (s *sysUserController) DeleteSysUser(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	err2 := s.sysUserService.DeleteSysUser(ctx, uint(id))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("删除成功")
}

func (s *sysUserController) GetSysUserList(ctx *gin.Context) {
	result := r.NewResult(ctx)
	pageQuery, err := utils.GetPageQueryByQuery(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	user := &po.SysUser{
		Username:     ctx.Query("username"),
		LoginName:    ctx.Query("loginName"),
		Email:        ctx.Query("email"),
		Phone:        ctx.Query("phone"),
		Introduction: ctx.Query("introduction"),
	}
	sexStr := ctx.Query("sex")
	if sexStr == "1" {
		user.Sex = 1
	} else if sexStr == "2" {
		user.Sex = 2
	}
	pageQuery.Query = user
	pageInfo, err2 := s.sysUserService.GetSysUserList(ctx, pageQuery)
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(pageInfo)
}

type updateUserRolesRequest struct {
	UserID  uint   `json:"userID"`
	RoleIDs []uint `json:"roleIDs"`
}

func (s *sysUserController) UpdateUserRoles(ctx *gin.Context) {
	result := r.NewResult(ctx)
	var request updateUserRolesRequest
	if err := ctx.BindJSON(&request); err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err2 := s.sysUserService.UpdateUserRoles(ctx, request.UserID, request.RoleIDs); err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessMessage("更新成功")
}

func (s *sysUserController) GetRoleIDsByUserID(ctx *gin.Context) {
	result := r.NewResult(ctx)
	id := utils.GetIntParamOrDefault(ctx, "id", 0)
	roleIDs, err2 := s.sysUserService.GetRoleIDsByUserID(ctx, uint(id))
	if err2 != nil {
		result.Error(err2)
		return
	}
	result.SuccessData(roleIDs)
}

func (s *sysUserController) GetAllSimpleRole(ctx *gin.Context) {
	result := r.NewResult(ctx)
	roles, err := s.sysUserService.GetAllSimpleRole(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(roles)
}
