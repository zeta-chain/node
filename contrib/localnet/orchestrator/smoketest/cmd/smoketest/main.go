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

	DeployerAddress    = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263" // #nosec G101 - used for testing

	FungibleAdminMnemonic = "snow grace federal cupboard arrive fancy gym lady uniform rotate exercise either leave alien grass" // #nosec G101 - used for testing
)

func main() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// 0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6	fda3be1b1517bdf48615bdadacc1e6463d2865868dc8077d2cdcfa4709a16894
// 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC   d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263

// 0x5cC2fBb200A929B372e3016F1925DcF988E081fd   729a6cdc5c925242e7df92fdeeb94dadbf2d0b9950d4db8f034ab27a3b114ba7
