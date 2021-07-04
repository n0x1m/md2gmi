package pipe

type Pipelines []Pipeline

type Stream struct {
	nodes []Pipeline
}

func New() *Stream {
	return &Stream{}
}

// Use appends a pipeline processor to the Stream pipeline stack.
func (s *Stream) Use(nodes ...Pipeline) {
	s.nodes = append(s.nodes, nodes...)
}

func Chain(middlewares ...Pipeline) Pipelines {
	return Pipelines(middlewares)
}

// chain builds a Connector composed of an inline pipeline stack and endpoint
// processor in the order they are passed.
func chain(nodes []Pipeline, src Connector) Connector {
	c := nodes[0](src)
	for i := 1; i < len(nodes); i++ {
		c = nodes[i](c)
	}

	return c
}

// Handle registers a source and maps it to a sink.
func (s *Stream) Handle(src Source, dest Sink) {
	dest(chain(s.nodes, src()))
}
