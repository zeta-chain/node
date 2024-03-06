import { BinaryReader, BinaryWriter } from "../../binary";
/** Params defines the parameters for the module. */
export interface Params {
  maxBondFactor: string;
  minBondFactor: string;
  avgBlockTime: string;
  targetBondRatio: string;
  validatorEmissionPercentage: string;
  observerEmissionPercentage: string;
  tssSignerEmissionPercentage: string;
  durationFactorConstant: string;
  observerSlashAmount: string;
}
export interface ParamsProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.Params";
  value: Uint8Array;
}
/** Params defines the parameters for the module. */
export interface ParamsAmino {
  max_bond_factor?: string;
  min_bond_factor?: string;
  avg_block_time?: string;
  target_bond_ratio?: string;
  validator_emission_percentage?: string;
  observer_emission_percentage?: string;
  tss_signer_emission_percentage?: string;
  duration_factor_constant?: string;
  observer_slash_amount?: string;
}
export interface ParamsAminoMsg {
  type: "/zetachain.zetacore.emissions.Params";
  value: ParamsAmino;
}
/** Params defines the parameters for the module. */
export interface ParamsSDKType {
  max_bond_factor: string;
  min_bond_factor: string;
  avg_block_time: string;
  target_bond_ratio: string;
  validator_emission_percentage: string;
  observer_emission_percentage: string;
  tss_signer_emission_percentage: string;
  duration_factor_constant: string;
  observer_slash_amount: string;
}
function createBaseParams(): Params {
  return {
    maxBondFactor: "",
    minBondFactor: "",
    avgBlockTime: "",
    targetBondRatio: "",
    validatorEmissionPercentage: "",
    observerEmissionPercentage: "",
    tssSignerEmissionPercentage: "",
    durationFactorConstant: "",
    observerSlashAmount: ""
  };
}
export const Params = {
  typeUrl: "/zetachain.zetacore.emissions.Params",
  encode(message: Params, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.maxBondFactor !== "") {
      writer.uint32(10).string(message.maxBondFactor);
    }
    if (message.minBondFactor !== "") {
      writer.uint32(18).string(message.minBondFactor);
    }
    if (message.avgBlockTime !== "") {
      writer.uint32(26).string(message.avgBlockTime);
    }
    if (message.targetBondRatio !== "") {
      writer.uint32(34).string(message.targetBondRatio);
    }
    if (message.validatorEmissionPercentage !== "") {
      writer.uint32(42).string(message.validatorEmissionPercentage);
    }
    if (message.observerEmissionPercentage !== "") {
      writer.uint32(50).string(message.observerEmissionPercentage);
    }
    if (message.tssSignerEmissionPercentage !== "") {
      writer.uint32(58).string(message.tssSignerEmissionPercentage);
    }
    if (message.durationFactorConstant !== "") {
      writer.uint32(66).string(message.durationFactorConstant);
    }
    if (message.observerSlashAmount !== "") {
      writer.uint32(74).string(message.observerSlashAmount);
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
          message.maxBondFactor = reader.string();
          break;
        case 2:
          message.minBondFactor = reader.string();
          break;
        case 3:
          message.avgBlockTime = reader.string();
          break;
        case 4:
          message.targetBondRatio = reader.string();
          break;
        case 5:
          message.validatorEmissionPercentage = reader.string();
          break;
        case 6:
          message.observerEmissionPercentage = reader.string();
          break;
        case 7:
          message.tssSignerEmissionPercentage = reader.string();
          break;
        case 8:
          message.durationFactorConstant = reader.string();
          break;
        case 9:
          message.observerSlashAmount = reader.string();
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
    message.maxBondFactor = object.maxBondFactor ?? "";
    message.minBondFactor = object.minBondFactor ?? "";
    message.avgBlockTime = object.avgBlockTime ?? "";
    message.targetBondRatio = object.targetBondRatio ?? "";
    message.validatorEmissionPercentage = object.validatorEmissionPercentage ?? "";
    message.observerEmissionPercentage = object.observerEmissionPercentage ?? "";
    message.tssSignerEmissionPercentage = object.tssSignerEmissionPercentage ?? "";
    message.durationFactorConstant = object.durationFactorConstant ?? "";
    message.observerSlashAmount = object.observerSlashAmount ?? "";
    return message;
  },
  fromAmino(object: ParamsAmino): Params {
    const message = createBaseParams();
    if (object.max_bond_factor !== undefined && object.max_bond_factor !== null) {
      message.maxBondFactor = object.max_bond_factor;
    }
    if (object.min_bond_factor !== undefined && object.min_bond_factor !== null) {
      message.minBondFactor = object.min_bond_factor;
    }
    if (object.avg_block_time !== undefined && object.avg_block_time !== null) {
      message.avgBlockTime = object.avg_block_time;
    }
    if (object.target_bond_ratio !== undefined && object.target_bond_ratio !== null) {
      message.targetBondRatio = object.target_bond_ratio;
    }
    if (object.validator_emission_percentage !== undefined && object.validator_emission_percentage !== null) {
      message.validatorEmissionPercentage = object.validator_emission_percentage;
    }
    if (object.observer_emission_percentage !== undefined && object.observer_emission_percentage !== null) {
      message.observerEmissionPercentage = object.observer_emission_percentage;
    }
    if (object.tss_signer_emission_percentage !== undefined && object.tss_signer_emission_percentage !== null) {
      message.tssSignerEmissionPercentage = object.tss_signer_emission_percentage;
    }
    if (object.duration_factor_constant !== undefined && object.duration_factor_constant !== null) {
      message.durationFactorConstant = object.duration_factor_constant;
    }
    if (object.observer_slash_amount !== undefined && object.observer_slash_amount !== null) {
      message.observerSlashAmount = object.observer_slash_amount;
    }
    return message;
  },
  toAmino(message: Params): ParamsAmino {
    const obj: any = {};
    obj.max_bond_factor = message.maxBondFactor;
    obj.min_bond_factor = message.minBondFactor;
    obj.avg_block_time = message.avgBlockTime;
    obj.target_bond_ratio = message.targetBondRatio;
    obj.validator_emission_percentage = message.validatorEmissionPercentage;
    obj.observer_emission_percentage = message.observerEmissionPercentage;
    obj.tss_signer_emission_percentage = message.tssSignerEmissionPercentage;
    obj.duration_factor_constant = message.durationFactorConstant;
    obj.observer_slash_amount = message.observerSlashAmount;
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
      typeUrl: "/zetachain.zetacore.emissions.Params",
      value: Params.encode(message).finish()
    };
  }
};