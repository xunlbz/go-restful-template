package collector

import (
	"encoding/json"
	"fmt"

	"github.com/xunlbz/go-restful-template/pkg/lib"
)

var (
	label_docker                 = "docker"
	name_container_list          = "container_list"
	desc_container_list          = "容器列表"
	name_container_running_count = "container_running_count"
	desc_container_running_count = "运行中容器数量"
	name_container_count         = "container_count"
	desc_container_count         = "容器数量"
	name_image_count             = "image_count"
	desc_image_count             = "镜像数量"
	name_image_disk_usage        = "image_disk_usage"
	desc_image_disk_usage        = "镜像大小"
)

type dockerCollector struct {
	metrics      Metrics
	dockerClient lib.DockerClient
}

func init() {
	registerCollector(label_docker, NewDockerCollector)
}

func NewDockerCollector() (Collector, error) {
	c := new(dockerCollector)
	c.dockerClient = lib.NewDockerClient()
	return c, nil
}

func (c *dockerCollector) Update() Metrics {
	c.metrics = NewMetrics(label_docker, make([]Metric, 0))
	c.metrics.Entries = append(c.metrics.Entries, c.getContainerList())
	c.metrics.Entries = append(c.metrics.Entries, c.getInfo()...)
	c.metrics.Entries = append(c.metrics.Entries, c.getDiskUsage()...)
	return c.metrics
}

func (c *dockerCollector) Read() Metrics {
	c.Update()
	return c.metrics
}

func (c *dockerCollector) getContainerList() Metric {
	var metric Metric
	val := c.dockerClient.GetContainerList()
	data, err := json.Marshal(val)
	if err != nil {
		return metric
	}
	metric = NewMetric(label_docker, name_container_list, string(data), desc_container_list)
	return metric
}

func (c *dockerCollector) getInfo() (metrics []Metric) {
	info := c.dockerClient.GetInfo()
	metrics = append(metrics, NewMetric(label_docker, name_container_running_count, fmt.Sprint(info.ContainersRunning), desc_container_running_count))
	metrics = append(metrics, NewMetric(label_docker, name_container_count, fmt.Sprint(info.Containers), desc_container_count))
	metrics = append(metrics, NewMetric(label_docker, name_image_count, fmt.Sprint(info.Images), desc_image_count))
	return
}

func (c *dockerCollector) getDiskUsage() (metrics []Metric) {
	diskUsage := c.dockerClient.GetDiskUsage()
	metrics = append(metrics, NewMetric(label_docker, name_image_disk_usage, fmt.Sprintf("%.3fGB", float64(diskUsage.LayersSize*1.0)/1e9), desc_image_disk_usage))
	return
}
