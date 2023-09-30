// Copyright 2023 The Spin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// spinHelpHelp explains the spin help command.
const spinHelpHelp = `Get more help about a topic or Spin command.

Usage:

    spin help <topic>

The arguments are:

    <topic>    topic for help

The topic may be a Spin command. To list these run "spin"`

// spinHelpMain handles printing information about help topics.
func spinHelpMain(s *state) {
	if len(s.Args) < 3 {
		s.println(spinHelpHelp)
		return
	}

	switch v := s.Args[2]; v {
	case "help":
		s.println(spinHelpHelp)
	case "keys":
		s.println(keysHelp)
	case "store":
		s.println(storeHelp)
	case "tool":
		s.println(toolHelp)
	case "version":
		s.println(spinVersionHelp)
	default:
		s.println("no detailed help for %q", v)
	}
}
