package parse

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

var ErrGoModNotFound = errors.New("go mod not found")

func getModulePath(root string) (string, error) {
	gomodPath := filepath.Join(root, "go.mod")
	if _, err := os.Stat(gomodPath); err != nil {
		if os.IsNotExist(err) {
			return "", ErrGoModNotFound
		}
		return "", fmt.Errorf("os.Stat:%w", err)
	}
	open, err := os.Open(gomodPath)
	if err != nil {
		return "", fmt.Errorf("os.Open:%w", err)
	}
	bytes, err := io.ReadAll(open)
	if err != nil {
		return "", fmt.Errorf("io.ReadAll:%w", err)
	}
	parse, err := modfile.Parse(gomodPath, bytes, nil)
	if err != nil {
		return "", fmt.Errorf("modfile.Parse:%w", err)
	}
	return parse.Module.Mod.Path, nil
}
