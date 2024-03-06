import { BinaryReader, BinaryWriter } from "../../binary";
/** store key is tss+chainid+nonce */
export interface NonceToCctx {
  chainId: bigint;
  nonce: bigint;
  cctxIndex: string;
  tss: string;
}
export interface NonceToCctxProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.NonceToCctx";
  value: Uint8Array;
}
/** store key is tss+chainid+nonce */
export interface NonceToCctxAmino {
  chain_id?: string;
  nonce?: string;
  cctxIndex?: string;
  tss?: string;
}
export interface NonceToCctxAminoMsg {
  type: "/zetachain.zetacore.observer.NonceToCctx";
  value: NonceToCctxAmino;
}
/** store key is tss+chainid+nonce */
export interface NonceToCctxSDKType {
  chain_id: bigint;
  nonce: bigint;
  cctxIndex: string;
  tss: string;
}
function createBaseNonceToCctx(): NonceToCctx {
  return {
    chainId: BigInt(0),
    nonce: BigInt(0),
    cctxIndex: "",
    tss: ""
  };
}
export const NonceToCctx = {
  typeUrl: "/zetachain.zetacore.observer.NonceToCctx",
  encode(message: NonceToCctx, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.nonce !== BigInt(0)) {
      writer.uint32(16).int64(message.nonce);
    }
    if (message.cctxIndex !== "") {
      writer.uint32(26).string(message.cctxIndex);
    }
    if (message.tss !== "") {
      writer.uint32(34).string(message.tss);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): NonceToCctx {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNonceToCctx();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.nonce = reader.int64();
          break;
        case 3:
          message.cctxIndex = reader.string();
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
  fromPartial(object: Partial<NonceToCctx>): NonceToCctx {
    const message = createBaseNonceToCctx();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.nonce = object.nonce !== undefined && object.nonce !== null ? BigInt(object.nonce.toString()) : BigInt(0);
    message.cctxIndex = object.cctxIndex ?? "";
    message.tss = object.tss ?? "";
    return message;
  },
  fromAmino(object: NonceToCctxAmino): NonceToCctx {
    const message = createBaseNonceToCctx();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = BigInt(object.nonce);
    }
    if (object.cctxIndex !== undefined && object.cctxIndex !== null) {
      message.cctxIndex = object.cctxIndex;
    }
    if (object.tss !== undefined && object.tss !== null) {
      message.tss = object.tss;
    }
    return message;
  },
  toAmino(message: NonceToCctx): NonceToCctxAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.nonce = message.nonce ? message.nonce.toString() : undefined;
    obj.cctxIndex = message.cctxIndex;
    obj.tss = message.tss;
    return obj;
  },
  fromAminoMsg(object: NonceToCctxAminoMsg): NonceToCctx {
    return NonceToCctx.fromAmino(object.value);
  },
  fromProtoMsg(message: NonceToCctxProtoMsg): NonceToCctx {
    return NonceToCctx.decode(message.value);
  },
  toProto(message: NonceToCctx): Uint8Array {
    return NonceToCctx.encode(message).finish();
  },
  toProtoMsg(message: NonceToCctx): NonceToCctxProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.NonceToCctx",
      value: NonceToCctx.encode(message).finish()
    };
  }
};