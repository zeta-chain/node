package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) VerifyProof(ctx sdk.Context, proof *common.Proof, hash string, txIndex int64, chainID int64) (ethtypes.Transaction, error) {
	var txx ethtypes.Transaction
	crosschainFlags, found := k.zetaObserverKeeper.GetCrosschainFlags(ctx)
	if !found {
		return txx, fmt.Errorf("crosschain flags not found")
	}
	if common.IsBitcoinChain(chainID) && !crosschainFlags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled {
		return txx, fmt.Errorf("cannot verify proof for bitcoin chain %d", chainID)
	}

	if common.IsEVMChain(chainID) && !crosschainFlags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled {
		return txx, fmt.Errorf("cannot verify proof for bitcoin chain %d ", chainID)
	}

	senderChain := common.GetChainFromChainID(chainID)
	if senderChain == nil {
		return txx, types.ErrUnsupportedChain
	}

	if !senderChain.IsProvable() {
		return txx, fmt.Errorf("chain %d does not support block header verification", chainID)
	}

	blockHash := eth.HexToHash(hash)

	res, found := k.zetaObserverKeeper.GetBlockHeader(ctx, blockHash.Bytes())
	if !found {
		return txx, fmt.Errorf("block header not found %s", blockHash)
	}

	// verify and process the proof
	val, err := proof.Verify(res.Header, int(txIndex))
	if err != nil && !common.IsErrorInvalidProof(err) {
		return txx, err
	}
	err = txx.UnmarshalBinary(val)
	if err != nil {
		return txx, err
	}
	return txx, nil
}
