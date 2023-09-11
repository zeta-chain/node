package zetaclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math"
	"math/big"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/txscript"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/common"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

const (
	satoshiPerBitcoin = 1e8
)

func getSatoshis(btc float64) (int64, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type.  This may happen if f is NaN or +-Infinity.
	// BTC max amount is 21 mil and its at least 10^(-8) or one satoshi.
	switch {
	case math.IsNaN(btc):
		fallthrough
	case math.IsInf(btc, 1):
		fallthrough
	case math.IsInf(btc, -1):
		return 0, errors.New("invalid bitcoin amount")
	case btc > 21000000.0:
		return 0, errors.New("exceeded max bitcoin amount")
	case btc < 0.00000001:
		return 0, errors.New("cannot be less than 1 satoshi")
	}
	return round(btc * satoshiPerBitcoin), nil
}

func round(f float64) int64 {
	if f < 0 {
		return int64(f - 0.5)
	}
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
		t.interval = newInterval
		t.impl = time.NewTicker(time.Duration(t.interval) * time.Second)
		logger.Info().Msgf("%s ticker interval changed from %d to %d", t.name, t.interval, newInterval)
	}
}

func (t *DynamicTicker) Stop() {
	t.impl.Stop()
}

func (ob *EVMChainClient) CheckReceiptForCointypeZeta(txHash string) error {
	hash := ethcommon.HexToHash(txHash)
	receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return err
	}
	connector, err := ob.GetConnectorContract()
	if err != nil {
		return err
	}
	for _, log := range receipt.Logs {
		event, err := connector.ParseZetaSent(*log)
		if err != nil {
			ob.PostInboundVoteForZetaSentEvent(event)
		}
	}
	return nil
}

func (ob *EVMChainClient) CheckReceiptForCointypeERC20(txHash string, client *ethclient.Client, custody *erc20custody.ERC20Custody) error {
	hash := ethcommon.HexToHash(txHash)
	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return err
	}
	for _, log := range receipt.Logs {
		zetaDeposited, err := custody.ParseDeposited(*log)
		if err != nil {
			ob.PostInboundVoteForDepositedEvents(zetaDeposited)
		}
	}
	return nil
}

func (ob *EVMChainClient) CheckReceiptForCointypeGas(tx *ethtypes.Transaction, block *ethtypes.Block) {
	receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msg("TransactionReceipt error")
		return
	}
	if receipt.Status != 1 { // 1: successful, 0: failed
		ob.logger.ExternalChainWatcher.Info().Msgf("tx %s failed; don't act", tx.Hash().Hex())
		return
	}

	from, err := ob.EvmClient.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msg("TransactionSender error; trying local recovery (assuming LondonSigner dynamic fee tx type) of sender address")
		signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
		from, err = signer.Sender(tx)
		if err != nil {
			ob.logger.ExternalChainWatcher.Err(err).Msg("local recovery of sender address failed")
			return
		}
	}
	zetaHash, err := ob.ReportTokenSentToTSS(tx.Hash(), tx.Value(), receipt, from, tx.Data())
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("Gas Deposit detected and reported: PostSend zeta tx: %s", zetaHash)
}

func (ob *EVMChainClient) PostInboundVoteForDepositedEvents(event *erc20custody.ERC20CustodyDeposited) {
	ob.logger.ExternalChainWatcher.Info().Msgf("TxBlockNumber %d Transaction Hash: %s Message : %s", event.Raw.BlockNumber, event.Raw.TxHash, event.Message)
	if bytes.Compare(event.Message, []byte(DonationMessage)) == 0 {
		ob.logger.ExternalChainWatcher.Info().Msgf("thank you rich folk for your donation!: %s", event.Raw.TxHash.Hex())
		return
	}
	zetaHash, err := ob.zetaClient.PostSend(
		"",
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
		PostSendEVMGasLimit,
		event.Asset.String(),
	)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("ZRC20Custody Deposited event detected and reported: PostSend zeta tx: %s", zetaHash)
}

func (ob *EVMChainClient) PostInboundVoteForZetaSentEvent(event *zetaconnector.ZetaConnectorNonEthZetaSent) {
	ob.logger.ExternalChainWatcher.Info().Msgf("TxBlockNumber %d Transaction Hash: %s Message : %s", event.Raw.BlockNumber, event.Raw.TxHash, event.Message)
	destChain := common.GetChainFromChainID(event.DestinationChainId.Int64())
	if destChain == nil {
		ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not supported  %d", event.DestinationChainId.Int64())
		return
	}
	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)
	cfgDest, found := ob.cfg.GetEVMConfig(destChain.ChainId)
	if !found {
		ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not present in EVMChainConfigs  %d", event.DestinationChainId.Int64())
		return
	}
	if strings.EqualFold(destAddr, cfgDest.ZetaTokenContractAddress) {
		ob.logger.ExternalChainWatcher.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
		return
	}
	zetaHash, err := ob.zetaClient.PostSend(
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
		PostSendNonEVMGasLimit,
		"",
	)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
}
