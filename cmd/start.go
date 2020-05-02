package cmd

import (
	"context"
	"flag"
	"github.com/google/subcommands"
)

type Start struct {
}

func (*Start) Name() string {
	return "start [container id]"
}

func (*Start) Synopsis() string {
	return "start container"
}

func (*Start) Usage() string {
	return "start container"
}

func (*Start) SetFlags(f *flag.FlagSet) {

}

func (s *Start) Execute(_ context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {

	return subcommands.ExitSuccess
}
