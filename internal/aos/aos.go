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

// Create creates a new file with the given name.
// If the file already exists, it will be truncated.
// Returns the created file and an error, if any.
func Create(name string) (*os.File, error) {
	log.Debugf("creating %s", name)
	return os.Create(name)
}

// Remove removes the file or directory with the given name.
// Returns an error if the removal fails.
func Remove(name string) error {
	log.Debugf("removing %s", name)
	return os.Remove(name)
}

// Mkdir creates a directory with the specified path and permissions.
// Returns an error if the creation fails.
func Mkdir(path string, perm os.FileMode) error {
	log.WithField("permission", perm.String()).Debugf("creating %s", path)
	return os.Mkdir(path, perm)
}

// MkdirAll creates a directory and any necessary parent directories.
// It sets the permissions for the directory and its parents if they don't exist.
// Returns an error if the creation fails.
func MkdirAll(path string, perm os.FileMode) error {
	log.WithField("permission", perm.String()).Debugf("creating %s", path)
	return os.MkdirAll(path, perm)
}

// ReadYAML reads the contents of a YAML file and unmarshals it into the provided value.
// Returns an error if reading or unmarshaling fails.
func ReadYAML(filename string, v any) error {
	data, err := ReadFile(filename)
	if err != nil {
		return err
	}
	log.Debug("parsing configure")
	return yaml.Unmarshal(data, v)
}

// ReadJSON reads the contents of a JSON file and unmarshals it into the provided value.
// Returns an error if reading or unmarshaling fails.
func ReadJSON(filename string, v any) error {
	data, err := ReadFile(filename)
	if err != nil {
		return err
	}
	log.Debug("parsing configure")
	return json.Unmarshal(data, v)
}

// ReadFile reads the contents of the specified file and returns the data.
// Returns an error if reading the file fails.
func ReadFile(name string) ([]byte, error) {
	log.Debugf("reading %s", name)
	return os.ReadFile(name)
}

// Exist checks if a file or directory exists at the given path.
// Returns true if it exists, false if it does not, and an error if there is a failure checking the path.
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
