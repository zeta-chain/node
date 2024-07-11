package runner

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

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

// ListDeployerUTXOs list the deployer's UTXOs
func (r *E2ERunner) ListDeployerUTXOs() ([]btcjson.ListUnspentResult, error) {
	// query UTXOs from node
	utxos, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(
		1,
		9999999,
		[]btcutil.Address{r.BTCDeployerAddress},
	)
	if err != nil {
		return nil, err
	}

	// filter big-enough UTXOs for test if running on Regtest
	if r.IsLocalBitcoin() {
		utxosFiltered := []btcjson.ListUnspentResult{}
		for _, utxo := range utxos {
			if utxo.Amount >= 1.0 {
				utxosFiltered = append(utxosFiltered, utxo)
			}
		}
		return utxosFiltered, nil
	}

	return utxos, nil
}

// GetTop20UTXOsForTssAddress returns the top 20 UTXOs for the TSS address.
// Top 20 utxos should be used for TSS migration, as we can only migrate at max 20 utxos at a time.
func (r *E2ERunner) GetTop20UTXOsForTssAddress() ([]btcjson.ListUnspentResult, error) {
	// query UTXOs from node
	utxos, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(
		0,
		9999999,
		[]btcutil.Address{r.BTCTSSAddress},
	)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(utxos, func(i, j int) bool {
		return utxos[i].Amount < utxos[j].Amount
	})

	if len(utxos) > 20 {
		utxos = utxos[:20]
	}
	return utxos, nil
}

// DepositBTCWithAmount deposits BTC on ZetaChain with a specific amount
func (r *E2ERunner) DepositBTCWithAmount(amount float64) *chainhash.Hash {
	r.Logger.Print("⏳ depositing BTC into ZEVM")

	// list deployer utxos
	utxos, err := r.ListDeployerUTXOs()
	require.NoError(r, err)

	spendableAmount := 0.0
	spendableUTXOs := 0
	for _, utxo := range utxos {
		if utxo.Spendable {
			spendableAmount += utxo.Amount
			spendableUTXOs++
		}
	}

	require.LessOrEqual(r, amount, spendableAmount, "not enough spendable BTC to run the test")

	r.Logger.Info("ListUnspent:")
	r.Logger.Info("  spendableAmount: %f", spendableAmount)
	r.Logger.Info("  spendableUTXOs: %d", spendableUTXOs)
	r.Logger.Info("Now sending two txs to TSS address...")

	amount += zetabitcoin.DefaultDepositorFee
	txHash, err := r.SendToTSSFromDeployerToDeposit(amount, utxos)
	require.NoError(r, err)

	r.Logger.Info("send BTC to TSS txHash: %s", txHash.String())

	return txHash
}

// DepositBTC deposits BTC on ZetaChain
func (r *E2ERunner) DepositBTC(testHeader bool) {
	r.Logger.Print("⏳ depositing BTC into ZEVM")
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ BTC deposited in %s", time.Since(startTime))
	}()

	// list deployer utxos
	utxos, err := r.ListDeployerUTXOs()
	require.NoError(r, err)

	spendableAmount := 0.0
	spendableUTXOs := 0
	for _, utxo := range utxos {
		if utxo.Spendable {
			spendableAmount += utxo.Amount
			spendableUTXOs++
		}
	}

	require.GreaterOrEqual(r, spendableAmount, 1.15, "not enough spendable BTC to run the test")
	require.GreaterOrEqual(r, spendableUTXOs, 5, "not enough spendable BTC UTXOs to run the test")

	r.Logger.Info("ListUnspent:")
	r.Logger.Info("  spendableAmount: %f", spendableAmount)
	r.Logger.Info("  spendableUTXOs: %d", spendableUTXOs)
	r.Logger.Info("Now sending two txs to TSS address...")

	// send two transactions to the TSS address
	amount1 := 1.1 + zetabitcoin.DefaultDepositorFee
	txHash1, err := r.SendToTSSFromDeployerToDeposit(amount1, utxos[:2])
	require.NoError(r, err)

	amount2 := 0.05 + zetabitcoin.DefaultDepositorFee
	txHash2, err := r.SendToTSSFromDeployerToDeposit(amount2, utxos[2:4])
	require.NoError(r, err)

	// send a donation to the TSS address to compensate for the funds minted automatically during pool creation
	// and prevent accounting errors
	_, err = r.SendToTSSFromDeployerWithMemo(0.11, utxos[4:5], []byte(constant.DonationMessage))
	require.NoError(r, err)

	r.Logger.Info("testing if the deposit into BTC ZRC20 is successful...")

	cctx := utils.WaitCctxMinedByInboundHash(
		r.Ctx,
		txHash2.String(),
		r.CctxClient,
		r.Logger,
		r.CctxTimeout,
	)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	balance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	require.Equal(r, 1, balance.Sign(), "balance should be positive")

	// due to the high block throughput in localnet, ZetaClient might catch up slowly with the blocks
	// to optimize block header proof test, this test is directly executed here on the first deposit instead of having a separate test
	if testHeader {
		r.ProveBTCTransaction(txHash1)
	}
}

func (r *E2ERunner) SendToTSSFromDeployerToDeposit(amount float64, inputUTXOs []btcjson.ListUnspentResult) (
	*chainhash.Hash,
	error,
) {
	return r.SendToTSSFromDeployerWithMemo(amount, inputUTXOs, r.EVMAddress().Bytes())
}

func (r *E2ERunner) SendToTSSFromDeployerWithMemo(
	amount float64,
	inputUTXOs []btcjson.ListUnspentResult,
	memo []byte,
) (*chainhash.Hash, error) {
	btcRPC := r.BtcRPCClient
	to := r.BTCTSSAddress
	btcDeployerAddress := r.BTCDeployerAddress
	require.NotNil(r, r.BTCDeployerAddress, "btcDeployerAddress is nil")

	// prepare inputs
	inputs := make([]btcjson.TransactionInput, len(inputUTXOs))
	inputSats := btcutil.Amount(0)
	amounts := make([]float64, len(inputUTXOs))
	scriptPubkeys := make([]string, len(inputUTXOs))

	for i, utxo := range inputUTXOs {
		inputs[i] = btcjson.TransactionInput{
			Txid: utxo.TxID,
			Vout: utxo.Vout,
		}
		inputSats += btcutil.Amount(utxo.Amount * btcutil.SatoshiPerBitcoin)
		amounts[i] = utxo.Amount
		scriptPubkeys[i] = utxo.ScriptPubKey
	}

	feeSats := btcutil.Amount(0.0005 * btcutil.SatoshiPerBitcoin)
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
	r.Logger.Info("ADDRESS: %s, %s", btcDeployerAddress.EncodeAddress(), to.EncodeAddress())
	tx, err := btcRPC.CreateRawTransaction(inputs, amountMap, nil)
	require.NoError(r, err)

	// this adds a OP_RETURN + single BYTE len prefix to the data
	nullData, err := txscript.NullDataScript(memo)
	require.NoError(r, err)
	r.Logger.Info("nulldata (len %d): %x", len(nullData), nullData)
	require.NoError(r, err)
	memoOutput := wire.TxOut{Value: 0, PkScript: nullData}
	tx.TxOut = append(tx.TxOut, &memoOutput)
	tx.TxOut[1], tx.TxOut[2] = tx.TxOut[2], tx.TxOut[1]

	// make sure that TxOut[0] is sent to "to" address; TxOut[2] is change to oneself. TxOut[1] is memo.
	if !bytes.Equal(tx.TxOut[0].PkScript[2:], to.ScriptAddress()) {
		r.Logger.Info("tx.TxOut[0].PkScript: %x", tx.TxOut[0].PkScript)
		r.Logger.Info("to.ScriptAddress():   %x", to.ScriptAddress())
		r.Logger.Info("swapping txout[0] with txout[2]")
		tx.TxOut[0], tx.TxOut[2] = tx.TxOut[2], tx.TxOut[0]
	}

	r.Logger.Info("raw transaction: \n")
	for idx, txout := range tx.TxOut {
		r.Logger.Info("txout %d", idx)
		r.Logger.Info("  value: %d", txout.Value)
		r.Logger.Info("  PkScript: %x", txout.PkScript)
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
	require.NoError(r, err)
	require.True(r, signed, "btc transaction is not signed")

	txid, err := btcRPC.SendRawTransaction(stx, true)
	require.NoError(r, err)
	r.Logger.Info("txid: %+v", txid)
	_, err = r.GenerateToAddressIfLocalBitcoin(6, btcDeployerAddress)
	require.NoError(r, err)
	gtx, err := btcRPC.GetTransaction(txid)
	require.NoError(r, err)
	r.Logger.Info("rawtx confirmation: %d", gtx.BlockIndex)
	rawtx, err := btcRPC.GetRawTransactionVerbose(txid)
	require.NoError(r, err)

	depositorFee := zetabitcoin.DefaultDepositorFee
	events, err := btcobserver.FilterAndParseIncomingTx(
		btcRPC,
		[]btcjson.TxRawResult{*rawtx},
		0,
		r.BTCTSSAddress.EncodeAddress(),
		log.Logger,
		r.BitcoinParams,
		depositorFee,
	)
	require.NoError(r, err)
	r.Logger.Info("bitcoin inbound events:")
	for _, event := range events {
		r.Logger.Info("  TxHash: %s", event.TxHash)
		r.Logger.Info("  From: %s", event.FromAddress)
		r.Logger.Info("  To: %s", event.ToAddress)
		r.Logger.Info("  Amount: %f", event.Value)
		r.Logger.Info("  Memo: %x", event.MemoBytes)
	}
	return txid, nil
}

// GetBitcoinChainID gets the bitcoin chain ID from the network params
func (r *E2ERunner) GetBitcoinChainID() int64 {
	chainID, err := chains.BitcoinChainIDFromNetworkName(r.BitcoinParams.Name)
	require.NoError(r, err)
	return chainID
}

// IsLocalBitcoin returns true if the runner is running on a local bitcoin network
func (r *E2ERunner) IsLocalBitcoin() bool {
	return r.BitcoinParams.Name == chains.BitcoinRegnetParams.Name
}

// GenerateToAddressIfLocalBitcoin generates blocks to an address if the runner is interacting
// with a local bitcoin network
func (r *E2ERunner) GenerateToAddressIfLocalBitcoin(
	numBlocks int64,
	address btcutil.Address,
) ([]*chainhash.Hash, error) {
	// if not local bitcoin network, do nothing
	if r.IsLocalBitcoin() {
		return r.BtcRPCClient.GenerateToAddress(numBlocks, address, nil)
	}
	return nil, nil
}

// MineBlocksIfLocalBitcoin mines blocks on the local BTC chain at a rate of 1 blocks every 5 seconds
// and returns a channel that can be used to stop the mining
// If the chain is not local, the function does nothing
func (r *E2ERunner) MineBlocksIfLocalBitcoin() func() {
	stopChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopChan:
				return
			default:
				_, err := r.GenerateToAddressIfLocalBitcoin(1, r.BTCDeployerAddress)
				require.NoError(r, err)

				time.Sleep(3 * time.Second)
			}
		}
	}()

	return func() {
		close(stopChan)
	}
}

// ProveBTCTransaction proves that a BTC transaction is in a block header and that the block header is in ZetaChain
func (r *E2ERunner) ProveBTCTransaction(txHash *chainhash.Hash) {
	// get tx result
	btc := r.BtcRPCClient
	txResult, err := btc.GetTransaction(txHash)
	require.NoError(r, err, "should get tx result")
	require.True(r, txResult.Confirmations > 0, "tx should have already confirmed")

	txBytes, err := hex.DecodeString(txResult.Hex)
	require.NoError(r, err)

	// get the block with verbose transactions
	blockHash, err := chainhash.NewHashFromStr(txResult.BlockHash)
	require.NoError(r, err)

	blockVerbose, err := btc.GetBlockVerboseTx(blockHash)
	require.NoError(r, err, "should get block verbose tx")

	// get the block header
	header, err := btc.GetBlockHeader(blockHash)
	require.NoError(r, err, "should get block header")

	// collect all the txs in the block
	txns := []*btcutil.Tx{}
	for _, res := range blockVerbose.Tx {
		txBytes, err := hex.DecodeString(res.Hex)
		require.NoError(r, err)

		tx, err := btcutil.NewTxFromBytes(txBytes)
		require.NoError(r, err)

		txns = append(txns, tx)
	}

	// build merkle proof
	mk := bitcoin.NewMerkle(txns)
	path, index, err := mk.BuildMerkleProof(int(txResult.BlockIndex))
	require.NoError(r, err, "should build merkle proof")

	// verify merkle proof statically
	pass := bitcoin.Prove(*txHash, header.MerkleRoot, path, index)
	require.True(r, pass, "should verify merkle proof")

	// wait for block header to show up in ZetaChain
	startTime := time.Now()
	hash := header.BlockHash()
	for {
		// timeout
		reachedTimeout := time.Since(startTime) > blockHeaderBTCTimeout
		require.False(r, reachedTimeout, "timed out waiting for block header to show up in observer")

		_, err := r.LightclientClient.BlockHeader(r.Ctx, &lightclienttypes.QueryGetBlockHeaderRequest{
			BlockHash: hash.CloneBytes(),
		})
		if err != nil {
			r.Logger.Info(
				"waiting for block header to show up in observer... current hash %s; err %s",
				hash.String(),
				err.Error(),
			)
		}
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// verify merkle proof through RPC
	res, err := r.LightclientClient.Prove(r.Ctx, &lightclienttypes.QueryProveRequest{
		ChainId:   chains.BitcoinRegtest.ChainId,
		TxHash:    txHash.String(),
		BlockHash: blockHash.String(),
		Proof:     proofs.NewBitcoinProof(txBytes, path, index),
		TxIndex:   0, // bitcoin doesn't use txIndex
	})
	require.NoError(r, err)
	require.True(r, res.Valid, "txProof should be valid")

	r.Logger.Info("OK: txProof verified for inTx: %s", txHash.String())
}
