package utils

import (
	"fmt"
	"net"
	"sync"
)

var DebugPortManager = NewPortManager(5000, 5100)

// PortManager 用于分配docker调试端口
type PortManager struct {
	startPort int
	endPort   int
	usedPorts map[int]bool
	mu        sync.Mutex
}

// NewPortManager 创建一个新的端口管理器
func NewPortManager(start, end int) PortManager {
	return PortManager{
		startPort: start,
		endPort:   end,
		usedPorts: make(map[int]bool),
	}

}

// GetPort 获取一个可用端口
func (pm *PortManager) GetPort() (int, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for port := pm.startPort; port <= pm.endPort; port++ {
		if !pm.usedPorts[port] {
			if isPortAvailable(port) {
				pm.usedPorts[port] = true
				return port, nil
			}
		}
	}
	return 0, fmt.Errorf("no available ports in the range %d-%d", pm.startPort, pm.endPort)
}

// ReleasePort 释放一个端口
func (pm *PortManager) ReleasePort(port int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.usedPorts[port] {
		delete(pm.usedPorts, port)
	}
}

// isPortAvailable 检查端口是否可用
func isPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}
