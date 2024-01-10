package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdListOutTxTracker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-out-tx-tracker",
		Short: "list all OutTxTracker",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllOutTxTrackerRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.OutTxTrackerAll(context.Background(), params)
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

func CmdShowOutTxTracker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-out-tx-tracker [chainId] [nonce]",
		Short: "shows a OutTxTracker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			argNonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetOutTxTrackerRequest{
				ChainID: argChain,
				Nonce:   argNonce,
			}

			res, err := queryClient.OutTxTracker(context.Background(), params)
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
