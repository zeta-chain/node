package cctx

import (
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	zetatoolchains "github.com/zeta-chain/node/cmd/zetatool/chains"
	"github.com/zeta-chain/node/cmd/zetatool/context"
	zetaevmclient "github.com/zeta-chain/node/zetaclient/chains/evm/client"
)

func (c *CCTXDetails) CheckOutbound(ctx *context.Context) error {
	var (
		outboundChain = ctx.GetInboundChain()
		err           error
	)

	switch {
	case outboundChain.IsEVMChain():
		err = c.checkOutboundTx(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkOutboundTx checks if the outbound transaction is confirmed on the outbound chain.
// If it's confirmed, we update the status to PendingOutboundVoting or PendingRevertVoting. Which means that the confirmation is done and we are not waiting for observers to vote
// Transition Status PendingConfirmation -> Status PendingVoting
func (c *CCTXDetails) checkOutboundTx(ctx *context.Context) error {
	var (
		txHashList     = c.OutboundTrackerHashList
		outboundChain  = c.OutboundChain
		zetacoreClient = ctx.GetZetaCoreClient()
		goCtx          = ctx.GetContext()
	)

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, outboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %v", err)
	}

	// create evm client for the observation chain
	evmClient, err := zetatoolchains.GetEvmClient(ctx, outboundChain)
	if err != nil {
		return fmt.Errorf("failed to create evm client: %v", err)
	}

	foundConfirmedTx := false

	// If one of the hash is confirmed, we update the status to pending voting
	// There might be a condition where we have multiple txs and the wrong tx is confirmed.
	//To verify that we need, check CCTX data
	for _, hash := range txHashList {
		tx, _, err := zetatoolchains.GetEvmTx(ctx, evmClient, hash, outboundChain)
		if err != nil {
			continue
		}
		// Signer is unused
		c := zetaevmclient.New(evmClient, ethtypes.NewLondonSigner(tx.ChainId()))
		confirmed, err := c.IsTxConfirmed(goCtx, hash, chainParams.ConfirmationCount)
		if err != nil {
			continue
		}
		if confirmed {
			foundConfirmedTx = true
			break
		}
	}
	if foundConfirmedTx {
		c.UpdateOutboundVoting()
	}
	return nil
}
