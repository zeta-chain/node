package keeper

import (
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/wzeta.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// SetSystemContract set system contract in the store
func (k Keeper) SetSystemContract(ctx sdk.Context, sytemContract types.SystemContract) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SystemContractKey))
	b := k.cdc.MustMarshal(&sytemContract)
	store.Set([]byte{0}, b)
}

// GetSystemContract returns system contract from the store
func (k Keeper) GetSystemContract(ctx sdk.Context) (val types.SystemContract, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SystemContractKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveSystemContract removes system contract from the store
func (k Keeper) RemoveSystemContract(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SystemContractKey))
	store.Delete([]byte{0})
}

// GetSystemContractAddress returns the system contract address
// TODO : wzetaContractAddress and other constant strings , can be declared as a constant string in types
// TODO Remove repetitive code
func (k *Keeper) GetSystemContractAddress(ctx sdk.Context) (ethcommon.Address, error) {
	// set the system contract
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrStateVariableNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	return systemAddress, nil
}

func (k *Keeper) GetWZetaContractAddress(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrStateVariableNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, _ := systemcontract.SystemContractMetaData.GetAbi()

	res, err := k.CallEVM(ctx, *sysABI, types.ModuleAddressEVM, systemAddress, BigIntZero, nil, false, false, "wZetaContractAddress")
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to call wZetaContractAddress")
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var wzetaResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&wzetaResponse, "wZetaContractAddress", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrABIUnpack, "failed to unpack wZetaContractAddress: %s", err.Error())
	}
	return wzetaResponse.Value, nil
}

func (k *Keeper) GetUniswapV2FactoryAddress(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrStateVariableNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, _ := systemcontract.SystemContractMetaData.GetAbi()

	res, err := k.CallEVM(ctx, *sysABI, types.ModuleAddressEVM, systemAddress, BigIntZero, nil, false, false, "uniswapv2FactoryAddress")
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to call uniswapv2FactoryAddress")
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var wzetaResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&wzetaResponse, "uniswapv2FactoryAddress", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrABIUnpack, "failed to unpack uniswapv2FactoryAddress: %s", err.Error())
	}
	return wzetaResponse.Value, nil
}

func (k *Keeper) GetUniswapV2Router02Address(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrStateVariableNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, _ := systemcontract.SystemContractMetaData.GetAbi()

	res, err := k.CallEVM(ctx, *sysABI, types.ModuleAddressEVM, systemAddress, BigIntZero, nil, false, false, "uniswapv2Router02Address")
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to call uniswapv2Router02Address")
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var routerResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&routerResponse, "uniswapv2Router02Address", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrABIUnpack, "failed to unpack uniswapv2Router02Address: %s", err.Error())
	}
	return routerResponse.Value, nil
}

func (k *Keeper) CallWZetaDeposit(ctx sdk.Context, sender ethcommon.Address, amount *big.Int) error {
	wzetaAddress, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to get wzeta contract address")
	}
	abi, err := wzeta.WETH9MetaData.GetAbi()
	if err != nil {
		return err
	}
	gasLimit := big.NewInt(70_000) // for some reason, GasEstimate for this contract call is always insufficient
	_, err = k.CallEVM(ctx, *abi, sender, wzetaAddress, amount, gasLimit, true, false, "deposit")
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to call wzeta deposit")
	}
	return nil
}

func (k *Keeper) QueryWZetaBalanceOf(ctx sdk.Context, addr ethcommon.Address) (*big.Int, error) {
	wzetaAddress, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get wzeta contract address")
	}
	wzetaABI, err := wzeta.WETH9MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get ABI")
	}
	res, err := k.CallEVM(ctx, *wzetaABI, addr, wzetaAddress, big.NewInt(0), nil, false, false, "balanceOf", addr)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to call balanceOf")
	}
	type BigIntResponse struct {
		Value *big.Int
	}
	var balanceResponse BigIntResponse
	if err := wzetaABI.UnpackIntoInterface(&balanceResponse, "balanceOf", res.Ret); err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrABIUnpack, "failed to unpack balanceOf: %s", err.Error())
	}
	return balanceResponse.Value, nil
}

func (k *Keeper) QuerySystemContractGasCoinZRC20(ctx sdk.Context, chainid *big.Int) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrStateVariableNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, _ := systemcontract.SystemContractMetaData.GetAbi()

	res, err := k.CallEVM(ctx, *sysABI, types.ModuleAddressEVM, systemAddress, BigIntZero, nil, false, false, "gasCoinZRC20ByChainId", chainid)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to call gasCoinZRC20ByChainId")
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var zrc20Res AddressResponse
	if err := sysABI.UnpackIntoInterface(&zrc20Res, "gasCoinZRC20ByChainId", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrABIUnpack, "failed to unpack gasCoinZRC20ByChainId: %s", err.Error())
	}
	return zrc20Res.Value, nil
}

// returns the amount [in, out]
func (k *Keeper) CallUniswapv2RouterSwapExactETHForToken(ctx sdk.Context, sender ethcommon.Address,
	to ethcommon.Address, amountIn *big.Int, outZRC4 ethcommon.Address, noEthereumTxEvent bool) ([]*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wzeta, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}
	//function swapExactETHForTokens(uint amountOutMin, address[] calldata path, address to, uint deadline)  external payable
	//returns (uint[] memory amounts);
	res, err := k.CallEVM(ctx, *routerABI, sender, routerAddress, amountIn, big.NewInt(300_000), true, noEthereumTxEvent,
		"swapExactETHForTokens", BigIntZero, []ethcommon.Address{wzeta, outZRC4}, to, big.NewInt(1e17))
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to CallEVM method swapExactETHForTokens")
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapExactETHForTokens", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to UnpackIntoInterface swapExactETHForTokens")
	}
	return (*amounts)[:], nil
}

func (k *Keeper) CallUniswapv2RouterSwapEthForExactToken(ctx sdk.Context, sender ethcommon.Address, to ethcommon.Address, maxAmountIn *big.Int, amountOut *big.Int, outZRC4 ethcommon.Address) ([]*big.Int, error) {

	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wzeta, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}
	//function swapETHForExactTokens(uint amountOut, address[] calldata path, address to, uint deadline)
	//returns (uint[] memory amounts);
	res, err := k.CallEVM(ctx, *routerABI, sender, routerAddress, maxAmountIn, big.NewInt(300_000), true, false,
		"swapETHForExactTokens", amountOut, []ethcommon.Address{wzeta, outZRC4}, to, big.NewInt(1e17))
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to CallEVM method swapETHForExactTokens")
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapETHForExactTokens", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack swapETHForExactTokens")
	}
	return (*amounts)[:], nil
}

func (k *Keeper) QueryUniswapv2RouterGetAmountsIn(ctx sdk.Context, amountOut *big.Int, outZRC4 ethcommon.Address) (*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wzeta, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}

	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}
	//function getAmountsIn(uint amountOut, address[] memory path) public view returns (uint[] memory amounts);
	k.Logger(ctx).Info("getAmountsIn", "outZRC20", outZRC4.Hex(), "amountOut", amountOut, "wzeta", wzeta.Hex())
	res, err := k.CallEVM(ctx, *routerABI, types.ModuleAddressEVM, routerAddress, BigIntZero, nil, false, false,
		"getAmountsIn", amountOut, []ethcommon.Address{wzeta, outZRC4})
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to CallEVM method getAmountsIn")
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "getAmountsIn", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack getAmountsIn")
	}
	return (*amounts)[0], nil
}

func (k *Keeper) CallZRC20Burn(ctx sdk.Context, sender ethcommon.Address, zrc20address ethcommon.Address,
	amount *big.Int, noEthereumTxEvent bool) error {
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to get zrc20 abi")
	}
	_, err = k.CallEVM(ctx, *zrc20ABI, sender, zrc20address, big.NewInt(0), big.NewInt(100_000), true, noEthereumTxEvent,
		"burn", amount)
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to CallEVM method burn")
	}
	return nil
}

func (k *Keeper) CallZRC20Deposit(
	ctx sdk.Context,
	sender ethcommon.Address,
	zrc20address ethcommon.Address,
	to ethcommon.Address,
	amount *big.Int) error {
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to get zrc20 abi")
	}
	_, err = k.CallEVM(
		ctx,
		*zrc20ABI,
		sender,
		zrc20address,
		big.NewInt(0),
		big.NewInt(100_000),
		true,
		false,
		"deposit",
		to,
		amount,
	)
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to CallEVM method burn")
	}
	return nil
}
