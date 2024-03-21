import { BinaryReader, BinaryWriter } from "../../binary";
export interface TssFundMigratorInfo {
  chainId: bigint;
  migrationCctxIndex: string;
}
export interface TssFundMigratorInfoProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.TssFundMigratorInfo";
  value: Uint8Array;
}
export interface TssFundMigratorInfoAmino {
  chain_id?: string;
  migration_cctx_index?: string;
}
export interface TssFundMigratorInfoAminoMsg {
  type: "/zetachain.zetacore.observer.TssFundMigratorInfo";
  value: TssFundMigratorInfoAmino;
}
export interface TssFundMigratorInfoSDKType {
  chain_id: bigint;
  migration_cctx_index: string;
}
function createBaseTssFundMigratorInfo(): TssFundMigratorInfo {
  return {
    chainId: BigInt(0),
    migrationCctxIndex: ""
  };
}
export const TssFundMigratorInfo = {
  typeUrl: "/zetachain.zetacore.observer.TssFundMigratorInfo",
  encode(message: TssFundMigratorInfo, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.migrationCctxIndex !== "") {
      writer.uint32(18).string(message.migrationCctxIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): TssFundMigratorInfo {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTssFundMigratorInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.migrationCctxIndex = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<TssFundMigratorInfo>): TssFundMigratorInfo {
    const message = createBaseTssFundMigratorInfo();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.migrationCctxIndex = object.migrationCctxIndex ?? "";
    return message;
  },
  fromAmino(object: TssFundMigratorInfoAmino): TssFundMigratorInfo {
    const message = createBaseTssFundMigratorInfo();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.migration_cctx_index !== undefined && object.migration_cctx_index !== null) {
      message.migrationCctxIndex = object.migration_cctx_index;
    }
    return message;
  },
  toAmino(message: TssFundMigratorInfo): TssFundMigratorInfoAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.migration_cctx_index = message.migrationCctxIndex;
    return obj;
  },
  fromAminoMsg(object: TssFundMigratorInfoAminoMsg): TssFundMigratorInfo {
    return TssFundMigratorInfo.fromAmino(object.value);
  },
  fromProtoMsg(message: TssFundMigratorInfoProtoMsg): TssFundMigratorInfo {
    return TssFundMigratorInfo.decode(message.value);
  },
  toProto(message: TssFundMigratorInfo): Uint8Array {
    return TssFundMigratorInfo.encode(message).finish();
  },
  toProtoMsg(message: TssFundMigratorInfo): TssFundMigratorInfoProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.TssFundMigratorInfo",
      value: TssFundMigratorInfo.encode(message).finish()
    };
  }
};