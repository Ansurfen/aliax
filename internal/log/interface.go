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

// Interface represents the API of both Logger and Entry.
type Interface interface {
	WithField(string, any) *Entry
	WithError(error) *Entry
	WithoutPadding() *Entry
	Trace(string)
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
	Tracef(string, ...any)
	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
	ResetPadding()
	IncreasePadding()
	DecreasePadding()
}
