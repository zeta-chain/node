package evm

import (
	"fmt"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"
)

func TestEVMBlockCache(t *testing.T) {
	// create client
	blockCache, err := lru.New(1000)
	require.NoError(t, err)
	blockCacheV3, err := lru.New(1000)
	require.NoError(t, err)
	ob := ChainClient{
		blockCache:   blockCache,
		blockCacheV3: blockCacheV3,
	}

	// delete non-existing block should not panic
	blockNumber := int64(10388180)
	// #nosec G701 possible nummber
	ob.RemoveCachedBlock(uint64(blockNumber))

	// add a block
	header := &ethtypes.Header{
		Number: big.NewInt(blockNumber),
	}
	block := ethtypes.NewBlock(header, nil, nil, nil, nil)
	ob.blockCache.Add(blockNumber, block)

	// delete the block should not panic
	ob.RemoveCachedBlock(uint64(blockNumber))
}

func TestCheckEvmTxLog(t *testing.T) {
	// test data
	connectorAddr := ethcommon.HexToAddress("0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea")
	txHash := "0xb252c9e77feafdeeae25cc1f037a16c4b50fa03c494754b99a7339d816c79626"
	topics := []ethcommon.Hash{
		// https://goerli.etherscan.io/tx/0xb252c9e77feafdeeae25cc1f037a16c4b50fa03c494754b99a7339d816c79626#eventlog
		ethcommon.HexToHash("0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4"),
		ethcommon.HexToHash("0x00000000000000000000000023856df5d563bd893fc7df864102d8bbfe7fc487"),
		ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000061"),
	}

	tests := []struct {
		name string
		vLog *ethtypes.Log
		fail bool
	}{
		{
			name: "chain reorganization",
			vLog: &ethtypes.Log{
				Removed: true,
				Address: connectorAddr,
				TxHash:  ethcommon.HexToHash(txHash),
				Topics:  topics,
			},
			fail: true,
		},
		{
			name: "emitter address mismatch",
			vLog: &ethtypes.Log{
				Removed: false,
				Address: ethcommon.HexToAddress("0x184ba627DB853244c9f17f3Cb4378cB8B39bf147"),
				TxHash:  ethcommon.HexToHash(txHash),
				Topics:  topics,
			},
			fail: true,
		},
		{
			name: "tx hash mismatch",
			vLog: &ethtypes.Log{
				Removed: false,
				Address: connectorAddr,
				TxHash:  ethcommon.HexToHash("0x781c018d604af4dad0fe5e3cea4ad9fb949a996d8cd0cd04a92cadd7f08c05f2"),
				Topics:  topics,
			},
			fail: true,
		},
		{
			name: "topics mismatch",
			vLog: &ethtypes.Log{
				Removed: false,
				Address: connectorAddr,
				TxHash:  ethcommon.HexToHash(txHash),
				Topics: []ethcommon.Hash{
					// https://goerli.etherscan.io/tx/0xb252c9e77feafdeeae25cc1f037a16c4b50fa03c494754b99a7339d816c79626#eventlog
					ethcommon.HexToHash("0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4"),
					ethcommon.HexToHash("0x00000000000000000000000023856df5d563bd893fc7df864102d8bbfe7fc487"),
				},
			},
			fail: true,
		},
		{
			name: "should pass",
			vLog: &ethtypes.Log{
				Removed: false,
				Address: connectorAddr,
				TxHash:  ethcommon.HexToHash(txHash),
				Topics:  topics,
			},
			fail: false,
		},
	}

	evmClient := ChainClient{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("check test: %s\n", tt.name)
			err := evmClient.CheckEvmTxLog(
				tt.vLog,
				connectorAddr,
				"0xb252c9e77feafdeeae25cc1f037a16c4b50fa03c494754b99a7339d816c79626",
				TopicsZetaSent,
			)
			if tt.fail {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}
