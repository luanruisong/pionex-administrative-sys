package daemon

import (
	"os"
	"os/exec"
	"syscall"
)

type daemon struct {
	Pid  int
	Args []string
}

func Daemon() (*daemon, error) {
	args := make([]string, 0)
	cmdStr := os.Args[0]
	hasFL := false
	for i := 1; i < len(os.Args); i++ {
		v := os.Args[i]
		if v != "-d" {
			args = append(args, v)
			if v == "-fl" {
				hasFL = true
			}
		}
	}
	if !hasFL {
		args = append(args, "-fl")
	}

	cmd := exec.Command(cmdStr, args...)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	// daemon 进程不绑定父进程的 stdin/stdout/stderr
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &daemon{
		Pid:  cmd.Process.Pid,
		Args: args,
	}, nil
}
