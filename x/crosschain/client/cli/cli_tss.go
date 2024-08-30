package cli

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func CmdUpdateTss() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-tss-address [pubkey]",
		Short: "Create a new TSSVoter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsPubkey, err := cast.ToStringE(args[0])
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateTssAddress(clientCtx.GetFromAddress().String(), argsPubkey)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdMigrateTssFunds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-tss-funds [chainID] [amount]",
		Short: "Migrate TSS funds to the latest TSS address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsChainID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			argsAmount := math.NewUintFromString(args[1])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMigrateTssFunds(clientCtx.GetFromAddress().String(), argsChainID, argsAmount)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
