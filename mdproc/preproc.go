package mdproc

import (
	"bytes"
	"regexp"

	"github.com/n0x1m/md2gmi/pipe"
)

// state function.
type stateFn func(*fsm, []byte) stateFn

// state machine.
type fsm struct {
	state stateFn

	i   int
	out chan pipe.StreamItem

	// combining multiple input lines
	blockBuffer []byte
	sendBuffer  []byte
	// if we have a termination rule to abide, e.g. implied code fences
	pending []byte
}

func Preproc() pipe.Pipeline {
	return (&fsm{}).pipeline
}

func (m *fsm) pipeline(in chan pipe.StreamItem) chan pipe.StreamItem {
	m.out = make(chan pipe.StreamItem)

	go func() {
		for m.state = normal; m.state != nil; {
			b, ok := <-in
			if !ok {
				m.blockFlush()
				m.sync()
				close(m.out)
				m.state = nil

				continue
			}

			m.state = m.state(m, b.Payload())
			m.sync()
		}
	}()

	return m.out
}

func (m *fsm) sync() {
	if len(m.sendBuffer) > 0 {
		m.sendBuffer = append(m.sendBuffer, '\n')
		m.out <- pipe.NewItem(m.i, m.sendBuffer)
		m.sendBuffer = m.sendBuffer[:0]
		m.i++
	}
}

func (m *fsm) blockFlush() {
	// blockBuffer to sendbuffer
	m.sendBuffer = append(m.sendBuffer, m.blockBuffer...)
	m.blockBuffer = m.blockBuffer[:0]

	// pending to sendbuffer too
	if len(m.pending) > 0 {
		m.sendBuffer = append(m.sendBuffer, m.pending...)
		m.pending = m.pending[:0]
	}
}

func triggerBreak(data []byte) bool {
	return len(data) == 0 || data[len(data)-1] == '.'
}

func isTerminated(data []byte) bool {
	return len(data) > 0 && data[len(data)-1] != '.'
}

func handleList(data []byte) ([]byte, bool) {
	re := regexp.MustCompile(`^([ \t]*[-*^]{1,1})[^*-]`)
	sub := re.FindSubmatch(data)
	// if lists, collapse to single level
	if len(sub) > 1 {
		return bytes.Replace(data, sub[1], []byte("-"), 1), true
	}

	return data, false
}

func isFence(data []byte) bool {
	return len(data) >= 3 && string(data[0:3]) == "```"
}

func needsFence(data []byte) bool {
	return len(data) >= 4 && string(data[0:4]) == "    "
}

func normal(m *fsm, data []byte) stateFn {
	if len(bytes.TrimSpace(data)) == 0 {
		return normal
	}
	if data, isList := handleList(data); isList {
		//m.blockBuffer = append(data, '\n')
		m.blockBuffer = append(m.blockBuffer, data...)
		m.blockFlush()

		return list
	}

	if isFence(data) {
		m.blockBuffer = append(data, '\n')

		return fence
	}

	if needsFence(data) {
		m.blockBuffer = append(m.blockBuffer, []byte("```\n")...)
		m.blockBuffer = append(m.blockBuffer, append(data[4:], '\n')...)
		m.pending = []byte("```\n")

		return toFence
	}

	if isTerminated(data) {
		m.blockBuffer = append(m.blockBuffer, data...)
		m.blockBuffer = append(m.blockBuffer, ' ')

		return paragraph
	}

	m.blockBuffer = append(m.blockBuffer, append(data, '\n')...)
	m.blockFlush()

	return normal
}

func list(m *fsm, data []byte) stateFn {
	if data, isList := handleList(data); isList {
		m.blockBuffer = append(m.blockBuffer, data...)
		m.blockFlush()

		return list
	}

	m.blockBuffer = append(m.blockBuffer, append(data, '\n')...)
	m.blockFlush()

	return normal
}

func fence(m *fsm, data []byte) stateFn {
	m.blockBuffer = append(m.blockBuffer, append(data, '\n')...)
	// second fence returns to normal
	if isFence(data) {
		m.blockFlush()

		return normal
	}

	return fence
}

func toFence(m *fsm, data []byte) stateFn {
	if needsFence(data) {
		m.blockBuffer = append(m.blockBuffer, append(data[4:], '\n')...)

		return toFence
	}

	m.blockFlush()
	m.blockBuffer = append(m.blockBuffer, append(data, '\n')...)

	return normal
}

func paragraph(m *fsm, data []byte) stateFn {
	if triggerBreak(data) {
		m.blockBuffer = append(m.blockBuffer, data...)
		m.blockBuffer = bytes.TrimSpace(m.blockBuffer)
		// TODO, remove double spaces inside paragraphs
		m.blockBuffer = append(m.blockBuffer, '\n')
		m.blockFlush()

		return normal
	}

	m.blockBuffer = append(m.blockBuffer, data...)
	m.blockBuffer = append(m.blockBuffer, []byte(" ")...)

	return paragraph
}
