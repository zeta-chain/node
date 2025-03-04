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
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabtc "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
)

const (
	// BTCRegnetBlockTime is the block time for the Bitcoin regnet
	BTCRegnetBlockTime = 5 * time.Second
)

// ListDeployerUTXOs list the deployer's UTXOs
func (r *E2ERunner) ListDeployerUTXOs() []btcjson.ListUnspentResult {
	// query UTXOs from node
	utxos, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(
		r.Ctx,
		1,
		9999999,
		[]btcutil.Address{r.BTCDeployerAddress},
	)
	require.NoError(r, err)

	// filter big-enough UTXOs for test if running on Regtest
	if r.IsLocalBitcoin() {
		spendableAmount := 0.0
		spendableUTXOs := []btcjson.ListUnspentResult{}
		for _, utxo := range utxos {
			// 'Spendable' indicates whether we have the private keys to spend this output
			if utxo.Spendable && utxo.Amount >= 1.0 {
				spendableAmount += utxo.Amount
				spendableUTXOs = append(spendableUTXOs, utxo)
			}
		}
		r.Logger.Info("ListUnspent:")
		r.Logger.Info("  spendableUTXOs: %d", len(spendableUTXOs))
		r.Logger.Info("  spendableAmount: %f", spendableAmount)

		require.GreaterOrEqual(r, spendableAmount, 1.5, "not enough spendable BTC to run E2E test")
		require.GreaterOrEqual(r, len(spendableUTXOs), 5, "not enough spendable BTC UTXOs to run E2E test")

		return spendableUTXOs
	}
	require.NotEmpty(r, utxos)

	return utxos
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

// DepositBTCWithExactAmount deposits exact 'amount' of BTC to receiver ZEVM address
// It automatically adds the depositor fee so that the receiver gets the exact 'amount' in ZetaChain
func (r *E2ERunner) DepositBTCWithExactAmount(amount float64, memo *memo.InboundMemo) *chainhash.Hash {
	amount += zetabtc.DefaultDepositorFee

	return r.DepositBTCWithAmount(amount, memo)
}

// DepositBTCWithAmount deposits 'amount' of BTC to TSS address with the given memo
func (r *E2ERunner) DepositBTCWithAmount(amount float64, memo *memo.InboundMemo) *chainhash.Hash {
	// list deployer utxos
	utxos := r.ListDeployerUTXOs()

	var (
		err    error
		txHash *chainhash.Hash
	)

	// deposit BTC into ZEVM
	if memo != nil {
		r.Logger.Info("⏳ depositing BTC into ZEVM with standard memo")

		// encode memo to bytes
		memoBytes, err := memo.EncodeToBytes()
		require.NoError(r, err)

		txHash, err = r.SendToTSSFromDeployerWithMemo(amount, utxos[:1], memoBytes)
		require.NoError(r, err)
	} else {
		// the legacy memo layout: [20-byte receiver] + [payload]
		r.Logger.Info("⏳ depositing BTC into ZEVM with legacy memo")

		// encode 20-byte receiver, no payload
		memoBytes := r.EVMAddress().Bytes()

		txHash, err = r.SendToTSSFromDeployerWithMemo(amount, utxos[:1], memoBytes)
		require.NoError(r, err)
	}
	r.Logger.Info("deposited BTC to TSS txHash: %s", txHash.String())

	return txHash
}

// DonateBTC donates BTC from the Bitcoin node wallet to the TSS address.
func (r *E2ERunner) DonateBTC() {
	r.Logger.Info("⏳ donating BTC to TSS address")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("✅ BTC donated in %s", time.Since(startTime))
	}()

	// list deployer utxos
	utxos := r.ListDeployerUTXOs()

	r.Logger.Info("Now donating BTC to TSS address...")

	// send a donation to the TSS address to compensate for the funds minted automatically during pool creation
	// and prevent accounting errors
	// it also serves as gas fee for the TSS to send BTC to other addresses
	_, err := r.SendToTSSFromDeployerWithMemo(0.11, utxos[:2], []byte(constant.DonationMessage))
	require.NoError(r, err)
}

// DepositBTC deposits BTC from the Bitcoin node wallet into ZEVM address.
// Note: This function only works for node wallet based deployer account.
func (r *E2ERunner) DepositBTC(receiver common.Address) {
	r.Logger.Print("⏳ depositing BTC into ZEVM")
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ BTC deposited in %s", time.Since(startTime))
	}()

	// list deployer utxos
	utxos := r.ListDeployerUTXOs()
	r.Logger.Info("Now depositing BTC to ZEVM address...")

	// send initial BTC to the tester ZEVM address
	amount := 1.15 + zetabtc.DefaultDepositorFee
	txHash, err := r.SendToTSSFromDeployerWithMemo(amount, utxos[:1], receiver.Bytes())
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

	return txid, nil
}

// InscribeToTSSFromDeployerWithMemo creates an inscription that is sent to the tss address with the corresponding memo
func (r *E2ERunner) InscribeToTSSFromDeployerWithMemo(
	amount float64,
	memo []byte,
	feeRate int64,
) (*chainhash.Hash, int64) {
	// list deployer utxos
	utxos := r.ListDeployerUTXOs()

	// generate commit address
	builder := NewTapscriptSpender(r.BitcoinParams)
	receiver, err := builder.GenerateCommitAddress(memo)
	require.NoError(r, err)
	r.Logger.Info("received inscription commit address: %s", receiver)

	// send funds to the commit address
	commitTxHash, err := r.sendToAddrFromDeployerWithMemo(amount, receiver, utxos[:1], nil)
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

				time.Sleep(BTCRegnetBlockTime)
			}
		}
	}()

	return func() {
		close(stopChan)
	}
}
