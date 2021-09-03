//操作系统
package collector

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"runtime"
)

var (
	label_os                         = "os"
	name_os_system_arch              = "system_arch"
	desc_os_system_arch              = "操作系统平台架构"
	name_os_operation_system         = "operation_system"
	desc_os_operation_system         = "操作系统"
	name_os_operation_system_version = "operation_system_version"
	desc_os_operation_system_version = "操作系统版本号"
	name_os_cpu_mem_info             = "cpu_mem_info"
	desc_os_cpu_mem_info             = "系统CPU与内存"
)

type osCollector struct {
	metrics Metrics
}

func init() {
	registerCollector(label_os, NewOSCollector)
}

func NewOSCollector() (Collector, error) {
	return new(osCollector), nil
}

func (c *osCollector) Update() Metrics {
	c.metrics = NewMetrics(label_os, make([]Metric, 0))
	c.metrics.Entries = append(c.metrics.Entries, getSystemInfo())
	c.metrics.Entries = append(c.metrics.Entries, getCurrentVersion())
	c.metrics.Entries = append(c.metrics.Entries, getSystemArch())
	c.metrics.Entries = append(c.metrics.Entries, getSystemCPUAndMemInfo())
	return c.metrics
}

func (c *osCollector) Read() Metrics {
	c.Update()
	return c.metrics
}

func getSystemArch() Metric {
	return NewMetric(label_os, name_os_system_arch, runtime.GOARCH, desc_os_system_arch)
}

func getSystemCPUAndMemInfo() Metric {
	memory, _ := mem.VirtualMemory()
	value := fmt.Sprintf("%d核 %vG", runtime.NumCPU(), memory.Total/1e9)
	return NewMetric(label_os, name_os_cpu_mem_info, value, desc_os_cpu_mem_info)
}
