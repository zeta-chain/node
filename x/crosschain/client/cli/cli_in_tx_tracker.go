package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdAddToInTxTracker() *cobra.Command {
	cmd := &cobra.Command{
		Use: "add-to-in-tx-tracker [chain-id] [tx-hash] [coin-type]",
		Short: `Add a out-tx-tracker 
				Use 0:Zeta,1:Gas,2:ERC20`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			argTxHash := args[1]
			argsCoinType, err := common.GetCoinType(args[2])
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgAddToInTxTracker(
				clientCtx.GetFromAddress().String(),
				argChain,
				argsCoinType,
				argTxHash,
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

func CmdListInTxTrackerByChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-in-tx-tracker [chainId]",
		Short: "shows a list of in tx tracker by chainId",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params := &types.QueryAllInTxTrackerByChainRequest{
				ChainId:    argChain,
				Pagination: pageReq,
			}
			res, err := queryClient.InTxTrackerAllByChain(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	return cmd
}

func CmdListInTxTrackers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-all-in-tx-trackers",
		Short: "shows all inTxTrackers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryAllInTxTrackersRequest{}
			res, err := queryClient.InTxTrackerAll(context.Background(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
