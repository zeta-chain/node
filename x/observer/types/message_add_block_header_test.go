package types_test

import (
	"bytes"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgAddBlockHeader_ValidateBasic(t *testing.T) {
	keeper.SetConfig(false)
	var header ethtypes.Header
	file, err := os.Open("../../../common/testdata/eth_header_18495266.json")
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
	headerData := common.NewEthereumHeader(buffer.Bytes())
	tests := []struct {
		name  string
		msg   *types.MsgAddBlockHeader
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgAddBlockHeader(
				"invalid_address",
				1,
				[]byte{},
				6,
				common.HeaderData{},
			),
			error: true,
		},
		{
			name: "invalid chain id",
			msg: types.NewMsgAddBlockHeader(
				sample.AccAddress(),
				-1,
				[]byte{},
				6,
				common.HeaderData{},
			),
			error: true,
		},
		{
			name: "invalid header",
			msg: types.NewMsgAddBlockHeader(
				sample.AccAddress(),
				5,
				sample.Hash().Bytes(),
				6,
				common.HeaderData{},
			),
			error: true,
		},
		{
			name: "invalid blockHash length",
			msg: types.NewMsgAddBlockHeader(
				sample.AccAddress(),
				5,
				sample.Hash().Bytes()[:31],
				6,
				common.HeaderData{},
			),
			error: true,
		},
		{
			name: "valid",
			msg: types.NewMsgAddBlockHeader(
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

func TestMsgAddBlockHeader_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgAddBlockHeader
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgAddBlockHeader{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgAddBlockHeader{
				Creator: "invalid",
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				assert.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				assert.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgAddBlockHeader_Type(t *testing.T) {
	msg := types.MsgAddBlockHeader{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.TypeMsgAddBlockHeader, msg.Type())
}

func TestMsgAddBlockHeader_Route(t *testing.T) {
	msg := types.MsgAddBlockHeader{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgAddBlockHeader_GetSignBytes(t *testing.T) {
	msg := types.MsgAddBlockHeader{
		Creator: sample.AccAddress(),
	}
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}

func TestMsgAddBlockHeader_Digest(t *testing.T) {
	msg := types.MsgAddBlockHeader{
		Creator: sample.AccAddress(),
	}

	digest := msg.Digest()
	msg.Creator = ""
	expectedDigest := crypto.Keccak256Hash([]byte(msg.String()))
	assert.Equal(t, expectedDigest.Hex(), digest)
}
