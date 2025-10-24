package keeper

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	connectorzevm "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnectorzevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
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
func (h Hooks) PostTxProcessing(
	ctx sdk.Context,
	_ ethcommon.Address,
	msg core.Message,
	receipt *ethtypes.Receipt,
) error {
	return h.k.PostTxProcessing(ctx, msg, receipt)
}

// PostTxProcessing implements EvmHooks.PostTxProcessing.
func (k Keeper) PostTxProcessing(
	ctx sdk.Context,
	msg core.Message,
	receipt *ethtypes.Receipt,
) error {
	var emittingContract ethcommon.Address
	if msg.To != nil {
		emittingContract = *msg.To
	}
	return k.ProcessLogs(ctx, receipt.Logs, emittingContract, msg.From.Hex())
}

// ProcessLogs post-processes logs emitted by a zEVM contract; if the log contains Withdrawal event
// from registered ZRC20 contract, new CCTX will be created to trigger and track outbound
// transaction.
// Returning error from process logs does the following:
// - revert the whole tx.
// - clear the logs
// TODO: implement unit tests
// https://github.com/zeta-chain/node/issues/1759
// TODO: refactor and simplify
// https://github.com/zeta-chain/node/issues/2627
func (k Keeper) ProcessLogs(
	ctx sdk.Context,
	logs []*ethtypes.Log,
	emittingAddress ethcommon.Address,
	txOrigin string,
) error {
	system, found := k.fungibleKeeper.GetSystemContract(ctx)
	if !found {
		return fmt.Errorf("cannot find system contract")
	}
	connectorZEVMAddr := ethcommon.HexToAddress(system.ConnectorZevm)
	if connectorZEVMAddr == (ethcommon.Address{}) {
		return fmt.Errorf("connectorZEVM address is empty")
	}
	gatewayAddr := ethcommon.HexToAddress(system.Gateway)

	// read the logs and process inbounds from emitted events
	// run the processing for the v1 and the v2 protocol contracts
	for _, log := range logs {
		if !crypto.IsEmptyAddress(gatewayAddr) {
			if err := k.ProcessZEVMInboundV2(ctx, log, gatewayAddr, txOrigin); err != nil {
				// Emit an error event so the reason for the failure can be tracked
				EmitInboundProcessingFailure(ctx, log.TxHash.String(), err.Error())

				return errors.Wrap(err, "failed to process ZEVM inbound V2")
			}
		}
		if err := k.ProcessZEVMInboundV1(ctx, log, connectorZEVMAddr, emittingAddress, txOrigin); err != nil {
			return errors.Wrap(err, "failed to process ZEVM inbound V1")
		}
	}

	return nil
}

// ProcessZEVMInboundV1 processes the logs emitted by the zEVM contract for V1 protocol contracts
// it parses logs from Connector and ZRC20 contracts and processes them accordingly
func (k Keeper) ProcessZEVMInboundV1(
	ctx sdk.Context,
	log *ethtypes.Log,
	connectorZEVMAddr,
	emittingAddress ethcommon.Address,
	txOrigin string,
) error {
	eventZRC20Withdrawal, errZrc20 := ParseZRC20WithdrawalEvent(*log)
	eventZETASent, errZetaSent := ParseZetaSentEvent(*log, connectorZEVMAddr)
	if errZrc20 != nil && errZetaSent != nil {
		// This log does not contain any of the two events
		return nil
	}
	if eventZRC20Withdrawal != nil && eventZETASent != nil {
		// This log contains both events, this is not possible
		ctx.Logger().
			Error(fmt.Sprintf("ProcessLogs: log contains both ZRC20Withdrawal and ZetaSent events, %s , %s", log.Topics, log.Data))
		return nil
	}

	// if eventZrc20Withdrawal is not nil we will try to validate it and see if it can be processed
	if eventZRC20Withdrawal != nil {
		// Check if the contract is a registered ZRC20 contract. If its not a registered ZRC20 contract, we can discard this event as it is not relevant
		coin, foundCoin := k.fungibleKeeper.GetForeignCoins(ctx, eventZRC20Withdrawal.Raw.Address.Hex())
		if !foundCoin {
			ctx.Logger().
				Info(fmt.Sprintf("cannot find foreign coin with contract address %s", eventZRC20Withdrawal.Raw.Address.Hex()))
			return nil
		}

		// If Validation fails, we will not process the event and return and error. This condition means that the event was correct, and emitted from a registered ZRC20 contract
		// But the information entered by the user is incorrect. In this case we can return an error and roll back the transaction
		if err := k.ValidateZRC20WithdrawEvent(ctx, eventZRC20Withdrawal, coin.ForeignChainId, coin.CoinType); err != nil {
			return err
		}
		// If the event is valid, we will process it and create a new CCTX
		// If the process fails, we will return an error and roll back the transaction
		if err := k.ProcessZRC20WithdrawalEvent(ctx, eventZRC20Withdrawal, emittingAddress, txOrigin); err != nil {
			return err
		}
	}
	// if eventZetaSent is not nil we will try to validate it and see if it can be processed
	if eventZETASent != nil {
		if err := k.ProcessZetaSentEvent(ctx, eventZETASent, emittingAddress, txOrigin); err != nil {
			return err
		}
	}
	return nil
}

// ProcessZRC20WithdrawalEvent creates a new CCTX to process the withdrawal event
// error indicates system error and non-recoverable; should abort
func (k Keeper) ProcessZRC20WithdrawalEvent(
	ctx sdk.Context,
	event *zrc20.ZRC20Withdrawal,
	emittingContract ethcommon.Address,
	txOrigin string,
) error {
	ctx.Logger().Info(fmt.Sprintf("ZRC20 withdrawal to %s amount %d", hex.EncodeToString(event.To), event.Value))
	foreignCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, event.Raw.Address.Hex())
	if !found {
		return fmt.Errorf("cannot find foreign coin with emittingContract address %s", event.Raw.Address.Hex())
	}

	receiverChain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, foreignCoin.ForeignChainId)
	if !found {
		return errorsmod.Wrapf(
			observertypes.ErrSupportedChains,
			"chain with chainID %d not supported",
			foreignCoin.ForeignChainId,
		)
	}

	senderChain, err := chains.ZetaChainFromCosmosChainID(ctx.ChainID())
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
		sdkmath.NewUintFromBigInt(event.Value),
		"",
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		gasLimit.Uint64(),
		foreignCoin.CoinType,
		foreignCoin.Asset,
		uint64(event.Raw.Index),
		types.ProtocolContractVersion_V1,
		false, // not relevant for v1
		types.InboundStatus_SUCCESS,
		types.ConfirmationMode_SAFE,
	)

	cctx, err := k.ValidateInbound(ctx, msg, false)
	if err != nil {
		return err
	}

	if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
		return errors.New("cctx aborted")
	}

	EmitZRCWithdrawCreated(ctx, *cctx)

	return nil
}

func (k Keeper) ProcessZetaSentEvent(
	ctx sdk.Context,
	event *connectorzevm.ZetaConnectorZEVMZetaSent,
	emittingContract ethcommon.Address,
	txOrigin string,
) error {
	ctx.Logger().Info(fmt.Sprintf(
		"Zeta withdrawal to %s amount %d to chain with chainId %d",
		hex.EncodeToString(event.DestinationAddress),
		event.ZetaValueAndGas,
		event.DestinationChainId,
	))

	if err := k.bankKeeper.BurnCoins(
		ctx,
		fungibletypes.ModuleName,
		sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(event.ZetaValueAndGas))),
	); err != nil {
		ctx.Logger().Error(fmt.Sprintf("ProcessZetaSentEvent: failed to burn coins from fungible: %s", err.Error()))
		return fmt.Errorf("ProcessZetaSentEvent: failed to burn coins from fungible: %s", err.Error())
	}

	receiverChainID := event.DestinationChainId

	receiverChain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, receiverChainID.Int64())
	if !found {
		return observertypes.ErrSupportedChains
	}

	// Validation if we want to send ZETA to an external chain, but there is no ZETA token.
	chainParams, found := k.zetaObserverKeeper.GetChainParamsByChainID(ctx, receiverChain.ChainId)
	if !found {
		return observertypes.ErrChainParamsNotFound
	}

	if receiverChain.IsExternalChain() &&
		(chainParams.ZetaTokenContractAddress == "" || chainParams.ZetaTokenContractAddress == constant.EVMZeroAddress) {
		return types.ErrUnableToSendCoinType
	}

	toAddr := "0x" + hex.EncodeToString(event.DestinationAddress)
	senderChain, err := chains.ZetaChainFromCosmosChainID(ctx.ChainID())
	if err != nil {
		return fmt.Errorf("ProcessZetaSentEvent: failed to convert chainID: %s", err.Error())
	}

	amount := sdkmath.NewUintFromBigInt(event.ZetaValueAndGas)
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
		uint64(event.Raw.Index),
		types.ProtocolContractVersion_V1,
		false, // not relevant for v1
		types.InboundStatus_SUCCESS,
		types.ConfirmationMode_SAFE,
	)

	cctx, err := k.ValidateInbound(ctx, msg, true)
	if err != nil {
		return err
	}

	if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
		return errors.New("cctx aborted")
	}

	EmitZetaWithdrawCreated(ctx, *cctx)
	return nil
}

// ValidateZRC20WithdrawEvent checks if the ZRC20Withdrawal event is valid
// It verifies event information for BTC chains and returns an error if the event is invalid
func (k Keeper) ValidateZRC20WithdrawEvent(
	ctx sdk.Context,
	event *zrc20.ZRC20Withdrawal,
	chainID int64,
	coinType coin.CoinType,
) error {
	// The event was parsed; that means the user has deposited tokens to the contract.
	return k.validateOutbound(ctx, chainID, coinType, event.Value, event.To)
}

// validateOutbound validates the data of a ZRC20 Withdrawals and Call event (version 1 or 2)
// it checks if the withdrawal amount is valid and the destination address is supported depending on the chain
func (k Keeper) validateOutbound(
	ctx sdk.Context,
	chainID int64,
	coinType coin.CoinType,
	value *big.Int,
	to []byte,
) error {
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)
	if chains.IsBitcoinChain(chainID, additionalChains) {
		if value.Cmp(big.NewInt(constant.BTCWithdrawalDustAmount)) < 0 {
			return errorsmod.Wrapf(
				types.ErrInvalidWithdrawalAmount,
				"withdraw amount %s is less than dust amount %d",
				value.String(),
				constant.BTCWithdrawalDustAmount,
			)
		}
		addr, err := chains.DecodeBtcAddress(string(to), chainID)
		if err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "invalid Bitcoin address %s", string(to))
		}
		if !chains.IsBtcAddressSupported(addr) {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "unsupported Bitcoin address %s", string(to))
		}
	} else if chains.IsSolanaChain(chainID, additionalChains) {
		// The rent exempt check is not needed for ZRC20 (SPL) tokens because withdrawing SPL token
		// already needs a non-trivial amount of SOL for potential ATA creation so we can skip the check,
		// and also not needed for simple no asset call.
		if coinType == coin.CoinType_Gas && value.Cmp(big.NewInt(constant.SolanaWalletRentExempt)) < 0 {
			return errorsmod.Wrapf(
				types.ErrInvalidWithdrawalAmount,
				"withdraw amount %s is less than rent exempt %d",
				value.String(),
				constant.SolanaWalletRentExempt,
			)
		}
		_, err := chains.DecodeSolanaWalletAddress(string(to))
		if err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "invalid Solana address %s", string(to))
		}
	} else if chains.IsSuiChain(chainID, additionalChains) {
		// check the string format of the address is valid
		addr := sui.DecodeAddress(to)
		if err := sui.ValidateAddress(addr); err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "invalid Sui address %s", string(to))
		}
	}

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

// ParseZetaSentEvent tries extracting ZetaSent event from connectorZEVM contract;
// returns error if the log entry is not a ZetaSent event, or is not emitted from connectorZEVM
// It only returns a not-nil event if all the error checks pass
func ParseZetaSentEvent(
	log ethtypes.Log,
	connectorZEVM ethcommon.Address,
) (*connectorzevm.ZetaConnectorZEVMZetaSent, error) {
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
		return nil, fmt.Errorf(
			"ParseZetaSentEvent: event address %s does not match connectorZEVM %s",
			event.Raw.Address.Hex(),
			connectorZEVM.Hex(),
		)
	}
	return event, nil
}
