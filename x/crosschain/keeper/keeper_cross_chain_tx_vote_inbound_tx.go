package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"math/big"
)

// FIXME: use more specific error types & codes
func (k msgServer) VoteOnObservedInboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedInboundTx) (*types.MsgVoteOnObservedInboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := zetaObserverTypes.ObservationType_InBoundTx
	observationChain, found := k.zetaObserverKeeper.GetChainFromChainID(ctx, msg.SenderChain)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Observation %s", msg.SenderChain, observationType.String()))
	}
	receiverChain, found := k.zetaObserverKeeper.GetChainFromChainID(ctx, msg.ReceiverChain)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Observation %s", msg.ReceiverChain, observationType.String()))
	}

	ok, err := k.IsAuthorized(ctx, msg.Creator, observationChain, observationType)
	if !ok {
		return nil, err
	}

	index := msg.Digest()
	// Add votes and Set Ballot
	ballot, isNew, err := k.GetBallot(ctx, index, observationChain, observationType)
	if err != nil {
		return nil, err
	}
	if isNew {
		EmitEventBallotCreated(ctx, ballot, msg.InTxHash, observationChain.String())
	}
	// AddVoteToBallot adds a vote and sets the ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, zetaObserverTypes.VoteType_SuccessObservation)
	if err != nil {
		return nil, err
	}
	// CheckIfBallotIsFinalized checks status and sets the ballot if finalized

	ballot, isFinalized := k.CheckIfBallotIsFinalized(ctx, ballot)
	if !isFinalized {
		return &types.MsgVoteOnObservedInboundTxResponse{}, nil
	}

	// ******************************************************************************
	// below only happens when ballot is finalized: exactly when threshold vote is in
	// ******************************************************************************

	// Inbound Ballot has been finalized , Create CCTX
	// New CCTX can only set either to Aborted or PendingOutbound
	cctx := k.CreateNewCCTX(ctx, msg, index, types.CctxStatus_PendingOutbound, observationChain, receiverChain)
	// FinalizeInbound updates CCTX Prices and Nonce
	// Aborts is any of the updates fail
	switch receiverChain.ChainName {
	case zetaObserverTypes.ChainName_ZetaChain:
		err = k.HandleEVMDeposit(ctx, &cctx, *msg, observationChain)
		if err != nil {
			cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
			k.SetCrossChainTx(ctx, cctx)
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}
		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_PendingOutbound, "Status Changed to Pending Outbound", cctx.LogIdentifierForCCTX())
	default:
		err = k.FinalizeInbound(ctx, &cctx, *receiverChain, len(ballot.VoterList))
		if err != nil {
			cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
			k.SetCrossChainTx(ctx, cctx)
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}

		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_PendingOutbound, "Status Changed to Pending Outbound", cctx.LogIdentifierForCCTX())
	}
	EmitEventInboundFinalized(ctx, &cctx)
	k.SetCrossChainTx(ctx, cctx)
	return &types.MsgVoteOnObservedInboundTxResponse{}, nil
}

func (k msgServer) FinalizeInbound(ctx sdk.Context, cctx *types.CrossChainTx, receiveChain zetaObserverTypes.Chain, numberofobservers int) error {
	cctx.InBoundTxParams.InBoundTxFinalizedZetaHeight = uint64(ctx.BlockHeader().Height)
	k.UpdateLastBlockHeight(ctx, cctx)
	bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
	cctx.OutBoundTxParams.Broadcaster = uint64(bftTime.Nanosecond() % numberofobservers)

	err := k.UpdatePrices(ctx, receiveChain.ChainName.String(), cctx)
	if err != nil {
		return err
	}
	err = k.UpdateNonce(ctx, receiveChain.ChainName.String(), cctx)
	if err != nil {
		return err
	}
	return nil
}

func (k msgServer) UpdateLastBlockHeight(ctx sdk.Context, msg *types.CrossChainTx) {
	lastblock, isFound := k.GetLastBlockHeight(ctx, msg.InBoundTxParams.SenderChain)
	if !isFound {
		lastblock = types.LastBlockHeight{
			Creator:           msg.Creator,
			Index:             msg.InBoundTxParams.SenderChain, // ?
			Chain:             msg.InBoundTxParams.SenderChain,
			LastSendHeight:    msg.InBoundTxParams.InBoundTxObservedExternalHeight,
			LastReceiveHeight: 0,
		}
	} else {
		lastblock.LastSendHeight = msg.InBoundTxParams.InBoundTxObservedExternalHeight
	}
	k.SetLastBlockHeight(ctx, lastblock)
}

func (k msgServer) HandleEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx, msg types.MsgVoteOnObservedInboundTx, senderChain *zetaObserverTypes.Chain) error {

	gasCoin, found := k.fungibleKeeper.GetGasCoinForForeignCoin(ctx, senderChain.ChainName.String())
	if !found {
		return types.ErrGasCoinNotFound
	}
	to := ethcommon.HexToAddress(msg.Receiver)
	amount, ok := big.NewInt(0).SetString(msg.ZetaBurnt, 10)
	if !ok {
		return errors.Wrap(types.ErrFloatParseError, fmt.Sprintf("cannot parse zetaBurnt: %s", msg.ZetaBurnt))
	}
	depositContract := ethcommon.Address{}
	switch msg.CoinType {
	// We are only using Gas Type right now
	case common.CoinType_Gas:
		{
			depositContract = ethcommon.HexToAddress(gasCoin.Zrc20ContractAddress)
		}

	}
	var tx *evmtypes.MsgEthereumTxResponse
	if len(msg.Message) == 0 { // no message; transfer
		var txNoWithdraw *evmtypes.MsgEthereumTxResponse
		txNoWithdraw, err := k.fungibleKeeper.DepositZRC20(ctx, ethcommon.HexToAddress(gasCoin.Zrc20ContractAddress), to, amount)
		if err != nil {
			return errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
		}
		tx = txNoWithdraw
	} else { // non-empty message = [contractaddress, calldata]
		var txWithWithdraw *evmtypes.MsgEthereumTxResponse
		contract, data, err := parseContractAndData(msg.Message)
		if err != nil {
			return errors.Wrap(types.ErrUnableToParseContract, err.Error())
		}
		txWithWithdraw, err = k.fungibleKeeper.DepositZRC20AndCallContract(ctx, depositContract, amount, contract, data)
		if err != nil { // prepare to revert
			return errors.Wrap(types.ErrUnableToDepositZRC20, err.Error())
		}
		if !txWithWithdraw.Failed() {
			logs := evmtypes.LogsToEthereum(txWithWithdraw.Logs)
			ctx = ctx.WithValue("inCctxIndex", cctx.Index)
			err = k.ProcessWithdrawalEvent(ctx, logs, contract, msg.TxOrigin)
			if err != nil {
				return errors.Wrap(types.ErrCannotProcessWithdrawal, err.Error())
			}
			// TODO Add Event Types as constants
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, "crosschain"),
					sdk.NewAttribute("action", "depositZRC4AndCallContract"),
					sdk.NewAttribute("contract", contract.String()),
					sdk.NewAttribute("data", hex.EncodeToString(data)),
					sdk.NewAttribute("cctxIndex", cctx.Index),
				),
			)
		}
		tx = txWithWithdraw
	}

	cctx.OutBoundTxParams.OutBoundTxHash = tx.Hash
	cctx.CctxStatus.Status = types.CctxStatus_OutboundMined
	return nil
}

// message is hex encoded byte array
// [ contractAddress calldata ]
// [ 20B, variable]
func parseContractAndData(message string) (ethcommon.Address, []byte, error) {
	data, err := hex.DecodeString(message)
	if err != nil {
		return ethcommon.Address{}, nil, err
	}
	if len(data) < 20 {
		err = fmt.Errorf("invalid message length")
		return ethcommon.Address{}, nil, err
	}
	contractAddress := ethcommon.BytesToAddress(data[:20])
	data = data[20:]
	return contractAddress, data, nil
}
