package evm

import (
	"bytes"
	"context"
	sdkmath "cosmossdk.io/math"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
	"math/big"
	"strings"
)

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

	return zetabridge.GetInBoundVoteMessage(
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

	return zetabridge.GetInBoundVoteMessage(
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

	return zetabridge.GetInBoundVoteMessage(
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
