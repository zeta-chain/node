package memo

import (
	"fmt"

	"github.com/pkg/errors"
)

// Enum for non-EVM chain memo encoding format (2 bits)
const (
	// EncodingFmtABI represents ABI encoding format
	EncodingFmtABI uint8 = 0b00

	// EncodingFmtCompactShort represents 'compact short' encoding format
	EncodingFmtCompactShort uint8 = 0b01

	// EncodingFmtCompactLong represents 'compact long' encoding format
	EncodingFmtCompactLong uint8 = 0b10

	// EncodingFmtMax is the max value of encoding format
	EncodingFmtMax uint8 = 0b11
)

// Enum for length of bytes used to encode compact data
const (
	LenBytesShort = 1
	LenBytesLong  = 2
)

// Codec is the interface for a codec
type Codec interface {
	// AddArguments adds a list of arguments to the codec
	AddArguments(args ...CodecArg)

	// PackArguments packs the arguments into the encoded data
	PackArguments() ([]byte, error)

	// UnpackArguments unpacks the encoded data into the arguments
	UnpackArguments(data []byte) error
}

// GetLenBytes returns the number of bytes used to encode the length of the data
func GetLenBytes(encodingFmt uint8) (int, error) {
	switch encodingFmt {
	case EncodingFmtCompactShort:
		return LenBytesShort, nil
	case EncodingFmtCompactLong:
		return LenBytesLong, nil
	default:
		return 0, fmt.Errorf("invalid compact encoding format %d", encodingFmt)
	}
}

// GetCodec returns the codec based on the encoding format
func GetCodec(encodingFormat uint8) (Codec, error) {
	switch encodingFormat {
	case EncodingFmtABI:
		return NewCodecABI(), nil
	case EncodingFmtCompactShort, EncodingFmtCompactLong:
		return NewCodecCompact(encodingFormat)
	default:
		return nil, errors.New("unsupported encoding format")
	}
}
