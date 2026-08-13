package tss

import (
	"context"
	"io"
	"math"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	pubkeyA = "zetapub1addwnpepqglunjrgl3qg08duxq9pf28jmvrer3crwnnfzp6m0u0yh9jk9mnn5p76utc"
	pubkeyB = "zetapub1addwnpepqwwpjwwnes7cywfkr0afme7ymk8rf5jzhn8pfr6qqvfm9v342486qsrh4f5"
	tssKey  = "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc"

	// The libp2p peer IDs the three keys above convert to. Hardcoded rather than derived in
	// the test, so a change in the conversion shows up as a failure instead of agreeing with
	// itself.
	peerIDA = "16Uiu2HAkyig859BKphpgkyiJAE3wmsFAJMf541DRbFosai13ECbK"
	peerIDB = "16Uiu2HAmPALG7YS5PNAsbHpjeg5WwSrRYqhZFkKjyPmwwEWjTP35"
)

func TestResolveTSSPeers(t *testing.T) {

	finalized := observertypes.TSS{
		TssPubkey:          tssKey,
		TssParticipantList: []string{pubkeyA, pubkeyB},
	}

	t.Run("no TSS yet, keygen record is authoritative", func(t *testing.T) {
		keygen := observertypes.Keygen{
			Status:         observertypes.KeygenStatus_PendingKeygen,
			GranteePubkeys: []string{pubkeyA, pubkeyB},
			BlockNumber:    100,
		}

		pubkeys, finalizedKey := resolveTSSPeers(keygen, observertypes.TSS{})

		assert.False(t, finalizedKey)
		assert.Equal(t, []string{pubkeyA, pubkeyB}, pubkeys)
	})

	t.Run("blanked keygen record does not empty the whitelist", func(t *testing.T) {
		// The exact state zetacore writes on an observer set change: no grantees, pending,
		// scheduled for a block that never arrives.
		blanked := observertypes.Keygen{
			Status:      observertypes.KeygenStatus_PendingKeygen,
			BlockNumber: math.MaxInt64,
		}
		require.Empty(t, blanked.GranteePubkeys)

		pubkeys, finalizedKey := resolveTSSPeers(blanked, finalized)

		assert.True(t, finalizedKey, "a finalized key must skip the ceremony")
		assert.Equal(t, []string{pubkeyA, pubkeyB}, pubkeys, "whitelist must come from the TSS participants")
	})

	t.Run("finalized key wins over a populated keygen record", func(t *testing.T) {
		// A scheduled keygen must not strand the signer either: it would otherwise block on
		// its own future block while the existing key is perfectly usable.
		scheduled := observertypes.Keygen{
			Status:         observertypes.KeygenStatus_PendingKeygen,
			GranteePubkeys: []string{pubkeyA},
			BlockNumber:    26_400_000,
		}

		pubkeys, finalizedKey := resolveTSSPeers(scheduled, finalized)

		assert.True(t, finalizedKey)
		assert.Equal(t, []string{pubkeyA, pubkeyB}, pubkeys)
	})

	t.Run("genesis TSS without participants falls back to the grantees", func(t *testing.T) {
		// x/observer/genesis.go imports a TSS verbatim, so it can carry a real key with no
		// participant list. The key must still count as finalized, or startup would try to
		// generate a replacement.
		imported := observertypes.TSS{TssPubkey: tssKey}

		keygen := observertypes.Keygen{
			Status:         observertypes.KeygenStatus_KeyGenSuccess,
			GranteePubkeys: []string{pubkeyA, pubkeyB},
		}

		pubkeys, finalizedKey := resolveTSSPeers(keygen, imported)

		assert.True(t, finalizedKey, "an imported key is still a key; do not regenerate it")
		assert.Equal(t, []string{pubkeyA, pubkeyB}, pubkeys, "must not whitelist nobody")
	})

	t.Run("successful keygen still resolves to the TSS participants", func(t *testing.T) {
		succeeded := observertypes.Keygen{
			Status:         observertypes.KeygenStatus_KeyGenSuccess,
			GranteePubkeys: []string{pubkeyA, pubkeyB},
			BlockNumber:    math.MaxInt64,
		}

		pubkeys, finalizedKey := resolveTSSPeers(succeeded, finalized)

		assert.True(t, finalizedKey)
		assert.Equal(t, finalized.TssParticipantList, pubkeys)
	})
}

// stubTSSFetcher stands in for the zetacore client. The generated mock cannot be used here:
// zetaclient/testutils/mocks imports this package, so an in-package test importing it back
// would be an import cycle.
type stubTSSFetcher struct {
	tss   observertypes.TSS
	err   error
	calls int
}

func (s *stubTSSFetcher) GetTSS(context.Context) (observertypes.TSS, error) {
	s.calls++
	return s.tss, s.err
}

func peerIDStrings(t *testing.T, peers []peer.ID) []string {
	t.Helper()

	out := make([]string, len(peers))
	for i, p := range peers {
		out[i] = p.String()
	}

	return out
}

func TestResolveWhitelist(t *testing.T) {
	logger := zerolog.New(io.Discard)

	// The exact record zetacore writes on an observer set change.
	blanked := observertypes.Keygen{
		Status:      observertypes.KeygenStatus_PendingKeygen,
		BlockNumber: math.MaxInt64,
	}

	populated := observertypes.Keygen{
		Status:         observertypes.KeygenStatus_PendingKeygen,
		GranteePubkeys: []string{pubkeyA, pubkeyB},
		BlockNumber:    100,
	}

	finalized := observertypes.TSS{
		TssPubkey:          tssKey,
		TssParticipantList: []string{pubkeyA, pubkeyB},
	}

	// The mainnet case: the record has been erased, but the key it says nothing about is still
	// there and still carries the participants to whitelist.
	t.Run("finalized key survives a blanked record", func(t *testing.T) {
		client := &stubTSSFetcher{tss: finalized}

		currentTSS, peers, keyFinalized, err := resolveWhitelist(context.Background(), client, blanked, logger)

		require.NoError(t, err)
		assert.True(t, keyFinalized, "a finalized key must skip the ceremony")
		assert.Equal(t, tssKey, currentTSS.TssPubkey)
		assert.Equal(t, []string{peerIDA, peerIDB}, peerIDStrings(t, peers))
	})

	t.Run("no TSS yet whitelists the keygen grantees", func(t *testing.T) {
		client := &stubTSSFetcher{tss: observertypes.TSS{}}

		currentTSS, peers, keyFinalized, err := resolveWhitelist(context.Background(), client, populated, logger)

		require.NoError(t, err)
		assert.False(t, keyFinalized, "a first-ever node still has to run the ceremony")
		assert.Empty(t, currentTSS.TssPubkey)
		assert.Equal(t, []string{peerIDA, peerIDB}, peerIDStrings(t, peers))
	})

	// A failed query is not the same as "there is no key", but with an intact record the
	// fallback lands on the same set anyway.
	t.Run("query failure falls back to the keygen record", func(t *testing.T) {
		client := &stubTSSFetcher{err: context.Canceled}

		_, peers, keyFinalized, err := resolveWhitelist(context.Background(), client, populated, logger)

		require.NoError(t, err)
		assert.False(t, keyFinalized)
		assert.Equal(t, []string{peerIDA, peerIDB}, peerIDStrings(t, peers))

		// retry.Retry treats context errors as non-retryable, which is what keeps this test
		// instant instead of sitting through the 5s x10 constant backoff.
		assert.Equal(t, 1, client.calls)
	})

	// Both sources empty is the state that took mainnet down. Startup must stop here rather
	// than hand go-tss an empty whitelist and let it fail as "missing bootstrap peers".
	t.Run("query failure on a blanked record stops startup", func(t *testing.T) {
		client := &stubTSSFetcher{err: context.Canceled}

		_, peers, keyFinalized, err := resolveWhitelist(context.Background(), client, blanked, logger)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no TSS peers to whitelist")
		assert.Nil(t, peers)
		assert.False(t, keyFinalized)
	})

	t.Run("a pubkey that cannot become a peer ID is reported", func(t *testing.T) {
		client := &stubTSSFetcher{tss: observertypes.TSS{}}
		keygen := observertypes.Keygen{GranteePubkeys: []string{pubkeyA, "not-a-pubkey"}}

		_, peers, _, err := resolveWhitelist(context.Background(), client, keygen, logger)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not-a-pubkey")
		assert.Nil(t, peers)
	})
}
