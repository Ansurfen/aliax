// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"encoding/json"
	"errors"

	"github.com/ansurfen/globalenv"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:     "env",
	Short:   "Read and manage the ALIAXPATH environment variable",
	Long:    `The 'env' command allows you to read and display the ALIAXPATH environment variable,
which is a serialized JSON string. It provides key-value pairs that can be accessed and displayed in a human-readable format.`,
	Example: `aliax env`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("reading ALIAXPATH from env")
		path, err := getAliaxPath()
		if err != nil {
			log.WithError(err).Fatal("parsing environment")
		}
		log.IncreasePadding()
		for k, v := range path {
			log.Infof("%s → %s", k, v)
		}
		log.DecreasePadding()
	},
}

func init() {
	aliaxCmd.AddCommand(envCmd)
}

var errAliaxPathNotFound = errors.New("ALIAXPATH not set")

func getAliaxPath() (map[string]string, error) {
	serializedData, _ := globalenv.Get("ALIAXPATH")
	if serializedData == "" {
		return map[string]string{}, errAliaxPathNotFound
	}

	var data map[string]string
	err := json.Unmarshal([]byte(serializedData), &data)
	if err != nil {
		return map[string]string{}, err
	}

	return data, nil
}

func setAliaxPath(data map[string]string) ([]byte, error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return globalenv.Set("ALIAXPATH", string(serializedData))
}
