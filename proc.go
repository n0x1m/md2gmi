package main

import (
	"bytes"
	"fmt"
	"regexp"
)

// state machine
type proc struct {
	out chan []byte
}

func NewProc() *proc {
	return &proc{}
}

func (m *proc) Process(in chan []byte) chan []byte {
	m.out = make(chan []byte)
	go func() {
		for b := range in {
			m.out <- m.process(b)
		}
		close(m.out)
	}()
	return m.out
}

func (m *proc) process(data []byte) []byte {
	// find link name and url
	var buffer []byte
	re := regexp.MustCompile(`\[([^\]*]*)\]\(([^)]*)\)`)
	for i, match := range re.FindAllSubmatch(data, -1) {
		replaceWithIndex := append(match[1], fmt.Sprintf("[%d]", i+1)...)
		data = bytes.Replace(data, match[0], replaceWithIndex, 1)
		link := fmt.Sprintf("=> %s %d: %s\n", match[2], i+1, match[1])
		buffer = append(buffer, link...)
	}
	if len(buffer) > 0 {
		data = append(data, []byte("\n")...)
		data = append(data, buffer...)
	}

	// remove comments
	re2 := regexp.MustCompile(`<!--.*-->`)
	data = re2.ReplaceAll(data, []byte{})

	// collapse headings
	re3 := regexp.MustCompile(`^[#]{4,}`)
	data = re3.ReplaceAll(data, []byte("###"))

	// heading without spacing
	re4 := regexp.MustCompile(`^(#+)[^# ]`)
	sub := re4.FindSubmatch(data)
	if len(sub) > 0 {
		data = bytes.Replace(data, sub[1], append(sub[1], []byte(" ")...), 1)
	}

	return data
}
