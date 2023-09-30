// Copyright 2023 The Spin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

// spinVersionHelp explains the version command.
// It is used by the "spin help" command.
const spinVersionHelp = `Print Spin tool version information.

Usage: 

    spin version

Output is of the form:

    spin version spin<version> <gitsha>/<builtat>

where:
   <version>  tool version
   <gitsha>   tool build commit sha (shortened) if clean, else XXXXXXX
   <builtat>  tool build time in RFC3339 format`

// spinVersionMain prints information about the tool's version and build.
func spinVersionMain(s *state) {
	if len(s.GitSHA) < 7 {
		panic(fmt.Sprintf("state.GitSHA is %q, less than 7 ASCII characters", s.GitSHA)) // untested
	}
	s.println("spin version spin%s %s/%s", s.ToolVersion, s.GitSHA[:7], s.BuiltAtRFC3339)
}
