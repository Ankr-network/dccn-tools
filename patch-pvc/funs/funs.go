package funs

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

var timeout = 10 * time.Second

// Exec 执行 command
func Exec(name string, arg ...string) (string, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctxt, name, arg...)
	//当经过Timeout时间后，程序依然没有运行完，则会杀掉进程，ctxt也会有err信息
	if out, err := cmd.Output(); err != nil {
		//检测报错是否是因为超时引起的
		if ctxt.Err() != nil && ctxt.Err() == context.DeadlineExceeded {
			return "", errors.New("command timeout")
		}
		return string(out), err
	} else {
		return string(out), nil
	}
}
