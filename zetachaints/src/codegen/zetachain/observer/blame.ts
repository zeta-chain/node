import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
export interface Node {
  pubKey: string;
  blameData: Uint8Array;
  blameSignature: Uint8Array;
}
export interface NodeProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.Node";
  value: Uint8Array;
}
export interface NodeAmino {
  pub_key?: string;
  blame_data?: string;
  blame_signature?: string;
}
export interface NodeAminoMsg {
  type: "/zetachain.zetacore.observer.Node";
  value: NodeAmino;
}
export interface NodeSDKType {
  pub_key: string;
  blame_data: Uint8Array;
  blame_signature: Uint8Array;
}
export interface Blame {
  index: string;
  failureReason: string;
  nodes: Node[];
}
export interface BlameProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.Blame";
  value: Uint8Array;
}
export interface BlameAmino {
  index?: string;
  failure_reason?: string;
  nodes?: NodeAmino[];
}
export interface BlameAminoMsg {
  type: "/zetachain.zetacore.observer.Blame";
  value: BlameAmino;
}
export interface BlameSDKType {
  index: string;
  failure_reason: string;
  nodes: NodeSDKType[];
}
function createBaseNode(): Node {
  return {
    pubKey: "",
    blameData: new Uint8Array(),
    blameSignature: new Uint8Array()
  };
}
export const Node = {
  typeUrl: "/zetachain.zetacore.observer.Node",
  encode(message: Node, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pubKey !== "") {
      writer.uint32(10).string(message.pubKey);
    }
    if (message.blameData.length !== 0) {
      writer.uint32(18).bytes(message.blameData);
    }
    if (message.blameSignature.length !== 0) {
      writer.uint32(26).bytes(message.blameSignature);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Node {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNode();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pubKey = reader.string();
          break;
        case 2:
          message.blameData = reader.bytes();
          break;
        case 3:
          message.blameSignature = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Node>): Node {
    const message = createBaseNode();
    message.pubKey = object.pubKey ?? "";
    message.blameData = object.blameData ?? new Uint8Array();
    message.blameSignature = object.blameSignature ?? new Uint8Array();
    return message;
  },
  fromAmino(object: NodeAmino): Node {
    const message = createBaseNode();
    if (object.pub_key !== undefined && object.pub_key !== null) {
      message.pubKey = object.pub_key;
    }
    if (object.blame_data !== undefined && object.blame_data !== null) {
      message.blameData = bytesFromBase64(object.blame_data);
    }
    if (object.blame_signature !== undefined && object.blame_signature !== null) {
      message.blameSignature = bytesFromBase64(object.blame_signature);
    }
    return message;
  },
  toAmino(message: Node): NodeAmino {
    const obj: any = {};
    obj.pub_key = message.pubKey;
    obj.blame_data = message.blameData ? base64FromBytes(message.blameData) : undefined;
    obj.blame_signature = message.blameSignature ? base64FromBytes(message.blameSignature) : undefined;
    return obj;
  },
  fromAminoMsg(object: NodeAminoMsg): Node {
    return Node.fromAmino(object.value);
  },
  fromProtoMsg(message: NodeProtoMsg): Node {
    return Node.decode(message.value);
  },
  toProto(message: Node): Uint8Array {
    return Node.encode(message).finish();
  },
  toProtoMsg(message: Node): NodeProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.Node",
      value: Node.encode(message).finish()
    };
  }
};
function createBaseBlame(): Blame {
  return {
    index: "",
    failureReason: "",
    nodes: []
  };
}
export const Blame = {
  typeUrl: "/zetachain.zetacore.observer.Blame",
  encode(message: Blame, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    if (message.failureReason !== "") {
      writer.uint32(18).string(message.failureReason);
    }
    for (const v of message.nodes) {
      Node.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Blame {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlame();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        case 2:
          message.failureReason = reader.string();
          break;
        case 3:
          message.nodes.push(Node.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Blame>): Blame {
    const message = createBaseBlame();
    message.index = object.index ?? "";
    message.failureReason = object.failureReason ?? "";
    message.nodes = object.nodes?.map(e => Node.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: BlameAmino): Blame {
    const message = createBaseBlame();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    if (object.failure_reason !== undefined && object.failure_reason !== null) {
      message.failureReason = object.failure_reason;
    }
    message.nodes = object.nodes?.map(e => Node.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: Blame): BlameAmino {
    const obj: any = {};
    obj.index = message.index;
    obj.failure_reason = message.failureReason;
    if (message.nodes) {
      obj.nodes = message.nodes.map(e => e ? Node.toAmino(e) : undefined);
    } else {
      obj.nodes = [];
    }
    return obj;
  },
  fromAminoMsg(object: BlameAminoMsg): Blame {
    return Blame.fromAmino(object.value);
  },
  fromProtoMsg(message: BlameProtoMsg): Blame {
    return Blame.decode(message.value);
  },
  toProto(message: Blame): Uint8Array {
    return Blame.encode(message).finish();
  },
  toProtoMsg(message: Blame): BlameProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.Blame",
      value: Blame.encode(message).finish()
    };
  }
};