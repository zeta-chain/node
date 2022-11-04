package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdListNodeAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-node-account",
		Short: "list all NodeAccount",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllNodeAccountRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.NodeAccountAll(context.Background(), params)
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

func CmdShowNodeAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-node-account [index]",
		Short: "shows a NodeAccount",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetNodeAccountRequest{
				Index: args[0],
			}

			res, err := queryClient.NodeAccount(context.Background(), params)
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

func CmdSetNodeKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-node-keys [secp256k1] [ed25519] [validatorConsensusPubkey]",
		Short: "Broadcast message SetNodeKeys",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			secp256k1Key, err := common.NewPubKey(args[0])
			if err != nil {
				return fmt.Errorf("fail to parse secp256k1 pub key ,err:%w", err)
			}
			ed25519Key, err := common.NewPubKey(args[1])
			//TODO: re-enable the check. THis is for test when ed25519 key is not supported
			if err != nil {
				//return fmt.Errorf("fail to parse ed25519 pub key ,err:%w", err)
				fmt.Printf("fail to parse ed25519 pub key ,err:%s", err)
			}
			pk := common.NewPubKeySet(secp256k1Key, ed25519Key)
			validatorConsPubKey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeConsPub, args[2])
			if err != nil {
				return fmt.Errorf("fail to parse validator consensus public key: %w", err)
			}
			validatorConsPubKeyStr, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeConsPub, validatorConsPubKey)
			if err != nil {
				return fmt.Errorf("fail to convert public key to string: %w", err)
			}

			msg := types.NewMsgSetNodeKeys(clientCtx.GetFromAddress().String(), pk, validatorConsPubKeyStr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
