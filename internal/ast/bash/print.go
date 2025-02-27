// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package bashast

import (
	"fmt"
	"io"
	"strings"
)

// Print writes the string representation of the given AST node to the provided writer.
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
		for _, s := range node.List {
			print(w, s, space+"  ")
		}
	case *IfStmt:
		fmt.Fprintf(w, space+"if [[ %s ]]; then\n", node.Cond)
		print(w, node.Body, space)
		for node.Else != nil {
			switch el := node.Else.(type) {
			case *IfStmt:
				fmt.Fprintf(w, space+"elif [[ %s ]]; then\n", el.Cond)
				print(w, el.Body, space)
				node.Else = el.Else
			case *BlockStmt:
				fmt.Fprint(w, space+"else")
				print(w, el, space)
				node.Else = nil
			}
		}
		fmt.Fprintln(w, space+"fi")
	case *ForStmt:
		fmt.Fprintf(w, space+"for ((%s; %s; %s)); do\n", node.Init, node.Cond, node.Post)
		print(w, node.Body, space)
		fmt.Fprintln(w, space+"done")
	case *SwitchStmt:
		fmt.Fprintf(w, space+"case %s in\n", node.Cond)

		for _, c := range node.Cases {
			fmt.Fprintf(w, space+"  %s)\n", c.Cond)
			for _, s := range c.Body.List {
				print(w, s, space+"    ")
			}
			fmt.Fprintln(w, space+"    ;;")
		}
		if node.Default != nil {
			fmt.Fprintln(w, space+"  *)")
			for _, s := range node.Default.Body.List {
				print(w, s, space+"    ")
			}
			fmt.Fprintln(w, space+"    ;;")
		}
		fmt.Fprintln(w, space+"esac")
	case *ExprStmt:
		fmt.Fprintf(w, space+"%s\n", node.X)
	case *AssignStmt:
		fmt.Fprintf(w, space+"%s=%s\n", node.Lhs, node.Rhs)
	case *CallStmt:
		recv := []string{}
		for _, r := range node.Recv {
			recv = append(recv, r.String())
		}
		fmt.Fprintf(w, space+"%s %s\n", node.Func, strings.Join(recv, " "))
	case *Comment:
		fmt.Fprintf(w, space+"#%s\n", node.Text)
	}
}
