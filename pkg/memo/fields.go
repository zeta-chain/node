package memo

// Fields is the interface for memo fields
type Fields interface {
	// Pack encodes the memo fields
	Pack(opCode, encodingFormat uint8) (byte, []byte, error)

	// Unpack decodes the memo fields
	Unpack(opCode, encodingFormat, dataFlags uint8, data []byte) error

	// Validate checks if the fields are valid
	Validate(opCode uint8) error
}
