package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	cgCPUPathPrefix    = "/sys/fs/cgroup/cpu/"
	cgMemoryPathPrefix = "/sys/fs/cgroup/memory/"
)

type CGroup struct {
	containerID string
}

/**
 * Cgroup 对用户代码的内存和cpu做限制
 */
func NewCGroup(containerID string) (*CGroup, error) {
	preFixs := []string{cgCPUPathPrefix, cgMemoryPathPrefix}
	for _, prefix := range preFixs {
		cgroupDir := filepath.Join(prefix, containerID)
		if err := os.MkdirAll(cgroupDir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return &CGroup{
		containerID: containerID,
	}, nil
}

// AddPID 添加进程到cgroup组
func (c *CGroup) AddPID(pid int) error {
	preFixs := []string{cgCPUPathPrefix, cgMemoryPathPrefix}
	for _, prefix := range preFixs {
		cgroupDir := filepath.Join(prefix, c.containerID)
		path := filepath.Join(cgroupDir, "tasks")
		if err := os.WriteFile(path, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
			c.Release()
			fmt.Println(err)
			return err
		}
	}
	return nil
}

// SetCPUQuota 设置CPU配额
func (c *CGroup) SetCPUQuota(quota int64) error {
	cgroupDir := filepath.Join(cgCPUPathPrefix, c.containerID)
	cpuQuotaFile := filepath.Join(cgroupDir, "cpu.cfs_quota_us")
	if err := os.WriteFile(cpuQuotaFile, []byte(fmt.Sprintf("%d", quota)), 0644); err != nil {
		return err
	}
	return nil
}

// SetMemoryLimit 设置内存限制
func (c *CGroup) SetMemoryLimit(limit int64) error {
	cgroupDir := filepath.Join(cgMemoryPathPrefix, c.containerID)
	memoryLimitFile := filepath.Join(cgroupDir, "memory.limit_in_bytes")
	if err := os.WriteFile(memoryLimitFile, []byte(fmt.Sprintf("%d", limit)), 0644); err != nil {
		return err
	}
	return nil
}

func (c *CGroup) Release() error {
	preFixs := []string{cgCPUPathPrefix, cgMemoryPathPrefix}
	for _, prefix := range preFixs {
		cgroupDir := filepath.Join(prefix, c.containerID)
		err := os.Remove(cgroupDir)
		if err != nil {
			return err
		}
	}
	return nil
}
