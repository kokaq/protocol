package tcp

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type KokaqWireClient struct {
	address string
	timeout int
}

func NewKokaqWireClientFromHostPort(host string, port int, timeout int) *KokaqWireClient {
	return &KokaqWireClient{
		address: fmt.Sprintf("%s:%d", host, port),
		timeout: timeout,
	}
}

func NewKokaqWireClientFromAddress(address string, timeout int) *KokaqWireClient {
	return &KokaqWireClient{
		address: address,
		timeout: timeout,
	}
}

func (client *KokaqWireClient) SendToWire(request *KokaqWireRequest) (*KokaqWireResponse, error) {
	conn, err := net.DialTimeout("tcp", client.address, time.Duration(client.timeout)*time.Second)
	if err != nil {
		fmt.Printf("Socket Error: %v\n", err)
		return nil, err
	}
	defer conn.Close()

	//stream := conn
	bufWriter := bufio.NewWriter(conn)
	bufReader := bufio.NewReader(conn)

	// Write the request to the stream
	err = request.WriteToStream(bufWriter)
	bufWriter.Flush() // Ensure data is sent
	if err != nil {
		fmt.Printf("Error writing request: %v\n", err)
		return nil, err
	}

	fmt.Println("Request sent, awaiting response...")

	// Read the response from the stream
	response := &KokaqWireResponse{}
	err = response.ReadFromStream(bufReader)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return nil, err
	}

	return response, nil
}
