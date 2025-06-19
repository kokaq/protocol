package tcp

import (
	"encoding/binary"
	"io"
)

type Response struct {
	commonHeader      CommonHeader
	operationalHeader ResponseHeader
	namespaceId       uint32
	queueId           uint32
	payload           []byte
}

func NewResponseFromStream(stream io.Reader) (*Response, error) {
	var commonHeader *CommonHeader
	var operationalHeader *ResponseHeader
	var err error
	commonHeader, err = CommonHeaderFromStream(stream)
	if err != nil {
		return nil, err
	}
	operationalHeader, err = ResponseHeaderFromStream(stream)
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
	return &Response{
		commonHeader:      *commonHeader,
		operationalHeader: *operationalHeader,
		namespaceId:       namespaceId,
		queueId:           queueuId,
		payload:           payload,
	}, err
}

func NewResponse(commonHeader CommonHeader, operationalHeader ResponseHeader) (*Response, error) {
	return &Response{
		commonHeader:      commonHeader,
		operationalHeader: operationalHeader,
	}, nil
}

func (response *Response) SetPayload(payload []byte) {
	response.payload = payload
}

func (response *Response) GetPayload() []byte {
	return response.payload
}

func (response *Response) SetStatus(status uint8) {
	response.operationalHeader.Status = status
}

func (response *Response) GetStatus() uint8 {
	return response.operationalHeader.Status
}

func (response *Response) GetOpcode() uint8 {
	return response.operationalHeader.Opcode
}

func (response *Response) GetNamespaceId() uint32 {
	return response.namespaceId
}

func (response *Response) GetQueueId() uint32 {
	return response.queueId
}

func (response *Response) ToStream(stream io.Writer) error {
	buffer := response.ToBytes()
	_, err := stream.Write(buffer)
	return err
}
func (response *Response) ToBytes() []byte {
	cbuffer := response.commonHeader.ToBytes()
	obuffer := response.operationalHeader.ToBytes()
	buffer := make([]byte, 12+uint32(len(response.payload)))
	binary.BigEndian.PutUint32(buffer[0:], response.namespaceId)
	binary.BigEndian.PutUint32(buffer[4:], response.queueId)
	binary.BigEndian.PutUint32(buffer[8:], uint32(len(response.payload)))
	copy(buffer[12:], response.payload)
	return append(append(cbuffer, obuffer...), buffer...)
}
