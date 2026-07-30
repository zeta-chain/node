package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	pkgchains "github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/drain"
	"github.com/zeta-chain/node/pkg/rpc"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestLatestTSS(t *testing.T) {
	t.Run("picks highest finalized height", func(t *testing.T) {
		list := []observertypes.TSS{
			{TssPubkey: "a", FinalizedZetaHeight: 100},
			{TssPubkey: "c", FinalizedZetaHeight: 300},
			{TssPubkey: "b", FinalizedZetaHeight: 200},
		}
		got, err := latestTSS(list)
		require.NoError(t, err)
		require.Equal(t, "c", got.TssPubkey)
	})

	t.Run("errors on empty", func(t *testing.T) {
		_, err := latestTSS(nil)
		require.Error(t, err)
	})
}

func TestPickMedian(t *testing.T) {
	t.Run("returns median-indexed price", func(t *testing.T) {
		gp := &crosschaintypes.GasPrice{Prices: []uint64{10, 20, 30}, MedianIndex: 1}
		got, err := pickMedian(gp, 1)
		require.NoError(t, err)
		require.Equal(t, "20", got.String())
	})

	t.Run("errors on nil / empty / out-of-range", func(t *testing.T) {
		_, err := pickMedian(nil, 1)
		require.Error(t, err)
		_, err = pickMedian(&crosschaintypes.GasPrice{}, 1)
		require.Error(t, err)
		_, err = pickMedian(&crosschaintypes.GasPrice{Prices: []uint64{1}, MedianIndex: 5}, 1)
		require.Error(t, err)
	})
}

// mockCrosschain serves only GasPrice; every other QueryClient method is unused here.
type mockCrosschain struct {
	crosschaintypes.QueryClient
}

func (mockCrosschain) GasPrice(
	_ context.Context,
	_ *crosschaintypes.QueryGetGasPriceRequest,
	_ ...grpc.CallOption,
) (*crosschaintypes.QueryGetGasPriceResponse, error) {
	return &crosschaintypes.QueryGetGasPriceResponse{
		GasPrice: &crosschaintypes.GasPrice{Prices: []uint64{100_000}, MedianIndex: 0},
	}, nil
}

// evmRPCServer serves a minimal JSON-RPC endpoint. If ok is false it errors every request,
// modeling an unreachable/broken chain RPC.
func evmRPCServer(t *testing.T, ok bool) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ok {
			http.Error(w, "rpc down", http.StatusInternalServerError)
			return
		}
		var req struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		result := "0x0"
		switch req.Method {
		case "eth_getBalance":
			result = "0xde0b6b3a7640000" // 1e18 wei
		case "eth_getTransactionCount":
			result = "0x5"
		}
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(req.ID) + `,"result":"` + result + `"}`))
		require.NoError(t, err)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestBuildEVMTxsSkipsChainOnRPCError(t *testing.T) {
	// ARRANGE: chain A's RPC is down (balance errors), chain B's RPC is healthy.
	down := evmRPCServer(t, false)
	up := evmRPCServer(t, true)
	cfg := &config.Config{EthereumRPC: down.URL, BscRPC: up.URL}
	chainsIn := []pkgchains.Chain{
		{ChainId: 1, Network: pkgchains.Network_eth, IsExternal: true, Vm: pkgchains.Vm_evm},
		{ChainId: 56, Network: pkgchains.Network_bsc, IsExternal: true, Vm: pkgchains.Vm_evm},
	}
	zetacore := rpc.Clients{Crosschain: mockCrosschain{}}

	// ACT
	txs, err := buildEVMTxs(
		context.Background(),
		cfg,
		zetacore,
		chainsIn,
		ethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
		"0x000000000000000000000000000000000000dEaD",
		chainFilter{},
	)

	// ASSERT: the broken chain is skipped, not fatal; the healthy chain is included.
	require.NoError(t, err)
	require.Len(t, txs, 1)
	require.EqualValues(t, 56, txs[0].ChainID)
}

func TestDrainNetwork(t *testing.T) {
	require.Equal(t, drain.NetworkMainnet, drainNetwork(config.NetworkMainnet))
	require.Equal(t, drain.NetworkLocalnet, drainNetwork(config.NetworkLocalnet))
	require.Equal(t, drain.NetworkTestnet, drainNetwork(config.NetworkTestnet))
	require.Equal(t, drain.NetworkTestnet, drainNetwork(config.NetworkSignet))
}
