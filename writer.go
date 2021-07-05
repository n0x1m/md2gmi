package main

import (
	"fmt"
	"io"
	"os"

	"github.com/n0x1m/md2gmi/pipe"
)

func writer(out string) (io.Writer, error) {
	if out != "" {
		file, err := os.Create(out)
		if err != nil {
			return nil, fmt.Errorf("writer: %w", err)
		}

		return file, nil
	}

	return os.Stdout, nil
}

func sink(w io.Writer) pipe.Sink {
	return func(dest chan pipe.StreamItem) {
		for b := range dest {
			fmt.Fprint(w, string(b.Payload()))
		}
	}
}
