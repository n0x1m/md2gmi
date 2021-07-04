package main

type Node interface {
	Pipeline(chan []byte) chan []byte
}

type Source interface {
	Output() chan []byte
}

type Sink interface {
	Input(chan []byte)
}
