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

	name string `yaml:"-"`
}

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
	usage := ""
	switch {
	case len(c.Command) == 0 && len(c.Flags) > 0:
		usage = fmt.Sprintf("Usage:\n  %s [flags]", c.name)
	case len(c.Command) == 0 && len(c.Flags) == 0:
		usage = fmt.Sprintf("Usage:\n  %s", c.name)
	case len(c.Command) > 0 && len(c.Flags) > 0:
		usage = fmt.Sprintf("Usage:\n  %s [command] [flags]", c.name)
	}

	example := ""
	if len(c.Example) > 0 {
		example = fmt.Sprintf("Example:\n%s", c.Example)
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
		flags = append(flags, fmt.Sprintf("  %s\t%s", strings.Join(flag.Alias, ", "), flag.Usage))
	}

	availableflags := ""
	if len(flags) > 0 {
		availableflags = fmt.Sprintf("Flags:\n%s\n", strings.Join(flags, "\n"))
	}
	return strings.TrimSpace(fmt.Sprintf(`%s

%s
%s
%s
%s`, c.Long, usage, example, availableCommands, availableflags))
}

type Aliax struct {
	Executable string              `yaml:"executable"`
	RunPath    string              `yaml:"runPath"`
	Variable   map[string]string   `yaml:"variable"`
	Extend     map[string]*Command `yaml:"extend"`
	Command    map[string]*Command `yaml:"command"`
	Script     map[string]Script   `yaml:"script"`
}

// TODO
type Commands map[string]*Command

func (c *Commands) LoadBinPath() {

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
