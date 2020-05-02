package container

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/terassyi/mycon/pkg/spec"
)

type Config struct {
	Id     string
	Bundle string
	Spec   *specs.Spec
}

func NewConfig(id string, bundle string) (*Config, error) {
	s, err := spec.LoadSpec(bundle)
	if err != nil {
		return nil, err
	}
	return &Config{
		Id:     id,
		Bundle: bundle,
		Spec:   s,
	}, nil
}
