// Package dry provides a dry client for TSS, and dry-wrappers for the zetacore client and for the
// standard clients of the connected chains.
//
// A dry client implements all of the functions of the standard client, while a dry-client wrapper
// simply overrides mutating functions from the underlying client.
// The non-mutating functions behave as its counterparts from a standard client.
// The mutating functions panic with MsgUnreacheable when called.
//
// Dry-clients are redundant.
// They serve as an additional safeguard layer that guarantees that dry-mode zetaclient nodes never
// participate in signing, never mutate ZetaChain state, and never mutate the state of the
// connected chains.
package dry

import (
	"context"
	"fmt"

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

// TSSClient is a dry TSS client.
type TSSClient struct {
	pubKey tss.PubKey
}

func NewTSSClient(tssAddress string) (*TSSClient, error) {
	pubKey, err := tss.NewPubKeyFromBech32(tssAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid TSS pub key: %w", err)
	}
	return &TSSClient{pubKey}, nil
}

func (client *TSSClient) PubKey() tss.PubKey {
	return client.pubKey
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
	bitcoin.BitcoinClient
}

func WrapBitcoinClient(client bitcoin.BitcoinClient) *BitcoinClient {
	return &BitcoinClient{BitcoinClient: client}
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
	evm.EVMClient
}

func WrapEVMClient(client evm.EVMClient) *EVMClient {
	return &EVMClient{EVMClient: client}
}

func (*EVMClient) SendTransaction(context.Context, *eth.Transaction) error {
	panic(MsgUnreachable)
}

// ------------------------------------------------------------------------------------------------
// Solana
// ------------------------------------------------------------------------------------------------

// SolanaClient is a dry-wrapper for Solana clients.
type SolanaClient struct {
	solana.SolanaClient
}

func WrapSolanaClient(client solana.SolanaClient) *SolanaClient {
	return &SolanaClient{SolanaClient: client}
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
	sui.SuiClient
}

func WrapSuiClient(client sui.SuiClient) *SuiClient {
	return &SuiClient{SuiClient: client}
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
	ton.TONClient
}

func WrapTONClient(client ton.TONClient) *TONClient {
	return &TONClient{TONClient: client}
}

func (*TONClient) SendMessage(context.Context, []byte) (uint32, error) {
	panic(MsgUnreachable)
}
