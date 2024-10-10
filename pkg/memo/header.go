package memo

import (
	"fmt"

	"github.com/pkg/errors"

	zetamath "github.com/zeta-chain/node/pkg/math"
)

const (
	// Identifier is the ASCII code of 'Z' (0x5A)
	Identifier byte = 0x5A

	// HeaderSize is the size of the memo header: [identifier + ctrlByte1+ ctrlByte2 + dataFlags]
	HeaderSize = 4

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

// Header represent the memo header
type Header struct {
	// Version is the memo Version
	Version uint8

	// EncodingFormat is the memo encoding format
	EncodingFormat uint8

	// OpCode is the inbound operation code
	OpCode uint8

	// Reserved is the reserved control bits
	Reserved uint8

	// DataFlags is the data flags
	DataFlags uint8
}

// EncodeToBytes encodes the memo header to raw bytes
func (h *Header) EncodeToBytes() ([]byte, error) {
	// validate header
	if err := h.Validate(); err != nil {
		return nil, err
	}

	// create buffer for the header
	data := make([]byte, HeaderSize)

	// set byte-0 as memo identifier
	data[0] = Identifier

	// set version #, encoding format
	var ctrlByte1 byte
	ctrlByte1 = zetamath.SetBits(ctrlByte1, MaskVersion, h.Version)
	ctrlByte1 = zetamath.SetBits(ctrlByte1, MaskEncodingFormat, h.EncodingFormat)
	data[1] = ctrlByte1

	// set operation code, reserved bits
	var ctrlByte2 byte
	ctrlByte2 = zetamath.SetBits(ctrlByte2, MaskOpCode, h.OpCode)
	ctrlByte2 = zetamath.SetBits(ctrlByte2, MaskCtrlReserved, h.Reserved)
	data[2] = ctrlByte2

	// set data flags
	data[3] = h.DataFlags

	return data, nil
}

// DecodeFromBytes decodes the memo header from the given data
func (h *Header) DecodeFromBytes(data []byte) error {
	// memo data must be longer than the header size
	if len(data) <= HeaderSize {
		return errors.New("memo is too short")
	}

	// byte-0 is the memo identifier
	if data[0] != Identifier {
		return fmt.Errorf("invalid memo identifier: %d", data[0])
	}

	// extract version #, encoding format
	ctrlByte1 := data[1]
	h.Version = zetamath.GetBits(ctrlByte1, MaskVersion)
	h.EncodingFormat = zetamath.GetBits(ctrlByte1, MaskEncodingFormat)

	// extract operation code, reserved bits
	ctrlByte2 := data[2]
	h.OpCode = zetamath.GetBits(ctrlByte2, MaskOpCode)
	h.Reserved = zetamath.GetBits(ctrlByte2, MaskCtrlReserved)

	// extract data flags
	h.DataFlags = data[3]

	// validate header
	return h.Validate()
}

// Validate checks if the memo header is valid
func (h *Header) Validate() error {
	if h.Version != 0 {
		return fmt.Errorf("invalid memo version: %d", h.Version)
	}

	if h.EncodingFormat >= EncodingFmtMax {
		return fmt.Errorf("invalid encoding format: %d", h.EncodingFormat)
	}

	if h.OpCode >= OpCodeMax {
		return fmt.Errorf("invalid operation code: %d", h.OpCode)
	}

	// reserved control bits must be zero
	if h.Reserved != 0 {
		return fmt.Errorf("reserved control bits are not zero: %d", h.Reserved)
	}

	return nil
}
