package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/cmd/zetacored/config"
)

const node = "tcp://3.218.170.198:26657"
const signer = ""
const chainID = "athens_7001-1"
const amount = "100000000000000000000"
const broadcastMode = "sync"

//const node = "tcp://localhost:26657"
//const signer = "zeta"
//const chain_id = "localnet_101-1"
//const amount = "100000000" // Amount in azeta
//const broadcast_mode = "block"

type TokenDistribution struct {
	Address           string   `json:"address"`
	BalanceBefore     sdk.Coin `json:"balance_before"`
	BalanceAfter      sdk.Coin `json:"balance_after"`
	TokensDistributed sdk.Coin `json:"tokens_distributed"`
}

func main() {
	file, err := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "address-list.json"))
	if err != nil {
		fmt.Printf("error getting absolute path of address-list.json: %s\n", err)
		os.Exit(1)
	}
	addresses, err := readLines(file)
	if err != nil {
		fmt.Printf("error read file: %s\n", err)
		os.Exit(1)
	}
	addresses = removeDuplicates(addresses)
	fileS, err := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "successful_address.json"))
	if err != nil {
		fmt.Printf("error getting absolute path of successful_address.json: %s\n", err)
		os.Exit(1)
	}
	fileF, err := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "failed_address.json"))
	if err != nil {
		fmt.Printf("error getting absolute path of failed_address.json: %s\n", err)
		os.Exit(1)
	}

	distributionList := make([]TokenDistribution, len(addresses))
	for i, address := range addresses {
		// #nosec G204
		cmd := exec.Command(
			"zetacored",
			"q",
			"bank",
			"balances",
			address,
			"--output",
			"json",
			"--denom",
			"azeta",
			"--node",
			node,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(cmd.String())
			fmt.Printf("error getting balance for address %s: %s\n", address, string(output))
			os.Exit(1)
		}
		balance := sdk.Coin{}
		err = json.Unmarshal(output, &balance)
		if err != nil {
			fmt.Printf("error unmarshal balance for address %s: %s\n", address, err)
			os.Exit(1)
		}
		distributionAmount, ok := sdkmath.NewIntFromString(amount)
		if !ok {
			fmt.Printf("error unmarshalling amount: %s\n", amount)
			os.Exit(1)
		}
		distributionList[i] = TokenDistribution{
			Address:           address,
			BalanceBefore:     balance,
			TokensDistributed: sdk.NewCoin(config.BaseDenom, distributionAmount),
		}
	}

	args := make([]string, 0, 17+len(addresses))
	args = append(args, "tx", "bank", "multi-send", signer)
	args = append(args, addresses...)
	args = append(
		args,
		[]string{
			distributionList[0].TokensDistributed.String(),
			"--keyring-backend",
			"test",
			"--chain-id",
			chainID,
			"--yes",
			"--broadcast-mode",
			broadcastMode,
			"--gas=auto",
			"--gas-adjustment=2",
			"--gas-prices=0.001azeta",
			"--node",
			node,
		}...)

	// #nosec G204
	cmd := exec.Command("zetacored", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd.String())
		fmt.Printf("error distributing tokens: %s\n", string(output))
		os.Exit(1)
	}
	fmt.Println(string(output))

	time.Sleep(7 * time.Second)

	for i, address := range addresses {
		// #nosec G204
		cmd := exec.Command(
			"zetacored",
			"q",
			"bank",
			"balances",
			address,
			"--output",
			"json",
			"--denom",
			"azeta",
			"--node",
			node,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(cmd.String())
			fmt.Printf("error getting balance for address %s: %s\n", address, string(output))
			os.Exit(1)
		}
		balance := sdk.Coin{}
		err = json.Unmarshal(output, &balance)
		if err != nil {
			fmt.Printf("error unmarshal balance for address %s: %s\n", address, err)
			os.Exit(1)
		}
		distributionList[i].BalanceAfter = balance
	}
	var successfullDistributions []TokenDistribution
	var failedDistributions []TokenDistribution
	for _, distribution := range distributionList {
		if distribution.BalanceAfter.Sub(distribution.BalanceBefore).IsEqual(distribution.TokensDistributed) {
			successfullDistributions = append(successfullDistributions, distribution)
		} else {
			failedDistributions = append(failedDistributions, distribution)
		}
	}
	successFile, err := json.MarshalIndent(successfullDistributions, "", " ")
	if err != nil {
		fmt.Printf("error marshalling successful distributions: %s\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(fileS, successFile, 0600)
	if err != nil {
		fmt.Printf("error writing successful distributions to file: %s\n", err)
		os.Exit(1)
	}
	failedFile, err := json.MarshalIndent(failedDistributions, "", " ")
	if err != nil {
		fmt.Printf("error marshalling failed distributions: %s\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(fileF, failedFile, 0600)
	if err != nil {
		fmt.Printf("error writing failed distributions to file: %s\n", err)
		os.Exit(1)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path) // #nosec G304
	if err != nil {
		return nil, err
	}
	/* #nosec G307 */
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func removeDuplicates(s []string) []string {
	bucket := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}
