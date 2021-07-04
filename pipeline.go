package main

type Node interface {
	Pipeline(<-chan []byte) <-chan []byte
}

type Source interface {
	Input() <-chan []byte
}

type Sink interface {
	Output(<-chan []byte)
}
