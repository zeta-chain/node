package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmcontracts "github.com/zeta-chain/zetacore/contracts/evm"
	"github.com/zeta-chain/zetacore/x/mirror/types"
	"math/big"
)

func (k msgServer) DepoistERC20(goCtx context.Context, msg *types.MsgDepoistERC20) (*types.MsgDepoistERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenPairs, found := k.GetERC20TokenPairs(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrTokenPairNotFound, "token pair not found")
	}

	var pair *types.ERC20TokenPair
	for _, tokenPair := range tokenPairs.TokenPairs {
		if tokenPair.HomeERC20ContractAddress == msg.HomeERC20ContractAddress {
			found = true
			pair = tokenPair
			break
		}
	}

	if !found {
		return nil, sdkerrors.Wrap(types.ErrTokenPairNotFound, "token pair not found")
	}

	erc20 := evmcontracts.ERC20MinterBurnerDecimalsContract.ABI
	contract := ethcommon.HexToAddress(pair.MirrorERC20ContractAddress)
	receiver := ethcommon.HexToAddress(msg.RecipientAddress)
	balanceToken := k.BalanceOf(ctx, erc20, contract, receiver)
	amount, ok := big.NewInt(0).SetString(msg.Amount, 10)

	if !ok {
		return nil, sdkerrors.Wrap(types.ErrInvalidAmount, "invalid amount")
	}

	// Mint tokens and send to receiver
	_, err := k.CallEVM(ctx, erc20, types.ModuleAddress, contract, true, "mint", receiver, amount)
	if err != nil {
		return nil, err
	}

	// Check expected receiver balance after transfer
	tokens := amount
	balanceTokenAfter := k.BalanceOf(ctx, erc20, contract, receiver)
	if balanceTokenAfter == nil {
		return nil, sdkerrors.Wrap(types.ErrEVMCall, "failed to retrieve balance")
	}
	expToken := big.NewInt(0).Add(balanceToken, tokens)

	if r := balanceTokenAfter.Cmp(expToken); r != 0 {
		return nil, sdkerrors.Wrapf(
			types.ErrBalanceInvariance,
			"invalid token balance - expected: %v, actual: %v", expToken, balanceTokenAfter,
		)
	}

	return &types.MsgDepoistERC20Response{}, nil
}
