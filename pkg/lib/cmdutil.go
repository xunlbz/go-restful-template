package lib

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/xunlbz/go-restful-template/pkg/log"
)

var isContainer bool

func init() {
	isContainer = checkContainer()
}

func checkContainer() bool {
	f, err := os.Open("/.dockerenv")
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}

func ExecCommand(name string, arg ...string) ([]byte, error) {
	var cmd *exec.Cmd
	timeout := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if len(arg) == 0 {
		if isContainer {
			cmd = exec.CommandContext(ctx, "nsenter", "--mount=/host/1/ns/mnt", "--ipc=/host/1/ns/ipc", "--net=/host/1/ns/net", name)
		} else {
			cmd = exec.CommandContext(ctx, name)
		}
	} else {
		if isContainer {
			args := make([]string, 0)
			args = append(args, "--mount=/host/1/ns/mnt")
			args = append(args, "--ipc=/host/1/ns/ipc")
			args = append(args, "--net=/host/1/ns/net")
			args = append(args, "--") //nsenter命令不解析后边参数 如 sh -c "ls /srv"  -c为 sh 参数 不由nsenter解析否则不识别
			args = append(args, name)
			args = append(args, arg...)
			cmd = exec.CommandContext(ctx, "nsenter", args...)
		} else {
			cmd = exec.CommandContext(ctx, name, arg...)
		}
	}
	log.Debugf("run cmd: %s", cmd.String())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	res, err := cmd.Output()
	if err != nil {
		log.Errorf(err.Error(), stderr.String())
	}
	log.Debugf("run cmd result: %s", string(res))
	return res, err
}

//  some command will not stop, so this method must run in goroutine, unless you want to wait for result
func StartCommand(name string, arg ...string) error {
	var cmd *exec.Cmd
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if len(arg) == 0 {
		cmd = exec.CommandContext(ctx, name)
	} else {
		cmd = exec.CommandContext(ctx, name, arg...)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Errorf(err.Error(), stderr.String())

	}
	log.Infof("start cmd: %s", cmd.String())
	if err = cmd.Wait(); err != nil {
		log.Error(err.Error())
	} else {
		log.Infof("cmd exit pid is %v,", cmd.ProcessState.Pid())
	}
	return err
}

func StopCommand(name string) error {
	var cmd *exec.Cmd
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bin := "killall"
	args := []string{"-q", name}
	cmd = exec.CommandContext(ctx, bin, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	log.Infof("stop cmd: %s", cmd.String())
	if err != nil {
		log.Error(err.Error())
		return err
	} else {
		log.Infof("cmd exit")
	}
	return nil
}
