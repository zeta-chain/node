// Package dry provides dry-client wrappers for the TSS signer, for the zetacore client, and for
// the standard clients of the connected chains.
//
// A dry-client wrapper overrides mutating functions from the underlying client.
// These overridden functions panic with MsgUnreacheable when called.
//
// Dry-client wrappers are redundant.
// They serve as an additional safeguard layer that guarantees that dry-mode zetaclient nodes never
// participate in signing, never mutate ZetaChain state, and never mutate the state of the
// connected chains.
package dry

import (
	"context"

	suimodel "github.com/block-vision/sui-go-sdk/models"
	btcchainhash "github.com/btcsuite/btcd/chaincfg/chainhash"
	btcwire "github.com/btcsuite/btcd/wire"
	eth "github.com/ethereum/go-ethereum/core/types"
	sol "github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/zeta-chain/go-tss/blame"

	"github.com/zeta-chain/node/pkg/chains"
	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	crosschain "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/evm"
	"github.com/zeta-chain/node/zetaclient/chains/solana"
	"github.com/zeta-chain/node/zetaclient/chains/sui"
	"github.com/zeta-chain/node/zetaclient/chains/ton"
	"github.com/zeta-chain/node/zetaclient/chains/tssrepo"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/tss"
)

// MsgUnreachable is the panic message returned by this module's functions when they get called.
const MsgUnreachable = "called an unreachable dry-mode function"

// ------------------------------------------------------------------------------------------------
// ZetaCore
// ------------------------------------------------------------------------------------------------

// ZetacoreClient is a dry-wrapper for zetacore clients.
type ZetacoreClient struct {
	// We only embed the reader client. The writer interface is deliberately not embedded so the
	// compiler can ensure that all mutating methods are explicitly overridden.
	zrepo.ZetacoreReaderClient
}

func WrapZetacoreClient(client zrepo.ZetacoreClient) *ZetacoreClient {
	return &ZetacoreClient{ZetacoreReaderClient: client}
}

func (*ZetacoreClient) PostVoteGasPrice(context.Context, chains.Chain, uint64, uint64, uint64,
) (string, error) {
	panic(MsgUnreachable)
}

func (*ZetacoreClient) PostVoteTSS(context.Context, string, int64, chains.ReceiveStatus,
) (string, error) {
	panic(MsgUnreachable)
}

func (*ZetacoreClient) PostVoteBlameData(context.Context, *blame.Blame, int64, string,
) (string, error) {
	panic(MsgUnreachable)
}

func (*ZetacoreClient) PostVoteOutbound(context.Context, uint64, uint64,
	*crosschain.MsgVoteOutbound,
) (string, string, error) {
	panic(MsgUnreachable)
}

func (*ZetacoreClient) PostVoteInbound(context.Context, uint64, uint64,
	*crosschain.MsgVoteInbound, chan<- zetaerrors.ErrTxMonitor,
) (string, string, error) {
	panic(MsgUnreachable)
}

func (*ZetacoreClient) PostOutboundTracker(context.Context, int64, uint64, string,
) (string, error) {
	panic(MsgUnreachable)
}

// ------------------------------------------------------------------------------------------------
// TSS
// ------------------------------------------------------------------------------------------------

// TSSClient is a dry-wrapper for TSS clients.
type TSSClient struct {
	// client is deliberately not embedded so the compiler can ensure that all mutating
	// methods are explicitly overridden.
	client tssrepo.TSSClient
}

func WrapTSSClient(client tssrepo.TSSClient) *TSSClient {
	return &TSSClient{client}
}

func (signer *TSSClient) PubKey() tss.PubKey {
	return signer.client.PubKey()
}

func (*TSSClient) Sign(context.Context, []byte, uint64, uint64, int64) ([65]byte, error) {
	panic(MsgUnreachable)
}

func (*TSSClient) SignBatch(context.Context, [][]byte, uint64, uint64, int64,
) ([][65]byte, error) {
	panic(MsgUnreachable)
}

func (*TSSClient) IsSignatureCached(int64, [][]byte) bool {
	panic(MsgUnreachable)
}

// ------------------------------------------------------------------------------------------------
// Bitcoin
// ------------------------------------------------------------------------------------------------

// BitcoinClient is a dry-wrapper for Bitcoin clients.
type BitcoinClient struct {
	bitcoin.Client
}

func WrapBitcoinClient(client bitcoin.Client) *BitcoinClient {
	return &BitcoinClient{Client: client}
}

func (*BitcoinClient) SendRawTransaction(context.Context,
	*btcwire.MsgTx, bool,
) (*btcchainhash.Hash, error) {
	panic(MsgUnreachable)
}

// ------------------------------------------------------------------------------------------------
// EVM
// ------------------------------------------------------------------------------------------------

// EVMClient is a dry-wrapper for EVM clients.
type EVMClient struct {
	evm.Client
}

func WrapEVMClient(client evm.Client) *EVMClient {
	return &EVMClient{Client: client}
}

func (*EVMClient) SendTransaction(context.Context, *eth.Transaction) error {
	panic(MsgUnreachable)
}

// ------------------------------------------------------------------------------------------------
// Solana
// ------------------------------------------------------------------------------------------------

// SolanaClient is a dry-wrapper for Solana clients.
type SolanaClient struct {
	solana.Client
}

func WrapSolanaClient(client solana.Client) *SolanaClient {
	return &SolanaClient{Client: client}
}

func (*SolanaClient) SendTransactionWithOpts(context.Context, *sol.Transaction,
	solrpc.TransactionOpts,
) (sol.Signature, error) {
	panic(MsgUnreachable)
}

// ------------------------------------------------------------------------------------------------
// Sui
// ------------------------------------------------------------------------------------------------

// SuiClient is a dry-wrapper for Sui clients.
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

// TONClient is a dry-wrapper for TON clients.
type TONClient struct {
	ton.Client
}

func WrapTONClient(client ton.Client) *TONClient {
	return &TONClient{Client: client}
}

func (*TONClient) SendMessage(context.Context, []byte) (uint32, error) {
	panic(MsgUnreachable)
}
