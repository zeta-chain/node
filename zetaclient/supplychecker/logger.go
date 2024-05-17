package supplychecker

import (
	sdkmath "cosmossdk.io/math"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
)

// ZetaSupplyCheckLogs is a struct to log the output of the ZetaSupplyChecker
type ZetaSupplyCheckLogs struct {
	Logger                   zerolog.Logger
	AbortedTxAmounts         sdkmath.Int `json:"aborted_tx_amounts"`
	ZetaInTransit            sdkmath.Int `json:"zeta_in_transit"`
	ExternalChainTotalSupply sdkmath.Int `json:"external_chain_total_supply"`
	ZetaTokenSupplyOnNode    sdkmath.Int `json:"zeta_token_supply_on_node"`
	EthLockedAmount          sdkmath.Int `json:"eth_locked_amount"`
	NodeAmounts              sdkmath.Int `json:"node_amounts"`
	LHS                      sdkmath.Int `json:"LHS"`
	RHS                      sdkmath.Int `json:"RHS"`
	SupplyCheckSuccess       bool        `json:"supply_check_success"`
}

func (z ZetaSupplyCheckLogs) LogOutput() {
	output, err := bitcoin.PrettyPrintStruct(z)
	if err != nil {
		z.Logger.Error().Err(err).Msgf("error pretty printing struct")
	}
	z.Logger.Info().Msgf(output)
}
