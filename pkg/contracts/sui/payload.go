package sui

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

// abiType is the ABI type for the withdraw and call payload
// error is ignored as it is a constant
var abiType, _ = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
	{Name: "typeArguments", Type: "string[]"},
	{Name: "objects", Type: "bytes32[]"},
	{Name: "message", Type: "bytes"},
})

var abiArgs = abi.Arguments{{Type: abiType}}

// CallPayload represents the payload provided for a custom contract call in the Sui blockchain
// it contains data allowing to build a Sui PTB transaction to call on_call on a target contract
type CallPayload struct {
	// TypeArgs are custom type arguments provided for on_call
	TypeArgs []string

	// ObjectIDs are the object IDs to be used in the on_call
	ObjectIDs []string

	// Message is a generic byte array passed as last argument to on_call
	Message []byte
}

func NewCallPayload(typeArgs []string, objectIDs []string, message []byte) CallPayload {
	return CallPayload{
		TypeArgs:  typeArgs,
		ObjectIDs: objectIDs,
		Message:   message,
	}
}

// UnpackABI parses the ABI encoded payload for calls
func (c *CallPayload) UnpackABI(payload []byte) error {
	unpacked, err := abiArgs.Unpack(payload)
	if err != nil {
		return errors.Wrapf(ErrInvalidPayload, "unable to unpack ABI encoded payload (%x): %s", payload, err)
	}

	jsonData, err := json.Marshal(unpacked[0])
	if err != nil {
		return errors.Wrapf(ErrInvalidPayload, "unable to marshal unpacked payload (%x): %s", payload, err)
	}

	// raw payload format parsed from ABI
	var rawPayload struct {
		TypeArguments []string `json:"typeArguments"`
		Objects       [][]byte `json:"objects"`
		Message       string   `json:"message"` // base64-encoded
	}
	if err := json.Unmarshal(jsonData, &rawPayload); err != nil {
		return errors.Wrapf(ErrInvalidPayload, "unable to unmarshal formatted JSON for payload (%x): %s", payload, err)
	}

	// Convert object [][]byte to []string (hex-encoded)
	objects := make([]string, len(rawPayload.Objects))
	for i, obj := range rawPayload.Objects {
		objects[i] = "0x" + hex.EncodeToString(obj)
	}

	// Decode base64 message
	messageBytes, err := base64.StdEncoding.DecodeString(rawPayload.Message)
	if err != nil {
		// should never happen, guaranteed by json.Marshal()
		return errors.Wrapf(ErrInvalidPayload, "unable to decode base64 message: %s", err)
	}

	// Set the fields
	c.TypeArgs = rawPayload.TypeArguments
	c.ObjectIDs = objects
	c.Message = messageBytes

	return nil
}

// PackABI formats the call payload using ABI encoding
func (c *CallPayload) PackABI() ([]byte, error) {
	objects := c.ObjectIDs

	// build fixed [32]byte array
	objectsBytes := make([][32]byte, len(objects))
	for i, obj := range objects {
		objBytes, err := hex.DecodeString(strings.TrimPrefix(obj, "0x"))
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode hex object")
		}
		if len(objBytes) != 32 {
			return nil, fmt.Errorf("object at index %d is not 32 bytes", i)
		}
		copy(objectsBytes[i][:], objBytes)
	}

	// format
	payload := struct {
		TypeArguments []string
		Objects       [][32]byte
		Message       []byte
	}{
		TypeArguments: c.TypeArgs,
		Objects:       objectsBytes,
		Message:       c.Message,
	}
	return abiArgs.Pack(payload)
}
