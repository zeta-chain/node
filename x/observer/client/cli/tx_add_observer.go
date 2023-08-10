package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdAddObserver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-observer [observer-chain-id] [observation-type]",
		Short: "Broadcast message add-observer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argObservationType := args[1]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			chainID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgAddObserver(
				clientCtx.GetFromAddress().String(),
				int64(chainID),
				types.ParseStringToObservationType(argObservationType),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
