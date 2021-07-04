package pipe

import (
	"bytes"
	"encoding/gob"
)

type StreamItem struct {
	index   int
	payload []byte
}

func NewItem(index int, payload []byte) StreamItem {
	var buf bytes.Buffer
	w := StreamItem{index: index}
	if err := gob.NewEncoder(&buf).Encode(payload); err != nil {
		// assert no broken pipes
		panic(err)
	}
	w.payload = buf.Bytes()
	return w
}

func (w *StreamItem) Index() int {
	return w.index
}

func (w *StreamItem) Payload() []byte {
	var dec []byte
	buf := bytes.NewReader(w.payload)
	if err := gob.NewDecoder(buf).Decode(&dec); err != nil {
		// assert no broken pipes
		panic(err)
	}

	return dec
}