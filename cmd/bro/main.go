package main

import (
	"os"
)

func main() {
	s := &state{
		Args: os.Args,

		Out: os.Stdout,
		Err: os.Stderr,

		ToolVersion:    ToolVersion,
		GitSHA:         GitSHA,
		BuiltAtRFC3339: BuiltAtRFC3339,
	}

	spinMain(s)
}

// spinHelp explains the Spin tool.
const spinHelp = `Spin is a tool for interacting with a Spin system.

Usage: 

    spin <command> [arguments]

The commands are:

    dir        interact with dir servers
    help       more info about a command
    keys       manage keys used by the tool
    run        interact with spin compute
    store      interact with spin storage
    tool       run specified spin tool
    version    print Spin tool version

Use "spin help <command>" for more information about a command.`

// TODO: add run, snap

func spinMain(s *state) {
	if len(s.Args) < 2 {
		s.println(spinHelp)
		return
	}

	switch cmd := s.Args[1]; cmd {
	case "run":
		runMain(s)
	case "help":
		spinHelpMain(s)
	case "version":
		spinVersionMain(s)
	default:
		s.println("spin %s: unknown command\nRun \"spin help\" for usage.", cmd)
	}
}
