import { Duration, DurationAmino, DurationSDKType } from "../../google/protobuf/duration";
import { BinaryReader, BinaryWriter } from "../../binary";
export interface GasPriceIncreaseFlags {
  epochLength: bigint;
  retryInterval: Duration;
  gasPriceIncreasePercent: number;
  /**
   * Maximum gas price increase in percent of the median gas price
   * Default is used if 0
   */
  gasPriceIncreaseMax: number;
  /** Maximum number of pending crosschain transactions to check for gas price increase */
  maxPendingCctxs: number;
}
export interface GasPriceIncreaseFlagsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.GasPriceIncreaseFlags";
  value: Uint8Array;
}
export interface GasPriceIncreaseFlagsAmino {
  epochLength?: string;
  retryInterval?: DurationAmino;
  gasPriceIncreasePercent?: number;
  /**
   * Maximum gas price increase in percent of the median gas price
   * Default is used if 0
   */
  gasPriceIncreaseMax?: number;
  /** Maximum number of pending crosschain transactions to check for gas price increase */
  maxPendingCctxs?: number;
}
export interface GasPriceIncreaseFlagsAminoMsg {
  type: "/zetachain.zetacore.observer.GasPriceIncreaseFlags";
  value: GasPriceIncreaseFlagsAmino;
}
export interface GasPriceIncreaseFlagsSDKType {
  epochLength: bigint;
  retryInterval: DurationSDKType;
  gasPriceIncreasePercent: number;
  gasPriceIncreaseMax: number;
  maxPendingCctxs: number;
}
export interface BlockHeaderVerificationFlags {
  isEthTypeChainEnabled: boolean;
  isBtcTypeChainEnabled: boolean;
}
export interface BlockHeaderVerificationFlagsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.BlockHeaderVerificationFlags";
  value: Uint8Array;
}
export interface BlockHeaderVerificationFlagsAmino {
  isEthTypeChainEnabled?: boolean;
  isBtcTypeChainEnabled?: boolean;
}
export interface BlockHeaderVerificationFlagsAminoMsg {
  type: "/zetachain.zetacore.observer.BlockHeaderVerificationFlags";
  value: BlockHeaderVerificationFlagsAmino;
}
export interface BlockHeaderVerificationFlagsSDKType {
  isEthTypeChainEnabled: boolean;
  isBtcTypeChainEnabled: boolean;
}
export interface CrosschainFlags {
  isInboundEnabled: boolean;
  isOutboundEnabled: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlags;
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlags;
}
export interface CrosschainFlagsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.CrosschainFlags";
  value: Uint8Array;
}
export interface CrosschainFlagsAmino {
  isInboundEnabled?: boolean;
  isOutboundEnabled?: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlagsAmino;
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlagsAmino;
}
export interface CrosschainFlagsAminoMsg {
  type: "/zetachain.zetacore.observer.CrosschainFlags";
  value: CrosschainFlagsAmino;
}
export interface CrosschainFlagsSDKType {
  isInboundEnabled: boolean;
  isOutboundEnabled: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlagsSDKType;
  blockHeaderVerificationFlags?: BlockHeaderVerificationFlagsSDKType;
}
export interface LegacyCrosschainFlags {
  isInboundEnabled: boolean;
  isOutboundEnabled: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlags;
}
export interface LegacyCrosschainFlagsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.LegacyCrosschainFlags";
  value: Uint8Array;
}
export interface LegacyCrosschainFlagsAmino {
  isInboundEnabled?: boolean;
  isOutboundEnabled?: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlagsAmino;
}
export interface LegacyCrosschainFlagsAminoMsg {
  type: "/zetachain.zetacore.observer.LegacyCrosschainFlags";
  value: LegacyCrosschainFlagsAmino;
}
export interface LegacyCrosschainFlagsSDKType {
  isInboundEnabled: boolean;
  isOutboundEnabled: boolean;
  gasPriceIncreaseFlags?: GasPriceIncreaseFlagsSDKType;
}
function createBaseGasPriceIncreaseFlags(): GasPriceIncreaseFlags {
  return {
    epochLength: BigInt(0),
    retryInterval: Duration.fromPartial({}),
    gasPriceIncreasePercent: 0,
    gasPriceIncreaseMax: 0,
    maxPendingCctxs: 0
  };
}
export const GasPriceIncreaseFlags = {
  typeUrl: "/zetachain.zetacore.observer.GasPriceIncreaseFlags",
  encode(message: GasPriceIncreaseFlags, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.epochLength !== BigInt(0)) {
      writer.uint32(8).int64(message.epochLength);
    }
    if (message.retryInterval !== undefined) {
      Duration.encode(message.retryInterval, writer.uint32(18).fork()).ldelim();
    }
    if (message.gasPriceIncreasePercent !== 0) {
      writer.uint32(24).uint32(message.gasPriceIncreasePercent);
    }
    if (message.gasPriceIncreaseMax !== 0) {
      writer.uint32(32).uint32(message.gasPriceIncreaseMax);
    }
    if (message.maxPendingCctxs !== 0) {
      writer.uint32(40).uint32(message.maxPendingCctxs);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): GasPriceIncreaseFlags {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGasPriceIncreaseFlags();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.epochLength = reader.int64();
          break;
        case 2:
          message.retryInterval = Duration.decode(reader, reader.uint32());
          break;
        case 3:
          message.gasPriceIncreasePercent = reader.uint32();
          break;
        case 4:
          message.gasPriceIncreaseMax = reader.uint32();
          break;
        case 5:
          message.maxPendingCctxs = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<GasPriceIncreaseFlags>): GasPriceIncreaseFlags {
    const message = createBaseGasPriceIncreaseFlags();
    message.epochLength = object.epochLength !== undefined && object.epochLength !== null ? BigInt(object.epochLength.toString()) : BigInt(0);
    message.retryInterval = object.retryInterval !== undefined && object.retryInterval !== null ? Duration.fromPartial(object.retryInterval) : undefined;
    message.gasPriceIncreasePercent = object.gasPriceIncreasePercent ?? 0;
    message.gasPriceIncreaseMax = object.gasPriceIncreaseMax ?? 0;
    message.maxPendingCctxs = object.maxPendingCctxs ?? 0;
    return message;
  },
  fromAmino(object: GasPriceIncreaseFlagsAmino): GasPriceIncreaseFlags {
    const message = createBaseGasPriceIncreaseFlags();
    if (object.epochLength !== undefined && object.epochLength !== null) {
      message.epochLength = BigInt(object.epochLength);
    }
    if (object.retryInterval !== undefined && object.retryInterval !== null) {
      message.retryInterval = Duration.fromAmino(object.retryInterval);
    }
    if (object.gasPriceIncreasePercent !== undefined && object.gasPriceIncreasePercent !== null) {
      message.gasPriceIncreasePercent = object.gasPriceIncreasePercent;
    }
    if (object.gasPriceIncreaseMax !== undefined && object.gasPriceIncreaseMax !== null) {
      message.gasPriceIncreaseMax = object.gasPriceIncreaseMax;
    }
    if (object.maxPendingCctxs !== undefined && object.maxPendingCctxs !== null) {
      message.maxPendingCctxs = object.maxPendingCctxs;
    }
    return message;
  },
  toAmino(message: GasPriceIncreaseFlags): GasPriceIncreaseFlagsAmino {
    const obj: any = {};
    obj.epochLength = message.epochLength ? message.epochLength.toString() : undefined;
    obj.retryInterval = message.retryInterval ? Duration.toAmino(message.retryInterval) : undefined;
    obj.gasPriceIncreasePercent = message.gasPriceIncreasePercent;
    obj.gasPriceIncreaseMax = message.gasPriceIncreaseMax;
    obj.maxPendingCctxs = message.maxPendingCctxs;
    return obj;
  },
  fromAminoMsg(object: GasPriceIncreaseFlagsAminoMsg): GasPriceIncreaseFlags {
    return GasPriceIncreaseFlags.fromAmino(object.value);
  },
  fromProtoMsg(message: GasPriceIncreaseFlagsProtoMsg): GasPriceIncreaseFlags {
    return GasPriceIncreaseFlags.decode(message.value);
  },
  toProto(message: GasPriceIncreaseFlags): Uint8Array {
    return GasPriceIncreaseFlags.encode(message).finish();
  },
  toProtoMsg(message: GasPriceIncreaseFlags): GasPriceIncreaseFlagsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.GasPriceIncreaseFlags",
      value: GasPriceIncreaseFlags.encode(message).finish()
    };
  }
};
function createBaseBlockHeaderVerificationFlags(): BlockHeaderVerificationFlags {
  return {
    isEthTypeChainEnabled: false,
    isBtcTypeChainEnabled: false
  };
}
export const BlockHeaderVerificationFlags = {
  typeUrl: "/zetachain.zetacore.observer.BlockHeaderVerificationFlags",
  encode(message: BlockHeaderVerificationFlags, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.isEthTypeChainEnabled === true) {
      writer.uint32(8).bool(message.isEthTypeChainEnabled);
    }
    if (message.isBtcTypeChainEnabled === true) {
      writer.uint32(16).bool(message.isBtcTypeChainEnabled);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): BlockHeaderVerificationFlags {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockHeaderVerificationFlags();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.isEthTypeChainEnabled = reader.bool();
          break;
        case 2:
          message.isBtcTypeChainEnabled = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<BlockHeaderVerificationFlags>): BlockHeaderVerificationFlags {
    const message = createBaseBlockHeaderVerificationFlags();
    message.isEthTypeChainEnabled = object.isEthTypeChainEnabled ?? false;
    message.isBtcTypeChainEnabled = object.isBtcTypeChainEnabled ?? false;
    return message;
  },
  fromAmino(object: BlockHeaderVerificationFlagsAmino): BlockHeaderVerificationFlags {
    const message = createBaseBlockHeaderVerificationFlags();
    if (object.isEthTypeChainEnabled !== undefined && object.isEthTypeChainEnabled !== null) {
      message.isEthTypeChainEnabled = object.isEthTypeChainEnabled;
    }
    if (object.isBtcTypeChainEnabled !== undefined && object.isBtcTypeChainEnabled !== null) {
      message.isBtcTypeChainEnabled = object.isBtcTypeChainEnabled;
    }
    return message;
  },
  toAmino(message: BlockHeaderVerificationFlags): BlockHeaderVerificationFlagsAmino {
    const obj: any = {};
    obj.isEthTypeChainEnabled = message.isEthTypeChainEnabled;
    obj.isBtcTypeChainEnabled = message.isBtcTypeChainEnabled;
    return obj;
  },
  fromAminoMsg(object: BlockHeaderVerificationFlagsAminoMsg): BlockHeaderVerificationFlags {
    return BlockHeaderVerificationFlags.fromAmino(object.value);
  },
  fromProtoMsg(message: BlockHeaderVerificationFlagsProtoMsg): BlockHeaderVerificationFlags {
    return BlockHeaderVerificationFlags.decode(message.value);
  },
  toProto(message: BlockHeaderVerificationFlags): Uint8Array {
    return BlockHeaderVerificationFlags.encode(message).finish();
  },
  toProtoMsg(message: BlockHeaderVerificationFlags): BlockHeaderVerificationFlagsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.BlockHeaderVerificationFlags",
      value: BlockHeaderVerificationFlags.encode(message).finish()
    };
  }
};
function createBaseCrosschainFlags(): CrosschainFlags {
  return {
    isInboundEnabled: false,
    isOutboundEnabled: false,
    gasPriceIncreaseFlags: undefined,
    blockHeaderVerificationFlags: undefined
  };
}
export const CrosschainFlags = {
  typeUrl: "/zetachain.zetacore.observer.CrosschainFlags",
  encode(message: CrosschainFlags, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.isInboundEnabled === true) {
      writer.uint32(8).bool(message.isInboundEnabled);
    }
    if (message.isOutboundEnabled === true) {
      writer.uint32(16).bool(message.isOutboundEnabled);
    }
    if (message.gasPriceIncreaseFlags !== undefined) {
      GasPriceIncreaseFlags.encode(message.gasPriceIncreaseFlags, writer.uint32(26).fork()).ldelim();
    }
    if (message.blockHeaderVerificationFlags !== undefined) {
      BlockHeaderVerificationFlags.encode(message.blockHeaderVerificationFlags, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): CrosschainFlags {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCrosschainFlags();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.isInboundEnabled = reader.bool();
          break;
        case 2:
          message.isOutboundEnabled = reader.bool();
          break;
        case 3:
          message.gasPriceIncreaseFlags = GasPriceIncreaseFlags.decode(reader, reader.uint32());
          break;
        case 4:
          message.blockHeaderVerificationFlags = BlockHeaderVerificationFlags.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<CrosschainFlags>): CrosschainFlags {
    const message = createBaseCrosschainFlags();
    message.isInboundEnabled = object.isInboundEnabled ?? false;
    message.isOutboundEnabled = object.isOutboundEnabled ?? false;
    message.gasPriceIncreaseFlags = object.gasPriceIncreaseFlags !== undefined && object.gasPriceIncreaseFlags !== null ? GasPriceIncreaseFlags.fromPartial(object.gasPriceIncreaseFlags) : undefined;
    message.blockHeaderVerificationFlags = object.blockHeaderVerificationFlags !== undefined && object.blockHeaderVerificationFlags !== null ? BlockHeaderVerificationFlags.fromPartial(object.blockHeaderVerificationFlags) : undefined;
    return message;
  },
  fromAmino(object: CrosschainFlagsAmino): CrosschainFlags {
    const message = createBaseCrosschainFlags();
    if (object.isInboundEnabled !== undefined && object.isInboundEnabled !== null) {
      message.isInboundEnabled = object.isInboundEnabled;
    }
    if (object.isOutboundEnabled !== undefined && object.isOutboundEnabled !== null) {
      message.isOutboundEnabled = object.isOutboundEnabled;
    }
    if (object.gasPriceIncreaseFlags !== undefined && object.gasPriceIncreaseFlags !== null) {
      message.gasPriceIncreaseFlags = GasPriceIncreaseFlags.fromAmino(object.gasPriceIncreaseFlags);
    }
    if (object.blockHeaderVerificationFlags !== undefined && object.blockHeaderVerificationFlags !== null) {
      message.blockHeaderVerificationFlags = BlockHeaderVerificationFlags.fromAmino(object.blockHeaderVerificationFlags);
    }
    return message;
  },
  toAmino(message: CrosschainFlags): CrosschainFlagsAmino {
    const obj: any = {};
    obj.isInboundEnabled = message.isInboundEnabled;
    obj.isOutboundEnabled = message.isOutboundEnabled;
    obj.gasPriceIncreaseFlags = message.gasPriceIncreaseFlags ? GasPriceIncreaseFlags.toAmino(message.gasPriceIncreaseFlags) : undefined;
    obj.blockHeaderVerificationFlags = message.blockHeaderVerificationFlags ? BlockHeaderVerificationFlags.toAmino(message.blockHeaderVerificationFlags) : undefined;
    return obj;
  },
  fromAminoMsg(object: CrosschainFlagsAminoMsg): CrosschainFlags {
    return CrosschainFlags.fromAmino(object.value);
  },
  fromProtoMsg(message: CrosschainFlagsProtoMsg): CrosschainFlags {
    return CrosschainFlags.decode(message.value);
  },
  toProto(message: CrosschainFlags): Uint8Array {
    return CrosschainFlags.encode(message).finish();
  },
  toProtoMsg(message: CrosschainFlags): CrosschainFlagsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.CrosschainFlags",
      value: CrosschainFlags.encode(message).finish()
    };
  }
};
function createBaseLegacyCrosschainFlags(): LegacyCrosschainFlags {
  return {
    isInboundEnabled: false,
    isOutboundEnabled: false,
    gasPriceIncreaseFlags: undefined
  };
}
export const LegacyCrosschainFlags = {
  typeUrl: "/zetachain.zetacore.observer.LegacyCrosschainFlags",
  encode(message: LegacyCrosschainFlags, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.isInboundEnabled === true) {
      writer.uint32(8).bool(message.isInboundEnabled);
    }
    if (message.isOutboundEnabled === true) {
      writer.uint32(16).bool(message.isOutboundEnabled);
    }
    if (message.gasPriceIncreaseFlags !== undefined) {
      GasPriceIncreaseFlags.encode(message.gasPriceIncreaseFlags, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): LegacyCrosschainFlags {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLegacyCrosschainFlags();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.isInboundEnabled = reader.bool();
          break;
        case 2:
          message.isOutboundEnabled = reader.bool();
          break;
        case 3:
          message.gasPriceIncreaseFlags = GasPriceIncreaseFlags.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<LegacyCrosschainFlags>): LegacyCrosschainFlags {
    const message = createBaseLegacyCrosschainFlags();
    message.isInboundEnabled = object.isInboundEnabled ?? false;
    message.isOutboundEnabled = object.isOutboundEnabled ?? false;
    message.gasPriceIncreaseFlags = object.gasPriceIncreaseFlags !== undefined && object.gasPriceIncreaseFlags !== null ? GasPriceIncreaseFlags.fromPartial(object.gasPriceIncreaseFlags) : undefined;
    return message;
  },
  fromAmino(object: LegacyCrosschainFlagsAmino): LegacyCrosschainFlags {
    const message = createBaseLegacyCrosschainFlags();
    if (object.isInboundEnabled !== undefined && object.isInboundEnabled !== null) {
      message.isInboundEnabled = object.isInboundEnabled;
    }
    if (object.isOutboundEnabled !== undefined && object.isOutboundEnabled !== null) {
      message.isOutboundEnabled = object.isOutboundEnabled;
    }
    if (object.gasPriceIncreaseFlags !== undefined && object.gasPriceIncreaseFlags !== null) {
      message.gasPriceIncreaseFlags = GasPriceIncreaseFlags.fromAmino(object.gasPriceIncreaseFlags);
    }
    return message;
  },
  toAmino(message: LegacyCrosschainFlags): LegacyCrosschainFlagsAmino {
    const obj: any = {};
    obj.isInboundEnabled = message.isInboundEnabled;
    obj.isOutboundEnabled = message.isOutboundEnabled;
    obj.gasPriceIncreaseFlags = message.gasPriceIncreaseFlags ? GasPriceIncreaseFlags.toAmino(message.gasPriceIncreaseFlags) : undefined;
    return obj;
  },
  fromAminoMsg(object: LegacyCrosschainFlagsAminoMsg): LegacyCrosschainFlags {
    return LegacyCrosschainFlags.fromAmino(object.value);
  },
  fromProtoMsg(message: LegacyCrosschainFlagsProtoMsg): LegacyCrosschainFlags {
    return LegacyCrosschainFlags.decode(message.value);
  },
  toProto(message: LegacyCrosschainFlags): Uint8Array {
    return LegacyCrosschainFlags.encode(message).finish();
  },
  toProtoMsg(message: LegacyCrosschainFlags): LegacyCrosschainFlagsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.LegacyCrosschainFlags",
      value: LegacyCrosschainFlags.encode(message).finish()
    };
  }
};