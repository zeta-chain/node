package app

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	feemarkettypes "github.com/zeta-chain/ethermint/x/feemarket/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionsModuleTypes "github.com/zeta-chain/node/x/emissions/types"
	fungibleModuleTypes "github.com/zeta-chain/node/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// InitGenesisModuleList returns the module list for genesis initialization
// NOTE: Capability module must occur first so that it can initialize any capabilities
// TODO: enable back IBC
// all commented lines in this function are modules related to IBC
// https://github.com/zeta-chain/node/issues/2573
func InitGenesisModuleList() []string {
	return []string{
		//capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		//ibcexported.ModuleName,
		//ibctransfertypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		paramstypes.ModuleName,
		group.ModuleName,
		genutiltypes.ModuleName,
		upgradetypes.ModuleName,
		evidencetypes.ModuleName,
		vestingtypes.ModuleName,
		observertypes.ModuleName,
		crosschaintypes.ModuleName,
		//ibccrosschaintypes.ModuleName,
		fungibleModuleTypes.ModuleName,
		emissionsModuleTypes.ModuleName,
		authz.ModuleName,
		authoritytypes.ModuleName,
		lightclienttypes.ModuleName,
		consensusparamtypes.ModuleName,
		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
	}
}
