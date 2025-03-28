package sui

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type rawPayload struct {
	TypeArguments []string `json:"typeArguments"`
	Objects       [][]byte `json:"objects"`
	Message       string   `json:"message"` // base64-encoded
}

// ParseWithdrawAndCallPayload parses the ABI encoded payload for withdraw and call
func ParseWithdrawAndCallPayload(payload []byte) ([]string, []string, []byte, error) {
	abiType, err := abi.NewType("tuple", "struct Payload", []abi.ArgumentMarshaling{
		{Name: "typeArguments", Type: "string[]"},
		{Name: "objects", Type: "bytes32[]"},
		{Name: "message", Type: "bytes"},
	})
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to create ABI type")
	}
	abiArgs := abi.Arguments{
		{Type: abiType, Name: "payload"},
	}

	unpacked, err := abiArgs.Unpack(payload)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to unpack ABI arguments")
	}

	jsonData, err := json.Marshal(unpacked[0])
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to marshal ABI arguments")
	}

	var rawPayload rawPayload
	if err := json.Unmarshal(jsonData, &rawPayload); err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to unmarshal formatted JSON")
	}

	// Convert object [][]byte to []string (hex-encoded)
	objects := make([]string, len(rawPayload.Objects))
	for i, obj := range rawPayload.Objects {
		objects[i] = "0x" + hex.EncodeToString(obj)
	}

	// Decode base64 message
	messageBytes, err := base64.StdEncoding.DecodeString(rawPayload.Message)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to decode base64 message")
	}

	return rawPayload.TypeArguments, objects, messageBytes, nil
}

// FormatWithdrawAndCallPayload formats the withdraw and call payload using ABI encoding
func FormatWithdrawAndCallPayload(typeArguments []string, objects []string, message []byte) ([]byte, error) {
	// prepare type
	abiType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "typeArguments", Type: "string[]"},
		{Name: "objects", Type: "bytes32[]"},
		{Name: "message", Type: "bytes"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create ABI type")
	}
	abiArgs := abi.Arguments{{Type: abiType}}

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
		TypeArguments: typeArguments,
		Objects:       objectsBytes,
		Message:       message,
	}
	return abiArgs.Pack(payload)
}
