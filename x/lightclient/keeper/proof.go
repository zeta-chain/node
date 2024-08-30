package keeper

import (
	cosmoserror "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/proofs"
	"github.com/zeta-chain/node/x/lightclient/types"
)

// VerifyProof verifies the merkle proof for a given chain and block header
// It returns the transaction bytes if the proof is valid
func (k Keeper) VerifyProof(
	ctx sdk.Context,
	proof *proofs.Proof,
	chainID int64,
	blockHash string,
	txIndex int64,
) ([]byte, error) {
	// check block header verification is set
	if err := k.CheckBlockHeaderVerificationEnabled(ctx, chainID); err != nil {
		return nil, err
	}

	// additionalChains is a list of additional chains to search from
	// it is used in the protocol to dynamically support new chains without doing an upgrade
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)

	// get block header from the store
	hashBytes, err := chains.StringToHash(chainID, blockHash, additionalChains)
	if err != nil {
		return nil, cosmoserror.Wrapf(
			types.ErrInvalidBlockHash,
			"block hash %s conversion failed %s",
			blockHash,
			err.Error(),
		)
	}
	res, found := k.GetBlockHeader(ctx, hashBytes)
	if !found {
		return nil, cosmoserror.Wrapf(types.ErrBlockHeaderNotFound, "block header not found %s", blockHash)
	}

	// verify merkle proof
	txBytes, err := proof.Verify(res.Header, int(txIndex))
	if err != nil {
		return nil, cosmoserror.Wrapf(
			types.ErrProofVerificationFailed,
			"failed to verify merkle proof: %s",
			err.Error(),
		)
	}
	return txBytes, nil
}
