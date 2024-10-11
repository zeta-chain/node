package memo

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/crypto"
	zetamath "github.com/zeta-chain/node/pkg/math"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// Enum of the bit position of each memo fields
const (
	bitPosReceiver      uint8 = 0 // receiver
	bitPosPayload       uint8 = 1 // payload
	bitPosRevertAddress uint8 = 2 // revertAddress
	bitPosAbortAddress  uint8 = 3 // abortAddress
	bitPosRevertMessage uint8 = 4 // revertMessage
)

const (
	// MaskFlagsReserved is the mask for reserved data flags
	MaskFlagsReserved = 0b11100000
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
func (f *FieldsV0) Pack(opCode uint8, encodingFormat uint8, dataFlags uint8) ([]byte, error) {
	// validate fields
	err := f.Validate(opCode, dataFlags)
	if err != nil {
		return nil, err
	}

	codec, err := GetCodec(encodingFormat)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get codec")
	}

	return f.packFields(codec, dataFlags)
}

// Unpack decodes the memo fields
func (f *FieldsV0) Unpack(opCode uint8, encodingFormat uint8, dataFlags uint8, data []byte) error {
	codec, err := GetCodec(encodingFormat)
	if err != nil {
		return errors.Wrap(err, "unable to get codec")
	}

	err = f.unpackFields(codec, dataFlags, data)
	if err != nil {
		return err
	}

	return f.Validate(opCode, dataFlags)
}

// Validate checks if the fields are valid
func (f *FieldsV0) Validate(opCode uint8, dataFlags uint8) error {
	// receiver address must be a valid address
	if zetamath.IsBitSet(dataFlags, bitPosReceiver) && crypto.IsEmptyAddress(f.Receiver) {
		return errors.New("receiver address is empty")
	}

	// payload is not allowed for deposit operation
	if opCode == OpCodeDeposit && len(f.Payload) > 0 {
		return errors.New("payload is not allowed for deposit operation")
	}

	// abort address must be a valid address
	if zetamath.IsBitSet(dataFlags, bitPosAbortAddress) && !common.IsHexAddress(f.RevertOptions.AbortAddress) {
		return errors.New("invalid abort address")
	}

	// revert message is not allowed when CallOnRevert is false
	// 1. it's a good-to-have check to make the fields semantically correct.
	// 2. unpacking won't hit this error as the codec will catch it earlier.
	if !f.RevertOptions.CallOnRevert && len(f.RevertOptions.RevertMessage) > 0 {
		return errors.New("revert message is not allowed when CallOnRevert is false")
	}

	return nil
}

// DataFlags build the data flags from actual fields
func (f *FieldsV0) DataFlags() uint8 {
	var dataFlags uint8

	// set 'receiver' flag if provided
	if !crypto.IsEmptyAddress(f.Receiver) {
		zetamath.SetBit(&dataFlags, bitPosReceiver)
	}

	// set 'payload' flag if provided
	if len(f.Payload) > 0 {
		zetamath.SetBit(&dataFlags, bitPosPayload)
	}

	// set 'revertAddress' flag if provided
	if f.RevertOptions.RevertAddress != "" {
		zetamath.SetBit(&dataFlags, bitPosRevertAddress)
	}

	// set 'abortAddress' flag if provided
	if f.RevertOptions.AbortAddress != "" {
		zetamath.SetBit(&dataFlags, bitPosAbortAddress)
	}

	// set 'revertMessage' flag if provided
	if f.RevertOptions.CallOnRevert {
		zetamath.SetBit(&dataFlags, bitPosRevertMessage)
	}

	return dataFlags
}

// packFieldsV0 packs the memo fields for version 0
func (f *FieldsV0) packFields(codec Codec, dataFlags uint8) ([]byte, error) {
	// add 'receiver' argument optionally
	if zetamath.IsBitSet(dataFlags, bitPosReceiver) {
		codec.AddArguments(ArgReceiver(f.Receiver))
	}

	// add 'payload' argument optionally
	if zetamath.IsBitSet(dataFlags, bitPosPayload) {
		codec.AddArguments(ArgPayload(f.Payload))
	}

	// add 'revertAddress' argument optionally
	if zetamath.IsBitSet(dataFlags, bitPosRevertAddress) {
		codec.AddArguments(ArgRevertAddress(f.RevertOptions.RevertAddress))
	}

	// add 'abortAddress' argument optionally
	abortAddress := common.HexToAddress(f.RevertOptions.AbortAddress)
	if zetamath.IsBitSet(dataFlags, bitPosAbortAddress) {
		codec.AddArguments(ArgAbortAddress(abortAddress))
	}

	// add 'revertMessage' argument optionally
	if f.RevertOptions.CallOnRevert {
		codec.AddArguments(ArgRevertMessage(f.RevertOptions.RevertMessage))
	}

	// pack the codec arguments into data
	data, err := codec.PackArguments()
	if err != nil { // never happens
		return nil, errors.Wrap(err, "failed to pack arguments")
	}

	return data, nil
}

// unpackFields unpacks the memo fields for version 0
func (f *FieldsV0) unpackFields(codec Codec, dataFlags byte, data []byte) error {
	// add 'receiver' argument optionally
	if zetamath.IsBitSet(dataFlags, bitPosReceiver) {
		codec.AddArguments(ArgReceiver(&f.Receiver))
	}

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

	// unpack the data (after flags) into codec arguments
	err := codec.UnpackArguments(data)
	if err != nil {
		return errors.Wrap(err, "failed to unpack arguments")
	}

	// convert abort address to string
	if !crypto.IsEmptyAddress(abortAddress) {
		f.RevertOptions.AbortAddress = abortAddress.Hex()
	}

	return nil
}
