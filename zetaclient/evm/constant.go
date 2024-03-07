package evm

import "time"

const (
	// ZetaBlockTime is the block time of the Zeta network
	ZetaBlockTime = 6500 * time.Millisecond

	// OutTxInclusionTimeout is the timeout for waiting for an outtx to be included in a block
	OutTxInclusionTimeout = 20 * time.Minute

	// OutTxTrackerReportTimeout is the timeout for waiting for an outtx tracker report
	OutTxTrackerReportTimeout = 10 * time.Minute

	// [signature, zetaTxSenderAddress, destinationChainId]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ZetaConnector.base.sol#L34
	TopicsZetaSent = 3

	// [signature, sourceChainId, destinationAddress, internalSendHash]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ZetaConnector.base.sol#L45
	TopicsZetaReceived = 4

	// [signature, destinationChainId, internalSendHash]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ZetaConnector.base.sol#L54
	TopicsZetaReverted = 3

	// [signature, recipient, asset]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ERC20Custody.sol#L43
	TopicsWithdrawn = 3

	// [signature, asset]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ERC20Custody.sol#L42
	TopicsDeposited = 2
)
