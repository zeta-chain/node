import { Chain, ChainAmino, ChainSDKType } from "../common/common";
import { BinaryReader, BinaryWriter } from "../../binary";
import { Decimal } from "@cosmjs/math";
/** Deprecated(v14):Moved into the authority module */
export enum Policy_Type {
  group1 = 0,
  group2 = 1,
  UNRECOGNIZED = -1,
}
export const Policy_TypeSDKType = Policy_Type;
export const Policy_TypeAmino = Policy_Type;
export function policy_TypeFromJSON(object: any): Policy_Type {
  switch (object) {
    case 0:
    case "group1":
      return Policy_Type.group1;
    case 1:
    case "group2":
      return Policy_Type.group2;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Policy_Type.UNRECOGNIZED;
  }
}
export function policy_TypeToJSON(object: Policy_Type): string {
  switch (object) {
    case Policy_Type.group1:
      return "group1";
    case Policy_Type.group2:
      return "group2";
    case Policy_Type.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export interface ChainParamsList {
  chainParams: ChainParams[];
}
export interface ChainParamsListProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.ChainParamsList";
  value: Uint8Array;
}
export interface ChainParamsListAmino {
  chain_params?: ChainParamsAmino[];
}
export interface ChainParamsListAminoMsg {
  type: "/zetachain.zetacore.observer.ChainParamsList";
  value: ChainParamsListAmino;
}
export interface ChainParamsListSDKType {
  chain_params: ChainParamsSDKType[];
}
export interface ChainParams {
  chainId: bigint;
  confirmationCount: bigint;
  gasPriceTicker: bigint;
  inTxTicker: bigint;
  outTxTicker: bigint;
  watchUtxoTicker: bigint;
  zetaTokenContractAddress: string;
  connectorContractAddress: string;
  erc20CustodyContractAddress: string;
  outboundTxScheduleInterval: bigint;
  outboundTxScheduleLookahead: bigint;
  ballotThreshold: string;
  minObserverDelegation: string;
  isSupported: boolean;
}
export interface ChainParamsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.ChainParams";
  value: Uint8Array;
}
export interface ChainParamsAmino {
  chain_id?: string;
  confirmation_count?: string;
  gas_price_ticker?: string;
  in_tx_ticker?: string;
  out_tx_ticker?: string;
  watch_utxo_ticker?: string;
  zeta_token_contract_address?: string;
  connector_contract_address?: string;
  erc20_custody_contract_address?: string;
  outbound_tx_schedule_interval?: string;
  outbound_tx_schedule_lookahead?: string;
  ballot_threshold?: string;
  min_observer_delegation?: string;
  is_supported?: boolean;
}
export interface ChainParamsAminoMsg {
  type: "/zetachain.zetacore.observer.ChainParams";
  value: ChainParamsAmino;
}
export interface ChainParamsSDKType {
  chain_id: bigint;
  confirmation_count: bigint;
  gas_price_ticker: bigint;
  in_tx_ticker: bigint;
  out_tx_ticker: bigint;
  watch_utxo_ticker: bigint;
  zeta_token_contract_address: string;
  connector_contract_address: string;
  erc20_custody_contract_address: string;
  outbound_tx_schedule_interval: bigint;
  outbound_tx_schedule_lookahead: bigint;
  ballot_threshold: string;
  min_observer_delegation: string;
  is_supported: boolean;
}
/** Deprecated(v13): Use ChainParamsList */
export interface ObserverParams {
  chain?: Chain;
  ballotThreshold: string;
  minObserverDelegation: string;
  isSupported: boolean;
}
export interface ObserverParamsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.ObserverParams";
  value: Uint8Array;
}
/** Deprecated(v13): Use ChainParamsList */
export interface ObserverParamsAmino {
  chain?: ChainAmino;
  ballot_threshold?: string;
  min_observer_delegation?: string;
  is_supported?: boolean;
}
export interface ObserverParamsAminoMsg {
  type: "/zetachain.zetacore.observer.ObserverParams";
  value: ObserverParamsAmino;
}
/** Deprecated(v13): Use ChainParamsList */
export interface ObserverParamsSDKType {
  chain?: ChainSDKType;
  ballot_threshold: string;
  min_observer_delegation: string;
  is_supported: boolean;
}
/** Deprecated(v14):Moved into the authority module */
export interface Admin_Policy {
  policyType: Policy_Type;
  address: string;
}
export interface Admin_PolicyProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.Admin_Policy";
  value: Uint8Array;
}
/** Deprecated(v14):Moved into the authority module */
export interface Admin_PolicyAmino {
  policy_type?: Policy_Type;
  address?: string;
}
export interface Admin_PolicyAminoMsg {
  type: "/zetachain.zetacore.observer.Admin_Policy";
  value: Admin_PolicyAmino;
}
/** Deprecated(v14):Moved into the authority module */
export interface Admin_PolicySDKType {
  policy_type: Policy_Type;
  address: string;
}
/** Params defines the parameters for the module. */
export interface Params {
  /** Deprecated(v13): Use ChainParamsList */
  observerParams: ObserverParams[];
  /** Deprecated(v14):Moved into the authority module */
  adminPolicy: Admin_Policy[];
  ballotMaturityBlocks: bigint;
}
export interface ParamsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.Params";
  value: Uint8Array;
}
/** Params defines the parameters for the module. */
export interface ParamsAmino {
  /** Deprecated(v13): Use ChainParamsList */
  observer_params?: ObserverParamsAmino[];
  /** Deprecated(v14):Moved into the authority module */
  admin_policy?: Admin_PolicyAmino[];
  ballot_maturity_blocks?: string;
}
export interface ParamsAminoMsg {
  type: "/zetachain.zetacore.observer.Params";
  value: ParamsAmino;
}
/** Params defines the parameters for the module. */
export interface ParamsSDKType {
  observer_params: ObserverParamsSDKType[];
  admin_policy: Admin_PolicySDKType[];
  ballot_maturity_blocks: bigint;
}
function createBaseChainParamsList(): ChainParamsList {
  return {
    chainParams: []
  };
}
export const ChainParamsList = {
  typeUrl: "/zetachain.zetacore.observer.ChainParamsList",
  encode(message: ChainParamsList, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.chainParams) {
      ChainParams.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ChainParamsList {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChainParamsList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainParams.push(ChainParams.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ChainParamsList>): ChainParamsList {
    const message = createBaseChainParamsList();
    message.chainParams = object.chainParams?.map(e => ChainParams.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: ChainParamsListAmino): ChainParamsList {
    const message = createBaseChainParamsList();
    message.chainParams = object.chain_params?.map(e => ChainParams.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: ChainParamsList): ChainParamsListAmino {
    const obj: any = {};
    if (message.chainParams) {
      obj.chain_params = message.chainParams.map(e => e ? ChainParams.toAmino(e) : undefined);
    } else {
      obj.chain_params = [];
    }
    return obj;
  },
  fromAminoMsg(object: ChainParamsListAminoMsg): ChainParamsList {
    return ChainParamsList.fromAmino(object.value);
  },
  fromProtoMsg(message: ChainParamsListProtoMsg): ChainParamsList {
    return ChainParamsList.decode(message.value);
  },
  toProto(message: ChainParamsList): Uint8Array {
    return ChainParamsList.encode(message).finish();
  },
  toProtoMsg(message: ChainParamsList): ChainParamsListProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.ChainParamsList",
      value: ChainParamsList.encode(message).finish()
    };
  }
};
function createBaseChainParams(): ChainParams {
  return {
    chainId: BigInt(0),
    confirmationCount: BigInt(0),
    gasPriceTicker: BigInt(0),
    inTxTicker: BigInt(0),
    outTxTicker: BigInt(0),
    watchUtxoTicker: BigInt(0),
    zetaTokenContractAddress: "",
    connectorContractAddress: "",
    erc20CustodyContractAddress: "",
    outboundTxScheduleInterval: BigInt(0),
    outboundTxScheduleLookahead: BigInt(0),
    ballotThreshold: "",
    minObserverDelegation: "",
    isSupported: false
  };
}
export const ChainParams = {
  typeUrl: "/zetachain.zetacore.observer.ChainParams",
  encode(message: ChainParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(88).int64(message.chainId);
    }
    if (message.confirmationCount !== BigInt(0)) {
      writer.uint32(8).uint64(message.confirmationCount);
    }
    if (message.gasPriceTicker !== BigInt(0)) {
      writer.uint32(16).uint64(message.gasPriceTicker);
    }
    if (message.inTxTicker !== BigInt(0)) {
      writer.uint32(24).uint64(message.inTxTicker);
    }
    if (message.outTxTicker !== BigInt(0)) {
      writer.uint32(32).uint64(message.outTxTicker);
    }
    if (message.watchUtxoTicker !== BigInt(0)) {
      writer.uint32(40).uint64(message.watchUtxoTicker);
    }
    if (message.zetaTokenContractAddress !== "") {
      writer.uint32(66).string(message.zetaTokenContractAddress);
    }
    if (message.connectorContractAddress !== "") {
      writer.uint32(74).string(message.connectorContractAddress);
    }
    if (message.erc20CustodyContractAddress !== "") {
      writer.uint32(82).string(message.erc20CustodyContractAddress);
    }
    if (message.outboundTxScheduleInterval !== BigInt(0)) {
      writer.uint32(96).int64(message.outboundTxScheduleInterval);
    }
    if (message.outboundTxScheduleLookahead !== BigInt(0)) {
      writer.uint32(104).int64(message.outboundTxScheduleLookahead);
    }
    if (message.ballotThreshold !== "") {
      writer.uint32(114).string(Decimal.fromUserInput(message.ballotThreshold, 18).atomics);
    }
    if (message.minObserverDelegation !== "") {
      writer.uint32(122).string(Decimal.fromUserInput(message.minObserverDelegation, 18).atomics);
    }
    if (message.isSupported === true) {
      writer.uint32(128).bool(message.isSupported);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ChainParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChainParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 11:
          message.chainId = reader.int64();
          break;
        case 1:
          message.confirmationCount = reader.uint64();
          break;
        case 2:
          message.gasPriceTicker = reader.uint64();
          break;
        case 3:
          message.inTxTicker = reader.uint64();
          break;
        case 4:
          message.outTxTicker = reader.uint64();
          break;
        case 5:
          message.watchUtxoTicker = reader.uint64();
          break;
        case 8:
          message.zetaTokenContractAddress = reader.string();
          break;
        case 9:
          message.connectorContractAddress = reader.string();
          break;
        case 10:
          message.erc20CustodyContractAddress = reader.string();
          break;
        case 12:
          message.outboundTxScheduleInterval = reader.int64();
          break;
        case 13:
          message.outboundTxScheduleLookahead = reader.int64();
          break;
        case 14:
          message.ballotThreshold = Decimal.fromAtomics(reader.string(), 18).toString();
          break;
        case 15:
          message.minObserverDelegation = Decimal.fromAtomics(reader.string(), 18).toString();
          break;
        case 16:
          message.isSupported = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ChainParams>): ChainParams {
    const message = createBaseChainParams();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.confirmationCount = object.confirmationCount !== undefined && object.confirmationCount !== null ? BigInt(object.confirmationCount.toString()) : BigInt(0);
    message.gasPriceTicker = object.gasPriceTicker !== undefined && object.gasPriceTicker !== null ? BigInt(object.gasPriceTicker.toString()) : BigInt(0);
    message.inTxTicker = object.inTxTicker !== undefined && object.inTxTicker !== null ? BigInt(object.inTxTicker.toString()) : BigInt(0);
    message.outTxTicker = object.outTxTicker !== undefined && object.outTxTicker !== null ? BigInt(object.outTxTicker.toString()) : BigInt(0);
    message.watchUtxoTicker = object.watchUtxoTicker !== undefined && object.watchUtxoTicker !== null ? BigInt(object.watchUtxoTicker.toString()) : BigInt(0);
    message.zetaTokenContractAddress = object.zetaTokenContractAddress ?? "";
    message.connectorContractAddress = object.connectorContractAddress ?? "";
    message.erc20CustodyContractAddress = object.erc20CustodyContractAddress ?? "";
    message.outboundTxScheduleInterval = object.outboundTxScheduleInterval !== undefined && object.outboundTxScheduleInterval !== null ? BigInt(object.outboundTxScheduleInterval.toString()) : BigInt(0);
    message.outboundTxScheduleLookahead = object.outboundTxScheduleLookahead !== undefined && object.outboundTxScheduleLookahead !== null ? BigInt(object.outboundTxScheduleLookahead.toString()) : BigInt(0);
    message.ballotThreshold = object.ballotThreshold ?? "";
    message.minObserverDelegation = object.minObserverDelegation ?? "";
    message.isSupported = object.isSupported ?? false;
    return message;
  },
  fromAmino(object: ChainParamsAmino): ChainParams {
    const message = createBaseChainParams();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.confirmation_count !== undefined && object.confirmation_count !== null) {
      message.confirmationCount = BigInt(object.confirmation_count);
    }
    if (object.gas_price_ticker !== undefined && object.gas_price_ticker !== null) {
      message.gasPriceTicker = BigInt(object.gas_price_ticker);
    }
    if (object.in_tx_ticker !== undefined && object.in_tx_ticker !== null) {
      message.inTxTicker = BigInt(object.in_tx_ticker);
    }
    if (object.out_tx_ticker !== undefined && object.out_tx_ticker !== null) {
      message.outTxTicker = BigInt(object.out_tx_ticker);
    }
    if (object.watch_utxo_ticker !== undefined && object.watch_utxo_ticker !== null) {
      message.watchUtxoTicker = BigInt(object.watch_utxo_ticker);
    }
    if (object.zeta_token_contract_address !== undefined && object.zeta_token_contract_address !== null) {
      message.zetaTokenContractAddress = object.zeta_token_contract_address;
    }
    if (object.connector_contract_address !== undefined && object.connector_contract_address !== null) {
      message.connectorContractAddress = object.connector_contract_address;
    }
    if (object.erc20_custody_contract_address !== undefined && object.erc20_custody_contract_address !== null) {
      message.erc20CustodyContractAddress = object.erc20_custody_contract_address;
    }
    if (object.outbound_tx_schedule_interval !== undefined && object.outbound_tx_schedule_interval !== null) {
      message.outboundTxScheduleInterval = BigInt(object.outbound_tx_schedule_interval);
    }
    if (object.outbound_tx_schedule_lookahead !== undefined && object.outbound_tx_schedule_lookahead !== null) {
      message.outboundTxScheduleLookahead = BigInt(object.outbound_tx_schedule_lookahead);
    }
    if (object.ballot_threshold !== undefined && object.ballot_threshold !== null) {
      message.ballotThreshold = object.ballot_threshold;
    }
    if (object.min_observer_delegation !== undefined && object.min_observer_delegation !== null) {
      message.minObserverDelegation = object.min_observer_delegation;
    }
    if (object.is_supported !== undefined && object.is_supported !== null) {
      message.isSupported = object.is_supported;
    }
    return message;
  },
  toAmino(message: ChainParams): ChainParamsAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.confirmation_count = message.confirmationCount ? message.confirmationCount.toString() : undefined;
    obj.gas_price_ticker = message.gasPriceTicker ? message.gasPriceTicker.toString() : undefined;
    obj.in_tx_ticker = message.inTxTicker ? message.inTxTicker.toString() : undefined;
    obj.out_tx_ticker = message.outTxTicker ? message.outTxTicker.toString() : undefined;
    obj.watch_utxo_ticker = message.watchUtxoTicker ? message.watchUtxoTicker.toString() : undefined;
    obj.zeta_token_contract_address = message.zetaTokenContractAddress;
    obj.connector_contract_address = message.connectorContractAddress;
    obj.erc20_custody_contract_address = message.erc20CustodyContractAddress;
    obj.outbound_tx_schedule_interval = message.outboundTxScheduleInterval ? message.outboundTxScheduleInterval.toString() : undefined;
    obj.outbound_tx_schedule_lookahead = message.outboundTxScheduleLookahead ? message.outboundTxScheduleLookahead.toString() : undefined;
    obj.ballot_threshold = message.ballotThreshold;
    obj.min_observer_delegation = message.minObserverDelegation;
    obj.is_supported = message.isSupported;
    return obj;
  },
  fromAminoMsg(object: ChainParamsAminoMsg): ChainParams {
    return ChainParams.fromAmino(object.value);
  },
  fromProtoMsg(message: ChainParamsProtoMsg): ChainParams {
    return ChainParams.decode(message.value);
  },
  toProto(message: ChainParams): Uint8Array {
    return ChainParams.encode(message).finish();
  },
  toProtoMsg(message: ChainParams): ChainParamsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.ChainParams",
      value: ChainParams.encode(message).finish()
    };
  }
};
function createBaseObserverParams(): ObserverParams {
  return {
    chain: undefined,
    ballotThreshold: "",
    minObserverDelegation: "",
    isSupported: false
  };
}
export const ObserverParams = {
  typeUrl: "/zetachain.zetacore.observer.ObserverParams",
  encode(message: ObserverParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chain !== undefined) {
      Chain.encode(message.chain, writer.uint32(10).fork()).ldelim();
    }
    if (message.ballotThreshold !== "") {
      writer.uint32(26).string(Decimal.fromUserInput(message.ballotThreshold, 18).atomics);
    }
    if (message.minObserverDelegation !== "") {
      writer.uint32(34).string(Decimal.fromUserInput(message.minObserverDelegation, 18).atomics);
    }
    if (message.isSupported === true) {
      writer.uint32(40).bool(message.isSupported);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ObserverParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObserverParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chain = Chain.decode(reader, reader.uint32());
          break;
        case 3:
          message.ballotThreshold = Decimal.fromAtomics(reader.string(), 18).toString();
          break;
        case 4:
          message.minObserverDelegation = Decimal.fromAtomics(reader.string(), 18).toString();
          break;
        case 5:
          message.isSupported = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ObserverParams>): ObserverParams {
    const message = createBaseObserverParams();
    message.chain = object.chain !== undefined && object.chain !== null ? Chain.fromPartial(object.chain) : undefined;
    message.ballotThreshold = object.ballotThreshold ?? "";
    message.minObserverDelegation = object.minObserverDelegation ?? "";
    message.isSupported = object.isSupported ?? false;
    return message;
  },
  fromAmino(object: ObserverParamsAmino): ObserverParams {
    const message = createBaseObserverParams();
    if (object.chain !== undefined && object.chain !== null) {
      message.chain = Chain.fromAmino(object.chain);
    }
    if (object.ballot_threshold !== undefined && object.ballot_threshold !== null) {
      message.ballotThreshold = object.ballot_threshold;
    }
    if (object.min_observer_delegation !== undefined && object.min_observer_delegation !== null) {
      message.minObserverDelegation = object.min_observer_delegation;
    }
    if (object.is_supported !== undefined && object.is_supported !== null) {
      message.isSupported = object.is_supported;
    }
    return message;
  },
  toAmino(message: ObserverParams): ObserverParamsAmino {
    const obj: any = {};
    obj.chain = message.chain ? Chain.toAmino(message.chain) : undefined;
    obj.ballot_threshold = message.ballotThreshold;
    obj.min_observer_delegation = message.minObserverDelegation;
    obj.is_supported = message.isSupported;
    return obj;
  },
  fromAminoMsg(object: ObserverParamsAminoMsg): ObserverParams {
    return ObserverParams.fromAmino(object.value);
  },
  fromProtoMsg(message: ObserverParamsProtoMsg): ObserverParams {
    return ObserverParams.decode(message.value);
  },
  toProto(message: ObserverParams): Uint8Array {
    return ObserverParams.encode(message).finish();
  },
  toProtoMsg(message: ObserverParams): ObserverParamsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.ObserverParams",
      value: ObserverParams.encode(message).finish()
    };
  }
};
function createBaseAdmin_Policy(): Admin_Policy {
  return {
    policyType: 0,
    address: ""
  };
}
export const Admin_Policy = {
  typeUrl: "/zetachain.zetacore.observer.Admin_Policy",
  encode(message: Admin_Policy, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.policyType !== 0) {
      writer.uint32(8).int32(message.policyType);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Admin_Policy {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAdmin_Policy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policyType = (reader.int32() as any);
          break;
        case 2:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Admin_Policy>): Admin_Policy {
    const message = createBaseAdmin_Policy();
    message.policyType = object.policyType ?? 0;
    message.address = object.address ?? "";
    return message;
  },
  fromAmino(object: Admin_PolicyAmino): Admin_Policy {
    const message = createBaseAdmin_Policy();
    if (object.policy_type !== undefined && object.policy_type !== null) {
      message.policyType = policy_TypeFromJSON(object.policy_type);
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    return message;
  },
  toAmino(message: Admin_Policy): Admin_PolicyAmino {
    const obj: any = {};
    obj.policy_type = message.policyType;
    obj.address = message.address;
    return obj;
  },
  fromAminoMsg(object: Admin_PolicyAminoMsg): Admin_Policy {
    return Admin_Policy.fromAmino(object.value);
  },
  fromProtoMsg(message: Admin_PolicyProtoMsg): Admin_Policy {
    return Admin_Policy.decode(message.value);
  },
  toProto(message: Admin_Policy): Uint8Array {
    return Admin_Policy.encode(message).finish();
  },
  toProtoMsg(message: Admin_Policy): Admin_PolicyProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.Admin_Policy",
      value: Admin_Policy.encode(message).finish()
    };
  }
};
function createBaseParams(): Params {
  return {
    observerParams: [],
    adminPolicy: [],
    ballotMaturityBlocks: BigInt(0)
  };
}
export const Params = {
  typeUrl: "/zetachain.zetacore.observer.Params",
  encode(message: Params, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.observerParams) {
      ObserverParams.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.adminPolicy) {
      Admin_Policy.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.ballotMaturityBlocks !== BigInt(0)) {
      writer.uint32(24).int64(message.ballotMaturityBlocks);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Params {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.observerParams.push(ObserverParams.decode(reader, reader.uint32()));
          break;
        case 2:
          message.adminPolicy.push(Admin_Policy.decode(reader, reader.uint32()));
          break;
        case 3:
          message.ballotMaturityBlocks = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Params>): Params {
    const message = createBaseParams();
    message.observerParams = object.observerParams?.map(e => ObserverParams.fromPartial(e)) || [];
    message.adminPolicy = object.adminPolicy?.map(e => Admin_Policy.fromPartial(e)) || [];
    message.ballotMaturityBlocks = object.ballotMaturityBlocks !== undefined && object.ballotMaturityBlocks !== null ? BigInt(object.ballotMaturityBlocks.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: ParamsAmino): Params {
    const message = createBaseParams();
    message.observerParams = object.observer_params?.map(e => ObserverParams.fromAmino(e)) || [];
    message.adminPolicy = object.admin_policy?.map(e => Admin_Policy.fromAmino(e)) || [];
    if (object.ballot_maturity_blocks !== undefined && object.ballot_maturity_blocks !== null) {
      message.ballotMaturityBlocks = BigInt(object.ballot_maturity_blocks);
    }
    return message;
  },
  toAmino(message: Params): ParamsAmino {
    const obj: any = {};
    if (message.observerParams) {
      obj.observer_params = message.observerParams.map(e => e ? ObserverParams.toAmino(e) : undefined);
    } else {
      obj.observer_params = [];
    }
    if (message.adminPolicy) {
      obj.admin_policy = message.adminPolicy.map(e => e ? Admin_Policy.toAmino(e) : undefined);
    } else {
      obj.admin_policy = [];
    }
    obj.ballot_maturity_blocks = message.ballotMaturityBlocks ? message.ballotMaturityBlocks.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: ParamsAminoMsg): Params {
    return Params.fromAmino(object.value);
  },
  fromProtoMsg(message: ParamsProtoMsg): Params {
    return Params.decode(message.value);
  },
  toProto(message: Params): Uint8Array {
    return Params.encode(message).finish();
  },
  toProtoMsg(message: Params): ParamsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.Params",
      value: Params.encode(message).finish()
    };
  }
};