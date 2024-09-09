package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ethermint "github.com/zeta-chain/ethermint/types"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func NewCodec() *codec.ProtoCodec {
	registry := codectypes.NewInterfaceRegistry()

	cryptocodec.RegisterInterfaces(registry)
	authtypes.RegisterInterfaces(registry)
	authz.RegisterInterfaces(registry)
	banktypes.RegisterInterfaces(registry)
	stakingtypes.RegisterInterfaces(registry)
	slashingtypes.RegisterInterfaces(registry)
	upgradetypes.RegisterInterfaces(registry)
	distrtypes.RegisterInterfaces(registry)
	evidencetypes.RegisterInterfaces(registry)
	crisistypes.RegisterInterfaces(registry)
	evmtypes.RegisterInterfaces(registry)
	ethermint.RegisterInterfaces(registry)
	crosschaintypes.RegisterInterfaces(registry)
	emissionstypes.RegisterInterfaces(registry)
	fungibletypes.RegisterInterfaces(registry)
	observertypes.RegisterInterfaces(registry)

	return codec.NewProtoCodec(registry)
}
