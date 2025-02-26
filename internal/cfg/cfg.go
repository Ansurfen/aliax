// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cfg

import (
	"aliax/internal/aos"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
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
	DisableHelp bool                `yaml:"disableHelp"`
	Short       string              `yaml:"short"`
	Long        string              `yaml:"long"`
	Example     string              `yaml:"example"`
	Flags       []Flag              `yaml:"flags"`
	Match       []Case              `yaml:"match"`
	Command     map[string]*Command `yaml:"command"`
	Bin         string              `yaml:"bin"`

	name     string              `yaml:"-"`
	flagDict map[string]flagType `yaml:"-"`
}

type flagType uint8

const (
	flagTypeString flagType = iota
	flagTypeBool
)

func (c *Command) SetName(name string) {
	c.name = name
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Preload(name string) error {
	if !c.DisableHelp {
		flags := map[string]struct{}{}
		for _, flag := range c.Flags {
			flags[flag.Name] = struct{}{}
		}
		if _, ok := flags["help"]; ok {
			// add suggestion, e.g. try disable or instead of name
			return fmt.Errorf("help flag ready exist")
		}
		c.Flags = append(c.Flags, Flag{
			Name:  "help",
			Alias: []string{"-h", "--help"},
			Type:  "bool",
			Usage: fmt.Sprintf("help for %s", name),
		})
	}
	return nil
}

func (c *Command) HelpCmd(name string) string {
	if c.DisableHelp {
		return ""
	}
	cmds := []string{}
	for cmdName, cmd := range c.Command {
		cmds = append(cmds, fmt.Sprintf("  %s\t%s", cmdName, cmd.Short))
	}
	availableCommands := ""
	if len(cmds) > 0 {
		availableCommands = fmt.Sprintf("Available Commands:\n%s\n", strings.Join(cmds, "\n"))
	}

	flags := []string{}
	for _, flag := range c.Flags {
		flags = append(flags, fmt.Sprintf("  %s\t%s", flag.Name, flag.Usage))
	}

	availableflags := ""
	if len(flags) > 0 {
		availableflags = fmt.Sprintf("Flags:\n%s\n", strings.Join(flags, "\n"))
	}
	return fmt.Sprintf(`%s

Usage:
  %s [command]

%s
%s`, c.Long, name, availableCommands, availableflags)
}

type Aliax struct {
	Variable map[string]string   `yaml:"variable"`
	Extend   map[string]*Command `yaml:"extend"`
	Script   map[string]Script   `yaml:"script"`
	Command  map[string]*Command `yaml:"command"`
}

const work = "aliax.work"

var target = ""

func Name() string {
	if len(target) > 0 {
		return target
	}
	if ok, _ := aos.Exist(work); ok {
		data, err := aos.ReadFile(work)
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

type Script struct {
	Cmd *Command
	Run *string
}

func (sc *Script) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		var run string
		if err := value.Decode(&run); err != nil {
			return err
		}
		sc.Run = &run
		return nil
	}

	var cmd Command
	if err := value.Decode(&cmd); err != nil {
		return err
	}
	sc.Cmd = &cmd
	return nil
}
