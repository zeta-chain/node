package mode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	require := require.New(t)

	mode, err := New("invalid")
	require.Error(err)
	require.ErrorIs(err, ErrInvalidModeString)
	require.Equal(InvalidMode, mode)

	mode, err = New("standard")
	require.NoError(err)
	require.Equal(StandardMode, mode)

	mode, err = New("dry")
	require.NoError(err)
	require.Equal(DryMode, mode)

	mode, err = New("chaos")
	require.NoError(err)
	require.Equal(ChaosMode, mode)
}

func TestString(t *testing.T) {
	require := require.New(t)
	require.Equal("standard", StandardMode.String())
	require.Equal("dry", DryMode.String())
	require.Equal("chaos", ChaosMode.String())
	require.Equal("invalid mode: 10", ClientMode(10).String())
}

func TestIsDryMode(t *testing.T) {
	require := require.New(t)
	require.False(StandardMode.IsDryMode())
	require.True(DryMode.IsDryMode())
	require.False(ChaosMode.IsDryMode())
}

func TestIsChaosMode(t *testing.T) {
	require := require.New(t)
	require.False(StandardMode.IsChaosMode())
	require.False(DryMode.IsChaosMode())
	require.True(ChaosMode.IsChaosMode())
}
