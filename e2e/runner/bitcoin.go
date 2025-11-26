package runner

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/memo"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	zetabtc "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
)

const (
	// BTCBlockTime is the block time for the Bitcoin mainnet and testnet
	BTCBlockTime = 10 * time.Minute

	// BTCRegnetBlockTime is the block time for the Bitcoin regnet
	BTCRegnetBlockTime = 6 * time.Second

	// BTCDepositTxFee is the fixed deposit transaction fee (0.00003 BTC) for E2E tests
	// Given one UTXO input, the deposit transaction fee rate is approximately 10 sat/vB
	BTCDepositTxFee = 0.00003

	// feeRateCap is the maximum fee rate to avoid excessive fees in E2E tests
	feeRateCap = 100
)

// ListUTXOs list the deployer's UTXOs
func (r *E2ERunner) ListUTXOs() []btcjson.ListUnspentResult {
	address, _ := r.GetBtcKeypair()

	// query UTXOs from node
	utxos, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(
		r.Ctx,
		1,
		9999999,
		[]btcutil.Address{address},
	)
	require.NoError(r, err)

	// filter big-enough UTXOs for test if running on Regtest
	if r.IsLocalBitcoin() {
		spendableAmount := 0.0
		spendableUTXOs := []btcjson.ListUnspentResult{}
		for _, utxo := range utxos {
			if utxo.Amount >= 1.0 {
				spendableAmount += utxo.Amount
				spendableUTXOs = append(spendableUTXOs, utxo)
			}
		}
		r.Logger.Info("ListUnspent(%s):", address.EncodeAddress())
		r.Logger.Info("  spendableUTXOs: %d", len(spendableUTXOs))
		r.Logger.Info("  spendableAmount: %f", spendableAmount)

		require.GreaterOrEqual(r, spendableAmount, 1.5, "not enough spendable BTC to run E2E test")

		return spendableUTXOs
	}
	require.NotEmpty(r, utxos)

	return utxos
}

// SelectUTXOs selects a subset of deployer's UTXOs to cover the given amount
func (r *E2ERunner) SelectUTXOs(amount btcutil.Amount) ([]btcjson.ListUnspentResult, btcutil.Amount) {
	var (
		allUTXOs       = r.ListUTXOs()
		regnet         = r.IsLocalBitcoin()
		selectedUTXOs  = make([]btcjson.ListUnspentResult, 0, len(allUTXOs))
		selectedAmount = btcutil.Amount(0)
	)

	// select UTXOs until we get enough funds
	for _, utxo := range allUTXOs {
		selectedUTXOs = append(selectedUTXOs, utxo)
		selectedAmount += btcutil.Amount(utxo.Amount * btcutil.SatoshiPerBitcoin)

		// in regnet, deployer owns too many UTXOs, so we stop when we have enough funds
		// in others, it is a good idea to consume all UTXOs and keep a minimum UTXO list
		if regnet && selectedAmount >= amount {
			break
		}
	}
	r.Logger.Info("selected %d UTXOs", len(selectedUTXOs))

	// check if we have enough funds
	require.GreaterOrEqual(
		r,
		selectedAmount,
		amount,
		"not enough funds in sats; wanted %d, got %d",
		amount,
		selectedAmount,
	)

	return selectedUTXOs, selectedAmount
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
	var (
		err    error
		txHash *chainhash.Hash
	)

	// deposit BTC into ZEVM
	if memo != nil {
		r.Logger.Info("⏳ depositing BTC into ZEVM with standard memo (amount: %.4f)", amount)

		// encode memo to bytes
		memoBytes, err := memo.EncodeToBytes()
		require.NoError(r, err)

		txHash, err = r.SendToTSSWithMemo(amount, memoBytes)
		require.NoError(r, err)
	} else {
		// the legacy memo layout: [20-byte receiver] + [payload]
		r.Logger.Info("⏳ depositing BTC into ZEVM with legacy memo (amount: %.4f)", amount)

		// encode 20-byte receiver, no payload
		memoBytes := r.EVMAddress().Bytes()

		txHash, err = r.SendToTSSWithMemo(amount, memoBytes)
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

	r.Logger.Info("Now donating BTC to TSS address...")

	// send a donation to the TSS address to compensate for the funds minted automatically during pool creation
	// and prevent accounting errors
	// it also serves as gas fee for the TSS to send BTC to other addresses
	_, err := r.SendToTSSWithMemo(0.11, []byte(constant.DonationMessage))
	require.NoError(r, err)
}

// DepositBTC deposits BTC from the Bitcoin node wallet into ZEVM address.
func (r *E2ERunner) DepositBTC(receiver common.Address) {
	r.Logger.Print("⏳ depositing BTC into ZEVM")
	startTime := time.Now()
	defer func() {
		r.Logger.Print("✅ BTC deposited in %s", time.Since(startTime))
	}()

	r.Logger.Info("Now depositing BTC to ZEVM address...")

	// send initial BTC to the tester ZEVM address
	amount := 1.15 + zetabtc.DefaultDepositorFee
	txHash, err := r.SendToTSSWithMemo(amount, receiver.Bytes())
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

// WithdrawBTC is a helper function to call 'withdraw' on BTCZRC20 contract with optional 'approve'
func (r *E2ERunner) WithdrawBTC(
	to btcutil.Address,
	amount *big.Int,
	revertOptions gatewayzevm.RevertOptions,
	approve bool,
) *ethtypes.Transaction {
	// ensure enough balance to cover the withdrawal
	_, gasFee, err := r.BTCZRC20.WithdrawGasFee(&bind.CallOpts{})
	require.NoError(r, err)
	minimumAmount := new(big.Int).Add(amount, gasFee)
	currentBalance, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, r.ZEVMAuth.From)
	require.NoError(r, err)
	require.Greater(
		r,
		currentBalance.Int64(),
		minimumAmount.Int64(),
		"current balance must be greater than amount + gasFee",
	)

	// approve more to cover withdraw fee
	if approve {
		r.ApproveBTCZRC20(r.GatewayZEVMAddr)
	}

	// convert P2WPKH/P2WSH addresses to uppercase, it should work without issue
	var receiverStr string
	switch to.(type) {
	case *btcutil.AddressWitnessPubKeyHash,
		*btcutil.AddressWitnessScriptHash:
		receiverStr = strings.ToUpper(to.EncodeAddress())
	default:
		receiverStr = to.EncodeAddress()
	}

	// withdraw 'amount' of BTC through ZEVM gateway
	tx, err := r.GatewayZEVM.Withdraw0(
		r.ZEVMAuth,
		[]byte(receiverStr),
		amount,
		r.BTCZRC20Addr,
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// WithdrawBTCAndWaitCCTX withdraws BTC from ZRC20 contract and waits for the CCTX to be finalized
func (r *E2ERunner) WithdrawBTCAndWaitCCTX(
	to btcutil.Address,
	amount *big.Int,
	revertOptions gatewayzevm.RevertOptions,
	expectedCCTXStatus crosschaintypes.CctxStatus,
) *btcjson.TxRawResult {
	// approve and withdraw on ZRC20 contract
	tx := r.WithdrawBTC(to, amount, revertOptions, true)

	// mine blocks if testing on regnet
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// get cctx and check status
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, expectedCCTXStatus)

	// get outbound hash according to status
	// note: the first outbound param contains the cancel tx hash for reverted cctx
	outboundHash := cctx.GetCurrentOutboundParam().Hash
	if expectedCCTXStatus == crosschaintypes.CctxStatus_Reverted {
		outboundHash = cctx.OutboundParams[0].Hash
	}

	// get bitcoin tx by outbound hash
	hash, err := chainhash.NewHashFromStr(outboundHash)
	require.NoError(r, err)

	rawTx, err := r.BtcRPCClient.GetRawTransactionVerbose(r.Ctx, hash)
	require.NoError(r, err)

	r.Logger.Info("raw tx:")
	r.Logger.Info("  TxIn: %d", len(rawTx.Vin))
	for idx, txIn := range rawTx.Vin {
		r.Logger.Info("  TxIn %d:", idx)
		r.Logger.Info("    TxID:Vout:  %s:%d", txIn.Txid, txIn.Vout)
		r.Logger.Info("    ScriptSig: %s", txIn.ScriptSig.Hex)
	}
	r.Logger.Info("  TxOut: %d", len(rawTx.Vout))
	for _, txOut := range rawTx.Vout {
		r.Logger.Info("  TxOut %d:", txOut.N)
		r.Logger.Info("    Value: %.8f", txOut.Value)
		r.Logger.Info("    ScriptPubKey: %s", txOut.ScriptPubKey.Hex)
	}

	return rawTx
}

func (r *E2ERunner) SendToTSSWithMemo(
	amount float64,
	memo []byte,
) (*chainhash.Hash, error) {
	return r.sendToAddrWithMemo(amount, r.BTCTSSAddress, memo)
}

func (r *E2ERunner) sendToAddrWithMemo(
	amount float64,
	to btcutil.Address,
	memo []byte,
) (*chainhash.Hash, error) {
	btcRPC := r.BtcRPCClient
	address, wifKey := r.GetBtcKeypair()

	// Calculate required amount with gas fee and then add
	// additional dust amount to avoid dust 'change' output
	feeSats := btcutil.Amount(BTCDepositTxFee * btcutil.SatoshiPerBitcoin)
	amountInt, err := zetabtc.GetSatoshis(amount)
	require.NoError(r, err)
	amountSats := btcutil.Amount(amountInt)
	requiredSats := amountSats + feeSats + constant.BTCWithdrawalDustAmount

	// Select UTXOs to cover the required amount
	inputUTXOs, inputSats := r.SelectUTXOs(requiredSats)

	// Create a new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add inputs
	for _, utxo := range inputUTXOs {
		txHash, err := chainhash.NewHashFromStr(utxo.TxID)
		require.NoError(r, err)
		outPoint := wire.NewOutPoint(txHash, utxo.Vout)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}

	change := inputSats - feeSats - amountSats

	// Create output to recipient
	pkScript, err := txscript.PayToAddrScript(to)
	require.NoError(r, err)
	tx.AddTxOut(wire.NewTxOut(int64(amountSats), pkScript))

	// Add memo output if provided
	if memo != nil {
		nullData, err := txscript.NullDataScript(memo)
		require.NoError(r, err)
		r.Logger.Info("nulldata (len %d): %x", len(nullData), nullData)
		memoOutput := wire.NewTxOut(0, nullData)
		tx.AddTxOut(memoOutput)
	}

	// Add change output
	changePkScript, err := txscript.PayToAddrScript(address)
	require.NoError(r, err)
	tx.AddTxOut(wire.NewTxOut(int64(change), changePkScript))

	r.Logger.Info("raw transaction: \n")
	for idx, txout := range tx.TxOut {
		r.Logger.Info("txout %d", idx)
		r.Logger.Info("  value: %d", txout.Value)
		r.Logger.Info("  PkScript: %x", txout.PkScript)
	}

	// Sign each input
	for i, utxo := range inputUTXOs {
		pkScript, err := hex.DecodeString(utxo.ScriptPubKey)
		require.NoError(r, err)

		satoshis, err := btcutil.NewAmount(utxo.Amount)
		require.NoError(r, err)
		prevOutputFetcher := txscript.NewCannedPrevOutputFetcher(pkScript, int64(satoshis))

		// Create witness
		witnessScript, err := txscript.WitnessSignature(
			tx,
			txscript.NewTxSigHashes(tx, prevOutputFetcher),
			i,
			int64(satoshis),
			pkScript,
			txscript.SigHashAll,
			wifKey.PrivKey,
			true,
		)
		require.NoError(r, err)

		// For P2WPKH, scriptSig must be empty and signature goes in witness
		tx.TxIn[i].SignatureScript = nil
		tx.TxIn[i].Witness = witnessScript
	}

	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	err = tx.Serialize(buf)
	require.NoError(r, err)
	r.Logger.Info("raw tx hex: %s", hex.EncodeToString(buf.Bytes()))

	txid, err := btcRPC.SendRawTransaction(r.Ctx, tx, true)
	require.NoError(r, err)
	r.Logger.Info("txid: %+v", txid)

	// mine 1 block to confirm the transaction
	_, err = r.GenerateToAddressIfLocalBitcoin(1, address)
	require.NoError(r, err)

	// on live networks, need to wait for the transaction to be included
	_, err = r.WaitForBitcoinTxInclusion(txid, BTCBlockTime*2)
	require.NoError(r, err)

	return txid, nil
}

func (r *E2ERunner) InscribeToTSSWithMemo(
	amount float64,
	memo []byte,
	feeRate int64,
) (*btcjson.TxRawResult, int64, string) {
	address, _ := r.GetBtcKeypair()

	// generate commit address
	builder := NewTapscriptSpender(r.BitcoinParams)
	receiver, err := builder.GenerateCommitAddress(memo)
	require.NoError(r, err)
	r.Logger.Info("received inscription commit address: %s", receiver)

	// send funds to the commit address
	commitTxHash, err := r.sendToAddrWithMemo(amount, receiver, nil)
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

	// mine 1 block to confirm the reveal transaction
	_, err = r.GenerateToAddressIfLocalBitcoin(1, address)
	require.NoError(r, err)

	// on live networks, need to wait for the transaction to be included
	rawTx, err := r.WaitForBitcoinTxInclusion(txid, BTCBlockTime*2)
	require.NoError(r, err)

	return rawTx, revealTx.TxOut[0].Value, receiver.EncodeAddress()
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

// WaitForBitcoinTxInclusion waits for the given transaction to be included either in the mempool or a block
func (r *E2ERunner) WaitForBitcoinTxInclusion(
	txHash *chainhash.Hash,
	timeout time.Duration,
) (*btcjson.TxRawResult, error) {
	start := time.Now()

	for {
		time.Sleep(5 * time.Second)
		if time.Since(start) > timeout {
			return nil, fmt.Errorf("timeout waiting for tx inclusion: %s", txHash.String())
		}

		// error may occue if the tx is not yet included
		getTxResult, err := r.BtcRPCClient.GetTransaction(r.Ctx, txHash)
		if err != nil {
			r.Logger.Error("unable to get tx %s: %s", txHash.String(), err)
			continue
		}
		r.Logger.Info("tx %s got %d confirmations", txHash.String(), getTxResult.Confirmations)

		// get raw transaction for depositor fee calculation
		rawTx, err := r.BtcRPCClient.GetRawTransactionVerbose(r.Ctx, txHash)
		if err != nil {
			r.Logger.Error("unable to get raw tx %s: %s", txHash.String(), err)
			continue
		}

		// as far as tx gets into the mempool (0 confirmations), we consider it included
		// there's no need to wait for 1 block to be able to send out subsequent transactions
		if getTxResult.Confirmations >= 0 {
			return rawTx, nil
		}
	}
}

// BitcoinCalcReceivedAmount calculates the amount received by the receiver after deducting the depositor fee
func (r *E2ERunner) BitcoinCalcReceivedAmount(depositTx *btcjson.TxRawResult, depositedAmount int64) int64 {
	// calculate depositor fee
	depositorFee, err := zetabtc.CalcDepositorFee(r.Ctx, r.BtcRPCClient, depositTx, r.BitcoinParams)
	require.NoError(r, err)

	// convert depositor fee to satoshis
	depositFeeSats, err := zetabtc.GetSatoshis(depositorFee)
	require.NoError(r, err)

	return depositedAmount - depositFeeSats
}

// BitcoinEstimateFeeRate returns the estimated fee rate in sat/vB for live networks
func (r *E2ERunner) BitcoinEstimateFeeRate(confTarget int64) uint64 {
	// if not local bitcoin network, do nothing
	if r.IsLocalBitcoin() {
		return 1
	}

	// query live network fee rate
	feeResult, err := r.BtcRPCClient.EstimateSmartFee(r.Ctx, confTarget, &btcjson.EstimateModeEconomical)
	require.NoError(r, err)
	require.Empty(r, feeResult.Errors)
	require.NotNil(r, feeResult.FeeRate)

	// sanity check
	feeRate := *feeResult.FeeRate
	require.True(r, feeRate > 0)

	satPerByte, err := zetabtc.FeeRateToSatPerByte(feeRate)
	require.NoError(r, err)

	// ensure the fee rate is within the cap
	require.LessOrEqual(r, satPerByte, feeRateCap)

	return satPerByte
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
	address, _ := r.GetBtcKeypair()

	stopChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopChan:
				return
			default:
				_, err := r.GenerateToAddressIfLocalBitcoin(1, address)
				require.NoError(r, err)

				time.Sleep(BTCRegnetBlockTime)
			}
		}
	}()

	return func() {
		close(stopChan)
	}
}
