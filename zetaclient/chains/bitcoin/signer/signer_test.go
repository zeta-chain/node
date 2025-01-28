package signer_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/outboundprocessor"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// the relative path to the testdata directory
var TestDataDir = "../../../"

type testSuite struct {
	*signer.Signer
	tss            *mocks.TSS
	client         *mocks.BitcoinClient
	zetacoreClient *mocks.ZetacoreClient
}

func newTestSuite(t *testing.T, chain chains.Chain) *testSuite {
	// mock BTC RPC client
	rpcClient := mocks.NewBitcoinClient(t)
	rpcClient.On("GetBlockCount", mock.Anything).Maybe().Return(int64(101), nil)

	// mock TSS
	var tss *mocks.TSS
	if chains.IsBitcoinMainnet(chain.ChainId) {
		tss = mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet)
	} else {
		tss = mocks.NewTSS(t).FakePubKey(testutils.TSSPubkeyAthens3)
	}

	// mock Zetacore client
	zetacoreClient := mocks.NewZetacoreClient(t).
		WithKeys(&keys.Keys{}).
		WithZetaChain()

	// create logger
	testLogger := zerolog.New(zerolog.NewTestWriter(t))
	logger := base.Logger{Std: testLogger, Compliance: testLogger}

	// create signer
	baseSigner := base.NewSigner(chain, tss, logger)
	signer := signer.New(baseSigner, rpcClient)

	return &testSuite{
		Signer:         signer,
		tss:            tss,
		client:         rpcClient,
		zetacoreClient: zetacoreClient,
	}
}

func Test_NewSigner(t *testing.T) {
	// test private key with EVM address
	// EVM: 0x236C7f53a90493Bb423411fe4117Cb4c2De71DfB
	// BTC testnet: muGe9prUBjQwEnX19zG26fVRHNi8z7kSPo
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	require.NoError(t, err)
	tss := mocks.NewTSSFromPrivateKey(t, privateKey)

	baseSigner := base.NewSigner(chains.BitcoinMainnet, tss, base.DefaultLogger())
	signer := signer.New(baseSigner, mocks.NewBitcoinClient(t))
	require.NotNil(t, signer)
}

func Test_BroadcastOutbound(t *testing.T) {
	// test cases
	tests := []struct {
		name        string
		chain       chains.Chain
		nonce       uint64
		rbfTx       bool
		skipRBFTx   bool
		failTracker bool
	}{
		{
			name:  "should successfully broadcast and include outbound",
			chain: chains.BitcoinMainnet,
			nonce: uint64(148),
		},
		{
			name:  "should successfully broadcast and include RBF outbound",
			chain: chains.BitcoinMainnet,
			nonce: uint64(148),
			rbfTx: true,
		},
		{
			name:        "should successfully broadcast and include outbound, but fail to post outbound tracker",
			chain:       chains.BitcoinMainnet,
			nonce:       uint64(148),
			failTracker: true,
		},
		{
			name:      "should skip broadcasting RBF tx",
			chain:     chains.BitcoinMainnet,
			nonce:     uint64(148),
			rbfTx:     true,
			skipRBFTx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// setup signer and observer
			s := newTestSuite(t, tt.chain)
			observer := s.getNewObserver(t)

			// load tx and result
			chainID := tt.chain.ChainId
			rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, tt.nonce)
			txResult := testutils.LoadBTCTransaction(t, TestDataDir, chainID, rawResult.Txid)
			msgTx := testutils.LoadBTCMsgTx(t, TestDataDir, chainID, rawResult.Txid)
			hash := hashFromTXID(t, rawResult.Txid)

			// mock RPC response
			s.client.On("SendRawTransaction", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil, nil)
			s.client.On("GetTransactionByStr", mock.Anything, mock.Anything).Maybe().Return(hash, txResult, nil)
			s.client.On("GetRawTransactionResult", mock.Anything, mock.Anything, mock.Anything).
				Maybe().
				Return(*rawResult, nil)

			// mock Zetacore client response
			if tt.failTracker {
				s.zetacoreClient.WithPostOutboundTracker("")
			} else {
				s.zetacoreClient.WithPostOutboundTracker("0x123")
			}

			// mock the previous tx as included
			// this is necessary to allow the 'checkTSSVin' function to pass
			observer.SetIncludedTx(tt.nonce-1, &btcjson.GetTransactionResult{
				TxID: rawResult.Vin[0].Txid,
			})

			// set a higher pending nonce so the RBF tx is not the last tx
			if tt.rbfTx && tt.skipRBFTx {
				observer.SetPendingNonce(tt.nonce + 2)
			}

			// ACT
			ctx := makeCtx(t)
			s.BroadcastOutbound(
				ctx,
				msgTx,
				tt.nonce,
				tt.rbfTx,
				cctx,
				observer,
				s.zetacoreClient,
			)

			// ASSERT
			// check if outbound is included
			gotResult := observer.GetIncludedTx(tt.nonce)
			if tt.skipRBFTx {
				require.Nil(t, gotResult)
			} else {
				require.Equal(t, txResult, gotResult)
			}
		})
	}
}

func Test_P2PH(t *testing.T) {
	// Ordinarily the private key would come from whatever storage mechanism
	// is being used, but for this example just hard code it.
	privKeyBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2" +
		"d4f8720ee63e502ee2869afab7de234b80c")
	require.NoError(t, err)

	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
	require.NoError(t, err)

	// For this example, create a fake transaction that represents what
	// would ordinarily be the real transaction that is being spent.  It
	// contains a single output that pays to address in the amount of 1 BTC.
	originTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	originTx.AddTxIn(txIn)
	pkScript, err := txscript.PayToAddrScript(addr)
	require.NoError(t, err)

	txOut := wire.NewTxOut(100000000, pkScript)
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()

	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut = wire.NewOutPoint(&originTxHash, 0)
	txIn = wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(txIn)

	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)

	// Sign the redeeming transaction.
	lookupKey := func(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
		return privKey, true, nil
	}
	// Notice that the script database parameter is nil here since it isn't
	// used.  It must be specified when pay-to-script-hash transactions are
	// being signed.
	sigScript, err := txscript.SignTxOutput(&chaincfg.MainNetParams,
		redeemTx, 0, originTx.TxOut[0].PkScript, txscript.SigHashAll,
		txscript.KeyClosure(lookupKey), nil, nil)
	require.NoError(t, err)

	redeemTx.TxIn[0].SignatureScript = sigScript

	// Prove that the transaction has been validly signed by executing the
	// script pair.
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
		txscript.ScriptStrictMultiSig |
		txscript.ScriptDiscourageUpgradableNops
	vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
		flags, nil, nil, -1, txscript.NewMultiPrevOutFetcher(nil))
	require.NoError(t, err)

	err = vm.Execute()
	require.NoError(t, err)
}

func Test_P2WPH(t *testing.T) {
	// Ordinarily the private key would come from whatever storage mechanism
	// is being used, but for this example just hard code it.
	privKeyBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2" +
		"d4f8720ee63e502ee2869afab7de234b80c")
	require.NoError(t, err)

	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	//addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
	require.NoError(t, err)

	// For this example, create a fake transaction that represents what
	// would ordinarily be the real transaction that is being spent.  It
	// contains a single output that pays to address in the amount of 1 BTC.
	originTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	originTx.AddTxIn(txIn)
	pkScript, err := txscript.PayToAddrScript(addr)
	require.NoError(t, err)
	txOut := wire.NewTxOut(100000000, pkScript)
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()

	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut = wire.NewOutPoint(&originTxHash, 0)
	txIn = wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(txIn)

	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)
	txSigHashes := txscript.NewTxSigHashes(redeemTx, txscript.NewCannedPrevOutputFetcher([]byte{}, 0))
	pkScript, err = txscript.PayToAddrScript(addr)
	require.NoError(t, err)

	{
		txWitness, err := txscript.WitnessSignature(
			redeemTx,
			txSigHashes,
			0,
			100000000,
			pkScript,
			txscript.SigHashAll,
			privKey,
			true,
		)
		require.NoError(t, err)
		redeemTx.TxIn[0].Witness = txWitness
		// Prove that the transaction has been validly signed by executing the
		// script pair.
		flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
			txscript.ScriptStrictMultiSig |
			txscript.ScriptDiscourageUpgradableNops
		vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
			flags, nil, nil, -1, txscript.NewCannedPrevOutputFetcher([]byte{}, 0))
		require.NoError(t, err)

		err = vm.Execute()
		require.NoError(t, err)
	}

	{
		witnessHash, err := txscript.CalcWitnessSigHash(
			pkScript,
			txSigHashes,
			txscript.SigHashAll,
			redeemTx,
			0,
			100000000,
		)
		require.NoError(t, err)
		sig := ecdsa.Sign(privKey, witnessHash)
		txWitness := wire.TxWitness{append(sig.Serialize(), byte(txscript.SigHashAll)), pubKeyHash}
		redeemTx.TxIn[0].Witness = txWitness

		flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
			txscript.ScriptStrictMultiSig |
			txscript.ScriptDiscourageUpgradableNops
		vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
			flags, nil, nil, -1, txscript.NewMultiPrevOutFetcher(nil))
		require.NoError(t, err)

		err = vm.Execute()
		require.NoError(t, err)
	}
}

func makeCtx(t *testing.T) context.Context {
	app := zctx.New(config.New(false), nil, zerolog.Nop())

	chain := chains.BitcoinMainnet
	btcParams := mocks.MockChainParams(chain.ChainId, 2)

	err := app.Update(
		[]chains.Chain{chain, chains.ZetaChainMainnet},
		nil,
		map[int64]*observertypes.ChainParams{
			chain.ChainId: &btcParams,
		},
		observertypes.CrosschainFlags{},
	)
	require.NoError(t, err, "unable to update app context")

	return zctx.WithAppContext(context.Background(), app)
}

// getCCTX returns a CCTX for testing
func getCCTX(t *testing.T) *crosschaintypes.CrossChainTx {
	return testutils.LoadCctxByNonce(t, 8332, 148)
}

// getNewOutboundProcessor creates a new outbound processor for testing
func getNewOutboundProcessor() *outboundprocessor.Processor {
	logger := zerolog.Logger{}
	return outboundprocessor.NewProcessor(logger)
}

// getNewObserver creates a new BTC chain observer for testing
func (s *testSuite) getNewObserver(t *testing.T) *observer.Observer {
	// prepare mock arguments to create observer
	params := mocks.MockChainParams(s.Chain().ChainId, 2)
	ts := &metrics.TelemetryServer{}

	// create in-memory db
	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	// create logger
	testLogger := zerolog.New(zerolog.NewTestWriter(t))
	logger := base.Logger{Std: testLogger, Compliance: testLogger}

	// create observer
	baseObserver, err := base.NewObserver(s.Chain(), params, s.zetacoreClient, s.tss, 100, ts, database, logger)
	require.NoError(t, err)

	ob, err := observer.New(s.Chain(), baseObserver, s.client)
	require.NoError(t, err)
	return ob
}
