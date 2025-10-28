package runner

import (
	"encoding/hex"
	"fmt"
	"sync"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

const (
	loggerSeparator = " | "
	padding         = 12
)

// Logger is a wrapper around log.Logger that adds verbosity
type Logger struct {
	verbose bool
	logger  *color.Color
	prefix  string
	mu      sync.Mutex
}

// NewLogger creates a new Logger
func NewLogger(verbose bool, printColor color.Attribute, prefix string) *Logger {
	// trim prefix to padding
	if len(prefix) > padding {
		prefix = prefix[:padding]
	}

	return &Logger{
		verbose: verbose,
		logger:  color.New(printColor),
		prefix:  prefix,
	}
}

// SetColor sets the color of the logger
func (l *Logger) SetColor(printColor color.Attribute) {
	l.logger = color.New(printColor)
}

// Prefix returns the prefix of the logger
func (l *Logger) Prefix() string {
	return l.getPrefixWithPadding() + loggerSeparator
}

// Print prints a message to the logger
func (l *Logger) Print(message string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	text := fmt.Sprintf(message, args...)
	// #nosec G104 - we are not using user input
	l.logger.Print(l.getPrefixWithPadding() + loggerSeparator + text + "\n")
}

// PrintNoPrefix prints a message to the logger without the prefix
func (l *Logger) PrintNoPrefix(message string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	text := fmt.Sprintf(message, args...)
	// #nosec G104 - we are not using user input
	_, _ = l.logger.Print(text + "\n")
}

// Info prints a message to the logger if verbose is true
func (l *Logger) Info(message string, args ...any) {
	if !l.verbose {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var (
		content = fmt.Sprintf(message, args...)
		line    = fmt.Sprintf("%s%s[INFO] %s \n", l.getPrefixWithPadding(), loggerSeparator, content)
	)

	// #nosec G104 - we are not using user input
	_, _ = l.logger.Print(line)
}

// InfoLoud prints a message to the logger if verbose is true
func (l *Logger) InfoLoud(message string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.verbose {
		text := fmt.Sprintf(message, args...)
		// #nosec G104 - we are not using user input
		l.logger.Print(l.getPrefixWithPadding() + loggerSeparator + "[INFO] =======================================")
		// #nosec G104 - we are not using user input
		l.logger.Print(l.getPrefixWithPadding() + loggerSeparator + "[INFO]" + text + "\n")
		// #nosec G104 - we are not using user input
		l.logger.Print(l.getPrefixWithPadding() + loggerSeparator + "[INFO] =======================================")
	}
}

// Error prints an error message to the logger
func (l *Logger) Error(message string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	text := fmt.Sprintf(message, args...)
	// #nosec G104 - we are not using user input
	l.logger.Print(l.getPrefixWithPadding() + loggerSeparator + "[ERROR]" + text + "\n")
}

// CCTX prints a CCTX
func (l *Logger) CCTX(cctx crosschaintypes.CrossChainTx, name string) {
	l.Info(" %s cross-chain transaction: %s", name, cctx.Index)
	if cctx.CctxStatus != nil {
		l.Info(" CctxStatus:")
		l.Info("  Status: %s", cctx.CctxStatus.Status.String())
		if cctx.CctxStatus.StatusMessage != "" {
			l.Info("  StatusMessage: %s", cctx.CctxStatus.StatusMessage)
		}
	}
	if cctx.InboundParams != nil {
		l.Info(" InboundParams:")
		l.Info("  TxHash: %s", cctx.InboundParams.ObservedHash)
		l.Info("  TxHeight: %d", cctx.InboundParams.ObservedExternalHeight)
		l.Info("  BallotIndex: %s", cctx.InboundParams.BallotIndex)
		l.Info("  Amount: %s", cctx.InboundParams.Amount.String())
		l.Info("  FungibleTokenCoinType: %s", cctx.InboundParams.CoinType.String())
		l.Info("  SenderChainID: %d", cctx.InboundParams.SenderChainId)
		l.Info("  Origin: %s", cctx.InboundParams.TxOrigin)
		if cctx.InboundParams.Sender != "" {
			l.Info("  Sender: %s", cctx.InboundParams.Sender)
		}
		if cctx.InboundParams.Asset != "" {
			l.Info("  Asset: %s", cctx.InboundParams.Asset)
		}
	}
	if cctx.RelayedMessage != "" {
		l.Info("  RelayedMessage: %s", cctx.RelayedMessage)
	}
	for i, outTxParam := range cctx.OutboundParams {
		if i == 0 {
			l.Info(" OutboundTxParams:")
		} else {
			l.Info(" RevertTxParams:")
		}
		l.Info("  TxHash: %s", outTxParam.Hash)
		l.Info("  TxHeight: %d", outTxParam.ObservedExternalHeight)
		l.Info("  BallotIndex: %s", outTxParam.BallotIndex)
		l.Info("  TSSNonce: %d", outTxParam.TssNonce)
		l.Info("  CallOptions: %+v", outTxParam.CallOptions)
		l.Info("  GasLimit: %d", outTxParam.GasLimit)
		l.Info("  GasPrice: %s", outTxParam.GasPrice)
		l.Info("  GasUsed: %d", outTxParam.GasUsed)
		l.Info("  EffectiveGasPrice: %s", outTxParam.EffectiveGasPrice.String())
		l.Info("  EffectiveGasLimit: %d", outTxParam.EffectiveGasLimit)
		l.Info("  Amount: %s", outTxParam.Amount.String())
		l.Info("  FungibleTokenCoinType: %s", outTxParam.CoinType.String())
		l.Info("  Receiver: %s", outTxParam.Receiver)
		l.Info("  ReceiverChainID: %d", outTxParam.ReceiverChainId)
	}
}

// EVMTransaction prints a transaction
func (l *Logger) EVMTransaction(tx *ethtypes.Transaction, name string) {
	l.Info(" %s EVM transaction: %s", name, tx.Hash().Hex())
	if tx.To() != nil {
		l.Info("  To: %s", tx.To().Hex())
	} else {
		l.Info("  To: <nil>")
	}
	l.Info("  Value: %d", tx.Value())
	l.Info("  Gas: %d", tx.Gas())
	l.Info("  GasPrice: %d", tx.GasPrice())
}

// EVMReceipt prints a receipt
func (l *Logger) EVMReceipt(receipt ethtypes.Receipt, name string) {
	l.Info(" %s EVM receipt: %s", name, receipt.TxHash.Hex())
	l.Info("  BlockNumber: %d", receipt.BlockNumber)
	l.Info("  GasUsed: %d", receipt.GasUsed)
	l.Info("  ContractAddress: %s", receipt.ContractAddress.Hex())
	l.Info("  Status: %d", receipt.Status)
}

// ZRC20Withdrawal prints a ZRC20Withdrawal event
func (l *Logger) ZRC20Withdrawal(
	contract interface {
		ParseWithdrawal(ethtypes.Log) (*zrc20.ZRC20Withdrawal, error)
	},
	receipt ethtypes.Receipt,
	name string,
) {
	for _, log := range receipt.Logs {
		event, err := contract.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		l.Info(
			" %s ZRC20Withdrawal: from %s, to %x, value %d, gasfee %d",
			name,
			event.From.Hex(),
			event.To,
			event.Value,
			event.GasFee,
		)
	}
}

type depositParser interface {
	ParseDeposited(ethtypes.Log) (*gatewayevm.GatewayEVMDeposited, error)
}

// GatewayDeposit prints a GatewayDeposit event
func (l *Logger) GatewayDeposit(
	contract depositParser,
	receipt ethtypes.Receipt,
	name string,
) {
	for _, log := range receipt.Logs {
		event, err := contract.ParseDeposited(*log)
		if err != nil {
			continue
		}

		l.Info(" Gateway Deposit: %s", name)
		l.Info("  Sender: %s", event.Sender.Hex())
		l.Info("  Receiver: %s", event.Receiver.Hex())
		l.Info("  Amount: %s", event.Amount.String())
		l.Info("  Asset: %s", event.Asset.Hex())
		l.Info("  Payload: %s", hex.EncodeToString(event.Payload))
	}
}

func (l *Logger) getPrefixWithPadding() string {
	// add padding to prefix
	prefix := l.prefix
	for i := len(l.prefix); i < padding; i++ {
		prefix += " "
	}
	return prefix
}
