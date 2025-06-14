package tcp

type Handler interface {
	HandleRequest(*KokaqWireRequest) (*KokaqWireResponse, error)
}
