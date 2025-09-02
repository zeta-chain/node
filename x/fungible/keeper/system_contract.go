package keeper

import (
	"fmt"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/wzeta.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	"github.com/zeta-chain/node/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/node/x/fungible/types"
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
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	return systemAddress, nil
}

// GetWZetaContractAddress returns the wzeta contract address on ZetaChain
func (k *Keeper) GetWZetaContractAddress(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}

	res, err := k.CallEVM(
		ctx,
		*sysABI,
		types.ModuleAddressEVM,
		systemAddress,
		BigIntZero,
		nil,
		false,
		false,
		"wZetaContractAddress",
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to call wZetaContractAddress (%s)",
			err.Error(),
		)
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var wzetaResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&wzetaResponse, "wZetaContractAddress", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrABIUnpack,
			"failed to unpack wZetaContractAddress: %s",
			err.Error(),
		)
	}

	if wzetaResponse.Value == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrContractNotFound, "wzeta contract invalid address")
	}
	return wzetaResponse.Value, nil
}

// GetUniswapV2FactoryAddress returns the uniswapv2 factory contract address on ZetaChain
func (k *Keeper) GetUniswapV2FactoryAddress(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}

	res, err := k.CallEVM(
		ctx,
		*sysABI,
		types.ModuleAddressEVM,
		systemAddress,
		BigIntZero,
		nil,
		false,
		false,
		"uniswapv2FactoryAddress",
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to call uniswapv2FactoryAddress (%s)",
			err.Error(),
		)
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var uniswapFactoryResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&uniswapFactoryResponse, "uniswapv2FactoryAddress", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrABIUnpack,
			"failed to unpack uniswapv2FactoryAddress: %s",
			err.Error(),
		)
	}

	if uniswapFactoryResponse.Value == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractNotFound,
			"uniswap factory contract invalid address",
		)
	}
	return uniswapFactoryResponse.Value, nil
}

// GetUniswapV2Router02Address returns the uniswapv2 router02 address on ZetaChain
func (k *Keeper) GetUniswapV2Router02Address(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}

	res, err := k.CallEVM(
		ctx,
		*sysABI,
		types.ModuleAddressEVM,
		systemAddress,
		BigIntZero,
		nil,
		false,
		false,
		"uniswapv2Router02Address",
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to call uniswapv2Router02Address (%s)",
			err.Error(),
		)
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var routerResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&routerResponse, "uniswapv2Router02Address", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrABIUnpack,
			"failed to unpack uniswapv2Router02Address: %s",
			err.Error(),
		)
	}

	if routerResponse.Value == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractNotFound,
			"uniswap router contract invalid address",
		)
	}
	return routerResponse.Value, nil
}

// CallWZetaDeposit calls the deposit method of the wzeta contract
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

	_, err = k.CallEVM(
		ctx,
		*abi,
		sender,
		wzetaAddress,
		amount,
		gasLimit,
		true,
		false,
		"deposit",
	)
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrContractCall, "failed to call wzeta deposit (%s)", err.Error())
	}
	return nil
}

// QueryWZetaBalanceOf returns the balance of the given address in the wzeta contract
func (k *Keeper) QueryWZetaBalanceOf(ctx sdk.Context, addr ethcommon.Address) (*big.Int, error) {
	wzetaAddress, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get wzeta contract address")
	}

	wzetaABI, err := wzeta.WETH9MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get ABI")
	}

	res, err := k.CallEVM(
		ctx,
		*wzetaABI,
		addr,
		wzetaAddress,
		big.NewInt(0),
		nil,
		false,
		false,
		"balanceOf",
		addr,
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to call balanceOf (%s)", err.Error())
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

// QuerySystemContractGasCoinZRC20 returns the gas coin zrc20 address for the given chain id
func (k *Keeper) QuerySystemContractGasCoinZRC20(ctx sdk.Context, chainid *big.Int) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}

	res, err := k.CallEVM(
		ctx,
		*sysABI,
		types.ModuleAddressEVM,
		systemAddress,
		BigIntZero,
		nil,
		false,
		false,
		"gasCoinZRC20ByChainId",
		chainid,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to call gasCoinZRC20ByChainId (%s)",
			err.Error(),
		)
	}

	type AddressResponse struct {
		Value ethcommon.Address
	}
	var zrc20Res AddressResponse
	if err := sysABI.UnpackIntoInterface(&zrc20Res, "gasCoinZRC20ByChainId", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrABIUnpack,
			"failed to unpack gasCoinZRC20ByChainId: %s",
			err.Error(),
		)
	}
	if zrc20Res.Value == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrContractNotFound, "gas coin contract invalid address")
	}
	return zrc20Res.Value, nil
}

// CallUniswapV2RouterSwapExactTokensForTokens calls the swapExactTokensForETH method of the uniswapv2 router contract
// to swap tokens to another tokens using wZeta as intermediary
func (k *Keeper) CallUniswapV2RouterSwapExactTokensForTokens(
	ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	amountIn *big.Int,
	inZRC4,
	outZRC4 ethcommon.Address,
	noEthereumTxEvent bool,
) (ret []*big.Int, err error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wzetaAddr, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function swapExactTokensForTokens(
	//	uint amountIn,
	//	uint amountOutMin,
	//	address[] calldata path,
	//	address to,
	//	uint deadline
	//)
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		sender,
		routerAddress,
		BigIntZero,
		big.NewInt(1_000_000),
		true,
		noEthereumTxEvent,
		"swapExactTokensForTokens",
		amountIn,
		BigIntZero,
		[]ethcommon.Address{inZRC4, wzetaAddr, outZRC4},
		to,
		big.NewInt(1e17),
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to CallEVM method swapExactTokensForTokens (%s)",
			err.Error(),
		)
	}

	amounts := new([3]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapExactTokensForTokens", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to UnpackIntoInterface swapExactTokensForTokens")
	}
	return (*amounts)[:], nil
}

// CallUniswapV2RouterSwapExactTokensForETH calls the swapExactTokensForETH method of the uniswapv2 router contract
func (k *Keeper) CallUniswapV2RouterSwapExactTokensForETH(
	ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	amountIn *big.Int,
	inZRC4 ethcommon.Address,
	noEthereumTxEvent bool,
) (ret []*big.Int, err error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wzetaAddr, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function swapExactTokensForETH(
	//	uint amountIn,
	//	uint amountOutMin,
	//	address[] calldata path,
	//	address to,
	//	uint deadline
	//)
	ctx.Logger().Error("Calling swapExactTokensForETH")
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		sender,
		routerAddress,
		BigIntZero,
		big.NewInt(300_000),
		true,
		noEthereumTxEvent,
		"swapExactTokensForETH",
		amountIn,
		BigIntZero,
		[]ethcommon.Address{inZRC4, wzetaAddr},
		to,
		big.NewInt(1e17),
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to CallEVM method swapExactTokensForETH (%s)",
			err.Error(),
		)
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapExactTokensForETH", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to UnpackIntoInterface swapExactTokensForETH")
	}
	return (*amounts)[:], nil
}

// CallUniswapV2RouterSwapExactETHForToken calls the swapExactETHForTokens method of the uniswapv2 router contract
func (k *Keeper) CallUniswapV2RouterSwapExactETHForToken(
	ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	amountIn *big.Int,
	outZRC4 ethcommon.Address,
	noEthereumTxEvent bool,
) ([]*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}

	wzetaAddr, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function swapExactETHForTokens(uint amountOutMin, address[] calldata path, address to, uint deadline)  external payable
	//returns (uint[] memory amounts);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		sender,
		routerAddress,
		amountIn,
		big.NewInt(300_000),
		true,
		noEthereumTxEvent,
		"swapExactETHForTokens",
		BigIntZero,
		[]ethcommon.Address{wzetaAddr, outZRC4},
		to,
		big.NewInt(1e17),
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to CallEVM method swapExactETHForTokens (%s)",
			err.Error(),
		)
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapExactETHForTokens", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to UnpackIntoInterface swapExactETHForTokens")
	}
	return (*amounts)[:], nil
}

// CallUniswapV2RouterSwapEthForExactToken calls the swapETHForExactTokens method of the uniswapv2 router contract
func (k *Keeper) CallUniswapV2RouterSwapEthForExactToken(
	ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	maxAmountIn *big.Int,
	amountOut *big.Int,
	outZRC4 ethcommon.Address,
) ([]*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wzetaAddr, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function swapETHForExactTokens(uint amountOut, address[] calldata path, address to, uint deadline)
	//returns (uint[] memory amounts);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		sender,
		routerAddress,
		maxAmountIn,
		big.NewInt(300_000),
		true,
		false,
		"swapETHForExactTokens",
		amountOut,
		[]ethcommon.Address{wzetaAddr, outZRC4},
		to,
		big.NewInt(1e17),
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to CallEVM method swapETHForExactTokens (%s)",
			err.Error(),
		)
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapETHForExactTokens", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack swapETHForExactTokens")
	}
	return (*amounts)[:], nil
}

// QueryUniswapV2RouterGetZetaAmountsIn returns the amount of zeta needed to buy the given amount of ZRC4 tokens
func (k *Keeper) QueryUniswapV2RouterGetZetaAmountsIn(
	ctx sdk.Context,
	amountOut *big.Int,
	outZRC4 ethcommon.Address,
) (*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wzetaAddr, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function getAmountsIn(uint amountOut, address[] memory path) public view returns (uint[] memory amounts);
	k.Logger(ctx).Info("getAmountsIn", "outZRC20", outZRC4.Hex(), "amountOut", amountOut, "wzeta", wzetaAddr.Hex())
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		BigIntZero,
		nil,
		false,
		false,
		"getAmountsIn",
		amountOut,
		[]ethcommon.Address{wzetaAddr, outZRC4},
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, cosmoserrors.Wrap(
			types.ErrContractCall,
			fmt.Sprintf("failed to CallEVM method getAmountsIn (%s)", err.Error()),
		)
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "getAmountsIn", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack getAmountsIn")
	}
	return (*amounts)[0], nil
}

// QueryUniswapV2RouterGetZRC4AmountsIn returns the amount of ZRC4 tokens needed to buy the given amount of zeta
func (k *Keeper) QueryUniswapV2RouterGetZRC4AmountsIn(
	ctx sdk.Context,
	amountOut *big.Int,
	inZRC4 ethcommon.Address,
) (*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wzetaAddr, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function getAmountsIn(uint amountOut, address[] memory path) public view returns (uint[] memory amounts);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		BigIntZero,
		nil,
		false,
		false,
		"getAmountsIn",
		amountOut,
		[]ethcommon.Address{inZRC4, wzetaAddr},
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method getAmountsIn (%s)", err.Error())
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "getAmountsIn", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack getAmountsIn")
	}
	return (*amounts)[0], nil
}

// QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn returns the amount of ZRC4 tokens needed to buy another ZRC4 token, it uses the WZeta contract as a bridge
func (k *Keeper) QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(
	ctx sdk.Context,
	amountOut *big.Int,
	inZRC4, outZRC4 ethcommon.Address,
) (*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wzetaAddr, err := k.GetWZetaContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWZetaContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function getAmountsIn(uint amountOut, address[] memory path) public view returns (uint[] memory amounts);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		BigIntZero,
		nil,
		false,
		false,
		"getAmountsIn",
		amountOut,
		[]ethcommon.Address{inZRC4, wzetaAddr, outZRC4},
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method getAmountsIn (%s)", err.Error())
	}

	amounts := new([3]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "getAmountsIn", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack getAmountsIn")
	}
	return (*amounts)[0], nil
}

// CallZRC20Burn calls the burn method of the zrc20 contract
func (k *Keeper) CallZRC20Burn(
	ctx sdk.Context,
	sender ethcommon.Address,
	zrc20address ethcommon.Address,
	amount *big.Int,
	noEthereumTxEvent bool,
) error {
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
		DefaultGasLimit,
		true,
		noEthereumTxEvent,
		"burn",
		amount,
	)
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method burn (%s)", err.Error())
	}

	return nil
}

// CallZRC20Deposit calls the deposit method of the zrc20 contract
func (k *Keeper) CallZRC20Deposit(
	ctx sdk.Context,
	sender ethcommon.Address,
	zrc20address ethcommon.Address,
	to ethcommon.Address,
	amount *big.Int,
) error {
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
		DefaultGasLimit,
		true,
		false,
		"deposit",
		to,
		amount,
	)
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method burn (%s)", err.Error())
	}
	return nil
}

// CallZRC20Approve calls the approve method of the zrc20 contract
func (k *Keeper) CallZRC20Approve(
	ctx sdk.Context,
	owner ethcommon.Address,
	zrc20address ethcommon.Address,
	spender ethcommon.Address,
	amount *big.Int,
	noEthereumTxEvent bool,
) error {
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to get zrc20 abi")
	}

	_, err = k.CallEVM(
		ctx,
		*zrc20ABI,
		owner,
		zrc20address,
		BigIntZero,
		DefaultGasLimit,
		true,
		noEthereumTxEvent,
		"approve",
		spender,
		amount,
	)
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method approve (%s)", err.Error())
	}

	return nil
}

func (k *Keeper) GetGatewayGasLimit(ctx sdk.Context) (*big.Int, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return &big.Int{}, types.ErrSystemContractNotFound
	}

	return system.GatewayGasLimit.BigInt(), nil
}

func (k *Keeper) MustGetGatewayGasLimit(ctx sdk.Context) *big.Int {
	gasLimit, err := k.GetGatewayGasLimit(ctx)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("failed to get gateway gas limit, using default %s", err.Error()))
		return types.GatewayGasLimit
	}
	return gasLimit
}

func (k *Keeper) SetGatewayGasLimit(ctx sdk.Context, gasLimit sdkmath.Int) error {
	system := types.SystemContract{}
	existingSystemContract, found := k.GetSystemContract(ctx)
	if found {
		system = existingSystemContract
	}
	system.GatewayGasLimit = gasLimit
	k.SetSystemContract(ctx, system)
	return nil
}
