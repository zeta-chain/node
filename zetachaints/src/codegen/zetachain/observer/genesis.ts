import { Ballot, BallotAmino, BallotSDKType } from "./ballot";
import { ObserverSet, ObserverSetAmino, ObserverSetSDKType, LastObserverCount, LastObserverCountAmino, LastObserverCountSDKType } from "./observer";
import { NodeAccount, NodeAccountAmino, NodeAccountSDKType } from "./node_account";
import { CrosschainFlags, CrosschainFlagsAmino, CrosschainFlagsSDKType } from "./crosschain_flags";
import { Params, ParamsAmino, ParamsSDKType, ChainParamsList, ChainParamsListAmino, ChainParamsListSDKType } from "./params";
import { Keygen, KeygenAmino, KeygenSDKType } from "./keygen";
import { TSS, TSSAmino, TSSSDKType } from "./tss";
import { TssFundMigratorInfo, TssFundMigratorInfoAmino, TssFundMigratorInfoSDKType } from "./tss_funds_migrator";
import { Blame, BlameAmino, BlameSDKType } from "./blame";
import { PendingNonces, PendingNoncesAmino, PendingNoncesSDKType } from "./pending_nonces";
import { ChainNonces, ChainNoncesAmino, ChainNoncesSDKType } from "./chain_nonces";
import { NonceToCctx, NonceToCctxAmino, NonceToCctxSDKType } from "./nonce_to_cctx";
import { BinaryReader, BinaryWriter } from "../../binary";
export interface GenesisState {
  ballots: Ballot[];
  observers: ObserverSet;
  nodeAccountList: NodeAccount[];
  crosschainFlags?: CrosschainFlags;
  params?: Params;
  keygen?: Keygen;
  lastObserverCount?: LastObserverCount;
  chainParamsList: ChainParamsList;
  tss?: TSS;
  tssHistory: TSS[];
  tssFundMigrators: TssFundMigratorInfo[];
  blameList: Blame[];
  pendingNonces: PendingNonces[];
  chainNonces: ChainNonces[];
  nonceToCctx: NonceToCctx[];
}
export interface GenesisStateProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.GenesisState";
  value: Uint8Array;
}
export interface GenesisStateAmino {
  ballots?: BallotAmino[];
  observers?: ObserverSetAmino;
  nodeAccountList?: NodeAccountAmino[];
  crosschain_flags?: CrosschainFlagsAmino;
  params?: ParamsAmino;
  keygen?: KeygenAmino;
  last_observer_count?: LastObserverCountAmino;
  chain_params_list?: ChainParamsListAmino;
  tss?: TSSAmino;
  tss_history?: TSSAmino[];
  tss_fund_migrators?: TssFundMigratorInfoAmino[];
  blame_list?: BlameAmino[];
  pending_nonces?: PendingNoncesAmino[];
  chain_nonces?: ChainNoncesAmino[];
  nonce_to_cctx?: NonceToCctxAmino[];
}
export interface GenesisStateAminoMsg {
  type: "/zetachain.zetacore.observer.GenesisState";
  value: GenesisStateAmino;
}
export interface GenesisStateSDKType {
  ballots: BallotSDKType[];
  observers: ObserverSetSDKType;
  nodeAccountList: NodeAccountSDKType[];
  crosschain_flags?: CrosschainFlagsSDKType;
  params?: ParamsSDKType;
  keygen?: KeygenSDKType;
  last_observer_count?: LastObserverCountSDKType;
  chain_params_list: ChainParamsListSDKType;
  tss?: TSSSDKType;
  tss_history: TSSSDKType[];
  tss_fund_migrators: TssFundMigratorInfoSDKType[];
  blame_list: BlameSDKType[];
  pending_nonces: PendingNoncesSDKType[];
  chain_nonces: ChainNoncesSDKType[];
  nonce_to_cctx: NonceToCctxSDKType[];
}
function createBaseGenesisState(): GenesisState {
  return {
    ballots: [],
    observers: ObserverSet.fromPartial({}),
    nodeAccountList: [],
    crosschainFlags: undefined,
    params: undefined,
    keygen: undefined,
    lastObserverCount: undefined,
    chainParamsList: ChainParamsList.fromPartial({}),
    tss: undefined,
    tssHistory: [],
    tssFundMigrators: [],
    blameList: [],
    pendingNonces: [],
    chainNonces: [],
    nonceToCctx: []
  };
}
export const GenesisState = {
  typeUrl: "/zetachain.zetacore.observer.GenesisState",
  encode(message: GenesisState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.ballots) {
      Ballot.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.observers !== undefined) {
      ObserverSet.encode(message.observers, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.nodeAccountList) {
      NodeAccount.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.crosschainFlags !== undefined) {
      CrosschainFlags.encode(message.crosschainFlags, writer.uint32(34).fork()).ldelim();
    }
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(42).fork()).ldelim();
    }
    if (message.keygen !== undefined) {
      Keygen.encode(message.keygen, writer.uint32(50).fork()).ldelim();
    }
    if (message.lastObserverCount !== undefined) {
      LastObserverCount.encode(message.lastObserverCount, writer.uint32(58).fork()).ldelim();
    }
    if (message.chainParamsList !== undefined) {
      ChainParamsList.encode(message.chainParamsList, writer.uint32(66).fork()).ldelim();
    }
    if (message.tss !== undefined) {
      TSS.encode(message.tss, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.tssHistory) {
      TSS.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.tssFundMigrators) {
      TssFundMigratorInfo.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    for (const v of message.blameList) {
      Blame.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    for (const v of message.pendingNonces) {
      PendingNonces.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    for (const v of message.chainNonces) {
      ChainNonces.encode(v!, writer.uint32(114).fork()).ldelim();
    }
    for (const v of message.nonceToCctx) {
      NonceToCctx.encode(v!, writer.uint32(122).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ballots.push(Ballot.decode(reader, reader.uint32()));
          break;
        case 2:
          message.observers = ObserverSet.decode(reader, reader.uint32());
          break;
        case 3:
          message.nodeAccountList.push(NodeAccount.decode(reader, reader.uint32()));
          break;
        case 4:
          message.crosschainFlags = CrosschainFlags.decode(reader, reader.uint32());
          break;
        case 5:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 6:
          message.keygen = Keygen.decode(reader, reader.uint32());
          break;
        case 7:
          message.lastObserverCount = LastObserverCount.decode(reader, reader.uint32());
          break;
        case 8:
          message.chainParamsList = ChainParamsList.decode(reader, reader.uint32());
          break;
        case 9:
          message.tss = TSS.decode(reader, reader.uint32());
          break;
        case 10:
          message.tssHistory.push(TSS.decode(reader, reader.uint32()));
          break;
        case 11:
          message.tssFundMigrators.push(TssFundMigratorInfo.decode(reader, reader.uint32()));
          break;
        case 12:
          message.blameList.push(Blame.decode(reader, reader.uint32()));
          break;
        case 13:
          message.pendingNonces.push(PendingNonces.decode(reader, reader.uint32()));
          break;
        case 14:
          message.chainNonces.push(ChainNonces.decode(reader, reader.uint32()));
          break;
        case 15:
          message.nonceToCctx.push(NonceToCctx.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.ballots = object.ballots?.map(e => Ballot.fromPartial(e)) || [];
    message.observers = object.observers !== undefined && object.observers !== null ? ObserverSet.fromPartial(object.observers) : undefined;
    message.nodeAccountList = object.nodeAccountList?.map(e => NodeAccount.fromPartial(e)) || [];
    message.crosschainFlags = object.crosschainFlags !== undefined && object.crosschainFlags !== null ? CrosschainFlags.fromPartial(object.crosschainFlags) : undefined;
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    message.keygen = object.keygen !== undefined && object.keygen !== null ? Keygen.fromPartial(object.keygen) : undefined;
    message.lastObserverCount = object.lastObserverCount !== undefined && object.lastObserverCount !== null ? LastObserverCount.fromPartial(object.lastObserverCount) : undefined;
    message.chainParamsList = object.chainParamsList !== undefined && object.chainParamsList !== null ? ChainParamsList.fromPartial(object.chainParamsList) : undefined;
    message.tss = object.tss !== undefined && object.tss !== null ? TSS.fromPartial(object.tss) : undefined;
    message.tssHistory = object.tssHistory?.map(e => TSS.fromPartial(e)) || [];
    message.tssFundMigrators = object.tssFundMigrators?.map(e => TssFundMigratorInfo.fromPartial(e)) || [];
    message.blameList = object.blameList?.map(e => Blame.fromPartial(e)) || [];
    message.pendingNonces = object.pendingNonces?.map(e => PendingNonces.fromPartial(e)) || [];
    message.chainNonces = object.chainNonces?.map(e => ChainNonces.fromPartial(e)) || [];
    message.nonceToCctx = object.nonceToCctx?.map(e => NonceToCctx.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    const message = createBaseGenesisState();
    message.ballots = object.ballots?.map(e => Ballot.fromAmino(e)) || [];
    if (object.observers !== undefined && object.observers !== null) {
      message.observers = ObserverSet.fromAmino(object.observers);
    }
    message.nodeAccountList = object.nodeAccountList?.map(e => NodeAccount.fromAmino(e)) || [];
    if (object.crosschain_flags !== undefined && object.crosschain_flags !== null) {
      message.crosschainFlags = CrosschainFlags.fromAmino(object.crosschain_flags);
    }
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromAmino(object.params);
    }
    if (object.keygen !== undefined && object.keygen !== null) {
      message.keygen = Keygen.fromAmino(object.keygen);
    }
    if (object.last_observer_count !== undefined && object.last_observer_count !== null) {
      message.lastObserverCount = LastObserverCount.fromAmino(object.last_observer_count);
    }
    if (object.chain_params_list !== undefined && object.chain_params_list !== null) {
      message.chainParamsList = ChainParamsList.fromAmino(object.chain_params_list);
    }
    if (object.tss !== undefined && object.tss !== null) {
      message.tss = TSS.fromAmino(object.tss);
    }
    message.tssHistory = object.tss_history?.map(e => TSS.fromAmino(e)) || [];
    message.tssFundMigrators = object.tss_fund_migrators?.map(e => TssFundMigratorInfo.fromAmino(e)) || [];
    message.blameList = object.blame_list?.map(e => Blame.fromAmino(e)) || [];
    message.pendingNonces = object.pending_nonces?.map(e => PendingNonces.fromAmino(e)) || [];
    message.chainNonces = object.chain_nonces?.map(e => ChainNonces.fromAmino(e)) || [];
    message.nonceToCctx = object.nonce_to_cctx?.map(e => NonceToCctx.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    if (message.ballots) {
      obj.ballots = message.ballots.map(e => e ? Ballot.toAmino(e) : undefined);
    } else {
      obj.ballots = [];
    }
    obj.observers = message.observers ? ObserverSet.toAmino(message.observers) : undefined;
    if (message.nodeAccountList) {
      obj.nodeAccountList = message.nodeAccountList.map(e => e ? NodeAccount.toAmino(e) : undefined);
    } else {
      obj.nodeAccountList = [];
    }
    obj.crosschain_flags = message.crosschainFlags ? CrosschainFlags.toAmino(message.crosschainFlags) : undefined;
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    obj.keygen = message.keygen ? Keygen.toAmino(message.keygen) : undefined;
    obj.last_observer_count = message.lastObserverCount ? LastObserverCount.toAmino(message.lastObserverCount) : undefined;
    obj.chain_params_list = message.chainParamsList ? ChainParamsList.toAmino(message.chainParamsList) : undefined;
    obj.tss = message.tss ? TSS.toAmino(message.tss) : undefined;
    if (message.tssHistory) {
      obj.tss_history = message.tssHistory.map(e => e ? TSS.toAmino(e) : undefined);
    } else {
      obj.tss_history = [];
    }
    if (message.tssFundMigrators) {
      obj.tss_fund_migrators = message.tssFundMigrators.map(e => e ? TssFundMigratorInfo.toAmino(e) : undefined);
    } else {
      obj.tss_fund_migrators = [];
    }
    if (message.blameList) {
      obj.blame_list = message.blameList.map(e => e ? Blame.toAmino(e) : undefined);
    } else {
      obj.blame_list = [];
    }
    if (message.pendingNonces) {
      obj.pending_nonces = message.pendingNonces.map(e => e ? PendingNonces.toAmino(e) : undefined);
    } else {
      obj.pending_nonces = [];
    }
    if (message.chainNonces) {
      obj.chain_nonces = message.chainNonces.map(e => e ? ChainNonces.toAmino(e) : undefined);
    } else {
      obj.chain_nonces = [];
    }
    if (message.nonceToCctx) {
      obj.nonce_to_cctx = message.nonceToCctx.map(e => e ? NonceToCctx.toAmino(e) : undefined);
    } else {
      obj.nonce_to_cctx = [];
    }
    return obj;
  },
  fromAminoMsg(object: GenesisStateAminoMsg): GenesisState {
    return GenesisState.fromAmino(object.value);
  },
  fromProtoMsg(message: GenesisStateProtoMsg): GenesisState {
    return GenesisState.decode(message.value);
  },
  toProto(message: GenesisState): Uint8Array {
    return GenesisState.encode(message).finish();
  },
  toProtoMsg(message: GenesisState): GenesisStateProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};