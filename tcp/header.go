package tcp

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ------------------------Common Header------------------------

type CommonHeader struct {
	Magic   uint16
	Version uint8
	MsgType uint8
	RQ      uint8
	Opaque1 uint32
}

func (header *CommonHeader) ToBytes() []byte {
	buffer := make([]byte, 12)
	binary.BigEndian.PutUint16(buffer[0:], Magic)
	buffer[2] = Version
	buffer[3] = header.MsgType<<2 | header.RQ
	binary.BigEndian.PutUint32(buffer[4:], header.Opaque1)
	return buffer
}

func CommonHeaderFromStream(stream io.Reader) (*CommonHeader, error) {
	data := make([]byte, 12)
	if _, err := stream.Read(data); err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}
	return CommonHeaderFromBytes(data)
}

func CommonHeaderFromBytes(data []byte) (*CommonHeader, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("invalid common header length")
	}
	header := &CommonHeader{}
	header.Magic = binary.BigEndian.Uint16(data[0:2])
	header.Version = data[2]
	header.MsgType = data[3] >> 4
	header.RQ = data[3] & 0x0F
	header.Opaque1 = binary.BigEndian.Uint32(data[4:8])
	return header, nil
}

// ------------------------Request Header------------------------

type RequestHeader struct {
	Opcode   uint8
	ClientId uint8
	Opaque2  uint16
}

func (header *RequestHeader) ToBytes() []byte {
	buffer := make([]byte, 4)
	buffer[0] = header.Opcode
	buffer[1] = header.ClientId
	binary.BigEndian.PutUint16(buffer[2:], header.Opaque2)
	return buffer
}

func RequestHeaderFromStream(stream io.Reader) (*RequestHeader, error) {
	data := make([]byte, 4)
	if _, err := stream.Read(data); err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}
	return RequestHeaderFromBytes(data)
}

func RequestHeaderFromBytes(data []byte) (*RequestHeader, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid request operational header length")
	}
	header := &RequestHeader{}
	header.Opcode = data[0]
	header.ClientId = data[1]
	header.Opaque2 = binary.BigEndian.Uint16(data[2:])
	return header, nil
}

// ------------------------Response Header------------------------

type ResponseHeader struct {
	Opcode  uint8
	Status  uint8
	Opaque2 uint16
}

func (header *ResponseHeader) ToBytes() []byte {
	buffer := make([]byte, 4)
	buffer[0] = header.Opcode
	buffer[1] = header.Status
	binary.BigEndian.PutUint16(buffer[2:], header.Opaque2)
	return buffer
}

func ResponseHeaderFromStream(stream io.Reader) (*ResponseHeader, error) {
	data := make([]byte, 4)
	if _, err := stream.Read(data); err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}
	return ResponseHeaderFromBytes(data)
}

func ResponseHeaderFromBytes(data []byte) (*ResponseHeader, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid response operational header length")
	}
	header := &ResponseHeader{}
	header.Opcode = data[0]
	header.Status = data[1]
	header.Opaque2 = binary.BigEndian.Uint16(data[2:])
	return header, nil
}
