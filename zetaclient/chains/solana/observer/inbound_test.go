package observer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/solana/observer"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

var (
	// the relative path to the testdata directory
	TestDataDir = "../../../"
)

func Test_FilterInboundEventAndVote(t *testing.T) {
	// load archived inbound vote tx result
	// https://explorer.solana.com/tx/5LuQMorgd11p8GWEw6pmyHCDtA26NUyeNFhLWPNk2oBoM9pkag1LzhwGSRos3j4TJLhKjswFhZkGtvSGdLDkmqsk?cluster=devnet
	txHash := "5LuQMorgd11p8GWEw6pmyHCDtA26NUyeNFhLWPNk2oBoM9pkag1LzhwGSRos3j4TJLhKjswFhZkGtvSGdLDkmqsk"
	chain := chains.SolanaDevnet
	txResult := testutils.LoadSolanaInboundTxResult(t, TestDataDir, chain.ChainId, txHash, false)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = "2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s"
	zetacoreClient := mocks.NewZetacoreClient(t)
	zetacoreClient.WithKeys(&keys.Keys{}).WithZetaChain().WithPostVoteInbound("", "")

	dbpath := sample.CreateTempDir(t)
	ob, err := observer.NewObserver(chain, nil, *chainParams, zetacoreClient, nil, dbpath, base.DefaultLogger(), nil)
	require.NoError(t, err)

	t.Run("should filter inbound events and vote", func(t *testing.T) {
		err := ob.FilterInboundEventsAndVote(context.TODO(), txResult)
		require.NoError(t, err)
	})
}

func Test_FilterInboundEvents(t *testing.T) {
	// load archived inbound deposit tx result
	// https://explorer.solana.com/tx/5LuQMorgd11p8GWEw6pmyHCDtA26NUyeNFhLWPNk2oBoM9pkag1LzhwGSRos3j4TJLhKjswFhZkGtvSGdLDkmqsk?cluster=devnet
	txHash := "5LuQMorgd11p8GWEw6pmyHCDtA26NUyeNFhLWPNk2oBoM9pkag1LzhwGSRos3j4TJLhKjswFhZkGtvSGdLDkmqsk"
	chain := chains.SolanaDevnet
	txResult := testutils.LoadSolanaInboundTxResult(t, TestDataDir, chain.ChainId, txHash, false)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = "2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s"
	dbpath := sample.CreateTempDir(t)
	ob, err := observer.NewObserver(chain, nil, *chainParams, nil, nil, dbpath, base.DefaultLogger(), nil)
	require.NoError(t, err)

	// expected result
	sender := "AKbG83jg2V65R7XvaPFrnUvUTWsFENEzDPbLJFEiAk6L"
	eventExpected := &clienttypes.InboundEvent{
		SenderChainID: chain.ChainId,
		Sender:        sender,
		Receiver:      sender,
		TxOrigin:      sender,
		Amount:        1280,
		Memo:          []byte("hello this is a good memo for you to enjoy"),
		BlockNumber:   txResult.Slot,
		TxHash:        txHash,
		Index:         0, // not a EVM smart contract call
		CoinType:      coin.CoinType_Gas,
		Asset:         "", // no asset for gas token SOL
	}

	t.Run("should filter inbound event deposit SOL", func(t *testing.T) {
		events, err := ob.FilterInboundEvents(txResult)
		require.NoError(t, err)

		// check result
		require.Len(t, events, 1)
		require.EqualValues(t, eventExpected, events[0])
	})
}

func Test_BuildInboundVoteMsgFromEvent(t *testing.T) {
	// create test observer
	chain := chains.SolanaDevnet
	params := sample.ChainParams(chain.ChainId)
	params.GatewayAddress = sample.SolanaAddress(t)
	zetacoreClient := mocks.NewZetacoreClient(t)
	zetacoreClient.WithKeys(&keys.Keys{}).WithZetaChain().WithPostVoteInbound("", "")

	dbpath := sample.CreateTempDir(t)
	ob, err := observer.NewObserver(chain, nil, *params, zetacoreClient, nil, dbpath, base.DefaultLogger(), nil)
	require.NoError(t, err)

	// create test compliance config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for valid event", func(t *testing.T) {
		sender := sample.SolanaAddress(t)
		memo := sample.EthAddress().Bytes()
		event := sample.InboundEvent(chain.ChainId, sender, sender, 1280, []byte(memo))

		msg := ob.BuildInboundVoteMsgFromEvent(event)
		require.NotNil(t, msg)
	})
	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		sender := sample.SolanaAddress(t)
		receiver := sample.SolanaAddress(t)
		event := sample.InboundEvent(chain.ChainId, sender, receiver, 1280, nil)

		// restrict sender
		cfg.ComplianceConfig.RestrictedAddresses = []string{sender}
		config.LoadComplianceConfig(cfg)

		msg := ob.BuildInboundVoteMsgFromEvent(event)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		sender := sample.SolanaAddress(t)
		receiver := sample.SolanaAddress(t)
		memo := sample.EthAddress().Bytes()
		event := sample.InboundEvent(chain.ChainId, sender, receiver, 1280, []byte(memo))

		// restrict receiver
		cfg.ComplianceConfig.RestrictedAddresses = []string{receiver}
		config.LoadComplianceConfig(cfg)

		msg := ob.BuildInboundVoteMsgFromEvent(event)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg on donation transaction", func(t *testing.T) {
		// create event with donation memo
		sender := sample.SolanaAddress(t)
		event := sample.InboundEvent(chain.ChainId, sender, sender, 1280, []byte(constant.DonationMessage))

		msg := ob.BuildInboundVoteMsgFromEvent(event)
		require.Nil(t, msg)
	})
}

func Test_ParseInboundAsDeposit(t *testing.T) {
	// load archived inbound deposit tx result
	// https://explorer.solana.com/tx/5LuQMorgd11p8GWEw6pmyHCDtA26NUyeNFhLWPNk2oBoM9pkag1LzhwGSRos3j4TJLhKjswFhZkGtvSGdLDkmqsk?cluster=devnet
	txHash := "5LuQMorgd11p8GWEw6pmyHCDtA26NUyeNFhLWPNk2oBoM9pkag1LzhwGSRos3j4TJLhKjswFhZkGtvSGdLDkmqsk"
	chain := chains.SolanaDevnet

	txResult := testutils.LoadSolanaInboundTxResult(t, TestDataDir, chain.ChainId, txHash, false)
	tx, err := txResult.Transaction.GetTransaction()
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = "2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s"
	dbpath := sample.CreateTempDir(t)
	ob, err := observer.NewObserver(chain, nil, *chainParams, nil, nil, dbpath, base.DefaultLogger(), nil)
	require.NoError(t, err)

	// expected result
	sender := "AKbG83jg2V65R7XvaPFrnUvUTWsFENEzDPbLJFEiAk6L"
	eventExpected := &clienttypes.InboundEvent{
		SenderChainID: chain.ChainId,
		Sender:        sender,
		Receiver:      sender,
		TxOrigin:      sender,
		Amount:        1280,
		Memo:          []byte("hello this is a good memo for you to enjoy"),
		BlockNumber:   txResult.Slot,
		TxHash:        txHash,
		Index:         0, // not a EVM smart contract call
		CoinType:      coin.CoinType_Gas,
		Asset:         "", // no asset for gas token SOL
	}

	t.Run("should parse inbound event deposit SOL", func(t *testing.T) {
		event, err := ob.ParseInboundAsDeposit(tx, 0, txResult.Slot)
		require.NoError(t, err)

		// check result
		require.EqualValues(t, eventExpected, event)
	})
}
