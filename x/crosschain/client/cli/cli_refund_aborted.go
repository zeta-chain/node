package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdRefundAborted() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-aborted [cctx-index] [refund-address]",
		Short: `Refund an aborted tx , the refund address is optional, if not provided, the refund will be sent to the sender/tx origin of the cctx.`,
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgRefundAbortedCCTX(clientCtx.GetFromAddress().String(), args[0], args[1])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
