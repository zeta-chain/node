package zetaclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/txscript"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

const (
	satoshiPerBitcoin = 1e8
)

func getSatoshis(btc float64) (int64, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type.  This may happen if f is NaN or +-Infinity.
	// BTC max amount is 21 mil and its at least 0 (Note: bitcoin allows creating 0-value outputs)
	switch {
	case math.IsNaN(btc):
		fallthrough
	case math.IsInf(btc, 1):
		fallthrough
	case math.IsInf(btc, -1):
		return 0, errors.New("invalid bitcoin amount")
	case btc > 21000000.0:
		return 0, errors.New("exceeded max bitcoin amount")
	case btc < 0.0:
		return 0, errors.New("cannot be less than zero")
	}
	return round(btc * satoshiPerBitcoin), nil
}

func round(f float64) int64 {
	if f < 0 {
		// #nosec G701 always in range
		return int64(f - 0.5)
	}
	// #nosec G701 always in range
	return int64(f + 0.5)
}

func payToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

type DynamicTicker struct {
	name     string
	interval uint64
	impl     *time.Ticker
}

func NewDynamicTicker(name string, interval uint64) *DynamicTicker {
	return &DynamicTicker{
		name:     name,
		interval: interval,
		impl:     time.NewTicker(time.Duration(interval) * time.Second),
	}
}

func (t *DynamicTicker) C() <-chan time.Time {
	return t.impl.C
}

func (t *DynamicTicker) UpdateInterval(newInterval uint64, logger zerolog.Logger) {
	if newInterval > 0 && t.interval != newInterval {
		t.impl.Stop()
		oldInterval := t.interval
		t.interval = newInterval
		t.impl = time.NewTicker(time.Duration(t.interval) * time.Second)
		logger.Info().Msgf("%s ticker interval changed from %d to %d", t.name, oldInterval, newInterval)
	}
}

func (t *DynamicTicker) Stop() {
	t.impl.Stop()
}

func (ob *EVMChainClient) GetInboundVoteMsgForDepositedEvent(event *erc20custody.ERC20CustodyDeposited) (types.MsgVoteOnObservedInboundTx, error) {
	ob.logger.ExternalChainWatcher.Info().Msgf("TxBlockNumber %d Transaction Hash: %s Message : %s", event.Raw.BlockNumber, event.Raw.TxHash, event.Message)
	if bytes.Compare(event.Message, []byte(DonationMessage)) == 0 {
		ob.logger.ExternalChainWatcher.Info().Msgf("thank you rich folk for your donation!: %s", event.Raw.TxHash.Hex())
		return types.MsgVoteOnObservedInboundTx{}, fmt.Errorf("thank you rich folk for your donation!: %s", event.Raw.TxHash.Hex())
	}
	// get the sender of the event's transaction
	tx, _, err := ob.evmClient.TransactionByHash(context.Background(), event.Raw.TxHash)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg(fmt.Sprintf("failed to get transaction by hash: %s", event.Raw.TxHash.Hex()))
		return types.MsgVoteOnObservedInboundTx{}, errors.Wrap(err, fmt.Sprintf("failed to get transaction by hash: %s", event.Raw.TxHash.Hex()))
	}
	signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
	sender, err := signer.Sender(tx)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg(fmt.Sprintf("can't recover the sender from the tx hash: %s", event.Raw.TxHash.Hex()))
		return types.MsgVoteOnObservedInboundTx{}, errors.Wrap(err, fmt.Sprintf("can't recover the sender from the tx hash: %s", event.Raw.TxHash.Hex()))

	}
	return *GetInBoundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		"",
		clienttypes.BytesToEthHex(event.Recipient),
		common.ZetaChain().ChainId,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		common.CoinType_ERC20,
		event.Asset.String(),
		ob.zetaClient.keys.GetOperatorAddress().String(),
	), nil
}

func (ob *EVMChainClient) GetInboundVoteMsgForZetaSentEvent(event *zetaconnector.ZetaConnectorNonEthZetaSent) (types.MsgVoteOnObservedInboundTx, error) {
	ob.logger.ExternalChainWatcher.Info().Msgf("TxBlockNumber %d Transaction Hash: %s Message : %s", event.Raw.BlockNumber, event.Raw.TxHash, event.Message)
	destChain := common.GetChainFromChainID(event.DestinationChainId.Int64())
	if destChain == nil {
		ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not supported  %d", event.DestinationChainId.Int64())
		return types.MsgVoteOnObservedInboundTx{}, fmt.Errorf("chain id not supported  %d", event.DestinationChainId.Int64())
	}
	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)
	if *destChain != common.ZetaChain() {
		cfgDest, found := ob.cfg.GetEVMConfig(destChain.ChainId)
		if !found {
			return types.MsgVoteOnObservedInboundTx{}, fmt.Errorf("chain id not present in EVMChainConfigs  %d", event.DestinationChainId.Int64())
		}
		if strings.EqualFold(destAddr, cfgDest.ZetaTokenContractAddress) {
			ob.logger.ExternalChainWatcher.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
			return types.MsgVoteOnObservedInboundTx{}, fmt.Errorf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
		}
	}
	return *GetInBoundVoteMessage(
		event.ZetaTxSenderAddress.Hex(),
		ob.chain.ChainId,
		event.SourceTxOriginAddress.Hex(),
		clienttypes.BytesToEthHex(event.DestinationAddress),
		destChain.ChainId,
		sdkmath.NewUintFromBigInt(event.ZetaValueAndGas),
		base64.StdEncoding.EncodeToString(event.Message),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		event.DestinationGasLimit.Uint64(),
		common.CoinType_Zeta,
		"",
		ob.zetaClient.keys.GetOperatorAddress().String(),
	), nil
}

func (ob *EVMChainClient) GetInboundVoteMsgForTokenSentToTSS(txhash ethcommon.Hash, value *big.Int, receipt *ethtypes.Receipt, from ethcommon.Address, data []byte) *types.MsgVoteOnObservedInboundTx {
	ob.logger.ExternalChainWatcher.Info().Msgf("TSS inTx detected: %s, blocknum %d", txhash.Hex(), receipt.BlockNumber)
	ob.logger.ExternalChainWatcher.Info().Msgf("TSS inTx value: %s", value.String())
	ob.logger.ExternalChainWatcher.Info().Msgf("TSS inTx from: %s", from.Hex())
	message := ""
	if len(data) != 0 {
		message = hex.EncodeToString(data)
	}
	return GetInBoundVoteMessage(
		from.Hex(),
		ob.chain.ChainId,
		from.Hex(),
		from.Hex(),
		common.ZetaChain().ChainId,
		sdkmath.NewUintFromBigInt(value),
		message,
		txhash.Hex(),
		receipt.BlockNumber.Uint64(),
		90_000,
		common.CoinType_Gas,
		"",
		ob.zetaClient.keys.GetOperatorAddress().String(),
	)
}
