package keeper

import (
	storetypes "cosmossdk.io/store/types"
	evidencetypes "cosmossdk.io/x/evidence/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func StoreKeys() (
	map[string]*storetypes.KVStoreKey,
	map[string]*storetypes.MemoryStoreKey,
	map[string]*storetypes.TransientStoreKey,
	map[string]storetypes.StoreKey,
) {
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		group.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		ibcexported.StoreKey,
		ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey,
		authzkeeper.StoreKey,
		evmtypes.StoreKey,
		erc20types.StoreKey,
		feemarkettypes.StoreKey,
		authoritytypes.StoreKey,
		lightclienttypes.StoreKey,
		crosschaintypes.StoreKey,
		observertypes.StoreKey,
		fungibletypes.StoreKey,
		emissionstypes.StoreKey,
		consensusparamtypes.StoreKey,
		crisistypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	allKeys := make(map[string]storetypes.StoreKey, len(keys)+len(tkeys)+len(memKeys))
	for k, v := range keys {
		allKeys[k] = v
	}
	for k, v := range tkeys {
		allKeys[k] = v
	}
	for k, v := range memKeys {
		allKeys[k] = v
	}

	return keys, memKeys, tkeys, allKeys
}
