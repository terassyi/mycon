package spec

import (
	"encoding/json"
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"path/filepath"
)

const (
	configPath string = "config.json"
)

func LoadSpec(path string) (*specs.Spec, error) {
	filePath := filepath.Join(path, configPath)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config.json: %v", err)
	}
	var spec *specs.Spec
	if err := json.NewDecoder(f).Decode(&spec); err != nil {
		return nil, fmt.Errorf("failed to decode config.json: %v", err)
	}
	return spec, nil
}
