package lib

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/xunlbz/go-restful-template/pkg/log"
)

var mutex sync.Mutex

type DockerClient struct {
	cli *client.Client
}

func NewDockerClient() (dc DockerClient) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error(err)
		return
	}
	dc.cli = cli
	return dc
}

func (dc DockerClient) GetContainerList() []types.Container {

	containers, err := dc.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
	}
	for _, container := range containers {
		log.Debugf("get docker containers: %s %s %s", container.ID[:10], container.Names, container.Image)
	}
	return containers
}

func (dc DockerClient) GetInfo() types.Info {

	info, err := dc.cli.Info(context.Background())
	if err != nil {
		log.Error(err)
	}
	return info
}

func (dc DockerClient) GetImages() []types.ImageSummary {

	images, err := dc.cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Error(err)
	}
	return images
}

func (dc DockerClient) GetDiskUsage() types.DiskUsage {

	diskUsage, err := dc.cli.DiskUsage(context.Background())
	if err != nil {
		log.Error(err)
	}
	return diskUsage
}

func (dc DockerClient) ContainerRestart(containerId string, timeout time.Duration) error {
	return dc.cli.ContainerRestart(context.Background(), containerId, &timeout)
}

func (dc DockerClient) ContainerLog(containerId string) (io.ReadCloser, error) {
	options := types.ContainerLogsOptions{}
	options.Since = "1s" //GetNowDateTimeString()
	options.Follow = true
	options.ShowStderr = true
	options.ShowStdout = true
	options.Timestamps = true
	return dc.cli.ContainerLogs(context.Background(), containerId, options)
}
