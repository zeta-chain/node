// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file crosschain/genesis.proto (package zetachain.zetacore.crosschain, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { OutTxTracker } from "./out_tx_tracker_pb.js";
import type { GasPrice } from "./gas_price_pb.js";
import type { CrossChainTx, ZetaAccounting } from "./cross_chain_tx_pb.js";
import type { LastBlockHeight } from "./last_block_height_pb.js";
import type { InTxHashToCctx } from "./in_tx_hash_to_cctx_pb.js";
import type { InTxTracker } from "./in_tx_tracker_pb.js";
import type { RateLimiterFlags } from "./rate_limiter_flags_pb.js";

/**
 * GenesisState defines the crosschain modules genesis state.
 *
 * @generated from message zetachain.zetacore.crosschain.GenesisState
 */
export declare class GenesisState extends Message<GenesisState> {
  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.OutTxTracker outTxTrackerList = 2;
   */
  outTxTrackerList: OutTxTracker[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.GasPrice gasPriceList = 5;
   */
  gasPriceList: GasPrice[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.CrossChainTx CrossChainTxs = 7;
   */
  CrossChainTxs: CrossChainTx[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.LastBlockHeight lastBlockHeightList = 8;
   */
  lastBlockHeightList: LastBlockHeight[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.InTxHashToCctx inTxHashToCctxList = 9;
   */
  inTxHashToCctxList: InTxHashToCctx[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.InTxTracker in_tx_tracker_list = 11;
   */
  inTxTrackerList: InTxTracker[];

  /**
   * @generated from field: zetachain.zetacore.crosschain.ZetaAccounting zeta_accounting = 12;
   */
  zetaAccounting?: ZetaAccounting;

  /**
   * @generated from field: repeated string FinalizedInbounds = 16;
   */
  FinalizedInbounds: string[];

  /**
   * @generated from field: zetachain.zetacore.crosschain.RateLimiterFlags rate_limiter_flags = 17;
   */
  rateLimiterFlags?: RateLimiterFlags;

  constructor(data?: PartialMessage<GenesisState>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.crosschain.GenesisState";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GenesisState;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GenesisState;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GenesisState;

  static equals(a: GenesisState | PlainMessage<GenesisState> | undefined, b: GenesisState | PlainMessage<GenesisState> | undefined): boolean;
}

/**
 * Remove legacy types
 * https://github.com/zeta-chain/node/issues/2139
 *
 * @generated from message zetachain.zetacore.crosschain.GenesisState_legacy
 */
export declare class GenesisState_legacy extends Message<GenesisState_legacy> {
  /**
   * @generated from field: zetachain.zetacore.crosschain.Params params = 1;
   */
  params?: Params;

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.OutTxTracker outTxTrackerList = 2;
   */
  outTxTrackerList: OutTxTracker[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.GasPrice gasPriceList = 5;
   */
  gasPriceList: GasPrice[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.CrossChainTx CrossChainTxs = 7;
   */
  CrossChainTxs: CrossChainTx[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.LastBlockHeight lastBlockHeightList = 8;
   */
  lastBlockHeightList: LastBlockHeight[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.InTxHashToCctx inTxHashToCctxList = 9;
   */
  inTxHashToCctxList: InTxHashToCctx[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.InTxTracker in_tx_tracker_list = 11;
   */
  inTxTrackerList: InTxTracker[];

  /**
   * @generated from field: zetachain.zetacore.crosschain.ZetaAccounting zeta_accounting = 12;
   */
  zetaAccounting?: ZetaAccounting;

  /**
   * @generated from field: repeated string FinalizedInbounds = 16;
   */
  FinalizedInbounds: string[];

  constructor(data?: PartialMessage<GenesisState_legacy>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.crosschain.GenesisState_legacy";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GenesisState_legacy;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GenesisState_legacy;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GenesisState_legacy;

  static equals(a: GenesisState_legacy | PlainMessage<GenesisState_legacy> | undefined, b: GenesisState_legacy | PlainMessage<GenesisState_legacy> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.crosschain.Params
 */
export declare class Params extends Message<Params> {
  /**
   * @generated from field: bool enabled = 1;
   */
  enabled: boolean;

  constructor(data?: PartialMessage<Params>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.crosschain.Params";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Params;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Params;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Params;

  static equals(a: Params | PlainMessage<Params> | undefined, b: Params | PlainMessage<Params> | undefined): boolean;
}

