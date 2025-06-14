package tcp

import (
	"encoding/binary"
	"errors"
	"io"
)

type MetadataField struct {
	FieldTag   byte // 1 byte
	FieldSize  byte // 1 byte
	FieldValue []byte
}

/*
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | magic                         | version       |msg type   |RQ |
------+---------------+---------------+---------------+---------------+
    4 | opaque                                                        |
------+---------------+---------------+---------------+---------------+
*/

// KokaqWire represents the abstract wire protocol.
type KokaqWire struct {
	MessageType      byte
	RequestType      byte
	Opaque           uint32
	NamespaceName    uint32
	QueueName        uint32
	NamespaceNameLen byte
	QueueNameLen     byte
	Payload          []byte
	PayloadLength    uint16
	FieldCount       uint32
	Fields           []MetadataField
}

// NewKokaqWire creates a new instance of KokaqWire.
func NewKokaqWire(messageType byte, requestType byte, opaque uint32) *KokaqWire {
	return &KokaqWire{
		MessageType: messageType,
		RequestType: requestType,
		Opaque:      opaque,
		Fields:      []MetadataField{},
	}
}

// SetPayloadInternal sets the payload with optional namespace and queue names.
func (p *KokaqWire) SetPayloadInternal(payload []byte, namespaceName *uint32, queueName *uint32) {
	if namespaceName != nil {
		p.NamespaceNameLen = 4
		p.NamespaceName = *namespaceName
	}
	if queueName != nil {
		p.QueueNameLen = 4
		p.QueueName = *queueName
	}
	p.Payload = payload
	p.PayloadLength = uint16(len(payload))
}

// SetMetadataInternal sets metadata internally for the KokaqWire.
func (p *KokaqWire) SetMetadataInternal(fieldTag byte, value []byte) {
	field := MetadataField{
		FieldTag:   fieldTag,
		FieldValue: value,
		FieldSize:  byte(len(value)),
	}
	p.Fields = append(p.Fields, field)
	p.FieldCount = uint32(len(p.Fields))
}

// ToBytes serializes the KokaqWire object to a byte array.
func (p *KokaqWire) ToBytes() []byte {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint16(buffer[0:], Magic)    // Magic
	buffer[2] = Version                                 // Version
	buffer[3] = p.MessageType<<2 | p.RequestType        // MessageType and RequestType
	binary.LittleEndian.PutUint32(buffer[4:], p.Opaque) // Opaque
	return buffer
}

// WriteToStream writes the KokaqWire to a stream.
func (p *KokaqWire) WriteToStream(stream io.Writer) error {
	buffer := p.ToBytes()
	_, err := stream.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

// ReadFromStream reads KokaqWire from a stream.
func (p *KokaqWire) ReadFromStream(stream io.Reader) error {
	buffer := make([]byte, 8)
	_, err := stream.Read(buffer)
	if err != nil {
		return err
	}

	if binary.LittleEndian.Uint16(buffer[0:]) != Magic || buffer[2] != Version {
		return errors.New("invalid byte stream")
	}

	flags := buffer[3]
	p.RequestType = flags & 0x3                       // Extract the lower 2 bits for RequestType
	p.MessageType = flags >> 2                        // Extract the upper bits for MessageType
	p.Opaque = binary.LittleEndian.Uint32(buffer[4:]) // Opaque field

	return nil
}
