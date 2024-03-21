import { BinaryReader, BinaryWriter } from "../../binary";
/** Params defines the parameters for the module. */
export interface Params {
  enabled: boolean;
}
export interface ParamsProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.Params";
  value: Uint8Array;
}
/** Params defines the parameters for the module. */
export interface ParamsAmino {
  enabled?: boolean;
}
export interface ParamsAminoMsg {
  type: "/zetachain.zetacore.crosschain.Params";
  value: ParamsAmino;
}
/** Params defines the parameters for the module. */
export interface ParamsSDKType {
  enabled: boolean;
}
function createBaseParams(): Params {
  return {
    enabled: false
  };
}
export const Params = {
  typeUrl: "/zetachain.zetacore.crosschain.Params",
  encode(message: Params, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.enabled === true) {
      writer.uint32(8).bool(message.enabled);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Params {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.enabled = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Params>): Params {
    const message = createBaseParams();
    message.enabled = object.enabled ?? false;
    return message;
  },
  fromAmino(object: ParamsAmino): Params {
    const message = createBaseParams();
    if (object.enabled !== undefined && object.enabled !== null) {
      message.enabled = object.enabled;
    }
    return message;
  },
  toAmino(message: Params): ParamsAmino {
    const obj: any = {};
    obj.enabled = message.enabled;
    return obj;
  },
  fromAminoMsg(object: ParamsAminoMsg): Params {
    return Params.fromAmino(object.value);
  },
  fromProtoMsg(message: ParamsProtoMsg): Params {
    return Params.decode(message.value);
  },
  toProto(message: Params): Uint8Array {
    return Params.encode(message).finish();
  },
  toProtoMsg(message: Params): ParamsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.Params",
      value: Params.encode(message).finish()
    };
  }
};