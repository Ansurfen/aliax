// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package text

func In(s string, substrs []string) bool {
	for _, sub := range substrs {
		if sub == s {
			return true
		}
	}
	return false
}
