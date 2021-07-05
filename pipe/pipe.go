package pipe

type Connector chan StreamItem

type Source func() chan StreamItem
type Sink func(chan StreamItem)
type Pipeline func(chan StreamItem) chan StreamItem
