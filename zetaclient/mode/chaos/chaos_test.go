package chaos

import (
	"bytes"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/mode"
)

const (
	seed            int64       = 12345
	perm            os.FileMode = 0600
	percentagesPath             = "percentages.json"
)

func TestNewSource(t *testing.T) {
	err := os.WriteFile(percentagesPath, []byte(`{"some": {"year": 1984}}`), perm)
	require.NoError(t, err)
	defer func() {
		err := os.Remove(percentagesPath)
		require.NoError(t, err)
	}()

	t.Run("ok (with seed)", func(t *testing.T) {
		log, source, err := newSource(mode.ChaosMode, seed, percentagesPath)
		require.NoError(t, err)
		require.NotNil(t, source)
		require.Empty(t, log)

		require.NotNil(t, source.percentages)
		require.NotNil(t, source.rand)

		require.Equal(t, 1984, source.percentages["some"]["year"])
		require.Len(t, source.percentages, 1)
		require.Len(t, source.percentages["some"], 1)
	})

	t.Run("ok (with no seed)", func(t *testing.T) {
		log, source, err := newSource(mode.ChaosMode, 0, percentagesPath)
		require.NoError(t, err)
		require.NotNil(t, source)
		require.Contains(t, log.String(), "using a random chaos seed")

		require.NotNil(t, source.percentages)
		require.NotNil(t, source.rand)

		require.Equal(t, 1984, source.percentages["some"]["year"])
		require.Len(t, source.percentages, 1)
		require.Len(t, source.percentages["some"], 1)
	})

	t.Run("invalid mode", func(t *testing.T) {
		log, source, err := newSource(mode.StandardMode, seed, percentagesPath)
		require.Error(t, err)
		require.Nil(t, source)
		require.Empty(t, log)

		require.ErrorIs(t, err, ErrNotChaosMode)
	})

	t.Run("invalid file", func(t *testing.T) {
		log, source, err := newSource(mode.ChaosMode, seed, "does_not_exist.json")
		require.Error(t, err)
		require.Nil(t, source)
		require.Empty(t, log)

		require.ErrorIs(t, err, ErrReadPercentages)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		percentagesPath := "bad_percentages.json"
		err := os.WriteFile(percentagesPath, []byte("invalid"), perm)
		require.NoError(t, err)
		defer func() {
			err := os.Remove(percentagesPath)
			require.NoError(t, err)
		}()

		log, source, err := newSource(mode.ChaosMode, seed, percentagesPath)
		require.Error(t, err)
		require.Nil(t, source)
		require.Empty(t, log)

		require.ErrorIs(t, err, ErrParsePercentages)
	})
}

func TestShouldFail(t *testing.T) {
	err := os.WriteFile(percentagesPath, []byte(`{
		"A": {
			"X": 100
		},
		"B": {
			"X": 100,
			"Y": 0,
			"Z": 50
		}
	}`), perm)
	require.NoError(t, err)
	defer func() {
		err := os.Remove(percentagesPath)
		require.NoError(t, err)
	}()

	_, source, err := newSource(mode.ChaosMode, 0, percentagesPath)
	require.NoError(t, err)
	require.NotNil(t, source)

	var shouldFail bool

	shouldFail = source.shouldFail("A", "W")
	require.False(t, shouldFail)

	shouldFail = source.shouldFail("A", "X")
	require.True(t, shouldFail)

	shouldFail = source.shouldFail("B", "X")
	require.True(t, shouldFail)

	shouldFail = source.shouldFail("B", "Y")
	require.False(t, shouldFail)

	yes, no := 0, 0
	for range 1000 {
		if source.shouldFail("B", "Z") {
			yes++
		} else {
			no++
		}
	}
	require.InDelta(t, yes, 500, 100)
	require.InDelta(t, no, 500, 100)
}

// ------------------------------------------------------------------------------------------------

func newSource(mode mode.ClientMode, seed int64, path string) (*bytes.Buffer, *Source, error) {
	log := new(bytes.Buffer)
	logger := zerolog.New(log)
	config := config.Config{
		ClientMode:           mode,
		ChaosSeed:            seed,
		ChaosPercentagesPath: path,
	}
	source, err := NewSource(logger, config)
	return log, source, err
}
