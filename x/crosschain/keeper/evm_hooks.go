package keeper

import (
	"cosmossdk.io/math"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
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
	target := receipt.ContractAddress
	if msg.To() != nil {
		target = *msg.To()
	}
	for _, log := range receipt.Logs {
		eZRC20, err := ParseZRC20WithdrawalEvent(*log)
		if err == nil {
			if err := k.ProcessZRC20WithdrawalEvent(ctx, eZRC20, target, ""); err != nil {
				return err
			}
		}
		eZeta, err := ParseZetaSentEvent(*log)
		if err == nil {
			if err := k.ProcessZetaSentEvent(ctx, eZeta, target, ""); err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) ProcessWithdrawalLogs(ctx sdk.Context, logs []*ethtypes.Log, contract ethcommon.Address, txOrigin string) error {
	for _, log := range logs {
		var event *contracts.ZRC20Withdrawal
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

// FIXME: authenticate the emitting contract with foreign_coins
func (k Keeper) ProcessZRC20WithdrawalEvent(ctx sdk.Context, event *contracts.ZRC20Withdrawal, contract ethcommon.Address, txOrigin string) error {
	fmt.Printf("#############################\n")
	fmt.Printf("ZRC20 withdrawal to %s amount %d\n", hex.EncodeToString(event.To), event.Value)
	fmt.Printf("#############################\n")

	// TODO , change to using GetAllForeignCoins for Chain .
	// TODO , Add receiver chain in the message
	foreignCoinList, err := k.GetAllForeignCoins(ctx)
	if err != nil {
		return err
	}
	foundCoin := false
	var receiverChainName common.ChainName
	coinType := common.CoinType_Zeta
	asset := ""
	for _, coin := range foreignCoinList {
		if coin.Zrc20ContractAddress == event.Raw.Address.Hex() {
			receiverChainName = common.ParseChainName(coin.ForeignChain)
			foundCoin = true
			coinType = coin.CoinType
			asset = coin.Erc20ContractAddress
		}
	}
	if !foundCoin {
		return fmt.Errorf("cannot find foreign coin with contract address %s", event.Raw.Address.Hex())
	}
	receiverChain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainName(receiverChainName)
	if receiverChain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}
	senderChain := common.ZetaChain()

	// FIXME: the following gas limit etc does not make sense for bitcoin
	// FIXME: use the foreign coin's gaslimit
	toAddr := "0x" + hex.EncodeToString(event.To)
	msg := zetacoretypes.NewMsgSendVoter("", contract.Hex(), senderChain.ChainId, txOrigin, toAddr, receiverChain.ChainId, math.NewUintFromBigInt(event.Value),
		"", event.Raw.TxHash.String(), event.Raw.BlockNumber, 90000, coinType, asset)
	sendHash := msg.Digest()
	cctx := k.CreateNewCCTX(ctx, msg, sendHash, zetacoretypes.CctxStatus_PendingOutbound, &senderChain, receiverChain)
	EmitZRCWithdrawCreated(ctx, cctx)
	return k.ProcessCCTX(ctx, cctx, receiverChain)
}

func (k Keeper) ProcessZetaSentEvent(ctx sdk.Context, event *contracts.ZetaConnectorZEVMZetaSent, contract ethcommon.Address, txOrigin string) error {
	fmt.Printf("#############################\n")
	fmt.Printf("Zeta withdrawal to %s amount %d to chain with chainId %d\n", hex.EncodeToString(event.DestinationAddress), event.ZetaValueAndGas, event.DestinationChainId)
	fmt.Printf("#############################\n")

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

// FIXME: add check for event emitting contracts
func ParseZRC20WithdrawalEvent(log ethtypes.Log) (*contracts.ZRC20Withdrawal, error) {
	zrc20Abi, err := contracts.ZRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	event := new(contracts.ZRC20Withdrawal)
	eventName := "Withdrawal"
	if log.Topics[0] != zrc20Abi.Events[eventName].ID {
		return nil, fmt.Errorf("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err := zrc20Abi.UnpackIntoInterface(event, eventName, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range zrc20Abi.Events[eventName].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	err = abi.ParseTopics(event, indexed, log.Topics[1:])
	if err != nil {
		return nil, err
	}
	event.Raw = log

	return event, nil
}

// FIXME: add check for event emitting contracts
// TODO: use the abigen'd filter instead of manual parsing for other events above
func ParseZetaSentEvent(log ethtypes.Log) (*contracts.ZetaConnectorZEVMZetaSent, error) {
	zetaConnectorZEVM, err := contracts.NewZetaConnectorZEVMFilterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
	}
	event, err := zetaConnectorZEVM.ParseZetaSent(log)
	if err != nil {
		return nil, err
	}

	return event, nil
}
