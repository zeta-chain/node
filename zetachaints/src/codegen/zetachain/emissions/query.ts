import { Params, ParamsAmino, ParamsSDKType } from "./params";
import { BinaryReader, BinaryWriter } from "../../binary";
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {}
export interface QueryParamsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.QueryParamsRequest";
  value: Uint8Array;
}
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequestAmino {}
export interface QueryParamsRequestAminoMsg {
  type: "/zetachain.zetacore.emissions.QueryParamsRequest";
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
  typeUrl: "/zetachain.zetacore.emissions.QueryParamsResponse";
  value: Uint8Array;
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponseAmino {
  /** params holds all the parameters of this module. */
  params?: ParamsAmino;
}
export interface QueryParamsResponseAminoMsg {
  type: "/zetachain.zetacore.emissions.QueryParamsResponse";
  value: QueryParamsResponseAmino;
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponseSDKType {
  params: ParamsSDKType;
}
export interface QueryListPoolAddressesRequest {}
export interface QueryListPoolAddressesRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.QueryListPoolAddressesRequest";
  value: Uint8Array;
}
export interface QueryListPoolAddressesRequestAmino {}
export interface QueryListPoolAddressesRequestAminoMsg {
  type: "/zetachain.zetacore.emissions.QueryListPoolAddressesRequest";
  value: QueryListPoolAddressesRequestAmino;
}
export interface QueryListPoolAddressesRequestSDKType {}
export interface QueryListPoolAddressesResponse {
  undistributedObserverBalancesAddress: string;
  undistributedTssBalancesAddress: string;
  emissionModuleAddress: string;
}
export interface QueryListPoolAddressesResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.QueryListPoolAddressesResponse";
  value: Uint8Array;
}
export interface QueryListPoolAddressesResponseAmino {
  undistributed_observer_balances_address?: string;
  undistributed_tss_balances_address?: string;
  emission_module_address?: string;
}
export interface QueryListPoolAddressesResponseAminoMsg {
  type: "/zetachain.zetacore.emissions.QueryListPoolAddressesResponse";
  value: QueryListPoolAddressesResponseAmino;
}
export interface QueryListPoolAddressesResponseSDKType {
  undistributed_observer_balances_address: string;
  undistributed_tss_balances_address: string;
  emission_module_address: string;
}
export interface QueryGetEmissionsFactorsRequest {}
export interface QueryGetEmissionsFactorsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.QueryGetEmissionsFactorsRequest";
  value: Uint8Array;
}
export interface QueryGetEmissionsFactorsRequestAmino {}
export interface QueryGetEmissionsFactorsRequestAminoMsg {
  type: "/zetachain.zetacore.emissions.QueryGetEmissionsFactorsRequest";
  value: QueryGetEmissionsFactorsRequestAmino;
}
export interface QueryGetEmissionsFactorsRequestSDKType {}
export interface QueryGetEmissionsFactorsResponse {
  reservesFactor: string;
  bondFactor: string;
  durationFactor: string;
}
export interface QueryGetEmissionsFactorsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.QueryGetEmissionsFactorsResponse";
  value: Uint8Array;
}
export interface QueryGetEmissionsFactorsResponseAmino {
  reservesFactor?: string;
  bondFactor?: string;
  durationFactor?: string;
}
export interface QueryGetEmissionsFactorsResponseAminoMsg {
  type: "/zetachain.zetacore.emissions.QueryGetEmissionsFactorsResponse";
  value: QueryGetEmissionsFactorsResponseAmino;
}
export interface QueryGetEmissionsFactorsResponseSDKType {
  reservesFactor: string;
  bondFactor: string;
  durationFactor: string;
}
export interface QueryShowAvailableEmissionsRequest {
  address: string;
}
export interface QueryShowAvailableEmissionsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.QueryShowAvailableEmissionsRequest";
  value: Uint8Array;
}
export interface QueryShowAvailableEmissionsRequestAmino {
  address?: string;
}
export interface QueryShowAvailableEmissionsRequestAminoMsg {
  type: "/zetachain.zetacore.emissions.QueryShowAvailableEmissionsRequest";
  value: QueryShowAvailableEmissionsRequestAmino;
}
export interface QueryShowAvailableEmissionsRequestSDKType {
  address: string;
}
export interface QueryShowAvailableEmissionsResponse {
  amount: string;
}
export interface QueryShowAvailableEmissionsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.QueryShowAvailableEmissionsResponse";
  value: Uint8Array;
}
export interface QueryShowAvailableEmissionsResponseAmino {
  amount?: string;
}
export interface QueryShowAvailableEmissionsResponseAminoMsg {
  type: "/zetachain.zetacore.emissions.QueryShowAvailableEmissionsResponse";
  value: QueryShowAvailableEmissionsResponseAmino;
}
export interface QueryShowAvailableEmissionsResponseSDKType {
  amount: string;
}
function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}
export const QueryParamsRequest = {
  typeUrl: "/zetachain.zetacore.emissions.QueryParamsRequest",
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
      typeUrl: "/zetachain.zetacore.emissions.QueryParamsRequest",
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
  typeUrl: "/zetachain.zetacore.emissions.QueryParamsResponse",
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
      typeUrl: "/zetachain.zetacore.emissions.QueryParamsResponse",
      value: QueryParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryListPoolAddressesRequest(): QueryListPoolAddressesRequest {
  return {};
}
export const QueryListPoolAddressesRequest = {
  typeUrl: "/zetachain.zetacore.emissions.QueryListPoolAddressesRequest",
  encode(_: QueryListPoolAddressesRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryListPoolAddressesRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryListPoolAddressesRequest();
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
  fromPartial(_: Partial<QueryListPoolAddressesRequest>): QueryListPoolAddressesRequest {
    const message = createBaseQueryListPoolAddressesRequest();
    return message;
  },
  fromAmino(_: QueryListPoolAddressesRequestAmino): QueryListPoolAddressesRequest {
    const message = createBaseQueryListPoolAddressesRequest();
    return message;
  },
  toAmino(_: QueryListPoolAddressesRequest): QueryListPoolAddressesRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryListPoolAddressesRequestAminoMsg): QueryListPoolAddressesRequest {
    return QueryListPoolAddressesRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryListPoolAddressesRequestProtoMsg): QueryListPoolAddressesRequest {
    return QueryListPoolAddressesRequest.decode(message.value);
  },
  toProto(message: QueryListPoolAddressesRequest): Uint8Array {
    return QueryListPoolAddressesRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryListPoolAddressesRequest): QueryListPoolAddressesRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.QueryListPoolAddressesRequest",
      value: QueryListPoolAddressesRequest.encode(message).finish()
    };
  }
};
function createBaseQueryListPoolAddressesResponse(): QueryListPoolAddressesResponse {
  return {
    undistributedObserverBalancesAddress: "",
    undistributedTssBalancesAddress: "",
    emissionModuleAddress: ""
  };
}
export const QueryListPoolAddressesResponse = {
  typeUrl: "/zetachain.zetacore.emissions.QueryListPoolAddressesResponse",
  encode(message: QueryListPoolAddressesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.undistributedObserverBalancesAddress !== "") {
      writer.uint32(10).string(message.undistributedObserverBalancesAddress);
    }
    if (message.undistributedTssBalancesAddress !== "") {
      writer.uint32(18).string(message.undistributedTssBalancesAddress);
    }
    if (message.emissionModuleAddress !== "") {
      writer.uint32(26).string(message.emissionModuleAddress);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryListPoolAddressesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryListPoolAddressesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.undistributedObserverBalancesAddress = reader.string();
          break;
        case 2:
          message.undistributedTssBalancesAddress = reader.string();
          break;
        case 3:
          message.emissionModuleAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryListPoolAddressesResponse>): QueryListPoolAddressesResponse {
    const message = createBaseQueryListPoolAddressesResponse();
    message.undistributedObserverBalancesAddress = object.undistributedObserverBalancesAddress ?? "";
    message.undistributedTssBalancesAddress = object.undistributedTssBalancesAddress ?? "";
    message.emissionModuleAddress = object.emissionModuleAddress ?? "";
    return message;
  },
  fromAmino(object: QueryListPoolAddressesResponseAmino): QueryListPoolAddressesResponse {
    const message = createBaseQueryListPoolAddressesResponse();
    if (object.undistributed_observer_balances_address !== undefined && object.undistributed_observer_balances_address !== null) {
      message.undistributedObserverBalancesAddress = object.undistributed_observer_balances_address;
    }
    if (object.undistributed_tss_balances_address !== undefined && object.undistributed_tss_balances_address !== null) {
      message.undistributedTssBalancesAddress = object.undistributed_tss_balances_address;
    }
    if (object.emission_module_address !== undefined && object.emission_module_address !== null) {
      message.emissionModuleAddress = object.emission_module_address;
    }
    return message;
  },
  toAmino(message: QueryListPoolAddressesResponse): QueryListPoolAddressesResponseAmino {
    const obj: any = {};
    obj.undistributed_observer_balances_address = message.undistributedObserverBalancesAddress;
    obj.undistributed_tss_balances_address = message.undistributedTssBalancesAddress;
    obj.emission_module_address = message.emissionModuleAddress;
    return obj;
  },
  fromAminoMsg(object: QueryListPoolAddressesResponseAminoMsg): QueryListPoolAddressesResponse {
    return QueryListPoolAddressesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryListPoolAddressesResponseProtoMsg): QueryListPoolAddressesResponse {
    return QueryListPoolAddressesResponse.decode(message.value);
  },
  toProto(message: QueryListPoolAddressesResponse): Uint8Array {
    return QueryListPoolAddressesResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryListPoolAddressesResponse): QueryListPoolAddressesResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.QueryListPoolAddressesResponse",
      value: QueryListPoolAddressesResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetEmissionsFactorsRequest(): QueryGetEmissionsFactorsRequest {
  return {};
}
export const QueryGetEmissionsFactorsRequest = {
  typeUrl: "/zetachain.zetacore.emissions.QueryGetEmissionsFactorsRequest",
  encode(_: QueryGetEmissionsFactorsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetEmissionsFactorsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetEmissionsFactorsRequest();
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
  fromPartial(_: Partial<QueryGetEmissionsFactorsRequest>): QueryGetEmissionsFactorsRequest {
    const message = createBaseQueryGetEmissionsFactorsRequest();
    return message;
  },
  fromAmino(_: QueryGetEmissionsFactorsRequestAmino): QueryGetEmissionsFactorsRequest {
    const message = createBaseQueryGetEmissionsFactorsRequest();
    return message;
  },
  toAmino(_: QueryGetEmissionsFactorsRequest): QueryGetEmissionsFactorsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryGetEmissionsFactorsRequestAminoMsg): QueryGetEmissionsFactorsRequest {
    return QueryGetEmissionsFactorsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetEmissionsFactorsRequestProtoMsg): QueryGetEmissionsFactorsRequest {
    return QueryGetEmissionsFactorsRequest.decode(message.value);
  },
  toProto(message: QueryGetEmissionsFactorsRequest): Uint8Array {
    return QueryGetEmissionsFactorsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetEmissionsFactorsRequest): QueryGetEmissionsFactorsRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.QueryGetEmissionsFactorsRequest",
      value: QueryGetEmissionsFactorsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetEmissionsFactorsResponse(): QueryGetEmissionsFactorsResponse {
  return {
    reservesFactor: "",
    bondFactor: "",
    durationFactor: ""
  };
}
export const QueryGetEmissionsFactorsResponse = {
  typeUrl: "/zetachain.zetacore.emissions.QueryGetEmissionsFactorsResponse",
  encode(message: QueryGetEmissionsFactorsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.reservesFactor !== "") {
      writer.uint32(10).string(message.reservesFactor);
    }
    if (message.bondFactor !== "") {
      writer.uint32(18).string(message.bondFactor);
    }
    if (message.durationFactor !== "") {
      writer.uint32(26).string(message.durationFactor);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetEmissionsFactorsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetEmissionsFactorsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.reservesFactor = reader.string();
          break;
        case 2:
          message.bondFactor = reader.string();
          break;
        case 3:
          message.durationFactor = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetEmissionsFactorsResponse>): QueryGetEmissionsFactorsResponse {
    const message = createBaseQueryGetEmissionsFactorsResponse();
    message.reservesFactor = object.reservesFactor ?? "";
    message.bondFactor = object.bondFactor ?? "";
    message.durationFactor = object.durationFactor ?? "";
    return message;
  },
  fromAmino(object: QueryGetEmissionsFactorsResponseAmino): QueryGetEmissionsFactorsResponse {
    const message = createBaseQueryGetEmissionsFactorsResponse();
    if (object.reservesFactor !== undefined && object.reservesFactor !== null) {
      message.reservesFactor = object.reservesFactor;
    }
    if (object.bondFactor !== undefined && object.bondFactor !== null) {
      message.bondFactor = object.bondFactor;
    }
    if (object.durationFactor !== undefined && object.durationFactor !== null) {
      message.durationFactor = object.durationFactor;
    }
    return message;
  },
  toAmino(message: QueryGetEmissionsFactorsResponse): QueryGetEmissionsFactorsResponseAmino {
    const obj: any = {};
    obj.reservesFactor = message.reservesFactor;
    obj.bondFactor = message.bondFactor;
    obj.durationFactor = message.durationFactor;
    return obj;
  },
  fromAminoMsg(object: QueryGetEmissionsFactorsResponseAminoMsg): QueryGetEmissionsFactorsResponse {
    return QueryGetEmissionsFactorsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetEmissionsFactorsResponseProtoMsg): QueryGetEmissionsFactorsResponse {
    return QueryGetEmissionsFactorsResponse.decode(message.value);
  },
  toProto(message: QueryGetEmissionsFactorsResponse): Uint8Array {
    return QueryGetEmissionsFactorsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetEmissionsFactorsResponse): QueryGetEmissionsFactorsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.QueryGetEmissionsFactorsResponse",
      value: QueryGetEmissionsFactorsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryShowAvailableEmissionsRequest(): QueryShowAvailableEmissionsRequest {
  return {
    address: ""
  };
}
export const QueryShowAvailableEmissionsRequest = {
  typeUrl: "/zetachain.zetacore.emissions.QueryShowAvailableEmissionsRequest",
  encode(message: QueryShowAvailableEmissionsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryShowAvailableEmissionsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryShowAvailableEmissionsRequest();
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
  fromPartial(object: Partial<QueryShowAvailableEmissionsRequest>): QueryShowAvailableEmissionsRequest {
    const message = createBaseQueryShowAvailableEmissionsRequest();
    message.address = object.address ?? "";
    return message;
  },
  fromAmino(object: QueryShowAvailableEmissionsRequestAmino): QueryShowAvailableEmissionsRequest {
    const message = createBaseQueryShowAvailableEmissionsRequest();
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    return message;
  },
  toAmino(message: QueryShowAvailableEmissionsRequest): QueryShowAvailableEmissionsRequestAmino {
    const obj: any = {};
    obj.address = message.address;
    return obj;
  },
  fromAminoMsg(object: QueryShowAvailableEmissionsRequestAminoMsg): QueryShowAvailableEmissionsRequest {
    return QueryShowAvailableEmissionsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryShowAvailableEmissionsRequestProtoMsg): QueryShowAvailableEmissionsRequest {
    return QueryShowAvailableEmissionsRequest.decode(message.value);
  },
  toProto(message: QueryShowAvailableEmissionsRequest): Uint8Array {
    return QueryShowAvailableEmissionsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryShowAvailableEmissionsRequest): QueryShowAvailableEmissionsRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.QueryShowAvailableEmissionsRequest",
      value: QueryShowAvailableEmissionsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryShowAvailableEmissionsResponse(): QueryShowAvailableEmissionsResponse {
  return {
    amount: ""
  };
}
export const QueryShowAvailableEmissionsResponse = {
  typeUrl: "/zetachain.zetacore.emissions.QueryShowAvailableEmissionsResponse",
  encode(message: QueryShowAvailableEmissionsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.amount !== "") {
      writer.uint32(10).string(message.amount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryShowAvailableEmissionsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryShowAvailableEmissionsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.amount = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryShowAvailableEmissionsResponse>): QueryShowAvailableEmissionsResponse {
    const message = createBaseQueryShowAvailableEmissionsResponse();
    message.amount = object.amount ?? "";
    return message;
  },
  fromAmino(object: QueryShowAvailableEmissionsResponseAmino): QueryShowAvailableEmissionsResponse {
    const message = createBaseQueryShowAvailableEmissionsResponse();
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    return message;
  },
  toAmino(message: QueryShowAvailableEmissionsResponse): QueryShowAvailableEmissionsResponseAmino {
    const obj: any = {};
    obj.amount = message.amount;
    return obj;
  },
  fromAminoMsg(object: QueryShowAvailableEmissionsResponseAminoMsg): QueryShowAvailableEmissionsResponse {
    return QueryShowAvailableEmissionsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryShowAvailableEmissionsResponseProtoMsg): QueryShowAvailableEmissionsResponse {
    return QueryShowAvailableEmissionsResponse.decode(message.value);
  },
  toProto(message: QueryShowAvailableEmissionsResponse): Uint8Array {
    return QueryShowAvailableEmissionsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryShowAvailableEmissionsResponse): QueryShowAvailableEmissionsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.QueryShowAvailableEmissionsResponse",
      value: QueryShowAvailableEmissionsResponse.encode(message).finish()
    };
  }
};