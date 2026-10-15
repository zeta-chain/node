package snapshot

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Genesis holds the module genesis states needed for the snapshot. Bank and
// staking drive attribution; distribution is only checked, to prove the export
// was taken with --for-zero-height and left no unpaid rewards behind.
type Genesis struct {
	Bank         banktypes.GenesisState
	Staking      stakingtypes.GenesisState
	Distribution distrtypes.GenesisState
}

// ParseAppState pulls the bank, staking and distribution genesis states out of a
// decoded `app_state` map. It is deterministic and does no file or RPC I/O; the
// caller reads the export file and hands over the raw module JSON.
func ParseAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) (*Genesis, error) {
	bankRaw, ok := appState[banktypes.ModuleName]
	if !ok {
		return nil, fmt.Errorf("app_state is missing the %q module", banktypes.ModuleName)
	}
	stakingRaw, ok := appState[stakingtypes.ModuleName]
	if !ok {
		return nil, fmt.Errorf("app_state is missing the %q module", stakingtypes.ModuleName)
	}
	distrRaw, ok := appState[distrtypes.ModuleName]
	if !ok {
		return nil, fmt.Errorf("app_state is missing the %q module", distrtypes.ModuleName)
	}

	var gen Genesis
	if err := cdc.UnmarshalJSON(bankRaw, &gen.Bank); err != nil {
		return nil, fmt.Errorf("unmarshal bank genesis: %w", err)
	}
	if err := cdc.UnmarshalJSON(stakingRaw, &gen.Staking); err != nil {
		return nil, fmt.Errorf("unmarshal staking genesis: %w", err)
	}
	if err := cdc.UnmarshalJSON(distrRaw, &gen.Distribution); err != nil {
		return nil, fmt.Errorf("unmarshal distribution genesis: %w", err)
	}
	return &gen, nil
}
