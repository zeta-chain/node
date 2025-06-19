package ante_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/app/ante"
	serverconfig "github.com/zeta-chain/node/server/config"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
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
	//		*cctxtypes.MsgAddInboundTracker,
	//		*observertypes.MsgVoteBlockHeader,
	//		*observertypes.MsgVoteTSS,
	//		*observertypes.MsgVoteBlame:
	buildTxFromMsg := func(msg sdk.Msg) sdk.Tx {
		txBuilder := app.MakeEncodingConfig(serverconfig.DefaultEVMChainID).TxConfig.NewTxBuilder()
		txBuilder.SetMsgs(msg)
		return txBuilder.GetTx()
	}
	buildAuthzTxFromMsg := func(msg sdk.Msg) sdk.Tx {
		txBuilder := app.MakeEncodingConfig(serverconfig.DefaultEVMChainID).TxConfig.NewTxBuilder()
		msgExec := authz.NewMsgExec(sample.Bech32AccAddress(), []sdk.Msg{msg})
		txBuilder.SetMsgs(&msgExec)
		return txBuilder.GetTx()
	}
	isAuthorized := func(_ string) error {
		return nil
	}
	isAuthorizedFalse := func(_ string) error {
		return errors.New("not authorized")
	}

	tests := []struct {
		name         string
		tx           sdk.Tx
		isAuthorized func(string) error
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
			"MsgAddInboundTracker",
			buildTxFromMsg(&crosschaintypes.MsgAddInboundTracker{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgAddInboundTracker}",
			buildAuthzTxFromMsg(&crosschaintypes.MsgAddInboundTracker{
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
			"MsgVoteBlame",
			buildTxFromMsg(&observertypes.MsgVoteBlame{
				Creator: sample.AccAddress(),
			}),
			isAuthorized,

			true,
		},
		{
			"MsgExec{MsgVoteBlame}",
			buildAuthzTxFromMsg(&observertypes.MsgVoteBlame{
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
