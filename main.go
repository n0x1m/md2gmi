package main

import (
	"bufio"
	"flag"
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

func write(w io.Writer, b []byte) {
	fmt.Fprint(w, string(b))
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

func sink(w io.Writer) pipe.Sink {
	return func(dest chan pipe.StreamItem) {
		for b := range dest {
			write(w, b.Payload())
		}
	}
}

func main() {
	var in, out string

	flag.StringVar(&in, "f", "", "specify a .md (Markdown) file to read from, otherwise stdin (default)")
	flag.StringVar(&out, "o", "", "specify a .gmi (gemtext) file to write to, otherwise stdout (default)")
	flag.Parse()

	r, err := reader(in)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	w, err := writer(out)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	s := pipe.New()
	s.Use(newPreproc().Process)
	s.Use(RemoveFrontMatter)
	s.Use(RemoveComments)
	s.Use(FormatHeadings)
	s.Use(FormatLinks)
	s.Handle(source(r), sink(w))
}
