package factory

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/terassyi/mycon/pkg/container"
	"github.com/terassyi/mycon/pkg/spec"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	rootPath = "/run/mycon"
	fifoName = "fifo.exe"
)

// Factory
type Factory struct {
	Id       string
	Pid      int
	Root     string
	InitPath string
	InitArgs []string
}

// New returns new Factory to make a new container initialization process
func New(id string, root string) (*Factory, error) {
	if root != "" {
		if err := os.MkdirAll(root, 0700); err != nil {
			return nil, err
		}
	} else {
		root = rootPath
	}
	factory := &Factory{
		Id:       id,
		Pid:      -1,
		Root:     root,
		InitPath: "/proc/self/exe",
		InitArgs: []string{os.Args[0], "-debug", "init", id}, // path to mycon init
		//InitArgs: []string{os.Args[0], "init"},
	}
	return factory, nil
}

// Create creates a new container structure
func (f *Factory) Create(config *container.Config) (*container.Container, error) {
	// check container directory already exists
	containerRootPath := filepath.Join(f.Root, f.Id)
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

	// create the init process
	logrus.Debugf("create the init process")
	logrus.Debug("move to bundle dir: %v", config.Bundle)
	if err := os.Chdir(config.Bundle); err != nil {
		logrus.Debug("failed to chdir bundle dir: %v", err)
		return nil, err
	}
	dir, err := os.Getwd()
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}
	logrus.Debugf("working dir: %v", dir)
	if err := f.create(); err != nil {
		return nil, err
	}

	return &container.Container{
		Id:     f.Id,
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
		logrus.Debug(err)
		return err
	}
	logrus.Debugf("mycon create fifo fd->%v", fd)
	pid, err := f.exec(cmd)
	if err != nil {
		return err
	}
	f.Pid = pid
	logrus.Debugf("container init process is called [pid=%d]", pid)
	if _, err := os.Stat(fmt.Sprintf("/proc/self/fd/%d", fd)); err != nil {
		logrus.Debugf("failed to open fifo file from create process: %v", err)
		return err
	}
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
	cmd := exec.Command(f.InitPath, f.InitArgs[1:]...)
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
	logrus.Debugf(cmd.String())
	return cmd
}

func (f *Factory) createFifo() error {
	path := filepath.Join(f.Root, f.Id, fifoName)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("fifo.exe already exists")
	}
	if err := unix.Mkfifo(path, 0744); err != nil {
		return fmt.Errorf("failed to create fifo file: %v", err)
	}
	logrus.Debugf("create fifo file to %v", path)
	return nil
}

func (f *Factory) deleteFifo() error {
	path := filepath.Join(f.Root, f.Id, fifoName)
	return os.RemoveAll(path)
}

func (f *Factory) setFifoFd(cmd *exec.Cmd) (int, error) {
	path := filepath.Join(f.Root, f.Id, fifoName)
	fd, err := unix.Open(path, unix.O_PATH|unix.O_CLOEXEC, 0)
	if err != nil {
		logrus.Debug(err)
		return -1, err
	}
	defer unix.Close(fd)
	cmd.ExtraFiles = append(cmd.ExtraFiles, os.NewFile(uintptr(fd), fifoName))
	cmd.Env = append(cmd.Env, fmt.Sprintf("_MYCON_FIFOFD=%v", fd+3+len(cmd.ExtraFiles)-1))
	return fd, err
}

func (f *Factory) Initialize() error {
	// get fifo fd
	fd := os.Getenv("_MYCON_FIFOFD")
	if fd == "" {
		return fmt.Errorf("fd to fifo.exe [_MYCON_FIFOFD] is not set.")
	}
	fifoFd, err := strconv.Atoi(fd)
	if err != nil {
		return err
	}
	logrus.Debugf("container fifo fd = %v", fd)
	s, err := spec.LoadSpec(".")
	if err != nil {
		logrus.Debug(err)
		return err
	}
	initer, err := container.NewInitializer(s, fifoFd, f.Id)
	if err != nil {
		logrus.Debug(err)
		return err
	}
	return initer.Init()
}
