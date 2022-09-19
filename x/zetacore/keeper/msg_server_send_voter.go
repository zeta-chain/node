package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/common"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"math/big"
)

var (
	OneEighteen, _ = big.NewInt(0).SetString("1000000000000000000", 10)
)

func (k msgServer) SendVoter(goCtx context.Context, msg *types.MsgSendVoter) (*types.MsgSendVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !IsBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(types.ErrNotBondedValidator, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	index := msg.Digest()
	send, isFound := k.GetSend(ctx, index)

	if isDuplicateSigner(msg.Creator, send.Signers) {
		return nil, sdkerrors.Wrap(types.ErrDuplicateMsg, fmt.Sprintf("signer %s double signing!!", msg.Creator))
	}

	if isFound { // send exists; add creator to signers
		send.Signers = append(send.Signers, msg.Creator)
	} else {
		send = types.Send{
			Creator:             msg.Creator,
			Index:               index,
			Sender:              msg.Sender,
			SenderChain:         msg.SenderChain,
			Receiver:            msg.Receiver,
			ReceiverChain:       msg.ReceiverChain,
			ZetaBurnt:           msg.ZetaBurnt,
			Message:             msg.Message,
			InTxHash:            msg.InTxHash,
			InBlockHeight:       msg.InBlockHeight,
			GasLimit:            msg.GasLimit,
			FinalizedMetaHeight: 0,
			Signers:             []string{msg.Creator},
			Status:              types.SendStatus_PendingInbound,
			Nonce:               0,
			RecvHash:            "",
			IndexTxList:         -1,
			LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
			ZetaMint:            "",
			CoinType:            msg.CoinType,
		}
		k.EmitEventSendCreated(ctx, &send)
	}

	if hasSuperMajorityValidators(len(send.Signers), validators) {
		send.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()

		send.FinalizedMetaHeight = uint64(ctx.BlockHeader().Height)
		send.Status = types.SendStatus_PendingOutbound
		k.UpdateLastBlockHigh(ctx, msg)

		bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
		send.Broadcaster = uint64(bftTime.Nanosecond() % len(send.Signers))

		abort := false
		// validate receiver address & chain; abort if failed
		recvChain, err := parseChainAndAddress(send.ReceiverChain, send.Receiver)
		if err != nil {
			send.StatusMessage = err.Error()
			send.Status = types.SendStatus_PendingRevert
			abort = true
		}

		var chain common.Chain // the chain for outbound
		if abort {
			chain, err = common.ParseChain(send.SenderChain)
			if err != nil {
				send.StatusMessage = fmt.Sprintf("cannot parse sender chain: %s", send.SenderChain)
				send.Status = types.SendStatus_Aborted
				goto EPILOGUE
			}
			send.Status = types.SendStatus_PendingRevert
		} else {
			chain = recvChain
		}

		if recvChain == common.ZETAChain { // if to zEVM, directly call EVM
			if send.CoinType == common.CoinType_Gas {
				foreignCoinList := k.fungibleKeeper.GetAllForeignCoins(ctx)
				found := false
				var gasCoin fungibletypes.ForeignCoins
				for _, coin := range foreignCoinList {
					if coin.CoinType == common.CoinType_Gas && coin.ForeignChain == send.SenderChain {
						found = true
						gasCoin = coin
						break
					}
				}
				if !found {
					send.StatusMessage = fmt.Sprintf("cannot get gas coin on chain %s", send.SenderChain)
					send.Status = types.SendStatus_Aborted
					goto EPILOGUE
				}
				to := ethcommon.HexToAddress(send.Receiver)
				amount, ok := big.NewInt(0).SetString(send.ZetaBurnt, 10)
				if !ok {
					send.StatusMessage = fmt.Sprintf("cannot parse zetaBurnt: %s", send.ZetaBurnt)
					send.Status = types.SendStatus_Aborted
					goto EPILOGUE
				}
				var tx *evmtypes.MsgEthereumTxResponse
				if len(send.Message) == 0 { // no message; transfer
					tx, err = k.fungibleKeeper.DepositZRC4(ctx, ethcommon.HexToAddress(gasCoin.ZRC4ContractAddress), to, amount)
					if err != nil {
						send.StatusMessage = fmt.Sprintf("cannot deposit zetaMint: %s", err.Error())
						send.Status = types.SendStatus_Aborted
						goto EPILOGUE
					}
				} else { // non-empty message = [contractaddress, calldata]
					contract, data, err := parseContractAndData(send.Message)
					tx, err = k.fungibleKeeper.DepositZRC4AndCallContract(ctx, ethcommon.HexToAddress(gasCoin.ZRC4ContractAddress), amount, contract, data)
					if err != nil {
						send.StatusMessage = fmt.Sprintf("cannot DepositZRC4AndCallContract zetaMint: %s", err.Error())
						send.Status = types.SendStatus_PendingRevert
						chain, err = common.ParseChain(send.SenderChain)
						if err != nil {
							send.StatusMessage = fmt.Sprintf("cannot parse sender chain: %s", send.SenderChain)
							send.Status = types.SendStatus_Aborted
							goto EPILOGUE
						}
						k.updateSend(ctx, chain.String(), &send)
						goto EPILOGUE
					}
					if !tx.Failed() {
						logs := evmtypes.LogsToEthereum(tx.Logs)
						k.ProcessWithdrawalEvent(ctx, logs, contract)
						ctx.EventManager().EmitEvent(
							sdk.NewEvent(sdk.EventTypeMessage,
								sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
								sdk.NewAttribute("action", "depositZRC4AndCallContract"),
								sdk.NewAttribute("contract", contract.String()),
								sdk.NewAttribute("data", hex.EncodeToString(data)),
							),
						)
					}
				}
				fmt.Printf("=======  tx: %s\n", tx.Hash)
				fmt.Printf("vmerror: %s\n", tx.VmError)
				fmt.Printf("=======  tx: %s\n", tx.Hash)

				send.OutTxHash = tx.Hash
				if tx.Failed() {
					send.StatusMessage = fmt.Sprintf("deposit zetaMint failed: %s, error %s", tx.Hash, tx.VmError)
					send.Status = types.SendStatus_Aborted
					goto EPILOGUE
				}
				send.Status = types.SendStatus_OutboundMined
			} else if send.CoinType == common.CoinType_Zeta {
				toBytes := ethcommon.HexToAddress(send.Receiver).Bytes()
				to := sdk.AccAddress(toBytes)
				amount, ok := big.NewInt(0).SetString(send.ZetaBurnt, 10)
				if !ok {
					send.StatusMessage = fmt.Sprintf("cannot parse zetaBurnt: %s", send.ZetaBurnt)
					send.Status = types.SendStatus_Aborted
					goto EPILOGUE
				}
				err := k.fungibleKeeper.MintZetaToEVMAccount(ctx, to, amount)
				if err != nil {
					send.StatusMessage = fmt.Sprintf("cannot MintZetaToEVMAccount: %s", err.Error())
					send.Status = types.SendStatus_Aborted
					goto EPILOGUE
				}
				send.Status = types.SendStatus_OutboundMined
			}

		} else {
			k.updateSend(ctx, chain.String(), &send)
			k.EmitEventSendFinalized(ctx, &send)
		}
	}

EPILOGUE:
	k.SetSend(ctx, send)
	return &types.MsgSendVoterResponse{}, nil
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

// updates gas price, gas fee, zeta to mint, and nonce
// returns ok?
func (k msgServer) updateSend(ctx sdk.Context, chain string, send *types.Send) bool {
	gasPrice, isFound := k.GetGasPrice(ctx, chain)
	if !isFound {
		send.StatusMessage = fmt.Sprintf("no gas price found: chain %s", send.ReceiverChain)
		send.Status = types.SendStatus_Aborted
		return false
	}
	mi := gasPrice.MedianIndex
	medianPrice := gasPrice.Prices[mi]
	send.GasPrice = fmt.Sprintf("%d", medianPrice)
	price, ok := big.NewInt(0).SetString(send.GasPrice, 10)
	if !ok {
		send.StatusMessage = fmt.Sprintf("GasPrice cannot parse")
		send.Status = types.SendStatus_Aborted
		return false
	}
	gasLimit := big.NewInt(0).SetUint64(send.GasLimit)

	gasFeeInZeta, abort := k.computeFeeInZeta(ctx, price, gasLimit, chain, send)
	if abort {
		send.Status = types.SendStatus_Aborted
		return false
	}
	zetaBurntInt, ok := big.NewInt(0).SetString(send.ZetaBurnt, 0)
	if !ok {
		send.StatusMessage = fmt.Sprintf("ZetaBurnt cannot parse")
		send.Status = types.SendStatus_Aborted
		return false
	}
	if gasFeeInZeta.Cmp(zetaBurntInt) > 0 {
		send.StatusMessage = fmt.Sprintf("feeInZeta(%d) more than mBurnt (%d)", gasFeeInZeta, zetaBurntInt)
		send.Status = types.SendStatus_Aborted
		return false
	}
	send.ZetaMint = fmt.Sprintf("%d", big.NewInt(0).Sub(zetaBurntInt, gasFeeInZeta))

	nonce, found := k.GetChainNonces(ctx, chain)
	if !found {
		send.StatusMessage = fmt.Sprintf("cannot find receiver chain nonce: %s", chain)
		send.Status = types.SendStatus_Aborted
		return false
	}

	send.Nonce = nonce.Nonce
	nonce.Nonce++
	k.SetChainNonces(ctx, nonce)
	return true
}

// returns (chain,error)
// chain: the receiverChain if ok
func parseChainAndAddress(chain string, addr string) (common.Chain, error) {
	recvChain, err := common.ParseChain(chain)
	if err != nil {
		return recvChain, fmt.Errorf("cannot parse receiver chain %s", chain)
	}
	_, err = common.NewAddress(addr, recvChain)
	if err != nil {
		return recvChain, fmt.Errorf("cannot parse receiver addr %s", addr)
	}
	return recvChain, nil
}

func (k msgServer) UpdateLastBlockHigh(ctx sdk.Context, msg *types.MsgSendVoter) {
	lastblock, isFound := k.GetLastBlockHeight(ctx, msg.SenderChain)
	if !isFound {
		lastblock = types.LastBlockHeight{
			Creator:           msg.Creator,
			Index:             msg.SenderChain,
			Chain:             msg.SenderChain,
			LastSendHeight:    msg.InBlockHeight,
			LastReceiveHeight: 0,
		}
	} else {
		lastblock.LastSendHeight = msg.InBlockHeight
	}
	k.SetLastBlockHeight(ctx, lastblock)
}

func (k msgServer) EmitEventSendCreated(ctx sdk.Context, send *types.Send) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
			sdk.NewAttribute(types.SubTypeKey, types.InboundCreated),
			sdk.NewAttribute(types.SendHash, send.Index),
			sdk.NewAttribute(types.NewStatus, send.Status.String()),
		),
	)
}

func (k msgServer) EmitEventSendFinalized(ctx sdk.Context, send *types.Send) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
			sdk.NewAttribute(types.SubTypeKey, types.InboundFinalized),
			sdk.NewAttribute(types.SendHash, send.Index),
			sdk.NewAttribute(types.Sender, send.Sender),
			sdk.NewAttribute(types.SenderChain, send.SenderChain),
			sdk.NewAttribute(types.Receiver, send.Receiver),
			sdk.NewAttribute(types.ReceiverChain, send.ReceiverChain),
			sdk.NewAttribute(types.ZetaBurnt, send.ZetaBurnt),
			sdk.NewAttribute(types.ZetaMint, send.ZetaMint),
			sdk.NewAttribute(types.Message, send.Message),
			sdk.NewAttribute(types.InTxHash, send.InTxHash),
			sdk.NewAttribute(types.InBlockHeight, fmt.Sprintf("%d", send.InBlockHeight)),
			sdk.NewAttribute(types.NewStatus, send.Status.String()),
			sdk.NewAttribute(types.StatusMessage, send.StatusMessage),
		),
	)
}

// returns feeInZeta (uint uuzeta), and whether to abort zeta-tx
func (k msgServer) computeFeeInZeta(ctx sdk.Context, price *big.Int, gasLimit *big.Int, chain string, send *types.Send) (*big.Int, bool) {
	abort := false
	rate, isFound := k.GetZetaConversionRate(ctx, chain)
	if !isFound {
		send.StatusMessage = fmt.Sprintf("Zeta conversion rate not found")
		abort = true
	}
	exchangeRateInt, ok := big.NewInt(0).SetString(rate.ZetaConversionRates[rate.MedianIndex], 0)
	if !ok {
		send.StatusMessage = fmt.Sprintf("median exchange rate %s cannot parse into float", rate.ZetaConversionRates[rate.MedianIndex])
		abort = true
	}

	// price*gasLimit*exchangeRate/1e18
	gasFeeInZeta := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(0).Mul(price, gasLimit), exchangeRateInt), OneEighteen)
	// add protocol flat fee: 1 ZETA
	gasFeeInZeta.Add(gasFeeInZeta, OneEighteen)
	return gasFeeInZeta, abort
}
