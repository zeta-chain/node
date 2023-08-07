package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdGetSupportedChains() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-chains",
		Short: "list all SupportedChains",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySupportedChains{}

			res, err := queryClient.SupportedChains(context.Background(), params)
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

// Transaction CLI /////////////////////////

//func CmdSetSupportedChains() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "add-supported-chains chainID chainName ",
//		Short: "Broadcast message set-supported-chains",
//		Args:  cobra.ExactArgs(2),
//		RunE: func(cmd *cobra.Command, args []string) error {
//
//			chainName := common.ParseStringToObserverChain(args[1])
//			if chainName == 0 {
//				return errors.New("ChainName type not supported\"")
//			}
//			chainID, err := strconv.Atoi(args[0])
//			if err != nil {
//				return err
//			}
//			clientCtx, err := client.GetClientTxContext(cmd)
//			if err != nil {
//				return err
//			}
//
//			msg := &types.MsgSetSupportedChains{
//				Creator:   clientCtx.GetFromAddress().String(),
//				ChainId:   int64(chainID),
//				ChainName: chainName,
//			}
//			if err := msg.ValidateBasic(); err != nil {
//				return err
//			}
//			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
//		},
//	}
//
//	flags.AddTxFlagsToCmd(cmd)
//
//	return cmd
//}
