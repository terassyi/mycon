package cgroups

import (
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const (
	//cgroupMemorySwapLimit = "memory.memsw.limit_in_bytes"
	//cgroupMemoryLimit     = "memory.limit_in_bytes"
	cgroupKernelMemoryLimit = "memory.kmem.limit_in_bytes"
	cpuShares               = "cpu.shares"
	cpuCfsPeriodUs          = "cpu.cfs_period_us"
	cpuCfsQuotaUs           = "cpu.cfs_quota_us"
	cpuRtPeriodUs           = "cpu.rt_period_us"
	cpuRtRuntimeUs          = "cpu.rt_runtime_us"
)

type Cgroups struct {
	Root      string
	Resources *specs.LinuxResources
	Pid       int
}

func New(resources *specs.LinuxResources) (*Cgroups, error) {
	if _, err := os.Stat("/sys/fs/cgroup"); err != nil {
		return nil, err
	}
	return &Cgroups{
		Root:      "/sys/fs/cgroup",
		Resources: resources,
		Pid:       os.Getpid(),
	}, nil
}

func (cg *Cgroups) Limit() error {
	logrus.Debugf("uid=%v gid=%v pid=%v", os.Getuid(), os.Getgid(), os.Getpid())
	if err := cg.limitCpu(); err != nil {
		logrus.Debugf("failed to limit cpu: %v", err)
		return err
	}
	if err := cg.limitKernelMemory(); err != nil {
		logrus.Debugf("failed to limit kernel memory: %v", err)
		return err
	}
	return nil
}

func (cg *Cgroups) limitCpu() error {
	if cg.Resources == nil || cg.Resources.CPU == nil {
		logrus.Debugf("cpu limitation is not set")
		return nil
	}
	dir := filepath.Join(cg.Root, "cpu", "mycon")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	if cg.Resources.CPU.Shares != nil {
		if err := writeFile(dir, cpuShares, strconv.FormatUint(*cg.Resources.CPU.Shares, 10)); err != nil {
			return err
		}
	}
	if cg.Resources.CPU.Quota != nil {
		if err := writeFile(dir, cpuCfsQuotaUs, strconv.FormatInt(*cg.Resources.CPU.Quota, 10)); err != nil {
			return err
		}
	}
	if cg.Resources.CPU.Period != nil {
		if err := writeFile(dir, cpuCfsPeriodUs, strconv.FormatUint(*cg.Resources.CPU.Period, 10)); err != nil {
			return err
		}
	}
	if cg.Resources.CPU.RealtimePeriod != nil {
		if err := writeFile(dir, cpuRtPeriodUs, strconv.FormatUint(*cg.Resources.CPU.RealtimePeriod, 10)); err != nil {
			return err
		}
	}
	if cg.Resources.CPU.RealtimeRuntime != nil {
		if err := writeFile(dir, cpuRtRuntimeUs, strconv.FormatInt(*cg.Resources.CPU.RealtimeRuntime, 10)); err != nil {
			return err
		}
	}
	return nil
}

func (cg *Cgroups) limitKernelMemory() error {
	if cg.Resources == nil || cg.Resources.Memory == nil {
		logrus.Debugf("kernel memory limitation is not set")
		return nil
	}
	dir := filepath.Join(cg.Root, "memory", "mycon")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	if cg.Resources.Memory.Kernel != nil {
		if err := writeFile(dir, cgroupKernelMemoryLimit, strconv.FormatInt(*cg.Resources.Memory.Kernel, 10)); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(dir, file, data string) error {
	if dir == "" {
		return fmt.Errorf("no such directory for %s", file)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, file), []byte(data), 0700); err != nil {
		return fmt.Errorf("failed to write %s to %s: %v", data, file, err)
	}
	return nil
}
