import { BinaryReader, BinaryWriter } from "../../binary";
export interface TSS {
  tssPubkey: string;
  tssParticipantList: string[];
  operatorAddressList: string[];
  finalizedZetaHeight: bigint;
  keyGenZetaHeight: bigint;
}
export interface TSSProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.TSS";
  value: Uint8Array;
}
export interface TSSAmino {
  tss_pubkey?: string;
  tss_participant_list?: string[];
  operator_address_list?: string[];
  finalizedZetaHeight?: string;
  keyGenZetaHeight?: string;
}
export interface TSSAminoMsg {
  type: "/zetachain.zetacore.observer.TSS";
  value: TSSAmino;
}
export interface TSSSDKType {
  tss_pubkey: string;
  tss_participant_list: string[];
  operator_address_list: string[];
  finalizedZetaHeight: bigint;
  keyGenZetaHeight: bigint;
}
function createBaseTSS(): TSS {
  return {
    tssPubkey: "",
    tssParticipantList: [],
    operatorAddressList: [],
    finalizedZetaHeight: BigInt(0),
    keyGenZetaHeight: BigInt(0)
  };
}
export const TSS = {
  typeUrl: "/zetachain.zetacore.observer.TSS",
  encode(message: TSS, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.tssPubkey !== "") {
      writer.uint32(26).string(message.tssPubkey);
    }
    for (const v of message.tssParticipantList) {
      writer.uint32(34).string(v!);
    }
    for (const v of message.operatorAddressList) {
      writer.uint32(42).string(v!);
    }
    if (message.finalizedZetaHeight !== BigInt(0)) {
      writer.uint32(48).int64(message.finalizedZetaHeight);
    }
    if (message.keyGenZetaHeight !== BigInt(0)) {
      writer.uint32(56).int64(message.keyGenZetaHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): TSS {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTSS();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 3:
          message.tssPubkey = reader.string();
          break;
        case 4:
          message.tssParticipantList.push(reader.string());
          break;
        case 5:
          message.operatorAddressList.push(reader.string());
          break;
        case 6:
          message.finalizedZetaHeight = reader.int64();
          break;
        case 7:
          message.keyGenZetaHeight = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<TSS>): TSS {
    const message = createBaseTSS();
    message.tssPubkey = object.tssPubkey ?? "";
    message.tssParticipantList = object.tssParticipantList?.map(e => e) || [];
    message.operatorAddressList = object.operatorAddressList?.map(e => e) || [];
    message.finalizedZetaHeight = object.finalizedZetaHeight !== undefined && object.finalizedZetaHeight !== null ? BigInt(object.finalizedZetaHeight.toString()) : BigInt(0);
    message.keyGenZetaHeight = object.keyGenZetaHeight !== undefined && object.keyGenZetaHeight !== null ? BigInt(object.keyGenZetaHeight.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: TSSAmino): TSS {
    const message = createBaseTSS();
    if (object.tss_pubkey !== undefined && object.tss_pubkey !== null) {
      message.tssPubkey = object.tss_pubkey;
    }
    message.tssParticipantList = object.tss_participant_list?.map(e => e) || [];
    message.operatorAddressList = object.operator_address_list?.map(e => e) || [];
    if (object.finalizedZetaHeight !== undefined && object.finalizedZetaHeight !== null) {
      message.finalizedZetaHeight = BigInt(object.finalizedZetaHeight);
    }
    if (object.keyGenZetaHeight !== undefined && object.keyGenZetaHeight !== null) {
      message.keyGenZetaHeight = BigInt(object.keyGenZetaHeight);
    }
    return message;
  },
  toAmino(message: TSS): TSSAmino {
    const obj: any = {};
    obj.tss_pubkey = message.tssPubkey;
    if (message.tssParticipantList) {
      obj.tss_participant_list = message.tssParticipantList.map(e => e);
    } else {
      obj.tss_participant_list = [];
    }
    if (message.operatorAddressList) {
      obj.operator_address_list = message.operatorAddressList.map(e => e);
    } else {
      obj.operator_address_list = [];
    }
    obj.finalizedZetaHeight = message.finalizedZetaHeight ? message.finalizedZetaHeight.toString() : undefined;
    obj.keyGenZetaHeight = message.keyGenZetaHeight ? message.keyGenZetaHeight.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: TSSAminoMsg): TSS {
    return TSS.fromAmino(object.value);
  },
  fromProtoMsg(message: TSSProtoMsg): TSS {
    return TSS.decode(message.value);
  },
  toProto(message: TSS): Uint8Array {
    return TSS.encode(message).finish();
  },
  toProtoMsg(message: TSS): TSSProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.TSS",
      value: TSS.encode(message).finish()
    };
  }
};