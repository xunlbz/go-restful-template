//操作系统
package collector

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"runtime"
	"time"
)

var (
	label_osuse                 = "osuse"
	name_osuse_cpu_used_percent = "cpu_used_percent"
	desc_osuse_cpu_used_percent = "CPU使用率"
	name_osuse_mem_used_info    = "mem_used_info"
	desc_osuse_mem_used_info    = "内存使用"
	name_osuse_mem_used_percent = "mem_used_percent"
	desc_osuse_mem_used_percent = "内存使用率"
	name_osuse_system_time      = "system_time"
	desc_osuse_system_time      = "系统时间"
	name_osuse_disk_usage       = "disk_usage "
	desc_osuse_disk_usage       = "硬盘使用情况"
)

type osuseCollector struct {
	metrics Metrics
}

func init() {
	registerCollector(label_osuse, NewOSUSECollector)
}

func NewOSUSECollector() (Collector, error) {
	return new(osuseCollector), nil
}

func (c *osuseCollector) Update() Metrics {
	c.metrics = NewMetrics(label_osuse, make([]Metric, 0))
	c.metrics.Entries = append(c.metrics.Entries, getSystemMemUsedInfo())
	c.metrics.Entries = append(c.metrics.Entries, getSystemMemUsedPercent())
	c.metrics.Entries = append(c.metrics.Entries, getSystemCpuUsedPercent())
	c.metrics.Entries = append(c.metrics.Entries, getSystemTime())
	c.metrics.Entries = append(c.metrics.Entries, getDiskUsage())
	return c.metrics
}

func (c *osuseCollector) Read() Metrics {
	c.Update()
	return c.metrics
}

func getSystemMemUsedInfo() Metric {
	memory, _ := mem.VirtualMemory()
	value := float64(memory.Used) / 1e9
	return NewMetric(label_osuse, name_osuse_mem_used_info, fmt.Sprintf("%.1f", value), desc_osuse_mem_used_info)
}

func getSystemMemUsedPercent() Metric {
	memory, _ := mem.VirtualMemory()
	return NewMetric(label_osuse, name_osuse_mem_used_percent, fmt.Sprintf("%.1f", memory.UsedPercent), desc_osuse_mem_used_percent)
}

func getSystemCpuUsedPercent() Metric {
	percent, _ := cpu.Percent(time.Second, false)
	return NewMetric(label_osuse, name_osuse_cpu_used_percent, fmt.Sprintf("%.1f", percent[len(percent)-1]), desc_osuse_cpu_used_percent)
}

func getSystemTime() Metric {

	return NewMetric(label_osuse, name_osuse_system_time, fmt.Sprintln(time.Now().Format("2006-01-02 15:04:05")), desc_osuse_system_time)
}

func getDiskUsage() Metric {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:"
	}
	ust, err := disk.Usage(path)
	str := fmt.Sprintf("%.1fG(已用)/%dG", float64(ust.Used*1.0)/1e9, ust.Total/1e9)
	if err != nil {
		str = ""
	}
	return NewMetric(label_osuse, name_osuse_disk_usage, str, desc_osuse_disk_usage)
}
