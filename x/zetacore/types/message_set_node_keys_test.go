package types_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"testing"
)

func TestMsgSetNodeKeys_ValidateBasic(t *testing.T) {
	kb := keyring.NewInMemory()
	path := sdk.GetConfig().GetFullBIP44Path()
	_, err := kb.NewAccount("signerName", testdata.TestMnemonic, "", path, hd.Secp256k1)
	require.NoError(t, err)

	k := mc.NewKeysWithKeybase(kb, "signerName", "")
	pubKeySet, err := k.GetPubKeySet()
	assert.NoError(t, err)
	msg := types.MsgSetNodeKeys{
		Creator:                  k.GetSignerInfo().GetAddress().String(),
		PubkeySet:                &pubKeySet,
		ValidatorConsensusPubkey: "",
	}
	err = msg.ValidateBasic()
	assert.NoError(t, err)
}
