package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"github.com/terassyi/mycon/cmd"
	"os"
)

var (
	debug bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "debug mode")
}

func main() {
	subcommands.Register(subcommands.FlagsCommand(), "")
	//subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(new(cmd.Create), "")
	subcommands.Register(new(cmd.Start), "")

	const internalOnly = "internal only"
	subcommands.Register(new(cmd.Init), internalOnly)

	flag.Parse()
	setDebugMode(debug)

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

func setDebugMode(debug bool) {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		fmt.Println("! debug mode !")
	}
}
