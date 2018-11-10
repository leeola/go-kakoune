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

const (
	pidFilename    = "tabnine.pid"
	stdinFilename  = "PIPE_STDIN"
	stdoutFilename = "PIPE_STDOUT"
)

type Process struct {
	binPath     string
	workingPath string
}

func NewProcess(c Config) (*Process, error) {
	return &Process{
		binPath:     c.TabnineBin,
		workingPath: c.ConfigDir,
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
	stdin, err := t.openStdinPipe()
	if err != nil {
		return fmt.Errorf("openstdinpipe: %v", err)
	}

	stdout, err := t.openStdoutPipe()
	if err != nil {
		return fmt.Errorf("openstdoutpipe: %v", err)
	}

	cmd := exec.Command(t.binPath)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	// should the named pipe files be attached? It seems not required,
	// though we could be leaking file descriptors here?
	cmd.ExtraFiles = []*os.File{stdin, stdout}
	cmd.Start()

	if err := t.writePID(cmd.Process.Pid); err != nil {
		return fmt.Errorf("writepid %d: %v", cmd.Process.Pid, err)
	}

	return nil
}

func (t Process) readPID() (int, error) {
	pidPath := filepath.Join(t.workingPath, pidFilename)

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

	pidPath := filepath.Join(t.workingPath, pidFilename)

	b := []byte(strconv.Itoa(pid))

	if err := ioutil.WriteFile(pidPath, b, 0644); err != nil {
		return fmt.Errorf("writefile: %v", err)
	}

	return nil
}

func (t Process) openStdinPipe() (*os.File, error) {
	p := filepath.Join(t.workingPath, stdinFilename)
	flag := os.O_RDWR

	// TODO(leeola): is there a way to open a read or write only pipe
	// *(for stdin and stdout respectively)*, but still detatched from this
	// process? The goal is not to have this process block, but do block the
	// TabNine process waiting for io over the named pipes.
	stdin, err := os.OpenFile(p, flag, os.ModeNamedPipe)
	if os.IsNotExist(err) {
		if err := syscall.Mkfifo(p, 0644); err != nil {
			return nil, fmt.Errorf("mkfifo: %v", err)
		}

		stdin, err = os.OpenFile(p, flag, os.ModeNamedPipe)
	}
	if err != nil {
		return nil, fmt.Errorf("openfile: %v", err)
	}

	return stdin, nil
}

func (t Process) openStdoutPipe() (*os.File, error) {
	p := filepath.Join(t.workingPath, stdoutFilename)
	flag := os.O_RDWR

	// TODO(leeola): is there a way to open a read or write only pipe
	// *(for stdin and stdout respectively)*, but still detatched from this
	// process? The goal is not to have this process block, but do block the
	// TabNine process waiting for io over the named pipes.
	stdout, err := os.OpenFile(p, flag, os.ModeNamedPipe)
	if os.IsNotExist(err) {
		if err := syscall.Mkfifo(p, 0644); err != nil {
			return nil, fmt.Errorf("mkfifo: %v", err)
		}

		stdout, err = os.OpenFile(p, flag, os.ModeNamedPipe)
	}
	if err != nil {
		return nil, fmt.Errorf("openfile: %v", err)
	}

	return stdout, nil
}

func openStdinClient(workingPath string) (*os.File, error) {
	p := filepath.Join(workingPath, stdinFilename)
	stdin, err := os.OpenFile(p, os.O_WRONLY|syscall.O_NONBLOCK, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("openfile: %v", err)
	}
	return stdin, nil
}

func openStdoutClient(workingPath string) (*os.File, error) {
	p := filepath.Join(workingPath, stdoutFilename)
	stdout, err := os.OpenFile(p, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("openfile: %v", err)
	}
	return stdout, nil
}
