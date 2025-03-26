package signer

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type rawPayload struct {
	TypeArguments []string `json:"typeArguments"`
	Objects       [][]byte `json:"objects"`
	Message       string   `json:"message"` // base64-encoded
}

// parseWithdrawAndCallPayload parses the ABI encoded payload for withdraw and call
func parseWithdrawAndCallPayload(payload []byte) ([]string, []string, []byte, error) {
	abiType, err := abi.NewType("tuple", "struct Payload", []abi.ArgumentMarshaling{
		{Name: "typeArguments", Type: "string[]"},
		{Name: "objects", Type: "bytes32[]"},
		{Name: "message", Type: "bytes"},
	})

	if err != nil {
		return nil, nil, nil, err
	}

	abiArgs := abi.Arguments{
		{Type: abiType, Name: "payload"},
	}

	unpacked, err := abiArgs.Unpack(payload)
	if err != nil {
		return nil, nil, nil, err
	}

	jsonData, err := json.Marshal(unpacked[0])
	if err != nil {
		return nil, nil, nil, err
	}

	var rawPayload rawPayload
	if err := json.Unmarshal(jsonData, &rawPayload); err != nil {
		return nil, nil, nil, err
	}

	// Convert object [][]byte to []string (hex-encoded)
	objects := make([]string, len(rawPayload.Objects))
	for i, obj := range rawPayload.Objects {
		objects[i] = "0x" + hex.EncodeToString(obj)
	}

	// Decode base64 message
	messageBytes, err := base64.StdEncoding.DecodeString(rawPayload.Message)
	if err != nil {
		return nil, nil, nil, err
	}

	return rawPayload.TypeArguments, objects, messageBytes, nil
}
