package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k Keeper) ZevmOnReceive(ctx sdk.Context,
	zetaTxSender []byte,
	zetaTxReceiver eth.Address,
	senderChainID *big.Int,
	amount *big.Int,
	data []byte,
	cctxIndexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error) {
	acc := k.evmKeeper.GetAccount(ctx, zetaTxReceiver)
	if acc == nil {
		return nil, errors.Wrap(types.ErrAccountNotFound, fmt.Sprintf("address: %s", zetaTxReceiver.String()))
	}
	if !acc.IsContract() {
		return nil, errors.Wrap(types.ErrCallNonContract, fmt.Sprintf("address is not a contract: %s", zetaTxReceiver.String()))
	}
	evmCallResponse, err := k.CallOnReceiveZevmConnector(ctx, zetaTxSender, senderChainID, zetaTxReceiver, amount, data, cctxIndexBytes)
	if err != nil {
		return nil, err
	}
	return evmCallResponse, nil
}

func (k Keeper) ZevmOnRevert(ctx sdk.Context,
	zetaTxSender eth.Address,
	zetaTxReceiver []byte,
	senderChainID *big.Int,
	destinationChainID *big.Int,
	amount *big.Int,
	data []byte,
	cctxIndexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error) {
	acc := k.evmKeeper.GetAccount(ctx, zetaTxSender)
	if acc == nil {
		return nil, errors.Wrap(types.ErrAccountNotFound, fmt.Sprintf("address: %s", zetaTxSender.String()))

	}
	if !acc.IsContract() {
		return nil, errors.Wrap(types.ErrCallNonContract, fmt.Sprintf("to address is not a contract: %s", zetaTxSender.String()))
	}

	evmCallResponse, err := k.CallOnRevertZevmConnector(ctx, zetaTxSender, senderChainID, zetaTxReceiver, destinationChainID, amount, data, cctxIndexBytes)
	if err != nil {
		return nil, err
	}
	return evmCallResponse, nil
}
