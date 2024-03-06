import { PubKeySet, PubKeySetAmino, PubKeySetSDKType } from "../common/common";
import { BinaryReader, BinaryWriter } from "../../binary";
export enum NodeStatus {
  Unknown = 0,
  Whitelisted = 1,
  Standby = 2,
  Ready = 3,
  Active = 4,
  Disabled = 5,
  UNRECOGNIZED = -1,
}
export const NodeStatusSDKType = NodeStatus;
export const NodeStatusAmino = NodeStatus;
export function nodeStatusFromJSON(object: any): NodeStatus {
  switch (object) {
    case 0:
    case "Unknown":
      return NodeStatus.Unknown;
    case 1:
    case "Whitelisted":
      return NodeStatus.Whitelisted;
    case 2:
    case "Standby":
      return NodeStatus.Standby;
    case 3:
    case "Ready":
      return NodeStatus.Ready;
    case 4:
    case "Active":
      return NodeStatus.Active;
    case 5:
    case "Disabled":
      return NodeStatus.Disabled;
    case -1:
    case "UNRECOGNIZED":
    default:
      return NodeStatus.UNRECOGNIZED;
  }
}
export function nodeStatusToJSON(object: NodeStatus): string {
  switch (object) {
    case NodeStatus.Unknown:
      return "Unknown";
    case NodeStatus.Whitelisted:
      return "Whitelisted";
    case NodeStatus.Standby:
      return "Standby";
    case NodeStatus.Ready:
      return "Ready";
    case NodeStatus.Active:
      return "Active";
    case NodeStatus.Disabled:
      return "Disabled";
    case NodeStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export interface NodeAccount {
  operator: string;
  granteeAddress: string;
  granteePubkey?: PubKeySet;
  nodeStatus: NodeStatus;
}
export interface NodeAccountProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.NodeAccount";
  value: Uint8Array;
}
export interface NodeAccountAmino {
  operator?: string;
  granteeAddress?: string;
  granteePubkey?: PubKeySetAmino;
  nodeStatus?: NodeStatus;
}
export interface NodeAccountAminoMsg {
  type: "/zetachain.zetacore.observer.NodeAccount";
  value: NodeAccountAmino;
}
export interface NodeAccountSDKType {
  operator: string;
  granteeAddress: string;
  granteePubkey?: PubKeySetSDKType;
  nodeStatus: NodeStatus;
}
function createBaseNodeAccount(): NodeAccount {
  return {
    operator: "",
    granteeAddress: "",
    granteePubkey: undefined,
    nodeStatus: 0
  };
}
export const NodeAccount = {
  typeUrl: "/zetachain.zetacore.observer.NodeAccount",
  encode(message: NodeAccount, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.operator !== "") {
      writer.uint32(10).string(message.operator);
    }
    if (message.granteeAddress !== "") {
      writer.uint32(18).string(message.granteeAddress);
    }
    if (message.granteePubkey !== undefined) {
      PubKeySet.encode(message.granteePubkey, writer.uint32(26).fork()).ldelim();
    }
    if (message.nodeStatus !== 0) {
      writer.uint32(32).int32(message.nodeStatus);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): NodeAccount {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNodeAccount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.operator = reader.string();
          break;
        case 2:
          message.granteeAddress = reader.string();
          break;
        case 3:
          message.granteePubkey = PubKeySet.decode(reader, reader.uint32());
          break;
        case 4:
          message.nodeStatus = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<NodeAccount>): NodeAccount {
    const message = createBaseNodeAccount();
    message.operator = object.operator ?? "";
    message.granteeAddress = object.granteeAddress ?? "";
    message.granteePubkey = object.granteePubkey !== undefined && object.granteePubkey !== null ? PubKeySet.fromPartial(object.granteePubkey) : undefined;
    message.nodeStatus = object.nodeStatus ?? 0;
    return message;
  },
  fromAmino(object: NodeAccountAmino): NodeAccount {
    const message = createBaseNodeAccount();
    if (object.operator !== undefined && object.operator !== null) {
      message.operator = object.operator;
    }
    if (object.granteeAddress !== undefined && object.granteeAddress !== null) {
      message.granteeAddress = object.granteeAddress;
    }
    if (object.granteePubkey !== undefined && object.granteePubkey !== null) {
      message.granteePubkey = PubKeySet.fromAmino(object.granteePubkey);
    }
    if (object.nodeStatus !== undefined && object.nodeStatus !== null) {
      message.nodeStatus = nodeStatusFromJSON(object.nodeStatus);
    }
    return message;
  },
  toAmino(message: NodeAccount): NodeAccountAmino {
    const obj: any = {};
    obj.operator = message.operator;
    obj.granteeAddress = message.granteeAddress;
    obj.granteePubkey = message.granteePubkey ? PubKeySet.toAmino(message.granteePubkey) : undefined;
    obj.nodeStatus = message.nodeStatus;
    return obj;
  },
  fromAminoMsg(object: NodeAccountAminoMsg): NodeAccount {
    return NodeAccount.fromAmino(object.value);
  },
  fromProtoMsg(message: NodeAccountProtoMsg): NodeAccount {
    return NodeAccount.decode(message.value);
  },
  toProto(message: NodeAccount): Uint8Array {
    return NodeAccount.encode(message).finish();
  },
  toProtoMsg(message: NodeAccount): NodeAccountProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.NodeAccount",
      value: NodeAccount.encode(message).finish()
    };
  }
};