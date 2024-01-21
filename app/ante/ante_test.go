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

	tests := []struct {
		name   string
		tx     sdk.Tx
		wantIs bool
	}{
		{
			"MsgCreateTSSVoter",
			buildTxFromMsg(&crosschaintypes.MsgCreateTSSVoter{
				Creator:   sample.AccAddress(),
				TssPubkey: "pubkey1234",
			}),
			true,
		},
		{
			"MsgExec{MsgCreateTSSVoter}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgCreateTSSVoter{
				Creator:   sample.AccAddress(),
				TssPubkey: "pubkey1234",
			}),
			true,
		},
		{
			"MsgSend",
			buildTxFromMsg(&banktypes.MsgSend{}),
			false,
		},
		{
			"MsgExec{MsgSend}",
			buildAuthzTxFromMsg(&banktypes.MsgSend{}),
			false,
		},
		{
			"MsgCreateValidator",
			buildTxFromMsg(&stakingtypes.MsgCreateValidator{}),
			false,
		},

		{
			"MsgVoteOnObservedInboundTx",
			buildTxFromMsg(&crosschaintypes.MsgVoteOnObservedInboundTx{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgExec{MsgVoteOnObservedInboundTx}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgVoteOnObservedInboundTx{
				Creator: sample.AccAddress(),
			}),
			true,
		},

		{
			"MsgVoteOnObservedOutboundTx",
			buildTxFromMsg(&crosschaintypes.MsgVoteOnObservedOutboundTx{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgExec{MsgVoteOnObservedOutboundTx}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgVoteOnObservedOutboundTx{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgAddToOutTxTracker",
			buildTxFromMsg(&crosschaintypes.MsgAddToOutTxTracker{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgExec{MsgAddToOutTxTracker}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgAddToOutTxTracker{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgCreateTSSVoter",
			buildTxFromMsg(&crosschaintypes.MsgCreateTSSVoter{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgExec{MsgCreateTSSVoter}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgCreateTSSVoter{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgAddBlockHeader",
			buildTxFromMsg(&observertypes.MsgAddBlockHeader{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgExec{MsgAddBlockHeader}",
			buildAuthzTxFromMsg(&observertypes.MsgAddBlockHeader{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgAddBlameVote",
			buildTxFromMsg(&observertypes.MsgAddBlameVote{
				Creator: sample.AccAddress(),
			}),
			true,
		},
		{
			"MsgExec{MsgAddBlameVote}",
			buildAuthzTxFromMsg(&observertypes.MsgAddBlameVote{
				Creator: sample.AccAddress(),
			}),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := ante.IsSystemTx(tt.tx, isAuthorized)
			require.Equal(t, tt.wantIs, is)
		})
	}
}
