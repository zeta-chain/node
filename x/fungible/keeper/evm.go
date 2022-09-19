package keeper

import (
	"encoding/json"
	"fmt"
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
	contracts "github.com/zeta-chain/zetacore/contracts/evm"
)

// DeployERC20Contract creates and deploys an ERC20 contract on the EVM with the
// erc20 module account as owner.
func (k Keeper) DeployZRC4Contract(
	ctx sdk.Context,
	name, symbol string,
	decimals uint8,
	chain string,
	coinType zetacommon.CoinType,
	erc20Contract string,
) (common.Address, error) { // FIXME: geneeralized beyond ETH
	abi, err := contracts.ZRC4MetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get ZRC4 ABI: %s", err.Error())
	}
	ctorArgs, err := abi.Pack(
		"",                     // function--empty string for constructor
		name,                   // name
		symbol,                 // symbol
		decimals,               // decimals
		types.ModuleAddressEVM, // owner
	)

	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "coin metadata is invalid %s: %s", name, err.Error())
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

	coinIndex := fmt.Sprintf("%s-%s", chain, name)
	coin, _ := k.GetForeignCoins(ctx, coinIndex)
	coin.CoinType = coinType
	coin.Name = name
	coin.Symbol = symbol
	coin.Decimals = uint32(decimals)
	coin.ERC20ContractAddress = erc20Contract
	coin.ZRC4ContractAddress = contractAddr.String()
	coin.Index = coinIndex
	coin.ForeignChain = chain
	k.SetForeignCoins(ctx, coin)

	// update ZetaDepositAndCall system contract addr
	depositCaller, _ := k.GetZetaDepositAndCallContract(ctx)
	if len(depositCaller.Address) != 0 {
		ZDCAddr := common.HexToAddress(depositCaller.Address)
		_, err = k.CallEVM(ctx, *abi, types.ModuleAddressEVM, contractAddr, true, "updateZetaDepositAndCallAddress", ZDCAddr)
		if err != nil {
			return common.Address{}, sdkerrors.Wrapf(err, "failed to update ZetaDepositAndCall contract address for %s", name)
		}
	}

	return contractAddr, nil
}

// Deploy the ZetaDepositAndCall system contract
func (k Keeper) DeployZetaDepositAndCall(ctx sdk.Context) (common.Address, error) {
	abi, err := contracts.ZetaDepositAndCallMetaData.GetAbi()
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIGet, "failed to get ZetaDepositAndCallMetaData ABI: %s", err.Error())
	}
	ctorArgs, err := abi.Pack(
		"",                     // function--empty string for constructor
		types.ModuleAddressEVM, // owner: fungible module
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

	depositCaller, _ := k.GetZetaDepositAndCallContract(ctx)
	depositCaller.Address = contractAddr.String()
	k.SetZetaDepositAndCallContract(ctx, depositCaller)

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
	to common.Address,
	amount *big.Int,
	contract common.Address,
	from common.Address,
	data []byte) (*evmtypes.MsgEthereumTxResponse, error) {

	abi, err := contracts.ZRC4MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, zrc4Contract, true, "deposit", to, amount)
	if err != nil {
		return nil, err
	}
	res, err = k.CallEVMWithData(ctx, from, &contract, data, true)
	if err != nil {
		return nil, err
	}
	return res, err
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
