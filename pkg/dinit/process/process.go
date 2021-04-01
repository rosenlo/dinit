package process

import (
	"os"
	"os/exec"
	"syscall"
)

func NewCommand(path string, arg ...string) *exec.Cmd {
	cmd := exec.Command(path, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}
