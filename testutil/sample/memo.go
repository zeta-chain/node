package sample

import (
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/memo"
)

// MemoHead is a helper function to create a memo head
// Note: all arguments are assume to be <= 0b1111 for simplicity.
func MemoHead(version, encodingFmt, opCode, reserved, flags uint8) []byte {
	head := make([]byte, memo.MemoHeaderSize)
	head[0] = memo.MemoIdentifier
	head[1] = version<<4 | encodingFmt
	head[2] = opCode<<4 | reserved
	head[3] = flags
	return head
}

// ABIPack is a helper function to simulates the abi.Pack function.
// Note: all arguments are assumed to be <= 32 bytes for simplicity.
func ABIPack(t *testing.T, args ...memo.CodecArg) []byte {
	packedData := make([]byte, 0)

	// data offset for 1st dynamic-length field
	offset := memo.ABIAlignment * len(args)

	// 1. pack 32-byte offset for each dynamic-length field (bytes, string)
	// 2. pack actual data for each fixed-length field (address)
	for _, arg := range args {
		switch arg.Type {
		case memo.ArgTypeBytes:
			// left-pad length as uint16
			buff := make([]byte, 2)
			binary.BigEndian.PutUint16(buff, uint16(offset))
			offsetData := abiPad32(t, buff, true)
			packedData = append(packedData, offsetData...)

			argLen := len(arg.Arg.([]byte))
			if argLen > 0 {
				offset += memo.ABIAlignment * 2 // [length + data]
			} else {
				offset += memo.ABIAlignment // only [length]
			}

		case memo.ArgTypeString:
			// left-pad length as uint16
			buff := make([]byte, 2)
			binary.BigEndian.PutUint16(buff, uint16(offset))
			offsetData := abiPad32(t, buff, true)
			packedData = append(packedData, offsetData...)

			argLen := len([]byte(arg.Arg.(string)))
			if argLen > 0 {
				offset += memo.ABIAlignment * 2 // [length + data]
			} else {
				offset += memo.ABIAlignment // only [length]
			}

		case memo.ArgTypeAddress: // left-pad for address
			data := abiPad32(t, arg.Arg.(common.Address).Bytes(), true)
			packedData = append(packedData, data...)
		}
	}

	// pack dynamic-length fields
	dynamicData := abiPackDynamicData(t, args...)
	packedData = append(packedData, dynamicData...)

	return packedData
}

// CompactPack is a helper function to pack arguments into compact encoded data
// Note: all arguments are assumed to be <= 65535 bytes for simplicity.
func CompactPack(encodingFmt uint8, args ...memo.CodecArg) []byte {
	var (
		length     int
		packedData []byte
	)

	for _, arg := range args {
		// get length of argument
		switch arg.Type {
		case memo.ArgTypeBytes:
			length = len(arg.Arg.([]byte))
		case memo.ArgTypeString:
			length = len([]byte(arg.Arg.(string)))
		default:
			// skip length for other types
			length = -1
		}

		// append length in bytes
		if length != -1 {
			switch encodingFmt {
			case memo.EncodingFmtCompactShort:
				packedData = append(packedData, byte(length))
			case memo.EncodingFmtCompactLong:
				buff := make([]byte, 2)
				binary.LittleEndian.PutUint16(buff, uint16(length))
				packedData = append(packedData, buff...)
			}
		}

		// append actual data in bytes
		switch arg.Type {
		case memo.ArgTypeBytes:
			packedData = append(packedData, arg.Arg.([]byte)...)
		case memo.ArgTypeAddress:
			packedData = append(packedData, arg.Arg.(common.Address).Bytes()...)
		case memo.ArgTypeString:
			packedData = append(packedData, []byte(arg.Arg.(string))...)
		}
	}

	return packedData
}

// abiPad32 is a helper function to pad a byte slice to 32 bytes
func abiPad32(t *testing.T, data []byte, left bool) []byte {
	// nothing needs to be encoded, return empty bytes
	if len(data) == 0 {
		return []byte{}
	}

	require.LessOrEqual(t, len(data), memo.ABIAlignment)
	padded := make([]byte, 32)

	if left {
		// left-pad the data for fixed-size types
		copy(padded[32-len(data):], data)
	} else {
		// right-pad the data for dynamic types
		copy(padded, data)
	}
	return padded
}

// apiPackDynamicData is a helper function to pack dynamic-length data
func abiPackDynamicData(t *testing.T, args ...memo.CodecArg) []byte {
	packedData := make([]byte, 0)

	// pack with ABI format: length + data
	for _, arg := range args {
		// get length
		var length int
		switch arg.Type {
		case memo.ArgTypeBytes:
			length = len(arg.Arg.([]byte))
		case memo.ArgTypeString:
			length = len([]byte(arg.Arg.(string)))
		default:
			continue
		}

		// append length in bytes
		lengthData := abiPad32(t, []byte{byte(length)}, true)
		packedData = append(packedData, lengthData...)

		// append actual data in bytes
		switch arg.Type {
		case memo.ArgTypeBytes: // right-pad for bytes
			data := abiPad32(t, arg.Arg.([]byte), false)
			packedData = append(packedData, data...)
		case memo.ArgTypeString: // right-pad for string
			data := abiPad32(t, []byte(arg.Arg.(string)), false)
			packedData = append(packedData, data...)
		}
	}

	return packedData
}
