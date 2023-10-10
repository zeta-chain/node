package common

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	DustUTXOOffset = 2000
)

// A very special value to mark current nonce in UTXO
func NonceMarkAmount(nonce uint64) int64 {
	// #nosec G701 always in range
	return int64(nonce) + DustUTXOOffset // +2000 to avoid being a dust rejection
}

// HashToString convert hash bytes to string
func HashToString(chainID int64, blockHash []byte) (string, error) {
	if IsEVMChain(chainID) {
		return hex.EncodeToString(blockHash), nil
	} else if IsBitcoinChain(chainID) {
		hash, err := chainhash.NewHash(blockHash)
		if err != nil {
			return "", err
		}
		return hash.String(), nil
	}
	return "", fmt.Errorf("cannot convert hash to string for chain %d", chainID)
}

// StringToHash convert string to hash bytes
func StringToHash(chainID int64, hash string) ([]byte, error) {
	if IsEVMChain(chainID) {
		return ethcommon.HexToHash(hash).Bytes(), nil
	} else if IsBitcoinChain(chainID) {
		hash, err := chainhash.NewHashFromStr(hash)
		if err != nil {
			return nil, err
		}
		return hash.CloneBytes(), nil
	}
	return nil, fmt.Errorf("cannot convert hash to bytes for chain %d", chainID)
}
