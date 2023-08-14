package keeper

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) UpdateZRC20WithdrawFee(goCtx context.Context, msg *types.MsgUpdateZRC20WithdrawFee) (*types.MsgUpdateZRC20WithdrawFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_deploy_fungible_coin) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Deploy can only be executed by the correct policy account")
	}
	zrc20Addr := ethcommon.HexToAddress(msg.Zrc20Address)
	if zrc20Addr == (ethcommon.Address{}) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid zrc20 contract address (%s)", msg.Zrc20Address)
	}

	// update contracts
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrABIGet, "failed to get zrc20 abi")
	}

	foreignCoins := k.GetAllForeignCoins(ctx)
	found := false
	var coin types.ForeignCoins
	for _, fcoin := range foreignCoins {
		coinZRC20Addr := ethcommon.HexToAddress(fcoin.Zrc20ContractAddress)
		if coinZRC20Addr == (ethcommon.Address{}) {
			k.Logger(ctx).Error("invalid zrc20 contract address", "address", fcoin.Zrc20ContractAddress)
			continue
		}
		if coinZRC20Addr == zrc20Addr {
			coin = fcoin
			found = true
			break
		}
	}

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAddress, "no foreign coin match requested zrc20 address (%s)", msg.Zrc20Address)
	}

	res, err := k.CallEVM(ctx, *zrc20ABI, types.ModuleAddressEVM, zrc20Addr, BigIntZero, nil, false, "PROTOCOL_FLAT_FEE")
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrContractCall, "failed to call zrc20 contract PROTOCOL_FLAT_FEE method (%s)", err.Error())
	}
	unpacked, err := zrc20ABI.Unpack("PROTOCOL_FLAT_FEE", res.Ret)
	if err != nil || len(unpacked) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrContractCall, "failed to unpack zrc20 contract PROTOCOL_FLAT_FEE method (%s)", err.Error())
	}
	oldWithdrawFee, ok := unpacked[0].(*big.Int)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrContractCall, "failed to interpret the returned unpacked zrc20 contract PROTOCOL_FLAT_FEE method; ret %x", res.Ret)
	}

	tmpCtx, commit := ctx.CacheContext()
	_, err = k.CallEVM(tmpCtx, *zrc20ABI, types.ModuleAddressEVM, zrc20Addr, BigIntZero, nil, true, "updateProtocolFlatFee", msg.NewWithdrawFee.BigInt())

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventZRC20WithdrawFeeUpdated{
			MsgTypeUrl:     sdk.MsgTypeURL(&types.MsgUpdateZRC20WithdrawFee{}),
			ChainId:        coin.ForeignChainId,
			CoinType:       coin.CoinType,
			Zrc20Address:   zrc20Addr.Hex(),
			OldWithdrawFee: oldWithdrawFee.String(),
			NewWithdrawFee: msg.NewWithdrawFee.BigInt().String(),
			Signer:         msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event", "error", err.Error())
		return nil, sdkerrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}
	commit()
	return &types.MsgUpdateZRC20WithdrawFeeResponse{}, nil
}
