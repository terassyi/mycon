package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
	"github.com/terassyi/mycon/cmd"
	"os"
	"syscall"
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
	// Call the subcommand and pass in the configuration.
	var ws syscall.WaitStatus
	subcmdCode := subcommands.Execute(context.Background(), &ws)
	if subcmdCode == subcommands.ExitSuccess {
		logrus.Debugf("Exiting with status: %v", ws)
		if ws.Signaled() {
			// No good way to return it, emulate what the shell does. Maybe raise
			// signall to self?
			os.Exit(128 + int(ws.Signal()))
		}
		os.Exit(ws.ExitStatus())
	}
	// Return an error that is unlikely to be used by the application.
	logrus.Warningf("Failure to execute command, err: %v", subcmdCode)
	os.Exit(128)
}

func setDebugMode(debug bool) {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		fmt.Println("! debug mode !")
	}
}
