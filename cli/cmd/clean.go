// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"aliax/internal/cfg"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	cleanCmd = &cobra.Command{
		Use: "clean",
		Run: func(cmd *cobra.Command, args []string) {
			data, err := os.ReadFile(config)
			if err != nil {
				panic(err)
			}
			var file cfg.Aliax
			err = yaml.Unmarshal(data, &file)
			if err != nil {
				panic(err)
			}
			bins := map[string]struct{}{}
			for name := range file.Extend {
				bins[name] = struct{}{}
			}
			err = filepath.Walk("run-scripts", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.HasSuffix(filepath.Base(path), ".ps1") {
					os.Remove(path)
				}

				if strings.HasSuffix(filepath.Base(path), ".sh") {
					os.Remove(path)
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	aliaxCmd.AddCommand(cleanCmd)
}
