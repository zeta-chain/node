package observer_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

var (
	// the relative path to the testdata directory
	TestDataDir = "../../../"
)

func Test_FilterInboundEventAndVote(t *testing.T) {
	// load archived inbound vote tx result from localnet
	txHash := "QSoSLxcJAFAzxWnHVJ4s2d5k2LyjC83YaLwbMUHYcEvVnCfERsowNb6Nj55GiTXNTbNF9fzF5F8JHUEpAGMrV5k"
	chain := chains.SolanaDevnet
	txResult := testutils.LoadSolanaInboundTxResult(t, TestDataDir, chain.ChainId, txHash, false)

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	// create observer
	chainParams := sample.ChainParams(chain.ChainId)
	chainParams.GatewayAddress = testutils.OldSolanaGatewayAddressDevnet
	zetacoreClient := mocks.NewZetacoreClient(t)
	zetacoreClient.
		WithKeys(&keys.Keys{OperatorAddress: []byte("something")}).
		WithZetaChain().
		WithPostVoteInbound("", "")

	zetacoreClient.MockGetCctxByHash("anything", nil)
	zetacoreClient.MockGetBallotByID(mock.Anything, nil)
	zetacoreClient.WithPostVoteInbound(sample.ZetaIndex(t), mock.Anything)

	baseObserver, err := base.NewObserver(
		chain,
		*chainParams,
		zrepo.New(zetacoreClient, chain, mode.StandardMode),
		nil,
		1000,
		nil,
		database,
		base.DefaultLogger(),
	)
	require.NoError(t, err)

	ob, err := observer.New(baseObserver, nil, chainParams.GatewayAddress)
	require.NoError(t, err)

	t.Run("should filter inbound events and vote", func(t *testing.T) {
		// expected result
		sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
		eventExpected := &clienttypes.InboundEvent{
			SenderChainID:    chain.ChainId,
			Sender:           sender,
			Receiver:         "0x103FD9224F00ce3013e95629e52DFc31D805D68d",
			TxOrigin:         sender,
			Amount:           24000000,
			Memo:             []byte{},
			BlockNumber:      txResult.Slot,
			TxHash:           txHash,
			Index:            0, // not a EVM smart contract call
			CoinType:         coin.CoinType_Gas,
			Asset:            "", // no asset for gas token SOL
			IsCrossChainCall: false,
		}

		events, err := observer.FilterInboundEvents(
			txResult,
			solana.MustPublicKeyFromBase58(testutils.OldSolanaGatewayAddressDevnet),
			chain.ChainId,
			zerolog.Nop(),
		)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.EqualValues(t, eventExpected, events[0])

		err = ob.VoteInboundEvents(context.TODO(), events, false, false)
		require.NoError(t, err)
	})
}

func Test_FilterInboundEvents(t *testing.T) {
	// ARRANGE
	// load archived inbound vote tx result from localnet
	txHash := "QSoSLxcJAFAzxWnHVJ4s2d5k2LyjC83YaLwbMUHYcEvVnCfERsowNb6Nj55GiTXNTbNF9fzF5F8JHUEpAGMrV5k"
	chain := chains.SolanaDevnet
	txResult := testutils.LoadSolanaInboundTxResult(t, TestDataDir, chain.ChainId, txHash, false)

	// load archived inbound vote tx result from localnet with inner deposit instruction
	txHashInner := "2TH3fMqFEULjavgmYEXtQrX6qSeMwKMPPgXEmr6QRFomgpVKhj8LNJNDwyqC1dVeSBU1Av2o4TWn75PvNfjncfUH"
	txResultInner := testutils.LoadSolanaInboundTxResult(t, TestDataDir, chain.ChainId, txHashInner, false)

	// load reverted inbound tx
	txHashReverted := "2M5hpf4CNdfV4a44Ra8mYAyQRZr3UX61FjAU6qBmV6K7HxyS6CsSPcvq2fS7eB9QqT8rx8jE2wMoMYauTmuvPgrx"
	txResultRevert := testutils.LoadSolanaInboundTxResult(t, TestDataDir, chain.ChainId, txHashReverted, false)

	// given gateway ID
	gatewayID, _, err := contracts.ParseGatewayWithPDA(testutils.OldSolanaGatewayAddressDevnet)
	require.NoError(t, err)

	t.Run("should filter inbound event deposit SOL", func(t *testing.T) {
		// expected result
		sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
		eventExpected := &clienttypes.InboundEvent{
			SenderChainID:    chain.ChainId,
			Sender:           sender,
			Receiver:         "0x103FD9224F00ce3013e95629e52DFc31D805D68d",
			TxOrigin:         sender,
			Amount:           24000000,
			Memo:             []byte{},
			BlockNumber:      txResult.Slot,
			TxHash:           txHash,
			Index:            0,
			CoinType:         coin.CoinType_Gas,
			Asset:            "", // no asset for gas token SOL
			IsCrossChainCall: false,
		}

		// ACT
		events, err := observer.FilterInboundEvents(txResult, gatewayID, chain.ChainId, zerolog.Nop())
		require.NoError(t, err)

		// ASSERT
		require.Len(t, events, 1)
		require.EqualValues(t, eventExpected, events[0])
	})

	t.Run("should filter inner inbound event deposit SOL", func(t *testing.T) {
		// expected result
		sender := "37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ"
		eventExpected := &clienttypes.InboundEvent{
			SenderChainID:    chain.ChainId,
			Sender:           sender,
			Receiver:         "0x103FD9224F00ce3013e95629e52DFc31D805D68d",
			TxOrigin:         sender,
			Amount:           24000000,
			Memo:             []byte{},
			BlockNumber:      txResultInner.Slot,
			TxHash:           txHashInner,
			Index:            0,
			CoinType:         coin.CoinType_Gas,
			Asset:            "", // no asset for gas token SOL
			IsCrossChainCall: false,
		}

		// ACT
		events, err := observer.FilterInboundEvents(txResultInner, gatewayID, chain.ChainId, zerolog.Nop())
		require.NoError(t, err)

		// ASSERT
		require.Len(t, events, 1)
		require.EqualValues(t, eventExpected, events[0])
	})

	t.Run("should not filter reverted inbound deposit SOL", func(t *testing.T) {
		_, err := observer.FilterInboundEvents(txResultRevert, gatewayID, chain.ChainId, zerolog.Nop())
		require.Error(t, err)
	})
}

func Test_BuildInboundVoteMsgFromEvent(t *testing.T) {
	// create test observer
	chain := chains.SolanaDevnet
	params := sample.ChainParams(chain.ChainId)
	params.GatewayAddress = sample.SolanaAddress(t)
	zetacoreClient := mocks.NewZetacoreClient(t)
	zetacoreClient.WithKeys(&keys.Keys{}).WithZetaChain().WithPostVoteInbound("", "")

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	baseObserver, err := base.NewObserver(
		chain,
		*params,
		zrepo.New(zetacoreClient, chain, mode.StandardMode),
		nil,
		1000,
		nil,
		database,
		base.DefaultLogger(),
	)
	require.NoError(t, err)

	ob, err := observer.New(baseObserver, nil, params.GatewayAddress)
	require.NoError(t, err)

	// create test compliance config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for valid event", func(t *testing.T) {
		sender := sample.SolanaAddress(t)
		receiver := sample.EthAddress()
		message := sample.Bytes()
		event := sample.InboundEvent(chain.ChainId, sender, receiver.Hex(), 1280, message)

		msg := ob.BuildInboundVoteMsgFromEvent(event)
		require.NotNil(t, msg)
		require.Equal(t, sender, msg.Sender)
		require.Equal(t, receiver.Hex(), msg.Receiver)
		require.Equal(t, hex.EncodeToString(message), msg.Message)
	})

	t.Run("should return nil if event is not processable", func(t *testing.T) {
		sender := sample.SolanaAddress(t)
		receiver := sample.SolanaAddress(t)
		event := sample.InboundEvent(chain.ChainId, sender, receiver, 1280, nil)

		// restrict sender
		cfg.ComplianceConfig.RestrictedAddresses = []string{sender}
		config.SetRestrictedAddressesFromConfig(cfg)

		msg := ob.BuildInboundVoteMsgFromEvent(event)
		require.Nil(t, msg)
	})
}

func Test_IsEventProcessable(t *testing.T) {
	// prepare params
	chain := chains.SolanaDevnet
	params := sample.ChainParams(chain.ChainId)
	params.GatewayAddress = sample.SolanaAddress(t)

	// create test observer
	ob := MockSolanaObserver(t, chain, nil, *params, nil, nil)

	// setup compliance config
	cfg := config.Config{
		ComplianceConfig: sample.ComplianceConfig(),
	}
	config.SetRestrictedAddressesFromConfig(cfg)

	// test cases
	tests := []struct {
		name   string
		event  clienttypes.InboundEvent
		result bool
	}{
		{
			name:   "should return true for processable event",
			event:  clienttypes.InboundEvent{Sender: sample.SolanaAddress(t), Receiver: sample.SolanaAddress(t)},
			result: true,
		},
		{
			name:   "should return false on donation message",
			event:  clienttypes.InboundEvent{Memo: []byte(constant.DonationMessage)},
			result: false,
		},
		{
			name: "should return false on compliance violation",
			event: clienttypes.InboundEvent{
				Sender:   sample.RestrictedSolAddressTest,
				Receiver: sample.EthAddress().Hex(),
			},
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ob.IsEventProcessable(tt.event)
			require.Equal(t, tt.result, result)
		})
	}
}
