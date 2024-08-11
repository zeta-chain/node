package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/systemcontract.sol"

	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// CallDepositAndCallZRC20 calls the depositAndCall (ZRC20 version) function on the gateway contract
// Callable only by the fungible module account
// returns directly CallEVM()
// function depositAndCall(
//
//	    zContext calldata context,
//	    address zrc20,
//	    uint256 amount,
//	    address target,
//	    bytes calldata message
//	)
func (k Keeper) CallDepositAndCallZRC20(
	ctx sdk.Context,
	context systemcontract.ZContext,
	zrc20 common.Address,
	amount *big.Int,
	target common.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	gatewayABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	systemContract, found := k.GetSystemContract(ctx)
	if !found {
		return nil, types.ErrSystemContractNotFound
	}
	gatewayAddr := common.HexToAddress(systemContract.Gateway)
	if gatewayAddr == (common.Address{}) {
		return nil, types.ErrGatewayContractNotSet
	}

	// NOTE:
	// depositAndCall: ZETA version for depositAndCall method
	// depositAndCall0: ZRC20 version for depositAndCall method
	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		nil,
		true,
		false,
		"depositAndCall0",
		context,
		zrc20,
		amount,
		target,
		message,
	)
}

// CallExecute calls the execute function on the gateway contract
// function execute(
//
//	zContext calldata context,
//	address zrc20,
//	uint256 amount,
//	address target,
//	bytes calldata message
//
// )
func (k Keeper) CallExecute(
	ctx sdk.Context,
	context systemcontract.ZContext,
	zrc20 common.Address,
	amount *big.Int,
	target common.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	gatewayABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	systemContract, found := k.GetSystemContract(ctx)
	if !found {
		return nil, types.ErrSystemContractNotFound
	}
	gatewayAddr := common.HexToAddress(systemContract.Gateway)
	if gatewayAddr == (common.Address{}) {
		return nil, types.ErrGatewayContractNotSet
	}

	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		nil,
		true,
		false,
		"execute",
		context,
		zrc20,
		amount,
		target,
		message,
	)
}

// CallExecuteRevert calls the executeRevert function on the gateway contract
// function executeRevert(
//
//	revertContext calldata context,
//	address zrc20,
//	uint256 amount,
//	address target,
//	bytes calldata message
//
// )
// TODO: complete once interface finalized
// https://github.com/zeta-chain/node/issues/2660
func (k Keeper) CallExecuteRevert(
	ctx sdk.Context,
	context systemcontract.ZContext,
	zrc20 common.Address,
	amount *big.Int,
	target common.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	gatewayABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	systemContract, found := k.GetSystemContract(ctx)
	if !found {
		return nil, types.ErrSystemContractNotFound
	}
	gatewayAddr := common.HexToAddress(systemContract.Gateway)
	if gatewayAddr == (common.Address{}) {
		return nil, types.ErrGatewayContractNotSet
	}

	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		nil,
		true,
		false,
		"executeRevert",
		context,
		zrc20,
		amount,
		target,
		message,
	)
}

// CallDepositAndRevert calls the depositAndRevert function on the gateway contract
// function depositAndRevert(
//
//	revertContext calldata context,
//	address zrc20,
//	uint256 amount,
//	address target,
//	bytes calldata message
//
// )
func (k Keeper) CallDepositAndRevert(
	ctx sdk.Context,
	context systemcontract.ZContext,
	zrc20 common.Address,
	amount *big.Int,
	target common.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	gatewayABI, err := gatewayzevm.GatewayZEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	systemContract, found := k.GetSystemContract(ctx)
	if !found {
		return nil, types.ErrSystemContractNotFound
	}
	gatewayAddr := common.HexToAddress(systemContract.Gateway)
	if gatewayAddr == (common.Address{}) {
		return nil, types.ErrGatewayContractNotSet
	}

	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		nil,
		true,
		false,
		"depositAndRevert",
		context,
		zrc20,
		amount,
		target,
		message,
	)
}
