package main

// state function
type stateFn func(*fsm, []byte) stateFn

// state machine
type fsm struct {
	state stateFn
	out   chan []byte

	// combining multiple input lines
	buffer []byte
	// if we have a termination rule to abide, e.g. implied code fences
	pending []byte
}

func NewParser() *fsm {
	return &fsm{}
}

func (m *fsm) Parse(in chan []byte) chan []byte {
	m.out = make(chan []byte)
	go func() {
		for m.state = normal; m.state != nil; {
			b, ok := <-in
			if !ok {
				m.flush()
				close(m.out)
				m.state = nil
				continue
			}
			m.state = m.state(m, b)
		}
	}()
	return m.out
}

func (m *fsm) flush() {
	if len(m.pending) > 0 {
		m.out <- append(m.pending, '\n')
		m.pending = m.pending[:0]
	}
}

func isBlank(data []byte) bool {
	return len(data) == 0
}

func isHeader(data []byte) bool {
	return len(data) > 0 && data[0] == '#'
}

func triggerBreak(data []byte) bool {
	return len(data) == 0 || data[len(data)-1] == '.'
}

func isFence(data []byte) bool {
	return len(data) >= 3 && string(data[0:3]) == "```"
}

func needsFence(data []byte) bool {
	return len(data) >= 4 && string(data[0:4]) == "    "
}

func normal(m *fsm, data []byte) stateFn {
	m.flush()
	// blank line
	if isBlank(data) {
		m.out <- []byte("\n")
		return normal
	}
	// header
	if isHeader(data) {
		m.out <- append(data, '\n')
		return normal
	}
	if isFence(data) {
		m.out <- append(data, '\n')
		return fence
	}
	if needsFence(data) {
		m.out <- []byte("```\n")
		m.out <- append(data, '\n')
		m.pending = []byte("```")
		return toFence
	}
	if data[len(data)-1] != '.' {
		m.buffer = append(m.buffer, data...)
		m.buffer = append(m.buffer, []byte(" ")...)
		return paragraph
	}
	// TODO
	// find links
	// collapse lists
	m.out <- append(data, '\n')

	return normal
}

func fence(m *fsm, data []byte) stateFn {
	m.out <- append(data, '\n')
	if isFence(data) {
		return normal
	}
	return fence
}

func toFence(m *fsm, data []byte) stateFn {
	m.out <- append(data, '\n')
	if needsFence(data) {
		return toFence
	}
	return normal
}

func paragraph(m *fsm, data []byte) stateFn {
	if triggerBreak(data) {
		m.buffer = append(m.buffer, data...)
		m.out <- append(m.buffer, '\n')
		m.buffer = m.buffer[:0]
		return normal
	}
	m.buffer = append(m.buffer, data...)
	m.buffer = append(m.buffer, []byte(" ")...)
	return paragraph
}

func link(m *fsm, data []byte) stateFn {

	return link
}
