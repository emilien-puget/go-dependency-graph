package package_list

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	goFileExtension = ".go"
)

func GetPackagesToParse(pathDir string, skipDirs []string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Dir:   pathDir,
		Mode:  packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Tests: false,
	}
	dirs, err := findGoSourceDirectories(pathDir, skipDirs)
	if err != nil {
		return nil, fmt.Errorf("findGoSourceDirectories: %w", err)
	}
	pkgs, err := packages.Load(cfg, dirs...)
	if err != nil {
		return nil, fmt.Errorf("packages.Load: %w", err)
	}
	return pkgs, nil
}

func findGoSourceDirectories(pathDir string, skipDirs []string) ([]string, error) {
	uniqueDirs := make(map[string]bool)

	err := filepath.WalkDir(pathDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("filepath.WalkDir: %w", err)
		}

		for i := range skipDirs {
			if strings.Contains(p, string(os.PathSeparator)+skipDirs[i]+string(os.PathSeparator)) {
				return nil
			}
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), goFileExtension) {
			dir, _ := filepath.Split(p)
			absDir, absErr := filepath.Abs(dir)
			if absErr != nil {
				return fmt.Errorf("filepath.Abs: %w", absErr)
			}
			uniqueDirs[absDir] = true
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("filepath.WalkDir: %w", err)
	}

	dirs := make([]string, 0, len(uniqueDirs))
	for dir := range uniqueDirs {
		dirs = append(dirs, dir)
	}

	return dirs, nil
}
