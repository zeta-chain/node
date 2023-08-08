package types_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	mc "github.com/zeta-chain/zetacore/zetaclient"
)

func TestMsgSetNodeKeys_ValidateBasic(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	kb := keyring.NewInMemory(cdc)
	path := sdk.GetConfig().GetFullBIP44Path()
	//_, err := kb.NewAccount("signerName", testdata.TestMnemonic, "", path, hd.Secp256k1)
	//require.NoError(t, err)
	_, err := kb.NewAccount(mc.GetGranteeKeyName("signerName"), testdata.TestMnemonic, "", path, hd.Secp256k1)
	require.NoError(t, err)
	granterAddress := sdk.AccAddress(crypto.AddressHash([]byte("granterAddress")))
	k := mc.NewKeysWithKeybase(kb, granterAddress, "signerName", "")
	pubKeySet, err := k.GetPubKeySet()
	assert.NoError(t, err)
	addr, err := k.GetSignerInfo().GetAddress()
	assert.NoError(t, err)
	msg := types.MsgSetNodeKeys{
		Creator:           addr.String(),
		TssSigner_Address: addr.String(),
		PubkeySet:         &pubKeySet,
	}
	err = msg.ValidateBasic()
	assert.NoError(t, err)
}
