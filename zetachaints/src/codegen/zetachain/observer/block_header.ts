import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
export interface BlockHeaderState {
  chainId: bigint;
  latestHeight: bigint;
  earliestHeight: bigint;
  latestBlockHash: Uint8Array;
}
export interface BlockHeaderStateProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.BlockHeaderState";
  value: Uint8Array;
}
export interface BlockHeaderStateAmino {
  chain_id?: string;
  latest_height?: string;
  earliest_height?: string;
  latest_block_hash?: string;
}
export interface BlockHeaderStateAminoMsg {
  type: "/zetachain.zetacore.observer.BlockHeaderState";
  value: BlockHeaderStateAmino;
}
export interface BlockHeaderStateSDKType {
  chain_id: bigint;
  latest_height: bigint;
  earliest_height: bigint;
  latest_block_hash: Uint8Array;
}
function createBaseBlockHeaderState(): BlockHeaderState {
  return {
    chainId: BigInt(0),
    latestHeight: BigInt(0),
    earliestHeight: BigInt(0),
    latestBlockHash: new Uint8Array()
  };
}
export const BlockHeaderState = {
  typeUrl: "/zetachain.zetacore.observer.BlockHeaderState",
  encode(message: BlockHeaderState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.latestHeight !== BigInt(0)) {
      writer.uint32(16).int64(message.latestHeight);
    }
    if (message.earliestHeight !== BigInt(0)) {
      writer.uint32(24).int64(message.earliestHeight);
    }
    if (message.latestBlockHash.length !== 0) {
      writer.uint32(34).bytes(message.latestBlockHash);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): BlockHeaderState {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockHeaderState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.latestHeight = reader.int64();
          break;
        case 3:
          message.earliestHeight = reader.int64();
          break;
        case 4:
          message.latestBlockHash = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<BlockHeaderState>): BlockHeaderState {
    const message = createBaseBlockHeaderState();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.latestHeight = object.latestHeight !== undefined && object.latestHeight !== null ? BigInt(object.latestHeight.toString()) : BigInt(0);
    message.earliestHeight = object.earliestHeight !== undefined && object.earliestHeight !== null ? BigInt(object.earliestHeight.toString()) : BigInt(0);
    message.latestBlockHash = object.latestBlockHash ?? new Uint8Array();
    return message;
  },
  fromAmino(object: BlockHeaderStateAmino): BlockHeaderState {
    const message = createBaseBlockHeaderState();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.latest_height !== undefined && object.latest_height !== null) {
      message.latestHeight = BigInt(object.latest_height);
    }
    if (object.earliest_height !== undefined && object.earliest_height !== null) {
      message.earliestHeight = BigInt(object.earliest_height);
    }
    if (object.latest_block_hash !== undefined && object.latest_block_hash !== null) {
      message.latestBlockHash = bytesFromBase64(object.latest_block_hash);
    }
    return message;
  },
  toAmino(message: BlockHeaderState): BlockHeaderStateAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.latest_height = message.latestHeight ? message.latestHeight.toString() : undefined;
    obj.earliest_height = message.earliestHeight ? message.earliestHeight.toString() : undefined;
    obj.latest_block_hash = message.latestBlockHash ? base64FromBytes(message.latestBlockHash) : undefined;
    return obj;
  },
  fromAminoMsg(object: BlockHeaderStateAminoMsg): BlockHeaderState {
    return BlockHeaderState.fromAmino(object.value);
  },
  fromProtoMsg(message: BlockHeaderStateProtoMsg): BlockHeaderState {
    return BlockHeaderState.decode(message.value);
  },
  toProto(message: BlockHeaderState): Uint8Array {
    return BlockHeaderState.encode(message).finish();
  },
  toProtoMsg(message: BlockHeaderState): BlockHeaderStateProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.BlockHeaderState",
      value: BlockHeaderState.encode(message).finish()
    };
  }
};