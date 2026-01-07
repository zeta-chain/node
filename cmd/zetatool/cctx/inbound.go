package cctx

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	"github.com/rs/zerolog"

	zetatoolclients "github.com/zeta-chain/node/cmd/zetatool/clients"
	"github.com/zeta-chain/node/cmd/zetatool/context"
	"github.com/zeta-chain/node/cmd/zetatool/legacy"
	"github.com/zeta-chain/node/pkg/chains"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	evmobserver "github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	solobserver "github.com/zeta-chain/node/zetaclient/chains/solana/observer"
)

// CheckInbound checks the inbound chain,gets the inbound ballot identifier and updates the TrackingDetails
func (c *TrackingDetails) CheckInbound(ctx *context.Context) error {
	var (
		inboundChain = ctx.GetInboundChain()
		err          error
	)

	switch {
	case inboundChain.IsZetaChain():
		{
			err = c.zevmInboundBallotIdentifier(ctx)
			if err != nil {
				return fmt.Errorf(
					"failed to get inbound ballot for zeta chain %d, %w",
					inboundChain.ChainId,
					err,
				)
			}
		}

	case inboundChain.IsEVMChain():
		{
			err = c.evmInboundBallotIdentifier(ctx)
			if err != nil {
				return fmt.Errorf(
					"failed to get inbound ballot for evm chain %d, %w",
					inboundChain.ChainId,
					err,
				)
			}
		}
	case inboundChain.IsBitcoinChain():
		{
			err = c.btcInboundBallotIdentifier(ctx)
			if err != nil {
				return fmt.Errorf(
					"failed to get inbound ballot for bitcoin chain %d, %w",
					inboundChain.ChainId,
					err,
				)
			}
		}
	case inboundChain.IsSolanaChain():
		{
			err = c.solanaInboundBallotIdentifier(ctx)
			if err != nil {
				return fmt.Errorf(
					"failed to get inbound ballot for solana chain %d, %w",
					inboundChain.ChainId,
					err,
				)
			}
		}
	default:
		return fmt.Errorf("unsupported chain type %d", inboundChain.ChainId)
	}
	return nil
}

// btcInboundBallotIdentifier gets the inbound ballot identifier for the inbound hash from bitcoin chain
func (c *TrackingDetails) btcInboundBallotIdentifier(ctx *context.Context) error {
	var (
		inboundHash    = ctx.GetInboundHash()
		inboundChain   = ctx.GetInboundChain()
		zetacoreClient = ctx.GetZetacoreClient()
		zetaChainID    = ctx.GetConfig().ZetaChainID
		cfg            = ctx.GetConfig()
		logger         = ctx.GetLogger()
		goCtx          = ctx.GetContext()
	)

	params, err := chains.BitcoinNetParamsFromChainID(inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("unable to get bitcoin net params from chain id: %w", err)
	}

	btcClient, err := zetatoolclients.NewBitcoinClientAdapter(cfg, inboundChain, logger)
	if err != nil {
		return fmt.Errorf("unable to create rpc client: %w", err)
	}

	err = btcClient.Ping(goCtx)
	if err != nil {
		return fmt.Errorf("error ping the bitcoin server: %w", err)
	}

	tssBtcAddress, err := zetacoreClient.GetBTCTSSAddress(goCtx, inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get tss address: %w", err)
	}

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}

	feeRateMultiplier := observertypes.DefaultGasPriceMultiplier.MustFloat64()
	if !chainParams.GasPriceMultiplier.IsNil() && chainParams.GasPriceMultiplier.IsPositive() {
		feeRateMultiplier = chainParams.GasPriceMultiplier.MustFloat64()
	}

	confirmationCount := chainParams.InboundConfirmationSafe()

	// Fetch transaction from Bitcoin RPC
	hash, err := chainhash.NewHashFromStr(inboundHash)
	if err != nil {
		return fmt.Errorf("invalid tx hash: %w", err)
	}

	tx, err := btcClient.GetRawTransactionVerbose(goCtx, hash)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	isConfirmed := tx.Confirmations >= confirmationCount

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return fmt.Errorf("invalid block hash: %w", err)
	}

	blockVb, err := btcClient.GetBlockVerbose(goCtx, blockHash)
	if err != nil {
		return fmt.Errorf("failed to get block: %w", err)
	}

	// Build inbound event using the adapter method
	event, err := btcClient.GetBtcEventWithWitness(
		goCtx,
		*tx,
		tssBtcAddress,
		uint64(blockVb.Height), // #nosec G115 always positive
		feeRateMultiplier,
		zerolog.New(zerolog.Nop()),
		params,
	)
	if err != nil {
		return fmt.Errorf("failed to build btc event: %w", err)
	}
	if event == nil {
		return fmt.Errorf("no event built for btc sent to TSS")
	}

	// Decode memo and resolve amount using zetaclient's event methods
	if err := event.DecodeMemoBytes(inboundChain.ChainId); err != nil {
		return fmt.Errorf("failed to decode memo: %w", err)
	}
	if err := event.ResolveAmountForMsgVoteInbound(); err != nil {
		return fmt.Errorf("failed to resolve amount: %w", err)
	}

	// Build vote message using zetaclient's standalone function
	msg := btcobserver.NewBtcInboundVote(
		event,
		inboundChain.ChainId,
		zetaChainID,
		"",
		crosschaintypes.ConfirmationMode_SAFE,
	)
	if msg == nil {
		return fmt.Errorf("failed to create vote message for bitcoin inbound")
	}

	c.CCTXIdentifier = msg.Digest()
	c.updateInboundConfirmation(isConfirmed)
	return nil
}

// evmInboundBallotIdentifier gets the inbound ballot identifier for the inbound hash from evm chain
func (c *TrackingDetails) evmInboundBallotIdentifier(ctx *context.Context) error {
	var (
		inboundHash    = ctx.GetInboundHash()
		inboundChain   = ctx.GetInboundChain()
		zetacoreClient = ctx.GetZetacoreClient()
		zetaChainID    = ctx.GetConfig().ZetaChainID
		goCtx          = ctx.GetContext()
	)

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}

	evmClient, err := zetatoolclients.NewEVMClientAdapter(inboundChain, ctx.GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create evm client: %w", err)
	}

	tx, receipt, err := zetatoolclients.GetEvmTx(goCtx, evmClient, inboundHash, inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get tx: %w", err)
	}

	isConfirmed, err := zetatoolclients.IsTxConfirmed(goCtx, evmClient, inboundHash, chainParams.InboundConfirmationSafe())
	if err != nil {
		return fmt.Errorf("unable to confirm tx: %w", err)
	}

	tssEthAddress, err := zetacoreClient.GetEVMTSSAddress(goCtx)
	if err != nil {
		return fmt.Errorf("failed to get tss address: %w", err)
	}

	if tx.To() == nil {
		return fmt.Errorf("invalid transaction,to field is empty %s", inboundHash)
	}

	msg := &crosschaintypes.MsgVoteInbound{}

	switch {
	case compareAddress(tx.To().Hex(), chainParams.ConnectorContractAddress):
		{
			for _, log := range receipt.Logs {
				event, err := evmClient.ParseConnectorZetaSent(*log, chainParams.ConnectorContractAddress)
				if err == nil && event != nil {
					msg = legacy.ZetaTokenVoteV1(event, inboundChain.ChainId)
				}
			}
		}
	case compareAddress(tx.To().Hex(), chainParams.Erc20CustodyContractAddress):
		{
			sender, err := evmClient.TransactionSender(goCtx, tx, receipt.BlockHash, receipt.TransactionIndex)
			if err != nil {
				return fmt.Errorf("failed to get tx sender: %w", err)
			}
			for _, log := range receipt.Logs {
				zetaDeposited, err := evmClient.ParseCustodyDeposited(*log, chainParams.Erc20CustodyContractAddress)
				if err == nil && zetaDeposited != nil {
					msg = legacy.Erc20VoteV1(zetaDeposited, sender, inboundChain.ChainId, zetaChainID)
				}
			}
		}
	case compareAddress(tx.To().Hex(), tssEthAddress):
		{
			if receipt.Status != ethtypes.ReceiptStatusSuccessful {
				return fmt.Errorf("tx failed on chain %d", inboundChain.ChainId)
			}
			sender, err := evmClient.TransactionSender(goCtx, tx, receipt.BlockHash, receipt.TransactionIndex)
			if err != nil {
				return fmt.Errorf("failed to get tx sender: %w", err)
			}
			msg = legacy.GasVoteV1(tx, sender, receipt.BlockNumber.Uint64(), inboundChain.ChainId, zetaChainID)
		}
	default:
		{
			gatewayAddr := ethcommon.HexToAddress(chainParams.GatewayAddress)
			foundLog := false
			for _, log := range receipt.Logs {
				if log == nil || log.Address != gatewayAddr {
					continue
				}
				eventDeposit, err := evmClient.ParseGatewayDeposited(*log, chainParams.GatewayAddress)
				if err == nil {
					voteMsg := evmobserver.NewDepositInboundVote(
						eventDeposit,
						inboundChain.ChainId,
						zetaChainID,
						"",
						chainParams.ZetaTokenContractAddress,
						crosschaintypes.ConfirmationMode_SAFE,
					)
					msg = &voteMsg
					foundLog = true
					break
				}
				eventDepositAndCall, err := evmClient.ParseGatewayDepositedAndCalled(*log, chainParams.GatewayAddress)
				if err == nil {
					voteMsg := evmobserver.NewDepositAndCallInboundVote(
						eventDepositAndCall,
						inboundChain.ChainId,
						zetaChainID,
						"",
						chainParams.ZetaTokenContractAddress,
					)
					msg = &voteMsg
					foundLog = true
					break
				}
				eventCall, err := evmClient.ParseGatewayCalled(*log, chainParams.GatewayAddress)
				if err == nil {
					voteMsg := evmobserver.NewCallInboundVote(
						eventCall,
						inboundChain.ChainId,
						zetaChainID,
						"",
					)
					msg = &voteMsg
					foundLog = true
					break
				}
			}
			if !foundLog {
				return fmt.Errorf("no valid gateway event found for tx %s", inboundHash)
			}
		}
	}
	c.CCTXIdentifier = msg.Digest()
	c.updateInboundConfirmation(isConfirmed)
	return nil
}

// solanaInboundBallotIdentifier gets the inbound ballot identifier for the inbound hash from solana chain
func (c *TrackingDetails) solanaInboundBallotIdentifier(ctx *context.Context) error {
	var (
		inboundHash    = ctx.GetInboundHash()
		inboundChain   = ctx.GetInboundChain()
		zetacoreClient = ctx.GetZetacoreClient()
		zetaChainID    = ctx.GetConfig().ZetaChainID
		cfg            = ctx.GetConfig()
		logger         = ctx.GetLogger()
		goCtx          = ctx.GetContext()
	)

	solClient, err := zetatoolclients.NewSolanaClientAdapter(cfg.SolanaRPC)
	if err != nil {
		return fmt.Errorf("error creating rpc client: %w", err)
	}

	signature, err := solana.SignatureFromBase58(inboundHash)
	if err != nil {
		return fmt.Errorf("error parsing signature: %w", err)
	}

	txResult, err := solClient.GetTransaction(goCtx, signature)
	if err != nil {
		return fmt.Errorf("error getting transaction: %w", err)
	}

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}

	gatewayID, _, err := solanacontracts.ParseGatewayWithPDA(chainParams.GatewayAddress)
	if err != nil {
		return fmt.Errorf("cannot parse gateway address: %s, err: %w", chainParams.GatewayAddress, err)
	}

	resolvedTx := solClient.ProcessTransactionResultWithAddressLookups(goCtx, txResult, logger, signature)

	events, err := solClient.FilterInboundEvents(txResult,
		gatewayID,
		inboundChain.ChainId,
		logger,
		resolvedTx,
	)

	if err != nil {
		return fmt.Errorf("failed to filter solana inbound events: %w", err)
	}

	var msg *crosschaintypes.MsgVoteInbound

	for _, event := range events {
		msg = solobserver.NewSolanaInboundVote(
			event,
			zetaChainID,
			"",
		)
	}

	if msg == nil {
		return fmt.Errorf("no valid solana inbound event found")
	}

	c.CCTXIdentifier = msg.Digest()
	c.Status = PendingInboundVoting

	return nil
}

// zevmInboundBallotIdentifier gets the inbound ballot identifier for the inbound hash from zetachain
func (c *TrackingDetails) zevmInboundBallotIdentifier(ctx *context.Context) error {
	var (
		inboundHash    = ctx.GetInboundHash()
		zetacoreClient = ctx.GetZetacoreClient()
		goCtx          = ctx.GetContext()
	)

	inboundHashToCCTX, err := zetacoreClient.InboundHashToCctxData(goCtx, inboundHash)
	if err != nil {
		return fmt.Errorf("inbound chain is zetachain , cctx should be available in the same block: %w", err)
	}
	if len(inboundHashToCCTX.CrossChainTxs) < 1 {
		return fmt.Errorf("inbound hash does not have any cctx linked %s", inboundHash)
	}

	c.CCTXIdentifier = inboundHashToCCTX.CrossChainTxs[0].Index
	c.Status = PendingOutbound
	return nil
}

func compareAddress(a, b string) bool {
	return strings.EqualFold(a, b)
}
