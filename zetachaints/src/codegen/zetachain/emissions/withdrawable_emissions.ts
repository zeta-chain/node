import { BinaryReader, BinaryWriter } from "../../binary";
export interface WithdrawableEmissions {
  address: string;
  amount: string;
}
export interface WithdrawableEmissionsProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.WithdrawableEmissions";
  value: Uint8Array;
}
export interface WithdrawableEmissionsAmino {
  address?: string;
  amount?: string;
}
export interface WithdrawableEmissionsAminoMsg {
  type: "/zetachain.zetacore.emissions.WithdrawableEmissions";
  value: WithdrawableEmissionsAmino;
}
export interface WithdrawableEmissionsSDKType {
  address: string;
  amount: string;
}
function createBaseWithdrawableEmissions(): WithdrawableEmissions {
  return {
    address: "",
    amount: ""
  };
}
export const WithdrawableEmissions = {
  typeUrl: "/zetachain.zetacore.emissions.WithdrawableEmissions",
  encode(message: WithdrawableEmissions, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.amount !== "") {
      writer.uint32(18).string(message.amount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): WithdrawableEmissions {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWithdrawableEmissions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        case 2:
          message.amount = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<WithdrawableEmissions>): WithdrawableEmissions {
    const message = createBaseWithdrawableEmissions();
    message.address = object.address ?? "";
    message.amount = object.amount ?? "";
    return message;
  },
  fromAmino(object: WithdrawableEmissionsAmino): WithdrawableEmissions {
    const message = createBaseWithdrawableEmissions();
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    return message;
  },
  toAmino(message: WithdrawableEmissions): WithdrawableEmissionsAmino {
    const obj: any = {};
    obj.address = message.address;
    obj.amount = message.amount;
    return obj;
  },
  fromAminoMsg(object: WithdrawableEmissionsAminoMsg): WithdrawableEmissions {
    return WithdrawableEmissions.fromAmino(object.value);
  },
  fromProtoMsg(message: WithdrawableEmissionsProtoMsg): WithdrawableEmissions {
    return WithdrawableEmissions.decode(message.value);
  },
  toProto(message: WithdrawableEmissions): Uint8Array {
    return WithdrawableEmissions.encode(message).finish();
  },
  toProtoMsg(message: WithdrawableEmissions): WithdrawableEmissionsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.WithdrawableEmissions",
      value: WithdrawableEmissions.encode(message).finish()
    };
  }
};