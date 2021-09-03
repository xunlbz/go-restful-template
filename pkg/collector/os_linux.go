//操作系统
package collector

import (
	"github.com/xunlbz/go-restful-template/pkg/lib"
	"github.com/xunlbz/go-restful-template/pkg/log"
	"io/ioutil"
	"os"
	"strings"
)

func getSystemInfo() (m Metric) {
	op := ""
	if res, err := lib.ExecCommand("lsb_release", "-d"); err == nil {
		op = formatString(strings.Split(string(res), ":")[1])
	}
	return NewMetric(label_os, name_os_operation_system, op, desc_os_operation_system)
}

func getCurrentVersion() Metric {
	version := ""
	file, err := os.Open(procFilePath("sys/kernel/osrelease"))
	if err != nil {
		log.Error(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err == nil {
		version = formatString(string(data))
	} else {
		log.Error(err)
	}
	return NewMetric(label_os, name_os_operation_system_version, version, desc_os_operation_system_version)
}

func formatString(str string) string {
	s := strings.TrimLeft(str, " ")
	s = strings.TrimRight(str, " ")
	return s
}
