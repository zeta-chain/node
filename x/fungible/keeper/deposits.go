package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k Keeper) DepositCoinZeta(ctx sdk.Context, to eth.Address, amount *big.Int) error {
	zetaToAddress := sdk.AccAddress(to.Bytes())
	return k.MintZetaToEVMAccount(ctx, zetaToAddress, amount)
}

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
	var Zrc20Contract eth.Address
	var coin fungibletypes.ForeignCoins
	var found bool
	if coinType == common.CoinType_Gas {
		coin, found = k.GetGasCoinForForeignCoin(ctx, senderChain.ChainId)
		if !found {
			return nil, types.ErrGasCoinNotFound
		}
	} else {
		coin, found = k.GetForeignCoinFromAsset(ctx, asset, senderChain.ChainId)
		if !found {
			return nil, types.ErrForeignCoinNotFound
		}
	}
	Zrc20Contract = eth.HexToAddress(coin.Zrc20ContractAddress)
	if len(message) == 0 { // no message; transfer
		return k.DepositZRC20(ctx, Zrc20Contract, to, amount)
	}
	// non-empty message = [contractaddress, calldata]
	if len(data) == 0 {
		return k.DepositZRC20(ctx, Zrc20Contract, contract, amount)
	}
	context := systemcontract.ZContext{
		Origin:  from,
		Sender:  eth.Address{},
		ChainID: big.NewInt(senderChain.ChainId),
	}
	return k.DepositZRC20AndCallContract(ctx, context, Zrc20Contract, contract, amount, data)

}
