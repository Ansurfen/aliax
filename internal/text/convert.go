// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package text

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GBK2UTF8 converts a GBK-encoded byte slice to UTF-8 encoded byte slice.
// It takes a GBK-encoded byte slice as input and returns the UTF-8 encoded byte slice.
// If there is an error during the conversion, it returns an error.
func GBK2UTF8(str []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(str), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}
