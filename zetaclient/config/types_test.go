package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

// the relative path to the testdata directory
var TestDataDir = "../"

func Test_Clone(t *testing.T) {
	// read archived zetaclient config file
	cfg := testutils.LoadZetaclientConfig(t, TestDataDir)

	// clone the config
	clone := cfg.Clone()

	// assert that the cloned config is equal to the original config
	require.Equal(t, cfg, clone)
}
