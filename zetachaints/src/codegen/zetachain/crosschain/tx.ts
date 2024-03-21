import { Proof, ProofAmino, ProofSDKType, CoinType, ReceiveStatus, receiveStatusFromJSON, coinTypeFromJSON } from "../common/common";
import { BinaryReader, BinaryWriter } from "../../binary";
export interface MsgCreateTSSVoter {
  creator: string;
  tssPubkey: string;
  keyGenZetaHeight: bigint;
  status: ReceiveStatus;
}
export interface MsgCreateTSSVoterProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgCreateTSSVoter";
  value: Uint8Array;
}
export interface MsgCreateTSSVoterAmino {
  creator?: string;
  tss_pubkey?: string;
  keyGenZetaHeight?: string;
  status?: ReceiveStatus;
}
export interface MsgCreateTSSVoterAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgCreateTSSVoter";
  value: MsgCreateTSSVoterAmino;
}
export interface MsgCreateTSSVoterSDKType {
  creator: string;
  tss_pubkey: string;
  keyGenZetaHeight: bigint;
  status: ReceiveStatus;
}
export interface MsgCreateTSSVoterResponse {}
export interface MsgCreateTSSVoterResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgCreateTSSVoterResponse";
  value: Uint8Array;
}
export interface MsgCreateTSSVoterResponseAmino {}
export interface MsgCreateTSSVoterResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgCreateTSSVoterResponse";
  value: MsgCreateTSSVoterResponseAmino;
}
export interface MsgCreateTSSVoterResponseSDKType {}
export interface MsgMigrateTssFunds {
  creator: string;
  chainId: bigint;
  amount: string;
}
export interface MsgMigrateTssFundsProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgMigrateTssFunds";
  value: Uint8Array;
}
export interface MsgMigrateTssFundsAmino {
  creator?: string;
  chain_id?: string;
  amount?: string;
}
export interface MsgMigrateTssFundsAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgMigrateTssFunds";
  value: MsgMigrateTssFundsAmino;
}
export interface MsgMigrateTssFundsSDKType {
  creator: string;
  chain_id: bigint;
  amount: string;
}
export interface MsgMigrateTssFundsResponse {}
export interface MsgMigrateTssFundsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgMigrateTssFundsResponse";
  value: Uint8Array;
}
export interface MsgMigrateTssFundsResponseAmino {}
export interface MsgMigrateTssFundsResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgMigrateTssFundsResponse";
  value: MsgMigrateTssFundsResponseAmino;
}
export interface MsgMigrateTssFundsResponseSDKType {}
export interface MsgUpdateTssAddress {
  creator: string;
  tssPubkey: string;
}
export interface MsgUpdateTssAddressProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgUpdateTssAddress";
  value: Uint8Array;
}
export interface MsgUpdateTssAddressAmino {
  creator?: string;
  tss_pubkey?: string;
}
export interface MsgUpdateTssAddressAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgUpdateTssAddress";
  value: MsgUpdateTssAddressAmino;
}
export interface MsgUpdateTssAddressSDKType {
  creator: string;
  tss_pubkey: string;
}
export interface MsgUpdateTssAddressResponse {}
export interface MsgUpdateTssAddressResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgUpdateTssAddressResponse";
  value: Uint8Array;
}
export interface MsgUpdateTssAddressResponseAmino {}
export interface MsgUpdateTssAddressResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgUpdateTssAddressResponse";
  value: MsgUpdateTssAddressResponseAmino;
}
export interface MsgUpdateTssAddressResponseSDKType {}
export interface MsgAddToInTxTracker {
  creator: string;
  chainId: bigint;
  txHash: string;
  coinType: CoinType;
  proof?: Proof;
  blockHash: string;
  txIndex: bigint;
}
export interface MsgAddToInTxTrackerProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAddToInTxTracker";
  value: Uint8Array;
}
export interface MsgAddToInTxTrackerAmino {
  creator?: string;
  chain_id?: string;
  tx_hash?: string;
  coin_type?: CoinType;
  proof?: ProofAmino;
  block_hash?: string;
  tx_index?: string;
}
export interface MsgAddToInTxTrackerAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgAddToInTxTracker";
  value: MsgAddToInTxTrackerAmino;
}
export interface MsgAddToInTxTrackerSDKType {
  creator: string;
  chain_id: bigint;
  tx_hash: string;
  coin_type: CoinType;
  proof?: ProofSDKType;
  block_hash: string;
  tx_index: bigint;
}
export interface MsgAddToInTxTrackerResponse {}
export interface MsgAddToInTxTrackerResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAddToInTxTrackerResponse";
  value: Uint8Array;
}
export interface MsgAddToInTxTrackerResponseAmino {}
export interface MsgAddToInTxTrackerResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgAddToInTxTrackerResponse";
  value: MsgAddToInTxTrackerResponseAmino;
}
export interface MsgAddToInTxTrackerResponseSDKType {}
export interface MsgWhitelistERC20 {
  creator: string;
  erc20Address: string;
  chainId: bigint;
  name: string;
  symbol: string;
  decimals: number;
  gasLimit: bigint;
}
export interface MsgWhitelistERC20ProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgWhitelistERC20";
  value: Uint8Array;
}
export interface MsgWhitelistERC20Amino {
  creator?: string;
  erc20_address?: string;
  chain_id?: string;
  name?: string;
  symbol?: string;
  decimals?: number;
  gas_limit?: string;
}
export interface MsgWhitelistERC20AminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgWhitelistERC20";
  value: MsgWhitelistERC20Amino;
}
export interface MsgWhitelistERC20SDKType {
  creator: string;
  erc20_address: string;
  chain_id: bigint;
  name: string;
  symbol: string;
  decimals: number;
  gas_limit: bigint;
}
export interface MsgWhitelistERC20Response {
  zrc20Address: string;
  cctxIndex: string;
}
export interface MsgWhitelistERC20ResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgWhitelistERC20Response";
  value: Uint8Array;
}
export interface MsgWhitelistERC20ResponseAmino {
  zrc20_address?: string;
  cctx_index?: string;
}
export interface MsgWhitelistERC20ResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgWhitelistERC20Response";
  value: MsgWhitelistERC20ResponseAmino;
}
export interface MsgWhitelistERC20ResponseSDKType {
  zrc20_address: string;
  cctx_index: string;
}
export interface MsgAddToOutTxTracker {
  creator: string;
  chainId: bigint;
  nonce: bigint;
  txHash: string;
  proof?: Proof;
  blockHash: string;
  txIndex: bigint;
}
export interface MsgAddToOutTxTrackerProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAddToOutTxTracker";
  value: Uint8Array;
}
export interface MsgAddToOutTxTrackerAmino {
  creator?: string;
  chain_id?: string;
  nonce?: string;
  tx_hash?: string;
  proof?: ProofAmino;
  block_hash?: string;
  tx_index?: string;
}
export interface MsgAddToOutTxTrackerAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgAddToOutTxTracker";
  value: MsgAddToOutTxTrackerAmino;
}
export interface MsgAddToOutTxTrackerSDKType {
  creator: string;
  chain_id: bigint;
  nonce: bigint;
  tx_hash: string;
  proof?: ProofSDKType;
  block_hash: string;
  tx_index: bigint;
}
export interface MsgAddToOutTxTrackerResponse {
  /** if the tx was removed from the tracker due to no pending cctx */
  isRemoved: boolean;
}
export interface MsgAddToOutTxTrackerResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAddToOutTxTrackerResponse";
  value: Uint8Array;
}
export interface MsgAddToOutTxTrackerResponseAmino {
  /** if the tx was removed from the tracker due to no pending cctx */
  is_removed?: boolean;
}
export interface MsgAddToOutTxTrackerResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgAddToOutTxTrackerResponse";
  value: MsgAddToOutTxTrackerResponseAmino;
}
export interface MsgAddToOutTxTrackerResponseSDKType {
  is_removed: boolean;
}
export interface MsgRemoveFromOutTxTracker {
  creator: string;
  chainId: bigint;
  nonce: bigint;
}
export interface MsgRemoveFromOutTxTrackerProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker";
  value: Uint8Array;
}
export interface MsgRemoveFromOutTxTrackerAmino {
  creator?: string;
  chain_id?: string;
  nonce?: string;
}
export interface MsgRemoveFromOutTxTrackerAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker";
  value: MsgRemoveFromOutTxTrackerAmino;
}
export interface MsgRemoveFromOutTxTrackerSDKType {
  creator: string;
  chain_id: bigint;
  nonce: bigint;
}
export interface MsgRemoveFromOutTxTrackerResponse {}
export interface MsgRemoveFromOutTxTrackerResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTrackerResponse";
  value: Uint8Array;
}
export interface MsgRemoveFromOutTxTrackerResponseAmino {}
export interface MsgRemoveFromOutTxTrackerResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTrackerResponse";
  value: MsgRemoveFromOutTxTrackerResponseAmino;
}
export interface MsgRemoveFromOutTxTrackerResponseSDKType {}
export interface MsgGasPriceVoter {
  creator: string;
  chainId: bigint;
  price: bigint;
  blockNumber: bigint;
  supply: string;
}
export interface MsgGasPriceVoterProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgGasPriceVoter";
  value: Uint8Array;
}
export interface MsgGasPriceVoterAmino {
  creator?: string;
  chain_id?: string;
  price?: string;
  block_number?: string;
  supply?: string;
}
export interface MsgGasPriceVoterAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgGasPriceVoter";
  value: MsgGasPriceVoterAmino;
}
export interface MsgGasPriceVoterSDKType {
  creator: string;
  chain_id: bigint;
  price: bigint;
  block_number: bigint;
  supply: string;
}
export interface MsgGasPriceVoterResponse {}
export interface MsgGasPriceVoterResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgGasPriceVoterResponse";
  value: Uint8Array;
}
export interface MsgGasPriceVoterResponseAmino {}
export interface MsgGasPriceVoterResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgGasPriceVoterResponse";
  value: MsgGasPriceVoterResponseAmino;
}
export interface MsgGasPriceVoterResponseSDKType {}
export interface MsgVoteOnObservedOutboundTx {
  creator: string;
  cctxHash: string;
  observedOutTxHash: string;
  observedOutTxBlockHeight: bigint;
  observedOutTxGasUsed: bigint;
  observedOutTxEffectiveGasPrice: string;
  observedOutTxEffectiveGasLimit: bigint;
  valueReceived: string;
  status: ReceiveStatus;
  outTxChain: bigint;
  outTxTssNonce: bigint;
  coinType: CoinType;
}
export interface MsgVoteOnObservedOutboundTxProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx";
  value: Uint8Array;
}
export interface MsgVoteOnObservedOutboundTxAmino {
  creator?: string;
  cctx_hash?: string;
  observed_outTx_hash?: string;
  observed_outTx_blockHeight?: string;
  observed_outTx_gas_used?: string;
  observed_outTx_effective_gas_price?: string;
  observed_outTx_effective_gas_limit?: string;
  value_received?: string;
  status?: ReceiveStatus;
  outTx_chain?: string;
  outTx_tss_nonce?: string;
  coin_type?: CoinType;
}
export interface MsgVoteOnObservedOutboundTxAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx";
  value: MsgVoteOnObservedOutboundTxAmino;
}
export interface MsgVoteOnObservedOutboundTxSDKType {
  creator: string;
  cctx_hash: string;
  observed_outTx_hash: string;
  observed_outTx_blockHeight: bigint;
  observed_outTx_gas_used: bigint;
  observed_outTx_effective_gas_price: string;
  observed_outTx_effective_gas_limit: bigint;
  value_received: string;
  status: ReceiveStatus;
  outTx_chain: bigint;
  outTx_tss_nonce: bigint;
  coin_type: CoinType;
}
export interface MsgVoteOnObservedOutboundTxResponse {}
export interface MsgVoteOnObservedOutboundTxResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTxResponse";
  value: Uint8Array;
}
export interface MsgVoteOnObservedOutboundTxResponseAmino {}
export interface MsgVoteOnObservedOutboundTxResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTxResponse";
  value: MsgVoteOnObservedOutboundTxResponseAmino;
}
export interface MsgVoteOnObservedOutboundTxResponseSDKType {}
export interface MsgVoteOnObservedInboundTx {
  creator: string;
  sender: string;
  senderChainId: bigint;
  receiver: string;
  receiverChain: bigint;
  /** string zeta_burnt = 6; */
  amount: string;
  /** string mMint = 7; */
  message: string;
  inTxHash: string;
  inBlockHeight: bigint;
  gasLimit: bigint;
  coinType: CoinType;
  txOrigin: string;
  asset: string;
  /** event index of the sent asset in the observed tx */
  eventIndex: bigint;
}
export interface MsgVoteOnObservedInboundTxProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx";
  value: Uint8Array;
}
export interface MsgVoteOnObservedInboundTxAmino {
  creator?: string;
  sender?: string;
  sender_chain_id?: string;
  receiver?: string;
  receiver_chain?: string;
  /** string zeta_burnt = 6; */
  amount?: string;
  /** string mMint = 7; */
  message?: string;
  in_tx_hash?: string;
  in_block_height?: string;
  gas_limit?: string;
  coin_type?: CoinType;
  tx_origin?: string;
  asset?: string;
  /** event index of the sent asset in the observed tx */
  event_index?: string;
}
export interface MsgVoteOnObservedInboundTxAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx";
  value: MsgVoteOnObservedInboundTxAmino;
}
export interface MsgVoteOnObservedInboundTxSDKType {
  creator: string;
  sender: string;
  sender_chain_id: bigint;
  receiver: string;
  receiver_chain: bigint;
  amount: string;
  message: string;
  in_tx_hash: string;
  in_block_height: bigint;
  gas_limit: bigint;
  coin_type: CoinType;
  tx_origin: string;
  asset: string;
  event_index: bigint;
}
export interface MsgVoteOnObservedInboundTxResponse {}
export interface MsgVoteOnObservedInboundTxResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTxResponse";
  value: Uint8Array;
}
export interface MsgVoteOnObservedInboundTxResponseAmino {}
export interface MsgVoteOnObservedInboundTxResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTxResponse";
  value: MsgVoteOnObservedInboundTxResponseAmino;
}
export interface MsgVoteOnObservedInboundTxResponseSDKType {}
export interface MsgAbortStuckCCTX {
  creator: string;
  cctxIndex: string;
}
export interface MsgAbortStuckCCTXProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX";
  value: Uint8Array;
}
export interface MsgAbortStuckCCTXAmino {
  creator?: string;
  cctx_index?: string;
}
export interface MsgAbortStuckCCTXAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX";
  value: MsgAbortStuckCCTXAmino;
}
export interface MsgAbortStuckCCTXSDKType {
  creator: string;
  cctx_index: string;
}
export interface MsgAbortStuckCCTXResponse {}
export interface MsgAbortStuckCCTXResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTXResponse";
  value: Uint8Array;
}
export interface MsgAbortStuckCCTXResponseAmino {}
export interface MsgAbortStuckCCTXResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTXResponse";
  value: MsgAbortStuckCCTXResponseAmino;
}
export interface MsgAbortStuckCCTXResponseSDKType {}
export interface MsgRefundAbortedCCTX {
  creator: string;
  cctxIndex: string;
  /** if not provided, the refund will be sent to the sender/txOrgin */
  refundAddress: string;
}
export interface MsgRefundAbortedCCTXProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX";
  value: Uint8Array;
}
export interface MsgRefundAbortedCCTXAmino {
  creator?: string;
  cctx_index?: string;
  /** if not provided, the refund will be sent to the sender/txOrgin */
  refund_address?: string;
}
export interface MsgRefundAbortedCCTXAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX";
  value: MsgRefundAbortedCCTXAmino;
}
export interface MsgRefundAbortedCCTXSDKType {
  creator: string;
  cctx_index: string;
  refund_address: string;
}
export interface MsgRefundAbortedCCTXResponse {}
export interface MsgRefundAbortedCCTXResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTXResponse";
  value: Uint8Array;
}
export interface MsgRefundAbortedCCTXResponseAmino {}
export interface MsgRefundAbortedCCTXResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTXResponse";
  value: MsgRefundAbortedCCTXResponseAmino;
}
export interface MsgRefundAbortedCCTXResponseSDKType {}
function createBaseMsgCreateTSSVoter(): MsgCreateTSSVoter {
  return {
    creator: "",
    tssPubkey: "",
    keyGenZetaHeight: BigInt(0),
    status: 0
  };
}
export const MsgCreateTSSVoter = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgCreateTSSVoter",
  encode(message: MsgCreateTSSVoter, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.tssPubkey !== "") {
      writer.uint32(18).string(message.tssPubkey);
    }
    if (message.keyGenZetaHeight !== BigInt(0)) {
      writer.uint32(24).int64(message.keyGenZetaHeight);
    }
    if (message.status !== 0) {
      writer.uint32(32).int32(message.status);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateTSSVoter {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateTSSVoter();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.tssPubkey = reader.string();
          break;
        case 3:
          message.keyGenZetaHeight = reader.int64();
          break;
        case 4:
          message.status = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgCreateTSSVoter>): MsgCreateTSSVoter {
    const message = createBaseMsgCreateTSSVoter();
    message.creator = object.creator ?? "";
    message.tssPubkey = object.tssPubkey ?? "";
    message.keyGenZetaHeight = object.keyGenZetaHeight !== undefined && object.keyGenZetaHeight !== null ? BigInt(object.keyGenZetaHeight.toString()) : BigInt(0);
    message.status = object.status ?? 0;
    return message;
  },
  fromAmino(object: MsgCreateTSSVoterAmino): MsgCreateTSSVoter {
    const message = createBaseMsgCreateTSSVoter();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.tss_pubkey !== undefined && object.tss_pubkey !== null) {
      message.tssPubkey = object.tss_pubkey;
    }
    if (object.keyGenZetaHeight !== undefined && object.keyGenZetaHeight !== null) {
      message.keyGenZetaHeight = BigInt(object.keyGenZetaHeight);
    }
    if (object.status !== undefined && object.status !== null) {
      message.status = receiveStatusFromJSON(object.status);
    }
    return message;
  },
  toAmino(message: MsgCreateTSSVoter): MsgCreateTSSVoterAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.tss_pubkey = message.tssPubkey;
    obj.keyGenZetaHeight = message.keyGenZetaHeight ? message.keyGenZetaHeight.toString() : undefined;
    obj.status = message.status;
    return obj;
  },
  fromAminoMsg(object: MsgCreateTSSVoterAminoMsg): MsgCreateTSSVoter {
    return MsgCreateTSSVoter.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreateTSSVoterProtoMsg): MsgCreateTSSVoter {
    return MsgCreateTSSVoter.decode(message.value);
  },
  toProto(message: MsgCreateTSSVoter): Uint8Array {
    return MsgCreateTSSVoter.encode(message).finish();
  },
  toProtoMsg(message: MsgCreateTSSVoter): MsgCreateTSSVoterProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgCreateTSSVoter",
      value: MsgCreateTSSVoter.encode(message).finish()
    };
  }
};
function createBaseMsgCreateTSSVoterResponse(): MsgCreateTSSVoterResponse {
  return {};
}
export const MsgCreateTSSVoterResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgCreateTSSVoterResponse",
  encode(_: MsgCreateTSSVoterResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCreateTSSVoterResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateTSSVoterResponse();
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
  fromPartial(_: Partial<MsgCreateTSSVoterResponse>): MsgCreateTSSVoterResponse {
    const message = createBaseMsgCreateTSSVoterResponse();
    return message;
  },
  fromAmino(_: MsgCreateTSSVoterResponseAmino): MsgCreateTSSVoterResponse {
    const message = createBaseMsgCreateTSSVoterResponse();
    return message;
  },
  toAmino(_: MsgCreateTSSVoterResponse): MsgCreateTSSVoterResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgCreateTSSVoterResponseAminoMsg): MsgCreateTSSVoterResponse {
    return MsgCreateTSSVoterResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCreateTSSVoterResponseProtoMsg): MsgCreateTSSVoterResponse {
    return MsgCreateTSSVoterResponse.decode(message.value);
  },
  toProto(message: MsgCreateTSSVoterResponse): Uint8Array {
    return MsgCreateTSSVoterResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgCreateTSSVoterResponse): MsgCreateTSSVoterResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgCreateTSSVoterResponse",
      value: MsgCreateTSSVoterResponse.encode(message).finish()
    };
  }
};
function createBaseMsgMigrateTssFunds(): MsgMigrateTssFunds {
  return {
    creator: "",
    chainId: BigInt(0),
    amount: ""
  };
}
export const MsgMigrateTssFunds = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
  encode(message: MsgMigrateTssFunds, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.amount !== "") {
      writer.uint32(26).string(message.amount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgMigrateTssFunds {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgMigrateTssFunds();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.amount = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgMigrateTssFunds>): MsgMigrateTssFunds {
    const message = createBaseMsgMigrateTssFunds();
    message.creator = object.creator ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.amount = object.amount ?? "";
    return message;
  },
  fromAmino(object: MsgMigrateTssFundsAmino): MsgMigrateTssFunds {
    const message = createBaseMsgMigrateTssFunds();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    return message;
  },
  toAmino(message: MsgMigrateTssFunds): MsgMigrateTssFundsAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.amount = message.amount;
    return obj;
  },
  fromAminoMsg(object: MsgMigrateTssFundsAminoMsg): MsgMigrateTssFunds {
    return MsgMigrateTssFunds.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgMigrateTssFundsProtoMsg): MsgMigrateTssFunds {
    return MsgMigrateTssFunds.decode(message.value);
  },
  toProto(message: MsgMigrateTssFunds): Uint8Array {
    return MsgMigrateTssFunds.encode(message).finish();
  },
  toProtoMsg(message: MsgMigrateTssFunds): MsgMigrateTssFundsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
      value: MsgMigrateTssFunds.encode(message).finish()
    };
  }
};
function createBaseMsgMigrateTssFundsResponse(): MsgMigrateTssFundsResponse {
  return {};
}
export const MsgMigrateTssFundsResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgMigrateTssFundsResponse",
  encode(_: MsgMigrateTssFundsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgMigrateTssFundsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgMigrateTssFundsResponse();
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
  fromPartial(_: Partial<MsgMigrateTssFundsResponse>): MsgMigrateTssFundsResponse {
    const message = createBaseMsgMigrateTssFundsResponse();
    return message;
  },
  fromAmino(_: MsgMigrateTssFundsResponseAmino): MsgMigrateTssFundsResponse {
    const message = createBaseMsgMigrateTssFundsResponse();
    return message;
  },
  toAmino(_: MsgMigrateTssFundsResponse): MsgMigrateTssFundsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgMigrateTssFundsResponseAminoMsg): MsgMigrateTssFundsResponse {
    return MsgMigrateTssFundsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgMigrateTssFundsResponseProtoMsg): MsgMigrateTssFundsResponse {
    return MsgMigrateTssFundsResponse.decode(message.value);
  },
  toProto(message: MsgMigrateTssFundsResponse): Uint8Array {
    return MsgMigrateTssFundsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgMigrateTssFundsResponse): MsgMigrateTssFundsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgMigrateTssFundsResponse",
      value: MsgMigrateTssFundsResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateTssAddress(): MsgUpdateTssAddress {
  return {
    creator: "",
    tssPubkey: ""
  };
}
export const MsgUpdateTssAddress = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgUpdateTssAddress",
  encode(message: MsgUpdateTssAddress, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.tssPubkey !== "") {
      writer.uint32(18).string(message.tssPubkey);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateTssAddress {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateTssAddress();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.tssPubkey = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateTssAddress>): MsgUpdateTssAddress {
    const message = createBaseMsgUpdateTssAddress();
    message.creator = object.creator ?? "";
    message.tssPubkey = object.tssPubkey ?? "";
    return message;
  },
  fromAmino(object: MsgUpdateTssAddressAmino): MsgUpdateTssAddress {
    const message = createBaseMsgUpdateTssAddress();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.tss_pubkey !== undefined && object.tss_pubkey !== null) {
      message.tssPubkey = object.tss_pubkey;
    }
    return message;
  },
  toAmino(message: MsgUpdateTssAddress): MsgUpdateTssAddressAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.tss_pubkey = message.tssPubkey;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateTssAddressAminoMsg): MsgUpdateTssAddress {
    return MsgUpdateTssAddress.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateTssAddressProtoMsg): MsgUpdateTssAddress {
    return MsgUpdateTssAddress.decode(message.value);
  },
  toProto(message: MsgUpdateTssAddress): Uint8Array {
    return MsgUpdateTssAddress.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateTssAddress): MsgUpdateTssAddressProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgUpdateTssAddress",
      value: MsgUpdateTssAddress.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateTssAddressResponse(): MsgUpdateTssAddressResponse {
  return {};
}
export const MsgUpdateTssAddressResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgUpdateTssAddressResponse",
  encode(_: MsgUpdateTssAddressResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateTssAddressResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateTssAddressResponse();
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
  fromPartial(_: Partial<MsgUpdateTssAddressResponse>): MsgUpdateTssAddressResponse {
    const message = createBaseMsgUpdateTssAddressResponse();
    return message;
  },
  fromAmino(_: MsgUpdateTssAddressResponseAmino): MsgUpdateTssAddressResponse {
    const message = createBaseMsgUpdateTssAddressResponse();
    return message;
  },
  toAmino(_: MsgUpdateTssAddressResponse): MsgUpdateTssAddressResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateTssAddressResponseAminoMsg): MsgUpdateTssAddressResponse {
    return MsgUpdateTssAddressResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateTssAddressResponseProtoMsg): MsgUpdateTssAddressResponse {
    return MsgUpdateTssAddressResponse.decode(message.value);
  },
  toProto(message: MsgUpdateTssAddressResponse): Uint8Array {
    return MsgUpdateTssAddressResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateTssAddressResponse): MsgUpdateTssAddressResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgUpdateTssAddressResponse",
      value: MsgUpdateTssAddressResponse.encode(message).finish()
    };
  }
};
function createBaseMsgAddToInTxTracker(): MsgAddToInTxTracker {
  return {
    creator: "",
    chainId: BigInt(0),
    txHash: "",
    coinType: 0,
    proof: undefined,
    blockHash: "",
    txIndex: BigInt(0)
  };
}
export const MsgAddToInTxTracker = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAddToInTxTracker",
  encode(message: MsgAddToInTxTracker, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.txHash !== "") {
      writer.uint32(26).string(message.txHash);
    }
    if (message.coinType !== 0) {
      writer.uint32(32).int32(message.coinType);
    }
    if (message.proof !== undefined) {
      Proof.encode(message.proof, writer.uint32(42).fork()).ldelim();
    }
    if (message.blockHash !== "") {
      writer.uint32(50).string(message.blockHash);
    }
    if (message.txIndex !== BigInt(0)) {
      writer.uint32(56).int64(message.txIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddToInTxTracker {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddToInTxTracker();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.txHash = reader.string();
          break;
        case 4:
          message.coinType = (reader.int32() as any);
          break;
        case 5:
          message.proof = Proof.decode(reader, reader.uint32());
          break;
        case 6:
          message.blockHash = reader.string();
          break;
        case 7:
          message.txIndex = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgAddToInTxTracker>): MsgAddToInTxTracker {
    const message = createBaseMsgAddToInTxTracker();
    message.creator = object.creator ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.txHash = object.txHash ?? "";
    message.coinType = object.coinType ?? 0;
    message.proof = object.proof !== undefined && object.proof !== null ? Proof.fromPartial(object.proof) : undefined;
    message.blockHash = object.blockHash ?? "";
    message.txIndex = object.txIndex !== undefined && object.txIndex !== null ? BigInt(object.txIndex.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgAddToInTxTrackerAmino): MsgAddToInTxTracker {
    const message = createBaseMsgAddToInTxTracker();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.tx_hash !== undefined && object.tx_hash !== null) {
      message.txHash = object.tx_hash;
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    if (object.proof !== undefined && object.proof !== null) {
      message.proof = Proof.fromAmino(object.proof);
    }
    if (object.block_hash !== undefined && object.block_hash !== null) {
      message.blockHash = object.block_hash;
    }
    if (object.tx_index !== undefined && object.tx_index !== null) {
      message.txIndex = BigInt(object.tx_index);
    }
    return message;
  },
  toAmino(message: MsgAddToInTxTracker): MsgAddToInTxTrackerAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.tx_hash = message.txHash;
    obj.coin_type = message.coinType;
    obj.proof = message.proof ? Proof.toAmino(message.proof) : undefined;
    obj.block_hash = message.blockHash;
    obj.tx_index = message.txIndex ? message.txIndex.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgAddToInTxTrackerAminoMsg): MsgAddToInTxTracker {
    return MsgAddToInTxTracker.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddToInTxTrackerProtoMsg): MsgAddToInTxTracker {
    return MsgAddToInTxTracker.decode(message.value);
  },
  toProto(message: MsgAddToInTxTracker): Uint8Array {
    return MsgAddToInTxTracker.encode(message).finish();
  },
  toProtoMsg(message: MsgAddToInTxTracker): MsgAddToInTxTrackerProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgAddToInTxTracker",
      value: MsgAddToInTxTracker.encode(message).finish()
    };
  }
};
function createBaseMsgAddToInTxTrackerResponse(): MsgAddToInTxTrackerResponse {
  return {};
}
export const MsgAddToInTxTrackerResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAddToInTxTrackerResponse",
  encode(_: MsgAddToInTxTrackerResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddToInTxTrackerResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddToInTxTrackerResponse();
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
  fromPartial(_: Partial<MsgAddToInTxTrackerResponse>): MsgAddToInTxTrackerResponse {
    const message = createBaseMsgAddToInTxTrackerResponse();
    return message;
  },
  fromAmino(_: MsgAddToInTxTrackerResponseAmino): MsgAddToInTxTrackerResponse {
    const message = createBaseMsgAddToInTxTrackerResponse();
    return message;
  },
  toAmino(_: MsgAddToInTxTrackerResponse): MsgAddToInTxTrackerResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgAddToInTxTrackerResponseAminoMsg): MsgAddToInTxTrackerResponse {
    return MsgAddToInTxTrackerResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddToInTxTrackerResponseProtoMsg): MsgAddToInTxTrackerResponse {
    return MsgAddToInTxTrackerResponse.decode(message.value);
  },
  toProto(message: MsgAddToInTxTrackerResponse): Uint8Array {
    return MsgAddToInTxTrackerResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgAddToInTxTrackerResponse): MsgAddToInTxTrackerResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgAddToInTxTrackerResponse",
      value: MsgAddToInTxTrackerResponse.encode(message).finish()
    };
  }
};
function createBaseMsgWhitelistERC20(): MsgWhitelistERC20 {
  return {
    creator: "",
    erc20Address: "",
    chainId: BigInt(0),
    name: "",
    symbol: "",
    decimals: 0,
    gasLimit: BigInt(0)
  };
}
export const MsgWhitelistERC20 = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgWhitelistERC20",
  encode(message: MsgWhitelistERC20, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.erc20Address !== "") {
      writer.uint32(18).string(message.erc20Address);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(24).int64(message.chainId);
    }
    if (message.name !== "") {
      writer.uint32(34).string(message.name);
    }
    if (message.symbol !== "") {
      writer.uint32(42).string(message.symbol);
    }
    if (message.decimals !== 0) {
      writer.uint32(48).uint32(message.decimals);
    }
    if (message.gasLimit !== BigInt(0)) {
      writer.uint32(56).int64(message.gasLimit);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgWhitelistERC20 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWhitelistERC20();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.erc20Address = reader.string();
          break;
        case 3:
          message.chainId = reader.int64();
          break;
        case 4:
          message.name = reader.string();
          break;
        case 5:
          message.symbol = reader.string();
          break;
        case 6:
          message.decimals = reader.uint32();
          break;
        case 7:
          message.gasLimit = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgWhitelistERC20>): MsgWhitelistERC20 {
    const message = createBaseMsgWhitelistERC20();
    message.creator = object.creator ?? "";
    message.erc20Address = object.erc20Address ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.name = object.name ?? "";
    message.symbol = object.symbol ?? "";
    message.decimals = object.decimals ?? 0;
    message.gasLimit = object.gasLimit !== undefined && object.gasLimit !== null ? BigInt(object.gasLimit.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgWhitelistERC20Amino): MsgWhitelistERC20 {
    const message = createBaseMsgWhitelistERC20();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.erc20_address !== undefined && object.erc20_address !== null) {
      message.erc20Address = object.erc20_address;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    }
    if (object.symbol !== undefined && object.symbol !== null) {
      message.symbol = object.symbol;
    }
    if (object.decimals !== undefined && object.decimals !== null) {
      message.decimals = object.decimals;
    }
    if (object.gas_limit !== undefined && object.gas_limit !== null) {
      message.gasLimit = BigInt(object.gas_limit);
    }
    return message;
  },
  toAmino(message: MsgWhitelistERC20): MsgWhitelistERC20Amino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.erc20_address = message.erc20Address;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.name = message.name;
    obj.symbol = message.symbol;
    obj.decimals = message.decimals;
    obj.gas_limit = message.gasLimit ? message.gasLimit.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgWhitelistERC20AminoMsg): MsgWhitelistERC20 {
    return MsgWhitelistERC20.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgWhitelistERC20ProtoMsg): MsgWhitelistERC20 {
    return MsgWhitelistERC20.decode(message.value);
  },
  toProto(message: MsgWhitelistERC20): Uint8Array {
    return MsgWhitelistERC20.encode(message).finish();
  },
  toProtoMsg(message: MsgWhitelistERC20): MsgWhitelistERC20ProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgWhitelistERC20",
      value: MsgWhitelistERC20.encode(message).finish()
    };
  }
};
function createBaseMsgWhitelistERC20Response(): MsgWhitelistERC20Response {
  return {
    zrc20Address: "",
    cctxIndex: ""
  };
}
export const MsgWhitelistERC20Response = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgWhitelistERC20Response",
  encode(message: MsgWhitelistERC20Response, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.zrc20Address !== "") {
      writer.uint32(10).string(message.zrc20Address);
    }
    if (message.cctxIndex !== "") {
      writer.uint32(18).string(message.cctxIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgWhitelistERC20Response {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgWhitelistERC20Response();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.zrc20Address = reader.string();
          break;
        case 2:
          message.cctxIndex = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgWhitelistERC20Response>): MsgWhitelistERC20Response {
    const message = createBaseMsgWhitelistERC20Response();
    message.zrc20Address = object.zrc20Address ?? "";
    message.cctxIndex = object.cctxIndex ?? "";
    return message;
  },
  fromAmino(object: MsgWhitelistERC20ResponseAmino): MsgWhitelistERC20Response {
    const message = createBaseMsgWhitelistERC20Response();
    if (object.zrc20_address !== undefined && object.zrc20_address !== null) {
      message.zrc20Address = object.zrc20_address;
    }
    if (object.cctx_index !== undefined && object.cctx_index !== null) {
      message.cctxIndex = object.cctx_index;
    }
    return message;
  },
  toAmino(message: MsgWhitelistERC20Response): MsgWhitelistERC20ResponseAmino {
    const obj: any = {};
    obj.zrc20_address = message.zrc20Address;
    obj.cctx_index = message.cctxIndex;
    return obj;
  },
  fromAminoMsg(object: MsgWhitelistERC20ResponseAminoMsg): MsgWhitelistERC20Response {
    return MsgWhitelistERC20Response.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgWhitelistERC20ResponseProtoMsg): MsgWhitelistERC20Response {
    return MsgWhitelistERC20Response.decode(message.value);
  },
  toProto(message: MsgWhitelistERC20Response): Uint8Array {
    return MsgWhitelistERC20Response.encode(message).finish();
  },
  toProtoMsg(message: MsgWhitelistERC20Response): MsgWhitelistERC20ResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgWhitelistERC20Response",
      value: MsgWhitelistERC20Response.encode(message).finish()
    };
  }
};
function createBaseMsgAddToOutTxTracker(): MsgAddToOutTxTracker {
  return {
    creator: "",
    chainId: BigInt(0),
    nonce: BigInt(0),
    txHash: "",
    proof: undefined,
    blockHash: "",
    txIndex: BigInt(0)
  };
}
export const MsgAddToOutTxTracker = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAddToOutTxTracker",
  encode(message: MsgAddToOutTxTracker, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.nonce !== BigInt(0)) {
      writer.uint32(24).uint64(message.nonce);
    }
    if (message.txHash !== "") {
      writer.uint32(34).string(message.txHash);
    }
    if (message.proof !== undefined) {
      Proof.encode(message.proof, writer.uint32(42).fork()).ldelim();
    }
    if (message.blockHash !== "") {
      writer.uint32(50).string(message.blockHash);
    }
    if (message.txIndex !== BigInt(0)) {
      writer.uint32(56).int64(message.txIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddToOutTxTracker {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddToOutTxTracker();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.nonce = reader.uint64();
          break;
        case 4:
          message.txHash = reader.string();
          break;
        case 5:
          message.proof = Proof.decode(reader, reader.uint32());
          break;
        case 6:
          message.blockHash = reader.string();
          break;
        case 7:
          message.txIndex = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgAddToOutTxTracker>): MsgAddToOutTxTracker {
    const message = createBaseMsgAddToOutTxTracker();
    message.creator = object.creator ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.nonce = object.nonce !== undefined && object.nonce !== null ? BigInt(object.nonce.toString()) : BigInt(0);
    message.txHash = object.txHash ?? "";
    message.proof = object.proof !== undefined && object.proof !== null ? Proof.fromPartial(object.proof) : undefined;
    message.blockHash = object.blockHash ?? "";
    message.txIndex = object.txIndex !== undefined && object.txIndex !== null ? BigInt(object.txIndex.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgAddToOutTxTrackerAmino): MsgAddToOutTxTracker {
    const message = createBaseMsgAddToOutTxTracker();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = BigInt(object.nonce);
    }
    if (object.tx_hash !== undefined && object.tx_hash !== null) {
      message.txHash = object.tx_hash;
    }
    if (object.proof !== undefined && object.proof !== null) {
      message.proof = Proof.fromAmino(object.proof);
    }
    if (object.block_hash !== undefined && object.block_hash !== null) {
      message.blockHash = object.block_hash;
    }
    if (object.tx_index !== undefined && object.tx_index !== null) {
      message.txIndex = BigInt(object.tx_index);
    }
    return message;
  },
  toAmino(message: MsgAddToOutTxTracker): MsgAddToOutTxTrackerAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.nonce = message.nonce ? message.nonce.toString() : undefined;
    obj.tx_hash = message.txHash;
    obj.proof = message.proof ? Proof.toAmino(message.proof) : undefined;
    obj.block_hash = message.blockHash;
    obj.tx_index = message.txIndex ? message.txIndex.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgAddToOutTxTrackerAminoMsg): MsgAddToOutTxTracker {
    return MsgAddToOutTxTracker.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddToOutTxTrackerProtoMsg): MsgAddToOutTxTracker {
    return MsgAddToOutTxTracker.decode(message.value);
  },
  toProto(message: MsgAddToOutTxTracker): Uint8Array {
    return MsgAddToOutTxTracker.encode(message).finish();
  },
  toProtoMsg(message: MsgAddToOutTxTracker): MsgAddToOutTxTrackerProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgAddToOutTxTracker",
      value: MsgAddToOutTxTracker.encode(message).finish()
    };
  }
};
function createBaseMsgAddToOutTxTrackerResponse(): MsgAddToOutTxTrackerResponse {
  return {
    isRemoved: false
  };
}
export const MsgAddToOutTxTrackerResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAddToOutTxTrackerResponse",
  encode(message: MsgAddToOutTxTrackerResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.isRemoved === true) {
      writer.uint32(8).bool(message.isRemoved);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAddToOutTxTrackerResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAddToOutTxTrackerResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.isRemoved = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgAddToOutTxTrackerResponse>): MsgAddToOutTxTrackerResponse {
    const message = createBaseMsgAddToOutTxTrackerResponse();
    message.isRemoved = object.isRemoved ?? false;
    return message;
  },
  fromAmino(object: MsgAddToOutTxTrackerResponseAmino): MsgAddToOutTxTrackerResponse {
    const message = createBaseMsgAddToOutTxTrackerResponse();
    if (object.is_removed !== undefined && object.is_removed !== null) {
      message.isRemoved = object.is_removed;
    }
    return message;
  },
  toAmino(message: MsgAddToOutTxTrackerResponse): MsgAddToOutTxTrackerResponseAmino {
    const obj: any = {};
    obj.is_removed = message.isRemoved;
    return obj;
  },
  fromAminoMsg(object: MsgAddToOutTxTrackerResponseAminoMsg): MsgAddToOutTxTrackerResponse {
    return MsgAddToOutTxTrackerResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAddToOutTxTrackerResponseProtoMsg): MsgAddToOutTxTrackerResponse {
    return MsgAddToOutTxTrackerResponse.decode(message.value);
  },
  toProto(message: MsgAddToOutTxTrackerResponse): Uint8Array {
    return MsgAddToOutTxTrackerResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgAddToOutTxTrackerResponse): MsgAddToOutTxTrackerResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgAddToOutTxTrackerResponse",
      value: MsgAddToOutTxTrackerResponse.encode(message).finish()
    };
  }
};
function createBaseMsgRemoveFromOutTxTracker(): MsgRemoveFromOutTxTracker {
  return {
    creator: "",
    chainId: BigInt(0),
    nonce: BigInt(0)
  };
}
export const MsgRemoveFromOutTxTracker = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker",
  encode(message: MsgRemoveFromOutTxTracker, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.nonce !== BigInt(0)) {
      writer.uint32(24).uint64(message.nonce);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgRemoveFromOutTxTracker {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveFromOutTxTracker();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.nonce = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgRemoveFromOutTxTracker>): MsgRemoveFromOutTxTracker {
    const message = createBaseMsgRemoveFromOutTxTracker();
    message.creator = object.creator ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.nonce = object.nonce !== undefined && object.nonce !== null ? BigInt(object.nonce.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgRemoveFromOutTxTrackerAmino): MsgRemoveFromOutTxTracker {
    const message = createBaseMsgRemoveFromOutTxTracker();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = BigInt(object.nonce);
    }
    return message;
  },
  toAmino(message: MsgRemoveFromOutTxTracker): MsgRemoveFromOutTxTrackerAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.nonce = message.nonce ? message.nonce.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgRemoveFromOutTxTrackerAminoMsg): MsgRemoveFromOutTxTracker {
    return MsgRemoveFromOutTxTracker.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRemoveFromOutTxTrackerProtoMsg): MsgRemoveFromOutTxTracker {
    return MsgRemoveFromOutTxTracker.decode(message.value);
  },
  toProto(message: MsgRemoveFromOutTxTracker): Uint8Array {
    return MsgRemoveFromOutTxTracker.encode(message).finish();
  },
  toProtoMsg(message: MsgRemoveFromOutTxTracker): MsgRemoveFromOutTxTrackerProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker",
      value: MsgRemoveFromOutTxTracker.encode(message).finish()
    };
  }
};
function createBaseMsgRemoveFromOutTxTrackerResponse(): MsgRemoveFromOutTxTrackerResponse {
  return {};
}
export const MsgRemoveFromOutTxTrackerResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTrackerResponse",
  encode(_: MsgRemoveFromOutTxTrackerResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgRemoveFromOutTxTrackerResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRemoveFromOutTxTrackerResponse();
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
  fromPartial(_: Partial<MsgRemoveFromOutTxTrackerResponse>): MsgRemoveFromOutTxTrackerResponse {
    const message = createBaseMsgRemoveFromOutTxTrackerResponse();
    return message;
  },
  fromAmino(_: MsgRemoveFromOutTxTrackerResponseAmino): MsgRemoveFromOutTxTrackerResponse {
    const message = createBaseMsgRemoveFromOutTxTrackerResponse();
    return message;
  },
  toAmino(_: MsgRemoveFromOutTxTrackerResponse): MsgRemoveFromOutTxTrackerResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgRemoveFromOutTxTrackerResponseAminoMsg): MsgRemoveFromOutTxTrackerResponse {
    return MsgRemoveFromOutTxTrackerResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRemoveFromOutTxTrackerResponseProtoMsg): MsgRemoveFromOutTxTrackerResponse {
    return MsgRemoveFromOutTxTrackerResponse.decode(message.value);
  },
  toProto(message: MsgRemoveFromOutTxTrackerResponse): Uint8Array {
    return MsgRemoveFromOutTxTrackerResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgRemoveFromOutTxTrackerResponse): MsgRemoveFromOutTxTrackerResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTrackerResponse",
      value: MsgRemoveFromOutTxTrackerResponse.encode(message).finish()
    };
  }
};
function createBaseMsgGasPriceVoter(): MsgGasPriceVoter {
  return {
    creator: "",
    chainId: BigInt(0),
    price: BigInt(0),
    blockNumber: BigInt(0),
    supply: ""
  };
}
export const MsgGasPriceVoter = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgGasPriceVoter",
  encode(message: MsgGasPriceVoter, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chainId !== BigInt(0)) {
      writer.uint32(16).int64(message.chainId);
    }
    if (message.price !== BigInt(0)) {
      writer.uint32(24).uint64(message.price);
    }
    if (message.blockNumber !== BigInt(0)) {
      writer.uint32(32).uint64(message.blockNumber);
    }
    if (message.supply !== "") {
      writer.uint32(42).string(message.supply);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgGasPriceVoter {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgGasPriceVoter();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chainId = reader.int64();
          break;
        case 3:
          message.price = reader.uint64();
          break;
        case 4:
          message.blockNumber = reader.uint64();
          break;
        case 5:
          message.supply = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgGasPriceVoter>): MsgGasPriceVoter {
    const message = createBaseMsgGasPriceVoter();
    message.creator = object.creator ?? "";
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.price = object.price !== undefined && object.price !== null ? BigInt(object.price.toString()) : BigInt(0);
    message.blockNumber = object.blockNumber !== undefined && object.blockNumber !== null ? BigInt(object.blockNumber.toString()) : BigInt(0);
    message.supply = object.supply ?? "";
    return message;
  },
  fromAmino(object: MsgGasPriceVoterAmino): MsgGasPriceVoter {
    const message = createBaseMsgGasPriceVoter();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.price !== undefined && object.price !== null) {
      message.price = BigInt(object.price);
    }
    if (object.block_number !== undefined && object.block_number !== null) {
      message.blockNumber = BigInt(object.block_number);
    }
    if (object.supply !== undefined && object.supply !== null) {
      message.supply = object.supply;
    }
    return message;
  },
  toAmino(message: MsgGasPriceVoter): MsgGasPriceVoterAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.price = message.price ? message.price.toString() : undefined;
    obj.block_number = message.blockNumber ? message.blockNumber.toString() : undefined;
    obj.supply = message.supply;
    return obj;
  },
  fromAminoMsg(object: MsgGasPriceVoterAminoMsg): MsgGasPriceVoter {
    return MsgGasPriceVoter.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgGasPriceVoterProtoMsg): MsgGasPriceVoter {
    return MsgGasPriceVoter.decode(message.value);
  },
  toProto(message: MsgGasPriceVoter): Uint8Array {
    return MsgGasPriceVoter.encode(message).finish();
  },
  toProtoMsg(message: MsgGasPriceVoter): MsgGasPriceVoterProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgGasPriceVoter",
      value: MsgGasPriceVoter.encode(message).finish()
    };
  }
};
function createBaseMsgGasPriceVoterResponse(): MsgGasPriceVoterResponse {
  return {};
}
export const MsgGasPriceVoterResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgGasPriceVoterResponse",
  encode(_: MsgGasPriceVoterResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgGasPriceVoterResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgGasPriceVoterResponse();
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
  fromPartial(_: Partial<MsgGasPriceVoterResponse>): MsgGasPriceVoterResponse {
    const message = createBaseMsgGasPriceVoterResponse();
    return message;
  },
  fromAmino(_: MsgGasPriceVoterResponseAmino): MsgGasPriceVoterResponse {
    const message = createBaseMsgGasPriceVoterResponse();
    return message;
  },
  toAmino(_: MsgGasPriceVoterResponse): MsgGasPriceVoterResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgGasPriceVoterResponseAminoMsg): MsgGasPriceVoterResponse {
    return MsgGasPriceVoterResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgGasPriceVoterResponseProtoMsg): MsgGasPriceVoterResponse {
    return MsgGasPriceVoterResponse.decode(message.value);
  },
  toProto(message: MsgGasPriceVoterResponse): Uint8Array {
    return MsgGasPriceVoterResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgGasPriceVoterResponse): MsgGasPriceVoterResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgGasPriceVoterResponse",
      value: MsgGasPriceVoterResponse.encode(message).finish()
    };
  }
};
function createBaseMsgVoteOnObservedOutboundTx(): MsgVoteOnObservedOutboundTx {
  return {
    creator: "",
    cctxHash: "",
    observedOutTxHash: "",
    observedOutTxBlockHeight: BigInt(0),
    observedOutTxGasUsed: BigInt(0),
    observedOutTxEffectiveGasPrice: "",
    observedOutTxEffectiveGasLimit: BigInt(0),
    valueReceived: "",
    status: 0,
    outTxChain: BigInt(0),
    outTxTssNonce: BigInt(0),
    coinType: 0
  };
}
export const MsgVoteOnObservedOutboundTx = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx",
  encode(message: MsgVoteOnObservedOutboundTx, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.cctxHash !== "") {
      writer.uint32(18).string(message.cctxHash);
    }
    if (message.observedOutTxHash !== "") {
      writer.uint32(26).string(message.observedOutTxHash);
    }
    if (message.observedOutTxBlockHeight !== BigInt(0)) {
      writer.uint32(32).uint64(message.observedOutTxBlockHeight);
    }
    if (message.observedOutTxGasUsed !== BigInt(0)) {
      writer.uint32(80).uint64(message.observedOutTxGasUsed);
    }
    if (message.observedOutTxEffectiveGasPrice !== "") {
      writer.uint32(90).string(message.observedOutTxEffectiveGasPrice);
    }
    if (message.observedOutTxEffectiveGasLimit !== BigInt(0)) {
      writer.uint32(96).uint64(message.observedOutTxEffectiveGasLimit);
    }
    if (message.valueReceived !== "") {
      writer.uint32(42).string(message.valueReceived);
    }
    if (message.status !== 0) {
      writer.uint32(48).int32(message.status);
    }
    if (message.outTxChain !== BigInt(0)) {
      writer.uint32(56).int64(message.outTxChain);
    }
    if (message.outTxTssNonce !== BigInt(0)) {
      writer.uint32(64).uint64(message.outTxTssNonce);
    }
    if (message.coinType !== 0) {
      writer.uint32(72).int32(message.coinType);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgVoteOnObservedOutboundTx {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgVoteOnObservedOutboundTx();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.cctxHash = reader.string();
          break;
        case 3:
          message.observedOutTxHash = reader.string();
          break;
        case 4:
          message.observedOutTxBlockHeight = reader.uint64();
          break;
        case 10:
          message.observedOutTxGasUsed = reader.uint64();
          break;
        case 11:
          message.observedOutTxEffectiveGasPrice = reader.string();
          break;
        case 12:
          message.observedOutTxEffectiveGasLimit = reader.uint64();
          break;
        case 5:
          message.valueReceived = reader.string();
          break;
        case 6:
          message.status = (reader.int32() as any);
          break;
        case 7:
          message.outTxChain = reader.int64();
          break;
        case 8:
          message.outTxTssNonce = reader.uint64();
          break;
        case 9:
          message.coinType = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgVoteOnObservedOutboundTx>): MsgVoteOnObservedOutboundTx {
    const message = createBaseMsgVoteOnObservedOutboundTx();
    message.creator = object.creator ?? "";
    message.cctxHash = object.cctxHash ?? "";
    message.observedOutTxHash = object.observedOutTxHash ?? "";
    message.observedOutTxBlockHeight = object.observedOutTxBlockHeight !== undefined && object.observedOutTxBlockHeight !== null ? BigInt(object.observedOutTxBlockHeight.toString()) : BigInt(0);
    message.observedOutTxGasUsed = object.observedOutTxGasUsed !== undefined && object.observedOutTxGasUsed !== null ? BigInt(object.observedOutTxGasUsed.toString()) : BigInt(0);
    message.observedOutTxEffectiveGasPrice = object.observedOutTxEffectiveGasPrice ?? "";
    message.observedOutTxEffectiveGasLimit = object.observedOutTxEffectiveGasLimit !== undefined && object.observedOutTxEffectiveGasLimit !== null ? BigInt(object.observedOutTxEffectiveGasLimit.toString()) : BigInt(0);
    message.valueReceived = object.valueReceived ?? "";
    message.status = object.status ?? 0;
    message.outTxChain = object.outTxChain !== undefined && object.outTxChain !== null ? BigInt(object.outTxChain.toString()) : BigInt(0);
    message.outTxTssNonce = object.outTxTssNonce !== undefined && object.outTxTssNonce !== null ? BigInt(object.outTxTssNonce.toString()) : BigInt(0);
    message.coinType = object.coinType ?? 0;
    return message;
  },
  fromAmino(object: MsgVoteOnObservedOutboundTxAmino): MsgVoteOnObservedOutboundTx {
    const message = createBaseMsgVoteOnObservedOutboundTx();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.cctx_hash !== undefined && object.cctx_hash !== null) {
      message.cctxHash = object.cctx_hash;
    }
    if (object.observed_outTx_hash !== undefined && object.observed_outTx_hash !== null) {
      message.observedOutTxHash = object.observed_outTx_hash;
    }
    if (object.observed_outTx_blockHeight !== undefined && object.observed_outTx_blockHeight !== null) {
      message.observedOutTxBlockHeight = BigInt(object.observed_outTx_blockHeight);
    }
    if (object.observed_outTx_gas_used !== undefined && object.observed_outTx_gas_used !== null) {
      message.observedOutTxGasUsed = BigInt(object.observed_outTx_gas_used);
    }
    if (object.observed_outTx_effective_gas_price !== undefined && object.observed_outTx_effective_gas_price !== null) {
      message.observedOutTxEffectiveGasPrice = object.observed_outTx_effective_gas_price;
    }
    if (object.observed_outTx_effective_gas_limit !== undefined && object.observed_outTx_effective_gas_limit !== null) {
      message.observedOutTxEffectiveGasLimit = BigInt(object.observed_outTx_effective_gas_limit);
    }
    if (object.value_received !== undefined && object.value_received !== null) {
      message.valueReceived = object.value_received;
    }
    if (object.status !== undefined && object.status !== null) {
      message.status = receiveStatusFromJSON(object.status);
    }
    if (object.outTx_chain !== undefined && object.outTx_chain !== null) {
      message.outTxChain = BigInt(object.outTx_chain);
    }
    if (object.outTx_tss_nonce !== undefined && object.outTx_tss_nonce !== null) {
      message.outTxTssNonce = BigInt(object.outTx_tss_nonce);
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    return message;
  },
  toAmino(message: MsgVoteOnObservedOutboundTx): MsgVoteOnObservedOutboundTxAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.cctx_hash = message.cctxHash;
    obj.observed_outTx_hash = message.observedOutTxHash;
    obj.observed_outTx_blockHeight = message.observedOutTxBlockHeight ? message.observedOutTxBlockHeight.toString() : undefined;
    obj.observed_outTx_gas_used = message.observedOutTxGasUsed ? message.observedOutTxGasUsed.toString() : undefined;
    obj.observed_outTx_effective_gas_price = message.observedOutTxEffectiveGasPrice;
    obj.observed_outTx_effective_gas_limit = message.observedOutTxEffectiveGasLimit ? message.observedOutTxEffectiveGasLimit.toString() : undefined;
    obj.value_received = message.valueReceived;
    obj.status = message.status;
    obj.outTx_chain = message.outTxChain ? message.outTxChain.toString() : undefined;
    obj.outTx_tss_nonce = message.outTxTssNonce ? message.outTxTssNonce.toString() : undefined;
    obj.coin_type = message.coinType;
    return obj;
  },
  fromAminoMsg(object: MsgVoteOnObservedOutboundTxAminoMsg): MsgVoteOnObservedOutboundTx {
    return MsgVoteOnObservedOutboundTx.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgVoteOnObservedOutboundTxProtoMsg): MsgVoteOnObservedOutboundTx {
    return MsgVoteOnObservedOutboundTx.decode(message.value);
  },
  toProto(message: MsgVoteOnObservedOutboundTx): Uint8Array {
    return MsgVoteOnObservedOutboundTx.encode(message).finish();
  },
  toProtoMsg(message: MsgVoteOnObservedOutboundTx): MsgVoteOnObservedOutboundTxProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx",
      value: MsgVoteOnObservedOutboundTx.encode(message).finish()
    };
  }
};
function createBaseMsgVoteOnObservedOutboundTxResponse(): MsgVoteOnObservedOutboundTxResponse {
  return {};
}
export const MsgVoteOnObservedOutboundTxResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTxResponse",
  encode(_: MsgVoteOnObservedOutboundTxResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgVoteOnObservedOutboundTxResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgVoteOnObservedOutboundTxResponse();
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
  fromPartial(_: Partial<MsgVoteOnObservedOutboundTxResponse>): MsgVoteOnObservedOutboundTxResponse {
    const message = createBaseMsgVoteOnObservedOutboundTxResponse();
    return message;
  },
  fromAmino(_: MsgVoteOnObservedOutboundTxResponseAmino): MsgVoteOnObservedOutboundTxResponse {
    const message = createBaseMsgVoteOnObservedOutboundTxResponse();
    return message;
  },
  toAmino(_: MsgVoteOnObservedOutboundTxResponse): MsgVoteOnObservedOutboundTxResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgVoteOnObservedOutboundTxResponseAminoMsg): MsgVoteOnObservedOutboundTxResponse {
    return MsgVoteOnObservedOutboundTxResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgVoteOnObservedOutboundTxResponseProtoMsg): MsgVoteOnObservedOutboundTxResponse {
    return MsgVoteOnObservedOutboundTxResponse.decode(message.value);
  },
  toProto(message: MsgVoteOnObservedOutboundTxResponse): Uint8Array {
    return MsgVoteOnObservedOutboundTxResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgVoteOnObservedOutboundTxResponse): MsgVoteOnObservedOutboundTxResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTxResponse",
      value: MsgVoteOnObservedOutboundTxResponse.encode(message).finish()
    };
  }
};
function createBaseMsgVoteOnObservedInboundTx(): MsgVoteOnObservedInboundTx {
  return {
    creator: "",
    sender: "",
    senderChainId: BigInt(0),
    receiver: "",
    receiverChain: BigInt(0),
    amount: "",
    message: "",
    inTxHash: "",
    inBlockHeight: BigInt(0),
    gasLimit: BigInt(0),
    coinType: 0,
    txOrigin: "",
    asset: "",
    eventIndex: BigInt(0)
  };
}
export const MsgVoteOnObservedInboundTx = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx",
  encode(message: MsgVoteOnObservedInboundTx, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.sender !== "") {
      writer.uint32(18).string(message.sender);
    }
    if (message.senderChainId !== BigInt(0)) {
      writer.uint32(24).int64(message.senderChainId);
    }
    if (message.receiver !== "") {
      writer.uint32(34).string(message.receiver);
    }
    if (message.receiverChain !== BigInt(0)) {
      writer.uint32(40).int64(message.receiverChain);
    }
    if (message.amount !== "") {
      writer.uint32(50).string(message.amount);
    }
    if (message.message !== "") {
      writer.uint32(66).string(message.message);
    }
    if (message.inTxHash !== "") {
      writer.uint32(74).string(message.inTxHash);
    }
    if (message.inBlockHeight !== BigInt(0)) {
      writer.uint32(80).uint64(message.inBlockHeight);
    }
    if (message.gasLimit !== BigInt(0)) {
      writer.uint32(88).uint64(message.gasLimit);
    }
    if (message.coinType !== 0) {
      writer.uint32(96).int32(message.coinType);
    }
    if (message.txOrigin !== "") {
      writer.uint32(106).string(message.txOrigin);
    }
    if (message.asset !== "") {
      writer.uint32(114).string(message.asset);
    }
    if (message.eventIndex !== BigInt(0)) {
      writer.uint32(120).uint64(message.eventIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgVoteOnObservedInboundTx {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgVoteOnObservedInboundTx();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.sender = reader.string();
          break;
        case 3:
          message.senderChainId = reader.int64();
          break;
        case 4:
          message.receiver = reader.string();
          break;
        case 5:
          message.receiverChain = reader.int64();
          break;
        case 6:
          message.amount = reader.string();
          break;
        case 8:
          message.message = reader.string();
          break;
        case 9:
          message.inTxHash = reader.string();
          break;
        case 10:
          message.inBlockHeight = reader.uint64();
          break;
        case 11:
          message.gasLimit = reader.uint64();
          break;
        case 12:
          message.coinType = (reader.int32() as any);
          break;
        case 13:
          message.txOrigin = reader.string();
          break;
        case 14:
          message.asset = reader.string();
          break;
        case 15:
          message.eventIndex = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgVoteOnObservedInboundTx>): MsgVoteOnObservedInboundTx {
    const message = createBaseMsgVoteOnObservedInboundTx();
    message.creator = object.creator ?? "";
    message.sender = object.sender ?? "";
    message.senderChainId = object.senderChainId !== undefined && object.senderChainId !== null ? BigInt(object.senderChainId.toString()) : BigInt(0);
    message.receiver = object.receiver ?? "";
    message.receiverChain = object.receiverChain !== undefined && object.receiverChain !== null ? BigInt(object.receiverChain.toString()) : BigInt(0);
    message.amount = object.amount ?? "";
    message.message = object.message ?? "";
    message.inTxHash = object.inTxHash ?? "";
    message.inBlockHeight = object.inBlockHeight !== undefined && object.inBlockHeight !== null ? BigInt(object.inBlockHeight.toString()) : BigInt(0);
    message.gasLimit = object.gasLimit !== undefined && object.gasLimit !== null ? BigInt(object.gasLimit.toString()) : BigInt(0);
    message.coinType = object.coinType ?? 0;
    message.txOrigin = object.txOrigin ?? "";
    message.asset = object.asset ?? "";
    message.eventIndex = object.eventIndex !== undefined && object.eventIndex !== null ? BigInt(object.eventIndex.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: MsgVoteOnObservedInboundTxAmino): MsgVoteOnObservedInboundTx {
    const message = createBaseMsgVoteOnObservedInboundTx();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.sender !== undefined && object.sender !== null) {
      message.sender = object.sender;
    }
    if (object.sender_chain_id !== undefined && object.sender_chain_id !== null) {
      message.senderChainId = BigInt(object.sender_chain_id);
    }
    if (object.receiver !== undefined && object.receiver !== null) {
      message.receiver = object.receiver;
    }
    if (object.receiver_chain !== undefined && object.receiver_chain !== null) {
      message.receiverChain = BigInt(object.receiver_chain);
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    }
    if (object.message !== undefined && object.message !== null) {
      message.message = object.message;
    }
    if (object.in_tx_hash !== undefined && object.in_tx_hash !== null) {
      message.inTxHash = object.in_tx_hash;
    }
    if (object.in_block_height !== undefined && object.in_block_height !== null) {
      message.inBlockHeight = BigInt(object.in_block_height);
    }
    if (object.gas_limit !== undefined && object.gas_limit !== null) {
      message.gasLimit = BigInt(object.gas_limit);
    }
    if (object.coin_type !== undefined && object.coin_type !== null) {
      message.coinType = coinTypeFromJSON(object.coin_type);
    }
    if (object.tx_origin !== undefined && object.tx_origin !== null) {
      message.txOrigin = object.tx_origin;
    }
    if (object.asset !== undefined && object.asset !== null) {
      message.asset = object.asset;
    }
    if (object.event_index !== undefined && object.event_index !== null) {
      message.eventIndex = BigInt(object.event_index);
    }
    return message;
  },
  toAmino(message: MsgVoteOnObservedInboundTx): MsgVoteOnObservedInboundTxAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.sender = message.sender;
    obj.sender_chain_id = message.senderChainId ? message.senderChainId.toString() : undefined;
    obj.receiver = message.receiver;
    obj.receiver_chain = message.receiverChain ? message.receiverChain.toString() : undefined;
    obj.amount = message.amount;
    obj.message = message.message;
    obj.in_tx_hash = message.inTxHash;
    obj.in_block_height = message.inBlockHeight ? message.inBlockHeight.toString() : undefined;
    obj.gas_limit = message.gasLimit ? message.gasLimit.toString() : undefined;
    obj.coin_type = message.coinType;
    obj.tx_origin = message.txOrigin;
    obj.asset = message.asset;
    obj.event_index = message.eventIndex ? message.eventIndex.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgVoteOnObservedInboundTxAminoMsg): MsgVoteOnObservedInboundTx {
    return MsgVoteOnObservedInboundTx.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgVoteOnObservedInboundTxProtoMsg): MsgVoteOnObservedInboundTx {
    return MsgVoteOnObservedInboundTx.decode(message.value);
  },
  toProto(message: MsgVoteOnObservedInboundTx): Uint8Array {
    return MsgVoteOnObservedInboundTx.encode(message).finish();
  },
  toProtoMsg(message: MsgVoteOnObservedInboundTx): MsgVoteOnObservedInboundTxProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx",
      value: MsgVoteOnObservedInboundTx.encode(message).finish()
    };
  }
};
function createBaseMsgVoteOnObservedInboundTxResponse(): MsgVoteOnObservedInboundTxResponse {
  return {};
}
export const MsgVoteOnObservedInboundTxResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTxResponse",
  encode(_: MsgVoteOnObservedInboundTxResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgVoteOnObservedInboundTxResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgVoteOnObservedInboundTxResponse();
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
  fromPartial(_: Partial<MsgVoteOnObservedInboundTxResponse>): MsgVoteOnObservedInboundTxResponse {
    const message = createBaseMsgVoteOnObservedInboundTxResponse();
    return message;
  },
  fromAmino(_: MsgVoteOnObservedInboundTxResponseAmino): MsgVoteOnObservedInboundTxResponse {
    const message = createBaseMsgVoteOnObservedInboundTxResponse();
    return message;
  },
  toAmino(_: MsgVoteOnObservedInboundTxResponse): MsgVoteOnObservedInboundTxResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgVoteOnObservedInboundTxResponseAminoMsg): MsgVoteOnObservedInboundTxResponse {
    return MsgVoteOnObservedInboundTxResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgVoteOnObservedInboundTxResponseProtoMsg): MsgVoteOnObservedInboundTxResponse {
    return MsgVoteOnObservedInboundTxResponse.decode(message.value);
  },
  toProto(message: MsgVoteOnObservedInboundTxResponse): Uint8Array {
    return MsgVoteOnObservedInboundTxResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgVoteOnObservedInboundTxResponse): MsgVoteOnObservedInboundTxResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTxResponse",
      value: MsgVoteOnObservedInboundTxResponse.encode(message).finish()
    };
  }
};
function createBaseMsgAbortStuckCCTX(): MsgAbortStuckCCTX {
  return {
    creator: "",
    cctxIndex: ""
  };
}
export const MsgAbortStuckCCTX = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX",
  encode(message: MsgAbortStuckCCTX, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.cctxIndex !== "") {
      writer.uint32(18).string(message.cctxIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAbortStuckCCTX {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAbortStuckCCTX();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.cctxIndex = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgAbortStuckCCTX>): MsgAbortStuckCCTX {
    const message = createBaseMsgAbortStuckCCTX();
    message.creator = object.creator ?? "";
    message.cctxIndex = object.cctxIndex ?? "";
    return message;
  },
  fromAmino(object: MsgAbortStuckCCTXAmino): MsgAbortStuckCCTX {
    const message = createBaseMsgAbortStuckCCTX();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.cctx_index !== undefined && object.cctx_index !== null) {
      message.cctxIndex = object.cctx_index;
    }
    return message;
  },
  toAmino(message: MsgAbortStuckCCTX): MsgAbortStuckCCTXAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.cctx_index = message.cctxIndex;
    return obj;
  },
  fromAminoMsg(object: MsgAbortStuckCCTXAminoMsg): MsgAbortStuckCCTX {
    return MsgAbortStuckCCTX.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAbortStuckCCTXProtoMsg): MsgAbortStuckCCTX {
    return MsgAbortStuckCCTX.decode(message.value);
  },
  toProto(message: MsgAbortStuckCCTX): Uint8Array {
    return MsgAbortStuckCCTX.encode(message).finish();
  },
  toProtoMsg(message: MsgAbortStuckCCTX): MsgAbortStuckCCTXProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX",
      value: MsgAbortStuckCCTX.encode(message).finish()
    };
  }
};
function createBaseMsgAbortStuckCCTXResponse(): MsgAbortStuckCCTXResponse {
  return {};
}
export const MsgAbortStuckCCTXResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTXResponse",
  encode(_: MsgAbortStuckCCTXResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAbortStuckCCTXResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAbortStuckCCTXResponse();
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
  fromPartial(_: Partial<MsgAbortStuckCCTXResponse>): MsgAbortStuckCCTXResponse {
    const message = createBaseMsgAbortStuckCCTXResponse();
    return message;
  },
  fromAmino(_: MsgAbortStuckCCTXResponseAmino): MsgAbortStuckCCTXResponse {
    const message = createBaseMsgAbortStuckCCTXResponse();
    return message;
  },
  toAmino(_: MsgAbortStuckCCTXResponse): MsgAbortStuckCCTXResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgAbortStuckCCTXResponseAminoMsg): MsgAbortStuckCCTXResponse {
    return MsgAbortStuckCCTXResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAbortStuckCCTXResponseProtoMsg): MsgAbortStuckCCTXResponse {
    return MsgAbortStuckCCTXResponse.decode(message.value);
  },
  toProto(message: MsgAbortStuckCCTXResponse): Uint8Array {
    return MsgAbortStuckCCTXResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgAbortStuckCCTXResponse): MsgAbortStuckCCTXResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTXResponse",
      value: MsgAbortStuckCCTXResponse.encode(message).finish()
    };
  }
};
function createBaseMsgRefundAbortedCCTX(): MsgRefundAbortedCCTX {
  return {
    creator: "",
    cctxIndex: "",
    refundAddress: ""
  };
}
export const MsgRefundAbortedCCTX = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX",
  encode(message: MsgRefundAbortedCCTX, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.cctxIndex !== "") {
      writer.uint32(18).string(message.cctxIndex);
    }
    if (message.refundAddress !== "") {
      writer.uint32(26).string(message.refundAddress);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgRefundAbortedCCTX {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRefundAbortedCCTX();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.cctxIndex = reader.string();
          break;
        case 3:
          message.refundAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgRefundAbortedCCTX>): MsgRefundAbortedCCTX {
    const message = createBaseMsgRefundAbortedCCTX();
    message.creator = object.creator ?? "";
    message.cctxIndex = object.cctxIndex ?? "";
    message.refundAddress = object.refundAddress ?? "";
    return message;
  },
  fromAmino(object: MsgRefundAbortedCCTXAmino): MsgRefundAbortedCCTX {
    const message = createBaseMsgRefundAbortedCCTX();
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    }
    if (object.cctx_index !== undefined && object.cctx_index !== null) {
      message.cctxIndex = object.cctx_index;
    }
    if (object.refund_address !== undefined && object.refund_address !== null) {
      message.refundAddress = object.refund_address;
    }
    return message;
  },
  toAmino(message: MsgRefundAbortedCCTX): MsgRefundAbortedCCTXAmino {
    const obj: any = {};
    obj.creator = message.creator;
    obj.cctx_index = message.cctxIndex;
    obj.refund_address = message.refundAddress;
    return obj;
  },
  fromAminoMsg(object: MsgRefundAbortedCCTXAminoMsg): MsgRefundAbortedCCTX {
    return MsgRefundAbortedCCTX.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRefundAbortedCCTXProtoMsg): MsgRefundAbortedCCTX {
    return MsgRefundAbortedCCTX.decode(message.value);
  },
  toProto(message: MsgRefundAbortedCCTX): Uint8Array {
    return MsgRefundAbortedCCTX.encode(message).finish();
  },
  toProtoMsg(message: MsgRefundAbortedCCTX): MsgRefundAbortedCCTXProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX",
      value: MsgRefundAbortedCCTX.encode(message).finish()
    };
  }
};
function createBaseMsgRefundAbortedCCTXResponse(): MsgRefundAbortedCCTXResponse {
  return {};
}
export const MsgRefundAbortedCCTXResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTXResponse",
  encode(_: MsgRefundAbortedCCTXResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgRefundAbortedCCTXResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgRefundAbortedCCTXResponse();
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
  fromPartial(_: Partial<MsgRefundAbortedCCTXResponse>): MsgRefundAbortedCCTXResponse {
    const message = createBaseMsgRefundAbortedCCTXResponse();
    return message;
  },
  fromAmino(_: MsgRefundAbortedCCTXResponseAmino): MsgRefundAbortedCCTXResponse {
    const message = createBaseMsgRefundAbortedCCTXResponse();
    return message;
  },
  toAmino(_: MsgRefundAbortedCCTXResponse): MsgRefundAbortedCCTXResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgRefundAbortedCCTXResponseAminoMsg): MsgRefundAbortedCCTXResponse {
    return MsgRefundAbortedCCTXResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgRefundAbortedCCTXResponseProtoMsg): MsgRefundAbortedCCTXResponse {
    return MsgRefundAbortedCCTXResponse.decode(message.value);
  },
  toProto(message: MsgRefundAbortedCCTXResponse): Uint8Array {
    return MsgRefundAbortedCCTXResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgRefundAbortedCCTXResponse): MsgRefundAbortedCCTXResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTXResponse",
      value: MsgRefundAbortedCCTXResponse.encode(message).finish()
    };
  }
};