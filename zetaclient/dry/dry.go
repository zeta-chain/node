// Package dry provides dry-client wrappers for the standard clients of the connected chains.
//
// A dry-client wrapper overrides mutating functions from the underlying client.
// These overridden functions panic with MsgUnreacheable when called.
//
// Dry-client wrappers are redundant.
// They and serve as an additional safeguard layer that guarantees that dry-mode zetaclient nodes
// never mutate the state of the connected chains.
package dry

import (
	"context"

	suimodel "github.com/block-vision/sui-go-sdk/models"

	"github.com/zeta-chain/node/zetaclient/chains/sui"
	"github.com/zeta-chain/node/zetaclient/chains/ton"
)

const MsgUnreachable = "unreachable"

// ------------------------------------------------------------------------------------------------
// Bitcoin
// ------------------------------------------------------------------------------------------------

// TODO
// See: https://github.com/zeta-chain/node/issues/4232

// ------------------------------------------------------------------------------------------------
// EVM
// ------------------------------------------------------------------------------------------------

// TODO
// See: https://github.com/zeta-chain/node/issues/4231

// ------------------------------------------------------------------------------------------------
// Solana
// ------------------------------------------------------------------------------------------------

// TODO
// See: https://github.com/zeta-chain/node/issues/4230

// ------------------------------------------------------------------------------------------------
// Sui
// ------------------------------------------------------------------------------------------------

type SuiClient struct {
	sui.Client
}

func WrapSuiClient(client sui.Client) *SuiClient {
	return &SuiClient{Client: client}
}

func (*SuiClient) SuiExecuteTransactionBlock(context.Context,
	suimodel.SuiExecuteTransactionBlockRequest,
) (suimodel.SuiTransactionBlockResponse, error) {
	panic(MsgUnreachable)
}

// ------------------------------------------------------------------------------------------------
// TON
// ------------------------------------------------------------------------------------------------

type TONClient struct {
	ton.Client
}

func WrapTONClient(client ton.Client) *TONClient {
	return &TONClient{Client: client}
}

func (*TONClient) SendMessage(context.Context, []byte) (uint32, error) {
	panic(MsgUnreachable)
}
