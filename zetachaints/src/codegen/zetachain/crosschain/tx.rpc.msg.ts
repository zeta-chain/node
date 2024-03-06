import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { MsgAddToOutTxTracker, MsgAddToOutTxTrackerResponse, MsgAddToInTxTracker, MsgAddToInTxTrackerResponse, MsgRemoveFromOutTxTracker, MsgRemoveFromOutTxTrackerResponse, MsgGasPriceVoter, MsgGasPriceVoterResponse, MsgVoteOnObservedOutboundTx, MsgVoteOnObservedOutboundTxResponse, MsgVoteOnObservedInboundTx, MsgVoteOnObservedInboundTxResponse, MsgWhitelistERC20, MsgWhitelistERC20Response, MsgUpdateTssAddress, MsgUpdateTssAddressResponse, MsgMigrateTssFunds, MsgMigrateTssFundsResponse, MsgCreateTSSVoter, MsgCreateTSSVoterResponse, MsgAbortStuckCCTX, MsgAbortStuckCCTXResponse, MsgRefundAbortedCCTX, MsgRefundAbortedCCTXResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
  addToOutTxTracker(request: MsgAddToOutTxTracker): Promise<MsgAddToOutTxTrackerResponse>;
  addToInTxTracker(request: MsgAddToInTxTracker): Promise<MsgAddToInTxTrackerResponse>;
  removeFromOutTxTracker(request: MsgRemoveFromOutTxTracker): Promise<MsgRemoveFromOutTxTrackerResponse>;
  gasPriceVoter(request: MsgGasPriceVoter): Promise<MsgGasPriceVoterResponse>;
  voteOnObservedOutboundTx(request: MsgVoteOnObservedOutboundTx): Promise<MsgVoteOnObservedOutboundTxResponse>;
  voteOnObservedInboundTx(request: MsgVoteOnObservedInboundTx): Promise<MsgVoteOnObservedInboundTxResponse>;
  whitelistERC20(request: MsgWhitelistERC20): Promise<MsgWhitelistERC20Response>;
  updateTssAddress(request: MsgUpdateTssAddress): Promise<MsgUpdateTssAddressResponse>;
  migrateTssFunds(request: MsgMigrateTssFunds): Promise<MsgMigrateTssFundsResponse>;
  createTSSVoter(request: MsgCreateTSSVoter): Promise<MsgCreateTSSVoterResponse>;
  abortStuckCCTX(request: MsgAbortStuckCCTX): Promise<MsgAbortStuckCCTXResponse>;
  refundAbortedCCTX(request: MsgRefundAbortedCCTX): Promise<MsgRefundAbortedCCTXResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.addToOutTxTracker = this.addToOutTxTracker.bind(this);
    this.addToInTxTracker = this.addToInTxTracker.bind(this);
    this.removeFromOutTxTracker = this.removeFromOutTxTracker.bind(this);
    this.gasPriceVoter = this.gasPriceVoter.bind(this);
    this.voteOnObservedOutboundTx = this.voteOnObservedOutboundTx.bind(this);
    this.voteOnObservedInboundTx = this.voteOnObservedInboundTx.bind(this);
    this.whitelistERC20 = this.whitelistERC20.bind(this);
    this.updateTssAddress = this.updateTssAddress.bind(this);
    this.migrateTssFunds = this.migrateTssFunds.bind(this);
    this.createTSSVoter = this.createTSSVoter.bind(this);
    this.abortStuckCCTX = this.abortStuckCCTX.bind(this);
    this.refundAbortedCCTX = this.refundAbortedCCTX.bind(this);
  }
  addToOutTxTracker(request: MsgAddToOutTxTracker): Promise<MsgAddToOutTxTrackerResponse> {
    const data = MsgAddToOutTxTracker.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "AddToOutTxTracker", data);
    return promise.then(data => MsgAddToOutTxTrackerResponse.decode(new BinaryReader(data)));
  }
  addToInTxTracker(request: MsgAddToInTxTracker): Promise<MsgAddToInTxTrackerResponse> {
    const data = MsgAddToInTxTracker.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "AddToInTxTracker", data);
    return promise.then(data => MsgAddToInTxTrackerResponse.decode(new BinaryReader(data)));
  }
  removeFromOutTxTracker(request: MsgRemoveFromOutTxTracker): Promise<MsgRemoveFromOutTxTrackerResponse> {
    const data = MsgRemoveFromOutTxTracker.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "RemoveFromOutTxTracker", data);
    return promise.then(data => MsgRemoveFromOutTxTrackerResponse.decode(new BinaryReader(data)));
  }
  gasPriceVoter(request: MsgGasPriceVoter): Promise<MsgGasPriceVoterResponse> {
    const data = MsgGasPriceVoter.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "GasPriceVoter", data);
    return promise.then(data => MsgGasPriceVoterResponse.decode(new BinaryReader(data)));
  }
  voteOnObservedOutboundTx(request: MsgVoteOnObservedOutboundTx): Promise<MsgVoteOnObservedOutboundTxResponse> {
    const data = MsgVoteOnObservedOutboundTx.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "VoteOnObservedOutboundTx", data);
    return promise.then(data => MsgVoteOnObservedOutboundTxResponse.decode(new BinaryReader(data)));
  }
  voteOnObservedInboundTx(request: MsgVoteOnObservedInboundTx): Promise<MsgVoteOnObservedInboundTxResponse> {
    const data = MsgVoteOnObservedInboundTx.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "VoteOnObservedInboundTx", data);
    return promise.then(data => MsgVoteOnObservedInboundTxResponse.decode(new BinaryReader(data)));
  }
  whitelistERC20(request: MsgWhitelistERC20): Promise<MsgWhitelistERC20Response> {
    const data = MsgWhitelistERC20.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "WhitelistERC20", data);
    return promise.then(data => MsgWhitelistERC20Response.decode(new BinaryReader(data)));
  }
  updateTssAddress(request: MsgUpdateTssAddress): Promise<MsgUpdateTssAddressResponse> {
    const data = MsgUpdateTssAddress.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "UpdateTssAddress", data);
    return promise.then(data => MsgUpdateTssAddressResponse.decode(new BinaryReader(data)));
  }
  migrateTssFunds(request: MsgMigrateTssFunds): Promise<MsgMigrateTssFundsResponse> {
    const data = MsgMigrateTssFunds.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "MigrateTssFunds", data);
    return promise.then(data => MsgMigrateTssFundsResponse.decode(new BinaryReader(data)));
  }
  createTSSVoter(request: MsgCreateTSSVoter): Promise<MsgCreateTSSVoterResponse> {
    const data = MsgCreateTSSVoter.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "CreateTSSVoter", data);
    return promise.then(data => MsgCreateTSSVoterResponse.decode(new BinaryReader(data)));
  }
  abortStuckCCTX(request: MsgAbortStuckCCTX): Promise<MsgAbortStuckCCTXResponse> {
    const data = MsgAbortStuckCCTX.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "AbortStuckCCTX", data);
    return promise.then(data => MsgAbortStuckCCTXResponse.decode(new BinaryReader(data)));
  }
  refundAbortedCCTX(request: MsgRefundAbortedCCTX): Promise<MsgRefundAbortedCCTXResponse> {
    const data = MsgRefundAbortedCCTX.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Msg", "RefundAbortedCCTX", data);
    return promise.then(data => MsgRefundAbortedCCTXResponse.decode(new BinaryReader(data)));
  }
}