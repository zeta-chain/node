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

func CmdListChainNonces() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-chain-nonces",
		Short: "list all chainNonces",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllChainNoncesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ChainNoncesAll(context.Background(), params)
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

func CmdShowChainNonces() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-chain-nonces [index]",
		Short: "shows a chainNonces",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetChainNoncesRequest{
				Index: args[0],
			}

			res, err := queryClient.ChainNonces(context.Background(), params)
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

func CmdNonceVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nonce-voter [chain] [nonce]",
		Short: "Broadcast message nonceVoter",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			argsNonce, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgNonceVoter(clientCtx.GetFromAddress().String(), argsChain, uint64(argsNonce))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
