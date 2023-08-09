package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/trie"
	"math/big"
	"testing"
)

func TestProofGeneration(t *testing.T) {
	RPC_URL := "https://rpc.ankr.com/eth_goerli"
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		t.Fatal(err)
	}
	bn := int64(9485814)
	block, err := client.BlockByNumber(context.Background(), big.NewInt(bn))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("block %d\n", block.Number())
	t.Logf("  tx root %x\n", block.Header().TxHash)

	tr := DeriveSha(block.Transactions(), trie.NewStackTrie(nil))
	t.Logf("  sha2    %x\n", tr.Hash())

}
