import { Policies, PoliciesAmino, PoliciesSDKType } from "./policies";
import { BinaryReader, BinaryWriter } from "../../binary";
/** GenesisState defines the authority module's genesis state. */
export interface GenesisState {
  policies: Policies;
}
export interface GenesisStateProtoMsg {
  typeUrl: "/zetachain.zetacore.authority.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the authority module's genesis state. */
export interface GenesisStateAmino {
  policies?: PoliciesAmino;
}
export interface GenesisStateAminoMsg {
  type: "/zetachain.zetacore.authority.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the authority module's genesis state. */
export interface GenesisStateSDKType {
  policies: PoliciesSDKType;
}
function createBaseGenesisState(): GenesisState {
  return {
    policies: Policies.fromPartial({})
  };
}
export const GenesisState = {
  typeUrl: "/zetachain.zetacore.authority.GenesisState",
  encode(message: GenesisState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.policies !== undefined) {
      Policies.encode(message.policies, writer.uint32(10).fork()).ldelim();
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
          message.policies = Policies.decode(reader, reader.uint32());
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
    message.policies = object.policies !== undefined && object.policies !== null ? Policies.fromPartial(object.policies) : undefined;
    return message;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    const message = createBaseGenesisState();
    if (object.policies !== undefined && object.policies !== null) {
      message.policies = Policies.fromAmino(object.policies);
    }
    return message;
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    obj.policies = message.policies ? Policies.toAmino(message.policies) : undefined;
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
      typeUrl: "/zetachain.zetacore.authority.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};