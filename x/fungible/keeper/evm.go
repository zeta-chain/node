package keeper

import (
	"encoding/json"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/server/config"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	zetacommon "github.com/zeta-chain/zetacore/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	clientconfig "github.com/zeta-chain/zetacore/zetaclient/config"
)

// DeployERC20Contract creates and deploys an ERC20 contract on the EVM with the
// erc20 module account as owner. Also adds itself to ForeignCoins fungible module state variable
func (k Keeper) DeployZRC4Contract(
	ctx sdk.Context,
	name, symbol string,
	decimals uint8,
	chainStr string,
	coinType zetacommon.CoinType,
	erc20Contract string,
	gasLimit *big.Int,
) (common.Address, error) { // FIXME: geneeralized beyond ETH
	abi, err := contracts.ZRC4MetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get ZRC4 ABI: %s", err.Error())
	}
	chain, found := clientconfig.Chains[chainStr]
	if !found {
		return common.Address{}, sdkerrors.Wrapf(types.ErrChainNotFound, "chain %s not found", chainStr)
	}

	system, found := k.GetSystemContract(ctx)
	if !found {
		return common.Address{}, sdkerrors.Wrapf(types.ErrSystemContractNotFound, "system contract not found")
	}
	ctorArgs, err := abi.Pack(
		"",              // function--empty string for constructor
		name,            // name
		symbol,          // symbol
		decimals,        // decimals
		chain.ChainID,   // chainID
		uint8(coinType), // coinType: 0: Zeta 1: gas 2 ERC20
		gasLimit,        //gas limit for transfer; 21k for gas asset; around 70k for ERC20
		common.HexToAddress(system.ZetaDepositAndCallContract),
		common.HexToAddress(system.GasPriceOracleContract),
	)

	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "coin constructor is invalid %s: %s", name, err.Error())
	}

	data := make([]byte, len(contracts.ZRC4Contract.Bin)+len(ctorArgs))
	copy(data[:len(contracts.ZRC4Contract.Bin)], contracts.ZRC4Contract.Bin)
	copy(data[len(contracts.ZRC4Contract.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy contract for %s", name)
	}

	coinIndex := name
	coin, _ := k.GetForeignCoins(ctx, coinIndex)
	coin.CoinType = coinType
	coin.Name = name
	coin.Symbol = symbol
	coin.Decimals = uint32(decimals)
	coin.ERC20ContractAddress = erc20Contract
	coin.ZRC4ContractAddress = contractAddr.String()
	coin.Index = coinIndex
	coin.ForeignChain = chainStr
	k.SetForeignCoins(ctx, coin)

	return contractAddr, nil
}

// Deploy the ZetaDepositAndCall system contract
func (k Keeper) DeployZetaDepositAndCall(ctx sdk.Context) (common.Address, error) {
	abi, err := contracts.ZetaDepositAndCallMetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get ZetaDepositAndCallMetaData ABI: %s", err.Error())
	}
	ctorArgs, err := abi.Pack(
		"", // function--empty string for constructor
	)

	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "error packing ZetaDepositAndCallMetaData constructor arguments: %s", err.Error())
	}

	data := make([]byte, len(contracts.ZetaDepositAndCallContract.Bin)+len(ctorArgs))
	copy(data[:len(contracts.ZetaDepositAndCallContract.Bin)], contracts.ZetaDepositAndCallContract.Bin)
	copy(data[len(contracts.ZetaDepositAndCallContract.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy ZetaDepositAndCall system contract")
	}

	system, _ := k.GetSystemContract(ctx)
	system.ZetaDepositAndCallContract = contractAddr.String()
	k.SetSystemContract(ctx, system)

	// go update all addr on ZRC-4 contracts
	zrc4ABI, err := contracts.ZRC4MetaData.GetAbi()
	coins := k.GetAllForeignCoins(ctx)
	for _, coin := range coins {
		if len(coin.ZRC4ContractAddress) != 0 {
			zrc4Address := common.HexToAddress(coin.ZRC4ContractAddress)
			_, err = k.CallEVM(ctx, *zrc4ABI, types.ModuleAddressEVM, zrc4Address, true, "updateZetaDepositAndCallAddress", contractAddr)
			if err != nil {
				return common.Address{}, sdkerrors.Wrapf(err, "failed to update ZetaDepositAndCall contract address for %s", coin.Name)
			}
		}
	}

	return contractAddr, nil
}

// Deploy the DeployGasPriceOracle system contract
func (k Keeper) DeployGasPriceOracleContract(ctx sdk.Context) (common.Address, error) {
	abi, err := contracts.GasPriceOracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get GasPriceOracleMetaData ABI: %s", err.Error())
	}
	ctorArgs, err := abi.Pack(
		"", // function--empty string for constructor
	)

	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "error packing GasPriceOracleMetaData constructor arguments: %s", err.Error())
	}

	data := make([]byte, len(contracts.GasPriceOracleContract.Bin)+len(ctorArgs))
	copy(data[:len(contracts.GasPriceOracleContract.Bin)], contracts.GasPriceOracleContract.Bin)
	copy(data[len(contracts.GasPriceOracleContract.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	res, err := k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy GasPriceOracleMetaData system contract")
	}
	k.Logger(ctx).Info("DeployGasPriceOracleContract", "res.log", res.Logs)

	system, _ := k.GetSystemContract(ctx)
	system.GasPriceOracleContract = contractAddr.String()
	k.SetSystemContract(ctx, system)

	// go update all addr on ZRC-4 contracts
	zrc4ABI, err := contracts.ZRC4MetaData.GetAbi()
	coins := k.GetAllForeignCoins(ctx)
	for _, coin := range coins {
		if len(coin.ZRC4ContractAddress) != 0 {
			zrc4Address := common.HexToAddress(coin.ZRC4ContractAddress)
			_, err = k.CallEVM(ctx, *zrc4ABI, types.ModuleAddressEVM, zrc4Address, true, "updateGasPriceOracleAddress", contractAddr)
			if err != nil {
				//return common.Address{}, sdkerrors.Wrapf(err, "failed to update GasPriceOracleAddress contract address for %s", coin.Name)
				k.Logger(ctx).Error("failed to update GasPriceOracleAddress contract address for %s", coin.Name)
			}
		}
	}

	return contractAddr, nil
}

func (k Keeper) DeployUniswapV2Factory(ctx sdk.Context) (common.Address, error) {
	abi, err := contracts.UniswapV2FactoryMetaData.GetAbi()
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

	data := make([]byte, len(contracts.UniswapV2FactoryContract.Bin)+len(ctorArgs))
	copy(data[:len(contracts.UniswapV2FactoryContract.Bin)], contracts.UniswapV2FactoryContract.Bin)
	copy(data[len(contracts.UniswapV2FactoryContract.Bin):], ctorArgs)

	nonce, err := k.authKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddressEVM, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddressEVM, nil, data, true)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy UniswapV2FactoryContract contract")
	}

	//verify that factory is correct--hashOfPairCode must be: 96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f
	// this is important because router02 needs exactly this build to compute correct pair address
	// Name
	_, err = k.CallEVM(ctx, *abi, types.ModuleAddressEVM, contractAddr, false, "hashOfPairCode")
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to call hashOfPairCode() contract")
	}

	return contractAddr, nil
}

// Depoisit ZRC4 tokens into to account;
// Callable only by the fungible module account
func (k Keeper) DepositZRC4(
	ctx sdk.Context,
	contract common.Address,
	to common.Address,
	amount *big.Int,
) (*evmtypes.MsgEthereumTxResponse, error) {
	//abi, err := contracts.ZRC4MetaData.GetAbi()
	abi, err := contracts.ZRC4MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, contract, true, "deposit", to, amount)
	if err != nil {
		return nil, err
	}

	return res, err
}

// Deposit into ZRC4 and call contract function in a single tx
// callable from fungible module
func (k Keeper) DepositZRC4AndCallContract(ctx sdk.Context,
	zrc4Contract common.Address,
	amount *big.Int,
	targetContract common.Address,
	message []byte) (*evmtypes.MsgEthereumTxResponse, error) {

	system, found := k.GetSystemContract(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrContractNotFound, "ZetaDepositAndCall address not found")
	}
	ZDCAddress := common.HexToAddress(system.ZetaDepositAndCallContract)

	abi, err := contracts.ZetaDepositAndCallMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, ZDCAddress, true,
		"DepositAndCall", zrc4Contract, amount, targetContract, message)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// QueryZRC4Data returns the data of a deployed ZRC4 contract
func (k Keeper) QueryZRC4Data(
	ctx sdk.Context,
	contract common.Address,
) (types.ZRC4Data, error) {
	var (
		nameRes    types.ZRC4StringResponse
		symbolRes  types.ZRC4StringResponse
		decimalRes types.ZRC4Uint8Response
	)

	zrc4 := contracts.ZRC4Contract.ABI

	// Name
	res, err := k.CallEVM(ctx, zrc4, types.ModuleAddressEVM, contract, false, "name")
	if err != nil {
		return types.ZRC4Data{}, err
	}

	if err := zrc4.UnpackIntoInterface(&nameRes, "name", res.Ret); err != nil {
		return types.ZRC4Data{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack name: %s", err.Error(),
		)
	}

	// Symbol
	res, err = k.CallEVM(ctx, zrc4, types.ModuleAddressEVM, contract, false, "symbol")
	if err != nil {
		return types.ZRC4Data{}, err
	}

	if err := zrc4.UnpackIntoInterface(&symbolRes, "symbol", res.Ret); err != nil {
		return types.ZRC4Data{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack symbol: %s", err.Error(),
		)
	}

	// Decimals
	res, err = k.CallEVM(ctx, zrc4, types.ModuleAddressEVM, contract, false, "decimals")
	if err != nil {
		return types.ZRC4Data{}, err
	}

	if err := zrc4.UnpackIntoInterface(&decimalRes, "decimals", res.Ret); err != nil {
		return types.ZRC4Data{}, sdkerrors.Wrapf(
			types.ErrABIUnpack, "failed to unpack decimals: %s", err.Error(),
		)
	}

	return types.NewZRC4Data(nameRes.Value, symbolRes.Value, decimalRes.Value), nil
}

// BalanceOfZRC4 queries an account's balance for a given ZRC4 contract
func (k Keeper) BalanceOfZRC4(
	ctx sdk.Context,
	contract, account common.Address,
) *big.Int {
	abi, err := contracts.ZRC4MetaData.GetAbi()
	if err != nil {
		return nil
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, contract, false, "balanceOf", account)
	if err != nil {
		return nil
	}

	unpacked, err := abi.Unpack("balanceOf", res.Ret)
	if err != nil || len(unpacked) == 0 {
		return nil
	}

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

	resp, err := k.CallEVMWithData(ctx, from, &contract, data, commit)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
	}
	return resp, nil
}

// CallEVMWithData performs a smart contract method call using contract data
func (k Keeper) CallEVMWithData(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	data []byte,
	commit bool,
) (*evmtypes.MsgEthereumTxResponse, error) {
	nonce, err := k.authKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	gasCap := config.DefaultGasCap
	if commit {
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
	}

	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		big.NewInt(0), // amount
		gasCap,        // gasLimit
		big.NewInt(0), // gasFeeCap
		big.NewInt(0), // gasTipCap
		big.NewInt(0), // gasPrice
		data,
		ethtypes.AccessList{}, // AccessList
		!commit,               // isFake
	)

	res, err := k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), commit)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	return res, nil
}

// monitorApprovalEvent returns an error if the given transactions logs include
// an unexpected `Approval` event
func (k Keeper) monitorApprovalEvent(res *evmtypes.MsgEthereumTxResponse) error {
	if res == nil || len(res.Logs) == 0 {
		return nil
	}

	logApprovalSig := []byte("Approval(address,address,uint256)")
	logApprovalSigHash := crypto.Keccak256Hash(logApprovalSig)

	for _, log := range res.Logs {
		if log.Topics[0] == logApprovalSigHash.Hex() {
			return sdkerrors.Wrapf(
				types.ErrUnexpectedEvent, "unexpected Approval event",
			)
		}
	}

	return nil
}
