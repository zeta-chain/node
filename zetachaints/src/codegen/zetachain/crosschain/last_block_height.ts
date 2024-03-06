import { BinaryReader, BinaryWriter } from "../../binary";
export interface LastBlockHeight {
  creator: string;
  index: string;
  chain: string;
  lastSendHeight: bigint;
  lastReceiveHeight: bigint;
}
export interface LastBlockHeightProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.LastBlockHeight";
  value: Uint8Array;
}
export interface LastBlockHeightAmino {
  creator?: string;
  index?: string;
  chain?: string;
  lastSendHeight?: string;
  lastReceiveHeight?: string;
}
export interface LastBlockHeightAminoMsg {
  type: "/zetachain.zetacore.crosschain.LastBlockHeight";
  value: LastBlockHeightAmino;
}
export interface LastBlockHeightSDKType {
  creator: string;
  index: string;
  chain: string;
  lastSendHeight: bigint;
  lastReceiveHeight: bigint;
}
function createBaseLastBlockHeight(): LastBlockHeight {
  return {
    creator: "",
    index: "",
    chain: "",
    lastSendHeight: BigInt(0),
    lastReceiveHeight: BigInt(0)
  };
}
export const LastBlockHeight = {
  typeUrl: "/zetachain.zetacore.crosschain.LastBlockHeight",
  encode(message: LastBlockHeight, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.index !== "") {
      writer.uint32(18).string(message.index);
    }
    if (message.chain !== "") {
      writer.uint32(26).string(message.chain);
    }
    if (message.lastSendHeight !== BigInt(0)) {
      writer.uint32(32).uint64(message.lastSendHeight);
    }
    if (message.lastReceiveHeight !== BigInt(0)) {
      writer.uint32(40).uint64(message.lastReceiveHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): LastBlockHeight {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLastBlockHeight();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.index = reader.string();
          break;
        case 3:
          message.chain = reader.string();
          break;
        case 4:
          message.lastSendHeight = reader.uint64();
          break;
        case 5:
          message.lastReceiveHeight = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<LastBlockHeight>): LastBlockHeight {
    const message = createBaseLastBlockHeight();
    message.creator = object.creator ?? "";
    message.index = object.index ?? "";
    message.chain = object.chain ?? "";
    message.lastSendHeight = object.lastSendHeight !== undefined && object.lastSendHeight !== null ? BigInt(object.lastSendHeight.toString()) : BigInt(0);
    message.lastReceiveHeight = object.lastReceiveHeight !== undefined && object.lastReceiveHeight !== null ? BigInt(object.lastReceiveHeight.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: LastBlockHeightAmino): LastBlockHeight {
    const message = createBaseLastBlockHeight();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    if (object.chain !== undefined && object.chain !== null) {
      message.chain = object.chain;
    }
    if (object.lastSendHeight !== undefined && object.lastSendHeight !== null) {
      message.lastSendHeight = BigInt(object.lastSendHeight);
    }
    if (object.lastReceiveHeight !== undefined && object.lastReceiveHeight !== null) {
      message.lastReceiveHeight = BigInt(object.lastReceiveHeight);
    }
    return message;
  },
  toAmino(message: LastBlockHeight): LastBlockHeightAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.index = message.index;
    obj.chain = message.chain;
    obj.lastSendHeight = message.lastSendHeight ? message.lastSendHeight.toString() : undefined;
    obj.lastReceiveHeight = message.lastReceiveHeight ? message.lastReceiveHeight.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: LastBlockHeightAminoMsg): LastBlockHeight {
    return LastBlockHeight.fromAmino(object.value);
  },
  fromProtoMsg(message: LastBlockHeightProtoMsg): LastBlockHeight {
    return LastBlockHeight.decode(message.value);
  },
  toProto(message: LastBlockHeight): Uint8Array {
    return LastBlockHeight.encode(message).finish();
  },
  toProtoMsg(message: LastBlockHeight): LastBlockHeightProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.LastBlockHeight",
      value: LastBlockHeight.encode(message).finish()
    };
  }
};