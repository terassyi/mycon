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
	return "create"
}

func (*Create) Synopsis() string {
	return "create new container based on config.json"
}

// TODO make usage
func (*Create) Usage() string {
	return "create container"
}

func (c *Create) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.bundle, "bundle", "", "bundle directory")
}

func (c *Create) Execute(_ context.Context, flag *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	id := flag.Arg(0)
	f, err := factory.New(id, "")
	if err != nil {
		logrus.Debug(err)
		return subcommands.ExitFailure
	}

	config, err := container.NewConfig(id, c.bundle)
	if err != nil {
		logrus.Debug(err)
		return subcommands.ExitFailure
	}
	logrus.Debug(config.String())
	logrus.Debugf("new container create id=%v", id)
	_, err = f.Create(config)
	if err != nil {
		logrus.Debug(err)
		return subcommands.ExitFailure
	}
	logrus.Debugf("return success")
	return subcommands.ExitSuccess
}
