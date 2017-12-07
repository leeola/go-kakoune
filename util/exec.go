package util

import (
	"bytes"
	"os/exec"
	"syscall"
)

func Exec(bin string, args ...string) (stdout, stderr string, exit int, err error) {
	cmd := exec.Command(bin, args...)

	var stdoutB bytes.Buffer
	cmd.Stdout = &stdoutB

	var stderrB bytes.Buffer
	cmd.Stderr = &stderrB

	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exit = status.ExitStatus()
			} else {
				return "", "", 0, err
			}
		} else {
			return "", "", 0, err
		}
	}

	return stdoutB.String(), stderrB.String(), exit, nil
}
