package observer

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	zetasolrpc "github.com/zeta-chain/node/zetaclient/chains/solana/rpc"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// Observer is the observer for the Solana chain
type Observer struct {
	// base.Observer implements the base chain observer
	*base.Observer

	// solClient is the Solana RPC client that interacts with the Solana chain
	solClient interfaces.SolanaRPCClient

	// gatewayID is the program ID of gateway program on Solana chain
	gatewayID solana.PublicKey

	// pda is the program derived address of the gateway program
	pda solana.PublicKey

	// finalizedTxResults indexes tx results with the outbound hash
	finalizedTxResults map[string]*rpc.GetTransactionResult
}

// New Observer constructor
func New(baseObserver *base.Observer, solClient interfaces.SolanaRPCClient, gatewayAddress string) (*Observer, error) {
	// parse gateway ID and PDA
	gatewayID, pda, err := contracts.ParseGatewayWithPDA(gatewayAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse gateway address %s", gatewayAddress)
	}

	// create solana observer
	ob := &Observer{
		Observer:           baseObserver,
		solClient:          solClient,
		gatewayID:          gatewayID,
		pda:                pda,
		finalizedTxResults: make(map[string]*rpc.GetTransactionResult),
	}

	ob.Observer.LoadLastTxScanned()

	return ob, nil
}

// LoadLastTxScanned loads the last scanned tx from the database.
func (ob *Observer) LoadLastTxScanned() error {
	ob.Observer.LoadLastTxScanned()
	ob.Logger().Chain.Info().
		Str(logs.FieldTx, ob.LastTxScanned()).
		Msg("chain starts scanning from tx")

	return nil
}

// SetTxResult sets the tx result for the given nonce
func (ob *Observer) SetTxResult(nonce uint64, result *rpc.GetTransactionResult) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.finalizedTxResults[ob.OutboundID(nonce)] = result
}

// GetTxResult returns the tx result for the given nonce
func (ob *Observer) GetTxResult(nonce uint64) *rpc.GetTransactionResult {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.finalizedTxResults[ob.OutboundID(nonce)]
}

// IsTxFinalized returns true if there is a finalized tx for nonce
func (ob *Observer) IsTxFinalized(nonce uint64) bool {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.finalizedTxResults[ob.OutboundID(nonce)] != nil
}

// CheckRPCStatus checks the RPC status of the Solana chain
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := zetasolrpc.HealthCheck(ctx, ob.solClient)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc status")
	}

	metrics.ReportBlockLatency(ob.Chain().Name, blockTime)

	return nil
}
