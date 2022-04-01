package btc

import (
	"fmt"
	"github.com/zeta-chain/zetacore/cmd/mockmpi/common"
	"github.com/zeta-chain/zetacore/cmd/mockmpi/eth"
	"os"
	"testing"
	"time"
)

func init() {
	BlocksAPIKey = "f68d3e47-564b-4125-ae7f-42010c833965"
	PrivateKey = "269b36b89a1d92d4a82938ac08d3b16445a28895192c144d05a195fc695a8e36"
}

func TestTrial(t *testing.T) {
	t.Skipf("do not run in CI")

	os.Setenv("PRIVKEY", "08f2f8501c4f3f7859a37950588aac989ed8427947498ed27e50882073ee6b72")
	os.Setenv("ETH_ENDPOINT", "wss://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/eth/goerli/archive/ws")
	os.Setenv("BSC_ENDPOINT", "wss://speedy-nodes-nyc.moralis.io/eb13a7dfda3e4b15212356f9/bsc/testnet/archive/ws")
	eth.RegisterChains()
	RegisterChain()

	for _, chain := range common.ALL_CHAINS {
		chain.Start()
	}

	time.Sleep(time.Minute * 20)
}

func TestTrackBlockchain(t *testing.T) {
	t.Skipf("do not run in CI")

	tracker := &BlockTracker{
		OnDeposit: func(chain int, to string, amount int) {
			fmt.Println("Depositing", amount, "satoshi", "on chain", chain, "to", to)
		},
	}

	tracker.Start(2192662)

	time.Sleep(time.Minute * 5)
	tracker.Stop()
}

func TestSendBitcoin(t *testing.T) {
	t.Skipf("do not run in CI")

	out := utxo{
		Address:     FundsTestAddress,
		TxID:        "12800ab1f432246c48ad2c5edab815ccb4787ad4f5171b3e6ac252d5925cefc9",
		OutputIndex: 1,
		Script:      GetPayToAddrScript(FundsTestAddress),
		Satoshis:    3368557,
	}

	DepositPayment(out, 10*1000)
}
