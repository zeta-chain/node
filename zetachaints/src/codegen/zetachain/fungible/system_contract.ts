import { BinaryReader, BinaryWriter } from "../../binary";
export interface SystemContract {
  systemContract: string;
  connectorZevm: string;
}
export interface SystemContractProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.SystemContract";
  value: Uint8Array;
}
export interface SystemContractAmino {
  system_contract?: string;
  connector_zevm?: string;
}
export interface SystemContractAminoMsg {
  type: "/zetachain.zetacore.fungible.SystemContract";
  value: SystemContractAmino;
}
export interface SystemContractSDKType {
  system_contract: string;
  connector_zevm: string;
}
function createBaseSystemContract(): SystemContract {
  return {
    systemContract: "",
    connectorZevm: ""
  };
}
export const SystemContract = {
  typeUrl: "/zetachain.zetacore.fungible.SystemContract",
  encode(message: SystemContract, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.systemContract !== "") {
      writer.uint32(10).string(message.systemContract);
    }
    if (message.connectorZevm !== "") {
      writer.uint32(18).string(message.connectorZevm);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): SystemContract {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSystemContract();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.systemContract = reader.string();
          break;
        case 2:
          message.connectorZevm = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<SystemContract>): SystemContract {
    const message = createBaseSystemContract();
    message.systemContract = object.systemContract ?? "";
    message.connectorZevm = object.connectorZevm ?? "";
    return message;
  },
  fromAmino(object: SystemContractAmino): SystemContract {
    const message = createBaseSystemContract();
    if (object.system_contract !== undefined && object.system_contract !== null) {
      message.systemContract = object.system_contract;
    }
    if (object.connector_zevm !== undefined && object.connector_zevm !== null) {
      message.connectorZevm = object.connector_zevm;
    }
    return message;
  },
  toAmino(message: SystemContract): SystemContractAmino {
    const obj: any = {};
    obj.system_contract = message.systemContract;
    obj.connector_zevm = message.connectorZevm;
    return obj;
  },
  fromAminoMsg(object: SystemContractAminoMsg): SystemContract {
    return SystemContract.fromAmino(object.value);
  },
  fromProtoMsg(message: SystemContractProtoMsg): SystemContract {
    return SystemContract.decode(message.value);
  },
  toProto(message: SystemContract): Uint8Array {
    return SystemContract.encode(message).finish();
  },
  toProtoMsg(message: SystemContract): SystemContractProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.SystemContract",
      value: SystemContract.encode(message).finish()
    };
  }
};