package writer

import (
	"bufio"
	"os"
)

func GetWriter(path *string) (*bufio.Writer, func(), error) {
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
