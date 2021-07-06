package pipe

import (
	"context"
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
	return newItem(context.Background(), index, payload)
}

func NewItemWithContext(ctx context.Context, index int, payload []byte) StreamItem {
	return newItem(ctx, index, payload)
}

func newItem(ctx context.Context, index int, payload []byte) StreamItem {
	s := StreamItem{
		ctx:     ctx,
		index:   index,
		payload: make([]byte, len(payload)),
	}
	copy(s.payload, payload)
	return s
}

func (s *StreamItem) Index() int {
	return s.index
}

func (s *StreamItem) Payload() []byte {
	return s.payload
}
