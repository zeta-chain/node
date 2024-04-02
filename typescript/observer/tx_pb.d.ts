// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file observer/tx.proto (package zetachain.zetacore.observer, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { ObserverUpdateReason } from "./observer_pb.js";
import type { HeaderData } from "../pkg/proofs/proofs_pb.js";
import type { ChainParams } from "./params_pb.js";
import type { Blame } from "./blame_pb.js";
import type { BlockHeaderVerificationFlags, GasPriceIncreaseFlags } from "./crosschain_flags_pb.js";

/**
 * @generated from message zetachain.zetacore.observer.MsgUpdateObserver
 */
export declare class MsgUpdateObserver extends Message<MsgUpdateObserver> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: string old_observer_address = 2;
   */
  oldObserverAddress: string;

  /**
   * @generated from field: string new_observer_address = 3;
   */
  newObserverAddress: string;

  /**
   * @generated from field: zetachain.zetacore.observer.ObserverUpdateReason update_reason = 4;
   */
  updateReason: ObserverUpdateReason;

  constructor(data?: PartialMessage<MsgUpdateObserver>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateObserver";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateObserver;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateObserver;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateObserver;

  static equals(a: MsgUpdateObserver | PlainMessage<MsgUpdateObserver> | undefined, b: MsgUpdateObserver | PlainMessage<MsgUpdateObserver> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgUpdateObserverResponse
 */
export declare class MsgUpdateObserverResponse extends Message<MsgUpdateObserverResponse> {
  constructor(data?: PartialMessage<MsgUpdateObserverResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateObserverResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateObserverResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateObserverResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateObserverResponse;

  static equals(a: MsgUpdateObserverResponse | PlainMessage<MsgUpdateObserverResponse> | undefined, b: MsgUpdateObserverResponse | PlainMessage<MsgUpdateObserverResponse> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgVoteBlockHeader
 */
export declare class MsgVoteBlockHeader extends Message<MsgVoteBlockHeader> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: int64 chain_id = 2;
   */
  chainId: bigint;

  /**
   * @generated from field: bytes block_hash = 3;
   */
  blockHash: Uint8Array;

  /**
   * @generated from field: int64 height = 4;
   */
  height: bigint;

  /**
   * @generated from field: proofs.HeaderData header = 5;
   */
  header?: HeaderData;

  constructor(data?: PartialMessage<MsgVoteBlockHeader>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgVoteBlockHeader";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgVoteBlockHeader;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgVoteBlockHeader;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgVoteBlockHeader;

  static equals(a: MsgVoteBlockHeader | PlainMessage<MsgVoteBlockHeader> | undefined, b: MsgVoteBlockHeader | PlainMessage<MsgVoteBlockHeader> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgVoteBlockHeaderResponse
 */
export declare class MsgVoteBlockHeaderResponse extends Message<MsgVoteBlockHeaderResponse> {
  constructor(data?: PartialMessage<MsgVoteBlockHeaderResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgVoteBlockHeaderResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgVoteBlockHeaderResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgVoteBlockHeaderResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgVoteBlockHeaderResponse;

  static equals(a: MsgVoteBlockHeaderResponse | PlainMessage<MsgVoteBlockHeaderResponse> | undefined, b: MsgVoteBlockHeaderResponse | PlainMessage<MsgVoteBlockHeaderResponse> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgUpdateChainParams
 */
export declare class MsgUpdateChainParams extends Message<MsgUpdateChainParams> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: zetachain.zetacore.observer.ChainParams chainParams = 2;
   */
  chainParams?: ChainParams;

  constructor(data?: PartialMessage<MsgUpdateChainParams>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateChainParams";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateChainParams;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateChainParams;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateChainParams;

  static equals(a: MsgUpdateChainParams | PlainMessage<MsgUpdateChainParams> | undefined, b: MsgUpdateChainParams | PlainMessage<MsgUpdateChainParams> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgUpdateChainParamsResponse
 */
export declare class MsgUpdateChainParamsResponse extends Message<MsgUpdateChainParamsResponse> {
  constructor(data?: PartialMessage<MsgUpdateChainParamsResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateChainParamsResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateChainParamsResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateChainParamsResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateChainParamsResponse;

  static equals(a: MsgUpdateChainParamsResponse | PlainMessage<MsgUpdateChainParamsResponse> | undefined, b: MsgUpdateChainParamsResponse | PlainMessage<MsgUpdateChainParamsResponse> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgRemoveChainParams
 */
export declare class MsgRemoveChainParams extends Message<MsgRemoveChainParams> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: int64 chain_id = 2;
   */
  chainId: bigint;

  constructor(data?: PartialMessage<MsgRemoveChainParams>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgRemoveChainParams";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgRemoveChainParams;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgRemoveChainParams;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgRemoveChainParams;

  static equals(a: MsgRemoveChainParams | PlainMessage<MsgRemoveChainParams> | undefined, b: MsgRemoveChainParams | PlainMessage<MsgRemoveChainParams> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgRemoveChainParamsResponse
 */
export declare class MsgRemoveChainParamsResponse extends Message<MsgRemoveChainParamsResponse> {
  constructor(data?: PartialMessage<MsgRemoveChainParamsResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgRemoveChainParamsResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgRemoveChainParamsResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgRemoveChainParamsResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgRemoveChainParamsResponse;

  static equals(a: MsgRemoveChainParamsResponse | PlainMessage<MsgRemoveChainParamsResponse> | undefined, b: MsgRemoveChainParamsResponse | PlainMessage<MsgRemoveChainParamsResponse> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgAddObserver
 */
export declare class MsgAddObserver extends Message<MsgAddObserver> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: string observer_address = 2;
   */
  observerAddress: string;

  /**
   * @generated from field: string zetaclient_grantee_pubkey = 3;
   */
  zetaclientGranteePubkey: string;

  /**
   * @generated from field: bool add_node_account_only = 4;
   */
  addNodeAccountOnly: boolean;

  constructor(data?: PartialMessage<MsgAddObserver>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgAddObserver";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgAddObserver;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgAddObserver;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgAddObserver;

  static equals(a: MsgAddObserver | PlainMessage<MsgAddObserver> | undefined, b: MsgAddObserver | PlainMessage<MsgAddObserver> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgAddObserverResponse
 */
export declare class MsgAddObserverResponse extends Message<MsgAddObserverResponse> {
  constructor(data?: PartialMessage<MsgAddObserverResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgAddObserverResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgAddObserverResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgAddObserverResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgAddObserverResponse;

  static equals(a: MsgAddObserverResponse | PlainMessage<MsgAddObserverResponse> | undefined, b: MsgAddObserverResponse | PlainMessage<MsgAddObserverResponse> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgAddBlameVote
 */
export declare class MsgAddBlameVote extends Message<MsgAddBlameVote> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: int64 chain_id = 2;
   */
  chainId: bigint;

  /**
   * @generated from field: zetachain.zetacore.observer.Blame blame_info = 3;
   */
  blameInfo?: Blame;

  constructor(data?: PartialMessage<MsgAddBlameVote>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgAddBlameVote";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgAddBlameVote;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgAddBlameVote;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgAddBlameVote;

  static equals(a: MsgAddBlameVote | PlainMessage<MsgAddBlameVote> | undefined, b: MsgAddBlameVote | PlainMessage<MsgAddBlameVote> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgAddBlameVoteResponse
 */
export declare class MsgAddBlameVoteResponse extends Message<MsgAddBlameVoteResponse> {
  constructor(data?: PartialMessage<MsgAddBlameVoteResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgAddBlameVoteResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgAddBlameVoteResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgAddBlameVoteResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgAddBlameVoteResponse;

  static equals(a: MsgAddBlameVoteResponse | PlainMessage<MsgAddBlameVoteResponse> | undefined, b: MsgAddBlameVoteResponse | PlainMessage<MsgAddBlameVoteResponse> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgUpdateCrosschainFlags
 */
export declare class MsgUpdateCrosschainFlags extends Message<MsgUpdateCrosschainFlags> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: bool isInboundEnabled = 3;
   */
  isInboundEnabled: boolean;

  /**
   * @generated from field: bool isOutboundEnabled = 4;
   */
  isOutboundEnabled: boolean;

  /**
   * @generated from field: zetachain.zetacore.observer.GasPriceIncreaseFlags gasPriceIncreaseFlags = 5;
   */
  gasPriceIncreaseFlags?: GasPriceIncreaseFlags;

  /**
   * @generated from field: zetachain.zetacore.observer.BlockHeaderVerificationFlags blockHeaderVerificationFlags = 6;
   */
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlags;

  constructor(data?: PartialMessage<MsgUpdateCrosschainFlags>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateCrosschainFlags";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateCrosschainFlags;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateCrosschainFlags;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateCrosschainFlags;

  static equals(a: MsgUpdateCrosschainFlags | PlainMessage<MsgUpdateCrosschainFlags> | undefined, b: MsgUpdateCrosschainFlags | PlainMessage<MsgUpdateCrosschainFlags> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgUpdateCrosschainFlagsResponse
 */
export declare class MsgUpdateCrosschainFlagsResponse extends Message<MsgUpdateCrosschainFlagsResponse> {
  constructor(data?: PartialMessage<MsgUpdateCrosschainFlagsResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateCrosschainFlagsResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateCrosschainFlagsResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateCrosschainFlagsResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateCrosschainFlagsResponse;

  static equals(a: MsgUpdateCrosschainFlagsResponse | PlainMessage<MsgUpdateCrosschainFlagsResponse> | undefined, b: MsgUpdateCrosschainFlagsResponse | PlainMessage<MsgUpdateCrosschainFlagsResponse> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgUpdateKeygen
 */
export declare class MsgUpdateKeygen extends Message<MsgUpdateKeygen> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: int64 block = 2;
   */
  block: bigint;

  constructor(data?: PartialMessage<MsgUpdateKeygen>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateKeygen";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateKeygen;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateKeygen;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateKeygen;

  static equals(a: MsgUpdateKeygen | PlainMessage<MsgUpdateKeygen> | undefined, b: MsgUpdateKeygen | PlainMessage<MsgUpdateKeygen> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgUpdateKeygenResponse
 */
export declare class MsgUpdateKeygenResponse extends Message<MsgUpdateKeygenResponse> {
  constructor(data?: PartialMessage<MsgUpdateKeygenResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateKeygenResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateKeygenResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateKeygenResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateKeygenResponse;

  static equals(a: MsgUpdateKeygenResponse | PlainMessage<MsgUpdateKeygenResponse> | undefined, b: MsgUpdateKeygenResponse | PlainMessage<MsgUpdateKeygenResponse> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgResetChainNonces
 */
export declare class MsgResetChainNonces extends Message<MsgResetChainNonces> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: int64 chain_id = 2;
   */
  chainId: bigint;

  /**
   * @generated from field: int64 chain_nonce_low = 3;
   */
  chainNonceLow: bigint;

  /**
   * @generated from field: int64 chain_nonce_high = 4;
   */
  chainNonceHigh: bigint;

  constructor(data?: PartialMessage<MsgResetChainNonces>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgResetChainNonces";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgResetChainNonces;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgResetChainNonces;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgResetChainNonces;

  static equals(a: MsgResetChainNonces | PlainMessage<MsgResetChainNonces> | undefined, b: MsgResetChainNonces | PlainMessage<MsgResetChainNonces> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgResetChainNoncesResponse
 */
export declare class MsgResetChainNoncesResponse extends Message<MsgResetChainNoncesResponse> {
  constructor(data?: PartialMessage<MsgResetChainNoncesResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgResetChainNoncesResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgResetChainNoncesResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgResetChainNoncesResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgResetChainNoncesResponse;

  static equals(a: MsgResetChainNoncesResponse | PlainMessage<MsgResetChainNoncesResponse> | undefined, b: MsgResetChainNoncesResponse | PlainMessage<MsgResetChainNoncesResponse> | undefined): boolean;
}

