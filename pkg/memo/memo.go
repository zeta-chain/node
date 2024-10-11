package memo

import (
	"fmt"

	"github.com/pkg/errors"
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
	case 0:
		data, err = m.FieldsV0.Pack(m.OpCode, m.EncodingFormat, dataFlags)
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
// Returns an error if given data is not a valid memo
func DecodeFromBytes(data []byte) (*InboundMemo, error) {
	memo := &InboundMemo{}

	// decode header
	err := memo.Header.DecodeFromBytes(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode memo header")
	}

	// decode fields based on version
	switch memo.Version {
	case 0:
		err = memo.FieldsV0.Unpack(memo.OpCode, memo.EncodingFormat, memo.Header.DataFlags, data[HeaderSize:])
	default:
		return nil, fmt.Errorf("invalid memo version: %d", memo.Version)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unpack memo fields version: %d", memo.Version)
	}

	return memo, nil
}
