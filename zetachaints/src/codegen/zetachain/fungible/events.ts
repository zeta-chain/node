import { CoinType, coinTypeFromJSON } from "../common/common";
import { UpdatePausedStatusAction, updatePausedStatusActionFromJSON } from "./tx";
import { BinaryReader, BinaryWriter } from "../../binary";
export interface EventSystemContractUpdated {
  msgTypeUrl: string;
  newContractAddress: string;
  oldContractAddress: string;
  signer: string;
}
export interface EventSystemContractUpdatedProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.EventSystemContractUpdated";
  value: Uint8Array;
}
export interface EventSystemContractUpdatedAmino {
  msg_type_url?: string;
  new_contract_address?: string;
  old_contract_address?: string;
  signer?: string;
}
export interface EventSystemContractUpdatedAminoMsg {
  type: "/zetachain.zetacore.fungible.EventSystemContractUpdated";
  value: EventSystemContractUpdatedAmino;
}
export interface EventSystemContractUpdatedSDKType {
  msg_type_url: string;
  new_contract_address: string;
  old_contract_address: string;
  signer: string;
}
export interface EventZRC20Deployed {
  msgTypeUrl: string;
  chainId: bigint;
  contract: string;
  name: string;
  symbol: string;
  decimals: bigint;
  coinType: CoinType;
  erc20: string;
  gasLimit: bigint;
}
export interface EventZRC20DeployedProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.EventZRC20Deployed";
  value: Uint8Array;
}
export interface EventZRC20DeployedAmino {
  msg_type_url?: string;
  chain_id?: string;
  contract?: string;
  name?: string;
  symbol?: string;
  decimals?: string;
  coin_type?: CoinType;
  erc20?: string;
  gas_limit?: string;
}
export interface EventZRC20DeployedAminoMsg {
  type: "/zetachain.zetacore.fungible.EventZRC20Deployed";
  value: EventZRC20DeployedAmino;
}
export interface EventZRC20DeployedSDKType {
  msg_type_url: string;
  chain_id: bigint;
  contract: string;
  name: string;
  symbol: string;
  decimals: bigint;
  coin_type: CoinType;
  erc20: string;
  gas_limit: bigint;
}
export interface EventZRC20WithdrawFeeUpdated {
  msgTypeUrl: string;
  chainId: bigint;
  coinType: CoinType;
  zrc20Address: string;
  oldWithdrawFee: string;
  newWithdrawFee: string;
  signer: string;
  oldGasLimit: string;
  newGasLimit: string;
}
export interface EventZRC20WithdrawFeeUpdatedProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.EventZRC20WithdrawFeeUpdated";
  value: Uint8Array;
}
export interface EventZRC20WithdrawFeeUpdatedAmino {
  msg_type_url?: string;
  chain_id?: string;
  coin_type?: CoinType;
  zrc20_address?: string;
  old_withdraw_fee?: string;
  new_withdraw_fee?: string;
  signer?: string;
  old_gas_limit?: string;
  new_gas_limit?: string;
}
export interface EventZRC20WithdrawFeeUpdatedAminoMsg {
  type: "/zetachain.zetacore.fungible.EventZRC20WithdrawFeeUpdated";
  value: EventZRC20WithdrawFeeUpdatedAmino;
}
export interface EventZRC20WithdrawFeeUpdatedSDKType {
  msg_type_url: string;
  chain_id: bigint;
  coin_type: CoinType;
  zrc20_address: string;
  old_withdraw_fee: string;
  new_withdraw_fee: string;
  signer: string;
  old_gas_limit: string;
  new_gas_limit: string;
}
export interface EventZRC20PausedStatusUpdated {
  msgTypeUrl: string;
  zrc20Addresses: string[];
  action: UpdatePausedStatusAction;
  signer: string;
}
export interface EventZRC20PausedStatusUpdatedProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.EventZRC20PausedStatusUpdated";
  value: Uint8Array;
}
export interface EventZRC20PausedStatusUpdatedAmino {
  msg_type_url?: string;
  zrc20_addresses?: string[];
  action?: UpdatePausedStatusAction;
  signer?: string;
}
export interface EventZRC20PausedStatusUpdatedAminoMsg {
  type: "/zetachain.zetacore.fungible.EventZRC20PausedStatusUpdated";
  value: EventZRC20PausedStatusUpdatedAmino;
}
export interface EventZRC20PausedStatusUpdatedSDKType {
  msg_type_url: string;
  zrc20_addresses: string[];
  action: UpdatePausedStatusAction;
  signer: string;
}
export interface EventSystemContractsDeployed {
  msgTypeUrl: string;
  uniswapV2Factory: string;
  wzeta: string;
  uniswapV2Router: string;
  connectorZevm: string;
  systemContract: string;
  signer: string;
}
export interface EventSystemContractsDeployedProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.EventSystemContractsDeployed";
  value: Uint8Array;
}
export interface EventSystemContractsDeployedAmino {
  msg_type_url?: string;
  uniswap_v2_factory?: string;
  wzeta?: string;
  uniswap_v2_router?: string;
  connector_zevm?: string;
  system_contract?: string;
  signer?: string;
}
export interface EventSystemContractsDeployedAminoMsg {
  type: "/zetachain.zetacore.fungible.EventSystemContractsDeployed";
  value: EventSystemContractsDeployedAmino;
}
export interface EventSystemContractsDeployedSDKType {
  msg_type_url: string;
  uniswap_v2_factory: string;
  wzeta: string;
  uniswap_v2_router: string;
  connector_zevm: string;
  system_contract: string;
  signer: string;
}
export interface EventBytecodeUpdated {
  msgTypeUrl: string;
  contractAddress: string;
  newBytecodeHash: string;
  oldBytecodeHash: string;
  signer: string;
}
export interface EventBytecodeUpdatedProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.EventBytecodeUpdated";
  value: Uint8Array;
}
export interface EventBytecodeUpdatedAmino {
  msg_type_url?: string;
  contract_address?: string;
  new_bytecode_hash?: string;
  old_bytecode_hash?: string;
  signer?: string;
}
export interface EventBytecodeUpdatedAminoMsg {
  type: "/zetachain.zetacore.fungible.EventBytecodeUpdated";
  value: EventBytecodeUpdatedAmino;
}
export interface EventBytecodeUpdatedSDKType {
  msg_type_url: string;
  contract_address: string;
  new_bytecode_hash: string;
  old_bytecode_hash: string;
  signer: string;
}
function createBaseEventSystemContractUpdated(): EventSystemContractUpdated {
  return {
    msgTypeUrl: "",
    newContractAddress: "",
    oldContractAddress: "",
    signer: ""
  };
}
export const EventSystemContractUpdated = {
  typeUrl: "/zetachain.zetacore.fungible.EventSystemContractUpdated",
  encode(message: EventSystemContractUpdated, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.newContractAddress !== "") {
      writer.uint32(18).string(message.newContractAddress);
    }
    if (message.oldContractAddress !== "") {
      writer.uint32(26).string(message.oldContractAddress);
    }
    if (message.signer !== "") {
      writer.uint32(34).string(message.signer);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventSystemContractUpdated {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventSystemContractUpdated();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.newContractAddress = reader.string();
          break;
        case 3:
          message.oldContractAddress = reader.string();
          break;
        case 4:
          message.signer = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventSystemContractUpdated>): EventSystemContractUpdated {
    const message = createBaseEventSystemContractUpdated();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.newContractAddress = object.newContractAddress ?? "";
    message.oldContractAddress = object.oldContractAddress ?? "";
    message.signer = object.signer ?? "";
    return message;
  },
  fromAmino(object: EventSystemContractUpdatedAmino): EventSystemContractUpdated {
    const message = createBaseEventSystemContractUpdated();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.new_contract_address !== undefined && object.new_contract_address !== null) {
      message.newContractAddress = object.new_contract_address;
    }
    if (object.old_contract_address !== undefined && object.old_contract_address !== null) {
      message.oldContractAddress = object.old_contract_address;
    }
    if (object.signer !== undefined && object.signer !== null) {
      message.signer = object.signer;
    }
    return message;
  },
  toAmino(message: EventSystemContractUpdated): EventSystemContractUpdatedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.new_contract_address = message.newContractAddress;
    obj.old_contract_address = message.oldContractAddress;
    obj.signer = message.signer;
    return obj;
  },
  fromAminoMsg(object: EventSystemContractUpdatedAminoMsg): EventSystemContractUpdated {
    return EventSystemContractUpdated.fromAmino(object.value);
  },
  fromProtoMsg(message: EventSystemContractUpdatedProtoMsg): EventSystemContractUpdated {
    return EventSystemContractUpdated.decode(message.value);
  },
  toProto(message: EventSystemContractUpdated): Uint8Array {
    return EventSystemContractUpdated.encode(message).finish();
  },
  toProtoMsg(message: EventSystemContractUpdated): EventSystemContractUpdatedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.EventSystemContractUpdated",
      value: EventSystemContractUpdated.encode(message).finish()
    };
  }
};
function createBaseEventZRC20Deployed(): EventZRC20Deployed {
  return {
    msgTypeUrl: "",
    chainId: BigInt(0),
    contract: "",
    name: "",
    symbol: "",
    decimals: BigInt(0),
    coinType: 0,
    erc20: "",
    gasLimit: BigInt(0)
  };
}
export const EventZRC20Deployed = {
  typeUrl: "/zetachain.zetacore.fungible.EventZRC20Deployed",
  encode(message: EventZRC20Deployed, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.contract !== "") {
      writer.uint32(26).string(message.contract);
    }
    if (message.name !== "") {
      writer.uint32(34).string(message.name);
    }
    if (message.symbol !== "") {
      writer.uint32(42).string(message.symbol);
    }
    if (message.decimals !== BigInt(0)) {
      writer.uint32(48).int64(message.decimals);
    }
    if (message.coinType !== 0) {
      writer.uint32(56).int32(message.coinType);
    }
    if (message.erc20 !== "") {
      writer.uint32(66).string(message.erc20);
    }
    if (message.gasLimit !== BigInt(0)) {
      writer.uint32(72).int64(message.gasLimit);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventZRC20Deployed {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventZRC20Deployed();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.contract = reader.string();
          break;
        case 4:
          message.name = reader.string();
          break;
        case 5:
          message.symbol = reader.string();
          break;
        case 6:
          message.decimals = reader.int64();
          break;
        case 7:
          message.coinType = (reader.int32() as any);
          break;
        case 8:
          message.erc20 = reader.string();
          break;
        case 9:
          message.gasLimit = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventZRC20Deployed>): EventZRC20Deployed {
    const message = createBaseEventZRC20Deployed();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.contract = object.contract ?? "";
    message.name = object.name ?? "";
    message.symbol = object.symbol ?? "";
    message.decimals = object.decimals !== undefined && object.decimals !== null ? BigInt(object.decimals.toString()) : BigInt(0);
    message.coinType = object.coinType ?? 0;
    message.erc20 = object.erc20 ?? "";
    message.gasLimit = object.gasLimit !== undefined && object.gasLimit !== null ? BigInt(object.gasLimit.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: EventZRC20DeployedAmino): EventZRC20Deployed {
    const message = createBaseEventZRC20Deployed();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.contract !== undefined && object.contract !== null) {
      message.contract = object.contract;
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    }
    if (object.symbol !== undefined && object.symbol !== null) {
      message.symbol = object.symbol;
    }
    if (object.decimals !== undefined && object.decimals !== null) {
      message.decimals = BigInt(object.decimals);
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    if (object.erc20 !== undefined && object.erc20 !== null) {
      message.erc20 = object.erc20;
    }
    if (object.gas_limit !== undefined && object.gas_limit !== null) {
      message.gasLimit = BigInt(object.gas_limit);
    }
    return message;
  },
  toAmino(message: EventZRC20Deployed): EventZRC20DeployedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.contract = message.contract;
    obj.name = message.name;
    obj.symbol = message.symbol;
    obj.decimals = message.decimals ? message.decimals.toString() : undefined;
    obj.coin_type = message.coinType;
    obj.erc20 = message.erc20;
    obj.gas_limit = message.gasLimit ? message.gasLimit.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: EventZRC20DeployedAminoMsg): EventZRC20Deployed {
    return EventZRC20Deployed.fromAmino(object.value);
  },
  fromProtoMsg(message: EventZRC20DeployedProtoMsg): EventZRC20Deployed {
    return EventZRC20Deployed.decode(message.value);
  },
  toProto(message: EventZRC20Deployed): Uint8Array {
    return EventZRC20Deployed.encode(message).finish();
  },
  toProtoMsg(message: EventZRC20Deployed): EventZRC20DeployedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.EventZRC20Deployed",
      value: EventZRC20Deployed.encode(message).finish()
    };
  }
};
function createBaseEventZRC20WithdrawFeeUpdated(): EventZRC20WithdrawFeeUpdated {
  return {
    msgTypeUrl: "",
    chainId: BigInt(0),
    coinType: 0,
    zrc20Address: "",
    oldWithdrawFee: "",
    newWithdrawFee: "",
    signer: "",
    oldGasLimit: "",
    newGasLimit: ""
  };
}
export const EventZRC20WithdrawFeeUpdated = {
  typeUrl: "/zetachain.zetacore.fungible.EventZRC20WithdrawFeeUpdated",
  encode(message: EventZRC20WithdrawFeeUpdated, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.coinType !== 0) {
      writer.uint32(24).int32(message.coinType);
    }
    if (message.zrc20Address !== "") {
      writer.uint32(34).string(message.zrc20Address);
    }
    if (message.oldWithdrawFee !== "") {
      writer.uint32(42).string(message.oldWithdrawFee);
    }
    if (message.newWithdrawFee !== "") {
      writer.uint32(50).string(message.newWithdrawFee);
    }
    if (message.signer !== "") {
      writer.uint32(58).string(message.signer);
    }
    if (message.oldGasLimit !== "") {
      writer.uint32(66).string(message.oldGasLimit);
    }
    if (message.newGasLimit !== "") {
      writer.uint32(74).string(message.newGasLimit);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventZRC20WithdrawFeeUpdated {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventZRC20WithdrawFeeUpdated();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.coinType = (reader.int32() as any);
          break;
        case 4:
          message.zrc20Address = reader.string();
          break;
        case 5:
          message.oldWithdrawFee = reader.string();
          break;
        case 6:
          message.newWithdrawFee = reader.string();
          break;
        case 7:
          message.signer = reader.string();
          break;
        case 8:
          message.oldGasLimit = reader.string();
          break;
        case 9:
          message.newGasLimit = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventZRC20WithdrawFeeUpdated>): EventZRC20WithdrawFeeUpdated {
    const message = createBaseEventZRC20WithdrawFeeUpdated();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.coinType = object.coinType ?? 0;
    message.zrc20Address = object.zrc20Address ?? "";
    message.oldWithdrawFee = object.oldWithdrawFee ?? "";
    message.newWithdrawFee = object.newWithdrawFee ?? "";
    message.signer = object.signer ?? "";
    message.oldGasLimit = object.oldGasLimit ?? "";
    message.newGasLimit = object.newGasLimit ?? "";
    return message;
  },
  fromAmino(object: EventZRC20WithdrawFeeUpdatedAmino): EventZRC20WithdrawFeeUpdated {
    const message = createBaseEventZRC20WithdrawFeeUpdated();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    if (object.zrc20_address !== undefined && object.zrc20_address !== null) {
      message.zrc20Address = object.zrc20_address;
    }
    if (object.old_withdraw_fee !== undefined && object.old_withdraw_fee !== null) {
      message.oldWithdrawFee = object.old_withdraw_fee;
    }
    if (object.new_withdraw_fee !== undefined && object.new_withdraw_fee !== null) {
      message.newWithdrawFee = object.new_withdraw_fee;
    }
    if (object.signer !== undefined && object.signer !== null) {
      message.signer = object.signer;
    }
    if (object.old_gas_limit !== undefined && object.old_gas_limit !== null) {
      message.oldGasLimit = object.old_gas_limit;
    }
    if (object.new_gas_limit !== undefined && object.new_gas_limit !== null) {
      message.newGasLimit = object.new_gas_limit;
    }
    return message;
  },
  toAmino(message: EventZRC20WithdrawFeeUpdated): EventZRC20WithdrawFeeUpdatedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.coin_type = message.coinType;
    obj.zrc20_address = message.zrc20Address;
    obj.old_withdraw_fee = message.oldWithdrawFee;
    obj.new_withdraw_fee = message.newWithdrawFee;
    obj.signer = message.signer;
    obj.old_gas_limit = message.oldGasLimit;
    obj.new_gas_limit = message.newGasLimit;
    return obj;
  },
  fromAminoMsg(object: EventZRC20WithdrawFeeUpdatedAminoMsg): EventZRC20WithdrawFeeUpdated {
    return EventZRC20WithdrawFeeUpdated.fromAmino(object.value);
  },
  fromProtoMsg(message: EventZRC20WithdrawFeeUpdatedProtoMsg): EventZRC20WithdrawFeeUpdated {
    return EventZRC20WithdrawFeeUpdated.decode(message.value);
  },
  toProto(message: EventZRC20WithdrawFeeUpdated): Uint8Array {
    return EventZRC20WithdrawFeeUpdated.encode(message).finish();
  },
  toProtoMsg(message: EventZRC20WithdrawFeeUpdated): EventZRC20WithdrawFeeUpdatedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.EventZRC20WithdrawFeeUpdated",
      value: EventZRC20WithdrawFeeUpdated.encode(message).finish()
    };
  }
};
function createBaseEventZRC20PausedStatusUpdated(): EventZRC20PausedStatusUpdated {
  return {
    msgTypeUrl: "",
    zrc20Addresses: [],
    action: 0,
    signer: ""
  };
}
export const EventZRC20PausedStatusUpdated = {
  typeUrl: "/zetachain.zetacore.fungible.EventZRC20PausedStatusUpdated",
  encode(message: EventZRC20PausedStatusUpdated, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    for (const v of message.zrc20Addresses) {
      writer.uint32(18).string(v!);
    }
    if (message.action !== 0) {
      writer.uint32(24).int32(message.action);
    }
    if (message.signer !== "") {
      writer.uint32(34).string(message.signer);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventZRC20PausedStatusUpdated {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventZRC20PausedStatusUpdated();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.zrc20Addresses.push(reader.string());
          break;
        case 3:
          message.action = (reader.int32() as any);
          break;
        case 4:
          message.signer = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventZRC20PausedStatusUpdated>): EventZRC20PausedStatusUpdated {
    const message = createBaseEventZRC20PausedStatusUpdated();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.zrc20Addresses = object.zrc20Addresses?.map(e => e) || [];
    message.action = object.action ?? 0;
    message.signer = object.signer ?? "";
    return message;
  },
  fromAmino(object: EventZRC20PausedStatusUpdatedAmino): EventZRC20PausedStatusUpdated {
    const message = createBaseEventZRC20PausedStatusUpdated();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    message.zrc20Addresses = object.zrc20_addresses?.map(e => e) || [];
    if (object.action !== undefined && object.action !== null) {
      message.action = updatePausedStatusActionFromJSON(object.action);
    }
    if (object.signer !== undefined && object.signer !== null) {
      message.signer = object.signer;
    }
    return message;
  },
  toAmino(message: EventZRC20PausedStatusUpdated): EventZRC20PausedStatusUpdatedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    if (message.zrc20Addresses) {
      obj.zrc20_addresses = message.zrc20Addresses.map(e => e);
    } else {
      obj.zrc20_addresses = [];
    }
    obj.action = message.action;
    obj.signer = message.signer;
    return obj;
  },
  fromAminoMsg(object: EventZRC20PausedStatusUpdatedAminoMsg): EventZRC20PausedStatusUpdated {
    return EventZRC20PausedStatusUpdated.fromAmino(object.value);
  },
  fromProtoMsg(message: EventZRC20PausedStatusUpdatedProtoMsg): EventZRC20PausedStatusUpdated {
    return EventZRC20PausedStatusUpdated.decode(message.value);
  },
  toProto(message: EventZRC20PausedStatusUpdated): Uint8Array {
    return EventZRC20PausedStatusUpdated.encode(message).finish();
  },
  toProtoMsg(message: EventZRC20PausedStatusUpdated): EventZRC20PausedStatusUpdatedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.EventZRC20PausedStatusUpdated",
      value: EventZRC20PausedStatusUpdated.encode(message).finish()
    };
  }
};
function createBaseEventSystemContractsDeployed(): EventSystemContractsDeployed {
  return {
    msgTypeUrl: "",
    uniswapV2Factory: "",
    wzeta: "",
    uniswapV2Router: "",
    connectorZevm: "",
    systemContract: "",
    signer: ""
  };
}
export const EventSystemContractsDeployed = {
  typeUrl: "/zetachain.zetacore.fungible.EventSystemContractsDeployed",
  encode(message: EventSystemContractsDeployed, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.uniswapV2Factory !== "") {
      writer.uint32(18).string(message.uniswapV2Factory);
    }
    if (message.wzeta !== "") {
      writer.uint32(26).string(message.wzeta);
    }
    if (message.uniswapV2Router !== "") {
      writer.uint32(34).string(message.uniswapV2Router);
    }
    if (message.connectorZevm !== "") {
      writer.uint32(42).string(message.connectorZevm);
    }
    if (message.systemContract !== "") {
      writer.uint32(50).string(message.systemContract);
    }
    if (message.signer !== "") {
      writer.uint32(58).string(message.signer);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventSystemContractsDeployed {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventSystemContractsDeployed();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.uniswapV2Factory = reader.string();
          break;
        case 3:
          message.wzeta = reader.string();
          break;
        case 4:
          message.uniswapV2Router = reader.string();
          break;
        case 5:
          message.connectorZevm = reader.string();
          break;
        case 6:
          message.systemContract = reader.string();
          break;
        case 7:
          message.signer = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventSystemContractsDeployed>): EventSystemContractsDeployed {
    const message = createBaseEventSystemContractsDeployed();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.uniswapV2Factory = object.uniswapV2Factory ?? "";
    message.wzeta = object.wzeta ?? "";
    message.uniswapV2Router = object.uniswapV2Router ?? "";
    message.connectorZevm = object.connectorZevm ?? "";
    message.systemContract = object.systemContract ?? "";
    message.signer = object.signer ?? "";
    return message;
  },
  fromAmino(object: EventSystemContractsDeployedAmino): EventSystemContractsDeployed {
    const message = createBaseEventSystemContractsDeployed();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.uniswap_v2_factory !== undefined && object.uniswap_v2_factory !== null) {
      message.uniswapV2Factory = object.uniswap_v2_factory;
    }
    if (object.wzeta !== undefined && object.wzeta !== null) {
      message.wzeta = object.wzeta;
    }
    if (object.uniswap_v2_router !== undefined && object.uniswap_v2_router !== null) {
      message.uniswapV2Router = object.uniswap_v2_router;
    }
    if (object.connector_zevm !== undefined && object.connector_zevm !== null) {
      message.connectorZevm = object.connector_zevm;
    }
    if (object.system_contract !== undefined && object.system_contract !== null) {
      message.systemContract = object.system_contract;
    }
    if (object.signer !== undefined && object.signer !== null) {
      message.signer = object.signer;
    }
    return message;
  },
  toAmino(message: EventSystemContractsDeployed): EventSystemContractsDeployedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.uniswap_v2_factory = message.uniswapV2Factory;
    obj.wzeta = message.wzeta;
    obj.uniswap_v2_router = message.uniswapV2Router;
    obj.connector_zevm = message.connectorZevm;
    obj.system_contract = message.systemContract;
    obj.signer = message.signer;
    return obj;
  },
  fromAminoMsg(object: EventSystemContractsDeployedAminoMsg): EventSystemContractsDeployed {
    return EventSystemContractsDeployed.fromAmino(object.value);
  },
  fromProtoMsg(message: EventSystemContractsDeployedProtoMsg): EventSystemContractsDeployed {
    return EventSystemContractsDeployed.decode(message.value);
  },
  toProto(message: EventSystemContractsDeployed): Uint8Array {
    return EventSystemContractsDeployed.encode(message).finish();
  },
  toProtoMsg(message: EventSystemContractsDeployed): EventSystemContractsDeployedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.EventSystemContractsDeployed",
      value: EventSystemContractsDeployed.encode(message).finish()
    };
  }
};
function createBaseEventBytecodeUpdated(): EventBytecodeUpdated {
  return {
    msgTypeUrl: "",
    contractAddress: "",
    newBytecodeHash: "",
    oldBytecodeHash: "",
    signer: ""
  };
}
export const EventBytecodeUpdated = {
  typeUrl: "/zetachain.zetacore.fungible.EventBytecodeUpdated",
  encode(message: EventBytecodeUpdated, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.msgTypeUrl !== "") {
      writer.uint32(10).string(message.msgTypeUrl);
    }
    if (message.contractAddress !== "") {
      writer.uint32(18).string(message.contractAddress);
    }
    if (message.newBytecodeHash !== "") {
      writer.uint32(26).string(message.newBytecodeHash);
    }
    if (message.oldBytecodeHash !== "") {
      writer.uint32(34).string(message.oldBytecodeHash);
    }
    if (message.signer !== "") {
      writer.uint32(42).string(message.signer);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EventBytecodeUpdated {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventBytecodeUpdated();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.msgTypeUrl = reader.string();
          break;
        case 2:
          message.contractAddress = reader.string();
          break;
        case 3:
          message.newBytecodeHash = reader.string();
          break;
        case 4:
          message.oldBytecodeHash = reader.string();
          break;
        case 5:
          message.signer = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EventBytecodeUpdated>): EventBytecodeUpdated {
    const message = createBaseEventBytecodeUpdated();
    message.msgTypeUrl = object.msgTypeUrl ?? "";
    message.contractAddress = object.contractAddress ?? "";
    message.newBytecodeHash = object.newBytecodeHash ?? "";
    message.oldBytecodeHash = object.oldBytecodeHash ?? "";
    message.signer = object.signer ?? "";
    return message;
  },
  fromAmino(object: EventBytecodeUpdatedAmino): EventBytecodeUpdated {
    const message = createBaseEventBytecodeUpdated();
    if (object.msg_type_url !== undefined && object.msg_type_url !== null) {
      message.msgTypeUrl = object.msg_type_url;
    }
    if (object.contract_address !== undefined && object.contract_address !== null) {
      message.contractAddress = object.contract_address;
    }
    if (object.new_bytecode_hash !== undefined && object.new_bytecode_hash !== null) {
      message.newBytecodeHash = object.new_bytecode_hash;
    }
    if (object.old_bytecode_hash !== undefined && object.old_bytecode_hash !== null) {
      message.oldBytecodeHash = object.old_bytecode_hash;
    }
    if (object.signer !== undefined && object.signer !== null) {
      message.signer = object.signer;
    }
    return message;
  },
  toAmino(message: EventBytecodeUpdated): EventBytecodeUpdatedAmino {
    const obj: any = {};
    obj.msg_type_url = message.msgTypeUrl;
    obj.contract_address = message.contractAddress;
    obj.new_bytecode_hash = message.newBytecodeHash;
    obj.old_bytecode_hash = message.oldBytecodeHash;
    obj.signer = message.signer;
    return obj;
  },
  fromAminoMsg(object: EventBytecodeUpdatedAminoMsg): EventBytecodeUpdated {
    return EventBytecodeUpdated.fromAmino(object.value);
  },
  fromProtoMsg(message: EventBytecodeUpdatedProtoMsg): EventBytecodeUpdated {
    return EventBytecodeUpdated.decode(message.value);
  },
  toProto(message: EventBytecodeUpdated): Uint8Array {
    return EventBytecodeUpdated.encode(message).finish();
  },
  toProtoMsg(message: EventBytecodeUpdated): EventBytecodeUpdatedProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.EventBytecodeUpdated",
      value: EventBytecodeUpdated.encode(message).finish()
    };
  }
};