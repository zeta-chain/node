package zetaclient

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

const (
	tssPubkey    = "zetapub1addwnpepqde0ztz2agdt0ss47dhdj2867ad63ju82f87a7h97memasegvnr3xehkryd"
	tssPubkeyNew = "zetapub1addwnpepqfapt52wqw6k2kv0kvkuf8u0e8l37q57ntau7qu5ppz9sh690cs9cg0yxzs"
)

// type alias for testing
type CCTX = crosschaintypes.CrossChainTx
type OutTxParam = crosschaintypes.OutboundTxParams

func SetupTest(t *testing.T) *CctxScanner {
	logger := zerolog.New(os.Stdout)
	sc, err := NewCctxScanner(nil, TempSQLiteDbPath, true, tssPubkey, &logger)
	require.NoError(t, err)
	return sc
}

func SaveNLoadNonces(t *testing.T, sc *CctxScanner, goerliNonce uint64, bsctestNonce uint64, mumbaiNonce uint64, btctestNonce uint64) {
	goerli := clienttypes.ToFirstNonceToScanSQLType(5, goerliNonce)
	bsctest := clienttypes.ToFirstNonceToScanSQLType(97, bsctestNonce)
	mumbai := clienttypes.ToFirstNonceToScanSQLType(80001, mumbaiNonce)
	btctest := clienttypes.ToFirstNonceToScanSQLType(18332, btctestNonce)
	firstNonces := []*clienttypes.FirstNonceToScanSQLType{goerli, bsctest, mumbai, btctest}
	for _, firstNonce := range firstNonces {
		dbc := sc.db.Save(firstNonce)
		require.NoError(t, dbc.Error)
	}
	err := sc.LoadDB(TempSQLiteDbPath, true)
	require.NoError(t, err)
}

// Restart the scanner by reloading the DB
func Restart(t *testing.T, tssPubkey string) *CctxScanner {
	logger := zerolog.New(os.Stdout)
	sc, err := NewCctxScanner(nil, TempSQLiteDbPath, true, tssPubkey, &logger)
	require.NoError(t, err)
	return sc
}

func AddMissedCctxBatch1(sc *CctxScanner) (map[int64]map[uint64]*CCTX, map[int64]uint64, map[int64]uint64) {
	// expected nonce map
	expNextNonceToScanRestart := make(map[int64]uint64)
	expNextNonceToScan := make(map[int64]uint64)
	missedCctxMap := make(map[int64]map[uint64]*CCTX)

	// range [0, 1000]
	cctx0 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 361}}}
	goerliMissed := []*CCTX{cctx0}
	expNextNonceToScanRestart[5] = 361
	expNextNonceToScan[5] = 1000
	missedCctxMap[5] = make(map[uint64]*CCTX)
	missedCctxMap[5][361] = cctx0

	// range [12000, 13000]
	cctx1 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 12359}}}
	cctx2 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 12007}}}
	bscMissed := []*CCTX{cctx1, cctx2}
	expNextNonceToScanRestart[97] = 12007
	expNextNonceToScan[97] = 13000
	missedCctxMap[97] = make(map[uint64]*CCTX)
	missedCctxMap[97][12359] = cctx1
	missedCctxMap[97][12007] = cctx2

	// range [4000, 5000]
	cctx3 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 4081}}}
	cctx4 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 4602}}}
	cctx5 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 4600}}}
	mumbaiMissed1 := []*CCTX{cctx3, cctx4, cctx5}
	expNextNonceToScanRestart[80001] = 4081
	expNextNonceToScan[80001] = 5000
	missedCctxMap[80001] = make(map[uint64]*CCTX)
	missedCctxMap[80001][4081] = cctx3
	missedCctxMap[80001][4602] = cctx4
	missedCctxMap[80001][4600] = cctx5

	// range [11000, 12000]
	cctx6 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 11475}}}
	cctx7 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 11292}}}
	cctx8 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 11528}}}
	mumbaiMissed2 := []*CCTX{cctx6, cctx7, cctx8}
	expNextNonceToScan[80001] = 12000
	missedCctxMap[80001][11475] = cctx6
	missedCctxMap[80001][11292] = cctx7
	missedCctxMap[80001][11528] = cctx8

	sc.addMissedPendingCctx(5, 0, 1000, goerliMissed)
	sc.addMissedPendingCctx(97, 12000, 13000, bscMissed)
	sc.addMissedPendingCctx(80001, 4000, 5000, mumbaiMissed1)
	sc.addMissedPendingCctx(80001, 11000, 12000, mumbaiMissed2)

	return missedCctxMap, expNextNonceToScanRestart, expNextNonceToScan
}

func AddMissedCctxBatch2(sc *CctxScanner, missedCctxMap map[int64]map[uint64]*CCTX) map[int64]uint64 {
	// expected nonce map
	expNextNonceToScan := make(map[int64]uint64)

	// range [60000, 61000]
	cctx0 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 60953}}}
	cctx1 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 60437}}}
	goerliMissed := []*CCTX{cctx0, cctx1}
	expNextNonceToScan[5] = 61000
	missedCctxMap[5][60953] = cctx0
	missedCctxMap[5][60437] = cctx1

	// range [14000, 15000]
	cctx2 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 14271}}}
	bscMissed := []*CCTX{cctx2}
	expNextNonceToScan[97] = 15000
	missedCctxMap[97][14271] = cctx2

	// range [23000, 24000]
	cctx3 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 23651}}}
	cctx4 := &CCTX{OutboundTxParams: []*OutTxParam{{OutboundTxTssNonce: 23494}}}
	mumbaiMissed := []*CCTX{cctx3, cctx4}
	expNextNonceToScan[80001] = 24000
	missedCctxMap[80001][23651] = cctx3
	missedCctxMap[80001][23494] = cctx4

	sc.addMissedPendingCctx(5, 60000, 61000, goerliMissed)
	sc.addMissedPendingCctx(97, 14000, 15000, bscMissed)
	sc.addMissedPendingCctx(80001, 23000, 24000, mumbaiMissed)

	return expNextNonceToScan
}

func CheckEmptyNonces(t *testing.T, sc *CctxScanner) {
	require.Equal(t, uint64(0), sc.nextNonceToScan[5])
	require.Equal(t, uint64(0), sc.nextNonceToScan[97])
	require.Equal(t, uint64(0), sc.nextNonceToScan[80001])
}

func TestScannerDB(t *testing.T) {
	sc := SetupTest(t)

	// Make sure all maps are empty
	CheckEmptyNonces(t, sc)

	// Create some entries in the DB
	SaveNLoadNonces(t, sc, 1, 41806, 17490, 138)

	// Check the DB nonces
	var firstNonces1 []clienttypes.FirstNonceToScanSQLType
	err := sc.db.Find(&firstNonces1).Error
	require.NoError(t, err)
	for _, firstNonce := range firstNonces1 {
		want := sc.nextNonceToScan[firstNonce.ID]
		have := firstNonce.FirstNonce
		require.Equal(t, want, have)
	}

	// Update entries in the DB
	SaveNLoadNonces(t, sc, 2349, 51570, 21086, 259)

	// Check the DB nonces again
	var firstNonces2 []clienttypes.FirstNonceToScanSQLType
	err = sc.db.Find(&firstNonces2).Error
	require.NoError(t, err)
	for _, firstNonce := range firstNonces2 {
		want := sc.nextNonceToScan[firstNonce.ID]
		have := firstNonce.FirstNonce
		require.Equal(t, want, have)
	}

	// Tear down
	sc.Reset("")
}

func TestScannerDBReset(t *testing.T) {
	sc := SetupTest(t)

	// Create some entries in the DB
	SaveNLoadNonces(t, sc, 1, 41806, 17490, 138)

	// Restart scanner with different tss pubkey
	sc = Restart(t, tssPubkeyNew)

	// Make sure all maps are empty again
	CheckEmptyNonces(t, sc)

	// Tear down
	sc.Reset("")
}

func TestCctxNonces(t *testing.T) {
	sc := SetupTest(t)

	// Add some missed pending cctx
	allMissedMap, expFirstNonceMapRestart, expNextNonceMap := AddMissedCctxBatch1(sc)

	// Check the next nonce to scan
	for chainID, want := range expNextNonceMap {
		have := sc.nextNonceToScan[chainID]
		require.Equal(t, want, have)
	}

	// Add some more missed pending cctx
	expNextNonceMap = AddMissedCctxBatch2(sc, allMissedMap)

	// Check the next nonce to scan
	for chainID, want := range expNextNonceMap {
		have := sc.nextNonceToScan[chainID]
		require.Equal(t, want, have) // next nonce should change
	}

	// Restart the scanner
	sc = Restart(t, tssPubkey)

	// Check the next nonce to scan again
	for chainID, want := range expFirstNonceMapRestart {
		have := sc.nextNonceToScan[chainID]
		require.Equal(t, want, have) // next nonce should fall back to first nonce after restart
	}

	// Tear down
	sc.Reset("")
}

func CheckMissedCctxByChain(t *testing.T, sc *CctxScanner, allMissedMap map[int64]map[uint64]*crosschaintypes.CrossChainTx, chainID int64) {
	chainMissed := sc.AllMissedPendingCctxByChain(chainID)
	require.Equal(t, len(allMissedMap[chainID]), len(chainMissed))
	for _, have := range chainMissed {
		want := allMissedMap[chainID][have.OutboundTxParams[0].OutboundTxTssNonce]
		require.Equal(t, *want, *have)
	}
}

func TestGetMissedPendingCctxByChain(t *testing.T) {
	sc := SetupTest(t)

	// Add some missed pending cctx
	allMissedMap, _, _ := AddMissedCctxBatch1(sc)

	// Check missed cctx list for goerli, bsc, mumbai
	CheckMissedCctxByChain(t, sc, allMissedMap, 5)
	CheckMissedCctxByChain(t, sc, allMissedMap, 97)
	CheckMissedCctxByChain(t, sc, allMissedMap, 80001)

	// Add some more missed pending cctx
	_ = AddMissedCctxBatch2(sc, allMissedMap)

	// Check missed cctx list for goerli, bsc, mumbai again
	CheckMissedCctxByChain(t, sc, allMissedMap, 5)
	CheckMissedCctxByChain(t, sc, allMissedMap, 97)
	CheckMissedCctxByChain(t, sc, allMissedMap, 80001)

	// Tear down
	sc.Reset("")
}

func TestRemoveMissedPendingCctx(t *testing.T) {
	sc := SetupTest(t)

	// Add some missed pending cctx
	_, expNextNonceMapRestart, expNextNonceMap := AddMissedCctxBatch1(sc)

	// Remove a goerli missed cctx, edge case: delete the only cctx
	sc.removeMissedPendingCctx(5, 361)
	require.Nil(t, sc.missedPendingCctx[5][361])
	require.Equal(t, expNextNonceMap[5], sc.nextNonceToScan[5]) // won't affect next nonce

	// Remove some bsc missed cctx
	sc.removeMissedPendingCctx(97, 12359)
	require.Nil(t, sc.missedPendingCctx[97][12359])
	require.Equal(t, expNextNonceMap[97], sc.nextNonceToScan[97]) // won't affect next nonce

	// Remove some mumbai missed cctx
	sc.removeMissedPendingCctx(80001, 4600)
	sc.removeMissedPendingCctx(80001, 11528)
	require.Nil(t, sc.missedPendingCctx[80001][4600])
	require.Nil(t, sc.missedPendingCctx[80001][11528])
	require.Equal(t, expNextNonceMap[80001], sc.nextNonceToScan[80001]) // won't affect next nonce

	// Restart the scanner anc check nonces
	sc = Restart(t, tssPubkey)
	require.Equal(t, expNextNonceMap[5], sc.nextNonceToScan[5])                // next nonce fall back to first nonce for goerli, , EDGE CASE: next nonce should be 1000
	require.Equal(t, expNextNonceMapRestart[97], sc.nextNonceToScan[97])       // next nonce fall back to first nonce for bsc
	require.Equal(t, expNextNonceMapRestart[80001], sc.nextNonceToScan[80001]) // next nonce fall back to first nonce for mumbai

	// Tear down
	sc.Reset("") // clean db after each test
}
