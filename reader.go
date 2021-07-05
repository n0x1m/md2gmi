package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/n0x1m/md2gmi/pipe"
)

func reader(in string) (io.Reader, error) {
	if in != "" {
		file, err := os.Open(in)
		if err != nil {
			return nil, fmt.Errorf("reader: %w", err)
		}

		return file, nil
	}

	return os.Stdin, nil
}

func source(r io.Reader) pipe.Source {
	return func() chan pipe.StreamItem {
		data := make(chan pipe.StreamItem)
		s := bufio.NewScanner(r)

		go func() {
			i := 0

			for s.Scan() {
				data <- pipe.NewItem(i, s.Bytes())
				i++
			}
			close(data)
		}()

		return data
	}
}
