package staking

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	precompiletypes "github.com/zeta-chain/node/precompiles/types"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
)

var (
	ABI                 abi.ABI
	ContractAddress     = common.HexToAddress("0x0000000000000000000000000000000000000066")
	GasRequiredByMethod = map[[4]byte]uint64{}
	ViewMethod          = map[[4]byte]bool{}
)

func init() {
	initABI()
}

func initABI() {
	if err := ABI.UnmarshalJSON([]byte(IStakingMetaData.ABI)); err != nil {
		panic(err)
	}

	GasRequiredByMethod = map[[4]byte]uint64{}
	for methodName := range ABI.Methods {
		var methodID [4]byte
		copy(methodID[:], ABI.Methods[methodName].ID[:4])
		switch methodName {
		// TODO: https://github.com/zeta-chain/node/issues/2812
		// just temporary flat values, double check these flat values
		// can we just use WriteCostFlat/ReadCostFlat from gas config for flat values?
		case StakeMethodName:
			GasRequiredByMethod[methodID] = StakeMethodGas
		case UnstakeMethodName:
			GasRequiredByMethod[methodID] = UnstakeMethodGas
		case MoveStakeMethodName:
			GasRequiredByMethod[methodID] = MoveStakeMethodGas
		case DistributeMethodName:
			GasRequiredByMethod[methodID] = DistributeMethodGas
		case ClaimRewardsMethodName:
			GasRequiredByMethod[methodID] = ClaimRewardsMethodGas
		case GetAllValidatorsMethodName:
			GasRequiredByMethod[methodID] = 0
			ViewMethod[methodID] = true
		case GetSharesMethodName:
			GasRequiredByMethod[methodID] = 0
			ViewMethod[methodID] = true
		case GetRewardsMethodName:
			GasRequiredByMethod[methodID] = 0
			ViewMethod[methodID] = true
		case GetValidatorsMethodName:
			GasRequiredByMethod[methodID] = 0
			ViewMethod[methodID] = true
		default:
			GasRequiredByMethod[methodID] = 0
		}
	}
}

type Contract struct {
	precompiletypes.BaseContract

	stakingKeeper      stakingkeeper.Keeper
	fungibleKeeper     fungiblekeeper.Keeper
	bankKeeper         bankkeeper.Keeper
	distributionKeeper distrkeeper.Keeper
	cdc                codec.Codec
	kvGasConfig        storetypes.GasConfig
}

func NewIStakingContract(
	ctx sdk.Context,
	stakingKeeper *stakingkeeper.Keeper,
	fungibleKeeper fungiblekeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	distributionKeeper distrkeeper.Keeper,
	cdc codec.Codec,
	kvGasConfig storetypes.GasConfig,
) *Contract {
	accAddress := sdk.AccAddress(ContractAddress.Bytes())
	if !fungibleKeeper.GetAuthKeeper().HasAccount(ctx, accAddress) {
		acc := fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ctx, accAddress)
		fungibleKeeper.GetAuthKeeper().SetAccount(ctx, acc)
	}

	return &Contract{
		BaseContract:       precompiletypes.NewBaseContract(ContractAddress),
		stakingKeeper:      *stakingKeeper,
		fungibleKeeper:     fungibleKeeper,
		bankKeeper:         bankKeeper,
		distributionKeeper: distributionKeeper,
		cdc:                cdc,
		kvGasConfig:        kvGasConfig,
	}
}

// Address() is required to implement the PrecompiledContract interface.
func (c *Contract) Address() common.Address {
	return ContractAddress
}

// Abi() is required to implement the PrecompiledContract interface.
func (c *Contract) Abi() abi.ABI {
	return ABI
}

// RequiredGas is required to implement the PrecompiledContract interface.
// The gas has to be calculated deterministically based on the input.
func (c *Contract) RequiredGas(input []byte) uint64 {
	// get methodID (first 4 bytes)
	var methodID [4]byte
	copy(methodID[:], input[:4])
	// base cost to prevent large input size
	baseCost := uint64(len(input)) * c.kvGasConfig.WriteCostPerByte
	if ViewMethod[methodID] {
		baseCost = uint64(len(input)) * c.kvGasConfig.ReadCostPerByte
	}

	if requiredGas, ok := GasRequiredByMethod[methodID]; ok {
		return requiredGas + baseCost
	}

	// Can not happen, but return 0 if the method is not found.
	return 0
}

// Run is the entrypoint of the precompiled contract, it switches over the input method,
// and execute them accordingly.
func (c *Contract) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) ([]byte, error) {
	method, err := ABI.MethodById(contract.Input[:4])
	if err != nil {
		return nil, err
	}

	args, err := method.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(precompiletypes.ExtStateDB)

	// If the method is not a view method, it should not be executed in read-only mode.
	if _, isViewMethod := ViewMethod[[4]byte(method.ID)]; !isViewMethod && readOnly {
		return nil, precompiletypes.ErrWriteMethod{
			Method: method.Name,
		}
	}

	switch method.Name {
	case GetAllValidatorsMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.GetAllValidators(ctx, method)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case GetSharesMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.GetShares(ctx, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case StakeMethodName:
		// Disabled until further notice, check https://github.com/zeta-chain/node/issues/3005.
		return nil, precompiletypes.ErrDisabledMethod{
			Method: method.Name,
		}

		//nolint:govet
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.Stake(ctx, evm, contract, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case UnstakeMethodName:
		// Disabled until further notice, check https://github.com/zeta-chain/node/issues/3005.
		return nil, precompiletypes.ErrDisabledMethod{
			Method: method.Name,
		}

		//nolint:govet
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.Unstake(ctx, evm, contract, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case MoveStakeMethodName:
		// Disabled until further notice, check https://github.com/zeta-chain/node/issues/3005.
		return nil, precompiletypes.ErrDisabledMethod{
			Method: method.Name,
		}

		//nolint:govet
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.MoveStake(ctx, evm, contract, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case DistributeMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.distribute(ctx, evm, contract, method, args)
			return err
		})
		if execErr != nil {
			res, errPack := method.Outputs.Pack(false)
			if errPack != nil {
				return nil, errPack
			}

			return res, err
		}
		return res, nil
	case GetRewardsMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.getRewards(ctx, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case GetValidatorsMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.getValidatorListForDelegator(ctx, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case ClaimRewardsMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.claimRewards(ctx, evm, contract, method, args)
			return err
		})
		if execErr != nil {
			res, errPack := method.Outputs.Pack(false)
			if errPack != nil {
				return nil, errPack
			}

			return res, err
		}
		return res, nil
	default:
		return nil, precompiletypes.ErrInvalidMethod{
			Method: method.Name,
		}
	}
}
