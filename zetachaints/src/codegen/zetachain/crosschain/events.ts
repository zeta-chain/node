import { BinaryReader, BinaryWriter } from "../../binary";
export interface EventInboundFinalized {
  msgTypeUrl: string;
  cctxIndex: string;
  sender: string;
  txOrgin: string;
  asset: string;
  inTxHash: string;
  inBlockHeight: string;
  receiver: string;
  receiverChain: string;
  amount: string;
  relayedMessage: string;
  newStatus: string;
  statusMessage: string;
  senderChain: string;
}
export interface EventInboundFinalizedProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.EventInboundFinalized";
  value: Uint8Array;
}
export interface EventInboundFinalizedAmino {
  msg_type_url?: string;
  cctx_index?: string;
  sender?: string;
  tx_orgin?: string;
  asset?: string;
  in_tx_hash?: string;
  in_block_height?: string;
  receiver?: string;
  receiver_chain?: string;
  amount?: string;
  relayed_message?: string;
  new_status?: string;
  status_message?: string;
  sender_chain?: string;
}
export interface EventInboundFinalizedAminoMsg {
  type: "/zetachain.zetacore.crosschain.EventInboundFinalized";
  value: EventInboundFinalizedAmino;
}
export interface EventInboundFinalizedSDKType {
  msg_type_url: string;
  cctx_index: string;
  sender: string;
  tx_orgin: string;
  asset: string;
  in_tx_hash: string;
  in_block_height: string;
  receiver: string;
  receiver_chain: string;
  amount: string;
  relayed_message: string;
  new_status: string;
  status_message: string;
  sender_chain: string;
}
export interface EventZrcWithdrawCreated {
  msgTypeUrl: string;
  cctxIndex: string;
  sender: string;
  senderChain: string;
  inTxHash: string;
  receiver: string;
  receiverChain: string;
  amount: string;
  newStatus: string;
}
export interface EventZrcWithdrawCreatedProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.EventZrcWithdrawCreated";
  value: Uint8Array;
}
export interface EventZrcWithdrawCreatedAmino {
  msg_type_url?: string;
  cctx_index?: string;
  sender?: string;
  sender_chain?: string;
  in_tx_hash?: string;
  receiver?: string;
  receiver_chain?: string;
  amount?: string;
  new_status?: string;
}
export interface EventZrcWithdrawCreatedAminoMsg {
  type: "/zetachain.zetacore.crosschain.EventZrcWithdrawCreated";
  value: EventZrcWithdrawCreatedAmino;
}
export interface EventZrcWithdrawCreatedSDKType {
  msg_type_url: string;
  cctx_index: string;
  sender: string;
  sender_chain: string;
  in_tx_hash: string;
  receiver: string;
  receiver_chain: string;
  amount: string;
  new_status: string;
}
export interface EventZetaWithdrawCreated {
  msgTypeUrl: string;
  cctxIndex: string;
  sender: string;
  inTxHash: string;
  newStatus: string;
}
export interface EventZetaWithdrawCreatedProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.EventZetaWithdrawCreated";
  value: Uint8Array;
}
export interface EventZetaWithdrawCreatedAmino {
  msg_type_url?: string;
  cctx_index?: string;
  sender?: string;
  in_tx_hash?: string;
  new_status?: string;
}
export interface EventZetaWithdrawCreatedAminoMsg {
  type: "/zetachain.zetacore.crosschain.EventZetaWithdrawCreated";
  value: EventZetaWithdrawCreatedAmino;
}
export interface EventZetaWithdrawCreatedSDKType {
  msg_type_url: string;
  cctx_index: string;
  sender: string;
  in_tx_hash: string;
  new_status: string;
}
export interface EventOutboundFailure {
  msgTypeUrl: string;
  cctxIndex: string;
  oldStatus: string;
  newStatus: string;
  valueReceived: string;
}
export interface EventOutboundFailureProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.EventOutboundFailure";
  value: Uint8Array;
}
export interface EventOutboundFailureAmino {
  msg_type_url?: string;
  cctx_index?: string;
  old_status?: string;
  new_status?: string;
  value_received?: string;
}
export interface EventOutboundFailureAminoMsg {
  type: "/zetachain.zetacore.crosschain.EventOutboundFailure";
  value: EventOutboundFailureAmino;
}
export interface EventOutboundFailureSDKType {
  msg_type_url: string;
  cctx_index: string;
  old_status: string;
  new_status: string;
  value_received: string;
}
export interface EventOutboundSuccess {
  msgTypeUrl: string;
  cctxIndex: string;
  oldStatus: string;
  newStatus: string;
  valueReceived: string;
}
export interface EventOutboundSuccessProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.EventOutboundSuccess";
  value: Uint8Array;
}
export interface EventOutboundSuccessAmino {
  msg_type_url?: string;
  cctx_index?: string;
  old_status?: string;
  new_status?: string;
  value_received?: string;
}
export interface EventOutboundSuccessAminoMsg {
  type: "/zetachain.zetacore.crosschain.EventOutboundSuccess";
  value: EventOutboundSuccessAmino;
}
export interface EventOutboundSuccessSDKType {
  msg_type_url: string;
  cctx_index: string;
  old_status: string;
  new_status: string;
  value_received: string;
}
export interface EventCCTXGasPriceIncreased {
  cctxIndex: string;
  gasPriceIncrease: string;
  additionalFees: string;
}
export interface EventCCTXGasPriceIncreasedProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.EventCCTXGasPriceIncreased";
  value: Uint8Array;
}
export interface EventCCTXGasPriceIncreasedAmino {
  cctx_index?: string;
  gas_price_increase?: string;
  additional_fees?: string;
}
export interface EventCCTXGasPriceIncreasedAminoMsg {
  type: "/zetachain.zetacore.crosschain.EventCCTXGasPriceIncreased";
  value: EventCCTXGasPriceIncreasedAmino;
}
export interface EventCCTXGasPriceIncreasedSDKType {
  cctx_index: string;
  gas_price_increase: string;
  additional_fees: string;
}
function createBaseEventInboundFinalized(): EventInboundFinalized {
  return {
    msgTypeUrl: "",
    cctxIndex: "",
    sender: "",
    txOrgin: "",
    asset: "",
    inTxHash: "",
    inBlockHeight: "",
    receiver: "",
    receiverChain: "",
    amount: "",
    relayedMessage: "",
    newStatus: "",
    statusMessage: "",
    senderChain: ""
  };
}
export const EventInboundFinalized = {
  typeUrl: "/zetachain.zetacore.crosschain.EventInboundFinalized",
  encode(message: EventInboundFinalized, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.cctxIndex !== "") {
      writer.uint32(18).string(message.cctxIndex);
    }
    if (message.sender !== "") {
      writer.uint32(26).string(message.sender);
    }
    if (message.txOrgin !== "") {
      writer.uint32(34).string(message.txOrgin);
    }
    if (message.asset !== "") {
      writer.uint32(42).string(message.asset);
    }
    if (message.inTxHash !== "") {
      writer.uint32(50).string(message.inTxHash);
    }
    if (message.inBlockHeight !== "") {
      writer.uint32(58).string(message.inBlockHeight);
    }
    if (message.receiver !== "") {
      writer.uint32(66).string(message.receiver);
    }
    if (message.receiverChain !== "") {
      writer.uint32(74).string(message.receiverChain);
    }
    if (message.amount !== "") {
      writer.uint32(82).string(message.amount);
    }
    if (message.relayedMessage !== "") {
      writer.uint32(90).string(message.relayedMessage);
    }
    if (message.newStatus !== "") {
      writer.uint32(98).string(message.newStatus);
    }
    if (message.statusMessage !== "") {
      writer.uint32(106).string(message.statusMessage);
    }
    if (message.senderChain !== "") {
      writer.uint32(114).string(message.senderChain);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventInboundFinalized {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventInboundFinalized();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.cctxIndex = reader.string();
          break;
        case 3:
          message.sender = reader.string();
          break;
        case 4:
          message.txOrgin = reader.string();
          break;
        case 5:
          message.asset = reader.string();
          break;
        case 6:
          message.inTxHash = reader.string();
          break;
        case 7:
          message.inBlockHeight = reader.string();
          break;
        case 8:
          message.receiver = reader.string();
          break;
        case 9:
          message.receiverChain = reader.string();
          break;
        case 10:
          message.amount = reader.string();
          break;
        case 11:
          message.relayedMessage = reader.string();
          break;
        case 12:
          message.newStatus = reader.string();
          break;
        case 13:
          message.statusMessage = reader.string();
          break;
        case 14:
          message.senderChain = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventInboundFinalized>): EventInboundFinalized {
    const message = createBaseEventInboundFinalized();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.cctxIndex = object.cctxIndex ?? "";
    message.sender = object.sender ?? "";
    message.txOrgin = object.txOrgin ?? "";
    message.asset = object.asset ?? "";
    message.inTxHash = object.inTxHash ?? "";
    message.inBlockHeight = object.inBlockHeight ?? "";
    message.receiver = object.receiver ?? "";
    message.receiverChain = object.receiverChain ?? "";
    message.amount = object.amount ?? "";
    message.relayedMessage = object.relayedMessage ?? "";
    message.newStatus = object.newStatus ?? "";
    message.statusMessage = object.statusMessage ?? "";
    message.senderChain = object.senderChain ?? "";
    return message;
  },
  fromAmino(object: EventInboundFinalizedAmino): EventInboundFinalized {
    const message = createBaseEventInboundFinalized();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.cctx_index !== undefined && object.cctx_index !== null) {
      message.cctxIndex = object.cctx_index;
    }
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = object.sender;
    }
    if (object.tx_orgin !== undefined && object.tx_orgin !== null) {
      message.txOrgin = object.tx_orgin;
    }
    if (object.asset !== undefined && object.asset !== null) {
      message.asset = object.asset;
    }
    if (object.in_tx_hash !== undefined && object.in_tx_hash !== null) {
      message.inTxHash = object.in_tx_hash;
    }
    if (object.in_block_height !== undefined && object.in_block_height !== null) {
      message.inBlockHeight = object.in_block_height;
    }
    if (object.receiver !== undefined && object.receiver !== null) {
      message.receiver = object.receiver;
    }
    if (object.receiver_chain !== undefined && object.receiver_chain !== null) {
      message.receiverChain = object.receiver_chain;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    if (object.relayed_message !== undefined && object.relayed_message !== null) {
      message.relayedMessage = object.relayed_message;
    }
    if (object.new_status !== undefined && object.new_status !== null) {
      message.newStatus = object.new_status;
    }
    if (object.status_message !== undefined && object.status_message !== null) {
      message.statusMessage = object.status_message;
    }
    if (object.sender_chain !== undefined && object.sender_chain !== null) {
      message.senderChain = object.sender_chain;
    }
    return message;
  },
  toAmino(message: EventInboundFinalized): EventInboundFinalizedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.cctx_index = message.cctxIndex;
    obj.sender = message.sender;
    obj.tx_orgin = message.txOrgin;
    obj.asset = message.asset;
    obj.in_tx_hash = message.inTxHash;
    obj.in_block_height = message.inBlockHeight;
    obj.receiver = message.receiver;
    obj.receiver_chain = message.receiverChain;
    obj.amount = message.amount;
    obj.relayed_message = message.relayedMessage;
    obj.new_status = message.newStatus;
    obj.status_message = message.statusMessage;
    obj.sender_chain = message.senderChain;
    return obj;
  },
  fromAminoMsg(object: EventInboundFinalizedAminoMsg): EventInboundFinalized {
    return EventInboundFinalized.fromAmino(object.value);
  },
  fromProtoMsg(message: EventInboundFinalizedProtoMsg): EventInboundFinalized {
    return EventInboundFinalized.decode(message.value);
  },
  toProto(message: EventInboundFinalized): Uint8Array {
    return EventInboundFinalized.encode(message).finish();
  },
  toProtoMsg(message: EventInboundFinalized): EventInboundFinalizedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.EventInboundFinalized",
      value: EventInboundFinalized.encode(message).finish()
    };
  }
};
function createBaseEventZrcWithdrawCreated(): EventZrcWithdrawCreated {
  return {
    msgTypeUrl: "",
    cctxIndex: "",
    sender: "",
    senderChain: "",
    inTxHash: "",
    receiver: "",
    receiverChain: "",
    amount: "",
    newStatus: ""
  };
}
export const EventZrcWithdrawCreated = {
  typeUrl: "/zetachain.zetacore.crosschain.EventZrcWithdrawCreated",
  encode(message: EventZrcWithdrawCreated, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.cctxIndex !== "") {
      writer.uint32(18).string(message.cctxIndex);
    }
    if (message.sender !== "") {
      writer.uint32(26).string(message.sender);
    }
    if (message.senderChain !== "") {
      writer.uint32(34).string(message.senderChain);
    }
    if (message.inTxHash !== "") {
      writer.uint32(42).string(message.inTxHash);
    }
    if (message.receiver !== "") {
      writer.uint32(50).string(message.receiver);
    }
    if (message.receiverChain !== "") {
      writer.uint32(58).string(message.receiverChain);
    }
    if (message.amount !== "") {
      writer.uint32(66).string(message.amount);
    }
    if (message.newStatus !== "") {
      writer.uint32(74).string(message.newStatus);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventZrcWithdrawCreated {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventZrcWithdrawCreated();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.cctxIndex = reader.string();
          break;
        case 3:
          message.sender = reader.string();
          break;
        case 4:
          message.senderChain = reader.string();
          break;
        case 5:
          message.inTxHash = reader.string();
          break;
        case 6:
          message.receiver = reader.string();
          break;
        case 7:
          message.receiverChain = reader.string();
          break;
        case 8:
          message.amount = reader.string();
          break;
        case 9:
          message.newStatus = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventZrcWithdrawCreated>): EventZrcWithdrawCreated {
    const message = createBaseEventZrcWithdrawCreated();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.cctxIndex = object.cctxIndex ?? "";
    message.sender = object.sender ?? "";
    message.senderChain = object.senderChain ?? "";
    message.inTxHash = object.inTxHash ?? "";
    message.receiver = object.receiver ?? "";
    message.receiverChain = object.receiverChain ?? "";
    message.amount = object.amount ?? "";
    message.newStatus = object.newStatus ?? "";
    return message;
  },
  fromAmino(object: EventZrcWithdrawCreatedAmino): EventZrcWithdrawCreated {
    const message = createBaseEventZrcWithdrawCreated();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.cctx_index !== undefined && object.cctx_index !== null) {
      message.cctxIndex = object.cctx_index;
    }
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = object.sender;
    }
    if (object.sender_chain !== undefined && object.sender_chain !== null) {
      message.senderChain = object.sender_chain;
    }
    if (object.in_tx_hash !== undefined && object.in_tx_hash !== null) {
      message.inTxHash = object.in_tx_hash;
    }
    if (object.receiver !== undefined && object.receiver !== null) {
      message.receiver = object.receiver;
    }
    if (object.receiver_chain !== undefined && object.receiver_chain !== null) {
      message.receiverChain = object.receiver_chain;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    if (object.new_status !== undefined && object.new_status !== null) {
      message.newStatus = object.new_status;
    }
    return message;
  },
  toAmino(message: EventZrcWithdrawCreated): EventZrcWithdrawCreatedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.cctx_index = message.cctxIndex;
    obj.sender = message.sender;
    obj.sender_chain = message.senderChain;
    obj.in_tx_hash = message.inTxHash;
    obj.receiver = message.receiver;
    obj.receiver_chain = message.receiverChain;
    obj.amount = message.amount;
    obj.new_status = message.newStatus;
    return obj;
  },
  fromAminoMsg(object: EventZrcWithdrawCreatedAminoMsg): EventZrcWithdrawCreated {
    return EventZrcWithdrawCreated.fromAmino(object.value);
  },
  fromProtoMsg(message: EventZrcWithdrawCreatedProtoMsg): EventZrcWithdrawCreated {
    return EventZrcWithdrawCreated.decode(message.value);
  },
  toProto(message: EventZrcWithdrawCreated): Uint8Array {
    return EventZrcWithdrawCreated.encode(message).finish();
  },
  toProtoMsg(message: EventZrcWithdrawCreated): EventZrcWithdrawCreatedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.EventZrcWithdrawCreated",
      value: EventZrcWithdrawCreated.encode(message).finish()
    };
  }
};
function createBaseEventZetaWithdrawCreated(): EventZetaWithdrawCreated {
  return {
    msgTypeUrl: "",
    cctxIndex: "",
    sender: "",
    inTxHash: "",
    newStatus: ""
  };
}
export const EventZetaWithdrawCreated = {
  typeUrl: "/zetachain.zetacore.crosschain.EventZetaWithdrawCreated",
  encode(message: EventZetaWithdrawCreated, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.cctxIndex !== "") {
      writer.uint32(18).string(message.cctxIndex);
    }
    if (message.sender !== "") {
      writer.uint32(26).string(message.sender);
    }
    if (message.inTxHash !== "") {
      writer.uint32(34).string(message.inTxHash);
    }
    if (message.newStatus !== "") {
      writer.uint32(42).string(message.newStatus);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventZetaWithdrawCreated {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventZetaWithdrawCreated();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.cctxIndex = reader.string();
          break;
        case 3:
          message.sender = reader.string();
          break;
        case 4:
          message.inTxHash = reader.string();
          break;
        case 5:
          message.newStatus = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventZetaWithdrawCreated>): EventZetaWithdrawCreated {
    const message = createBaseEventZetaWithdrawCreated();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.cctxIndex = object.cctxIndex ?? "";
    message.sender = object.sender ?? "";
    message.inTxHash = object.inTxHash ?? "";
    message.newStatus = object.newStatus ?? "";
    return message;
  },
  fromAmino(object: EventZetaWithdrawCreatedAmino): EventZetaWithdrawCreated {
    const message = createBaseEventZetaWithdrawCreated();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.cctx_index !== undefined && object.cctx_index !== null) {
      message.cctxIndex = object.cctx_index;
    }
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = object.sender;
    }
    if (object.in_tx_hash !== undefined && object.in_tx_hash !== null) {
      message.inTxHash = object.in_tx_hash;
    }
    if (object.new_status !== undefined && object.new_status !== null) {
      message.newStatus = object.new_status;
    }
    return message;
  },
  toAmino(message: EventZetaWithdrawCreated): EventZetaWithdrawCreatedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.cctx_index = message.cctxIndex;
    obj.sender = message.sender;
    obj.in_tx_hash = message.inTxHash;
    obj.new_status = message.newStatus;
    return obj;
  },
  fromAminoMsg(object: EventZetaWithdrawCreatedAminoMsg): EventZetaWithdrawCreated {
    return EventZetaWithdrawCreated.fromAmino(object.value);
  },
  fromProtoMsg(message: EventZetaWithdrawCreatedProtoMsg): EventZetaWithdrawCreated {
    return EventZetaWithdrawCreated.decode(message.value);
  },
  toProto(message: EventZetaWithdrawCreated): Uint8Array {
    return EventZetaWithdrawCreated.encode(message).finish();
  },
  toProtoMsg(message: EventZetaWithdrawCreated): EventZetaWithdrawCreatedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.EventZetaWithdrawCreated",
      value: EventZetaWithdrawCreated.encode(message).finish()
    };
  }
};
function createBaseEventOutboundFailure(): EventOutboundFailure {
  return {
    msgTypeUrl: "",
    cctxIndex: "",
    oldStatus: "",
    newStatus: "",
    valueReceived: ""
  };
}
export const EventOutboundFailure = {
  typeUrl: "/zetachain.zetacore.crosschain.EventOutboundFailure",
  encode(message: EventOutboundFailure, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.cctxIndex !== "") {
      writer.uint32(18).string(message.cctxIndex);
    }
    if (message.oldStatus !== "") {
      writer.uint32(26).string(message.oldStatus);
    }
    if (message.newStatus !== "") {
      writer.uint32(34).string(message.newStatus);
    }
    if (message.valueReceived !== "") {
      writer.uint32(42).string(message.valueReceived);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventOutboundFailure {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventOutboundFailure();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.cctxIndex = reader.string();
          break;
        case 3:
          message.oldStatus = reader.string();
          break;
        case 4:
          message.newStatus = reader.string();
          break;
        case 5:
          message.valueReceived = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventOutboundFailure>): EventOutboundFailure {
    const message = createBaseEventOutboundFailure();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.cctxIndex = object.cctxIndex ?? "";
    message.oldStatus = object.oldStatus ?? "";
    message.newStatus = object.newStatus ?? "";
    message.valueReceived = object.valueReceived ?? "";
    return message;
  },
  fromAmino(object: EventOutboundFailureAmino): EventOutboundFailure {
    const message = createBaseEventOutboundFailure();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.cctx_index !== undefined && object.cctx_index !== null) {
      message.cctxIndex = object.cctx_index;
    }
    if (object.old_status !== undefined && object.old_status !== null) {
      message.oldStatus = object.old_status;
    }
    if (object.new_status !== undefined && object.new_status !== null) {
      message.newStatus = object.new_status;
    }
    if (object.value_received !== undefined && object.value_received !== null) {
      message.valueReceived = object.value_received;
    }
    return message;
  },
  toAmino(message: EventOutboundFailure): EventOutboundFailureAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.cctx_index = message.cctxIndex;
    obj.old_status = message.oldStatus;
    obj.new_status = message.newStatus;
    obj.value_received = message.valueReceived;
    return obj;
  },
  fromAminoMsg(object: EventOutboundFailureAminoMsg): EventOutboundFailure {
    return EventOutboundFailure.fromAmino(object.value);
  },
  fromProtoMsg(message: EventOutboundFailureProtoMsg): EventOutboundFailure {
    return EventOutboundFailure.decode(message.value);
  },
  toProto(message: EventOutboundFailure): Uint8Array {
    return EventOutboundFailure.encode(message).finish();
  },
  toProtoMsg(message: EventOutboundFailure): EventOutboundFailureProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.EventOutboundFailure",
      value: EventOutboundFailure.encode(message).finish()
    };
  }
};
function createBaseEventOutboundSuccess(): EventOutboundSuccess {
  return {
    msgTypeUrl: "",
    cctxIndex: "",
    oldStatus: "",
    newStatus: "",
    valueReceived: ""
  };
}
export const EventOutboundSuccess = {
  typeUrl: "/zetachain.zetacore.crosschain.EventOutboundSuccess",
  encode(message: EventOutboundSuccess, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.cctxIndex !== "") {
      writer.uint32(18).string(message.cctxIndex);
    }
    if (message.oldStatus !== "") {
      writer.uint32(26).string(message.oldStatus);
    }
    if (message.newStatus !== "") {
      writer.uint32(34).string(message.newStatus);
    }
    if (message.valueReceived !== "") {
      writer.uint32(42).string(message.valueReceived);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventOutboundSuccess {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventOutboundSuccess();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.cctxIndex = reader.string();
          break;
        case 3:
          message.oldStatus = reader.string();
          break;
        case 4:
          message.newStatus = reader.string();
          break;
        case 5:
          message.valueReceived = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventOutboundSuccess>): EventOutboundSuccess {
    const message = createBaseEventOutboundSuccess();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.cctxIndex = object.cctxIndex ?? "";
    message.oldStatus = object.oldStatus ?? "";
    message.newStatus = object.newStatus ?? "";
    message.valueReceived = object.valueReceived ?? "";
    return message;
  },
  fromAmino(object: EventOutboundSuccessAmino): EventOutboundSuccess {
    const message = createBaseEventOutboundSuccess();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.cctx_index !== undefined && object.cctx_index !== null) {
      message.cctxIndex = object.cctx_index;
    }
    if (object.old_status !== undefined && object.old_status !== null) {
      message.oldStatus = object.old_status;
    }
    if (object.new_status !== undefined && object.new_status !== null) {
      message.newStatus = object.new_status;
    }
    if (object.value_received !== undefined && object.value_received !== null) {
      message.valueReceived = object.value_received;
    }
    return message;
  },
  toAmino(message: EventOutboundSuccess): EventOutboundSuccessAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.cctx_index = message.cctxIndex;
    obj.old_status = message.oldStatus;
    obj.new_status = message.newStatus;
    obj.value_received = message.valueReceived;
    return obj;
  },
  fromAminoMsg(object: EventOutboundSuccessAminoMsg): EventOutboundSuccess {
    return EventOutboundSuccess.fromAmino(object.value);
  },
  fromProtoMsg(message: EventOutboundSuccessProtoMsg): EventOutboundSuccess {
    return EventOutboundSuccess.decode(message.value);
  },
  toProto(message: EventOutboundSuccess): Uint8Array {
    return EventOutboundSuccess.encode(message).finish();
  },
  toProtoMsg(message: EventOutboundSuccess): EventOutboundSuccessProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.EventOutboundSuccess",
      value: EventOutboundSuccess.encode(message).finish()
    };
  }
};
function createBaseEventCCTXGasPriceIncreased(): EventCCTXGasPriceIncreased {
  return {
    cctxIndex: "",
    gasPriceIncrease: "",
    additionalFees: ""
  };
}
export const EventCCTXGasPriceIncreased = {
  typeUrl: "/zetachain.zetacore.crosschain.EventCCTXGasPriceIncreased",
  encode(message: EventCCTXGasPriceIncreased, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.cctxIndex !== "") {
      writer.uint32(10).string(message.cctxIndex);
    }
    if (message.gasPriceIncrease !== "") {
      writer.uint32(18).string(message.gasPriceIncrease);
    }
    if (message.additionalFees !== "") {
      writer.uint32(26).string(message.additionalFees);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventCCTXGasPriceIncreased {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventCCTXGasPriceIncreased();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.cctxIndex = reader.string();
          break;
        case 2:
          message.gasPriceIncrease = reader.string();
          break;
        case 3:
          message.additionalFees = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventCCTXGasPriceIncreased>): EventCCTXGasPriceIncreased {
    const message = createBaseEventCCTXGasPriceIncreased();
    message.cctxIndex = object.cctxIndex ?? "";
    message.gasPriceIncrease = object.gasPriceIncrease ?? "";
    message.additionalFees = object.additionalFees ?? "";
    return message;
  },
  fromAmino(object: EventCCTXGasPriceIncreasedAmino): EventCCTXGasPriceIncreased {
    const message = createBaseEventCCTXGasPriceIncreased();
    if (object.cctx_index !== undefined && object.cctx_index !== null) {
      message.cctxIndex = object.cctx_index;
    }
    if (object.gas_price_increase !== undefined && object.gas_price_increase !== null) {
      message.gasPriceIncrease = object.gas_price_increase;
    }
    if (object.additional_fees !== undefined && object.additional_fees !== null) {
      message.additionalFees = object.additional_fees;
    }
    return message;
  },
  toAmino(message: EventCCTXGasPriceIncreased): EventCCTXGasPriceIncreasedAmino {
    const obj: any = {};
    obj.cctx_index = message.cctxIndex;
    obj.gas_price_increase = message.gasPriceIncrease;
    obj.additional_fees = message.additionalFees;
    return obj;
  },
  fromAminoMsg(object: EventCCTXGasPriceIncreasedAminoMsg): EventCCTXGasPriceIncreased {
    return EventCCTXGasPriceIncreased.fromAmino(object.value);
  },
  fromProtoMsg(message: EventCCTXGasPriceIncreasedProtoMsg): EventCCTXGasPriceIncreased {
    return EventCCTXGasPriceIncreased.decode(message.value);
  },
  toProto(message: EventCCTXGasPriceIncreased): Uint8Array {
    return EventCCTXGasPriceIncreased.encode(message).finish();
  },
  toProtoMsg(message: EventCCTXGasPriceIncreased): EventCCTXGasPriceIncreasedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.EventCCTXGasPriceIncreased",
      value: EventCCTXGasPriceIncreased.encode(message).finish()
    };
  }
};