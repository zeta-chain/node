package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/thorchain/tss/tss-lib/ecdsa/keygen"
)

func init() {
	RootCmd.AddCommand(GenPrePramsCmd)
}

var GenPrePramsCmd = &cobra.Command{
	Use:   "gen-pre-params <path>",
	Short: "Generate pre parameters for TSS",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		startTime := time.Now()
		preParams, err := keygen.GeneratePreParams(time.Second * 300)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(args[0], os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
		err = json.NewEncoder(file).Encode(preParams)
		if err != nil {
			return err
		}
		fmt.Printf("Generated new pre-parameters in %v\n", time.Since(startTime))
		return nil
	},
}
