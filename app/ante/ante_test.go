package ante_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
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
	//      *cctxtypes.MsgGasPriceVoter,
	//		*cctxtypes.MsgVoteOnObservedInboundTx,
	//		*cctxtypes.MsgVoteOnObservedOutboundTx,
	//		*cctxtypes.MsgAddToOutTxTracker,
	//		*cctxtypes.MsgCreateTSSVoter,
	//		*observertypes.MsgAddBlockHeader,
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
			"MsgCreateTSSVoter",
			buildTxFromMsg(&crosschaintypes.MsgCreateTSSVoter{
				Creator:   sample.AccAddress(),
				TssPubkey: "pubkey1234",
			}),
			isAuthorizedFalse,
			false,
		},
		{
			"MsgCreateTSSVoter",
			buildTxFromMsg(&crosschaintypes.MsgCreateTSSVoter{
				Creator:   sample.AccAddress(),
				TssPubkey: "pubkey1234",
			}),
			isAuthorized,
			true,
		},
		{
			"MsgExec{MsgCreateTSSVoter}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgCreateTSSVoter{
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
			"MsgVoteOnObservedInboundTx",
			buildTxFromMsg(&crosschaintypes.MsgVoteOnObservedInboundTx{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgVoteOnObservedInboundTx}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgVoteOnObservedInboundTx{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},

		{
			"MsgVoteOnObservedOutboundTx",
			buildTxFromMsg(&crosschaintypes.MsgVoteOnObservedOutboundTx{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgVoteOnObservedOutboundTx}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgVoteOnObservedOutboundTx{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgAddToOutTxTracker",
			buildTxFromMsg(&crosschaintypes.MsgAddToOutTxTracker{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgAddToOutTxTracker}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgAddToOutTxTracker{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgCreateTSSVoter",
			buildTxFromMsg(&crosschaintypes.MsgCreateTSSVoter{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgCreateTSSVoter}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgCreateTSSVoter{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgAddBlockHeader",
			buildTxFromMsg(&observertypes.MsgAddBlockHeader{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgAddBlockHeader}",
			buildAuthzTxFromMsg(&observertypes.MsgAddBlockHeader{
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
			assert.Equal(t, tt.wantIs, is)
		})
	}
}
