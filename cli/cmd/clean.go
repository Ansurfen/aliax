// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"aliax/internal/cfg"
	"aliax/internal/io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

var (
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Remove generated scripts and clean up workspace",
		Long: `The "clean" command removes all auto-generated scripts from the workspace.
It scans the "run-scripts" directory and deletes scripts (e.g. .ps1, .sh).
Additionally, it ensures that outdated extended commands are cleared.`,
		Example: "  aliax clean",
		Run: func(cmd *cobra.Command, args []string) {
			var file cfg.Aliax
			err := io.ReadYAML(config, &file)
			if err != nil {
				log.WithError(err).Fatalf("parsing %s", config)
			}
			log.Infof("parsing %s", config)
			bins := map[string]struct{}{}
			for name := range file.Extend {
				bins[name] = struct{}{}
			}
			log.Info("cleaning")
			log.IncreasePadding()
			err = filepath.Walk("run-scripts", func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.HasSuffix(filepath.Base(path), ".ps1") {
					err = os.Remove(path)
					if err != nil {
						log.WithError(err).Errorf("removing %s", path)
						return err
					} else {
						log.Infof("removing %s", path)
					}
				}

				if strings.HasSuffix(filepath.Base(path), ".sh") {
					err = os.Remove(path)
					if err != nil {
						log.WithError(err).Errorf("removing %s", path)
						return err
					} else {
						log.Infof("removing %s", path)
					}
				}

				if filepath.Base(filepath.Dir(path)) == "bash" {
					err = os.Remove(path)
					if err != nil {
						log.WithError(err).Errorf("removing %s", path)
						return err
					} else {
						log.Infof("removing %s", path)
					}
				}
				return nil
			})
			log.DecreasePadding()
			if err != nil {
				log.WithError(err).Fatal("walk run-scripts")
			}
			log.Info("thanks for using aliax!")
		},
	}
)

func init() {
	aliaxCmd.AddCommand(cleanCmd)
}
