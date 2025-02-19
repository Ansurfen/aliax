// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package shell

import (
	"aliax/internal/text"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/caarlos0/log"
)

func Run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	raw, err := cmd.CombinedOutput()

	output, err2 := text.GBK2UTF8(raw)
	if err2 != nil {
		if err != nil {
			log.WithError(err).WithField("script", cmd).WithField("output", string(raw)).Error("running command")
			return err
		}
		log.WithField("script", cmd).WithField("output", string(raw)).Info("running command")
	} else {
		if err != nil {
			log.WithError(err).WithField("script", cmd).WithField("output", string(output)).Error("running command")
			return err
		}
		log.WithField("script", cmd).WithField("output", string(output)).Info("running command")
	}
	return nil
}

func OnceScript(s string) error {
	suffix := ".sh"
	if runtime.GOOS == "windows" {
		suffix = ".ps1"
	}
	tmpFile, err := os.CreateTemp(".", fmt.Sprintf("aliax_temp_*.%s", suffix))
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(s)
	if err != nil {
		return err
	}

	tmpFile.Close()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", tmpFile.Name())
	} else {
		cmd = exec.Command("bash", tmpFile.Name())
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func LookPath(file string) (string, error) {
	log.Debugf("looking path %s", file)
	return exec.LookPath(file)
}
