package bank

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/zeta-chain/ethermint/x/evm/types"

	ptypes "github.com/zeta-chain/node/precompiles/types"
)

const (
	// ZEVM cosmos coins prefix.
	ZEVMDenom = "zevm/"

	// Write methods.
	DepositMethodName  = "deposit"
	WithdrawMethodName = "withdraw"

	// Read methods.
	BalanceOfMethodName = "balanceOf"
)

var (
	ABI                 abi.ABI
	ContractAddress     = common.HexToAddress("0x0000000000000000000000000000000000000067")
	GasRequiredByMethod = map[[4]byte]uint64{}
	ViewMethod          = map[[4]byte]bool{}
)

func init() {
	initABI()
}

func initABI() {
	if err := ABI.UnmarshalJSON([]byte(IBankMetaData.ABI)); err != nil {
		panic(err)
	}

	GasRequiredByMethod = map[[4]byte]uint64{}
	for methodName := range ABI.Methods {
		var methodID [4]byte
		copy(methodID[:], ABI.Methods[methodName].ID[:4])
		switch methodName {
		case DepositMethodName:
			GasRequiredByMethod[methodID] = 200000
		case WithdrawMethodName:
			GasRequiredByMethod[methodID] = 200000
		case BalanceOfMethodName:
			GasRequiredByMethod[methodID] = 10000
		default:
			GasRequiredByMethod[methodID] = 0
		}
	}
}

type Contract struct {
	ptypes.BaseContract

	bankKeeper  bank.Keeper
	cdc         codec.Codec
	kvGasConfig storetypes.GasConfig
}

func NewIBankContract(
	bankKeeper bank.Keeper,
	cdc codec.Codec,
	kvGasConfig storetypes.GasConfig,
) *Contract {
	return &Contract{
		BaseContract: ptypes.NewBaseContract(ContractAddress),
		bankKeeper:   bankKeeper,
		cdc:          cdc,
		kvGasConfig:  kvGasConfig,
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

	stateDB := evm.StateDB.(ptypes.ExtStateDB)

	switch method.Name {
	case DepositMethodName:
		if readOnly {
			return nil, nil
		}

		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.deposit(ctx, method, contract.CallerAddress, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil

	case WithdrawMethodName:
		if readOnly {
			return nil, nil
		}

		return nil, nil
		// TODO

	case BalanceOfMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.balanceOf(ctx, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil

	default:
		return nil, ptypes.ErrInvalidMethod{
			Method: method.Name,
		}
	}
}

func ZRC20ToCosmosDenom(ZRC20Address common.Address) string {
	return ZEVMDenom + ZRC20Address.String()
}

func (c *Contract) deposit(
	ctx sdk.Context,
	method *abi.Method,
	caller common.Address,
	args []interface{},
) (result []byte, err error) {
	if len(args) != 2 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// TODO: The origin tokens have to be:
	// 1. Checked the caller has the right amount of original tokens.
	// 2. burned or locked.
	// Otherwise this deposit functions has the ability to infinite mint coins.

	// function deposit(address zrc20, uint256 amount) external returns (bool success);
	ZRC20Addr, amount := args[0].(common.Address), args[1].(*big.Int)

	// Handle the toAddr:
	// check it's valid and not blocked.
	toAddr := sdk.AccAddress(caller.Bytes())
	if toAddr.Empty() {
		return nil, &ptypes.ErrInvalidAddr{
			Got:    toAddr.String(),
			Reason: "empty address",
		}
	}

	if c.bankKeeper.BlockedAddr(toAddr) {
		return nil, &ptypes.ErrInvalidAddr{
			Got:    toAddr.String(),
			Reason: "blocked by bank keeper",
		}
	}

	// The process of creating a new cosmos coin is:
	// - Generate the new coin denom using ZRC20 address,
	//   this way we map ZRC20 addresses to cosmos denoms "zevm/0x12345".
	// - Mint coins.
	// - Send coins to the caller.
	tokenDenom := ZRC20ToCosmosDenom(ZRC20Addr)
	coin := sdk.NewCoin(tokenDenom, math.NewIntFromBigInt(amount))
	if !coin.IsValid() {
		return nil, &ptypes.ErrInvalidCoin{
			Got:      coin.GetDenom(),
			Negative: coin.IsNegative(),
			Nil:      coin.IsNil(),
		}
	}

	// A sdk.Coins (type []sdk.Coin) has to be created because it's the type expected by MintCoins
	// and SendCoinsFromModuleToAccount.
	// But sdk.Coins will only contain one coin, always.
	coinSet := sdk.NewCoins(coin)
	if !coinSet.IsValid() {
		return nil, &ptypes.ErrInvalidCoin{
			Got:      coinSet.Sort().GetDenomByIndex(0),
			Negative: coinSet.IsAnyNegative(),
			Nil:      coinSet.IsAnyNil(),
		}
	}

	if !c.bankKeeper.IsSendEnabledCoin(ctx, coin) {
		return nil, &ptypes.ErrUnexpected{
			When: "IsSendEnabledCoins",
			Got:  "coin not enabled to be sent",
		}
	}

	if err := c.bankKeeper.MintCoins(ctx, types.ModuleName, coinSet); err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "MintCoins",
			Got:  err.Error(),
		}
	}

	if err := c.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coinSet); err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "SendCoinsFromModuleToAccount",
			Got:  err.Error(),
		}
	}

	return method.Outputs.Pack(true)
}

func (c *Contract) balanceOf(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) (result []byte, err error) {
	if len(args) != 2 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// function balanceOf(address zrc20, address user) external view returns (uint256 balance);
	tokenAddr, addr := args[0].(common.Address), args[1].(common.Address)

	// common.Address has to be converted to AccAddress.
	accAddr := sdk.AccAddress(addr.Bytes())
	if accAddr.Empty() {
		return nil, &ptypes.ErrInvalidAddr{
			Got: accAddr.String(),
		}
	}

	// Convert ZRC20 address to a Cosmos denom formatted as "zevm/0x12345".
	tokenDenom := ZRC20ToCosmosDenom(tokenAddr)

	// Bank Keeper GetBalance returns the specified Cosmos coin balance for a given address.
	// Check explicitly the balance is a non-negative non-nil value.
	coin := c.bankKeeper.GetBalance(ctx, accAddr, tokenDenom)
	if !coin.IsValid() {
		return nil, &ptypes.ErrInvalidCoin{
			Got:      coin.GetDenom(),
			Negative: coin.IsNegative(),
			Nil:      coin.IsNil(),
		}
	}

	return method.Outputs.Pack(coin.Amount.BigInt())
}
