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
	"github.com/zeta-chain/protocol-contracts/v2/pkg/zrc20.sol"

	ptypes "github.com/zeta-chain/node/precompiles/types"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
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

	bankKeeper     bank.Keeper
	fungibleKeeper fungiblekeeper.Keeper
	cdc            codec.Codec
	kvGasConfig    storetypes.GasConfig
}

func NewIBankContract(
	bankKeeper bank.Keeper,
	fungibleKeeper fungiblekeeper.Keeper,
	cdc codec.Codec,
	kvGasConfig storetypes.GasConfig,
) *Contract {
	return &Contract{
		BaseContract:   ptypes.NewBaseContract(ContractAddress),
		bankKeeper:     bankKeeper,
		fungibleKeeper: fungibleKeeper,
		cdc:            cdc,
		kvGasConfig:    kvGasConfig,
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
	// This function is developed using the
	// Check - Effects - Interactions pattern:
	// 1. Check everything is correct.
	if len(args) != 2 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// function deposit(address zrc20, uint256 amount) external returns (bool success);
	ZRC20Addr, amount := args[0].(common.Address), args[1].(*big.Int)

	// Initialize the ZRC20 ABI, as we need to call the balanceOf and allowance methods.
	ZRC20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// Check for enough balance.
	// function balanceOf(address account) public view virtual override returns (uint256)
	argsBalanceOf := []interface{}{caller}

	resBalanceOf, err := c.CallContract(ctx, ZRC20ABI, ZRC20Addr, "balanceOf", argsBalanceOf)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "CallContract",
			Got:  err.Error(),
		}
	}

	balance := resBalanceOf[0].(*big.Int)
	if balance.Cmp(amount) < 0 {
		return nil, &ptypes.ErrUnexpected{
			When: "balance0f",
			Got:  "not enough balance",
		}
	}

	// Check for enough allowance.
	// function allowance(address owner, address spender) public view virtual override returns (uint256)
	argsAllowance := []interface{}{caller, ContractAddress}

	resAllowance, err := c.CallContract(ctx, ZRC20ABI, ZRC20Addr, "allowance", argsAllowance)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "CallContract",
			Got:  err.Error(),
		}
	}

	allowance := resAllowance[0].(*big.Int)
	if allowance.Cmp(amount) < 0 {
		return nil, &ptypes.ErrUnexpected{
			When: "allowance",
			Got:  "not enough allowance",
		}
	}

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

	// 2. Effect: subtract balance.
	// function transferFrom(address sender, address recipient, uint256 amount) public virtual override returns (bool)
	argsTransferFrom := []interface{}{caller, ContractAddress, amount}

	resTransferFrom, err := c.CallContract(ctx, ZRC20ABI, ZRC20Addr, "transferFrom", argsTransferFrom)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "CallContract",
			Got:  err.Error(),
		}
	}

	transferred := resTransferFrom[0].(bool)
	if !transferred {
		return nil, &ptypes.ErrUnexpected{
			When: "TransferFrom",
			Got:  "transfer not successful",
		}
	}

	// 3. Interactions: create cosmos coin and send.
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

// CallContract calls a given contract on behalf of the precompiled contract.
// Note that the precompile contract address is hardcoded.
func (c *Contract) CallContract(
	ctx sdk.Context,
	abi *abi.ABI,
	dst common.Address,
	method string,
	args []interface{},
) ([]interface{}, error) {
	input, err := abi.Methods[method].Inputs.Pack(args)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "Pack " + method,
			Got:  err.Error(),
		}
	}

	res, err := c.fungibleKeeper.CallEVM(
		ctx,
		*abi,
		ContractAddress,
		dst,
		big.NewInt(0),
		nil,
		true,
		false,
		method,
		input,
	)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "CallEVM " + method,
			Got:  err.Error(),
		}
	}

	if res.VmError != "" {
		return nil, &ptypes.ErrUnexpected{
			When: "VmError " + method,
			Got:  res.VmError,
		}
	}

	ret, err := abi.Methods[method].Outputs.Unpack(res.Ret)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "Unpack " + method,
			Got:  err.Error(),
		}
	}

	return ret, nil
}
