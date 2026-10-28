package main

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"
)

type layout struct {
	source    []byte
	positions *token.FileSet
	file      *ast.File
	edits     map[int]string
}

func newLayout(source []byte, positions *token.FileSet, file *ast.File) *layout {
	return &layout{
		source:    source,
		positions: positions,
		file:      file,
		edits:     map[int]string{},
	}
}

func (l *layout) offset(position token.Pos) int {
	return l.positions.Position(position).Offset
}

func (l *layout) gap(position token.Pos, lines int) {
	offset := l.offset(position)
	start := offset

	for start > 0 && strings.ContainsRune(" \t\r\n", rune(l.source[start-1])) {
		start--
	}

	count := strings.Count(string(l.source[start:offset]), "\n")
	if count < lines {
		addition := strings.Repeat("\n", lines-count)
		if !strings.HasSuffix(l.edits[offset], addition) {
			l.edits[offset] += addition
		}
	}
}

func (l *layout) visit(file *ast.File) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch current := node.(type) {
		case *ast.FuncDecl:
			position := current.Pos()
			if current.Doc != nil {
				position = current.Doc.Pos()
			}

			l.gap(position, 2)
		case *ast.BlockStmt:
			l.block(current)
		case *ast.CompositeLit:
			l.literal(current)
		}

		return true
	})
}

func (l *layout) block(block *ast.BlockStmt) {
	for index, statement := range block.List {
		lines := 1
		if index > 0 && separate(block.List[index-1], statement) {
			lines = 2
		}

		position := statement.Pos()

		for _, comment := range l.file.Comments {
			if comment.End() < position && comment.Pos() > block.Lbrace && (index == 0 || comment.Pos() > block.List[index-1].End()) {
				position = comment.Pos()
				break
			}
		}

		l.gap(position, lines)
	}

	if len(block.List) > 0 {
		l.gap(block.Rbrace, 1)
	}
}

func separate(previous, current ast.Stmt) bool {
	if _, cleanup := current.(*ast.DeferStmt); cleanup {
		return false
	}

	switch previous.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt, *ast.DeferStmt:
		return true
	}

	switch current.(type) {
	case *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.ReturnStmt:
		return true
	}

	_, previousCall := previous.(*ast.ExprStmt)
	_, nextAssignment := current.(*ast.AssignStmt)

	return previousCall && nextAssignment
}

func (l *layout) apply() []byte {
	offsets := make([]int, 0, len(l.edits))

	for offset := range l.edits {
		offsets = append(offsets, offset)
	}

	slices.Sort(offsets)

	result := append([]byte{}, l.source...)

	for _, offset := range slices.Backward(offsets) {
		result = append(result[:offset], append([]byte(l.edits[offset]), result[offset:]...)...)
	}

	return result
}
