package writer

import (
	"bufio"
	"os"
)

func GetWriter(path *string) (*bufio.Writer, func(), error) {
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
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
