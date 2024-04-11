package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k Keeper) ZEVMDepositAndCallContract(ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	inboundSenderChainID int64,
	inboundAmount *big.Int,
	data []byte,
	indexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error) {
	acc := k.evmKeeper.GetAccount(ctx, to)
	if acc == nil {
		return nil, errors.Wrap(types.ErrAccountNotFound, fmt.Sprintf("address: %s", to.String()))
	}
	if !acc.IsContract() {
		err := k.DepositCoinZeta(ctx, to, inboundAmount)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return k.ZevmOnReceive(ctx, sender.Bytes(), to, big.NewInt(inboundSenderChainID), inboundAmount, data, indexBytes)

}
func (k Keeper) ZevmOnReceive(ctx sdk.Context,
	zetaTxSender []byte,
	zetaTxReceiver ethcommon.Address,
	senderChainID *big.Int,
	amount *big.Int,
	data []byte,
	cctxIndexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error) {
	return k.CallOnReceiveZevmConnector(ctx, zetaTxSender, senderChainID, zetaTxReceiver, amount, data, cctxIndexBytes)
}

func (k Keeper) ZEVMRevertAndCallContract(ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	inboundSenderChainID int64,
	destinationChainID int64,
	remainingAmount *big.Int,
	data []byte,
	indexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error) {
	acc := k.evmKeeper.GetAccount(ctx, sender)
	if acc == nil {
		return nil, errors.Wrap(types.ErrAccountNotFound, fmt.Sprintf("address: %s", to.String()))
	}
	if !acc.IsContract() {
		err := k.DepositCoinZeta(ctx, sender, remainingAmount)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return k.ZevmOnRevert(ctx, sender, to.Bytes(), big.NewInt(inboundSenderChainID), big.NewInt(destinationChainID), remainingAmount, data, indexBytes)

}
func (k Keeper) ZevmOnRevert(ctx sdk.Context,
	zetaTxSender ethcommon.Address,
	zetaTxReceiver []byte,
	senderChainID *big.Int,
	destinationChainID *big.Int,
	amount *big.Int,
	data []byte,
	cctxIndexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error) {
	return k.CallOnRevertZevmConnector(ctx, zetaTxSender, senderChainID, zetaTxReceiver, destinationChainID, amount, data, cctxIndexBytes)
}
