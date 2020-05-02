package cmd

import (
	"context"
	"flag"
	"github.com/google/subcommands"
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
	return "init"
}

func (*Init) SetFlags(f *flag.FlagSet) {

}

func (i *Init) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {

	return subcommands.ExitSuccess
}
