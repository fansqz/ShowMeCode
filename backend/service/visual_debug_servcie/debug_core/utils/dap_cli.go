package utils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/utils/gosync"
	"github.com/google/go-dap"
)

// EventCallback event回调通知
type EventCallback func(event dap.EventMessage)

// RequestCallback request回调通知，这部分是dap调试器主动发送的请求
type RequestCallback func(request dap.RequestMessage)

// AsyncCallback SendAsync方法异步处理命令时的回调函数
type AsyncCallback func(notification dap.ResponseMessage)

// DapClient DAP 客户端接口
type DapClient interface {
	// Close 关闭客户端连接
	Close()
	// SendSync 同步发送请求到dap
	SendSync(request dap.Message) (dap.Message, error)
	// SendWithTimeout 同步发送请求到dap，带超时
	SendWithTimeout(request dap.Message, timeout time.Duration) (dap.Message, error)
	// SendAsync 异步发送请求至dap
	SendAsync(request dap.Message, callback AsyncCallback) error
	// Send 发送请求
	Send(request dap.Message) error
	// ReadMessage 从dap client中读取响应
	ReadMessage() (dap.Message, error)
	// NewRequest 创建新请求
	NewRequest(command string) *dap.Request
}

// dapClient DapClient 接口的实现
type dapClient struct {
	conn   net.Conn
	reader *bufio.Reader
	// seq 用于跟踪每个请求的序列号
	seq          int
	pending      sync.Map
	pendingAsync sync.Map
	// eventCallback 事件回调
	eventCallback EventCallback
	// requestCallback 请求回调（dap主动发起的请求）
	requestCallback RequestCallback

	seqMutex sync.RWMutex
}

// NewDapClient creates a new DapClient over a TCP connection.
// Call Close() to close the connection.
func NewDapClient(ctx context.Context, addr string, callback EventCallback) (DapClient, error) {
	fmt.Println("Connecting to server at:", addr)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("new conn fail, err = %v", err))
	}
	return newDapClientFromConn(ctx, conn, callback), nil
}

// newDapClientFromConn creates a new dapClient with the given TCP connection.
func newDapClientFromConn(ctx context.Context, conn net.Conn, callback EventCallback) *dapClient {
	c := &dapClient{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
	c.seq = 1 // match VS Code numbering
	c.eventCallback = callback
	go c.recordReader(ctx)
	return c
}

// Close closes the client connection.
func (c *dapClient) Close() {
	c.conn.Close()
}

// SendSync 同步发送请求到dap
func (c *dapClient) SendSync(request dap.Message) (dap.Message, error) {
	pending := make(chan dap.ResponseMessage)
	c.pending.Store(request.GetSeq(), pending)
	if err := c.Send(request); err != nil {
		return nil, err
	}

	// 同步等待结果
	result := <-pending
	c.pending.Delete(request.GetSeq())
	return result, nil
}

// SendWithTimeout 同步发送请求到dap，带超时
func (c *dapClient) SendWithTimeout(request dap.Message, timeout time.Duration) (dap.Message, error) {
	pending := make(chan dap.ResponseMessage)
	c.pending.Store(request.GetSeq(), pending)
	if err := c.Send(request); err != nil {
		return nil, err
	}

	// 启动一个定时器，当超时时间到达时，定时器的通道会接收到信号
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// 同步等待结果或超时
	var result dap.Message
	select {
	case result = <-pending:
		// 收到响应，从 pending 映射中移除该请求序列号对应的通道
		c.pending.Delete(request.GetSeq())
		return result, nil
	case <-timer.C:
		// 超时，从 pending 映射中移除该请求序列号对应的通道
		c.pending.Delete(request.GetSeq())
		return nil, fmt.Errorf("timeout: %v", timeout)
	}
}

// SendAsync 异步发送请求至dap
func (c *dapClient) SendAsync(request dap.Message, callback AsyncCallback) error {
	c.pendingAsync.Store(request.GetSeq(), callback)
	if err := c.Send(request); err != nil {
		return err
	}
	return nil
}

// Send 发送请求
func (c *dapClient) Send(request dap.Message) error {
	return dap.WriteProtocolMessage(c.conn, request)
}

// recordReader 循环接收dap响应
func (c *dapClient) recordReader(ctx context.Context) {
	for {
		message, err := c.ReadMessage()
		if err != nil {
			logger.WithCtx(ctx).Warnf("read message fail, err = %v", err)
			break
		}
		if event, ok := message.(dap.EventMessage); ok {
			if c.eventCallback != nil {
				gosync.Go(ctx, func(ctx context.Context) {
					c.eventCallback(event)
				})
			}
			continue
		}
		if response, ok := message.(dap.ResponseMessage); ok {
			seq := response.GetResponse().RequestSeq
			// 检查是否是同步发送
			if pending, ok := c.pending.Load(seq); ok {
				if p, ok := pending.(chan dap.ResponseMessage); ok {
					p <- response
				}
			}
			// 异步发送
			if callback, ok := c.pendingAsync.Load(seq); ok {
				if cb, ok := callback.(AsyncCallback); ok {
					cb(response)
				}
			}
		}
		if request, ok := message.(dap.RequestMessage); ok {
			if c.requestCallback != nil {
				c.requestCallback(request)
			}
			continue
		}
	}
}

// ReadMessage 从dap client中读取响应，包括 response和event
func (c *dapClient) ReadMessage() (dap.Message, error) {
	return dap.ReadProtocolMessage(c.reader)
}

// NewRequest 创建新请求
func (c *dapClient) NewRequest(command string) *dap.Request {
	request := &dap.Request{}
	request.Type = "request"
	request.Command = command
	c.seqMutex.Lock()
	request.Seq = c.seq
	c.seq++
	c.seqMutex.Unlock()
	return request
}
