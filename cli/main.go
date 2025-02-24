// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"aliax/cli/cmd"
	"aliax/internal/aos"
	"aliax/internal/cfg"
	"aliax/internal/errors"
	"aliax/internal/shell"
	"aliax/internal/style"
	"aliax/internal/template"
	"aliax/internal/text"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caarlos0/log"
	"github.com/google/shlex"
)

func executeCustomCmd() error {
	sub_cmd := os.Args[1]
	subCmd := map[string]struct{}{}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	for _, sub := range cmd.Root().Commands() {
		subCmd[sub.Use] = struct{}{}
	}
	if _, ok := subCmd[sub_cmd]; !ok {
		cfgName := cfg.Name()
		var file cfg.Aliax
		err := aos.ReadYAML(cfgName, &file)
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
					if aos.IsWindows {
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
						log.WithField("script", c.Run).Infof("running command: %s", sub_cmd)
						if err = executeCommand(c.Run); err != nil {
							log.WithError(err).Fatalf("running command: %s", c.Run)
						}
						return nil
					}
				}
			}
			return nil
		}
	}
	return errors.ErrCmdNotFinish
}

var (
	verbose bool
	dry     bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&dry, "dry", false, "")
	flag.BoolVar(&dry, "d", false, "")
}

func main() {
	log.SetLevel(log.InfoLevel)
	logFileName := filepath.Join(aos.LogPath, time.Now().Format("2006-01-02_15-04-05")+".log")
	logFile, err := aos.Create(logFileName)
	if err != nil {
		log.WithError(err).Fatal("fail to create log")
	}

	log.Log = log.New(io.MultiWriter(os.Stderr, logFile))

	flag.Parse()

	if len(os.Args) > 1 {
		err := executeCustomCmd()
		if err == nil {
			return
		}
	}

	cmd.Execute()

	if len(os.Args) > 1 {
		if text.In(os.Args[1], []string{"clean", "init", "env"}) {
			log.Info("thanks for using aliax!")
		}
	}
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

	cmds := strings.Join(parts, " ")

	cmdExec := shell.StartCmd(cmds)

	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("error executing command: %v", err)
	}

	return nil
}
