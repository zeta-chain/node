package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/revert.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/zrc20.sol"

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
		nil,
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
	if crypto.IsEmptyAddress(gatewayAddr) {
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
//
//	function executeRevert(
//	address target,
//	RevertContext revertContext,
//	)
func (k Keeper) CallExecuteRevert(
	ctx sdk.Context,
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
		nil,
		true,
		false,
		"executeRevert",
		target,
		revert.RevertContext{
			Asset:         zrc20,
			Amount:        amount.Uint64(),
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
		nil,
		true,
		false,
		"depositAndRevert",
		zrc20,
		amount,
		target,
		revert.RevertContext{
			Asset:         zrc20,
			Amount:        amount.Uint64(),
			RevertMessage: message,
		},
	)
}
