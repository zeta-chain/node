package keeper

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	errorspkg "github.com/pkg/errors"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/systemcontract.sol"

	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// DepositCoinZeta immediately mints ZETA to the EVM account
func (k Keeper) DepositCoinZeta(ctx sdk.Context, to ethcommon.Address, amount *big.Int) error {
	zetaToAddress := sdk.AccAddress(to.Bytes())
	return k.MintZetaToEVMAccount(ctx, zetaToAddress, amount)
}

func (k Keeper) DepositCoinsToFungibleModule(ctx sdk.Context, amount *big.Int) error {
	return k.MintZetaToFungibleModule(ctx, amount)
}

// ZRC20DepositAndCallContract deposits ZRC20 to the EVM account and calls the contract
// returns [txResponse, isContractCall, error]
func (k Keeper) ZRC20DepositAndCallContract(
	ctx sdk.Context,
	from []byte,
	to ethcommon.Address,
	amount *big.Int,
	senderChainID int64,
	message []byte,
	coinType coin.CoinType,
	asset string,
	protocolContractVersion crosschaintypes.ProtocolContractVersion,
	isCrossChainCall bool,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	// get ZRC20 contract
	zrc20Contract, _, err := k.getAndCheckZRC20(ctx, amount, senderChainID, coinType, asset)
	if err != nil {
		return nil, false, err
	}

	// handle the deposit for protocol contract version 2
	if protocolContractVersion == crosschaintypes.ProtocolContractVersion_V2 {
		return k.ProcessDeposit(
			ctx,
			from,
			senderChainID,
			zrc20Contract,
			to,
			amount,
			message,
			coinType,
			isCrossChainCall,
		)
	}

	// check if the receiver is a contract
	// if it is, then the hook onCrossChainCall() will be called
	// if not, the zrc20 are simply transferred to the receiver
	acc := k.evmKeeper.GetAccount(ctx, to)
	if acc != nil && acc.IsContract() {
		context := systemcontract.ZContext{
			Origin:  from,
			Sender:  ethcommon.Address{},
			ChainID: big.NewInt(senderChainID),
		}
		res, err := k.CallDepositAndCall(ctx, context, zrc20Contract, to, amount, message)
		return res, true, err
	}

	// if the account is a EOC, no contract call can be made with the data
	if len(message) > 0 {
		return nil, false, types.ErrCallNonContract
	}

	res, err := k.DepositZRC20(ctx, zrc20Contract, to, amount)
	return res, false, err
}

// ProcessDeposit handles a deposit from an inbound tx with protocol version 2
// returns [txResponse, isContractCall, error]
// isContractCall is true if the message is non empty
func (k Keeper) ProcessDeposit(
	ctx sdk.Context,
	from []byte,
	senderChainID int64,
	zrc20Addr ethcommon.Address,
	to ethcommon.Address,
	amount *big.Int,
	message []byte,
	coinType coin.CoinType,
	isCrossChainCall bool,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	if coinType == coin.CoinType_Zeta {
		return nil, false, errors.New("ZETA asset is currently unsupported for deposit with V2 protocol contracts")
	}

	context := gatewayzevm.MessageContext{
		Sender:    from,
		SenderEVM: ethcommon.BytesToAddress(from),
		ChainID:   big.NewInt(senderChainID),
	}

	if coinType == coin.CoinType_NoAssetCall {
		// simple call
		res, err := k.CallExecute(ctx, context, zrc20Addr, amount, to, message)
		return res, true, err
	} else if isCrossChainCall {
		// call with asset
		res, err := k.CallDepositAndCallZRC20(ctx, context, zrc20Addr, amount, to, message)
		return res, true, err
	}

	// simple deposit
	res, err := k.DepositZRC20(ctx, zrc20Addr, to, amount)
	return res, false, err
}

// ProcessRevert handles a revert deposit from an inbound tx with protocol version 2
func (k Keeper) ProcessRevert(
	ctx sdk.Context,
	inboundSender string,
	amount *big.Int,
	chainID int64,
	coinType coin.CoinType,
	asset string,
	revertAddress ethcommon.Address,
	callOnRevert bool,
	revertMessage []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	if coinType == coin.CoinType_Zeta {
		return nil, errors.New("ZETA asset is currently unsupported for revert with V2 protocol contracts")
	}

	// get the zrc20 contract
	zrc20Addr, _, err := k.getAndCheckZRC20(
		ctx,
		amount,
		chainID,
		coinType,
		asset,
	)
	if err != nil {
		return nil, err
	}

	switch coinType {
	case coin.CoinType_NoAssetCall:
		if callOnRevert {
			// no asset, call simple revert
			res, err := k.CallExecuteRevert(ctx, inboundSender, zrc20Addr, amount, revertAddress, revertMessage)
			return res, err
		} else {
			// no asset, no call, do nothing
			return nil, nil
		}
	case coin.CoinType_ERC20, coin.CoinType_Gas:
		if callOnRevert {
			// revert with a ZRC20 asset
			res, err := k.CallDepositAndRevert(
				ctx,
				inboundSender,
				zrc20Addr,
				amount,
				revertAddress,
				revertMessage,
			)
			return res, err
		} else {
			// simply deposit back to the revert address
			res, err := k.DepositZRC20(ctx, zrc20Addr, revertAddress, amount)
			return res, err
		}
	}

	return nil, fmt.Errorf("unsupported coin type for revert %s", coinType)
}

// ProcessAbort handles an abort deposit from an inbound tx with protocol version 2
func (k Keeper) ProcessAbort(
	ctx sdk.Context,
	inboundSender string,
	amount *big.Int,
	outgoing bool,
	chainID int64,
	coinType coin.CoinType,
	asset string,
	abortAddress ethcommon.Address,
	revertMessage []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	if coinType == coin.CoinType_Zeta {
		return nil, errors.New("ZETA asset is currently unsupported for abort with V2 protocol contracts")
	}

	// get the zrc20 contract
	zrc20Addr, _, err := k.getAndCheckZRC20(
		ctx,
		amount,
		chainID,
		coinType,
		asset,
	)
	if err != nil {
		return nil, err
	}

	// if the cctx contains asset, the asset is first deposited to the abort address, separately from onAbort call
	if coinType == coin.CoinType_ERC20 || coinType == coin.CoinType_Gas {
		// simply deposit back to the revert address
		// if the deposit fails, processing the abort entirely fails
		// MsgRefundAbort can still be used to retry the operation later on
		if _, err := k.DepositZRC20(ctx, zrc20Addr, abortAddress, amount); err != nil {
			return nil, err
		}
	}

	// call onAbort
	txRes, err := k.CallExecuteAbort(
		ctx,
		inboundSender,
		zrc20Addr,
		amount,
		outgoing,
		big.NewInt(chainID),
		abortAddress,
		revertMessage,
	)
	if err != nil {
		return txRes, errors.Join(err, types.ErrOnAbortFailed)
	}
	return txRes, nil
}

// getAndCheckZRC20 returns the ZRC20 contract address and foreign coin for the given chainID and asset
// it also checks if the foreign coin is paused and if the cap is reached
func (k Keeper) getAndCheckZRC20(
	ctx sdk.Context,
	amount *big.Int,
	chainID int64,
	coinType coin.CoinType,
	asset string,
) (ethcommon.Address, types.ForeignCoins, error) {
	var zrc20Contract ethcommon.Address
	var foreignCoin types.ForeignCoins
	var found bool

	// get foreign coin
	// retrieve the gas token of the chain for no asset call
	// this simplify the current workflow and allow to pause calls by pausing the gas token
	// TODO: refactor this logic and create specific workflow for no asset call
	// https://github.com/zeta-chain/node/issues/2627

	if coinType == coin.CoinType_Gas || coinType == coin.CoinType_NoAssetCall {
		foreignCoin, found = k.GetGasCoinForForeignCoin(ctx, chainID)
		if !found {
			return ethcommon.Address{}, types.ForeignCoins{}, crosschaintypes.ErrGasCoinNotFound
		}
	} else {
		foreignCoin, found = k.GetForeignCoinFromAsset(ctx, asset, chainID)
		if !found {
			return ethcommon.Address{}, types.ForeignCoins{}, errorspkg.Wrapf(
				crosschaintypes.ErrForeignCoinNotFound,
				"asset: %s, chainID %d", asset, chainID,
			)
		}
	}
	zrc20Contract = ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress)

	// check if foreign coin is paused
	if foreignCoin.Paused {
		return ethcommon.Address{}, types.ForeignCoins{}, types.ErrPausedZRC20
	}

	// check foreign coins cap if it has a cap
	if !foreignCoin.LiquidityCap.IsNil() && !foreignCoin.LiquidityCap.IsZero() {
		liquidityCap := foreignCoin.LiquidityCap.BigInt()
		totalSupply, err := k.TotalSupplyZRC4(ctx, zrc20Contract)
		if err != nil {
			return ethcommon.Address{}, types.ForeignCoins{}, err
		}
		newSupply := new(big.Int).Add(totalSupply, amount)
		if newSupply.Cmp(liquidityCap) > 0 {
			return ethcommon.Address{}, types.ForeignCoins{}, types.ErrForeignCoinCapReached
		}
	}

	return zrc20Contract, foreignCoin, nil
}
