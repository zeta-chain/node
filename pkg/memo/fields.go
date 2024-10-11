package memo

// Fields is the interface for memo fields
type Fields interface {
	// Pack encodes the memo fields
	Pack(opCode, encodingFormat, dataFlags uint8) ([]byte, error)

	// Unpack decodes the memo fields
	Unpack(opCode, encodingFormat, dataFlags uint8, data []byte) error

	// Validate checks if the fields are valid
	Validate(opCode, dataFlags uint8) error

	// DataFlags build the data flags for the fields
	DataFlags() uint8
}
