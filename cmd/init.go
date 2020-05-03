package cmd

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"github.com/terassyi/mycon/pkg/factory"
)

type Init struct {
}

func (*Init) Name() string {
	return "init"
}

func (*Init) Synopsis() string {
	return "initialize container resources"
}

func (*Init) Usage() string {
	return "init [container id]"
}

func (*Init) SetFlags(f *flag.FlagSet) {

}

func (i *Init) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	logrus.Debugf("init process starts")
	fac, err := factory.New("", "")
	if err != nil {
		return subcommands.ExitFailure
	}
	id := f.Arg(0)
	fac.Id = id
	logrus.Debugf("process initialize")
	if err := fac.Initialize(); err != nil {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
