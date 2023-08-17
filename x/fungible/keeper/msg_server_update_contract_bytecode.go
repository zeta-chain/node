package keeper

import (
	"bytes"
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
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
		return nil, sdkerrors.Wrapf(types.ErrContractNotFound, "contract (%s) not found", contractAddress.String())
	}
	oldCodeHash := acct.CodeHash

	newByteCode := msg.NewBytecode
	if len(newByteCode) == 0 { // empty bytecode will disable the contract; turning it into an EOA?
		acct.CodeHash = evmtypes.EmptyCodeHash
	}

	newCodeHash := crypto.Keccak256(newByteCode)
	if bytes.Compare(oldCodeHash, newCodeHash) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrDeployContract, "contract (%s) bytecode not changed; code hash %x", contractAddress.String(), oldCodeHash)
	}

	k.evmKeeper.SetCode(ctx, newCodeHash, newByteCode)
	acct.CodeHash = newCodeHash
	err := k.evmKeeper.SetAccount(ctx, contractAddress, *acct)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDeployContract, "failed to update contract (%s) bytecode; code hash %x", contractAddress.String(), oldCodeHash)
	}

	//err = ctx.EventManager().EmitTypedEvent(
	//	&types.EventSystemContractUpdated{
	//		MsgTypeUrl:         sdk.MsgTypeURL(&types.MsgUpdateSystemContract{}),
	//		NewContractAddress: msg.NewSystemContractAddress,
	//		OldContractAddress: oldSystemContractAddress,
	//		Signer:             msg.Creator,
	//	},
	//)
	//if err != nil {
	//	k.Logger(ctx).Error("failed to emit event", "error", err.Error())
	//	return nil, sdkerrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	//}
	//commit()
	return &types.MsgUpdateContractBytecodeResponse{}, nil
}
