package keeper

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zetaconnectorzevm.sol"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

var _ evmtypes.EvmHooks = Hooks{}

type Hooks struct {
	k Keeper
}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing is a wrapper for calling the EVM PostTxProcessing hook on
// the module keeper
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	return h.k.PostTxProcessing(ctx, msg, receipt)
}

// PostTxProcessing implements EvmHooks.PostTxProcessing.
func (k Keeper) PostTxProcessing(
	ctx sdk.Context,
	msg core.Message,
	receipt *ethtypes.Receipt,
) error {
	var emittingContract ethcommon.Address
	if msg.To() != nil {
		emittingContract = *msg.To()
	}
	return k.ProcessLogs(ctx, receipt.Logs, emittingContract, msg.From().Hex())
}

// ProcessLogs post-processes logs emitted by a zEVM contract; if the log contains Withdrawal event
// from registered ZRC20 contract, new CCTX will be created to trigger and track outbound
// transaction.
// Returning error from process logs does the following:
// - revert the whole tx.
// - clear the logs
// TODO: implement unit tests
// https://github.com/zeta-chain/node/issues/1759
func (k Keeper) ProcessLogs(ctx sdk.Context, logs []*ethtypes.Log, emittingContract ethcommon.Address, txOrigin string) error {
	system, found := k.fungibleKeeper.GetSystemContract(ctx)
	if !found {
		return fmt.Errorf("cannot find system contract")
	}
	connectorZEVMAddr := ethcommon.HexToAddress(system.ConnectorZevm)
	if connectorZEVMAddr == (ethcommon.Address{}) {
		return fmt.Errorf("connectorZEVM address is empty")
	}
	for _, log := range logs {
		eventZrc20Withdrawal, errZrc20 := ParseZRC20WithdrawalEvent(*log)
		eventZetaSent, errZetaSent := ParseZetaSentEvent(*log, connectorZEVMAddr)
		if errZrc20 != nil && errZetaSent != nil {
			// This log does not contain any of the two events
			continue
		}
		if eventZrc20Withdrawal != nil && eventZetaSent != nil {
			// This log contains both events, this is not possible
			ctx.Logger().Error(fmt.Sprintf("ProcessLogs: log contains both ZRC20Withdrawal and ZetaSent events, %s , %s", log.Topics, log.Data))
			continue
		}

		// We have found either eventZrc20Withdrawal or eventZetaSent
		// These cannot be processed without TSS keys, return an error if TSS is not found
		tss, found := k.zetaObserverKeeper.GetTSS(ctx)
		if !found {
			return errorsmod.Wrap(types.ErrCannotFindTSSKeys, "Cannot process logs without TSS keys")
		}

		// Do not process withdrawal events if inbound is disabled
		if !k.zetaObserverKeeper.IsInboundEnabled(ctx) {
			return observertypes.ErrInboundDisabled
		}

		// if eventZrc20Withdrawal is not nil we will try to validate it and see if it can be processed
		if eventZrc20Withdrawal != nil {
			// Check if the contract is a registered ZRC20 contract. If its not a registered ZRC20 contract, we can discard this event as it is not relevant
			coin, foundCoin := k.fungibleKeeper.GetForeignCoins(ctx, eventZrc20Withdrawal.Raw.Address.Hex())
			if !foundCoin {
				ctx.Logger().Info(fmt.Sprintf("cannot find foreign coin with contract address %s", eventZrc20Withdrawal.Raw.Address.Hex()))
				continue
			}

			// If Validation fails, we will not process the event and return and error. This condition means that the event was correct, and emitted from a registered ZRC20 contract
			// But the information entered by the user is incorrect. In this case we can return an error and roll back the transaction
			if err := ValidateZrc20WithdrawEvent(eventZrc20Withdrawal, coin.ForeignChainId); err != nil {
				return err
			}
			// If the event is valid, we will process it and create a new CCTX
			// If the process fails, we will return an error and roll back the transaction
			if err := k.ProcessZRC20WithdrawalEvent(ctx, eventZrc20Withdrawal, emittingContract, txOrigin, tss); err != nil {
				return err
			}
		}
		// if eventZetaSent is not nil we will try to validate it and see if it can be processed
		if eventZetaSent != nil {
			if err := k.ProcessZetaSentEvent(ctx, eventZetaSent, emittingContract, txOrigin, tss); err != nil {
				return err
			}
		}
	}
	return nil
}

// ProcessZRC20WithdrawalEvent creates a new CCTX to process the withdrawal event
// error indicates system error and non-recoverable; should abort
func (k Keeper) ProcessZRC20WithdrawalEvent(ctx sdk.Context, event *zrc20.ZRC20Withdrawal, emittingContract ethcommon.Address, txOrigin string, tss observertypes.TSS) error {

	ctx.Logger().Info(fmt.Sprintf("ZRC20 withdrawal to %s amount %d", hex.EncodeToString(event.To), event.Value))
	foreignCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, event.Raw.Address.Hex())
	if !found {
		return fmt.Errorf("cannot find foreign coin with emittingContract address %s", event.Raw.Address.Hex())
	}
	receiverChain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, foreignCoin.ForeignChainId)
	if receiverChain == nil {
		return errorsmod.Wrapf(observertypes.ErrSupportedChains, "chain with chainID %d not supported", foreignCoin.ForeignChainId)
	}
	senderChain, err := chains.ZetaChainFromChainID(ctx.ChainID())
	if err != nil {
		return fmt.Errorf("ProcessZRC20WithdrawalEvent: failed to convert chainID: %s", err.Error())
	}
	toAddr, err := receiverChain.EncodeAddress(event.To)
	if err != nil {
		return fmt.Errorf("cannot encode address %s: %s", event.To, err.Error())
	}

	gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress))
	if err != nil {
		return fmt.Errorf("cannot query gas limit: %s", err.Error())
	}

	// gasLimit+uint64(event.Raw.Index) to generate different cctx for multiple events in the same tx.
	msg := types.NewMsgVoteInbound(
		"",
		emittingContract.Hex(),
		senderChain.ChainId,
		txOrigin,
		toAddr,
		foreignCoin.ForeignChainId,
		math.NewUintFromBigInt(event.Value),
		"",
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		gasLimit.Uint64(),
		foreignCoin.CoinType,
		foreignCoin.Asset,
		event.Raw.Index,
	)

	// Create a new cctx with status as pending Inbound, this is created directly from the event without waiting for any observer votes
	cctx, err := types.NewCCTX(ctx, *msg, tss.TssPubkey)
	if err != nil {
		return fmt.Errorf("ProcessZRC20WithdrawalEvent: failed to initialize cctx: %s", err.Error())
	}
	cctx.SetPendingOutbound("ZRC20 withdrawal event setting to pending outbound directly")
	// Get gas price and amount
	gasprice, found := k.GetGasPrice(ctx, receiverChain.ChainId)
	if !found {
		return fmt.Errorf("gasprice not found for %s", receiverChain)
	}
	cctx.GetCurrentOutboundParam().GasPrice = fmt.Sprintf("%d", gasprice.Prices[gasprice.MedianIndex])
	cctx.GetCurrentOutboundParam().Amount = cctx.InboundParams.Amount

	EmitZRCWithdrawCreated(ctx, cctx)
	return k.ProcessCCTX(ctx, cctx, receiverChain)
}

func (k Keeper) ProcessZetaSentEvent(ctx sdk.Context, event *connectorzevm.ZetaConnectorZEVMZetaSent, emittingContract ethcommon.Address, txOrigin string, tss observertypes.TSS) error {
	ctx.Logger().Info(fmt.Sprintf(
		"Zeta withdrawal to %s amount %d to chain with chainId %d",
		hex.EncodeToString(event.DestinationAddress),
		event.ZetaValueAndGas,
		event.DestinationChainId,
	))

	if err := k.bankKeeper.BurnCoins(
		ctx,
		fungibletypes.ModuleName,
		sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(event.ZetaValueAndGas))),
	); err != nil {
		ctx.Logger().Error(fmt.Sprintf("ProcessZetaSentEvent: failed to burn coins from fungible: %s", err.Error()))
		return fmt.Errorf("ProcessZetaSentEvent: failed to burn coins from fungible: %s", err.Error())
	}

	receiverChainID := event.DestinationChainId

	receiverChain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, receiverChainID.Int64())
	if receiverChain == nil {
		return observertypes.ErrSupportedChains
	}
	// Validation if we want to send ZETA to an external chain, but there is no ZETA token.
	chainParams, found := k.zetaObserverKeeper.GetChainParamsByChainID(ctx, receiverChain.ChainId)
	if !found {
		return observertypes.ErrChainParamsNotFound
	}
	if receiverChain.IsExternalChain() && chainParams.ZetaTokenContractAddress == "" {
		return types.ErrUnableToSendCoinType
	}
	toAddr := "0x" + hex.EncodeToString(event.DestinationAddress)
	senderChain, err := chains.ZetaChainFromChainID(ctx.ChainID())
	if err != nil {
		return fmt.Errorf("ProcessZetaSentEvent: failed to convert chainID: %s", err.Error())
	}
	amount := math.NewUintFromBigInt(event.ZetaValueAndGas)
	messageString := base64.StdEncoding.EncodeToString(event.Message)
	// Bump gasLimit by event index (which is very unlikely to be larger than 1000) to always have different ZetaSent events msgs.
	msg := types.NewMsgVoteInbound(
		"",
		emittingContract.Hex(),
		senderChain.ChainId,
		txOrigin, toAddr,
		receiverChain.ChainId,
		amount,
		messageString,
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		90000,
		coin.CoinType_Zeta,
		"",
		event.Raw.Index,
	)

	// create a new cctx with status as pending Inbound,
	// this is created directly from the event without waiting for any observer votes
	cctx, err := types.NewCCTX(ctx, *msg, tss.TssPubkey)
	if err != nil {
		return fmt.Errorf("ProcessZetaSentEvent: failed to initialize cctx: %s", err.Error())
	}
	cctx.SetPendingOutbound("ZetaSent event setting to pending outbound directly")

	if err := k.PayGasAndUpdateCctx(
		ctx,
		receiverChain.ChainId,
		&cctx,
		amount,
		true,
	); err != nil {
		return fmt.Errorf("ProcessWithdrawalEvent: pay gas failed: %s", err.Error())
	}

	EmitZetaWithdrawCreated(ctx, cctx)
	return k.ProcessCCTX(ctx, cctx, receiverChain)
}

func (k Keeper) ProcessCCTX(ctx sdk.Context, cctx types.CrossChainTx, receiverChain *chains.Chain) error {
	inCctxIndex, ok := ctx.Value("inCctxIndex").(string)
	if ok {
		cctx.InboundParams.ObservedHash = inCctxIndex
	}

	if err := k.UpdateNonce(ctx, receiverChain.ChainId, &cctx); err != nil {
		return fmt.Errorf("ProcessWithdrawalEvent: update nonce failed: %s", err.Error())
	}

	k.SetCctxAndNonceToCctxAndInboundHashToCctx(ctx, cctx)
	ctx.Logger().Debug("ProcessCCTX successful \n")
	return nil
}

// ParseZRC20WithdrawalEvent tries extracting ZRC20Withdrawal event from the input logs using the zrc20 contract;
// It only returns a not-nil event if the event has been correctly validated as a valid withdrawal event
func ParseZRC20WithdrawalEvent(log ethtypes.Log) (*zrc20.ZRC20Withdrawal, error) {
	zrc20ZEVM, err := zrc20.NewZRC20Filterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
	}
	if len(log.Topics) == 0 {
		return nil, fmt.Errorf("ParseZRC20WithdrawalEvent: invalid log - no topics")
	}
	event, err := zrc20ZEVM.ParseWithdrawal(log)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// ValidateZrc20WithdrawEvent checks if the ZRC20Withdrawal event is valid
// It verifies event information for BTC chains and returns an error if the event is invalid
func ValidateZrc20WithdrawEvent(event *zrc20.ZRC20Withdrawal, chainID int64) error {
	// The event was parsed; that means the user has deposited tokens to the contract.

	if chains.IsBitcoinChain(chainID) {
		if event.Value.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("ParseZRC20WithdrawalEvent: invalid amount %s", event.Value.String())
		}
		addr, err := chains.DecodeBtcAddress(string(event.To), chainID)
		if err != nil {
			return fmt.Errorf("ParseZRC20WithdrawalEvent: invalid address %s: %s", event.To, err)
		}
		if !chains.IsBtcAddressSupported(addr) {
			return fmt.Errorf("ParseZRC20WithdrawalEvent: unsupported address %s", string(event.To))
		}
	}
	return nil
}

// ParseZetaSentEvent tries extracting ZetaSent event from connectorZEVM contract;
// returns error if the log entry is not a ZetaSent event, or is not emitted from connectorZEVM
// It only returns a not-nil event if all the error checks pass
func ParseZetaSentEvent(log ethtypes.Log, connectorZEVM ethcommon.Address) (*connectorzevm.ZetaConnectorZEVMZetaSent, error) {
	zetaConnectorZEVM, err := connectorzevm.NewZetaConnectorZEVMFilterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
	}
	if len(log.Topics) == 0 {
		return nil, fmt.Errorf("ParseZetaSentEvent: invalid log - no topics")
	}
	event, err := zetaConnectorZEVM.ParseZetaSent(log)
	if err != nil {
		return nil, err
	}

	if event.Raw.Address != connectorZEVM {
		return nil, fmt.Errorf("ParseZetaSentEvent: event address %s does not match connectorZEVM %s", event.Raw.Address.Hex(), connectorZEVM.Hex())
	}
	return event, nil
}
