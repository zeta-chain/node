// Package ton provider bindings for TON blockchain including Gateway contract wrapper.
package ton

import (
	"context"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// Gateway represents bindings for Zeta Gateway contract on TON
//
// Gateway.ParseTransaction parses Gateway transaction.
// The parser reads tx body cell and decodes it based on Operation code (op)
//   - inbound transactions: deposit, donate, depositAndCall
//   - outbound transactions: not implemented yet
//   - errors for all other transactions
//
// `Send*` methods work the same way by constructing (& signing) tx body cell that is expected by the contract
//
// @see https://github.com/zeta-chain/protocol-contracts-ton/blob/main/contracts/gateway.fc
type Gateway struct {
	accountID ton.AccountID
}

type MethodRunner interface {
	RunSmcMethod(ctx context.Context, acc ton.AccountID, method string, params tlb.VmStack) (uint32, tlb.VmStack, error)
}

type Filter func(*Transaction) bool

const (
	sizeOpCode  = 32
	sizeQueryID = 64
	sizeSeqno   = 32
)

var (
	ErrParse     = errors.New("unable to parse tx")
	ErrUnknownOp = errors.New("unknown op")
	ErrCast      = errors.New("unable to cast tx content")
)

// NewGateway Gateway constructor
func NewGateway(accountID ton.AccountID) *Gateway {
	return &Gateway{accountID}
}

// AccountID returns gateway address
func (gw *Gateway) AccountID() ton.AccountID {
	return gw.accountID
}

// ParseTransaction parses transaction to Transaction
func (gw *Gateway) ParseTransaction(tx ton.Transaction) (*Transaction, error) {
	if isOutbound(tx) {
		return gw.parseOutbound(tx)
	}

	return gw.parseInbound(tx)
}

// ParseAndFilter parses transaction and applies filter to it. Returns (tx, skip?, error)
// If parse fails due to known error, skip is set to true
func (gw *Gateway) ParseAndFilter(tx ton.Transaction, filter Filter) (*Transaction, bool, error) {
	parsedTX, err := gw.ParseTransaction(tx)
	switch {
	case errors.Is(err, ErrParse):
		return nil, true, nil
	case errors.Is(err, ErrUnknownOp):
		return nil, true, nil
	case err != nil:
		return nil, false, err
	}

	if !filter(parsedTX) {
		return parsedTX, true, nil
	}

	return parsedTX, false, nil
}

// ParseAndFilterMany parses and filters many txs.
func (gw *Gateway) ParseAndFilterMany(txs []ton.Transaction, filter Filter) []*Transaction {
	//goland:noinspection GoPreferNilSlice
	out := []*Transaction{}

	for i := range txs {
		tx, skip, err := gw.ParseAndFilter(txs[i], filter)
		if skip || err != nil {
			continue
		}

		out = append(out, tx)
	}

	return out
}

// FilterInbounds filters transactions with deposit operations
func FilterInbounds(tx *Transaction) bool { return tx.IsInbound() }

func parseInternalMessageBody(tx ton.Transaction) (*boc.Cell, error) {
	if !tx.Msgs.InMsg.Exists {
		return nil, errors.Wrap(ErrParse, "tx should have an internal message")
	}

	var (
		inMsg = tx.Msgs.InMsg.Value.Value
		body  = boc.Cell(inMsg.Body.Value)
	)

	return &body, nil
}

var zero = math.NewUint(0)

// GetTxFee returns maximum transaction fee for the given operation.
// Real fee may be lower.
func (gw *Gateway) GetTxFee(ctx context.Context, client MethodRunner, op Op) (math.Uint, error) {
	const (
		method  = "calculate_gas_fee"
		sumType = "VmStkTinyInt"
	)

	query := tlb.VmStack{{SumType: sumType, VmStkTinyInt: int64(op)}}

	exitCode, res, err := client.RunSmcMethod(ctx, gw.accountID, method, query)
	switch {
	case err != nil:
		return zero, err
	case exitCode != 0:
		return zero, errors.Errorf("calculate_gas_fee failed with exit code %d", exitCode)
	case len(res) == 0:
		return zero, errors.New("empty result")
	case res[0].SumType != sumType:
		return zero, errors.Errorf("res is not %s (got %s)", sumType, res[0].SumType)
	case res[0].VmStkTinyInt <= 0:
		return zero, errors.New("fee is zero or negative")
	}

	return math.NewUint(uint64(res[0].VmStkTinyInt)), nil
}
