package main

import (
	"go/ast"
	"strings"
)

func (l *layout) literal(literal *ast.CompositeLit) {
	if len(literal.Elts) < 2 {
		return
	}

	for _, element := range literal.Elts {
		if _, keyed := element.(*ast.KeyValueExpr); !keyed {
			return
		}
	}

	last := literal.Elts[len(literal.Elts)-1]
	after := l.offset(last.End())
	closing := l.offset(literal.Rbrace)
	if !strings.HasPrefix(strings.TrimSpace(string(l.source[after:closing])), ",") {
		l.edits[after] = "," + l.edits[after]
	}

	for _, element := range literal.Elts {
		l.gap(element.Pos(), 1)
	}

	l.gap(literal.Rbrace, 1)
}
