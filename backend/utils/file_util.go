package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

// 通过page和pageSize读取文件列表
func readFilePage(dirPath string, page, pageSize int) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	startIndex := (page - 1) * pageSize
	endIndex := page * pageSize

	if startIndex >= len(files) {
		return nil, fmt.Errorf("页码超出范围")
	}
	if endIndex > len(files) {
		endIndex = len(files)
	}

	return files[startIndex:endIndex], nil
}

// 读取文件总数
func countFiles(dirPath string) (int, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return 0, err
	}
	defer dir.Close()

	fileInfoList, err := dir.Readdir(-1)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, fileInfo := range fileInfoList {
		if fileInfo.Mode().IsRegular() {
			count++
		}
	}

	return count, nil
}

// 检测文件夹是否存在
func CheckFolderExists(folderPath string) bool {
	fileInfo, err := os.Stat(folderPath)
	if err == nil && fileInfo.IsDir() {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		fmt.Println("发生错误:", err)
		return false
	}
}

// CheckAndDeletePath 检测一个文件或文件夹是否存在，如果存在则删除
func CheckAndDeletePath(filename string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("文件 %s 不存在", filename)
	}

	// 存在则删除文件
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}
