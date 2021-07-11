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

<!-- a <!-- double -->
comment -->

<!--
multi line comment block, terminated.

With a break line, <!-- and multi line
 inline comments -->
and [with a](link);
-->
> this is
a quote

*not a list*<!-- comment
between two things -->
**also not a list**

 - this
	- is
* an unordered
  * list

This is
a paragraph with a link to the [gemini protocol](https://en.wikipedia.org/wiki/Gemini_(protocol)).

` + "```" + `
this is multi
line code
` + "```" + `

and<!-- inline comment-->

    this is code too`

	preproc = `--- title: "This is the Title!" categories: [a,b] ---

> this is a quote

*not a list*

**also not a list**

- this
- is
- an unordered
- list

This is a paragraph with a link to the [gemini protocol](https://en.wikipedia.org/wiki/Gemini_(protocol)).

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

*not a list*

**also not a list**

- this
- is
- an unordered
- list

This is a paragraph with a link to the gemini protocol[1].

=> https://en.wikipedia.org/wiki/Gemini_(protocol) 1: gemini protocol

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

func source(t *testing.T, in string) func() chan pipe.StreamItem {
	t.Helper()

	return func() chan pipe.StreamItem {
		data := make(chan pipe.StreamItem, len(strings.Split(in, "\n")))
		for _, line := range strings.Split(in, "\n") {
			data <- pipe.NewItem(0, []byte(line))
		}

		close(data)

		return data
	}
}

func sink(t *testing.T, expected string) func(dest chan pipe.StreamItem) {
	t.Helper()

	return func(dest chan pipe.StreamItem) {
		var data []byte
		for in := range dest {
			data = append(data, in.Payload()...)
		}

		if string(data) != expected {
			t.Errorf("mismatch, expected '%s' but was '%s'", expected, data)
		}
	}
}

func TestPreproc(t *testing.T) {
	t.Parallel()

	s := pipe.New()
	s.Use(mdproc.Preprocessor())
	s.Handle(source(t, input), sink(t, preproc))
}

func TestMd2Gmi(t *testing.T) {
	t.Parallel()

	s := pipe.New()
	s.Use(mdproc.Preprocessor())
	s.Use(mdproc.RemoveFrontMatter)
	s.Use(mdproc.FormatHeadings)
	s.Use(mdproc.FormatLinks)
	s.Handle(source(t, input), sink(t, gmi))
}
