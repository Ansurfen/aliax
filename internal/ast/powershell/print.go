// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package psast

import (
	token "aliax/internal/token/powershell"
	"fmt"
	"io"
	"strings"
)

func Print(node Node, w io.Writer) {
	print(w, node, "")
}

func print(w io.Writer, node Node, space string) {
	switch node := node.(type) {
	case *File:
		for _, s := range node.Stmts {
			print(w, s, space)
		}
	case *BlockStmt:
		fmt.Fprintf(w, " {\n")
		for _, s := range node.List {
			print(w, s, space+"  ")
		}
		fmt.Fprintln(w, space+"}")
	case *IfStmt:
		fmt.Fprintf(w, space+"if (%s)", node.Cond)
		print(w, node.Body, space)
		for node.Else != nil {
			switch el := node.Else.(type) {
			case *IfStmt:
				fmt.Fprintf(w, space+"elseif (%s)", el.Cond)
				print(w, el.Body, space)
				node.Else = el.Else
			case *BlockStmt:
				fmt.Fprint(w, space+"else")
				print(w, el, space)
				node.Else = nil
			}
		}
	case *ForStmt:
		fmt.Fprintf(w, space+"for (%s; %s; %s)", node.Init, node.Cond, node.Post)
		print(w, node.Body, space)
	case *SwitchStmt:
		if node.Mode != MatchModeNone {
			fmt.Fprintf(w, space+"switch %s (%s) {\n", node.Mode, node.Cond)
		} else {
			fmt.Fprintf(w, space+"switch (%s) {\n", node.Cond)
		}
		for _, c := range node.Cases {
			fmt.Fprintf(w, space+"  %s", c.Cond)
			print(w, c.Body, space+"  ")
		}
		if node.Default != nil {
			fmt.Fprint(w, space+"  default")
			print(w, node.Default.Body, space+"  ")
		}
		fmt.Fprintln(w, space+"}")
	case *ExprStmt:
		fmt.Fprintf(w, space+"%s\n", node.X)
	case *AssignStmt:
		fmt.Fprintf(w, space+"%s = %s\n", node.Lhs, node.Rhs)
	case *CallStmt:
		recv := []string{}
		for _, r := range node.Recv {
			recv = append(recv, r.String())
		}
		if node.Op == token.BITAND {
			fmt.Fprintf(w, space+"& %s %s\n", node.Func, strings.Join(recv, " "))
		} else {
			fmt.Fprintf(w, space+"%s %s\n", node.Func, strings.Join(recv, " "))
		}
	case *Comment:
		fmt.Fprintf(w, space+"#%s\n", node.Text)
	default:
		panic(node)
	}
}
