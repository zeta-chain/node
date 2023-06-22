package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const node = "tcp://3.218.170.198:26657"
const signer = "tanmay"
const chain_id = "athens_7001-1"
const amount = "1000azeta"
const broadcast_mode = "sync"

//const node = "tcp://localhost:26657"
//const signer = "zeta"
//const chain_id = "localnet_101-1"
//const amount = "100000000azeta"
//const broadcast_mode = "sync"

func main() {
	file, _ := filepath.Abs(filepath.Join("cmd", "zetacore_utils", "address-list.json"))
	addresses, err := readLines(file)
	if err != nil {
		panic(err)
	}
	args := []string{"tx", "bank", "multi-send", signer}
	for _, address := range addresses {
		args = append(args, address)
	}
	fmt.Println("Distributing to :", addresses)
	args = append(args, []string{amount, "--keyring-backend", "test", "--chain-id", chain_id, "--yes",
		"--broadcast-mode", broadcast_mode, "--gas=auto", "--gas-adjustment=2", "--gas-prices=0.001azeta", "--node", node}...)

	cmd := exec.Command("zetacored", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd.String())
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
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
