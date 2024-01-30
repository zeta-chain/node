package zetaclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
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
	bytesPerKB        = 1000
	bytesEmptyTx      = 10  // an empty tx is about 10 bytes
	bytesPerInput     = 41  // each input is about 41 bytes
	bytesPerOutput    = 31  // each output is about 31 bytes
	bytes1stWitness   = 110 // the 1st witness incurs about 110 bytes and it may vary
	bytesPerWitness   = 108 // each additional witness incurs about 108 bytes and it may vary
)

var (
	BtcOutTxBytesMin        uint64
	BtcOutTxBytesMax        uint64
	BtcOutTxBytesDepositor  uint64
	BtcOutTxBytesWithdrawer uint64
	BtcDepositorFeeMin      float64
)

func init() {
	BtcOutTxBytesMin = EstimateSegWitTxSize(2, 3)      // 403B, estimated size for a 2-input, 3-output SegWit tx
	BtcOutTxBytesMax = EstimateSegWitTxSize(21, 3)     // 3234B, estimated size for a 21-input, 3-output SegWit tx
	BtcOutTxBytesDepositor = SegWitTxSizeDepositor()   // 149B, the outtx size incurred by the depositor
	BtcOutTxBytesWithdrawer = SegWitTxSizeWithdrawer() // 254B, the outtx size incurred by the withdrawer

	// depositor fee calculation is based on a fixed fee rate of 5 sat/byte just for simplicity.
	// In reality, the fee rate on UTXO deposit is different from the fee rate when the UTXO is spent.
	BtcDepositorFeeMin = DepositorFee(5) // 0.00000745 (5 * 149B / 100000000), the minimum deposit fee in BTC for 5 sat/byte
}

func IsEnvFlagEnabled(flag string) bool {
	value := os.Getenv(flag)
	return value == "true" || value == "1"
}

func PrettyPrintStruct(val interface{}) (string, error) {
	prettyStruct, err := json.MarshalIndent(
		val,
		"",
		" ",
	)
	if err != nil {
		return "", err
	}
	return string(prettyStruct), nil
}

// FeeRateToSatPerByte converts a fee rate in BTC/KB to sat/byte.
func FeeRateToSatPerByte(rate float64) *big.Int {
	// #nosec G701 always in range
	satPerKB := new(big.Int).SetInt64(int64(rate * satoshiPerBitcoin))
	return new(big.Int).Div(satPerKB, big.NewInt(bytesPerKB))
}

// EstimateSegWitTxSize estimates SegWit tx size
func EstimateSegWitTxSize(numInputs uint64, numOutputs uint64) uint64 {
	if numInputs == 0 {
		return 0
	}
	bytesInput := numInputs * bytesPerInput
	bytesOutput := numOutputs * bytesPerOutput
	bytesWitness := bytes1stWitness + (numInputs-1)*bytesPerWitness
	return bytesEmptyTx + bytesInput + bytesOutput + bytesWitness
}

// SegWitTxSizeDepositor returns SegWit tx size (149B) incurred by the depositor
func SegWitTxSizeDepositor() uint64 {
	return bytesPerInput + bytesPerWitness
}

// SegWitTxSizeWithdrawer returns SegWit tx size (254B) incurred by the withdrawer
func SegWitTxSizeWithdrawer() uint64 {
	bytesInput := uint64(1) * bytesPerInput   // nonce mark
	bytesOutput := uint64(3) * bytesPerOutput // 3 outputs: new nonce mark, payment, change
	return bytesEmptyTx + bytesInput + bytesOutput + bytes1stWitness
}

// DepositorFee calculates the depositor fee in BTC for a given sat/byte fee rate
// Note: the depositor fee is charged in order to cover the cost of spending the deposited UTXO in the future
func DepositorFee(satPerByte int64) float64 {
	return float64(satPerByte) * float64(BtcOutTxBytesDepositor) / satoshiPerBitcoin
}

func GetSatoshis(btc float64) (int64, error) {
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

func NewDynamicTicker(name string, interval uint64) (*DynamicTicker, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("non-positive ticker interval %d for %s", interval, name)
	}

	return &DynamicTicker{
		name:     name,
		interval: interval,
		impl:     time.NewTicker(time.Duration(interval) * time.Second),
	}, nil
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

// CheckEvmTxLog checks the basics of an EVM tx log
func (ob *EVMChainClient) CheckEvmTxLog(vLog *ethtypes.Log, wantAddress ethcommon.Address, wantHash string, wantTopics int) error {
	if vLog.Removed {
		return fmt.Errorf("log is removed, chain reorg?")
	}
	if vLog.Address != wantAddress {
		return fmt.Errorf("log emitter address mismatch: want %s got %s", wantAddress.Hex(), vLog.Address.Hex())
	}
	if vLog.TxHash.Hex() == "" {
		return fmt.Errorf("log tx hash is empty: %d %s", vLog.BlockNumber, vLog.TxHash.Hex())
	}
	if wantHash != "" && vLog.TxHash.Hex() != wantHash {
		return fmt.Errorf("log tx hash mismatch: want %s got %s", wantHash, vLog.TxHash.Hex())
	}
	if len(vLog.Topics) != wantTopics {
		return fmt.Errorf("number of topics mismatch: want %d got %d", wantTopics, len(vLog.Topics))
	}
	return nil
}

// HasEnoughConfirmations checks if the given receipt has enough confirmations
func (ob *EVMChainClient) HasEnoughConfirmations(receipt *ethtypes.Receipt, lastHeight uint64) bool {
	confHeight := receipt.BlockNumber.Uint64() + ob.GetChainParams().ConfirmationCount
	return lastHeight >= confHeight
}

// GetTransactionSender returns the sender of the given transaction
func (ob *EVMChainClient) GetTransactionSender(tx *ethtypes.Transaction, blockHash ethcommon.Hash, txIndex uint) (ethcommon.Address, error) {
	sender, err := ob.evmClient.TransactionSender(context.Background(), tx, blockHash, txIndex)
	if err != nil {
		// trying local recovery (assuming LondonSigner dynamic fee tx type) of sender address
		signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
		sender, err = signer.Sender(tx)
		if err != nil {
			ob.logger.ExternalChainWatcher.Err(err).Msgf("can't recover the sender from tx hash %s chain %d", tx.Hash().Hex(), ob.chain.ChainId)
			return ethcommon.Address{}, err
		}
	}
	return sender, nil
}

func (ob *EVMChainClient) GetInboundVoteMsgForDepositedEvent(event *erc20custody.ERC20CustodyDeposited, sender ethcommon.Address) *types.MsgVoteOnObservedInboundTx {
	if bytes.Equal(event.Message, []byte(DonationMessage)) {
		ob.logger.ExternalChainWatcher.Info().Msgf("thank you rich folk for your donation! tx %s chain %d", event.Raw.TxHash.Hex(), ob.chain.ChainId)
		return nil
	}
	message := hex.EncodeToString(event.Message)
	ob.logger.ExternalChainWatcher.Info().Msgf("ERC20CustodyDeposited inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, sender.Hex(), event.Amount.String(), message)

	return GetInBoundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		"",
		clienttypes.BytesToEthHex(event.Recipient),
		ob.zetaClient.ZetaChain().ChainId,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		common.CoinType_ERC20,
		event.Asset.String(),
		ob.zetaClient.GetKeys().GetOperatorAddress().String(),
		event.Raw.Index,
	)
}

func (ob *EVMChainClient) GetInboundVoteMsgForZetaSentEvent(event *zetaconnector.ZetaConnectorNonEthZetaSent) *types.MsgVoteOnObservedInboundTx {
	destChain := common.GetChainFromChainID(event.DestinationChainId.Int64())
	if destChain == nil {
		ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not supported  %d", event.DestinationChainId.Int64())
		return nil
	}
	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)
	if !destChain.IsZetaChain() {
		cfgDest, found := ob.cfg.GetEVMConfig(destChain.ChainId)
		if !found {
			ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not present in EVMChainConfigs  %d", event.DestinationChainId.Int64())
			return nil
		}
		if strings.EqualFold(destAddr, cfgDest.ZetaTokenContractAddress) {
			ob.logger.ExternalChainWatcher.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
			return nil
		}
	}
	message := base64.StdEncoding.EncodeToString(event.Message)
	ob.logger.ExternalChainWatcher.Info().Msgf("ZetaSent inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, event.ZetaTxSenderAddress.Hex(), event.ZetaValueAndGas.String(), message)

	return GetInBoundVoteMessage(
		event.ZetaTxSenderAddress.Hex(),
		ob.chain.ChainId,
		event.SourceTxOriginAddress.Hex(),
		clienttypes.BytesToEthHex(event.DestinationAddress),
		destChain.ChainId,
		sdkmath.NewUintFromBigInt(event.ZetaValueAndGas),
		message,
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		event.DestinationGasLimit.Uint64(),
		common.CoinType_Zeta,
		"",
		ob.zetaClient.GetKeys().GetOperatorAddress().String(),
		event.Raw.Index,
	)
}

func (ob *EVMChainClient) GetInboundVoteMsgForTokenSentToTSS(tx *ethtypes.Transaction, sender ethcommon.Address, blockNumber uint64) *types.MsgVoteOnObservedInboundTx {
	if bytes.Equal(tx.Data(), []byte(DonationMessage)) {
		ob.logger.ExternalChainWatcher.Info().Msgf("thank you rich folk for your donation! tx %s chain %d", tx.Hash().Hex(), ob.chain.ChainId)
		return nil
	}
	message := ""
	if len(tx.Data()) != 0 {
		message = hex.EncodeToString(tx.Data())
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("TSS inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, tx.Hash().Hex(), blockNumber, sender.Hex(), tx.Value().String(), hex.EncodeToString(tx.Data()))

	return GetInBoundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		sender.Hex(),
		sender.Hex(),
		ob.zetaClient.ZetaChain().ChainId,
		sdkmath.NewUintFromBigInt(tx.Value()),
		message,
		tx.Hash().Hex(),
		blockNumber,
		90_000,
		common.CoinType_Gas,
		"",
		ob.zetaClient.GetKeys().GetOperatorAddress().String(),
		0, // not a smart contract call
	)
}
