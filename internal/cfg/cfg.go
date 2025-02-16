// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cfg

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
