package zetacore

import (
	"bytes"
	"context"
	"encoding/hex"
	"net"
	"os"
	"testing"

	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/go-tss/blame"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

const (
	testSigner   = "jack"
	sampleHash   = "FA51DB4412144F1130669F2BAE8CB44AADBD8D85958DBFFCB0FE236878097E1A"
	ethBlockHash = "1a17bcc359e84ba8ae03b17ec425f97022cd11c3e279f6bdf7a96fcffa12b366"
)

func Test_GasPriceMultiplier(t *testing.T) {
	tt := []struct {
		name       string
		chain      chains.Chain
		multiplier float64
		fail       bool
	}{
		{
			name:       "get Ethereum multiplier",
			chain:      chains.Ethereum,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Goerli multiplier",
			chain:      chains.Goerli,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get BSC multiplier",
			chain:      chains.BscMainnet,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get BSC Testnet multiplier",
			chain:      chains.BscTestnet,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Polygon multiplier",
			chain:      chains.Polygon,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Mumbai Testnet multiplier",
			chain:      chains.Mumbai,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Bitcoin multiplier",
			chain:      chains.BitcoinMainnet,
			multiplier: 2.0,
			fail:       false,
		},
		{
			name:       "get Bitcoin Testnet multiplier",
			chain:      chains.BitcoinTestnet,
			multiplier: 2.0,
			fail:       false,
		},
		{
			name: "get unknown chain gas price multiplier",
			chain: chains.Chain{
				Consensus: chains.Consensus_tendermint,
			},
			multiplier: 1.0,
			fail:       true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			multiplier, err := GasPriceMultiplier(tc.chain)
			if tc.fail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.multiplier, multiplier)
		})
	}
}

func getHeaderData(t *testing.T) proofs.HeaderData {
	var header ethtypes.Header
	file, err := os.Open("../../testutil/testdata/eth_header_18495266.json")
	require.NoError(t, err)
	defer file.Close()
	headerBytes := make([]byte, 4096)
	n, err := file.Read(headerBytes)
	require.NoError(t, err)
	err = header.UnmarshalJSON(headerBytes[:n])
	require.NoError(t, err)
	var buffer bytes.Buffer
	err = header.EncodeRLP(&buffer)
	require.NoError(t, err)
	return proofs.NewEthereumHeader(buffer.Bytes())
}

func TestZetacore_PostGasPrice(t *testing.T) {
	ctx := context.Background()

	extraGRPC := withDummyServer(100)
	setupMockServer(t, observertypes.RegisterQueryServer, skipMethod, nil, nil, extraGRPC...)

	client := setupZetacoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	t.Run("post gas price success", func(t *testing.T) {
		hash, err := client.PostVoteGasPrice(ctx, chains.BscMainnet, 1000000, "100", 1234)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})

	// Test for failed broadcast, it will take several seconds to complete. Excluding to reduce runtime.
	//
	//t.Run("post gas price fail", func(t *testing.T) {
	//	zetacoreBroadcast = MockBroadcastError
	//	hash, err := client.PostGasPrice(chains.BscMainnet, 1000000, "100", 1234)
	//	require.ErrorContains(t, err, "post gasprice failed")
	//	require.Equal(t, "", hash)
	//})
}

func TestZetacore_AddOutboundTracker(t *testing.T) {
	ctx := context.Background()

	const nonce = 123
	chainID := chains.BscMainnet.ChainId

	method := "/zetachain.zetacore.crosschain.Query/OutboundTracker"
	input := &crosschaintypes.QueryGetOutboundTrackerRequest{
		ChainID: chains.BscMainnet.ChainId,
		Nonce:   nonce,
	}
	output := &crosschaintypes.QueryGetOutboundTrackerResponse{
		OutboundTracker: crosschaintypes.OutboundTracker{
			Index:    "456",
			ChainId:  chainID,
			Nonce:    nonce,
			HashList: nil,
		},
	}

	extraGRPC := withDummyServer(100)
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, output, extraGRPC...)

	tendermintMock := mocks.NewSDKClientWithErr(t, nil, 0)

	client := setupZetacoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withTendermint(tendermintMock),
	)

	t.Run("add tx hash success", func(t *testing.T) {
		tendermintMock.SetBroadcastTxHash(sampleHash)
		hash, err := client.AddOutboundTracker(ctx, chainID, nonce, "", nil, "", 456)
		assert.NoError(t, err)
		assert.Equal(t, sampleHash, hash)
	})

	t.Run("add tx hash fail", func(t *testing.T) {
		tendermintMock.SetError(errors.New("broadcast error"))
		hash, err := client.AddOutboundTracker(ctx, chainID, nonce, "", nil, "", 456)
		assert.Error(t, err)
		assert.Empty(t, hash)
	})
}

func TestZetacore_SetTSS(t *testing.T) {
	ctx := context.Background()

	extraGRPC := withDummyServer(100)
	setupMockServer(t, crosschaintypes.RegisterMsgServer, skipMethod, nil, nil, extraGRPC...)

	client := setupZetacoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	t.Run("set tss success", func(t *testing.T) {
		hash, err := client.PostVoteTSS(
			ctx,
			"zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			9987,
			chains.ReceiveStatus_success,
		)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestZetacore_UpdateZetacoreContext(t *testing.T) {
	ctx := context.Background()

	//Setup server for multiple grpc calls
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)

	server := grpcmock.MockUnstartedServer(
		grpcmock.RegisterService(crosschaintypes.RegisterQueryServer),
		grpcmock.RegisterService(upgradetypes.RegisterQueryServer),
		grpcmock.RegisterService(observertypes.RegisterQueryServer),
		grpcmock.RegisterService(lightclienttypes.RegisterQueryServer),
		grpcmock.RegisterService(authoritytypes.RegisterQueryServer),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.WithListener(listener),
		func(s *grpcmock.Server) {
			method := "/zetachain.zetacore.crosschain.Query/LastZetaHeight"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(crosschaintypes.QueryLastZetaHeightRequest{}).
				Return(crosschaintypes.QueryLastZetaHeightResponse{Height: 12345})

			method = "/cosmos.upgrade.v1beta1.Query/CurrentPlan"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(upgradetypes.QueryCurrentPlanRequest{}).
				Return(upgradetypes.QueryCurrentPlanResponse{
					Plan: &upgradetypes.Plan{
						Name:   "big upgrade",
						Height: 100,
					},
				})

			method = "/zetachain.zetacore.observer.Query/GetChainParams"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(observertypes.QueryGetChainParamsRequest{}).
				Return(observertypes.QueryGetChainParamsResponse{ChainParams: &observertypes.ChainParamsList{
					ChainParams: []*observertypes.ChainParams{
						{
							ChainId: 7000,
						},
					},
				}})

			method = "/zetachain.zetacore.observer.Query/SupportedChains"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(observertypes.QuerySupportedChains{}).
				Return(observertypes.QuerySupportedChainsResponse{
					Chains: []chains.Chain{
						{
							chains.BitcoinMainnet.ChainId,
							chains.BitcoinMainnet.ChainName,
							chains.BscMainnet.Network,
							chains.BscMainnet.NetworkType,
							chains.BscMainnet.Vm,
							chains.BscMainnet.Consensus,
							chains.BscMainnet.IsExternal,
							chains.BscMainnet.CctxGateway,
						},
						{
							chains.Ethereum.ChainId,
							chains.Ethereum.ChainName,
							chains.Ethereum.Network,
							chains.Ethereum.NetworkType,
							chains.Ethereum.Vm,
							chains.Ethereum.Consensus,
							chains.Ethereum.IsExternal,
							chains.Ethereum.CctxGateway,
						},
					},
				})

			method = "/zetachain.zetacore.observer.Query/Keygen"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(observertypes.QueryGetKeygenRequest{}).
				Return(observertypes.QueryGetKeygenResponse{
					Keygen: &observertypes.Keygen{
						Status:         observertypes.KeygenStatus_KeyGenSuccess,
						GranteePubkeys: nil,
						BlockNumber:    5646,
					}})

			method = "/zetachain.zetacore.observer.Query/TSS"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(observertypes.QueryGetTSSRequest{}).
				Return(observertypes.QueryGetTSSResponse{
					TSS: observertypes.TSS{
						TssPubkey:           "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
						TssParticipantList:  nil,
						OperatorAddressList: nil,
						FinalizedZetaHeight: 1000,
						KeyGenZetaHeight:    900,
					},
				})

			method = "/zetachain.zetacore.observer.Query/CrosschainFlags"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(observertypes.QueryGetCrosschainFlagsRequest{}).
				Return(observertypes.QueryGetCrosschainFlagsResponse{CrosschainFlags: observertypes.CrosschainFlags{
					IsInboundEnabled:      true,
					IsOutboundEnabled:     false,
					GasPriceIncreaseFlags: nil,
				}})

			method = "/zetachain.zetacore.lightclient.Query/HeaderEnabledChains"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(lightclienttypes.QueryHeaderEnabledChainsRequest{}).
				Return(lightclienttypes.QueryHeaderEnabledChainsResponse{HeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{
					{
						ChainId: chains.Ethereum.ChainId,
						Enabled: true,
					},
					{
						ChainId: chains.BitcoinMainnet.ChainId,
						Enabled: false,
					},
				}})

			method = "/zetachain.zetacore.authority.Query/ChainInfo"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(authoritytypes.QueryGetChainInfoRequest{}).
				Return(authoritytypes.QueryGetChainInfoResponse{
					ChainInfo: authoritytypes.ChainInfo{
						Chains: []chains.Chain{
							sample.Chain(1000),
							sample.Chain(1001),
							sample.Chain(1002),
						},
					},
				})
		},
	)(t)

	server.Serve()
	defer server.Close()

	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())
	client := setupZetacoreClient(t,
		withObserverKeys(keys.NewKeysWithKeybase(mocks.NewKeyring(), address, testSigner, "")),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0)),
	)

	t.Run("zetacore update success", func(t *testing.T) {
		cfg := config.New(false)
		appContext := zctx.New(cfg, zerolog.Nop())
		err := client.UpdateZetacoreContext(ctx, appContext, false, zerolog.Logger{})
		require.NoError(t, err)
	})
}

func TestZetacore_PostBlameData(t *testing.T) {
	ctx := context.Background()

	extraGRPC := withDummyServer(100)
	setupMockServer(t, observertypes.RegisterQueryServer, skipMethod, nil, nil, extraGRPC...)

	client := setupZetacoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	t.Run("post blame data success", func(t *testing.T) {
		hash, err := client.PostVoteBlameData(
			ctx,
			&blame.Blame{
				FailReason: "",
				IsUnicast:  false,
				BlameNodes: nil,
			},
			chains.BscMainnet.ChainId,
			"102394876-bsc",
		)
		assert.NoError(t, err)
		assert.Equal(t, sampleHash, hash)
	})
}

func TestZetacore_PostVoteBlockHeader(t *testing.T) {
	ctx := context.Background()

	extraGRPC := withDummyServer(100)
	setupMockServer(t, observertypes.RegisterQueryServer, skipMethod, nil, nil, extraGRPC...)

	client := setupZetacoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	blockHash, err := hex.DecodeString(ethBlockHash)
	require.NoError(t, err)

	t.Run("post add block header success", func(t *testing.T) {
		hash, err := client.PostVoteBlockHeader(
			ctx,
			chains.Ethereum.ChainId,
			blockHash,
			18495266,
			getHeaderData(t),
		)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestZetacore_PostVoteInbound(t *testing.T) {
	ctx := context.Background()

	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())

	expectedOutput := observertypes.QueryHasVotedResponse{HasVoted: false}
	input := observertypes.QueryHasVotedRequest{
		BallotIdentifier: "0x2d10e9b7ce7921fa6b61ada3020d1c797d5ec52424cdcf86ef31cbbbcd45db58",
		VoterAddress:     address.String(),
	}
	method := "/zetachain.zetacore.observer.Query/HasVoted"

	extraGRPC := withDummyServer(100)
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput, extraGRPC...)

	client := setupZetacoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	t.Run("post inbound vote already voted", func(t *testing.T) {
		hash, _, err := client.PostVoteInbound(ctx, 100, 200, &crosschaintypes.MsgVoteInbound{
			Creator: address.String(),
		})
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestZetacore_GetInboundVoteMessage(t *testing.T) {
	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())
	t.Run("get inbound vote message", func(t *testing.T) {
		msg := GetInboundVoteMessage(
			address.String(),
			chains.Ethereum.ChainId,
			"",
			address.String(),
			chains.ZetaChainMainnet.ChainId,
			math.NewUint(500),
			"",
			"", 12345,
			1000,
			coin.CoinType_Gas,
			"azeta",
			address.String(),
			0)
		require.Equal(t, address.String(), msg.Creator)
	})
}

func TestZetacore_MonitorVoteInboundResult(t *testing.T) {
	ctx := context.Background()

	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())
	client := setupZetacoreClient(t,
		withObserverKeys(keys.NewKeysWithKeybase(mocks.NewKeyring(), address, testSigner, "")),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0)),
	)

	t.Run("monitor inbound vote", func(t *testing.T) {
		err := client.MonitorVoteInboundResult(ctx, sampleHash, 1000, &crosschaintypes.MsgVoteInbound{
			Creator: address.String(),
		})

		require.NoError(t, err)
	})
}

func TestZetacore_PostVoteOutbound(t *testing.T) {
	const (
		blockHeight = 1234
		accountNum  = 10
		accountSeq  = 10
	)

	ctx := context.Background()

	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())

	expectedOutput := observertypes.QueryHasVotedResponse{HasVoted: false}
	input := observertypes.QueryHasVotedRequest{
		BallotIdentifier: "0xf52f379287561dd07869de72b09fb56b7f6dfdda65b01c25882722e315f333f1",
		VoterAddress:     address.String(),
	}
	method := "/zetachain.zetacore.observer.Query/HasVoted"

	extraGRPC := withDummyServer(blockHeight)

	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput, extraGRPC...)
	require.NotNil(t, server)

	client := setupZetacoreClient(t,
		withDefaultObserverKeys(),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
		withAccountRetriever(t, accountNum, accountSeq),
	)

	msg := crosschaintypes.NewMsgVoteOutbound(
		address.String(),
		sampleHash,
		sampleHash,
		blockHeight,
		1000,
		math.NewInt(100),
		1200,
		math.NewUint(500),
		chains.ReceiveStatus_success,
		chains.Ethereum.ChainId,
		10001,
		coin.CoinType_Gas,
	)

	hash, ballot, err := client.PostVoteOutbound(ctx, 100_000, 200_000, msg)

	assert.NoError(t, err)
	assert.Equal(t, sampleHash, hash)
	assert.Equal(t, "0xf52f379287561dd07869de72b09fb56b7f6dfdda65b01c25882722e315f333f1", ballot)
}

func TestZetacore_MonitorVoteOutboundResult(t *testing.T) {
	ctx := context.Background()

	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())
	client := setupZetacoreClient(t,
		withObserverKeys(keys.NewKeysWithKeybase(mocks.NewKeyring(), address, testSigner, "")),
		withTendermint(mocks.NewSDKClientWithErr(t, nil, 0)),
	)

	t.Run("monitor outbound vote", func(t *testing.T) {
		msg := &crosschaintypes.MsgVoteOutbound{Creator: address.String()}

		err := client.MonitorVoteOutboundResult(ctx, sampleHash, 1000, msg)
		assert.NoError(t, err)
	})
}
