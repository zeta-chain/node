import { BinaryReader, BinaryWriter } from "../../binary";
export interface ChainNonces {
  creator: string;
  index: string;
  chainId: bigint;
  nonce: bigint;
  signers: string[];
  finalizedHeight: bigint;
}
export interface ChainNoncesProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.ChainNonces";
  value: Uint8Array;
}
export interface ChainNoncesAmino {
  creator?: string;
  index?: string;
  chain_id?: string;
  nonce?: string;
  signers?: string[];
  finalizedHeight?: string;
}
export interface ChainNoncesAminoMsg {
  type: "/zetachain.zetacore.observer.ChainNonces";
  value: ChainNoncesAmino;
}
export interface ChainNoncesSDKType {
  creator: string;
  index: string;
  chain_id: bigint;
  nonce: bigint;
  signers: string[];
  finalizedHeight: bigint;
}
function createBaseChainNonces(): ChainNonces {
  return {
    creator: "",
    index: "",
    chainId: BigInt(0),
    nonce: BigInt(0),
    signers: [],
    finalizedHeight: BigInt(0)
  };
}
export const ChainNonces = {
  typeUrl: "/zetachain.zetacore.observer.ChainNonces",
  encode(message: ChainNonces, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.index !== "") {
      writer.uint32(18).string(message.index);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(24).int64(message.chainId);
    }
    if (message.nonce !== BigInt(0)) {
      writer.uint32(32).uint64(message.nonce);
    }
    for (const v of message.signers) {
      writer.uint32(42).string(v!);
    }
    if (message.finalizedHeight !== BigInt(0)) {
      writer.uint32(48).uint64(message.finalizedHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ChainNonces {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChainNonces();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.index = reader.string();
          break;
        case 3:
          message.chainId = reader.int64();
          break;
        case 4:
          message.nonce = reader.uint64();
          break;
        case 5:
          message.signers.push(reader.string());
          break;
        case 6:
          message.finalizedHeight = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ChainNonces>): ChainNonces {
    const message = createBaseChainNonces();
    message.creator = object.creator ?? "";
    message.index = object.index ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.nonce = object.nonce !== undefined && object.nonce !== null ? BigInt(object.nonce.toString()) : BigInt(0);
    message.signers = object.signers?.map(e => e) || [];
    message.finalizedHeight = object.finalizedHeight !== undefined && object.finalizedHeight !== null ? BigInt(object.finalizedHeight.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: ChainNoncesAmino): ChainNonces {
    const message = createBaseChainNonces();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = BigInt(object.nonce);
    }
    message.signers = object.signers?.map(e => e) || [];
    if (object.finalizedHeight !== undefined && object.finalizedHeight !== null) {
      message.finalizedHeight = BigInt(object.finalizedHeight);
    }
    return message;
  },
  toAmino(message: ChainNonces): ChainNoncesAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.index = message.index;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.nonce = message.nonce ? message.nonce.toString() : undefined;
    if (message.signers) {
      obj.signers = message.signers.map(e => e);
    } else {
      obj.signers = [];
    }
    obj.finalizedHeight = message.finalizedHeight ? message.finalizedHeight.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: ChainNoncesAminoMsg): ChainNonces {
    return ChainNonces.fromAmino(object.value);
  },
  fromProtoMsg(message: ChainNoncesProtoMsg): ChainNonces {
    return ChainNonces.decode(message.value);
  },
  toProto(message: ChainNonces): Uint8Array {
    return ChainNonces.encode(message).finish();
  },
  toProtoMsg(message: ChainNonces): ChainNoncesProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.ChainNonces",
      value: ChainNonces.encode(message).finish()
    };
  }
};