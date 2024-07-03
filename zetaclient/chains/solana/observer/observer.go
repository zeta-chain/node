package observer

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	clientcontext "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

type Observer struct {
	Tss            interfaces.TSSSigner
	zetacoreClient interfaces.ZetacoreClient
	Mu             *sync.Mutex

	chain        chains.Chain
	solanaClient *rpc.Client

	stop        chan struct{}
	logger      zerolog.Logger
	coreContext *clientcontext.ZetacoreContext
	chainParams observertypes.ChainParams
	programId   solana.PublicKey
	ts          *metrics.TelemetryServer

	lastTxSig solana.Signature
}

var _ interfaces.ChainObserver = &Observer{}

// NewObserver returns a new EVM chain observer
// TODO: read config for testnet and mainnet
func NewObserver(
	appContext *clientcontext.AppContext,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	programIdStr string,
	dbpath string,
	ts *metrics.TelemetryServer,
) (*Observer, error) {
	ob := Observer{
		ts: ts,
	}

	logger := log.With().Str("chain", "solana").Logger()
	ob.logger = logger

	ob.coreContext = appContext.ZetacoreContext()
	chainParams := observertypes.ChainParams{
		IsSupported: true,
	}
	ob.chainParams = chainParams
	ob.stop = make(chan struct{})
	ob.Mu = &sync.Mutex{}
	ob.zetacoreClient = zetacoreClient
	ob.Tss = tss
	ob.programId = solana.MustPublicKeyFromBase58(programIdStr)

	endpoint := "http://solana:8899"
	logger.Info().Msgf("Chain solana endpoint %s", endpoint)
	client := rpc.New(endpoint)
	if client == nil {
		logger.Error().Msg("solana Client new error")
		return nil, fmt.Errorf("solana Client new error")
	}

	ob.solanaClient = client
	{
		res1, err := client.GetVersion(context.TODO())
		if err != nil {
			logger.Error().Err(err).Msg("solana GetVersion error")
			return nil, err
		}
		logger.Info().Msgf("solana GetVersion %+v", res1)
		res2, err := client.GetHealth(context.TODO())
		if err != nil {
			logger.Error().Err(err).Msg("solana GetHealth error")
			return nil, err
		}
		logger.Info().Msgf("solana GetHealth %v", res2)

		logger.Info().Msgf("getting program info for %s", ob.programId.String())
		res3, err := client.GetAccountInfo(context.TODO(), ob.programId)
		if err != nil {
			logger.Error().Err(err).Msg("solana GetProgramAccounts error")
			return nil, err
		}
		//logger.Info().Msgf("solana GetProgramAccounts %v", res3)
		logger.Info().Msg(spew.Sprintf("%+v", res3))
	}
	return &ob, nil
}

func (o *Observer) IsOutboundProcessed(cctx *types.CrossChainTx, logger zerolog.Logger) (bool, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (o *Observer) SetChainParams(params observertypes.ChainParams) {
	//TODO implement me
	panic("implement me")
}

func (o *Observer) GetChainParams() observertypes.ChainParams {
	//TODO implement me
	return observertypes.ChainParams{
		IsSupported: true,
	}
}

func (o *Observer) GetTxID(nonce uint64) string {
	//TODO implement me
	panic("implement me")
}

func (o *Observer) WatchInboundTracker() {
	//TODO implement me
	panic("implement me")
}

func (o *Observer) Start() {
	o.logger.Info().Msgf("observer starting...")
	go o.WatchInbound()
}

func (o *Observer) Stop() {
	o.logger.Info().Msgf("observer stopping...")
}

func (o *Observer) WatchInbound() {
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("Solana_WatchInbound ticker"),
		10,
	)
	if err != nil {
		o.logger.Error().Err(err).Msg("error creating ticker")
		return
	}
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C():
			if !clientcontext.IsInboundObservationEnabled(o.coreContext, o.GetChainParams()) {
				o.logger.Info().
					Msgf("WatchInbound: inbound observation is disabled for chain solana")
				continue
			}
			err := o.ObserveInbound()
			if err != nil {
				o.logger.Err(err).Msg("WatchInbound: observeInbound error")
			}

		case <-o.stop:
			o.logger.Info().Msgf("WatchInbound stopped for chain %d", o.chain.ChainId)
			return
		}
	}
}

func (o *Observer) ObserveInbound() error {
	limit := 1000

	out, err := o.solanaClient.GetSignaturesForAddressWithOpts(
		context.TODO(),
		o.programId,
		&rpc.GetSignaturesForAddressOpts{
			Limit: &limit,
			//Before: solana.MustSignatureFromBase58("5pLBywq74Nc6jYrWUqn9KjnYXHbQEY2UPkhWefZF5u4NYaUvEwz1Cirqaym9wDeHNAjiQwuLBfrdhXo8uFQA45jL"),
			Until:      o.lastTxSig,
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		o.logger.Err(err).Msg("GetSignaturesForAddressWithOpts error")
		return err
	}
	o.logger.Info().Msgf("GetSignaturesForAddressWithOpts length %d", len(out))

	for i := len(out) - 1; i >= 0; i-- { // iterate txs from oldest to latest
		sig := out[i]
		o.logger.Info().Msgf("found sig: %s", sig.Signature)
		if sig.Err != nil { // ignore "failed" tx
			continue
		}
		tx, err := o.solanaClient.GetTransaction(context.TODO(), sig.Signature, &rpc.GetTransactionOpts{})
		if err != nil {
			o.logger.Err(err).Msg("GetTransaction error")
			return err // abort this observe operation in order to restart in next ticker trigger
		}
		o.lastTxSig = sig.Signature
		type DepositInstructionParams struct {
			Discriminator [8]byte
			Amount        uint64
			Memo          []byte
		}
		transaction, _ := tx.Transaction.GetTransaction()
		instruction := transaction.Message.Instructions[0] // TODO: parse not only the first instruction
		data := instruction.Data
		pk, _ := transaction.Message.Program(instruction.ProgramIDIndex)
		log.Info().Msgf("Program ID: %s", pk)
		var inst DepositInstructionParams
		err = borsh.Deserialize(&inst, data)
		if err != nil {
			log.Warn().Msgf("borsh.Deserialize error: %v", err)
			continue
		}
		// TODO: read discriminator from the IDL json file
		discriminator := []byte{242, 35, 198, 137, 82, 225, 242, 182}
		if !bytes.Equal(inst.Discriminator[:], discriminator) {
			continue
		}
		o.logger.Info().Msgf("  Amount Parameter: %d", inst.Amount)
		o.logger.Info().Msgf("  Memo (%d): %x", len(inst.Memo), inst.Memo)
		//var accounts []solana.PublicKey
		//for _, accIndex := range instruction.Accounts {
		//	accKey := transaction.Message.AccountKeys[accIndex]
		//	accounts = append(accounts, accKey)
		//}

	}
	return nil
}
