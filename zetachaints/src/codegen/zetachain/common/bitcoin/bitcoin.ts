import { BinaryReader, BinaryWriter } from "../../../binary";
import { bytesFromBase64, base64FromBytes } from "../../../helpers";
export interface Proof {
  txBytes: Uint8Array;
  path: Uint8Array;
  index: number;
}
export interface ProofProtoMsg {
  typeUrl: "/bitcoin.Proof";
  value: Uint8Array;
}
export interface ProofAmino {
  tx_bytes?: string;
  path?: string;
  index?: number;
}
export interface ProofAminoMsg {
  type: "/bitcoin.Proof";
  value: ProofAmino;
}
export interface ProofSDKType {
  tx_bytes: Uint8Array;
  path: Uint8Array;
  index: number;
}
function createBaseProof(): Proof {
  return {
    txBytes: new Uint8Array(),
    path: new Uint8Array(),
    index: 0
  };
}
export const Proof = {
  typeUrl: "/bitcoin.Proof",
  encode(message: Proof, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.txBytes.length !== 0) {
      writer.uint32(10).bytes(message.txBytes);
    }
    if (message.path.length !== 0) {
      writer.uint32(18).bytes(message.path);
    }
    if (message.index !== 0) {
      writer.uint32(24).uint32(message.index);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Proof {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProof();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.txBytes = reader.bytes();
          break;
        case 2:
          message.path = reader.bytes();
          break;
        case 3:
          message.index = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Proof>): Proof {
    const message = createBaseProof();
    message.txBytes = object.txBytes ?? new Uint8Array();
    message.path = object.path ?? new Uint8Array();
    message.index = object.index ?? 0;
    return message;
  },
  fromAmino(object: ProofAmino): Proof {
    const message = createBaseProof();
    if (object.tx_bytes !== undefined && object.tx_bytes !== null) {
      message.txBytes = bytesFromBase64(object.tx_bytes);
    }
    if (object.path !== undefined && object.path !== null) {
      message.path = bytesFromBase64(object.path);
    }
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    return message;
  },
  toAmino(message: Proof): ProofAmino {
    const obj: any = {};
    obj.tx_bytes = message.txBytes ? base64FromBytes(message.txBytes) : undefined;
    obj.path = message.path ? base64FromBytes(message.path) : undefined;
    obj.index = message.index;
    return obj;
  },
  fromAminoMsg(object: ProofAminoMsg): Proof {
    return Proof.fromAmino(object.value);
  },
  fromProtoMsg(message: ProofProtoMsg): Proof {
    return Proof.decode(message.value);
  },
  toProto(message: Proof): Uint8Array {
    return Proof.encode(message).finish();
  },
  toProtoMsg(message: Proof): ProofProtoMsg {
    return {
      typeUrl: "/bitcoin.Proof",
      value: Proof.encode(message).finish()
    };
  }
};