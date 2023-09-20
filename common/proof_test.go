package common_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
)

func Test_IsErrorInvalidProof(t *testing.T) {
	require.False(t, common.IsErrorInvalidProof(nil))
	require.False(t, common.IsErrorInvalidProof(errors.New("foo")))
	require.True(t, common.IsErrorInvalidProof(common.NewErrInvalidProof(errors.New("foo"))))
}
