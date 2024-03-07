package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

const (
	APIURLBlocks        = "https://mempool.space/api/v1/blocks"
	APIUrlBlocksTestnet = "https://mempool.space/testnet/api/v1/blocks"
)

type MempoolBlock struct {
	ID                string     `json:"id"`
	Height            int        `json:"height"`
	Version           int        `json:"version"`
	Time              int        `json:"timestamp"`
	Bits              int        `json:"bits"`
	Nonce             int        `json:"nonce"`
	Difficulty        float64    `json:"difficulty"`
	MerkleRoot        string     `json:"merkle_root"`
	TxCount           int        `json:"tx_count"`
	Size              int        `json:"size"`
	Weight            int        `json:"weight"`
	PreviousBlockHash string     `json:"previousblockhash"`
	MedianTime        int        `json:"mediantime"`
	Extras            BlockExtra `json:"extras"`
}

type BlockExtra struct {
	TotalFees int       `json:"totalFees"`
	MedianFee float64   `json:"medianFee"`
	FeeRange  []float64 `json:"feeRange"`
	Reward    int       `json:"reward"`
	Pool      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"pool"`
	AvgFee                 int     `json:"avgFee"`
	AvgFeeRate             int     `json:"avgFeeRate"`
	CoinbaseRaw            string  `json:"coinbaseRaw"`
	CoinbaseAddress        string  `json:"coinbaseAddress"`
	CoinbaseSignature      string  `json:"coinbaseSignature"`
	CoinbaseSignatureASCII string  `json:"coinbaseSignatureAscii"`
	AvgTxSize              float64 `json:"avgTxSize"`
	TotalInputs            int     `json:"totalInputs"`
	TotalOutputs           int     `json:"totalOutputs"`
	TotalOutputAmt         int     `json:"totalOutputAmt"`
	MedianFeeAmt           int     `json:"medianFeeAmt"`
	FeePercentiles         []int   `json:"feePercentiles"`
	SegwitTotalTxs         int     `json:"segwitTotalTxs"`
	SegwitTotalSize        int     `json:"segwitTotalSize"`
	SegwitTotalWeight      int     `json:"segwitTotalWeight"`
	Header                 string  `json:"header"`
	UTXOSetChange          int     `json:"utxoSetChange"`
	UTXOSetSize            int     `json:"utxoSetSize"`
	TotalInputAmt          int     `json:"totalInputAmt"`
	VirtualSize            float64 `json:"virtualSize"`
	Orphans                []struct {
		Height int    `json:"height"`
		Hash   string `json:"hash"`
		Status string `json:"status"`
	} `json:"orphans"`
	MatchRate      float64 `json:"matchRate"`
	EXpectedFees   big.Int `json:"expectedFees"`
	ExpectedWeight int     `json:"expectedWeight"`
}

func Get(ctx context.Context, path string, v interface{}) error {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	if s := r.StatusCode; s != 200 {
		return fmt.Errorf("status %d", s)
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func GetBlocks(ctx context.Context, n int, testnet bool) ([]MempoolBlock, error) {
	path := fmt.Sprintf("%s/%d", APIURLBlocks, n)
	if testnet {
		path = fmt.Sprintf("%s/%d", APIUrlBlocksTestnet, n)
	}
	blocks := make([]MempoolBlock, 0)
	if err := Get(ctx, path, &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}
