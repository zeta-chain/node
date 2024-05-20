package mocks

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func MockChainParams(chainID int64, confirmation uint64) observertypes.ChainParams {
	return observertypes.ChainParams{
		ChainId:                     chainID,
		ConfirmationCount:           confirmation,
		ConnectorContractAddress:    testutils.ConnectorAddresses[chainID].Hex(),
		Erc20CustodyContractAddress: testutils.CustodyAddresses[chainID].Hex(),
		IsSupported:                 true,
	}
}

func MockConnectorNonEth(chainID int64) *zetaconnector.ZetaConnectorNonEth {
	connector, err := zetaconnector.NewZetaConnectorNonEth(testutils.ConnectorAddresses[chainID], &ethclient.Client{})
	if err != nil {
		panic(err)
	}
	return connector
}

func MockERC20Custody(chainID int64) *erc20custody.ERC20Custody {
	custody, err := erc20custody.NewERC20Custody(testutils.CustodyAddresses[chainID], &ethclient.Client{})
	if err != nil {
		panic(err)
	}
	return custody
}
