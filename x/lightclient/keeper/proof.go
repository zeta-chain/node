package keeper

import (
	cosmoserror "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// VerifyProof verifies the merkle proof for a given chain and block header
// It returns the transaction bytes if the proof is valid
func (k Keeper) VerifyProof(ctx sdk.Context, proof *proofs.Proof, chainID int64, blockHash string, txIndex int64) ([]byte, error) {
	// check verification flags are set
	if err := k.CheckVerificationFlagsEnabled(ctx, chainID); err != nil {
		return nil, err
	}

	// chain must support header-based merkle proof verification
	senderChain := chains.GetChainFromChainID(chainID)
	if senderChain == nil {
		return nil, cosmoserror.Wrapf(types.ErrChainNotSupported, "chain id %d doesn't exist", chainID)
	}
	if !senderChain.SupportMerkleProof() {
		return nil, cosmoserror.Wrapf(types.ErrChainNotSupported, "chain id %d doesn't support merkle proof", chainID)
	}

	// get block header from the store
	hashBytes, err := chains.StringToHash(chainID, blockHash)
	if err != nil {
		return nil, cosmoserror.Wrapf(types.ErrInvalidBlockHash, "block hash %s conversion failed %s", blockHash, err.Error())
	}
	res, found := k.GetBlockHeader(ctx, hashBytes)
	if !found {
		return nil, cosmoserror.Wrapf(types.ErrBlockHeaderNotFound, "block header not found %s", blockHash)
	}

	// verify merkle proof
	txBytes, err := proof.Verify(res.Header, int(txIndex))
	if err != nil {
		return nil, cosmoserror.Wrapf(types.ErrProofVerificationFailed, "failed to verify merkle proof: %s", err.Error())
	}
	return txBytes, nil
}
