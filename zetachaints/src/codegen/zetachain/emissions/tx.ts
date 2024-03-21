import { BinaryReader, BinaryWriter } from "../../binary";
export interface MsgWithdrawEmission {
  creator: string;
  amount: string;
}
export interface MsgWithdrawEmissionProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.MsgWithdrawEmission";
  value: Uint8Array;
}
export interface MsgWithdrawEmissionAmino {
  creator?: string;
  amount?: string;
}
export interface MsgWithdrawEmissionAminoMsg {
  type: "/zetachain.zetacore.emissions.MsgWithdrawEmission";
  value: MsgWithdrawEmissionAmino;
}
export interface MsgWithdrawEmissionSDKType {
  creator: string;
  amount: string;
}
export interface MsgWithdrawEmissionResponse {}
export interface MsgWithdrawEmissionResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.emissions.MsgWithdrawEmissionResponse";
  value: Uint8Array;
}
export interface MsgWithdrawEmissionResponseAmino {}
export interface MsgWithdrawEmissionResponseAminoMsg {
  type: "/zetachain.zetacore.emissions.MsgWithdrawEmissionResponse";
  value: MsgWithdrawEmissionResponseAmino;
}
export interface MsgWithdrawEmissionResponseSDKType {}
function createBaseMsgWithdrawEmission(): MsgWithdrawEmission {
  return {
    creator: "",
    amount: ""
  };
}
export const MsgWithdrawEmission = {
  typeUrl: "/zetachain.zetacore.emissions.MsgWithdrawEmission",
  encode(message: MsgWithdrawEmission, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.amount !== "") {
      writer.uint32(18).string(message.amount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgWithdrawEmission {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawEmission();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
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
  fromPartial(object: Partial<MsgWithdrawEmission>): MsgWithdrawEmission {
    const message = createBaseMsgWithdrawEmission();
    message.creator = object.creator ?? "";
    message.amount = object.amount ?? "";
    return message;
  },
  fromAmino(object: MsgWithdrawEmissionAmino): MsgWithdrawEmission {
    const message = createBaseMsgWithdrawEmission();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    return message;
  },
  toAmino(message: MsgWithdrawEmission): MsgWithdrawEmissionAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.amount = message.amount;
    return obj;
  },
  fromAminoMsg(object: MsgWithdrawEmissionAminoMsg): MsgWithdrawEmission {
    return MsgWithdrawEmission.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgWithdrawEmissionProtoMsg): MsgWithdrawEmission {
    return MsgWithdrawEmission.decode(message.value);
  },
  toProto(message: MsgWithdrawEmission): Uint8Array {
    return MsgWithdrawEmission.encode(message).finish();
  },
  toProtoMsg(message: MsgWithdrawEmission): MsgWithdrawEmissionProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.MsgWithdrawEmission",
      value: MsgWithdrawEmission.encode(message).finish()
    };
  }
};
function createBaseMsgWithdrawEmissionResponse(): MsgWithdrawEmissionResponse {
  return {};
}
export const MsgWithdrawEmissionResponse = {
  typeUrl: "/zetachain.zetacore.emissions.MsgWithdrawEmissionResponse",
  encode(_: MsgWithdrawEmissionResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgWithdrawEmissionResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWithdrawEmissionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(_: Partial<MsgWithdrawEmissionResponse>): MsgWithdrawEmissionResponse {
    const message = createBaseMsgWithdrawEmissionResponse();
    return message;
  },
  fromAmino(_: MsgWithdrawEmissionResponseAmino): MsgWithdrawEmissionResponse {
    const message = createBaseMsgWithdrawEmissionResponse();
    return message;
  },
  toAmino(_: MsgWithdrawEmissionResponse): MsgWithdrawEmissionResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgWithdrawEmissionResponseAminoMsg): MsgWithdrawEmissionResponse {
    return MsgWithdrawEmissionResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgWithdrawEmissionResponseProtoMsg): MsgWithdrawEmissionResponse {
    return MsgWithdrawEmissionResponse.decode(message.value);
  },
  toProto(message: MsgWithdrawEmissionResponse): Uint8Array {
    return MsgWithdrawEmissionResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgWithdrawEmissionResponse): MsgWithdrawEmissionResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.emissions.MsgWithdrawEmissionResponse",
      value: MsgWithdrawEmissionResponse.encode(message).finish()
    };
  }
};