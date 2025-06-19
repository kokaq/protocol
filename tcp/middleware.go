package tcp

import (
	"net"
	"sync"
)

type Middleware interface {
	Handle(conn net.Conn, request *Request)
}

type MiddlewareRegistry struct {
	Middlewares []Middleware
	Mu          sync.RWMutex
}

func NewMiddlewareRegistry() *MiddlewareRegistry {
	return &MiddlewareRegistry{
		Middlewares: make([]Middleware, 0),
		Mu:          sync.RWMutex{},
	}
}

func (registry *MiddlewareRegistry) RegisterMiddleware(middleware Middleware) {
	registry.Mu.Lock()
	defer registry.Mu.Unlock()
	if len(registry.Middlewares) == 0 {
		registry.Middlewares = make([]Middleware, 1)
		registry.Middlewares[0] = middleware
	} else {
		registry.Middlewares = append(registry.Middlewares, middleware)
	}
}
