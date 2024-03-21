import { ObserverUpdateReason, observerUpdateReasonFromJSON } from "./observer";
import { ChainParams, ChainParamsAmino, ChainParamsSDKType } from "./params";
import { Blame, BlameAmino, BlameSDKType } from "./blame";
import { GasPriceIncreaseFlags, GasPriceIncreaseFlagsAmino, GasPriceIncreaseFlagsSDKType, BlockHeaderVerificationFlags, BlockHeaderVerificationFlagsAmino, BlockHeaderVerificationFlagsSDKType } from "./crosschain_flags";
import { HeaderData, HeaderDataAmino, HeaderDataSDKType } from "../common/common";
import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
export interface MsgUpdateObserver {
  creator: string;
  oldObserverAddress: string;
  newObserverAddress: string;
  updateReason: ObserverUpdateReason;
}
export interface MsgUpdateObserverProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateObserver";
  value: Uint8Array;
}
export interface MsgUpdateObserverAmino {
  creator?: string;
  old_observer_address?: string;
  new_observer_address?: string;
  update_reason?: ObserverUpdateReason;
}
export interface MsgUpdateObserverAminoMsg {
  type: "/zetachain.zetacore.observer.MsgUpdateObserver";
  value: MsgUpdateObserverAmino;
}
export interface MsgUpdateObserverSDKType {
  creator: string;
  old_observer_address: string;
  new_observer_address: string;
  update_reason: ObserverUpdateReason;
}
export interface MsgUpdateObserverResponse {}
export interface MsgUpdateObserverResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateObserverResponse";
  value: Uint8Array;
}
export interface MsgUpdateObserverResponseAmino {}
export interface MsgUpdateObserverResponseAminoMsg {
  type: "/zetachain.zetacore.observer.MsgUpdateObserverResponse";
  value: MsgUpdateObserverResponseAmino;
}
export interface MsgUpdateObserverResponseSDKType {}
export interface MsgAddBlockHeader {
  creator: string;
  chainId: bigint;
  blockHash: Uint8Array;
  height: bigint;
  header: HeaderData;
}
export interface MsgAddBlockHeaderProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgAddBlockHeader";
  value: Uint8Array;
}
export interface MsgAddBlockHeaderAmino {
  creator?: string;
  chain_id?: string;
  block_hash?: string;
  height?: string;
  header?: HeaderDataAmino;
}
export interface MsgAddBlockHeaderAminoMsg {
  type: "/zetachain.zetacore.observer.MsgAddBlockHeader";
  value: MsgAddBlockHeaderAmino;
}
export interface MsgAddBlockHeaderSDKType {
  creator: string;
  chain_id: bigint;
  block_hash: Uint8Array;
  height: bigint;
  header: HeaderDataSDKType;
}
export interface MsgAddBlockHeaderResponse {}
export interface MsgAddBlockHeaderResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgAddBlockHeaderResponse";
  value: Uint8Array;
}
export interface MsgAddBlockHeaderResponseAmino {}
export interface MsgAddBlockHeaderResponseAminoMsg {
  type: "/zetachain.zetacore.observer.MsgAddBlockHeaderResponse";
  value: MsgAddBlockHeaderResponseAmino;
}
export interface MsgAddBlockHeaderResponseSDKType {}
export interface MsgUpdateChainParams {
  creator: string;
  chainParams?: ChainParams;
}
export interface MsgUpdateChainParamsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateChainParams";
  value: Uint8Array;
}
export interface MsgUpdateChainParamsAmino {
  creator?: string;
  chainParams?: ChainParamsAmino;
}
export interface MsgUpdateChainParamsAminoMsg {
  type: "/zetachain.zetacore.observer.MsgUpdateChainParams";
  value: MsgUpdateChainParamsAmino;
}
export interface MsgUpdateChainParamsSDKType {
  creator: string;
  chainParams?: ChainParamsSDKType;
}
export interface MsgUpdateChainParamsResponse {}
export interface MsgUpdateChainParamsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateChainParamsResponse";
  value: Uint8Array;
}
export interface MsgUpdateChainParamsResponseAmino {}
export interface MsgUpdateChainParamsResponseAminoMsg {
  type: "/zetachain.zetacore.observer.MsgUpdateChainParamsResponse";
  value: MsgUpdateChainParamsResponseAmino;
}
export interface MsgUpdateChainParamsResponseSDKType {}
export interface MsgRemoveChainParams {
  creator: string;
  chainId: bigint;
}
export interface MsgRemoveChainParamsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgRemoveChainParams";
  value: Uint8Array;
}
export interface MsgRemoveChainParamsAmino {
  creator?: string;
  chain_id?: string;
}
export interface MsgRemoveChainParamsAminoMsg {
  type: "/zetachain.zetacore.observer.MsgRemoveChainParams";
  value: MsgRemoveChainParamsAmino;
}
export interface MsgRemoveChainParamsSDKType {
  creator: string;
  chain_id: bigint;
}
export interface MsgRemoveChainParamsResponse {}
export interface MsgRemoveChainParamsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgRemoveChainParamsResponse";
  value: Uint8Array;
}
export interface MsgRemoveChainParamsResponseAmino {}
export interface MsgRemoveChainParamsResponseAminoMsg {
  type: "/zetachain.zetacore.observer.MsgRemoveChainParamsResponse";
  value: MsgRemoveChainParamsResponseAmino;
}
export interface MsgRemoveChainParamsResponseSDKType {}
export interface MsgAddObserver {
  creator: string;
  observerAddress: string;
  zetaclientGranteePubkey: string;
  addNodeAccountOnly: boolean;
}
export interface MsgAddObserverProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgAddObserver";
  value: Uint8Array;
}
export interface MsgAddObserverAmino {
  creator?: string;
  observer_address?: string;
  zetaclient_grantee_pubkey?: string;
  add_node_account_only?: boolean;
}
export interface MsgAddObserverAminoMsg {
  type: "/zetachain.zetacore.observer.MsgAddObserver";
  value: MsgAddObserverAmino;
}
export interface MsgAddObserverSDKType {
  creator: string;
  observer_address: string;
  zetaclient_grantee_pubkey: string;
  add_node_account_only: boolean;
}
export interface MsgAddObserverResponse {}
export interface MsgAddObserverResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgAddObserverResponse";
  value: Uint8Array;
}
export interface MsgAddObserverResponseAmino {}
export interface MsgAddObserverResponseAminoMsg {
  type: "/zetachain.zetacore.observer.MsgAddObserverResponse";
  value: MsgAddObserverResponseAmino;
}
export interface MsgAddObserverResponseSDKType {}
export interface MsgAddBlameVote {
  creator: string;
  chainId: bigint;
  blameInfo: Blame;
}
export interface MsgAddBlameVoteProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgAddBlameVote";
  value: Uint8Array;
}
export interface MsgAddBlameVoteAmino {
  creator?: string;
  chain_id?: string;
  blame_info?: BlameAmino;
}
export interface MsgAddBlameVoteAminoMsg {
  type: "/zetachain.zetacore.observer.MsgAddBlameVote";
  value: MsgAddBlameVoteAmino;
}
export interface MsgAddBlameVoteSDKType {
  creator: string;
  chain_id: bigint;
  blame_info: BlameSDKType;
}
export interface MsgAddBlameVoteResponse {}
export interface MsgAddBlameVoteResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgAddBlameVoteResponse";
  value: Uint8Array;
}
export interface MsgAddBlameVoteResponseAmino {}
export interface MsgAddBlameVoteResponseAminoMsg {
  type: "/zetachain.zetacore.observer.MsgAddBlameVoteResponse";
  value: MsgAddBlameVoteResponseAmino;
}
export interface MsgAddBlameVoteResponseSDKType {}
export interface MsgUpdateCrosschainFlags {
  creator: string;
  isInboundEnabled: boolean;
  isOutboundEnabled: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlags;
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlags;
}
export interface MsgUpdateCrosschainFlagsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlags";
  value: Uint8Array;
}
export interface MsgUpdateCrosschainFlagsAmino {
  creator?: string;
  isInboundEnabled?: boolean;
  isOutboundEnabled?: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlagsAmino;
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlagsAmino;
}
export interface MsgUpdateCrosschainFlagsAminoMsg {
  type: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlags";
  value: MsgUpdateCrosschainFlagsAmino;
}
export interface MsgUpdateCrosschainFlagsSDKType {
  creator: string;
  isInboundEnabled: boolean;
  isOutboundEnabled: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlagsSDKType;
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlagsSDKType;
}
export interface MsgUpdateCrosschainFlagsResponse {}
export interface MsgUpdateCrosschainFlagsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlagsResponse";
  value: Uint8Array;
}
export interface MsgUpdateCrosschainFlagsResponseAmino {}
export interface MsgUpdateCrosschainFlagsResponseAminoMsg {
  type: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlagsResponse";
  value: MsgUpdateCrosschainFlagsResponseAmino;
}
export interface MsgUpdateCrosschainFlagsResponseSDKType {}
export interface MsgUpdateKeygen {
  creator: string;
  block: bigint;
}
export interface MsgUpdateKeygenProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateKeygen";
  value: Uint8Array;
}
export interface MsgUpdateKeygenAmino {
  creator?: string;
  block?: string;
}
export interface MsgUpdateKeygenAminoMsg {
  type: "/zetachain.zetacore.observer.MsgUpdateKeygen";
  value: MsgUpdateKeygenAmino;
}
export interface MsgUpdateKeygenSDKType {
  creator: string;
  block: bigint;
}
export interface MsgUpdateKeygenResponse {}
export interface MsgUpdateKeygenResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateKeygenResponse";
  value: Uint8Array;
}
export interface MsgUpdateKeygenResponseAmino {}
export interface MsgUpdateKeygenResponseAminoMsg {
  type: "/zetachain.zetacore.observer.MsgUpdateKeygenResponse";
  value: MsgUpdateKeygenResponseAmino;
}
export interface MsgUpdateKeygenResponseSDKType {}
function createBaseMsgUpdateObserver(): MsgUpdateObserver {
  return {
    creator: "",
    oldObserverAddress: "",
    newObserverAddress: "",
    updateReason: 0
  };
}
export const MsgUpdateObserver = {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateObserver",
  encode(message: MsgUpdateObserver, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.oldObserverAddress !== "") {
      writer.uint32(18).string(message.oldObserverAddress);
    }
    if (message.newObserverAddress !== "") {
      writer.uint32(26).string(message.newObserverAddress);
    }
    if (message.updateReason !== 0) {
      writer.uint32(32).int32(message.updateReason);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateObserver {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateObserver();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.oldObserverAddress = reader.string();
          break;
        case 3:
          message.newObserverAddress = reader.string();
          break;
        case 4:
          message.updateReason = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateObserver>): MsgUpdateObserver {
    const message = createBaseMsgUpdateObserver();
    message.creator = object.creator ?? "";
    message.oldObserverAddress = object.oldObserverAddress ?? "";
    message.newObserverAddress = object.newObserverAddress ?? "";
    message.updateReason = object.updateReason ?? 0;
    return message;
  },
  fromAmino(object: MsgUpdateObserverAmino): MsgUpdateObserver {
    const message = createBaseMsgUpdateObserver();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.old_observer_address !== undefined && object.old_observer_address !== null) {
      message.oldObserverAddress = object.old_observer_address;
    }
    if (object.new_observer_address !== undefined && object.new_observer_address !== null) {
      message.newObserverAddress = object.new_observer_address;
    }
    if (object.update_reason !== undefined && object.update_reason !== null) {
      message.updateReason = observerUpdateReasonFromJSON(object.update_reason);
    }
    return message;
  },
  toAmino(message: MsgUpdateObserver): MsgUpdateObserverAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.old_observer_address = message.oldObserverAddress;
    obj.new_observer_address = message.newObserverAddress;
    obj.update_reason = message.updateReason;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateObserverAminoMsg): MsgUpdateObserver {
    return MsgUpdateObserver.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateObserverProtoMsg): MsgUpdateObserver {
    return MsgUpdateObserver.decode(message.value);
  },
  toProto(message: MsgUpdateObserver): Uint8Array {
    return MsgUpdateObserver.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateObserver): MsgUpdateObserverProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgUpdateObserver",
      value: MsgUpdateObserver.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateObserverResponse(): MsgUpdateObserverResponse {
  return {};
}
export const MsgUpdateObserverResponse = {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateObserverResponse",
  encode(_: MsgUpdateObserverResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateObserverResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateObserverResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<MsgUpdateObserverResponse>): MsgUpdateObserverResponse {
    const message = createBaseMsgUpdateObserverResponse();
    return message;
  },
  fromAmino(_: MsgUpdateObserverResponseAmino): MsgUpdateObserverResponse {
    const message = createBaseMsgUpdateObserverResponse();
    return message;
  },
  toAmino(_: MsgUpdateObserverResponse): MsgUpdateObserverResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateObserverResponseAminoMsg): MsgUpdateObserverResponse {
    return MsgUpdateObserverResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateObserverResponseProtoMsg): MsgUpdateObserverResponse {
    return MsgUpdateObserverResponse.decode(message.value);
  },
  toProto(message: MsgUpdateObserverResponse): Uint8Array {
    return MsgUpdateObserverResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateObserverResponse): MsgUpdateObserverResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgUpdateObserverResponse",
      value: MsgUpdateObserverResponse.encode(message).finish()
    };
  }
};
function createBaseMsgAddBlockHeader(): MsgAddBlockHeader {
  return {
    creator: "",
    chainId: BigInt(0),
    blockHash: new Uint8Array(),
    height: BigInt(0),
    header: HeaderData.fromPartial({})
  };
}
export const MsgAddBlockHeader = {
  typeUrl: "/zetachain.zetacore.observer.MsgAddBlockHeader",
  encode(message: MsgAddBlockHeader, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.blockHash.length !== 0) {
      writer.uint32(26).bytes(message.blockHash);
    }
    if (message.height !== BigInt(0)) {
      writer.uint32(32).int64(message.height);
    }
    if (message.header !== undefined) {
      HeaderData.encode(message.header, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddBlockHeader {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddBlockHeader();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.blockHash = reader.bytes();
          break;
        case 4:
          message.height = reader.int64();
          break;
        case 5:
          message.header = HeaderData.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgAddBlockHeader>): MsgAddBlockHeader {
    const message = createBaseMsgAddBlockHeader();
    message.creator = object.creator ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.blockHash = object.blockHash ?? new Uint8Array();
    message.height = object.height !== undefined && object.height !== null ? BigInt(object.height.toString()) : BigInt(0);
    message.header = object.header !== undefined && object.header !== null ? HeaderData.fromPartial(object.header) : undefined;
    return message;
  },
  fromAmino(object: MsgAddBlockHeaderAmino): MsgAddBlockHeader {
    const message = createBaseMsgAddBlockHeader();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.block_hash !== undefined && object.block_hash !== null) {
      message.blockHash = bytesFromBase64(object.block_hash);
    }
    if (object.height !== undefined && object.height !== null) {
      message.height = BigInt(object.height);
    }
    if (object.header !== undefined && object.header !== null) {
      message.header = HeaderData.fromAmino(object.header);
    }
    return message;
  },
  toAmino(message: MsgAddBlockHeader): MsgAddBlockHeaderAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.block_hash = message.blockHash ? base64FromBytes(message.blockHash) : undefined;
    obj.height = message.height ? message.height.toString() : undefined;
    obj.header = message.header ? HeaderData.toAmino(message.header) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgAddBlockHeaderAminoMsg): MsgAddBlockHeader {
    return MsgAddBlockHeader.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddBlockHeaderProtoMsg): MsgAddBlockHeader {
    return MsgAddBlockHeader.decode(message.value);
  },
  toProto(message: MsgAddBlockHeader): Uint8Array {
    return MsgAddBlockHeader.encode(message).finish();
  },
  toProtoMsg(message: MsgAddBlockHeader): MsgAddBlockHeaderProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgAddBlockHeader",
      value: MsgAddBlockHeader.encode(message).finish()
    };
  }
};
function createBaseMsgAddBlockHeaderResponse(): MsgAddBlockHeaderResponse {
  return {};
}
export const MsgAddBlockHeaderResponse = {
  typeUrl: "/zetachain.zetacore.observer.MsgAddBlockHeaderResponse",
  encode(_: MsgAddBlockHeaderResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddBlockHeaderResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddBlockHeaderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<MsgAddBlockHeaderResponse>): MsgAddBlockHeaderResponse {
    const message = createBaseMsgAddBlockHeaderResponse();
    return message;
  },
  fromAmino(_: MsgAddBlockHeaderResponseAmino): MsgAddBlockHeaderResponse {
    const message = createBaseMsgAddBlockHeaderResponse();
    return message;
  },
  toAmino(_: MsgAddBlockHeaderResponse): MsgAddBlockHeaderResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgAddBlockHeaderResponseAminoMsg): MsgAddBlockHeaderResponse {
    return MsgAddBlockHeaderResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddBlockHeaderResponseProtoMsg): MsgAddBlockHeaderResponse {
    return MsgAddBlockHeaderResponse.decode(message.value);
  },
  toProto(message: MsgAddBlockHeaderResponse): Uint8Array {
    return MsgAddBlockHeaderResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgAddBlockHeaderResponse): MsgAddBlockHeaderResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgAddBlockHeaderResponse",
      value: MsgAddBlockHeaderResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateChainParams(): MsgUpdateChainParams {
  return {
    creator: "",
    chainParams: undefined
  };
}
export const MsgUpdateChainParams = {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateChainParams",
  encode(message: MsgUpdateChainParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chainParams !== undefined) {
      ChainParams.encode(message.chainParams, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateChainParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateChainParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chainParams = ChainParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateChainParams>): MsgUpdateChainParams {
    const message = createBaseMsgUpdateChainParams();
    message.creator = object.creator ?? "";
    message.chainParams = object.chainParams !== undefined && object.chainParams !== null ? ChainParams.fromPartial(object.chainParams) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateChainParamsAmino): MsgUpdateChainParams {
    const message = createBaseMsgUpdateChainParams();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.chainParams !== undefined && object.chainParams !== null) {
      message.chainParams = ChainParams.fromAmino(object.chainParams);
    }
    return message;
  },
  toAmino(message: MsgUpdateChainParams): MsgUpdateChainParamsAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.chainParams = message.chainParams ? ChainParams.toAmino(message.chainParams) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateChainParamsAminoMsg): MsgUpdateChainParams {
    return MsgUpdateChainParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateChainParamsProtoMsg): MsgUpdateChainParams {
    return MsgUpdateChainParams.decode(message.value);
  },
  toProto(message: MsgUpdateChainParams): Uint8Array {
    return MsgUpdateChainParams.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateChainParams): MsgUpdateChainParamsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgUpdateChainParams",
      value: MsgUpdateChainParams.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateChainParamsResponse(): MsgUpdateChainParamsResponse {
  return {};
}
export const MsgUpdateChainParamsResponse = {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateChainParamsResponse",
  encode(_: MsgUpdateChainParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateChainParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateChainParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<MsgUpdateChainParamsResponse>): MsgUpdateChainParamsResponse {
    const message = createBaseMsgUpdateChainParamsResponse();
    return message;
  },
  fromAmino(_: MsgUpdateChainParamsResponseAmino): MsgUpdateChainParamsResponse {
    const message = createBaseMsgUpdateChainParamsResponse();
    return message;
  },
  toAmino(_: MsgUpdateChainParamsResponse): MsgUpdateChainParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateChainParamsResponseAminoMsg): MsgUpdateChainParamsResponse {
    return MsgUpdateChainParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateChainParamsResponseProtoMsg): MsgUpdateChainParamsResponse {
    return MsgUpdateChainParamsResponse.decode(message.value);
  },
  toProto(message: MsgUpdateChainParamsResponse): Uint8Array {
    return MsgUpdateChainParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateChainParamsResponse): MsgUpdateChainParamsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgUpdateChainParamsResponse",
      value: MsgUpdateChainParamsResponse.encode(message).finish()
    };
  }
};
function createBaseMsgRemoveChainParams(): MsgRemoveChainParams {
  return {
    creator: "",
    chainId: BigInt(0)
  };
}
export const MsgRemoveChainParams = {
  typeUrl: "/zetachain.zetacore.observer.MsgRemoveChainParams",
  encode(message: MsgRemoveChainParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgRemoveChainParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveChainParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgRemoveChainParams>): MsgRemoveChainParams {
    const message = createBaseMsgRemoveChainParams();
    message.creator = object.creator ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgRemoveChainParamsAmino): MsgRemoveChainParams {
    const message = createBaseMsgRemoveChainParams();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    return message;
  },
  toAmino(message: MsgRemoveChainParams): MsgRemoveChainParamsAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgRemoveChainParamsAminoMsg): MsgRemoveChainParams {
    return MsgRemoveChainParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRemoveChainParamsProtoMsg): MsgRemoveChainParams {
    return MsgRemoveChainParams.decode(message.value);
  },
  toProto(message: MsgRemoveChainParams): Uint8Array {
    return MsgRemoveChainParams.encode(message).finish();
  },
  toProtoMsg(message: MsgRemoveChainParams): MsgRemoveChainParamsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgRemoveChainParams",
      value: MsgRemoveChainParams.encode(message).finish()
    };
  }
};
function createBaseMsgRemoveChainParamsResponse(): MsgRemoveChainParamsResponse {
  return {};
}
export const MsgRemoveChainParamsResponse = {
  typeUrl: "/zetachain.zetacore.observer.MsgRemoveChainParamsResponse",
  encode(_: MsgRemoveChainParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgRemoveChainParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveChainParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<MsgRemoveChainParamsResponse>): MsgRemoveChainParamsResponse {
    const message = createBaseMsgRemoveChainParamsResponse();
    return message;
  },
  fromAmino(_: MsgRemoveChainParamsResponseAmino): MsgRemoveChainParamsResponse {
    const message = createBaseMsgRemoveChainParamsResponse();
    return message;
  },
  toAmino(_: MsgRemoveChainParamsResponse): MsgRemoveChainParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgRemoveChainParamsResponseAminoMsg): MsgRemoveChainParamsResponse {
    return MsgRemoveChainParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRemoveChainParamsResponseProtoMsg): MsgRemoveChainParamsResponse {
    return MsgRemoveChainParamsResponse.decode(message.value);
  },
  toProto(message: MsgRemoveChainParamsResponse): Uint8Array {
    return MsgRemoveChainParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgRemoveChainParamsResponse): MsgRemoveChainParamsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgRemoveChainParamsResponse",
      value: MsgRemoveChainParamsResponse.encode(message).finish()
    };
  }
};
function createBaseMsgAddObserver(): MsgAddObserver {
  return {
    creator: "",
    observerAddress: "",
    zetaclientGranteePubkey: "",
    addNodeAccountOnly: false
  };
}
export const MsgAddObserver = {
  typeUrl: "/zetachain.zetacore.observer.MsgAddObserver",
  encode(message: MsgAddObserver, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.observerAddress !== "") {
      writer.uint32(18).string(message.observerAddress);
    }
    if (message.zetaclientGranteePubkey !== "") {
      writer.uint32(26).string(message.zetaclientGranteePubkey);
    }
    if (message.addNodeAccountOnly === true) {
      writer.uint32(32).bool(message.addNodeAccountOnly);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddObserver {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddObserver();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.observerAddress = reader.string();
          break;
        case 3:
          message.zetaclientGranteePubkey = reader.string();
          break;
        case 4:
          message.addNodeAccountOnly = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgAddObserver>): MsgAddObserver {
    const message = createBaseMsgAddObserver();
    message.creator = object.creator ?? "";
    message.observerAddress = object.observerAddress ?? "";
    message.zetaclientGranteePubkey = object.zetaclientGranteePubkey ?? "";
    message.addNodeAccountOnly = object.addNodeAccountOnly ?? false;
    return message;
  },
  fromAmino(object: MsgAddObserverAmino): MsgAddObserver {
    const message = createBaseMsgAddObserver();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.observer_address !== undefined && object.observer_address !== null) {
      message.observerAddress = object.observer_address;
    }
    if (object.zetaclient_grantee_pubkey !== undefined && object.zetaclient_grantee_pubkey !== null) {
      message.zetaclientGranteePubkey = object.zetaclient_grantee_pubkey;
    }
    if (object.add_node_account_only !== undefined && object.add_node_account_only !== null) {
      message.addNodeAccountOnly = object.add_node_account_only;
    }
    return message;
  },
  toAmino(message: MsgAddObserver): MsgAddObserverAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.observer_address = message.observerAddress;
    obj.zetaclient_grantee_pubkey = message.zetaclientGranteePubkey;
    obj.add_node_account_only = message.addNodeAccountOnly;
    return obj;
  },
  fromAminoMsg(object: MsgAddObserverAminoMsg): MsgAddObserver {
    return MsgAddObserver.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddObserverProtoMsg): MsgAddObserver {
    return MsgAddObserver.decode(message.value);
  },
  toProto(message: MsgAddObserver): Uint8Array {
    return MsgAddObserver.encode(message).finish();
  },
  toProtoMsg(message: MsgAddObserver): MsgAddObserverProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgAddObserver",
      value: MsgAddObserver.encode(message).finish()
    };
  }
};
function createBaseMsgAddObserverResponse(): MsgAddObserverResponse {
  return {};
}
export const MsgAddObserverResponse = {
  typeUrl: "/zetachain.zetacore.observer.MsgAddObserverResponse",
  encode(_: MsgAddObserverResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddObserverResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddObserverResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<MsgAddObserverResponse>): MsgAddObserverResponse {
    const message = createBaseMsgAddObserverResponse();
    return message;
  },
  fromAmino(_: MsgAddObserverResponseAmino): MsgAddObserverResponse {
    const message = createBaseMsgAddObserverResponse();
    return message;
  },
  toAmino(_: MsgAddObserverResponse): MsgAddObserverResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgAddObserverResponseAminoMsg): MsgAddObserverResponse {
    return MsgAddObserverResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddObserverResponseProtoMsg): MsgAddObserverResponse {
    return MsgAddObserverResponse.decode(message.value);
  },
  toProto(message: MsgAddObserverResponse): Uint8Array {
    return MsgAddObserverResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgAddObserverResponse): MsgAddObserverResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgAddObserverResponse",
      value: MsgAddObserverResponse.encode(message).finish()
    };
  }
};
function createBaseMsgAddBlameVote(): MsgAddBlameVote {
  return {
    creator: "",
    chainId: BigInt(0),
    blameInfo: Blame.fromPartial({})
  };
}
export const MsgAddBlameVote = {
  typeUrl: "/zetachain.zetacore.observer.MsgAddBlameVote",
  encode(message: MsgAddBlameVote, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.blameInfo !== undefined) {
      Blame.encode(message.blameInfo, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddBlameVote {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddBlameVote();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.blameInfo = Blame.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgAddBlameVote>): MsgAddBlameVote {
    const message = createBaseMsgAddBlameVote();
    message.creator = object.creator ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.blameInfo = object.blameInfo !== undefined && object.blameInfo !== null ? Blame.fromPartial(object.blameInfo) : undefined;
    return message;
  },
  fromAmino(object: MsgAddBlameVoteAmino): MsgAddBlameVote {
    const message = createBaseMsgAddBlameVote();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.blame_info !== undefined && object.blame_info !== null) {
      message.blameInfo = Blame.fromAmino(object.blame_info);
    }
    return message;
  },
  toAmino(message: MsgAddBlameVote): MsgAddBlameVoteAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.blame_info = message.blameInfo ? Blame.toAmino(message.blameInfo) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgAddBlameVoteAminoMsg): MsgAddBlameVote {
    return MsgAddBlameVote.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddBlameVoteProtoMsg): MsgAddBlameVote {
    return MsgAddBlameVote.decode(message.value);
  },
  toProto(message: MsgAddBlameVote): Uint8Array {
    return MsgAddBlameVote.encode(message).finish();
  },
  toProtoMsg(message: MsgAddBlameVote): MsgAddBlameVoteProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgAddBlameVote",
      value: MsgAddBlameVote.encode(message).finish()
    };
  }
};
function createBaseMsgAddBlameVoteResponse(): MsgAddBlameVoteResponse {
  return {};
}
export const MsgAddBlameVoteResponse = {
  typeUrl: "/zetachain.zetacore.observer.MsgAddBlameVoteResponse",
  encode(_: MsgAddBlameVoteResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddBlameVoteResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddBlameVoteResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<MsgAddBlameVoteResponse>): MsgAddBlameVoteResponse {
    const message = createBaseMsgAddBlameVoteResponse();
    return message;
  },
  fromAmino(_: MsgAddBlameVoteResponseAmino): MsgAddBlameVoteResponse {
    const message = createBaseMsgAddBlameVoteResponse();
    return message;
  },
  toAmino(_: MsgAddBlameVoteResponse): MsgAddBlameVoteResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgAddBlameVoteResponseAminoMsg): MsgAddBlameVoteResponse {
    return MsgAddBlameVoteResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddBlameVoteResponseProtoMsg): MsgAddBlameVoteResponse {
    return MsgAddBlameVoteResponse.decode(message.value);
  },
  toProto(message: MsgAddBlameVoteResponse): Uint8Array {
    return MsgAddBlameVoteResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgAddBlameVoteResponse): MsgAddBlameVoteResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgAddBlameVoteResponse",
      value: MsgAddBlameVoteResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateCrosschainFlags(): MsgUpdateCrosschainFlags {
  return {
    creator: "",
    isInboundEnabled: false,
    isOutboundEnabled: false,
    gasPriceIncreaseFlags: undefined,
    blockHeaderVerificationFlags: undefined
  };
}
export const MsgUpdateCrosschainFlags = {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlags",
  encode(message: MsgUpdateCrosschainFlags, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.isInboundEnabled === true) {
      writer.uint32(24).bool(message.isInboundEnabled);
    }
    if (message.isOutboundEnabled === true) {
      writer.uint32(32).bool(message.isOutboundEnabled);
    }
    if (message.gasPriceIncreaseFlags !== undefined) {
      GasPriceIncreaseFlags.encode(message.gasPriceIncreaseFlags, writer.uint32(42).fork()).ldelim();
    }
    if (message.blockHeaderVerificationFlags !== undefined) {
      BlockHeaderVerificationFlags.encode(message.blockHeaderVerificationFlags, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateCrosschainFlags {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateCrosschainFlags();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 3:
          message.isInboundEnabled = reader.bool();
          break;
        case 4:
          message.isOutboundEnabled = reader.bool();
          break;
        case 5:
          message.gasPriceIncreaseFlags = GasPriceIncreaseFlags.decode(reader, reader.uint32());
          break;
        case 6:
          message.blockHeaderVerificationFlags = BlockHeaderVerificationFlags.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateCrosschainFlags>): MsgUpdateCrosschainFlags {
    const message = createBaseMsgUpdateCrosschainFlags();
    message.creator = object.creator ?? "";
    message.isInboundEnabled = object.isInboundEnabled ?? false;
    message.isOutboundEnabled = object.isOutboundEnabled ?? false;
    message.gasPriceIncreaseFlags = object.gasPriceIncreaseFlags !== undefined && object.gasPriceIncreaseFlags !== null ? GasPriceIncreaseFlags.fromPartial(object.gasPriceIncreaseFlags) : undefined;
    message.blockHeaderVerificationFlags = object.blockHeaderVerificationFlags !== undefined && object.blockHeaderVerificationFlags !== null ? BlockHeaderVerificationFlags.fromPartial(object.blockHeaderVerificationFlags) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateCrosschainFlagsAmino): MsgUpdateCrosschainFlags {
    const message = createBaseMsgUpdateCrosschainFlags();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.isInboundEnabled !== undefined && object.isInboundEnabled !== null) {
      message.isInboundEnabled = object.isInboundEnabled;
    }
    if (object.isOutboundEnabled !== undefined && object.isOutboundEnabled !== null) {
      message.isOutboundEnabled = object.isOutboundEnabled;
    }
    if (object.gasPriceIncreaseFlags !== undefined && object.gasPriceIncreaseFlags !== null) {
      message.gasPriceIncreaseFlags = GasPriceIncreaseFlags.fromAmino(object.gasPriceIncreaseFlags);
    }
    if (object.blockHeaderVerificationFlags !== undefined && object.blockHeaderVerificationFlags !== null) {
      message.blockHeaderVerificationFlags = BlockHeaderVerificationFlags.fromAmino(object.blockHeaderVerificationFlags);
    }
    return message;
  },
  toAmino(message: MsgUpdateCrosschainFlags): MsgUpdateCrosschainFlagsAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.isInboundEnabled = message.isInboundEnabled;
    obj.isOutboundEnabled = message.isOutboundEnabled;
    obj.gasPriceIncreaseFlags = message.gasPriceIncreaseFlags ? GasPriceIncreaseFlags.toAmino(message.gasPriceIncreaseFlags) : undefined;
    obj.blockHeaderVerificationFlags = message.blockHeaderVerificationFlags ? BlockHeaderVerificationFlags.toAmino(message.blockHeaderVerificationFlags) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateCrosschainFlagsAminoMsg): MsgUpdateCrosschainFlags {
    return MsgUpdateCrosschainFlags.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateCrosschainFlagsProtoMsg): MsgUpdateCrosschainFlags {
    return MsgUpdateCrosschainFlags.decode(message.value);
  },
  toProto(message: MsgUpdateCrosschainFlags): Uint8Array {
    return MsgUpdateCrosschainFlags.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateCrosschainFlags): MsgUpdateCrosschainFlagsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlags",
      value: MsgUpdateCrosschainFlags.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateCrosschainFlagsResponse(): MsgUpdateCrosschainFlagsResponse {
  return {};
}
export const MsgUpdateCrosschainFlagsResponse = {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlagsResponse",
  encode(_: MsgUpdateCrosschainFlagsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateCrosschainFlagsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateCrosschainFlagsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<MsgUpdateCrosschainFlagsResponse>): MsgUpdateCrosschainFlagsResponse {
    const message = createBaseMsgUpdateCrosschainFlagsResponse();
    return message;
  },
  fromAmino(_: MsgUpdateCrosschainFlagsResponseAmino): MsgUpdateCrosschainFlagsResponse {
    const message = createBaseMsgUpdateCrosschainFlagsResponse();
    return message;
  },
  toAmino(_: MsgUpdateCrosschainFlagsResponse): MsgUpdateCrosschainFlagsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateCrosschainFlagsResponseAminoMsg): MsgUpdateCrosschainFlagsResponse {
    return MsgUpdateCrosschainFlagsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateCrosschainFlagsResponseProtoMsg): MsgUpdateCrosschainFlagsResponse {
    return MsgUpdateCrosschainFlagsResponse.decode(message.value);
  },
  toProto(message: MsgUpdateCrosschainFlagsResponse): Uint8Array {
    return MsgUpdateCrosschainFlagsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateCrosschainFlagsResponse): MsgUpdateCrosschainFlagsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlagsResponse",
      value: MsgUpdateCrosschainFlagsResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateKeygen(): MsgUpdateKeygen {
  return {
    creator: "",
    block: BigInt(0)
  };
}
export const MsgUpdateKeygen = {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateKeygen",
  encode(message: MsgUpdateKeygen, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.block !== BigInt(0)) {
      writer.uint32(16).int64(message.block);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateKeygen {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateKeygen();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.block = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateKeygen>): MsgUpdateKeygen {
    const message = createBaseMsgUpdateKeygen();
    message.creator = object.creator ?? "";
    message.block = object.block !== undefined && object.block !== null ? BigInt(object.block.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgUpdateKeygenAmino): MsgUpdateKeygen {
    const message = createBaseMsgUpdateKeygen();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.block !== undefined && object.block !== null) {
      message.block = BigInt(object.block);
    }
    return message;
  },
  toAmino(message: MsgUpdateKeygen): MsgUpdateKeygenAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.block = message.block ? message.block.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateKeygenAminoMsg): MsgUpdateKeygen {
    return MsgUpdateKeygen.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateKeygenProtoMsg): MsgUpdateKeygen {
    return MsgUpdateKeygen.decode(message.value);
  },
  toProto(message: MsgUpdateKeygen): Uint8Array {
    return MsgUpdateKeygen.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateKeygen): MsgUpdateKeygenProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgUpdateKeygen",
      value: MsgUpdateKeygen.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateKeygenResponse(): MsgUpdateKeygenResponse {
  return {};
}
export const MsgUpdateKeygenResponse = {
  typeUrl: "/zetachain.zetacore.observer.MsgUpdateKeygenResponse",
  encode(_: MsgUpdateKeygenResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateKeygenResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateKeygenResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<MsgUpdateKeygenResponse>): MsgUpdateKeygenResponse {
    const message = createBaseMsgUpdateKeygenResponse();
    return message;
  },
  fromAmino(_: MsgUpdateKeygenResponseAmino): MsgUpdateKeygenResponse {
    const message = createBaseMsgUpdateKeygenResponse();
    return message;
  },
  toAmino(_: MsgUpdateKeygenResponse): MsgUpdateKeygenResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateKeygenResponseAminoMsg): MsgUpdateKeygenResponse {
    return MsgUpdateKeygenResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateKeygenResponseProtoMsg): MsgUpdateKeygenResponse {
    return MsgUpdateKeygenResponse.decode(message.value);
  },
  toProto(message: MsgUpdateKeygenResponse): Uint8Array {
    return MsgUpdateKeygenResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateKeygenResponse): MsgUpdateKeygenResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.MsgUpdateKeygenResponse",
      value: MsgUpdateKeygenResponse.encode(message).finish()
    };
  }
};