package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

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
		var contractAddress ethcommon.Address
		var receiverChainID *big.Int
		var callOptions gatewayzevm.CallOptions
		if withdrawalEvent != nil {
			zrc20 = withdrawalEvent.Zrc20
			value = withdrawalEvent.Value
			receiver = withdrawalEvent.Receiver
			contractAddress = withdrawalEvent.Raw.Address
			receiverChainID = withdrawalEvent.ChainId
			callOptions = withdrawalEvent.CallOptions
		} else if callEvent != nil {
			zrc20 = callEvent.Zrc20
			value = big.NewInt(0)
			receiver = callEvent.Receiver
			contractAddress = callEvent.Raw.Address
			callOptions = callEvent.CallOptions
			receiverChainID = big.NewInt(0) // Receiver chain ID is only used for withdraws when coin type is ZETA.
		} else {
			zrc20 = withdrawalAndCallEvent.Zrc20
			value = withdrawalAndCallEvent.Value
			receiver = withdrawalAndCallEvent.Receiver
			contractAddress = withdrawalAndCallEvent.Raw.Address
			receiverChainID = withdrawalAndCallEvent.ChainId
			callOptions = withdrawalAndCallEvent.CallOptions
		}

		wzeta, err := k.fungibleKeeper.GetWZetaContractAddress(ctx)
		if err != nil {
			fmt.Println("Failed to get WZeta contract address:", err)
		}

		coinType := coin.CoinType_ERC20
		receiverChain := chains.Chain{}
		asset := ""
		gasLimitQueried := big.NewInt(0)
		foundChain := true

		if callEvent != nil {
			coinType = coin.CoinType_NoAssetCall
		} else {
			switch {
			case zrc20 == wzeta:
				coinType = coin.CoinType_Zeta
			default:
				coinType = coin.CoinType_ERC20
			}
		}

		if coinType == coin.CoinType_Zeta {
			if receiverChainID == nil || receiverChainID.Int64() == 0 {
				return errorsmod.Wrapf(types.ErrInvalidWithdrawalEvent, " receiver chain ID is nil or zero for ZETA withdrawal")
			}

			// Fetch details from event data directly
			receiverChain, foundChain = k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, receiverChainID.Int64())
			if !foundChain {
				fmt.Println("Cannot find supported chain for receiver chain ID:", receiverChainID.Int64())
			}
			asset = wzeta.String()
			if callOptions.GasLimit.Int64() == 0 {
				callOptions.GasLimit = big.NewInt(100_000) // default gas limit for WZeta
			}
			gasLimitQueried = callOptions.GasLimit
		}

		if coinType == coin.CoinType_ERC20 || coinType == coin.CoinType_NoAssetCall {
			foreignCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, zrc20.Hex())
			if !found {
				ctx.Logger().
					Info(fmt.Sprintf("cannot find foreign coin with contract address %s", contractAddress.Hex()))
				return nil
			}
			receiverChain, found = k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, foreignCoin.ForeignChainId)
			if !found {
				return errorsmod.Wrapf(
					observertypes.ErrSupportedChains,
					"chain with chainID %d not supported",
					foreignCoin.ForeignChainId,
				)
			}

			// Fetch details from foreign coin
			gasLimitQueried, err = k.fungibleKeeper.QueryGasLimit(
				ctx,
				ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress),
			)
			if err != nil {
				return err
			}
			coinType = foreignCoin.CoinType
			asset = foreignCoin.Asset
		}

		//// validate data of the withdrawal event
		//if callEvent != nil {
		//	coinType = coin.CoinType_NoAssetCall
		//}
		if err := k.validateOutbound(ctx, receiverChain.ChainId, coinType, value, receiver); err != nil {
			return err
		}

		// create inbound object depending on the event type
		if withdrawalEvent != nil {
			inbound, err = types.NewWithdrawalInbound(
				ctx,
				txOrigin,
				coinType,
				asset,
				withdrawalEvent,
				receiverChain,
				gasLimitQueried,
			)
			if err != nil {
				return err
			}
		} else if callEvent != nil {
			inbound, err = types.NewCallInbound(
				ctx,
				txOrigin,
				callEvent,
				receiverChain,
				gasLimitQueried,
			)
			if err != nil {
				return err
			}
		} else {
			inbound, err = types.NewWithdrawAndCallInbound(
				ctx,
				txOrigin,
				coinType,
				asset,
				withdrawalAndCallEvent,
				receiverChain,
				gasLimitQueried,
			)
			if err != nil {
				return err
			}
		}

		if inbound == nil {
			return errors.New("ParseGatewayEvent: invalid log - no event found")
		}

		// validate inbound for processing
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
