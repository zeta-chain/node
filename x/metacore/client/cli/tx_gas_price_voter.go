package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

var _ = strconv.Itoa(0)

func CmdGasPriceVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-price-voter [chain] [price] [blockNumber]",
		Short: "Broadcast message gasPriceVoter",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsChain := (args[0])
			argsPrice, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			argsSupply := args[2]

			argsBlockNumber, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgGasPriceVoter(clientCtx.GetFromAddress().String(), (argsChain), uint64(argsPrice), (argsSupply),  uint64(argsBlockNumber))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
