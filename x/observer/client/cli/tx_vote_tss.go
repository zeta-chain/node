package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/types"
)

func CmdVoteTSS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-tss [pubkey] [keygen-block] [status]",
		Short: "Vote for a new TSS creation",
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

			status, err := chains.ReceiveStatusFromString(args[2])
			if err != nil {
				return err
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
