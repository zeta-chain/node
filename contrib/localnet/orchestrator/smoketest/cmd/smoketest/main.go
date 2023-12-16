package main

import (
	"fmt"
	"os"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

var (
	// TODO: make these variables configurable
	// https://github.com/zeta-chain/node-private/issues/41

	SmokeTestTimeout = 30 * time.Minute

	DeployerAddress       = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey    = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263"                                   // #nosec G101 - used for testing
	FungibleAdminMnemonic = "snow grace federal cupboard arrive fancy gym lady uniform rotate exercise either leave alien grass" // #nosec G101 - used for testing
)

func main() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
