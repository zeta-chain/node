package keeper_test

import (
	"encoding/json"
	"os"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/proofs"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
)

// readHeader reads a header from a file.
// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func readHeader(filename string) (*ethtypes.Header, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var NewHeader ethtypes.Header
	err = decoder.Decode(&NewHeader)
	return &NewHeader, err
}

func ethHeaders() (*ethtypes.Header, *ethtypes.Header, *ethtypes.Header, error) {
	header1, err := readHeader("./testdata/header_sepolia_5000000.json")
	if err != nil {
		return nil, nil, nil, err
	}
	header2, err := readHeader("./testdata/header_sepolia_5000001.json")
	if err != nil {
		return nil, nil, nil, err
	}
	header3, err := readHeader("./testdata/header_sepolia_5000002.json")
	if err != nil {
		return nil, nil, nil, err
	}
	return header1, header2, header3, nil
}

// sepoliaBlockHeaders returns three block headers for the Sepolia chain.
func sepoliaBlockHeaders(t *testing.T) (proofs.BlockHeader, proofs.BlockHeader, proofs.BlockHeader) {
	header1, header2, header3, err := ethHeaders()
	if err != nil {
		panic(err)
	}

	headerRLP1, err := rlp.EncodeToBytes(header1)
	require.NoError(t, err)
	headerRLP2, err := rlp.EncodeToBytes(header1)
	require.NoError(t, err)
	headerRLP3, err := rlp.EncodeToBytes(header1)
	require.NoError(t, err)

	return proofs.BlockHeader{
			Height:     5000000,
			Hash:       header1.Hash().Bytes(),
			ParentHash: header1.ParentHash.Bytes(),
			ChainId:    chains.Sepolia.ChainId,
			Header:     proofs.NewEthereumHeader(headerRLP1),
		},
		proofs.BlockHeader{
			Height:     5000001,
			Hash:       header2.Hash().Bytes(),
			ParentHash: header2.ParentHash.Bytes(),
			ChainId:    chains.Sepolia.ChainId,
			Header:     proofs.NewEthereumHeader(headerRLP2),
		},
		proofs.BlockHeader{
			Height:     5000002,
			Hash:       header3.Hash().Bytes(),
			ParentHash: header3.ParentHash.Bytes(),
			ChainId:    chains.Sepolia.ChainId,
			Header:     proofs.NewEthereumHeader(headerRLP3),
		}
}

// TestKeeper_GetBlockHeader tests get, set, and remove block header
func TestKeeper_GetBlockHeader(t *testing.T) {
	k, ctx, _, _ := keepertest.LightclientKeeper(t)
	blockHash := sample.Hash().Bytes()
	_, found := k.GetBlockHeader(ctx, blockHash)
	require.False(t, found)

	k.SetBlockHeader(ctx, sample.BlockHeader(blockHash))
	_, found = k.GetBlockHeader(ctx, blockHash)
	require.True(t, found)

	k.RemoveBlockHeader(ctx, blockHash)
	_, found = k.GetBlockHeader(ctx, blockHash)
	require.False(t, found)
}

func TestKeeper_GetAllBlockHeaders(t *testing.T) {
	k, ctx, _, _ := keepertest.LightclientKeeper(t)
	b1 := sample.BlockHeader(sample.Hash().Bytes())
	b2 := sample.BlockHeader(sample.Hash().Bytes())
	b3 := sample.BlockHeader(sample.Hash().Bytes())

	k.SetBlockHeader(ctx, b1)
	k.SetBlockHeader(ctx, b2)
	k.SetBlockHeader(ctx, b3)

	list := k.GetAllBlockHeaders(ctx)
	require.Len(t, list, 3)
	require.Contains(t, list, b1)
	require.Contains(t, list, b2)
	require.Contains(t, list, b3)
}

// TODO: Test with bitcoin headers
// https://github.com/zeta-chain/node/issues/1994

func TestKeeper_CheckNewBlockHeader(t *testing.T) {
	t.Run("should succeed if block header is valid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Sepolia.ChainId,
					Enabled: true,
				},
			},
		})

		bh, _, _ := sepoliaBlockHeaders(t)

		parentHash, err := k.CheckNewBlockHeader(ctx, bh.ChainId, bh.Hash, bh.Height, bh.Header)
		require.NoError(t, err)
		require.Equal(t, bh.ParentHash, parentHash)
	})

	t.Run("should fail if verification flag not enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Sepolia.ChainId,
					Enabled: false,
				},
			},
		})

		bh, _, _ := sepoliaBlockHeaders(t)

		_, err := k.CheckNewBlockHeader(ctx, bh.ChainId, bh.Hash, bh.Height, bh.Header)
		require.ErrorIs(t, err, types.ErrBlockHeaderVerificationDisabled)
	})

	// TODO: Fix the Sepolia sample headers to enable this test:
	// https://github.com/zeta-chain/node/issues/1996
	//t.Run("should succeed if block header is valid with chain state existing", func(t *testing.T) {
	//	k, ctx, _, _ := keepertest.LightclientKeeper(t)
	//
	//	k.SetHeaderSupportedChain(ctx, types.HeaderSupportedChain{
	//		EthTypeChainEnabled: true,
	//	})
	//
	//	bh1, bh2, _ := sepoliaBlockHeaders(t)
	//	k.SetChainState(ctx, types.ChainState{
	//		ChainId:         bh1.ChainId,
	//		LatestHeight:    bh1.Height,
	//		EarliestHeight:  bh1.Height - 100,
	//		LatestBlockHash: bh1.Hash,
	//	})
	//	k.SetBlockHeader(ctx, bh1)
	//
	//	parentHash, err := k.CheckNewBlockHeader(ctx, bh2.ChainId, bh2.Hash, bh2.Height, bh2.Header)
	//	require.NoError(t, err)
	//	require.Equal(t, bh2.ParentHash, parentHash)
	//})

	t.Run("fail if block already exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Sepolia.ChainId,
					Enabled: true,
				},
			},
		})
		bh, _, _ := sepoliaBlockHeaders(t)
		k.SetBlockHeader(ctx, bh)

		_, err := k.CheckNewBlockHeader(ctx, bh.ChainId, bh.Hash, bh.Height, bh.Header)
		require.ErrorIs(t, err, types.ErrBlockAlreadyExist)
	})

	t.Run("fail if chain state and invalid height", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Sepolia.ChainId,
					Enabled: true,
				},
			},
		})
		bh, _, _ := sepoliaBlockHeaders(t)

		k.SetChainState(ctx, types.ChainState{
			ChainId:         bh.ChainId,
			LatestHeight:    bh.Height - 2,
			EarliestHeight:  bh.Height - 100,
			LatestBlockHash: bh.Hash,
		})

		_, err := k.CheckNewBlockHeader(ctx, bh.ChainId, bh.Hash, bh.Height, bh.Header)
		require.ErrorIs(t, err, types.ErrInvalidHeight)
	})

	t.Run("fail if chain state and no parent", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Sepolia.ChainId,
					Enabled: true,
				},
			},
		})

		bh, _, _ := sepoliaBlockHeaders(t)

		k.SetChainState(ctx, types.ChainState{
			ChainId:         bh.ChainId,
			LatestHeight:    bh.Height - 1,
			EarliestHeight:  bh.Height - 100,
			LatestBlockHash: bh.Hash,
		})

		_, err := k.CheckNewBlockHeader(ctx, bh.ChainId, bh.Hash, bh.Height, bh.Header)
		require.ErrorIs(t, err, types.ErrNoParentHash)
	})
}

func TestKeeper_AddBlockHeader(t *testing.T) {
	t.Run("should add a block header and create chain state if doesn't exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Sepolia.ChainId,
					Enabled: true,
				},
			},
		})
		bh, _, _ := sepoliaBlockHeaders(t)

		k.AddBlockHeader(ctx, bh.ChainId, bh.Height, bh.Hash, bh.Header, bh.ParentHash)

		retrieved, found := k.GetBlockHeader(ctx, bh.Hash)
		require.True(t, found)
		require.EqualValues(t, bh.Header, retrieved.Header)
		require.EqualValues(t, bh.Height, retrieved.Height)
		require.EqualValues(t, bh.Hash, retrieved.Hash)
		require.EqualValues(t, bh.ParentHash, retrieved.ParentHash)
		require.EqualValues(t, bh.ChainId, retrieved.ChainId)

		// Check chain state
		chainState, found := k.GetChainState(ctx, bh.ChainId)
		require.True(t, found)
		require.EqualValues(t, bh.Height, chainState.LatestHeight)
		require.EqualValues(t, bh.Height, chainState.EarliestHeight)
		require.EqualValues(t, bh.Hash, chainState.LatestBlockHash)
		require.EqualValues(t, bh.ChainId, chainState.ChainId)
	})

	t.Run("should add a block header and update chain state if exists", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Sepolia.ChainId,
					Enabled: true,
				},
			},
		})

		bh, _, _ := sepoliaBlockHeaders(t)

		k.SetChainState(ctx, types.ChainState{
			ChainId:         bh.ChainId,
			LatestHeight:    bh.Height - 1,
			EarliestHeight:  bh.Height - 100,
			LatestBlockHash: bh.ParentHash,
		})

		k.AddBlockHeader(ctx, bh.ChainId, bh.Height, bh.Hash, bh.Header, bh.ParentHash)

		retrieved, found := k.GetBlockHeader(ctx, bh.Hash)
		require.True(t, found)
		require.EqualValues(t, bh.Header, retrieved.Header)
		require.EqualValues(t, bh.Height, retrieved.Height)
		require.EqualValues(t, bh.Hash, retrieved.Hash)
		require.EqualValues(t, bh.ParentHash, retrieved.ParentHash)
		require.EqualValues(t, bh.ChainId, retrieved.ChainId)

		// Check chain state
		chainState, found := k.GetChainState(ctx, bh.ChainId)
		require.True(t, found)
		require.EqualValues(t, bh.Height, chainState.LatestHeight)
		require.EqualValues(t, bh.Height-100, chainState.EarliestHeight)
		require.EqualValues(t, bh.Hash, chainState.LatestBlockHash)
		require.EqualValues(t, bh.ChainId, chainState.ChainId)
	})

	t.Run(
		"should add a block header and update chain state if exists and set earliest height if 0",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.LightclientKeeper(t)

			k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
				HeaderSupportedChains: []types.HeaderSupportedChain{
					{
						ChainId: chains.Sepolia.ChainId,
						Enabled: true,
					},
				},
			})

			bh, _, _ := sepoliaBlockHeaders(t)

			k.SetChainState(ctx, types.ChainState{
				ChainId:         bh.ChainId,
				LatestHeight:    bh.Height - 1,
				EarliestHeight:  0,
				LatestBlockHash: bh.ParentHash,
			})

			k.AddBlockHeader(ctx, bh.ChainId, bh.Height, bh.Hash, bh.Header, bh.ParentHash)

			retrieved, found := k.GetBlockHeader(ctx, bh.Hash)
			require.True(t, found)
			require.EqualValues(t, bh.Header, retrieved.Header)
			require.EqualValues(t, bh.Height, retrieved.Height)
			require.EqualValues(t, bh.Hash, retrieved.Hash)
			require.EqualValues(t, bh.ParentHash, retrieved.ParentHash)
			require.EqualValues(t, bh.ChainId, retrieved.ChainId)

			// Check chain state
			chainState, found := k.GetChainState(ctx, bh.ChainId)
			require.True(t, found)
			require.EqualValues(t, bh.Height, chainState.LatestHeight)
			require.EqualValues(t, bh.Height, chainState.EarliestHeight)
			require.EqualValues(t, bh.Hash, chainState.LatestBlockHash)
			require.EqualValues(t, bh.ChainId, chainState.ChainId)
		},
	)
}
