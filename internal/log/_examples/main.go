// This file is based on the original work by caarlos0 from https://github.com/caarlos0/log.
// The original code is licensed under the MIT License.
//
// Copyright (c) 2022 Carlos Alexandro Becker
// 
// Modifications and enhancements by ansurfen/aliax.
// We appreciate the contributions of caarlos0 and thank them for their work.
//
// SPDX-License-Identifier: MIT
package main

import (
	"aliax/internal/log"
	"errors"
)

func main() {
	log.SetLevel(log.TraceLevel)
	log.Trace("enter main")
	log.
		WithField("field5", "value5").
		WithField("field2", "value2").
		WithField("field1", "value1").
		WithField("field4", "value4").
		WithField("FOO", "bar").
		WithField("field3", "value3").
		Info("AQUI")
	log.WithField("foo", "bar").Debug("debug")
	log.WithField("foo", "bar").Info("info")
	log.WithField("foo", "bar").Warn("warn")
	log.WithField("multiple", "fields").
		WithField("yes", true).
		Info("a longer line in this particular log")
	log.IncreasePadding()
	log.WithField("foo", "bar").Info("info with increased padding")
	log.IncreasePadding()
	log.WithField("foo", "bar").
		WithField("text", "a multi\nline text going\non for multiple lines\nhello\nworld!").
		Info("info with a more increased padding")
	log.WithoutPadding().WithField("foo", "bar").Info("info without padding")
	log.WithField("foo", "bar").Info("info with a more increased padding")
	log.ResetPadding()
	log.WithField("foo", "bar").
		WithField("text", "a multi\nline text going\non for multiple lines\nhello\nworld!").
		WithField("another", "bar").
		WithField("lalalal", "bar").
		Info("info with a more increased padding")
	log.WithError(errors.New("some error")).Error("error")
	log.WithError(errors.New("some fatal error")).Fatal("fatal")
}
