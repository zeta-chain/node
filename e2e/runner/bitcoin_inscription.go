package runner

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// TapscriptSpender is a utility struct that helps create Taproot address and reveal transaction
type TapscriptSpender struct {
	// internalKey is a local-generated private key used for signing the Taproot script path.
	internalKey *btcec.PrivateKey

	// taprootOutputKey is the Taproot output key derived from the internal key and the merkle root.
	// It is used to create Taproot addresses that can be funded.
	taprootOutputKey *btcec.PublicKey

	// taprootOutputAddr is the Taproot address derived from the taprootOutputKey.
	taprootOutputAddr *btcutil.AddressTaproot

	// tapLeaf represents the Taproot leaf node script (tapscript) that contains the embedded inscription data.
	tapLeaf txscript.TapLeaf

	// ctrlBlockBytes contains the control block data required for spending the Taproot output via the script path.
	// This includes the internal key and proof for the tapLeaf used to authenticate spending.
	ctrlBlockBytes []byte

	net *chaincfg.Params
}

// NewTapscriptSpender creates a new NewTapscriptSpender instance
func NewTapscriptSpender(net *chaincfg.Params) *TapscriptSpender {
	return &TapscriptSpender{
		net: net,
	}
}

// GenerateCommitAddress generates a Taproot commit address for the given receiver and payload
func (s *TapscriptSpender) GenerateCommitAddress(memo []byte) (*btcutil.AddressTaproot, error) {
	// OP_RETURN is a better choice for memo <= 80 bytes
	if len(memo) <= txscript.MaxDataCarrierSize {
		return nil, fmt.Errorf("OP_RETURN is a better choice for memo <= 80 bytes")
	}

	// generate internal private key, leaf script and Taproot output key
	err := s.genTaprootLeafAndKeys(memo)
	if err != nil {
		return nil, errors.Wrap(err, "genTaprootLeafAndKeys failed")
	}

	return s.taprootOutputAddr, nil
}

// BuildRevealTxn returns a signed reveal transaction that spends the commit transaction
func (s *TapscriptSpender) BuildRevealTxn(
	to btcutil.Address,
	commitTxn wire.OutPoint,
	commitAmount int64,
	feeRate int64,
) (*wire.MsgTx, error) {
	// Step 1: create tx message
	revealTx := wire.NewMsgTx(2)

	// Step 2: add input (the commit tx)
	outpoint := wire.NewOutPoint(&commitTxn.Hash, commitTxn.Index)
	revealTx.AddTxIn(wire.NewTxIn(outpoint, nil, nil))

	// Step 3: add output (to TSS)
	pkScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create receiver pkScript")
	}
	fee, err := s.estimateFee(revealTx, to, commitAmount, feeRate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to estimate fee for reveal txn")
	}
	revealTx.AddTxOut(wire.NewTxOut(commitAmount-fee, pkScript))

	// Step 4: compute the sighash for the P2TR input to be spent using script path
	commitScript, err := txscript.PayToAddrScript(s.taprootOutputAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create commit pkScript")
	}
	prevOutFetcher := txscript.NewCannedPrevOutputFetcher(commitScript, commitAmount)
	sigHashes := txscript.NewTxSigHashes(revealTx, prevOutFetcher)
	// sigHash, err := txscript.CalcTapscriptSignaturehash(
	// 	sigHashes,
	// 	txscript.SigHashDefault,
	// 	revealTx,
	// 	int(commitTxn.Index),
	// 	prevOutFetcher,
	// 	s.tapLeaf,                // this used to be the script content, but now we want to retrieve the funds in commit UTXO
	// )
	sigHash, err := txscript.CalcTaprootSignatureHash(
		sigHashes,
		txscript.SigHashDefault,
		revealTx,
		int(commitTxn.Index),
		prevOutFetcher,
		// now we don't pass the script content, instead we try taproot-spending
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate taproot sighash")
	}

	// tweak the internal key with the tapscript root hash
	tapScriptTree := txscript.AssembleTaprootScriptTree(s.tapLeaf)
	tapScriptRoot := tapScriptTree.RootNode.TapHash()
	tapTweakedPrivKey := txscript.TweakTaprootPrivKey(*s.internalKey, tapScriptRoot[:])

	// Step 5: sign the sighash with the internal key
	sig, err := schnorr.Sign(tapTweakedPrivKey, sigHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign sighash")
	}
	revealTx.TxIn[0].Witness = wire.TxWitness{sig.Serialize()}

	return revealTx, nil
}

// genTaprootLeafAndKeys generates internal private key, leaf script and Taproot output key
func (s *TapscriptSpender) genTaprootLeafAndKeys(data []byte) error {
	// generate an internal private key
	internalKey, err := btcec.NewPrivateKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate internal private key")
	}

	// generate the leaf script
	leafScript, err := genLeafScript(internalKey.PubKey(), data)
	if err != nil {
		return errors.Wrap(err, "failed to generate leaf script")
	}

	// assemble Taproot tree
	tapLeaf := txscript.NewBaseTapLeaf(leafScript)
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapLeaf)

	// compute the Taproot output key and address
	tapScriptRoot := tapScriptTree.RootNode.TapHash()
	taprootOutputKey := txscript.ComputeTaprootOutputKey(internalKey.PubKey(), tapScriptRoot[:])
	taprootOutputAddr, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(taprootOutputKey), s.net)
	if err != nil {
		return errors.Wrap(err, "failed to create Taproot address")
	}

	// construct the control block for the Taproot leaf script.
	ctrlBlock := tapScriptTree.LeafMerkleProofs[0].ToControlBlock(internalKey.PubKey())
	ctrlBlockBytes, err := ctrlBlock.ToBytes()
	if err != nil {
		return errors.Wrap(err, "failed to serialize control block")
	}

	// save generated keys, script and control block for later use
	s.internalKey = internalKey
	s.taprootOutputKey = taprootOutputKey
	s.taprootOutputAddr = taprootOutputAddr
	s.tapLeaf = tapLeaf
	s.ctrlBlockBytes = ctrlBlockBytes

	return nil
}

// estimateFee estimates the tx fee based given fee rate and estimated tx virtual size
func (s *TapscriptSpender) estimateFee(
	tx *wire.MsgTx,
	to btcutil.Address,
	amount int64,
	feeRate int64,
) (int64, error) {
	txCopy := tx.Copy()

	// add output to the copied transaction
	pkScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		return 0, err
	}
	txCopy.AddTxOut(wire.NewTxOut(amount, pkScript))

	// create 64-byte fake Schnorr signature
	sigBytes := make([]byte, 64)

	// set the witness for the first input
	txWitness := wire.TxWitness{sigBytes, s.tapLeaf.Script, s.ctrlBlockBytes}
	txCopy.TxIn[0].Witness = txWitness

	// calculate the fee based on the estimated virtual size
	fee := mempool.GetTxVirtualSize(btcutil.NewTx(txCopy)) * feeRate

	return fee, nil
}

//=================================================================================================
//=================================================================================================

// LeafScriptBuilder represents a builder for Taproot leaf scripts
type LeafScriptBuilder struct {
	script txscript.ScriptBuilder
}

// NewLeafScriptBuilder initializes a new LeafScriptBuilder with a public key and `OP_CHECKSIG`
func NewLeafScriptBuilder(pubKey *btcec.PublicKey) *LeafScriptBuilder {
	builder := txscript.NewScriptBuilder()
	builder.AddData(schnorr.SerializePubKey(pubKey))
	builder.AddOp(txscript.OP_CHECKSIG)

	return &LeafScriptBuilder{script: *builder}
}

// PushData adds a large data to the Taproot leaf script following OP_FALSE and OP_IF structure
func (b *LeafScriptBuilder) PushData(data []byte) {
	// start the inscription envelope
	b.script.AddOp(txscript.OP_FALSE)
	b.script.AddOp(txscript.OP_IF)

	// break data into chunks and push each one
	dataLen := len(data)
	for i := 0; i < dataLen; i += txscript.MaxScriptElementSize {
		if dataLen-i >= txscript.MaxScriptElementSize {
			b.script.AddData(data[i : i+txscript.MaxScriptElementSize])
		} else {
			b.script.AddData(data[i:])
		}
	}

	// end the inscription envelope
	b.script.AddOp(txscript.OP_ENDIF)
}

// Script returns the current script
func (b *LeafScriptBuilder) Script() ([]byte, error) {
	return b.script.Script()
}

// genLeafScript creates a Taproot leaf script using provided pubkey and data
func genLeafScript(pubKey *btcec.PublicKey, data []byte) ([]byte, error) {
	builder := NewLeafScriptBuilder(pubKey)
	builder.PushData(data)
	return builder.Script()
}
