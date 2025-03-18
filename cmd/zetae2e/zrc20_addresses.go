package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/pkg/chains"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const flagZRC20Network = "zrc20-network"
const flagZRC20Symbol = "zrc20-symbol"

func registerERC20Flags(cmd *cobra.Command) {
	cmd.Flags().String(flagZRC20Network, "", "network from /zeta-chain/observer/supportedChains")
	cmd.Flags().String(flagZRC20Symbol, "", "symbol from /zeta-chain/fungible/foreign_coins")
}

func processZRC20Flags(cmd *cobra.Command, conf *config.Config) error {
	zrc20ChainName, err := cmd.Flags().GetString(flagZRC20Network)
	if err != nil {
		return err
	}
	zrc20Symbol, err := cmd.Flags().GetString(flagZRC20Symbol)
	if err != nil {
		return err
	}
	if zrc20ChainName != "" && zrc20Symbol != "" {
		erc20Asset, zrc20ContractAddress, chain, err := findZRC20(
			cmd.Context(),
			conf,
			zrc20ChainName,
			zrc20Symbol,
		)
		if err != nil {
			return err
		}
		if chain.IsEVMChain() {
			conf.Contracts.EVM.ERC20 = config.DoubleQuotedString(erc20Asset)
			conf.Contracts.ZEVM.ERC20ZRC20Addr = config.DoubleQuotedString(zrc20ContractAddress)
		} else if chain.IsSolanaChain() {
			conf.Contracts.Solana.SPLAddr = config.DoubleQuotedString(erc20Asset)
			conf.Contracts.ZEVM.SPLZRC20Addr = config.DoubleQuotedString(zrc20ContractAddress)
		}
	}
	return nil
}

// findZRC20 loads ERC20/SPL/etc addresses via gRPC given CLI flags
func findZRC20(
	ctx context.Context,
	conf *config.Config,
	networkName, zrc20Symbol string,
) (string, string, chains.Chain, error) {
	clients, err := zetae2econfig.GetZetacoreClient(*conf)
	if err != nil {
		return "", "", chains.Chain{}, fmt.Errorf("get zeta clients: %w", err)
	}

	supportedChainsRes, err := clients.Observer.SupportedChains(ctx, &observertypes.QuerySupportedChains{})
	if err != nil {
		return "", "", chains.Chain{}, fmt.Errorf("get chain params: %w", err)
	}

	chainID := int64(0)
	for _, chain := range supportedChainsRes.Chains {
		if strings.EqualFold(chain.Network.String(), networkName) {
			chainID = chain.ChainId
			break
		}
	}
	if chainID == 0 {
		return "", "", chains.Chain{}, fmt.Errorf("chain %s not found", networkName)
	}

	chain, ok := chains.GetChainFromChainID(chainID, nil)
	if !ok {
		return "", "", chains.Chain{}, fmt.Errorf("invalid/unknown chain ID %d", chainID)
	}

	foreignCoinsRes, err := clients.Fungible.ForeignCoinsAll(ctx, &fungibletypes.QueryAllForeignCoinsRequest{})
	if err != nil {
		return "", "", chain, fmt.Errorf("get foreign coins: %w", err)
	}

	for _, coin := range foreignCoinsRes.ForeignCoins {
		if coin.ForeignChainId != chainID {
			continue
		}
		// sometimes symbol is USDT, sometimes it's like USDT.SEPOLIA
		if strings.HasPrefix(coin.Symbol, zrc20Symbol) || strings.HasSuffix(coin.Symbol, zrc20Symbol) {
			return coin.Asset, coin.Zrc20ContractAddress, chain, nil
		}
	}
	return "", "", chain, fmt.Errorf("zrc20 %s not found on %s", zrc20Symbol, networkName)
}
