package zetabridge

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/zeta-chain/zetacore/app/ante"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/zetaclient/authz"
	"github.com/zeta-chain/zetacore/zetaclient/hsm"
)

// BroadcastInterface defines the signature of the broadcast function used by zetabridge transactions
type BroadcastInterface = func(bridge *ZetaCoreBridge, gaslimit uint64, authzWrappedMsg sdktypes.Msg, authzSigner authz.Signer) (string, error)

const (
	// DefaultBaseGasPrice is the default base gas price
	DefaultBaseGasPrice = 1_000_000
)

var (
	// paying 50% more than the current base gas price to buffer for potential block-by-block
	// gas price increase due to EIP1559 feemarket on ZetaChain
	bufferMultiplier = sdktypes.MustNewDecFromStr("1.5")

	// Variable function used by transactions to broadcast a message to ZetaCore. This will create enough flexibility
	// in the implementation to allow for more comprehensive unit testing.
	zetaBridgeBroadcast BroadcastInterface = BroadcastToZetaCore
)

// BroadcastToZetaCore is the default broadcast function used to send transactions to Zeta Core
func BroadcastToZetaCore(bridge *ZetaCoreBridge, gasLimit uint64, authzWrappedMsg sdktypes.Msg, authzSigner authz.Signer) (string, error) {
	return bridge.Broadcast(gasLimit, authzWrappedMsg, authzSigner)
}

// Broadcast Broadcasts tx to metachain. Returns txHash and error
func (b *ZetaCoreBridge) Broadcast(gaslimit uint64, authzWrappedMsg sdktypes.Msg, authzSigner authz.Signer) (string, error) {
	b.broadcastLock.Lock()
	defer b.broadcastLock.Unlock()
	var err error

	blockHeight, err := b.GetZetaBlockHeight()
	if err != nil {
		return "", err
	}
	baseGasPrice, err := b.GetBaseGasPrice()
	if err != nil {
		return "", err
	}
	if baseGasPrice == 0 {
		baseGasPrice = DefaultBaseGasPrice // shouldn't happen, but just in case
	}
	reductionRate := sdktypes.MustNewDecFromStr(ante.GasPriceReductionRate)
	// multiply gas price by the system tx reduction rate
	adjustedBaseGasPrice := sdktypes.NewDec(baseGasPrice).Mul(reductionRate).Mul(bufferMultiplier)

	if blockHeight > b.blockHeight {
		b.blockHeight = blockHeight
		accountNumber, seqNumber, err := b.GetAccountNumberAndSequenceNumber(authzSigner.KeyType)
		if err != nil {
			return "", err
		}
		b.accountNumber[authzSigner.KeyType] = accountNumber
		if b.seqNumber[authzSigner.KeyType] < seqNumber {
			b.seqNumber[authzSigner.KeyType] = seqNumber
		}
	}

	flags := flag.NewFlagSet("zetabridge", 0)

	ctx, err := b.GetContext()
	if err != nil {
		return "", err
	}
	factory, err := clienttx.NewFactoryCLI(ctx, flags)
	if err != nil {
		return "", err
	}

	factory = factory.WithAccountNumber(b.accountNumber[authzSigner.KeyType])
	factory = factory.WithSequence(b.seqNumber[authzSigner.KeyType])
	factory = factory.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	builder, err := factory.BuildUnsignedTx(authzWrappedMsg)
	if err != nil {
		return "", err
	}

	builder.SetGasLimit(gaslimit)

	// #nosec G701 always in range
	fee := sdktypes.NewCoins(sdktypes.NewCoin(config.BaseDenom,
		sdktypes.NewInt(int64(gaslimit)).Mul(adjustedBaseGasPrice.Ceil().RoundInt())))
	builder.SetFeeAmount(fee)
	err = b.SignTx(factory, ctx.GetFromName(), builder, true, ctx.TxConfig)
	if err != nil {
		return "", err
	}
	txBytes, err := ctx.TxConfig.TxEncoder()(builder.GetTx())
	if err != nil {
		return "", err
	}

	// broadcast to a Tendermint node
	commit, err := ctx.BroadcastTxSync(txBytes)
	if err != nil {
		b.logger.Error().Err(err).Msgf("fail to broadcast tx %s", err.Error())
		return "", err
	}

	// Code will be the tendermint ABICode , it start at 1 , so if it is an error , code will not be zero
	if commit.Code > 0 {
		if commit.Code == 32 {
			errMsg := commit.RawLog
			re := regexp.MustCompile(`account sequence mismatch, expected ([0-9]*), got ([0-9]*)`)
			matches := re.FindStringSubmatch(errMsg)
			if len(matches) != 3 {
				return "", err
			}
			expectedSeq, err := strconv.ParseUint(matches[1], 10, 64)
			if err != nil {
				b.logger.Warn().Msgf("cannot parse expected seq %s", matches[1])
				return "", err
			}
			gotSeq, err := strconv.Atoi(matches[2])
			if err != nil {
				b.logger.Warn().Msgf("cannot parse got seq %s", matches[2])
				return "", err
			}
			b.seqNumber[authzSigner.KeyType] = expectedSeq
			b.logger.Warn().Msgf("Reset seq number to %d (from err msg) from %d", b.seqNumber[authzSigner.KeyType], gotSeq)
		}
		return commit.TxHash, fmt.Errorf("fail to broadcast to zetachain,code:%d, log:%s", commit.Code, commit.RawLog)
	}

	// increment seqNum
	b.seqNumber[authzSigner.KeyType] = b.seqNumber[authzSigner.KeyType] + 1

	return commit.TxHash, nil
}

// GetContext return a valid context with all relevant values set
func (b *ZetaCoreBridge) GetContext() (client.Context, error) {
	ctx := client.Context{}
	addr, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		b.logger.Error().Err(err).Msg("fail to get address from key")
		return ctx, err
	}

	// if password is needed, set it as input
	password := b.keys.GetHotkeyPassword()
	if password != "" {
		ctx = ctx.WithInput(strings.NewReader(fmt.Sprintf("%[1]s\n%[1]s\n", password)))
	}

	ctx = ctx.WithKeyring(b.keys.GetKeybase())
	ctx = ctx.WithChainID(b.zetaChainID)
	ctx = ctx.WithHomeDir(b.cfg.ChainHomeFolder)
	ctx = ctx.WithFromName(b.cfg.SignerName)
	ctx = ctx.WithFromAddress(addr)
	ctx = ctx.WithBroadcastMode("sync")

	ctx = ctx.WithCodec(b.encodingCfg.Codec)
	ctx = ctx.WithInterfaceRegistry(b.encodingCfg.InterfaceRegistry)
	ctx = ctx.WithTxConfig(b.encodingCfg.TxConfig)
	ctx = ctx.WithLegacyAmino(b.encodingCfg.Amino)
	ctx = ctx.WithAccountRetriever(authtypes.AccountRetriever{})

	if b.enableMockSDKClient {
		ctx = ctx.WithClient(b.mockSDKClient)
	} else {
		remote := b.cfg.ChainRPC
		if !strings.HasPrefix(b.cfg.ChainHost, "http") {
			remote = fmt.Sprintf("tcp://%s", remote)
		}

		ctx = ctx.WithNodeURI(remote)
		wsClient, err := rpchttp.New(remote, "/websocket")
		if err != nil {
			return ctx, err
		}

		ctx = ctx.WithClient(wsClient)
	}

	return ctx, nil
}

func (b *ZetaCoreBridge) SignTx(
	txf clienttx.Factory,
	name string,
	txBuilder client.TxBuilder,
	overwriteSig bool,
	txConfig client.TxConfig,
) error {
	if b.cfg.HsmMode {
		return hsm.SignWithHSM(txf, name, txBuilder, overwriteSig, txConfig)
	}
	return clienttx.Sign(txf, name, txBuilder, overwriteSig)
}

// QueryTxResult query the result of a tx
func (b *ZetaCoreBridge) QueryTxResult(hash string) (*sdktypes.TxResponse, error) {
	ctx, err := b.GetContext()
	if err != nil {
		return nil, err
	}
	return authtx.QueryTx(ctx, hash)
}

// HandleBroadcastError returns whether to retry in a few seconds, and whether to report via AddTxHashToOutboundTracker
// returns (bool retry, bool report)
func HandleBroadcastError(err error, nonce, toChain, outboundHash string) (bool, bool) {
	if strings.Contains(err.Error(), "nonce too low") {
		log.Warn().Err(err).Msgf("nonce too low! this might be a unnecessary key-sign. increase re-try interval and awaits outTx confirmation")
		return false, false
	}
	if strings.Contains(err.Error(), "replacement transaction underpriced") {
		log.Warn().Err(err).Msgf("Broadcast replacement: nonce %s chain %s outboundHash %s", nonce, toChain, outboundHash)
		return false, false
	} else if strings.Contains(err.Error(), "already known") { // this is error code from QuickNode
		log.Warn().Err(err).Msgf("Broadcast duplicates: nonce %s chain %s outboundHash %s", nonce, toChain, outboundHash)
		return false, true // report to tracker, because there's possibilities a successful broadcast gets this error code
	}

	log.Error().Err(err).Msgf("Broadcast error: nonce %s chain %s outboundHash %s; retrying...", nonce, toChain, outboundHash)
	return true, false
}
