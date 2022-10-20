package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
	"strings"
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

func CmdSetSupportedChains() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-supported-chains [chains separated by comma] ",
		Short: "Broadcast message gasPriceVoter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chains := strings.Split(args[0], ",")
			observerChainList := make([]types.ObserverChain, len(chains))
			for i, chain := range chains {
				observerChainList[i] = types.ParseCommonChaintoObservationChain(chain)
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgSetSupportedChains{
				Creator:   clientCtx.GetFromAddress().String(),
				Chainlist: observerChainList,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
