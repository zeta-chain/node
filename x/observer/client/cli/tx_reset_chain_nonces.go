package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/observer/types"
)

func CmdResetChainNonces() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset-chain-nonces [chain-id] [chain-nonce-low] [chain-nonce-high]",
		Short: "Broadcast message to reset chain nonces",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// get chainID as int64
			chainID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			// get chainNonceLow as int64
			chainNonceLow, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			// get chainNonceHigh as int64
			chainNonceHigh, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgResetChainNonces(
				clientCtx.GetFromAddress().String(),
				chainID,
				chainNonceLow,
				chainNonceHigh,
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
