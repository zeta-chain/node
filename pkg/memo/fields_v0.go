package memo

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/crypto"
	zetamath "github.com/zeta-chain/node/pkg/math"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// Enum of the bit position of each memo fields
const (
	bitPosPayload       uint8 = 0 // payload
	bitPosRevertAddress uint8 = 1 // revertAddress
	bitPosAbortAddress  uint8 = 2 // abortAddress
	bitPosRevertMessage uint8 = 3 // revertMessage
)

const (
	// MaskFlagReserved is the mask for reserved data flags
	MaskFlagReserved = 0b11110000
)

var _ Fields = (*FieldsV0)(nil)

// FieldsV0 contains the data fields of the inbound memo V0
type FieldsV0 struct {
	// Receiver is the ZEVM receiver address
	Receiver common.Address

	// Payload is the calldata passed to ZEVM contract call
	Payload []byte

	// RevertOptions is the options for cctx revert handling
	RevertOptions crosschaintypes.RevertOptions
}

// Pack encodes the memo fields
func (f *FieldsV0) Pack(encodingFormat uint8) ([]byte, error) {
	codec, err := GetCodec(encodingFormat)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get codec")
	}

	return f.packFields(codec)
}

// Unpack decodes the memo fields
func (f *FieldsV0) Unpack(data []byte, encodingFormat uint8) error {
	codec, err := GetCodec(encodingFormat)
	if err != nil {
		return errors.Wrap(err, "unable to get codec")
	}

	return f.unpackFields(codec, data)
}

// packFieldsV0 packs the memo fields for version 0
func (f *FieldsV0) packFields(codec Codec) ([]byte, error) {
	// create data flags byte
	dataFlags := byte(0)

	// add 'receiver' as the first argument
	codec.AddArguments(ArgReceiver(f.Receiver))

	// add 'payload' argument optionally
	if len(f.Payload) > 0 {
		zetamath.SetBit(&dataFlags, bitPosPayload)
		codec.AddArguments(ArgPayload(f.Payload))
	}

	// add 'revertAddress' argument optionally
	if f.RevertOptions.RevertAddress != "" {
		zetamath.SetBit(&dataFlags, bitPosRevertAddress)
		codec.AddArguments(ArgRevertAddress(f.RevertOptions.RevertAddress))
	}

	// add 'abortAddress' argument optionally
	abortAddress := common.HexToAddress(f.RevertOptions.AbortAddress)
	if !crypto.IsEmptyAddress(abortAddress) {
		zetamath.SetBit(&dataFlags, bitPosAbortAddress)
		codec.AddArguments(ArgAbortAddress(abortAddress))
	}

	// add 'revertMessage' argument optionally
	if f.RevertOptions.CallOnRevert {
		zetamath.SetBit(&dataFlags, bitPosRevertMessage)
		codec.AddArguments(ArgRevertMessage(f.RevertOptions.RevertMessage))
	}

	// pack the codec arguments into data
	data, err := codec.PackArguments()
	if err != nil { // never happens
		return nil, errors.Wrap(err, "failed to pack arguments")
	}

	return append([]byte{dataFlags}, data...), nil
}

// unpackFields unpacks the memo fields for version 0
func (f *FieldsV0) unpackFields(codec Codec, data []byte) error {
	// get data flags
	dataFlags := data[0]

	// add 'receiver' as the first argument
	codec.AddArguments(ArgReceiver(&f.Receiver))

	// add 'payload' argument optionally
	if zetamath.IsBitSet(dataFlags, bitPosPayload) {
		codec.AddArguments(ArgPayload(&f.Payload))
	}

	// add 'revertAddress' argument optionally
	if zetamath.IsBitSet(dataFlags, bitPosRevertAddress) {
		codec.AddArguments(ArgRevertAddress(&f.RevertOptions.RevertAddress))
	}

	// add 'abortAddress' argument optionally
	var abortAddress common.Address
	if zetamath.IsBitSet(dataFlags, bitPosAbortAddress) {
		codec.AddArguments(ArgAbortAddress(&abortAddress))
	}

	// add 'revertMessage' argument optionally
	f.RevertOptions.CallOnRevert = zetamath.IsBitSet(dataFlags, bitPosRevertMessage)
	if f.RevertOptions.CallOnRevert {
		codec.AddArguments(ArgRevertMessage(&f.RevertOptions.RevertMessage))
	}

	// all reserved flag bits must be zero
	reserved := zetamath.GetBits(dataFlags, MaskFlagReserved)
	if reserved != 0 {
		return fmt.Errorf("reserved flag bits are not zero: %d", reserved)
	}

	// unpack the data (after flags) into codec arguments
	err := codec.UnpackArguments(data[1:])
	if err != nil {
		return errors.Wrap(err, "failed to unpack arguments")
	}

	// convert abort address to string
	if !crypto.IsEmptyAddress(abortAddress) {
		f.RevertOptions.AbortAddress = abortAddress.Hex()
	}

	return nil
}
