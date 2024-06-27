package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/e2e/config"
)

type v1AddressItem struct {
	Address   string `json:"address"`
	Category  string `json:"category"`
	ChainID   int    `json:"chain_id"`
	ChainName string `json:"chain_name"`
	Type      string `json:"type"`
}

func NewServeAddressesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve-addresses",
		Short: "Serve addresses of deployed contracts from config file",
		RunE:  runServeAddresses,
	}
	cmd.Flags().StringVarP(&configFile, flagConfig, "c", "", "path to the configuration file")
	if err := cmd.MarkFlagRequired(flagConfig); err != nil {
		fmt.Println("Error marking flag as required")
		os.Exit(1)
	}
	chainName, ok := os.LookupEnv("CHAIN_NAME")
	if !ok {
		chainName = "localnet"
	}
	cmd.Flags().String("chain-name", chainName, "name of the chain to return in results")
	return cmd
}

func runServeAddresses(cmd *cobra.Command, args []string) error {
	// read the config file
	configPath, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return err
	}
	chainName, err := cmd.Flags().GetString("chain-name")
	if err != nil {
		return err
	}

	ethChainName := fmt.Sprintf("eth_%s", chainName)
	zetaChainName := fmt.Sprintf("zeta_%s", chainName)

	// we load the config file in the request because it may not be populated on start

	http.HandleFunc("/v1/addresses", func(w http.ResponseWriter, r *http.Request) {
		conf, err := config.ReadConfig(configPath)
		if err != nil {
			log.Error("unable to read config: %v", err)
			w.WriteHeader(500)
			return
		}

		// TODO: pull TSS addresses, ZRC20 addresses from zetacored
		res := []v1AddressItem{
			// eth
			{
				Address:   conf.Contracts.EVM.ConnectorEthAddr.String(),
				Category:  "messaging",
				ChainID:   1337,
				ChainName: ethChainName,
				Type:      "connector",
			},
			{
				Address:   conf.Contracts.EVM.CustodyAddr.String(),
				Category:  "omnichain",
				ChainID:   1337,
				ChainName: ethChainName,
				Type:      "erc20Custody",
			},
			// zeta
			{
				Address:   conf.Contracts.ZEVM.ConnectorZEVMAddr.String(),
				Category:  "messaging",
				ChainID:   101,
				ChainName: zetaChainName,
				Type:      "connector",
			},
			{
				Address:   conf.Contracts.ZEVM.SystemContractAddr.String(),
				Category:  "omnichain",
				ChainID:   101,
				ChainName: zetaChainName,
				Type:      "systemContract",
			},
		}

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Error("unable to read config: %v", err)
			w.WriteHeader(500)
			return
		}
	})

	return http.ListenAndServe(":9991", nil)
}
