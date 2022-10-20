package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdListTSS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-tss",
		Short: "list all TSS",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllTSSRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.TSSAll(context.Background(), params)
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

func CmdShowTSS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-tss [index]",
		Short: "shows a TSS",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetTSSRequest{
				Index: args[0],
			}

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

func CmdListTSSVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-tss-voter",
		Short: "list all TSSVoter",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllTSSVoterRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.TSSVoterAll(context.Background(), params)
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

func CmdShowTSSVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-tss-voter [index]",
		Short: "shows a TSSVoter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetTSSVoterRequest{
				Index: args[0],
			}

			res, err := queryClient.TSSVoter(context.Background(), params)
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
		Use:   "create-tss-voter [chain] [address] [pubkey]",
		Short: "Create a new TSSVoter",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			argsChain, err := cast.ToStringE(args[0])
			if err != nil {
				return err
			}
			argsAddress, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}
			argsPubkey, err := cast.ToStringE(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateTSSVoter(clientCtx.GetFromAddress().String(), argsChain, argsAddress, argsPubkey)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
