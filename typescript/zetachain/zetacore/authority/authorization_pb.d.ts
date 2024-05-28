// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file zetachain/zetacore/authority/authorization.proto (package zetachain.zetacore.authority, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { PolicyType } from "./policies_pb.js";

/**
 * @generated from message zetachain.zetacore.authority.Authorization
 */
export declare class Authorization extends Message<Authorization> {
  /**
   * The URL of the message that needs to be authorized
   *
   * @generated from field: string msg_url = 1;
   */
  msgUrl: string;

  /**
   * The policy that is authorized to access the message
   *
   * @generated from field: zetachain.zetacore.authority.PolicyType authorized_policy = 2;
   */
  authorizedPolicy: PolicyType;

  constructor(data?: PartialMessage<Authorization>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.authority.Authorization";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Authorization;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Authorization;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Authorization;

  static equals(a: Authorization | PlainMessage<Authorization> | undefined, b: Authorization | PlainMessage<Authorization> | undefined): boolean;
}

/**
 * AuthorizationList holds the list of authorizations on zetachain
 *
 * @generated from message zetachain.zetacore.authority.AuthorizationList
 */
export declare class AuthorizationList extends Message<AuthorizationList> {
  /**
   * @generated from field: repeated zetachain.zetacore.authority.Authorization authorizations = 1;
   */
  authorizations: Authorization[];

  constructor(data?: PartialMessage<AuthorizationList>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.authority.AuthorizationList";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): AuthorizationList;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): AuthorizationList;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): AuthorizationList;

  static equals(a: AuthorizationList | PlainMessage<AuthorizationList> | undefined, b: AuthorizationList | PlainMessage<AuthorizationList> | undefined): boolean;
}

