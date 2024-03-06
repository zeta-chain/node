import { Chain, ChainAmino, ChainSDKType } from "../common/common";
import { BinaryReader, BinaryWriter } from "../../binary";
export enum ObservationType {
  EmptyObserverType = 0,
  InBoundTx = 1,
  OutBoundTx = 2,
  TSSKeyGen = 3,
  TSSKeySign = 4,
  UNRECOGNIZED = -1,
}
export const ObservationTypeSDKType = ObservationType;
export const ObservationTypeAmino = ObservationType;
export function observationTypeFromJSON(object: any): ObservationType {
  switch (object) {
    case 0:
    case "EmptyObserverType":
      return ObservationType.EmptyObserverType;
    case 1:
    case "InBoundTx":
      return ObservationType.InBoundTx;
    case 2:
    case "OutBoundTx":
      return ObservationType.OutBoundTx;
    case 3:
    case "TSSKeyGen":
      return ObservationType.TSSKeyGen;
    case 4:
    case "TSSKeySign":
      return ObservationType.TSSKeySign;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ObservationType.UNRECOGNIZED;
  }
}
export function observationTypeToJSON(object: ObservationType): string {
  switch (object) {
    case ObservationType.EmptyObserverType:
      return "EmptyObserverType";
    case ObservationType.InBoundTx:
      return "InBoundTx";
    case ObservationType.OutBoundTx:
      return "OutBoundTx";
    case ObservationType.TSSKeyGen:
      return "TSSKeyGen";
    case ObservationType.TSSKeySign:
      return "TSSKeySign";
    case ObservationType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export enum ObserverUpdateReason {
  Undefined = 0,
  Tombstoned = 1,
  AdminUpdate = 2,
  UNRECOGNIZED = -1,
}
export const ObserverUpdateReasonSDKType = ObserverUpdateReason;
export const ObserverUpdateReasonAmino = ObserverUpdateReason;
export function observerUpdateReasonFromJSON(object: any): ObserverUpdateReason {
  switch (object) {
    case 0:
    case "Undefined":
      return ObserverUpdateReason.Undefined;
    case 1:
    case "Tombstoned":
      return ObserverUpdateReason.Tombstoned;
    case 2:
    case "AdminUpdate":
      return ObserverUpdateReason.AdminUpdate;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ObserverUpdateReason.UNRECOGNIZED;
  }
}
export function observerUpdateReasonToJSON(object: ObserverUpdateReason): string {
  switch (object) {
    case ObserverUpdateReason.Undefined:
      return "Undefined";
    case ObserverUpdateReason.Tombstoned:
      return "Tombstoned";
    case ObserverUpdateReason.AdminUpdate:
      return "AdminUpdate";
    case ObserverUpdateReason.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export interface ObserverMapper {
  index: string;
  observerChain?: Chain;
  observerList: string[];
}
export interface ObserverMapperProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.ObserverMapper";
  value: Uint8Array;
}
export interface ObserverMapperAmino {
  index?: string;
  observer_chain?: ChainAmino;
  observer_list?: string[];
}
export interface ObserverMapperAminoMsg {
  type: "/zetachain.zetacore.observer.ObserverMapper";
  value: ObserverMapperAmino;
}
export interface ObserverMapperSDKType {
  index: string;
  observer_chain?: ChainSDKType;
  observer_list: string[];
}
export interface ObserverSet {
  observerList: string[];
}
export interface ObserverSetProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.ObserverSet";
  value: Uint8Array;
}
export interface ObserverSetAmino {
  observer_list?: string[];
}
export interface ObserverSetAminoMsg {
  type: "/zetachain.zetacore.observer.ObserverSet";
  value: ObserverSetAmino;
}
export interface ObserverSetSDKType {
  observer_list: string[];
}
export interface LastObserverCount {
  count: bigint;
  lastChangeHeight: bigint;
}
export interface LastObserverCountProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.LastObserverCount";
  value: Uint8Array;
}
export interface LastObserverCountAmino {
  count?: string;
  last_change_height?: string;
}
export interface LastObserverCountAminoMsg {
  type: "/zetachain.zetacore.observer.LastObserverCount";
  value: LastObserverCountAmino;
}
export interface LastObserverCountSDKType {
  count: bigint;
  last_change_height: bigint;
}
function createBaseObserverMapper(): ObserverMapper {
  return {
    index: "",
    observerChain: undefined,
    observerList: []
  };
}
export const ObserverMapper = {
  typeUrl: "/zetachain.zetacore.observer.ObserverMapper",
  encode(message: ObserverMapper, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    if (message.observerChain !== undefined) {
      Chain.encode(message.observerChain, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.observerList) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ObserverMapper {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObserverMapper();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        case 2:
          message.observerChain = Chain.decode(reader, reader.uint32());
          break;
        case 4:
          message.observerList.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ObserverMapper>): ObserverMapper {
    const message = createBaseObserverMapper();
    message.index = object.index ?? "";
    message.observerChain = object.observerChain !== undefined && object.observerChain !== null ? Chain.fromPartial(object.observerChain) : undefined;
    message.observerList = object.observerList?.map(e => e) || [];
    return message;
  },
  fromAmino(object: ObserverMapperAmino): ObserverMapper {
    const message = createBaseObserverMapper();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    if (object.observer_chain !== undefined && object.observer_chain !== null) {
      message.observerChain = Chain.fromAmino(object.observer_chain);
    }
    message.observerList = object.observer_list?.map(e => e) || [];
    return message;
  },
  toAmino(message: ObserverMapper): ObserverMapperAmino {
    const obj: any = {};
    obj.index = message.index;
    obj.observer_chain = message.observerChain ? Chain.toAmino(message.observerChain) : undefined;
    if (message.observerList) {
      obj.observer_list = message.observerList.map(e => e);
    } else {
      obj.observer_list = [];
    }
    return obj;
  },
  fromAminoMsg(object: ObserverMapperAminoMsg): ObserverMapper {
    return ObserverMapper.fromAmino(object.value);
  },
  fromProtoMsg(message: ObserverMapperProtoMsg): ObserverMapper {
    return ObserverMapper.decode(message.value);
  },
  toProto(message: ObserverMapper): Uint8Array {
    return ObserverMapper.encode(message).finish();
  },
  toProtoMsg(message: ObserverMapper): ObserverMapperProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.ObserverMapper",
      value: ObserverMapper.encode(message).finish()
    };
  }
};
function createBaseObserverSet(): ObserverSet {
  return {
    observerList: []
  };
}
export const ObserverSet = {
  typeUrl: "/zetachain.zetacore.observer.ObserverSet",
  encode(message: ObserverSet, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.observerList) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ObserverSet {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObserverSet();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.observerList.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ObserverSet>): ObserverSet {
    const message = createBaseObserverSet();
    message.observerList = object.observerList?.map(e => e) || [];
    return message;
  },
  fromAmino(object: ObserverSetAmino): ObserverSet {
    const message = createBaseObserverSet();
    message.observerList = object.observer_list?.map(e => e) || [];
    return message;
  },
  toAmino(message: ObserverSet): ObserverSetAmino {
    const obj: any = {};
    if (message.observerList) {
      obj.observer_list = message.observerList.map(e => e);
    } else {
      obj.observer_list = [];
    }
    return obj;
  },
  fromAminoMsg(object: ObserverSetAminoMsg): ObserverSet {
    return ObserverSet.fromAmino(object.value);
  },
  fromProtoMsg(message: ObserverSetProtoMsg): ObserverSet {
    return ObserverSet.decode(message.value);
  },
  toProto(message: ObserverSet): Uint8Array {
    return ObserverSet.encode(message).finish();
  },
  toProtoMsg(message: ObserverSet): ObserverSetProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.ObserverSet",
      value: ObserverSet.encode(message).finish()
    };
  }
};
function createBaseLastObserverCount(): LastObserverCount {
  return {
    count: BigInt(0),
    lastChangeHeight: BigInt(0)
  };
}
export const LastObserverCount = {
  typeUrl: "/zetachain.zetacore.observer.LastObserverCount",
  encode(message: LastObserverCount, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.count !== BigInt(0)) {
      writer.uint32(8).uint64(message.count);
    }
    if (message.lastChangeHeight !== BigInt(0)) {
      writer.uint32(16).int64(message.lastChangeHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): LastObserverCount {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLastObserverCount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.count = reader.uint64();
          break;
        case 2:
          message.lastChangeHeight = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<LastObserverCount>): LastObserverCount {
    const message = createBaseLastObserverCount();
    message.count = object.count !== undefined && object.count !== null ? BigInt(object.count.toString()) : BigInt(0);
    message.lastChangeHeight = object.lastChangeHeight !== undefined && object.lastChangeHeight !== null ? BigInt(object.lastChangeHeight.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: LastObserverCountAmino): LastObserverCount {
    const message = createBaseLastObserverCount();
    if (object.count !== undefined && object.count !== null) {
      message.count = BigInt(object.count);
    }
    if (object.last_change_height !== undefined && object.last_change_height !== null) {
      message.lastChangeHeight = BigInt(object.last_change_height);
    }
    return message;
  },
  toAmino(message: LastObserverCount): LastObserverCountAmino {
    const obj: any = {};
    obj.count = message.count ? message.count.toString() : undefined;
    obj.last_change_height = message.lastChangeHeight ? message.lastChangeHeight.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: LastObserverCountAminoMsg): LastObserverCount {
    return LastObserverCount.fromAmino(object.value);
  },
  fromProtoMsg(message: LastObserverCountProtoMsg): LastObserverCount {
    return LastObserverCount.decode(message.value);
  },
  toProto(message: LastObserverCount): Uint8Array {
    return LastObserverCount.encode(message).finish();
  },
  toProtoMsg(message: LastObserverCount): LastObserverCountProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.LastObserverCount",
      value: LastObserverCount.encode(message).finish()
    };
  }
};