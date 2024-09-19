package zetacore

import (
	"context"
	"net"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	feemarkettypes "github.com/zeta-chain/ethermint/x/feemarket/types"
	keyinterfaces "github.com/zeta-chain/node/zetaclient/keys/interfaces"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

const skipMethod = "skip"

// setupMockServer setup mock zetacore GRPC server
func setupMockServer(
	t *testing.T,
	serviceFunc any, method string, input any, expectedOutput any,
	extra ...grpcmock.ServerOption,
) *grpcmock.Server {
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)

	opts := []grpcmock.ServerOption{
		grpcmock.RegisterService(serviceFunc),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.WithListener(listener),
	}

	opts = append(opts, extra...)

	if method != skipMethod {
		opts = append(opts, func(s *grpcmock.Server) {
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(input).
				Return(expectedOutput)
		})
	}

	server := grpcmock.MockUnstartedServer(opts...)(t)

	server.Serve()

	t.Cleanup(func() {
		require.NoError(t, server.Close())
	})

	return server
}

func withDummyServer(zetaBlockHeight int64) []grpcmock.ServerOption {
	return []grpcmock.ServerOption{
		grpcmock.RegisterService(crosschaintypes.RegisterQueryServer),
		grpcmock.RegisterService(crosschaintypes.RegisterMsgServer),
		grpcmock.RegisterService(feemarkettypes.RegisterQueryServer),
		grpcmock.RegisterService(authtypes.RegisterQueryServer),
		grpcmock.RegisterService(abci.RegisterABCIApplicationServer),
		func(s *grpcmock.Server) {
			// Block Height
			s.ExpectUnary("/zetachain.zetacore.crosschain.Query/LastZetaHeight").
				UnlimitedTimes().
				Return(crosschaintypes.QueryLastZetaHeightResponse{Height: zetaBlockHeight})

			// London Base Fee
			s.ExpectUnary("/ethermint.feemarket.v1.Query/Params").
				UnlimitedTimes().
				Return(feemarkettypes.QueryParamsResponse{
					Params: feemarkettypes.Params{BaseFee: types.NewInt(100)},
				})
		},
	}
}

type clientTestConfig struct {
	keys keyinterfaces.ObserverKeys
	opts []Opt
}

type clientTestOpt func(*clientTestConfig)

func withObserverKeys(keys keyinterfaces.ObserverKeys) clientTestOpt {
	return func(cfg *clientTestConfig) { cfg.keys = keys }
}

func withDefaultObserverKeys() clientTestOpt {
	var (
		key     = mocks.TestKeyringPair
		address = types.AccAddress(key.PubKey().Address().Bytes())
		keyRing = mocks.NewKeyring()
	)

	return withObserverKeys(keys.NewKeysWithKeybase(keyRing, address, testSigner, ""))
}

func withTendermint(client cosmosclient.TendermintRPC) clientTestOpt {
	return func(cfg *clientTestConfig) { cfg.opts = append(cfg.opts, WithTendermintClient(client)) }
}

func withAccountRetriever(t *testing.T, accNum uint64, accSeq uint64) clientTestOpt {
	ctrl := gomock.NewController(t)
	ac := mock.NewMockAccountRetriever(ctrl)
	ac.EXPECT().
		GetAccountNumberSequence(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(accNum, accSeq, nil)

	return func(cfg *clientTestConfig) {
		cfg.opts = append(cfg.opts, WithCustomAccountRetriever(ac))
	}
}

func setupZetacoreClient(t *testing.T, opts ...clientTestOpt) *Client {
	const (
		chainIP = "127.0.0.1"
		signer  = testSigner
		chainID = "zetachain_7000-1"
	)

	var cfg clientTestConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.keys == nil {
		cfg.keys = &keys.Keys{}
	}

	c, err := NewClient(
		cfg.keys,
		chainIP, signer,
		chainID,
		false,
		zerolog.Nop(),
		cfg.opts...,
	)

	require.NoError(t, err)

	return c
}

// Need to test after refactor
func TestZetacore_GetGenesisSupply(t *testing.T) {
}

func TestZetacore_GetZetaHotKeyBalance(t *testing.T) {
	ctx := context.Background()

	expectedOutput := banktypes.QueryBalanceResponse{
		Balance: &types.Coin{
			Denom:  config.BaseDenom,
			Amount: types.NewInt(55646484),
		},
	}
	input := banktypes.QueryBalanceRequest{
		Address: types.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes()).String(),
		Denom:   config.BaseDenom,
	}
	method := "/cosmos.bank.v1beta1.Query/Balance"
	setupMockServer(t, banktypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	// should be able to get balance of signer
	client.keys = keys.NewKeysWithKeybase(mocks.NewKeyring(), types.AccAddress{}, "bob", "")
	resp, err := client.GetZetaHotKeyBalance(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Balance.Amount, resp)

	// should return error on empty signer
	client.keys = keys.NewKeysWithKeybase(mocks.NewKeyring(), types.AccAddress{}, "", "")
	resp, err = client.GetZetaHotKeyBalance(ctx)
	require.Error(t, err)
	require.Equal(t, types.ZeroInt(), resp)
}
