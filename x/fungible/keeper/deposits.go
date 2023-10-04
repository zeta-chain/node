package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/zetacore/common"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// DepositCoinZeta immediately mints ZETA to the EVM account
func (k Keeper) DepositCoinZeta(ctx sdk.Context, to eth.Address, amount *big.Int) error {
	zetaToAddress := sdk.AccAddress(to.Bytes())
	return k.MintZetaToEVMAccount(ctx, zetaToAddress, amount)
}

// ZRC20DepositAndCallContract deposits ZRC20 to the EVM account and calls the contract
func (k Keeper) ZRC20DepositAndCallContract(
	ctx sdk.Context,
	from []byte,
	to eth.Address,
	amount *big.Int,
	senderChain *common.Chain,
	data []byte,
	coinType common.CoinType,
	asset string,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	var ZRC20Contract eth.Address
	var coin types.ForeignCoins
	var found bool

	// get foreign coin
	if coinType == common.CoinType_Gas {
		coin, found = k.GetGasCoinForForeignCoin(ctx, senderChain.ChainId)
		if !found {
			return nil, false, crosschaintypes.ErrGasCoinNotFound
		}
	} else {
		coin, found = k.GetForeignCoinFromAsset(ctx, asset, senderChain.ChainId)
		if !found {
			return nil, false, crosschaintypes.ErrForeignCoinNotFound
		}
	}
	ZRC20Contract = eth.HexToAddress(coin.Zrc20ContractAddress)

	// check foreign coins cap if it has a cap
	if !coin.LiquidityCap.IsNil() && !coin.LiquidityCap.IsZero() {
		liquidityCap := coin.LiquidityCap.BigInt()
		totalSupply, err := k.TotalSupplyZRC4(ctx, ZRC20Contract)
		if err != nil {
			return nil, false, err
		}
		newSupply := new(big.Int).Add(totalSupply, amount)
		if newSupply.Cmp(liquidityCap) > 0 {
			return nil, false, types.ErrForeignCoinCapReached
		}
	}

	// check if the receiver is a contract
	// if it is, then the hook onCrossChainCall() will be called
	// if not, the zrc20 are simply transferred to the receiver
	acc := k.evmKeeper.GetAccount(ctx, to)
	if acc != nil && acc.IsContract() {
		context := systemcontract.ZContext{
			Origin:  from,
			Sender:  eth.Address{},
			ChainID: big.NewInt(senderChain.ChainId),
		}
		res, err := k.DepositZRC20AndCallContract(ctx, context, ZRC20Contract, to, amount, data)
		return res, true, err
	}

	// if the account is a EOC, no contract call can be made with the data
	if len(data) > 0 {
		return nil, false, types.ErrCallNonContract
	}

	res, err := k.DepositZRC20(ctx, ZRC20Contract, to, amount)
	return res, false, err
}
