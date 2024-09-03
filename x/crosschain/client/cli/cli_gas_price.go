package cli

import (
	"context"
	"strconv"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func CmdListGasPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-gas-price",
		Short: "list all gasPrice",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllGasPriceRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.GasPriceAll(context.Background(), params)
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

func CmdShowGasPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-gas-price [index]",
		Short: "shows a gasPrice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetGasPriceRequest{
				Index: args[0],
			}

			res, err := queryClient.GasPrice(context.Background(), params)
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

func CmdVoteGasPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-gas-price [chain] [price] [priorityFee] [blockNumber]",
		Short: "Broadcast message to vote gas price",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "invalid chain id %q", args[0])
			}

			argsPrice, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "invalid price %q", args[1])
			}

			argsPriorityFee, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "invalid priorityFee %q", args[2])
			}

			argsBlockNumber, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "invalid blockNumber %q", args[3])
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return errors.Wrap(err, "failed to get client context")
			}

			msg := types.NewMsgVoteGasPrice(
				clientCtx.GetFromAddress().String(),
				argsChain,
				argsPrice,
				argsPriorityFee,
				argsBlockNumber,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
