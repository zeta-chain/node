package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	errorspkg "github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/systemcontract.sol"

	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// DepositCoinZeta immediately mints ZETA to the EVM account
func (k Keeper) DepositCoinZeta(ctx sdk.Context, to ethcommon.Address, amount *big.Int) error {
	zetaToAddress := sdk.AccAddress(to.Bytes())
	return k.MintZetaToEVMAccount(ctx, zetaToAddress, amount)
}

// ZRC20DepositAndCallContract deposits ZRC20 to the EVM account and calls the contract
// returns [txResponse, isContractCall, error]
// This function should be renamed to DepositAndCallContract as it now handles both ZRC20 and ZETA deposits
// It would be better to split into two functions V1 and Legacy logic flow
// TODO : https://github.com/zeta-chain/node/issues/3988
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
	context := gatewayzevm.MessageContext{
		Sender:    from,
		SenderEVM: ethcommon.BytesToAddress(from),
		ChainID:   big.NewInt(senderChainID),
	}

	switch coinType {
	case coin.CoinType_NoAssetCall:
		return k.processNoAssetCall(ctx, context, zrc20Addr, amount, to, message)

	case coin.CoinType_Zeta:
		return k.processZetaDeposit(ctx, context, amount, to, message, isCrossChainCall)

	case coin.CoinType_ERC20, coin.CoinType_Gas:
		return k.processZRC20Deposit(ctx, context, zrc20Addr, amount, to, message, isCrossChainCall)
	default:
		return nil, false, errorspkg.Wrapf(types.ErrProcessDeposit, " unsupported coin type %s", coinType)
	}
}

// processNoAssetCall handles deposits with no asset (simple calls)
func (k Keeper) processNoAssetCall(
	ctx sdk.Context,
	context gatewayzevm.MessageContext,
	zrc20Addr ethcommon.Address,
	amount *big.Int,
	to ethcommon.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	res, err := k.CallExecute(ctx, context, zrc20Addr, amount, to, message)
	return res, true, err
}

// processZetaDeposit handles ZETA coin deposits
func (k Keeper) processZetaDeposit(
	ctx sdk.Context,
	context gatewayzevm.MessageContext,
	amount *big.Int,
	to ethcommon.Address,
	message []byte,
	isCrossChainCall bool,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	// Use a helper function to handle the mint + execute + commit pattern
	return k.ExecuteWithMintedZeta(
		ctx,
		amount,
		func(tmpCtx sdk.Context) (*evmtypes.MsgEthereumTxResponse, bool, error) {
			if isCrossChainCall {
				res, err := k.DepositAndCallZeta(tmpCtx, context, amount, to, message)
				return res, true, err
			}

			res, err := k.DepositZeta(tmpCtx, to, amount)
			return res, false, err
		},
	)
}

// processZRC20Deposit handles ZRC20 token deposits [ZRC20 tokens exist for ERC20 and GAS tokens]
func (k Keeper) processZRC20Deposit(
	ctx sdk.Context,
	context gatewayzevm.MessageContext,
	zrc20Addr ethcommon.Address,
	amount *big.Int,
	to ethcommon.Address,
	message []byte,
	isCrossChainCall bool,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	if isCrossChainCall {
		res, err := k.CallDepositAndCallZRC20(ctx, context, zrc20Addr, amount, to, message)
		return res, true, err
	}

	res, err := k.DepositZRC20(ctx, zrc20Addr, to, amount)
	return res, false, err
}
