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
	encodingConfig := app.MakeEncodingConfig(serverconfig.DefaultEVMChainID)

	buildTxFromMsg := func(msg sdk.Msg) sdk.Tx {
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		txBuilder.SetMsgs(msg)
		return txBuilder.GetTx()
	}
	buildAuthzTxFromMsgWithGrantee := func(grantee sdk.AccAddress, msg sdk.Msg) sdk.Tx {
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		msgExec := authz.NewMsgExec(grantee, []sdk.Msg{msg})
		txBuilder.SetMsgs(&msgExec)
		return txBuilder.GetTx()
	}
	buildMultiMsgTx := func(msgs ...sdk.Msg) sdk.Tx {
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		txBuilder.SetMsgs(msgs...)
		return txBuilder.GetTx()
	}

	// ARRANGE: fixed addresses for grantee verification tests
	observerAddr := sample.AccAddress()
	registeredGrantee := sample.Bech32AccAddress()
	attackerGrantee := sample.Bech32AccAddress()

	// isAuthorized mocks CheckSystemTxAuthorization: observer always authorized,
	// and if msgExecSigner is non-empty it must match the registeredGrantee.
	isAuthorized := func(observer string, msgExecSigner string) error {
		if msgExecSigner != "" && msgExecSigner != registeredGrantee.String() {
			return errors.New("unauthorized grantee")
		}
		return nil
	}
	isAuthorizedNoExec := func(_ string, _ string) error {
		return nil
	}
	isUnauthorized := func(_ string, _ string) error {
		return errors.New("not authorized")
	}

	tests := []struct {
		name   string
		tx     sdk.Tx
		isAuth func(string, string) error
		wantIs bool
	}{
		// --- direct system messages (no MsgExec) ---
		{
			"MsgVoteTSS unauthorized observer",
			buildTxFromMsg(&observertypes.MsgVoteTSS{
				Creator:   sample.AccAddress(),
				TssPubkey: "pubkey1234",
			}),
			isUnauthorized,
			false,
		},
		{
			"MsgVoteTSS authorized observer",
			buildTxFromMsg(&observertypes.MsgVoteTSS{
				Creator:   sample.AccAddress(),
				TssPubkey: "pubkey1234",
			}),
			isAuthorizedNoExec,
			true,
		},
		{
			"MsgVoteInbound",
			buildTxFromMsg(&crosschaintypes.MsgVoteInbound{
				Creator: sample.AccAddress(),
			}),
			isAuthorizedNoExec,
			true,
		},
		{
			"MsgVoteOutbound",
			buildTxFromMsg(&crosschaintypes.MsgVoteOutbound{
				Creator: sample.AccAddress(),
			}),
			isAuthorizedNoExec,
			true,
		},
		{
			"MsgVoteGasPrice",
			buildTxFromMsg(&crosschaintypes.MsgVoteGasPrice{
				Creator: sample.AccAddress(),
			}),
			isAuthorizedNoExec,
			true,
		},
		{
			"MsgAddOutboundTracker",
			buildTxFromMsg(&crosschaintypes.MsgAddOutboundTracker{
				Creator: sample.AccAddress(),
			}),
			isAuthorizedNoExec,
			true,
		},
		{
			"MsgAddInboundTracker",
			buildTxFromMsg(&crosschaintypes.MsgAddInboundTracker{
				Creator: sample.AccAddress(),
			}),
			isAuthorizedNoExec,
			true,
		},
		{
			"MsgVoteBlockHeader",
			buildTxFromMsg(&observertypes.MsgVoteBlockHeader{
				Creator: sample.AccAddress(),
			}),
			isAuthorizedNoExec,
			true,
		},
		{
			"MsgVoteBlame",
			buildTxFromMsg(&observertypes.MsgVoteBlame{
				Creator: sample.AccAddress(),
			}),
			isAuthorizedNoExec,
			true,
		},

		// --- non-system messages ---
		{
			"MsgSend is not system tx",
			buildTxFromMsg(&banktypes.MsgSend{}),
			isAuthorizedNoExec,
			false,
		},
		{
			"MsgCreateValidator is not system tx",
			buildTxFromMsg(&stakingtypes.MsgCreateValidator{}),
			isAuthorizedNoExec,
			false,
		},
		{
			"MsgExec{MsgSend} is not system tx",
			buildAuthzTxFromMsgWithGrantee(registeredGrantee, &banktypes.MsgSend{}),
			isAuthorized,
			false,
		},

		// --- MsgExec with registered hotkey (legitimate) ---
		{
			"MsgExec{MsgVoteTSS} with registered hotkey",
			buildAuthzTxFromMsgWithGrantee(registeredGrantee, &observertypes.MsgVoteTSS{
				Creator:   observerAddr,
				TssPubkey: "pubkey1234",
			}),
			isAuthorized,
			true,
		},
		{
			"MsgExec{MsgVoteInbound} with registered hotkey",
			buildAuthzTxFromMsgWithGrantee(registeredGrantee, &crosschaintypes.MsgVoteInbound{
				Creator: observerAddr,
			}),
			isAuthorized,
			true,
		},
		{
			"MsgExec{MsgVoteOutbound} with registered hotkey",
			buildAuthzTxFromMsgWithGrantee(registeredGrantee, &crosschaintypes.MsgVoteOutbound{
				Creator: observerAddr,
			}),
			isAuthorized,
			true,
		},
		{
			"MsgExec{MsgAddOutboundTracker} with registered hotkey",
			buildAuthzTxFromMsgWithGrantee(registeredGrantee, &crosschaintypes.MsgAddOutboundTracker{
				Creator: observerAddr,
			}),
			isAuthorized,
			true,
		},
		{
			"MsgExec{MsgAddInboundTracker} with registered hotkey",
			buildAuthzTxFromMsgWithGrantee(registeredGrantee, &crosschaintypes.MsgAddInboundTracker{
				Creator: observerAddr,
			}),
			isAuthorized,
			true,
		},
		{
			"MsgExec{MsgVoteBlockHeader} with registered hotkey",
			buildAuthzTxFromMsgWithGrantee(registeredGrantee, &observertypes.MsgVoteBlockHeader{
				Creator: observerAddr,
			}),
			isAuthorized,
			true,
		},
		{
			"MsgExec{MsgVoteBlame} with registered hotkey",
			buildAuthzTxFromMsgWithGrantee(registeredGrantee, &observertypes.MsgVoteBlame{
				Creator: observerAddr,
			}),
			isAuthorized,
			true,
		},

		// --- MsgExec attack: grantee is not registered hotkey ---
		{
			"MsgExec{MsgVoteTSS} with attacker grantee",
			buildAuthzTxFromMsgWithGrantee(attackerGrantee, &observertypes.MsgVoteTSS{
				Creator:   observerAddr,
				TssPubkey: "pubkey1234",
			}),
			isAuthorized,
			false,
		},
		{
			"MsgExec{MsgVoteInbound} with attacker grantee",
			buildAuthzTxFromMsgWithGrantee(attackerGrantee, &crosschaintypes.MsgVoteInbound{
				Creator: observerAddr,
			}),
			isAuthorized,
			false,
		},
		{
			"MsgExec{MsgVoteGasPrice} with attacker grantee",
			buildAuthzTxFromMsgWithGrantee(attackerGrantee, &crosschaintypes.MsgVoteGasPrice{
				Creator: observerAddr,
			}),
			isAuthorized,
			false,
		},

		// --- MsgExec with unauthorized observer ---
		{
			"MsgExec{MsgVoteTSS} with unauthorized observer",
			buildAuthzTxFromMsgWithGrantee(registeredGrantee, &observertypes.MsgVoteTSS{
				Creator:   observerAddr,
				TssPubkey: "pubkey1234",
			}),
			isUnauthorized,
			false,
		},

		// --- multi-msg tx ---
		{
			"multi-msg tx is not system tx",
			buildMultiMsgTx(
				&crosschaintypes.MsgVoteInbound{Creator: sample.AccAddress()},
				&crosschaintypes.MsgVoteOutbound{Creator: sample.AccAddress()},
			),
			isAuthorizedNoExec,
			false,
		},

		// --- nested MsgExec ---
		{
			"MsgExec{MsgExec{MsgVoteTSS}} nested exec rejected",
			func() sdk.Tx {
				txBuilder := encodingConfig.TxConfig.NewTxBuilder()
				innerExec := authz.NewMsgExec(registeredGrantee, []sdk.Msg{&observertypes.MsgVoteTSS{
					Creator:   observerAddr,
					TssPubkey: "pubkey1234",
				}})
				outerExec := authz.NewMsgExec(attackerGrantee, []sdk.Msg{&innerExec})
				txBuilder.SetMsgs(&outerExec)
				return txBuilder.GetTx()
			}(),
			isAuthorized,
			false,
		},

		// --- MsgExec with multiple inner messages ---
		{
			"MsgExec with multiple inner messages rejected",
			func() sdk.Tx {
				txBuilder := encodingConfig.TxConfig.NewTxBuilder()
				msgExec := authz.NewMsgExec(registeredGrantee, []sdk.Msg{
					&observertypes.MsgVoteTSS{Creator: observerAddr, TssPubkey: "pubkey1234"},
					&observertypes.MsgVoteBlame{Creator: observerAddr},
				})
				txBuilder.SetMsgs(&msgExec)
				return txBuilder.GetTx()
			}(),
			isAuthorized,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := ante.IsSystemTx(tt.tx, tt.isAuth)
			require.Equal(t, tt.wantIs, is)
		})
	}
}
