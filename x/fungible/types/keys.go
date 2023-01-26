package types

import (
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
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	ModuleAddress = authtypes.NewModuleAddress(ModuleName)
	//ModuleAddressEVM common.EVMAddress
	ModuleAddressEVM = common.BytesToAddress(ModuleAddress.Bytes())
	AdminAddress     = "zeta1rx9r8hff0adaqhr5tuadkzj4e7ns2ntg446vtt"
)

func init() {
	//fmt.Printf("ModuleAddressEVM of %s: %s\n", ModuleName, ModuleAddressEVM.String())
	// 0x735b14BB79463307AAcBED86DAf3322B1e6226aB
}

const (
	SystemContractKey = "SystemContract-value-"
)
