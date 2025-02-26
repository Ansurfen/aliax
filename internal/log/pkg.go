// This file is based on the original work by caarlos0 from https://github.com/caarlos0/log.
// The original code is licensed under the MIT License.
//
// # Copyright (c) 2022 Carlos Alexandro Becker
//
// Modifications and enhancements by ansurfen/aliax.
// We appreciate the contributions of caarlos0 and thank them for their work.
//
// SPDX-License-Identifier: MIT
package log

import (
	"io"
	"os"
)

// singletons ftw?
var Log Interface = New(os.Stderr)

// New creates a new logger.
func New(w io.Writer) *Logger {
	return &Logger{
		Writer:  w,
		Padding: defaultPadding,
		Level:   InfoLevel,
	}
}

// SetLevel sets the log level. This is not thread-safe.
func SetLevel(l Level) {
	if logger, ok := Log.(*Logger); ok {
		logger.Level = l
	}
}

// SetLevelFromString sets the log level from a string, panicing when invalid. This is not thread-safe.
func SetLevelFromString(s string) {
	if logger, ok := Log.(*Logger); ok {
		logger.Level = MustParseLevel(s)
	}
}

// ResetPadding resets the padding to default.
func ResetPadding() {
	Log.ResetPadding()
}

// IncreasePadding increases the padding 1 times.
func IncreasePadding() {
	Log.IncreasePadding()
}

// DecreasePadding decreases the padding 1 times.
func DecreasePadding() {
	Log.DecreasePadding()
}

// WithField returns a new entry with the `key` and `value` set.
func WithField(key string, value any) *Entry {
	return Log.WithField(key, value)
}

// WithError returns a new entry with the "error" set to `err`.
func WithError(err error) *Entry {
	return Log.WithError(err)
}

// WithoutPadding returns a new entry with padding set to default.
func WithoutPadding() *Entry {
	return Log.WithoutPadding()
}

// Trace level message.
func Trace(msg string) {
	Log.Trace(msg)
}

// Debug level message.
func Debug(msg string) {
	Log.Debug(msg)
}

// Info level message.
func Info(msg string) {
	Log.Info(msg)
}

// Warn level message.
func Warn(msg string) {
	Log.Warn(msg)
}

// Error level message.
func Error(msg string) {
	Log.Error(msg)
}

// Fatal level message, followed by an exit.
func Fatal(msg string) {
	Log.Fatal(msg)
}

// Trace level message.
func Tracef(msg string, v ...any) {
	Log.Tracef(msg, v...)
}

// Debugf level formatted message.
func Debugf(msg string, v ...any) {
	Log.Debugf(msg, v...)
}

// Infof level formatted message.
func Infof(msg string, v ...any) {
	Log.Infof(msg, v...)
}

// Warnf level formatted message.
func Warnf(msg string, v ...any) {
	Log.Warnf(msg, v...)
}

// Errorf level formatted message.
func Errorf(msg string, v ...any) {
	Log.Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...any) {
	Log.Fatalf(msg, v...)
}
