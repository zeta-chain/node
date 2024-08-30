package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	"github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/zevm/systemcontract.sol"

	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// DepositCoinZeta immediately mints ZETA to the EVM account
func (k Keeper) DepositCoinZeta(ctx sdk.Context, to eth.Address, amount *big.Int) error {
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
	to eth.Address,
	amount *big.Int,
	senderChainID int64,
	message []byte,
	coinType coin.CoinType,
	asset string,
	protocolContractVersion crosschaintypes.ProtocolContractVersion,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	// get ZRC20 contract
	zrc20Contract, _, err := k.getAndCheckZRC20(ctx, amount, senderChainID, coinType, asset)
	if err != nil {
		return nil, false, err
	}

	// handle the deposit for protocol contract version 2
	if protocolContractVersion == crosschaintypes.ProtocolContractVersion_V2 {
		return k.ProcessV2Deposit(ctx, from, senderChainID, zrc20Contract, to, amount, message, coinType)
	}

	// check if the receiver is a contract
	// if it is, then the hook onCrossChainCall() will be called
	// if not, the zrc20 are simply transferred to the receiver
	acc := k.evmKeeper.GetAccount(ctx, to)
	if acc != nil && acc.IsContract() {
		context := systemcontract.ZContext{
			Origin:  from,
			Sender:  eth.Address{},
			ChainID: big.NewInt(senderChainID),
		}
		res, err := k.DepositZRC20AndCallContract(ctx, context, zrc20Contract, to, amount, message)
		return res, true, err
	}

	// if the account is a EOC, no contract call can be made with the data
	if len(message) > 0 {
		return nil, false, types.ErrCallNonContract
	}

	res, err := k.DepositZRC20(ctx, zrc20Contract, to, amount)
	return res, false, err
}

// getAndCheckZRC20 returns the ZRC20 contract address and foreign coin for the given chainID and asset
// it also checks if the foreign coin is paused and if the cap is reached
func (k Keeper) getAndCheckZRC20(
	ctx sdk.Context,
	amount *big.Int,
	chainID int64,
	coinType coin.CoinType,
	asset string,
) (eth.Address, types.ForeignCoins, error) {
	var zrc20Contract eth.Address
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
			return eth.Address{}, types.ForeignCoins{}, crosschaintypes.ErrGasCoinNotFound
		}
	} else {
		foreignCoin, found = k.GetForeignCoinFromAsset(ctx, asset, chainID)
		if !found {
			return eth.Address{}, types.ForeignCoins{}, crosschaintypes.ErrForeignCoinNotFound
		}
	}
	zrc20Contract = eth.HexToAddress(foreignCoin.Zrc20ContractAddress)

	// check if foreign coin is paused
	if foreignCoin.Paused {
		return eth.Address{}, types.ForeignCoins{}, types.ErrPausedZRC20
	}

	// check foreign coins cap if it has a cap
	if !foreignCoin.LiquidityCap.IsNil() && !foreignCoin.LiquidityCap.IsZero() {
		liquidityCap := foreignCoin.LiquidityCap.BigInt()
		totalSupply, err := k.TotalSupplyZRC4(ctx, zrc20Contract)
		if err != nil {
			return eth.Address{}, types.ForeignCoins{}, err
		}
		newSupply := new(big.Int).Add(totalSupply, amount)
		if newSupply.Cmp(liquidityCap) > 0 {
			return eth.Address{}, types.ForeignCoins{}, types.ErrForeignCoinCapReached
		}
	}

	return zrc20Contract, foreignCoin, nil
}
