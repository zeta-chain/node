package keeper

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// Authorized: admin policy group 2.
func (k msgServer) UpdateSystemContract(goCtx context.Context, msg *types.MsgUpdateSystemContract) (*types.MsgUpdateSystemContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.observerKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_group2) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Deploy can only be executed by the correct policy account")
	}
	newSystemContractAddr := ethcommon.HexToAddress(msg.NewSystemContractAddress)
	if newSystemContractAddr == (ethcommon.Address{}) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid system contract address (%s)", msg.NewSystemContractAddress)
	}

	// update contracts
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrABIGet, "failed to get zrc20 abi")
	}
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrABIGet, "failed to get system contract abi")
	}
	foreignCoins := k.GetAllForeignCoins(ctx)
	tmpCtx, commit := ctx.CacheContext()
	for _, fcoin := range foreignCoins {
		zrc20Addr := ethcommon.HexToAddress(fcoin.Zrc20ContractAddress)
		if zrc20Addr == (ethcommon.Address{}) {
			k.Logger(ctx).Error("invalid zrc20 contract address", "address", fcoin.Zrc20ContractAddress)
			continue
		}
		_, err = k.CallEVM(tmpCtx, *zrc20ABI, types.ModuleAddressEVM, zrc20Addr, BigIntZero, nil, true, false, "updateSystemContractAddress", newSystemContractAddr)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrContractCall, "failed to call zrc20 contract method updateSystemContractAddress (%s)", err.Error())
		}
		if fcoin.CoinType == common.CoinType_Gas {
			_, err = k.CallEVM(tmpCtx, *sysABI, types.ModuleAddressEVM, newSystemContractAddr, BigIntZero, nil, true, false, "setGasCoinZRC20", big.NewInt(fcoin.ForeignChainId), zrc20Addr)
			if err != nil {
				return nil, sdkerrors.Wrapf(types.ErrContractCall, "failed to call system contract method setGasCoinZRC20 (%s)", err.Error())
			}
			_, err = k.CallEVM(tmpCtx, *sysABI, types.ModuleAddressEVM, newSystemContractAddr, BigIntZero, nil, true, false, "setGasZetaPool", big.NewInt(fcoin.ForeignChainId), zrc20Addr)
			if err != nil {
				return nil, sdkerrors.Wrapf(types.ErrContractCall, "failed to call system contract method setGasZetaPool (%s)", err.Error())
			}
		}
	}

	sys, found := k.GetSystemContract(ctx)
	if !found {
		k.Logger(ctx).Error("system contract not found")
	}
	oldSystemContractAddress := sys.SystemContract
	sys.SystemContract = newSystemContractAddr.Hex()
	k.SetSystemContract(ctx, sys)
	err = ctx.EventManager().EmitTypedEvent(
		&types.EventSystemContractUpdated{
			MsgTypeUrl:         sdk.MsgTypeURL(&types.MsgUpdateSystemContract{}),
			NewContractAddress: msg.NewSystemContractAddress,
			OldContractAddress: oldSystemContractAddress,
			Signer:             msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event", "error", err.Error())
		return nil, sdkerrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}
	commit()
	return &types.MsgUpdateSystemContractResponse{}, nil
}
