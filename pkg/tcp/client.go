package tcp

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type KokaqWireClient struct {
	address string
	timeout int
	logger  *logrus.Logger
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
		client.logger.WithError(err).Error("Socket Error")
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
		client.logger.WithError(err).Error("Error writing request")
		return nil, err
	}

	client.logger.Info("Request sent, awaiting response...")

	// Read the response from the stream
	response := &KokaqWireResponse{}
	err = response.ReadFromStream(bufReader)
	if err != nil {
		client.logger.WithError(err).Error("Error reading response")
		return nil, err
	}

	return response, nil
}
