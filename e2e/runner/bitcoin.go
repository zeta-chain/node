package runner

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabtc "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
)

// ListDeployerUTXOs list the deployer's UTXOs
func (r *E2ERunner) ListDeployerUTXOs() ([]btcjson.ListUnspentResult, error) {
	// query UTXOs from node
	utxos, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(
		r.Ctx,
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
		r.Ctx,
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

	if len(utxos) > signer.MaxNoOfInputsPerTx {
		utxos = utxos[:signer.MaxNoOfInputsPerTx]
	}
	return utxos, nil
}

// DepositBTCWithAmount deposits BTC into ZetaChain with a specific amount and memo
func (r *E2ERunner) DepositBTCWithAmount(amount float64, memo *memo.InboundMemo) *chainhash.Hash {
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

	// add depositor fee so that receiver gets the exact given 'amount' in ZetaChain
	amount += zetabtc.DefaultDepositorFee

	// deposit to TSS address
	var txHash *chainhash.Hash
	if memo != nil {
		txHash, err = r.DepositBTCWithStandardMemo(amount, utxos, memo)
	} else {
		txHash, err = r.DepositBTCWithLegacyMemo(amount, utxos, r.EVMAddress())
	}
	require.NoError(r, err)

	r.Logger.Info("send BTC to TSS txHash: %s", txHash.String())

	return txHash
}

// DepositBTC deposits BTC from the Bitcoin node wallet into ZetaChain.
// Note: This function only works for node wallet based deployer account.
func (r *E2ERunner) DepositBTC(receiver common.Address) {
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
		// 'Spendable' indicates whether we have the private keys to spend this output
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
	r.Logger.Info("Now sending two txs to TSS address and tester ZEVM address...")

	// send initial BTC to the tester ZEVM address
	amount := 1.15 + zetabtc.DefaultDepositorFee
	txHash, err := r.DepositBTCWithLegacyMemo(amount, utxos[:2], receiver)
	require.NoError(r, err)

	// send a donation to the TSS address to compensate for the funds minted automatically during pool creation
	// and prevent accounting errors
	// it also serves as gas fee for the TSS to send BTC to other addresses
	_, err = r.SendToTSSFromDeployerWithMemo(0.11, utxos[2:4], []byte(constant.DonationMessage))
	require.NoError(r, err)

	r.Logger.Info("testing if the deposit into BTC ZRC20 is successful...")

	cctx := utils.WaitCctxMinedByInboundHash(
		r.Ctx,
		txHash.String(),
		r.CctxClient,
		r.Logger,
		r.CctxTimeout,
	)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	balance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	require.Equal(r, 1, balance.Sign(), "balance should be positive")
}

// DepositBTCWithLegacyMemo deposits BTC from the deployer address to the TSS using legacy memo
//
// The legacy memo layout: [20-byte receiver] + [payload]
func (r *E2ERunner) DepositBTCWithLegacyMemo(
	amount float64,
	inputUTXOs []btcjson.ListUnspentResult,
	receiver common.Address,
) (*chainhash.Hash, error) {
	r.Logger.Info("⏳ depositing BTC into ZEVM with legacy memo")

	// payload is not needed for pure deposit
	memoBytes := receiver.Bytes()

	return r.SendToTSSFromDeployerWithMemo(amount, inputUTXOs, memoBytes)
}

// DepositBTCWithStandardMemo deposits BTC from the deployer address to the TSS using standard `InboundMemo` struct
func (r *E2ERunner) DepositBTCWithStandardMemo(
	amount float64,
	inputUTXOs []btcjson.ListUnspentResult,
	memoStd *memo.InboundMemo,
) (*chainhash.Hash, error) {
	r.Logger.Info("⏳ depositing BTC into ZEVM with standard memo")

	// encode memo to bytes
	memoBytes, err := memoStd.EncodeToBytes()
	require.NoError(r, err)

	return r.SendToTSSFromDeployerWithMemo(amount, inputUTXOs, memoBytes)
}

func (r *E2ERunner) SendToTSSFromDeployerWithMemo(
	amount float64,
	inputUTXOs []btcjson.ListUnspentResult,
	memo []byte,
) (*chainhash.Hash, error) {
	return r.sendToAddrFromDeployerWithMemo(amount, r.BTCTSSAddress, inputUTXOs, memo)
}

func (r *E2ERunner) sendToAddrFromDeployerWithMemo(
	amount float64,
	to btcutil.Address,
	inputUTXOs []btcjson.ListUnspentResult,
	memo []byte,
) (*chainhash.Hash, error) {
	btcRPC := r.BtcRPCClient
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

	// use static fee 0.0005 BTC to calculate change
	feeSats := btcutil.Amount(0.0005 * btcutil.SatoshiPerBitcoin)
	amountInt, err := zetabtc.GetSatoshis(amount)
	require.NoError(r, err)
	amountSats := btcutil.Amount(amountInt)
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
	tx, err := btcRPC.CreateRawTransaction(r.Ctx, inputs, amountMap, nil)
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

	stx, signed, err := btcRPC.SignRawTransactionWithWallet2(r.Ctx, tx, inputsForSign)
	require.NoError(r, err)
	require.True(r, signed, "btc transaction is not signed")

	txid, err := btcRPC.SendRawTransaction(r.Ctx, stx, true)
	require.NoError(r, err)
	r.Logger.Info("txid: %+v", txid)
	_, err = r.GenerateToAddressIfLocalBitcoin(6, btcDeployerAddress)
	require.NoError(r, err)
	gtx, err := btcRPC.GetTransaction(r.Ctx, txid)
	require.NoError(r, err)
	r.Logger.Info("rawtx confirmation: %d", gtx.BlockIndex)
	rawtx, err := btcRPC.GetRawTransactionVerbose(r.Ctx, txid)
	require.NoError(r, err)

	events, err := btcobserver.FilterAndParseIncomingTx(
		r.Ctx,
		btcRPC,
		[]btcjson.TxRawResult{*rawtx},
		0,
		r.BTCTSSAddress.EncodeAddress(),
		log.Logger,
		r.BitcoinParams,
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

// InscribeToTSSFromDeployerWithMemo creates an inscription that is sent to the tss address with the corresponding memo
func (r *E2ERunner) InscribeToTSSFromDeployerWithMemo(
	amount float64,
	memo []byte,
	feeRate int64,
) (*chainhash.Hash, int64) {
	// list deployer utxos
	utxos, err := r.ListDeployerUTXOs()
	require.NoError(r, err)

	// generate commit address
	builder := NewTapscriptSpender(r.BitcoinParams)
	receiver, err := builder.GenerateCommitAddress(memo)
	require.NoError(r, err)
	r.Logger.Info("received inscription commit address: %s", receiver)

	// send funds to the commit address
	commitTxHash, err := r.sendToAddrFromDeployerWithMemo(amount, receiver, utxos, nil)
	require.NoError(r, err)
	r.Logger.Info("obtained inscription commit txn hash: %s", commitTxHash.String())

	// parameters to build the reveal transaction
	commitOutputIdx := uint32(0)
	commitAmount, err := zetabtc.GetSatoshis(amount)
	require.NoError(r, err)

	// build the reveal transaction to spend above funds
	revealTx, err := builder.BuildRevealTxn(
		r.BTCTSSAddress,
		wire.OutPoint{
			Hash:  *commitTxHash,
			Index: commitOutputIdx,
		},
		commitAmount,
		feeRate,
	)
	require.NoError(r, err)

	// submit the reveal transaction
	txid, err := r.BtcRPCClient.SendRawTransaction(r.Ctx, revealTx, true)
	require.NoError(r, err)
	r.Logger.Info("reveal txid: %s", txid.String())

	return txid, revealTx.TxOut[0].Value
}

// GetBitcoinChainID gets the bitcoin chain ID from the network params
func (r *E2ERunner) GetBitcoinChainID() int64 {
	chainID, err := chains.BitcoinChainIDFromNetworkName(r.BitcoinParams.Name)
	require.NoError(r, err)
	return chainID
}

// IsLocalBitcoin returns true if the runner is running on a local bitcoin network
func (r *E2ERunner) IsLocalBitcoin() bool {
	return r.BitcoinParams.Name == chaincfg.RegressionNetParams.Name
}

// GenerateToAddressIfLocalBitcoin generates blocks to an address if the runner is interacting
// with a local bitcoin network
func (r *E2ERunner) GenerateToAddressIfLocalBitcoin(
	numBlocks int64,
	address btcutil.Address,
) ([]*chainhash.Hash, error) {
	// if not local bitcoin network, do nothing
	if r.IsLocalBitcoin() {
		return r.BtcRPCClient.GenerateToAddress(r.Ctx, numBlocks, address, nil)
	}
	return nil, nil
}

// QueryOutboundReceiverAndAmount queries the outbound receiver and amount (in satoshis) from the given txid
func (r *E2ERunner) QueryOutboundReceiverAndAmount(txid string) (string, int64) {
	txHash, err := chainhash.NewHashFromStr(txid)
	require.NoError(r, err)

	// query outbound raw transaction
	revertTx, err := r.BtcRPCClient.GetRawTransaction(r.Ctx, txHash)
	require.NoError(r, err, revertTx)
	require.True(r, len(revertTx.MsgTx().TxOut) >= 2, "bitcoin outbound must have at least two outputs")

	// parse receiver address from pkScript
	txOutput := revertTx.MsgTx().TxOut[1]
	pkScript := txOutput.PkScript
	receiver, err := zetabtc.DecodeScriptP2WPKH(hex.EncodeToString(pkScript), r.BitcoinParams)
	require.NoError(r, err)

	return receiver, txOutput.Value
}

// MineBlocksIfLocalBitcoin mines blocks on the local BTC chain at a rate of 1 blocks every 5 seconds
// and returns a channel that can be used to stop the mining
// If the chain is not local, the function does nothing
func (r *E2ERunner) MineBlocksIfLocalBitcoin() func() {
	require.NotNil(r, r.BTCDeployerAddress, "E2ERunner.BTCDeployerAddress is nil")

	stopChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopChan:
				return
			default:
				_, err := r.GenerateToAddressIfLocalBitcoin(1, r.BTCDeployerAddress)
				require.NoError(r, err)

				time.Sleep(6 * time.Second)
			}
		}
	}()

	return func() {
		close(stopChan)
	}
}
