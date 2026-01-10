package file_store

import (
	"io"
)

type Store interface {
	// SaveFile 上传文件到对象存储
	SaveFile(storePath string, file io.Reader) error
	// DownloadFile 下载一个文件 storePath:对象存储的路劲，loadPath:本地路径
	DownloadFile(storePath, localPath string) error
	// DownloadFolder 下载文件夹, storePath：文件在对象存储中的位置，localPath：文件在本地中的位置
	DownloadFolder(storePath, localPath string) error
	// DownloadAndCompressFolder 下载并压缩文件夹，zipPath:压缩文件路劲，比如a/b/c.zip
	DownloadAndCompressFolder(storePath, localPath, zipPath string) error
	// ReadFile 读取文件
	ReadFile(storePath string) ([]byte, error)
	// UploadFolder 上传文件夹到对象存储
	UploadFolder(storePath string, localPath string)
	// DeleteFolder 删除文件夹
	DeleteFolder(storePath string) error
}
