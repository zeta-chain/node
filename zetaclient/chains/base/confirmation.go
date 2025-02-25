package base

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/constant"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// GetScanRangeInboundSafe calculates the block range to scan using inbound safe confirmation count.
// It returns a range of blocks [from, end (exclusive)) that need to be scanned.
func (ob *Observer) GetScanRangeInboundSafe(blockLimit uint64) (from uint64, end uint64) {
	lastBlock := ob.LastBlock()
	lastScanned := ob.LastBlockScanned()
	confirmation := ob.ChainParams().InboundConfirmationSafe()

	return calcUnscannedBlockRange(lastBlock, lastScanned, confirmation, blockLimit)
}

// GetScanRangeInboundFast calculates the block range to scan using inbound fast confirmation count.
// It returns a range of blocks [from, end (exclusive)) that need to be scanned.
func (ob *Observer) GetScanRangeInboundFast(blockLimit uint64) (from uint64, end uint64) {
	lastBlock := ob.LastBlock()
	lastScanned := ob.LastBlockScanned()
	confirmation := ob.ChainParams().InboundConfirmationFast()

	return calcUnscannedBlockRange(lastBlock, lastScanned, confirmation, blockLimit)
}

// IsBlockConfirmedForInboundSafe checks if the block number is confirmed using inbound safe confirmation count.
func (ob *Observer) IsBlockConfirmedForInboundSafe(blockNumber uint64) bool {
	lastBlock := ob.LastBlock()
	confirmation := ob.ChainParams().InboundConfirmationSafe()
	return isBlockConfirmed(blockNumber, confirmation, lastBlock)
}

// IsBlockConfirmedForInboundFast checks if the block number is confirmed using inbound fast confirmation count.
// It falls back to safe confirmation count if fast confirmation is disabled.
func (ob *Observer) IsBlockConfirmedForInboundFast(blockNumber uint64) bool {
	lastBlock := ob.LastBlock()
	confirmation := ob.ChainParams().InboundConfirmationFast()
	return isBlockConfirmed(blockNumber, confirmation, lastBlock)
}

// GetInboundConfirmationMode returns the confirmation mode for given inbound block number.
//   - It assumes the given block is at least FAST confirmed.
//   - It returns SAFE mode if the given block reaches SAFE confirmation count.
func (ob *Observer) GetInboundConfirmationMode(scannedBlock uint64) crosschaintypes.ConfirmationMode {
	if ob.IsBlockConfirmedForInboundSafe(scannedBlock) {
		return crosschaintypes.ConfirmationMode_SAFE
	}
	return crosschaintypes.ConfirmationMode_FAST
}

// IsBlockConfirmedForOutboundSafe checks if the block number is confirmed using outbound safe confirmation count.
func (ob *Observer) IsBlockConfirmedForOutboundSafe(blockNumber uint64) bool {
	lastBlock := ob.LastBlock()
	confirmation := ob.ChainParams().OutboundConfirmationSafe()
	return isBlockConfirmed(blockNumber, confirmation, lastBlock)
}

// IsBlockConfirmedForOutboundFast checks if the block number is confirmed using outbound fast confirmation count.
// It falls back to safe confirmation count if fast confirmation is disabled.
func (ob *Observer) IsBlockConfirmedForOutboundFast(blockNumber uint64) bool {
	lastBlock := ob.LastBlock()
	confirmation := ob.ChainParams().OutboundConfirmationFast()
	return isBlockConfirmed(blockNumber, confirmation, lastBlock)
}

// IsInboundEligibleForFastConfirmation determines if given inbound vote message is eligible for fast confirmation.
func (ob *Observer) IsInboundEligibleForFastConfirmation(
	ctx context.Context,
	msg *crosschaintypes.MsgVoteInbound,
) (bool, error) {
	// fast confirmation applies to chains that use confirmation count (e.g. EVM chains and Bitcoin)
	// these chains should have their liquidity cap multiplier set explicitly
	chainID := msg.SenderChainId
	multiplier, enabled := constant.GetInboundFastConfirmationLiquidityMultiplier(chainID)
	if !enabled {
		return false, nil
	}

	// check eligibility
	if !msg.EligibleForFastConfirmation() {
		return false, nil
	}

	// query liquidity cap for asset
	fCoins, err := ob.zetacoreClient.GetForeignCoinsFromAsset(ctx, chainID, ethcommon.HexToAddress(msg.Asset))
	if err != nil {
		return false, errors.Wrapf(err, "unable to get foreign coins for asset %s on chain %d", msg.Asset, chainID)
	}

	// ensure the deposit amount does not exceed amount cap
	fastAmountCap := constant.CalcInboundFastAmountCap(fCoins.LiquidityCap, multiplier)
	if msg.Amount.BigInt().Cmp(fastAmountCap) > 0 {
		return false, nil
	}

	return true, nil
}

// calcUnscannedBlockRange calculates the unscanned block range [from, end (exclusive)) within given block limit.
//
// example 1: given lastBlock =  99, lastScanned = 90, confirmation = 10, then no unscanned block
// example 2: given lastBlock = 100, lastScanned = 90, confirmation = 10, then 1 unscanned block (block 91)
func calcUnscannedBlockRange(lastBlock, lastScanned, confirmation, blockLimit uint64) (from uint64, end uint64) {
	// got unscanned blocks or not?
	// returning same values to indicate no unscanned block
	nextBlock := lastScanned + 1
	if lastBlock < lastScanned+confirmation {
		return nextBlock, nextBlock
	}

	// calculate the highest confirmed block
	// example: given lastBlock = 101, confirmation = 10, then the highest confirmed block is 92
	highestConfirmed := lastBlock - confirmation + 1

	// calculate a range of unscanned blocks within block limit
	// 'end' is exclusive, so ensure it is not greater than (highestConfirmed+1)
	from = nextBlock
	end = min(from+blockLimit, highestConfirmed+1)

	return from, end
}

// isBlockConfirmed checks if the block number is confirmed.
//
// Note: block 100 is confirmed if the last block is 100 and confirmation count is 1.
func isBlockConfirmed(blockNumber, confirmation, lastBlock uint64) bool {
	confHeight := blockNumber + confirmation - 1
	return lastBlock >= confHeight
}
