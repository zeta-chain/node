package keeper

import (
	"cosmossdk.io/math"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	connectorzevm "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/connectorzevm.sol"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
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
		eventWithdrawal, err := k.ParseZRC20WithdrawalEvent(ctx, *log)
		if err == nil {
			if err := k.ProcessZRC20WithdrawalEvent(ctx, eventWithdrawal, emittingContract, txOrigin); err != nil {
				return err
			}
		}
		eZeta, err := ParseZetaSentEvent(*log, connectorZEVMAddr)
		if err == nil {
			if err := k.ProcessZetaSentEvent(ctx, eZeta, emittingContract, txOrigin); err != nil {
				return err
			}
		}
	}
	return nil
}

// ProcessZRC20WithdrawalEvent creates a new CCTX to process the withdrawal event
// error indicates system error and non-recoverable; should abort
func (k Keeper) ProcessZRC20WithdrawalEvent(ctx sdk.Context, event *zrc20.ZRC20Withdrawal, emittingContract ethcommon.Address, txOrigin string) error {
	ctx.Logger().Info("ZRC20 withdrawal to %s amount %d\n", hex.EncodeToString(event.To), event.Value)

	foreignCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, event.Raw.Address.Hex())
	if !found {
		return fmt.Errorf("cannot find foreign coin with emittingContract address %s", event.Raw.Address.Hex())
	}

	recvChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(foreignCoin.ForeignChainId)
	senderChain := common.ZetaChain()
	toAddr, err := recvChain.EncodeAddress(event.To)
	if err != nil {
		return fmt.Errorf("cannot encode address %s: %s", event.To, err.Error())
	}
	gasLimit := foreignCoin.GasLimit
	msg := types.NewMsgSendVoter(
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
		gasLimit,
		foreignCoin.CoinType,
		foreignCoin.Asset,
	)
	sendHash := msg.Digest()

	cctx := k.CreateNewCCTX(ctx, msg, sendHash, types.CctxStatus_PendingOutbound, &senderChain, recvChain)
	EmitZRCWithdrawCreated(ctx, cctx)
	return k.ProcessCCTX(ctx, cctx, recvChain)
}

func (k Keeper) ProcessZetaSentEvent(ctx sdk.Context, event *connectorzevm.ZetaConnectorZEVMZetaSent, emittingContract ethcommon.Address, txOrigin string) error {
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
		fmt.Printf("burn coins failed: %s\n", err.Error())
		return fmt.Errorf("ProcessZetaSentEvent: failed to burn coins from fungible: %s", err.Error())
	}

	receiverChainID := event.DestinationChainId
	receiverChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(receiverChainID.Int64())
	if receiverChain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}
	// Validation if we want to send ZETA to external chain, but there is no ZETA token.
	coreParams, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, receiverChain.ChainId)
	if !found {
		return types.ErrNotFoundCoreParams
	}
	if receiverChain.IsExternalChain() && coreParams.ZetaTokenContractAddress == "" {
		return types.ErrUnableToSendCoinType
	}
	toAddr := "0x" + hex.EncodeToString(event.DestinationAddress)
	senderChain := common.ZetaChain()
	amount := math.NewUintFromBigInt(event.ZetaValueAndGas)

	// Bump gasLimit by event index (which is very unlikely to be larger than 1000) to always have different ZetaSent events msgs.
	msg := types.NewMsgSendVoter(
		"",
		emittingContract.Hex(),
		senderChain.ChainId,
		txOrigin, toAddr,
		receiverChain.ChainId,
		amount,
		"",
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		90000+uint64(event.Raw.Index),
		common.CoinType_Zeta,
		"",
	)
	sendHash := msg.Digest()

	// Create the CCTX
	cctx := k.CreateNewCCTX(ctx, msg, sendHash, types.CctxStatus_PendingOutbound, &senderChain, receiverChain)
	EmitZetaWithdrawCreated(ctx, cctx)
	return k.ProcessCCTX(ctx, cctx, receiverChain)
}

func (k Keeper) ProcessCCTX(ctx sdk.Context, cctx types.CrossChainTx, receiverChain *common.Chain) error {
	cctx.GetCurrentOutTxParam().Amount = cctx.InboundTxParams.Amount
	inCctxIndex, ok := ctx.Value("inCctxIndex").(string)
	if ok {
		cctx.InboundTxParams.InboundTxObservedHash = inCctxIndex
	}

	// Pay gas in Zeta and update the amount for the cctx
	if err := k.PayGasInZetaAndUpdateCctx(ctx, receiverChain.ChainId, &cctx); err != nil {
		return fmt.Errorf("ProcessWithdrawalEvent: pay gas failed: %s", err.Error())
	}

	if err := k.UpdateNonce(ctx, receiverChain.ChainId, &cctx); err != nil {
		return fmt.Errorf("ProcessWithdrawalEvent: update nonce failed: %s", err.Error())
	}

	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
	ctx.Logger().Debug("ProcessCCTX successful \n")
	return nil
}

// ParseZRC20WithdrawalEvent tries extracting Withdrawal event from registered ZRC20 contract;
// returns error if the log entry is not a Withdrawal event, or is not emitted from a
// registered ZRC20 contract
func (k Keeper) ParseZRC20WithdrawalEvent(ctx sdk.Context, log ethtypes.Log) (*zrc20.ZRC20Withdrawal, error) {
	zrc20ZEVM, err := zrc20.NewZRC20Filterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
	}
	event, err := zrc20ZEVM.ParseWithdrawal(log)
	if err != nil {
		return nil, err
	}

	_, found := k.fungibleKeeper.GetForeignCoins(ctx, event.Raw.Address.Hex())
	if !found {
		return nil, fmt.Errorf("ParseZRC20WithdrawalEvent: cannot find foreign coin with contract address %s", event.Raw.Address.Hex())
	}
	return event, nil
}

func ParseZetaSentEvent(log ethtypes.Log, connectorZEVM ethcommon.Address) (*connectorzevm.ZetaConnectorZEVMZetaSent, error) {
	zetaConnectorZEVM, err := connectorzevm.NewZetaConnectorZEVMFilterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
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
