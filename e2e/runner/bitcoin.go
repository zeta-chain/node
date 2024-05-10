package runner

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/pkg/proofs/bitcoin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	zetabitcoin "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
)

var blockHeaderBTCTimeout = 5 * time.Minute

// DepositBTCWithAmount deposits BTC on ZetaChain with a specific amount
func (runner *E2ERunner) DepositBTCWithAmount(amount float64) (txHash *chainhash.Hash) {
	runner.Logger.Print("⏳ depositing BTC into ZEVM")

	// fetch utxos
	utxos, err := runner.BtcRPCClient.ListUnspentMinMaxAddresses(1, 9999999, []btcutil.Address{runner.BTCDeployerAddress})
	if err != nil {
		panic(err)
	}

	spendableAmount := 0.0
	spendableUTXOs := 0
	for _, utxo := range utxos {
		if utxo.Spendable {
			spendableAmount += utxo.Amount
			spendableUTXOs++
		}
	}

	if spendableAmount < amount {
		panic(fmt.Errorf("not enough spendable BTC to run the test; have %f", spendableAmount))
	}

	runner.Logger.Info("ListUnspent:")
	runner.Logger.Info("  spendableAmount: %f", spendableAmount)
	runner.Logger.Info("  spendableUTXOs: %d", spendableUTXOs)
	runner.Logger.Info("Now sending two txs to TSS address...")

	amount = amount + zetabitcoin.DefaultDepositorFee
	txHash, err = runner.SendToTSSFromDeployerToDeposit(runner.BTCTSSAddress, amount, utxos, runner.BtcRPCClient, runner.BTCDeployerAddress)
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("send BTC to TSS txHash: %s", txHash.String())

	return txHash
}

// DepositBTC deposits BTC on ZetaChain
func (runner *E2ERunner) DepositBTC(testHeader bool) {
	runner.Logger.Print("⏳ depositing BTC into ZEVM")
	startTime := time.Now()
	defer func() {
		runner.Logger.Print("✅ BTC deposited in %s", time.Since(startTime))
	}()

	// fetch utxos
	btc := runner.BtcRPCClient
	utxos, err := runner.BtcRPCClient.ListUnspent()
	if err != nil {
		panic(err)
	}
	spendableAmount := 0.0
	spendableUTXOs := 0
	for _, utxo := range utxos {
		if utxo.Spendable {
			spendableAmount += utxo.Amount
			spendableUTXOs++
		}
	}

	if spendableAmount < 1.15 {
		panic(fmt.Errorf("not enough spendable BTC to run the test; have %f", spendableAmount))
	}
	if spendableUTXOs < 5 {
		panic(fmt.Errorf("not enough spendable BTC UTXOs to run the test; have %d", spendableUTXOs))
	}

	runner.Logger.Info("ListUnspent:")
	runner.Logger.Info("  spendableAmount: %f", spendableAmount)
	runner.Logger.Info("  spendableUTXOs: %d", spendableUTXOs)
	runner.Logger.Info("Now sending two txs to TSS address...")

	// send two transactions to the TSS address
	amount1 := 1.1 + zetabitcoin.DefaultDepositorFee
	txHash1, err := runner.SendToTSSFromDeployerToDeposit(runner.BTCTSSAddress, amount1, utxos[:2], btc, runner.BTCDeployerAddress)
	if err != nil {
		panic(err)
	}
	amount2 := 0.05 + zetabitcoin.DefaultDepositorFee
	txHash2, err := runner.SendToTSSFromDeployerToDeposit(runner.BTCTSSAddress, amount2, utxos[2:4], btc, runner.BTCDeployerAddress)
	if err != nil {
		panic(err)
	}

	// send a donation to the TSS address to compensate for the funds minted automatically during pool creation
	// and prevent accounting errors
	_, err = runner.SendToTSSFromDeployerWithMemo(
		runner.BTCTSSAddress,
		0.11,
		utxos[4:5],
		btc,
		[]byte(constant.DonationMessage),
		runner.BTCDeployerAddress,
	)
	if err != nil {
		panic(err)
	}

	runner.Logger.Info("testing if the deposit into BTC ZRC20 is successful...")

	cctx := utils.WaitCctxMinedByInTxHash(runner.Ctx, txHash2.String(), runner.CctxClient, runner.Logger, runner.CctxTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		panic(fmt.Sprintf(
			"expected mined status; got %s, message: %s",
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage),
		)
	}

	balance, err := runner.BTCZRC20.BalanceOf(&bind.CallOpts{}, runner.DeployerAddress)
	if err != nil {
		panic(err)
	}
	if balance.Cmp(big.NewInt(0)) != 1 {
		panic("balance should be positive")
	}

	// due to the high block throughput in localnet, ZetaClient might catch up slowly with the blocks
	// to optimize block header proof test, this test is directly executed here on the first deposit instead of having a separate test
	if testHeader {
		runner.ProveBTCTransaction(txHash1)
	}
}

func (runner *E2ERunner) SendToTSSFromDeployerToDeposit(
	to btcutil.Address,
	amount float64,
	inputUTXOs []btcjson.ListUnspentResult,
	btc *rpcclient.Client,
	btcDeployerAddress *btcutil.AddressWitnessPubKeyHash,
) (*chainhash.Hash, error) {
	return runner.SendToTSSFromDeployerWithMemo(to, amount, inputUTXOs, btc, runner.DeployerAddress.Bytes(), btcDeployerAddress)
}

func (runner *E2ERunner) SendToTSSFromDeployerWithMemo(
	to btcutil.Address,
	amount float64,
	inputUTXOs []btcjson.ListUnspentResult,
	btcRPC *rpcclient.Client,
	memo []byte,
	btcDeployerAddress *btcutil.AddressWitnessPubKeyHash,
) (*chainhash.Hash, error) {
	// prepare inputs
	inputs := make([]btcjson.TransactionInput, len(inputUTXOs))
	inputSats := btcutil.Amount(0)
	amounts := make([]float64, len(inputUTXOs))
	scriptPubkeys := make([]string, len(inputUTXOs))

	for i, utxo := range inputUTXOs {
		inputs[i] = btcjson.TransactionInput{utxo.TxID, utxo.Vout}
		inputSats += btcutil.Amount(utxo.Amount * btcutil.SatoshiPerBitcoin)
		amounts[i] = utxo.Amount
		scriptPubkeys[i] = utxo.ScriptPubKey
	}

	feeSats := btcutil.Amount(0.0001 * btcutil.SatoshiPerBitcoin)
	amountSats := btcutil.Amount(amount * btcutil.SatoshiPerBitcoin)
	change := inputSats - feeSats - amountSats

	if change < 0 {
		return nil, fmt.Errorf("not enough input amount in sats; wanted %d, got %d", amountSats+feeSats, inputSats)
	}
	amountMap := map[btcutil.Address]btcutil.Amount{
		to:                 amountSats,
		btcDeployerAddress: change,
	}

	// create raw
	runner.Logger.Info("ADDRESS: %s, %s", btcDeployerAddress.EncodeAddress(), to.EncodeAddress())
	tx, err := btcRPC.CreateRawTransaction(inputs, amountMap, nil)
	if err != nil {
		panic(err)
	}

	// this adds a OP_RETURN + single BYTE len prefix to the data
	nullData, err := txscript.NullDataScript(memo)
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("nulldata (len %d): %x", len(nullData), nullData)
	if err != nil {
		panic(err)
	}
	memoOutput := wire.TxOut{Value: 0, PkScript: nullData}
	tx.TxOut = append(tx.TxOut, &memoOutput)
	tx.TxOut[1], tx.TxOut[2] = tx.TxOut[2], tx.TxOut[1]

	// make sure that TxOut[0] is sent to "to" address; TxOut[2] is change to oneself. TxOut[1] is memo.
	if bytes.Compare(tx.TxOut[0].PkScript[2:], to.ScriptAddress()) != 0 {
		runner.Logger.Info("tx.TxOut[0].PkScript: %x", tx.TxOut[0].PkScript)
		runner.Logger.Info("to.ScriptAddress():   %x", to.ScriptAddress())
		runner.Logger.Info("swapping txout[0] with txout[2]")
		tx.TxOut[0], tx.TxOut[2] = tx.TxOut[2], tx.TxOut[0]
	}

	runner.Logger.Info("raw transaction: \n")
	for idx, txout := range tx.TxOut {
		runner.Logger.Info("txout %d", idx)
		runner.Logger.Info("  value: %d", txout.Value)
		runner.Logger.Info("  PkScript: %x", txout.PkScript)
	}

	inputsForSign := make([]btcjson.RawTxWitnessInput, len(inputs))
	for i, input := range inputs {
		inputsForSign[i] = btcjson.RawTxWitnessInput{
			Txid:         input.Txid,
			Vout:         input.Vout,
			Amount:       &amounts[i],
			ScriptPubKey: scriptPubkeys[i],
		}
	}

	stx, signed, err := btcRPC.SignRawTransactionWithWallet2(tx, inputsForSign)
	if err != nil {
		panic(err)
	}
	if !signed {
		panic("btc transaction not signed")
	}
	txid, err := btcRPC.SendRawTransaction(stx, true)
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("txid: %+v", txid)
	_, err = btcRPC.GenerateToAddress(6, btcDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	gtx, err := btcRPC.GetTransaction(txid)
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("rawtx confirmation: %d", gtx.BlockIndex)
	rawtx, err := btcRPC.GetRawTransactionVerbose(txid)
	if err != nil {
		panic(err)
	}

	depositorFee := zetabitcoin.DefaultDepositorFee
	events, err := btcobserver.FilterAndParseIncomingTx(
		btcRPC,
		[]btcjson.TxRawResult{*rawtx},
		0,
		runner.BTCTSSAddress.EncodeAddress(),
		log.Logger,
		runner.BitcoinParams,
		depositorFee,
	)
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("bitcoin intx events:")
	for _, event := range events {
		runner.Logger.Info("  TxHash: %s", event.TxHash)
		runner.Logger.Info("  From: %s", event.FromAddress)
		runner.Logger.Info("  To: %s", event.ToAddress)
		runner.Logger.Info("  Amount: %f", event.Value)
		runner.Logger.Info("  Memo: %x", event.MemoBytes)
	}
	return txid, nil
}

// MineBlocks mines blocks on the BTC chain at a rate of 1 blocks every 5 seconds
// and returns a channel that can be used to stop the mining
func (runner *E2ERunner) MineBlocks() chan struct{} {
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				_, err := runner.BtcRPCClient.GenerateToAddress(1, runner.BTCDeployerAddress, nil)
				if err != nil {
					panic(err)
				}
				time.Sleep(3 * time.Second)
			}
		}
	}()
	return stop
}

// ProveBTCTransaction proves that a BTC transaction is in a block header and that the block header is in ZetaChain
func (runner *E2ERunner) ProveBTCTransaction(txHash *chainhash.Hash) {
	// get tx result
	btc := runner.BtcRPCClient
	txResult, err := btc.GetTransaction(txHash)
	if err != nil {
		panic("should get outTx result")
	}
	if txResult.Confirmations <= 0 {
		panic("outTx should have already confirmed")
	}
	txBytes, err := hex.DecodeString(txResult.Hex)
	if err != nil {
		panic(err)
	}

	// get the block with verbose transactions
	blockHash, err := chainhash.NewHashFromStr(txResult.BlockHash)
	if err != nil {
		panic(err)
	}
	blockVerbose, err := btc.GetBlockVerboseTx(blockHash)
	if err != nil {
		panic("should get block verbose tx")
	}

	// get the block header
	header, err := btc.GetBlockHeader(blockHash)
	if err != nil {
		panic("should get block header")
	}

	// collect all the txs in the block
	txns := []*btcutil.Tx{}
	for _, res := range blockVerbose.Tx {
		txBytes, err := hex.DecodeString(res.Hex)
		if err != nil {
			panic(err)
		}
		tx, err := btcutil.NewTxFromBytes(txBytes)
		if err != nil {
			panic(err)
		}
		txns = append(txns, tx)
	}

	// build merkle proof
	mk := bitcoin.NewMerkle(txns)
	path, index, err := mk.BuildMerkleProof(int(txResult.BlockIndex))
	if err != nil {
		panic("should build merkle proof")
	}

	// verify merkle proof statically
	pass := bitcoin.Prove(*txHash, header.MerkleRoot, path, index)
	if !pass {
		panic("should verify merkle proof")
	}

	// wait for block header to show up in ZetaChain
	startTime := time.Now()
	hash := header.BlockHash()
	for {
		// timeout
		if time.Since(startTime) > blockHeaderBTCTimeout {
			panic("timed out waiting for block header to show up in observer")
		}

		_, err := runner.LightclientClient.BlockHeader(runner.Ctx, &lightclienttypes.QueryGetBlockHeaderRequest{
			BlockHash: hash.CloneBytes(),
		})
		if err != nil {
			runner.Logger.Info("waiting for block header to show up in observer... current hash %s; err %s", hash.String(), err.Error())
		}
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// verify merkle proof through RPC
	res, err := runner.LightclientClient.Prove(runner.Ctx, &lightclienttypes.QueryProveRequest{
		ChainId:   chains.BtcRegtestChain.ChainId,
		TxHash:    txHash.String(),
		BlockHash: blockHash.String(),
		Proof:     proofs.NewBitcoinProof(txBytes, path, index),
		TxIndex:   0, // bitcoin doesn't use txIndex
	})
	if err != nil {
		panic(err)
	}
	if !res.Valid {
		panic("txProof should be valid")
	}
	runner.Logger.Info("OK: txProof verified for inTx: %s", txHash.String())
}
