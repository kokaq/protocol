package tcp

import (
	"fmt"
	"net"
	"sync"
)

type KokaqWireServer struct {
	port     int
	listener net.Listener
	wg       sync.WaitGroup
	timeout  int
}

type IKokaqServerhandler interface {
	HandleRequest(*KokaqWireRequest) (*KokaqWireResponse, error)
}

// NewKokaqWireServer initializes a new KokaqWireServer with the provided port.
func NewKokaqWireServer(port int, timeout int) *KokaqWireServer {
	return &KokaqWireServer{
		port:    port,
		timeout: timeout,
	}
}

// Start begins the TCP listener and handles incoming client requests.
func (server *KokaqWireServer) Start(handler IKokaqServerhandler) error {
	address := fmt.Sprintf(":%d", server.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return err
	}
	server.listener = listener
	defer server.listener.Close()

	fmt.Printf("Listening on port %d...\n", server.port)

	for {
		fmt.Println("Waiting for a connection...")
		client, err := server.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		fmt.Println("Receiving a request...")

		// Handle each client in a separate goroutine
		server.wg.Add(1)
		go func(client net.Conn) {
			fmt.Println("Received a request")
			defer server.wg.Done()
			defer client.Close()
			//client.SetDeadline(time.Now().Add(time.Second * time.Duration(server.timeout)))

			var req = &KokaqWireRequest{}

			err := req.ReadFromStream(client)
			if err != nil {
				fmt.Printf("Error reading request: %v\n", err)
				return
			}

			switch req.MessageType {
			case MessageTypeOperational:
				fmt.Println("Received an operational message.")
			case MessageTypeControl:
				fmt.Println("Received a control message.")
			case MessageTypeAdmin:
				fmt.Println("Received a adminops message.")
			default:
				fmt.Println("Received an unknown message type.")
				return
			}

			switch req.OpCode {
			case OpCodeCreate:
				fmt.Println("Received a create request.")
			case OpCodeDelete:
				fmt.Println("Received a delete request.")
			case OpCodeGet:
				fmt.Println("Received a get request.")
			case OpCodePeek:
				fmt.Println("Received a peek request.")
			case OpCodePop:
				fmt.Println("Received a pop request.")
			case OpCodePush:
				fmt.Println("Received a push request.")
			case OpCodeAcquirePeekLock:
				fmt.Println("Received an acquire peek lock request.")
			default:
				fmt.Println("Received an unknown operation.")
				return
			}

			res, err := handler.HandleRequest(req)
			err = res.WriteToStream(client)
			if err != nil {
				fmt.Printf("Error writing response: %v\n", err)
				return
			}
			fmt.Println("Done handling the request.")
		}(client)
	}
}

// Stop stops the server and waits for all connections to be handled.
func (server *KokaqWireServer) Stop() {
	server.listener.Close()
	server.wg.Wait()
}
