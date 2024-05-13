package keeper

import (
	"context"

	cosmoserror "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// UpdateContractBytecode updates the bytecode of a contract from the bytecode
// of an existing contract Only a ZRC20 contract or the WZeta connector contract
// can be updated IMPORTANT: the new contract bytecode must have the same
// storage layout as the old contract bytecode the new contract can add new
// variable but cannot remove any existing variable
//
// Authozied: admin policy group 2
func (k msgServer) UpdateContractBytecode(goCtx context.Context, msg *types.MsgUpdateContractBytecode) (*types.MsgUpdateContractBytecodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check authorization
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserror.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// fetch account to update
	if !ethcommon.IsHexAddress(msg.ContractAddress) {
		return nil, cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", msg.ContractAddress)
	}
	contractAddress := ethcommon.HexToAddress(msg.ContractAddress)
	acct := k.evmKeeper.GetAccount(ctx, contractAddress)
	if acct == nil {
		return nil, cosmoserror.Wrapf(types.ErrContractNotFound, "contract (%s) not found", contractAddress.Hex())
	}

	// check the contract is a zrc20
	_, found := k.GetForeignCoins(ctx, msg.ContractAddress)
	if !found {
		// check contract is wzeta connector contract
		systemContract, found := k.GetSystemContract(ctx)
		if !found {
			return nil, types.ErrSystemContractNotFound
		}
		if msg.ContractAddress != systemContract.ConnectorZevm {
			// not a zrc20 or wzeta connector contract, can't be updated
			return nil, cosmoserror.Wrapf(types.ErrInvalidContract, "contract (%s) is neither a zrc20 nor wzeta connector", msg.ContractAddress)
		}
	}

	// set the new CodeHash to the account
	oldCodeHash := acct.CodeHash
	acct.CodeHash = ethcommon.HexToHash(msg.NewCodeHash).Bytes()
	err = k.evmKeeper.SetAccount(ctx, contractAddress, *acct)
	if err != nil {
		return nil, cosmoserror.Wrapf(
			types.ErrSetBytecode,
			"failed to update contract (%s) bytecode (%s)",
			contractAddress.Hex(),
			err.Error(),
		)
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventBytecodeUpdated{
			MsgTypeUrl:      sdk.MsgTypeURL(&types.MsgUpdateContractBytecode{}),
			ContractAddress: msg.ContractAddress,
			OldBytecodeHash: ethcommon.BytesToHash(oldCodeHash).Hex(),
			NewBytecodeHash: msg.NewCodeHash,
			Signer:          msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event", "error", err.Error())
		return nil, cosmoserror.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgUpdateContractBytecodeResponse{}, nil
}
