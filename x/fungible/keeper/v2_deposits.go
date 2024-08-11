package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/systemcontract.sol"

	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/crypto"
)

// ProcessV2Deposit handles a deposit from an inbound tx with protocol version 2
// returns [txResponse, isContractCall, error]
// isContractCall is true if the message is non empty
func (k Keeper) ProcessV2Deposit(
	ctx sdk.Context,
	from []byte,
	senderChainID int64,
	zrc20Addr ethcommon.Address,
	to ethcommon.Address,
	amount *big.Int,
	message []byte,
	coinType coin.CoinType,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	context := systemcontract.ZContext{
		Origin:  from,
		Sender:  ethcommon.Address{},
		ChainID: big.NewInt(senderChainID),
	}

	if len(message) == 0 {
		// simple deposit
		res, err := k.DepositZRC20(ctx, zrc20Addr, to, amount)
		return res, false, err
	} else if coinType == coin.CoinType_NoAssetCall {
		// simple call
		res, err := k.CallExecute(ctx, context, zrc20Addr, amount, to, message)
		return res, true, err
	}
	// deposit and call
	res, err := k.CallDepositAndCallZRC20(ctx, context, zrc20Addr, amount, to, message)
	return res, true, err
}

// ProcessV2RevertDeposit handles a revert deposit from an inbound tx with protocol version 2
// TODO: implement revert deposit
// https://github.com/zeta-chain/node/issues/2660
func (k Keeper) ProcessV2RevertDeposit(
	ctx sdk.Context,
	zrc20Addr ethcommon.Address,
	amount *big.Int,
	revertAddress ethcommon.Address,
	callOnRevert bool,
) error {
	// zrc20 empty means no asset
	zrc20Defined := !crypto.IsEmptyAddress(zrc20Addr)

	switch {
	case !callOnRevert && !zrc20Defined:
		// no asset, no call, do nothing
		return nil
	case !callOnRevert && zrc20Defined:
		// simply deposit back to the revert address
		_, err := k.DepositZRC20(ctx, zrc20Addr, revertAddress, amount)
		return err
	case callOnRevert && !zrc20Defined:
		// no asset, call simple revert
		// CallExecuteRevert
	case callOnRevert && zrc20Defined:
		// deposit asset and revert
		// CallDepositAndRevert
	}

	return nil
}
