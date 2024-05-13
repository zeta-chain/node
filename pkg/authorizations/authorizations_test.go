package authorizations_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/authorizations"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
)

func TestGetRequiredPolicy(t *testing.T) {
	tt := []struct {
		name     string
		msgURl   string
		expected authoritytypes.PolicyType
	}{
		{
			name:     "Admin policy",
			msgURl:   "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
			expected: authoritytypes.PolicyType_groupAdmin,
		},
		{
			name:     "Operational policy",
			msgURl:   "/zetachain.zetacore.crosschain.MsgUpdateRateLimiterFlags",
			expected: authoritytypes.PolicyType_groupOperational,
		},
		{
			name:     "Emergency policy",
			msgURl:   "/zetachain.zetacore.crosschain.MsgAddToInTxTracker",
			expected: authoritytypes.PolicyType_groupEmergency,
		},
		{
			name:     "No policy",
			msgURl:   "/zetachain.zetacore.crosschain.MsgNoPolicy",
			expected: authoritytypes.PolicyType_emptyPolicyType,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, authorizations.GetRequiredPolicy(tc.msgURl))
		})
	}
}

func TestCheckPolicyList(t *testing.T) {
	tt := []struct {
		name      string
		msgURl    string
		msgList   []string
		assertion require.BoolAssertionFunc
	}{
		{
			name:      "Found",
			msgURl:    "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
			msgList:   authorizations.AdminPolicyMessageList,
			assertion: require.True,
		},
		{
			name:      "Not found",
			msgURl:    "/zetachain.zetacore.crosschain.MsgNoPolicy",
			msgList:   authorizations.AdminPolicyMessageList,
			assertion: require.False,
		},
		{
			name:      "Not found in wrong list",
			msgURl:    "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
			msgList:   authorizations.OperationalPolicyMessageList,
			assertion: require.False,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.assertion(t, authorizations.CheckPolicyList(tc.msgURl, tc.msgList))
		})
	}
}
