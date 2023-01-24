package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// FIXME: use more specific error types & codes
func (k msgServer) VoteOnObservedInboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedInboundTx) (*types.MsgVoteOnObservedInboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := zetaObserverTypes.ObservationType_InBoundTx
	observationChain := zetaObserverTypes.ParseCommonChaintoObservationChain(msg.SenderChain)
	//Check is msg.Creator is authorized to vote
	ok, err := k.IsAuthorized(ctx, msg.Creator, observationChain, observationType.String())
	if !ok {
		return nil, err
	}

	index := msg.Digest()
	// Add votes and Set Ballot
	ballot, err := k.GetBallot(ctx, index, observationChain, observationType)
	if err != nil {
		return nil, err
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
	cctx := k.CreateNewCCTX(ctx, msg, index, types.CctxStatus_PendingInbound)

	EmitEventCCTXCreated(ctx, cctx)
	// FinalizeInbound updates CCTX Prices and Nonce
	// Aborts is any of the updates fail
	//TODO : move to a separate function
	toChain, err := common.ParseChain(msg.ReceiverChain)
	if err != nil {
		return nil, err
	}
	if toChain == common.ZETAChain {
		cctx.InBoundTxParams.InBoundTxFinalizedZetaHeight = uint64(ctx.BlockHeader().Height)
		if msg.CoinType != common.CoinType_Zeta {
			foreignCoinList := k.fungibleKeeper.GetAllForeignCoins(ctx)
			found := false
			//TODO :  Foreign coins to use type-foreign chain , It's not a good idea to iterate here,as this handler might be called frequently
			var coin fungibletypes.ForeignCoins
			for _, foreignCoin := range foreignCoinList {
				fmt.Printf("%s %s %s %s\n", cctx.InBoundTxParams.Asset, msg.SenderChain, foreignCoin.Erc20ContractAddress, foreignCoin.ForeignChain)
				if (msg.CoinType == common.CoinType_Gas || foreignCoin.Erc20ContractAddress == cctx.InBoundTxParams.Asset) && foreignCoin.ForeignChain == msg.SenderChain {
					found = true
					coin = foreignCoin
					break
				}
			}
			// TODO Break it into subfunctions, and only call Changestatus to aborted , when the subfunction returns an error
			if !found {
				errMsg := fmt.Sprintf("cannot get gas coin on chain %s", msg.SenderChain)
				cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, errMsg, cctx.LogIdentifierForCCTX())
				k.SetCrossChainTx(ctx, cctx)
				return &types.MsgVoteOnObservedInboundTxResponse{}, nil
			}
			to := ethcommon.HexToAddress(msg.Receiver)
			amount, ok := big.NewInt(0).SetString(msg.ZetaBurnt, 10)
			if !ok {
				errMsg := fmt.Sprintf("cannot parse zetaBurnt: %s", msg.ZetaBurnt)
				cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, errMsg, cctx.LogIdentifierForCCTX())
				k.SetCrossChainTx(ctx, cctx)
				return &types.MsgVoteOnObservedInboundTxResponse{}, nil
			}
			var tx *evmtypes.MsgEthereumTxResponse
			if len(msg.Message) == 0 { // no message; transfer
				tx, err = k.fungibleKeeper.DepositZRC20(ctx, ethcommon.HexToAddress(coin.Zrc20ContractAddress), to, amount)
				if err != nil {
					errMsg := fmt.Sprintf("cannot DepositZRC20, %s", err.Error())
					cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, errMsg, cctx.LogIdentifierForCCTX())
					k.SetCrossChainTx(ctx, cctx)
					return &types.MsgVoteOnObservedInboundTxResponse{}, nil
				}
			} else { // non-empty message = [contractaddress, calldata]
				contract, data, err := parseContractAndData(msg.Message)
				if err != nil {
					errMsg := fmt.Sprintf("cannot parse contract and data, %s", err.Error())
					cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, errMsg, cctx.LogIdentifierForCCTX())
					k.SetCrossChainTx(ctx, cctx)
					return &types.MsgVoteOnObservedInboundTxResponse{}, nil
				}
				if len(data) > 0 {
					tx, err = k.fungibleKeeper.DepositZRC20AndCallContract(ctx, ethcommon.HexToAddress(coin.Zrc20ContractAddress), amount, contract, data)
				} else { // data = empty
					tx, err = k.fungibleKeeper.DepositZRC20(ctx, ethcommon.HexToAddress(coin.Zrc20ContractAddress), contract, amount)
				}
				if err != nil { // prepare to revert
					errMsg := fmt.Sprintf("cannot DepositZRC20: %s", err.Error())
					cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, errMsg, cctx.LogIdentifierForCCTX())
					k.SetCrossChainTx(ctx, cctx)
					return &types.MsgVoteOnObservedInboundTxResponse{}, nil
				}
				if !tx.Failed() {
					logs := evmtypes.LogsToEthereum(tx.Logs)
					ctx = ctx.WithValue("inCctxIndex", cctx.Index)
					txOrigin := msg.TxOrigin
					if txOrigin == "" {
						txOrigin = msg.Sender
					}
					err = k.ProcessWithdrawalLogs(ctx, logs, contract, txOrigin)
					if err != nil {
						errMsg := fmt.Sprintf("cannot process withdrawal event: %s", err.Error())
						cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, errMsg, cctx.LogIdentifierForCCTX())
						k.SetCrossChainTx(ctx, cctx)
						return &types.MsgVoteOnObservedInboundTxResponse{}, nil
					}
					ctx.EventManager().EmitEvent(
						sdk.NewEvent(sdk.EventTypeMessage,
							sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
							sdk.NewAttribute("action", "depositZRC4AndCallContract"),
							sdk.NewAttribute("contract", contract.String()),
							sdk.NewAttribute("data", hex.EncodeToString(data)),
							sdk.NewAttribute("cctxIndex", cctx.Index),
						),
					)
				}
			}
			fmt.Printf("=======  tx: %s\n", tx.Hash)
			fmt.Printf("vmerror: %s\n", tx.VmError)
			fmt.Printf("=======  tx: %s\n", tx.Hash)

			cctx.OutBoundTxParams.OutBoundTxHash = tx.Hash
			cctx.CctxStatus.Status = types.CctxStatus_OutboundMined
			k.SetCrossChainTx(ctx, cctx)
		} else {
			toBytes := ethcommon.HexToAddress(msg.Receiver).Bytes()
			to := sdk.AccAddress(toBytes)
			amount, ok := big.NewInt(0).SetString(msg.ZetaBurnt, 10)
			if !ok {
				errMsg := fmt.Sprintf("cannot parse zetaBurnt: %s", msg.ZetaBurnt)
				cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, errMsg, cctx.LogIdentifierForCCTX())
				k.SetCrossChainTx(ctx, cctx)
				return &types.MsgVoteOnObservedInboundTxResponse{}, nil
			}
			err := k.fungibleKeeper.MintZetaToEVMAccount(ctx, to, amount)
			if err != nil {
				errMsg := fmt.Sprintf("cannot MintZetaToEVMAccount: %s", err.Error())
				cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, errMsg, cctx.LogIdentifierForCCTX())
				k.SetCrossChainTx(ctx, cctx)
				return &types.MsgVoteOnObservedInboundTxResponse{}, nil
			}
			cctx.CctxStatus.Status = types.CctxStatus_OutboundMined
			k.SetCrossChainTx(ctx, cctx)
		}
	} else {
		err = k.FinalizeInbound(ctx, &cctx, msg.ReceiverChain, len(ballot.VoterList))
		if err != nil {
			cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
			ctx.Logger().Error(err.Error())
			k.SetCrossChainTx(ctx, cctx)
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}

		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_PendingOutbound, "Status Changed to Pending Outbound", cctx.LogIdentifierForCCTX())
		EmitEventInboundFinalized(ctx, &cctx)
		k.SetCrossChainTx(ctx, cctx)
	}

	return &types.MsgVoteOnObservedInboundTxResponse{}, nil
}

func (k msgServer) FinalizeInbound(ctx sdk.Context, cctx *types.CrossChainTx, receiveChain string, numberofobservers int) error {
	cctx.InBoundTxParams.InBoundTxFinalizedZetaHeight = uint64(ctx.BlockHeader().Height)
	k.UpdateLastBlockHeight(ctx, cctx)
	bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
	cctx.OutBoundTxParams.Broadcaster = uint64(bftTime.Nanosecond() % numberofobservers)

	err := k.UpdatePrices(ctx, receiveChain, cctx)
	if err != nil {
		return err
	}
	err = k.UpdateNonce(ctx, receiveChain, cctx)
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
