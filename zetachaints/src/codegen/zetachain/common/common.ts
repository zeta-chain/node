import { Proof as Proof1 } from "./ethereum/ethereum";
import { ProofAmino as Proof1Amino } from "./ethereum/ethereum";
import { ProofSDKType as Proof1SDKType } from "./ethereum/ethereum";
import { Proof as Proof2 } from "./bitcoin/bitcoin";
import { ProofAmino as Proof2Amino } from "./bitcoin/bitcoin";
import { ProofSDKType as Proof2SDKType } from "./bitcoin/bitcoin";
import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
export enum ReceiveStatus {
  /** Created - some observer sees inbound tx */
  Created = 0,
  Success = 1,
  Failed = 2,
  UNRECOGNIZED = -1,
}
export const ReceiveStatusSDKType = ReceiveStatus;
export const ReceiveStatusAmino = ReceiveStatus;
export function receiveStatusFromJSON(object: any): ReceiveStatus {
  switch (object) {
    case 0:
    case "Created":
      return ReceiveStatus.Created;
    case 1:
    case "Success":
      return ReceiveStatus.Success;
    case 2:
    case "Failed":
      return ReceiveStatus.Failed;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ReceiveStatus.UNRECOGNIZED;
  }
}
export function receiveStatusToJSON(object: ReceiveStatus): string {
  switch (object) {
    case ReceiveStatus.Created:
      return "Created";
    case ReceiveStatus.Success:
      return "Success";
    case ReceiveStatus.Failed:
      return "Failed";
    case ReceiveStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export enum CoinType {
  Zeta = 0,
  /** Gas - Ether, BNB, Matic, Klay, BTC, etc */
  Gas = 1,
  /** ERC20 - ERC20 token */
  ERC20 = 2,
  /** Cmd - not a real coin, rather a command */
  Cmd = 3,
  UNRECOGNIZED = -1,
}
export const CoinTypeSDKType = CoinType;
export const CoinTypeAmino = CoinType;
export function coinTypeFromJSON(object: any): CoinType {
  switch (object) {
    case 0:
    case "Zeta":
      return CoinType.Zeta;
    case 1:
    case "Gas":
      return CoinType.Gas;
    case 2:
    case "ERC20":
      return CoinType.ERC20;
    case 3:
    case "Cmd":
      return CoinType.Cmd;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CoinType.UNRECOGNIZED;
  }
}
export function coinTypeToJSON(object: CoinType): string {
  switch (object) {
    case CoinType.Zeta:
      return "Zeta";
    case CoinType.Gas:
      return "Gas";
    case CoinType.ERC20:
      return "ERC20";
    case CoinType.Cmd:
      return "Cmd";
    case CoinType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export enum ChainName {
  empty = 0,
  eth_mainnet = 1,
  zeta_mainnet = 2,
  btc_mainnet = 3,
  polygon_mainnet = 4,
  bsc_mainnet = 5,
  /** goerli_testnet - Testnet */
  goerli_testnet = 6,
  mumbai_testnet = 7,
  ganache_testnet = 8,
  baobab_testnet = 9,
  bsc_testnet = 10,
  zeta_testnet = 11,
  btc_testnet = 12,
  sepolia_testnet = 13,
  /**
   * goerli_localnet - LocalNet
   *  zeta_localnet = 13;
   */
  goerli_localnet = 14,
  btc_regtest = 15,
  UNRECOGNIZED = -1,
}
export const ChainNameSDKType = ChainName;
export const ChainNameAmino = ChainName;
export function chainNameFromJSON(object: any): ChainName {
  switch (object) {
    case 0:
    case "empty":
      return ChainName.empty;
    case 1:
    case "eth_mainnet":
      return ChainName.eth_mainnet;
    case 2:
    case "zeta_mainnet":
      return ChainName.zeta_mainnet;
    case 3:
    case "btc_mainnet":
      return ChainName.btc_mainnet;
    case 4:
    case "polygon_mainnet":
      return ChainName.polygon_mainnet;
    case 5:
    case "bsc_mainnet":
      return ChainName.bsc_mainnet;
    case 6:
    case "goerli_testnet":
      return ChainName.goerli_testnet;
    case 7:
    case "mumbai_testnet":
      return ChainName.mumbai_testnet;
    case 8:
    case "ganache_testnet":
      return ChainName.ganache_testnet;
    case 9:
    case "baobab_testnet":
      return ChainName.baobab_testnet;
    case 10:
    case "bsc_testnet":
      return ChainName.bsc_testnet;
    case 11:
    case "zeta_testnet":
      return ChainName.zeta_testnet;
    case 12:
    case "btc_testnet":
      return ChainName.btc_testnet;
    case 13:
    case "sepolia_testnet":
      return ChainName.sepolia_testnet;
    case 14:
    case "goerli_localnet":
      return ChainName.goerli_localnet;
    case 15:
    case "btc_regtest":
      return ChainName.btc_regtest;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ChainName.UNRECOGNIZED;
  }
}
export function chainNameToJSON(object: ChainName): string {
  switch (object) {
    case ChainName.empty:
      return "empty";
    case ChainName.eth_mainnet:
      return "eth_mainnet";
    case ChainName.zeta_mainnet:
      return "zeta_mainnet";
    case ChainName.btc_mainnet:
      return "btc_mainnet";
    case ChainName.polygon_mainnet:
      return "polygon_mainnet";
    case ChainName.bsc_mainnet:
      return "bsc_mainnet";
    case ChainName.goerli_testnet:
      return "goerli_testnet";
    case ChainName.mumbai_testnet:
      return "mumbai_testnet";
    case ChainName.ganache_testnet:
      return "ganache_testnet";
    case ChainName.baobab_testnet:
      return "baobab_testnet";
    case ChainName.bsc_testnet:
      return "bsc_testnet";
    case ChainName.zeta_testnet:
      return "zeta_testnet";
    case ChainName.btc_testnet:
      return "btc_testnet";
    case ChainName.sepolia_testnet:
      return "sepolia_testnet";
    case ChainName.goerli_localnet:
      return "goerli_localnet";
    case ChainName.btc_regtest:
      return "btc_regtest";
    case ChainName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** PubKeySet contains two pub keys , secp256k1 and ed25519 */
export interface PubKeySet {
  secp256k1: string;
  ed25519: string;
}
export interface PubKeySetProtoMsg {
  typeUrl: "/common.PubKeySet";
  value: Uint8Array;
}
/** PubKeySet contains two pub keys , secp256k1 and ed25519 */
export interface PubKeySetAmino {
  secp256k1?: string;
  ed25519?: string;
}
export interface PubKeySetAminoMsg {
  type: "/common.PubKeySet";
  value: PubKeySetAmino;
}
/** PubKeySet contains two pub keys , secp256k1 and ed25519 */
export interface PubKeySetSDKType {
  secp256k1: string;
  ed25519: string;
}
export interface Chain {
  chainName: ChainName;
  chainId: bigint;
}
export interface ChainProtoMsg {
  typeUrl: "/common.Chain";
  value: Uint8Array;
}
export interface ChainAmino {
  chain_name?: ChainName;
  chain_id?: string;
}
export interface ChainAminoMsg {
  type: "/common.Chain";
  value: ChainAmino;
}
export interface ChainSDKType {
  chain_name: ChainName;
  chain_id: bigint;
}
export interface BlockHeader {
  height: bigint;
  hash: Uint8Array;
  parentHash: Uint8Array;
  chainId: bigint;
  /** chain specific header */
  header: HeaderData;
}
export interface BlockHeaderProtoMsg {
  typeUrl: "/common.BlockHeader";
  value: Uint8Array;
}
export interface BlockHeaderAmino {
  height?: string;
  hash?: string;
  parent_hash?: string;
  chain_id?: string;
  /** chain specific header */
  header?: HeaderDataAmino;
}
export interface BlockHeaderAminoMsg {
  type: "/common.BlockHeader";
  value: BlockHeaderAmino;
}
export interface BlockHeaderSDKType {
  height: bigint;
  hash: Uint8Array;
  parent_hash: Uint8Array;
  chain_id: bigint;
  header: HeaderDataSDKType;
}
export interface HeaderData {
  /** binary encoded headers; RLP for ethereum */
  ethereumHeader?: Uint8Array;
  /** 80-byte little-endian encoded binary data */
  bitcoinHeader?: Uint8Array;
}
export interface HeaderDataProtoMsg {
  typeUrl: "/common.HeaderData";
  value: Uint8Array;
}
export interface HeaderDataAmino {
  /** binary encoded headers; RLP for ethereum */
  ethereum_header?: string;
  /** 80-byte little-endian encoded binary data */
  bitcoin_header?: string;
}
export interface HeaderDataAminoMsg {
  type: "/common.HeaderData";
  value: HeaderDataAmino;
}
export interface HeaderDataSDKType {
  ethereum_header?: Uint8Array;
  bitcoin_header?: Uint8Array;
}
export interface Proof {
  ethereumProof?: Proof1;
  bitcoinProof?: Proof2;
}
export interface ProofProtoMsg {
  typeUrl: "/common.Proof";
  value: Uint8Array;
}
export interface ProofAmino {
  ethereum_proof?: Proof1Amino;
  bitcoin_proof?: Proof2Amino;
}
export interface ProofAminoMsg {
  type: "/common.Proof";
  value: ProofAmino;
}
export interface ProofSDKType {
  ethereum_proof?: Proof1SDKType;
  bitcoin_proof?: Proof2SDKType;
}
function createBasePubKeySet(): PubKeySet {
  return {
    secp256k1: "",
    ed25519: ""
  };
}
export const PubKeySet = {
  typeUrl: "/common.PubKeySet",
  encode(message: PubKeySet, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.secp256k1 !== "") {
      writer.uint32(10).string(message.secp256k1);
    }
    if (message.ed25519 !== "") {
      writer.uint32(18).string(message.ed25519);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): PubKeySet {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePubKeySet();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.secp256k1 = reader.string();
          break;
        case 2:
          message.ed25519 = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<PubKeySet>): PubKeySet {
    const message = createBasePubKeySet();
    message.secp256k1 = object.secp256k1 ?? "";
    message.ed25519 = object.ed25519 ?? "";
    return message;
  },
  fromAmino(object: PubKeySetAmino): PubKeySet {
    const message = createBasePubKeySet();
    if (object.secp256k1 !== undefined && object.secp256k1 !== null) {
      message.secp256k1 = object.secp256k1;
    }
    if (object.ed25519 !== undefined && object.ed25519 !== null) {
      message.ed25519 = object.ed25519;
    }
    return message;
  },
  toAmino(message: PubKeySet): PubKeySetAmino {
    const obj: any = {};
    obj.secp256k1 = message.secp256k1;
    obj.ed25519 = message.ed25519;
    return obj;
  },
  fromAminoMsg(object: PubKeySetAminoMsg): PubKeySet {
    return PubKeySet.fromAmino(object.value);
  },
  fromProtoMsg(message: PubKeySetProtoMsg): PubKeySet {
    return PubKeySet.decode(message.value);
  },
  toProto(message: PubKeySet): Uint8Array {
    return PubKeySet.encode(message).finish();
  },
  toProtoMsg(message: PubKeySet): PubKeySetProtoMsg {
    return {
      typeUrl: "/common.PubKeySet",
      value: PubKeySet.encode(message).finish()
    };
  }
};
function createBaseChain(): Chain {
  return {
    chainName: 0,
    chainId: BigInt(0)
  };
}
export const Chain = {
  typeUrl: "/common.Chain",
  encode(message: Chain, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainName !== 0) {
      writer.uint32(8).int32(message.chainName);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Chain {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChain();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainName = (reader.int32() as any);
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Chain>): Chain {
    const message = createBaseChain();
    message.chainName = object.chainName ?? 0;
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: ChainAmino): Chain {
    const message = createBaseChain();
    if (object.chain_name !== undefined && object.chain_name !== null) {
      message.chainName = chainNameFromJSON(object.chain_name);
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    return message;
  },
  toAmino(message: Chain): ChainAmino {
    const obj: any = {};
    obj.chain_name = message.chainName;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: ChainAminoMsg): Chain {
    return Chain.fromAmino(object.value);
  },
  fromProtoMsg(message: ChainProtoMsg): Chain {
    return Chain.decode(message.value);
  },
  toProto(message: Chain): Uint8Array {
    return Chain.encode(message).finish();
  },
  toProtoMsg(message: Chain): ChainProtoMsg {
    return {
      typeUrl: "/common.Chain",
      value: Chain.encode(message).finish()
    };
  }
};
function createBaseBlockHeader(): BlockHeader {
  return {
    height: BigInt(0),
    hash: new Uint8Array(),
    parentHash: new Uint8Array(),
    chainId: BigInt(0),
    header: HeaderData.fromPartial({})
  };
}
export const BlockHeader = {
  typeUrl: "/common.BlockHeader",
  encode(message: BlockHeader, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.height !== BigInt(0)) {
      writer.uint32(8).int64(message.height);
    }
    if (message.hash.length !== 0) {
      writer.uint32(18).bytes(message.hash);
    }
    if (message.parentHash.length !== 0) {
      writer.uint32(26).bytes(message.parentHash);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(32).int64(message.chainId);
    }
    if (message.header !== undefined) {
      HeaderData.encode(message.header, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): BlockHeader {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockHeader();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.height = reader.int64();
          break;
        case 2:
          message.hash = reader.bytes();
          break;
        case 3:
          message.parentHash = reader.bytes();
          break;
        case 4:
          message.chainId = reader.int64();
          break;
        case 5:
          message.header = HeaderData.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<BlockHeader>): BlockHeader {
    const message = createBaseBlockHeader();
    message.height = object.height !== undefined && object.height !== null ? BigInt(object.height.toString()) : BigInt(0);
    message.hash = object.hash ?? new Uint8Array();
    message.parentHash = object.parentHash ?? new Uint8Array();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.header = object.header !== undefined && object.header !== null ? HeaderData.fromPartial(object.header) : undefined;
    return message;
  },
  fromAmino(object: BlockHeaderAmino): BlockHeader {
    const message = createBaseBlockHeader();
    if (object.height !== undefined && object.height !== null) {
      message.height = BigInt(object.height);
    }
    if (object.hash !== undefined && object.hash !== null) {
      message.hash = bytesFromBase64(object.hash);
    }
    if (object.parent_hash !== undefined && object.parent_hash !== null) {
      message.parentHash = bytesFromBase64(object.parent_hash);
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.header !== undefined && object.header !== null) {
      message.header = HeaderData.fromAmino(object.header);
    }
    return message;
  },
  toAmino(message: BlockHeader): BlockHeaderAmino {
    const obj: any = {};
    obj.height = message.height ? message.height.toString() : undefined;
    obj.hash = message.hash ? base64FromBytes(message.hash) : undefined;
    obj.parent_hash = message.parentHash ? base64FromBytes(message.parentHash) : undefined;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.header = message.header ? HeaderData.toAmino(message.header) : undefined;
    return obj;
  },
  fromAminoMsg(object: BlockHeaderAminoMsg): BlockHeader {
    return BlockHeader.fromAmino(object.value);
  },
  fromProtoMsg(message: BlockHeaderProtoMsg): BlockHeader {
    return BlockHeader.decode(message.value);
  },
  toProto(message: BlockHeader): Uint8Array {
    return BlockHeader.encode(message).finish();
  },
  toProtoMsg(message: BlockHeader): BlockHeaderProtoMsg {
    return {
      typeUrl: "/common.BlockHeader",
      value: BlockHeader.encode(message).finish()
    };
  }
};
function createBaseHeaderData(): HeaderData {
  return {
    ethereumHeader: undefined,
    bitcoinHeader: undefined
  };
}
export const HeaderData = {
  typeUrl: "/common.HeaderData",
  encode(message: HeaderData, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.ethereumHeader !== undefined) {
      writer.uint32(10).bytes(message.ethereumHeader);
    }
    if (message.bitcoinHeader !== undefined) {
      writer.uint32(18).bytes(message.bitcoinHeader);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): HeaderData {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHeaderData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ethereumHeader = reader.bytes();
          break;
        case 2:
          message.bitcoinHeader = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<HeaderData>): HeaderData {
    const message = createBaseHeaderData();
    message.ethereumHeader = object.ethereumHeader ?? undefined;
    message.bitcoinHeader = object.bitcoinHeader ?? undefined;
    return message;
  },
  fromAmino(object: HeaderDataAmino): HeaderData {
    const message = createBaseHeaderData();
    if (object.ethereum_header !== undefined && object.ethereum_header !== null) {
      message.ethereumHeader = bytesFromBase64(object.ethereum_header);
    }
    if (object.bitcoin_header !== undefined && object.bitcoin_header !== null) {
      message.bitcoinHeader = bytesFromBase64(object.bitcoin_header);
    }
    return message;
  },
  toAmino(message: HeaderData): HeaderDataAmino {
    const obj: any = {};
    obj.ethereum_header = message.ethereumHeader ? base64FromBytes(message.ethereumHeader) : undefined;
    obj.bitcoin_header = message.bitcoinHeader ? base64FromBytes(message.bitcoinHeader) : undefined;
    return obj;
  },
  fromAminoMsg(object: HeaderDataAminoMsg): HeaderData {
    return HeaderData.fromAmino(object.value);
  },
  fromProtoMsg(message: HeaderDataProtoMsg): HeaderData {
    return HeaderData.decode(message.value);
  },
  toProto(message: HeaderData): Uint8Array {
    return HeaderData.encode(message).finish();
  },
  toProtoMsg(message: HeaderData): HeaderDataProtoMsg {
    return {
      typeUrl: "/common.HeaderData",
      value: HeaderData.encode(message).finish()
    };
  }
};
function createBaseProof(): Proof {
  return {
    ethereumProof: undefined,
    bitcoinProof: undefined
  };
}
export const Proof = {
  typeUrl: "/common.Proof",
  encode(message: Proof, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.ethereumProof !== undefined) {
      Proof1.encode(message.ethereumProof, writer.uint32(10).fork()).ldelim();
    }
    if (message.bitcoinProof !== undefined) {
      Proof2.encode(message.bitcoinProof, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Proof {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProof();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ethereumProof = Proof1.decode(reader, reader.uint32());
          break;
        case 2:
          message.bitcoinProof = Proof2.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Proof>): Proof {
    const message = createBaseProof();
    message.ethereumProof = object.ethereumProof !== undefined && object.ethereumProof !== null ? Proof1.fromPartial(object.ethereumProof) : undefined;
    message.bitcoinProof = object.bitcoinProof !== undefined && object.bitcoinProof !== null ? Proof2.fromPartial(object.bitcoinProof) : undefined;
    return message;
  },
  fromAmino(object: ProofAmino): Proof {
    const message = createBaseProof();
    if (object.ethereum_proof !== undefined && object.ethereum_proof !== null) {
      message.ethereumProof = Proof1.fromAmino(object.ethereum_proof);
    }
    if (object.bitcoin_proof !== undefined && object.bitcoin_proof !== null) {
      message.bitcoinProof = Proof2.fromAmino(object.bitcoin_proof);
    }
    return message;
  },
  toAmino(message: Proof): ProofAmino {
    const obj: any = {};
    obj.ethereum_proof = message.ethereumProof ? Proof1.toAmino(message.ethereumProof) : undefined;
    obj.bitcoin_proof = message.bitcoinProof ? Proof2.toAmino(message.bitcoinProof) : undefined;
    return obj;
  },
  fromAminoMsg(object: ProofAminoMsg): Proof {
    return Proof.fromAmino(object.value);
  },
  fromProtoMsg(message: ProofProtoMsg): Proof {
    return Proof.decode(message.value);
  },
  toProto(message: Proof): Uint8Array {
    return Proof.encode(message).finish();
  },
  toProtoMsg(message: Proof): ProofProtoMsg {
    return {
      typeUrl: "/common.Proof",
      value: Proof.encode(message).finish()
    };
  }
};