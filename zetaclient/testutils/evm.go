package testutils

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.non-eth.sol"
)

// ParseReceiptZetaSent parses a ZetaSent event from a receipt
func ParseReceiptZetaSent(
	receipt *ethtypes.Receipt,
	connector *zetaconnector.ZetaConnectorNonEth,
) *zetaconnector.ZetaConnectorNonEthZetaSent {
	for _, log := range receipt.Logs {
		event, err := connector.ParseZetaSent(*log)
		if err == nil && event != nil {
			return event // found
		}
	}
	return nil
}

// ParseReceiptERC20Deposited parses an Deposited event from a receipt
func ParseReceiptERC20Deposited(
	receipt *ethtypes.Receipt,
	custody *erc20custody.ERC20Custody,
) *erc20custody.ERC20CustodyDeposited {
	for _, log := range receipt.Logs {
		event, err := custody.ParseDeposited(*log)
		if err == nil && event != nil {
			return event // found
		}
	}
	return nil
}
