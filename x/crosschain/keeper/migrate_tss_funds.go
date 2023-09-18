package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) MigrateTSSFunds(ctx sdk.Context) error {
	currentTssAddress, found := k.GetTSS(ctx)
	if !found {
		return types.ErrCannotFindTSSKeys
	}
	tssList := k.GetAllTSS(ctx)
	if len(tssList) < 2 {
		return errors.New("crosschain", 1999999, "Cannot migrate TSS funds, only one TSS found")
	}
	// Sort tssList by FinalizedZetaHeight
	sort.SliceStable(tssList, func(i, j int) bool {
		return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
	})
	newTssAdress := tssList[len(tssList)-1]

	//ethAddressOld, err := getTssAddrEVM(currentTssAddress.TssPubkey)
	//if err != nil {
	//	return err
	//}
	//btcAddressOld, err := getTssAddrBTC(currentTssAddress.TssPubkey)
	//if err != nil {
	//	return err
	//}
	ethAddressNew, err := getTssAddrEVM(newTssAdress.TssPubkey)
	if err != nil {
		return err
	}
	btcAddressNew, err := getTssAddrBTC(newTssAdress.TssPubkey)
	if err != nil {
		return err
	}
	for _, chain := range common.DefaultChainsList() {
		medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chain.ChainId)
		if !isFound {
			return types.ErrUnableToGetGasPrice
		}
		cctx := types.CrossChainTx{
			Creator:        "",
			Index:          "",
			ZetaFees:       sdk.Uint{},
			RelayedMessage: "",
			CctxStatus: &types.Status{
				Status: types.CctxStatus_PendingOutbound,
			},
			InboundTxParams: nil,
			OutboundTxParams: []*types.OutboundTxParams{{
				Receiver:                         "",
				ReceiverChainId:                  chain.ChainId,
				CoinType:                         common.CoinType_Gas,
				Amount:                           sdk.Uint{},
				OutboundTxTssNonce:               0,
				OutboundTxGasLimit:               1_000_000,
				OutboundTxGasPrice:               medianGasPrice.MulUint64(2).String(),
				OutboundTxHash:                   "",
				OutboundTxBallotIndex:            "",
				OutboundTxObservedExternalHeight: 0,
				OutboundTxGasUsed:                0,
				OutboundTxEffectiveGasPrice:      sdk.Int{},
				OutboundTxEffectiveGasLimit:      0,
				TssPubkey:                        currentTssAddress.TssPubkey,
			}},
		}
		if common.IsEVMChain(chain.ChainId) {
			cctx.GetCurrentOutTxParam().Receiver = ethAddressNew.String()
		}
		if common.IsBitcoinChain(chain.ChainId) {
			cctx.GetCurrentOutTxParam().Receiver = btcAddressNew
		}
		err := k.UpdateNonce(ctx, chain.ChainId, &cctx)
		if err != nil {
			return err
		}
		k.SetCrossChainTx(ctx, cctx)
	}
	return nil
}
