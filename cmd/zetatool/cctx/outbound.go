package cctx

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gagliardetto/solana-go"

	zetatoolclients "github.com/zeta-chain/node/cmd/zetatool/clients"
	"github.com/zeta-chain/node/cmd/zetatool/context"
)

func (c *TrackingDetails) CheckOutbound(ctx *context.Context) error {
	var (
		outboundChain = ctx.GetInboundChain()
	)

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

func (c *TrackingDetails) checkEvmOutboundTx(ctx *context.Context) error {
	var (
		txHashList     = c.OutboundTrackerHashList
		outboundChain  = c.OutboundChain
		zetacoreClient = ctx.GetZetacoreClient()
		goCtx          = ctx.GetContext()
	)

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, outboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}

	evmClient, err := zetatoolclients.NewEVMClientAdapter(outboundChain, ctx.GetConfig())
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
	solClient, err := zetatoolclients.NewSolanaClientAdapter(cfg.SolanaRPC)
	if err != nil {
		return fmt.Errorf("error creating rpc client: %w", err)
	}

	for _, hash := range txHashList {
		signature := solana.MustSignatureFromBase58(hash)
		_, err := solClient.GetTransaction(goCtx, signature)
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
		zetacoreClient = ctx.GetZetacoreClient()
		goCtx          = ctx.GetContext()
		cfg            = ctx.GetConfig()
		logger         = ctx.GetLogger()
	)

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, outboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}
	confirmationCount := chainParams.OutboundConfirmationSafe()

	btcClient, err := zetatoolclients.NewBitcoinClientAdapter(cfg, outboundChain, logger)
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
