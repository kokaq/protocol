package tcp

// Global constants
const (
	None    byte   = 0x00  // Represents 'None' value.
	One     byte   = 0x01  // Represents 'One' value.
	Magic   uint16 = 0x420 // Magic constant (ushort equivalent in Go is uint16).
	Version byte   = 0x1   // Version constant.

	// MessageType represents different types of messages.
	MessageTypeOperational byte = 0x1 // Operational message.
	MessageTypeAdmin       byte = 0x2 // Admin message.
	MessageTypeControl     byte = 0x3 // Control message.

	// RequestTypeTwoWay represents different types of requests.
	RequestTypeResponse byte = 0x1 // Response.
	RequestTypeTwoWay   byte = 0x1 // Two-way request.
	RequestTypeOneWay   byte = 0x3 // One-way request.

	// OpCode represents operation codes for various actions.
	OpCodeNoOp            byte = 0x00 // No operation.
	OpCodeCreate          byte = 0x01 // Create operation.
	OpCodeDelete          byte = 0x02 // Delete operation.
	OpCodeGet             byte = 0x03 // Get operation.
	OpCodePeek            byte = 0x04 // Peek operation.
	OpCodePop             byte = 0x05 // Pop operation.
	OpCodePush            byte = 0x06 // Push operation.
	OpCodeAcquirePeekLock byte = 0x07 // Acquire peek lock operation.
	OpCodeReleasePeekLock byte = 0x08 // Release peek lock operation.

	// ClientId represents different types of clients.
	ClientProxyHttp      byte = 0x01 // HTTP proxy client.
	ClientProxyAmqp      byte = 0x02 // AMQP proxy client.
	ClientQueueService   byte = 0x03 // Queue service client.
	ClientStorageService byte = 0x04 // Storage service client.
	ClientHealthService  byte = 0x05 // Health service client.

	// ResponseStatus represents various response statuses.
	ResponseStatusSuccess        byte = 0x01 // Success response status.
	ResponseStatusFail           byte = 0x02 // Fail response status.
	ResponseStatusPartialSuccess byte = 0x03 // Partial success response status.
	ResponseStatusUnknown        byte = 0x04 // Unknown response status.
	ResponseStatusReasonOk       byte = 0x01 // Ok
	ResponseStatusReasonBad      byte = 0x02 // Bad
	ResponseStatusReasonExists   byte = 0x03 // Item Exists
	ResponseStatusNotAllowed     byte = 0x04 // NotAllowed
	ResponseStatusReasonInfra    byte = 0x05 // Infra

	// BodyType represents types of message bodies.
	BodyTypeMetadata byte = 0x01 // Metadata body type.
	BodyTypePayload  byte = 0x02 // Payload body type.

	// TagMeta represents metadata tags for various operations.
	TagMetaTTL                  byte = 0x01 // Time to live tag.
	TagMetaVersionMeta          byte = 0x02 // Version tag.
	TagMetaCreationTime         byte = 0x03 // Creation time tag.
	TagMetaExpTime              byte = 0x04 // Expiration time tag.
	TagMetaUUID                 byte = 0x05 // UUID tag.
	TagMetaLastModificationTime byte = 0x06 // Last modification time tag.
	TagMetaOriginatorRequestId  byte = 0x07 // Originator request ID tag.
	TagMetaCorrelationId        byte = 0x08 // Correlation ID tag.
	TagMetaRequestTime          byte = 0x09 // Request time tag.
)
