import { Policies, PoliciesAmino, PoliciesSDKType } from "./policies";
import { BinaryReader, BinaryWriter } from "../../binary";
/** QueryGetPoliciesRequest is the request type for the Query/Policies RPC method. */
export interface QueryGetPoliciesRequest {}
export interface QueryGetPoliciesRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.authority.QueryGetPoliciesRequest";
  value: Uint8Array;
}
/** QueryGetPoliciesRequest is the request type for the Query/Policies RPC method. */
export interface QueryGetPoliciesRequestAmino {}
export interface QueryGetPoliciesRequestAminoMsg {
  type: "/zetachain.zetacore.authority.QueryGetPoliciesRequest";
  value: QueryGetPoliciesRequestAmino;
}
/** QueryGetPoliciesRequest is the request type for the Query/Policies RPC method. */
export interface QueryGetPoliciesRequestSDKType {}
/** QueryGetPoliciesResponse is the response type for the Query/Policies RPC method. */
export interface QueryGetPoliciesResponse {
  policies: Policies;
}
export interface QueryGetPoliciesResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.authority.QueryGetPoliciesResponse";
  value: Uint8Array;
}
/** QueryGetPoliciesResponse is the response type for the Query/Policies RPC method. */
export interface QueryGetPoliciesResponseAmino {
  policies?: PoliciesAmino;
}
export interface QueryGetPoliciesResponseAminoMsg {
  type: "/zetachain.zetacore.authority.QueryGetPoliciesResponse";
  value: QueryGetPoliciesResponseAmino;
}
/** QueryGetPoliciesResponse is the response type for the Query/Policies RPC method. */
export interface QueryGetPoliciesResponseSDKType {
  policies: PoliciesSDKType;
}
function createBaseQueryGetPoliciesRequest(): QueryGetPoliciesRequest {
  return {};
}
export const QueryGetPoliciesRequest = {
  typeUrl: "/zetachain.zetacore.authority.QueryGetPoliciesRequest",
  encode(_: QueryGetPoliciesRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetPoliciesRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetPoliciesRequest();
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
  fromPartial(_: Partial<QueryGetPoliciesRequest>): QueryGetPoliciesRequest {
    const message = createBaseQueryGetPoliciesRequest();
    return message;
  },
  fromAmino(_: QueryGetPoliciesRequestAmino): QueryGetPoliciesRequest {
    const message = createBaseQueryGetPoliciesRequest();
    return message;
  },
  toAmino(_: QueryGetPoliciesRequest): QueryGetPoliciesRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryGetPoliciesRequestAminoMsg): QueryGetPoliciesRequest {
    return QueryGetPoliciesRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetPoliciesRequestProtoMsg): QueryGetPoliciesRequest {
    return QueryGetPoliciesRequest.decode(message.value);
  },
  toProto(message: QueryGetPoliciesRequest): Uint8Array {
    return QueryGetPoliciesRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetPoliciesRequest): QueryGetPoliciesRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.authority.QueryGetPoliciesRequest",
      value: QueryGetPoliciesRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetPoliciesResponse(): QueryGetPoliciesResponse {
  return {
    policies: Policies.fromPartial({})
  };
}
export const QueryGetPoliciesResponse = {
  typeUrl: "/zetachain.zetacore.authority.QueryGetPoliciesResponse",
  encode(message: QueryGetPoliciesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.policies !== undefined) {
      Policies.encode(message.policies, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetPoliciesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetPoliciesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policies = Policies.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetPoliciesResponse>): QueryGetPoliciesResponse {
    const message = createBaseQueryGetPoliciesResponse();
    message.policies = object.policies !== undefined && object.policies !== null ? Policies.fromPartial(object.policies) : undefined;
    return message;
  },
  fromAmino(object: QueryGetPoliciesResponseAmino): QueryGetPoliciesResponse {
    const message = createBaseQueryGetPoliciesResponse();
    if (object.policies !== undefined && object.policies !== null) {
      message.policies = Policies.fromAmino(object.policies);
    }
    return message;
  },
  toAmino(message: QueryGetPoliciesResponse): QueryGetPoliciesResponseAmino {
    const obj: any = {};
    obj.policies = message.policies ? Policies.toAmino(message.policies) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetPoliciesResponseAminoMsg): QueryGetPoliciesResponse {
    return QueryGetPoliciesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetPoliciesResponseProtoMsg): QueryGetPoliciesResponse {
    return QueryGetPoliciesResponse.decode(message.value);
  },
  toProto(message: QueryGetPoliciesResponse): Uint8Array {
    return QueryGetPoliciesResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetPoliciesResponse): QueryGetPoliciesResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.authority.QueryGetPoliciesResponse",
      value: QueryGetPoliciesResponse.encode(message).finish()
    };
  }
};