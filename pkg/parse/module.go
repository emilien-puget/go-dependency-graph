package parse

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func getModulePath(root string) (string, error) {
	gomodPath := filepath.Join(root, "go.mod")
	open, err := os.Open(gomodPath)
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(open)
	if err != nil {
		return "", err
	}
	parse, err := modfile.Parse(gomodPath, bytes, nil)
	if err != nil {
		return "", err
	}
	return parse.Module.Mod.Path, nil
}
