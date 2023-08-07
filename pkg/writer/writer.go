package writer

import (
	"bufio"
	"errors"
	"os"
)

var ErrMissingResult = errors.New("result is required")

func GetWriter(path *string) (*bufio.Writer, func(), error) {
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		if path == nil || *path == "" {
			return nil, nil, ErrMissingResult
		}
		file, err := os.Create(*path)
		if err != nil {
			return nil, nil, err
		}
		writer := bufio.NewWriter(file)
		return writer, func() {
			_ = writer.Flush()
			_ = file.Close()
		}, nil
	}
	writer := bufio.NewWriter(os.Stdout)
	return writer, func() {
		_ = writer.Flush()
	}, nil
}
