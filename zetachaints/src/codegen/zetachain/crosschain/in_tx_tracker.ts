import { CoinType, coinTypeFromJSON } from "../common/common";
import { BinaryReader, BinaryWriter } from "../../binary";
export interface InTxTracker {
  chainId: bigint;
  txHash: string;
  coinType: CoinType;
}
export interface InTxTrackerProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.InTxTracker";
  value: Uint8Array;
}
export interface InTxTrackerAmino {
  chain_id?: string;
  tx_hash?: string;
  coin_type?: CoinType;
}
export interface InTxTrackerAminoMsg {
  type: "/zetachain.zetacore.crosschain.InTxTracker";
  value: InTxTrackerAmino;
}
export interface InTxTrackerSDKType {
  chain_id: bigint;
  tx_hash: string;
  coin_type: CoinType;
}
function createBaseInTxTracker(): InTxTracker {
  return {
    chainId: BigInt(0),
    txHash: "",
    coinType: 0
  };
}
export const InTxTracker = {
  typeUrl: "/zetachain.zetacore.crosschain.InTxTracker",
  encode(message: InTxTracker, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.txHash !== "") {
      writer.uint32(18).string(message.txHash);
    }
    if (message.coinType !== 0) {
      writer.uint32(24).int32(message.coinType);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): InTxTracker {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInTxTracker();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.txHash = reader.string();
          break;
        case 3:
          message.coinType = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<InTxTracker>): InTxTracker {
    const message = createBaseInTxTracker();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.txHash = object.txHash ?? "";
    message.coinType = object.coinType ?? 0;
    return message;
  },
  fromAmino(object: InTxTrackerAmino): InTxTracker {
    const message = createBaseInTxTracker();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.tx_hash !== undefined && object.tx_hash !== null) {
      message.txHash = object.tx_hash;
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    return message;
  },
  toAmino(message: InTxTracker): InTxTrackerAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.tx_hash = message.txHash;
    obj.coin_type = message.coinType;
    return obj;
  },
  fromAminoMsg(object: InTxTrackerAminoMsg): InTxTracker {
    return InTxTracker.fromAmino(object.value);
  },
  fromProtoMsg(message: InTxTrackerProtoMsg): InTxTracker {
    return InTxTracker.decode(message.value);
  },
  toProto(message: InTxTracker): Uint8Array {
    return InTxTracker.encode(message).finish();
  },
  toProtoMsg(message: InTxTracker): InTxTrackerProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.InTxTracker",
      value: InTxTracker.encode(message).finish()
    };
  }
};