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
    4 | message size in bytes                                         |
------+---------------+---------------+---------------+---------------+
    8 | opaque                                                        |
------+---------------+---------------+---------------+---------------+
   12 | opcode        | status|reason | opaque        | Tag/ID        |
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
// KokaqWireResponse struct
type KokaqWireResponse struct {
	KokaqWire
	OpCode         byte
	Status         byte
	Reason         byte
	ResponseOpaque byte
	ResponseTag    byte
	MetadataOpaque byte
	PayloadOpaque  byte
}

func NewKokaqWireResponseFromRequest(request *KokaqWireRequest) *KokaqWireResponse {
	return &KokaqWireResponse{
		KokaqWire:      *NewKokaqWire(request.MessageType, RequestTypeResponse, request.Opaque),
		OpCode:         request.OpCode,
		ResponseOpaque: request.RequestOpaque,
		ResponseTag:    None,
	}
}

func NewKokaqWireResponseWithMetadata(messageType, requestType byte, opaque uint32, opCode, status, reason, responseOpaque, metaDataOpaque byte) *KokaqWireResponse {
	return &KokaqWireResponse{
		KokaqWire:      *NewKokaqWire(messageType, requestType, opaque),
		OpCode:         opCode,
		Status:         status,
		Reason:         reason,
		ResponseOpaque: responseOpaque,
		ResponseTag:    BodyTypeMetadata,
		MetadataOpaque: metaDataOpaque,
	}
}

func NewKokaqWireResponseWithPayload(messageType, requestType byte, opaque uint32, opCode, status, reason, responseOpaque, payloadOpaque byte) *KokaqWireResponse {
	return &KokaqWireResponse{
		KokaqWire:      *NewKokaqWire(messageType, requestType, opaque),
		OpCode:         opCode,
		Status:         status,
		Reason:         reason,
		ResponseOpaque: responseOpaque,
		ResponseTag:    BodyTypePayload,
		PayloadOpaque:  payloadOpaque,
	}
}

func NewKokaqWireResponse(messageType byte, requestType byte, opaque uint32, opCode, status, reason, responseOpaque byte) *KokaqWireResponse {
	return &KokaqWireResponse{
		KokaqWire:      *NewKokaqWire(messageType, requestType, opaque),
		OpCode:         opCode,
		Status:         status,
		Reason:         reason,
		ResponseTag:    None,
		ResponseOpaque: responseOpaque,
	}
}

// ToBytes serializes the request into a byte array.
func (p *KokaqWireResponse) ToBytes() []byte {
	// Initial buffer size is 12 bytes (could grow based on payload/metadata)
	buffer := make([]byte, 12)
	copy(buffer, p.KokaqWire.ToBytes())
	buffer[8] = p.OpCode
	buffer[9] = p.Status<<4 | p.Reason
	buffer[10] = p.ResponseOpaque
	buffer[11] = p.ResponseTag

	if p.ResponseTag != None {
		switch p.ResponseTag {
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
func (p *KokaqWireResponse) serializePayload() []byte {
	// Calculate extra size based on NamespaceNameLength, QueueNameLength, and Payload
	extraSize := int(p.KokaqWire.NamespaceNameLen) + int(p.QueueNameLen) + len(p.Payload)
	buffer := make([]byte, 20+extraSize)
	buffer[0] = p.ResponseTag
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
func (p *KokaqWireResponse) serializeMetadata() []byte {
	// Calculate extra size based on fields
	extraSize := 0
	for _, field := range p.Fields {
		extraSize += int(field.FieldSize) + 2
	}
	buffer := make([]byte, 20+extraSize)
	buffer[0] = p.ResponseTag
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
func (p *KokaqWireResponse) WriteToStream(stream io.Writer) error {
	buffer := p.ToBytes()
	_, err := stream.Write(buffer)
	return err
}

// ReadFromStream reads the request data from the stream.
func (p *KokaqWireResponse) ReadFromStream(stream io.Reader) error {
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
	p.Status = buffer[1] >> 4
	p.Reason = buffer[1] & 0x0F
	p.ResponseOpaque = buffer[2]
	p.ResponseTag = buffer[3]

	if p.ResponseTag != None {
		buffer = make([]byte, 8)
		_, err = stream.Read(buffer)
		if err != nil {
			return err
		}

		if buffer[0] != p.ResponseTag || binary.LittleEndian.Uint16(buffer[1:]) != Magic {
			return errors.New("invalid request type")
		}

		switch p.ResponseTag {
		case BodyTypePayload:
			return p.readPayloadFromStream(stream, buffer)
		case BodyTypeMetadata:
			return p.readMetadataFromStream(stream, buffer)
		}
	}
	return nil
}

// readPayloadFromStream reads the payload data from the stream.
func (p *KokaqWireResponse) readPayloadFromStream(stream io.Reader, buffer []byte) error {
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

	p.Payload = make([]byte, p.PayloadLength)
	_, err := stream.Read(p.Payload)
	return err
}

// readMetadataFromStream reads the metadata from the stream.
func (p *KokaqWireResponse) readMetadataFromStream(stream io.Reader, buffer []byte) error {
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
