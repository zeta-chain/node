package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

var _ = strconv.Itoa(0)

func CmdDeployERC20Mirror() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy-erc-20-mirror [home-chain] [home-erc-20-contract-address] [name] [symbol] [decimals]",
		Short: "Broadcast message DeployERC20Mirror",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argHomeChain := args[0]
			argHomeERC20ContractAddress := args[1]
			argName := args[2]
			argSymbol := args[3]
			argDecimals, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeployERC20Mirror(
				clientCtx.GetFromAddress().String(),
				argHomeChain,
				argHomeERC20ContractAddress,
				argName,
				argSymbol,
				uint32(argDecimals),
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
