package memo

import (
	"fmt"
)

type EncodingFormat uint8

// Enum for non-EVM chain memo encoding format (2 bits)
const (
	// EncodingFmtABI represents ABI encoding format
	EncodingFmtABI EncodingFormat = 0b0000

	// EncodingFmtCompactShort represents 'compact short' encoding format
	EncodingFmtCompactShort EncodingFormat = 0b0001

	// EncodingFmtCompactLong represents 'compact long' encoding format
	EncodingFmtCompactLong EncodingFormat = 0b0010

	// EncodingFmtInvalid represents invalid encoding format
	EncodingFmtInvalid EncodingFormat = 0b0011
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
func GetLenBytes(encodingFmt EncodingFormat) (int, error) {
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
func GetCodec(encodingFmt EncodingFormat) (Codec, error) {
	switch encodingFmt {
	case EncodingFmtABI:
		return NewCodecABI(), nil
	case EncodingFmtCompactShort, EncodingFmtCompactLong:
		return NewCodecCompact(encodingFmt)
	default:
		return nil, fmt.Errorf("invalid encoding format %d", encodingFmt)
	}
}
