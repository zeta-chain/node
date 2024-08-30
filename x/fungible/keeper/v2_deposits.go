package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/systemcontract.sol"

	"github.com/zeta-chain/node/pkg/coin"
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
func (k Keeper) ProcessV2RevertDeposit(
	ctx sdk.Context,
	amount *big.Int,
	chainID int64,
	coinType coin.CoinType,
	asset string,
	revertAddress ethcommon.Address,
	callOnRevert bool,
	revertMessage []byte,
) error {
	// get the zrc20 contract
	zrc20Addr, _, err := k.getAndCheckZRC20(
		ctx,
		amount,
		chainID,
		coinType,
		asset,
	)
	if err != nil {
		return err
	}

	switch coinType {
	case coin.CoinType_NoAssetCall:

		if callOnRevert {
			// no asset, call simple revert
			_, err := k.CallExecuteRevert(ctx, zrc20Addr, amount, revertAddress, revertMessage)
			return err
		} else {
			// no asset, no call, do nothing
			return nil
		}
	case coin.CoinType_Zeta:
		return errors.New("ZETA asset is currently unsupported for revert with V2 protocol contracts")
	case coin.CoinType_ERC20, coin.CoinType_Gas:
		if callOnRevert {
			// revert with a ZRC20 asset
			_, err := k.CallDepositAndRevert(
				ctx,
				zrc20Addr,
				amount,
				revertAddress,
				revertMessage,
			)
			return err
		} else {
			// simply deposit back to the revert address
			_, err := k.DepositZRC20(ctx, zrc20Addr, revertAddress, amount)
			return err
		}
	}

	return fmt.Errorf("unsupported coin type for revert %s", coinType)
}
