// This file is based on the original work by caarlos0 from https://github.com/caarlos0/log.
// The original code is licensed under the MIT License.
//
// Copyright (c) 2022 Carlos Alexandro Becker
// 
// Modifications and enhancements by ansurfen/aliax.
// We appreciate the contributions of caarlos0 and thank them for their work.
//
// SPDX-License-Identifier: MIT
package log_test

import (
	"context"
	"testing"

	"github.com/caarlos0/log"
	"github.com/stretchr/testify/require"
)

func TestFromContext(t *testing.T) {
	ctx := context.Background()

	logger := log.FromContext(ctx)
	require.Equal(t, log.Log, logger)

	logs := log.WithField("foo", "bar")
	ctx = log.NewContext(ctx, logs)

	logger = log.FromContext(ctx)
	require.Equal(t, logs, logger)
}
