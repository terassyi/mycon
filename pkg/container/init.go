package container

import (
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
)

const (
	rootPath = "/run/mycon/"
	fifoName = "fifo.exe"
)

type Initializer struct {
	Id     string
	FifoFd int
	Spec   *specs.Spec
}

func NewInitializer(spec *specs.Spec, fd int, id string) (*Initializer, error) {
	return &Initializer{
		Id:     id,
		FifoFd: fd,
		Spec:   spec,
	}, nil
}

func (i *Initializer) Init() error {
	if err := i.prepareRootfs(); err != nil {
		logrus.Debug(err)
		return err
	}
	name, err := exec.LookPath(i.Spec.Process.Args[0])
	if err != nil {
		logrus.Debugf("cannot find container exec command: %v", err)
		return err
	}

	//if err := <- i.waitToStart(); err != nil {
	//	logrus.Debug(err)
	//	return err
	//}
	logrus.Debugf("execute container start commands")
	if err := unix.Exec(name, i.Spec.Process.Args[0:], os.Environ()); err != nil {
		logrus.Debug(err)
		return err
	}
	return nil
}

func (i *Initializer) waitToStart() <-chan error {
	errCh := make(chan error)
	go func() {
		if err := i.writeToFifo(); err != nil {
			errCh <- err
			return
		}
	}()
	return errCh
}

func (i *Initializer) writeToFifo() error {
	path := fmt.Sprintf("/proc/self/fd/%d", i.FifoFd)
	file, err := unix.Open(path, unix.O_WRONLY, 0700)
	//path := filepath.Join(rootPath, i.Id, fifoName)
	logrus.Debugf("fifo path: %v", path)
	//file, err := unix.Open(path, unix.O_WRONLY, 0700)
	if err != nil {
		logrus.Debugf("failed to open fifo file from init process: %v", err)
		return err
	}
	if _, err := unix.Write(file, []byte("0")); err != nil {
		return err
	}
	defer unix.Close(file)
	return nil
}

func (i *Initializer) prepareRootfs() error {
	if err := i.prepareRoot(); err != nil {
		logrus.Debugf("failed to mount: %v", err)
		return err
	}
	logrus.Debugf("mount finished.")
	//if err := unix.Chdir(i.Spec.Root.Path); err != nil {
	//	logrus.Debugf("failed to chdir to rootfs: %v", err)
	//	return err
	//}
	//logrus.Debugf("chdir: %v", i.Spec.Root.Path)
	// pivot_root
	if err := i.pivotRoot(); err != nil {
		logrus.Debugf("failed to pivot_root: %v", err)
		return err
	}
	return nil
}

func (i *Initializer) prepareRoot() error {
	// mount
	if err := unix.Mount("", "/", "", unix.MS_SLAVE|unix.MS_REC, ""); err != nil {
		return err
	}
	return unix.Mount(i.Spec.Root.Path, i.Spec.Root.Path, "bind", unix.MS_BIND|unix.MS_REC, "")
}

func (i *Initializer) pivotRoot() error {
	oldroot, err := unix.Open("/", unix.O_DIRECTORY|unix.O_RDONLY, 0)
	if err != nil {
		logrus.Debugf("failed to open old root")
		return err
	}
	defer unix.Close(oldroot)
	newroot, err := unix.Open(i.Spec.Root.Path, unix.O_DIRECTORY|unix.O_RDONLY, 0)
	if err != nil {
		logrus.Debug("failed to open new root: ", i.Spec.Root.Path)
		cd, _ := os.Getwd()
		logrus.Debug("now in ", cd)
		return err
	}
	defer unix.Close(newroot)
	// fetch new root file system
	if err := unix.Fchdir(newroot); err != nil {
		logrus.Debug("failed to fetch new root")
		return err
	}
	if err := unix.PivotRoot(".", "."); err != nil {
		logrus.Debugf("failed to pivot_root: %v", err)
		return err
	}
	if err := unix.Fchdir(oldroot); err != nil {
		logrus.Debug("failed to fetch old root")
		return err
	}
	if err := unix.Mount("", ".", "", unix.MS_SLAVE|unix.MS_REC, ""); err != nil {
		logrus.Debug("failed to mount .")
		return err
	}
	if err := unix.Unmount(".", unix.MNT_DETACH); err != nil {
		logrus.Debug("failed to unmount .")
		return err
	}
	if err := unix.Chdir("/"); err != nil {
		logrus.Debug("failed to chdir /")
		return fmt.Errorf("failed to chdir: %v", err)
	}
	return nil
}
