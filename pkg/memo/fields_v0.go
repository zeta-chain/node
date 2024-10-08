package memo

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

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

// FieldsV0 contains the data fields of the inbound memo V0
type FieldsV0 struct {
	// Receiver is the ZEVM receiver address
	Receiver common.Address

	// Payload is the calldata passed to ZEVM contract call
	Payload []byte

	// RevertOptions is the options for cctx revert handling
	RevertOptions *crosschaintypes.RevertOptions
}

// FieldsEncoderV0 is the encoder for outbound memo fields V0
func FieldsEncoderV0(memo *InboundMemo) ([]byte, error) {
	codec, err := GetCodec(memo.EncodingFormat)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get codec")
	}

	return PackMemoFieldsV0(codec, memo)
}

// FieldsDecoderV0 is the decoder for inbound memo fields V0
func FieldsDecoderV0(data []byte, memo *InboundMemo) error {
	codec, err := GetCodec(memo.EncodingFormat)
	if err != nil {
		return errors.Wrap(err, "unable to get codec")
	}
	return UnpackMemoFieldsV0(codec, data[HeaderSize:], memo)
}

// PackMemoFieldsV0 packs the memo fields for version 0
func PackMemoFieldsV0(codec Codec, memo *InboundMemo) ([]byte, error) {
	// create data flags byte
	dataFlags := byte(0)

	// add 'receiver' as the first argument
	codec.AddArguments(ArgReceiver(memo.Receiver))

	// add 'payload' argument optionally
	if len(memo.Payload) > 0 {
		zetamath.SetBit(&dataFlags, bitPosPayload)
		codec.AddArguments(ArgPayload(memo.Payload))
	}

	if memo.RevertOptions != nil {
		// add 'revertAddress' argument optionally
		if memo.RevertOptions.RevertAddress != "" {
			zetamath.SetBit(&dataFlags, bitPosRevertAddress)
			codec.AddArguments(ArgRevertAddress(memo.RevertOptions.RevertAddress))
		}

		// add 'abortAddress' argument optionally
		if memo.RevertOptions.AbortAddress != "" {
			zetamath.SetBit(&dataFlags, bitPosAbortAddress)
			codec.AddArguments(ArgAbortAddress(common.HexToAddress(memo.RevertOptions.AbortAddress)))
		}

		// add 'revertMessage' argument optionally
		if memo.RevertOptions.CallOnRevert {
			zetamath.SetBit(&dataFlags, bitPosRevertMessage)
			codec.AddArguments(ArgRevertMessage(memo.RevertOptions.RevertMessage))
		}
	}

	// pack the codec arguments into data
	data, err := codec.PackArguments()
	if err != nil {
		return nil, err
	}

	return append([]byte{dataFlags}, data...), nil
}

// UnpackMemoFieldsV0 unpacks the memo fields for version 0
func UnpackMemoFieldsV0(codec Codec, data []byte, memo *InboundMemo) error {
	// byte-2 contains data flags
	dataFlags := data[2]

	// add 'receiver' as the first argument
	codec.AddArguments(ArgReceiver(&memo.Receiver))

	// add 'payload' argument optionally
	if zetamath.IsBitSet(dataFlags, bitPosPayload) {
		codec.AddArguments(ArgPayload(&memo.Payload))
	}

	// add 'revertAddress' argument optionally
	if zetamath.IsBitSet(dataFlags, bitPosRevertAddress) {
		codec.AddArguments(ArgRevertAddress(&memo.RevertOptions.RevertAddress))
	}

	// add 'abortAddress' argument optionally
	var abortAddress common.Address
	if zetamath.IsBitSet(dataFlags, bitPosRevertMessage) {
		codec.AddArguments(ArgAbortAddress(&abortAddress))
	}

	// add 'revertMessage' argument optionally
	memo.RevertOptions.CallOnRevert = zetamath.IsBitSet(dataFlags, bitPosAbortAddress)
	if memo.RevertOptions.CallOnRevert {
		codec.AddArguments(ArgRevertMessage(&memo.RevertOptions.RevertMessage))
	}

	// unpack the data (after header) into codec arguments
	err := codec.UnpackArguments(data[HeaderSize:])
	if err != nil {
		return err
	}

	// convert abort address to string
	memo.RevertOptions.AbortAddress = abortAddress.Hex()

	return nil
}
