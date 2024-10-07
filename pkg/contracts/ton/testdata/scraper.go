package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func main() {
	var testnet bool

	flag.BoolVar(&testnet, "testnet", false, "Use testnet network")
	flag.Parse()

	if len(flag.Args()) < 3 {
		log.Fatalf("Usage: go run scraper.go [-testnet] <account> <lt> <hash>")
	}

	// Parse account
	acc, err := ton.ParseAccountID(flag.Arg(0))
	must(err, "Unable to parse account")

	// Parse LT
	lt, err := strconv.ParseUint(flag.Arg(1), 10, 64)
	must(err, "Unable to parse logical time")

	// Parse hash
	var hash ton.Bits256

	must(hash.FromHex(flag.Arg(2)), "Unable to parse hash")

	ctx, client := context.Background(), getClient(testnet)

	state, err := client.GetAccountState(ctx, acc)
	must(err, "Unable to get account state")

	if state.Account.Status() != tlb.AccountActive {
		fail("account %s is not active", acc.ToRaw())
	}

	txs, err := client.GetTransactions(ctx, 1, acc, lt, hash)
	must(err, "Unable to get transactions")

	switch {
	case len(txs) == 0:
		fail("Not found")
	case len(txs) > 1:
		fail("invalid tx list length (got %d, want 1); lt %d, hash %s", len(txs), lt, hash.Hex())
	}

	// Print the transaction
	tx := txs[0]

	cell, err := transactionToCell(tx)
	must(err, "unable to convert tx to cell")

	bocRaw, err := cell.MarshalJSON()
	must(err, "unable to marshal cell to JSON")

	printAny(map[string]any{
		"test":        testnet,
		"account":     acc.ToRaw(),
		"description": "todo",
		"logicalTime": lt,
		"hash":        hash.Hex(),
		"boc":         json.RawMessage(bocRaw),
	})
}

func getClient(testnet bool) *liteapi.Client {
	if testnet {
		c, err := liteapi.NewClientWithDefaultTestnet()
		must(err, "unable to create testnet lite client")

		return c
	}

	c, err := liteapi.NewClientWithDefaultMainnet()
	must(err, "unable to create mainnet lite client")

	return c
}

func printAny(v any) {
	b, err := json.MarshalIndent(v, "", " ")
	must(err, "unable marshal data")

	fmt.Println(string(b))
}

func transactionToCell(tx ton.Transaction) (*boc.Cell, error) {
	b, err := tx.SourceBoc()
	if err != nil {
		return nil, err
	}

	cells, err := boc.DeserializeBoc(b)
	if err != nil {
		return nil, err
	}

	if len(cells) != 1 {
		return nil, fmt.Errorf("invalid cell count: %d", len(cells))
	}

	return cells[0], nil
}

func must(err error, msg string) {
	if err == nil {
		return
	}

	if msg == "" {
		log.Fatalf("Error: %s", err.Error())
	}

	log.Fatalf("%s; error: %s", msg, err.Error())
}

func fail(msg string, args ...any) {
	must(fmt.Errorf(msg, args...), "FAIL")
}
