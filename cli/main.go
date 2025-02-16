// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"aliax/cli/cmd"
	"aliax/internal/cfg"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/google/shlex"
	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) > 1 {
		subCmd := map[string]struct{}{}
		for _, sub := range cmd.Root().Commands() {
			subCmd[sub.Use] = struct{}{}
		}
		if _, ok := subCmd[os.Args[1]]; !ok {
			data, err := os.ReadFile(cfg.Name())
			if err != nil {
				panic(err)
			}
			var file cfg.Aliax
			err = yaml.Unmarshal(data, &file)
			if err != nil {
				panic(err)
			}
			if script, ok := file.Script[os.Args[1]]; ok {
				if err = ExecuteCommand(script); err != nil {
					panic(err)
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

	if isWindows {
		cmdExec = exec.Command("cmd", "/C", strings.Join(parts, " "))
	} else {
		cmdExec = exec.Command("bash", "-c", strings.Join(parts, " "))
	}

	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("error executing command: %v", err)
	}

	return nil
}
