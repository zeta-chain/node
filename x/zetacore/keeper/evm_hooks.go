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
	"github.com/zeta-chain/zetacore/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	zetacoretypes "github.com/zeta-chain/zetacore/x/zetacore/types"
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
	return k.ProcessWithdrawalEvent(ctx, receipt.Logs, receipt.ContractAddress)
}

// FIXME: authenticate the emitting contract with foreign_coins
func (k Keeper) ProcessWithdrawalEvent(ctx sdk.Context, logs []*ethtypes.Log, contract ethcommon.Address) error {
	var event *contracts.ZRC20Withdrawal

	found := false
	for _, log := range logs {
		e, err := ParseWithdrawalEvent(*log)
		if err != nil {
			fmt.Printf("######### skip log %s #########\n", log.Topics[0].String())
			continue
		} else {
			found = true
			event = e
		}
	}
	if found {
		fmt.Printf("#############################\n")
		fmt.Printf("withdrawal to %s amount %d\n", hex.EncodeToString(event.To), event.Value)
		fmt.Printf("#############################\n")
		foreignCoinList := k.fungibleKeeper.GetAllForeignCoins(ctx)
		foundCoin := false
		receiverChain := ""
		coinType := common.CoinType_Zeta
		for _, coin := range foreignCoinList {
			if coin.Zrc20ContractAddress == event.Raw.Address.Hex() {
				receiverChain = coin.ForeignChain
				foundCoin = true
				coinType = coin.CoinType
			}
		}
		if !foundCoin {
			return fmt.Errorf("cannot find foreign coin with contract address %s", event.Raw.Address.Hex())
		}

		toAddr := "0x" + hex.EncodeToString(event.To)
		msg := zetacoretypes.NewMsgSendVoter("", contract.Hex(), common.ZETAChain.String(), toAddr, receiverChain, event.Value.String(), "", "", event.Raw.TxHash.String(), event.Raw.BlockNumber, 90000, coinType)
		sendHash := msg.Digest()

		cctx := k.CreateNewCCTX(ctx, msg, sendHash)
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
		k.UpdateNonce(ctx, receiverChain, &cctx)

		k.SetCrossChainTx(ctx, cctx)
		fmt.Printf("####setting send... ###########\n")
	}
	return nil
}

// FIXME: add check for event emitting contracts
func ParseWithdrawalEvent(log ethtypes.Log) (*contracts.ZRC20Withdrawal, error) {
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
	abi.ParseTopics(event, indexed, log.Topics[1:])
	event.Raw = log

	return event, nil
}
