package main

import (
	"flag"
	"fmt"
	"github.com/fansqz/ShowMeCode/go-debugger/constants"
	"github.com/fansqz/ShowMeCode/go-debugger/debugger"
	"github.com/fansqz/ShowMeCode/go-debugger/debugger/c_debugger"
	"github.com/fansqz/ShowMeCode/go-debugger/debugger/cpp_debugger"
	"github.com/fansqz/ShowMeCode/go-debugger/utils"
	"github.com/google/go-dap"
	"log"
	"net"
	"os"
	"path"
)

var ConnList []net.Conn

// 定义版本号
const Version = "1.0.1"

func main() {
	//启动日志
	SetupLogger()
	defer CloseLogger()

	showVersion := flag.Bool("version", false, "Show the version number")
	port := flag.String("port", "8889", "TCP port to listen on")
	execFile := flag.String("file", "", "Exec file")
	codeFile := flag.String("codeFile", "", "main file code")
	language := flag.String("language", "c", "Program language")
	flag.Parse()

	// 检查是否需要显示版本信息
	if *showVersion {
		fmt.Printf("Version: %s\n", Version)
		return
	}
	if language == nil || *language == "" {
		fmt.Println("language cannot be empty")
		return
	}
	var code string
	if codeFile != nil && *codeFile != "" {
		// 编译文件
		data, _ := os.ReadFile(*codeFile)
		code = string(data)
	}
	if execFile == nil || *execFile == "" {
		fmt.Println("exec file cannot be empty")
		return
	}

	// 监听端口
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		fmt.Printf("listening at: %s\n", *port)
		return
	}
	defer listener.Close()
	fmt.Printf("started listening at: %s\n", listener.Addr().String())

	// 启动调试器
	debug, err := createDebugger(*language, *execFile, code)
	if err != nil {
		log.Printf("start debug fail, err = %s\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		ConnList = append(ConnList, conn)
		if err != nil {
			log.Printf("Connection failed: %v\n", err)
			continue
		}
		// Handle multiple client connections concurrently
		go handleConnection(conn, debug)
	}
}

// createDebugger 创建调试器
func createDebugger(language string, execFile string, code string) (debugger.Debugger, error) {
	var d debugger.Debugger
	switch language {
	case string(constants.LanguageC):
		d = c_debugger.NewCDebugger()
	case string(constants.LanguageCpp):
		d = cpp_debugger.NewCPPDebugger()
	}
	err := d.Start(&debugger.StartOption{
		ExecFile: execFile,
		MainCode: code,
		Callback: func(event dap.EventMessage) {
			for _, conn := range ConnList {
				dap.WriteProtocolMessage(conn, event)
			}
		},
	})
	return d, err
}

// 编译文件
func compileFile(language string, code string) (string, error) {
	switch language {
	case string(constants.LanguageC):
		workPath := path.Join("/var/fanCode/tempDir", utils.GetUUID())
		return c_debugger.CompileCFile(workPath, code)
	case string(constants.LanguageCpp):
		workPath := path.Join("/var/fanCode/tempDir", utils.GetUUID())
		return cpp_debugger.CompileCPPFile(workPath, code)
	}
	return "", fmt.Errorf("language not support")
}
