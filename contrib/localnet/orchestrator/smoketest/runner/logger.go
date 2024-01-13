package runner

import (
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"

	"github.com/fatih/color"
)

const (
	loggerSeparator = " | "
	padding         = 10
)

// Logger is a wrapper around log.Logger that adds verbosity
type Logger struct {
	verbose bool
	logger  *color.Color
	prefix  string
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

// Print prints a message to the logger
func (l *Logger) Print(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	// #nosec G104 - we are not using user input
	l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + text + "\n")
}

// Info prints a message to the logger if verbose is true
func (l *Logger) Info(message string, args ...interface{}) {
	if l.verbose {
		text := fmt.Sprintf(message, args...)
		// #nosec G104 - we are not using user input
		l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[INFO]" + text + "\n")
	}
}

// InfoLoud prints a message to the logger if verbose is true
func (l *Logger) InfoLoud(message string, args ...interface{}) {
	if l.verbose {
		text := fmt.Sprintf(message, args...)
		// #nosec G104 - we are not using user input
		l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[INFO] =======================================")
		// #nosec G104 - we are not using user input
		l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[INFO]" + text + "\n")
		// #nosec G104 - we are not using user input
		l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[INFO] =======================================")
	}
}

// Error prints an error message to the logger
func (l *Logger) Error(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	// #nosec G104 - we are not using user input
	l.logger.Printf(l.getPrefixWithPadding() + loggerSeparator + "[ERROR]" + text + "\n")
}

// CCTX prints a CCTX
func (l *Logger) CCTX(cctx crosschaintypes.CrossChainTx) {
	l.Info("🔁Cross-chain transaction: %s", cctx.Index)
	if cctx.CctxStatus != nil {
		l.Info(" CctxStatus:")
		l.Info("  Status: %s", cctx.CctxStatus.Status.String())
		if cctx.CctxStatus.StatusMessage != "" {
			l.Info("  StatusMessage: %s", cctx.CctxStatus.StatusMessage)
		}
	}
	if cctx.InboundTxParams != nil {
		l.Info(" InboundTxParams:")
		l.Info("  TxHash: %s", cctx.InboundTxParams.InboundTxObservedHash)
		l.Info("  TxHeight: %d", cctx.InboundTxParams.InboundTxObservedExternalHeight)
		l.Info("  BallotIndex: %s", cctx.InboundTxParams.InboundTxBallotIndex)
		l.Info("  Amount: %s", cctx.InboundTxParams.Amount.String())
		l.Info("  CoinType: %s", cctx.InboundTxParams.CoinType.String())
		l.Info("  SenderChainID: %d", cctx.InboundTxParams.SenderChainId)
		l.Info("  Origin: %s", cctx.InboundTxParams.TxOrigin)
		if cctx.InboundTxParams.Sender != "" {
			l.Info("  Sender: %s", cctx.InboundTxParams.Sender)
		}
		if cctx.InboundTxParams.Asset != "" {
			l.Info("  Asset: %s", cctx.InboundTxParams.Asset)
		}
	}
	if cctx.RelayedMessage != "" {
		l.Info("  RelayedMessage: %s", cctx.RelayedMessage)
	}
	for i, outTxParam := range cctx.OutboundTxParams {
		if i == 0 {
			l.Info(" OutboundTxParams:")
		} else {
			l.Info(" RevertTxParams:")
		}
		l.Info("  TxHash: %s", outTxParam.OutboundTxHash)
		l.Info("  TxHeight: %d", outTxParam.OutboundTxObservedExternalHeight)
		l.Info("  BallotIndex: %s", outTxParam.OutboundTxBallotIndex)
		l.Info("  TSSNonce: %d", outTxParam.OutboundTxTssNonce)
		l.Info("  GasLimit: %d", outTxParam.OutboundTxGasLimit)
		l.Info("  GasPrice: %s", outTxParam.OutboundTxGasPrice)
		l.Info("  GasUsed: %d", outTxParam.OutboundTxGasUsed)
		l.Info("  EffectiveGasPrice: %s", outTxParam.OutboundTxEffectiveGasPrice.String())
		l.Info("  EffectiveGasLimit: %d", outTxParam.OutboundTxEffectiveGasLimit)
		l.Info("  Amount: %s", outTxParam.Amount.String())
		l.Info("  CoinType: %s", outTxParam.CoinType.String())
		l.Info("  Receiver: %s", outTxParam.Receiver)
		l.Info("  ReceiverChainID: %d", outTxParam.ReceiverChainId)
	}
}

// EVMTransaction prints a transaction
func (l *Logger) EVMTransaction(tx ethtypes.Transaction) {
	l.Info("⤴️EVM transaction: %s", tx.Hash().Hex())
	l.Info("  To: %s", tx.To().Hex())
	l.Info("  Value: %d", tx.Value())
	l.Info("  Gas: %d", tx.Gas())
	l.Info("  GasPrice: %d", tx.GasPrice())
	l.Info("  Data: %s", tx.Data())
}

// EVMReceipt prints a receipt
func (l *Logger) EVMReceipt(receipt ethtypes.Receipt) {
	l.Info("📄EVM receipt: %s", receipt.TxHash.Hex())
	l.Info("  BlockNumber: %d", receipt.BlockNumber)
	l.Info("  GasUsed: %d", receipt.GasUsed)
	l.Info("  ContractAddress: %s", receipt.ContractAddress.Hex())
	l.Info("  Status: %d", receipt.Status)
}

func (l *Logger) getPrefixWithPadding() string {
	// add padding to prefix
	prefix := l.prefix
	for i := len(l.prefix); i < padding; i++ {
		prefix += " "
	}
	return prefix
}
