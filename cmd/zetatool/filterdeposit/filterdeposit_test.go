package filterdeposit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetatool/config"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestCheckForCCTX(t *testing.T) {
	t.Run("no missed inbound txns found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/zeta-chain/crosschain/in_tx_hash_to_cctx_data/0x093f4ca4c1884df0fd9dd59b75979342ded29d3c9b6861644287a2e1417b9a39" {
				t.Errorf("Expected to request '/zeta-chain', got: %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			//Return CCtx
			cctx := types.CrossChainTx{}
			bytes, err := json.Marshal(cctx)
			require.NoError(t, err)
			_, err = w.Write(bytes)
			require.NoError(t, err)
		}))
		defer server.Close()

		deposits := []Deposit{{
			TxID:   "0x093f4ca4c1884df0fd9dd59b75979342ded29d3c9b6861644287a2e1417b9a39",
			Amount: uint64(657177295293237048),
		}}
		cfg := config.DefaultConfig()
		cfg.ZetaURL = server.URL
		missedInbounds := CheckForCCTX(deposits, cfg)
		require.Equal(t, 0, len(missedInbounds))
	})

	t.Run("1 missed inbound txn found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("{\n  \"code\": 5,\n  \"message\": \"not found\",\n  \"details\": [\n  ]\n}"))
			require.NoError(t, err)
		}))
		defer server.Close()

		deposits := []Deposit{{
			TxID:   "0x093f4ca4c1884df0fd9dd59b75979342ded29d3c9b6861644287a2e1417b9a39",
			Amount: uint64(657177295293237048),
		}}
		cfg := config.DefaultConfig()
		cfg.ZetaURL = server.URL
		missedInbounds := CheckForCCTX(deposits, cfg)
		require.Equal(t, 1, len(missedInbounds))
	})
}
