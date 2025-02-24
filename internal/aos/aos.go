// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package aos

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/log"
	"gopkg.in/yaml.v3"
)

var (
	// RootPath to where the executable is stored
	RootPath     string
	LogPath      string
	TemplatePath string

	IsWindows = runtime.GOOS == "windows"
)

func init() {
	path, err := os.Executable()
	if err != nil {
		log.WithError(err).Fatal("fail to read root path")
	}

	RootPath = filepath.Dir(path)
	LogPath = filepath.Join(RootPath, "log")
	TemplatePath = filepath.Join(RootPath, "template")

	err = Mkdir(LogPath, 0755)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			log.WithError(err).Debug("fail to create log path")
		} else {
			log.WithError(err).Fatal("fail to create log path")
		}
	}

	err = Mkdir(TemplatePath, 0755)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			log.WithError(err).Debug("fail to create template path")
		} else {
			log.WithError(err).Fatal("fail to create template path")
		}
	}
}

func Create(name string) (*os.File, error) {
	log.Debugf("creating %s", name)
	return os.Create(name)
}

func Remove(name string) error {
	log.Debugf("removing %s", name)
	return os.Remove(name)
}

func Mkdir(path string, perm os.FileMode) error {
	log.WithField("permission", perm.String()).Debugf("creating %s", path)
	return os.Mkdir(path, perm)
}

func MkdirAll(path string, perm os.FileMode) error {
	log.WithField("permission", perm.String()).Debugf("creating %s", path)
	return os.MkdirAll(path, perm)
}

func ReadYAML(filename string, v any) error {
	data, err := ReadFile(filename)
	if err != nil {
		return err
	}
	log.Debug("parsing configure")
	return yaml.Unmarshal(data, v)
}

func ReadJSON(filename string, v any) error {
	data, err := ReadFile(filename)
	if err != nil {
		return err
	}
	log.Debug("parsing configure")
	return json.Unmarshal(data, v)
}

func ReadFile(name string) ([]byte, error) {
	log.Debugf("reading %s", name)
	return os.ReadFile(name)
}

func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
