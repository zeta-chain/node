package mocks

import (
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.non-eth.sol"

	"github.com/zeta-chain/node/pkg/constant"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

func MockChainParams(chainID int64, confirmation uint64) observertypes.ChainParams {
	connectorAddr := constant.EVMZeroAddress
	if a, ok := testutils.ConnectorAddresses[chainID]; ok {
		connectorAddr = a.Hex()
	}

	erc20CustodyAddr := constant.EVMZeroAddress
	if a, ok := testutils.CustodyAddresses[chainID]; ok {
		erc20CustodyAddr = a.Hex()
	}

	gwAddress := ""
	if gw, ok := testutils.GatewayAddresses[chainID]; ok {
		gwAddress = gw
	}

	return observertypes.ChainParams{
		ChainId:                     chainID,
		ZetaTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    connectorAddr,
		Erc20CustodyContractAddress: erc20CustodyAddr,
		InboundTicker:               12,
		OutboundTicker:              15,
		WatchUtxoTicker:             1,
		GasPriceTicker:              30,
		OutboundScheduleInterval:    30,
		OutboundScheduleLookahead:   60,
		BallotThreshold:             observertypes.DefaultBallotThreshold,
		MinObserverDelegation:       observertypes.DefaultMinObserverDelegation,
		GatewayAddress:              gwAddress,
		IsSupported:                 true,
		ConfirmationParams: &observertypes.ConfirmationParams{
			SafeInboundCount:  confirmation,
			SafeOutboundCount: confirmation,
		},
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
