// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"aliax/cli/cmd"
	"aliax/internal/cfg"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) > 1 {
		subCmd := map[string]struct{}{}
		for _, sub := range cmd.Root().Commands() {
			subCmd[sub.Use] = struct{}{}
		}
		if _, ok := subCmd[os.Args[1]]; !ok {
			data, err := os.ReadFile("aliax.yaml")
			if err != nil {
				panic(err)
			}
			var file cfg.Aliax
			err = yaml.Unmarshal(data, &file)
			if err != nil {
				panic(err)
			}
			if script, ok := file.Script[os.Args[1]]; ok {
				fmt.Println(script)
				return
			}
		}
	}

	cmd.Execute()
}
