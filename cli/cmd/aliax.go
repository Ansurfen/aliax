// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const config = "aliax.yaml"

var aliaxCmd = &cobra.Command{
	Use: "aliax",
}

func Execute() {
	err := aliaxCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Root() *cobra.Command {
	return aliaxCmd
}
