package capabilities

import (
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/syndtr/gocapability/capability"
	"strings"
)

const allCapabilityTypes = capability.CAPS | capability.BOUNDS | capability.AMBS

type Capabilities struct {
	CapMap      map[string]capability.Cap
	Pid         capability.Capabilities
	Bounding    []capability.Cap
	Inheritable []capability.Cap
	Effective   []capability.Cap
	Permitted   []capability.Cap
	Ambient     []capability.Cap
}

func New(caps *specs.LinuxCapabilities) (*Capabilities, error) {
	// init capabilities
	capMap := make(map[string]capability.Cap)
	for _, c := range capability.List() {
		key := fmt.Sprintf("CAP_%s", strings.ToUpper(c.String()))
		capMap[key] = c
	}

	bounding := []capability.Cap{}
	for _, c := range caps.Bounding {
		bounding = append(bounding, capMap[c])
	}
	effective := []capability.Cap{}
	for _, c := range caps.Effective {
		effective = append(effective, capMap[c])
	}
	inheritable := []capability.Cap{}
	for _, c := range caps.Inheritable {
		inheritable = append(inheritable, capMap[c])
	}
	permitted := []capability.Cap{}
	for _, c := range caps.Permitted {
		permitted = append(permitted, capMap[c])
	}
	ambient := []capability.Cap{}
	for _, c := range caps.Ambient {
		ambient = append(ambient, capMap[c])
	}
	pid, err := capability.NewPid2(0)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize capabilities: %v", err)
	}
	return &Capabilities{
		CapMap:      capMap,
		Pid:         pid,
		Bounding:    bounding,
		Inheritable: inheritable,
		Effective:   effective,
		Permitted:   permitted,
		Ambient:     ambient,
	}, nil
}

func (caps *Capabilities) ApplyCaps() error {
	caps.Pid.Clear(allCapabilityTypes)
	caps.Pid.Set(capability.BOUNDING, caps.Bounding...)
	caps.Pid.Set(capability.PERMITTED, caps.Permitted...)
	caps.Pid.Set(capability.EFFECTIVE, caps.Effective...)
	caps.Pid.Set(capability.INHERITABLE, caps.Inheritable...)
	caps.Pid.Set(capability.AMBIENT, caps.Ambient...)
	return caps.Pid.Apply(allCapabilityTypes)
}
