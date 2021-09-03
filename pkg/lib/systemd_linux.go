// ubuntu 需要手动安装 apt install libsystemd-dev
package lib

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/coreos/go-systemd/v22/sdjournal"
	"github.com/xunlbz/go-restful-template/pkg/log"
	"io"
	"time"
)

func newDbus() (*dbus.Conn, error) {
	return dbus.NewSystemdConnection()
}

// ListUnit 获取服务列表,匹配查询
func ListUnit(patterns []string) (units []dbus.UnitStatus, err error) {
	conn, err := newDbus()
	if err != nil {
		return
	}
	defer conn.Close()
	return conn.ListUnitsByPatterns([]string{}, patterns)
}

// RestartUnit 重启服务
func RestartUnit(name string) error {
	conn, err := newDbus()
	if err != nil {
		return err
	}
	defer conn.Close()
	reschan := make(chan string)
	_, err = conn.RestartUnit(name, "replace", reschan)
	if err != nil {
		log.Error(err)
		return err
	}

	job := <-reschan
	if job != "done" {
		return fmt.Errorf("job is not done: %s", job)
	}
	return nil
}

//JournalLog 获取服务日志
func JournalFollow(name string, writer io.Writer, end <-chan time.Time) (err error) {
	config := sdjournal.JournalReaderConfig{
		Since: time.Duration(-300) * time.Second,
		Matches: []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
				Value: name,
			},
		},
	}
	r, err := sdjournal.NewJournalReader(config)
	if err != nil {
		return
	}
	if r == nil {
		err = fmt.Errorf("got a nil reader")
		return
	}
	defer r.Close()

	// and follow the reader synchronously
	if err = r.Follow(end, writer); err != sdjournal.ErrExpired {
		err = fmt.Errorf("error during follow: %s", err)
		return
	}
	return nil
}
