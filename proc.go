package main

// state function
type stateFn2 func(*proc, []byte) stateFn2

// state machine
type proc struct {
	state stateFn2
	out   chan []byte

	// combining multiple input lines
	buffer []byte
	// if we have a termination rule to abide, e.g. implied code fences
	pending []byte
}

func NewProc() *proc {
	return &proc{}
}

func (m *proc) Process(in chan []byte) chan []byte {
	m.out = make(chan []byte)
	go func() {
		for m.state = line; m.state != nil; {
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

func (m *proc) flush() {
	if len(m.pending) > 0 {
		m.out <- append(m.pending, '\n')
		m.pending = m.pending[:0]
	}
}

func line(m *proc, data []byte) stateFn2 {
	// TODO
	// find links
	// collapse lists
	m.out <- data

	return line
}
