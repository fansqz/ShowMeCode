package user

import (
	"strconv"

	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/models/dto"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/user_saved_code_service"
	"github.com/gin-gonic/gin"
)

type UserSavedCodeHandler struct {
	savedCodeService user_saved_code_service.UserSavedCodeService
}

func NewUserSavedCodeHandler(savedCodeService user_saved_code_service.UserSavedCodeService) *UserSavedCodeHandler {
	return &UserSavedCodeHandler{
		savedCodeService: savedCodeService,
	}
}

// CreateUserSavedCode 创建用户保存的代码
func (h *UserSavedCodeHandler) CreateUserSavedCode(c *gin.Context) {
	result := r.NewResult(c)
	var req dto.UserSavedCodeDtoForCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithCtx(c).Errorf("[CreateUserSavedCode] bind json fail, err = %v", err)
		result.SimpleErrorMessage("参数错误: " + err.Error())
		return
	}

	data, err := h.savedCodeService.CreateUserSavedCode(c, &req)
	if err != nil {
		logger.WithCtx(c).Errorf("[CreateUserSavedCode] create fail, err = %v", err)
		result.Error(err)
		return
	}

	result.SuccessData(data)
}

// UpdateUserSavedCode 更新用户保存的代码
func (h *UserSavedCodeHandler) UpdateUserSavedCode(c *gin.Context) {
	result := r.NewResult(c)
	var req dto.UserSavedCodeDtoForUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithCtx(c).Errorf("[UpdateUserSavedCode] bind json fail, err = %v", err)
		result.SimpleErrorMessage("参数错误: " + err.Error())
		return
	}

	data, err := h.savedCodeService.UpdateUserSavedCode(c, &req)
	if err != nil {
		logger.WithCtx(c).Errorf("[UpdateUserSavedCode] update fail, err = %v", err)
		result.Error(err)
		return
	}

	result.SuccessData(data)
}

// DeleteUserSavedCode 删除用户保存的代码
func (h *UserSavedCodeHandler) DeleteUserSavedCode(c *gin.Context) {
	result := r.NewResult(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.WithCtx(c).Errorf("[DeleteUserSavedCode] parse id fail, err = %v", err)
		result.SimpleErrorMessage("无效的ID")
		return
	}

	err = h.savedCodeService.DeleteUserSavedCode(c, uint(id))
	if err != nil {
		logger.WithCtx(c).Errorf("[DeleteUserSavedCode] delete fail, err = %v", err)
		result.Error(err)
		return
	}

	result.SuccessMessage("删除成功")
}

// GetUserSavedCodeByID 根据ID获取用户保存的代码
func (h *UserSavedCodeHandler) GetUserSavedCodeByID(c *gin.Context) {
	result := r.NewResult(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.WithCtx(c).Errorf("[GetUserSavedCodeByID] parse id fail, err = %v", err)
		result.SimpleErrorMessage("无效的ID")
		return
	}

	data, err := h.savedCodeService.GetUserSavedCodeByID(c, uint(id))
	if err != nil {
		logger.WithCtx(c).Errorf("[GetUserSavedCodeByID] get fail, err = %v", err)
		result.Error(err)
		return
	}

	result.SuccessData(data)
}

// GetUserSavedCodeList 获取用户保存的代码列表
func (h *UserSavedCodeHandler) GetUserSavedCodeList(c *gin.Context) {
	result := r.NewResult(c)
	var req dto.UserSavedCodeDtoForQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.WithCtx(c).Errorf("[GetUserSavedCodeList] bind query fail, err = %v", err)
		result.SimpleErrorMessage("参数错误: " + err.Error())
		return
	}

	// 处理文档ID参数
	if documentIDStr := c.Query("documentID"); documentIDStr != "" {
		if documentID, err := strconv.ParseUint(documentIDStr, 10, 32); err == nil {
			req.DocumentID = &[]uint{uint(documentID)}[0]
		}
	}

	pageInfo, err := h.savedCodeService.GetUserSavedCodeList(c, &req)
	if err != nil {
		logger.WithCtx(c).Errorf("[GetUserSavedCodeList] get list fail, err = %v", err)
		result.Error(err)
		return
	}

	result.SuccessData(pageInfo)
}
