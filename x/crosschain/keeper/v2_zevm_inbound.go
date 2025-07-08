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
	fmt.Println("Processing zEVM inbound V2 log:")

	// Parse gateway event
	withdrawalEvent, callEvent, withdrawalAndCallEvent, err := types.ParseGatewayEvent(*log, gatewayAddr)
	if err != nil || (withdrawalEvent == nil && callEvent == nil && withdrawalAndCallEvent == nil) {
		return err
	}

	// Extract common data from events
	eventData := k.extractEventData(withdrawalEvent, callEvent, withdrawalAndCallEvent)

	// Get WZeta contract address
	wzeta, err := k.fungibleKeeper.GetWZetaContractAddress(ctx)
	if err != nil {
		fmt.Println("Failed to get WZeta contract address:", err)
	}

	// Process based on coin type
	var coinType coin.CoinType
	var receiverChain chains.Chain
	var asset string
	var gasLimitQueried *big.Int
	var found bool

	if eventData.zrc20 == wzeta {
		coinType, receiverChain, asset, gasLimitQueried, found = k.processZetaCoin(ctx, eventData, wzeta)
		if !found {
			fmt.Println("Cannot find supported chain for receiver chain ID:", eventData.receiverChainID.Int64())
		}
	} else {
		coinType, receiverChain, asset, gasLimitQueried, err = k.processZRC20Coin(ctx, eventData)
		if err != nil {
			return err
		}
	}

	// Handle call event (NoAssetCall)
	if callEvent != nil {
		coinType = coin.CoinType_NoAssetCall
	}

	// Validate outbound
	if err := k.validateOutbound(ctx, receiverChain.ChainId, coinType, eventData.value, eventData.receiver); err != nil {
		return err
	}

	// Create inbound based on event type
	inbound, err := k.createInbound(ctx, txOrigin, coinType, asset, withdrawalEvent, callEvent, withdrawalAndCallEvent, receiverChain, gasLimitQueried)
	if err != nil {
		return err
	}

	if inbound == nil {
		return errors.New("ParseGatewayEvent: invalid log - no event found")
	}

	// Validate and process inbound
	cctx, err := k.ValidateInbound(ctx, inbound, false)
	if err != nil {
		return err
	}
	if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
		return errors.New("cctx aborted")
	}

	EmitZRCWithdrawCreated(ctx, *cctx)
	return nil
}

// EventData holds common data extracted from gateway events
type EventData struct {
	zrc20           ethcommon.Address
	value           *big.Int
	receiver        []byte
	contractAddress ethcommon.Address
	receiverChainID *big.Int
	callOptions     gatewayzevm.CallOptions
}

// extractEventData extracts common data from different event types
func (k Keeper) extractEventData(withdrawalEvent *gatewayzevm.GatewayZEVMWithdrawn, callEvent *gatewayzevm.GatewayZEVMCalled, withdrawalAndCallEvent *gatewayzevm.GatewayZEVMWithdrawnAndCalled) EventData {
	var eventData EventData

	if withdrawalEvent != nil {
		eventData.zrc20 = withdrawalEvent.Zrc20
		eventData.value = withdrawalEvent.Value
		eventData.receiver = withdrawalEvent.Receiver
		eventData.contractAddress = withdrawalEvent.Raw.Address
		eventData.receiverChainID = withdrawalEvent.ChainId
		eventData.callOptions = withdrawalEvent.CallOptions
	} else if callEvent != nil {
		eventData.zrc20 = callEvent.Zrc20
		eventData.value = big.NewInt(0)
		eventData.receiver = callEvent.Receiver
		eventData.contractAddress = callEvent.Raw.Address
		eventData.callOptions = callEvent.CallOptions
	} else {
		eventData.zrc20 = withdrawalAndCallEvent.Zrc20
		eventData.value = withdrawalAndCallEvent.Value
		eventData.receiver = withdrawalAndCallEvent.Receiver
		eventData.contractAddress = withdrawalAndCallEvent.Raw.Address
		eventData.receiverChainID = withdrawalAndCallEvent.ChainId
		eventData.callOptions = withdrawalAndCallEvent.CallOptions
	}

	return eventData
}

// processZetaCoin handles processing for Zeta coin type
func (k Keeper) processZetaCoin(ctx sdk.Context, eventData EventData, wzeta ethcommon.Address) (coin.CoinType, chains.Chain, string, *big.Int, bool) {
	coinType := coin.CoinType_Zeta
	receiverChain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, eventData.receiverChainID.Int64())
	asset := wzeta.String()

	gasLimitQueried := eventData.callOptions.GasLimit
	if gasLimitQueried == big.NewInt(0) {
		gasLimitQueried = big.NewInt(100_000) // default gas limit for WZeta withdrawals
	}

	return coinType, receiverChain, asset, gasLimitQueried, found
}

// processZRC20Coin handles processing for ZRC20 coin type
func (k Keeper) processZRC20Coin(ctx sdk.Context, eventData EventData) (coin.CoinType, chains.Chain, string, *big.Int, error) {
	foreignCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, eventData.zrc20.Hex())
	if !found {
		ctx.Logger().
			Info(fmt.Sprintf("cannot find foreign coin with contract address %s", eventData.contractAddress.Hex()))
		return coin.CoinType_ERC20, chains.Chain{}, "", big.NewInt(0), nil
	}

	receiverChain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, foreignCoin.ForeignChainId)
	if !found {
		return coin.CoinType_ERC20, chains.Chain{}, "", big.NewInt(0), errorsmod.Wrapf(
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
		return coin.CoinType_ERC20, chains.Chain{}, "", big.NewInt(0), err
	}

	return foreignCoin.CoinType, receiverChain, foreignCoin.Asset, gasLimitQueried, nil
}

// createInbound creates the appropriate inbound object based on event type
func (k Keeper) createInbound(
	ctx sdk.Context,
	txOrigin string,
	coinType coin.CoinType,
	asset string,
	withdrawalEvent *gatewayzevm.GatewayZEVMWithdrawn,
	callEvent *gatewayzevm.GatewayZEVMCalled,
	withdrawalAndCallEvent *gatewayzevm.GatewayZEVMWithdrawnAndCalled,
	receiverChain chains.Chain,
	gasLimitQueried *big.Int,
) (*types.MsgVoteInbound, error) {
	if withdrawalEvent != nil {
		return types.NewWithdrawalInbound(
			ctx,
			txOrigin,
			coinType,
			asset,
			withdrawalEvent,
			receiverChain,
			gasLimitQueried,
		)
	} else if callEvent != nil {
		return types.NewCallInbound(
			ctx,
			txOrigin,
			callEvent,
			receiverChain,
			gasLimitQueried,
		)
	} else {
		return types.NewWithdrawAndCallInbound(
			ctx,
			txOrigin,
			coinType,
			asset,
			withdrawalAndCallEvent,
			receiverChain,
			gasLimitQueried,
		)
	}
}
