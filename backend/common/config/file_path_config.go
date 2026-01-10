package config

import "gopkg.in/ini.v1"

type FilePathConfig struct {
	ProblemFileDir             string `ini:"problemFileDir"`             //题目文件目录
	ProblemDescriptionTemplate string `ini:"problemDescriptionTemplate"` //题目描述模板文位置
	ProblemFileTemplate        string `ini:"problemFileTemplate"`        //题目编程文件的模板文件
	TempDir                    string `ini:"tmpDir"`                     //临时目录

}

func NewFilePathConfig(cfg *ini.File) *FilePathConfig {
	filePathConfig := &FilePathConfig{}
	cfg.Section("filePath").MapTo(filePathConfig)
	return filePathConfig
}
