package observer

import (
	"context"

	"github.com/gagliardetto/solana-go"

	"github.com/zeta-chain/zetacore/pkg/bg"
	"github.com/zeta-chain/zetacore/pkg/chains"
	solanacontract "github.com/zeta-chain/zetacore/pkg/contract/solana"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

var _ interfaces.ChainObserver = &Observer{}

// Observer is the observer for the Solana chain
type Observer struct {
	// base.Observer implements the base chain observer
	base.Observer

	// solClient is the Solana RPC client that interacts with the Solana chain
	solClient interfaces.SolanaRPCClient

	// gatewayID is the program ID of gateway program on Solana chain
	gatewayID solana.PublicKey

	// pda is the program derived address of the gateway program
	pdaID solana.PublicKey
}

// NewObserver returns a new Solana chain observer
func NewObserver(
	chain chains.Chain,
	solClient interfaces.SolanaRPCClient,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	dbpath string,
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
		ts,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// create solana observer
	ob := Observer{
		Observer:  *baseObserver,
		solClient: solClient,
		gatewayID: solana.MustPublicKeyFromBase58(chainParams.GatewayAddress),
	}

	// compute gateway PDA
	seed := []byte(solanacontract.PDASeed)
	ob.pdaID, _, err = solana.FindProgramAddress([][]byte{seed}, ob.gatewayID)
	if err != nil {
		return nil, err
	}

	// load btc chain observer DB
	err = ob.LoadDB(dbpath)
	if err != nil {
		return nil, err
	}

	return &ob, nil
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
	ob.Logger().Chain.Info().Msgf("observer is starting for chain %d", ob.Chain().ChainId)

	// watch Solana chain for incoming txs and post votes to zetacore
	bg.Work(ctx, ob.WatchInbound, bg.WithName("WatchInbound"), bg.WithLogger(ob.Logger().Inbound))

	// watch zetacore for Solana inbound trackers
	bg.Work(ctx, ob.WatchInboundTracker, bg.WithName("WatchInboundTracker"), bg.WithLogger(ob.Logger().Inbound))
}
