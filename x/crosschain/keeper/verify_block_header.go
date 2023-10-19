package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) VerifyProof(ctx sdk.Context, proof *common.Proof, chainID int64, blockHash string, txIndex int64) ([]byte, error) {
	// header-based merkle proof verification must be enabled
	crosschainFlags, found := k.zetaObserverKeeper.GetCrosschainFlags(ctx)
	if !found {
		return nil, fmt.Errorf("crosschain flags not found")
	}
	if crosschainFlags.BlockHeaderVerificationFlags == nil {
		return nil, fmt.Errorf("block header verification flags not found")
	}
	if common.IsBitcoinChain(chainID) && !crosschainFlags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled {
		return nil, fmt.Errorf("proof verification not enabled for bitcoin chain")
	}
	if common.IsEVMChain(chainID) && !crosschainFlags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled {
		return nil, fmt.Errorf("proof verification not enabled for evm chain")
	}

	// chain must support header-based merkle proof verification
	senderChain := common.GetChainFromChainID(chainID)
	if senderChain == nil {
		return nil, types.ErrUnsupportedChain
	}
	if !senderChain.SupportMerkleProof() {
		return nil, fmt.Errorf("chain %d does not support block header-based verification", chainID)
	}

	// get block header from the store
	hashBytes, err := common.StringToHash(chainID, blockHash)
	if err != nil {
		return nil, fmt.Errorf("block hash %s conversion failed %s", blockHash, err)
	}
	res, found := k.zetaObserverKeeper.GetBlockHeader(ctx, hashBytes)
	if !found {
		return nil, fmt.Errorf("block header not found %s", blockHash)
	}

	// verify merkle proof
	txBytes, err := proof.Verify(res.Header, int(txIndex))
	if err != nil {
		return nil, err
	}

	return txBytes, err
}
