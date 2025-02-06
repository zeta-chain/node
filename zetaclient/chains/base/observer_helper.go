package base

// CalcUnscannedBlockRangeInboundSafe calculates the unscanned block range using inbound safe confirmation count.
// It returns a range of blocks [from, end (exclusive)) that need to be scanned.
func (ob *Observer) CalcUnscannedBlockRangeInboundSafe(blockLimit uint64) (from uint64, end uint64) {
	lastBlock := ob.LastBlock()
	lastScanned := ob.LastBlockScanned()
	confirmation := ob.ChainParams().ConfirmationParams.SafeInboundCount

	return calcUnscannedBlockRange(lastBlock, lastScanned, confirmation, blockLimit)
}

// CalcUnscannedBlockRangeInboundFast calculates the unscanned block range using inbound fast confirmation count.
// It returns a range of blocks [from, end (exclusive)) that need to be scanned.
func (ob *Observer) CalcUnscannedBlockRangeInboundFast(blockLimit uint64) (from uint64, end uint64) {
	lastBlock := ob.LastBlock()
	lastScanned := ob.LastBlockScanned()
	confirmation := ob.ChainParams().InboundConfirmationFast()

	return calcUnscannedBlockRange(lastBlock, lastScanned, confirmation, blockLimit)
}

// IsBlockConfirmedForInboundSafe checks if the block number is confirmed using inbound safe confirmation count.
func (ob *Observer) IsBlockConfirmedForInboundSafe(blockNumber uint64) bool {
	lastBlock := ob.LastBlock()
	confirmation := ob.ChainParams().ConfirmationParams.SafeInboundCount
	return isBlockConfirmed(blockNumber, confirmation, lastBlock)
}

// IsBlockConfirmedForInboundFast checks if the block number is confirmed using inbound fast confirmation count.
// It falls back to safe confirmation count if fast confirmation is disabled.
func (ob *Observer) IsBlockConfirmedForInboundFast(blockNumber uint64) bool {
	lastBlock := ob.LastBlock()
	confirmation := ob.ChainParams().InboundConfirmationFast()
	return isBlockConfirmed(blockNumber, confirmation, lastBlock)
}

// IsBlockConfirmedForOutboundSafe checks if the block number is confirmed using outbound safe confirmation count.
func (ob *Observer) IsBlockConfirmedForOutboundSafe(blockNumber uint64) bool {
	lastBlock := ob.LastBlock()
	confirmation := ob.ChainParams().ConfirmationParams.SafeOutboundCount
	return isBlockConfirmed(blockNumber, confirmation, lastBlock)
}

// IsBlockConfirmedForOutboundFast checks if the block number is confirmed using outbound fast confirmation count.
// It falls back to safe confirmation count if fast confirmation is disabled.
func (ob *Observer) IsBlockConfirmedForOutboundFast(blockNumber uint64) bool {
	lastBlock := ob.LastBlock()
	confirmation := ob.ChainParams().OutboundConfirmationFast()
	return isBlockConfirmed(blockNumber, confirmation, lastBlock)
}

// calcUnscannedBlockRange calculates the unscanned block range [from, end (exclusive)) within given block limit.
//
// example 1: given lastBlock =  99, lastScanned = 90, confirmation = 10, then no unscanned block
// example 2: given lastBlock = 100, lastScanned = 90, confirmation = 10, then 1 unscanned block (block 91)
func calcUnscannedBlockRange(lastBlock, lastScanned, confirmation, blockLimit uint64) (from uint64, end uint64) {
	// got unscanned blocks or not?
	if lastBlock < lastScanned+confirmation {
		return 0, 0
	}

	// calculate the highest confirmed block
	// example: given lastBlock = 101, confirmation = 10, then the highest confirmed block is 92
	highestConfirmed := lastBlock - confirmation + 1

	// calculate a range of unscanned blocks within block limit
	from = lastScanned + 1
	end = from + blockLimit

	// 'end' is exclusive, so ensure it is not greater than (highestConfirmed+1)
	if end > highestConfirmed+1 {
		end = highestConfirmed + 1
	}

	return from, end
}

// isBlockConfirmed checks if the block number is confirmed.
//
// Note: block 100 is confirmed if the last block is 100 and confirmation count is 1.
func isBlockConfirmed(blockNumber uint64, confirmation uint64, lastBlock uint64) bool {
	confHeight := blockNumber + confirmation - 1
	return lastBlock >= confHeight
}
