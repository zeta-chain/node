package signer

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	"github.com/davecgh/go-spew/spew"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

	gatewayID, err := solana.PublicKeyFromBase58(chainParams.GatewayAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid gateway address %s", chainParams.GatewayAddress)
	}

	// compute gateway PDA
	seed := []byte(contract.PDASeed)
	pda, _, err := solana.FindProgramAddress([][]byte{seed}, gatewayID)
	if err != nil {
		return nil, err
	}

	// create solana observer
	signer := &Signer{
		Signer:    baseSigner,
		client:    solClient,
		gatewayID: gatewayID,
		pda:       pda,
	}

	return signer, nil
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
		return nil, errors.Wrap(err, "Key-ssign failed")
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
	logger := signer.Logger().Std.With().
		Str("OutboundID", outboundID).
		Str("SendHash", cctx.Index).
		Logger()

	params := cctx.GetCurrentOutboundParam()
	logger.Info().
		Msgf("Solana TryProcessOutbound: %s, value %d to %s", cctx.Index, params.Amount.BigInt(), params.Receiver)

	// support gas token only for Solana outbound
	coinType := cctx.InboundParams.CoinType
	if coinType == coin.CoinType_Zeta || coinType == coin.CoinType_ERC20 {
		logger.Error().Msgf("TryProcessOutbound: can only send SOL to the Solana network")
		return
	}

	chain := signer.Chain()
	outboundTssNonce := params.TssNonce
	// get size limit and gas price
	// fee := 5000 // FIXME: this is the fixed fee (for signatures), explore priority fee for compute units

	// check receiver address
	to, err := chains.DecodeSolanaWalletAddress(params.Receiver)
	if err != nil {
		logger.Error().Msgf("TryProcessOutbound: cannot decode receiver address %s", params.Receiver)
		return
	}
	amount := params.Amount.Uint64()

	{ // TODO: refactor this piece out to a separate (withdraw) function
		// FIXME: config this; right now it's the same privkey used by local e2e test_solana_*.go
		privkey := solana.MustPrivateKeyFromBase58(
			"4yqSQxDeTBvn86BuxcN5jmZb2gaobFXrBqu8kiE9rZxNkVMe3LfXmFigRsU4sRp7vk4vVP1ZCFiejDKiXBNWvs2C",
		)

		seed := []byte("meta")
		pdaComputed, bump, err := solana.FindProgramAddress([][]byte{seed}, signer.gatewayID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("computed pda: %s, bump %d\n", pdaComputed, bump)
		pdaInfo, err := signer.client.GetAccountInfo(context.TODO(), pdaComputed)
		if err != nil {
			panic(err)
		}
		fmt.Printf("pdainfo: %v\n", pdaInfo.Bytes())

		var pda contract.PdaInfo
		err = borsh.Deserialize(&pda, pdaInfo.Bytes())
		if err != nil {
			panic(err)
		}
		fmt.Printf("pda parsed: %+v\n", pda)

		recent, err := signer.client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}
		fmt.Println("recent blockhash:", recent.Value.Blockhash)
		var inst solana.GenericInstruction

		var message []byte
		bytes := make([]byte, 8)
		chainID := uint64(chain.ChainId)
		nonce := outboundTssNonce
		binary.BigEndian.PutUint64(bytes, chainID)
		message = append(message, bytes...)
		binary.BigEndian.PutUint64(bytes, nonce)
		message = append(message, bytes...)
		binary.BigEndian.PutUint64(bytes, amount)
		message = append(message, bytes...)
		message = append(message, to.Bytes()...)
		messageHash := crypto.Keccak256Hash(message)
		fmt.Printf(
			"solana msghash: chainid %d, nonce %d, amount %d, to %s, hash %s",
			chainID,
			nonce,
			amount,
			to.String(),
			messageHash.String(),
		)
		// this sig will be 65 bytes; R || S || V, where V is 0 or 1
		signature, err := signer.TSS().Sign(ctx, messageHash.Bytes(), height, nonce, chain.ChainId, "")
		if err != nil {
			signer.Logger().Std.Error().Err(err).Msg("cannot sign message")
			panic(err)
		}
		signer.Logger().Std.Info().
			Msgf("Key-sign success: %d => %s, nonce %d", cctx.InboundParams.SenderChainId, chain.ChainName, outboundTssNonce)

		signer.Logger().Std.Info().Msgf("recovery id %d", signature[64])
		var sig [64]byte
		copy(sig[:], signature[:64])

		inst.DataBytes, err = borsh.Serialize(contract.WithdrawInstructionParams{
			Discriminator: [8]byte{183, 18, 70, 156, 148, 109, 161, 34},
			Amount:        amount,
			Signature:     sig,
			RecoveryID:    signature[64],
			MessageHash:   messageHash,
			Nonce:         nonce,
		})
		if err != nil {
			signer.Logger().Std.Error().Err(err).Msg("cannot serialize instruction")
			panic(err)
		}
		var accountSlice []*solana.AccountMeta
		accountSlice = append(accountSlice, solana.Meta(privkey.PublicKey()).WRITE().SIGNER())
		accountSlice = append(accountSlice, solana.Meta(pdaComputed).WRITE())
		accountSlice = append(accountSlice, solana.Meta(to).WRITE())
		accountSlice = append(accountSlice, solana.Meta(signer.gatewayID))
		inst.ProgID = signer.gatewayID
		inst.AccountValues = accountSlice
		tx, err := solana.NewTransaction(
			[]solana.Instruction{&inst},
			recent.Value.Blockhash,
			solana.TransactionPayer(privkey.PublicKey()),
		)
		if err != nil {
			panic(err)
		}
		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if privkey.PublicKey().Equals(key) {
					return &privkey
				}
				return nil
			},
		)
		if err != nil {
			panic(fmt.Errorf("unable to sign transaction: %w", err))
		}
		spew.Dump(tx)
		// FIXME: simulate before broadcast!
		txsig, err := signer.client.SendTransactionWithOpts(
			ctx,
			tx,
			rpc.TransactionOpts{
				// default PreflightCommitment is "finalized" which is too conservative
				// and results in repeated broadcast attempts that only 1 will succeed
				// Setting a "processed" level will simulate tx against more recent state
				// thus fails faster after a tx is already broadcasted and processed in a block.
				// This reduces the number of "failed" txs due to repeated broadcast attempts.
				PreflightCommitment: rpc.CommitmentConfirmed,
			},
		)
		if err != nil {
			signer.Logger().Std.Warn().Err(err).Msg("broadcast error")
		} else {
			signer.Logger().Std.Info().Msgf("broadcast success! tx sig %s; waiting for confirmation...", txsig)
			// launch a go routine with timeout to check for tx confirmation;
			// repeatedly query until timeout or the transaction is included in a block, either with success or failure
			go func() {
				txsig := txsig // capture the value
				nonce := nonce
				t1 := time.Now()
				for {
					if time.Since(t1) > 2*time.Minute {
						return
					}
					out, err := signer.client.GetConfirmedTransactionWithOpts(context.TODO(), txsig, &rpc.GetTransactionOpts{
						// I'd like to use "CommitmentProcessed" but it seems not supported in RPC: see https://solana.com/docs/rpc/http/gettransaction
						Commitment: rpc.CommitmentConfirmed,
					})
					if err == nil {
						if out.Meta.Err == nil { // successfully included in a block; report and exit goroutine
							txhash, err := zetacoreClient.AddOutboundTracker(ctx, signer.Chain().ChainId, nonce, txsig.String(), nil, "", -1)
							if err != nil {
								signer.Logger().Std.Error().Err(err).Msgf("unable to add to tracker: tx %s", txsig)
							} else {
								signer.Logger().Std.Info().Msgf("added txsig %s to outbound tracker; zeta txhash %s", txsig, txhash)
							}
							return
						}
						// it's included by failed (likely competing txs succeeded). exit goroutine.
						// FIXME: we should report this failed tx ONLY IF it failed not due to nonce mismatch error
						// FIXME: add a check for nonce mismatch error
						signer.Logger().Std.Warn().Msgf("tx %s failed: %v", txsig, out.Meta.Err)
						return
					}
					time.Sleep(10 * time.Second)
				}
			}()
		}
	}
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
