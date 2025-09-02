package types

import (
	"math/big"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the module name
	ModuleName = "fungible"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_fungible"

	// GatewayGasLimit is the default gas limit for gateway contract calls
	DefaultGatewayGasLimit = 1_500_000
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	ModuleAddress    = authtypes.NewModuleAddress(ModuleName)
	ModuleAddressEVM = common.BytesToAddress(ModuleAddress.Bytes())
	GatewayGasLimit  = big.NewInt(DefaultGatewayGasLimit)
)

const (
	SystemContractKey = "SystemContract-value-"
)
