import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Params, ParamsAmino, ParamsSDKType } from "./params";
import { ForeignCoins, ForeignCoinsAmino, ForeignCoinsSDKType } from "./foreign_coins";
import { SystemContract, SystemContractAmino, SystemContractSDKType } from "./system_contract";
import { BinaryReader, BinaryWriter } from "../../binary";
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {}
export interface QueryParamsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryParamsRequest";
  value: Uint8Array;
}
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequestAmino {}
export interface QueryParamsRequestAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryParamsRequest";
  value: QueryParamsRequestAmino;
}
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params;
}
export interface QueryParamsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryParamsResponse";
  value: Uint8Array;
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponseAmino {
  /** params holds all the parameters of this module. */
  params?: ParamsAmino;
}
export interface QueryParamsResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryParamsResponse";
  value: QueryParamsResponseAmino;
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponseSDKType {
  params: ParamsSDKType;
}
export interface QueryGetForeignCoinsRequest {
  index: string;
}
export interface QueryGetForeignCoinsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetForeignCoinsRequest";
  value: Uint8Array;
}
export interface QueryGetForeignCoinsRequestAmino {
  index?: string;
}
export interface QueryGetForeignCoinsRequestAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryGetForeignCoinsRequest";
  value: QueryGetForeignCoinsRequestAmino;
}
export interface QueryGetForeignCoinsRequestSDKType {
  index: string;
}
export interface QueryGetForeignCoinsResponse {
  foreignCoins: ForeignCoins;
}
export interface QueryGetForeignCoinsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetForeignCoinsResponse";
  value: Uint8Array;
}
export interface QueryGetForeignCoinsResponseAmino {
  foreignCoins?: ForeignCoinsAmino;
}
export interface QueryGetForeignCoinsResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryGetForeignCoinsResponse";
  value: QueryGetForeignCoinsResponseAmino;
}
export interface QueryGetForeignCoinsResponseSDKType {
  foreignCoins: ForeignCoinsSDKType;
}
export interface QueryAllForeignCoinsRequest {
  pagination?: PageRequest;
}
export interface QueryAllForeignCoinsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryAllForeignCoinsRequest";
  value: Uint8Array;
}
export interface QueryAllForeignCoinsRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllForeignCoinsRequestAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryAllForeignCoinsRequest";
  value: QueryAllForeignCoinsRequestAmino;
}
export interface QueryAllForeignCoinsRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllForeignCoinsResponse {
  foreignCoins: ForeignCoins[];
  pagination?: PageResponse;
}
export interface QueryAllForeignCoinsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryAllForeignCoinsResponse";
  value: Uint8Array;
}
export interface QueryAllForeignCoinsResponseAmino {
  foreignCoins?: ForeignCoinsAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllForeignCoinsResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryAllForeignCoinsResponse";
  value: QueryAllForeignCoinsResponseAmino;
}
export interface QueryAllForeignCoinsResponseSDKType {
  foreignCoins: ForeignCoinsSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryGetSystemContractRequest {}
export interface QueryGetSystemContractRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetSystemContractRequest";
  value: Uint8Array;
}
export interface QueryGetSystemContractRequestAmino {}
export interface QueryGetSystemContractRequestAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryGetSystemContractRequest";
  value: QueryGetSystemContractRequestAmino;
}
export interface QueryGetSystemContractRequestSDKType {}
export interface QueryGetSystemContractResponse {
  SystemContract: SystemContract;
}
export interface QueryGetSystemContractResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetSystemContractResponse";
  value: Uint8Array;
}
export interface QueryGetSystemContractResponseAmino {
  SystemContract?: SystemContractAmino;
}
export interface QueryGetSystemContractResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryGetSystemContractResponse";
  value: QueryGetSystemContractResponseAmino;
}
export interface QueryGetSystemContractResponseSDKType {
  SystemContract: SystemContractSDKType;
}
export interface QueryGetGasStabilityPoolAddress {}
export interface QueryGetGasStabilityPoolAddressProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolAddress";
  value: Uint8Array;
}
export interface QueryGetGasStabilityPoolAddressAmino {}
export interface QueryGetGasStabilityPoolAddressAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolAddress";
  value: QueryGetGasStabilityPoolAddressAmino;
}
export interface QueryGetGasStabilityPoolAddressSDKType {}
export interface QueryGetGasStabilityPoolAddressResponse {
  cosmosAddress: string;
  evmAddress: string;
}
export interface QueryGetGasStabilityPoolAddressResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolAddressResponse";
  value: Uint8Array;
}
export interface QueryGetGasStabilityPoolAddressResponseAmino {
  cosmos_address?: string;
  evm_address?: string;
}
export interface QueryGetGasStabilityPoolAddressResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolAddressResponse";
  value: QueryGetGasStabilityPoolAddressResponseAmino;
}
export interface QueryGetGasStabilityPoolAddressResponseSDKType {
  cosmos_address: string;
  evm_address: string;
}
export interface QueryGetGasStabilityPoolBalance {
  chainId: bigint;
}
export interface QueryGetGasStabilityPoolBalanceProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolBalance";
  value: Uint8Array;
}
export interface QueryGetGasStabilityPoolBalanceAmino {
  chain_id?: string;
}
export interface QueryGetGasStabilityPoolBalanceAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolBalance";
  value: QueryGetGasStabilityPoolBalanceAmino;
}
export interface QueryGetGasStabilityPoolBalanceSDKType {
  chain_id: bigint;
}
export interface QueryGetGasStabilityPoolBalanceResponse {
  balance: string;
}
export interface QueryGetGasStabilityPoolBalanceResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolBalanceResponse";
  value: Uint8Array;
}
export interface QueryGetGasStabilityPoolBalanceResponseAmino {
  balance?: string;
}
export interface QueryGetGasStabilityPoolBalanceResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolBalanceResponse";
  value: QueryGetGasStabilityPoolBalanceResponseAmino;
}
export interface QueryGetGasStabilityPoolBalanceResponseSDKType {
  balance: string;
}
export interface QueryAllGasStabilityPoolBalance {}
export interface QueryAllGasStabilityPoolBalanceProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryAllGasStabilityPoolBalance";
  value: Uint8Array;
}
export interface QueryAllGasStabilityPoolBalanceAmino {}
export interface QueryAllGasStabilityPoolBalanceAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryAllGasStabilityPoolBalance";
  value: QueryAllGasStabilityPoolBalanceAmino;
}
export interface QueryAllGasStabilityPoolBalanceSDKType {}
export interface QueryAllGasStabilityPoolBalanceResponse {
  balances: QueryAllGasStabilityPoolBalanceResponse_Balance[];
}
export interface QueryAllGasStabilityPoolBalanceResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryAllGasStabilityPoolBalanceResponse";
  value: Uint8Array;
}
export interface QueryAllGasStabilityPoolBalanceResponseAmino {
  balances?: QueryAllGasStabilityPoolBalanceResponse_BalanceAmino[];
}
export interface QueryAllGasStabilityPoolBalanceResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryAllGasStabilityPoolBalanceResponse";
  value: QueryAllGasStabilityPoolBalanceResponseAmino;
}
export interface QueryAllGasStabilityPoolBalanceResponseSDKType {
  balances: QueryAllGasStabilityPoolBalanceResponse_BalanceSDKType[];
}
export interface QueryAllGasStabilityPoolBalanceResponse_Balance {
  chainId: bigint;
  balance: string;
}
export interface QueryAllGasStabilityPoolBalanceResponse_BalanceProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.Balance";
  value: Uint8Array;
}
export interface QueryAllGasStabilityPoolBalanceResponse_BalanceAmino {
  chain_id?: string;
  balance?: string;
}
export interface QueryAllGasStabilityPoolBalanceResponse_BalanceAminoMsg {
  type: "/zetachain.zetacore.fungible.Balance";
  value: QueryAllGasStabilityPoolBalanceResponse_BalanceAmino;
}
export interface QueryAllGasStabilityPoolBalanceResponse_BalanceSDKType {
  chain_id: bigint;
  balance: string;
}
export interface QueryCodeHashRequest {
  address: string;
}
export interface QueryCodeHashRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryCodeHashRequest";
  value: Uint8Array;
}
export interface QueryCodeHashRequestAmino {
  address?: string;
}
export interface QueryCodeHashRequestAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryCodeHashRequest";
  value: QueryCodeHashRequestAmino;
}
export interface QueryCodeHashRequestSDKType {
  address: string;
}
export interface QueryCodeHashResponse {
  codeHash: string;
}
export interface QueryCodeHashResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.QueryCodeHashResponse";
  value: Uint8Array;
}
export interface QueryCodeHashResponseAmino {
  code_hash?: string;
}
export interface QueryCodeHashResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.QueryCodeHashResponse";
  value: QueryCodeHashResponseAmino;
}
export interface QueryCodeHashResponseSDKType {
  code_hash: string;
}
function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}
export const QueryParamsRequest = {
  typeUrl: "/zetachain.zetacore.fungible.QueryParamsRequest",
  encode(_: QueryParamsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<QueryParamsRequest>): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
  fromAmino(_: QueryParamsRequestAmino): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
  toAmino(_: QueryParamsRequest): QueryParamsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryParamsRequestAminoMsg): QueryParamsRequest {
    return QueryParamsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryParamsRequestProtoMsg): QueryParamsRequest {
    return QueryParamsRequest.decode(message.value);
  },
  toProto(message: QueryParamsRequest): Uint8Array {
    return QueryParamsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryParamsRequest): QueryParamsRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryParamsRequest",
      value: QueryParamsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryParamsResponse(): QueryParamsResponse {
  return {
    params: Params.fromPartial({})
  };
}
export const QueryParamsResponse = {
  typeUrl: "/zetachain.zetacore.fungible.QueryParamsResponse",
  encode(message: QueryParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryParamsResponse>): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: QueryParamsResponseAmino): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: QueryParamsResponse): QueryParamsResponseAmino {
    const obj: any = {};
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryParamsResponseAminoMsg): QueryParamsResponse {
    return QueryParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryParamsResponseProtoMsg): QueryParamsResponse {
    return QueryParamsResponse.decode(message.value);
  },
  toProto(message: QueryParamsResponse): Uint8Array {
    return QueryParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryParamsResponse): QueryParamsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryParamsResponse",
      value: QueryParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetForeignCoinsRequest(): QueryGetForeignCoinsRequest {
  return {
    index: ""
  };
}
export const QueryGetForeignCoinsRequest = {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetForeignCoinsRequest",
  encode(message: QueryGetForeignCoinsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetForeignCoinsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetForeignCoinsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetForeignCoinsRequest>): QueryGetForeignCoinsRequest {
    const message = createBaseQueryGetForeignCoinsRequest();
    message.index = object.index ?? "";
    return message;
  },
  fromAmino(object: QueryGetForeignCoinsRequestAmino): QueryGetForeignCoinsRequest {
    const message = createBaseQueryGetForeignCoinsRequest();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    return message;
  },
  toAmino(message: QueryGetForeignCoinsRequest): QueryGetForeignCoinsRequestAmino {
    const obj: any = {};
    obj.index = message.index;
    return obj;
  },
  fromAminoMsg(object: QueryGetForeignCoinsRequestAminoMsg): QueryGetForeignCoinsRequest {
    return QueryGetForeignCoinsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetForeignCoinsRequestProtoMsg): QueryGetForeignCoinsRequest {
    return QueryGetForeignCoinsRequest.decode(message.value);
  },
  toProto(message: QueryGetForeignCoinsRequest): Uint8Array {
    return QueryGetForeignCoinsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetForeignCoinsRequest): QueryGetForeignCoinsRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryGetForeignCoinsRequest",
      value: QueryGetForeignCoinsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetForeignCoinsResponse(): QueryGetForeignCoinsResponse {
  return {
    foreignCoins: ForeignCoins.fromPartial({})
  };
}
export const QueryGetForeignCoinsResponse = {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetForeignCoinsResponse",
  encode(message: QueryGetForeignCoinsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.foreignCoins !== undefined) {
      ForeignCoins.encode(message.foreignCoins, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetForeignCoinsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetForeignCoinsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.foreignCoins = ForeignCoins.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetForeignCoinsResponse>): QueryGetForeignCoinsResponse {
    const message = createBaseQueryGetForeignCoinsResponse();
    message.foreignCoins = object.foreignCoins !== undefined && object.foreignCoins !== null ? ForeignCoins.fromPartial(object.foreignCoins) : undefined;
    return message;
  },
  fromAmino(object: QueryGetForeignCoinsResponseAmino): QueryGetForeignCoinsResponse {
    const message = createBaseQueryGetForeignCoinsResponse();
    if (object.foreignCoins !== undefined && object.foreignCoins !== null) {
      message.foreignCoins = ForeignCoins.fromAmino(object.foreignCoins);
    }
    return message;
  },
  toAmino(message: QueryGetForeignCoinsResponse): QueryGetForeignCoinsResponseAmino {
    const obj: any = {};
    obj.foreignCoins = message.foreignCoins ? ForeignCoins.toAmino(message.foreignCoins) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetForeignCoinsResponseAminoMsg): QueryGetForeignCoinsResponse {
    return QueryGetForeignCoinsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetForeignCoinsResponseProtoMsg): QueryGetForeignCoinsResponse {
    return QueryGetForeignCoinsResponse.decode(message.value);
  },
  toProto(message: QueryGetForeignCoinsResponse): Uint8Array {
    return QueryGetForeignCoinsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetForeignCoinsResponse): QueryGetForeignCoinsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryGetForeignCoinsResponse",
      value: QueryGetForeignCoinsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllForeignCoinsRequest(): QueryAllForeignCoinsRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllForeignCoinsRequest = {
  typeUrl: "/zetachain.zetacore.fungible.QueryAllForeignCoinsRequest",
  encode(message: QueryAllForeignCoinsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllForeignCoinsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllForeignCoinsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllForeignCoinsRequest>): QueryAllForeignCoinsRequest {
    const message = createBaseQueryAllForeignCoinsRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllForeignCoinsRequestAmino): QueryAllForeignCoinsRequest {
    const message = createBaseQueryAllForeignCoinsRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllForeignCoinsRequest): QueryAllForeignCoinsRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllForeignCoinsRequestAminoMsg): QueryAllForeignCoinsRequest {
    return QueryAllForeignCoinsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllForeignCoinsRequestProtoMsg): QueryAllForeignCoinsRequest {
    return QueryAllForeignCoinsRequest.decode(message.value);
  },
  toProto(message: QueryAllForeignCoinsRequest): Uint8Array {
    return QueryAllForeignCoinsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllForeignCoinsRequest): QueryAllForeignCoinsRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryAllForeignCoinsRequest",
      value: QueryAllForeignCoinsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllForeignCoinsResponse(): QueryAllForeignCoinsResponse {
  return {
    foreignCoins: [],
    pagination: undefined
  };
}
export const QueryAllForeignCoinsResponse = {
  typeUrl: "/zetachain.zetacore.fungible.QueryAllForeignCoinsResponse",
  encode(message: QueryAllForeignCoinsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.foreignCoins) {
      ForeignCoins.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllForeignCoinsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllForeignCoinsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.foreignCoins.push(ForeignCoins.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllForeignCoinsResponse>): QueryAllForeignCoinsResponse {
    const message = createBaseQueryAllForeignCoinsResponse();
    message.foreignCoins = object.foreignCoins?.map(e => ForeignCoins.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllForeignCoinsResponseAmino): QueryAllForeignCoinsResponse {
    const message = createBaseQueryAllForeignCoinsResponse();
    message.foreignCoins = object.foreignCoins?.map(e => ForeignCoins.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllForeignCoinsResponse): QueryAllForeignCoinsResponseAmino {
    const obj: any = {};
    if (message.foreignCoins) {
      obj.foreignCoins = message.foreignCoins.map(e => e ? ForeignCoins.toAmino(e) : undefined);
    } else {
      obj.foreignCoins = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllForeignCoinsResponseAminoMsg): QueryAllForeignCoinsResponse {
    return QueryAllForeignCoinsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllForeignCoinsResponseProtoMsg): QueryAllForeignCoinsResponse {
    return QueryAllForeignCoinsResponse.decode(message.value);
  },
  toProto(message: QueryAllForeignCoinsResponse): Uint8Array {
    return QueryAllForeignCoinsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllForeignCoinsResponse): QueryAllForeignCoinsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryAllForeignCoinsResponse",
      value: QueryAllForeignCoinsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetSystemContractRequest(): QueryGetSystemContractRequest {
  return {};
}
export const QueryGetSystemContractRequest = {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetSystemContractRequest",
  encode(_: QueryGetSystemContractRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetSystemContractRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSystemContractRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<QueryGetSystemContractRequest>): QueryGetSystemContractRequest {
    const message = createBaseQueryGetSystemContractRequest();
    return message;
  },
  fromAmino(_: QueryGetSystemContractRequestAmino): QueryGetSystemContractRequest {
    const message = createBaseQueryGetSystemContractRequest();
    return message;
  },
  toAmino(_: QueryGetSystemContractRequest): QueryGetSystemContractRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryGetSystemContractRequestAminoMsg): QueryGetSystemContractRequest {
    return QueryGetSystemContractRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetSystemContractRequestProtoMsg): QueryGetSystemContractRequest {
    return QueryGetSystemContractRequest.decode(message.value);
  },
  toProto(message: QueryGetSystemContractRequest): Uint8Array {
    return QueryGetSystemContractRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetSystemContractRequest): QueryGetSystemContractRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryGetSystemContractRequest",
      value: QueryGetSystemContractRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetSystemContractResponse(): QueryGetSystemContractResponse {
  return {
    SystemContract: SystemContract.fromPartial({})
  };
}
export const QueryGetSystemContractResponse = {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetSystemContractResponse",
  encode(message: QueryGetSystemContractResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.SystemContract !== undefined) {
      SystemContract.encode(message.SystemContract, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetSystemContractResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSystemContractResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.SystemContract = SystemContract.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetSystemContractResponse>): QueryGetSystemContractResponse {
    const message = createBaseQueryGetSystemContractResponse();
    message.SystemContract = object.SystemContract !== undefined && object.SystemContract !== null ? SystemContract.fromPartial(object.SystemContract) : undefined;
    return message;
  },
  fromAmino(object: QueryGetSystemContractResponseAmino): QueryGetSystemContractResponse {
    const message = createBaseQueryGetSystemContractResponse();
    if (object.SystemContract !== undefined && object.SystemContract !== null) {
      message.SystemContract = SystemContract.fromAmino(object.SystemContract);
    }
    return message;
  },
  toAmino(message: QueryGetSystemContractResponse): QueryGetSystemContractResponseAmino {
    const obj: any = {};
    obj.SystemContract = message.SystemContract ? SystemContract.toAmino(message.SystemContract) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetSystemContractResponseAminoMsg): QueryGetSystemContractResponse {
    return QueryGetSystemContractResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetSystemContractResponseProtoMsg): QueryGetSystemContractResponse {
    return QueryGetSystemContractResponse.decode(message.value);
  },
  toProto(message: QueryGetSystemContractResponse): Uint8Array {
    return QueryGetSystemContractResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetSystemContractResponse): QueryGetSystemContractResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryGetSystemContractResponse",
      value: QueryGetSystemContractResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetGasStabilityPoolAddress(): QueryGetGasStabilityPoolAddress {
  return {};
}
export const QueryGetGasStabilityPoolAddress = {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolAddress",
  encode(_: QueryGetGasStabilityPoolAddress, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetGasStabilityPoolAddress {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGasStabilityPoolAddress();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<QueryGetGasStabilityPoolAddress>): QueryGetGasStabilityPoolAddress {
    const message = createBaseQueryGetGasStabilityPoolAddress();
    return message;
  },
  fromAmino(_: QueryGetGasStabilityPoolAddressAmino): QueryGetGasStabilityPoolAddress {
    const message = createBaseQueryGetGasStabilityPoolAddress();
    return message;
  },
  toAmino(_: QueryGetGasStabilityPoolAddress): QueryGetGasStabilityPoolAddressAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryGetGasStabilityPoolAddressAminoMsg): QueryGetGasStabilityPoolAddress {
    return QueryGetGasStabilityPoolAddress.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetGasStabilityPoolAddressProtoMsg): QueryGetGasStabilityPoolAddress {
    return QueryGetGasStabilityPoolAddress.decode(message.value);
  },
  toProto(message: QueryGetGasStabilityPoolAddress): Uint8Array {
    return QueryGetGasStabilityPoolAddress.encode(message).finish();
  },
  toProtoMsg(message: QueryGetGasStabilityPoolAddress): QueryGetGasStabilityPoolAddressProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolAddress",
      value: QueryGetGasStabilityPoolAddress.encode(message).finish()
    };
  }
};
function createBaseQueryGetGasStabilityPoolAddressResponse(): QueryGetGasStabilityPoolAddressResponse {
  return {
    cosmosAddress: "",
    evmAddress: ""
  };
}
export const QueryGetGasStabilityPoolAddressResponse = {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolAddressResponse",
  encode(message: QueryGetGasStabilityPoolAddressResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.cosmosAddress !== "") {
      writer.uint32(10).string(message.cosmosAddress);
    }
    if (message.evmAddress !== "") {
      writer.uint32(18).string(message.evmAddress);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetGasStabilityPoolAddressResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGasStabilityPoolAddressResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.cosmosAddress = reader.string();
          break;
        case 2:
          message.evmAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetGasStabilityPoolAddressResponse>): QueryGetGasStabilityPoolAddressResponse {
    const message = createBaseQueryGetGasStabilityPoolAddressResponse();
    message.cosmosAddress = object.cosmosAddress ?? "";
    message.evmAddress = object.evmAddress ?? "";
    return message;
  },
  fromAmino(object: QueryGetGasStabilityPoolAddressResponseAmino): QueryGetGasStabilityPoolAddressResponse {
    const message = createBaseQueryGetGasStabilityPoolAddressResponse();
    if (object.cosmos_address !== undefined && object.cosmos_address !== null) {
      message.cosmosAddress = object.cosmos_address;
    }
    if (object.evm_address !== undefined && object.evm_address !== null) {
      message.evmAddress = object.evm_address;
    }
    return message;
  },
  toAmino(message: QueryGetGasStabilityPoolAddressResponse): QueryGetGasStabilityPoolAddressResponseAmino {
    const obj: any = {};
    obj.cosmos_address = message.cosmosAddress;
    obj.evm_address = message.evmAddress;
    return obj;
  },
  fromAminoMsg(object: QueryGetGasStabilityPoolAddressResponseAminoMsg): QueryGetGasStabilityPoolAddressResponse {
    return QueryGetGasStabilityPoolAddressResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetGasStabilityPoolAddressResponseProtoMsg): QueryGetGasStabilityPoolAddressResponse {
    return QueryGetGasStabilityPoolAddressResponse.decode(message.value);
  },
  toProto(message: QueryGetGasStabilityPoolAddressResponse): Uint8Array {
    return QueryGetGasStabilityPoolAddressResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetGasStabilityPoolAddressResponse): QueryGetGasStabilityPoolAddressResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolAddressResponse",
      value: QueryGetGasStabilityPoolAddressResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetGasStabilityPoolBalance(): QueryGetGasStabilityPoolBalance {
  return {
    chainId: BigInt(0)
  };
}
export const QueryGetGasStabilityPoolBalance = {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolBalance",
  encode(message: QueryGetGasStabilityPoolBalance, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetGasStabilityPoolBalance {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGasStabilityPoolBalance();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetGasStabilityPoolBalance>): QueryGetGasStabilityPoolBalance {
    const message = createBaseQueryGetGasStabilityPoolBalance();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryGetGasStabilityPoolBalanceAmino): QueryGetGasStabilityPoolBalance {
    const message = createBaseQueryGetGasStabilityPoolBalance();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    return message;
  },
  toAmino(message: QueryGetGasStabilityPoolBalance): QueryGetGasStabilityPoolBalanceAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetGasStabilityPoolBalanceAminoMsg): QueryGetGasStabilityPoolBalance {
    return QueryGetGasStabilityPoolBalance.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetGasStabilityPoolBalanceProtoMsg): QueryGetGasStabilityPoolBalance {
    return QueryGetGasStabilityPoolBalance.decode(message.value);
  },
  toProto(message: QueryGetGasStabilityPoolBalance): Uint8Array {
    return QueryGetGasStabilityPoolBalance.encode(message).finish();
  },
  toProtoMsg(message: QueryGetGasStabilityPoolBalance): QueryGetGasStabilityPoolBalanceProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolBalance",
      value: QueryGetGasStabilityPoolBalance.encode(message).finish()
    };
  }
};
function createBaseQueryGetGasStabilityPoolBalanceResponse(): QueryGetGasStabilityPoolBalanceResponse {
  return {
    balance: ""
  };
}
export const QueryGetGasStabilityPoolBalanceResponse = {
  typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolBalanceResponse",
  encode(message: QueryGetGasStabilityPoolBalanceResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.balance !== "") {
      writer.uint32(18).string(message.balance);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetGasStabilityPoolBalanceResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGasStabilityPoolBalanceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.balance = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetGasStabilityPoolBalanceResponse>): QueryGetGasStabilityPoolBalanceResponse {
    const message = createBaseQueryGetGasStabilityPoolBalanceResponse();
    message.balance = object.balance ?? "";
    return message;
  },
  fromAmino(object: QueryGetGasStabilityPoolBalanceResponseAmino): QueryGetGasStabilityPoolBalanceResponse {
    const message = createBaseQueryGetGasStabilityPoolBalanceResponse();
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = object.balance;
    }
    return message;
  },
  toAmino(message: QueryGetGasStabilityPoolBalanceResponse): QueryGetGasStabilityPoolBalanceResponseAmino {
    const obj: any = {};
    obj.balance = message.balance;
    return obj;
  },
  fromAminoMsg(object: QueryGetGasStabilityPoolBalanceResponseAminoMsg): QueryGetGasStabilityPoolBalanceResponse {
    return QueryGetGasStabilityPoolBalanceResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetGasStabilityPoolBalanceResponseProtoMsg): QueryGetGasStabilityPoolBalanceResponse {
    return QueryGetGasStabilityPoolBalanceResponse.decode(message.value);
  },
  toProto(message: QueryGetGasStabilityPoolBalanceResponse): Uint8Array {
    return QueryGetGasStabilityPoolBalanceResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetGasStabilityPoolBalanceResponse): QueryGetGasStabilityPoolBalanceResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryGetGasStabilityPoolBalanceResponse",
      value: QueryGetGasStabilityPoolBalanceResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllGasStabilityPoolBalance(): QueryAllGasStabilityPoolBalance {
  return {};
}
export const QueryAllGasStabilityPoolBalance = {
  typeUrl: "/zetachain.zetacore.fungible.QueryAllGasStabilityPoolBalance",
  encode(_: QueryAllGasStabilityPoolBalance, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllGasStabilityPoolBalance {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllGasStabilityPoolBalance();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<QueryAllGasStabilityPoolBalance>): QueryAllGasStabilityPoolBalance {
    const message = createBaseQueryAllGasStabilityPoolBalance();
    return message;
  },
  fromAmino(_: QueryAllGasStabilityPoolBalanceAmino): QueryAllGasStabilityPoolBalance {
    const message = createBaseQueryAllGasStabilityPoolBalance();
    return message;
  },
  toAmino(_: QueryAllGasStabilityPoolBalance): QueryAllGasStabilityPoolBalanceAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryAllGasStabilityPoolBalanceAminoMsg): QueryAllGasStabilityPoolBalance {
    return QueryAllGasStabilityPoolBalance.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllGasStabilityPoolBalanceProtoMsg): QueryAllGasStabilityPoolBalance {
    return QueryAllGasStabilityPoolBalance.decode(message.value);
  },
  toProto(message: QueryAllGasStabilityPoolBalance): Uint8Array {
    return QueryAllGasStabilityPoolBalance.encode(message).finish();
  },
  toProtoMsg(message: QueryAllGasStabilityPoolBalance): QueryAllGasStabilityPoolBalanceProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryAllGasStabilityPoolBalance",
      value: QueryAllGasStabilityPoolBalance.encode(message).finish()
    };
  }
};
function createBaseQueryAllGasStabilityPoolBalanceResponse(): QueryAllGasStabilityPoolBalanceResponse {
  return {
    balances: []
  };
}
export const QueryAllGasStabilityPoolBalanceResponse = {
  typeUrl: "/zetachain.zetacore.fungible.QueryAllGasStabilityPoolBalanceResponse",
  encode(message: QueryAllGasStabilityPoolBalanceResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.balances) {
      QueryAllGasStabilityPoolBalanceResponse_Balance.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllGasStabilityPoolBalanceResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllGasStabilityPoolBalanceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.balances.push(QueryAllGasStabilityPoolBalanceResponse_Balance.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllGasStabilityPoolBalanceResponse>): QueryAllGasStabilityPoolBalanceResponse {
    const message = createBaseQueryAllGasStabilityPoolBalanceResponse();
    message.balances = object.balances?.map(e => QueryAllGasStabilityPoolBalanceResponse_Balance.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: QueryAllGasStabilityPoolBalanceResponseAmino): QueryAllGasStabilityPoolBalanceResponse {
    const message = createBaseQueryAllGasStabilityPoolBalanceResponse();
    message.balances = object.balances?.map(e => QueryAllGasStabilityPoolBalanceResponse_Balance.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: QueryAllGasStabilityPoolBalanceResponse): QueryAllGasStabilityPoolBalanceResponseAmino {
    const obj: any = {};
    if (message.balances) {
      obj.balances = message.balances.map(e => e ? QueryAllGasStabilityPoolBalanceResponse_Balance.toAmino(e) : undefined);
    } else {
      obj.balances = [];
    }
    return obj;
  },
  fromAminoMsg(object: QueryAllGasStabilityPoolBalanceResponseAminoMsg): QueryAllGasStabilityPoolBalanceResponse {
    return QueryAllGasStabilityPoolBalanceResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllGasStabilityPoolBalanceResponseProtoMsg): QueryAllGasStabilityPoolBalanceResponse {
    return QueryAllGasStabilityPoolBalanceResponse.decode(message.value);
  },
  toProto(message: QueryAllGasStabilityPoolBalanceResponse): Uint8Array {
    return QueryAllGasStabilityPoolBalanceResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllGasStabilityPoolBalanceResponse): QueryAllGasStabilityPoolBalanceResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryAllGasStabilityPoolBalanceResponse",
      value: QueryAllGasStabilityPoolBalanceResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllGasStabilityPoolBalanceResponse_Balance(): QueryAllGasStabilityPoolBalanceResponse_Balance {
  return {
    chainId: BigInt(0),
    balance: ""
  };
}
export const QueryAllGasStabilityPoolBalanceResponse_Balance = {
  typeUrl: "/zetachain.zetacore.fungible.Balance",
  encode(message: QueryAllGasStabilityPoolBalanceResponse_Balance, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.balance !== "") {
      writer.uint32(18).string(message.balance);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllGasStabilityPoolBalanceResponse_Balance {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllGasStabilityPoolBalanceResponse_Balance();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.balance = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllGasStabilityPoolBalanceResponse_Balance>): QueryAllGasStabilityPoolBalanceResponse_Balance {
    const message = createBaseQueryAllGasStabilityPoolBalanceResponse_Balance();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.balance = object.balance ?? "";
    return message;
  },
  fromAmino(object: QueryAllGasStabilityPoolBalanceResponse_BalanceAmino): QueryAllGasStabilityPoolBalanceResponse_Balance {
    const message = createBaseQueryAllGasStabilityPoolBalanceResponse_Balance();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.balance !== undefined && object.balance !== null) {
      message.balance = object.balance;
    }
    return message;
  },
  toAmino(message: QueryAllGasStabilityPoolBalanceResponse_Balance): QueryAllGasStabilityPoolBalanceResponse_BalanceAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.balance = message.balance;
    return obj;
  },
  fromAminoMsg(object: QueryAllGasStabilityPoolBalanceResponse_BalanceAminoMsg): QueryAllGasStabilityPoolBalanceResponse_Balance {
    return QueryAllGasStabilityPoolBalanceResponse_Balance.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllGasStabilityPoolBalanceResponse_BalanceProtoMsg): QueryAllGasStabilityPoolBalanceResponse_Balance {
    return QueryAllGasStabilityPoolBalanceResponse_Balance.decode(message.value);
  },
  toProto(message: QueryAllGasStabilityPoolBalanceResponse_Balance): Uint8Array {
    return QueryAllGasStabilityPoolBalanceResponse_Balance.encode(message).finish();
  },
  toProtoMsg(message: QueryAllGasStabilityPoolBalanceResponse_Balance): QueryAllGasStabilityPoolBalanceResponse_BalanceProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.Balance",
      value: QueryAllGasStabilityPoolBalanceResponse_Balance.encode(message).finish()
    };
  }
};
function createBaseQueryCodeHashRequest(): QueryCodeHashRequest {
  return {
    address: ""
  };
}
export const QueryCodeHashRequest = {
  typeUrl: "/zetachain.zetacore.fungible.QueryCodeHashRequest",
  encode(message: QueryCodeHashRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryCodeHashRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCodeHashRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryCodeHashRequest>): QueryCodeHashRequest {
    const message = createBaseQueryCodeHashRequest();
    message.address = object.address ?? "";
    return message;
  },
  fromAmino(object: QueryCodeHashRequestAmino): QueryCodeHashRequest {
    const message = createBaseQueryCodeHashRequest();
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    return message;
  },
  toAmino(message: QueryCodeHashRequest): QueryCodeHashRequestAmino {
    const obj: any = {};
    obj.address = message.address;
    return obj;
  },
  fromAminoMsg(object: QueryCodeHashRequestAminoMsg): QueryCodeHashRequest {
    return QueryCodeHashRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryCodeHashRequestProtoMsg): QueryCodeHashRequest {
    return QueryCodeHashRequest.decode(message.value);
  },
  toProto(message: QueryCodeHashRequest): Uint8Array {
    return QueryCodeHashRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryCodeHashRequest): QueryCodeHashRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryCodeHashRequest",
      value: QueryCodeHashRequest.encode(message).finish()
    };
  }
};
function createBaseQueryCodeHashResponse(): QueryCodeHashResponse {
  return {
    codeHash: ""
  };
}
export const QueryCodeHashResponse = {
  typeUrl: "/zetachain.zetacore.fungible.QueryCodeHashResponse",
  encode(message: QueryCodeHashResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.codeHash !== "") {
      writer.uint32(10).string(message.codeHash);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryCodeHashResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCodeHashResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.codeHash = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryCodeHashResponse>): QueryCodeHashResponse {
    const message = createBaseQueryCodeHashResponse();
    message.codeHash = object.codeHash ?? "";
    return message;
  },
  fromAmino(object: QueryCodeHashResponseAmino): QueryCodeHashResponse {
    const message = createBaseQueryCodeHashResponse();
    if (object.code_hash !== undefined && object.code_hash !== null) {
      message.codeHash = object.code_hash;
    }
    return message;
  },
  toAmino(message: QueryCodeHashResponse): QueryCodeHashResponseAmino {
    const obj: any = {};
    obj.code_hash = message.codeHash;
    return obj;
  },
  fromAminoMsg(object: QueryCodeHashResponseAminoMsg): QueryCodeHashResponse {
    return QueryCodeHashResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryCodeHashResponseProtoMsg): QueryCodeHashResponse {
    return QueryCodeHashResponse.decode(message.value);
  },
  toProto(message: QueryCodeHashResponse): Uint8Array {
    return QueryCodeHashResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryCodeHashResponse): QueryCodeHashResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.QueryCodeHashResponse",
      value: QueryCodeHashResponse.encode(message).finish()
    };
  }
};