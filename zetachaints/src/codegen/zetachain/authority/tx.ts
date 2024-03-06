import { Policies, PoliciesAmino, PoliciesSDKType } from "./policies";
import { BinaryReader, BinaryWriter } from "../../binary";
/** MsgUpdatePolicies defines the MsgUpdatePolicies service. */
export interface MsgUpdatePolicies {
  signer: string;
  policies: Policies;
}
export interface MsgUpdatePoliciesProtoMsg {
  typeUrl: "/zetachain.zetacore.authority.MsgUpdatePolicies";
  value: Uint8Array;
}
/** MsgUpdatePolicies defines the MsgUpdatePolicies service. */
export interface MsgUpdatePoliciesAmino {
  signer?: string;
  policies?: PoliciesAmino;
}
export interface MsgUpdatePoliciesAminoMsg {
  type: "/zetachain.zetacore.authority.MsgUpdatePolicies";
  value: MsgUpdatePoliciesAmino;
}
/** MsgUpdatePolicies defines the MsgUpdatePolicies service. */
export interface MsgUpdatePoliciesSDKType {
  signer: string;
  policies: PoliciesSDKType;
}
/** MsgUpdatePoliciesResponse defines the MsgUpdatePoliciesResponse service. */
export interface MsgUpdatePoliciesResponse {}
export interface MsgUpdatePoliciesResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.authority.MsgUpdatePoliciesResponse";
  value: Uint8Array;
}
/** MsgUpdatePoliciesResponse defines the MsgUpdatePoliciesResponse service. */
export interface MsgUpdatePoliciesResponseAmino {}
export interface MsgUpdatePoliciesResponseAminoMsg {
  type: "/zetachain.zetacore.authority.MsgUpdatePoliciesResponse";
  value: MsgUpdatePoliciesResponseAmino;
}
/** MsgUpdatePoliciesResponse defines the MsgUpdatePoliciesResponse service. */
export interface MsgUpdatePoliciesResponseSDKType {}
function createBaseMsgUpdatePolicies(): MsgUpdatePolicies {
  return {
    signer: "",
    policies: Policies.fromPartial({})
  };
}
export const MsgUpdatePolicies = {
  typeUrl: "/zetachain.zetacore.authority.MsgUpdatePolicies",
  encode(message: MsgUpdatePolicies, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.signer !== "") {
      writer.uint32(10).string(message.signer);
    }
    if (message.policies !== undefined) {
      Policies.encode(message.policies, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdatePolicies {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePolicies();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.signer = reader.string();
          break;
        case 2:
          message.policies = Policies.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdatePolicies>): MsgUpdatePolicies {
    const message = createBaseMsgUpdatePolicies();
    message.signer = object.signer ?? "";
    message.policies = object.policies !== undefined && object.policies !== null ? Policies.fromPartial(object.policies) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdatePoliciesAmino): MsgUpdatePolicies {
    const message = createBaseMsgUpdatePolicies();
    if (object.signer !== undefined && object.signer !== null) {
      message.signer = object.signer;
    }
    if (object.policies !== undefined && object.policies !== null) {
      message.policies = Policies.fromAmino(object.policies);
    }
    return message;
  },
  toAmino(message: MsgUpdatePolicies): MsgUpdatePoliciesAmino {
    const obj: any = {};
    obj.signer = message.signer;
    obj.policies = message.policies ? Policies.toAmino(message.policies) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdatePoliciesAminoMsg): MsgUpdatePolicies {
    return MsgUpdatePolicies.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdatePoliciesProtoMsg): MsgUpdatePolicies {
    return MsgUpdatePolicies.decode(message.value);
  },
  toProto(message: MsgUpdatePolicies): Uint8Array {
    return MsgUpdatePolicies.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdatePolicies): MsgUpdatePoliciesProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.authority.MsgUpdatePolicies",
      value: MsgUpdatePolicies.encode(message).finish()
    };
  }
};
function createBaseMsgUpdatePoliciesResponse(): MsgUpdatePoliciesResponse {
  return {};
}
export const MsgUpdatePoliciesResponse = {
  typeUrl: "/zetachain.zetacore.authority.MsgUpdatePoliciesResponse",
  encode(_: MsgUpdatePoliciesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdatePoliciesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePoliciesResponse();
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
  fromPartial(_: Partial<MsgUpdatePoliciesResponse>): MsgUpdatePoliciesResponse {
    const message = createBaseMsgUpdatePoliciesResponse();
    return message;
  },
  fromAmino(_: MsgUpdatePoliciesResponseAmino): MsgUpdatePoliciesResponse {
    const message = createBaseMsgUpdatePoliciesResponse();
    return message;
  },
  toAmino(_: MsgUpdatePoliciesResponse): MsgUpdatePoliciesResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdatePoliciesResponseAminoMsg): MsgUpdatePoliciesResponse {
    return MsgUpdatePoliciesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdatePoliciesResponseProtoMsg): MsgUpdatePoliciesResponse {
    return MsgUpdatePoliciesResponse.decode(message.value);
  },
  toProto(message: MsgUpdatePoliciesResponse): Uint8Array {
    return MsgUpdatePoliciesResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdatePoliciesResponse): MsgUpdatePoliciesResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.authority.MsgUpdatePoliciesResponse",
      value: MsgUpdatePoliciesResponse.encode(message).finish()
    };
  }
};