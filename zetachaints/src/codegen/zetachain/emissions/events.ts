import { BinaryReader, BinaryWriter } from "../../binary";
export enum EmissionType {
  Slash = 0,
  Rewards = 1,
  UNRECOGNIZED = -1,
}
export const EmissionTypeSDKType = EmissionType;
export const EmissionTypeAmino = EmissionType;
export function emissionTypeFromJSON(object: any): EmissionType {
  switch (object) {
    case 0:
    case "Slash":
      return EmissionType.Slash;
    case 1:
    case "Rewards":
      return EmissionType.Rewards;
    case -1:
    case "UNRECOGNIZED":
    default:
      return EmissionType.UNRECOGNIZED;
  }
}
export function emissionTypeToJSON(object: EmissionType): string {
  switch (object) {
    case EmissionType.Slash:
      return "Slash";
    case EmissionType.Rewards:
      return "Rewards";
    case EmissionType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export interface ObserverEmission {
  emissionType: EmissionType;
  observerAddress: string;
  amount: string;
}
export interface ObserverEmissionProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.ObserverEmission";
  value: Uint8Array;
}
export interface ObserverEmissionAmino {
  emission_type?: EmissionType;
  observer_address?: string;
  amount?: string;
}
export interface ObserverEmissionAminoMsg {
  type: "/zetachain.zetacore.emissions.ObserverEmission";
  value: ObserverEmissionAmino;
}
export interface ObserverEmissionSDKType {
  emission_type: EmissionType;
  observer_address: string;
  amount: string;
}
export interface EventObserverEmissions {
  msgTypeUrl: string;
  emissions: ObserverEmission[];
}
export interface EventObserverEmissionsProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.EventObserverEmissions";
  value: Uint8Array;
}
export interface EventObserverEmissionsAmino {
  msg_type_url?: string;
  emissions?: ObserverEmissionAmino[];
}
export interface EventObserverEmissionsAminoMsg {
  type: "/zetachain.zetacore.emissions.EventObserverEmissions";
  value: EventObserverEmissionsAmino;
}
export interface EventObserverEmissionsSDKType {
  msg_type_url: string;
  emissions: ObserverEmissionSDKType[];
}
export interface EventBlockEmissions {
  msgTypeUrl: string;
  bondFactor: string;
  reservesFactor: string;
  durationFactor: string;
  validatorRewardsForBlock: string;
  observerRewardsForBlock: string;
  tssRewardsForBlock: string;
}
export interface EventBlockEmissionsProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.EventBlockEmissions";
  value: Uint8Array;
}
export interface EventBlockEmissionsAmino {
  msg_type_url?: string;
  bond_factor?: string;
  reserves_factor?: string;
  duration_factor?: string;
  validator_rewards_for_block?: string;
  observer_rewards_for_block?: string;
  tss_rewards_for_block?: string;
}
export interface EventBlockEmissionsAminoMsg {
  type: "/zetachain.zetacore.emissions.EventBlockEmissions";
  value: EventBlockEmissionsAmino;
}
export interface EventBlockEmissionsSDKType {
  msg_type_url: string;
  bond_factor: string;
  reserves_factor: string;
  duration_factor: string;
  validator_rewards_for_block: string;
  observer_rewards_for_block: string;
  tss_rewards_for_block: string;
}
function createBaseObserverEmission(): ObserverEmission {
  return {
    emissionType: 0,
    observerAddress: "",
    amount: ""
  };
}
export const ObserverEmission = {
  typeUrl: "/zetachain.zetacore.emissions.ObserverEmission",
  encode(message: ObserverEmission, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.emissionType !== 0) {
      writer.uint32(8).int32(message.emissionType);
    }
    if (message.observerAddress !== "") {
      writer.uint32(18).string(message.observerAddress);
    }
    if (message.amount !== "") {
      writer.uint32(26).string(message.amount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ObserverEmission {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObserverEmission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.emissionType = (reader.int32() as any);
          break;
        case 2:
          message.observerAddress = reader.string();
          break;
        case 3:
          message.amount = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ObserverEmission>): ObserverEmission {
    const message = createBaseObserverEmission();
    message.emissionType = object.emissionType ?? 0;
    message.observerAddress = object.observerAddress ?? "";
    message.amount = object.amount ?? "";
    return message;
  },
  fromAmino(object: ObserverEmissionAmino): ObserverEmission {
    const message = createBaseObserverEmission();
    if (object.emission_type !== undefined && object.emission_type !== null) {
      message.emissionType = emissionTypeFromJSON(object.emission_type);
    }
    if (object.observer_address !== undefined && object.observer_address !== null) {
      message.observerAddress = object.observer_address;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    return message;
  },
  toAmino(message: ObserverEmission): ObserverEmissionAmino {
    const obj: any = {};
    obj.emission_type = message.emissionType;
    obj.observer_address = message.observerAddress;
    obj.amount = message.amount;
    return obj;
  },
  fromAminoMsg(object: ObserverEmissionAminoMsg): ObserverEmission {
    return ObserverEmission.fromAmino(object.value);
  },
  fromProtoMsg(message: ObserverEmissionProtoMsg): ObserverEmission {
    return ObserverEmission.decode(message.value);
  },
  toProto(message: ObserverEmission): Uint8Array {
    return ObserverEmission.encode(message).finish();
  },
  toProtoMsg(message: ObserverEmission): ObserverEmissionProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.ObserverEmission",
      value: ObserverEmission.encode(message).finish()
    };
  }
};
function createBaseEventObserverEmissions(): EventObserverEmissions {
  return {
    msgTypeUrl: "",
    emissions: []
  };
}
export const EventObserverEmissions = {
  typeUrl: "/zetachain.zetacore.emissions.EventObserverEmissions",
  encode(message: EventObserverEmissions, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    for (const v of message.emissions) {
      ObserverEmission.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventObserverEmissions {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventObserverEmissions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.emissions.push(ObserverEmission.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventObserverEmissions>): EventObserverEmissions {
    const message = createBaseEventObserverEmissions();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.emissions = object.emissions?.map(e => ObserverEmission.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: EventObserverEmissionsAmino): EventObserverEmissions {
    const message = createBaseEventObserverEmissions();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    message.emissions = object.emissions?.map(e => ObserverEmission.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: EventObserverEmissions): EventObserverEmissionsAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    if (message.emissions) {
      obj.emissions = message.emissions.map(e => e ? ObserverEmission.toAmino(e) : undefined);
    } else {
      obj.emissions = [];
    }
    return obj;
  },
  fromAminoMsg(object: EventObserverEmissionsAminoMsg): EventObserverEmissions {
    return EventObserverEmissions.fromAmino(object.value);
  },
  fromProtoMsg(message: EventObserverEmissionsProtoMsg): EventObserverEmissions {
    return EventObserverEmissions.decode(message.value);
  },
  toProto(message: EventObserverEmissions): Uint8Array {
    return EventObserverEmissions.encode(message).finish();
  },
  toProtoMsg(message: EventObserverEmissions): EventObserverEmissionsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.EventObserverEmissions",
      value: EventObserverEmissions.encode(message).finish()
    };
  }
};
function createBaseEventBlockEmissions(): EventBlockEmissions {
  return {
    msgTypeUrl: "",
    bondFactor: "",
    reservesFactor: "",
    durationFactor: "",
    validatorRewardsForBlock: "",
    observerRewardsForBlock: "",
    tssRewardsForBlock: ""
  };
}
export const EventBlockEmissions = {
  typeUrl: "/zetachain.zetacore.emissions.EventBlockEmissions",
  encode(message: EventBlockEmissions, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.bondFactor !== "") {
      writer.uint32(18).string(message.bondFactor);
    }
    if (message.reservesFactor !== "") {
      writer.uint32(26).string(message.reservesFactor);
    }
    if (message.durationFactor !== "") {
      writer.uint32(34).string(message.durationFactor);
    }
    if (message.validatorRewardsForBlock !== "") {
      writer.uint32(42).string(message.validatorRewardsForBlock);
    }
    if (message.observerRewardsForBlock !== "") {
      writer.uint32(50).string(message.observerRewardsForBlock);
    }
    if (message.tssRewardsForBlock !== "") {
      writer.uint32(58).string(message.tssRewardsForBlock);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventBlockEmissions {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventBlockEmissions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.bondFactor = reader.string();
          break;
        case 3:
          message.reservesFactor = reader.string();
          break;
        case 4:
          message.durationFactor = reader.string();
          break;
        case 5:
          message.validatorRewardsForBlock = reader.string();
          break;
        case 6:
          message.observerRewardsForBlock = reader.string();
          break;
        case 7:
          message.tssRewardsForBlock = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventBlockEmissions>): EventBlockEmissions {
    const message = createBaseEventBlockEmissions();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.bondFactor = object.bondFactor ?? "";
    message.reservesFactor = object.reservesFactor ?? "";
    message.durationFactor = object.durationFactor ?? "";
    message.validatorRewardsForBlock = object.validatorRewardsForBlock ?? "";
    message.observerRewardsForBlock = object.observerRewardsForBlock ?? "";
    message.tssRewardsForBlock = object.tssRewardsForBlock ?? "";
    return message;
  },
  fromAmino(object: EventBlockEmissionsAmino): EventBlockEmissions {
    const message = createBaseEventBlockEmissions();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.bond_factor !== undefined && object.bond_factor !== null) {
      message.bondFactor = object.bond_factor;
    }
    if (object.reserves_factor !== undefined && object.reserves_factor !== null) {
      message.reservesFactor = object.reserves_factor;
    }
    if (object.duration_factor !== undefined && object.duration_factor !== null) {
      message.durationFactor = object.duration_factor;
    }
    if (object.validator_rewards_for_block !== undefined && object.validator_rewards_for_block !== null) {
      message.validatorRewardsForBlock = object.validator_rewards_for_block;
    }
    if (object.observer_rewards_for_block !== undefined && object.observer_rewards_for_block !== null) {
      message.observerRewardsForBlock = object.observer_rewards_for_block;
    }
    if (object.tss_rewards_for_block !== undefined && object.tss_rewards_for_block !== null) {
      message.tssRewardsForBlock = object.tss_rewards_for_block;
    }
    return message;
  },
  toAmino(message: EventBlockEmissions): EventBlockEmissionsAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.bond_factor = message.bondFactor;
    obj.reserves_factor = message.reservesFactor;
    obj.duration_factor = message.durationFactor;
    obj.validator_rewards_for_block = message.validatorRewardsForBlock;
    obj.observer_rewards_for_block = message.observerRewardsForBlock;
    obj.tss_rewards_for_block = message.tssRewardsForBlock;
    return obj;
  },
  fromAminoMsg(object: EventBlockEmissionsAminoMsg): EventBlockEmissions {
    return EventBlockEmissions.fromAmino(object.value);
  },
  fromProtoMsg(message: EventBlockEmissionsProtoMsg): EventBlockEmissions {
    return EventBlockEmissions.decode(message.value);
  },
  toProto(message: EventBlockEmissions): Uint8Array {
    return EventBlockEmissions.encode(message).finish();
  },
  toProtoMsg(message: EventBlockEmissions): EventBlockEmissionsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.EventBlockEmissions",
      value: EventBlockEmissions.encode(message).finish()
    };
  }
};