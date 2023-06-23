package main

import (
	"bufio"
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const node = "tcp://3.218.170.198:26657"
const signer = "tanmay"
const chain_id = "athens_7001-1"
const amount = "100000000"
const broadcast_mode = "sync"

//const node = "tcp://localhost:26657"
//const signer = "zeta"
//const chain_id = "localnet_101-1"
//const amount = "100000000" // Amount in azeta
//const broadcast_mode = "block"

type TokeDistribution struct {
	Address           string   `json:"address"`
	BalanceBefore     sdk.Coin `json:"balance_before"`
	BalanceAfter      sdk.Coin `json:"balance_after"`
	TokensDistributed sdk.Coin `json:"tokens_distributed"`
}

func main() {
	file, _ := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "address-list.json"))
	addresses, err := readLines(file)
	if err != nil {
		panic(err)
	}
	fileS, _ := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "successfull_address.json"))
	addresses, err = readLines(file)
	if err != nil {
		panic(err)
	}
	fileF, _ := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "failed_address.json"))
	addresses, err = readLines(file)
	if err != nil {
		panic(err)
	}

	distributionList := make([]TokeDistribution, len(addresses))
	for i, address := range addresses {
		cmd := exec.Command("zetacored", "q", "bank", "balances", address, "--output", "json", "--denom", "azeta", "--node", node)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(cmd.String())
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			return
		}
		balance := sdk.Coin{}
		err = json.Unmarshal(output, &balance)
		if err != nil {
			panic(err)
		}
		distributionAmount, ok := sdkmath.NewIntFromString(amount)
		if !ok {
			panic("parse error for amount")
		}
		distributionList[i] = TokeDistribution{
			Address:           address,
			BalanceBefore:     balance,
			TokensDistributed: sdk.NewCoin(config.BaseDenom, distributionAmount),
		}
	}

	args := []string{"tx", "bank", "multi-send", signer}
	for _, address := range addresses {
		args = append(args, address)
	}

	args = append(args, []string{distributionList[0].TokensDistributed.String(), "--keyring-backend", "test", "--chain-id", chain_id, "--yes",
		"--broadcast-mode", broadcast_mode, "--gas=auto", "--gas-adjustment=2", "--gas-prices=0.001azeta", "--node", node}...)

	cmd := exec.Command("zetacored", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd.String())
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))

	time.Sleep(7 * time.Second)

	for i, address := range addresses {
		cmd := exec.Command("zetacored", "q", "bank", "balances", address, "--output", "json", "--denom", "azeta", "--node", node)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(cmd.String())
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
			return
		}
		balance := sdk.Coin{}
		err = json.Unmarshal(output, &balance)
		if err != nil {
			panic(err)
		}
		distributionList[i].BalanceAfter = balance
	}
	var successfullDistributions []TokeDistribution
	var failedDistributions []TokeDistribution
	for _, distribution := range distributionList {
		if distribution.BalanceAfter.Sub(distribution.BalanceBefore).IsEqual(distribution.TokensDistributed) {
			successfullDistributions = append(successfullDistributions, distribution)
		} else {
			failedDistributions = append(failedDistributions, distribution)
		}
	}
	succesFile, _ := json.MarshalIndent(successfullDistributions, "", " ")
	_ = os.WriteFile(fileS, succesFile, 0600)
	failedFile, _ := json.MarshalIndent(failedDistributions, "", " ")
	_ = os.WriteFile(fileF, failedFile, 0600)

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
