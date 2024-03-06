import { CoinType, coinTypeFromJSON } from "../common/common";
import { BinaryReader, BinaryWriter } from "../../binary";
export enum UpdatePausedStatusAction {
  PAUSE = 0,
  UNPAUSE = 1,
  UNRECOGNIZED = -1,
}
export const UpdatePausedStatusActionSDKType = UpdatePausedStatusAction;
export const UpdatePausedStatusActionAmino = UpdatePausedStatusAction;
export function updatePausedStatusActionFromJSON(object: any): UpdatePausedStatusAction {
  switch (object) {
    case 0:
    case "PAUSE":
      return UpdatePausedStatusAction.PAUSE;
    case 1:
    case "UNPAUSE":
      return UpdatePausedStatusAction.UNPAUSE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return UpdatePausedStatusAction.UNRECOGNIZED;
  }
}
export function updatePausedStatusActionToJSON(object: UpdatePausedStatusAction): string {
  switch (object) {
    case UpdatePausedStatusAction.PAUSE:
      return "PAUSE";
    case UpdatePausedStatusAction.UNPAUSE:
      return "UNPAUSE";
    case UpdatePausedStatusAction.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export interface MsgDeploySystemContracts {
  creator: string;
}
export interface MsgDeploySystemContractsProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgDeploySystemContracts";
  value: Uint8Array;
}
export interface MsgDeploySystemContractsAmino {
  creator?: string;
}
export interface MsgDeploySystemContractsAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgDeploySystemContracts";
  value: MsgDeploySystemContractsAmino;
}
export interface MsgDeploySystemContractsSDKType {
  creator: string;
}
export interface MsgDeploySystemContractsResponse {
  uniswapV2Factory: string;
  wzeta: string;
  uniswapV2Router: string;
  connectorZEVM: string;
  systemContract: string;
}
export interface MsgDeploySystemContractsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgDeploySystemContractsResponse";
  value: Uint8Array;
}
export interface MsgDeploySystemContractsResponseAmino {
  uniswapV2Factory?: string;
  wzeta?: string;
  uniswapV2Router?: string;
  connectorZEVM?: string;
  systemContract?: string;
}
export interface MsgDeploySystemContractsResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgDeploySystemContractsResponse";
  value: MsgDeploySystemContractsResponseAmino;
}
export interface MsgDeploySystemContractsResponseSDKType {
  uniswapV2Factory: string;
  wzeta: string;
  uniswapV2Router: string;
  connectorZEVM: string;
  systemContract: string;
}
export interface MsgUpdateZRC20WithdrawFee {
  creator: string;
  /** zrc20 address */
  zrc20Address: string;
  newWithdrawFee: string;
  newGasLimit: string;
}
export interface MsgUpdateZRC20WithdrawFeeProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee";
  value: Uint8Array;
}
export interface MsgUpdateZRC20WithdrawFeeAmino {
  creator?: string;
  /** zrc20 address */
  zrc20_address?: string;
  new_withdraw_fee?: string;
  new_gas_limit?: string;
}
export interface MsgUpdateZRC20WithdrawFeeAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee";
  value: MsgUpdateZRC20WithdrawFeeAmino;
}
export interface MsgUpdateZRC20WithdrawFeeSDKType {
  creator: string;
  zrc20_address: string;
  new_withdraw_fee: string;
  new_gas_limit: string;
}
export interface MsgUpdateZRC20WithdrawFeeResponse {}
export interface MsgUpdateZRC20WithdrawFeeResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFeeResponse";
  value: Uint8Array;
}
export interface MsgUpdateZRC20WithdrawFeeResponseAmino {}
export interface MsgUpdateZRC20WithdrawFeeResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFeeResponse";
  value: MsgUpdateZRC20WithdrawFeeResponseAmino;
}
export interface MsgUpdateZRC20WithdrawFeeResponseSDKType {}
export interface MsgUpdateSystemContract {
  creator: string;
  newSystemContractAddress: string;
}
export interface MsgUpdateSystemContractProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateSystemContract";
  value: Uint8Array;
}
export interface MsgUpdateSystemContractAmino {
  creator?: string;
  new_system_contract_address?: string;
}
export interface MsgUpdateSystemContractAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateSystemContract";
  value: MsgUpdateSystemContractAmino;
}
export interface MsgUpdateSystemContractSDKType {
  creator: string;
  new_system_contract_address: string;
}
export interface MsgUpdateSystemContractResponse {}
export interface MsgUpdateSystemContractResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateSystemContractResponse";
  value: Uint8Array;
}
export interface MsgUpdateSystemContractResponseAmino {}
export interface MsgUpdateSystemContractResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateSystemContractResponse";
  value: MsgUpdateSystemContractResponseAmino;
}
export interface MsgUpdateSystemContractResponseSDKType {}
export interface MsgDeployFungibleCoinZRC20 {
  creator: string;
  ERC20: string;
  foreignChainId: bigint;
  decimals: number;
  name: string;
  symbol: string;
  coinType: CoinType;
  gasLimit: bigint;
}
export interface MsgDeployFungibleCoinZRC20ProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20";
  value: Uint8Array;
}
export interface MsgDeployFungibleCoinZRC20Amino {
  creator?: string;
  ERC20?: string;
  foreign_chain_id?: string;
  decimals?: number;
  name?: string;
  symbol?: string;
  coin_type?: CoinType;
  gas_limit?: string;
}
export interface MsgDeployFungibleCoinZRC20AminoMsg {
  type: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20";
  value: MsgDeployFungibleCoinZRC20Amino;
}
export interface MsgDeployFungibleCoinZRC20SDKType {
  creator: string;
  ERC20: string;
  foreign_chain_id: bigint;
  decimals: number;
  name: string;
  symbol: string;
  coin_type: CoinType;
  gas_limit: bigint;
}
export interface MsgDeployFungibleCoinZRC20Response {
  address: string;
}
export interface MsgDeployFungibleCoinZRC20ResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20Response";
  value: Uint8Array;
}
export interface MsgDeployFungibleCoinZRC20ResponseAmino {
  address?: string;
}
export interface MsgDeployFungibleCoinZRC20ResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20Response";
  value: MsgDeployFungibleCoinZRC20ResponseAmino;
}
export interface MsgDeployFungibleCoinZRC20ResponseSDKType {
  address: string;
}
export interface MsgRemoveForeignCoin {
  creator: string;
  name: string;
}
export interface MsgRemoveForeignCoinProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgRemoveForeignCoin";
  value: Uint8Array;
}
export interface MsgRemoveForeignCoinAmino {
  creator?: string;
  name?: string;
}
export interface MsgRemoveForeignCoinAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgRemoveForeignCoin";
  value: MsgRemoveForeignCoinAmino;
}
export interface MsgRemoveForeignCoinSDKType {
  creator: string;
  name: string;
}
export interface MsgRemoveForeignCoinResponse {}
export interface MsgRemoveForeignCoinResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgRemoveForeignCoinResponse";
  value: Uint8Array;
}
export interface MsgRemoveForeignCoinResponseAmino {}
export interface MsgRemoveForeignCoinResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgRemoveForeignCoinResponse";
  value: MsgRemoveForeignCoinResponseAmino;
}
export interface MsgRemoveForeignCoinResponseSDKType {}
export interface MsgUpdateContractBytecode {
  creator: string;
  contractAddress: string;
  newCodeHash: string;
}
export interface MsgUpdateContractBytecodeProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateContractBytecode";
  value: Uint8Array;
}
export interface MsgUpdateContractBytecodeAmino {
  creator?: string;
  contract_address?: string;
  new_code_hash?: string;
}
export interface MsgUpdateContractBytecodeAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateContractBytecode";
  value: MsgUpdateContractBytecodeAmino;
}
export interface MsgUpdateContractBytecodeSDKType {
  creator: string;
  contract_address: string;
  new_code_hash: string;
}
export interface MsgUpdateContractBytecodeResponse {}
export interface MsgUpdateContractBytecodeResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateContractBytecodeResponse";
  value: Uint8Array;
}
export interface MsgUpdateContractBytecodeResponseAmino {}
export interface MsgUpdateContractBytecodeResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateContractBytecodeResponse";
  value: MsgUpdateContractBytecodeResponseAmino;
}
export interface MsgUpdateContractBytecodeResponseSDKType {}
export interface MsgUpdateZRC20PausedStatus {
  creator: string;
  zrc20Addresses: string[];
  action: UpdatePausedStatusAction;
}
export interface MsgUpdateZRC20PausedStatusProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus";
  value: Uint8Array;
}
export interface MsgUpdateZRC20PausedStatusAmino {
  creator?: string;
  zrc20_addresses?: string[];
  action?: UpdatePausedStatusAction;
}
export interface MsgUpdateZRC20PausedStatusAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus";
  value: MsgUpdateZRC20PausedStatusAmino;
}
export interface MsgUpdateZRC20PausedStatusSDKType {
  creator: string;
  zrc20_addresses: string[];
  action: UpdatePausedStatusAction;
}
export interface MsgUpdateZRC20PausedStatusResponse {}
export interface MsgUpdateZRC20PausedStatusResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatusResponse";
  value: Uint8Array;
}
export interface MsgUpdateZRC20PausedStatusResponseAmino {}
export interface MsgUpdateZRC20PausedStatusResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatusResponse";
  value: MsgUpdateZRC20PausedStatusResponseAmino;
}
export interface MsgUpdateZRC20PausedStatusResponseSDKType {}
export interface MsgUpdateZRC20LiquidityCap {
  creator: string;
  zrc20Address: string;
  liquidityCap: string;
}
export interface MsgUpdateZRC20LiquidityCapProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap";
  value: Uint8Array;
}
export interface MsgUpdateZRC20LiquidityCapAmino {
  creator?: string;
  zrc20_address?: string;
  liquidity_cap?: string;
}
export interface MsgUpdateZRC20LiquidityCapAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap";
  value: MsgUpdateZRC20LiquidityCapAmino;
}
export interface MsgUpdateZRC20LiquidityCapSDKType {
  creator: string;
  zrc20_address: string;
  liquidity_cap: string;
}
export interface MsgUpdateZRC20LiquidityCapResponse {}
export interface MsgUpdateZRC20LiquidityCapResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCapResponse";
  value: Uint8Array;
}
export interface MsgUpdateZRC20LiquidityCapResponseAmino {}
export interface MsgUpdateZRC20LiquidityCapResponseAminoMsg {
  type: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCapResponse";
  value: MsgUpdateZRC20LiquidityCapResponseAmino;
}
export interface MsgUpdateZRC20LiquidityCapResponseSDKType {}
function createBaseMsgDeploySystemContracts(): MsgDeploySystemContracts {
  return {
    creator: ""
  };
}
export const MsgDeploySystemContracts = {
  typeUrl: "/zetachain.zetacore.fungible.MsgDeploySystemContracts",
  encode(message: MsgDeploySystemContracts, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDeploySystemContracts {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeploySystemContracts();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgDeploySystemContracts>): MsgDeploySystemContracts {
    const message = createBaseMsgDeploySystemContracts();
    message.creator = object.creator ?? "";
    return message;
  },
  fromAmino(object: MsgDeploySystemContractsAmino): MsgDeploySystemContracts {
    const message = createBaseMsgDeploySystemContracts();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    return message;
  },
  toAmino(message: MsgDeploySystemContracts): MsgDeploySystemContractsAmino {
    const obj: any = {};
    obj.creator = message.creator;
    return obj;
  },
  fromAminoMsg(object: MsgDeploySystemContractsAminoMsg): MsgDeploySystemContracts {
    return MsgDeploySystemContracts.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDeploySystemContractsProtoMsg): MsgDeploySystemContracts {
    return MsgDeploySystemContracts.decode(message.value);
  },
  toProto(message: MsgDeploySystemContracts): Uint8Array {
    return MsgDeploySystemContracts.encode(message).finish();
  },
  toProtoMsg(message: MsgDeploySystemContracts): MsgDeploySystemContractsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgDeploySystemContracts",
      value: MsgDeploySystemContracts.encode(message).finish()
    };
  }
};
function createBaseMsgDeploySystemContractsResponse(): MsgDeploySystemContractsResponse {
  return {
    uniswapV2Factory: "",
    wzeta: "",
    uniswapV2Router: "",
    connectorZEVM: "",
    systemContract: ""
  };
}
export const MsgDeploySystemContractsResponse = {
  typeUrl: "/zetachain.zetacore.fungible.MsgDeploySystemContractsResponse",
  encode(message: MsgDeploySystemContractsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.uniswapV2Factory !== "") {
      writer.uint32(10).string(message.uniswapV2Factory);
    }
    if (message.wzeta !== "") {
      writer.uint32(18).string(message.wzeta);
    }
    if (message.uniswapV2Router !== "") {
      writer.uint32(26).string(message.uniswapV2Router);
    }
    if (message.connectorZEVM !== "") {
      writer.uint32(34).string(message.connectorZEVM);
    }
    if (message.systemContract !== "") {
      writer.uint32(42).string(message.systemContract);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDeploySystemContractsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeploySystemContractsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uniswapV2Factory = reader.string();
          break;
        case 2:
          message.wzeta = reader.string();
          break;
        case 3:
          message.uniswapV2Router = reader.string();
          break;
        case 4:
          message.connectorZEVM = reader.string();
          break;
        case 5:
          message.systemContract = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgDeploySystemContractsResponse>): MsgDeploySystemContractsResponse {
    const message = createBaseMsgDeploySystemContractsResponse();
    message.uniswapV2Factory = object.uniswapV2Factory ?? "";
    message.wzeta = object.wzeta ?? "";
    message.uniswapV2Router = object.uniswapV2Router ?? "";
    message.connectorZEVM = object.connectorZEVM ?? "";
    message.systemContract = object.systemContract ?? "";
    return message;
  },
  fromAmino(object: MsgDeploySystemContractsResponseAmino): MsgDeploySystemContractsResponse {
    const message = createBaseMsgDeploySystemContractsResponse();
    if (object.uniswapV2Factory !== undefined && object.uniswapV2Factory !== null) {
      message.uniswapV2Factory = object.uniswapV2Factory;
    }
    if (object.wzeta !== undefined && object.wzeta !== null) {
      message.wzeta = object.wzeta;
    }
    if (object.uniswapV2Router !== undefined && object.uniswapV2Router !== null) {
      message.uniswapV2Router = object.uniswapV2Router;
    }
    if (object.connectorZEVM !== undefined && object.connectorZEVM !== null) {
      message.connectorZEVM = object.connectorZEVM;
    }
    if (object.systemContract !== undefined && object.systemContract !== null) {
      message.systemContract = object.systemContract;
    }
    return message;
  },
  toAmino(message: MsgDeploySystemContractsResponse): MsgDeploySystemContractsResponseAmino {
    const obj: any = {};
    obj.uniswapV2Factory = message.uniswapV2Factory;
    obj.wzeta = message.wzeta;
    obj.uniswapV2Router = message.uniswapV2Router;
    obj.connectorZEVM = message.connectorZEVM;
    obj.systemContract = message.systemContract;
    return obj;
  },
  fromAminoMsg(object: MsgDeploySystemContractsResponseAminoMsg): MsgDeploySystemContractsResponse {
    return MsgDeploySystemContractsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDeploySystemContractsResponseProtoMsg): MsgDeploySystemContractsResponse {
    return MsgDeploySystemContractsResponse.decode(message.value);
  },
  toProto(message: MsgDeploySystemContractsResponse): Uint8Array {
    return MsgDeploySystemContractsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgDeploySystemContractsResponse): MsgDeploySystemContractsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgDeploySystemContractsResponse",
      value: MsgDeploySystemContractsResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateZRC20WithdrawFee(): MsgUpdateZRC20WithdrawFee {
  return {
    creator: "",
    zrc20Address: "",
    newWithdrawFee: "",
    newGasLimit: ""
  };
}
export const MsgUpdateZRC20WithdrawFee = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee",
  encode(message: MsgUpdateZRC20WithdrawFee, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.zrc20Address !== "") {
      writer.uint32(18).string(message.zrc20Address);
    }
    if (message.newWithdrawFee !== "") {
      writer.uint32(50).string(message.newWithdrawFee);
    }
    if (message.newGasLimit !== "") {
      writer.uint32(58).string(message.newGasLimit);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateZRC20WithdrawFee {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateZRC20WithdrawFee();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.zrc20Address = reader.string();
          break;
        case 6:
          message.newWithdrawFee = reader.string();
          break;
        case 7:
          message.newGasLimit = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateZRC20WithdrawFee>): MsgUpdateZRC20WithdrawFee {
    const message = createBaseMsgUpdateZRC20WithdrawFee();
    message.creator = object.creator ?? "";
    message.zrc20Address = object.zrc20Address ?? "";
    message.newWithdrawFee = object.newWithdrawFee ?? "";
    message.newGasLimit = object.newGasLimit ?? "";
    return message;
  },
  fromAmino(object: MsgUpdateZRC20WithdrawFeeAmino): MsgUpdateZRC20WithdrawFee {
    const message = createBaseMsgUpdateZRC20WithdrawFee();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.zrc20_address !== undefined && object.zrc20_address !== null) {
      message.zrc20Address = object.zrc20_address;
    }
    if (object.new_withdraw_fee !== undefined && object.new_withdraw_fee !== null) {
      message.newWithdrawFee = object.new_withdraw_fee;
    }
    if (object.new_gas_limit !== undefined && object.new_gas_limit !== null) {
      message.newGasLimit = object.new_gas_limit;
    }
    return message;
  },
  toAmino(message: MsgUpdateZRC20WithdrawFee): MsgUpdateZRC20WithdrawFeeAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.zrc20_address = message.zrc20Address;
    obj.new_withdraw_fee = message.newWithdrawFee;
    obj.new_gas_limit = message.newGasLimit;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateZRC20WithdrawFeeAminoMsg): MsgUpdateZRC20WithdrawFee {
    return MsgUpdateZRC20WithdrawFee.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateZRC20WithdrawFeeProtoMsg): MsgUpdateZRC20WithdrawFee {
    return MsgUpdateZRC20WithdrawFee.decode(message.value);
  },
  toProto(message: MsgUpdateZRC20WithdrawFee): Uint8Array {
    return MsgUpdateZRC20WithdrawFee.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateZRC20WithdrawFee): MsgUpdateZRC20WithdrawFeeProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee",
      value: MsgUpdateZRC20WithdrawFee.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateZRC20WithdrawFeeResponse(): MsgUpdateZRC20WithdrawFeeResponse {
  return {};
}
export const MsgUpdateZRC20WithdrawFeeResponse = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFeeResponse",
  encode(_: MsgUpdateZRC20WithdrawFeeResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateZRC20WithdrawFeeResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateZRC20WithdrawFeeResponse();
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
  fromPartial(_: Partial<MsgUpdateZRC20WithdrawFeeResponse>): MsgUpdateZRC20WithdrawFeeResponse {
    const message = createBaseMsgUpdateZRC20WithdrawFeeResponse();
    return message;
  },
  fromAmino(_: MsgUpdateZRC20WithdrawFeeResponseAmino): MsgUpdateZRC20WithdrawFeeResponse {
    const message = createBaseMsgUpdateZRC20WithdrawFeeResponse();
    return message;
  },
  toAmino(_: MsgUpdateZRC20WithdrawFeeResponse): MsgUpdateZRC20WithdrawFeeResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateZRC20WithdrawFeeResponseAminoMsg): MsgUpdateZRC20WithdrawFeeResponse {
    return MsgUpdateZRC20WithdrawFeeResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateZRC20WithdrawFeeResponseProtoMsg): MsgUpdateZRC20WithdrawFeeResponse {
    return MsgUpdateZRC20WithdrawFeeResponse.decode(message.value);
  },
  toProto(message: MsgUpdateZRC20WithdrawFeeResponse): Uint8Array {
    return MsgUpdateZRC20WithdrawFeeResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateZRC20WithdrawFeeResponse): MsgUpdateZRC20WithdrawFeeResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFeeResponse",
      value: MsgUpdateZRC20WithdrawFeeResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateSystemContract(): MsgUpdateSystemContract {
  return {
    creator: "",
    newSystemContractAddress: ""
  };
}
export const MsgUpdateSystemContract = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateSystemContract",
  encode(message: MsgUpdateSystemContract, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.newSystemContractAddress !== "") {
      writer.uint32(18).string(message.newSystemContractAddress);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateSystemContract {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateSystemContract();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.newSystemContractAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateSystemContract>): MsgUpdateSystemContract {
    const message = createBaseMsgUpdateSystemContract();
    message.creator = object.creator ?? "";
    message.newSystemContractAddress = object.newSystemContractAddress ?? "";
    return message;
  },
  fromAmino(object: MsgUpdateSystemContractAmino): MsgUpdateSystemContract {
    const message = createBaseMsgUpdateSystemContract();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.new_system_contract_address !== undefined && object.new_system_contract_address !== null) {
      message.newSystemContractAddress = object.new_system_contract_address;
    }
    return message;
  },
  toAmino(message: MsgUpdateSystemContract): MsgUpdateSystemContractAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.new_system_contract_address = message.newSystemContractAddress;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateSystemContractAminoMsg): MsgUpdateSystemContract {
    return MsgUpdateSystemContract.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateSystemContractProtoMsg): MsgUpdateSystemContract {
    return MsgUpdateSystemContract.decode(message.value);
  },
  toProto(message: MsgUpdateSystemContract): Uint8Array {
    return MsgUpdateSystemContract.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateSystemContract): MsgUpdateSystemContractProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateSystemContract",
      value: MsgUpdateSystemContract.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateSystemContractResponse(): MsgUpdateSystemContractResponse {
  return {};
}
export const MsgUpdateSystemContractResponse = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateSystemContractResponse",
  encode(_: MsgUpdateSystemContractResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateSystemContractResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateSystemContractResponse();
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
  fromPartial(_: Partial<MsgUpdateSystemContractResponse>): MsgUpdateSystemContractResponse {
    const message = createBaseMsgUpdateSystemContractResponse();
    return message;
  },
  fromAmino(_: MsgUpdateSystemContractResponseAmino): MsgUpdateSystemContractResponse {
    const message = createBaseMsgUpdateSystemContractResponse();
    return message;
  },
  toAmino(_: MsgUpdateSystemContractResponse): MsgUpdateSystemContractResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateSystemContractResponseAminoMsg): MsgUpdateSystemContractResponse {
    return MsgUpdateSystemContractResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateSystemContractResponseProtoMsg): MsgUpdateSystemContractResponse {
    return MsgUpdateSystemContractResponse.decode(message.value);
  },
  toProto(message: MsgUpdateSystemContractResponse): Uint8Array {
    return MsgUpdateSystemContractResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateSystemContractResponse): MsgUpdateSystemContractResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateSystemContractResponse",
      value: MsgUpdateSystemContractResponse.encode(message).finish()
    };
  }
};
function createBaseMsgDeployFungibleCoinZRC20(): MsgDeployFungibleCoinZRC20 {
  return {
    creator: "",
    ERC20: "",
    foreignChainId: BigInt(0),
    decimals: 0,
    name: "",
    symbol: "",
    coinType: 0,
    gasLimit: BigInt(0)
  };
}
export const MsgDeployFungibleCoinZRC20 = {
  typeUrl: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20",
  encode(message: MsgDeployFungibleCoinZRC20, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.ERC20 !== "") {
      writer.uint32(18).string(message.ERC20);
    }
    if (message.foreignChainId !== BigInt(0)) {
      writer.uint32(24).int64(message.foreignChainId);
    }
    if (message.decimals !== 0) {
      writer.uint32(32).uint32(message.decimals);
    }
    if (message.name !== "") {
      writer.uint32(42).string(message.name);
    }
    if (message.symbol !== "") {
      writer.uint32(50).string(message.symbol);
    }
    if (message.coinType !== 0) {
      writer.uint32(56).int32(message.coinType);
    }
    if (message.gasLimit !== BigInt(0)) {
      writer.uint32(64).int64(message.gasLimit);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDeployFungibleCoinZRC20 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeployFungibleCoinZRC20();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.ERC20 = reader.string();
          break;
        case 3:
          message.foreignChainId = reader.int64();
          break;
        case 4:
          message.decimals = reader.uint32();
          break;
        case 5:
          message.name = reader.string();
          break;
        case 6:
          message.symbol = reader.string();
          break;
        case 7:
          message.coinType = (reader.int32() as any);
          break;
        case 8:
          message.gasLimit = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgDeployFungibleCoinZRC20>): MsgDeployFungibleCoinZRC20 {
    const message = createBaseMsgDeployFungibleCoinZRC20();
    message.creator = object.creator ?? "";
    message.ERC20 = object.ERC20 ?? "";
    message.foreignChainId = object.foreignChainId !== undefined && object.foreignChainId !== null ? BigInt(object.foreignChainId.toString()) : BigInt(0);
    message.decimals = object.decimals ?? 0;
    message.name = object.name ?? "";
    message.symbol = object.symbol ?? "";
    message.coinType = object.coinType ?? 0;
    message.gasLimit = object.gasLimit !== undefined && object.gasLimit !== null ? BigInt(object.gasLimit.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgDeployFungibleCoinZRC20Amino): MsgDeployFungibleCoinZRC20 {
    const message = createBaseMsgDeployFungibleCoinZRC20();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.ERC20 !== undefined && object.ERC20 !== null) {
      message.ERC20 = object.ERC20;
    }
    if (object.foreign_chain_id !== undefined && object.foreign_chain_id !== null) {
      message.foreignChainId = BigInt(object.foreign_chain_id);
    }
    if (object.decimals !== undefined && object.decimals !== null) {
      message.decimals = object.decimals;
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    }
    if (object.symbol !== undefined && object.symbol !== null) {
      message.symbol = object.symbol;
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    if (object.gas_limit !== undefined && object.gas_limit !== null) {
      message.gasLimit = BigInt(object.gas_limit);
    }
    return message;
  },
  toAmino(message: MsgDeployFungibleCoinZRC20): MsgDeployFungibleCoinZRC20Amino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.ERC20 = message.ERC20;
    obj.foreign_chain_id = message.foreignChainId ? message.foreignChainId.toString() : undefined;
    obj.decimals = message.decimals;
    obj.name = message.name;
    obj.symbol = message.symbol;
    obj.coin_type = message.coinType;
    obj.gas_limit = message.gasLimit ? message.gasLimit.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgDeployFungibleCoinZRC20AminoMsg): MsgDeployFungibleCoinZRC20 {
    return MsgDeployFungibleCoinZRC20.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDeployFungibleCoinZRC20ProtoMsg): MsgDeployFungibleCoinZRC20 {
    return MsgDeployFungibleCoinZRC20.decode(message.value);
  },
  toProto(message: MsgDeployFungibleCoinZRC20): Uint8Array {
    return MsgDeployFungibleCoinZRC20.encode(message).finish();
  },
  toProtoMsg(message: MsgDeployFungibleCoinZRC20): MsgDeployFungibleCoinZRC20ProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20",
      value: MsgDeployFungibleCoinZRC20.encode(message).finish()
    };
  }
};
function createBaseMsgDeployFungibleCoinZRC20Response(): MsgDeployFungibleCoinZRC20Response {
  return {
    address: ""
  };
}
export const MsgDeployFungibleCoinZRC20Response = {
  typeUrl: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20Response",
  encode(message: MsgDeployFungibleCoinZRC20Response, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDeployFungibleCoinZRC20Response {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeployFungibleCoinZRC20Response();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgDeployFungibleCoinZRC20Response>): MsgDeployFungibleCoinZRC20Response {
    const message = createBaseMsgDeployFungibleCoinZRC20Response();
    message.address = object.address ?? "";
    return message;
  },
  fromAmino(object: MsgDeployFungibleCoinZRC20ResponseAmino): MsgDeployFungibleCoinZRC20Response {
    const message = createBaseMsgDeployFungibleCoinZRC20Response();
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    return message;
  },
  toAmino(message: MsgDeployFungibleCoinZRC20Response): MsgDeployFungibleCoinZRC20ResponseAmino {
    const obj: any = {};
    obj.address = message.address;
    return obj;
  },
  fromAminoMsg(object: MsgDeployFungibleCoinZRC20ResponseAminoMsg): MsgDeployFungibleCoinZRC20Response {
    return MsgDeployFungibleCoinZRC20Response.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDeployFungibleCoinZRC20ResponseProtoMsg): MsgDeployFungibleCoinZRC20Response {
    return MsgDeployFungibleCoinZRC20Response.decode(message.value);
  },
  toProto(message: MsgDeployFungibleCoinZRC20Response): Uint8Array {
    return MsgDeployFungibleCoinZRC20Response.encode(message).finish();
  },
  toProtoMsg(message: MsgDeployFungibleCoinZRC20Response): MsgDeployFungibleCoinZRC20ResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20Response",
      value: MsgDeployFungibleCoinZRC20Response.encode(message).finish()
    };
  }
};
function createBaseMsgRemoveForeignCoin(): MsgRemoveForeignCoin {
  return {
    creator: "",
    name: ""
  };
}
export const MsgRemoveForeignCoin = {
  typeUrl: "/zetachain.zetacore.fungible.MsgRemoveForeignCoin",
  encode(message: MsgRemoveForeignCoin, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgRemoveForeignCoin {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveForeignCoin();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgRemoveForeignCoin>): MsgRemoveForeignCoin {
    const message = createBaseMsgRemoveForeignCoin();
    message.creator = object.creator ?? "";
    message.name = object.name ?? "";
    return message;
  },
  fromAmino(object: MsgRemoveForeignCoinAmino): MsgRemoveForeignCoin {
    const message = createBaseMsgRemoveForeignCoin();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    }
    return message;
  },
  toAmino(message: MsgRemoveForeignCoin): MsgRemoveForeignCoinAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.name = message.name;
    return obj;
  },
  fromAminoMsg(object: MsgRemoveForeignCoinAminoMsg): MsgRemoveForeignCoin {
    return MsgRemoveForeignCoin.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRemoveForeignCoinProtoMsg): MsgRemoveForeignCoin {
    return MsgRemoveForeignCoin.decode(message.value);
  },
  toProto(message: MsgRemoveForeignCoin): Uint8Array {
    return MsgRemoveForeignCoin.encode(message).finish();
  },
  toProtoMsg(message: MsgRemoveForeignCoin): MsgRemoveForeignCoinProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgRemoveForeignCoin",
      value: MsgRemoveForeignCoin.encode(message).finish()
    };
  }
};
function createBaseMsgRemoveForeignCoinResponse(): MsgRemoveForeignCoinResponse {
  return {};
}
export const MsgRemoveForeignCoinResponse = {
  typeUrl: "/zetachain.zetacore.fungible.MsgRemoveForeignCoinResponse",
  encode(_: MsgRemoveForeignCoinResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgRemoveForeignCoinResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveForeignCoinResponse();
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
  fromPartial(_: Partial<MsgRemoveForeignCoinResponse>): MsgRemoveForeignCoinResponse {
    const message = createBaseMsgRemoveForeignCoinResponse();
    return message;
  },
  fromAmino(_: MsgRemoveForeignCoinResponseAmino): MsgRemoveForeignCoinResponse {
    const message = createBaseMsgRemoveForeignCoinResponse();
    return message;
  },
  toAmino(_: MsgRemoveForeignCoinResponse): MsgRemoveForeignCoinResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgRemoveForeignCoinResponseAminoMsg): MsgRemoveForeignCoinResponse {
    return MsgRemoveForeignCoinResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRemoveForeignCoinResponseProtoMsg): MsgRemoveForeignCoinResponse {
    return MsgRemoveForeignCoinResponse.decode(message.value);
  },
  toProto(message: MsgRemoveForeignCoinResponse): Uint8Array {
    return MsgRemoveForeignCoinResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgRemoveForeignCoinResponse): MsgRemoveForeignCoinResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgRemoveForeignCoinResponse",
      value: MsgRemoveForeignCoinResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateContractBytecode(): MsgUpdateContractBytecode {
  return {
    creator: "",
    contractAddress: "",
    newCodeHash: ""
  };
}
export const MsgUpdateContractBytecode = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateContractBytecode",
  encode(message: MsgUpdateContractBytecode, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.contractAddress !== "") {
      writer.uint32(18).string(message.contractAddress);
    }
    if (message.newCodeHash !== "") {
      writer.uint32(26).string(message.newCodeHash);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateContractBytecode {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateContractBytecode();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.contractAddress = reader.string();
          break;
        case 3:
          message.newCodeHash = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateContractBytecode>): MsgUpdateContractBytecode {
    const message = createBaseMsgUpdateContractBytecode();
    message.creator = object.creator ?? "";
    message.contractAddress = object.contractAddress ?? "";
    message.newCodeHash = object.newCodeHash ?? "";
    return message;
  },
  fromAmino(object: MsgUpdateContractBytecodeAmino): MsgUpdateContractBytecode {
    const message = createBaseMsgUpdateContractBytecode();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.contract_address !== undefined && object.contract_address !== null) {
      message.contractAddress = object.contract_address;
    }
    if (object.new_code_hash !== undefined && object.new_code_hash !== null) {
      message.newCodeHash = object.new_code_hash;
    }
    return message;
  },
  toAmino(message: MsgUpdateContractBytecode): MsgUpdateContractBytecodeAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.contract_address = message.contractAddress;
    obj.new_code_hash = message.newCodeHash;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateContractBytecodeAminoMsg): MsgUpdateContractBytecode {
    return MsgUpdateContractBytecode.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateContractBytecodeProtoMsg): MsgUpdateContractBytecode {
    return MsgUpdateContractBytecode.decode(message.value);
  },
  toProto(message: MsgUpdateContractBytecode): Uint8Array {
    return MsgUpdateContractBytecode.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateContractBytecode): MsgUpdateContractBytecodeProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateContractBytecode",
      value: MsgUpdateContractBytecode.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateContractBytecodeResponse(): MsgUpdateContractBytecodeResponse {
  return {};
}
export const MsgUpdateContractBytecodeResponse = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateContractBytecodeResponse",
  encode(_: MsgUpdateContractBytecodeResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateContractBytecodeResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateContractBytecodeResponse();
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
  fromPartial(_: Partial<MsgUpdateContractBytecodeResponse>): MsgUpdateContractBytecodeResponse {
    const message = createBaseMsgUpdateContractBytecodeResponse();
    return message;
  },
  fromAmino(_: MsgUpdateContractBytecodeResponseAmino): MsgUpdateContractBytecodeResponse {
    const message = createBaseMsgUpdateContractBytecodeResponse();
    return message;
  },
  toAmino(_: MsgUpdateContractBytecodeResponse): MsgUpdateContractBytecodeResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateContractBytecodeResponseAminoMsg): MsgUpdateContractBytecodeResponse {
    return MsgUpdateContractBytecodeResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateContractBytecodeResponseProtoMsg): MsgUpdateContractBytecodeResponse {
    return MsgUpdateContractBytecodeResponse.decode(message.value);
  },
  toProto(message: MsgUpdateContractBytecodeResponse): Uint8Array {
    return MsgUpdateContractBytecodeResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateContractBytecodeResponse): MsgUpdateContractBytecodeResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateContractBytecodeResponse",
      value: MsgUpdateContractBytecodeResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateZRC20PausedStatus(): MsgUpdateZRC20PausedStatus {
  return {
    creator: "",
    zrc20Addresses: [],
    action: 0
  };
}
export const MsgUpdateZRC20PausedStatus = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus",
  encode(message: MsgUpdateZRC20PausedStatus, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    for (const v of message.zrc20Addresses) {
      writer.uint32(18).string(v!);
    }
    if (message.action !== 0) {
      writer.uint32(24).int32(message.action);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateZRC20PausedStatus {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateZRC20PausedStatus();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.zrc20Addresses.push(reader.string());
          break;
        case 3:
          message.action = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateZRC20PausedStatus>): MsgUpdateZRC20PausedStatus {
    const message = createBaseMsgUpdateZRC20PausedStatus();
    message.creator = object.creator ?? "";
    message.zrc20Addresses = object.zrc20Addresses?.map(e => e) || [];
    message.action = object.action ?? 0;
    return message;
  },
  fromAmino(object: MsgUpdateZRC20PausedStatusAmino): MsgUpdateZRC20PausedStatus {
    const message = createBaseMsgUpdateZRC20PausedStatus();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    message.zrc20Addresses = object.zrc20_addresses?.map(e => e) || [];
    if (object.action !== undefined && object.action !== null) {
      message.action = updatePausedStatusActionFromJSON(object.action);
    }
    return message;
  },
  toAmino(message: MsgUpdateZRC20PausedStatus): MsgUpdateZRC20PausedStatusAmino {
    const obj: any = {};
    obj.creator = message.creator;
    if (message.zrc20Addresses) {
      obj.zrc20_addresses = message.zrc20Addresses.map(e => e);
    } else {
      obj.zrc20_addresses = [];
    }
    obj.action = message.action;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateZRC20PausedStatusAminoMsg): MsgUpdateZRC20PausedStatus {
    return MsgUpdateZRC20PausedStatus.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateZRC20PausedStatusProtoMsg): MsgUpdateZRC20PausedStatus {
    return MsgUpdateZRC20PausedStatus.decode(message.value);
  },
  toProto(message: MsgUpdateZRC20PausedStatus): Uint8Array {
    return MsgUpdateZRC20PausedStatus.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateZRC20PausedStatus): MsgUpdateZRC20PausedStatusProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus",
      value: MsgUpdateZRC20PausedStatus.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateZRC20PausedStatusResponse(): MsgUpdateZRC20PausedStatusResponse {
  return {};
}
export const MsgUpdateZRC20PausedStatusResponse = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatusResponse",
  encode(_: MsgUpdateZRC20PausedStatusResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateZRC20PausedStatusResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateZRC20PausedStatusResponse();
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
  fromPartial(_: Partial<MsgUpdateZRC20PausedStatusResponse>): MsgUpdateZRC20PausedStatusResponse {
    const message = createBaseMsgUpdateZRC20PausedStatusResponse();
    return message;
  },
  fromAmino(_: MsgUpdateZRC20PausedStatusResponseAmino): MsgUpdateZRC20PausedStatusResponse {
    const message = createBaseMsgUpdateZRC20PausedStatusResponse();
    return message;
  },
  toAmino(_: MsgUpdateZRC20PausedStatusResponse): MsgUpdateZRC20PausedStatusResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateZRC20PausedStatusResponseAminoMsg): MsgUpdateZRC20PausedStatusResponse {
    return MsgUpdateZRC20PausedStatusResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateZRC20PausedStatusResponseProtoMsg): MsgUpdateZRC20PausedStatusResponse {
    return MsgUpdateZRC20PausedStatusResponse.decode(message.value);
  },
  toProto(message: MsgUpdateZRC20PausedStatusResponse): Uint8Array {
    return MsgUpdateZRC20PausedStatusResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateZRC20PausedStatusResponse): MsgUpdateZRC20PausedStatusResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatusResponse",
      value: MsgUpdateZRC20PausedStatusResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateZRC20LiquidityCap(): MsgUpdateZRC20LiquidityCap {
  return {
    creator: "",
    zrc20Address: "",
    liquidityCap: ""
  };
}
export const MsgUpdateZRC20LiquidityCap = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap",
  encode(message: MsgUpdateZRC20LiquidityCap, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.zrc20Address !== "") {
      writer.uint32(18).string(message.zrc20Address);
    }
    if (message.liquidityCap !== "") {
      writer.uint32(26).string(message.liquidityCap);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateZRC20LiquidityCap {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateZRC20LiquidityCap();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.zrc20Address = reader.string();
          break;
        case 3:
          message.liquidityCap = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateZRC20LiquidityCap>): MsgUpdateZRC20LiquidityCap {
    const message = createBaseMsgUpdateZRC20LiquidityCap();
    message.creator = object.creator ?? "";
    message.zrc20Address = object.zrc20Address ?? "";
    message.liquidityCap = object.liquidityCap ?? "";
    return message;
  },
  fromAmino(object: MsgUpdateZRC20LiquidityCapAmino): MsgUpdateZRC20LiquidityCap {
    const message = createBaseMsgUpdateZRC20LiquidityCap();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.zrc20_address !== undefined && object.zrc20_address !== null) {
      message.zrc20Address = object.zrc20_address;
    }
    if (object.liquidity_cap !== undefined && object.liquidity_cap !== null) {
      message.liquidityCap = object.liquidity_cap;
    }
    return message;
  },
  toAmino(message: MsgUpdateZRC20LiquidityCap): MsgUpdateZRC20LiquidityCapAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.zrc20_address = message.zrc20Address;
    obj.liquidity_cap = message.liquidityCap;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateZRC20LiquidityCapAminoMsg): MsgUpdateZRC20LiquidityCap {
    return MsgUpdateZRC20LiquidityCap.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateZRC20LiquidityCapProtoMsg): MsgUpdateZRC20LiquidityCap {
    return MsgUpdateZRC20LiquidityCap.decode(message.value);
  },
  toProto(message: MsgUpdateZRC20LiquidityCap): Uint8Array {
    return MsgUpdateZRC20LiquidityCap.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateZRC20LiquidityCap): MsgUpdateZRC20LiquidityCapProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap",
      value: MsgUpdateZRC20LiquidityCap.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateZRC20LiquidityCapResponse(): MsgUpdateZRC20LiquidityCapResponse {
  return {};
}
export const MsgUpdateZRC20LiquidityCapResponse = {
  typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCapResponse",
  encode(_: MsgUpdateZRC20LiquidityCapResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateZRC20LiquidityCapResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateZRC20LiquidityCapResponse();
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
  fromPartial(_: Partial<MsgUpdateZRC20LiquidityCapResponse>): MsgUpdateZRC20LiquidityCapResponse {
    const message = createBaseMsgUpdateZRC20LiquidityCapResponse();
    return message;
  },
  fromAmino(_: MsgUpdateZRC20LiquidityCapResponseAmino): MsgUpdateZRC20LiquidityCapResponse {
    const message = createBaseMsgUpdateZRC20LiquidityCapResponse();
    return message;
  },
  toAmino(_: MsgUpdateZRC20LiquidityCapResponse): MsgUpdateZRC20LiquidityCapResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateZRC20LiquidityCapResponseAminoMsg): MsgUpdateZRC20LiquidityCapResponse {
    return MsgUpdateZRC20LiquidityCapResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateZRC20LiquidityCapResponseProtoMsg): MsgUpdateZRC20LiquidityCapResponse {
    return MsgUpdateZRC20LiquidityCapResponse.decode(message.value);
  },
  toProto(message: MsgUpdateZRC20LiquidityCapResponse): Uint8Array {
    return MsgUpdateZRC20LiquidityCapResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateZRC20LiquidityCapResponse): MsgUpdateZRC20LiquidityCapResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCapResponse",
      value: MsgUpdateZRC20LiquidityCapResponse.encode(message).finish()
    };
  }
};