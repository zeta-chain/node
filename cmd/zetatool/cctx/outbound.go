package cctx

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"

	zetatoolclients "github.com/zeta-chain/node/cmd/zetatool/clients"
	"github.com/zeta-chain/node/cmd/zetatool/context"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	solrepo "github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	zetaclientConfig "github.com/zeta-chain/node/zetaclient/config"
)

func (c *TrackingDetails) CheckOutbound(ctx *context.Context) error {
	outboundChain := c.OutboundChain

	switch {
	case outboundChain.IsEVMChain():
		return c.checkEvmOutboundTx(ctx)
	case outboundChain.IsBitcoinChain():
		return c.checkBitcoinOutboundTx(ctx)
	case outboundChain.IsSolanaChain():
		return c.checkSolanaOutboundTx(ctx)
	default:
		return fmt.Errorf("unsupported outbound chain")
	}
}

// checkEvmOutboundTx checks if the outbound transaction is confirmed on the outbound chain.
// If it's confirmed, we update the status to PendingOutboundVoting or PendingRevertVoting. Which means that the confirmation is done and we are not waiting for observers to vote
// Transition Status PendingConfirmation -> Status PendingVoting
func (c *TrackingDetails) checkEvmOutboundTx(ctx *context.Context) error {
	var (
		txHashList     = c.OutboundTrackerHashList
		outboundChain  = c.OutboundChain
		zetacoreReader = ctx.GetZetacoreReader()
		goCtx          = ctx.GetContext()
	)

	chainParams, err := zetacoreReader.GetChainParamsForChainID(goCtx, outboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}

	evmClient, err := zetatoolclients.NewEVMClientForChain(outboundChain, ctx.GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create evm client: %w", err)
	}

	foundConfirmedTx := false
	for _, hash := range txHashList {
		confirmed, err := zetatoolclients.IsTxConfirmed(goCtx, evmClient, hash, chainParams.OutboundConfirmationSafe())
		if err != nil {
			continue
		}
		if confirmed {
			foundConfirmedTx = true
			break
		}
	}
	if foundConfirmedTx {
		c.updateOutboundVoting()
	}
	return nil
}

func (c *TrackingDetails) checkSolanaOutboundTx(ctx *context.Context) error {
	var (
		txHashList = c.OutboundTrackerHashList
		goCtx      = ctx.GetContext()
		cfg        = ctx.GetConfig()
	)

	foundConfirmedTx := false
	solClient := solrpc.New(cfg.SolanaRPC)
	if solClient == nil {
		return fmt.Errorf("error creating rpc client")
	}
	solRepo := solrepo.New(solClient)

	for _, hash := range txHashList {
		signature := solana.MustSignatureFromBase58(hash)
		_, err := solRepo.GetTransaction(goCtx, signature)
		if err != nil {
			continue
		}
		foundConfirmedTx = true
	}

	if foundConfirmedTx {
		c.updateOutboundVoting()
	}
	return nil
}

func (c *TrackingDetails) checkBitcoinOutboundTx(ctx *context.Context) error {
	var (
		txHashList     = c.OutboundTrackerHashList
		outboundChain  = c.OutboundChain
		zetacoreReader = ctx.GetZetacoreReader()
		goCtx          = ctx.GetContext()
		cfg            = ctx.GetConfig()
		logger         = ctx.GetLogger()
	)

	chainParams, err := zetacoreReader.GetChainParamsForChainID(goCtx, outboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}
	confirmationCount := chainParams.OutboundConfirmationSafe()

	params, err := chains.BitcoinNetParamsFromChainID(outboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("unable to get bitcoin net params from chain id: %w", err)
	}

	connCfg := zetaclientConfig.BTCConfig{
		RPCUsername: cfg.BtcUser,
		RPCPassword: cfg.BtcPassword,
		RPCHost:     cfg.BtcHost,
		RPCParams:   params.Name,
	}

	btcClient, err := client.New(connCfg, outboundChain.ChainId, logger)
	if err != nil {
		return fmt.Errorf("unable to create rpc client: %w", err)
	}

	err = btcClient.Ping(goCtx)
	if err != nil {
		return fmt.Errorf("error ping the bitcoin server: %w", err)
	}

	foundConfirmedTx := false

	for _, hash := range txHashList {
		txHash, err := chainhash.NewHashFromStr(hash)
		if err != nil {
			continue
		}
		tx, err := btcClient.GetRawTransactionVerbose(goCtx, txHash)
		if err != nil {
			continue
		}

		if tx.Confirmations >= confirmationCount {
			foundConfirmedTx = true
		}
	}

	if foundConfirmedTx {
		c.updateOutboundVoting()
	}
	return nil
}
