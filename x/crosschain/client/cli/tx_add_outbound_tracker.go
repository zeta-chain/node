package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"strconv"
)

func CmdAddToWatchList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-to-out-tx-tracker [chain] [nonce] [tx-hash]",
		Short: "Add a out-tx-tracker",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			argNonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			argTxHash := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddToOutTxTracker(
				clientCtx.GetFromAddress().String(),
				argChain,
				argNonce,
				argTxHash,
				nil, // TODO: add option to provide a proof from CLI arguments https://github.com/zeta-chain/node/issues/1134
				"",
				-1,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdRemoveFromWatchList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-from-out-tx-tracker [chain] [nonce]",
		Short: "Remove a out-tx-tracker",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			argNonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveFromOutTxTracker(
				clientCtx.GetFromAddress().String(),
				argChain,
				argNonce,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
