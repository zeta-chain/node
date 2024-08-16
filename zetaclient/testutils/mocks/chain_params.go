package mocks

import (
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"

	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func MockChainParams(chainID int64, confirmation uint64) observertypes.ChainParams {
	const zeroAddress = "0x0000000000000000000000000000000000000000"

	connectorAddr := zeroAddress
	if a, ok := testutils.ConnectorAddresses[chainID]; ok {
		connectorAddr = a.Hex()
	}

	erc20CustodyAddr := zeroAddress
	if a, ok := testutils.CustodyAddresses[chainID]; ok {
		erc20CustodyAddr = a.Hex()
	}

	return observertypes.ChainParams{
		ChainId:                     chainID,
		ConfirmationCount:           confirmation,
		ZetaTokenContractAddress:    zeroAddress,
		ConnectorContractAddress:    connectorAddr,
		Erc20CustodyContractAddress: erc20CustodyAddr,
		InboundTicker:               12,
		OutboundTicker:              15,
		WatchUtxoTicker:             0,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             observertypes.DefaultBallotThreshold,
		MinObserverDelegation:       observertypes.DefaultMinObserverDelegation,
		IsSupported:                 true,
	}
}

func MockConnectorNonEth(t *testing.T, chainID int64) *zetaconnector.ZetaConnectorNonEth {
	connector, err := zetaconnector.NewZetaConnectorNonEth(testutils.ConnectorAddresses[chainID], &ethclient.Client{})
	require.NoError(t, err)
	return connector
}

func MockERC20Custody(t *testing.T, chainID int64) *erc20custody.ERC20Custody {
	custody, err := erc20custody.NewERC20Custody(testutils.CustodyAddresses[chainID], &ethclient.Client{})
	require.NoError(t, err)
	return custody
}
