package common_service

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/fansqz/fancode-backend/common/config"
	conf "github.com/fansqz/fancode-backend/common/config"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/utils"
	"mime/multipart"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

// FileService 文件上传相关service
type FileService interface {
	// StartUpload 启动上传命令
	StartUpload(ctx context.Context) (string, error)
	// Upload 上传分片
	Upload(ctx context.Context, path string, file *multipart.FileHeader) error
	// CheckChunkSet 检测分片的文件名称集合
	CheckChunkSet(ctx context.Context, path string) ([]string, error)
	// CancelUpload 取消上传
	CancelUpload(ctx context.Context, path string) error
	// CompleteUpload 完成大文件上传功能
	CompleteUpload(ctx context.Context, path string, fileName string, hash string, hashType string) error
}

type fileService struct {
	config *conf.AppConfig
}

func NewFileService(config *conf.AppConfig) FileService {
	return &fileService{
		config: config,
	}
}

func (f *fileService) StartUpload(ctx context.Context) (string, error) {
	tempPath := getTempDir(f.config)
	err := os.MkdirAll(tempPath, 0755)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[StartUpload] make dir error, err = %v", err)
		return "", e.ErrServer
	}
	return tempPath, nil
}

func (f *fileService) Upload(ctx context.Context, path string, file *multipart.FileHeader) error {
	ginContext := ctx.(*gin.Context)
	err := ginContext.SaveUploadedFile(file, path)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[Upload] save upload test_file error, err = %v", err)
		return e.ErrServer
	}
	return nil
}

func (f *fileService) CheckChunkSet(ctx context.Context, path string) ([]string, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[CheckChunkSet] read dir error, err = %v", err)
		return nil, e.ErrServer
	}
	answer := make([]string, len(dirs))
	for i, a := range dirs {
		answer[i] = a.Name()
	}
	return answer, nil
}

// CancelUpload 取消上传
func (f *fileService) CancelUpload(ctx context.Context, path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[CancelUpload] remove test_file error, err = %v", err)
		return e.ErrServer
	}
	return nil
}

// CompleteUpload 完成大文件上传功能
func (f *fileService) CompleteUpload(ctx context.Context, p string, fileName string, h string, hashType string) error {
	// 读取path中的所有文件
	files, err := os.ReadDir(p)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[CompleteUpload] read dir error, err = %v", err)
		// 处理错误
		return e.ErrServer
	}

	// 创建结果文件
	resultFile, err := os.Create(fileName)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[CompleteUpload] create test_file error, err = %v", err)
		// 处理错误
		return e.ErrServer
	}
	defer resultFile.Close()

	// 遍历所有文件，逐个写入结果文件
	for _, file := range files {
		filePath := path.Join(p, file.Name())
		var fileData []byte
		fileData, err = os.ReadFile(filePath)
		if err != nil {
			logger.WithCtx(ctx).Errorf("[CompleteUpload] read test_file error, err = %v", err)
			// 处理错误
			return e.ErrServer
		}
		resultFile.Write(fileData)
	}

	hash2, err := hash(p, hashType)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[CompleteUpload] hash test_file error, err = %v", err)
		return err
	}
	if hash2 != h {
		return e.NewCustomMsg("hash miss match")
	}
	return nil
}

func hash(filePath string, hashType string) (string, error) {
	// 计算结果文件的哈希值，并与传入的哈希值进行比较
	switch hashType {
	case "md5":
		resultHash, err := calculateMD5(filePath)
		if err != nil {
			return "", e.ErrServer
		}
		return resultHash, nil
	case "sha1":
		resultHash, err := calculateSHA1(filePath)
		if err != nil {
			return "", e.ErrServer
		}
		return resultHash, nil
	case "sha256":
		resultHash, err := calculateSHA256(filePath)
		if err != nil {
			return "", e.ErrServer
		}
		return resultHash, nil
	default:
		// 不支持的哈希算法类型，处理错误
		return "", e.NewCustomMsg("hash type not support")
	}
}

// 计算MD5哈希值
func calculateMD5(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		// 处理错误
		return "", err
	}
	hash := md5.Sum(fileData)
	return hex.EncodeToString(hash[:]), nil
}

// 计算SHA1哈希值
func calculateSHA1(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		// 处理错误
		return "", err
	}
	hash := sha1.Sum(fileData)
	return hex.EncodeToString(hash[:]), nil
}

// 计算SHA256哈希值
func calculateSHA256(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		// 处理错误
		return "", err
	}
	hash := sha256.Sum256(fileData)
	return hex.EncodeToString(hash[:]), nil
}

// getTempDir 获取一个随机的临时文件夹
func getTempDir(config *config.AppConfig) string {
	uuid := utils.GetUUID()
	executePath := config.FilePathConfig.TempDir + "/" + uuid
	return executePath
}
