import { BinaryReader, BinaryWriter } from "../../binary";
/** PolicyType defines the type of policy */
export enum PolicyType {
  groupEmergency = 0,
  groupAdmin = 1,
  UNRECOGNIZED = -1,
}
export const PolicyTypeSDKType = PolicyType;
export const PolicyTypeAmino = PolicyType;
export function policyTypeFromJSON(object: any): PolicyType {
  switch (object) {
    case 0:
    case "groupEmergency":
      return PolicyType.groupEmergency;
    case 1:
    case "groupAdmin":
      return PolicyType.groupAdmin;
    case -1:
    case "UNRECOGNIZED":
    default:
      return PolicyType.UNRECOGNIZED;
  }
}
export function policyTypeToJSON(object: PolicyType): string {
  switch (object) {
    case PolicyType.groupEmergency:
      return "groupEmergency";
    case PolicyType.groupAdmin:
      return "groupAdmin";
    case PolicyType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export interface Policy {
  policyType: PolicyType;
  address: string;
}
export interface PolicyProtoMsg {
  typeUrl: "/zetachain.zetacore.authority.Policy";
  value: Uint8Array;
}
export interface PolicyAmino {
  policy_type?: PolicyType;
  address?: string;
}
export interface PolicyAminoMsg {
  type: "/zetachain.zetacore.authority.Policy";
  value: PolicyAmino;
}
export interface PolicySDKType {
  policy_type: PolicyType;
  address: string;
}
/** Policy contains info about authority policies */
export interface Policies {
  items: Policy[];
}
export interface PoliciesProtoMsg {
  typeUrl: "/zetachain.zetacore.authority.Policies";
  value: Uint8Array;
}
/** Policy contains info about authority policies */
export interface PoliciesAmino {
  items?: PolicyAmino[];
}
export interface PoliciesAminoMsg {
  type: "/zetachain.zetacore.authority.Policies";
  value: PoliciesAmino;
}
/** Policy contains info about authority policies */
export interface PoliciesSDKType {
  items: PolicySDKType[];
}
function createBasePolicy(): Policy {
  return {
    policyType: 0,
    address: ""
  };
}
export const Policy = {
  typeUrl: "/zetachain.zetacore.authority.Policy",
  encode(message: Policy, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.policyType !== 0) {
      writer.uint32(8).int32(message.policyType);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Policy {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policyType = (reader.int32() as any);
          break;
        case 2:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Policy>): Policy {
    const message = createBasePolicy();
    message.policyType = object.policyType ?? 0;
    message.address = object.address ?? "";
    return message;
  },
  fromAmino(object: PolicyAmino): Policy {
    const message = createBasePolicy();
    if (object.policy_type !== undefined && object.policy_type !== null) {
      message.policyType = policyTypeFromJSON(object.policy_type);
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    return message;
  },
  toAmino(message: Policy): PolicyAmino {
    const obj: any = {};
    obj.policy_type = message.policyType;
    obj.address = message.address;
    return obj;
  },
  fromAminoMsg(object: PolicyAminoMsg): Policy {
    return Policy.fromAmino(object.value);
  },
  fromProtoMsg(message: PolicyProtoMsg): Policy {
    return Policy.decode(message.value);
  },
  toProto(message: Policy): Uint8Array {
    return Policy.encode(message).finish();
  },
  toProtoMsg(message: Policy): PolicyProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.authority.Policy",
      value: Policy.encode(message).finish()
    };
  }
};
function createBasePolicies(): Policies {
  return {
    items: []
  };
}
export const Policies = {
  typeUrl: "/zetachain.zetacore.authority.Policies",
  encode(message: Policies, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.items) {
      Policy.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Policies {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePolicies();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.items.push(Policy.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Policies>): Policies {
    const message = createBasePolicies();
    message.items = object.items?.map(e => Policy.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: PoliciesAmino): Policies {
    const message = createBasePolicies();
    message.items = object.items?.map(e => Policy.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: Policies): PoliciesAmino {
    const obj: any = {};
    if (message.items) {
      obj.items = message.items.map(e => e ? Policy.toAmino(e) : undefined);
    } else {
      obj.items = [];
    }
    return obj;
  },
  fromAminoMsg(object: PoliciesAminoMsg): Policies {
    return Policies.fromAmino(object.value);
  },
  fromProtoMsg(message: PoliciesProtoMsg): Policies {
    return Policies.decode(message.value);
  },
  toProto(message: Policies): Uint8Array {
    return Policies.encode(message).finish();
  },
  toProtoMsg(message: Policies): PoliciesProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.authority.Policies",
      value: Policies.encode(message).finish()
    };
  }
};