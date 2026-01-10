package controller

import (
	e "github.com/fansqz/fancode-backend/common/error"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/common_service"
	"path"

	"github.com/gin-gonic/gin"
)

type FileController interface {
	// StartUpload 启动上传命令
	StartUpload(ctx *gin.Context)
	// Upload 上传分片
	Upload(ctx *gin.Context)
	// CheckChunkSet 检测分片的文件名称集合
	CheckChunkSet(ctx *gin.Context)
	// CancelUpload 取消上传
	CancelUpload(ctx *gin.Context)
	// CompleteUpload 完成大文件上传功能
	CompleteUpload(ctx *gin.Context)
}

type fileController struct {
	fileService common_service.FileService
}

func NewFileController(fileService common_service.FileService) FileController {
	return &fileController{
		fileService: fileService,
	}
}

func (f *fileController) StartUpload(ctx *gin.Context) {
	result := r.NewResult(ctx)
	temp, err := f.fileService.StartUpload(ctx)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(temp)
}

// Upload 上传分片
func (f *fileController) Upload(ctx *gin.Context) {
	result := r.NewResult(ctx)
	path := ctx.PostForm("path")
	fileHead, err := ctx.FormFile("chunk")
	if err != nil {
		result.Error(e.ErrBadRequest)
		return
	}
	if err2 := f.fileService.Upload(ctx, path, fileHead); err2 != nil {
		result.Error(err2)
	}
	result.SuccessMessage("success upload")
}

func (f *fileController) CheckChunkSet(ctx *gin.Context) {
	result := r.NewResult(ctx)
	path := ctx.PostForm("path")
	set, err := f.fileService.CheckChunkSet(ctx, path)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(set)
}

func (f *fileController) CancelUpload(ctx *gin.Context) {
	result := r.NewResult(ctx)
	path := ctx.PostForm("path")
	if err := f.fileService.CancelUpload(ctx, path); err != nil {
		result.Error(err)
		return
	}
	result.SuccessMessage("取消成功")
}

func (f *fileController) CompleteUpload(ctx *gin.Context) {
	result := r.NewResult(ctx)
	p := ctx.PostForm("path")
	fileName := ctx.PostForm("fileName")
	hash := ctx.PostForm("hash")
	hashType := ctx.PostForm("hashType")
	err := f.fileService.CompleteUpload(ctx, p, fileName, hash, hashType)
	if err != nil {
		result.Error(err)
		return
	}
	result.SuccessData(path.Join(p, fileName))
}
