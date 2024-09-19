package observer

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/pkg/chains"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

var _ interfaces.ChainObserver = (*Observer)(nil)

// Observer is the observer for the Solana chain
type Observer struct {
	// base.Observer implements the base chain observer
	base.Observer

	// solClient is the Solana RPC client that interacts with the Solana chain
	solClient interfaces.SolanaRPCClient

	// gatewayID is the program ID of gateway program on Solana chain
	gatewayID solana.PublicKey

	// pda is the program derived address of the gateway program
	pda solana.PublicKey

	// finalizedTxResults indexes tx results with the outbound hash
	finalizedTxResults map[string]*rpc.GetTransactionResult
}

// NewObserver returns a new Solana chain observer
func NewObserver(
	chain chains.Chain,
	solClient interfaces.SolanaRPCClient,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	rpcAlertLatency int64,
	db *db.DB,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (*Observer, error) {
	// create base observer
	baseObserver, err := base.NewObserver(
		chain,
		chainParams,
		zetacoreClient,
		tss,
		base.DefaultBlockCacheSize,
		base.DefaultHeaderCacheSize,
		rpcAlertLatency,
		ts,
		db,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// parse gateway ID and PDA
	gatewayID, pda, err := contracts.ParseGatewayIDAndPda(chainParams.GatewayAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse gateway address %s", chainParams.GatewayAddress)
	}

	// create solana observer
	ob := &Observer{
		Observer:           *baseObserver,
		solClient:          solClient,
		gatewayID:          gatewayID,
		pda:                pda,
		finalizedTxResults: make(map[string]*rpc.GetTransactionResult),
	}

	ob.Observer.LoadLastTxScanned()

	return ob, nil
}

// SolClient returns the solana rpc client
func (ob *Observer) SolClient() interfaces.SolanaRPCClient {
	return ob.solClient
}

// WithSolClient attaches a new solana rpc client to the observer
func (ob *Observer) WithSolClient(client interfaces.SolanaRPCClient) {
	ob.solClient = client
}

// SetChainParams sets the chain params for the observer
// Note: chain params is accessed concurrently
func (ob *Observer) SetChainParams(params observertypes.ChainParams) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.WithChainParams(params)
}

// GetChainParams returns the chain params for the observer
// Note: chain params is accessed concurrently
func (ob *Observer) GetChainParams() observertypes.ChainParams {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.ChainParams()
}

// Start starts the Go routine processes to observe the Solana chain
func (ob *Observer) Start(ctx context.Context) {
	if noop := ob.Observer.Start(); noop {
		ob.Logger().Chain.Info().Msgf("observer is already started for chain %d", ob.Chain().ChainId)
		return
	}

	ob.Logger().Chain.Info().Msgf("observer is starting for chain %d", ob.Chain().ChainId)

	// watch Solana chain for incoming txs and post votes to zetacore
	bg.Work(ctx, ob.WatchInbound, bg.WithName("WatchInbound"), bg.WithLogger(ob.Logger().Inbound))

	// watch Solana chain for outbound trackers
	bg.Work(ctx, ob.WatchOutbound, bg.WithName("WatchOutbound"), bg.WithLogger(ob.Logger().Outbound))

	// watch Solana chain for fee rate and post to zetacore
	bg.Work(ctx, ob.WatchGasPrice, bg.WithName("WatchGasPrice"), bg.WithLogger(ob.Logger().GasPrice))

	// watch zetacore for Solana inbound trackers
	bg.Work(ctx, ob.WatchInboundTracker, bg.WithName("WatchInboundTracker"), bg.WithLogger(ob.Logger().Inbound))

	// watch RPC status of the Solana chain
	bg.Work(ctx, ob.watchRPCStatus, bg.WithName("watchRPCStatus"), bg.WithLogger(ob.Logger().Chain))
}

// LoadLastTxScanned loads the last scanned tx from the database.
func (ob *Observer) LoadLastTxScanned() error {
	ob.Observer.LoadLastTxScanned()
	ob.Logger().Chain.Info().Msgf("chain %d starts scanning from tx %s", ob.Chain().ChainId, ob.LastTxScanned())

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
