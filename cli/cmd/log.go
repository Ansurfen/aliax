// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"strconv"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

// logCmdParameter stores parameters for the "log" command.
type logCmdParameter struct {
	level   string
	fields  map[string]string
	msg     string
	padding int
}

var (
	logParameter logCmdParameter
	logCmd       = &cobra.Command{
		Use:     "log",
		Short:   "",
		Long:    "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			for k, v := range logParameter.fields {
				if vv, err := strconv.Unquote(`"` + v + `"`); err != nil {
					log.Log = log.Log.WithField(k, v)
				} else {
					log.Log = log.Log.WithField(k, vv)
				}
			}

			for range logParameter.padding {
				log.Log.IncreasePadding()
			}

			switch logParameter.level {
			case "info":
				log.Log.Info(logParameter.msg)
			case "warn":
				log.Log.Warn(logParameter.msg)
			case "error":
				log.Log.Error(logParameter.msg)
			case "fatal":
				log.Log.Fatal(logParameter.msg)
			case "debug":
				log.Log.Debug(logParameter.msg)
			}
		},
	}
)

func init() {
	aliaxCmd.AddCommand(logCmd)
	logCmd.PersistentFlags().StringVarP(&logParameter.level, "level", "l", "info", "")
	logCmd.PersistentFlags().StringVarP(&logParameter.msg, "message", "m", "", "")
	logCmd.PersistentFlags().StringToStringVarP(&logParameter.fields, "field", "f", map[string]string{}, "")
	logCmd.PersistentFlags().IntVarP(&logParameter.padding, "padding", "p", 0, "")
}
