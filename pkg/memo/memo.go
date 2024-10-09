package memo

import (
	"fmt"

	"github.com/pkg/errors"

	zetamath "github.com/zeta-chain/node/pkg/math"
)

const (
	// MemoIdentifier is the ASCII code of 'Z' (0x5A)
	MemoIdentifier byte = 0x5A

	// MemoHeaderSize is the size of the memo header: [identifier + ctrlByte1+ ctrlByte2 + dataFlags]
	MemoHeaderSize = 4

	// MemoBasicsSize is the size of the memo basics: [identifier + ctrlByte1 + ctrlByte2]
	MemoBasicsSize = 3

	// MaskVersion is the mask for the version bits(upper 4 bits)
	MaskVersion byte = 0b11110000

	// MaskEncodingFormat is the mask for the encoding format bits(lower 4 bits)
	MaskEncodingFormat byte = 0b00001111

	// MaskOpCode is the mask for the operation code bits(upper 4 bits)
	MaskOpCode byte = 0b11110000

	// MaskCtrlReserved is the mask for reserved control bits (lower 4 bits)
	MaskCtrlReserved byte = 0b00001111
)

// Enum for non-EVM chain inbound operation code (4 bits)
const (
	OpCodeDeposit        uint8 = 0b0000 // operation 'deposit'
	OpCodeDepositAndCall uint8 = 0b0001 // operation 'deposit_and_call'
	OpCodeCall           uint8 = 0b0010 // operation 'call'
	OpCodeMax            uint8 = 0b0011 // operation max value
)

// InboundMemo represents the memo structure for non-EVM chains
type InboundMemo struct {
	// Version is the memo Version
	Version uint8

	// EncodingFormat is the memo encoding format
	EncodingFormat uint8

	// OpCode is the inbound operation code
	OpCode uint8

	// Reserved is the reserved control bits
	Reserved uint8

	// FieldsV0 contains the memo fields V0
	// Note: add a FieldsV1 if major update is needed in the future
	FieldsV0
}

// EncodeToBytes encodes a InboundMemo struct to raw bytes
func (m *InboundMemo) EncodeToBytes() ([]byte, error) {
	// validate memo basics
	err := m.ValidateBasics()
	if err != nil {
		return nil, err
	}

	// encode memo basics
	basics := m.EncodeBasics()

	// encode memo fields based on version
	var data []byte
	switch m.Version {
	case 0:
		data, err = m.FieldsV0.Pack(m.EncodingFormat)
	default:
		return nil, fmt.Errorf("invalid memo version: %d", m.Version)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to pack memo fields version: %d", m.Version)
	}

	return append(basics, data...), nil
}

// DecodeFromBytes decodes a InboundMemo struct from raw bytes
//
// Returns an error if given data is not a valid memo
func DecodeFromBytes(data []byte) (*InboundMemo, error) {
	memo := &InboundMemo{}

	// decode memo basics
	err := memo.DecodeBasics(data)
	if err != nil {
		return nil, err
	}

	// validate memo basics
	err = memo.ValidateBasics()
	if err != nil {
		return nil, err
	}

	// decode memo fields based on version
	switch memo.Version {
	case 0:
		err = memo.FieldsV0.Unpack(data[MemoBasicsSize:], memo.EncodingFormat)
	default:
		return nil, fmt.Errorf("invalid memo version: %d", memo.Version)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack memo fields version: %d", memo.Version)
	}

	return memo, nil
}

// Validate checks if the memo is valid
func (m *InboundMemo) ValidateBasics() error {
	if m.EncodingFormat >= EncodingFmtMax {
		return fmt.Errorf("invalid encoding format: %d", m.EncodingFormat)
	}

	if m.OpCode >= OpCodeMax {
		return fmt.Errorf("invalid operation code: %d", m.OpCode)
	}

	// reserved control bits must be zero
	if m.Reserved != 0 {
		return fmt.Errorf("reserved control bits are not zero: %d", m.Reserved)
	}

	return nil
}

// EncodeBasics encodes theidentifier, version, encoding format and operation code
func (m *InboundMemo) EncodeBasics() []byte {
	// basics: [identifier + ctrlByte1 + ctrlByte2]
	basics := make([]byte, MemoBasicsSize)

	// set byte-0 as memo identifier
	basics[0] = MemoIdentifier

	// set version #, encoding format
	var ctrlByte1 byte
	ctrlByte1 = zetamath.SetBits(ctrlByte1, MaskVersion, m.Version)
	ctrlByte1 = zetamath.SetBits(ctrlByte1, MaskEncodingFormat, m.EncodingFormat)
	basics[1] = ctrlByte1

	// set operation code, reserved bits
	var ctrlByte2 byte
	ctrlByte2 = zetamath.SetBits(ctrlByte2, MaskOpCode, m.OpCode)
	ctrlByte2 = zetamath.SetBits(ctrlByte2, MaskCtrlReserved, m.Reserved)
	basics[2] = ctrlByte2

	return basics
}

// DecodeBasics decodes the identifier, version, encoding format and operation code
func (m *InboundMemo) DecodeBasics(data []byte) error {
	// memo data must be longer than the header size
	if len(data) <= MemoHeaderSize {
		return errors.New("memo is too short")
	}

	// byte-0 is the memo identifier
	if data[0] != MemoIdentifier {
		return fmt.Errorf("invalid memo identifier: %d", data[0])
	}

	// extract version #, encoding format
	ctrlByte1 := data[1]
	m.Version = zetamath.GetBits(ctrlByte1, MaskVersion)
	m.EncodingFormat = zetamath.GetBits(ctrlByte1, MaskEncodingFormat)

	// extract operation code, reserved bits
	ctrlByte2 := data[2]
	m.OpCode = zetamath.GetBits(ctrlByte2, MaskOpCode)
	m.Reserved = zetamath.GetBits(ctrlByte2, MaskCtrlReserved)

	return nil
}
