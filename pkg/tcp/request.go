package tcp

import (
	"encoding/binary"
	"errors"
	"io"
)

/*

Common wireframe for all requests
      |0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|0|1|2|3|4|5|6|7|
  byte|              0|              1|              2|              3|
------+---------------+---------------+---------------+---------------+
    0 | magic                         | version       |msg type   |RQ |
------+---------------+---------------+---------------+---------------+
    4 | opaque                                                        |
------+---------------+---------------+---------------+---------------+
    8 | opcode        | clientId      | opaque        | Tag/ID        |
------+---------------+---------------+---------------+---------------+

Payload
------+---------------+---------------+---------------+---------------+
    0 | Tag/ID (0x01) |     magic                     |  opaque       |
------+---------------+---------------+-------------------------------+
    4 | namespace len | queue length  | payload len                   |
------+---------------+---------------+---------------+---------------+
    8 | nsname?                                                       |
------+---------------+---------------+---------------+---------------+
   12 | queuename?                                                    |
------+---------------+---------------+---------------+---------------+
   16 | payload?                                                      |
------+---------------+---------------+---------------+---------------+

Metadata
------+---------------+---------------+---------------+---------------+
    0 | Tag/ID (0x02) |     magic                     |  opaque       |
------+---------------+---------------+-------------------------------+
    4 | fieldcount len                                                |
------+---------------+---------------+---------------+---------------+
    8 | fieldtag      |  size         |               |               |
------+---------------+---------------+---------------+---------------+
*/

// KokaqWireRequest extends KokaqWire and adds more request-specific fields.
type KokaqWireRequest struct {
	KokaqWire
	OpCode         byte
	ClientId       byte
	RequestTag     byte
	RequestOpaque  byte
	PayloadOpaque  byte
	MetadataOpaque byte
}

func NewKokaqWireRequestWithMetadata(messageType byte, requestType byte, opaque uint32, opCode byte, clientId byte, requestOpaque byte, metadataOpaque byte) *KokaqWireRequest {
	return &KokaqWireRequest{
		KokaqWire: KokaqWire{
			MessageType: messageType,
			RequestType: requestType,
			Opaque:      opaque,
			Fields:      []MetadataField{},
		},
		OpCode:         opCode,
		ClientId:       clientId,
		RequestOpaque:  requestOpaque,
		RequestTag:     BodyTypeMetadata,
		MetadataOpaque: metadataOpaque,
	}
}

func NewKokaqWireRequestWithPayload(messageType byte, requestType byte, opaque uint32, opCode byte, clientId byte, requestOpaque byte, payloadOpaque byte) *KokaqWireRequest {
	return &KokaqWireRequest{
		KokaqWire: KokaqWire{
			MessageType: messageType,
			RequestType: requestType,
			Opaque:      opaque,
			Fields:      []MetadataField{},
		},
		OpCode:        opCode,
		ClientId:      clientId,
		RequestOpaque: requestOpaque,
		RequestTag:    BodyTypePayload,
		PayloadOpaque: payloadOpaque,
	}
}

func NewKokaqWireRequest(messageType byte, requestType byte, opaque uint32, opCode byte, clientId byte, requestOpaque byte) *KokaqWireRequest {
	return &KokaqWireRequest{
		KokaqWire: KokaqWire{
			MessageType: messageType,
			RequestType: requestType,
			Opaque:      opaque,
			Fields:      []MetadataField{},
		},
		OpCode:        opCode,
		ClientId:      clientId,
		RequestOpaque: requestOpaque,
		RequestTag:    None,
	}
}

// ToBytes serializes the request into a byte array.
func (p *KokaqWireRequest) ToBytes() []byte {
	// Initial buffer size is 12 bytes (could grow based on payload/metadata)
	buffer := make([]byte, 12)
	copy(buffer, p.KokaqWire.ToBytes())
	buffer[8] = p.OpCode
	buffer[9] = p.ClientId
	buffer[10] = p.RequestOpaque
	buffer[11] = p.RequestTag

	if p.RequestTag != None {
		switch p.RequestTag {
		case BodyTypePayload:
			payloadBuffer := p.serializePayload()
			return append(buffer, payloadBuffer...)
		case BodyTypeMetadata:
			metadataBuffer := p.serializeMetadata()
			return append(buffer, metadataBuffer...)
		}
	}
	return buffer
}

// Helper to serialize payload-specific data
func (p *KokaqWireRequest) serializePayload() []byte {
	// Calculate extra size based on NamespaceNameLength, QueueNameLength, and Payload
	extraSize := int(p.KokaqWire.NamespaceNameLen) + int(p.QueueNameLen) + len(p.Payload)
	buffer := make([]byte, 20+extraSize)
	buffer[0] = p.RequestTag
	binary.LittleEndian.PutUint16(buffer[1:], Magic)
	buffer[3] = p.PayloadOpaque
	buffer[4] = p.NamespaceNameLen
	buffer[5] = p.QueueNameLen
	binary.LittleEndian.PutUint16(buffer[6:], p.PayloadLength)

	offset := 8
	if p.NamespaceNameLen != None {
		binary.LittleEndian.PutUint32(buffer[offset:], p.NamespaceName)
		offset += 4
	}
	if p.QueueNameLen != None {
		binary.LittleEndian.PutUint32(buffer[offset:], p.QueueName)
		offset += 4
	}
	copy(buffer[offset:], p.Payload)
	return buffer
}

// Helper to serialize metadata-specific data
func (p *KokaqWireRequest) serializeMetadata() []byte {
	// Calculate extra size based on fields
	extraSize := 0
	for _, field := range p.Fields {
		extraSize += int(field.FieldSize) + 2
	}
	buffer := make([]byte, 20+extraSize)
	buffer[0] = p.RequestTag
	binary.LittleEndian.PutUint16(buffer[1:], Magic)
	buffer[3] = p.MetadataOpaque
	binary.LittleEndian.PutUint32(buffer[4:], p.FieldCount)

	offset := 8
	for _, field := range p.Fields {
		buffer[offset] = field.FieldTag
		buffer[offset+1] = field.FieldSize
		copy(buffer[offset+2:], field.FieldValue)
		offset += 2 + int(field.FieldSize)
	}
	return buffer
}

// WriteToStream writes the serialized request to the stream.
func (p *KokaqWireRequest) WriteToStream(stream io.Writer) error {
	buffer := p.ToBytes()
	_, err := stream.Write(buffer)
	return err
}

// ReadFromStream reads the request data from the stream.
func (p *KokaqWireRequest) ReadFromStream(stream io.Reader) error {
	err := p.KokaqWire.ReadFromStream(stream)
	if err != nil {
		return err
	}

	buffer := make([]byte, 4)
	_, err = stream.Read(buffer)
	if err != nil {
		return err
	}

	p.OpCode = buffer[0]
	p.ClientId = buffer[1]
	p.RequestOpaque = buffer[2]
	p.RequestTag = buffer[3]

	if p.RequestTag != None {
		buffer = make([]byte, 8)
		_, err = stream.Read(buffer)
		if err != nil {
			return err
		}

		if buffer[0] != p.RequestTag || binary.LittleEndian.Uint16(buffer[1:]) != Magic {
			return errors.New("invalid request type")
		}

		switch p.RequestTag {
		case BodyTypePayload:
			return p.readPayloadFromStream(stream, buffer)
		case BodyTypeMetadata:
			return p.readMetadataFromStream(stream, buffer)
		}
	}
	return nil
}

// readPayloadFromStream reads the payload data from the stream.
func (p *KokaqWireRequest) readPayloadFromStream(stream io.Reader, buffer []byte) error {
	p.PayloadOpaque = buffer[3]
	p.NamespaceNameLen = buffer[4]
	p.QueueNameLen = buffer[5]
	p.PayloadLength = binary.LittleEndian.Uint16(buffer[6:])

	if p.NamespaceNameLen != None {
		buffer = make([]byte, 4)
		_, err := stream.Read(buffer)
		if err != nil {
			return err
		}
		p.NamespaceName = binary.LittleEndian.Uint32(buffer)
	}

	if p.QueueNameLen != None {
		buffer = make([]byte, 4)
		_, err := stream.Read(buffer)
		if err != nil {
			return err
		}
		p.QueueName = binary.LittleEndian.Uint32(buffer)
	}

	if p.PayloadLength > 0 {

		p.Payload = make([]byte, p.PayloadLength)
		_, err := stream.Read(p.Payload)
		return err

	}
	return nil
}

// readMetadataFromStream reads the metadata from the stream.
func (p *KokaqWireRequest) readMetadataFromStream(stream io.Reader, buffer []byte) error {
	p.MetadataOpaque = buffer[3]
	p.FieldCount = binary.LittleEndian.Uint32(buffer[4:])

	for i := uint32(0); i < p.FieldCount; i++ {
		buffer = make([]byte, 2)
		_, err := stream.Read(buffer)
		if err != nil {
			return err
		}

		fieldTag := buffer[0]
		fieldSize := buffer[1]
		fieldValue := make([]byte, fieldSize)
		_, err = stream.Read(fieldValue)
		if err != nil {
			return err
		}

		p.Fields = append(p.Fields, MetadataField{
			FieldTag:   fieldTag,
			FieldValue: fieldValue,
			FieldSize:  fieldSize,
		})
	}
	return nil
}
