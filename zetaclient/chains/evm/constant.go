package evm

import "time"

const (
	// OutboundInclusionTimeout is the timeout for waiting for an outbound to be included in a block
	OutboundInclusionTimeout = 20 * time.Minute

	// ReorgProtectBlockCount is confirmations count to protect against reorg
	// Short 1~2 block reorgs could happen often on Ethereum due to network congestion or block production race conditions
	ReorgProtectBlockCount = 2

	// OutboundTrackerReportTimeout is the timeout for waiting for an outbound tracker report
	OutboundTrackerReportTimeout = 10 * time.Minute

	// EthTransferGasLimit is the gas limit for a standard ETH transfer
	EthTransferGasLimit = 21000

	// TopicsZetaSent is the number of topics for a Zeta sent event
	// [signature, zetaTxSenderAddress, destinationChainId]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ZetaConnector.base.sol#L34
	TopicsZetaSent = 3

	// TopicsZetaReceived is the number of topics for a Zeta received event
	// [signature, sourceChainId, destinationAddress, internalSendHash]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ZetaConnector.base.sol#L45
	TopicsZetaReceived = 4

	// TopicsZetaReverted is the number of topics for a Zeta reverted event
	// [signature, destinationChainId, internalSendHash]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ZetaConnector.base.sol#L54
	TopicsZetaReverted = 3

	// TopicsWithdrawn is the number of topics for a withdrawn event
	// [signature, recipient, asset]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ERC20Custody.sol#L43
	TopicsWithdrawn = 3

	// TopicsDeposited is the number of topics for a deposited event
	// [signature, asset]
	// https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ERC20Custody.sol#L42
	TopicsDeposited = 2
)
