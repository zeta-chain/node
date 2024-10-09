package memo

import (
	"fmt"

	"github.com/pkg/errors"

	zetamath "github.com/zeta-chain/node/pkg/math"
)

const (
	// MemoIdentifier is the ASCII code of 'Z' (0x5A)
	MemoIdentifier byte = 0x5A

	// HeaderSize is the size of the memo header
	HeaderSize = 3

	// MaskVersion is the mask for the version bits(5~7)
	BitMaskVersion byte = 0b11100000

	// BitMaskEncodingFormat is the mask for the encoding format bits(3~4)
	BitMaskEncodingFormat byte = 0b00011000

	// BitMaskOpCode is the mask for the operation code bits(0~2)
	BitMaskOpCode byte = 0b00000111
)

// Enum for non-EVM chain inbound operation code (3 bits)
const (
	InboundOpCodeDeposit        uint8 = 0b000 // operation 'deposit'
	InboundOpCodeDepositAndCall uint8 = 0b001 // operation 'deposit_and_call'
	InboundOpCodeCall           uint8 = 0b010 // operation 'call'
	InboundOpCodeMax            uint8 = 0b011 // operation max value
)

// InboundMemo represents the memo structure for non-EVM chains
type InboundMemo struct {
	// Version is the memo Version
	Version uint8

	// EncodingFormat is the memo encoding format
	EncodingFormat uint8

	// OpCode is the inbound operation code
	OpCode uint8

	// FieldsV0 contains the memo fields V0
	// Note: add a FieldsV1 if major update is needed in the future
	FieldsV0
}

// EncodeToBytes encodes a InboundMemo struct to raw bytes
func EncodeToBytes(memo *InboundMemo) ([]byte, error) {
	// encode memo basics
	basics := encodeBasics(memo)

	// encode memo fields based on version
	var data []byte
	var err error
	switch memo.Version {
	case 0:
		data, err = memo.FieldsV0.Pack(memo.EncodingFormat)
	default:
		return nil, fmt.Errorf("unsupported memo version: %d", memo.Version)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to pack memo fields version: %d", memo.Version)
	}

	return append(basics, data...), nil
}

// DecodeFromBytes decodes a InboundMemo struct from raw bytes
//
// Returns an error if given data is not a valid memo
func DecodeFromBytes(data []byte) (*InboundMemo, error) {
	memo := &InboundMemo{}

	// decode memo basics
	err := decodeBasics(data, memo)
	if err != nil {
		return nil, err
	}

	// decode memo fields based on version
	switch memo.Version {
	case 0:
		err = memo.FieldsV0.Unpack(data, memo.EncodingFormat)
	default:
		return nil, fmt.Errorf("unsupported memo version: %d", memo.Version)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack memo fields version: %d", memo.Version)
	}

	return memo, nil
}

// encodeBasics encodes the version, encoding format and operation code
func encodeBasics(memo *InboundMemo) []byte {
	// 2 bytes: [identifier + ctrlByte]
	basics := make([]byte, HeaderSize-1)

	// set byte-0 as memo identifier
	basics[0] = MemoIdentifier

	// set version # and encoding format
	var ctrlByte byte
	ctrlByte = zetamath.SetBits(ctrlByte, BitMaskVersion, memo.Version)
	ctrlByte = zetamath.SetBits(ctrlByte, BitMaskEncodingFormat, memo.EncodingFormat)
	ctrlByte = zetamath.SetBits(ctrlByte, BitMaskOpCode, memo.OpCode)

	// set ctrlByte to byte-1
	basics[1] = ctrlByte

	return basics
}

// decodeBasics decodes version, encoding format and operation code
func decodeBasics(data []byte, memo *InboundMemo) error {
	// memo data must be longer than the header size
	if len(data) <= HeaderSize {
		return errors.New("memo data too short")
	}

	// byte-0 is the memo identifier
	if data[0] != MemoIdentifier {
		return errors.New("memo identifier mismatch")
	}

	// byte-1 is the control byte
	ctrlByte := data[1]

	// extract version #
	memo.Version = zetamath.GetBits(ctrlByte, BitMaskVersion)
	if memo.Version != 0 {
		return fmt.Errorf("unsupported memo version: %d", memo.Version)
	}

	// extract encoding format
	memo.EncodingFormat = zetamath.GetBits(ctrlByte, BitMaskEncodingFormat)
	if memo.EncodingFormat >= EncodingFmtMax {
		return fmt.Errorf("invalid encoding format: %d", memo.EncodingFormat)
	}

	// extract operation code
	memo.OpCode = zetamath.GetBits(ctrlByte, BitMaskOpCode)
	if memo.OpCode >= InboundOpCodeMax {
		return fmt.Errorf("invalid operation code: %d", memo.OpCode)
	}

	return nil
}
