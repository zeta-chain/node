import { BinaryReader, BinaryWriter } from "../../binary";
export interface GasPrice {
  creator: string;
  index: string;
  chainId: bigint;
  signers: string[];
  blockNums: bigint[];
  prices: bigint[];
  medianIndex: bigint;
}
export interface GasPriceProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.GasPrice";
  value: Uint8Array;
}
export interface GasPriceAmino {
  creator?: string;
  index?: string;
  chain_id?: string;
  signers?: string[];
  block_nums?: string[];
  prices?: string[];
  median_index?: string;
}
export interface GasPriceAminoMsg {
  type: "/zetachain.zetacore.crosschain.GasPrice";
  value: GasPriceAmino;
}
export interface GasPriceSDKType {
  creator: string;
  index: string;
  chain_id: bigint;
  signers: string[];
  block_nums: bigint[];
  prices: bigint[];
  median_index: bigint;
}
function createBaseGasPrice(): GasPrice {
  return {
    creator: "",
    index: "",
    chainId: BigInt(0),
    signers: [],
    blockNums: [],
    prices: [],
    medianIndex: BigInt(0)
  };
}
export const GasPrice = {
  typeUrl: "/zetachain.zetacore.crosschain.GasPrice",
  encode(message: GasPrice, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.index !== "") {
      writer.uint32(18).string(message.index);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(24).int64(message.chainId);
    }
    for (const v of message.signers) {
      writer.uint32(34).string(v!);
    }
    writer.uint32(42).fork();
    for (const v of message.blockNums) {
      writer.uint64(v);
    }
    writer.ldelim();
    writer.uint32(50).fork();
    for (const v of message.prices) {
      writer.uint64(v);
    }
    writer.ldelim();
    if (message.medianIndex !== BigInt(0)) {
      writer.uint32(56).uint64(message.medianIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): GasPrice {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGasPrice();
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
          message.signers.push(reader.string());
          break;
        case 5:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.blockNums.push(reader.uint64());
            }
          } else {
            message.blockNums.push(reader.uint64());
          }
          break;
        case 6:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.prices.push(reader.uint64());
            }
          } else {
            message.prices.push(reader.uint64());
          }
          break;
        case 7:
          message.medianIndex = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<GasPrice>): GasPrice {
    const message = createBaseGasPrice();
    message.creator = object.creator ?? "";
    message.index = object.index ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.signers = object.signers?.map(e => e) || [];
    message.blockNums = object.blockNums?.map(e => BigInt(e.toString())) || [];
    message.prices = object.prices?.map(e => BigInt(e.toString())) || [];
    message.medianIndex = object.medianIndex !== undefined && object.medianIndex !== null ? BigInt(object.medianIndex.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: GasPriceAmino): GasPrice {
    const message = createBaseGasPrice();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    message.signers = object.signers?.map(e => e) || [];
    message.blockNums = object.block_nums?.map(e => BigInt(e)) || [];
    message.prices = object.prices?.map(e => BigInt(e)) || [];
    if (object.median_index !== undefined && object.median_index !== null) {
      message.medianIndex = BigInt(object.median_index);
    }
    return message;
  },
  toAmino(message: GasPrice): GasPriceAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.index = message.index;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    if (message.signers) {
      obj.signers = message.signers.map(e => e);
    } else {
      obj.signers = [];
    }
    if (message.blockNums) {
      obj.block_nums = message.blockNums.map(e => e.toString());
    } else {
      obj.block_nums = [];
    }
    if (message.prices) {
      obj.prices = message.prices.map(e => e.toString());
    } else {
      obj.prices = [];
    }
    obj.median_index = message.medianIndex ? message.medianIndex.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: GasPriceAminoMsg): GasPrice {
    return GasPrice.fromAmino(object.value);
  },
  fromProtoMsg(message: GasPriceProtoMsg): GasPrice {
    return GasPrice.decode(message.value);
  },
  toProto(message: GasPrice): Uint8Array {
    return GasPrice.encode(message).finish();
  },
  toProtoMsg(message: GasPrice): GasPriceProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.GasPrice",
      value: GasPrice.encode(message).finish()
    };
  }
};