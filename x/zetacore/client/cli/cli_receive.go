package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"strconv"
)

func CmdListReceive() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-receive",
		Short: "list all receive",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllReceiveRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ReceiveAll(context.Background(), params)
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

func CmdShowReceive() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-receive [index]",
		Short: "shows a receive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetReceiveRequest{
				Index: args[0],
			}

			res, err := queryClient.Receive(context.Background(), params)
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

func CmdReceiveConfirmation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receive-confirmation [sendHash] [outTxHash] [outBlockHeight] [ZetaMinted] [Status] [chain] [outTXNonce]",
		Short: "Broadcast message receiveConfirmation",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSendHash := (args[0])
			argsOutTxHash := (args[1])
			argsOutBlockHeight, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
			argsMMint := (args[3])
			var status common.ReceiveStatus
			if args[4] == "0" {
				status = common.ReceiveStatus_Success
			} else if args[4] == "1" {
				status = common.ReceiveStatus_Failed
			} else {
				return fmt.Errorf("wrong status")
			}
			chain := args[5]
			outTxNonce, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgReceiveConfirmation(clientCtx.GetFromAddress().String(), (argsSendHash), (argsOutTxHash), uint64(argsOutBlockHeight), sdk.NewUintFromString(argsMMint), status, chain, uint64(outTxNonce))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
