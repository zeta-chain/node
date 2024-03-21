import { CoinType, coinTypeFromJSON } from "../common/common";
import { BinaryReader, BinaryWriter } from "../../binary";
export interface ForeignCoins {
  /** string index = 1; */
  zrc20ContractAddress: string;
  asset: string;
  foreignChainId: bigint;
  decimals: number;
  name: string;
  symbol: string;
  coinType: CoinType;
  gasLimit: bigint;
  paused: boolean;
  liquidityCap: string;
}
export interface ForeignCoinsProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.ForeignCoins";
  value: Uint8Array;
}
export interface ForeignCoinsAmino {
  /** string index = 1; */
  zrc20_contract_address?: string;
  asset?: string;
  foreign_chain_id?: string;
  decimals?: number;
  name?: string;
  symbol?: string;
  coin_type?: CoinType;
  gas_limit?: string;
  paused?: boolean;
  liquidity_cap?: string;
}
export interface ForeignCoinsAminoMsg {
  type: "/zetachain.zetacore.fungible.ForeignCoins";
  value: ForeignCoinsAmino;
}
export interface ForeignCoinsSDKType {
  zrc20_contract_address: string;
  asset: string;
  foreign_chain_id: bigint;
  decimals: number;
  name: string;
  symbol: string;
  coin_type: CoinType;
  gas_limit: bigint;
  paused: boolean;
  liquidity_cap: string;
}
function createBaseForeignCoins(): ForeignCoins {
  return {
    zrc20ContractAddress: "",
    asset: "",
    foreignChainId: BigInt(0),
    decimals: 0,
    name: "",
    symbol: "",
    coinType: 0,
    gasLimit: BigInt(0),
    paused: false,
    liquidityCap: ""
  };
}
export const ForeignCoins = {
  typeUrl: "/zetachain.zetacore.fungible.ForeignCoins",
  encode(message: ForeignCoins, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.zrc20ContractAddress !== "") {
      writer.uint32(18).string(message.zrc20ContractAddress);
    }
    if (message.asset !== "") {
      writer.uint32(26).string(message.asset);
    }
    if (message.foreignChainId !== BigInt(0)) {
      writer.uint32(32).int64(message.foreignChainId);
    }
    if (message.decimals !== 0) {
      writer.uint32(40).uint32(message.decimals);
    }
    if (message.name !== "") {
      writer.uint32(50).string(message.name);
    }
    if (message.symbol !== "") {
      writer.uint32(58).string(message.symbol);
    }
    if (message.coinType !== 0) {
      writer.uint32(64).int32(message.coinType);
    }
    if (message.gasLimit !== BigInt(0)) {
      writer.uint32(72).uint64(message.gasLimit);
    }
    if (message.paused === true) {
      writer.uint32(80).bool(message.paused);
    }
    if (message.liquidityCap !== "") {
      writer.uint32(90).string(message.liquidityCap);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ForeignCoins {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseForeignCoins();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.zrc20ContractAddress = reader.string();
          break;
        case 3:
          message.asset = reader.string();
          break;
        case 4:
          message.foreignChainId = reader.int64();
          break;
        case 5:
          message.decimals = reader.uint32();
          break;
        case 6:
          message.name = reader.string();
          break;
        case 7:
          message.symbol = reader.string();
          break;
        case 8:
          message.coinType = (reader.int32() as any);
          break;
        case 9:
          message.gasLimit = reader.uint64();
          break;
        case 10:
          message.paused = reader.bool();
          break;
        case 11:
          message.liquidityCap = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ForeignCoins>): ForeignCoins {
    const message = createBaseForeignCoins();
    message.zrc20ContractAddress = object.zrc20ContractAddress ?? "";
    message.asset = object.asset ?? "";
    message.foreignChainId = object.foreignChainId !== undefined && object.foreignChainId !== null ? BigInt(object.foreignChainId.toString()) : BigInt(0);
    message.decimals = object.decimals ?? 0;
    message.name = object.name ?? "";
    message.symbol = object.symbol ?? "";
    message.coinType = object.coinType ?? 0;
    message.gasLimit = object.gasLimit !== undefined && object.gasLimit !== null ? BigInt(object.gasLimit.toString()) : BigInt(0);
    message.paused = object.paused ?? false;
    message.liquidityCap = object.liquidityCap ?? "";
    return message;
  },
  fromAmino(object: ForeignCoinsAmino): ForeignCoins {
    const message = createBaseForeignCoins();
    if (object.zrc20_contract_address !== undefined && object.zrc20_contract_address !== null) {
      message.zrc20ContractAddress = object.zrc20_contract_address;
    }
    if (object.asset !== undefined && object.asset !== null) {
      message.asset = object.asset;
    }
    if (object.foreign_chain_id !== undefined && object.foreign_chain_id !== null) {
      message.foreignChainId = BigInt(object.foreign_chain_id);
    }
    if (object.decimals !== undefined && object.decimals !== null) {
      message.decimals = object.decimals;
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    }
    if (object.symbol !== undefined && object.symbol !== null) {
      message.symbol = object.symbol;
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    if (object.gas_limit !== undefined && object.gas_limit !== null) {
      message.gasLimit = BigInt(object.gas_limit);
    }
    if (object.paused !== undefined && object.paused !== null) {
      message.paused = object.paused;
    }
    if (object.liquidity_cap !== undefined && object.liquidity_cap !== null) {
      message.liquidityCap = object.liquidity_cap;
    }
    return message;
  },
  toAmino(message: ForeignCoins): ForeignCoinsAmino {
    const obj: any = {};
    obj.zrc20_contract_address = message.zrc20ContractAddress;
    obj.asset = message.asset;
    obj.foreign_chain_id = message.foreignChainId ? message.foreignChainId.toString() : undefined;
    obj.decimals = message.decimals;
    obj.name = message.name;
    obj.symbol = message.symbol;
    obj.coin_type = message.coinType;
    obj.gas_limit = message.gasLimit ? message.gasLimit.toString() : undefined;
    obj.paused = message.paused;
    obj.liquidity_cap = message.liquidityCap;
    return obj;
  },
  fromAminoMsg(object: ForeignCoinsAminoMsg): ForeignCoins {
    return ForeignCoins.fromAmino(object.value);
  },
  fromProtoMsg(message: ForeignCoinsProtoMsg): ForeignCoins {
    return ForeignCoins.decode(message.value);
  },
  toProto(message: ForeignCoins): Uint8Array {
    return ForeignCoins.encode(message).finish();
  },
  toProtoMsg(message: ForeignCoins): ForeignCoinsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.ForeignCoins",
      value: ForeignCoins.encode(message).finish()
    };
  }
};