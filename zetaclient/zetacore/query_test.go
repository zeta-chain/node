package zetacore

import (
	"net"
	"testing"

	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschainTypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"
)

func setupMockServer(t *testing.T, serviceFunc any, method string, input any, expectedOutput any) *grpcmock.Server {
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)

	server := grpcmock.MockUnstartedServer(
		grpcmock.RegisterService(serviceFunc),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.WithListener(listener),
		func(s *grpcmock.Server) {
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(input).
				Return(expectedOutput)
		},
	)(t)

	return server
}

func closeMockServer(t *testing.T, server *grpcmock.Server) {
	err := server.Close()
	require.NoError(t, err)
}

func setupZetacoreClient() (*Client, error) {
	return NewClient(
		&keys.Keys{},
		"127.0.0.1",
		"",
		"zetachain_7000-1",
		false,
		&metrics.TelemetryServer{})
}

func TestZetacore_GetBallot(t *testing.T) {
	expectedOutput := observertypes.QueryBallotByIdentifierResponse{
		BallotIdentifier: "123",
		Voters:           nil,
		ObservationType:  0,
		BallotStatus:     0,
	}
	input := observertypes.QueryBallotByIdentifierRequest{BallotIdentifier: "123"}
	method := "/zetachain.zetacore.observer.Query/BallotByIdentifier"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetBallotByID("123")
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestZetacore_GetCrosschainFlags(t *testing.T) {
	expectedOutput := observertypes.QueryGetCrosschainFlagsResponse{CrosschainFlags: observertypes.CrosschainFlags{
		IsInboundEnabled:             true,
		IsOutboundEnabled:            false,
		GasPriceIncreaseFlags:        nil,
		BlockHeaderVerificationFlags: nil,
	}}
	input := observertypes.QueryGetCrosschainFlagsRequest{}
	method := "/zetachain.zetacore.observer.Query/CrosschainFlags"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetCrosschainFlags()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrosschainFlags, resp)
}

func TestZetacore_GetRateLimiterFlags(t *testing.T) {
	// create sample flags
	rateLimiterFlags := sample.RateLimiterFlags()
	expectedOutput := crosschainTypes.QueryRateLimiterFlagsResponse{
		RateLimiterFlags: rateLimiterFlags,
	}

	// setup mock server
	input := crosschainTypes.QueryRateLimiterFlagsRequest{}
	method := "/zetachain.zetacore.crosschain.Query/RateLimiterFlags"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	// query
	resp, err := client.GetRateLimiterFlags()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.RateLimiterFlags, resp)
}

func TestZetacore_HeaderEnabledChains(t *testing.T) {
	expectedOutput := lightclienttypes.QueryHeaderEnabledChainsResponse{HeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{
		{
			ChainId: chains.EthChain.ChainId,
			Enabled: true,
		},
		{
			ChainId: chains.BtcMainnetChain.ChainId,
			Enabled: true,
		},
	}}
	input := lightclienttypes.QueryHeaderEnabledChainsRequest{}
	method := "/zetachain.zetacore.lightclient.Query/HeaderEnabledChains"
	server := setupMockServer(t, lightclienttypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetBlockHeaderEnabledChains()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.HeaderEnabledChains, resp)
}

func TestZetacore_GetChainParamsForChainID(t *testing.T) {
	expectedOutput := observertypes.QueryGetChainParamsForChainResponse{ChainParams: &observertypes.ChainParams{
		ChainId:               123,
		BallotThreshold:       types.ZeroDec(),
		MinObserverDelegation: types.ZeroDec(),
	}}
	input := observertypes.QueryGetChainParamsForChainRequest{ChainId: 123}
	method := "/zetachain.zetacore.observer.Query/GetChainParamsForChain"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetChainParamsForChainID(123)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainParams, resp)
}

func TestZetacore_GetChainParams(t *testing.T) {
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
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetChainParams()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainParams.ChainParams, resp)
}

func TestZetacore_GetUpgradePlan(t *testing.T) {
	expectedOutput := upgradetypes.QueryCurrentPlanResponse{
		Plan: &upgradetypes.Plan{
			Name:   "big upgrade",
			Height: 100,
		},
	}
	input := upgradetypes.QueryCurrentPlanRequest{}
	method := "/cosmos.upgrade.v1beta1.Query/CurrentPlan"
	server := setupMockServer(t, upgradetypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetUpgradePlan()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Plan, resp)
}

func TestZetacore_GetAllCctx(t *testing.T) {
	expectedOutput := crosschainTypes.QueryAllCctxResponse{
		CrossChainTx: []*crosschainTypes.CrossChainTx{
			{
				Index: "cross-chain4456",
			},
		},
		Pagination: nil,
	}
	input := crosschainTypes.QueryAllCctxRequest{}
	method := "/zetachain.zetacore.crosschain.Query/CctxAll"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetAllCctx()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestZetacore_GetCctxByHash(t *testing.T) {
	expectedOutput := crosschainTypes.QueryGetCctxResponse{CrossChainTx: &crosschainTypes.CrossChainTx{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}}
	input := crosschainTypes.QueryGetCctxRequest{Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3"}
	method := "/zetachain.zetacore.crosschain.Query/Cctx"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetCctxByHash("9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3")
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestZetacore_GetCctxByNonce(t *testing.T) {
	expectedOutput := crosschainTypes.QueryGetCctxResponse{CrossChainTx: &crosschainTypes.CrossChainTx{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}}
	input := crosschainTypes.QueryGetCctxByNonceRequest{
		ChainID: 7000,
		Nonce:   55,
	}
	method := "/zetachain.zetacore.crosschain.Query/CctxByNonce"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetCctxByNonce(7000, 55)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestZetacore_GetObserverList(t *testing.T) {
	expectedOutput := observertypes.QueryObserverSetResponse{
		Observers: []string{
			"zeta19jr7nl82lrktge35f52x9g5y5prmvchmk40zhg",
			"zeta1cxj07f3ju484ry2cnnhxl5tryyex7gev0yzxtj",
			"zeta1hjct6q7npsspsg3dgvzk3sdf89spmlpf7rqmnw",
		},
	}
	input := observertypes.QueryObserverSet{}
	method := "/zetachain.zetacore.observer.Query/ObserverSet"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetObserverList()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Observers, resp)
}

func TestZetacore_GetRateLimiterInput(t *testing.T) {
	expectedOutput := crosschainTypes.QueryRateLimiterInputResponse{
		Height:                  10,
		CctxsMissed:             []*crosschainTypes.CrossChainTx{sample.CrossChainTx(t, "1-1")},
		CctxsPending:            []*crosschainTypes.CrossChainTx{sample.CrossChainTx(t, "1-2")},
		TotalPending:            1,
		PastCctxsValue:          "123456",
		PendingCctxsValue:       "1234",
		LowestPendingCctxHeight: 2,
	}
	input := crosschainTypes.QueryRateLimiterInputRequest{Window: 10}
	method := "/zetachain.zetacore.crosschain.Query/RateLimiterInput"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetRateLimiterInput(10)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, resp)
}

func TestZetacore_ListPendingCctx(t *testing.T) {
	expectedOutput := crosschainTypes.QueryListPendingCctxResponse{
		CrossChainTx: []*crosschainTypes.CrossChainTx{
			{
				Index: "cross-chain4456",
			},
		},
		TotalPending: 1,
	}
	input := crosschainTypes.QueryListPendingCctxRequest{ChainId: 7000}
	method := "/zetachain.zetacore.crosschain.Query/ListPendingCctx"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, totalPending, err := client.ListPendingCctx(7000)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
	require.Equal(t, expectedOutput.TotalPending, totalPending)
}

func TestZetacore_GetAbortedZetaAmount(t *testing.T) {
	expectedOutput := crosschainTypes.QueryZetaAccountingResponse{AbortedZetaAmount: "1080999"}
	input := crosschainTypes.QueryZetaAccountingRequest{}
	method := "/zetachain.zetacore.crosschain.Query/ZetaAccounting"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetAbortedZetaAmount()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.AbortedZetaAmount, resp)
}

// Need to test after refactor
func TestZetacore_GetGenesisSupply(t *testing.T) {
}

func TestZetacore_GetZetaTokenSupplyOnNode(t *testing.T) {
	expectedOutput := banktypes.QuerySupplyOfResponse{
		Amount: types.Coin{
			Denom:  config.BaseDenom,
			Amount: types.NewInt(329438),
		}}
	input := banktypes.QuerySupplyOfRequest{Denom: config.BaseDenom}
	method := "/cosmos.bank.v1beta1.Query/SupplyOf"
	server := setupMockServer(t, banktypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetZetaTokenSupplyOnNode()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.GetAmount().Amount, resp)
}

func TestZetacore_GetLastBlockHeight(t *testing.T) {
	expectedOutput := crosschainTypes.QueryAllLastBlockHeightResponse{
		LastBlockHeight: []*crosschainTypes.LastBlockHeight{
			{
				Index:             "test12345",
				Chain:             "7000",
				LastSendHeight:    32345,
				LastReceiveHeight: 23623,
			},
		},
	}
	input := crosschainTypes.QueryAllLastBlockHeightRequest{}
	method := "/zetachain.zetacore.crosschain.Query/LastBlockHeightAll"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	t.Run("last block height", func(t *testing.T) {
		resp, err := client.GetLastBlockHeight()
		require.NoError(t, err)
		require.Equal(t, expectedOutput.LastBlockHeight, resp)
	})
}

func TestZetacore_GetLatestZetaBlock(t *testing.T) {
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
	server := setupMockServer(t, tmservice.RegisterServiceServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetLatestZetaBlock()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.SdkBlock, resp)
}

func TestZetacore_GetNodeInfo(t *testing.T) {
	expectedOutput := tmservice.GetNodeInfoResponse{
		DefaultNodeInfo:    nil,
		ApplicationVersion: &tmservice.VersionInfo{},
	}
	input := tmservice.GetNodeInfoRequest{}
	method := "/cosmos.base.tendermint.v1beta1.Service/GetNodeInfo"
	server := setupMockServer(t, tmservice.RegisterServiceServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetNodeInfo()
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestZetacore_GetLastBlockHeightByChain(t *testing.T) {
	index := chains.BscMainnetChain
	expectedOutput := crosschainTypes.QueryGetLastBlockHeightResponse{
		LastBlockHeight: &crosschainTypes.LastBlockHeight{
			Index:             index.ChainName.String(),
			Chain:             "7000",
			LastSendHeight:    2134123,
			LastReceiveHeight: 1234333,
		},
	}
	input := crosschainTypes.QueryGetLastBlockHeightRequest{Index: index.ChainName.String()}
	method := "/zetachain.zetacore.crosschain.Query/LastBlockHeight"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetLastBlockHeightByChain(index)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.LastBlockHeight, resp)
}

func TestZetacore_GetZetaBlockHeight(t *testing.T) {
	expectedOutput := crosschainTypes.QueryLastZetaHeightResponse{Height: 12345}
	input := crosschainTypes.QueryLastZetaHeightRequest{}
	method := "/zetachain.zetacore.crosschain.Query/LastZetaHeight"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	t.Run("get zeta block height success", func(t *testing.T) {
		resp, err := client.GetBlockHeight()
		require.NoError(t, err)
		require.Equal(t, expectedOutput.Height, resp)
	})
}

func TestZetacore_GetBaseGasPrice(t *testing.T) {
	expectedOutput := feemarkettypes.QueryParamsResponse{
		Params: feemarkettypes.Params{
			BaseFee: types.NewInt(23455),
		},
	}
	input := feemarkettypes.QueryParamsRequest{}
	method := "/ethermint.feemarket.v1.Query/Params"
	server := setupMockServer(t, feemarkettypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetBaseGasPrice()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Params.BaseFee.Int64(), resp)
}

func TestZetacore_GetNonceByChain(t *testing.T) {
	chain := chains.BscMainnetChain
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
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetNonceByChain(chain)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainNonces, resp)
}

func TestZetacore_GetAllNodeAccounts(t *testing.T) {
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
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetAllNodeAccounts()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.NodeAccount, resp)
}

func TestZetacore_GetKeyGen(t *testing.T) {
	expectedOutput := observertypes.QueryGetKeygenResponse{
		Keygen: &observertypes.Keygen{
			Status:         observertypes.KeygenStatus_KeyGenSuccess,
			GranteePubkeys: nil,
			BlockNumber:    5646,
		}}
	input := observertypes.QueryGetKeygenRequest{}
	method := "/zetachain.zetacore.observer.Query/Keygen"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetKeyGen()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Keygen, resp)
}

func TestZetacore_GetBallotByID(t *testing.T) {
	expectedOutput := observertypes.QueryBallotByIdentifierResponse{
		BallotIdentifier: "ballot1235",
	}
	input := observertypes.QueryBallotByIdentifierRequest{BallotIdentifier: "ballot1235"}
	method := "/zetachain.zetacore.observer.Query/BallotByIdentifier"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetBallot("ballot1235")
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestZetacore_GetInboundTrackersForChain(t *testing.T) {
	chainID := chains.BscMainnetChain.ChainId
	expectedOutput := crosschainTypes.QueryAllInTxTrackerByChainResponse{
		InTxTracker: []crosschainTypes.InTxTracker{
			{
				ChainId:  chainID,
				TxHash:   "DC76A6DCCC3AA62E89E69042ADC44557C50D59E4D3210C37D78DC8AE49B3B27F",
				CoinType: coin.CoinType_Gas,
			},
		},
	}
	input := crosschainTypes.QueryAllInTxTrackerByChainRequest{ChainId: chainID}
	method := "/zetachain.zetacore.crosschain.Query/InTxTrackerAllByChain"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetInboundTrackersForChain(chainID)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.InTxTracker, resp)
}

func TestZetacore_GetCurrentTss(t *testing.T) {
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
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetCurrentTss()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.TSS, resp)
}

func TestZetacore_GetEthTssAddress(t *testing.T) {
	expectedOutput := observertypes.QueryGetTssAddressResponse{
		Eth: "0x70e967acfcc17c3941e87562161406d41676fd83",
		Btc: "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
	}
	input := observertypes.QueryGetTssAddressRequest{}
	method := "/zetachain.zetacore.observer.Query/GetTssAddress"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetEthTssAddress()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Eth, resp)
}

func TestZetacore_GetBtcTssAddress(t *testing.T) {
	expectedOutput := observertypes.QueryGetTssAddressResponse{
		Eth: "0x70e967acfcc17c3941e87562161406d41676fd83",
		Btc: "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
	}
	input := observertypes.QueryGetTssAddressRequest{BitcoinChainId: 8332}
	method := "/zetachain.zetacore.observer.Query/GetTssAddress"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetBtcTssAddress(8332)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Btc, resp)
}

func TestZetacore_GetTssHistory(t *testing.T) {
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
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetTssHistory()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.TssList, resp)
}

func TestZetacore_GetOutTxTracker(t *testing.T) {
	chain := chains.BscMainnetChain
	expectedOutput := crosschainTypes.QueryGetOutTxTrackerResponse{
		OutTxTracker: crosschainTypes.OutTxTracker{
			Index:    "tracker12345",
			ChainId:  chain.ChainId,
			Nonce:    456,
			HashList: nil,
		},
	}
	input := crosschainTypes.QueryGetOutTxTrackerRequest{
		ChainID: chain.ChainId,
		Nonce:   456,
	}
	method := "/zetachain.zetacore.crosschain.Query/OutTxTracker"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetOutTxTracker(chain, 456)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.OutTxTracker, *resp)
}

func TestZetacore_GetAllOutTxTrackerByChain(t *testing.T) {
	chain := chains.BscMainnetChain
	expectedOutput := crosschainTypes.QueryAllOutTxTrackerByChainResponse{
		OutTxTracker: []crosschainTypes.OutTxTracker{
			{
				Index:    "tracker23456",
				ChainId:  chain.ChainId,
				Nonce:    123456,
				HashList: nil,
			},
		},
	}
	input := crosschainTypes.QueryAllOutTxTrackerByChainRequest{
		Chain: chain.ChainId,
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      2000,
			CountTotal: false,
			Reverse:    false,
		},
	}
	method := "/zetachain.zetacore.crosschain.Query/OutTxTrackerAllByChain"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetAllOutTxTrackerByChain(chain.ChainId, interfaces.Ascending)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.OutTxTracker, resp)

	resp, err = client.GetAllOutTxTrackerByChain(chain.ChainId, interfaces.Descending)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.OutTxTracker, resp)
}

func TestZetacore_GetPendingNoncesByChain(t *testing.T) {
	expectedOutput := observertypes.QueryPendingNoncesByChainResponse{
		PendingNonces: observertypes.PendingNonces{
			NonceLow:  0,
			NonceHigh: 0,
			ChainId:   chains.EthChain.ChainId,
			Tss:       "",
		},
	}
	input := observertypes.QueryPendingNoncesByChainRequest{ChainId: chains.EthChain.ChainId}
	method := "/zetachain.zetacore.observer.Query/PendingNoncesByChain"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetPendingNoncesByChain(chains.EthChain.ChainId)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.PendingNonces, resp)
}

func TestZetacore_GetBlockHeaderChainState(t *testing.T) {
	chainID := chains.BscMainnetChain.ChainId
	expectedOutput := lightclienttypes.QueryGetChainStateResponse{ChainState: &lightclienttypes.ChainState{
		ChainId:         chainID,
		LatestHeight:    5566654,
		EarliestHeight:  4454445,
		LatestBlockHash: nil,
	}}
	input := lightclienttypes.QueryGetChainStateRequest{ChainId: chainID}
	method := "/zetachain.zetacore.lightclient.Query/ChainState"
	server := setupMockServer(t, lightclienttypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetBlockHeaderChainState(chainID)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, resp)
}

func TestZetacore_GetSupportedChains(t *testing.T) {
	expectedOutput := observertypes.QuerySupportedChainsResponse{
		Chains: []*chains.Chain{
			{
				ChainName:   chains.BtcMainnetChain.ChainName,
				ChainId:     chains.BtcMainnetChain.ChainId,
				Network:     chains.BscMainnetChain.Network,
				NetworkType: chains.BscMainnetChain.NetworkType,
				Vm:          chains.BscMainnetChain.Vm,
				Consensus:   chains.BscMainnetChain.Consensus,
				IsExternal:  chains.BscMainnetChain.IsExternal,
			},
			{
				ChainName:   chains.EthChain.ChainName,
				ChainId:     chains.EthChain.ChainId,
				Network:     chains.EthChain.Network,
				NetworkType: chains.EthChain.NetworkType,
				Vm:          chains.EthChain.Vm,
				Consensus:   chains.EthChain.Consensus,
				IsExternal:  chains.EthChain.IsExternal,
			},
		},
	}
	input := observertypes.QuerySupportedChains{}
	method := "/zetachain.zetacore.observer.Query/SupportedChains"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetSupportedChains()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Chains, resp)
}

func TestZetacore_GetPendingNonces(t *testing.T) {
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
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.GetPendingNonces()
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestZetacore_Prove(t *testing.T) {
	chainId := chains.BscMainnetChain.ChainId
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
	server := setupMockServer(t, lightclienttypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.Prove(blockHash, txHash, int64(txIndex), nil, chainId)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Valid, resp)
}

func TestZetacore_HasVoted(t *testing.T) {
	expectedOutput := observertypes.QueryHasVotedResponse{HasVoted: true}
	input := observertypes.QueryHasVotedRequest{
		BallotIdentifier: "123456asdf",
		VoterAddress:     "zeta1l40mm7meacx03r4lp87s9gkxfan32xnznp42u6",
	}
	method := "/zetachain.zetacore.observer.Query/HasVoted"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)

	resp, err := client.HasVoted("123456asdf", "zeta1l40mm7meacx03r4lp87s9gkxfan32xnznp42u6")
	require.NoError(t, err)
	require.Equal(t, expectedOutput.HasVoted, resp)
}

func TestZetacore_GetZetaHotKeyBalance(t *testing.T) {
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
	server := setupMockServer(t, banktypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	client, err := setupZetacoreClient()
	require.NoError(t, err)
	client.keys = keys.NewKeysWithKeybase(mocks.NewKeyring(), types.AccAddress{}, "", "")

	resp, err := client.GetZetaHotKeyBalance()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Balance.Amount, resp)
}
