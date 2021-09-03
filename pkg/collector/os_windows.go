//操作系统
//docs https://pkg.go.dev/golang.org/x/sys/windows/registry
package collector

import (
	"github.com/xunlbz/go-restful-template/pkg/log"
	"golang.org/x/sys/windows/registry"
)

func getValueFormRegistry(key string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	v, _, err := k.GetStringValue(key)
	return v, err
}

func getSystemInfo() Metric {
	pn, err := getValueFormRegistry("ProductName")
	if err != nil {
		log.Error(err)
	}
	return NewMetric(label_os, name_os_operation_system, pn, desc_os_operation_system)
}

func getCurrentVersion() Metric {
	cv, err := getValueFormRegistry("ReleaseId")
	if err != nil {
		log.Error(err)
	}
	return NewMetric(label_os, name_os_operation_system_version, cv, desc_os_operation_system_version)
}
