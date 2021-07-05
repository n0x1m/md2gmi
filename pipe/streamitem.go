package pipe

import (
	"bytes"
	"context"
	"encoding/gob"
)

type StreamItem struct {
	ctx     context.Context
	index   int
	payload []byte
}

func (s *StreamItem) Context() context.Context {
	return s.ctx
}

func NewItem(index int, payload []byte) StreamItem {
	var buf bytes.Buffer

	s := StreamItem{index: index}

	if err := gob.NewEncoder(&buf).Encode(payload); err != nil {
		// assert no broken pipes
		panic(err)
	}

	s.payload = buf.Bytes()

	return s
}

func (s *StreamItem) Index() int {
	return s.index
}

func (s *StreamItem) Payload() []byte {
	var dec []byte

	buf := bytes.NewReader(s.payload)

	if err := gob.NewDecoder(buf).Decode(&dec); err != nil {
		// assert no broken pipes
		panic(err)
	}

	return dec
}
