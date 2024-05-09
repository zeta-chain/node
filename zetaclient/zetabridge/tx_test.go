package zetabridge

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"net"
	"os"
	"testing"

	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/go-tss/blame"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"
)

const (
	sampleHash   = "fa51db4412144f1130669f2bae8cb44aadbd8d85958dbffcb0fe236878097e1a"
	ethBlockHash = "1a17bcc359e84ba8ae03b17ec425f97022cd11c3e279f6bdf7a96fcffa12b366"
)

func Test_GasPriceMultiplier(t *testing.T) {
	tt := []struct {
		name       string
		chainID    int64
		multiplier float64
		fail       bool
	}{
		{
			name:       "get Ethereum multiplier",
			chainID:    1,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Goerli multiplier",
			chainID:    5,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get BSC multiplier",
			chainID:    56,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get BSC Testnet multiplier",
			chainID:    97,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Polygon multiplier",
			chainID:    137,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Mumbai Testnet multiplier",
			chainID:    80001,
			multiplier: 1.2,
			fail:       false,
		},
		{
			name:       "get Bitcoin multiplier",
			chainID:    8332,
			multiplier: 2.0,
			fail:       false,
		},
		{
			name:       "get Bitcoin Testnet multiplier",
			chainID:    18332,
			multiplier: 2.0,
			fail:       false,
		},
		{
			name:       "get unknown chain gas price multiplier",
			chainID:    1234,
			multiplier: 1.0,
			fail:       true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			multiplier, err := GasPriceMultiplier(tc.chainID)
			if tc.fail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.multiplier, multiplier)
		})
	}
}

func ZetaBridgeBroadcastTest(_ *ZetaCoreBridge, _ uint64, _ sdktypes.Msg, _ authz.Signer) (string, error) {
	return sampleHash, nil
}

func ZetaBridgeBroadcastTestErr(_ *ZetaCoreBridge, _ uint64, _ sdktypes.Msg, _ authz.Signer) (string, error) {
	return sampleHash, errors.New("broadcast error")
}

func getHeaderData(t *testing.T) proofs.HeaderData {
	var header ethtypes.Header
	file, err := os.Open("../../pkg/testdata/eth_header_18495266.json")
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

func TestZetaCoreBridge_PostGasPrice(t *testing.T) {
	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")

	t.Run("post gas price success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.PostGasPrice(chains.BscMainnetChain, 1000000, "100", 1234)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})

	// Test for failed broadcast, it will take several seconds to complete. Excluding to reduce runtime.
	//
	//t.Run("post gas price fail", func(t *testing.T) {
	//	zetaBridgeBroadcast = ZetaBridgeBroadcastTestErr
	//	hash, err := zetabridge.PostGasPrice(chains.BscMainnetChain, 1000000, "100", 1234)
	//	require.ErrorContains(t, err, "post gasprice failed")
	//	require.Equal(t, "", hash)
	//})
}

func TestZetaCoreBridge_AddTxHashToOutboundTracker(t *testing.T) {
	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")

	t.Run("add tx hash success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.AddTxHashToOutboundTracker(chains.BscMainnetChain.ChainId, 123, "", nil, "", 456)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})

	t.Run("add tx hash fail", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTestErr
		hash, err := zetabridge.AddTxHashToOutboundTracker(chains.BscMainnetChain.ChainId, 123, "", nil, "", 456)
		require.Error(t, err)
		require.Equal(t, "", hash)
	})
}

func TestZetaCoreBridge_SetTSS(t *testing.T) {
	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")

	t.Run("set tss success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.SetTSS(
			"zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			9987,
			chains.ReceiveStatus_success,
		)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestZetaCoreBridge_UpdateZetaCoreContext(t *testing.T) {
	//Setup server for multiple grpc calls
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)

	server := grpcmock.MockUnstartedServer(
		grpcmock.RegisterService(crosschaintypes.RegisterQueryServer),
		grpcmock.RegisterService(upgradetypes.RegisterQueryServer),
		grpcmock.RegisterService(observertypes.RegisterQueryServer),
		grpcmock.RegisterService(lightclienttypes.RegisterQueryServer),
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
					Chains: []*chains.Chain{

						{
							chains.BtcMainnetChain.ChainName,
							chains.BtcMainnetChain.ChainId,
							chains.BscMainnetChain.Network,
							chains.BscMainnetChain.NetworkType,
							chains.BscMainnetChain.Vm,
							chains.BscMainnetChain.Consensus,
							chains.BscMainnetChain.IsExternal,
						},
						{
							chains.EthChain.ChainName,
							chains.EthChain.ChainId,
							chains.EthChain.Network,
							chains.EthChain.NetworkType,
							chains.EthChain.Vm,
							chains.EthChain.Consensus,
							chains.EthChain.IsExternal,
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
					IsInboundEnabled:             true,
					IsOutboundEnabled:            false,
					GasPriceIncreaseFlags:        nil,
					BlockHeaderVerificationFlags: nil,
				}})

			method = "/zetachain.zetacore.lightclient.Query/HeaderEnabledChains"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(lightclienttypes.QueryHeaderEnabledChainsRequest{}).
				Return(lightclienttypes.QueryHeaderEnabledChainsResponse{HeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{
					{
						ChainId: chains.EthChain.ChainId,
						Enabled: true,
					},
					{
						ChainId: chains.BtcMainnetChain.ChainId,
						Enabled: false,
					},
				}})
		},
	)(t)

	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")
	zetabridge.EnableMockSDKClient(stub.NewSDKClientWithErr(nil, 0))

	t.Run("core context update success", func(t *testing.T) {
		cfg := config.NewConfig()
		coreCtx := corecontext.NewZetaCoreContext(cfg)
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		err := zetabridge.UpdateZetaCoreContext(coreCtx, false, zerolog.Logger{})
		require.NoError(t, err)
	})
}

func TestZetaCoreBridge_PostBlameData(t *testing.T) {
	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")

	t.Run("post blame data success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.PostBlameData(
			&blame.Blame{
				FailReason: "",
				IsUnicast:  false,
				BlameNodes: nil,
			},
			chains.BscMainnetChain.ChainId,
			"102394876-bsc",
		)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestZetaCoreBridge_PostVoteBlockHeader(t *testing.T) {
	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")
	blockHash, err := hex.DecodeString(ethBlockHash)
	require.NoError(t, err)

	t.Run("post add block header success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.PostVoteBlockHeader(
			chains.EthChain.ChainId,
			blockHash,
			18495266,
			getHeaderData(t),
		)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestZetaCoreBridge_PostVoteInbound(t *testing.T) {
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())

	expectedOutput := observertypes.QueryHasVotedResponse{HasVoted: false}
	input := observertypes.QueryHasVotedRequest{
		BallotIdentifier: "0x2d10e9b7ce7921fa6b61ada3020d1c797d5ec52424cdcf86ef31cbbbcd45db58",
		VoterAddress:     address.String(),
	}
	method := "/zetachain.zetacore.observer.Query/HasVoted"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")
	zetabridge.EnableMockSDKClient(stub.NewSDKClientWithErr(nil, 0))

	t.Run("post inbound vote already voted", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, _, err := zetabridge.PostVoteInbound(100, 200, &crosschaintypes.MsgVoteInbound{
			Creator: address.String(),
		})
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestZetaCoreBridge_GetInboundVoteMessage(t *testing.T) {
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	t.Run("get inbound vote message", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		msg := GetInboundVoteMessage(
			address.String(),
			chains.EthChain.ChainId,
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

func TestZetaCoreBridge_MonitorVoteInboundTxResult(t *testing.T) {
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")
	zetabridge.EnableMockSDKClient(stub.NewSDKClientWithErr(nil, 0))

	t.Run("monitor inbound vote", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		zetabridge.MonitorVoteInboundTxResult(sampleHash, 1000, &crosschaintypes.MsgVoteInbound{
			Creator: address.String(),
		})
		// Nothing to verify against this function
		// Just running through without panic
	})
}

func TestZetaCoreBridge_PostVoteOutbound(t *testing.T) {
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())

	expectedOutput := observertypes.QueryHasVotedResponse{HasVoted: false}
	input := observertypes.QueryHasVotedRequest{
		BallotIdentifier: "0x6facd1e8a2008cd0acd7784628ff1ac326406f3bb355ea514bf290c6ce6f8355",
		VoterAddress:     address.String(),
	}
	method := "/zetachain.zetacore.observer.Query/HasVoted"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")
	zetabridge.EnableMockSDKClient(stub.NewSDKClientWithErr(nil, 0))

	zetaBridgeBroadcast = ZetaBridgeBroadcastTest
	hash, ballot, err := zetabridge.PostVoteOutbound(
		sampleHash,
		sampleHash,
		1234,
		1000,
		big.NewInt(100),
		1200,
		big.NewInt(500),
		chains.ReceiveStatus_success,
		chains.EthChain,
		10001,
		coin.CoinType_Gas)
	require.NoError(t, err)
	require.Equal(t, sampleHash, hash)
	require.Equal(t, "0x6facd1e8a2008cd0acd7784628ff1ac326406f3bb355ea514bf290c6ce6f8355", ballot)
}

func TestZetaCoreBridge_MonitorVoteOutboundTxResult(t *testing.T) {
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")
	zetabridge.EnableMockSDKClient(stub.NewSDKClientWithErr(nil, 0))

	t.Run("monitor outbound vote", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		zetabridge.MonitorVoteOutboundTxResult(sampleHash, 1000, &crosschaintypes.MsgVoteOutbound{
			Creator: address.String(),
		})
		// Nothing to verify against this function
		// Just running through without panic
	})
}
