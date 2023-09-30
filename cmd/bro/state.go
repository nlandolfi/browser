package main

import (
	"bufio"
	"fmt"
	"io"
)

// Set these using link flags; e.g., -X main.ToolVersion=...
var (
	// ToolVersion is the version of the Spin tool.
	ToolVersion string // e.g. 0.1.0

	// GitSHA is the SHA of the git commit, if the repo is clean.
	// If the repo is dirty, it will be all X characters.
	GitSHA string

	// BuiltAtRFC3339 is the build time of the Spin tool in RFC3339 format.
	BuiltAtRFC3339 string
)

// state is a helper type for items needed to process commands.
// It has convenient helper methods and aids with testing.
type state struct {
	Args []string  // equivalent of os.Args
	Out  io.Writer // equivalent of os.Stdout
	Err  io.Writer // equivalent of os.Stderr

	// For descriptions of these values, see above.
	ToolVersion    string
	GitSHA         string
	BuiltAtRFC3339 string

	In                *bufio.Reader
	GetPager          func() string
	GetTerminalHeight func() (int, error)
}

// info prints to the Out writer.
func (s *state) print(f string, vs ...interface{}) {
	fmt.Fprintf(s.Out, f, vs...)
}

// infoln prints to the Out writer, with a new line.
func (s *state) println(f string, vs ...interface{}) {
	fmt.Fprintf(s.Out, f+"\n", vs...)
}

// error prints to the Err writer.
func (s *state) error(f string, vs ...interface{}) {
	fmt.Fprintf(s.Err, f, vs...)
}

// errorln prints to the Err writer, with a newline.
func (s *state) errorln(f string, vs ...interface{}) {
	fmt.Fprintf(s.Err, f+"\n", vs...)
}
