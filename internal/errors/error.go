// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package errors

import "errors"

var (
	ErrCmdNotFinish = errors.New("command not finish")
	ErrCmdConflict  = errors.New("command conflict")
)
