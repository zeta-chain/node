import { BinaryReader, BinaryWriter } from "../../binary";
/** store key is tss+chainid */
export interface PendingNonces {
  nonceLow: bigint;
  nonceHigh: bigint;
  chainId: bigint;
  tss: string;
}
export interface PendingNoncesProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.PendingNonces";
  value: Uint8Array;
}
/** store key is tss+chainid */
export interface PendingNoncesAmino {
  nonce_low?: string;
  nonce_high?: string;
  chain_id?: string;
  tss?: string;
}
export interface PendingNoncesAminoMsg {
  type: "/zetachain.zetacore.observer.PendingNonces";
  value: PendingNoncesAmino;
}
/** store key is tss+chainid */
export interface PendingNoncesSDKType {
  nonce_low: bigint;
  nonce_high: bigint;
  chain_id: bigint;
  tss: string;
}
function createBasePendingNonces(): PendingNonces {
  return {
    nonceLow: BigInt(0),
    nonceHigh: BigInt(0),
    chainId: BigInt(0),
    tss: ""
  };
}
export const PendingNonces = {
  typeUrl: "/zetachain.zetacore.observer.PendingNonces",
  encode(message: PendingNonces, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.nonceLow !== BigInt(0)) {
      writer.uint32(8).int64(message.nonceLow);
    }
    if (message.nonceHigh !== BigInt(0)) {
      writer.uint32(16).int64(message.nonceHigh);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(24).int64(message.chainId);
    }
    if (message.tss !== "") {
      writer.uint32(34).string(message.tss);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): PendingNonces {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePendingNonces();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.nonceLow = reader.int64();
          break;
        case 2:
          message.nonceHigh = reader.int64();
          break;
        case 3:
          message.chainId = reader.int64();
          break;
        case 4:
          message.tss = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<PendingNonces>): PendingNonces {
    const message = createBasePendingNonces();
    message.nonceLow = object.nonceLow !== undefined && object.nonceLow !== null ? BigInt(object.nonceLow.toString()) : BigInt(0);
    message.nonceHigh = object.nonceHigh !== undefined && object.nonceHigh !== null ? BigInt(object.nonceHigh.toString()) : BigInt(0);
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.tss = object.tss ?? "";
    return message;
  },
  fromAmino(object: PendingNoncesAmino): PendingNonces {
    const message = createBasePendingNonces();
    if (object.nonce_low !== undefined && object.nonce_low !== null) {
      message.nonceLow = BigInt(object.nonce_low);
    }
    if (object.nonce_high !== undefined && object.nonce_high !== null) {
      message.nonceHigh = BigInt(object.nonce_high);
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.tss !== undefined && object.tss !== null) {
      message.tss = object.tss;
    }
    return message;
  },
  toAmino(message: PendingNonces): PendingNoncesAmino {
    const obj: any = {};
    obj.nonce_low = message.nonceLow ? message.nonceLow.toString() : undefined;
    obj.nonce_high = message.nonceHigh ? message.nonceHigh.toString() : undefined;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.tss = message.tss;
    return obj;
  },
  fromAminoMsg(object: PendingNoncesAminoMsg): PendingNonces {
    return PendingNonces.fromAmino(object.value);
  },
  fromProtoMsg(message: PendingNoncesProtoMsg): PendingNonces {
    return PendingNonces.decode(message.value);
  },
  toProto(message: PendingNonces): Uint8Array {
    return PendingNonces.encode(message).finish();
  },
  toProtoMsg(message: PendingNonces): PendingNoncesProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.PendingNonces",
      value: PendingNonces.encode(message).finish()
    };
  }
};