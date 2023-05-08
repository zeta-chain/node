package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
// flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
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

	cmd.AddCommand(CmdAddToWatchList())
	cmd.AddCommand(CmdCreateTSSVoter())
	cmd.AddCommand(CmdGasPriceVoter())
	cmd.AddCommand(CmdNonceVoter())
	cmd.AddCommand(CmdCCTXOutboundVoter())
	cmd.AddCommand(CmdCCTXInboundVoter())
	cmd.AddCommand(CmdSetNodeKeys())
	cmd.AddCommand(CmdRemoveFromWatchList())
	cmd.AddCommand(CmdUpdatePermissionFlags())
	cmd.AddCommand(CmdUpdateKeygen())
	// this line is used by starport scaffolding # 1

	return cmd
}
