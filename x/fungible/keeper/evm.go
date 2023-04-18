package keeper

import (
	"encoding/json"
	"math/big"
	"strconv"

	tmtypes "github.com/tendermint/tendermint/types"
	connectorzevm "github.com/zeta-chain/protocol/pkg/contracts/zevm/ConnectorZEVM.sol"
	systemcontract "github.com/zeta-chain/protocol/pkg/contracts/zevm/SystemContract.sol"
	zrc20 "github.com/zeta-chain/protocol/pkg/contracts/zevm/ZRC20.sol"
	uniswapv2factory "github.com/zeta-chain/protocol/pkg/uniswap/v2-core/contracts/UniswapV2Factory.sol"
	uniswapv2router02 "github.com/zeta-chain/protocol/pkg/uniswap/v2-periphery/contracts/UniswapV2Router02.sol"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/server/config"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	zetacommon "github.com/zeta-chain/zetacore/common"
)

// TODO USE string constant
var (
	BigIntZero                 = big.NewInt(0)
	ZEVMGasLimitDepositAndCall = big.NewInt(1_000_000)
)

// TODO Unit test for these funtions
// TODO Remove repetitive code
// DeployERC20Contract creates and deploys an ERC20 contract on the EVM with the
// erc20 module account as owner. Also adds itself to ForeignCoins fungible module state variable
func (k Keeper) DeployZRC20Contract(
	ctx sdk.Context,
	name, symbol string,
	decimals uint8,
	chainStr string,
	coinType zetacommon.CoinType,
	erc20Contract string,
	gasLimit *big.Int,
) (common.Address, error) {
	chainName := zetacommon.ParseChainName(chainStr)
	chain := k.zetaobserverKeeper.GetParams(ctx).GetChainFromChainName(chainName)
	if chain == nil {
		return common.Address{}, sdkerrors.Wrapf(zetaObserverTypes.ErrSupportedChains, "chain %s not found", chainStr)
	}
	abi, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get ZRC4 ABI: %s", err.Error())
	}
	system, found := k.GetSystemContract(ctx)
	if !found {
		return common.Address{}, sdkerrors.Wrapf(types.ErrSystemContractNotFound, "system contract not found")
	}
	ctorArgs, err := abi.Pack(
		"",                        // function--empty string for constructor
		name,                      // name
		symbol,                    // symbol
		decimals,                  // decimals
		big.NewInt(chain.ChainId), // chainID
		uint8(coinType),           // coinType: 0: Zeta 1: gas 2 ERC20
		gasLimit,                  //gas limit for transfer; 21k for gas asset; around 70k for ERC20
		common.HexToAddress(system.SystemContract),
	)

	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "coin constructor is invalid %s: %s", name, err.Error())
	}
	data := make([]byte, len(zrc20.ZRC20MetaData.Bin)+len(ctorArgs))
	copy(data[:len(zrc20.ZRC20MetaData.Bin)], zrc20.ZRC20MetaData.Bin)
	copy(data[len(zrc20.ZRC20MetaData.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true, BigIntZero, nil)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy contract for %s", name)
	}

	coinIndex := name
	coin, _ := k.GetForeignCoins(ctx, chain.ChainId, coinIndex)
	coin.CoinType = coinType
	coin.Name = name
	coin.Symbol = symbol
	coin.Decimals = uint32(decimals)
	coin.Asset = erc20Contract
	coin.Zrc20ContractAddress = contractAddr.String()
	coin.Index = coinIndex
	coin.ForeignChainId = chain.ChainId
	k.SetForeignCoins(ctx, coin)

	return contractAddr, nil
}

func (k Keeper) DeploySystemContract(ctx sdk.Context, wzeta common.Address, v2factory common.Address, router02 common.Address) (common.Address, error) {
	abi, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get SystemContract ABI: %s", err.Error())
	}
	ctorArgs, err := abi.Pack(
		"", // function--empty string for constructor,
		wzeta,
		v2factory,
		router02,
	)

	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "error packing SystemContract constructor arguments: %s", err.Error())
	}

	data := make([]byte, len(systemcontract.SystemContractMetaData.Bin)+len(ctorArgs))
	copy(data[:len(systemcontract.SystemContractMetaData.Bin)], systemcontract.SystemContractMetaData.Bin)
	copy(data[len(systemcontract.SystemContractMetaData.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true, BigIntZero, nil)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy SystemContractContract system contract")
	}

	system, _ := k.GetSystemContract(ctx)
	//system := types.SystemContract{}
	system.SystemContract = contractAddr.String()
	k.SetSystemContract(ctx, system)

	// go update all addr on ZRC-4 contracts

	// TODO : Change to
	// GET all supported chains
	// Get all coins for al chains
	//zrc4ABI, err := zrc20.ZRC20MetaData.GetAbi()
	//coins := k.GetAllForeignCoins(ctx)
	//for _, coin := range coins {
	//	if len(coin.Zrc20ContractAddress) != 0 {
	//		zrc4Address := common.HexToAddress(coin.Zrc20ContractAddress)
	//		_, err = k.CallEVM(ctx, *zrc4ABI, types.ModuleAddressEVM, zrc4Address, BigIntZero, nil, true, "updateSystemContractAddress", contractAddr)
	//		if err != nil {
	//			k.Logger(ctx).Error("failed to update updateSystemContractAddress contract address for %s: %s", coin.Name, contractAddr, err.Error())
	//		}
	//	}
	//}

	return contractAddr, nil
}

func (k Keeper) DeployUniswapV2Factory(ctx sdk.Context) (common.Address, error) {
	abi, err := uniswapv2factory.UniswapV2FactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get UniswapV2FactoryMetaData ABI: %s", err.Error())
	}
	ctorArgs, err := abi.Pack(
		"",                     // function--empty string for constructor
		types.ModuleAddressEVM, // feeToSetter
	)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "error packing UniswapV2Factory constructor arguments: %s", err.Error())
	}

	data := make([]byte, len(uniswapv2factory.UniswapV2FactoryMetaData.Bin)+len(ctorArgs))
	copy(data[:len(uniswapv2factory.UniswapV2FactoryMetaData.Bin)], uniswapv2factory.UniswapV2FactoryMetaData.Bin)
	copy(data[len(uniswapv2factory.UniswapV2FactoryMetaData.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true, BigIntZero, nil)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy UniswapV2FactoryContract contract")
	}

	//verify that factory is correct--hashOfPairCode must be: 96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f
	// this is important because router02 needs exactly this build to compute correct pair address
	// Name
	_, err = k.CallEVM(ctx, *abi, types.ModuleAddressEVM, contractAddr, BigIntZero, nil, false, "hashOfPairCode")
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to call hashOfPairCode() contract")
	}

	return contractAddr, nil
}

func (k Keeper) DeployUniswapV2Router02(ctx sdk.Context, factory common.Address, wzeta common.Address) (common.Address, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get UniswapV2Router02MetaData ABI: %s", err.Error())
	}
	ctorArgs, err := routerABI.Pack(
		"", // function--empty string for constructor
		factory, wzeta)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "error packing UniswapV2Router02 constructor arguments: %s", err.Error())
	}

	data := make([]byte, len(uniswapv2router02.UniswapV2Router02Contract.Bin)+len(ctorArgs))
	copy(data[:len(uniswapv2router02.UniswapV2Router02MetaData.Bin)], uniswapv2router02.UniswapV2Router02MetaData.Bin)
	copy(data[len(uniswapv2router02.UniswapV2Router02MetaData.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true, BigIntZero, nil)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy UniswapV2Router02Contract contract")
	}

	return contractAddr, nil
}

func (k Keeper) DeployWZETA(ctx sdk.Context) (common.Address, error) {
	abi, err := connectorzevm.WZETAMetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get WZETAMetaData ABI: %s", err.Error())
	}
	ctorArgs, err := abi.Pack(
		"", // function--empty string for constructor
	)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "error packing WZETA constructor arguments: %s", err.Error())
	}

	data := make([]byte, len(connectorzevm.WZETAMetaData.Bin)+len(ctorArgs))
	copy(data[:len(connectorzevm.WZETAMetaData.Bin)], connectorzevm.WZETAMetaData.Bin)
	copy(data[len(connectorzevm.WZETAMetaData.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true, BigIntZero, nil)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy WZETA contract")
	}

	return contractAddr, nil
}

func (k Keeper) DeployConnectorZEVM(ctx sdk.Context, wzeta common.Address) (common.Address, error) {
	abi, err := connectorzevm.ZetaConnectorZEVMMetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get ZetaConnectorZEVMMetaData ABI: %s", err.Error())
	}
	ctorArgs, err := abi.Pack(
		"", // function--empty string for constructor
		wzeta,
	)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "error packing ZetaConnectorZEVM constructor arguments: %s", err.Error())
	}

	data := make([]byte, len(connectorzevm.ZetaConnectorZEVMMetaData.Bin)+len(ctorArgs))
	copy(data[:len(connectorzevm.ZetaConnectorZEVMMetaData.Bin)], connectorzevm.ZetaConnectorZEVMMetaData.Bin)
	copy(data[len(connectorzevm.ZetaConnectorZEVMMetaData.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true, BigIntZero, nil)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy WZETA contract")
	}

	return contractAddr, nil
}

// Depoisit ZRC4 tokens into to account;
// Callable only by the fungible module account
func (k Keeper) DepositZRC20(
	ctx sdk.Context,
	contract common.Address,
	to common.Address,
	amount *big.Int,
) (*evmtypes.MsgEthereumTxResponse, error) {
	//abi, err := zrc20.ZRC4MetaData.GetAbi()
	abi, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, contract, BigIntZero, nil, true, "deposit", to, amount)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Deposit into ZRC4 and call contract function in a single tx
// callable from fungible module
func (k Keeper) DepositZRC20AndCallContract(ctx sdk.Context,
	zrc4Contract common.Address,
	targetContract common.Address,
	amount *big.Int,
	message []byte) (*evmtypes.MsgEthereumTxResponse, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrContractNotFound, "GetSystemContract address not found")
	}
	systemAddress := common.HexToAddress(system.SystemContract)

	abi, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, systemAddress, BigIntZero, ZEVMGasLimitDepositAndCall, true,
		"depositAndCall", zrc4Contract, amount, targetContract, message)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// QueryZRC4Data returns the data of a deployed ZRC4 contract
func (k Keeper) QueryZRC4Data(
	ctx sdk.Context,
	contract common.Address,
) (types.ZRC20Data, error) {
	var (
		nameRes    types.ZRC20StringResponse
		symbolRes  types.ZRC20StringResponse
		decimalRes types.ZRC20Uint8Response
	)

	zrc4 := zrc20.ZRC20Contract.ABI

	// Name
	res, err := k.CallEVM(ctx, zrc4, types.ModuleAddressEVM, contract, BigIntZero, nil, false, "name")
	if err != nil {
		return types.ZRC20Data{}, err
	}

	if err := zrc4.UnpackIntoInterface(&nameRes, "name", res.Ret); err != nil {
		return types.ZRC20Data{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack name: %s", err.Error(),
		)
	}

	// Symbol
	res, err = k.CallEVM(ctx, zrc4, types.ModuleAddressEVM, contract, BigIntZero, nil, false, "symbol")
	if err != nil {
		return types.ZRC20Data{}, err
	}

	if err := zrc4.UnpackIntoInterface(&symbolRes, "symbol", res.Ret); err != nil {
		return types.ZRC20Data{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack symbol: %s", err.Error(),
		)
	}

	// Decimals
	res, err = k.CallEVM(ctx, zrc4, types.ModuleAddressEVM, contract, BigIntZero, nil, false, "decimals")
	if err != nil {
		return types.ZRC20Data{}, err
	}

	if err := zrc4.UnpackIntoInterface(&decimalRes, "decimals", res.Ret); err != nil {
		return types.ZRC20Data{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack decimals: %s", err.Error(),
		)
	}

	return types.NewZRC20Data(nameRes.Value, symbolRes.Value, decimalRes.Value), nil
}

// BalanceOfZRC4 queries an account's balance for a given ZRC4 contract
func (k Keeper) BalanceOfZRC4(
	ctx sdk.Context,
	contract, account common.Address,
) *big.Int {
	abi, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return nil
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, contract, BigIntZero, nil, false, "balanceOf", account)
	if err != nil {
		return nil
	}
	// TODO :  return the error here, we loose the error message if we return a nil . Maube use (big.Int,error )
	unpacked, err := abi.Unpack("balanceOf", res.Ret)
	if err != nil || len(unpacked) == 0 {
		return nil
	}
	// TODO :  return the error here, we loose the error message if we return a nil . Maube use (big.Int,error )

	balance, ok := unpacked[0].(*big.Int)
	if !ok {
		return nil
	}

	return balance
}

// CallEVM performs a smart contract method call using given args
func (k Keeper) CallEVM(
	ctx sdk.Context,
	abi abi.ABI,
	from, contract common.Address,
	value *big.Int,
	gasLimit *big.Int,
	commit bool,
	method string,
	args ...interface{},
) (*evmtypes.MsgEthereumTxResponse, error) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			types.ErrABIPack,
			sdkerrors.Wrap(err, "failed to create transaction data").Error(),
		)
	}

	k.Logger(ctx).Info("calling EVM", "from", from, "contract", contract, "value", value, "method", method)
	resp, err := k.CallEVMWithData(ctx, from, &contract, data, commit, value, gasLimit)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
	}
	return resp, nil
}

// CallEVMWithData performs a smart contract method call using contract data
// value is the amount of wei to send; gaslimit is the custom gas limit, if nil EstimateGas is used
// to bisect the correct gas limit (this may sometimes results in insufficient gas limit; not sure why)
func (k Keeper) CallEVMWithData(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	data []byte,
	commit bool,
	value *big.Int,
	gasLimit *big.Int,
) (*evmtypes.MsgEthereumTxResponse, error) {
	nonce, err := k.authKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}
	gasCap := config.DefaultGasCap
	if commit && gasLimit == nil {
		args, err := json.Marshal(evmtypes.TransactionArgs{
			From: &from,
			To:   contract,
			Data: (*hexutil.Bytes)(&data),
		})
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "failed to marshal tx args: %s", err.Error())
		}

		gasRes, err := k.evmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evmtypes.EthCallRequest{
			Args:   args,
			GasCap: config.DefaultGasCap,
		})
		if err != nil {
			return nil, err
		}
		gasCap = gasRes.Gas
		k.Logger(ctx).Info("call evm", "EstimateGas", gasCap)
	}
	if gasLimit != nil {
		gasCap = gasLimit.Uint64()
	}

	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		value,         // amount
		gasCap,        // gasLimit
		big.NewInt(0), // gasFeeCap
		big.NewInt(0), // gasTipCap
		big.NewInt(0), // gasPrice
		data,
		ethtypes.AccessList{}, // AccessList
		!commit,               // isFake
	)
	k.evmKeeper.WithChainID(ctx) //FIXME:  set chainID for signer; should not need to do this; but seems necessary. Why?
	k.Logger(ctx).Info("call evm", "gasCap", gasCap, "chainid", k.evmKeeper.ChainID(), "ctx.chainid", ctx.ChainID())
	res, err := k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), commit)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	msgBytes, _ := json.Marshal(msg)
	ethTxHash := common.BytesToHash(crypto.Keccak256(msgBytes)) // NOTE(pwu): this is a fake txhash
	attrs := []sdk.Attribute{}
	if len(ctx.TxBytes()) > 0 {
		// add event for tendermint transaction hash format
		hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
		ethTxHash = common.BytesToHash(hash) // NOTE(pwu): use cosmos tx hash as eth tx hash if available
		attrs = append(attrs, sdk.NewAttribute(evmtypes.AttributeKeyTxHash, hash.String()))
	}
	attrs = append(attrs, []sdk.Attribute{
		sdk.NewAttribute(sdk.AttributeKeyAmount, value.String()),
		// add event for ethereum transaction hash format; NOTE(pwu): this is a fake txhash
		sdk.NewAttribute(evmtypes.AttributeKeyEthereumTxHash, ethTxHash.String()),
		// add event for index of valid ethereum tx; NOTE(pwu): fake txindex
		sdk.NewAttribute(evmtypes.AttributeKeyTxIndex, strconv.FormatUint(8888, 10)),
		// add event for eth tx gas used, we can't get it from cosmos tx result when it contains multiple eth tx msgs.
		sdk.NewAttribute(evmtypes.AttributeKeyTxGasUsed, strconv.FormatUint(res.GasUsed, 10)),
	}...)

	// receipient: contract address
	if contract != nil {
		attrs = append(attrs, sdk.NewAttribute(evmtypes.AttributeKeyRecipient, contract.Hex()))
	}
	if res.Failed() {
		attrs = append(attrs, sdk.NewAttribute(evmtypes.AttributeKeyEthereumTxFailed, res.VmError))
	}

	txLogAttrs := make([]sdk.Attribute, len(res.Logs))
	for i, log := range res.Logs {
		log.TxHash = ethTxHash.String()
		value, err := json.Marshal(log)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to encode log")
		}
		txLogAttrs[i] = sdk.NewAttribute(evmtypes.AttributeKeyTxLog, string(value))
	}

	// emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			evmtypes.EventTypeEthereumTx,
			attrs...,
		),
		sdk.NewEvent(
			evmtypes.EventTypeTxLog,
			txLogAttrs...,
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, from.Hex()),
			sdk.NewAttribute(evmtypes.AttributeKeyTxType, "88"), // type 88: synthetic Eth tx
		),
	})

	logs := evmtypes.LogsToEthereum(res.Logs)
	var bloomReceipt ethtypes.Bloom
	if len(logs) > 0 {
		bloom := k.evmKeeper.GetBlockBloomTransient(ctx)
		bloom.Or(bloom, big.NewInt(0).SetBytes(ethtypes.LogsBloom(logs)))
		bloomReceipt = ethtypes.BytesToBloom(bloom.Bytes())
		k.evmKeeper.SetBlockBloomTransient(ctx, bloomReceipt.Big())
		k.evmKeeper.SetLogSizeTransient(ctx, (k.evmKeeper.GetLogSizeTransient(ctx))+uint64(len(logs)))
	}

	return res, nil
}
