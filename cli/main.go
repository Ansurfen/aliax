// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"aliax/cli/cmd"
	"aliax/internal/aio"
	"aliax/internal/aos"
	"aliax/internal/cfg"
	"aliax/internal/errors"
	"aliax/internal/shell"
	"aliax/internal/style"
	"aliax/internal/template"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/caarlos0/log"
	"github.com/google/shlex"
	"github.com/spf13/cobra"
)

var customCmd = &cobra.Command{
	Use:                "aliax <sub_cmd>",
	Hidden:             true,
	DisableSuggestions: true,
	SilenceErrors:      true,
	SilenceUsage:       true,
	RunE: func(cmd *cobra.Command, args []string) error {
		sub_cmd := args[0]
		subCmd := map[string]struct{}{}
		dry, _ := cmd.Flags().GetBool("dry")
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
		for _, sub := range cmd.Root().Commands() {
			subCmd[sub.Use] = struct{}{}
		}
		if _, ok := subCmd[sub_cmd]; !ok {
			cfgName := cfg.Name()
			var file cfg.Aliax
			log.Infof("parsing %s", cfgName)
			err := aio.ReadYAML(cfgName, &file)
			if err != nil {
				log.WithError(err).Fatalf("fail to parse file")
			}
			for name := range file.Script {
				if _, ok := subCmd[name]; ok {
					log.WithError(errors.ErrCmdConflict).
						WithField("target", cfgName).
						WithField("command", name).
						WithField("suggestion", fmt.Sprintf(`please rename your custom command
the following command names are not allowed. they are built-in commands for Aliax:
%s`, style.Keyword("init、clean、env、log、version"))).Fatal("invalid script")
				}
			}
			if script, ok := file.Script[sub_cmd]; ok {
				log.Debugf("initializing variables")
				log.IncreasePadding()
				for k, v := range file.Variable {
					buf := &strings.Builder{}
					log.Debugf("initializing %s", k)
					err = template.Execute(buf, v, nil)
					if err != nil {
						log.WithError(err).Fatal("fail to execute template")
					}
					file.Variable[k] = buf.String()
				}
				log.DecreasePadding()
				if script.Run != nil {
					log.WithField("script", *script.Run).Infof("running command: %s", sub_cmd)
					if err = executeCommand(*script.Run); err != nil {
						log.WithError(err).Fatalf("running command: %s", *script.Run)
					}
				} else {
					matches := (*script.Cmd).Match
					// TODO map collect
					for _, c := range matches {
						if runtime.GOOS == "windows" {
							buf := &strings.Builder{}
							err = template.Execute(buf, c.Run, file.Variable)
							if err != nil {
								log.WithError(err).Fatal("fail to execute template")
							}
							if dry {
								log.WithField("script", buf.String()).Info("dry mode")
								return nil
							}
							log.WithField("script", buf.String()).Info("running command")
							err = execute(buf.String())
							if err != nil {
								log.WithError(err).Fatal("fail to executing command")
							}
							return nil
						} else {
							fmt.Println(c.Run)
							return nil
						}
					}
				}
				return nil
			}
		}
		return errors.ErrCmdNotFinish
	},
}

func main() {
	log.SetLevel(log.InfoLevel)
	logFileName := filepath.Join(aos.LogPath, time.Now().Format("2006-01-02_15-04-05")+".log")
	log.Debugf("creating %s", logFileName)
	logFile, err := os.Create(logFileName)
	if err != nil {
		log.WithError(err).Fatal("fail to create log")
	}

	log.Log = log.New(io.MultiWriter(os.Stderr, logFile))
	customCmd.Flags().BoolP("dry", "d", false, "")
	customCmd.Flags().BoolP("verbose", "v", false, "")

	if len(os.Args) > 1 {
		err := customCmd.Execute()
		if err == nil {
			return
		}
	}

	cmd.Execute()
}

func execute(cmdStr string) error {
	if strings.Contains(cmdStr, "\n") {
		return shell.OnceScript(cmdStr)
	}
	return executeCommand(cmdStr)
}

func executeCommand(cmdStr string) error {
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
