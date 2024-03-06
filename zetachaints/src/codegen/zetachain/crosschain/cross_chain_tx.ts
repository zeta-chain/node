import { CoinType, coinTypeFromJSON } from "../common/common";
import { BinaryReader, BinaryWriter } from "../../binary";
export enum CctxStatus {
  /** PendingInbound - some observer sees inbound tx */
  PendingInbound = 0,
  /** PendingOutbound - super majority observer see inbound tx */
  PendingOutbound = 1,
  /** OutboundMined - the corresponding outbound tx is mined */
  OutboundMined = 3,
  /** PendingRevert - outbound cannot succeed; should revert inbound */
  PendingRevert = 4,
  /** Reverted - inbound reverted. */
  Reverted = 5,
  /** Aborted - inbound tx error or invalid paramters and cannot revert; just abort. But the amount can be refunded to zetachain using and admin proposal */
  Aborted = 6,
  UNRECOGNIZED = -1,
}
export const CctxStatusSDKType = CctxStatus;
export const CctxStatusAmino = CctxStatus;
export function cctxStatusFromJSON(object: any): CctxStatus {
  switch (object) {
    case 0:
    case "PendingInbound":
      return CctxStatus.PendingInbound;
    case 1:
    case "PendingOutbound":
      return CctxStatus.PendingOutbound;
    case 3:
    case "OutboundMined":
      return CctxStatus.OutboundMined;
    case 4:
    case "PendingRevert":
      return CctxStatus.PendingRevert;
    case 5:
    case "Reverted":
      return CctxStatus.Reverted;
    case 6:
    case "Aborted":
      return CctxStatus.Aborted;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CctxStatus.UNRECOGNIZED;
  }
}
export function cctxStatusToJSON(object: CctxStatus): string {
  switch (object) {
    case CctxStatus.PendingInbound:
      return "PendingInbound";
    case CctxStatus.PendingOutbound:
      return "PendingOutbound";
    case CctxStatus.OutboundMined:
      return "OutboundMined";
    case CctxStatus.PendingRevert:
      return "PendingRevert";
    case CctxStatus.Reverted:
      return "Reverted";
    case CctxStatus.Aborted:
      return "Aborted";
    case CctxStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export enum TxFinalizationStatus {
  /** NotFinalized - the corresponding tx is not finalized */
  NotFinalized = 0,
  /** Finalized - the corresponding tx is finalized but not executed yet */
  Finalized = 1,
  /** Executed - the corresponding tx is executed */
  Executed = 2,
  UNRECOGNIZED = -1,
}
export const TxFinalizationStatusSDKType = TxFinalizationStatus;
export const TxFinalizationStatusAmino = TxFinalizationStatus;
export function txFinalizationStatusFromJSON(object: any): TxFinalizationStatus {
  switch (object) {
    case 0:
    case "NotFinalized":
      return TxFinalizationStatus.NotFinalized;
    case 1:
    case "Finalized":
      return TxFinalizationStatus.Finalized;
    case 2:
    case "Executed":
      return TxFinalizationStatus.Executed;
    case -1:
    case "UNRECOGNIZED":
    default:
      return TxFinalizationStatus.UNRECOGNIZED;
  }
}
export function txFinalizationStatusToJSON(object: TxFinalizationStatus): string {
  switch (object) {
    case TxFinalizationStatus.NotFinalized:
      return "NotFinalized";
    case TxFinalizationStatus.Finalized:
      return "Finalized";
    case TxFinalizationStatus.Executed:
      return "Executed";
    case TxFinalizationStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export interface InboundTxParams {
  /** this address is the immediate contract/EOA that calls the Connector.send() */
  sender: string;
  senderChainId: bigint;
  /** this address is the EOA that signs the inbound tx */
  txOrigin: string;
  coinType: CoinType;
  /** for ERC20 coin type, the asset is an address of the ERC20 contract */
  asset: string;
  amount: string;
  inboundTxObservedHash: string;
  inboundTxObservedExternalHeight: bigint;
  inboundTxBallotIndex: string;
  inboundTxFinalizedZetaHeight: bigint;
  txFinalizationStatus: TxFinalizationStatus;
}
export interface InboundTxParamsProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.InboundTxParams";
  value: Uint8Array;
}
export interface InboundTxParamsAmino {
  /** this address is the immediate contract/EOA that calls the Connector.send() */
  sender?: string;
  sender_chain_id?: string;
  /** this address is the EOA that signs the inbound tx */
  tx_origin?: string;
  coin_type?: CoinType;
  /** for ERC20 coin type, the asset is an address of the ERC20 contract */
  asset?: string;
  amount?: string;
  inbound_tx_observed_hash?: string;
  inbound_tx_observed_external_height?: string;
  inbound_tx_ballot_index?: string;
  inbound_tx_finalized_zeta_height?: string;
  tx_finalization_status?: TxFinalizationStatus;
}
export interface InboundTxParamsAminoMsg {
  type: "/zetachain.zetacore.crosschain.InboundTxParams";
  value: InboundTxParamsAmino;
}
export interface InboundTxParamsSDKType {
  sender: string;
  sender_chain_id: bigint;
  tx_origin: string;
  coin_type: CoinType;
  asset: string;
  amount: string;
  inbound_tx_observed_hash: string;
  inbound_tx_observed_external_height: bigint;
  inbound_tx_ballot_index: string;
  inbound_tx_finalized_zeta_height: bigint;
  tx_finalization_status: TxFinalizationStatus;
}
export interface ZetaAccounting {
  /** aborted_zeta_amount stores the total aborted amount for cctx of coin-type ZETA */
  abortedZetaAmount: string;
}
export interface ZetaAccountingProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.ZetaAccounting";
  value: Uint8Array;
}
export interface ZetaAccountingAmino {
  /** aborted_zeta_amount stores the total aborted amount for cctx of coin-type ZETA */
  aborted_zeta_amount?: string;
}
export interface ZetaAccountingAminoMsg {
  type: "/zetachain.zetacore.crosschain.ZetaAccounting";
  value: ZetaAccountingAmino;
}
export interface ZetaAccountingSDKType {
  aborted_zeta_amount: string;
}
export interface OutboundTxParams {
  receiver: string;
  receiverChainId: bigint;
  coinType: CoinType;
  amount: string;
  outboundTxTssNonce: bigint;
  outboundTxGasLimit: bigint;
  outboundTxGasPrice: string;
  /**
   * the above are commands for zetaclients
   * the following fields are used when the outbound tx is mined
   */
  outboundTxHash: string;
  outboundTxBallotIndex: string;
  outboundTxObservedExternalHeight: bigint;
  outboundTxGasUsed: bigint;
  outboundTxEffectiveGasPrice: string;
  outboundTxEffectiveGasLimit: bigint;
  tssPubkey: string;
  txFinalizationStatus: TxFinalizationStatus;
}
export interface OutboundTxParamsProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.OutboundTxParams";
  value: Uint8Array;
}
export interface OutboundTxParamsAmino {
  receiver?: string;
  receiver_chainId?: string;
  coin_type?: CoinType;
  amount?: string;
  outbound_tx_tss_nonce?: string;
  outbound_tx_gas_limit?: string;
  outbound_tx_gas_price?: string;
  /**
   * the above are commands for zetaclients
   * the following fields are used when the outbound tx is mined
   */
  outbound_tx_hash?: string;
  outbound_tx_ballot_index?: string;
  outbound_tx_observed_external_height?: string;
  outbound_tx_gas_used?: string;
  outbound_tx_effective_gas_price?: string;
  outbound_tx_effective_gas_limit?: string;
  tss_pubkey?: string;
  tx_finalization_status?: TxFinalizationStatus;
}
export interface OutboundTxParamsAminoMsg {
  type: "/zetachain.zetacore.crosschain.OutboundTxParams";
  value: OutboundTxParamsAmino;
}
export interface OutboundTxParamsSDKType {
  receiver: string;
  receiver_chainId: bigint;
  coin_type: CoinType;
  amount: string;
  outbound_tx_tss_nonce: bigint;
  outbound_tx_gas_limit: bigint;
  outbound_tx_gas_price: string;
  outbound_tx_hash: string;
  outbound_tx_ballot_index: string;
  outbound_tx_observed_external_height: bigint;
  outbound_tx_gas_used: bigint;
  outbound_tx_effective_gas_price: string;
  outbound_tx_effective_gas_limit: bigint;
  tss_pubkey: string;
  tx_finalization_status: TxFinalizationStatus;
}
export interface Status {
  status: CctxStatus;
  statusMessage: string;
  lastUpdateTimestamp: bigint;
  isAbortRefunded: boolean;
}
export interface StatusProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.Status";
  value: Uint8Array;
}
export interface StatusAmino {
  status?: CctxStatus;
  status_message?: string;
  lastUpdate_timestamp?: string;
  isAbortRefunded?: boolean;
}
export interface StatusAminoMsg {
  type: "/zetachain.zetacore.crosschain.Status";
  value: StatusAmino;
}
export interface StatusSDKType {
  status: CctxStatus;
  status_message: string;
  lastUpdate_timestamp: bigint;
  isAbortRefunded: boolean;
}
export interface CrossChainTx {
  creator: string;
  index: string;
  zetaFees: string;
  /** Not used by protocol , just relayed across */
  relayedMessage: string;
  cctxStatus?: Status;
  inboundTxParams?: InboundTxParams;
  outboundTxParams: OutboundTxParams[];
}
export interface CrossChainTxProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.CrossChainTx";
  value: Uint8Array;
}
export interface CrossChainTxAmino {
  creator?: string;
  index?: string;
  zeta_fees?: string;
  /** Not used by protocol , just relayed across */
  relayed_message?: string;
  cctx_status?: StatusAmino;
  inbound_tx_params?: InboundTxParamsAmino;
  outbound_tx_params?: OutboundTxParamsAmino[];
}
export interface CrossChainTxAminoMsg {
  type: "/zetachain.zetacore.crosschain.CrossChainTx";
  value: CrossChainTxAmino;
}
export interface CrossChainTxSDKType {
  creator: string;
  index: string;
  zeta_fees: string;
  relayed_message: string;
  cctx_status?: StatusSDKType;
  inbound_tx_params?: InboundTxParamsSDKType;
  outbound_tx_params: OutboundTxParamsSDKType[];
}
function createBaseInboundTxParams(): InboundTxParams {
  return {
    sender: "",
    senderChainId: BigInt(0),
    txOrigin: "",
    coinType: 0,
    asset: "",
    amount: "",
    inboundTxObservedHash: "",
    inboundTxObservedExternalHeight: BigInt(0),
    inboundTxBallotIndex: "",
    inboundTxFinalizedZetaHeight: BigInt(0),
    txFinalizationStatus: 0
  };
}
export const InboundTxParams = {
  typeUrl: "/zetachain.zetacore.crosschain.InboundTxParams",
  encode(message: InboundTxParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.sender !== "") {
      writer.uint32(10).string(message.sender);
    }
    if (message.senderChainId !== BigInt(0)) {
      writer.uint32(16).int64(message.senderChainId);
    }
    if (message.txOrigin !== "") {
      writer.uint32(26).string(message.txOrigin);
    }
    if (message.coinType !== 0) {
      writer.uint32(32).int32(message.coinType);
    }
    if (message.asset !== "") {
      writer.uint32(42).string(message.asset);
    }
    if (message.amount !== "") {
      writer.uint32(50).string(message.amount);
    }
    if (message.inboundTxObservedHash !== "") {
      writer.uint32(58).string(message.inboundTxObservedHash);
    }
    if (message.inboundTxObservedExternalHeight !== BigInt(0)) {
      writer.uint32(64).uint64(message.inboundTxObservedExternalHeight);
    }
    if (message.inboundTxBallotIndex !== "") {
      writer.uint32(74).string(message.inboundTxBallotIndex);
    }
    if (message.inboundTxFinalizedZetaHeight !== BigInt(0)) {
      writer.uint32(80).uint64(message.inboundTxFinalizedZetaHeight);
    }
    if (message.txFinalizationStatus !== 0) {
      writer.uint32(88).int32(message.txFinalizationStatus);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): InboundTxParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInboundTxParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sender = reader.string();
          break;
        case 2:
          message.senderChainId = reader.int64();
          break;
        case 3:
          message.txOrigin = reader.string();
          break;
        case 4:
          message.coinType = (reader.int32() as any);
          break;
        case 5:
          message.asset = reader.string();
          break;
        case 6:
          message.amount = reader.string();
          break;
        case 7:
          message.inboundTxObservedHash = reader.string();
          break;
        case 8:
          message.inboundTxObservedExternalHeight = reader.uint64();
          break;
        case 9:
          message.inboundTxBallotIndex = reader.string();
          break;
        case 10:
          message.inboundTxFinalizedZetaHeight = reader.uint64();
          break;
        case 11:
          message.txFinalizationStatus = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<InboundTxParams>): InboundTxParams {
    const message = createBaseInboundTxParams();
    message.sender = object.sender ?? "";
    message.senderChainId = object.senderChainId !== undefined && object.senderChainId !== null ? BigInt(object.senderChainId.toString()) : BigInt(0);
    message.txOrigin = object.txOrigin ?? "";
    message.coinType = object.coinType ?? 0;
    message.asset = object.asset ?? "";
    message.amount = object.amount ?? "";
    message.inboundTxObservedHash = object.inboundTxObservedHash ?? "";
    message.inboundTxObservedExternalHeight = object.inboundTxObservedExternalHeight !== undefined && object.inboundTxObservedExternalHeight !== null ? BigInt(object.inboundTxObservedExternalHeight.toString()) : BigInt(0);
    message.inboundTxBallotIndex = object.inboundTxBallotIndex ?? "";
    message.inboundTxFinalizedZetaHeight = object.inboundTxFinalizedZetaHeight !== undefined && object.inboundTxFinalizedZetaHeight !== null ? BigInt(object.inboundTxFinalizedZetaHeight.toString()) : BigInt(0);
    message.txFinalizationStatus = object.txFinalizationStatus ?? 0;
    return message;
  },
  fromAmino(object: InboundTxParamsAmino): InboundTxParams {
    const message = createBaseInboundTxParams();
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = object.sender;
    }
    if (object.sender_chain_id !== undefined && object.sender_chain_id !== null) {
      message.senderChainId = BigInt(object.sender_chain_id);
    }
    if (object.tx_origin !== undefined && object.tx_origin !== null) {
      message.txOrigin = object.tx_origin;
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    if (object.asset !== undefined && object.asset !== null) {
      message.asset = object.asset;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    if (object.inbound_tx_observed_hash !== undefined && object.inbound_tx_observed_hash !== null) {
      message.inboundTxObservedHash = object.inbound_tx_observed_hash;
    }
    if (object.inbound_tx_observed_external_height !== undefined && object.inbound_tx_observed_external_height !== null) {
      message.inboundTxObservedExternalHeight = BigInt(object.inbound_tx_observed_external_height);
    }
    if (object.inbound_tx_ballot_index !== undefined && object.inbound_tx_ballot_index !== null) {
      message.inboundTxBallotIndex = object.inbound_tx_ballot_index;
    }
    if (object.inbound_tx_finalized_zeta_height !== undefined && object.inbound_tx_finalized_zeta_height !== null) {
      message.inboundTxFinalizedZetaHeight = BigInt(object.inbound_tx_finalized_zeta_height);
    }
    if (object.tx_finalization_status !== undefined && object.tx_finalization_status !== null) {
      message.txFinalizationStatus = txFinalizationStatusFromJSON(object.tx_finalization_status);
    }
    return message;
  },
  toAmino(message: InboundTxParams): InboundTxParamsAmino {
    const obj: any = {};
    obj.sender = message.sender;
    obj.sender_chain_id = message.senderChainId ? message.senderChainId.toString() : undefined;
    obj.tx_origin = message.txOrigin;
    obj.coin_type = message.coinType;
    obj.asset = message.asset;
    obj.amount = message.amount;
    obj.inbound_tx_observed_hash = message.inboundTxObservedHash;
    obj.inbound_tx_observed_external_height = message.inboundTxObservedExternalHeight ? message.inboundTxObservedExternalHeight.toString() : undefined;
    obj.inbound_tx_ballot_index = message.inboundTxBallotIndex;
    obj.inbound_tx_finalized_zeta_height = message.inboundTxFinalizedZetaHeight ? message.inboundTxFinalizedZetaHeight.toString() : undefined;
    obj.tx_finalization_status = message.txFinalizationStatus;
    return obj;
  },
  fromAminoMsg(object: InboundTxParamsAminoMsg): InboundTxParams {
    return InboundTxParams.fromAmino(object.value);
  },
  fromProtoMsg(message: InboundTxParamsProtoMsg): InboundTxParams {
    return InboundTxParams.decode(message.value);
  },
  toProto(message: InboundTxParams): Uint8Array {
    return InboundTxParams.encode(message).finish();
  },
  toProtoMsg(message: InboundTxParams): InboundTxParamsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.InboundTxParams",
      value: InboundTxParams.encode(message).finish()
    };
  }
};
function createBaseZetaAccounting(): ZetaAccounting {
  return {
    abortedZetaAmount: ""
  };
}
export const ZetaAccounting = {
  typeUrl: "/zetachain.zetacore.crosschain.ZetaAccounting",
  encode(message: ZetaAccounting, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.abortedZetaAmount !== "") {
      writer.uint32(10).string(message.abortedZetaAmount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ZetaAccounting {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseZetaAccounting();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.abortedZetaAmount = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ZetaAccounting>): ZetaAccounting {
    const message = createBaseZetaAccounting();
    message.abortedZetaAmount = object.abortedZetaAmount ?? "";
    return message;
  },
  fromAmino(object: ZetaAccountingAmino): ZetaAccounting {
    const message = createBaseZetaAccounting();
    if (object.aborted_zeta_amount !== undefined && object.aborted_zeta_amount !== null) {
      message.abortedZetaAmount = object.aborted_zeta_amount;
    }
    return message;
  },
  toAmino(message: ZetaAccounting): ZetaAccountingAmino {
    const obj: any = {};
    obj.aborted_zeta_amount = message.abortedZetaAmount;
    return obj;
  },
  fromAminoMsg(object: ZetaAccountingAminoMsg): ZetaAccounting {
    return ZetaAccounting.fromAmino(object.value);
  },
  fromProtoMsg(message: ZetaAccountingProtoMsg): ZetaAccounting {
    return ZetaAccounting.decode(message.value);
  },
  toProto(message: ZetaAccounting): Uint8Array {
    return ZetaAccounting.encode(message).finish();
  },
  toProtoMsg(message: ZetaAccounting): ZetaAccountingProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.ZetaAccounting",
      value: ZetaAccounting.encode(message).finish()
    };
  }
};
function createBaseOutboundTxParams(): OutboundTxParams {
  return {
    receiver: "",
    receiverChainId: BigInt(0),
    coinType: 0,
    amount: "",
    outboundTxTssNonce: BigInt(0),
    outboundTxGasLimit: BigInt(0),
    outboundTxGasPrice: "",
    outboundTxHash: "",
    outboundTxBallotIndex: "",
    outboundTxObservedExternalHeight: BigInt(0),
    outboundTxGasUsed: BigInt(0),
    outboundTxEffectiveGasPrice: "",
    outboundTxEffectiveGasLimit: BigInt(0),
    tssPubkey: "",
    txFinalizationStatus: 0
  };
}
export const OutboundTxParams = {
  typeUrl: "/zetachain.zetacore.crosschain.OutboundTxParams",
  encode(message: OutboundTxParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.receiver !== "") {
      writer.uint32(10).string(message.receiver);
    }
    if (message.receiverChainId !== BigInt(0)) {
      writer.uint32(16).int64(message.receiverChainId);
    }
    if (message.coinType !== 0) {
      writer.uint32(24).int32(message.coinType);
    }
    if (message.amount !== "") {
      writer.uint32(34).string(message.amount);
    }
    if (message.outboundTxTssNonce !== BigInt(0)) {
      writer.uint32(40).uint64(message.outboundTxTssNonce);
    }
    if (message.outboundTxGasLimit !== BigInt(0)) {
      writer.uint32(48).uint64(message.outboundTxGasLimit);
    }
    if (message.outboundTxGasPrice !== "") {
      writer.uint32(58).string(message.outboundTxGasPrice);
    }
    if (message.outboundTxHash !== "") {
      writer.uint32(66).string(message.outboundTxHash);
    }
    if (message.outboundTxBallotIndex !== "") {
      writer.uint32(74).string(message.outboundTxBallotIndex);
    }
    if (message.outboundTxObservedExternalHeight !== BigInt(0)) {
      writer.uint32(80).uint64(message.outboundTxObservedExternalHeight);
    }
    if (message.outboundTxGasUsed !== BigInt(0)) {
      writer.uint32(160).uint64(message.outboundTxGasUsed);
    }
    if (message.outboundTxEffectiveGasPrice !== "") {
      writer.uint32(170).string(message.outboundTxEffectiveGasPrice);
    }
    if (message.outboundTxEffectiveGasLimit !== BigInt(0)) {
      writer.uint32(176).uint64(message.outboundTxEffectiveGasLimit);
    }
    if (message.tssPubkey !== "") {
      writer.uint32(90).string(message.tssPubkey);
    }
    if (message.txFinalizationStatus !== 0) {
      writer.uint32(96).int32(message.txFinalizationStatus);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OutboundTxParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOutboundTxParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.receiver = reader.string();
          break;
        case 2:
          message.receiverChainId = reader.int64();
          break;
        case 3:
          message.coinType = (reader.int32() as any);
          break;
        case 4:
          message.amount = reader.string();
          break;
        case 5:
          message.outboundTxTssNonce = reader.uint64();
          break;
        case 6:
          message.outboundTxGasLimit = reader.uint64();
          break;
        case 7:
          message.outboundTxGasPrice = reader.string();
          break;
        case 8:
          message.outboundTxHash = reader.string();
          break;
        case 9:
          message.outboundTxBallotIndex = reader.string();
          break;
        case 10:
          message.outboundTxObservedExternalHeight = reader.uint64();
          break;
        case 20:
          message.outboundTxGasUsed = reader.uint64();
          break;
        case 21:
          message.outboundTxEffectiveGasPrice = reader.string();
          break;
        case 22:
          message.outboundTxEffectiveGasLimit = reader.uint64();
          break;
        case 11:
          message.tssPubkey = reader.string();
          break;
        case 12:
          message.txFinalizationStatus = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OutboundTxParams>): OutboundTxParams {
    const message = createBaseOutboundTxParams();
    message.receiver = object.receiver ?? "";
    message.receiverChainId = object.receiverChainId !== undefined && object.receiverChainId !== null ? BigInt(object.receiverChainId.toString()) : BigInt(0);
    message.coinType = object.coinType ?? 0;
    message.amount = object.amount ?? "";
    message.outboundTxTssNonce = object.outboundTxTssNonce !== undefined && object.outboundTxTssNonce !== null ? BigInt(object.outboundTxTssNonce.toString()) : BigInt(0);
    message.outboundTxGasLimit = object.outboundTxGasLimit !== undefined && object.outboundTxGasLimit !== null ? BigInt(object.outboundTxGasLimit.toString()) : BigInt(0);
    message.outboundTxGasPrice = object.outboundTxGasPrice ?? "";
    message.outboundTxHash = object.outboundTxHash ?? "";
    message.outboundTxBallotIndex = object.outboundTxBallotIndex ?? "";
    message.outboundTxObservedExternalHeight = object.outboundTxObservedExternalHeight !== undefined && object.outboundTxObservedExternalHeight !== null ? BigInt(object.outboundTxObservedExternalHeight.toString()) : BigInt(0);
    message.outboundTxGasUsed = object.outboundTxGasUsed !== undefined && object.outboundTxGasUsed !== null ? BigInt(object.outboundTxGasUsed.toString()) : BigInt(0);
    message.outboundTxEffectiveGasPrice = object.outboundTxEffectiveGasPrice ?? "";
    message.outboundTxEffectiveGasLimit = object.outboundTxEffectiveGasLimit !== undefined && object.outboundTxEffectiveGasLimit !== null ? BigInt(object.outboundTxEffectiveGasLimit.toString()) : BigInt(0);
    message.tssPubkey = object.tssPubkey ?? "";
    message.txFinalizationStatus = object.txFinalizationStatus ?? 0;
    return message;
  },
  fromAmino(object: OutboundTxParamsAmino): OutboundTxParams {
    const message = createBaseOutboundTxParams();
    if (object.receiver !== undefined && object.receiver !== null) {
      message.receiver = object.receiver;
    }
    if (object.receiver_chainId !== undefined && object.receiver_chainId !== null) {
      message.receiverChainId = BigInt(object.receiver_chainId);
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    if (object.outbound_tx_tss_nonce !== undefined && object.outbound_tx_tss_nonce !== null) {
      message.outboundTxTssNonce = BigInt(object.outbound_tx_tss_nonce);
    }
    if (object.outbound_tx_gas_limit !== undefined && object.outbound_tx_gas_limit !== null) {
      message.outboundTxGasLimit = BigInt(object.outbound_tx_gas_limit);
    }
    if (object.outbound_tx_gas_price !== undefined && object.outbound_tx_gas_price !== null) {
      message.outboundTxGasPrice = object.outbound_tx_gas_price;
    }
    if (object.outbound_tx_hash !== undefined && object.outbound_tx_hash !== null) {
      message.outboundTxHash = object.outbound_tx_hash;
    }
    if (object.outbound_tx_ballot_index !== undefined && object.outbound_tx_ballot_index !== null) {
      message.outboundTxBallotIndex = object.outbound_tx_ballot_index;
    }
    if (object.outbound_tx_observed_external_height !== undefined && object.outbound_tx_observed_external_height !== null) {
      message.outboundTxObservedExternalHeight = BigInt(object.outbound_tx_observed_external_height);
    }
    if (object.outbound_tx_gas_used !== undefined && object.outbound_tx_gas_used !== null) {
      message.outboundTxGasUsed = BigInt(object.outbound_tx_gas_used);
    }
    if (object.outbound_tx_effective_gas_price !== undefined && object.outbound_tx_effective_gas_price !== null) {
      message.outboundTxEffectiveGasPrice = object.outbound_tx_effective_gas_price;
    }
    if (object.outbound_tx_effective_gas_limit !== undefined && object.outbound_tx_effective_gas_limit !== null) {
      message.outboundTxEffectiveGasLimit = BigInt(object.outbound_tx_effective_gas_limit);
    }
    if (object.tss_pubkey !== undefined && object.tss_pubkey !== null) {
      message.tssPubkey = object.tss_pubkey;
    }
    if (object.tx_finalization_status !== undefined && object.tx_finalization_status !== null) {
      message.txFinalizationStatus = txFinalizationStatusFromJSON(object.tx_finalization_status);
    }
    return message;
  },
  toAmino(message: OutboundTxParams): OutboundTxParamsAmino {
    const obj: any = {};
    obj.receiver = message.receiver;
    obj.receiver_chainId = message.receiverChainId ? message.receiverChainId.toString() : undefined;
    obj.coin_type = message.coinType;
    obj.amount = message.amount;
    obj.outbound_tx_tss_nonce = message.outboundTxTssNonce ? message.outboundTxTssNonce.toString() : undefined;
    obj.outbound_tx_gas_limit = message.outboundTxGasLimit ? message.outboundTxGasLimit.toString() : undefined;
    obj.outbound_tx_gas_price = message.outboundTxGasPrice;
    obj.outbound_tx_hash = message.outboundTxHash;
    obj.outbound_tx_ballot_index = message.outboundTxBallotIndex;
    obj.outbound_tx_observed_external_height = message.outboundTxObservedExternalHeight ? message.outboundTxObservedExternalHeight.toString() : undefined;
    obj.outbound_tx_gas_used = message.outboundTxGasUsed ? message.outboundTxGasUsed.toString() : undefined;
    obj.outbound_tx_effective_gas_price = message.outboundTxEffectiveGasPrice;
    obj.outbound_tx_effective_gas_limit = message.outboundTxEffectiveGasLimit ? message.outboundTxEffectiveGasLimit.toString() : undefined;
    obj.tss_pubkey = message.tssPubkey;
    obj.tx_finalization_status = message.txFinalizationStatus;
    return obj;
  },
  fromAminoMsg(object: OutboundTxParamsAminoMsg): OutboundTxParams {
    return OutboundTxParams.fromAmino(object.value);
  },
  fromProtoMsg(message: OutboundTxParamsProtoMsg): OutboundTxParams {
    return OutboundTxParams.decode(message.value);
  },
  toProto(message: OutboundTxParams): Uint8Array {
    return OutboundTxParams.encode(message).finish();
  },
  toProtoMsg(message: OutboundTxParams): OutboundTxParamsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.OutboundTxParams",
      value: OutboundTxParams.encode(message).finish()
    };
  }
};
function createBaseStatus(): Status {
  return {
    status: 0,
    statusMessage: "",
    lastUpdateTimestamp: BigInt(0),
    isAbortRefunded: false
  };
}
export const Status = {
  typeUrl: "/zetachain.zetacore.crosschain.Status",
  encode(message: Status, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.statusMessage !== "") {
      writer.uint32(18).string(message.statusMessage);
    }
    if (message.lastUpdateTimestamp !== BigInt(0)) {
      writer.uint32(24).int64(message.lastUpdateTimestamp);
    }
    if (message.isAbortRefunded === true) {
      writer.uint32(32).bool(message.isAbortRefunded);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Status {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.status = (reader.int32() as any);
          break;
        case 2:
          message.statusMessage = reader.string();
          break;
        case 3:
          message.lastUpdateTimestamp = reader.int64();
          break;
        case 4:
          message.isAbortRefunded = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Status>): Status {
    const message = createBaseStatus();
    message.status = object.status ?? 0;
    message.statusMessage = object.statusMessage ?? "";
    message.lastUpdateTimestamp = object.lastUpdateTimestamp !== undefined && object.lastUpdateTimestamp !== null ? BigInt(object.lastUpdateTimestamp.toString()) : BigInt(0);
    message.isAbortRefunded = object.isAbortRefunded ?? false;
    return message;
  },
  fromAmino(object: StatusAmino): Status {
    const message = createBaseStatus();
    if (object.status !== undefined && object.status !== null) {
      message.status = cctxStatusFromJSON(object.status);
    }
    if (object.status_message !== undefined && object.status_message !== null) {
      message.statusMessage = object.status_message;
    }
    if (object.lastUpdate_timestamp !== undefined && object.lastUpdate_timestamp !== null) {
      message.lastUpdateTimestamp = BigInt(object.lastUpdate_timestamp);
    }
    if (object.isAbortRefunded !== undefined && object.isAbortRefunded !== null) {
      message.isAbortRefunded = object.isAbortRefunded;
    }
    return message;
  },
  toAmino(message: Status): StatusAmino {
    const obj: any = {};
    obj.status = message.status;
    obj.status_message = message.statusMessage;
    obj.lastUpdate_timestamp = message.lastUpdateTimestamp ? message.lastUpdateTimestamp.toString() : undefined;
    obj.isAbortRefunded = message.isAbortRefunded;
    return obj;
  },
  fromAminoMsg(object: StatusAminoMsg): Status {
    return Status.fromAmino(object.value);
  },
  fromProtoMsg(message: StatusProtoMsg): Status {
    return Status.decode(message.value);
  },
  toProto(message: Status): Uint8Array {
    return Status.encode(message).finish();
  },
  toProtoMsg(message: Status): StatusProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.Status",
      value: Status.encode(message).finish()
    };
  }
};
function createBaseCrossChainTx(): CrossChainTx {
  return {
    creator: "",
    index: "",
    zetaFees: "",
    relayedMessage: "",
    cctxStatus: undefined,
    inboundTxParams: undefined,
    outboundTxParams: []
  };
}
export const CrossChainTx = {
  typeUrl: "/zetachain.zetacore.crosschain.CrossChainTx",
  encode(message: CrossChainTx, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.index !== "") {
      writer.uint32(18).string(message.index);
    }
    if (message.zetaFees !== "") {
      writer.uint32(42).string(message.zetaFees);
    }
    if (message.relayedMessage !== "") {
      writer.uint32(50).string(message.relayedMessage);
    }
    if (message.cctxStatus !== undefined) {
      Status.encode(message.cctxStatus, writer.uint32(66).fork()).ldelim();
    }
    if (message.inboundTxParams !== undefined) {
      InboundTxParams.encode(message.inboundTxParams, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.outboundTxParams) {
      OutboundTxParams.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): CrossChainTx {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCrossChainTx();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.index = reader.string();
          break;
        case 5:
          message.zetaFees = reader.string();
          break;
        case 6:
          message.relayedMessage = reader.string();
          break;
        case 8:
          message.cctxStatus = Status.decode(reader, reader.uint32());
          break;
        case 9:
          message.inboundTxParams = InboundTxParams.decode(reader, reader.uint32());
          break;
        case 10:
          message.outboundTxParams.push(OutboundTxParams.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<CrossChainTx>): CrossChainTx {
    const message = createBaseCrossChainTx();
    message.creator = object.creator ?? "";
    message.index = object.index ?? "";
    message.zetaFees = object.zetaFees ?? "";
    message.relayedMessage = object.relayedMessage ?? "";
    message.cctxStatus = object.cctxStatus !== undefined && object.cctxStatus !== null ? Status.fromPartial(object.cctxStatus) : undefined;
    message.inboundTxParams = object.inboundTxParams !== undefined && object.inboundTxParams !== null ? InboundTxParams.fromPartial(object.inboundTxParams) : undefined;
    message.outboundTxParams = object.outboundTxParams?.map(e => OutboundTxParams.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: CrossChainTxAmino): CrossChainTx {
    const message = createBaseCrossChainTx();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    if (object.zeta_fees !== undefined && object.zeta_fees !== null) {
      message.zetaFees = object.zeta_fees;
    }
    if (object.relayed_message !== undefined && object.relayed_message !== null) {
      message.relayedMessage = object.relayed_message;
    }
    if (object.cctx_status !== undefined && object.cctx_status !== null) {
      message.cctxStatus = Status.fromAmino(object.cctx_status);
    }
    if (object.inbound_tx_params !== undefined && object.inbound_tx_params !== null) {
      message.inboundTxParams = InboundTxParams.fromAmino(object.inbound_tx_params);
    }
    message.outboundTxParams = object.outbound_tx_params?.map(e => OutboundTxParams.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: CrossChainTx): CrossChainTxAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.index = message.index;
    obj.zeta_fees = message.zetaFees;
    obj.relayed_message = message.relayedMessage;
    obj.cctx_status = message.cctxStatus ? Status.toAmino(message.cctxStatus) : undefined;
    obj.inbound_tx_params = message.inboundTxParams ? InboundTxParams.toAmino(message.inboundTxParams) : undefined;
    if (message.outboundTxParams) {
      obj.outbound_tx_params = message.outboundTxParams.map(e => e ? OutboundTxParams.toAmino(e) : undefined);
    } else {
      obj.outbound_tx_params = [];
    }
    return obj;
  },
  fromAminoMsg(object: CrossChainTxAminoMsg): CrossChainTx {
    return CrossChainTx.fromAmino(object.value);
  },
  fromProtoMsg(message: CrossChainTxProtoMsg): CrossChainTx {
    return CrossChainTx.decode(message.value);
  },
  toProto(message: CrossChainTx): Uint8Array {
    return CrossChainTx.encode(message).finish();
  },
  toProtoMsg(message: CrossChainTx): CrossChainTxProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.CrossChainTx",
      value: CrossChainTx.encode(message).finish()
    };
  }
};