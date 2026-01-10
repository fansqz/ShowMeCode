package utils

import "sync"

const (
	// Init 调试初始化状态
	Init = "Init"
	// Stopped 用户程序暂停
	Stopped = "stopped"
	// Running 用户程序运行中
	Running = "running"
	// Finish 调试结束状态
	Finish = "finish"
)

// StatusManager 记录调试器的状态的
type StatusManager struct {
	lock   sync.RWMutex
	status string
}

func NewStatusManager() *StatusManager {
	return &StatusManager{
		status: Init,
	}
}

func (s *StatusManager) Lock() {
	s.lock.Lock()
}

func (s *StatusManager) UnLock() {
	s.lock.Unlock()
}

func (s *StatusManager) Set(status string) {
	defer s.lock.Unlock()
	s.lock.Lock()
	s.status = status
}

func (s *StatusManager) Is(statusList ...string) bool {
	defer s.lock.RUnlock()
	s.lock.RLock()
	for _, status := range statusList {
		if s.status == status {
			return true
		}
	}
	return false
}
