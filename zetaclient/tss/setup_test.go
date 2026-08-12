package tss

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestResolveTSSPeers(t *testing.T) {
	const (
		pubkeyA = "zetapub1addwnpepqglunjrgl3qg08duxq9pf28jmvrer3crwnnfzp6m0u0yh9jk9mnn5p76utc"
		pubkeyB = "zetapub1addwnpepqwwpjwwnes7cywfkr0afme7ymk8rf5jzhn8pfr6qqvfm9v342486qsrh4f5"
		tssKey  = "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc"
	)

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
