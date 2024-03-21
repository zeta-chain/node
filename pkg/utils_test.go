package pkg

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNonceMarkAmount(t *testing.T) {
	tests := []struct {
		name   string
		nonce  uint64
		expect int64
	}{
		{"base_case", 100, 2100},
		{"zero_nonce", 0, 2000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NonceMarkAmount(tt.nonce)
			require.Equal(t, tt.expect, result)
		})
	}
}

func TestHashToString(t *testing.T) {
	evmChainId := int64(5)
	btcChainId := int64(8332)
	unknownChainId := int64(3)
	mockEthBlockHash := []byte("0xc2339489a45f8976d45482ad6fa08751a1eae91f92d60645521ca0aff2422639")
	mockBtcBlockHash := []byte("00000000000000000002dcaa3853ac587d4cafdd0aa1fff45942ab5798f29afd")
	expectedBtcHash, err := chainhash.NewHashFromStr("00000000000000000002dcaa3853ac587d4cafdd0aa1fff45942ab5798f29afd")
	require.NoError(t, err)

	tests := []struct {
		name      string
		chainID   int64
		blockHash []byte
		expect    string
		wantErr   bool
	}{
		{"evm chain", evmChainId, mockEthBlockHash, hex.EncodeToString(mockEthBlockHash), false},
		{"btc chain", btcChainId, expectedBtcHash.CloneBytes(), expectedBtcHash.String(), false},
		{"btc chain invalid hash", btcChainId, mockBtcBlockHash, "", true},
		{"unknown chain", unknownChainId, mockEthBlockHash, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HashToString(tt.chainID, tt.blockHash)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expect, result)
			}
		})
	}
}

func TestStringToHash(t *testing.T) {
	evmChainId := int64(5)
	btcChainId := int64(8332)
	unknownChainId := int64(3)
	wrontBtcHash := "00000000000000000002dcaa3853ac587d4cafdd0aa1fff45942ab5798f29afd00000000000000000002dcaa3853ac587d4cafdd0aa1fff45942ab5798f29afd"
	expectedBtcHash, err := chainhash.NewHashFromStr("00000000000000000002dcaa3853ac587d4cafdd0aa1fff45942ab5798f29afd")
	require.NoError(t, err)

	tests := []struct {
		name    string
		chainID int64
		hash    string
		expect  []byte
		wantErr bool
	}{
		{"evm chain", evmChainId, "95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5", ethcommon.HexToHash("95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5").Bytes(), false},
		{"btc chain", btcChainId, expectedBtcHash.String(), expectedBtcHash.CloneBytes(), false},
		{"btc chain invalid hash", btcChainId, wrontBtcHash, nil, true},
		{"unknown chain", unknownChainId, "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StringToHash(tt.chainID, tt.hash)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expect, result)
			}
		})
	}
}

func TestParseAddressAndData(t *testing.T) {
	expectedShortMsgResult, err := hex.DecodeString("1a2b3c4d5e6f708192a3b4c5d6e7f808")
	require.NoError(t, err)
	tests := []struct {
		name       string
		message    string
		expectAddr ethcommon.Address
		expectData []byte
		wantErr    bool
	}{
		{"valid msg", "95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5", ethcommon.HexToAddress("95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"), []byte{}, false},
		{"empty msg", "", ethcommon.Address{}, nil, false},
		{"invalid hex", "invalidHex", ethcommon.Address{}, nil, true},
		{"short msg", "1a2b3c4d5e6f708192a3b4c5d6e7f808", ethcommon.Address{}, expectedShortMsgResult, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, data, err := ParseAddressAndData(tt.message)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectAddr, addr)
				require.Equal(t, tt.expectData, data)
			}
		})
	}
}
