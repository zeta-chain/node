package cli

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/zeta-chain/zetacore/common"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdShowTSS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-tss",
		Short: "shows a TSS",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetTSSRequest{}

			res, err := queryClient.TSS(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListTssHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-tss-history",
		Short: "show historical list of TSS",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryTssHistoryRequest{}

			res, err := queryClient.TssHistory(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// Transaction CLI /////////////////////////

func CmdCreateTSSVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-tss-voter [pubkey] [keygenBlock] [status]",
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
			var status common.ReceiveStatus
			if args[2] == "0" {
				status = common.ReceiveStatus_Success
			} else if args[2] == "1" {
				status = common.ReceiveStatus_Failed
			} else {
				return fmt.Errorf("wrong status")
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateTSSVoter(clientCtx.GetFromAddress().String(), argsPubkey, keygenBlock, status)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

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
