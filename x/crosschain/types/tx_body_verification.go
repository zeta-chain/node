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

// VerifyInboundBody validates the tx body for an inbound
func VerifyInboundBody(
	msg MsgAddInboundTracker,
	txBytes []byte,
	chainParams observertypes.ChainParams,
	tss observertypes.QueryGetTssAddressResponse,
) error {
	// verify message against transaction body
	if chains.IsEVMChain(msg.ChainId) {
		return verifyInboundBodyEVM(msg, txBytes, chainParams, tss)
	}

	// TODO: implement verifyInboundBodyBTC
	// https://github.com/zeta-chain/node/issues/1993

	return fmt.Errorf("cannot verify inbound body for chain %d", msg.ChainId)
}

// verifyInboundBodyEVM validates the chain id and connector contract address for Zeta, ERC20 custody contract address for ERC20 and TSS address for Gas.
func verifyInboundBodyEVM(
	msg MsgAddInboundTracker,
	txBytes []byte,
	chainParams observertypes.ChainParams,
	tss observertypes.QueryGetTssAddressResponse,
) error {
	var txx ethtypes.Transaction
	err := txx.UnmarshalBinary(txBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal transaction %s", err.Error())
	}
	if txx.Hash().Hex() != msg.TxHash {
		return fmt.Errorf("invalid hash, want tx hash %s, got %s", txx.Hash().Hex(), msg.TxHash)
	}
	if txx.ChainId().Cmp(big.NewInt(msg.ChainId)) != 0 {
		return fmt.Errorf("invalid chain id, want evm chain id %d, got %d", txx.ChainId(), msg.ChainId)
	}
	switch msg.CoinType {
	case coin.CoinType_Zeta:
		// Inbound depositing ZETA interacts with connector contract
		if txx.To().Hex() != chainParams.ConnectorContractAddress {
			return fmt.Errorf("receiver is not connector contract for coin type %s", msg.CoinType)
		}
	case coin.CoinType_ERC20:
		// Inbound depositing ERC20 interacts with ERC20 custody contract
		if txx.To().Hex() != chainParams.Erc20CustodyContractAddress {
			return fmt.Errorf("receiver is not erc20Custory contract for coin type %s", msg.CoinType)
		}
	case coin.CoinType_Gas:
		// Inbound depositing Gas interacts with TSS address
		tssAddr := eth.HexToAddress(tss.Eth)
		if tssAddr == (eth.Address{}) {
			return fmt.Errorf("tss address not found")
		}
		if txx.To().Hex() != tssAddr.Hex() {
			return fmt.Errorf("receiver is not tssAddress contract for coin type %s", msg.CoinType)
		}
	default:
		return fmt.Errorf("coin type not supported %s", msg.CoinType)
	}
	return nil
}

// VerifyOutboundBody verifies the tx body for an outbound
func VerifyOutboundBody(msg MsgAddOutboundTracker, txBytes []byte, tss observertypes.QueryGetTssAddressResponse) error {
	// verify message against transaction body
	if chains.IsEVMChain(msg.ChainId) {
		return verifyOutboundBodyEVM(msg, txBytes, tss.Eth)
	} else if chains.IsBitcoinChain(msg.ChainId) {
		return verifyOutboundBodyBTC(msg, txBytes, tss.Btc)
	}
	return fmt.Errorf("cannot verify outbound body for chain %d", msg.ChainId)
}

// verifyOutboundBodyEVM validates the sender address, nonce, chain id and tx hash.
func verifyOutboundBodyEVM(msg MsgAddOutboundTracker, txBytes []byte, tssEth string) error {
	var txx ethtypes.Transaction
	err := txx.UnmarshalBinary(txBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal transaction %s", err.Error())
	}
	signer := ethtypes.NewLondonSigner(txx.ChainId())
	sender, err := ethtypes.Sender(signer, &txx)
	if err != nil {
		return fmt.Errorf("failed to recover sender %s", err.Error())
	}
	tssAddr := eth.HexToAddress(tssEth)
	if tssAddr == (eth.Address{}) {
		return fmt.Errorf("tss address not found")
	}
	if sender != tssAddr {
		return fmt.Errorf("sender is not tss address %s", sender)
	}
	if txx.ChainId().Cmp(big.NewInt(msg.ChainId)) != 0 {
		return fmt.Errorf("invalid chain id, want evm chain id %d, got %d", txx.ChainId(), msg.ChainId)
	}
	if txx.Nonce() != msg.Nonce {
		return fmt.Errorf("invalid nonce, want nonce %d, got %d", txx.Nonce(), msg.Nonce)
	}
	if txx.Hash().Hex() != msg.TxHash {
		return fmt.Errorf("invalid tx hash, want tx hash %s, got %s", txx.Hash().Hex(), msg.TxHash)
	}
	return nil
}

// verifyOutboundBodyBTC validates the SegWit sender address, nonce and chain id and tx hash
// TODO: Implement tests for the function
// https://github.com/zeta-chain/node/issues/1994
func verifyOutboundBodyBTC(msg MsgAddOutboundTracker, txBytes []byte, tssBtc string) error {
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
