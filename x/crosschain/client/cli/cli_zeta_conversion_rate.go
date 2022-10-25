package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdListZetaConversionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-zeta-conversion-rate",
		Short: "list all zetaConversionRate",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllZetaConversionRateRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ZetaConversionRateAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowZetaConversionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-zeta-conversion-rate [index]",
		Short: "shows a zetaConversionRate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := args[0]

			params := &types.QueryGetZetaConversionRateRequest{
				Index: argIndex,
			}

			res, err := queryClient.ZetaConversionRate(context.Background(), params)
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

func CmdZetaConversionRateVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zeta-conversion-rate-voter [chain] [zeta-conversion-rate] [block-number]",
		Short: "Broadcast message ZetaConversionRateVoter",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain := args[0]
			argRate := args[1]

			argBlockNumber, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgZetaConversionRateVoter(
				clientCtx.GetFromAddress().String(),
				argChain,
				argRate,
				uint64(argBlockNumber),
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
