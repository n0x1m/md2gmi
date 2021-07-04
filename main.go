package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func read(in string, data chan<- []byte) error {
	if in != "" {
		file, err := os.Open(in)
		if err != nil {
			return err
		}
		defer file.Close()

		s := bufio.NewScanner(file)
		for s.Scan() {
			data <- s.Bytes()
		}

		return nil
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		data <- s.Bytes()
	}

	return nil
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
	fmt.Fprintf(w, string(b))
}

type ow struct {
	w    io.Writer
	quit chan struct{}
}

func NewOw(w io.Writer, quit chan struct{}) *ow {
	return &ow{w: w, quit: quit}
}

func (m *ow) Input(data <-chan []byte) {
	for {
		select {
		case <-m.quit:
			return
		case b := <-data:
			write(m.w, b)
		}
	}
}

func writesetup(out string, data <-chan []byte, quit chan struct{}) error {
	w, err := writer(out)
	if err != nil {
		return fmt.Errorf("writer: %w", err)
	}

	sink := NewOw(w, quit)
	sink.Input(data)

	//NewParser(data, w, quit).Parse()

	return nil
}

func main() {
	var in, out string

	flag.StringVar(&in, "in", "", "specify a .md (Markdown) file to read from, otherwise stdin (default)")
	flag.StringVar(&out, "out", "", "specify a .gmi (gemtext) file to write to, otherwise stdout (default)")
	flag.Parse()

	quit := make(chan struct{})
	data := make(chan []byte)

	go func() {
		if err := read(in, data); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		close(quit)
	}()

	if err := writesetup(out, data, quit); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}
