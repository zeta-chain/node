package types_test

import (
	"bytes"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/proofs"
	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgVoteBlockHeader_ValidateBasic(t *testing.T) {
	keeper.SetConfig(false)
	var header ethtypes.Header
	file, err := os.Open("../../../testutil/testdata/eth_header_18495266.json")
	require.NoError(t, err)
	defer file.Close()
	headerBytes := make([]byte, 4096)
	n, err := file.Read(headerBytes)
	require.NoError(t, err)
	err = header.UnmarshalJSON(headerBytes[:n])
	require.NoError(t, err)
	var buffer bytes.Buffer
	err = header.EncodeRLP(&buffer)
	require.NoError(t, err)
	headerData := proofs.NewEthereumHeader(buffer.Bytes())
	tests := []struct {
		name  string
		msg   *types.MsgVoteBlockHeader
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgVoteBlockHeader(
				"invalid_address",
				1,
				[]byte{},
				6,
				proofs.HeaderData{},
			),
			error: true,
		},
		{
			name: "invalid chain id",
			msg: types.NewMsgVoteBlockHeader(
				sample.AccAddress(),
				-1,
				[]byte{},
				6,
				proofs.HeaderData{},
			),
			error: true,
		},
		{
			name: "bitcoin chain id",
			msg: types.NewMsgVoteBlockHeader(
				sample.AccAddress(),
				chains.BitcoinMainnet.ChainId,
				[]byte{},
				6,
				proofs.HeaderData{},
			),
			error: true,
		},
		{
			name: "invalid header",
			msg: types.NewMsgVoteBlockHeader(
				sample.AccAddress(),
				5,
				sample.Hash().Bytes(),
				6,
				proofs.HeaderData{},
			),
			error: true,
		},
		{
			name: "invalid blockHash length",
			msg: types.NewMsgVoteBlockHeader(
				sample.AccAddress(),
				5,
				sample.Hash().Bytes()[:31],
				6,
				proofs.HeaderData{},
			),
			error: true,
		},
		{
			name: "valid",
			msg: types.NewMsgVoteBlockHeader(
				sample.AccAddress(),
				5,
				header.Hash().Bytes(),
				18495266,
				headerData,
			),
			error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.error {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgVoteBlockHeader_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgVoteBlockHeader
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgVoteBlockHeader{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgVoteBlockHeader{
				Creator: "invalid",
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgVoteBlockHeader_Type(t *testing.T) {
	msg := types.MsgVoteBlockHeader{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgVoteBlockHeader, msg.Type())
}

func TestMsgVoteBlockHeader_Route(t *testing.T) {
	msg := types.MsgVoteBlockHeader{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgVoteBlockHeader_GetSignBytes(t *testing.T) {
	msg := types.MsgVoteBlockHeader{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}

func TestMsgVoteBlockHeader_Digest(t *testing.T) {
	msg := types.MsgVoteBlockHeader{
		Creator: sample.AccAddress(),
	}

	digest := msg.Digest()
	msg.Creator = ""
	expectedDigest := crypto.Keccak256Hash([]byte(msg.String()))
	require.Equal(t, expectedDigest.Hex(), digest)
}
