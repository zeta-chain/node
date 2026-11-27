package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/types"
)

// btcNetworks lists all supported Bitcoin networks for TSS address display.
var btcNetworks = []struct {
	Name    string
	ChainID int64
}{
	{"mainnet", chains.BitcoinMainnet.ChainId},
	{"testnet3", chains.BitcoinTestnet.ChainId},
	{"signet", chains.BitcoinSignetTestnet.ChainId},
	{"testnet4", chains.BitcoinTestnet4.ChainId},
	{"regtest", chains.BitcoinRegtest.ChainId},
}

func CmdGetTssAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-tss-address [bitcoinChainId]",
		Short: "Query current tss address",
		Long: `Query current TSS addresses for EVM, Sui, and Bitcoin.

Without arguments, prints the EVM and Sui addresses once, and the Bitcoin address for all
supported networks (mainnet, testnet3, signet, testnet4, regtest).

With a bitcoinChainId argument, prints only the Bitcoin address for that specific chain.

Examples:
  zetacored query observer get-tss-address
  zetacored query observer get-tss-address 18333`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// If a specific chain ID is provided, use the original single-query behavior
			if len(args) == 1 {
				bitcoinChainID, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return err
				}
				res, err := queryClient.GetTssAddress(cmd.Context(), &types.QueryGetTssAddressRequest{
					BitcoinChainId: bitcoinChainID,
				})
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			}

			// No args: print EVM/Sui once, then BTC for all networks
			return printAllTSSAddresses(cmd, queryClient)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// printAllTSSAddresses queries TSS addresses for all BTC networks and prints them.
func printAllTSSAddresses(cmd *cobra.Command, queryClient types.QueryClient) error {
	// Query with first BTC network to get EVM and Sui addresses
	first := btcNetworks[0]
	res, err := queryClient.GetTssAddress(cmd.Context(), &types.QueryGetTssAddressRequest{
		BitcoinChainId: first.ChainID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("eth: %s\n", res.Eth)
	fmt.Printf("sui: %s\n", res.Sui)
	fmt.Printf("btc (%s, %d): %s\n", first.Name, first.ChainID, res.Btc)

	// Query remaining BTC networks
	for _, net := range btcNetworks[1:] {
		btcRes, err := queryClient.GetTssAddress(cmd.Context(), &types.QueryGetTssAddressRequest{
			BitcoinChainId: net.ChainID,
		})
		if err != nil {
			fmt.Printf("btc (%s, %d): error: %v\n", net.Name, net.ChainID, err)
			continue
		}
		fmt.Printf("btc (%s, %d): %s\n", net.Name, net.ChainID, btcRes.Btc)
	}

	return nil
}

func CmdGetTssAddressByFinalizedZetaHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-historical-tss-address [finalizedZetaHeight] [bitcoinChainId]",
		Short: "Query tss address by finalized zeta height (for historical tss addresses)",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			finalizedZetaHeight, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			params := &types.QueryGetTssAddressByFinalizedHeightRequest{
				FinalizedZetaHeight: finalizedZetaHeight,
			}
			if len(args) == 2 {
				bitcoinChainID, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					return err
				}
				params.BitcoinChainId = bitcoinChainID
			}

			res, err := queryClient.GetTssAddressByFinalizedHeight(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
