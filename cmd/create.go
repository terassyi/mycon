package cmd

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"github.com/terassyi/mycon/pkg/container"
	"github.com/terassyi/mycon/pkg/factory"
)

// Create implement google/subcommands.Command interface
type Create struct {
	bundle string
}

func (*Create) Name() string {
	return "create [container id]"
}

func (*Create) Synopsis() string {
	return "create new container based on config.json"
}

// TODO make usage
func (*Create) Usage() string {
	return "create container"
}

func (c *Create) SetFlags(f *flag.FlagSet) {
	flag.StringVar(&c.bundle, "bundle", "", "bundle directory")
}

func (c *Create) Execute(_ context.Context, flag *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	logrus.Debugf("new container create")
	f, err := factory.New("")
	if err != nil {
		return subcommands.ExitFailure
	}
	id := flag.Arg(1)
	config, err := container.NewConfig(id, c.bundle)
	if err != nil {
		return subcommands.ExitFailure
	}
	_, err = f.Create(id, config)
	if err != nil {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
