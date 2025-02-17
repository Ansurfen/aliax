// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"aliax/internal/cfg"
	"os"

	"github.com/spf13/cobra"
)

// Version will be set by the build process using -ldflags.
// Do not modify this variable directly. It is populated at build time with the desired version string.
var Version = ""

var (
	config   = cfg.Name()
	aliaxCmd = &cobra.Command{
		Use:   "aliax",
		Short: "A CLI tool for managing and extending commands",
		Long: `Aliax is a command-line tool designed to enhance workflow efficiency by:
- Extending existing commands with additional functionality.
- Managing command aliases within a workspace.
- Creating new custom commands to streamline repetitive tasks.
`,
	}
)

func Execute() {
	err := aliaxCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Root() *cobra.Command {
	return aliaxCmd
}
