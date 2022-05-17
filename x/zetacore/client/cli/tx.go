package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
//flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdZetaConversionRateVoter())
	// this line is used by starport scaffolding # 1
	cmd.AddCommand(CmdCreateTSSVoter())

	cmd.AddCommand(CmdGasBalanceVoter())

	cmd.AddCommand(CmdGasPriceVoter())

	cmd.AddCommand(CmdNonceVoter())

	cmd.AddCommand(CmdReceiveConfirmation())

	cmd.AddCommand(CmdSendVoter())

	cmd.AddCommand(CmdSetNodeKeys())

	return cmd
}
