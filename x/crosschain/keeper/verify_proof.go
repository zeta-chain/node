package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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

func (k Keeper) VerifyEVMInTxBody(ctx sdk.Context, msg *types.MsgAddToInTxTracker, txBytes []byte) error {
	var txx ethtypes.Transaction
	err := txx.UnmarshalBinary(txBytes)
	if err != nil {
		return err
	}
	switch msg.CoinType {
	case common.CoinType_Zeta:
		coreParams, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, msg.ChainId)
		if !found {
			return types.ErrUnsupportedChain.Wrapf("core params not found for chain %d", msg.ChainId)
		}
		if txx.To().Hex() != coreParams.ConnectorContractAddress {
			return fmt.Errorf("receiver is not connector contract for coin type %s", msg.CoinType)
		}
		return nil
	case common.CoinType_ERC20:
		coreParams, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, msg.ChainId)
		if !found {
			return types.ErrUnsupportedChain.Wrapf("core params not found for chain %d", msg.ChainId)
		}
		if txx.To().Hex() != coreParams.Erc20CustodyContractAddress {
			return fmt.Errorf("receiver is not erc20Custory contract for coin type %s", msg.CoinType)
		}
		return nil
	case common.CoinType_Gas:
		tss, err := k.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
		if err != nil {
			return err
		}
		tssAddr := eth.HexToAddress(tss.Eth)
		if tssAddr == (eth.Address{}) {
			return fmt.Errorf("tss address not found")
		}
		if txx.To().Hex() != tssAddr.Hex() {
			return fmt.Errorf("receiver is not tssAddress contract for coin type %s", msg.CoinType)
		}
		return nil
	default:
		return fmt.Errorf("coin type %s not supported", msg.CoinType)
	}
}
