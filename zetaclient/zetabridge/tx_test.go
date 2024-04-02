package zetabridge

import (
	"bytes"
	"encoding/hex"
	"errors"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/go-tss/blame"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
	zetabridge, err := setupCorBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewMockKeyring(), address, "", "")

	t.Run("post gas price success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.PostGasPrice(chains.BscMainnetChain(), 1000000, "100", 1234)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})

	// Test for failed broadcast, it will take several seconds to complete. Excluding to reduce runtime.
	//
	//t.Run("post gas price fail", func(t *testing.T) {
	//	zetaBridgeBroadcast = ZetaBridgeBroadcastTestErr
	//	hash, err := zetabridge.PostGasPrice(chains.BscMainnetChain(), 1000000, "100", 1234)
	//	require.ErrorContains(t, err, "post gasprice failed")
	//	require.Equal(t, "", hash)
	//})
}

func TestZetaCoreBridge_AddTxHashToOutTxTracker(t *testing.T) {
	zetabridge, err := setupCorBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewMockKeyring(), address, "", "")

	t.Run("add tx hash success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.AddTxHashToOutTxTracker(chains.BscMainnetChain().ChainId, 123, "", nil, "", 456)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})

	t.Run("add tx hash fail", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTestErr
		hash, err := zetabridge.AddTxHashToOutTxTracker(chains.BscMainnetChain().ChainId, 123, "", nil, "", 456)
		require.Error(t, err)
		require.Equal(t, "", hash)
	})
}

func TestZetaCoreBridge_SetTSS(t *testing.T) {
	zetabridge, err := setupCorBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewMockKeyring(), address, "", "")

	t.Run("set tss success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.SetTSS(
			"zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			9987,
			chains.ReceiveStatus_Success,
		)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestZetaCoreBridge_CoreContextUpdater(t *testing.T) {
	zetabridge, err := setupCorBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewMockKeyring(), address, "", "")

	//t.Run("core context update success", func(t *testing.T) {
	//	zetaBridgeBroadcast = ZetaBridgeBroadcastTest
	//	hash, err := zetabridge.CoreContextUpdater()
	//	require.NoError(t, err)
	//	require.Equal(t, sampleHash, hash)
	//})
}

func TestZetaCoreBridge_PostBlameData(t *testing.T) {
	zetabridge, err := setupCorBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewMockKeyring(), address, "", "")

	t.Run("post blame data success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.PostBlameData(
			&blame.Blame{
				FailReason: "",
				IsUnicast:  false,
				BlameNodes: nil,
			},
			chains.BscMainnetChain().ChainId,
			"102394876-bsc",
		)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestZetaCoreBridge_PostAddBlockHeader(t *testing.T) {
	zetabridge, err := setupCorBridge()
	require.NoError(t, err)
	address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewMockKeyring(), address, "", "")
	blockHash, err := hex.DecodeString(ethBlockHash)
	require.NoError(t, err)

	t.Run("post add block header success", func(t *testing.T) {
		zetaBridgeBroadcast = ZetaBridgeBroadcastTest
		hash, err := zetabridge.PostAddBlockHeader(
			chains.EthChain().ChainId,
			blockHash,
			18495266,
			getHeaderData(t),
		)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

// TODO: Rest of tests Requires zetabridge refactoring
func TestZetaCoreBridge_PostVoteInbound(t *testing.T) {
	//zetabridge, err := setupCorBridge()
	//require.NoError(t, err)
	//address := sdktypes.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	//zetabridge.keys = keys.NewKeysWithKeybase(stub.NewMockKeyring(), address, "", "")
	//
	//t.Run("post add block header success", func(t *testing.T) {
	//	zetaBridgeBroadcast = ZetaBridgeBroadcastTest
	//	hash, _, err := zetabridge.PostVoteInbound(100, 200, &crosschaintypes.MsgVoteOnObservedInboundTx{
	//		Creator: address.String(),
	//	})
	//	require.NoError(t, err)
	//	require.Equal(t, sampleHash, hash)
	//})
}
func TestZetaCoreBridge_MonitorVoteInboundTxResult(t *testing.T) {
}

func TestZetaCoreBridge_PostVoteOutbound(t *testing.T) {
}

func TestZetaCoreBridge_PostVoteOutboundFromMsg(t *testing.T) {
}

func TestZetaCoreBridge_MonitorVoteOutboundTxResult(t *testing.T) {
}
