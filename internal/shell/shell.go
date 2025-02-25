// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package shell

import (
	"aliax/internal/aos"
	"aliax/internal/text"
	"fmt"
	"os"
	"os/exec"

	"github.com/caarlos0/log"
)

// StartCmd prepares a command for execution. It handles Windows and Unix-based OSes separately.
// On Windows, it runs the command through `cmd` with `/C`.
// On Unix-based OSes (Linux/macOS), it runs the command through `bash` with the `-c` option.
func StartCmd(name string, arg ...string) *exec.Cmd {
	if aos.IsWindows {
		return exec.Command("cmd", append([]string{"/C", name}, arg...)...)
	}
	return exec.Command("bash", append([]string{"-c", name}, arg...)...)
}

// Run executes a command and captures its output (both stdout and stderr).
// It then converts the output from GBK encoding to UTF-8 using `GBK2UTF8` function.
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

// OnceScript creates a temporary script file, writes the provided script content to it,
// and then executes it once based on the operating system.
func OnceScript(s string) error {
	suffix := ".sh"
	if aos.IsWindows {
		suffix = ".ps1"
	}
	tmpFile, err := os.CreateTemp(".", fmt.Sprintf("aliax_temp_*.%s", suffix))
	if err != nil {
		return err
	}
	defer aos.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(s)
	if err != nil {
		return err
	}

	tmpFile.Close()

	var cmd *exec.Cmd
	if aos.IsWindows {
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

// LookPath searches for an executable file in the system's PATH and returns its full path.
func LookPath(file string) (string, error) {
	log.Debugf("looking path %s", file)
	return exec.LookPath(file)
}
