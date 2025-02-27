// This file is based on the original work by caarlos0 from https://github.com/caarlos0/log.
// The original code is licensed under the MIT License.
//
// Copyright (c) 2022 Carlos Alexandro Becker
// 
// Modifications and enhancements by ansurfen/aliax.
// We appreciate the contributions of caarlos0 and thank them for their work.
//
// SPDX-License-Identifier: MIT
package log

import "context"

// logKey is a private context key.
type logKey struct{}

// NewContext returns a new context with logger.
func NewContext(ctx context.Context, v Interface) context.Context {
	return context.WithValue(ctx, logKey{}, v)
}

// FromContext returns the logger from context, or log.Log.
func FromContext(ctx context.Context) Interface {
	if v, ok := ctx.Value(logKey{}).(Interface); ok {
		return v
	}
	return Log
}
