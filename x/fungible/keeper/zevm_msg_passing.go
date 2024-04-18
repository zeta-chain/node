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
	return k.CallOnReceiveZevmConnector(ctx, sender.Bytes(), big.NewInt(inboundSenderChainID), to, inboundAmount, data, indexBytes)

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
	return k.CallOnRevertZevmConnector(ctx, sender, big.NewInt(inboundSenderChainID), to.Bytes(), big.NewInt(destinationChainID), remainingAmount, data, indexBytes)
}
