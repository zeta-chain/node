package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"strconv"
)

var _ = strconv.Itoa(0)

func CmdListSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-cctx",
		Short: "list all CCTX",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllCctxRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.CctxAll(context.Background(), params)
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

func CmdShowSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-cctx [index]",
		Short: "shows a CCTX",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetCctxRequest{
				Index: args[0],
			}

			res, err := queryClient.Cctx(context.Background(), params)
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
//zetacored tx zetacore cctx-voter 0x96B05C238b99768F349135de0653b687f9c13fEE ETH 0x96B05C238b99768F349135de0653b687f9c13fEE ETH 1000000000000000000 0 message hash 100 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1

func CmdSendVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cctx-voter [sender] [senderChain] [receiver] [receiverChain] [mBurnt] [mMint] [message] [inTxHash] [inBlockHeight]",
		Short: "Broadcast message sendVoter",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSender := (args[0])
			argsSenderChain := (args[1])
			argsReceiver := (args[2])
			argsReceiverChain := (args[3])
			argsMBurnt := (args[4])
			argsMMint := (args[5])
			argsMessage := (args[6])
			argsInTxHash := (args[7])
			argsInBlockHeight, err := strconv.ParseInt(args[8], 10, 64)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSendVoter(clientCtx.GetFromAddress().String(), (argsSender), (argsSenderChain), (argsReceiver), (argsReceiverChain), (argsMBurnt), (argsMMint), (argsMessage), (argsInTxHash), uint64(argsInBlockHeight), 250_000)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
