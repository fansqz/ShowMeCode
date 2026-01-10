package controller

import (
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/common_service"

	"github.com/gin-gonic/gin"
)

// CommonController
type CommonController interface {
	// GetURL 根据path获取url
	GetURL(ctx *gin.Context)
}

type commonController struct {
	commonService common_service.CommonService
}

func NewCommonController(commonService common_service.CommonService) CommonController {
	return &commonController{
		commonService: commonService,
	}
}

func (c *commonController) GetURL(ctx *gin.Context) {
	result := r.NewResult(ctx)
	path := ctx.Query("path")
	url, err := c.commonService.GetURL(ctx, path)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(url)
}
