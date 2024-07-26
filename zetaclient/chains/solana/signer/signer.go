package signer

import (
	"context"

	"cosmossdk.io/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	contract "github.com/zeta-chain/zetacore/pkg/contract/solana"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"
)

var _ interfaces.ChainSigner = (*Signer)(nil)

// Signer deals with signing BTC transactions and implements the ChainSigner interface
type Signer struct {
	*base.Signer

	// client is the Solana RPC client that interacts with the Solana chain
	client interfaces.SolanaRPCClient

	// gatewayID is the program ID of gateway program on Solana chain
	gatewayID solana.PublicKey

	// pda is the program derived address of the gateway program
	pda solana.PublicKey
}

// NewSigner creates a new Bitcoin signer
func NewSigner(
	chain chains.Chain,
	chainParams observertypes.ChainParams,
	solClient interfaces.SolanaRPCClient,
	tss interfaces.TSSSigner,
	ts *metrics.TelemetryServer,
	logger base.Logger,
) (*Signer, error) {
	// create base signer
	baseSigner := base.NewSigner(chain, tss, ts, logger)

	// parse gateway ID and PDA
	gatewayID, pda, err := contract.ParseGatewayIDAndPda(chainParams.GatewayAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse gateway address %s", chainParams.GatewayAddress)
	}

	// create solana observer
	return &Signer{
		Signer:    baseSigner,
		client:    solClient,
		gatewayID: gatewayID,
		pda:       pda,
	}, nil
}

// SignWithdrawTx signs a message for Solana gateway 'withdraw' transaction
func (signer *Signer) SignMsgWithdraw(
	ctx context.Context,
	params *types.OutboundParams,
	height uint64,
) (*contract.MsgWithdraw, error) {
	// #nosec G115 always positive
	chain := signer.Chain()
	chainID := uint64(signer.Chain().ChainId)
	nonce := params.TssNonce
	amount := params.Amount.Uint64()

	// check receiver address
	to, err := chains.DecodeSolanaWalletAddress(params.Receiver)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot decode receiver address %s", params.Receiver)
	}

	// prepare withdraw msg and compute hash
	msg := contract.NewMsgWithdraw(chainID, nonce, amount, to)
	msgHash := msg.Hash()

	// sign the message with TSS to get an ECDSA signature.
	// the produced signature is in the [R || S || V] format where V is 0 or 1.
	signature, err := signer.TSS().Sign(ctx, msgHash[:], height, nonce, chain.ChainId, "")
	if err != nil {
		return nil, errors.Wrap(err, "Key-sign failed")
	}
	msg.WithSignature(signature)

	signer.Logger().Std.Info().Msgf("Key-sign succeed for chain %d nonce %d", chainID, nonce)
	return msg, nil
}

// SignWithdrawTx signs the Solana gateway 'withdraw' transaction specified by 'msg'
func (signer *Signer) SignWithdrawTx(ctx context.Context, msg contract.MsgWithdraw) (*solana.Transaction, error) {
	// FIXME: config this; right now it's the same privkey used by local e2e test_solana_*.go
	privkey := solana.MustPrivateKeyFromBase58(
		"4yqSQxDeTBvn86BuxcN5jmZb2gaobFXrBqu8kiE9rZxNkVMe3LfXmFigRsU4sRp7vk4vVP1ZCFiejDKiXBNWvs2C",
	)

	// create withdraw instruction with program call data
	var err error
	var inst solana.GenericInstruction
	inst.DataBytes, err = borsh.Serialize(contract.WithdrawInstructionParams{
		Discriminator: contract.DiscriminatorWithdraw(),
		Amount:        msg.Amount,
		Signature:     msg.SigRS(),
		RecoveryID:    msg.SigV(),
		MessageHash:   msg.Hash(),
		Nonce:         msg.Nonce,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot serialize withdraw instruction")
	}

	// attach required accounts to the instruction
	var accountSlice []*solana.AccountMeta
	accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
	accountSlice = append(accountSlice, solana.Meta(signer.pda).WRITE())
	accountSlice = append(accountSlice, solana.Meta(msg.To).WRITE())
	accountSlice = append(accountSlice, solana.Meta(signer.gatewayID))
	inst.ProgID = signer.gatewayID
	inst.AccountValues = accountSlice

	// get a recent blockhash
	recent, err := signer.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, errors.Wrap(err, "GetRecentBlockhash error")
	}

	// create a transaction that wraps the instruction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{&inst},
		recent.Value.Blockhash,
		solana.TransactionPayer(privkey.PublicKey()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "NewTransaction error")
	}

	// fee payer signs the transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(privkey.PublicKey()) {
			return &privkey
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "signer unable to sign transaction")
	}

	return tx, nil
}

// TryProcessOutbound - signer interface implementation
// This function will attempt to build and sign a Solana transaction using the TSS signer.
// It will then broadcast the signed transaction to the Solana chain.
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *types.CrossChainTx,
	outboundProc *outboundprocessor.Processor,
	outboundID string,
	_ interfaces.ChainObserver,
	zetacoreClient interfaces.ZetacoreClient,
	height uint64,
) {
	// end outbound process on panic
	defer func() {
		outboundProc.EndTryProcess(outboundID)
		if err := recover(); err != nil {
			signer.Logger().Std.Error().Msgf("TryProcessOutbound: %s, caught panic error: %v", cctx.Index, err)
		}
	}()

	// prepare logger
	params := cctx.GetCurrentOutboundParam()
	logger := signer.Logger().Std.With().
		Str("OutboundID", outboundID).
		Str("SendHash", cctx.Index).
		Logger()
	logger.Info().
		Msgf("Solana TryProcessOutbound: %s, value %d to %s", cctx.Index, params.Amount.BigInt(), params.Receiver)

	// support gas token only for Solana outbound
	coinType := cctx.InboundParams.CoinType
	if coinType == coin.CoinType_Zeta || coinType == coin.CoinType_ERC20 {
		logger.Error().Msgf("TryProcessOutbound: can only send SOL to the Solana network")
		return
	}

	// sign gateway withdraw message by TSS
	chainID := signer.Chain().ChainId
	nonce := params.TssNonce
	msgWithdraw, err := signer.SignMsgWithdraw(ctx, params, height)
	if err != nil {
		logger.Error().Err(err).Msgf("TryProcessOutbound: SignMsgWithdraw error for chain %d nonce %d", chainID, nonce)
		return
	}

	// sign the withdraw transaction by fee payer
	tx, err := signer.SignWithdrawTx(ctx, *msgWithdraw)
	if err != nil {
		logger.Error().Err(err).Msgf("TryProcessOutbound: SignWithdrawTx error for chain %d nonce %d", chainID, nonce)
		return
	}

	// fee := 5000 // FIXME: this is the fixed fee (for signatures), explore priority fee for compute units
	// broadcast the signed tx to the Solana network with preflight check
	txSig, err := signer.client.SendTransactionWithOpts(
		ctx,
		tx,
		// Commitment "finalized" is too conservative for preflight check and
		// it results in repeated broadcast attempts that only 1 will succeed.
		// Commitment "processed" will simulate tx against more recent state
		// thus fails faster once a tx is already broadcasted and processed by the cluster.
		// This reduces the number of "failed" txs due to repeated broadcast attempts.
		rpc.TransactionOpts{PreflightCommitment: rpc.CommitmentConfirmed},
	)
	if err != nil {
		signer.Logger().
			Std.Warn().
			Err(err).
			Msgf("TryProcessOutbound: broadcast error for chain %d nonce %d", chainID, nonce)
		return
	}

	// report the outbound to the outbound tracker
	signer.reportToOutboundTracker(ctx, zetacoreClient, chainID, nonce, txSig, logger)
}

// SetGatewayAddress sets the gateway address
func (signer *Signer) SetGatewayAddress(address string) {
	// parse gateway ID and PDA
	gatewayID, pda, err := contract.ParseGatewayIDAndPda(address)
	if err != nil {
		signer.Logger().Std.Error().Err(err).Msgf("cannot parse gateway address %s", address)
	}

	// update gateway ID and PDA
	signer.Lock()
	defer signer.Unlock()

	signer.gatewayID = gatewayID
	signer.pda = pda
}

// GetGatewayAddress returns the gateway address
func (signer *Signer) GetGatewayAddress() string {
	signer.Lock()
	defer signer.Unlock()
	return signer.gatewayID.String()
}

// TODO: get rid of below four functions for Solana and Bitcoin
// https://github.com/zeta-chain/node/issues/2532
func (signer *Signer) SetZetaConnectorAddress(_ ethcommon.Address) {
}

func (signer *Signer) SetERC20CustodyAddress(_ ethcommon.Address) {
}

func (signer *Signer) GetZetaConnectorAddress() ethcommon.Address {
	return ethcommon.Address{}
}

func (signer *Signer) GetERC20CustodyAddress() ethcommon.Address {
	return ethcommon.Address{}
}
