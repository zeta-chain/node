package ton

import (
	_ "embed"
	"encoding/json"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
)

// https://github.com/zeta-chain/protocol-contracts-ton
// `make compile`
//
//go:embed gateway.compiled.json
var tonGatewayCodeJSON []byte

// GetGatewayCodeAndState returns TON Gateway code and initial state cells.
// Returns (code, state, error).
func GetGatewayCodeAndState(tss eth.Address) (*boc.Cell, *boc.Cell, error) {
	code, err := getGatewayCode()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get TON Gateway code")
	}

	state, err := buildGatewayState(tss)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to build TON Gateway state")
	}

	return code, state, nil
}

func getGatewayCode() (*boc.Cell, error) {
	var code struct {
		Hex string `json:"hex"`
	}

	if err := json.Unmarshal(tonGatewayCodeJSON, &code); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal TON Gateway code")
	}

	cells, err := boc.DeserializeBocHex(code.Hex)
	if err != nil {
		return nil, errors.Wrap(err, "unable to deserialize TON Gateway code")
	}

	if len(cells) != 1 {
		return nil, errors.New("invalid cells count")
	}

	return cells[0], nil
}

// buildGatewayState returns TON Gateway initial state cell
func buildGatewayState(tss eth.Address) (*boc.Cell, error) {
	const evmAddressBits = 20 * 8

	tssSlice := boc.NewBitString(evmAddressBits)
	if err := tssSlice.WriteBytes(tss.Bytes()); err != nil {
		return nil, errors.Wrap(err, "unable to convert TSS address to ton slice")
	}

	cell := boc.NewCell()

	err := errCollect(
		cell.WriteBit(true),           // deposits_enabled
		cell.WriteUint(0, 4),          // total_locked varUint
		cell.WriteUint(0, 4),          // fees
		cell.WriteUint(0, 32),         // seqno
		cell.WriteBitString(tssSlice), // tss_address
	)

	if err != nil {
		return nil, errors.Wrap(err, "unable to write TON Gateway state cell")
	}

	return cell, nil
}

func errCollect(errs ...error) error {
	for i, err := range errs {
		if err != nil {
			return errors.Wrapf(err, "error at index %d", i)
		}
	}

	return nil
}
