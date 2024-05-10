package ante_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/app/ante"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

var _ sdk.AnteHandler = (&MockAnteHandler{}).AnteHandle

// MockAnteHandler mocks an AnteHandler
type MockAnteHandler struct {
	WasCalled bool
	CalledCtx sdk.Context
}

// AnteHandle implements AnteHandler
func (mah *MockAnteHandler) AnteHandle(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
	mah.WasCalled = true
	mah.CalledCtx = ctx
	return ctx, nil
}

func TestIsSystemTx(t *testing.T) {
	// system tx types:
	//      *cctxtypes.MsgVoteGasPrice,
	//		*cctxtypes.MsgVoteInbound,
	//		*cctxtypes.MsgVoteOutbound,
	//		*cctxtypes.MsgAddOutboundTracker,
	//		*observertypes.MsgVoteBlockHeader,
	//		*observertypes.MsgVoteTSS,
	//		*observertypes.MsgAddBlameVote:
	buildTxFromMsg := func(msg sdk.Msg) sdk.Tx {
		txBuilder := app.MakeEncodingConfig().TxConfig.NewTxBuilder()
		txBuilder.SetMsgs(msg)
		return txBuilder.GetTx()
	}
	buildAuthzTxFromMsg := func(msg sdk.Msg) sdk.Tx {
		txBuilder := app.MakeEncodingConfig().TxConfig.NewTxBuilder()
		msgExec := authz.NewMsgExec(sample.Bech32AccAddress(), []sdk.Msg{msg})
		txBuilder.SetMsgs(&msgExec)
		return txBuilder.GetTx()
	}
	isAuthorized := func(_ string) bool {
		return true
	}
	isAuthorizedFalse := func(_ string) bool {
		return false
	}

	tests := []struct {
		name         string
		tx           sdk.Tx
		isAuthorized func(string) bool
		wantIs       bool
	}{
		{
			"MsgVoteTSS",
			buildTxFromMsg(&observertypes.MsgVoteTSS{
				Creator:   sample.AccAddress(),
				TssPubkey: "pubkey1234",
			}),
			isAuthorizedFalse,
			false,
		},
		{
			"MsgVoteTSS",
			buildTxFromMsg(&observertypes.MsgVoteTSS{
				Creator:   sample.AccAddress(),
				TssPubkey: "pubkey1234",
			}),
			isAuthorized,
			true,
		},
		{
			"MsgExec{MsgVoteTSS}",
			buildAuthzTxFromMsg(&observertypes.MsgVoteTSS{
				Creator:   sample.AccAddress(),
				TssPubkey: "pubkey1234",
			}),
			isAuthorized,

			true,
		},
		{
			"MsgSend",
			buildTxFromMsg(&banktypes.MsgSend{}),
			isAuthorized,

			false,
		},
		{
			"MsgExec{MsgSend}",
			buildAuthzTxFromMsg(&banktypes.MsgSend{}),
			isAuthorized,

			false,
		},
		{
			"MsgCreateValidator",
			buildTxFromMsg(&stakingtypes.MsgCreateValidator{}),
			isAuthorized,

			false,
		},

		{
			"MsgVoteInbound",
			buildTxFromMsg(&crosschaintypes.MsgVoteInbound{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgVoteInbound}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgVoteInbound{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},

		{
			"MsgVoteOutbound",
			buildTxFromMsg(&crosschaintypes.MsgVoteOutbound{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgVoteOutbound}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgVoteOutbound{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgAddOutboundTracker",
			buildTxFromMsg(&crosschaintypes.MsgAddOutboundTracker{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgAddOutboundTracker}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgAddOutboundTracker{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgVoteTSS",
			buildTxFromMsg(&observertypes.MsgVoteTSS{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgVoteTSS}",
			buildAuthzTxFromMsg(&observertypes.MsgVoteTSS{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgVoteBlockHeader",
			buildTxFromMsg(&observertypes.MsgVoteBlockHeader{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgVoteBlockHeader}",
			buildAuthzTxFromMsg(&observertypes.MsgVoteBlockHeader{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgAddBlameVote",
			buildTxFromMsg(&observertypes.MsgAddBlameVote{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgAddBlameVote}",
			buildAuthzTxFromMsg(&observertypes.MsgAddBlameVote{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := ante.IsSystemTx(tt.tx, tt.isAuthorized)
			require.Equal(t, tt.wantIs, is)
		})
	}
}
