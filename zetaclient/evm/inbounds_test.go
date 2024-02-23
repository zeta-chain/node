package evm

import (
	"path"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func MockEVMClient(chain common.Chain) *ChainClient {
	return &ChainClient{
		chain:      chain,
		zetaClient: testutils.MockCoreBridge(),
	}
}

func MockConnectorNonEth() *zetaconnector.ZetaConnectorNonEth {
	connector, err := zetaconnector.NewZetaConnectorNonEth(ethcommon.Address{}, &ethclient.Client{})
	if err != nil {
		panic(err)
	}
	return connector
}

func ParseReceiptZetaSent(
	receipt *ethtypes.Receipt,
	ob *ChainClient,
	connector *zetaconnector.ZetaConnectorNonEth) *types.MsgVoteOnObservedInboundTx {
	var msg *types.MsgVoteOnObservedInboundTx
	for _, log := range receipt.Logs {
		event, err := connector.ParseZetaSent(*log)
		if err == nil && event != nil {
			msg = ob.GetInboundVoteMsgForZetaSentEvent(event)
			break // found
		}
	}
	return msg
}

func TestEthereum_GetInboundVoteMsgForZetaSentEvent(t *testing.T) {
	// load archived ZetaSent receipt
	// zeta-chain/crosschain/cctx/0x477544c4b8c8be544b23328b21286125c89cd6bb5d1d6d388d91eea8ea1a6f1f
	receipt := ethtypes.Receipt{}
	name := "chain_1_receipt_ZetaSent_0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76.json"
	err := testutils.LoadObjectFromJSONFile(&receipt, path.Join("../", testutils.TestDataPathEVM, name))
	require.NoError(t, err)

	// create mock client and connector
	ob := MockEVMClient(common.EthChain())
	connector := MockConnectorNonEth()

	// parse ZetaSent event
	msg := ParseReceiptZetaSent(&receipt, ob, connector)
	require.NotNil(t, msg)
	require.Equal(t, "0x477544c4b8c8be544b23328b21286125c89cd6bb5d1d6d388d91eea8ea1a6f1f", msg.Digest())

	// create config
	cfg := &config.Config{
		ComplianceConfig: &config.ComplianceConfig{},
	}

	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{msg.Sender}
		config.LoadComplianceConfig(cfg)
		msgRestricted := ParseReceiptZetaSent(&receipt, ob, connector)
		require.Nil(t, msgRestricted)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{msg.Receiver}
		config.LoadComplianceConfig(cfg)
		msgRestricted := ParseReceiptZetaSent(&receipt, ob, connector)
		require.Nil(t, msgRestricted)
	})
	t.Run("should return nil msg if txOrigin is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{msg.TxOrigin}
		config.LoadComplianceConfig(cfg)
		msgRestricted := ParseReceiptZetaSent(&receipt, ob, connector)
		require.Nil(t, msgRestricted)
	})
}
