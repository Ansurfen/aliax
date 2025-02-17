// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package shell

import (
	"os/exec"

	"github.com/caarlos0/log"
	"github.com/djimenez/iconv-go"
)

func Run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	raw, err := cmd.CombinedOutput()
	output, err2 := iconv.ConvertString(string(raw), "GBK", "UTF-8")
	if err2 != nil {
		if err != nil {
			log.WithError(err).WithField("script", cmd).WithField("output", string(raw)).Error("running command")
			return err
		}
		log.WithField("script", cmd).WithField("output", string(raw)).Info("running command")
	} else {
		if err != nil {
			log.WithError(err).WithField("script", cmd).WithField("output", output).Error("running command")
			return err
		}
		log.WithField("script", cmd).WithField("output", output).Info("running command")
	}
	return nil
}
