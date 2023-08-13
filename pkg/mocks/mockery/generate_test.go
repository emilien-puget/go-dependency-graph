package mockery

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/emilien-puget/go-dependency-graph/pkg/mocks/config"
	"github.com/emilien-puget/go-dependency-graph/pkg/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFromSchema_outofpackage(t *testing.T) {
	t.Parallel()
	as, err := parse.Parse("testdata/named_inter")
	require.NoError(t, err)

	dir := t.TempDir()
	err = GenerateFromSchema(
		config.Config{
			InPackage:                  false,
			OutOfPackageMocksDirectory: dir,
		},
		as,
	)
	require.NoError(t, err)

	assertDirectoriesEqual(t, dir, "testdata/expect_named_inter/out_of_package")
}

func TestGenerateFromSchema_inpackage(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	err := copyDir("testdata/named_inter", tempDir)
	require.NoError(t, err)
	as, err := parse.Parse(tempDir)
	require.NoError(t, err)

	err = GenerateFromSchema(
		config.Config{
			InPackage: true,
		},
		as,
	)
	require.NoError(t, err)

	assertDirectoriesEqual(t, tempDir, "testdata/expect_named_inter/in_package")
}

func assertDirectoriesEqual(t *testing.T, gotDir, expectDir string) {
	t.Helper()

	err := filepath.Walk(gotDir, func(gotPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the corresponding path in the second directory
		relativePath, _ := filepath.Rel(gotDir, gotPath)
		expectPath := filepath.Join(expectDir, relativePath)

		// Check if the file exists in expectPath
		_, err = os.Stat(expectPath)
		if err != nil {
			return fmt.Errorf("got file %s not expected:%w", relativePath, err)
		}
		// Handle subdirectories
		if info.IsDir() {
			return nil // Continue walking
		}

		// Compare file contents line by line
		assertFilesEqualLineByLine(t, gotPath, expectPath)

		return nil
	})
	if err != nil {
		t.Errorf("Error walking directories: %s", err)
	}
}

func assertFilesEqualLineByLine(t *testing.T, filePath1, filePath2 string) {
	t.Helper()

	file1, err := os.Open(filePath1)
	if err != nil {
		t.Fatalf("Error opening file: %s", err)
	}
	defer file1.Close()

	file2, err := os.Open(filePath2)
	if err != nil {
		t.Fatalf("Error opening file: %s", err)
	}
	defer file2.Close()

	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)

	lineNumber := 1
	for scanner1.Scan() {
		assert.True(t, scanner2.Scan(), "File contents do not match at line %d", lineNumber)

		line1 := scanner1.Text()
		line2 := scanner2.Text()

		assert.Equal(t, line1, line2, "File contents do not match at line %d", lineNumber)

		lineNumber++
	}

	assert.False(t, scanner2.Scan(), "File contents do not match at line %d", lineNumber)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
