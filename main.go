package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/n0x1m/md2gmi/mdproc"
	"github.com/n0x1m/md2gmi/pipe"
)

func main() {
	var in, out string

	flag.StringVar(&in, "i", "", "specify a .md (Markdown) file to read from, otherwise stdin (default)")
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
	s.Use(mdproc.Preprocessor())
	s.Use(mdproc.RemoveFrontMatter)
	s.Use(mdproc.FormatHeadings)
	s.Use(mdproc.FormatLinks)
	s.Handle(source(r), sink(w))
}
