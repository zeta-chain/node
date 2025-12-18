package cctx

import (
	"fmt"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.non-eth.sol"

	zetatoolchains "github.com/zeta-chain/node/cmd/zetatool/chains"
	"github.com/zeta-chain/node/cmd/zetatool/context"
	"github.com/zeta-chain/node/pkg/chains"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	zetaevmclient "github.com/zeta-chain/node/zetaclient/chains/evm/client"
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	solrepo "github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	zetaclientConfig "github.com/zeta-chain/node/zetaclient/config"
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
		zetacoreClient = ctx.GetZetaCoreClient()
		zetaChainID    = ctx.GetConfig().ZetaChainID
		cfg            = ctx.GetConfig()
		logger         = ctx.GetLogger()
		goCtx          = ctx.GetContext()
	)

	params, err := chains.BitcoinNetParamsFromChainID(inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("unable to get bitcoin net params from chain id: %w", err)
	}

	connCfg := zetaclientConfig.BTCConfig{
		RPCUsername: cfg.BtcUser,
		RPCPassword: cfg.BtcPassword,
		RPCHost:     cfg.BtcHost,
		RPCParams:   params.Name,
	}

	rpcClient, err := client.New(connCfg, inboundChain.ChainId, logger)
	if err != nil {
		return fmt.Errorf("unable to create rpc client: %w", err)
	}

	err = rpcClient.Ping(goCtx)
	if err != nil {
		return fmt.Errorf("error ping the bitcoin server: %w", err)
	}

	res, err := zetacoreClient.Observer.GetTssAddress(goCtx, &types.QueryGetTssAddressRequest{})
	if err != nil {
		return fmt.Errorf("failed to get tss address: %w", err)
	}
	tssBtcAddress := res.GetBtc()

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}
	createConfirmationParamsIfAbsent(chainParams)

	cctxIdentifier, isConfirmed, err := zetatoolchains.BitcoinBallotIdentifier(
		ctx,
		rpcClient,
		params,
		tssBtcAddress,
		inboundHash,
		inboundChain.ChainId,
		zetaChainID,
		chainParams.InboundConfirmationSafe(),
	)
	if err != nil {
		return fmt.Errorf("failed to get bitcoin ballot identifier: %w", err)
	}
	c.CCTXIdentifier = cctxIdentifier
	c.updateInboundConfirmation(isConfirmed)
	return nil
}

// evmInboundBallotIdentifier gets the inbound ballot identifier for the inbound hash from evm chain
func (c *TrackingDetails) evmInboundBallotIdentifier(ctx *context.Context) error {
	var (
		inboundHash    = ctx.GetInboundHash()
		inboundChain   = ctx.GetInboundChain()
		zetacoreClient = ctx.GetZetaCoreClient()
		zetaChainID    = ctx.GetConfig().ZetaChainID
		goCtx          = ctx.GetContext()
	)

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}
	createConfirmationParamsIfAbsent(chainParams)

	evmClient, err := zetatoolchains.GetEvmClient(ctx, inboundChain)
	if err != nil {
		return fmt.Errorf("failed to create evm client: %w", err)
	}
	// create evm client for the observation chain
	tx, receipt, err := zetatoolchains.GetEvmTx(ctx, evmClient, inboundHash, inboundChain)
	if err != nil {
		return fmt.Errorf("failed to get tx: %w", err)
	}
	// Signer is unused
	zetaEvmClient := zetaevmclient.New(evmClient, ethtypes.NewLondonSigner(tx.ChainId()))
	isConfirmed, err := zetaEvmClient.IsTxConfirmed(goCtx, inboundHash, chainParams.InboundConfirmationSafe())
	if err != nil {
		return fmt.Errorf("unable to confirm tx: %w", err)
	}
	res, err := zetacoreClient.Observer.GetTssAddress(goCtx, &types.QueryGetTssAddressRequest{})
	if err != nil {
		return fmt.Errorf("failed to get tss address: %w", err)
	}
	tssEthAddress := res.GetEth()

	if tx.To() == nil {
		return fmt.Errorf("invalid transaction,to field is empty %s", inboundHash)
	}

	msg := &crosschaintypes.MsgVoteInbound{}
	// Create inbound vote message based on the cointype and protocol version

	switch {
	case compareAddress(tx.To().Hex(), chainParams.ConnectorContractAddress):
		{
			// build inbound vote message and post vote
			addrConnector := ethcommon.HexToAddress(chainParams.ConnectorContractAddress)
			connector, err := zetaconnector.NewZetaConnectorNonEth(addrConnector, evmClient)
			if err != nil {
				return fmt.Errorf("failed to get connector contract: %w", err)
			}
			for _, log := range receipt.Logs {
				event, err := connector.ParseZetaSent(*log)
				if err == nil && event != nil {
					msg = zetatoolchains.ZetaTokenVoteV1(event, inboundChain.ChainId)
				}
			}
		}
	case compareAddress(tx.To().Hex(), chainParams.Erc20CustodyContractAddress):
		{
			addrCustody := ethcommon.HexToAddress(chainParams.Erc20CustodyContractAddress)
			custody, err := erc20custody.NewERC20Custody(addrCustody, evmClient)
			if err != nil {
				return fmt.Errorf("failed to get custody contract: %w", err)
			}
			sender, err := evmClient.TransactionSender(goCtx, tx, receipt.BlockHash, receipt.TransactionIndex)
			if err != nil {
				return fmt.Errorf("failed to get tx sender: %w", err)
			}
			for _, log := range receipt.Logs {
				zetaDeposited, err := custody.ParseDeposited(*log)
				if err == nil && zetaDeposited != nil {
					msg = zetatoolchains.Erc20VoteV1(zetaDeposited, sender, inboundChain.ChainId, zetaChainID)
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
			msg = zetatoolchains.GasVoteV1(tx, sender, receipt.BlockNumber.Uint64(), inboundChain.ChainId, zetaChainID)
		}
	case compareAddress(tx.To().Hex(), chainParams.GatewayAddress):
		{
			gatewayAddr := ethcommon.HexToAddress(chainParams.GatewayAddress)
			gateway, err := gatewayevm.NewGatewayEVM(gatewayAddr, evmClient)
			if err != nil {
				return fmt.Errorf("failed to get gateway contract: %w", err)
			}
			for _, log := range receipt.Logs {
				if log == nil || log.Address != gatewayAddr {
					continue
				}
				eventDeposit, err := gateway.ParseDeposited(*log)
				if err == nil {
					msg = zetatoolchains.DepositInboundVoteV2(eventDeposit, inboundChain.ChainId, zetaChainID)
					break
				}
				eventDepositAndCall, err := gateway.ParseDepositedAndCalled(*log)
				if err == nil {
					msg = zetatoolchains.DepositAndCallInboundVoteV2(
						eventDepositAndCall,
						inboundChain.ChainId,
						zetaChainID,
					)
					break
				}
				eventCall, err := gateway.ParseCalled(*log)
				if err == nil {
					msg = zetatoolchains.CallInboundVoteV2(eventCall, inboundChain.ChainId, zetaChainID)
					break
				}
			}
		}
	default:
		return fmt.Errorf(
			"irrelevant transaction , not sent to any known address txHash: %s to address %s",
			inboundHash,
			tx.To(),
		)
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
		zetacoreClient = ctx.GetZetaCoreClient()
		zetaChainID    = ctx.GetConfig().ZetaChainID
		cfg            = ctx.GetConfig()
		logger         = ctx.GetLogger()
		goCtx          = ctx.GetContext()
	)
	solClient := solrpc.New(cfg.SolanaRPC)
	if solClient == nil {
		return fmt.Errorf("error creating rpc client")
	}
	solRepo := solrepo.New(solClient)

	signature, err := solana.SignatureFromBase58(inboundHash)
	if err != nil {
		return fmt.Errorf("error parsing signature: %w", err)
	}

	txResult, err := solRepo.GetTransaction(goCtx, signature)
	if err != nil {
		return fmt.Errorf("error getting transaction: %w", err)
	}

	chainParams, err := zetacoreClient.GetChainParamsForChainID(goCtx, inboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}
	createConfirmationParamsIfAbsent(chainParams)

	gatewayID, _, err := solanacontracts.ParseGatewayWithPDA(chainParams.GatewayAddress)
	if err != nil {
		return fmt.Errorf("cannot parse gateway address: %s, err: %w", chainParams.GatewayAddress, err)
	}

	// Process address lookup tables before filtering events
	resolvedTx := observer.ProcessTransactionResultWithAddressLookups(goCtx, txResult, solClient, logger, signature)

	events, err := observer.FilterInboundEvents(txResult,
		gatewayID,
		inboundChain.ChainId,
		logger,
		resolvedTx,
	)

	if err != nil {
		return fmt.Errorf("failed to filter solana inbound events: %w", err)
	}

	msg := &crosschaintypes.MsgVoteInbound{}

	// build inbound vote message from events and post to zetacore
	for _, event := range events {
		msg, err = zetatoolchains.VoteMsgFromSolEvent(event, zetaChainID)
		if err != nil {
			return fmt.Errorf("failed to create vote message: %w", err)
		}
	}

	c.CCTXIdentifier = msg.Digest()
	c.Status = PendingInboundVoting

	return nil
}

// zevmInboundBallotIdentifier gets the inbound ballot identifier for the inbound hash from zetachain
func (c *TrackingDetails) zevmInboundBallotIdentifier(ctx *context.Context) error {
	var (
		inboundHash    = ctx.GetInboundHash()
		zetacoreClient = ctx.GetZetaCoreClient()
		goCtx          = ctx.GetContext()
	)

	inboundHashToCCTX, err := zetacoreClient.Crosschain.InboundHashToCctx(
		goCtx, &crosschaintypes.QueryGetInboundHashToCctxRequest{
			InboundHash: inboundHash,
		})
	if err != nil {
		return fmt.Errorf("inbound chain is zetachain , cctx should be available in the same block: %w", err)
	}
	if len(inboundHashToCCTX.InboundHashToCctx.CctxIndex) < 1 {
		return fmt.Errorf("inbound hash does not have any cctx linked %s", inboundHash)
	}

	c.CCTXIdentifier = inboundHashToCCTX.InboundHashToCctx.CctxIndex[0]
	c.Status = PendingOutbound
	return nil
}

func compareAddress(a string, b string) bool {
	lowerA := strings.ToLower(a)
	lowerB := strings.ToLower(b)
	return strings.EqualFold(lowerA, lowerB)
}

// createConfirmationParamsIfAbsent sets the confirmation params if they are not already set
// TODO: Remove this once the confirmation migration is done
// https://github.com/zeta-chain/node/issues/3466
func createConfirmationParamsIfAbsent(chainParams *types.ChainParams) {
	if chainParams != nil && chainParams.ConfirmationParams == nil {
		chainParams.ConfirmationParams = &types.ConfirmationParams{
			SafeInboundCount:  chainParams.ConfirmationCount,
			SafeOutboundCount: chainParams.ConfirmationCount,
		}
	}
}
