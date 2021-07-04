package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/n0x1m/md2gmi/pipe"
)

/*
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
*/

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

func source(r io.Reader) pipe.Source {
	return func() chan pipe.StreamItem {
		data := make(chan pipe.StreamItem)
		s := bufio.NewScanner(r)
		go func() {
			i := 0
			for s.Scan() {
				data <- pipe.NewItem(i, s.Bytes())
				i += 1
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

	preproc := NewPreproc()

	//sink.Input(preproc.Process(source.Output()))
	s := pipe.New()
	s.Use(preproc.Process)
	s.Use(RemoveFrontMatter)
	s.Use(RemoveComments)
	s.Use(FormatHeadings)
	s.Use(FormatLinks)
	s.Handle(source(r), sink(w))
}
