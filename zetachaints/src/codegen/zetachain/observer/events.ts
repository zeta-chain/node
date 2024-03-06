import { GasPriceIncreaseFlags, GasPriceIncreaseFlagsAmino, GasPriceIncreaseFlagsSDKType, BlockHeaderVerificationFlags, BlockHeaderVerificationFlagsAmino, BlockHeaderVerificationFlagsSDKType } from "./crosschain_flags";
import { BinaryReader, BinaryWriter } from "../../binary";
export interface EventBallotCreated {
  msgTypeUrl: string;
  ballotIdentifier: string;
  observationHash: string;
  observationChain: string;
  ballotType: string;
}
export interface EventBallotCreatedProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.EventBallotCreated";
  value: Uint8Array;
}
export interface EventBallotCreatedAmino {
  msg_type_url?: string;
  ballot_identifier?: string;
  observation_hash?: string;
  observation_chain?: string;
  ballot_type?: string;
}
export interface EventBallotCreatedAminoMsg {
  type: "/zetachain.zetacore.observer.EventBallotCreated";
  value: EventBallotCreatedAmino;
}
export interface EventBallotCreatedSDKType {
  msg_type_url: string;
  ballot_identifier: string;
  observation_hash: string;
  observation_chain: string;
  ballot_type: string;
}
export interface EventKeygenBlockUpdated {
  msgTypeUrl: string;
  keygenBlock: string;
  keygenPubkeys: string;
}
export interface EventKeygenBlockUpdatedProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.EventKeygenBlockUpdated";
  value: Uint8Array;
}
export interface EventKeygenBlockUpdatedAmino {
  msg_type_url?: string;
  keygen_block?: string;
  keygen_pubkeys?: string;
}
export interface EventKeygenBlockUpdatedAminoMsg {
  type: "/zetachain.zetacore.observer.EventKeygenBlockUpdated";
  value: EventKeygenBlockUpdatedAmino;
}
export interface EventKeygenBlockUpdatedSDKType {
  msg_type_url: string;
  keygen_block: string;
  keygen_pubkeys: string;
}
export interface EventNewObserverAdded {
  msgTypeUrl: string;
  observerAddress: string;
  zetaclientGranteeAddress: string;
  zetaclientGranteePubkey: string;
  observerLastBlockCount: bigint;
}
export interface EventNewObserverAddedProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.EventNewObserverAdded";
  value: Uint8Array;
}
export interface EventNewObserverAddedAmino {
  msg_type_url?: string;
  observer_address?: string;
  zetaclient_grantee_address?: string;
  zetaclient_grantee_pubkey?: string;
  observer_last_block_count?: string;
}
export interface EventNewObserverAddedAminoMsg {
  type: "/zetachain.zetacore.observer.EventNewObserverAdded";
  value: EventNewObserverAddedAmino;
}
export interface EventNewObserverAddedSDKType {
  msg_type_url: string;
  observer_address: string;
  zetaclient_grantee_address: string;
  zetaclient_grantee_pubkey: string;
  observer_last_block_count: bigint;
}
export interface EventCrosschainFlagsUpdated {
  msgTypeUrl: string;
  isInboundEnabled: boolean;
  isOutboundEnabled: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlags;
  signer: string;
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlags;
}
export interface EventCrosschainFlagsUpdatedProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.EventCrosschainFlagsUpdated";
  value: Uint8Array;
}
export interface EventCrosschainFlagsUpdatedAmino {
  msg_type_url?: string;
  isInboundEnabled?: boolean;
  isOutboundEnabled?: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlagsAmino;
  signer?: string;
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlagsAmino;
}
export interface EventCrosschainFlagsUpdatedAminoMsg {
  type: "/zetachain.zetacore.observer.EventCrosschainFlagsUpdated";
  value: EventCrosschainFlagsUpdatedAmino;
}
export interface EventCrosschainFlagsUpdatedSDKType {
  msg_type_url: string;
  isInboundEnabled: boolean;
  isOutboundEnabled: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlagsSDKType;
  signer: string;
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlagsSDKType;
}
function createBaseEventBallotCreated(): EventBallotCreated {
  return {
    msgTypeUrl: "",
    ballotIdentifier: "",
    observationHash: "",
    observationChain: "",
    ballotType: ""
  };
}
export const EventBallotCreated = {
  typeUrl: "/zetachain.zetacore.observer.EventBallotCreated",
  encode(message: EventBallotCreated, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.ballotIdentifier !== "") {
      writer.uint32(18).string(message.ballotIdentifier);
    }
    if (message.observationHash !== "") {
      writer.uint32(26).string(message.observationHash);
    }
    if (message.observationChain !== "") {
      writer.uint32(34).string(message.observationChain);
    }
    if (message.ballotType !== "") {
      writer.uint32(42).string(message.ballotType);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventBallotCreated {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventBallotCreated();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.ballotIdentifier = reader.string();
          break;
        case 3:
          message.observationHash = reader.string();
          break;
        case 4:
          message.observationChain = reader.string();
          break;
        case 5:
          message.ballotType = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventBallotCreated>): EventBallotCreated {
    const message = createBaseEventBallotCreated();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.ballotIdentifier = object.ballotIdentifier ?? "";
    message.observationHash = object.observationHash ?? "";
    message.observationChain = object.observationChain ?? "";
    message.ballotType = object.ballotType ?? "";
    return message;
  },
  fromAmino(object: EventBallotCreatedAmino): EventBallotCreated {
    const message = createBaseEventBallotCreated();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.ballot_identifier !== undefined && object.ballot_identifier !== null) {
      message.ballotIdentifier = object.ballot_identifier;
    }
    if (object.observation_hash !== undefined && object.observation_hash !== null) {
      message.observationHash = object.observation_hash;
    }
    if (object.observation_chain !== undefined && object.observation_chain !== null) {
      message.observationChain = object.observation_chain;
    }
    if (object.ballot_type !== undefined && object.ballot_type !== null) {
      message.ballotType = object.ballot_type;
    }
    return message;
  },
  toAmino(message: EventBallotCreated): EventBallotCreatedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.ballot_identifier = message.ballotIdentifier;
    obj.observation_hash = message.observationHash;
    obj.observation_chain = message.observationChain;
    obj.ballot_type = message.ballotType;
    return obj;
  },
  fromAminoMsg(object: EventBallotCreatedAminoMsg): EventBallotCreated {
    return EventBallotCreated.fromAmino(object.value);
  },
  fromProtoMsg(message: EventBallotCreatedProtoMsg): EventBallotCreated {
    return EventBallotCreated.decode(message.value);
  },
  toProto(message: EventBallotCreated): Uint8Array {
    return EventBallotCreated.encode(message).finish();
  },
  toProtoMsg(message: EventBallotCreated): EventBallotCreatedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.EventBallotCreated",
      value: EventBallotCreated.encode(message).finish()
    };
  }
};
function createBaseEventKeygenBlockUpdated(): EventKeygenBlockUpdated {
  return {
    msgTypeUrl: "",
    keygenBlock: "",
    keygenPubkeys: ""
  };
}
export const EventKeygenBlockUpdated = {
  typeUrl: "/zetachain.zetacore.observer.EventKeygenBlockUpdated",
  encode(message: EventKeygenBlockUpdated, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.keygenBlock !== "") {
      writer.uint32(18).string(message.keygenBlock);
    }
    if (message.keygenPubkeys !== "") {
      writer.uint32(26).string(message.keygenPubkeys);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventKeygenBlockUpdated {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventKeygenBlockUpdated();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.keygenBlock = reader.string();
          break;
        case 3:
          message.keygenPubkeys = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventKeygenBlockUpdated>): EventKeygenBlockUpdated {
    const message = createBaseEventKeygenBlockUpdated();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.keygenBlock = object.keygenBlock ?? "";
    message.keygenPubkeys = object.keygenPubkeys ?? "";
    return message;
  },
  fromAmino(object: EventKeygenBlockUpdatedAmino): EventKeygenBlockUpdated {
    const message = createBaseEventKeygenBlockUpdated();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.keygen_block !== undefined && object.keygen_block !== null) {
      message.keygenBlock = object.keygen_block;
    }
    if (object.keygen_pubkeys !== undefined && object.keygen_pubkeys !== null) {
      message.keygenPubkeys = object.keygen_pubkeys;
    }
    return message;
  },
  toAmino(message: EventKeygenBlockUpdated): EventKeygenBlockUpdatedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.keygen_block = message.keygenBlock;
    obj.keygen_pubkeys = message.keygenPubkeys;
    return obj;
  },
  fromAminoMsg(object: EventKeygenBlockUpdatedAminoMsg): EventKeygenBlockUpdated {
    return EventKeygenBlockUpdated.fromAmino(object.value);
  },
  fromProtoMsg(message: EventKeygenBlockUpdatedProtoMsg): EventKeygenBlockUpdated {
    return EventKeygenBlockUpdated.decode(message.value);
  },
  toProto(message: EventKeygenBlockUpdated): Uint8Array {
    return EventKeygenBlockUpdated.encode(message).finish();
  },
  toProtoMsg(message: EventKeygenBlockUpdated): EventKeygenBlockUpdatedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.EventKeygenBlockUpdated",
      value: EventKeygenBlockUpdated.encode(message).finish()
    };
  }
};
function createBaseEventNewObserverAdded(): EventNewObserverAdded {
  return {
    msgTypeUrl: "",
    observerAddress: "",
    zetaclientGranteeAddress: "",
    zetaclientGranteePubkey: "",
    observerLastBlockCount: BigInt(0)
  };
}
export const EventNewObserverAdded = {
  typeUrl: "/zetachain.zetacore.observer.EventNewObserverAdded",
  encode(message: EventNewObserverAdded, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.observerAddress !== "") {
      writer.uint32(18).string(message.observerAddress);
    }
    if (message.zetaclientGranteeAddress !== "") {
      writer.uint32(26).string(message.zetaclientGranteeAddress);
    }
    if (message.zetaclientGranteePubkey !== "") {
      writer.uint32(34).string(message.zetaclientGranteePubkey);
    }
    if (message.observerLastBlockCount !== BigInt(0)) {
      writer.uint32(40).uint64(message.observerLastBlockCount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventNewObserverAdded {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventNewObserverAdded();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.observerAddress = reader.string();
          break;
        case 3:
          message.zetaclientGranteeAddress = reader.string();
          break;
        case 4:
          message.zetaclientGranteePubkey = reader.string();
          break;
        case 5:
          message.observerLastBlockCount = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventNewObserverAdded>): EventNewObserverAdded {
    const message = createBaseEventNewObserverAdded();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.observerAddress = object.observerAddress ?? "";
    message.zetaclientGranteeAddress = object.zetaclientGranteeAddress ?? "";
    message.zetaclientGranteePubkey = object.zetaclientGranteePubkey ?? "";
    message.observerLastBlockCount = object.observerLastBlockCount !== undefined && object.observerLastBlockCount !== null ? BigInt(object.observerLastBlockCount.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: EventNewObserverAddedAmino): EventNewObserverAdded {
    const message = createBaseEventNewObserverAdded();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.observer_address !== undefined && object.observer_address !== null) {
      message.observerAddress = object.observer_address;
    }
    if (object.zetaclient_grantee_address !== undefined && object.zetaclient_grantee_address !== null) {
      message.zetaclientGranteeAddress = object.zetaclient_grantee_address;
    }
    if (object.zetaclient_grantee_pubkey !== undefined && object.zetaclient_grantee_pubkey !== null) {
      message.zetaclientGranteePubkey = object.zetaclient_grantee_pubkey;
    }
    if (object.observer_last_block_count !== undefined && object.observer_last_block_count !== null) {
      message.observerLastBlockCount = BigInt(object.observer_last_block_count);
    }
    return message;
  },
  toAmino(message: EventNewObserverAdded): EventNewObserverAddedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.observer_address = message.observerAddress;
    obj.zetaclient_grantee_address = message.zetaclientGranteeAddress;
    obj.zetaclient_grantee_pubkey = message.zetaclientGranteePubkey;
    obj.observer_last_block_count = message.observerLastBlockCount ? message.observerLastBlockCount.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: EventNewObserverAddedAminoMsg): EventNewObserverAdded {
    return EventNewObserverAdded.fromAmino(object.value);
  },
  fromProtoMsg(message: EventNewObserverAddedProtoMsg): EventNewObserverAdded {
    return EventNewObserverAdded.decode(message.value);
  },
  toProto(message: EventNewObserverAdded): Uint8Array {
    return EventNewObserverAdded.encode(message).finish();
  },
  toProtoMsg(message: EventNewObserverAdded): EventNewObserverAddedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.EventNewObserverAdded",
      value: EventNewObserverAdded.encode(message).finish()
    };
  }
};
function createBaseEventCrosschainFlagsUpdated(): EventCrosschainFlagsUpdated {
  return {
    msgTypeUrl: "",
    isInboundEnabled: false,
    isOutboundEnabled: false,
    gasPriceIncreaseFlags: undefined,
    signer: "",
    blockHeaderVerificationFlags: undefined
  };
}
export const EventCrosschainFlagsUpdated = {
  typeUrl: "/zetachain.zetacore.observer.EventCrosschainFlagsUpdated",
  encode(message: EventCrosschainFlagsUpdated, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.isInboundEnabled === true) {
      writer.uint32(16).bool(message.isInboundEnabled);
    }
    if (message.isOutboundEnabled === true) {
      writer.uint32(24).bool(message.isOutboundEnabled);
    }
    if (message.gasPriceIncreaseFlags !== undefined) {
      GasPriceIncreaseFlags.encode(message.gasPriceIncreaseFlags, writer.uint32(34).fork()).ldelim();
    }
    if (message.signer !== "") {
      writer.uint32(42).string(message.signer);
    }
    if (message.blockHeaderVerificationFlags !== undefined) {
      BlockHeaderVerificationFlags.encode(message.blockHeaderVerificationFlags, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventCrosschainFlagsUpdated {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventCrosschainFlagsUpdated();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.isInboundEnabled = reader.bool();
          break;
        case 3:
          message.isOutboundEnabled = reader.bool();
          break;
        case 4:
          message.gasPriceIncreaseFlags = GasPriceIncreaseFlags.decode(reader, reader.uint32());
          break;
        case 5:
          message.signer = reader.string();
          break;
        case 6:
          message.blockHeaderVerificationFlags = BlockHeaderVerificationFlags.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventCrosschainFlagsUpdated>): EventCrosschainFlagsUpdated {
    const message = createBaseEventCrosschainFlagsUpdated();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.isInboundEnabled = object.isInboundEnabled ?? false;
    message.isOutboundEnabled = object.isOutboundEnabled ?? false;
    message.gasPriceIncreaseFlags = object.gasPriceIncreaseFlags !== undefined && object.gasPriceIncreaseFlags !== null ? GasPriceIncreaseFlags.fromPartial(object.gasPriceIncreaseFlags) : undefined;
    message.signer = object.signer ?? "";
    message.blockHeaderVerificationFlags = object.blockHeaderVerificationFlags !== undefined && object.blockHeaderVerificationFlags !== null ? BlockHeaderVerificationFlags.fromPartial(object.blockHeaderVerificationFlags) : undefined;
    return message;
  },
  fromAmino(object: EventCrosschainFlagsUpdatedAmino): EventCrosschainFlagsUpdated {
    const message = createBaseEventCrosschainFlagsUpdated();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.isInboundEnabled !== undefined && object.isInboundEnabled !== null) {
      message.isInboundEnabled = object.isInboundEnabled;
    }
    if (object.isOutboundEnabled !== undefined && object.isOutboundEnabled !== null) {
      message.isOutboundEnabled = object.isOutboundEnabled;
    }
    if (object.gasPriceIncreaseFlags !== undefined && object.gasPriceIncreaseFlags !== null) {
      message.gasPriceIncreaseFlags = GasPriceIncreaseFlags.fromAmino(object.gasPriceIncreaseFlags);
    }
    if (object.signer !== undefined && object.signer !== null) {
      message.signer = object.signer;
    }
    if (object.blockHeaderVerificationFlags !== undefined && object.blockHeaderVerificationFlags !== null) {
      message.blockHeaderVerificationFlags = BlockHeaderVerificationFlags.fromAmino(object.blockHeaderVerificationFlags);
    }
    return message;
  },
  toAmino(message: EventCrosschainFlagsUpdated): EventCrosschainFlagsUpdatedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.isInboundEnabled = message.isInboundEnabled;
    obj.isOutboundEnabled = message.isOutboundEnabled;
    obj.gasPriceIncreaseFlags = message.gasPriceIncreaseFlags ? GasPriceIncreaseFlags.toAmino(message.gasPriceIncreaseFlags) : undefined;
    obj.signer = message.signer;
    obj.blockHeaderVerificationFlags = message.blockHeaderVerificationFlags ? BlockHeaderVerificationFlags.toAmino(message.blockHeaderVerificationFlags) : undefined;
    return obj;
  },
  fromAminoMsg(object: EventCrosschainFlagsUpdatedAminoMsg): EventCrosschainFlagsUpdated {
    return EventCrosschainFlagsUpdated.fromAmino(object.value);
  },
  fromProtoMsg(message: EventCrosschainFlagsUpdatedProtoMsg): EventCrosschainFlagsUpdated {
    return EventCrosschainFlagsUpdated.decode(message.value);
  },
  toProto(message: EventCrosschainFlagsUpdated): Uint8Array {
    return EventCrosschainFlagsUpdated.encode(message).finish();
  },
  toProtoMsg(message: EventCrosschainFlagsUpdated): EventCrosschainFlagsUpdatedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.EventCrosschainFlagsUpdated",
      value: EventCrosschainFlagsUpdated.encode(message).finish()
    };
  }
};