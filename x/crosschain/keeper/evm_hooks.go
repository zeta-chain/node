package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	zetacoretypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

var _ evmtypes.EvmHooks = Hooks{}

const ZetaBridgeAddress = "0x2DD9830f8Ac0E421aFF9B7c8f7E9DF6F65DBF6Ea"

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
	return k.ProcessWithdrawalEvent(ctx, receipt.Logs, target, "")
}

// FIXME: authenticate the emitting contract with foreign_coins
func (k Keeper) ProcessWithdrawalEvent(ctx sdk.Context, logs []*ethtypes.Log, contract ethcommon.Address, txOrigin string) error {
	var eventZRC20Withdrawal *contracts.ZRC20Withdrawal
	var eventZetaSent *contracts.ZETABridgeZetaSent

	foundZRC20Withdrawal := false
	foundZetaSent := false
	for _, log := range logs {
		var eZeta *contracts.ZETABridgeZetaSent
		var eZRC20 *contracts.ZRC20Withdrawal
		eZRC20, err := ParseZRC20WithdrawalEvent(*log)
		if err != nil {
			eZeta, err = ParseZetaSentEvent(*log)
			if err == nil {
				foundZetaSent = true
				eventZetaSent = eZeta
			}
		} else {
			foundZRC20Withdrawal = true
			eventZRC20Withdrawal = eZRC20
		}
	}

	if foundZRC20Withdrawal {
		fmt.Printf("#############################\n")
		fmt.Printf("ZRC20 withdrawal to %s amount %d\n", hex.EncodeToString(eventZRC20Withdrawal.To), eventZRC20Withdrawal.Value)
		fmt.Printf("#############################\n")
		foreignCoinList := k.fungibleKeeper.GetAllForeignCoins(ctx)
		foundCoin := false
		receiverChain := ""
		coinType := common.CoinType_Zeta
		for _, coin := range foreignCoinList {
			if coin.Zrc20ContractAddress == eventZRC20Withdrawal.Raw.Address.Hex() {
				receiverChain = coin.ForeignChain
				foundCoin = true
				coinType = coin.CoinType
			}
		}
		if !foundCoin {
			return fmt.Errorf("cannot find foreign coin with contract address %s", eventZRC20Withdrawal.Raw.Address.Hex())
		}

		toAddr := "0x" + hex.EncodeToString(eventZRC20Withdrawal.To)
		msg := zetacoretypes.NewMsgSendVoter("", eventZRC20Withdrawal.Raw.Address.Hex(), common.ZETAChain.String(), txOrigin, toAddr, receiverChain, eventZRC20Withdrawal.Value.String(), "", "", eventZRC20Withdrawal.Raw.TxHash.String(), eventZRC20Withdrawal.Raw.BlockNumber, 90000, coinType)
		sendHash := msg.Digest()
		cctx := k.CreateNewCCTX(ctx, msg, sendHash, zetacoretypes.CctxStatus_PendingOutbound)
		EmitZRCWithdrawCreated(ctx, cctx)
		return k.ProcessCCTX(ctx, cctx, receiverChain)
	}

	if foundZetaSent {
		fmt.Printf("#############################\n")
		fmt.Printf("Zeta withdrawal to %s amount %d to chain with chainId %d\n", hex.EncodeToString(eventZetaSent.To), eventZetaSent.Value, eventZetaSent.ToChainID)
		fmt.Printf("#############################\n")

		if err := k.bankKeeper.BurnCoins(ctx, "fungible", sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(eventZetaSent.Value)))); err != nil {
			return fmt.Errorf("ProcessWithdrawalEvent: failed to burn coins from fungible: %s", err.Error())
		}

		receiverChain := "BSCTESTNET" // TODO: parse with config.FindByChainID(eventZetaSent.ToChainID) after moving config to common
		toAddr := "0x" + hex.EncodeToString(eventZetaSent.To)
		msg := zetacoretypes.NewMsgSendVoter("", eventZetaSent.Raw.Address.Hex(), common.ZETAChain.String(), txOrigin, toAddr, receiverChain, eventZetaSent.Value.String(), "", "", eventZetaSent.Raw.TxHash.String(), eventZetaSent.Raw.BlockNumber, 90000, common.CoinType_Zeta)
		sendHash := msg.Digest()
		cctx := k.CreateNewCCTX(ctx, msg, sendHash, zetacoretypes.CctxStatus_PendingOutbound)
		EmitZetaWithdrawCreated(ctx, cctx)
		return k.ProcessCCTX(ctx, cctx, receiverChain)
	}
	return nil
}

func (k Keeper) ProcessCCTX(ctx sdk.Context, cctx zetacoretypes.CrossChainTx, receiverChain string) error {
	cctx.ZetaMint = cctx.ZetaBurnt
	cctx.OutBoundTxParams.OutBoundTxGasLimit = 90_000
	gasprice, found := k.GetGasPrice(ctx, receiverChain)
	if !found {
		fmt.Printf("gasprice not found for %s\n", receiverChain)
		return fmt.Errorf("gasprice not found for %s", receiverChain)
	}
	cctx.OutBoundTxParams.OutBoundTxGasPrice = fmt.Sprintf("%d", gasprice.Prices[gasprice.MedianIndex])
	cctx.CctxStatus.Status = zetacoretypes.CctxStatus_PendingOutbound
	inCctxIndex, ok := ctx.Value("inCctxIndex").(string)
	if ok {
		cctx.InBoundTxParams.InBoundTxObservedHash = inCctxIndex
	}
	err := k.UpdateNonce(ctx, receiverChain, &cctx)
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
func ParseZetaSentEvent(log ethtypes.Log) (*contracts.ZETABridgeZetaSent, error) {
	zetaBridgeABI, err := contracts.ZETABridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	event := new(contracts.ZETABridgeZetaSent)
	eventName := "ZetaSent"
	if log.Topics[0] != zetaBridgeABI.Events[eventName].ID {
		return nil, fmt.Errorf("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err := zetaBridgeABI.UnpackIntoInterface(event, eventName, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range zetaBridgeABI.Events[eventName].Inputs {
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
