package keeper

import (
	"encoding/hex"
	"fmt"

	"cosmossdk.io/math"

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

	zetacoretypes "github.com/zeta-chain/zetacore/x/crosschain/types"
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
	system, found := k.fungibleKeeper.GetSystemContract(ctx)
	if !found {
		return fmt.Errorf("cannot find system contract")
	}
	connectorZEVMAddr := ethcommon.HexToAddress(system.ConnectorZevm)
	if connectorZEVMAddr == (ethcommon.Address{}) {
		return fmt.Errorf("connectorZEVM address is empty")
	}

	target := receipt.ContractAddress
	if msg.To() != nil {
		target = *msg.To()
	}
	for _, log := range receipt.Logs {
		//TODO: should authenticate emitting contract address in the Parse* functions
		eZRC20, err := ParseZRC20WithdrawalEvent(*log)
		if err == nil {
			if err := k.ProcessZRC20WithdrawalEvent(ctx, eZRC20, target, ""); err != nil {
				return err
			}
		}
		eZeta, err := ParseZetaSentEvent(*log)
		if err == nil {
			if eZeta.Raw.Address != connectorZEVMAddr {
				k.Logger(ctx).Error("Warning: ZetaSent event address is not connectorZEVMAddr", "event address", eZeta.Raw.Address.Hex(), "connectorZEVMAddr", connectorZEVMAddr.Hex())
				continue
			}
			if err := k.ProcessZetaSentEvent(ctx, eZeta, target, ""); err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) ProcessWithdrawalLogs(ctx sdk.Context, logs []*ethtypes.Log, contract ethcommon.Address, txOrigin string) error {
	for _, log := range logs {
		var event *zrc20.ZRC20Withdrawal
		event, err := ParseZRC20WithdrawalEvent(*log)
		if err != nil {
			fmt.Printf("######### skip log %s #########\n", log.Topics[0].String())
		} else {
			if err = k.ProcessZRC20WithdrawalEvent(ctx, event, contract, txOrigin); err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) ProcessZRC20WithdrawalEvent(ctx sdk.Context, event *zrc20.ZRC20Withdrawal, contract ethcommon.Address, txOrigin string) error {
	ctx.Logger().Info("ZRC20 withdrawal to %s amount %d\n", hex.EncodeToString(event.To), event.Value)

	foreignCoinList, err := k.GetAllForeignCoins(ctx)
	if err != nil {
		return err
	}
	foundCoin := false
	coinType := common.CoinType_Zeta
	asset := ""
	var receiverChainID int64
	for _, coin := range foreignCoinList {
		zrc20Addr := ethcommon.HexToAddress(coin.Zrc20ContractAddress)
		if zrc20Addr == event.Raw.Address && event.Raw.Address != (ethcommon.Address{}) {
			receiverChainID = coin.ForeignChainId
			foundCoin = true
			coinType = coin.CoinType
			asset = coin.Asset
		}
	}
	if !foundCoin {
		return fmt.Errorf("cannot find foreign coin with contract address %s", event.Raw.Address.Hex())
	}

	recvChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(receiverChainID)
	senderChain := common.ZetaChain()
	toAddr := "0x" + hex.EncodeToString(event.To)
	// FIXME: use proper gas limit
	msg := zetacoretypes.NewMsgSendVoter("", contract.Hex(), senderChain.ChainId, txOrigin, toAddr, receiverChainID, math.NewUintFromBigInt(event.Value),
		"", event.Raw.TxHash.String(), event.Raw.BlockNumber, 90000, coinType, asset)
	sendHash := msg.Digest()
	cctx := k.CreateNewCCTX(ctx, msg, sendHash, zetacoretypes.CctxStatus_PendingOutbound, &senderChain, recvChain)
	EmitZRCWithdrawCreated(ctx, cctx)
	return k.ProcessCCTX(ctx, cctx, recvChain)
}

func (k Keeper) ProcessZetaSentEvent(ctx sdk.Context, event *connectorzevm.ZetaConnectorZEVMZetaSent, contract ethcommon.Address, txOrigin string) error {

	ctx.Logger().Info("Zeta withdrawal to %s amount %d to chain with chainId %d\n", hex.EncodeToString(event.DestinationAddress), event.ZetaValueAndGas, event.DestinationChainId)
	if err := k.bankKeeper.BurnCoins(ctx, "fungible", sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(event.ZetaValueAndGas)))); err != nil {
		fmt.Printf("burn coins failed: %s\n", err.Error())
		return fmt.Errorf("ProcessWithdrawalEvent: failed to burn coins from fungible: %s", err.Error())
	}
	receiverChainID := event.DestinationChainId
	receiverChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(receiverChainID.Int64())
	if receiverChain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}
	//receiverChain := "BSCTESTNET" // TODO: parse with config.FindByChainID(eventZetaSent.ToChainID) after moving config to common
	toAddr := "0x" + hex.EncodeToString(event.DestinationAddress)
	senderChain := common.ZetaChain()
	amount := math.NewUintFromBigInt(event.ZetaValueAndGas)
	msg := zetacoretypes.NewMsgSendVoter("", contract.Hex(), senderChain.ChainId, txOrigin, toAddr, receiverChain.ChainId, amount, "", event.Raw.TxHash.String(), event.Raw.BlockNumber, 90000, common.CoinType_Zeta, "")
	sendHash := msg.Digest()
	cctx := k.CreateNewCCTX(ctx, msg, sendHash, zetacoretypes.CctxStatus_PendingOutbound, &senderChain, receiverChain)
	EmitZetaWithdrawCreated(ctx, cctx)
	return k.ProcessCCTX(ctx, cctx, receiverChain)
}

func (k Keeper) ProcessCCTX(ctx sdk.Context, cctx zetacoretypes.CrossChainTx, receiverChain *common.Chain) error {
	cctx.GetCurrentOutTxParam().Amount = cctx.InboundTxParams.Amount
	cctx.GetCurrentOutTxParam().OutboundTxGasLimit = 90_000
	gasprice, found := k.GetGasPrice(ctx, receiverChain.ChainId)
	if !found {
		fmt.Printf("gasprice not found for %s\n", receiverChain)
		return fmt.Errorf("gasprice not found for %s", receiverChain)
	}
	cctx.GetCurrentOutTxParam().OutboundTxGasPrice = fmt.Sprintf("%d", gasprice.Prices[gasprice.MedianIndex])
	cctx.CctxStatus.Status = zetacoretypes.CctxStatus_PendingOutbound
	inCctxIndex, ok := ctx.Value("inCctxIndex").(string)
	if ok {
		cctx.InboundTxParams.InboundTxObservedHash = inCctxIndex
	}
	err := k.UpdateNonce(ctx, receiverChain.ChainId, &cctx)
	if err != nil {
		return fmt.Errorf("ProcessWithdrawalEvent: update nonce failed: %s", err.Error())
	}

	k.SetCrossChainTx(ctx, cctx)
	fmt.Printf("####setting send... ###########\n")
	return nil
}

func ParseZRC20WithdrawalEvent(log ethtypes.Log) (*zrc20.ZRC20Withdrawal, error) {
	zrc20ZEVM, err := zrc20.NewZRC20Filterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
	}
	event, err := zrc20ZEVM.ParseWithdrawal(log)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func ParseZetaSentEvent(log ethtypes.Log) (*connectorzevm.ZetaConnectorZEVMZetaSent, error) {
	zetaConnectorZEVM, err := connectorzevm.NewZetaConnectorZEVMFilterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
	}
	event, err := zetaConnectorZEVM.ParseZetaSent(log)
	if err != nil {
		return nil, err
	}

	return event, nil
}
