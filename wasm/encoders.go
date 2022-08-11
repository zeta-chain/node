package wasm

import (
	"encoding/json"
	"fmt"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	zetaCoreModuleTypes "github.com/zeta-chain/zetacore/x/zetacore/types"
)

func Encoders(cdc codec.Codec) *wasmkeeper.MessageEncoders {
	return &wasmkeeper.MessageEncoders{
		Custom: EncodeZetacoreMessage(cdc),
	}
}

func EncodeZetacoreMessage(cdc codec.Codec) wasmkeeper.CustomEncoder {
	return func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var zetaMsg ZetaCoreMsg
		err := json.Unmarshal(msg, &zetaMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		switch {
		case zetaMsg.AddToWatchList != nil:
			return EncodeAddToWatchListMsg(sender, zetaMsg.AddToWatchList)
		}
		return nil, fmt.Errorf("Unknown ZetacoreMSG type")
	}
}

func EncodeAddToWatchListMsg(sender sdk.AccAddress, msg *AddToWatchList) ([]sdk.Msg, error) {
	zetaMsg := zetaCoreModuleTypes.NewMsgAddToWatchList(sender.String(), msg.Chain, msg.Nonce, msg.TxHash)
	return []sdk.Msg{zetaMsg}, nil
}
