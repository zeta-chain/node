package memo

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	// version0 is the latest version of the memo
	version0 uint8 = 0
)

// InboundMemo represents the inbound memo structure for non-EVM chains
type InboundMemo struct {
	// Header contains the memo header
	Header

	// FieldsV0 contains the memo fields V0
	// Note: add a FieldsV1 if breaking change is needed in the future
	FieldsV0
}

// EncodeToBytes encodes a InboundMemo struct to raw bytes
//
// Note:
//   - Any provided 'DataFlags' is ignored as they are calculated based on the fields set in the memo.
//   - The 'RevertGasLimit' is not used for now for non-EVM chains.
func (m *InboundMemo) EncodeToBytes() ([]byte, error) {
	// build fields flags
	dataFlags := m.FieldsV0.DataFlags()
	m.Header.DataFlags = dataFlags

	// encode head
	head, err := m.Header.EncodeToBytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode memo header")
	}

	// encode fields based on version
	var data []byte
	switch m.Version {
	case version0:
		data, err = m.FieldsV0.Pack(m.OpCode, m.EncodingFmt, dataFlags)
	default:
		return nil, fmt.Errorf("invalid memo version: %d", m.Version)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to pack memo fields version: %d", m.Version)
	}

	return append(head, data...), nil
}

// DecodeFromBytes decodes a InboundMemo struct from raw bytes
//
// Returns:
//   - [memo, true, nil] if given data is successfully decoded as a memo.
//   - [nil,  true, err] if given data is successfully decoded as a memo but contains improper field values.
//   - [nil, false, nil] if given data can't be decoded as a standard memo, leave the error to nil.
//
// Note: we won't have to differentiate between the two 'true' cases if legacy memo phase out is completed.
func DecodeFromBytes(data []byte) (*InboundMemo, bool, error) {
	memo := &InboundMemo{}

	// decode header
	// the 'err' is interpreted as 'not a standard memo', so
	// there is no need to propagate the 'err' to the caller
	err := memo.Header.DecodeFromBytes(data)
	if err != nil {
		return nil, false, nil
	}

	// decode fields based on version
	switch memo.Version {
	case version0:
		// unpack fields
		err = memo.FieldsV0.Unpack(memo.EncodingFmt, memo.Header.DataFlags, data[HeaderSize:])
		if err != nil {
			return nil, true, errors.Wrap(err, "failed to unpack memo FieldsV0")
		}

		// validate fields
		// Note: a well-formatted memo may still contain improper field values
		err = memo.FieldsV0.Validate(memo.OpCode, memo.Header.DataFlags)
		if err != nil {
			return nil, true, errors.Wrap(err, "failed to validate memo FieldsV0")
		}
	default:
		// unreachable code
		// version is validated when decoding the header
		return nil, true, fmt.Errorf("invalid memo version: %d", memo.Version)
	}

	return memo, true, nil
}

// DecodeLegacyMemoHex decodes hex encoded memo message into address and calldata
//
// The layout of legacy memo is: [20-byte address, variable calldata]
func DecodeLegacyMemoHex(message string) (common.Address, []byte, error) {
	if len(message) == 0 {
		return common.Address{}, nil, nil
	}

	data, err := hex.DecodeString(message)
	if err != nil {
		return common.Address{}, nil, errors.Wrap(err, "message should be a hex encoded string")
	}

	if len(data) < common.AddressLength {
		return common.Address{}, data, nil
	}

	address := common.BytesToAddress(data[:common.AddressLength])
	data = data[common.AddressLength:]
	return address, data, nil
}
