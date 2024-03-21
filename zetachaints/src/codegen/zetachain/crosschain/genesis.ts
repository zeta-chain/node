import { Params, ParamsAmino, ParamsSDKType } from "./params";
import { OutTxTracker, OutTxTrackerAmino, OutTxTrackerSDKType } from "./out_tx_tracker";
import { GasPrice, GasPriceAmino, GasPriceSDKType } from "./gas_price";
import { CrossChainTx, CrossChainTxAmino, CrossChainTxSDKType, ZetaAccounting, ZetaAccountingAmino, ZetaAccountingSDKType } from "./cross_chain_tx";
import { LastBlockHeight, LastBlockHeightAmino, LastBlockHeightSDKType } from "./last_block_height";
import { InTxHashToCctx, InTxHashToCctxAmino, InTxHashToCctxSDKType } from "./in_tx_hash_to_cctx";
import { InTxTracker, InTxTrackerAmino, InTxTrackerSDKType } from "./in_tx_tracker";
import { BinaryReader, BinaryWriter } from "../../binary";
/** GenesisState defines the metacore module's genesis state. */
export interface GenesisState {
  params: Params;
  outTxTrackerList: OutTxTracker[];
  gasPriceList: GasPrice[];
  CrossChainTxs: CrossChainTx[];
  lastBlockHeightList: LastBlockHeight[];
  inTxHashToCctxList: InTxHashToCctx[];
  inTxTrackerList: InTxTracker[];
  zetaAccounting: ZetaAccounting;
  FinalizedInbounds: string[];
}
export interface GenesisStateProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the metacore module's genesis state. */
export interface GenesisStateAmino {
  params?: ParamsAmino;
  outTxTrackerList?: OutTxTrackerAmino[];
  gasPriceList?: GasPriceAmino[];
  CrossChainTxs?: CrossChainTxAmino[];
  lastBlockHeightList?: LastBlockHeightAmino[];
  inTxHashToCctxList?: InTxHashToCctxAmino[];
  in_tx_tracker_list?: InTxTrackerAmino[];
  zeta_accounting?: ZetaAccountingAmino;
  FinalizedInbounds?: string[];
}
export interface GenesisStateAminoMsg {
  type: "/zetachain.zetacore.crosschain.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the metacore module's genesis state. */
export interface GenesisStateSDKType {
  params: ParamsSDKType;
  outTxTrackerList: OutTxTrackerSDKType[];
  gasPriceList: GasPriceSDKType[];
  CrossChainTxs: CrossChainTxSDKType[];
  lastBlockHeightList: LastBlockHeightSDKType[];
  inTxHashToCctxList: InTxHashToCctxSDKType[];
  in_tx_tracker_list: InTxTrackerSDKType[];
  zeta_accounting: ZetaAccountingSDKType;
  FinalizedInbounds: string[];
}
function createBaseGenesisState(): GenesisState {
  return {
    params: Params.fromPartial({}),
    outTxTrackerList: [],
    gasPriceList: [],
    CrossChainTxs: [],
    lastBlockHeightList: [],
    inTxHashToCctxList: [],
    inTxTrackerList: [],
    zetaAccounting: ZetaAccounting.fromPartial({}),
    FinalizedInbounds: []
  };
}
export const GenesisState = {
  typeUrl: "/zetachain.zetacore.crosschain.GenesisState",
  encode(message: GenesisState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.outTxTrackerList) {
      OutTxTracker.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.gasPriceList) {
      GasPrice.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.CrossChainTxs) {
      CrossChainTx.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.lastBlockHeightList) {
      LastBlockHeight.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    for (const v of message.inTxHashToCctxList) {
      InTxHashToCctx.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.inTxTrackerList) {
      InTxTracker.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    if (message.zetaAccounting !== undefined) {
      ZetaAccounting.encode(message.zetaAccounting, writer.uint32(98).fork()).ldelim();
    }
    for (const v of message.FinalizedInbounds) {
      writer.uint32(130).string(v!);
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
          message.outTxTrackerList.push(OutTxTracker.decode(reader, reader.uint32()));
          break;
        case 5:
          message.gasPriceList.push(GasPrice.decode(reader, reader.uint32()));
          break;
        case 7:
          message.CrossChainTxs.push(CrossChainTx.decode(reader, reader.uint32()));
          break;
        case 8:
          message.lastBlockHeightList.push(LastBlockHeight.decode(reader, reader.uint32()));
          break;
        case 9:
          message.inTxHashToCctxList.push(InTxHashToCctx.decode(reader, reader.uint32()));
          break;
        case 11:
          message.inTxTrackerList.push(InTxTracker.decode(reader, reader.uint32()));
          break;
        case 12:
          message.zetaAccounting = ZetaAccounting.decode(reader, reader.uint32());
          break;
        case 16:
          message.FinalizedInbounds.push(reader.string());
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
    message.outTxTrackerList = object.outTxTrackerList?.map(e => OutTxTracker.fromPartial(e)) || [];
    message.gasPriceList = object.gasPriceList?.map(e => GasPrice.fromPartial(e)) || [];
    message.CrossChainTxs = object.CrossChainTxs?.map(e => CrossChainTx.fromPartial(e)) || [];
    message.lastBlockHeightList = object.lastBlockHeightList?.map(e => LastBlockHeight.fromPartial(e)) || [];
    message.inTxHashToCctxList = object.inTxHashToCctxList?.map(e => InTxHashToCctx.fromPartial(e)) || [];
    message.inTxTrackerList = object.inTxTrackerList?.map(e => InTxTracker.fromPartial(e)) || [];
    message.zetaAccounting = object.zetaAccounting !== undefined && object.zetaAccounting !== null ? ZetaAccounting.fromPartial(object.zetaAccounting) : undefined;
    message.FinalizedInbounds = object.FinalizedInbounds?.map(e => e) || [];
    return message;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    const message = createBaseGenesisState();
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromAmino(object.params);
    }
    message.outTxTrackerList = object.outTxTrackerList?.map(e => OutTxTracker.fromAmino(e)) || [];
    message.gasPriceList = object.gasPriceList?.map(e => GasPrice.fromAmino(e)) || [];
    message.CrossChainTxs = object.CrossChainTxs?.map(e => CrossChainTx.fromAmino(e)) || [];
    message.lastBlockHeightList = object.lastBlockHeightList?.map(e => LastBlockHeight.fromAmino(e)) || [];
    message.inTxHashToCctxList = object.inTxHashToCctxList?.map(e => InTxHashToCctx.fromAmino(e)) || [];
    message.inTxTrackerList = object.in_tx_tracker_list?.map(e => InTxTracker.fromAmino(e)) || [];
    if (object.zeta_accounting !== undefined && object.zeta_accounting !== null) {
      message.zetaAccounting = ZetaAccounting.fromAmino(object.zeta_accounting);
    }
    message.FinalizedInbounds = object.FinalizedInbounds?.map(e => e) || [];
    return message;
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    if (message.outTxTrackerList) {
      obj.outTxTrackerList = message.outTxTrackerList.map(e => e ? OutTxTracker.toAmino(e) : undefined);
    } else {
      obj.outTxTrackerList = [];
    }
    if (message.gasPriceList) {
      obj.gasPriceList = message.gasPriceList.map(e => e ? GasPrice.toAmino(e) : undefined);
    } else {
      obj.gasPriceList = [];
    }
    if (message.CrossChainTxs) {
      obj.CrossChainTxs = message.CrossChainTxs.map(e => e ? CrossChainTx.toAmino(e) : undefined);
    } else {
      obj.CrossChainTxs = [];
    }
    if (message.lastBlockHeightList) {
      obj.lastBlockHeightList = message.lastBlockHeightList.map(e => e ? LastBlockHeight.toAmino(e) : undefined);
    } else {
      obj.lastBlockHeightList = [];
    }
    if (message.inTxHashToCctxList) {
      obj.inTxHashToCctxList = message.inTxHashToCctxList.map(e => e ? InTxHashToCctx.toAmino(e) : undefined);
    } else {
      obj.inTxHashToCctxList = [];
    }
    if (message.inTxTrackerList) {
      obj.in_tx_tracker_list = message.inTxTrackerList.map(e => e ? InTxTracker.toAmino(e) : undefined);
    } else {
      obj.in_tx_tracker_list = [];
    }
    obj.zeta_accounting = message.zetaAccounting ? ZetaAccounting.toAmino(message.zetaAccounting) : undefined;
    if (message.FinalizedInbounds) {
      obj.FinalizedInbounds = message.FinalizedInbounds.map(e => e);
    } else {
      obj.FinalizedInbounds = [];
    }
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
      typeUrl: "/zetachain.zetacore.crosschain.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};