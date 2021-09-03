// ubuntu 需要手动安装 apt install libsystemd-dev
package lib

import (
	"io"
	"time"
)

type Service struct {
	Name        string
	LoadState   string
	ActiveState string
	SubState    string
	Description string
}

// ListUnit 获取服务列表,匹配查询
func ListUnit(patterns []string) (units []Service, err error) {

	return
}

// RestartUnit 重启服务
func RestartUnit(name string) error {
	return nil
}

//JournalLog 获取服务日志
func JournalFollow(name string, writer io.Writer, end <-chan time.Time) (err error) {

	return nil
}
