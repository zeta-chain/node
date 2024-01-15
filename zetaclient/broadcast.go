package zetaclient

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/zeta-chain/zetacore/cmd/zetacored/config"

	"github.com/zeta-chain/zetacore/common/cosmos"
	"github.com/zeta-chain/zetacore/zetaclient/hsm"

	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	flag "github.com/spf13/pflag"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	// DefaultBaseGasPrice is the default base gas price
	DefaultBaseGasPrice = 1_000_000
)

// Broadcast Broadcasts tx to metachain. Returns txHash and error
func (b *ZetaCoreBridge) Broadcast(gaslimit uint64, authzWrappedMsg sdktypes.Msg, authzSigner AuthZSigner) (string, error) {
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
		baseGasPrice = DefaultBaseGasPrice // shoudn't happen, but just in case
	}

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
	//b.logger.Info().Uint64("account_number", b.accountNumber).Uint64("sequence_number", b.seqNumber).Msg("account info")

	flags := flag.NewFlagSet("zetacore", 0)

	ctx, err := b.GetContext()
	if err != nil {
		return "", err
	}
	factory := clienttx.NewFactoryCLI(ctx, flags)
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
		cosmos.NewInt(int64(gaslimit)).Mul(cosmos.NewInt(baseGasPrice))))
	builder.SetFeeAmount(fee)
	//fmt.Printf("signing from name: %s\n", ctx.GetFromName())
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
	//b.logger.Debug().Msgf("Received a TxHash of %v from the metachain, Code %d, log %s", commit.TxHash, commit.Code, commit.Logs)

	// increment seqNum
	//seq := b.seqNumber[authzSigner.KeyType]
	//atomic.AddUint64(&seq, 1)
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
