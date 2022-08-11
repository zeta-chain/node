package app

import (
	zetaWasm "github.com/zeta-chain/zetacore/wasm"
	zetaCoreModuleKeeper "github.com/zeta-chain/zetacore/x/zetacore/keeper"
	"strings"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
)

const (
	ProposalsEnabled        = "true"
	EnableSpecificProposals = ""
	// DefaultJunoInstanceCost is initially set the same as in wasmd
	DefaultZetaInstanceCost uint64 = 60_000
	// DefaultJunoCompileCost set to a large number for testing
	DefaultZetaCompileCost uint64 = 100
	SupportedFeatures             = "iterator,staking,stargate"
)

func GetWasmOpts(codec codec.Codec, appOpts servertypes.AppOptions, keeper zetaCoreModuleKeeper.Keeper) []wasm.Option {
	var wasmOpts []wasm.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	wasmOpts = append(wasmOpts,
		wasmkeeper.WithGasRegister(NewZetaWasmGasRegister()),
		wasmkeeper.WithMessageEncoders(zetaWasm.Encoders(codec)),
		wasmkeeper.WithQueryPlugins(zetaWasm.Plugins(keeper)),
	)
	return wasmOpts
}

func GetEnabledProposals() []wasm.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

// JunoGasRegisterConfig is defaults plus a custom compile amount
func ZetaGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultZetaInstanceCost
	gasConfig.CompileCost = DefaultZetaCompileCost

	return gasConfig
}

func NewZetaWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(ZetaGasRegisterConfig())
}
