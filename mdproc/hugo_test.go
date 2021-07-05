package mdproc_test

import (
	"strings"
	"testing"

	"github.com/n0x1m/md2gmi/mdproc"
	"github.com/n0x1m/md2gmi/pipe"
)

const (
	input = `---
title: "This is the Title!"
categories: [a,b]
---

<!-- a
comment -->

> this is
a quote

This is
a paragraph.

` + "```" + `
this is multi
line code
` + "```" + `

and

    this is code too
`

	preproc = `--- title: "This is the Title!" categories: [a,b] ---

<!-- a comment -->

> this is a quote

This is a paragraph.

` + "```" + `
this is multi
line code
` + "```" + `

and

` + "```" + `
this is code too
` + "```" + `



`

	gmi = `# This is the Title!

> this is a quote

This is a paragraph.

` + "```" + `
this is multi
line code
` + "```" + `

and

` + "```" + `
this is code too
` + "```" + `



`
)

func TestPreproc(t *testing.T) {
	source := func() chan pipe.StreamItem {
		data := make(chan pipe.StreamItem, len(strings.Split(input, "\n")))
		for _, line := range strings.Split(input, "\n") {
			data <- pipe.NewItem(0, []byte(line))
		}
		close(data)

		return data
	}

	sink := func(dest chan pipe.StreamItem) {
		var data []byte
		for in := range dest {
			data = append(data, in.Payload()...)
		}
		if string(data) != preproc {
			t.Errorf("mismatch, expected '%s' but was '%s'", preproc, data)
		}
	}

	sink(mdproc.Preproc()(source()))
}

func TestMd2Gmi(t *testing.T) {
	source := func() chan pipe.StreamItem {
		data := make(chan pipe.StreamItem, len(strings.Split(input, "\n")))
		for _, line := range strings.Split(input, "\n") {
			data <- pipe.NewItem(0, []byte(line))
		}
		close(data)

		return data
	}

	sink := func(dest chan pipe.StreamItem) {
		var data []byte
		for in := range dest {
			data = append(data, in.Payload()...)
		}
		if string(data) != gmi {
			t.Errorf("mismatch, expected '%s' but was '%s'", gmi, data)
		}
	}

	s := pipe.New()
	s.Use(mdproc.Preproc())
	s.Use(mdproc.RemoveFrontMatter)
	s.Use(mdproc.RemoveComments)
	s.Use(mdproc.FormatHeadings)
	s.Use(mdproc.FormatLinks)
	s.Handle(source, sink)
}
