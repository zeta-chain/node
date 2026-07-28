package snapshot

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Genesis holds the module genesis states needed for the snapshot. Only bank
// and staking are consumed: distribution is exported (rewards are already
// folded into bank via --for-zero-height) but not read here.
type Genesis struct {
	Bank    banktypes.GenesisState
	Staking stakingtypes.GenesisState
}

// ParseAppState pulls the bank and staking genesis states out of a decoded
// `app_state` map. It is deterministic and does no file or RPC I/O; the caller
// reads the export file and hands over the raw module JSON.
func ParseAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) (*Genesis, error) {
	bankRaw, ok := appState[banktypes.ModuleName]
	if !ok {
		return nil, fmt.Errorf("app_state is missing the %q module", banktypes.ModuleName)
	}
	stakingRaw, ok := appState[stakingtypes.ModuleName]
	if !ok {
		return nil, fmt.Errorf("app_state is missing the %q module", stakingtypes.ModuleName)
	}

	var gen Genesis
	if err := cdc.UnmarshalJSON(bankRaw, &gen.Bank); err != nil {
		return nil, fmt.Errorf("unmarshal bank genesis: %w", err)
	}
	if err := cdc.UnmarshalJSON(stakingRaw, &gen.Staking); err != nil {
		return nil, fmt.Errorf("unmarshal staking genesis: %w", err)
	}
	return &gen, nil
}
