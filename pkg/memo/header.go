package memo

import (
	"fmt"

	"github.com/pkg/errors"

	zetabits "github.com/zeta-chain/node/pkg/math/bits"
)

type OpCode uint8

const (
	// Identifier is the ASCII code of 'Z' (0x5A)
	Identifier byte = 0x5A

	// HeaderSize is the size of the memo header: [identifier + ctrlByte1+ ctrlByte2 + dataFlags]
	HeaderSize = 4

	// maskVersion is the mask for the version bits(upper 4 bits)
	maskVersion byte = 0b11110000

	// maskEncodingFormat is the mask for the encoding format bits(lower 4 bits)
	maskEncodingFormat byte = 0b00001111

	// maskOpCode is the mask for the operation code bits(upper 4 bits)
	maskOpCode byte = 0b11110000

	// maskCtrlReserved is the mask for reserved control bits (lower 4 bits)
	maskCtrlReserved byte = 0b00001111
)

// Enum for non-EVM chain inbound operation code (4 bits)
const (
	OpCodeDeposit        OpCode = 0b0000 // operation 'deposit'
	OpCodeDepositAndCall OpCode = 0b0001 // operation 'deposit_and_call'
	OpCodeCall           OpCode = 0b0010 // operation 'call'
	OpCodeInvalid        OpCode = 0b0011 // invalid operation code
)

// Header represent the memo header
type Header struct {
	// Version is the memo Version
	Version uint8

	// EncodingFmt is the memo encoding format
	EncodingFmt EncodingFormat

	// OpCode is the inbound operation code
	OpCode OpCode

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
	ctrlByte1 = zetabits.SetBits(ctrlByte1, maskVersion, h.Version)
	ctrlByte1 = zetabits.SetBits(ctrlByte1, maskEncodingFormat, byte(h.EncodingFmt))
	data[1] = ctrlByte1

	// set operation code, reserved bits
	var ctrlByte2 byte
	ctrlByte2 = zetabits.SetBits(ctrlByte2, maskOpCode, byte(h.OpCode))
	ctrlByte2 = zetabits.SetBits(ctrlByte2, maskCtrlReserved, h.Reserved)
	data[2] = ctrlByte2

	// set data flags
	data[3] = h.DataFlags

	return data, nil
}

// DecodeFromBytes decodes the memo header from the given data
func (h *Header) DecodeFromBytes(data []byte) error {
	// memo data must be longer than the header size
	if len(data) < HeaderSize {
		return errors.New("memo is too short")
	}

	// byte-0 is the memo identifier
	if data[0] != Identifier {
		return fmt.Errorf("invalid memo identifier: %d", data[0])
	}

	// extract version #, encoding format
	ctrlByte1 := data[1]
	h.Version = zetabits.GetBits(ctrlByte1, maskVersion)
	h.EncodingFmt = EncodingFormat(zetabits.GetBits(ctrlByte1, maskEncodingFormat))

	// extract operation code, reserved bits
	ctrlByte2 := data[2]
	h.OpCode = OpCode(zetabits.GetBits(ctrlByte2, maskOpCode))
	h.Reserved = zetabits.GetBits(ctrlByte2, maskCtrlReserved)

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

	if h.EncodingFmt >= EncodingFmtInvalid {
		return fmt.Errorf("invalid encoding format: %d", h.EncodingFmt)
	}

	if h.OpCode >= OpCodeInvalid {
		return fmt.Errorf("invalid operation code: %d", h.OpCode)
	}

	// reserved control bits must be zero
	if h.Reserved != 0 {
		return fmt.Errorf("reserved control bits are not zero: %d", h.Reserved)
	}

	return nil
}
