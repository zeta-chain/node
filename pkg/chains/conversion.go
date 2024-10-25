package chains

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// NonceMarkAmount uses special value to mark current nonce in UTXO
func NonceMarkAmount(nonce uint64) int64 {
	// #nosec G115 always in range
	return int64(nonce) + BtcNonceMarkOffset()
}

// StringToHash convert string to hash bytes
func StringToHash(chainID int64, hash string, additionalChains []Chain) ([]byte, error) {
	if IsEVMChain(chainID, additionalChains) {
		return ethcommon.HexToHash(hash).Bytes(), nil
	} else if IsBitcoinChain(chainID, additionalChains) {
		hash, err := chainhash.NewHashFromStr(hash)
		if err != nil {
			return nil, err
		}
		return hash.CloneBytes(), nil
	}
	return nil, fmt.Errorf("cannot convert hash to bytes for chain %d", chainID)
}
