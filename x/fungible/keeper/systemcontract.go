package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"
)

// TODO : wzetaContractAddress and other constant strings , can be declared as a constant string in types
// TODO Remove repetitive code
func (k *Keeper) GetSystemContractAddress(ctx sdk.Context) (ethcommon.Address, error) {
	// set the system contract
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrStateVaraibleNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	return systemAddress, nil
}

func (k *Keeper) QuerySystemContract(ctx sdk.Context, method string, args ...interface{}) {

}

func (k *Keeper) GetWZetaContractAddress(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrStateVaraibleNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, _ := contracts.SystemContractMetaData.GetAbi()

	res, err := k.CallEVM(ctx, *sysABI, types.ModuleAddressEVM, systemAddress, ZERO_VALUE, nil, false, "wzetaContractAddress")
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to call wzetaContractAddress")
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var wzetaResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&wzetaResponse, "wzetaContractAddress", res.Ret); err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack wzetaContractAddress: %s", err.Error())
	}
	return wzetaResponse.Value, nil
}

func (k *Keeper) GetUniswapv2FacotryAddress(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrStateVaraibleNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, _ := contracts.SystemContractMetaData.GetAbi()

	res, err := k.CallEVM(ctx, *sysABI, types.ModuleAddressEVM, systemAddress, ZERO_VALUE, nil, false, "uniswapv2FactoryAddress")
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to call uniswapv2FactoryAddress")
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var wzetaResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&wzetaResponse, "uniswapv2FactoryAddress", res.Ret); err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack uniswapv2FactoryAddress: %s", err.Error())
	}
	return wzetaResponse.Value, nil
}

func (k *Keeper) GetUniswapV2Router02Address(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrStateVaraibleNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, _ := contracts.SystemContractMetaData.GetAbi()

	res, err := k.CallEVM(ctx, *sysABI, types.ModuleAddressEVM, systemAddress, ZERO_VALUE, nil, false, "uniswapv2Router02Address")
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to call uniswapv2Router02Address")
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var routerResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&routerResponse, "uniswapv2Router02Address", res.Ret); err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack uniswapv2Router02Address: %s", err.Error())
	}
	return routerResponse.Value, nil
}

func (k *Keeper) CallWZetaDeposit(ctx sdk.Context, sender ethcommon.Address, amount *big.Int) error {
	wzetaAddress, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to get wzeta contract address")
	}
	wzetaABI, _ := contracts.WZETAMetaData.GetAbi()
	gasLimit := big.NewInt(70_000) // for some reason, GasEstimate for this contract call is always insufficient
	_, err = k.CallEVM(ctx, *wzetaABI, sender, wzetaAddress, amount, gasLimit, true, "deposit")
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to call wzeta deposit")
	}
	return nil
}

func (k *Keeper) QueryWZetaBalanceOf(ctx sdk.Context, addr ethcommon.Address) (*big.Int, error) {
	wzetaAddress, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to get wzeta contract address")
	}
	wzetaABI, _ := contracts.WZETAMetaData.GetAbi()
	res, err := k.CallEVM(ctx, *wzetaABI, addr, wzetaAddress, big.NewInt(0), nil, false, "balanceOf", addr)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to call balanceOf")
	}
	type BigIntResponse struct {
		Value *big.Int
	}
	var balanceResponse BigIntResponse
	if err := wzetaABI.UnpackIntoInterface(&balanceResponse, "balanceOf", res.Ret); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack balanceOf: %s", err.Error())
	}
	return balanceResponse.Value, nil
}

func (k *Keeper) QuerySystemContractGasCoinZRC4(ctx sdk.Context, chainid *big.Int) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrStateVaraibleNotFound, "failed to get system contract variable")
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, _ := contracts.SystemContractMetaData.GetAbi()

	res, err := k.CallEVM(ctx, *sysABI, types.ModuleAddressEVM, systemAddress, ZERO_VALUE, nil, false, "uniswapv2Router02Address")
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to call uniswapv2Router02Address")
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var routerResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&routerResponse, "uniswapv2Router02Address", res.Ret); err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack uniswapv2Router02Address: %s", err.Error())
	}
	return routerResponse.Value, nil
}

func (k *Keeper) CallUniswapv2RouterSwapExactETHForToken(ctx sdk.Context, sender ethcommon.Address, to ethcommon.Address, amountIn *big.Int, outZRC4 ethcommon.Address) ([]*big.Int, error) {
	routerABI, err := contracts.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to get router abi")
	}
	wzeta, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	//function swapExactETHForTokens(uint amountOutMin, address[] calldata path, address to, uint deadline)  external payable
	//returns (uint[] memory amounts);
	res, err := k.CallEVM(ctx, *routerABI, sender, routerAddress, amountIn, big.NewInt(300_000), true,
		"swapExactETHForTokens", ZERO_VALUE, []ethcommon.Address{wzeta, outZRC4}, to, big.NewInt(1e17))
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to CallEVM method swapExactETHForTokens")
	}

	amounts := new([2]*big.Int)
	routerABI.UnpackIntoInterface(&amounts, "swapExactETHForTokens", res.Ret)
	return (*amounts)[:], nil
}

func (k *Keeper) CallUniswapv2RouterSwapEthForExactToken(ctx sdk.Context, sender ethcommon.Address, to ethcommon.Address, maxAmountIn *big.Int, amountOut *big.Int, outZRC4 ethcommon.Address) ([]*big.Int, error) {

	routerABI, err := contracts.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to get router abi")
	}
	wzeta, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	//function swapETHForExactTokens(uint amountOut, address[] calldata path, address to, uint deadline)
	//returns (uint[] memory amounts);
	res, err := k.CallEVM(ctx, *routerABI, sender, routerAddress, maxAmountIn, big.NewInt(300_000), true,
		"swapETHForExactTokens", amountOut, []ethcommon.Address{wzeta, outZRC4}, to, big.NewInt(1e17))
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to CallEVM method swapETHForExactTokens")
	}

	amounts := new([2]*big.Int)
	routerABI.UnpackIntoInterface(&amounts, "swapETHForExactTokens", res.Ret)
	return (*amounts)[:], nil
}

func (k *Keeper) QueryUniswapv2RouterGetAmountsIn(ctx sdk.Context, amountOut *big.Int, outZRC4 ethcommon.Address) ([]*big.Int, error) {
	routerABI, err := contracts.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to get router abi")
	}
	wzeta, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}

	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	//function getAmountsIn(uint amountOut, address[] memory path) public view returns (uint[] memory amounts);
	res, err := k.CallEVM(ctx, *routerABI, types.ModuleAddressEVM, routerAddress, ZERO_VALUE, nil, false,
		"getAmountsIn", amountOut, []ethcommon.Address{wzeta, outZRC4})
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to CallEVM method getAmountsIn")
	}

	amounts := new([2]*big.Int)
	routerABI.UnpackIntoInterface(&amounts, "getAmountsIn", res.Ret)
	return (*amounts)[:], nil
}
