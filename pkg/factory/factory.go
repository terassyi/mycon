package factory

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/terassyi/mycon/pkg/container"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	rootPath = "/run/mycon"
	fifoName = "fifo.exe"
)

// Factory
type Factory struct {
	//Id string
	Pid      int
	Root     string
	InitPath string
	InitArgs []string
}

// New returns new Factory to make a new container initialization process
func New(root string) (*Factory, error) {
	if root != "" {
		if err := os.MkdirAll(root, 0700); err != nil {
			return nil, err
		}
	} else {
		root = rootPath
	}
	factory := &Factory{
		Pid:      -1,
		Root:     root,
		InitPath: "/proc/self/exe",
		InitArgs: []string{os.Args[0], "init"}, // path to mycon init
	}
	return factory, nil
}

// Create creates a new container structure
func (f *Factory) Create(id string, config *container.Config) (*container.Container, error) {
	// check container directory already exists
	containerRootPath := filepath.Join(f.Root, id)
	if _, err := os.Stat(containerRootPath); err == nil {
		return nil, fmt.Errorf("container root dir is already exist")
	}
	// make container dir
	if err := os.MkdirAll(containerRootPath, 0711); err != nil {
		return nil, err
	}
	if err := os.Chown(containerRootPath, unix.Getuid(), unix.Getgid()); err != nil {
		return nil, err
	}

	return &container.Container{
		Id:     id,
		Root:   containerRootPath,
		Config: config,
	}, nil
}

func (f *Factory) create() error {
	if err := f.createFifo(); err != nil {
		return err
	}
	cmd := f.buildInitCommand()
	fd, err := f.setFifoFd(cmd)
	if err != nil {
		return err
	}
	logrus.Debugf("mycon create fifo fd->%v", fd)
	pid, err := f.exec(cmd)
	if err != nil {
		return err
	}
	f.Pid = pid
	logrus.Debugf("container init process is called [pid=%d]", pid)
	return nil
}

// exec starts init process and return pid and error
func (f *Factory) exec(cmd *exec.Cmd) (int, error) {
	if err := cmd.Start(); err != nil {
		return -1, err
	}
	pid := cmd.Process.Pid
	return pid, nil
}

// buildInitCommand builds a command to start init process
func (f *Factory) buildInitCommand() *exec.Cmd {
	cmd := exec.Command(f.InitPath, f.InitArgs...)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWIPC | unix.CLONE_NEWNET | unix.CLONE_NEWNS |
			unix.CLONE_NEWPID | unix.CLONE_NEWUSER | unix.CLONE_NEWUTS,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func (f *Factory) createFifo() error {
	path := filepath.Join(f.Root, fifoName)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("fifo.exe already exists")
	}
	if err := unix.Mkfifo(path, 0744); err != nil {
		return fmt.Errorf("failed to create fifo file: %v", err)
	}
	return nil
}

func (f *Factory) deleteFifo() error {
	path := filepath.Join(f.Root, fifoName)
	return os.RemoveAll(path)
}

func (f *Factory) setFifoFd(cmd *exec.Cmd) (int, error) {
	path := filepath.Join(f.Root, fifoName)
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, err
	}
	cmd.ExtraFiles = append(cmd.ExtraFiles, os.NewFile(uintptr(fd), fifoName))
	cmd.Env = append(cmd.Env, fmt.Sprintf("_MYCON_FIFOFD=%v", fd))
	return fd, err
}
