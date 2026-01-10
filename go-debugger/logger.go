package main

import (
	"log"
	"os"
)

var logFile *os.File
var logPath = "/var/godebugger.log"

func SetupLogger() {
	// 检查文件是否存在
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		// 文件不存在，创建文件
		_, createErr := os.Create(logPath)
		if createErr != nil {
			return
		}
	} else if err != nil {
		return
	}

	// 打开文件
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	log.SetOutput(logFile)
}

func CloseLogger() {
	if logFile != nil {
		_ = logFile.Close()
	}
}
