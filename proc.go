package main

import (
	"bytes"
	"fmt"
	"regexp"
)

func FormatLinks(in chan []byte) chan []byte {
	out := make(chan []byte)
	go func() {
		for b := range in {
			out <- formatLinks(b)
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

func RemoveComments(in chan []byte) chan []byte {
	out := make(chan []byte)
	go func() {
		re := regexp.MustCompile(`<!--.*-->`)
		for b := range in {
			out <- re.ReplaceAll(b, []byte{})
		}
		close(out)
	}()
	return out
}

func FormatHeadings(in chan []byte) chan []byte {
	out := make(chan []byte)
	go func() {
		re := regexp.MustCompile(`^[#]{4,}`)
		re2 := regexp.MustCompile(`^(#+)[^# ]`)
		for b := range in {
			// fix up more than 4 levels
			b = re.ReplaceAll(b, []byte("###"))
			// ensure we have a space
			sub := re2.FindSubmatch(b)
			if len(sub) > 0 {
				b = bytes.Replace(b, sub[1], append(sub[1], []byte(" ")...), 1)
			}
			// writeback
			out <- b

		}
		close(out)
	}()
	return out
}
