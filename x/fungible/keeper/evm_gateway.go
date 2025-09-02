package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/revert.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zrc20.sol"

	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/x/fungible/types"
)

// CallUpdateGatewayAddress calls the updateGatewayAddress function on the ZRC20 contract
// function updateGatewayAddress(address addr)
func (k Keeper) CallUpdateGatewayAddress(
	ctx sdk.Context,
	zrc20Address common.Address,
	newGatewayAddress common.Address,
) (*evmtypes.MsgEthereumTxResponse, error) {
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return k.CallEVM(
		ctx,
		*zrc20ABI,
		types.ModuleAddressEVM,
		zrc20Address,
		BigIntZero,
		k.MustGetGatewayGasLimit(ctx),
		true,
		false,
		"updateGatewayAddress",
		newGatewayAddress,
	)
}

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
	context gatewayzevm.MessageContext,
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
	if crypto.IsEmptyAddress(gatewayAddr) {
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
		k.MustGetGatewayGasLimit(ctx),
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

// DepositAndCallZeta calls the depositAndCall function on the gateway contract
// function depositAndCall(
// MessageContext calldata context,
// address target,
// bytes calldata message
// )
func (k Keeper) DepositAndCallZeta(
	ctx sdk.Context,
	context gatewayzevm.MessageContext,
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
	if crypto.IsEmptyAddress(gatewayAddr) {
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
		amount,
		k.MustGetGatewayGasLimit(ctx),
		true,
		false,
		"depositAndCall",
		context,
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
	context gatewayzevm.MessageContext,
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
	if crypto.IsEmptyAddress(gatewayAddr) {
		return nil, types.ErrGatewayContractNotSet
	}

	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		k.MustGetGatewayGasLimit(ctx),
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
//
//	function executeRevert(
//	address target,
//	RevertContext revertContext,
//	)
func (k Keeper) CallExecuteRevert(
	ctx sdk.Context,
	inboundSender string,
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
	if crypto.IsEmptyAddress(gatewayAddr) {
		return nil, types.ErrGatewayContractNotSet
	}

	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		k.MustGetGatewayGasLimit(ctx),
		true,
		false,
		"executeRevert",
		target,
		revert.RevertContext{
			Sender:        common.HexToAddress(inboundSender),
			Asset:         zrc20,
			Amount:        amount,
			RevertMessage: message,
		},
	)
}

// CallDepositAndRevert calls the depositAndRevert function on the gateway contract
//
// function depositAndRevert(
//
//	address zrc20,
//	uint256 amount,
//	address target,
//	RevertContext revertContext
//
// )
func (k Keeper) CallDepositAndRevert(
	ctx sdk.Context,
	inboundSender string,
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
	if crypto.IsEmptyAddress(gatewayAddr) {
		return nil, types.ErrGatewayContractNotSet
	}

	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		k.MustGetGatewayGasLimit(ctx),
		true,
		false,
		"depositAndRevert",
		zrc20,
		amount,
		target,
		revert.RevertContext{
			Sender:        common.HexToAddress(inboundSender),
			Asset:         zrc20,
			Amount:        amount,
			RevertMessage: message,
		},
	)
}

// CallExecuteAbort calls the executeAbort function on the gateway contract
//
//	function executeAbort(
//	address target,
//	AbortContext abortContext,
//	)
func (k Keeper) CallExecuteAbort(
	ctx sdk.Context,
	inboundSender string,
	zrc20 common.Address,
	amount *big.Int,
	outgoing bool,
	chainID *big.Int,
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
	if crypto.IsEmptyAddress(gatewayAddr) {
		return nil, types.ErrGatewayContractNotSet
	}

	// TODO: set correct sender for non EVM chains
	// https://github.com/zeta-chain/node/issues/3532
	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		k.MustGetGatewayGasLimit(ctx),
		true,
		false,
		"executeAbort",
		target,
		revert.AbortContext{
			Sender:        common.HexToAddress(inboundSender).Bytes(),
			Asset:         zrc20,
			Amount:        amount,
			Outgoing:      outgoing,
			ChainID:       chainID,
			RevertMessage: message,
		},
	)
}
