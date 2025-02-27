// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	cmd := aliaxCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"log", "-m Hello World"})

	err := cmd.Execute()
	assert.NoError(t, err)
	fmt.Println(buf.String(), 6)
	assert.Contains(t, buf.String(), "Hello World")
}
