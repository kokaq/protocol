package tcp

import (
	"sync"
)

type Handler interface {
	Handle(request *Request) (*Response, error)
}

func NewMultiplexer() *Multiplexer {
	return &Multiplexer{
		Routes: make(map[uint8]Handler),
		Mu:     sync.RWMutex{},
	}
}

type Multiplexer struct {
	Routes map[uint8]Handler // keyed by Opcode
	Mu     sync.RWMutex
}

func (mux *Multiplexer) RegisterHandler(opcode uint8, handler Handler) {
	mux.Mu.Lock()
	defer mux.Mu.Unlock()
	mux.Routes[opcode] = handler
}
