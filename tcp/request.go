package tcp

import (
	"encoding/binary"
	"io"
	"strconv"
)

type Request struct {
	commonHeader      CommonHeader
	operationalHeader RequestHeader
	namespaceId       uint32
	queueId           uint32
	payload           []byte
}

func NewRequestFromStream(stream io.Reader) (*Request, error) {
	var commonHeader *CommonHeader
	var operationalHeader *RequestHeader
	var err error
	commonHeader, err = CommonHeaderFromStream(stream)
	if err != nil {
		return nil, err
	}
	operationalHeader, err = RequestHeaderFromStream(stream)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 12)
	if _, err = stream.Read(buffer); err != nil {
		return nil, err
	}
	var namespaceId = binary.BigEndian.Uint32(buffer[0:])
	var queueuId = binary.BigEndian.Uint32(buffer[4:])
	var payloadLength = binary.BigEndian.Uint32(buffer[8:])
	var payload = make([]byte, int(payloadLength))
	if _, err = stream.Read(payload); err != nil {
		return nil, err
	}
	return &Request{
		commonHeader:      *commonHeader,
		operationalHeader: *operationalHeader,
		namespaceId:       namespaceId,
		queueId:           queueuId,
		payload:           payload,
	}, err
}

func NewRequest(commonHeader CommonHeader, operationalHeader RequestHeader) (*Request, error) {
	return &Request{
		commonHeader:      commonHeader,
		operationalHeader: operationalHeader,
	}, nil
}

func (request *Request) ToString() string {
	return strconv.FormatUint(uint64(request.namespaceId), 10) + "-" + strconv.FormatUint(uint64(request.queueId), 10)
}

func (request *Request) SetPayload(payload []byte) {
	request.payload = payload
}

func (request *Request) GetPayload() []byte {
	return request.payload
}

func (request *Request) GetOpcode() uint8 {
	return request.operationalHeader.Opcode
}

func (request *Request) GetNamespaceId() uint32 {
	return request.namespaceId
}

func (request *Request) GetQueueId() uint32 {
	return request.queueId
}

func (request *Request) ToResponse() *Response {
	var responseHeader = &ResponseHeader{
		Opcode:  request.operationalHeader.Opcode,
		Opaque2: request.operationalHeader.Opaque2,
		Status:  ResponseStatusUnknown,
	}
	return &Response{
		commonHeader:      request.commonHeader,
		operationalHeader: *responseHeader,
		namespaceId:       request.namespaceId,
		queueId:           request.queueId,
	}
}

func (request *Request) ToStream(stream io.Writer) error {
	buffer := request.ToBytes()
	_, err := stream.Write(buffer)
	return err
}

func (request *Request) ToBytes() []byte {
	cbuffer := request.commonHeader.ToBytes()
	obuffer := request.operationalHeader.ToBytes()
	buffer := make([]byte, 12+uint32(len(request.payload)))
	binary.BigEndian.PutUint32(buffer[0:], request.namespaceId)
	binary.BigEndian.PutUint32(buffer[4:], request.queueId)
	binary.BigEndian.PutUint32(buffer[8:], uint32(len(request.payload)))
	copy(buffer[12:], request.payload)
	return append(append(cbuffer, obuffer...), buffer...)
}
