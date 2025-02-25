// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package token

type Token int

const (
	None Token = iota
	ADD        // +
	SUB        // -

	ASSIGN // =

	ADD_ASSIGN // +=

	AND // &

	EQ // ==
	NE // !=
	LT // <
	GT // >

	Inc // ++
	Dec // --

	DOUBLE_DOT // ..

	STRING
	NUMBER
	BOOL
)

func (t Token) String() string {
	return []string{"", "+", "-", "=", "+=", "&", "-eq", "-ne", "-lt", "-gt", "++", "--", ".."}[t]
}
