package keeper

import (
	"context"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"

	"github.com/zeta-chain/node/pkg/coin"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// UpdateSystemContract updates the system contract
func (k msgServer) UpdateSystemContract(
	goCtx context.Context,
	msg *types.MsgUpdateSystemContract,
) (*types.MsgUpdateSystemContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}
	newSystemContractAddr := ethcommon.HexToAddress(msg.NewSystemContractAddress)
	if newSystemContractAddr == (ethcommon.Address{}) {
		return nil, cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid system contract address (%s)",
			msg.NewSystemContractAddress,
		)
	}

	// update contracts
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrABIGet, "failed to get zrc20 abi")
	}
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrABIGet, "failed to get system contract abi")
	}
	foreignCoins := k.GetAllForeignCoins(ctx)
	tmpCtx, commit := ctx.CacheContext()
	for _, fcoin := range foreignCoins {
		zrc20Addr := ethcommon.HexToAddress(fcoin.Zrc20ContractAddress)
		if zrc20Addr == (ethcommon.Address{}) {
			k.Logger(ctx).Error("invalid zrc20 contract address", "address", fcoin.Zrc20ContractAddress)
			continue
		}
		_, err = k.CallEVM(
			tmpCtx,
			*zrc20ABI,
			types.ModuleAddressEVM,
			zrc20Addr,
			BigIntZero,
			DefaultGasLimit,
			true,
			false,
			"updateSystemContractAddress",
			newSystemContractAddr,
		)
		if err != nil {
			return nil, cosmoserrors.Wrapf(
				types.ErrContractCall,
				"failed to call zrc20 contract method updateSystemContractAddress (%s)",
				err.Error(),
			)
		}
		if fcoin.CoinType == coin.CoinType_Gas {
			_, err = k.CallEVM(
				tmpCtx,
				*sysABI,
				types.ModuleAddressEVM,
				newSystemContractAddr,
				BigIntZero,
				DefaultGasLimit,
				true,
				false,
				"setGasCoinZRC20",
				big.NewInt(fcoin.ForeignChainId),
				zrc20Addr,
			)
			if err != nil {
				return nil, cosmoserrors.Wrapf(
					types.ErrContractCall,
					"failed to call system contract method setGasCoinZRC20 (%s)",
					err.Error(),
				)
			}
			_, err = k.CallEVM(
				tmpCtx,
				*sysABI,
				types.ModuleAddressEVM,
				newSystemContractAddr,
				BigIntZero,
				DefaultGasLimit,
				true,
				false,
				"setGasZetaPool",
				big.NewInt(fcoin.ForeignChainId),
				zrc20Addr,
			)
			if err != nil {
				return nil, cosmoserrors.Wrapf(
					types.ErrContractCall,
					"failed to call system contract method setGasZetaPool (%s)",
					err.Error(),
				)
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
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}
	commit()
	return &types.MsgUpdateSystemContractResponse{}, nil
}
