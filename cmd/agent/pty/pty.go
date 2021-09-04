//go:build !windows
//+build !windows

package pty

import (
	"errors"
	"os"
	"os/exec"

	opty "github.com/creack/pty"
)

var defaultShells = []string{"zsh", "fish", "bash", "sh"}

type Pty struct {
	tty *os.File
	cmd *exec.Cmd
}

func DownloadDependency() {
}

func Start() (*Pty, error) {
	var shellPath string
	for i := 0; i < len(defaultShells); i++ {
		shellPath, _ = exec.LookPath(defaultShells[i])
		if shellPath != "" {
			break
		}
	}
	if shellPath == "" {
		return nil, errors.New("没有可用终端")
	}
	cmd := exec.Command(shellPath) // #nosec
	cmd.Env = append(os.Environ(), "TERM=xterm")
	tty, err := opty.Start(cmd)
	return &Pty{tty: tty, cmd: cmd}, err
}

func (pty *Pty) Write(p []byte) (n int, err error) {
	return pty.tty.Write(p)
}

func (pty *Pty) Read(p []byte) (n int, err error) {
	return pty.tty.Read(p)
}

func (pty *Pty) Setsize(cols, rows uint32) error {
	return opty.Setsize(pty.tty, &opty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	})
}

func (pty *Pty) Close() error {
	if err := pty.tty.Close(); err != nil {
		return err
	}
	return pty.cmd.Process.Kill()
}
