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

// Encoder is the interface for outbound memo encoders
type Encoder func(memo *InboundMemo) ([]byte, error)

// Decoder is the interface for inbound memo decoders
type Decoder func(data []byte, memo *InboundMemo) error

// memoEncoderRegistry contains all registered memo encoders
var memoEncoderRegistry = map[uint8]Encoder{
	0: FieldsEncoderV0,
}

// memoDecoderRegistry contains all registered memo decoders
var memoDecoderRegistry = map[uint8]Decoder{
	0: FieldsDecoderV0,
}

// InboundMemo represents the memo structure for non-EVM chains
type InboundMemo struct {
	// Version is the memo Version
	Version uint8

	// EncodingFormat is the memo encoding format
	EncodingFormat uint8

	// OpCode is the inbound operation code
	OpCode uint8

	// FieldsV0 contains the memo fields V0
	// Note: add a MemoFieldsV1 if major change is needed in the future
	FieldsV0
}

// EncodeMemoToBytes encodes a InboundMemo struct to raw bytes
func EncodeMemoToBytes(memo *InboundMemo) ([]byte, error) {
	// get registered memo encoder by version
	encoder, found := memoEncoderRegistry[memo.Version]
	if !found {
		return nil, fmt.Errorf("encoder not found for memo version: %d", memo.Version)
	}

	// encode memo basics
	basics := EncodeMemoBasics(memo)

	// encode memo fields using the encoder
	data, err := encoder(memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode memo fields")
	}

	return append(basics, data...), nil
}

// EncodeMemoBasics encodes the version, encoding format and operation code
func EncodeMemoBasics(memo *InboundMemo) []byte {
	// create 3-byte header
	head := make([]byte, HeaderSize)

	// set byte-0 as memo identifier
	head[0] = MemoIdentifier

	// set version # and encoding format
	var ctrlByte byte
	ctrlByte = zetamath.SetBits(ctrlByte, BitMaskVersion, memo.Version)
	ctrlByte = zetamath.SetBits(ctrlByte, BitMaskEncodingFormat, memo.EncodingFormat)
	ctrlByte = zetamath.SetBits(ctrlByte, BitMaskOpCode, memo.OpCode)

	// set ctrlByte to byte-1
	head[1] = ctrlByte

	return head
}

// DecodeMemoFromBytes decodes a InboundMemo struct from raw bytes
//
// Returns an error if given data is not a valid memo
func DecodeMemoFromBytes(data []byte) (*InboundMemo, error) {
	memo := &InboundMemo{}

	// decode memo basics
	err := DecodeMemoBasics(data, memo)
	if err != nil {
		return nil, err
	}

	// get registered memo decoder by version
	decoder, found := memoDecoderRegistry[memo.Version]
	if !found {
		return nil, fmt.Errorf("decoder not found for memo version: %d", memo.Version)
	}

	// decode memo fields using the decoer
	err = decoder(data, memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode memo fields")
	}

	return memo, nil
}

// DecodeMemoBasics decodes version, encoding format and operation code
func DecodeMemoBasics(data []byte, memo *InboundMemo) error {
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
