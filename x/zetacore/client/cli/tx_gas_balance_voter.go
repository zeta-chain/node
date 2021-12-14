package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

var _ = strconv.Itoa(0)

func CmdGasBalanceVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-balance-voter [chain] [balance] [blockNumber]",
		Short: "Broadcast message gasBalanceVoter",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsChain := string(args[0])
			argsBalance := string(args[1])
			argsBlockNumber, _ := strconv.Atoi(args[2])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgGasBalanceVoter(clientCtx.GetFromAddress().String(), string(argsChain), string(argsBalance), uint64(argsBlockNumber))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
