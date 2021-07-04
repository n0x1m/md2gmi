package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"os"
)

type WorkItem struct {
	index   int
	payload []byte
}

func New(index int, payload []byte) WorkItem {
	w := WorkItem{index: index}
	var indexBuffer bytes.Buffer
	encoder := gob.NewEncoder(&indexBuffer)
	if err := encoder.Encode(payload); err != nil {
		panic(err)
	}
	w.payload = indexBuffer.Bytes()
	return w
}

func (w *WorkItem) Index() int {
	return w.index
}

func (w *WorkItem) Payload() []byte {
	buf := bytes.NewReader(w.payload)
	decoder := gob.NewDecoder(buf)
	var tmp []byte
	if err := decoder.Decode(&tmp); err != nil {
		panic(err)
	}
	return tmp
}

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
		file, err := os.Create(out)
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

func InputStream(r io.Reader) *ir {
	return &ir{r: r}
}

func (m *ir) Output() chan WorkItem {
	data := make(chan WorkItem)
	s := bufio.NewScanner(m.r)
	go func() {
		i := 0
		for s.Scan() {
			data <- New(i, s.Bytes())
			i += 1
		}
		close(data)
	}()
	return data
}

type ow struct {
	w io.Writer
}

func OutputStream(w io.Writer) *ow {
	return &ow{w: w}
}

func (m *ow) Input(data chan WorkItem) {
	for b := range data {
		write(m.w, b.Payload())
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

	source := InputStream(r)
	sink := OutputStream(w)
	preproc := NewPreproc()

	//sink.Input(preproc.Process(source.Output()))
	sink.Input(
		FormatLinks(
			FormatHeadings(
				RemoveComments(
					RemoveFrontMatter(
						preproc.Process(source.Output()),
					),
				),
			),
		),
	)
}
