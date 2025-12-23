package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// InboundDetails holds the details required to process an inbound transaction
// It is used internally to encapsulate the information needed for processing withdrawals and calls
type InboundDetails struct {
	coinType        coin.CoinType
	asset           string
	receiverChain   chains.Chain
	gasLimitQueried *big.Int
}

// ProcessZEVMInboundV2 processes the logs emitted by the zEVM contract for V2 protocol contracts
// it parses logs from GatewayZEVM contract and updates the crosschain module state
func (k Keeper) ProcessZEVMInboundV2(
	ctx sdk.Context,
	log *ethtypes.Log,
	gatewayAddr ethcommon.Address,
	txOrigin string,
) error {
	// try to parse a withdrawal event from the log
	withdrawalEvent, callEvent, withdrawalAndCallEvent, err := types.ParseGatewayEvent(*log, gatewayAddr)
	if err == nil && (withdrawalEvent != nil || callEvent != nil || withdrawalAndCallEvent != nil) {
		var inbound *types.MsgVoteInbound

		// parse data from event and validate
		var zrc20 ethcommon.Address
		var value *big.Int
		var receiver []byte
		if withdrawalEvent != nil {
			zrc20 = withdrawalEvent.Zrc20
			value = withdrawalEvent.Value
			receiver = withdrawalEvent.Receiver
		} else if callEvent != nil {
			zrc20 = callEvent.Zrc20
			value = big.NewInt(0)
			receiver = callEvent.Receiver
		} else {
			zrc20 = withdrawalAndCallEvent.Zrc20
			value = withdrawalAndCallEvent.Value
			receiver = withdrawalAndCallEvent.Receiver
		}

		wzetaContractAddress, err := k.fungibleKeeper.GetWZetaContractAddress(ctx)
		if err != nil {
			return errorsmod.Wrapf(
				types.ErrCannotProcessWithdrawal,
				"failed to get WZeta contract address: %v", err,
			)
		}

		var inboundDetails InboundDetails
		// The following condition checks if the withdrawal is for ZETA or ZRC20.
		// ZRC20 cointype can be further classified as ZRC20 or Gas , based on foreign coin or NoAssetCall.
		// Note: NoAssetCall is not supported for ZETA
		switch {
		case zrc20 == wzetaContractAddress:
			if !k.zetaObserverKeeper.IsV2ZetaEnabled(ctx) {
				return types.ErrZetaThroughGateway
			}
			var receiverChainID *big.Int
			var callOptions gatewayzevm.CallOptions
			if withdrawalEvent != nil {
				receiverChainID = withdrawalEvent.ChainId
				callOptions = withdrawalEvent.CallOptions
			} else if withdrawalAndCallEvent != nil {
				receiverChainID = withdrawalAndCallEvent.ChainId
				callOptions = withdrawalAndCallEvent.CallOptions
			} else {
				return errorsmod.Wrap(types.ErrInvalidWithdrawalEvent, "ZETA withdrawal requires withdrawal event")
			}
			inboundDetails, err = k.getZETAInboundDetails(ctx, receiverChainID, callOptions)
		default:
			inboundDetails, err = k.getZRC20InboundDetails(ctx, zrc20, callEvent != nil)
		}

		if err != nil {
			return errorsmod.Wrapf(
				types.ErrInvalidWithdrawalEvent,
				"failed to parse inbound details for withdraw: %v", err,
			)
		}

		// if the coin type is ZETA, we need to burn the coins as the GatewayZEVM contract transfers the ZETA to the fungible module account.
		// For ERC20 and GAS coin types; the GatewayZEVM directly burns the tokens, so we can skip this step.
		// NOTE: value does not include tokens paid for the gas which is burned by the GatewayZEVM contract for both cases
		if inboundDetails.coinType == coin.CoinType_Zeta {
			err := k.bankKeeper.BurnCoins(ctx,
				fungibletypes.ModuleName,
				sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(value))))
			if err != nil {
				return errorsmod.Wrapf(
					types.ErrInvalidWithdrawalEvent,
					"failed to burn ZETA coins: %v", err,
				)
			}
		}

		// validate data of the withdrawal event
		if err := k.validateOutbound(ctx, inboundDetails.receiverChain.ChainId, inboundDetails.coinType, value, receiver); err != nil {
			return errorsmod.Wrapf(
				types.ErrInvalidWithdrawalEvent,
				"failed to validate withdrawal event: %v", err,
			)
		}

		// create inbound object depending on the event type
		if withdrawalEvent != nil {
			inbound, err = types.NewWithdrawalInbound(
				ctx,
				txOrigin,
				inboundDetails.coinType,
				inboundDetails.asset,
				withdrawalEvent,
				inboundDetails.receiverChain,
				inboundDetails.gasLimitQueried,
			)
			if err != nil {
				return err
			}
		} else if callEvent != nil {
			inbound, err = types.NewCallInbound(
				ctx,
				txOrigin,
				callEvent,
				inboundDetails.receiverChain,
				inboundDetails.gasLimitQueried,
			)
			if err != nil {
				return err
			}
		} else {
			inbound, err = types.NewWithdrawAndCallInbound(
				ctx,
				txOrigin,
				inboundDetails.coinType,
				inboundDetails.asset,
				withdrawalAndCallEvent,
				inboundDetails.receiverChain,
				inboundDetails.gasLimitQueried,
			)
			if err != nil {
				return err
			}
		}

		if inbound == nil {
			return errors.New("ParseGatewayEvent: invalid log - no event found")
		}
		// validate inbound for processing
		// V2 inbounds always pay gas directly at the contract call
		cctx, err := k.ValidateInbound(ctx, inbound, false)
		if err != nil {
			return err
		}
		if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
			return errors.New("cctx aborted")
		}

		EmitZRCWithdrawCreated(ctx, *cctx)
	}
	return nil
}

// getZETAInboundDetails retrieves the details for a ZETA withdrawal event, it returns an InboundDetails object
func (k Keeper) getZETAInboundDetails(
	ctx sdk.Context,
	receiverChainID *big.Int,
	callOptions gatewayzevm.CallOptions,
) (InboundDetails, error) {
	if receiverChainID == nil || receiverChainID.Int64() == 0 {
		return InboundDetails{}, errorsmod.Wrap(
			types.ErrInvalidWithdrawalEvent,
			"receiver chain ID is nil or zero for ZETA withdrawal",
		)
	}
	parsedReceiverChain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(
		ctx,
		receiverChainID.Int64(),
	)
	if !found {
		return InboundDetails{}, errorsmod.Wrapf(
			observertypes.ErrSupportedChains,
			"chain with chainID %d not supported",
			receiverChainID.Int64(),
		)
	}
	// Validation if we want to send ZETA to an external chain, but there is no ZETA token.
	chainParams, found := k.zetaObserverKeeper.GetChainParamsByChainID(ctx, parsedReceiverChain.ChainId)
	if !found {
		return InboundDetails{}, errorsmod.Wrapf(
			observertypes.ErrChainParamsNotFound,
			"chaind ID :%d",
			parsedReceiverChain.ChainId,
		)
	}

	if parsedReceiverChain.IsExternalChain() &&
		(chainParams.ZetaTokenContractAddress == "" || chainParams.ZetaTokenContractAddress == constant.EVMZeroAddress) {
		return InboundDetails{}, errorsmod.Wrapf(
			types.ErrUnableToSendCoinType,
			" cannot send ZETA to external chain %d",
			parsedReceiverChain.ChainId,
		)
	}

	gasLimit := callOptions.GasLimit
	if gasLimit == nil || gasLimit.Int64() == 0 {
		return InboundDetails{}, errorsmod.Wrap(
			types.ErrInvalidWithdrawalEvent, "gas limit not provided for ZETA withdrawal")
	}

	return InboundDetails{
		coinType:        coin.CoinType_Zeta,
		asset:           ethcommon.Address{}.Hex(),
		receiverChain:   parsedReceiverChain,
		gasLimitQueried: gasLimit,
	}, nil
}

// getZRC20InboundDetails retrieves the details for a ZRC20 withdrawal event, it returns an InboundDetails object
func (k Keeper) getZRC20InboundDetails(
	ctx sdk.Context,
	zrc20 ethcommon.Address,
	callEvent bool,
) (InboundDetails, error) {
	foreignCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, zrc20.Hex())
	if !found {
		ctx.Logger().Info(fmt.Sprintf("cannot find foreign coin associated to the zrc20 address %s", zrc20.Hex()))
		return InboundDetails{}, nil
	}

	receiverChain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, foreignCoin.ForeignChainId)
	if !found {
		return InboundDetails{}, errorsmod.Wrapf(
			observertypes.ErrSupportedChains,
			"chain with chainID %d not supported",
			foreignCoin.ForeignChainId,
		)
	}

	gasLimitQueried, err := k.fungibleKeeper.QueryGasLimit(
		ctx,
		ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress),
	)
	if err != nil {
		return InboundDetails{}, err
	}

	coinType := foreignCoin.CoinType
	if callEvent {
		coinType = coin.CoinType_NoAssetCall
	}

	return InboundDetails{
		receiverChain:   receiverChain,
		gasLimitQueried: gasLimitQueried,
		coinType:        coinType,
		asset:           foreignCoin.Asset,
	}, nil
}
