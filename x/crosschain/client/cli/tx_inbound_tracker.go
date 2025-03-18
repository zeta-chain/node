package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// CmdAddInboundTracker returns the command to add an inbound tracker
func CmdAddInboundTracker() *cobra.Command {
	cmd := &cobra.Command{
		Use: "add-inbound-tracker [chain-id] [tx-hash] [coin-type]",
		Short: `Add an inbound tracker 
				Use 0:Zeta,1:Gas,2:ERC20`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			argTxHash := args[1]
			argsCoinType, err := coin.GetCoinType(args[2])
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgAddInboundTracker(
				clientCtx.GetFromAddress().String(),
				argChain,
				argsCoinType,
				argTxHash,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdRemoveInboundTracker returns the command to remove an inbound tracker
func CmdRemoveInboundTracker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-inbound-tracker [chain-id] [tx-hash]",
		Short: `Remove an inbound tracker`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			argTxHash := args[1]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgRemoveInboundTracker(
				clientCtx.GetFromAddress().String(),
				argChain,
				argTxHash,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
