import { Params, ParamsAmino, ParamsSDKType } from "./params";
import { ForeignCoins, ForeignCoinsAmino, ForeignCoinsSDKType } from "./foreign_coins";
import { SystemContract, SystemContractAmino, SystemContractSDKType } from "./system_contract";
import { BinaryReader, BinaryWriter } from "../../binary";
/** GenesisState defines the fungible module's genesis state. */
export interface GenesisState {
  params: Params;
  foreignCoinsList: ForeignCoins[];
  systemContract?: SystemContract;
}
export interface GenesisStateProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the fungible module's genesis state. */
export interface GenesisStateAmino {
  params?: ParamsAmino;
  foreignCoinsList?: ForeignCoinsAmino[];
  systemContract?: SystemContractAmino;
}
export interface GenesisStateAminoMsg {
  type: "/zetachain.zetacore.fungible.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the fungible module's genesis state. */
export interface GenesisStateSDKType {
  params: ParamsSDKType;
  foreignCoinsList: ForeignCoinsSDKType[];
  systemContract?: SystemContractSDKType;
}
function createBaseGenesisState(): GenesisState {
  return {
    params: Params.fromPartial({}),
    foreignCoinsList: [],
    systemContract: undefined
  };
}
export const GenesisState = {
  typeUrl: "/zetachain.zetacore.fungible.GenesisState",
  encode(message: GenesisState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.foreignCoinsList) {
      ForeignCoins.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.systemContract !== undefined) {
      SystemContract.encode(message.systemContract, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.foreignCoinsList.push(ForeignCoins.decode(reader, reader.uint32()));
          break;
        case 3:
          message.systemContract = SystemContract.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    message.foreignCoinsList = object.foreignCoinsList?.map(e => ForeignCoins.fromPartial(e)) || [];
    message.systemContract = object.systemContract !== undefined && object.systemContract !== null ? SystemContract.fromPartial(object.systemContract) : undefined;
    return message;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    const message = createBaseGenesisState();
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromAmino(object.params);
    }
    message.foreignCoinsList = object.foreignCoinsList?.map(e => ForeignCoins.fromAmino(e)) || [];
    if (object.systemContract !== undefined && object.systemContract !== null) {
      message.systemContract = SystemContract.fromAmino(object.systemContract);
    }
    return message;
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    if (message.foreignCoinsList) {
      obj.foreignCoinsList = message.foreignCoinsList.map(e => e ? ForeignCoins.toAmino(e) : undefined);
    } else {
      obj.foreignCoinsList = [];
    }
    obj.systemContract = message.systemContract ? SystemContract.toAmino(message.systemContract) : undefined;
    return obj;
  },
  fromAminoMsg(object: GenesisStateAminoMsg): GenesisState {
    return GenesisState.fromAmino(object.value);
  },
  fromProtoMsg(message: GenesisStateProtoMsg): GenesisState {
    return GenesisState.decode(message.value);
  },
  toProto(message: GenesisState): Uint8Array {
    return GenesisState.encode(message).finish();
  },
  toProtoMsg(message: GenesisState): GenesisStateProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};