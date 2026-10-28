package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	check := len(os.Args) > 1 && os.Args[1] == "--check"
	failed := false
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() && path != "." && (strings.HasPrefix(entry.Name(), ".") || entry.Name() == "node_modules") {
			return filepath.SkipDir
		}

		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		changed, err := process(path, check)
		if changed && check {
			fmt.Println(path + ": run npm run format")

			failed = true
		}

		return err
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if failed {
		os.Exit(1)
	}
}

func process(path string, check bool) (bool, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	positions := token.NewFileSet()
	file, err := parser.ParseFile(positions, path, source, parser.ParseComments)
	if err != nil {
		return false, err
	}

	layout := newLayout(source, positions, file)
	layout.visit(file)

	result, err := format.Source(layout.apply())
	if err != nil {
		return false, fmt.Errorf("format %s: %w", path, err)
	}

	changed := !bytes.Equal(source, result)
	if changed && !check {
		return true, os.WriteFile(path, result, 0644)
	}

	return changed, nil
}
