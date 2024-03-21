import { BinaryReader, BinaryWriter } from "../../binary";
export interface InTxHashToCctx {
  inTxHash: string;
  cctxIndex: string[];
}
export interface InTxHashToCctxProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.InTxHashToCctx";
  value: Uint8Array;
}
export interface InTxHashToCctxAmino {
  in_tx_hash?: string;
  cctx_index?: string[];
}
export interface InTxHashToCctxAminoMsg {
  type: "/zetachain.zetacore.crosschain.InTxHashToCctx";
  value: InTxHashToCctxAmino;
}
export interface InTxHashToCctxSDKType {
  in_tx_hash: string;
  cctx_index: string[];
}
function createBaseInTxHashToCctx(): InTxHashToCctx {
  return {
    inTxHash: "",
    cctxIndex: []
  };
}
export const InTxHashToCctx = {
  typeUrl: "/zetachain.zetacore.crosschain.InTxHashToCctx",
  encode(message: InTxHashToCctx, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.inTxHash !== "") {
      writer.uint32(10).string(message.inTxHash);
    }
    for (const v of message.cctxIndex) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): InTxHashToCctx {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInTxHashToCctx();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.inTxHash = reader.string();
          break;
        case 2:
          message.cctxIndex.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<InTxHashToCctx>): InTxHashToCctx {
    const message = createBaseInTxHashToCctx();
    message.inTxHash = object.inTxHash ?? "";
    message.cctxIndex = object.cctxIndex?.map(e => e) || [];
    return message;
  },
  fromAmino(object: InTxHashToCctxAmino): InTxHashToCctx {
    const message = createBaseInTxHashToCctx();
    if (object.in_tx_hash !== undefined && object.in_tx_hash !== null) {
      message.inTxHash = object.in_tx_hash;
    }
    message.cctxIndex = object.cctx_index?.map(e => e) || [];
    return message;
  },
  toAmino(message: InTxHashToCctx): InTxHashToCctxAmino {
    const obj: any = {};
    obj.in_tx_hash = message.inTxHash;
    if (message.cctxIndex) {
      obj.cctx_index = message.cctxIndex.map(e => e);
    } else {
      obj.cctx_index = [];
    }
    return obj;
  },
  fromAminoMsg(object: InTxHashToCctxAminoMsg): InTxHashToCctx {
    return InTxHashToCctx.fromAmino(object.value);
  },
  fromProtoMsg(message: InTxHashToCctxProtoMsg): InTxHashToCctx {
    return InTxHashToCctx.decode(message.value);
  },
  toProto(message: InTxHashToCctx): Uint8Array {
    return InTxHashToCctx.encode(message).finish();
  },
  toProtoMsg(message: InTxHashToCctx): InTxHashToCctxProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.InTxHashToCctx",
      value: InTxHashToCctx.encode(message).finish()
    };
  }
};