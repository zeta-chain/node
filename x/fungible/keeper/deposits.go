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
	message string,
	contract eth.Address,
	data []byte,
	coinType common.CoinType,
	asset string,
) (*evmtypes.MsgEthereumTxResponse, error) {
	var ZRC20Contract eth.Address
	var coin types.ForeignCoins
	var found bool

	// get foreign coin
	if coinType == common.CoinType_Gas {
		coin, found = k.GetGasCoinForForeignCoin(ctx, senderChain.ChainId)
		if !found {
			return nil, crosschaintypes.ErrGasCoinNotFound
		}
	} else {
		coin, found = k.GetForeignCoinFromAsset(ctx, asset, senderChain.ChainId)
		if !found {
			return nil, crosschaintypes.ErrForeignCoinNotFound
		}
	}

	ZRC20Contract = eth.HexToAddress(coin.Zrc20ContractAddress)

	// check foreign coins cap
	totalSupply, err := k.TotalSupplyZRC4(ctx, ZRC20Contract)
	if err != nil {
		return nil, err
	}
	coinCap := big.NewInt(1000) //TODOHERE: get coin cap from system contract
	newSupply := new(big.Int).Add(totalSupply, amount)
	if newSupply.Cmp(coinCap) > 0 {
		return nil, types.ErrForeignCoinCapReached
	}

	// no message: transfer to EVM account
	if len(message) == 0 {
		return k.DepositZRC20(ctx, ZRC20Contract, to, amount)
	}

	// non-empty message with empty data: deposit to contract
	if len(data) == 0 {
		return k.DepositZRC20(ctx, ZRC20Contract, contract, amount)
	}

	// non-empty message with non-empty data: contract call
	context := systemcontract.ZContext{
		Origin:  from,
		Sender:  eth.Address{},
		ChainID: big.NewInt(senderChain.ChainId),
	}
	return k.DepositZRC20AndCallContract(ctx, context, ZRC20Contract, contract, amount, data)

}
