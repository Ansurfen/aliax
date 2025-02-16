// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cfg

import (
	"aliax/internal/io"
	"os"
)

type Flag struct {
	Name  string   `yaml:"name"`
	Alias []string `yaml:"alias"`
	Type  string   `yaml:"type"`
	Usage string   `yaml:"usage"`
}

type Case struct {
	Pattern  any    `yaml:"pattern"`
	Platform string `yaml:"platform"`
	Run      string `yaml:"run"`
}

type Command struct {
	Short   string             `yaml:"short"`
	Long    string             `yaml:"long"`
	Example string             `yaml:"example"`
	Flags   []Flag             `yaml:"flags"`
	Match   []Case             `yaml:"match"`
	Command map[string]Command `yaml:"command"`
	Bin     string             `yaml:"bin"`
}

type Aliax struct {
	Extend  map[string]Command `yaml:"extend"`
	Script  map[string]string  `yaml:"script"`
	Command map[string]Command `yaml:"command"`
}

const work = "aliax.work"

var target = ""

func Name() string {
	if len(target) > 0 {
		return target
	}
	if ok, _ := io.Exist(work); ok {
		data, err := os.ReadFile(work)
		if err == nil {
			target = string(data)
		} else {
			target = "aliax.yaml"
		}
	} else {
		target = "aliax.yaml"
	}
	return target
}
