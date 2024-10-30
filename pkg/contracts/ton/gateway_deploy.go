package ton

import (
	_ "embed"
	"encoding/json"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

//go:embed gateway.compiled.json
var gatewayCode []byte

// GatewayCode returns Gateway's code as cell
func GatewayCode() *boc.Cell {
	c, err := getGatewayCode()
	if err != nil {
		panic(err)
	}

	return c
}

// GatewayStateInit returns Gateway's stateInit as cell
func GatewayStateInit(authority ton.AccountID, tss eth.Address, depositsEnabled bool) *boc.Cell {
	c, err := buildGatewayStateInit(authority, tss, depositsEnabled)
	if err != nil {
		panic(err)
	}

	return c
}

func getGatewayCode() (*boc.Cell, error) {
	var code struct {
		Hex string `json:"hex"`
	}

	if err := json.Unmarshal(gatewayCode, &code); err != nil {
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

func buildGatewayStateInit(authority ton.AccountID, tss eth.Address, depositsEnabled bool) (*boc.Cell, error) {
	cell := boc.NewCell()

	err := ErrCollect(
		cell.WriteBit(depositsEnabled),              // deposits_enabled
		tlb.Marshal(cell, tlb.Coins(0)),             // total_locked
		cell.WriteUint(0, 32),                       // seqno
		cell.WriteBytes(tss.Bytes()),                // tss_address
		tlb.Marshal(cell, authority.ToMsgAddress()), // authority_address (TON)
	)

	if err != nil {
		return nil, errors.Wrap(err, "unable to write TON Gateway state cell")
	}

	return cell, nil
}
