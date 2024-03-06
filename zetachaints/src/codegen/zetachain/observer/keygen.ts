import { BinaryReader, BinaryWriter } from "../../binary";
export enum KeygenStatus {
  PendingKeygen = 0,
  KeyGenSuccess = 1,
  KeyGenFailed = 3,
  UNRECOGNIZED = -1,
}
export const KeygenStatusSDKType = KeygenStatus;
export const KeygenStatusAmino = KeygenStatus;
export function keygenStatusFromJSON(object: any): KeygenStatus {
  switch (object) {
    case 0:
    case "PendingKeygen":
      return KeygenStatus.PendingKeygen;
    case 1:
    case "KeyGenSuccess":
      return KeygenStatus.KeyGenSuccess;
    case 3:
    case "KeyGenFailed":
      return KeygenStatus.KeyGenFailed;
    case -1:
    case "UNRECOGNIZED":
    default:
      return KeygenStatus.UNRECOGNIZED;
  }
}
export function keygenStatusToJSON(object: KeygenStatus): string {
  switch (object) {
    case KeygenStatus.PendingKeygen:
      return "PendingKeygen";
    case KeygenStatus.KeyGenSuccess:
      return "KeyGenSuccess";
    case KeygenStatus.KeyGenFailed:
      return "KeyGenFailed";
    case KeygenStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export interface Keygen {
  /** 0--to generate key; 1--generated; 2--error */
  status: KeygenStatus;
  granteePubkeys: string[];
  /** the blocknum that the key needs to be generated */
  blockNumber: bigint;
}
export interface KeygenProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.Keygen";
  value: Uint8Array;
}
export interface KeygenAmino {
  /** 0--to generate key; 1--generated; 2--error */
  status?: KeygenStatus;
  granteePubkeys?: string[];
  /** the blocknum that the key needs to be generated */
  blockNumber?: string;
}
export interface KeygenAminoMsg {
  type: "/zetachain.zetacore.observer.Keygen";
  value: KeygenAmino;
}
export interface KeygenSDKType {
  status: KeygenStatus;
  granteePubkeys: string[];
  blockNumber: bigint;
}
function createBaseKeygen(): Keygen {
  return {
    status: 0,
    granteePubkeys: [],
    blockNumber: BigInt(0)
  };
}
export const Keygen = {
  typeUrl: "/zetachain.zetacore.observer.Keygen",
  encode(message: Keygen, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.status !== 0) {
      writer.uint32(16).int32(message.status);
    }
    for (const v of message.granteePubkeys) {
      writer.uint32(26).string(v!);
    }
    if (message.blockNumber !== BigInt(0)) {
      writer.uint32(32).int64(message.blockNumber);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Keygen {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseKeygen();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.status = (reader.int32() as any);
          break;
        case 3:
          message.granteePubkeys.push(reader.string());
          break;
        case 4:
          message.blockNumber = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Keygen>): Keygen {
    const message = createBaseKeygen();
    message.status = object.status ?? 0;
    message.granteePubkeys = object.granteePubkeys?.map(e => e) || [];
    message.blockNumber = object.blockNumber !== undefined && object.blockNumber !== null ? BigInt(object.blockNumber.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: KeygenAmino): Keygen {
    const message = createBaseKeygen();
    if (object.status !== undefined && object.status !== null) {
      message.status = keygenStatusFromJSON(object.status);
    }
    message.granteePubkeys = object.granteePubkeys?.map(e => e) || [];
    if (object.blockNumber !== undefined && object.blockNumber !== null) {
      message.blockNumber = BigInt(object.blockNumber);
    }
    return message;
  },
  toAmino(message: Keygen): KeygenAmino {
    const obj: any = {};
    obj.status = message.status;
    if (message.granteePubkeys) {
      obj.granteePubkeys = message.granteePubkeys.map(e => e);
    } else {
      obj.granteePubkeys = [];
    }
    obj.blockNumber = message.blockNumber ? message.blockNumber.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: KeygenAminoMsg): Keygen {
    return Keygen.fromAmino(object.value);
  },
  fromProtoMsg(message: KeygenProtoMsg): Keygen {
    return Keygen.decode(message.value);
  },
  toProto(message: Keygen): Uint8Array {
    return Keygen.encode(message).finish();
  },
  toProtoMsg(message: Keygen): KeygenProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.Keygen",
      value: Keygen.encode(message).finish()
    };
  }
};