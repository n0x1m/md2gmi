package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func reader(in string) (io.Reader, error) {
	if in != "" {
		file, err := os.Open(in)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	return os.Stdin, nil
}

func writer(out string) (io.Writer, error) {
	if out != "" {
		file, err := os.Open(out)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	return os.Stdout, nil
}

func write(w io.Writer, b []byte) {
	fmt.Fprint(w, string(b))
}

type ir struct {
	r io.Reader
}

func NewIr(r io.Reader) *ir {
	return &ir{r: r}
}

func (m *ir) Output() chan []byte {
	data := make(chan []byte)
	s := bufio.NewScanner(m.r)
	go func() {
		for s.Scan() {
			data <- s.Bytes()
		}
		close(data)
	}()
	return data
}

type ow struct {
	w io.Writer
}

func NewOw(w io.Writer) *ow {
	return &ow{w: w}
}

func (m *ow) Input(data chan []byte) {
	for b := range data {
		write(m.w, b)
	}
}

func main() {
	var in, out string

	flag.StringVar(&in, "in", "", "specify a .md (Markdown) file to read from, otherwise stdin (default)")
	flag.StringVar(&out, "out", "", "specify a .gmi (gemtext) file to write to, otherwise stdout (default)")
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

	source := NewIr(r)
	sink := NewOw(w)
	preproc := NewParser()

	sink.Input(preproc.Parse(source.Output()))
}
