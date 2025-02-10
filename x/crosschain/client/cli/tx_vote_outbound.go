package cli

import (
	"errors"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func CmdVoteOutbound() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vote-outbound [sendHash] [outboundHash] [outBlockHeight] [outGasUsed] [outEffectiveGasPrice] [outEffectiveGasLimit] [valueReceived] [Status] [chain] [outTXNonce] [coinType] [confirmationMode]",
		Short:   "Broadcast message to vote an outbound",
		Example: `zetacored tx crosschain vote-outbound 0x12044bec3b050fb28996630e9f2e9cc8d6cf9ef0e911e73348ade46c7ba3417a 0x4f29f9199b10189c8d02b83568aba4cb23984f11adf23e7e5d2eb037ca309497 67773716 65646 30011221226 100000 297254 0 137 13812 ERC20 SAFE`,
		Args:    cobra.ExactArgs(12),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSendHash := args[0]
			argsOutboundHash := args[1]

			argsOutBlockHeight, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			argsOutGasUsed, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			argsOutEffectiveGasPrice, ok := math.NewIntFromString(args[4])
			if !ok {
				return errors.New("invalid effective gas price, enter 0 if unused")
			}

			argsOutEffectiveGasLimit, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}

			argsMMint := args[6]

			status, err := chains.ReceiveStatusFromString(args[7])
			if err != nil {
				return err
			}

			chain, err := strconv.ParseInt(args[8], 10, 64)
			if err != nil {
				return err
			}

			outTxNonce, err := strconv.ParseUint(args[9], 10, 64)
			if err != nil {
				return err
			}

			coinType, ok := coin.CoinType_value[args[10]]
			if !ok {
				return fmt.Errorf("wrong coin type %s", args[10])
			}
			argsCoinType := coin.CoinType(coinType)

			confirmationMode, ok := types.ConfirmationMode_value[args[11]]
			if !ok {
				return fmt.Errorf("wrong confirmation mode %s", args[11])
			}
			argsConfirmationMode := types.ConfirmationMode(confirmationMode)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgVoteOutbound(
				clientCtx.GetFromAddress().String(),
				argsSendHash,
				argsOutboundHash,
				argsOutBlockHeight,
				argsOutGasUsed,
				argsOutEffectiveGasPrice,
				argsOutEffectiveGasLimit,
				math.NewUintFromString(argsMMint),
				status,
				chain,
				outTxNonce,
				argsCoinType,
				argsConfirmationMode,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
