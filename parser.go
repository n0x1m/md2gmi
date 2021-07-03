package main

import (
	"fmt"
	"io"
)

// state function
type stateFn func(*fsm, []byte) stateFn

// state machine
type fsm struct {
	state stateFn
	out   io.Writer
	quit  chan struct{}
	data  <-chan []byte
}

func NewParser(data <-chan []byte, writer io.Writer, quit chan struct{}) *fsm {
	return &fsm{
		out:  writer,
		data: data,
		quit: quit,
	}
}

func (m *fsm) Parse() {
	var buffer []byte
	for m.state = initial; m.state != nil; {
		select {
		case <-m.quit:
			m.state = nil
		case buffer = <-m.data:
			m.state = m.state(m, buffer)
		}
	}
}

func initial(m *fsm, data []byte) stateFn {
	// TODO
	// find linebreaks
	// find code fences
	// find links
	// collapse lists
	fmt.Fprintf(m.out, string(data)+"\n")

	return initial
}
