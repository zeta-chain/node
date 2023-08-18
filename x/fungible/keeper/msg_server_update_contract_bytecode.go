package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) UpdateContractBytecode(goCtx context.Context, msg *types.MsgUpdateContractBytecode) (*types.MsgUpdateContractBytecodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_deploy_fungible_coin) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Deploy can only be executed by the correct policy account")
	}
	contractAddress := ethcommon.HexToAddress(msg.ContractAddress)
	if contractAddress == (ethcommon.Address{}) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", msg.ContractAddress)
	}

	acct := k.evmKeeper.GetAccount(ctx, contractAddress)
	if acct == nil {
		return nil, sdkerrors.Wrapf(types.ErrContractNotFound, "contract (%s) not found", contractAddress.Hex())
	}

	oldCodeHash := ethcommon.BytesToHash(acct.CodeHash)
	oldBytecode := k.evmKeeper.GetCode(ctx, oldCodeHash)
	newCodeHash := crypto.Keccak256(msg.NewBytecode)
	newBytecode := msg.NewBytecode
	k.evmKeeper.SetCode(ctx, newCodeHash, newBytecode)
	acct.CodeHash = newCodeHash
	err := k.evmKeeper.SetAccount(ctx, contractAddress, *acct)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrContractNotFound, "failed to update contract (%s) bytecode (%s)", contractAddress.Hex(), err.Error())
	}
	k.Logger(ctx).Info("updated contract bytecode", "contract", contractAddress.Hex(), "oldBytecode", len(oldBytecode), "newBytecode", len(msg.NewBytecode))

	return &types.MsgUpdateContractBytecodeResponse{}, nil
}
