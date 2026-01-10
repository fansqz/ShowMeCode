package utils

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Config 容器配置
type Config struct {
	ImageName     string
	ContainerName string
	Memory        int64
	CPUQuota      int64
	Binds         []string
	PortMapping   [][]string
}

// DockerClient Docker 客户端接口
type DockerClient interface {
	// Exec 在容器中执行命令
	Exec(ctx context.Context, cmd []string) (output string, err error, exitCode int)
	// CopyToContainer 将文件内容复制到容器指定目录
	CopyToContainer(ctx context.Context, content []byte, destPath string, filename string) error
	// RemoveContainer 强制关停并删除容器
	RemoveContainer(ctx context.Context) error
	// Interrupt 中断信号处理
	Interrupt() error
	// GetContainerID 获取容器ID
	GetContainerID() string
	// GetDebugAttach 获取调试连接
	GetDebugAttach() types.HijackedResponse
	// SetDebugAttach 设置调试连接
	SetDebugAttach(attach types.HijackedResponse)
	// GetClient 获取底层 Docker 客户端
	GetClient() *client.Client
}

// dockerClient DockerClient 接口的实现
type dockerClient struct {
	client      *client.Client
	debugAttach types.HijackedResponse
	containerID string
	signalChan  chan os.Signal
}

// NewDockerClient 创建 Docker 客户端实例
func NewDockerClient(ctx context.Context, config *Config) (DockerClient, error) {
	portMappings := map[nat.Port][]nat.PortBinding{}
	for _, portMap := range config.PortMapping {
		portMappings[nat.Port(portMap[0]+"/tcp")] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: portMap[1],
			},
		}
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory:   config.Memory,
			CPUQuota: config.CPUQuota,
		},
		PortBindings: portMappings,
		Binds:        config.Binds,
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: config.ImageName,
		Tty:   true,
	}, hostConfig, nil, nil, config.ContainerName)
	if err != nil {
		return nil, err
	}

	if err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, err
	}

	dc := &dockerClient{
		client:      cli,
		containerID: resp.ID,
		signalChan:  make(chan os.Signal, 1),
	}

	signal.Notify(dc.signalChan, os.Interrupt, syscall.SIGTERM)
	go dc.handleSignals()

	return dc, nil
}

// Exec 在容器中执行命令
func (d *dockerClient) Exec(ctx context.Context, cmd []string) (string, error, int) {
	execConfig := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execIDResp, err := d.client.ContainerExecCreate(ctx, d.containerID, execConfig)
	if err != nil {
		return "", err, 0
	}

	execResp, err := d.client.ContainerExecAttach(ctx, execIDResp.ID, container.ExecStartOptions{})
	if err != nil {
		return "", err, 0
	}
	defer execResp.Close()

	output, err := io.ReadAll(execResp.Reader)
	if err != nil {
		return "", err, 0
	}

	inspect, err := d.client.ContainerExecInspect(ctx, execIDResp.ID)
	if err != nil {
		return "", err, 0
	}

	return string(output), nil, inspect.ExitCode
}

// CopyToContainer 将文件内容复制到容器指定目录
func (d *dockerClient) CopyToContainer(ctx context.Context, content []byte, destPath string, filename string) error {
	tarBuf, err := createTar(filename, content)
	if err != nil {
		return err
	}

	return d.client.CopyToContainer(ctx, d.containerID, destPath, tarBuf, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}

// RemoveContainer 强制关停并删除容器
func (d *dockerClient) RemoveContainer(ctx context.Context) error {
	if err := d.client.ContainerKill(ctx, d.containerID, "SIGKILL"); err != nil {
		inspect, errInspect := d.client.ContainerInspect(ctx, d.containerID)
		if errInspect != nil {
			return fmt.Errorf("failed to inspect container: %v", errInspect)
		}
		if inspect.State.Running {
			return fmt.Errorf("failed to kill container: %v", err)
		}
	}

	if err := d.client.ContainerRemove(ctx, d.containerID, container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}); err != nil {
		return fmt.Errorf("failed to remove container: %v", err)
	}

	if err := d.client.Close(); err != nil {
		return fmt.Errorf("failed to close docker client: %v", err)
	}

	return nil
}

// Interrupt 中断信号处理
func (d *dockerClient) Interrupt() error {
	close(d.signalChan)
	return nil
}

// GetContainerID 获取容器ID
func (d *dockerClient) GetContainerID() string {
	return d.containerID
}

// GetDebugAttach 获取调试连接
func (d *dockerClient) GetDebugAttach() types.HijackedResponse {
	return d.debugAttach
}

// SetDebugAttach 设置调试连接
func (d *dockerClient) SetDebugAttach(attach types.HijackedResponse) {
	d.debugAttach = attach
}

// GetClient 获取底层 Docker 客户端
func (d *dockerClient) GetClient() *client.Client {
	return d.client
}

// handleSignals 处理系统信号
func (d *dockerClient) handleSignals() {
	// for sig := range d.signalChan {
	//     fmt.Printf("Received signal: %v\n", sig)
	// }
}

// createTar 创建包含单个文件的 tar 归档
func createTar(filename string, content []byte) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	hdr := &tar.Header{
		Name:    filename,
		Mode:    0644,
		Size:    int64(len(content)),
		ModTime: time.Now(),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}

	if _, err := tw.Write(content); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
