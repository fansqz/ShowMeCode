package debug_core

//
//import (
//	"github.com/fansqz/fancode-backend/constants"
//	"context"
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"testing"
//	"time"
//)
//
//func TestDebugger_Run(t *testing.T) {
//	ctx := context.Background()
//	docker := NewDebugger()
//	code := `package main
//
//import (
//	"fmt"
//)
//
//func main() {
//	a := 0
//	b := 0
//	fmt.Printf("start\n")
//   _, _ = fmt.Scan(&a, &b)
//	fmt.Printf("a + b = %b\n", a + b)
//}
//`
//	cha := make(chan interface{}, 10)
//	option := &Option{
//		Language:       constants.LanguageGo,
//		Code:           code,
//		CompileTimeout: 30 * time.Second,
//		OptionTimeout:  10 * time.Second,
//		DebugTimeout:   10 * 60 * time.Second,
//		BreakPoints:    []int{9},
//		// Callback 事件回调
//		Callback:    func(data interface{}) { cha <- data },
//		MemoryLimit: 1024 * 1024 * 1024,
//		CPUQuota:    100000,
//		TempDir:     "/var/fancode/debugger",
//	}
//	err := docker.Start(ctx, option)
//	assert.Nil(t, err)
//	event := <-cha
//	fmt.Printf("event: %v", event)
//	event = <-cha
//	fmt.Printf("event: %v", event)
//	// 设置并删除断点
//	sf, err := docker.GetStackTrace(ctx)
//	assert.Nil(t, err)
//	fmt.Println(sf)
//	v, err := docker.GetFrameVariables(ctx, sf[0].ID)
//	fmt.Println(v)
//	err = docker.Continue(ctx)
//	event = <-cha
//	event = <-cha
//	err = docker.Send(ctx, "2 3\n")
//	event = <-cha
//	fmt.Printf("event: %v", event)
//	assert.Nil(t, err)
//	event = <-cha
//	event = <-cha
//	event = <-cha
//	fmt.Printf("event: %v", event)
//}
