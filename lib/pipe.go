package lib

import (
	"io"
)

type Chunk struct {
	Data []byte
}

type Conduit struct {
	In  chan *Chunk
	Out chan *Chunk
}

func (s *Conduit) Pipe(piping Piping) {
	piping.Pipe(s)
}

type Piping interface {
	Pipe(*Conduit)
}

type Connector struct {
	Conduits []*Conduit
	Pipes    []Piping
	In       chan<- *Chunk
	Out      <-chan *Chunk
}

func NewConnector(pipes []Piping) *Connector {
	connector := &Connector{Conduits: make([]*Conduit, len(pipes)), Pipes: pipes}
	ch := make(chan *Chunk)
	connector.In = ch

	for i := range pipes {
		next := make(chan *Chunk)
		connector.Conduits[i] = &Conduit{In: ch, Out: next}
		ch = next
	}

	connector.Out = ch

	return connector
}

func (c *Connector) Connect(data []byte, w io.Writer) (int, error) {
	if len(c.Pipes) == 0 {
		return w.Write(data)
	}

	go func() {
		defer close(c.In)
		c.In <- &Chunk{Data: data}
	}()

	for i, piping := range c.Pipes {
		go c.Conduits[i].Pipe(piping)
	}

	t := 0
	for chunk := range c.Out {
		n, err := w.Write(chunk.Data)
		t += n

		if err != nil {
			return t, err
		}
	}

	return t, nil
}
