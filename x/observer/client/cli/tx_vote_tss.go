package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdVoteTSS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-tss-voter [pubkey] [keygen-block] [status]",
		Short: "Create a new TSSVoter",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			argsPubkey, err := cast.ToStringE(args[0])
			if err != nil {
				return err
			}
			keygenBlock, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			var status chains.ReceiveStatus
			if args[2] == "0" {
				status = chains.ReceiveStatus_Success
			} else if args[2] == "1" {
				status = chains.ReceiveStatus_Failed
			} else {
				return fmt.Errorf("wrong status")
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgVoteTSS(clientCtx.GetFromAddress().String(), argsPubkey, keygenBlock, status)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
