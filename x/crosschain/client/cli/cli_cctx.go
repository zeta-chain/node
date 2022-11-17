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
	"github.com/zeta-chain/zetacore/x/crosschain/types"
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

func CmdCCTXInboundVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inbound-voter [sender] [senderChain] [receiver] [receiverChain] [mBurnt] [mMint] [message] [inTxHash] [inBlockHeight] [coinType]",
		Short: "Broadcast message sendVoter",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSender := (args[0])
			argsSenderChain, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			argsReceiver := (args[2])
			argsReceiverChain, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}
			argsMBurnt := (args[4])
			argsMMint := (args[5])
			argsMessage := (args[6])
			argsInTxHash := (args[7])
			argsInBlockHeight, err := strconv.ParseInt(args[8], 10, 64)
			argsCoinType := common.CoinType(common.CoinType_value[args[9]])
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSendVoter(clientCtx.GetFromAddress().String(), argsSender, int64(argsSenderChain), argsReceiver, int64((argsReceiverChain)), (argsMBurnt), (argsMMint), (argsMessage), (argsInTxHash), uint64(argsInBlockHeight), 250_000, argsCoinType)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCCTXOutboundVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outbound-voter [sendHash] [outTxHash] [outBlockHeight] [ZetaMinted] [Status] [chain] [outTXNonce] [coinType]",
		Short: "Broadcast message receiveConfirmation",
		Args:  cobra.ExactArgs(8),
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
			chain, err := strconv.Atoi(args[5])
			if err != nil {
				return err
			}
			outTxNonce, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return err
			}
			argsCoinType := common.CoinType(common.CoinType_value[args[7]])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgReceiveConfirmation(clientCtx.GetFromAddress().String(), argsSendHash, argsOutTxHash, uint64(argsOutBlockHeight), sdk.NewUintFromString(argsMMint), status, int64(chain), uint64(outTxNonce), argsCoinType)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
