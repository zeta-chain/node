package zetaclient

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	flag "github.com/spf13/pflag"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/zeta-chain/zetacore/app"
)

// Broadcast Broadcasts tx to metachain. Returns txHash and error
func (b *ZetaCoreBridge) Broadcast(gaslimit uint64, authzWrappedMsg sdktypes.Msg, authzSigner AuthZSigner) (string, error) {
	gaslimit = gaslimit * 3
	b.broadcastLock.Lock()
	defer b.broadcastLock.Unlock()
	var err error

	blockHeight, err := b.GetZetaBlockHeight()
	if err != nil {
		return "", err
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

	ctx := b.GetContext()
	factory := clienttx.NewFactoryCLI(ctx, flags)
	factory = factory.WithAccountNumber(b.accountNumber[authzSigner.KeyType])
	factory = factory.WithSequence(b.seqNumber[authzSigner.KeyType])
	factory = factory.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	builder, err := factory.BuildUnsignedTx(authzWrappedMsg)
	if err != nil {
		return "", err
	}
	builder.SetGasLimit(gaslimit)
	fee := sdktypes.NewCoins(sdktypes.NewCoin("azeta", sdktypes.NewInt(40000)))
	builder.SetFeeAmount(fee)
	//fmt.Printf("signing from name: %s\n", ctx.GetFromName())
	err = clienttx.Sign(factory, ctx.GetFromName(), builder, true)
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
			expectedSeq, err := strconv.Atoi(matches[1])
			if err != nil {
				b.logger.Warn().Msgf("cannot parse expected seq %s", matches[1])
				return "", err
			}
			gotSeq, err := strconv.Atoi(matches[2])
			if err != nil {
				b.logger.Warn().Msgf("cannot parse got seq %s", matches[2])
				return "", err
			}
			b.seqNumber[authzSigner.KeyType] = uint64(expectedSeq)
			b.logger.Warn().Msgf("Reset seq number to %d (from err msg) from %d", b.seqNumber[authzSigner.KeyType], gotSeq)
		}
		b.logger.Info().Msgf("retrying message in 3s: %s", commit.RawLog)
		time.Sleep(3 * time.Second)
		return commit.TxHash, fmt.Errorf("fail to broadcast to zetachain,code:%d, log:%s", commit.Code, commit.RawLog)
	}
	//b.logger.Debug().Msgf("Received a TxHash of %v from the metachain, Code %d, log %s", commit.TxHash, commit.Code, commit.Logs)

	// increment seqNum
	//seq := b.seqNumber[authzSigner.KeyType]
	//atomic.AddUint64(&seq, 1)
	b.seqNumber[authzSigner.KeyType] = b.seqNumber[authzSigner.KeyType] + 1
	//b.logger.Debug().Msgf("b.sequence number increased to %d", b.seqNumber)

	return commit.TxHash, nil
}

// GetContext return a valid context with all relevant values set
func (b *ZetaCoreBridge) GetContext() client.Context {
	ctx := client.Context{}
	addr, _ := b.keys.GetSignerInfo().GetAddress()
	// TODO : Handle error
	ctx = ctx.WithKeyring(b.keys.GetKeybase())
	ctx = ctx.WithChainID(b.zetaChainID)
	ctx = ctx.WithHomeDir(b.cfg.ChainHomeFolder)
	ctx = ctx.WithFromName(b.cfg.SignerName)
	ctx = ctx.WithFromAddress(addr)
	ctx = ctx.WithBroadcastMode("sync")

	encodingConfig := app.MakeEncodingConfig()
	ctx = ctx.WithCodec(encodingConfig.Codec)
	ctx = ctx.WithInterfaceRegistry(encodingConfig.InterfaceRegistry)
	ctx = ctx.WithTxConfig(encodingConfig.TxConfig)
	ctx = ctx.WithLegacyAmino(encodingConfig.Amino)
	ctx = ctx.WithAccountRetriever(authtypes.AccountRetriever{})

	remote := b.cfg.ChainRPC
	if !strings.HasPrefix(b.cfg.ChainHost, "http") {
		remote = fmt.Sprintf("tcp://%s", remote)
	}
	//fmt.Println("ctx.remote ", remote)

	ctx = ctx.WithNodeURI(remote)
	client, err := rpchttp.New(remote, "/websocket")
	if err != nil {
		panic(err)
	}
	ctx = ctx.WithClient(client)
	return ctx
}
