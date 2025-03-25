package sui

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strconv"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

// Withdrawal represents data for a Sui withdraw event
type Withdrawal struct {
	CoinType CoinType
	Amount   math.Uint
	Sender   string
	Receiver string
	Nonce    uint64
}

func (d *Withdrawal) IsGas() bool {
	return d.CoinType == SUI
}

func parseWithdrawal(event models.SuiEventResponse, eventType EventType) (Withdrawal, error) {
	if eventType != WithdrawEvent {
		return Withdrawal{}, errors.Errorf("invalid event type %q", eventType)
	}

	parsedJSON := event.ParsedJson

	coinType, err := extractStr(parsedJSON, "coin_type")
	if err != nil {
		return Withdrawal{}, err
	}

	amountRaw, err := extractStr(parsedJSON, "amount")
	if err != nil {
		return Withdrawal{}, err
	}

	amount, err := math.ParseUint(amountRaw)
	if err != nil {
		return Withdrawal{}, errors.Wrap(err, "unable to parse amount")
	}

	sender, err := extractStr(parsedJSON, "sender")
	if err != nil {
		return Withdrawal{}, err
	}

	receiver, err := extractStr(parsedJSON, "receiver")
	if err != nil {
		return Withdrawal{}, err
	}

	nonceRaw, err := extractStr(parsedJSON, "nonce")
	if err != nil {
		return Withdrawal{}, err
	}

	nonce, err := strconv.ParseUint(nonceRaw, 10, 64)
	if err != nil {
		return Withdrawal{}, errors.Wrap(err, "unable to parse nonce")
	}

	return Withdrawal{
		CoinType: CoinType(coinType),
		Amount:   amount,
		Sender:   sender,
		Receiver: receiver,
		Nonce:    nonce,
	}, nil
}

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
