package types

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// VerifyInTxBody validates the tx body for a inbound tx
func VerifyInTxBody(
	msg MsgAddToInTxTracker,
	txBytes []byte,
	chainParams observertypes.ChainParams,
	tss observertypes.QueryGetTssAddressResponse,
) error {
	// verify message against transaction body
	if chains.IsEVMChain(msg.ChainId) {
		return VerifyInTxBodyEVM(msg, txBytes, chainParams, tss)
	} else {
		return fmt.Errorf("cannot verify inTx body for chain %d", msg.ChainId)
	}
}

// VerifyInTxBodyEVM validates the chain id and connector contract address for Zeta, ERC20 custody contract address for ERC20 and TSS address for Gas.
func VerifyInTxBodyEVM(
	msg MsgAddToInTxTracker,
	txBytes []byte,
	chainParams observertypes.ChainParams,
	tss observertypes.QueryGetTssAddressResponse,
) error {
	var txx ethtypes.Transaction
	err := txx.UnmarshalBinary(txBytes)
	if err != nil {
		return err
	}
	if txx.Hash().Hex() != msg.TxHash {
		return fmt.Errorf("want tx hash %s, got %s", txx.Hash().Hex(), msg.TxHash)
	}
	if txx.ChainId().Cmp(big.NewInt(msg.ChainId)) != 0 {
		return fmt.Errorf("want evm chain id %d, got %d", txx.ChainId(), msg.ChainId)
	}
	switch msg.CoinType {
	case coin.CoinType_Zeta:
		// Inbound depositing ZETA interacts with connector contract
		if txx.To().Hex() != chainParams.ConnectorContractAddress {
			return fmt.Errorf("receiver is not connector contract for coin type %s", msg.CoinType)
		}
		return nil
	case coin.CoinType_ERC20:
		// Inbound depositing ERC20 interacts with ERC20 custody contract
		if txx.To().Hex() != chainParams.Erc20CustodyContractAddress {
			return fmt.Errorf("receiver is not erc20Custory contract for coin type %s", msg.CoinType)
		}
		return nil
	case coin.CoinType_Gas:
		// Inbound depositing Gas interacts with TSS address
		tssAddr := eth.HexToAddress(tss.Eth)
		if tssAddr == (eth.Address{}) {
			return fmt.Errorf("tss address not found")
		}
		if txx.To().Hex() != tssAddr.Hex() {
			return fmt.Errorf("receiver is not tssAddress contract for coin type %s", msg.CoinType)
		}
		return nil
	default:
		return fmt.Errorf("coin type %s not supported", msg.CoinType)
	}
}

// VerifyOutTxBody verifies the tx body for a outbound tx
func VerifyOutTxBody(msg MsgAddToOutTxTracker, txBytes []byte, tss observertypes.QueryGetTssAddressResponse) error {
	// verify message against transaction body
	if chains.IsEVMChain(msg.ChainId) {
		return VerifyOutTxBodyEVM(msg, txBytes, tss.Eth)
	} else if chains.IsBitcoinChain(msg.ChainId) {
		return VerifyOutTxBodyBTC(msg, txBytes, tss.Btc)
	} else {
		return fmt.Errorf("cannot verify outTx body for chain %d", msg.ChainId)
	}
}

// VerifyOutTxBodyEVM validates the sender address, nonce, chain id and tx hash.
func VerifyOutTxBodyEVM(msg MsgAddToOutTxTracker, txBytes []byte, tssEth string) error {
	var txx ethtypes.Transaction
	err := txx.UnmarshalBinary(txBytes)
	if err != nil {
		return err
	}
	signer := ethtypes.NewLondonSigner(txx.ChainId())
	sender, err := ethtypes.Sender(signer, &txx)
	if err != nil {
		return err
	}
	tssAddr := eth.HexToAddress(tssEth)
	if tssAddr == (eth.Address{}) {
		return fmt.Errorf("tss address not found")
	}
	if sender != tssAddr {
		return fmt.Errorf("sender %s is not tss address", sender)
	}
	if txx.ChainId().Cmp(big.NewInt(msg.ChainId)) != 0 {
		return fmt.Errorf("want evm chain id %d, got %d", txx.ChainId(), msg.ChainId)
	}
	if txx.Nonce() != msg.Nonce {
		return fmt.Errorf("want nonce %d, got %d", txx.Nonce(), msg.Nonce)
	}
	if txx.Hash().Hex() != msg.TxHash {
		return fmt.Errorf("want tx hash %s, got %s", txx.Hash().Hex(), msg.TxHash)
	}
	return nil
}

// VerifyOutTxBodyBTC validates the SegWit sender address, nonce and chain id and tx hash
func VerifyOutTxBodyBTC(msg MsgAddToOutTxTracker, txBytes []byte, tssBtc string) error {
	if !chains.IsBitcoinChain(msg.ChainId) {
		return fmt.Errorf("not a Bitcoin chain ID %d", msg.ChainId)
	}
	tx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return err
	}
	for _, vin := range tx.MsgTx().TxIn {
		if len(vin.Witness) != 2 { // outTx is SegWit transaction for now
			return fmt.Errorf("not a SegWit transaction")
		}
		pubKey, err := btcec.ParsePubKey(vin.Witness[1], btcec.S256())
		if err != nil {
			return fmt.Errorf("failed to parse public key")
		}
		bitcoinNetParams, err := chains.BitcoinNetParamsFromChainID(msg.ChainId)
		if err != nil {
			return fmt.Errorf("failed to get Bitcoin net params, error %s", err.Error())
		}
		addrP2WPKH, err := btcutil.NewAddressWitnessPubKeyHash(
			btcutil.Hash160(pubKey.SerializeCompressed()),
			bitcoinNetParams,
		)
		if err != nil {
			return fmt.Errorf("failed to create P2WPKH address")
		}
		if addrP2WPKH.EncodeAddress() != tssBtc {
			return fmt.Errorf("sender %s is not tss address", addrP2WPKH.EncodeAddress())
		}
	}
	if len(tx.MsgTx().TxOut) < 1 {
		return fmt.Errorf("outTx should have at least one output")
	}
	if tx.MsgTx().TxOut[0].Value != chains.NonceMarkAmount(msg.Nonce) {
		return fmt.Errorf("want nonce mark %d, got %d", tx.MsgTx().TxOut[0].Value, chains.NonceMarkAmount(msg.Nonce))
	}
	if tx.MsgTx().TxHash().String() != msg.TxHash {
		return fmt.Errorf("want tx hash %s, got %s", tx.MsgTx().TxHash(), msg.TxHash)
	}
	return nil
}
