package main

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/n0x1m/md2gmi/pipe"
)

func FormatLinks(in chan pipe.StreamItem) chan pipe.StreamItem {
	out := make(chan pipe.StreamItem)

	go func() {
		fenceOn := false

		for b := range in {
			data := b.Payload()
			if isFence(data) {
				fenceOn = !fenceOn
			}

			if fenceOn {
				out <- pipe.NewItem(b.Index(), b.Payload())

				continue
			}
			out <- pipe.NewItem(b.Index(), formatLinks(b.Payload()))
		}

		close(out)
	}()

	return out
}

func formatLinks(data []byte) []byte {
	// find link name and url
	var buffer []byte

	re := regexp.MustCompile(`!?\[([^\]*]*)\]\(([^)]*)\)`)

	for i, match := range re.FindAllSubmatch(data, -1) {
		replaceWithIndex := append(match[1], fmt.Sprintf("[%d]", i+1)...)
		data = bytes.Replace(data, match[0], replaceWithIndex, 1)
		// append entry to buffer to be added later
		link := fmt.Sprintf("=> %s %d: %s\n", match[2], i+1, match[1])
		buffer = append(buffer, link...)
	}
	// append links to that paragraph
	if len(buffer) > 0 {
		data = append(data, []byte("\n")...)
		data = append(data, buffer...)
	}

	return data
}

func RemoveComments(in chan pipe.StreamItem) chan pipe.StreamItem {
	out := make(chan pipe.StreamItem)

	go func() {
		re := regexp.MustCompile(`<!--.*-->`)

		for b := range in {
			data := b.Payload()
			for _, match := range re.FindAllSubmatch(data, -1) {
				data = bytes.Replace(data, match[0], []byte(""), 1)
			}
			out <- pipe.NewItem(b.Index(), append(bytes.TrimSpace(data), '\n'))
		}

		close(out)
	}()

	return out
}

func RemoveFrontMatter(in chan pipe.StreamItem) chan pipe.StreamItem {
	out := make(chan pipe.StreamItem)

	go func() {
		re := regexp.MustCompile(`---.*---`)

		for b := range in {
			data := b.Payload()
			for _, match := range re.FindAllSubmatch(data, -1) {
				data = bytes.Replace(data, match[0], []byte(""), 1)
			}
			out <- pipe.NewItem(b.Index(), append(bytes.TrimSpace(data), '\n'))
		}

		close(out)
	}()

	return out
}

func FormatHeadings(in chan pipe.StreamItem) chan pipe.StreamItem {
	out := make(chan pipe.StreamItem)

	go func() {
		re := regexp.MustCompile(`^[#]{4,}`)
		re2 := regexp.MustCompile(`^(#+)[^# ]`)

		for b := range in {
			// fix up more than 4 levels
			data := re.ReplaceAll(b.Payload(), []byte("###"))
			// ensure we have a space
			sub := re2.FindSubmatch(data)
			if len(sub) > 0 {
				data = bytes.Replace(data, sub[1], append(sub[1], []byte(" ")...), 1)
			}
			// generally if we deal with a heading, add an extra blank line
			if bytes.HasPrefix(data, []byte("#")) {
				data = append(data, '\n')
			}
			// writeback
			out <- pipe.NewItem(b.Index(), data)
		}

		close(out)
	}()

	return out
}
