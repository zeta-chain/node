package zetacore

import (
	"context"
	"net"
	"testing"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	keyinterfaces "github.com/zeta-chain/zetacore/zetaclient/keys/interfaces"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"

	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
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

func TestZetacore_GetBallot(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryBallotByIdentifierResponse{
		BallotIdentifier: "123",
		Voters:           nil,
		ObservationType:  0,
		BallotStatus:     0,
	}
	input := observertypes.QueryBallotByIdentifierRequest{BallotIdentifier: "123"}
	method := "/zetachain.zetacore.observer.Query/BallotByIdentifier"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetBallotByID(ctx, "123")
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestZetacore_GetCrosschainFlags(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetCrosschainFlagsResponse{CrosschainFlags: observertypes.CrosschainFlags{
		IsInboundEnabled:      true,
		IsOutboundEnabled:     false,
		GasPriceIncreaseFlags: nil,
	}}
	input := observertypes.QueryGetCrosschainFlagsRequest{}
	method := "/zetachain.zetacore.observer.Query/CrosschainFlags"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetCrosschainFlags(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrosschainFlags, resp)
}

func TestZetacore_GetRateLimiterFlags(t *testing.T) {
	ctx := context.Background()

	// create sample flags
	rateLimiterFlags := sample.RateLimiterFlags()
	expectedOutput := crosschaintypes.QueryRateLimiterFlagsResponse{
		RateLimiterFlags: rateLimiterFlags,
	}

	// setup mock server
	input := crosschaintypes.QueryRateLimiterFlagsRequest{}
	method := "/zetachain.zetacore.crosschain.Query/RateLimiterFlags"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	// query
	resp, err := client.GetRateLimiterFlags(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.RateLimiterFlags, resp)
}

func TestZetacore_HeaderEnabledChains(t *testing.T) {
	ctx := context.Background()

	expectedOutput := lightclienttypes.QueryHeaderEnabledChainsResponse{
		HeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{
			{
				ChainId: chains.Ethereum.ChainId,
				Enabled: true,
			},
			{
				ChainId: chains.BitcoinMainnet.ChainId,
				Enabled: true,
			},
		},
	}
	input := lightclienttypes.QueryHeaderEnabledChainsRequest{}
	method := "/zetachain.zetacore.lightclient.Query/HeaderEnabledChains"
	setupMockServer(t, lightclienttypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetBlockHeaderEnabledChains(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.HeaderEnabledChains, resp)
}

func TestZetacore_GetChainParamsForChainID(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetChainParamsForChainResponse{ChainParams: &observertypes.ChainParams{
		ChainId:               123,
		BallotThreshold:       types.ZeroDec(),
		MinObserverDelegation: types.ZeroDec(),
	}}
	input := observertypes.QueryGetChainParamsForChainRequest{ChainId: 123}
	method := "/zetachain.zetacore.observer.Query/GetChainParamsForChain"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetChainParamsForChainID(ctx, 123)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainParams, resp)
}

func TestZetacore_GetChainParams(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetChainParamsResponse{ChainParams: &observertypes.ChainParamsList{
		ChainParams: []*observertypes.ChainParams{
			{
				ChainId:               123,
				MinObserverDelegation: types.ZeroDec(),
				BallotThreshold:       types.ZeroDec(),
			},
		},
	}}
	input := observertypes.QueryGetChainParamsRequest{}
	method := "/zetachain.zetacore.observer.Query/GetChainParams"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetChainParams(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainParams.ChainParams, resp)
}

func TestZetacore_GetUpgradePlan(t *testing.T) {
	ctx := context.Background()

	expectedOutput := upgradetypes.QueryCurrentPlanResponse{
		Plan: &upgradetypes.Plan{
			Name:   "big upgrade",
			Height: 100,
		},
	}
	input := upgradetypes.QueryCurrentPlanRequest{}
	method := "/cosmos.upgrade.v1beta1.Query/CurrentPlan"
	setupMockServer(t, upgradetypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetUpgradePlan(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Plan, resp)
}

func TestZetacore_GetAllCctx(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryAllCctxResponse{
		CrossChainTx: []*crosschaintypes.CrossChainTx{
			{
				Index: "cross-chain4456",
			},
		},
		Pagination: nil,
	}
	input := crosschaintypes.QueryAllCctxRequest{}
	method := "/zetachain.zetacore.crosschain.Query/CctxAll"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetAllCctx(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestZetacore_GetCctxByHash(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryGetCctxResponse{CrossChainTx: &crosschaintypes.CrossChainTx{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}}
	input := crosschaintypes.QueryGetCctxRequest{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}
	method := "/zetachain.zetacore.crosschain.Query/Cctx"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetCctxByHash(ctx, "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3")
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestZetacore_GetCctxByNonce(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryGetCctxResponse{CrossChainTx: &crosschaintypes.CrossChainTx{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}}
	input := crosschaintypes.QueryGetCctxByNonceRequest{
		ChainID: 7000,
		Nonce:   55,
	}
	method := "/zetachain.zetacore.crosschain.Query/CctxByNonce"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetCctxByNonce(ctx, 7000, 55)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestZetacore_GetObserverList(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryObserverSetResponse{
		Observers: []string{
			"zeta19jr7nl82lrktge35f52x9g5y5prmvchmk40zhg",
			"zeta1cxj07f3ju484ry2cnnhxl5tryyex7gev0yzxtj",
			"zeta1hjct6q7npsspsg3dgvzk3sdf89spmlpf7rqmnw",
		},
	}
	input := observertypes.QueryObserverSet{}
	method := "/zetachain.zetacore.observer.Query/ObserverSet"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetObserverList(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Observers, resp)
}

func TestZetacore_GetRateLimiterInput(t *testing.T) {
	ctx := context.Background()

	expectedOutput := &crosschaintypes.QueryRateLimiterInputResponse{
		Height:                  10,
		CctxsMissed:             []*crosschaintypes.CrossChainTx{sample.CrossChainTx(t, "1-1")},
		CctxsPending:            []*crosschaintypes.CrossChainTx{sample.CrossChainTx(t, "1-2")},
		TotalPending:            1,
		PastCctxsValue:          "123456",
		PendingCctxsValue:       "1234",
		LowestPendingCctxHeight: 2,
	}
	input := crosschaintypes.QueryRateLimiterInputRequest{Window: 10}
	method := "/zetachain.zetacore.crosschain.Query/RateLimiterInput"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetRateLimiterInput(ctx, 10)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, resp)
}

func TestZetacore_ListPendingCctx(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryListPendingCctxResponse{
		CrossChainTx: []*crosschaintypes.CrossChainTx{
			{
				Index: "cross-chain4456",
			},
		},
		TotalPending: 1,
	}
	input := crosschaintypes.QueryListPendingCctxRequest{ChainId: 7000}
	method := "/zetachain.zetacore.crosschain.Query/ListPendingCctx"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, totalPending, err := client.ListPendingCCTX(ctx, 7000)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
	require.Equal(t, expectedOutput.TotalPending, totalPending)
}

func TestZetacore_GetAbortedZetaAmount(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryZetaAccountingResponse{AbortedZetaAmount: "1080999"}
	input := crosschaintypes.QueryZetaAccountingRequest{}
	method := "/zetachain.zetacore.crosschain.Query/ZetaAccounting"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetAbortedZetaAmount(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.AbortedZetaAmount, resp)
}

// Need to test after refactor
func TestZetacore_GetGenesisSupply(t *testing.T) {
}

func TestZetacore_GetZetaTokenSupplyOnNode(t *testing.T) {
	ctx := context.Background()

	expectedOutput := banktypes.QuerySupplyOfResponse{
		Amount: types.Coin{
			Denom:  config.BaseDenom,
			Amount: types.NewInt(329438),
		}}
	input := banktypes.QuerySupplyOfRequest{Denom: config.BaseDenom}
	method := "/cosmos.bank.v1beta1.Query/SupplyOf"
	setupMockServer(t, banktypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetZetaTokenSupplyOnNode(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.GetAmount().Amount, resp)
}

func TestZetacore_GetBlockHeight(t *testing.T) {
	ctx := context.Background()

	method := "/zetachain.zetacore.crosschain.Query/LastZetaHeight"
	input := &crosschaintypes.QueryLastZetaHeightRequest{}
	output := &crosschaintypes.QueryLastZetaHeightResponse{Height: 12345}

	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, output)

	client := setupZetacoreClient(t,
		withDefaultObserverKeys(),
	)

	t.Run("last block height", func(t *testing.T) {
		height, err := client.GetBlockHeight(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(12345), height)
	})
}

func TestZetacore_GetLatestZetaBlock(t *testing.T) {
	ctx := context.Background()

	expectedOutput := tmservice.GetLatestBlockResponse{
		SdkBlock: &tmservice.Block{
			Header:     tmservice.Header{},
			Data:       tmtypes.Data{},
			Evidence:   tmtypes.EvidenceList{},
			LastCommit: nil,
		},
	}
	input := tmservice.GetLatestBlockRequest{}
	method := "/cosmos.base.tendermint.v1beta1.Service/GetLatestBlock"
	setupMockServer(t, tmservice.RegisterServiceServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetLatestZetaBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.SdkBlock, resp)
}

func TestZetacore_GetNodeInfo(t *testing.T) {
	ctx := context.Background()

	expectedOutput := tmservice.GetNodeInfoResponse{
		DefaultNodeInfo:    nil,
		ApplicationVersion: &tmservice.VersionInfo{},
	}
	input := tmservice.GetNodeInfoRequest{}
	method := "/cosmos.base.tendermint.v1beta1.Service/GetNodeInfo"
	setupMockServer(t, tmservice.RegisterServiceServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetNodeInfo(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestZetacore_GetBaseGasPrice(t *testing.T) {
	ctx := context.Background()

	expectedOutput := feemarkettypes.QueryParamsResponse{
		Params: feemarkettypes.Params{
			BaseFee: types.NewInt(23455),
		},
	}
	input := feemarkettypes.QueryParamsRequest{}
	method := "/ethermint.feemarket.v1.Query/Params"
	setupMockServer(t, feemarkettypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetBaseGasPrice(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Params.BaseFee.Int64(), resp)
}

func TestZetacore_GetNonceByChain(t *testing.T) {
	ctx := context.Background()

	chain := chains.BscMainnet
	expectedOutput := observertypes.QueryGetChainNoncesResponse{
		ChainNonces: observertypes.ChainNonces{
			Creator:         "",
			Index:           "",
			ChainId:         chain.ChainId,
			Nonce:           8446,
			Signers:         nil,
			FinalizedHeight: 0,
		},
	}
	input := observertypes.QueryGetChainNoncesRequest{Index: chain.ChainName.String()}
	method := "/zetachain.zetacore.observer.Query/ChainNonces"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetNonceByChain(ctx, chain)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainNonces, resp)
}

func TestZetacore_GetAllNodeAccounts(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryAllNodeAccountResponse{
		NodeAccount: []*observertypes.NodeAccount{
			{
				Operator:       "zeta19jr7nl82lrktge35f52x9g5y5prmvchmk40zhg",
				GranteeAddress: "zeta1kxhesgcvl6j5upupd9m3d3g3gfz4l3pcpqfnw6",
				GranteePubkey:  nil,
				NodeStatus:     0,
			},
		},
	}
	input := observertypes.QueryAllNodeAccountRequest{}
	method := "/zetachain.zetacore.observer.Query/NodeAccountAll"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetAllNodeAccounts(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.NodeAccount, resp)
}

func TestZetacore_GetKeyGen(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetKeygenResponse{
		Keygen: &observertypes.Keygen{
			Status:         observertypes.KeygenStatus_KeyGenSuccess,
			GranteePubkeys: nil,
			BlockNumber:    5646,
		}}
	input := observertypes.QueryGetKeygenRequest{}
	method := "/zetachain.zetacore.observer.Query/Keygen"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetKeyGen(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Keygen, resp)
}

func TestZetacore_GetBallotByID(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryBallotByIdentifierResponse{
		BallotIdentifier: "ballot1235",
	}
	input := observertypes.QueryBallotByIdentifierRequest{BallotIdentifier: "ballot1235"}
	method := "/zetachain.zetacore.observer.Query/BallotByIdentifier"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetBallot(ctx, "ballot1235")
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestZetacore_GetInboundTrackersForChain(t *testing.T) {
	ctx := context.Background()

	chainID := chains.BscMainnet.ChainId
	expectedOutput := crosschaintypes.QueryAllInboundTrackerByChainResponse{
		InboundTracker: []crosschaintypes.InboundTracker{
			{
				ChainId:  chainID,
				TxHash:   "DC76A6DCCC3AA62E89E69042ADC44557C50D59E4D3210C37D78DC8AE49B3B27F",
				CoinType: coin.CoinType_Gas,
			},
		},
	}
	input := crosschaintypes.QueryAllInboundTrackerByChainRequest{ChainId: chainID}
	method := "/zetachain.zetacore.crosschain.Query/InboundTrackerAllByChain"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetInboundTrackersForChain(ctx, chainID)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.InboundTracker, resp)
}

func TestZetacore_GetCurrentTss(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetTSSResponse{
		TSS: observertypes.TSS{
			TssPubkey:           "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TssParticipantList:  nil,
			OperatorAddressList: nil,
			FinalizedZetaHeight: 1000,
			KeyGenZetaHeight:    900,
		},
	}
	input := observertypes.QueryGetTSSRequest{}
	method := "/zetachain.zetacore.observer.Query/TSS"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetCurrentTSS(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.TSS, resp)
}

func TestZetacore_GetEthTssAddress(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetTssAddressResponse{
		Eth: "0x70e967acfcc17c3941e87562161406d41676fd83",
		Btc: "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
	}
	input := observertypes.QueryGetTssAddressRequest{}
	method := "/zetachain.zetacore.observer.Query/GetTssAddress"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetEVMTSSAddress(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Eth, resp)
}

func TestZetacore_GetBtcTssAddress(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetTssAddressResponse{
		Eth: "0x70e967acfcc17c3941e87562161406d41676fd83",
		Btc: "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
	}
	input := observertypes.QueryGetTssAddressRequest{BitcoinChainId: 8332}
	method := "/zetachain.zetacore.observer.Query/GetTssAddress"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetBTCTSSAddress(ctx, 8332)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Btc, resp)
}

func TestZetacore_GetTssHistory(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryTssHistoryResponse{
		TssList: []observertypes.TSS{
			{
				TssPubkey:           "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
				TssParticipantList:  nil,
				OperatorAddressList: nil,
				FinalizedZetaHeight: 46546,
				KeyGenZetaHeight:    6897,
			},
		},
	}
	input := observertypes.QueryTssHistoryRequest{}
	method := "/zetachain.zetacore.observer.Query/TssHistory"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetTSSHistory(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.TssList, resp)
}

func TestZetacore_GetOutboundTracker(t *testing.T) {
	chain := chains.BscMainnet
	expectedOutput := crosschaintypes.QueryGetOutboundTrackerResponse{
		OutboundTracker: crosschaintypes.OutboundTracker{
			Index:    "tracker12345",
			ChainId:  chain.ChainId,
			Nonce:    456,
			HashList: nil,
		},
	}
	input := crosschaintypes.QueryGetOutboundTrackerRequest{
		ChainID: chain.ChainId,
		Nonce:   456,
	}
	method := "/zetachain.zetacore.crosschain.Query/OutboundTracker"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	ctx := context.Background()
	resp, err := client.GetOutboundTracker(ctx, chain, 456)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.OutboundTracker, *resp)
}

func TestZetacore_GetAllOutboundTrackerByChain(t *testing.T) {
	ctx := context.Background()

	chain := chains.BscMainnet
	expectedOutput := crosschaintypes.QueryAllOutboundTrackerByChainResponse{
		OutboundTracker: []crosschaintypes.OutboundTracker{
			{
				Index:    "tracker23456",
				ChainId:  chain.ChainId,
				Nonce:    123456,
				HashList: nil,
			},
		},
	}
	input := crosschaintypes.QueryAllOutboundTrackerByChainRequest{
		Chain: chain.ChainId,
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      2000,
			CountTotal: false,
			Reverse:    false,
		},
	}
	method := "/zetachain.zetacore.crosschain.Query/OutboundTrackerAllByChain"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetAllOutboundTrackerByChain(ctx, chain.ChainId, interfaces.Ascending)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.OutboundTracker, resp)

	resp, err = client.GetAllOutboundTrackerByChain(ctx, chain.ChainId, interfaces.Descending)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.OutboundTracker, resp)
}

func TestZetacore_GetPendingNoncesByChain(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryPendingNoncesByChainResponse{
		PendingNonces: observertypes.PendingNonces{
			NonceLow:  0,
			NonceHigh: 0,
			ChainId:   chains.Ethereum.ChainId,
			Tss:       "",
		},
	}
	input := observertypes.QueryPendingNoncesByChainRequest{ChainId: chains.Ethereum.ChainId}
	method := "/zetachain.zetacore.observer.Query/PendingNoncesByChain"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetPendingNoncesByChain(ctx, chains.Ethereum.ChainId)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.PendingNonces, resp)
}

func TestZetacore_GetBlockHeaderChainState(t *testing.T) {
	ctx := context.Background()

	chainID := chains.BscMainnet.ChainId
	expectedOutput := lightclienttypes.QueryGetChainStateResponse{ChainState: &lightclienttypes.ChainState{
		ChainId:         chainID,
		LatestHeight:    5566654,
		EarliestHeight:  4454445,
		LatestBlockHash: nil,
	}}
	input := lightclienttypes.QueryGetChainStateRequest{ChainId: chainID}
	method := "/zetachain.zetacore.lightclient.Query/ChainState"
	setupMockServer(t, lightclienttypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetBlockHeaderChainState(ctx, chainID)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainState, resp)
}

func TestZetacore_GetSupportedChains(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QuerySupportedChainsResponse{
		Chains: []chains.Chain{
			{
				ChainName:   chains.BitcoinMainnet.ChainName,
				ChainId:     chains.BitcoinMainnet.ChainId,
				Network:     chains.BscMainnet.Network,
				NetworkType: chains.BscMainnet.NetworkType,
				Vm:          chains.BscMainnet.Vm,
				Consensus:   chains.BscMainnet.Consensus,
				IsExternal:  chains.BscMainnet.IsExternal,
			},
			{
				ChainName:   chains.Ethereum.ChainName,
				ChainId:     chains.Ethereum.ChainId,
				Network:     chains.Ethereum.Network,
				NetworkType: chains.Ethereum.NetworkType,
				Vm:          chains.Ethereum.Vm,
				Consensus:   chains.Ethereum.Consensus,
				IsExternal:  chains.Ethereum.IsExternal,
			},
		},
	}
	input := observertypes.QuerySupportedChains{}
	method := "/zetachain.zetacore.observer.Query/SupportedChains"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetSupportedChains(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Chains, resp)
}

func TestZetacore_GetAdditionalChains(t *testing.T) {
	ctx := context.Background()

	expectedOutput := authoritytypes.QueryGetChainInfoResponse{
		ChainInfo: authoritytypes.ChainInfo{
			Chains: []chains.Chain{
				chains.BitcoinMainnet,
				chains.Ethereum,
			},
		},
	}
	input := observertypes.QuerySupportedChains{}
	method := "/zetachain.zetacore.authority.Query/ChainInfo"

	setupMockServer(t, authoritytypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	resp, err := client.GetAdditionalChains(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainInfo.Chains, resp)
}

func TestZetacore_GetPendingNonces(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryAllPendingNoncesResponse{
		PendingNonces: []observertypes.PendingNonces{
			{
				NonceLow:  225,
				NonceHigh: 226,
				ChainId:   8332,
				Tss:       "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
			},
		},
	}
	input := observertypes.QueryAllPendingNoncesRequest{}
	method := "/zetachain.zetacore.observer.Query/PendingNoncesAll"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.GetPendingNonces(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestZetacore_Prove(t *testing.T) {
	ctx := context.Background()

	chainId := chains.BscMainnet.ChainId
	txHash := "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3"
	blockHash := "0000000000000000000172c9a64f86f208b867a84dc7a0b7c75be51e750ed8eb"
	txIndex := 555
	expectedOutput := lightclienttypes.QueryProveResponse{
		Valid: true,
	}
	input := lightclienttypes.QueryProveRequest{
		ChainId:   chainId,
		TxHash:    txHash,
		Proof:     nil,
		BlockHash: blockHash,
		TxIndex:   int64(txIndex),
	}
	method := "/zetachain.zetacore.lightclient.Query/Prove"
	setupMockServer(t, lightclienttypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.Prove(ctx, blockHash, txHash, int64(txIndex), nil, chainId)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Valid, resp)
}

func TestZetacore_HasVoted(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryHasVotedResponse{HasVoted: true}
	input := observertypes.QueryHasVotedRequest{
		BallotIdentifier: "123456asdf",
		VoterAddress:     "zeta1l40mm7meacx03r4lp87s9gkxfan32xnznp42u6",
	}
	method := "/zetachain.zetacore.observer.Query/HasVoted"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupZetacoreClient(t, withDefaultObserverKeys())

	resp, err := client.HasVoted(ctx, "123456asdf", "zeta1l40mm7meacx03r4lp87s9gkxfan32xnznp42u6")
	require.NoError(t, err)
	require.Equal(t, expectedOutput.HasVoted, resp)
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
