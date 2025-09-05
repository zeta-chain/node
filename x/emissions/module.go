package emissions

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/emissions/client/cli"
	"github.com/zeta-chain/node/x/emissions/exported"
	"github.com/zeta-chain/node/x/emissions/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
)

const consensusVersion = 6

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}

	_ appmodule.HasEndBlocker   = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the emissions module.
type AppModuleBasic struct {
	cdc codec.Codec
}

func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the emissions module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns the emissions module's default genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the emissions module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterRESTRoutes registers the emissions module's REST service handlers.
func (AppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *mux.Router) {
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		fmt.Println("RegisterQueryHandlerClient err: %w", err)
	}
}

// GetTxCmd returns the emissions module's root tx command.
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the emissions module's root query command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the emissions module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
	// legacySubspace is used solely for migration of x/params managed parameters
	legacySubspace exported.Subspace
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	ss exported.Subspace,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		legacySubspace: ss,
	}
}

// Name returns the emissions module's name.
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
	m := keeper.NewMigrator(am.keeper, am.legacySubspace)
	if err := cfg.RegisterMigration(types.ModuleName, 5, m.Migrate5to6); err != nil {
		panic(err)
	}
}

// RegisterInvariants registers the emissions module's invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the emissions module's genesis initialization It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	InitGenesis(ctx, am.keeper, genState)

	am.keeper.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	am.keeper.GetAuthKeeper().GetModuleAccount(ctx, types.UndistributedTSSRewardsPool)
	am.keeper.GetAuthKeeper().GetModuleAccount(ctx, types.UndistributedObserverRewardsPool)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the emissions module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion implements ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return consensusVersion }

// BeginBlock executes all ABCI BeginBlock logic respective to the emissions module.
func (am AppModule) BeginBlock(c context.Context) error {
	ctx := sdk.UnwrapSDKContext(c)
	BeginBlocker(ctx, am.keeper)

	return nil
}

// EndBlock executes all ABCI EndBlock logic respective to the emissions module. It
// returns no validator updates.
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}
