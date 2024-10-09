package memo

// Fields is the interface for memo fields
type Fields interface {
	// Pack encodes the memo fields
	Pack(encodingFormat uint8) ([]byte, error)

	// Unpack decodes the memo fields
	Unpack(data []byte, encodingFormat uint8) error
}
