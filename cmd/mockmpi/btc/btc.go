package btc

import (
	"bytes"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/zeta-chain/zetacore/cmd/mockmpi/common"
	"github.com/zeta-chain/zetacore/cmd/mockmpi/eth"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

var (
	BlocksAPIKey string
	PrivateKey   string
)

const (
	FundsTestAddress         = "mhmtT7ecFutdWiwEAzNPykAA8z9WFfa6r8"
	TSSAddress               = "n3RffohEjyrpEviP5e8hFNMS9KdFtQ3rLo"
	Fee                      = 1000
	MinConfirmationThreshold = 1
)

type Chain struct {
	tracker *BlockTracker
}

func RegisterChain() {
	common.ALL_CHAINS = append(common.ALL_CHAINS, &Chain{})
}

func (c *Chain) Start() {
	log("Started BTC chain")
	c.tracker = &BlockTracker{
		OnDeposit: c.handleDeposit,
	}

	// Just for test, let's use 2192650
	c.tracker.Start(2192662)
}

func (c *Chain) handleDeposit(chainID int, to string, satoshiAmount int) {
	ch, err := common.FindChainByID(uint16(chainID))
	if err != nil {
		panic(err)
	}

	// The below is an ugly hack and can only work one way.
	// If we have a similar approach in the opposite direction
	// we'll have a dependency cycle.
	// The right approach is to define a Payload that is generic and not ETH specific,
	// as well as a SendTransaction method that accepts a generic and chain independent payload.

	switch chain := ch.(type) {
	case *eth.ChainETHish:

		// 1 BTC is 10^3 Zeta,
		// and 1 BTC is 10^8 satoshi,
		// so 1 Zeta is 1^5 satoshi.
		// In the API, 1 Zeta amounts to 10^18 big.Int units,
		// so 1 Satoshi is 10^13 Zeta units.
		coefficient := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(13), nil)
		amountToSend := big.NewInt(0).Mul(big.NewInt(int64(satoshiAmount)), coefficient)

		chain.SendTransaction(eth.Payload{
			DestChainID:  uint16(chainID),
			DestContract: []byte(to),
			GasLimit:     big.NewInt(100000),
			SrcChainID:   1,
			ZetaAmount:   amountToSend,
		})
	}

}

func (c *Chain) ID() uint16 {
	return 1
}

func (c *Chain) Name() string {
	return "BTC"
}

type onDepositEvent func(chain int, to string, amount int)

func GetPayToAddrScript(address string) []byte {
	rcvAddress, _ := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	rcvScript, _ := txscript.PayToAddrScript(rcvAddress)
	return rcvScript
}

type utxo struct {
	Address     string
	TxID        string
	OutputIndex int
	Script      []byte
	Satoshis    int64
}

func base58AddressToHex(in string) string {
	addr, err := btcutil.DecodeAddress(in, &chaincfg.TestNet3Params)
	if err != nil {
		panic(err)
	}

	addressInHex := hex.EncodeToString(addr.ScriptAddress())
	return addressInHex
}

func DepositPayment(unspentTx utxo, funds int64) utxo {
	privKeyRaw, err := hex.DecodeString(PrivateKey)
	if err != nil {
		panic(err)
	}

	key, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKeyRaw)

	redemTx := wire.NewMsgTx(wire.TxVersion)

	hash, err := chainhash.NewHashFromStr(unspentTx.TxID)
	if err != nil {
		panic(err)
	}

	outPoint := wire.NewOutPoint(hash, uint32(unspentTx.OutputIndex))
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redemTx.AddTxIn(txIn)

	appendPayments(unspentTx, redemTx, int64(Fee), funds)
	memo := encodeMemo()
	redemTx.AddTxOut(wire.NewTxOut(0, memo))

	sig, err := txscript.SignatureScript(
		redemTx,
		0,
		unspentTx.Script,
		txscript.SigHashAll,
		key,
		false)
	if err != nil {
		panic(err)
	}
	redemTx.TxIn[0].SignatureScript = sig

	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(unspentTx.Script, redemTx, 0, flags, nil, nil, unspentTx.Satoshis)
	if err != nil {
		logf("err != nil: %v\n", err)
	}

	if err := vm.Execute(); err != nil {
		logf("vm.Execute > err != nil: %v\n", err)
	}

	txHex := txToHex(redemTx)
	logf("Sending transaction: %v\n", txHex)

	queryBlockAPI("sendrawtransaction", fmt.Sprintf("[\"%s\"]", txHex))

	// Get txID from above API call

	return utxo{
		Address:     FundsTestAddress,
		TxID:        redemTx.TxHash().String(),
		OutputIndex: 1,
		Script:      GetPayToAddrScript(FundsTestAddress),
		Satoshis:    unspentTx.Satoshis - Fee - funds,
	}

}

func appendPayments(unspentTx utxo, redemTx *wire.MsgTx, fee int64, funds int64) {
	for i, payment := range []int64{funds, unspentTx.Satoshis - (fee + funds)} {
		var rcvScript []byte
		if i == 0 {
			rcvScript = GetPayToAddrScript(TSSAddress)
		} else {
			rcvScript = GetPayToAddrScript(FundsTestAddress)
		}

		txOut := wire.NewTxOut(payment, rcvScript)
		redemTx.AddTxOut(txOut)
	}
}

type Memo struct {
	ChainID int
	Address string
}

func encodeMemo() []byte {
	m := Memo{
		ChainID: 5,
		Address: TSSAddress,
	}

	b, err := asn1.Marshal(m)
	if err != nil {
		panic(err)
	}

	memo, err := txscript.NullDataScript(b)
	if err != nil {
		panic(err)
	}

	return memo
}

func txToHex(tx *wire.MsgTx) string {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf.Bytes())
}

type BcInfo struct {
	Result ChainInfo
}

type ChainInfo struct {
	Blocks        int
	BestBlockHash string
}

type Tx struct {
	BlockHash    string
	IndexInBlock int
	Outputs      []Out
}

type Out struct {
	Script string
	Value  int64
}

type BlockTracker struct {
	OnDeposit   onDepositEvent
	running     sync.WaitGroup
	stopChannel stopChan
}

func (bt *BlockTracker) Stop() {
	select {
	case <-bt.stopChannel:
		return
	default:

	}

	// Wait for block tracker to stop before returning
	defer bt.running.Wait()

	close(bt.stopChannel)
}

func (bt *BlockTracker) Start(latestBlockHeightPersisted int) {
	bt.stopChannel = make(stopChan)
	bt.running.Add(1)
	go func() {
		defer bt.running.Done()
		bt.trackBlocks(latestBlockHeightPersisted)
	}()
}

func (bt *BlockTracker) trackBlocks(latestBlockHeightPersisted int) {
	nextBlockHeightToAnalyze := latestBlockHeightPersisted
	nextBlockHashToAnalyze := blockStats(nextBlockHeightToAnalyze).Result.BlockHash

	for {
		if bt.stopChannel.shouldStop() {
			return
		}

		// Next block has not formed yet
		if nextBlockHashToAnalyze == "" {
			nextBlockHashToAnalyze = blockStats(nextBlockHeightToAnalyze).Result.BlockHash
		}

		if nextBlockHashToAnalyze == "" {
			time.Sleep(time.Second * 30)
			continue
		}

		analyzed := bt.trackTransactions(TSSAddress, nextBlockHeightToAnalyze, nextBlockHashToAnalyze)
		if analyzed {
			nextBlockHeightToAnalyze++
			nextBlockHashToAnalyze = blockStats(nextBlockHeightToAnalyze).Result.BlockHash
			continue
		}

		// If we didn't analyze, it's due to the block not being buried deep enough,
		// so let's sleep and probe again if it was buried deep enough.

		time.Sleep(time.Second * 30)
	}
}

type stopChan chan struct{}

func (sc stopChan) shouldStop() bool {
	select {
	case <-sc:
		return true
	default:

	}
	return false
}

func blockStats(latestBlockHeight int) BlockStatsInfo {
	rawBlockStats := queryBlockAPI("getblockstats", fmt.Sprintf("[%d]", latestBlockHeight))

	if len(rawBlockStats) == 0 {
		return BlockStatsInfo{}
	}

	var bsi BlockStatsInfo
	err := json.Unmarshal(rawBlockStats, &bsi)
	if err != nil {
		panic(err)
	}

	return bsi
}

func (bt *BlockTracker) trackTransactions(depositAddress string, blockHeight int, blockHash string) bool {
	expectedOutScript := fmt.Sprintf("OP_DUP OP_HASH160 %s OP_EQUALVERIFY OP_CHECKSIG", base58AddressToHex(depositAddress))

	bd := blockDataFromRaw(queryBlockAPI("getblock", fmt.Sprintf("[\"%s\"]", blockHash)))

	if bd.Confirmations == 0 {
		log("Block", blockHash, "has not formed yet")
	}

	if bd.Confirmations < MinConfirmationThreshold {
		confirmations := formatConfirmations(bd)
		log("Block", blockHash, "has only", bd.Confirmations, confirmations, "will not bother to analyze it")
		return false
	}

	txns := transactionsOfBlock(blockHash, bd)
	log("Got", len(txns), "transactions for block", blockHeight, blockHash)

	for _, tx := range txns {
		log("Tx", tx.IndexInBlock, ":")
		if amount, isTSSDeposit := isDepositForZetaTSS(tx, expectedOutScript); isTSSDeposit {
			memo := extractMnemonic(tx)
			if memo == nil {
				// No memo, it's either a bug or a donation to Zeta!
				continue
			}
			bt.OnDeposit(memo.ChainID, memo.Address, int(amount))
		}
	}

	return true
}

func formatConfirmations(bd BlockData) string {
	confirmations := "confirmations"
	if bd.Confirmations == 1 {
		confirmations = "confirmation"
	}
	return confirmations
}

func extractMnemonic(tx Tx) *Memo {
	for _, out := range tx.Outputs {
		// Starts with an OP_RETURN
		if !strings.HasPrefix(out.Script, "OP_RETURN") {
			continue
		}

		log(">", out.Script)

		a := strings.Split(out.Script, " ")
		if len(a) < 2 {
			continue
		}

		rawMnemonic := a[1]

		log(">>", rawMnemonic)

		memBytes, err := hex.DecodeString(rawMnemonic)
		if err != nil {
			panic(err)
		}

		m := &Memo{}
		if _, err := asn1.Unmarshal(memBytes, m); err != nil {
			panic(err)
		}

		return m
	}
	return nil
}

func isDepositForZetaTSS(tx Tx, expectedOutScript string) (int64, bool) {
	for _, out := range tx.Outputs {
		logf("Value: %d, Script: [%s]\n", out.Value, out.Script)
		if out.Script == expectedOutScript {
			return out.Value, true
		}
	}
	return 0, false
}

func transactionsOfBlock(latestBlockHash string, bd BlockData) []Tx {
	var txns []Tx
	for i, TxID := range bd.Tx {
		tx := Tx{
			BlockHash:    latestBlockHash,
			IndexInBlock: i,
		}

		time.Sleep(time.Millisecond * 100)
		rawTxHex := queryBlockAPI("getrawtransaction", fmt.Sprintf("[\"%s\"]", TxID))
		rawTx, err := hex.DecodeString(rawTxFromRaw(rawTxHex))
		if err != nil {
			panic(err)
		}

		wTx := wire.NewMsgTx(wire.TxVersion)
		err = wTx.Deserialize(bytes.NewBuffer(rawTx))
		if err != nil {
			panic(err)
		}

		for _, out := range wTx.TxOut {
			script, err := txscript.DisasmString(out.PkScript)
			if err != nil {
				panic(err)
			}
			tx.Outputs = append(tx.Outputs, Out{
				Script: script,
				Value:  out.Value,
			})
		}

		txns = append(txns, tx)
	}

	return txns
}

type BlockStatsInfo struct {
	Result BlockStats
}

type BlockStats struct {
	BlockHash string
}

type RawTransactionInfo struct {
	Result string
}

func rawTxFromRaw(rawData []byte) string {
	var ti RawTransactionInfo
	err := json.Unmarshal(rawData, &ti)
	if err != nil {
		panic(err)
	}
	return ti.Result
}

type BlockDataInfo struct {
	Result BlockData
}

type BlockData struct {
	Confirmations int
	Tx            []string
}

func blockDataFromRaw(rawData []byte) BlockData {
	var bd BlockDataInfo
	err := json.Unmarshal(rawData, &bd)
	if err != nil {
		panic(err)
	}
	return bd.Result
}

func queryBlockAPI(method string, params string) []byte {
	url := "https://btc.getblock.io/testnet/"
	s := `{
    "jsonrpc": "2.0",
    "method": "%s",
    "params": %s,
    "id": "getblock.io"
}`
	s = fmt.Sprintf(s, method, params)

	var jsonStr = []byte(s)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}

	req.Header.Set("x-api-key", BlocksAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		return nil
	}

	return body
}

func log(a ...interface{}) {
	fmt.Println(a...)
}

func logf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
