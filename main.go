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

func write(out string, data <-chan []byte, quit chan struct{}) error {
	w, err := writer(out)
	if err != nil {
		return fmt.Errorf("writer: %w", err)
	}

	NewParser(data, w, quit).Parse()

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

	if err := write(out, data, quit); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}
