package tss

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestKeySignManager_StartMsgSign(t *testing.T) {
	ksman := NewKeysignsTracker(zerolog.Logger{})
	ksman.StartMsgSign()
	ksman.StartMsgSign()
	ksman.StartMsgSign()
	ksman.StartMsgSign()
	require.Equal(t, int64(4), ksman.GetNumActiveMessageSigns())
}

func TestKeySignManager_EndMsgSign(t *testing.T) {
	ksman := NewKeysignsTracker(zerolog.Logger{})
	end1 := ksman.StartMsgSign()
	end2 := ksman.StartMsgSign()
	end1(true)
	end2(false)
	require.Equal(t, int64(0), ksman.GetNumActiveMessageSigns())
}
