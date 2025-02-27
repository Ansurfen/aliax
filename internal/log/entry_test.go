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

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntry_WithField(t *testing.T) {
	a := NewEntry(New(io.Discard))
	b := a.WithField("foo", "bar")
	require.Empty(t, a.Fields.Keys())
	require.Equal(t, []string{"foo"}, b.Fields.Keys())
}

func TestEntry_WithError(t *testing.T) {
	a := NewEntry(New(io.Discard))
	b := a.WithError(fmt.Errorf("boom"))
	require.Empty(t, a.Fields.Keys())
	require.Equal(t, []string{"error"}, b.Fields.Keys())
}

func TestEntry_WithError_nil(t *testing.T) {
	a := NewEntry(New(io.Discard))
	b := a.WithError(nil)
	require.Empty(t, a.Fields.Keys())
	require.Empty(t, b.Fields.Keys())
}

func TestEntry_WithoutPadding(t *testing.T) {
	log := New(io.Discard)

	a := NewEntry(log)
	require.Equal(t, defaultPadding, a.Padding)

	log.IncreasePadding()
	b := NewEntry(log)
	require.Equal(t, defaultPadding+2, b.Padding)

	c := b.WithoutPadding()
	require.Equal(t, defaultPadding, c.Padding)
}
