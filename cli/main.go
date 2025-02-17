// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"aliax/cli/cmd"
	"aliax/internal/cfg"
	"aliax/internal/io"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/caarlos0/log"
	"github.com/google/shlex"
)

func main() {
	log.SetLevel(log.DebugLevel)
	if len(os.Args) > 1 {
		subCmd := map[string]struct{}{}
		for _, sub := range cmd.Root().Commands() {
			subCmd[sub.Use] = struct{}{}
		}
		if _, ok := subCmd[os.Args[1]]; !ok {
			cfgName := cfg.Name()
			var file cfg.Aliax
			err := io.ReadYAML(cfgName, &file)
			if err != nil {
				log.WithError(err).Fatalf("parsing %s", cfgName)
			}
			if script, ok := file.Script[os.Args[1]]; ok {
				log.WithField("script", script).Infof("running command: %s", os.Args[1])
				if err = ExecuteCommand(script); err != nil {
					log.WithError(err).Fatalf("running command: %s", script)
				}
				return
			}
		}
	}

	cmd.Execute()
}

func ExecuteCommand(cmdStr string) error {
	parts, err := shlex.Split(cmdStr)
	if err != nil {
		return fmt.Errorf("error splitting command: %v", err)
	}

	isWindows := runtime.GOOS == "windows"

	var cmdExec *exec.Cmd

	cmds := strings.Join(parts, " ")

	if isWindows {
		cmdExec = exec.Command("cmd", "/C", cmds)
	} else {
		cmdExec = exec.Command("bash", "-c", cmds)
	}

	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("error executing command: %v", err)
	}

	return nil
}
