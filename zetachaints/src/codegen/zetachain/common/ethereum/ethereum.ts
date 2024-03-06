import { BinaryReader, BinaryWriter } from "../../../binary";
import { bytesFromBase64, base64FromBytes } from "../../../helpers";
export interface Proof {
  keys: Uint8Array[];
  values: Uint8Array[];
}
export interface ProofProtoMsg {
  typeUrl: "/ethereum.Proof";
  value: Uint8Array;
}
export interface ProofAmino {
  keys?: string[];
  values?: string[];
}
export interface ProofAminoMsg {
  type: "/ethereum.Proof";
  value: ProofAmino;
}
export interface ProofSDKType {
  keys: Uint8Array[];
  values: Uint8Array[];
}
function createBaseProof(): Proof {
  return {
    keys: [],
    values: []
  };
}
export const Proof = {
  typeUrl: "/ethereum.Proof",
  encode(message: Proof, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.keys) {
      writer.uint32(10).bytes(v!);
    }
    for (const v of message.values) {
      writer.uint32(18).bytes(v!);
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
          message.keys.push(reader.bytes());
          break;
        case 2:
          message.values.push(reader.bytes());
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
    message.keys = object.keys?.map(e => e) || [];
    message.values = object.values?.map(e => e) || [];
    return message;
  },
  fromAmino(object: ProofAmino): Proof {
    const message = createBaseProof();
    message.keys = object.keys?.map(e => bytesFromBase64(e)) || [];
    message.values = object.values?.map(e => bytesFromBase64(e)) || [];
    return message;
  },
  toAmino(message: Proof): ProofAmino {
    const obj: any = {};
    if (message.keys) {
      obj.keys = message.keys.map(e => base64FromBytes(e));
    } else {
      obj.keys = [];
    }
    if (message.values) {
      obj.values = message.values.map(e => base64FromBytes(e));
    } else {
      obj.values = [];
    }
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
      typeUrl: "/ethereum.Proof",
      value: Proof.encode(message).finish()
    };
  }
};