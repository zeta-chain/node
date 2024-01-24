package v6

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// observerKeeper prevents circular dependency
type observerKeeper interface {
	SetKeygen(ctx sdk.Context, keygen types.Keygen)
	GetKeygen(ctx sdk.Context) (val types.Keygen, found bool)
	GetTSS(ctx sdk.Context) (val types.TSS, found bool)
	StoreKey() storetypes.StoreKey
	Codec() codec.BinaryCodec
}

func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	return SetKeyGenStatus(ctx, observerKeeper)
}

func SetKeyGenStatus(ctx sdk.Context, keeper observerKeeper) error {
	keygen, found := keeper.GetKeygen(ctx)
	if !found {
		return types.ErrKeygenNotFound
	}
	if keygen.Status == types.KeygenStatus_PendingKeygen {
		tss, foundTss := keeper.GetTSS(ctx)
		if !foundTss {
			return types.ErrTssNotFound
		}
		keygen.Status = types.KeygenStatus_KeyGenSuccess
		keygen.BlockNumber = tss.KeyGenZetaHeight
		keygen.GranteePubkeys = tss.TssParticipantList
		keeper.SetKeygen(ctx, keygen)
	}
	return nil
}
