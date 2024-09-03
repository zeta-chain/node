package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/testdata"
	"github.com/zeta-chain/node/x/authority/client/cli"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
)

func TestReadPoliciesFromFile(t *testing.T) {
	fs := testdata.TypesFiles

	policies, err := cli.ReadPoliciesFromFile(fs, "types/policies.json")
	require.NoError(t, err)

	require.Len(t, policies.Items, 3)
	require.EqualValues(t, &authoritytypes.Policy{
		PolicyType: authoritytypes.PolicyType_groupEmergency,
		Address:    "zeta1nl7550unvzyswx5ts9m338ufmfydjsz2g0xt74",
	}, policies.Items[0])
	require.EqualValues(t, &authoritytypes.Policy{
		PolicyType: authoritytypes.PolicyType_groupOperational,
		Address:    "zeta1n0rn6sne54hv7w2uu93fl48ncyqz97d3kty6sh",
	}, policies.Items[1])
	require.EqualValues(t, &authoritytypes.Policy{
		PolicyType: authoritytypes.PolicyType_groupAdmin,
		Address:    "zeta1srsq755t654agc0grpxj4y3w0znktrpr9tcdgk",
	}, policies.Items[2])
}
