// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package aos

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/caarlos0/log"
)

var (
	// RootPath to where the executable is stored
	RootPath string
	LogPath  string
)

func init() {
	path, err := os.Executable()
	if err != nil {
		log.WithError(err).Fatal("fail to read root path")
	}

	RootPath = filepath.Dir(path)
	LogPath = filepath.Join(RootPath, "log")

	log.Debugf("creating %s", LogPath)
	err = os.Mkdir(LogPath, 0755)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			log.WithError(err).Debug("fail to create log path")
		} else {
			log.WithError(err).Fatal("fail to create log path")
		}
	}
}
