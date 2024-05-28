// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file zetachain/zetacore/authority/genesis.proto (package zetachain.zetacore.authority, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { Policies } from "./policies_pb.js";
import type { AuthorizationList } from "./authorization_pb.js";

/**
 * GenesisState defines the authority module's genesis state.
 *
 * @generated from message zetachain.zetacore.authority.GenesisState
 */
export declare class GenesisState extends Message<GenesisState> {
  /**
   * @generated from field: zetachain.zetacore.authority.Policies policies = 1;
   */
  policies?: Policies;

  /**
   * @generated from field: zetachain.zetacore.authority.AuthorizationList authorization_list = 2;
   */
  authorizationList?: AuthorizationList;

  constructor(data?: PartialMessage<GenesisState>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.authority.GenesisState";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GenesisState;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GenesisState;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GenesisState;

  static equals(a: GenesisState | PlainMessage<GenesisState> | undefined, b: GenesisState | PlainMessage<GenesisState> | undefined): boolean;
}

