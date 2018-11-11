package tabnine

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

type Process struct {
	binPath   string
	binArgs   []string
	configDir string
}

// NewProcess creates a new BackgroundProcess instance.
func NewProcess(configDir, binPath string, binArgs []string) (Process, error) {
	return Process{
		binPath:   binPath,
		binArgs:   binArgs,
		configDir: configDir,
	}, nil
}

func (t Process) EnsureStarted() error {
	running, err := t.Running()
	if err != nil {
		return fmt.Errorf("running: %v", err)
	}

	if running {
		return nil
	}

	return t.start()
}

func (t Process) Running() (bool, error) {
	r, _, err := t.running()
	return r, err
}

func (t Process) running() (bool, int, error) {
	pid, err := t.readPID()
	if err != nil {
		return false, 0, fmt.Errorf("readpid: %v", err)
	}

	// this doesn't support a pid of actually zero,
	// not sure if that's possible. Assuming not,
	// call me an ass :)
	if pid == 0 {
		return false, 0, nil
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false, 0, fmt.Errorf("findprocess %d: %v", pid, err)
	}

	// not sure if this is a valid method to check for process
	// status, experimenting for now, we'll see if it sticks..
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return false, 0, nil
	}

	return true, pid, nil
}

func (t Process) Start() error {
	running, pid, err := t.running()
	if err != nil {
		return fmt.Errorf("running: %v", err)
	}
	if running {
		return fmt.Errorf("already running process: %d", pid)
	}
	return t.start()
}

func (t Process) start() error {
	cmd := exec.Command(t.binPath, t.binArgs...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd start: %v", err)
	}

	if err := t.writePID(cmd.Process.Pid); err != nil {
		return fmt.Errorf("writepid %d: %v", cmd.Process.Pid, err)
	}

	return nil
}

func (t Process) pidFilename() string {
	return filepath.Base(t.binPath) + ".pid"
}

func (t Process) readPID() (int, error) {
	pidPath := filepath.Join(t.configDir, t.pidFilename())

	b, err := ioutil.ReadFile(pidPath)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("readfile: %v", err)
	}

	pid, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, fmt.Errorf("strconv: %v", err)
	}

	return pid, nil
}

func (t Process) writePID(pid int) error {
	// rejecting 0 pids, since t.running() doesn't support 0 pid either,
	// this way we fail more quickly.
	if pid == 0 {
		return fmt.Errorf("0 pid not supported currently")
	}

	pidPath := filepath.Join(t.configDir, t.pidFilename())

	b := []byte(strconv.Itoa(pid))

	if err := ioutil.WriteFile(pidPath, b, 0644); err != nil {
		return fmt.Errorf("writefile: %v", err)
	}

	return nil
}
