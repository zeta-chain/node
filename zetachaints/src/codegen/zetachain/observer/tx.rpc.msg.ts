import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { MsgAddObserver, MsgAddObserverResponse, MsgUpdateObserver, MsgUpdateObserverResponse, MsgUpdateChainParams, MsgUpdateChainParamsResponse, MsgRemoveChainParams, MsgRemoveChainParamsResponse, MsgAddBlameVote, MsgAddBlameVoteResponse, MsgUpdateCrosschainFlags, MsgUpdateCrosschainFlagsResponse, MsgUpdateKeygen, MsgUpdateKeygenResponse, MsgAddBlockHeader, MsgAddBlockHeaderResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
  addObserver(request: MsgAddObserver): Promise<MsgAddObserverResponse>;
  updateObserver(request: MsgUpdateObserver): Promise<MsgUpdateObserverResponse>;
  updateChainParams(request: MsgUpdateChainParams): Promise<MsgUpdateChainParamsResponse>;
  removeChainParams(request: MsgRemoveChainParams): Promise<MsgRemoveChainParamsResponse>;
  addBlameVote(request: MsgAddBlameVote): Promise<MsgAddBlameVoteResponse>;
  updateCrosschainFlags(request: MsgUpdateCrosschainFlags): Promise<MsgUpdateCrosschainFlagsResponse>;
  updateKeygen(request: MsgUpdateKeygen): Promise<MsgUpdateKeygenResponse>;
  addBlockHeader(request: MsgAddBlockHeader): Promise<MsgAddBlockHeaderResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.addObserver = this.addObserver.bind(this);
    this.updateObserver = this.updateObserver.bind(this);
    this.updateChainParams = this.updateChainParams.bind(this);
    this.removeChainParams = this.removeChainParams.bind(this);
    this.addBlameVote = this.addBlameVote.bind(this);
    this.updateCrosschainFlags = this.updateCrosschainFlags.bind(this);
    this.updateKeygen = this.updateKeygen.bind(this);
    this.addBlockHeader = this.addBlockHeader.bind(this);
  }
  addObserver(request: MsgAddObserver): Promise<MsgAddObserverResponse> {
    const data = MsgAddObserver.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Msg", "AddObserver", data);
    return promise.then(data => MsgAddObserverResponse.decode(new BinaryReader(data)));
  }
  updateObserver(request: MsgUpdateObserver): Promise<MsgUpdateObserverResponse> {
    const data = MsgUpdateObserver.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Msg", "UpdateObserver", data);
    return promise.then(data => MsgUpdateObserverResponse.decode(new BinaryReader(data)));
  }
  updateChainParams(request: MsgUpdateChainParams): Promise<MsgUpdateChainParamsResponse> {
    const data = MsgUpdateChainParams.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Msg", "UpdateChainParams", data);
    return promise.then(data => MsgUpdateChainParamsResponse.decode(new BinaryReader(data)));
  }
  removeChainParams(request: MsgRemoveChainParams): Promise<MsgRemoveChainParamsResponse> {
    const data = MsgRemoveChainParams.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Msg", "RemoveChainParams", data);
    return promise.then(data => MsgRemoveChainParamsResponse.decode(new BinaryReader(data)));
  }
  addBlameVote(request: MsgAddBlameVote): Promise<MsgAddBlameVoteResponse> {
    const data = MsgAddBlameVote.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Msg", "AddBlameVote", data);
    return promise.then(data => MsgAddBlameVoteResponse.decode(new BinaryReader(data)));
  }
  updateCrosschainFlags(request: MsgUpdateCrosschainFlags): Promise<MsgUpdateCrosschainFlagsResponse> {
    const data = MsgUpdateCrosschainFlags.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Msg", "UpdateCrosschainFlags", data);
    return promise.then(data => MsgUpdateCrosschainFlagsResponse.decode(new BinaryReader(data)));
  }
  updateKeygen(request: MsgUpdateKeygen): Promise<MsgUpdateKeygenResponse> {
    const data = MsgUpdateKeygen.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Msg", "UpdateKeygen", data);
    return promise.then(data => MsgUpdateKeygenResponse.decode(new BinaryReader(data)));
  }
  addBlockHeader(request: MsgAddBlockHeader): Promise<MsgAddBlockHeaderResponse> {
    const data = MsgAddBlockHeader.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Msg", "AddBlockHeader", data);
    return promise.then(data => MsgAddBlockHeaderResponse.decode(new BinaryReader(data)));
  }
}