import { BinaryReader, BinaryWriter } from "../../binary";
export interface TxHashList {
  txHash: string;
  txSigner: string;
  proved: boolean;
}
export interface TxHashListProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.TxHashList";
  value: Uint8Array;
}
export interface TxHashListAmino {
  tx_hash?: string;
  tx_signer?: string;
  proved?: boolean;
}
export interface TxHashListAminoMsg {
  type: "/zetachain.zetacore.crosschain.TxHashList";
  value: TxHashListAmino;
}
export interface TxHashListSDKType {
  tx_hash: string;
  tx_signer: string;
  proved: boolean;
}
export interface OutTxTracker {
  /** format: "chain-nonce" */
  index: string;
  chainId: bigint;
  nonce: bigint;
  hashList: TxHashList[];
}
export interface OutTxTrackerProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.OutTxTracker";
  value: Uint8Array;
}
export interface OutTxTrackerAmino {
  /** format: "chain-nonce" */
  index?: string;
  chain_id?: string;
  nonce?: string;
  hash_list?: TxHashListAmino[];
}
export interface OutTxTrackerAminoMsg {
  type: "/zetachain.zetacore.crosschain.OutTxTracker";
  value: OutTxTrackerAmino;
}
export interface OutTxTrackerSDKType {
  index: string;
  chain_id: bigint;
  nonce: bigint;
  hash_list: TxHashListSDKType[];
}
function createBaseTxHashList(): TxHashList {
  return {
    txHash: "",
    txSigner: "",
    proved: false
  };
}
export const TxHashList = {
  typeUrl: "/zetachain.zetacore.crosschain.TxHashList",
  encode(message: TxHashList, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.txHash !== "") {
      writer.uint32(10).string(message.txHash);
    }
    if (message.txSigner !== "") {
      writer.uint32(18).string(message.txSigner);
    }
    if (message.proved === true) {
      writer.uint32(24).bool(message.proved);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): TxHashList {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTxHashList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.txHash = reader.string();
          break;
        case 2:
          message.txSigner = reader.string();
          break;
        case 3:
          message.proved = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<TxHashList>): TxHashList {
    const message = createBaseTxHashList();
    message.txHash = object.txHash ?? "";
    message.txSigner = object.txSigner ?? "";
    message.proved = object.proved ?? false;
    return message;
  },
  fromAmino(object: TxHashListAmino): TxHashList {
    const message = createBaseTxHashList();
    if (object.tx_hash !== undefined && object.tx_hash !== null) {
      message.txHash = object.tx_hash;
    }
    if (object.tx_signer !== undefined && object.tx_signer !== null) {
      message.txSigner = object.tx_signer;
    }
    if (object.proved !== undefined && object.proved !== null) {
      message.proved = object.proved;
    }
    return message;
  },
  toAmino(message: TxHashList): TxHashListAmino {
    const obj: any = {};
    obj.tx_hash = message.txHash;
    obj.tx_signer = message.txSigner;
    obj.proved = message.proved;
    return obj;
  },
  fromAminoMsg(object: TxHashListAminoMsg): TxHashList {
    return TxHashList.fromAmino(object.value);
  },
  fromProtoMsg(message: TxHashListProtoMsg): TxHashList {
    return TxHashList.decode(message.value);
  },
  toProto(message: TxHashList): Uint8Array {
    return TxHashList.encode(message).finish();
  },
  toProtoMsg(message: TxHashList): TxHashListProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.TxHashList",
      value: TxHashList.encode(message).finish()
    };
  }
};
function createBaseOutTxTracker(): OutTxTracker {
  return {
    index: "",
    chainId: BigInt(0),
    nonce: BigInt(0),
    hashList: []
  };
}
export const OutTxTracker = {
  typeUrl: "/zetachain.zetacore.crosschain.OutTxTracker",
  encode(message: OutTxTracker, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.nonce !== BigInt(0)) {
      writer.uint32(24).uint64(message.nonce);
    }
    for (const v of message.hashList) {
      TxHashList.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OutTxTracker {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOutTxTracker();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.nonce = reader.uint64();
          break;
        case 4:
          message.hashList.push(TxHashList.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OutTxTracker>): OutTxTracker {
    const message = createBaseOutTxTracker();
    message.index = object.index ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.nonce = object.nonce !== undefined && object.nonce !== null ? BigInt(object.nonce.toString()) : BigInt(0);
    message.hashList = object.hashList?.map(e => TxHashList.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: OutTxTrackerAmino): OutTxTracker {
    const message = createBaseOutTxTracker();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = BigInt(object.nonce);
    }
    message.hashList = object.hash_list?.map(e => TxHashList.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: OutTxTracker): OutTxTrackerAmino {
    const obj: any = {};
    obj.index = message.index;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.nonce = message.nonce ? message.nonce.toString() : undefined;
    if (message.hashList) {
      obj.hash_list = message.hashList.map(e => e ? TxHashList.toAmino(e) : undefined);
    } else {
      obj.hash_list = [];
    }
    return obj;
  },
  fromAminoMsg(object: OutTxTrackerAminoMsg): OutTxTracker {
    return OutTxTracker.fromAmino(object.value);
  },
  fromProtoMsg(message: OutTxTrackerProtoMsg): OutTxTracker {
    return OutTxTracker.decode(message.value);
  },
  toProto(message: OutTxTracker): Uint8Array {
    return OutTxTracker.encode(message).finish();
  },
  toProtoMsg(message: OutTxTracker): OutTxTrackerProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.OutTxTracker",
      value: OutTxTracker.encode(message).finish()
    };
  }
};