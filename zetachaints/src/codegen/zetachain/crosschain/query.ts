import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Params, ParamsAmino, ParamsSDKType } from "./params";
import { OutTxTracker, OutTxTrackerAmino, OutTxTrackerSDKType } from "./out_tx_tracker";
import { InTxTracker, InTxTrackerAmino, InTxTrackerSDKType } from "./in_tx_tracker";
import { InTxHashToCctx, InTxHashToCctxAmino, InTxHashToCctxSDKType } from "./in_tx_hash_to_cctx";
import { CrossChainTx, CrossChainTxAmino, CrossChainTxSDKType } from "./cross_chain_tx";
import { GasPrice, GasPriceAmino, GasPriceSDKType } from "./gas_price";
import { LastBlockHeight, LastBlockHeightAmino, LastBlockHeightSDKType } from "./last_block_height";
import { BinaryReader, BinaryWriter } from "../../binary";
/**
 * Deprecated: Moved to observer
 * TODO: remove after v12 once upgrade testing is no longer needed with v11
 * https://github.com/zeta-chain/node/issues/1547
 */
export interface QueryGetTssAddressRequest {}
export interface QueryGetTssAddressRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetTssAddressRequest";
  value: Uint8Array;
}
/**
 * Deprecated: Moved to observer
 * TODO: remove after v12 once upgrade testing is no longer needed with v11
 * https://github.com/zeta-chain/node/issues/1547
 */
export interface QueryGetTssAddressRequestAmino {}
export interface QueryGetTssAddressRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetTssAddressRequest";
  value: QueryGetTssAddressRequestAmino;
}
/**
 * Deprecated: Moved to observer
 * TODO: remove after v12 once upgrade testing is no longer needed with v11
 * https://github.com/zeta-chain/node/issues/1547
 */
export interface QueryGetTssAddressRequestSDKType {}
/**
 * Deprecated: Moved to observer
 * TODO: remove after v12 once upgrade testing is no longer needed with v11
 * https://github.com/zeta-chain/node/issues/1547
 */
export interface QueryGetTssAddressResponse {
  eth: string;
  btc: string;
}
export interface QueryGetTssAddressResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetTssAddressResponse";
  value: Uint8Array;
}
/**
 * Deprecated: Moved to observer
 * TODO: remove after v12 once upgrade testing is no longer needed with v11
 * https://github.com/zeta-chain/node/issues/1547
 */
export interface QueryGetTssAddressResponseAmino {
  eth?: string;
  btc?: string;
}
export interface QueryGetTssAddressResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetTssAddressResponse";
  value: QueryGetTssAddressResponseAmino;
}
/**
 * Deprecated: Moved to observer
 * TODO: remove after v12 once upgrade testing is no longer needed with v11
 * https://github.com/zeta-chain/node/issues/1547
 */
export interface QueryGetTssAddressResponseSDKType {
  eth: string;
  btc: string;
}
export interface QueryZetaAccountingRequest {}
export interface QueryZetaAccountingRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryZetaAccountingRequest";
  value: Uint8Array;
}
export interface QueryZetaAccountingRequestAmino {}
export interface QueryZetaAccountingRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryZetaAccountingRequest";
  value: QueryZetaAccountingRequestAmino;
}
export interface QueryZetaAccountingRequestSDKType {}
export interface QueryZetaAccountingResponse {
  abortedZetaAmount: string;
}
export interface QueryZetaAccountingResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryZetaAccountingResponse";
  value: Uint8Array;
}
export interface QueryZetaAccountingResponseAmino {
  aborted_zeta_amount?: string;
}
export interface QueryZetaAccountingResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryZetaAccountingResponse";
  value: QueryZetaAccountingResponseAmino;
}
export interface QueryZetaAccountingResponseSDKType {
  aborted_zeta_amount: string;
}
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {}
export interface QueryParamsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryParamsRequest";
  value: Uint8Array;
}
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequestAmino {}
export interface QueryParamsRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryParamsRequest";
  value: QueryParamsRequestAmino;
}
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params;
}
export interface QueryParamsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryParamsResponse";
  value: Uint8Array;
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponseAmino {
  /** params holds all the parameters of this module. */
  params?: ParamsAmino;
}
export interface QueryParamsResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryParamsResponse";
  value: QueryParamsResponseAmino;
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponseSDKType {
  params: ParamsSDKType;
}
export interface QueryGetOutTxTrackerRequest {
  chainID: bigint;
  nonce: bigint;
}
export interface QueryGetOutTxTrackerRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetOutTxTrackerRequest";
  value: Uint8Array;
}
export interface QueryGetOutTxTrackerRequestAmino {
  chainID?: string;
  nonce?: string;
}
export interface QueryGetOutTxTrackerRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetOutTxTrackerRequest";
  value: QueryGetOutTxTrackerRequestAmino;
}
export interface QueryGetOutTxTrackerRequestSDKType {
  chainID: bigint;
  nonce: bigint;
}
export interface QueryGetOutTxTrackerResponse {
  outTxTracker: OutTxTracker;
}
export interface QueryGetOutTxTrackerResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetOutTxTrackerResponse";
  value: Uint8Array;
}
export interface QueryGetOutTxTrackerResponseAmino {
  outTxTracker?: OutTxTrackerAmino;
}
export interface QueryGetOutTxTrackerResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetOutTxTrackerResponse";
  value: QueryGetOutTxTrackerResponseAmino;
}
export interface QueryGetOutTxTrackerResponseSDKType {
  outTxTracker: OutTxTrackerSDKType;
}
export interface QueryAllOutTxTrackerRequest {
  pagination?: PageRequest;
}
export interface QueryAllOutTxTrackerRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerRequest";
  value: Uint8Array;
}
export interface QueryAllOutTxTrackerRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllOutTxTrackerRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerRequest";
  value: QueryAllOutTxTrackerRequestAmino;
}
export interface QueryAllOutTxTrackerRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllOutTxTrackerResponse {
  outTxTracker: OutTxTracker[];
  pagination?: PageResponse;
}
export interface QueryAllOutTxTrackerResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerResponse";
  value: Uint8Array;
}
export interface QueryAllOutTxTrackerResponseAmino {
  outTxTracker?: OutTxTrackerAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllOutTxTrackerResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerResponse";
  value: QueryAllOutTxTrackerResponseAmino;
}
export interface QueryAllOutTxTrackerResponseSDKType {
  outTxTracker: OutTxTrackerSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryAllOutTxTrackerByChainRequest {
  chain: bigint;
  pagination?: PageRequest;
}
export interface QueryAllOutTxTrackerByChainRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerByChainRequest";
  value: Uint8Array;
}
export interface QueryAllOutTxTrackerByChainRequestAmino {
  chain?: string;
  pagination?: PageRequestAmino;
}
export interface QueryAllOutTxTrackerByChainRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerByChainRequest";
  value: QueryAllOutTxTrackerByChainRequestAmino;
}
export interface QueryAllOutTxTrackerByChainRequestSDKType {
  chain: bigint;
  pagination?: PageRequestSDKType;
}
export interface QueryAllOutTxTrackerByChainResponse {
  outTxTracker: OutTxTracker[];
  pagination?: PageResponse;
}
export interface QueryAllOutTxTrackerByChainResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerByChainResponse";
  value: Uint8Array;
}
export interface QueryAllOutTxTrackerByChainResponseAmino {
  outTxTracker?: OutTxTrackerAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllOutTxTrackerByChainResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerByChainResponse";
  value: QueryAllOutTxTrackerByChainResponseAmino;
}
export interface QueryAllOutTxTrackerByChainResponseSDKType {
  outTxTracker: OutTxTrackerSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryAllInTxTrackerByChainRequest {
  chainId: bigint;
  pagination?: PageRequest;
}
export interface QueryAllInTxTrackerByChainRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackerByChainRequest";
  value: Uint8Array;
}
export interface QueryAllInTxTrackerByChainRequestAmino {
  chain_id?: string;
  pagination?: PageRequestAmino;
}
export interface QueryAllInTxTrackerByChainRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllInTxTrackerByChainRequest";
  value: QueryAllInTxTrackerByChainRequestAmino;
}
export interface QueryAllInTxTrackerByChainRequestSDKType {
  chain_id: bigint;
  pagination?: PageRequestSDKType;
}
export interface QueryAllInTxTrackerByChainResponse {
  inTxTracker: InTxTracker[];
  pagination?: PageResponse;
}
export interface QueryAllInTxTrackerByChainResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackerByChainResponse";
  value: Uint8Array;
}
export interface QueryAllInTxTrackerByChainResponseAmino {
  inTxTracker?: InTxTrackerAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllInTxTrackerByChainResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllInTxTrackerByChainResponse";
  value: QueryAllInTxTrackerByChainResponseAmino;
}
export interface QueryAllInTxTrackerByChainResponseSDKType {
  inTxTracker: InTxTrackerSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryAllInTxTrackersRequest {
  pagination?: PageRequest;
}
export interface QueryAllInTxTrackersRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackersRequest";
  value: Uint8Array;
}
export interface QueryAllInTxTrackersRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllInTxTrackersRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllInTxTrackersRequest";
  value: QueryAllInTxTrackersRequestAmino;
}
export interface QueryAllInTxTrackersRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllInTxTrackersResponse {
  inTxTracker: InTxTracker[];
  pagination?: PageResponse;
}
export interface QueryAllInTxTrackersResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackersResponse";
  value: Uint8Array;
}
export interface QueryAllInTxTrackersResponseAmino {
  inTxTracker?: InTxTrackerAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllInTxTrackersResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllInTxTrackersResponse";
  value: QueryAllInTxTrackersResponseAmino;
}
export interface QueryAllInTxTrackersResponseSDKType {
  inTxTracker: InTxTrackerSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryGetInTxHashToCctxRequest {
  inTxHash: string;
}
export interface QueryGetInTxHashToCctxRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetInTxHashToCctxRequest";
  value: Uint8Array;
}
export interface QueryGetInTxHashToCctxRequestAmino {
  inTxHash?: string;
}
export interface QueryGetInTxHashToCctxRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetInTxHashToCctxRequest";
  value: QueryGetInTxHashToCctxRequestAmino;
}
export interface QueryGetInTxHashToCctxRequestSDKType {
  inTxHash: string;
}
export interface QueryGetInTxHashToCctxResponse {
  inTxHashToCctx: InTxHashToCctx;
}
export interface QueryGetInTxHashToCctxResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetInTxHashToCctxResponse";
  value: Uint8Array;
}
export interface QueryGetInTxHashToCctxResponseAmino {
  inTxHashToCctx?: InTxHashToCctxAmino;
}
export interface QueryGetInTxHashToCctxResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetInTxHashToCctxResponse";
  value: QueryGetInTxHashToCctxResponseAmino;
}
export interface QueryGetInTxHashToCctxResponseSDKType {
  inTxHashToCctx: InTxHashToCctxSDKType;
}
export interface QueryInTxHashToCctxDataRequest {
  inTxHash: string;
}
export interface QueryInTxHashToCctxDataRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryInTxHashToCctxDataRequest";
  value: Uint8Array;
}
export interface QueryInTxHashToCctxDataRequestAmino {
  inTxHash?: string;
}
export interface QueryInTxHashToCctxDataRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryInTxHashToCctxDataRequest";
  value: QueryInTxHashToCctxDataRequestAmino;
}
export interface QueryInTxHashToCctxDataRequestSDKType {
  inTxHash: string;
}
export interface QueryInTxHashToCctxDataResponse {
  CrossChainTxs: CrossChainTx[];
}
export interface QueryInTxHashToCctxDataResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryInTxHashToCctxDataResponse";
  value: Uint8Array;
}
export interface QueryInTxHashToCctxDataResponseAmino {
  CrossChainTxs?: CrossChainTxAmino[];
}
export interface QueryInTxHashToCctxDataResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryInTxHashToCctxDataResponse";
  value: QueryInTxHashToCctxDataResponseAmino;
}
export interface QueryInTxHashToCctxDataResponseSDKType {
  CrossChainTxs: CrossChainTxSDKType[];
}
export interface QueryAllInTxHashToCctxRequest {
  pagination?: PageRequest;
}
export interface QueryAllInTxHashToCctxRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxHashToCctxRequest";
  value: Uint8Array;
}
export interface QueryAllInTxHashToCctxRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllInTxHashToCctxRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllInTxHashToCctxRequest";
  value: QueryAllInTxHashToCctxRequestAmino;
}
export interface QueryAllInTxHashToCctxRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllInTxHashToCctxResponse {
  inTxHashToCctx: InTxHashToCctx[];
  pagination?: PageResponse;
}
export interface QueryAllInTxHashToCctxResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxHashToCctxResponse";
  value: Uint8Array;
}
export interface QueryAllInTxHashToCctxResponseAmino {
  inTxHashToCctx?: InTxHashToCctxAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllInTxHashToCctxResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllInTxHashToCctxResponse";
  value: QueryAllInTxHashToCctxResponseAmino;
}
export interface QueryAllInTxHashToCctxResponseSDKType {
  inTxHashToCctx: InTxHashToCctxSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryGetGasPriceRequest {
  index: string;
}
export interface QueryGetGasPriceRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetGasPriceRequest";
  value: Uint8Array;
}
export interface QueryGetGasPriceRequestAmino {
  index?: string;
}
export interface QueryGetGasPriceRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetGasPriceRequest";
  value: QueryGetGasPriceRequestAmino;
}
export interface QueryGetGasPriceRequestSDKType {
  index: string;
}
export interface QueryGetGasPriceResponse {
  GasPrice?: GasPrice;
}
export interface QueryGetGasPriceResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetGasPriceResponse";
  value: Uint8Array;
}
export interface QueryGetGasPriceResponseAmino {
  GasPrice?: GasPriceAmino;
}
export interface QueryGetGasPriceResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetGasPriceResponse";
  value: QueryGetGasPriceResponseAmino;
}
export interface QueryGetGasPriceResponseSDKType {
  GasPrice?: GasPriceSDKType;
}
export interface QueryAllGasPriceRequest {
  pagination?: PageRequest;
}
export interface QueryAllGasPriceRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllGasPriceRequest";
  value: Uint8Array;
}
export interface QueryAllGasPriceRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllGasPriceRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllGasPriceRequest";
  value: QueryAllGasPriceRequestAmino;
}
export interface QueryAllGasPriceRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllGasPriceResponse {
  GasPrice: GasPrice[];
  pagination?: PageResponse;
}
export interface QueryAllGasPriceResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllGasPriceResponse";
  value: Uint8Array;
}
export interface QueryAllGasPriceResponseAmino {
  GasPrice?: GasPriceAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllGasPriceResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllGasPriceResponse";
  value: QueryAllGasPriceResponseAmino;
}
export interface QueryAllGasPriceResponseSDKType {
  GasPrice: GasPriceSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryGetLastBlockHeightRequest {
  index: string;
}
export interface QueryGetLastBlockHeightRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetLastBlockHeightRequest";
  value: Uint8Array;
}
export interface QueryGetLastBlockHeightRequestAmino {
  index?: string;
}
export interface QueryGetLastBlockHeightRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetLastBlockHeightRequest";
  value: QueryGetLastBlockHeightRequestAmino;
}
export interface QueryGetLastBlockHeightRequestSDKType {
  index: string;
}
export interface QueryGetLastBlockHeightResponse {
  LastBlockHeight?: LastBlockHeight;
}
export interface QueryGetLastBlockHeightResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetLastBlockHeightResponse";
  value: Uint8Array;
}
export interface QueryGetLastBlockHeightResponseAmino {
  LastBlockHeight?: LastBlockHeightAmino;
}
export interface QueryGetLastBlockHeightResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetLastBlockHeightResponse";
  value: QueryGetLastBlockHeightResponseAmino;
}
export interface QueryGetLastBlockHeightResponseSDKType {
  LastBlockHeight?: LastBlockHeightSDKType;
}
export interface QueryAllLastBlockHeightRequest {
  pagination?: PageRequest;
}
export interface QueryAllLastBlockHeightRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllLastBlockHeightRequest";
  value: Uint8Array;
}
export interface QueryAllLastBlockHeightRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllLastBlockHeightRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllLastBlockHeightRequest";
  value: QueryAllLastBlockHeightRequestAmino;
}
export interface QueryAllLastBlockHeightRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllLastBlockHeightResponse {
  LastBlockHeight: LastBlockHeight[];
  pagination?: PageResponse;
}
export interface QueryAllLastBlockHeightResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllLastBlockHeightResponse";
  value: Uint8Array;
}
export interface QueryAllLastBlockHeightResponseAmino {
  LastBlockHeight?: LastBlockHeightAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllLastBlockHeightResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllLastBlockHeightResponse";
  value: QueryAllLastBlockHeightResponseAmino;
}
export interface QueryAllLastBlockHeightResponseSDKType {
  LastBlockHeight: LastBlockHeightSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryGetCctxRequest {
  index: string;
}
export interface QueryGetCctxRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetCctxRequest";
  value: Uint8Array;
}
export interface QueryGetCctxRequestAmino {
  index?: string;
}
export interface QueryGetCctxRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetCctxRequest";
  value: QueryGetCctxRequestAmino;
}
export interface QueryGetCctxRequestSDKType {
  index: string;
}
export interface QueryGetCctxByNonceRequest {
  chainID: bigint;
  nonce: bigint;
}
export interface QueryGetCctxByNonceRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetCctxByNonceRequest";
  value: Uint8Array;
}
export interface QueryGetCctxByNonceRequestAmino {
  chainID?: string;
  nonce?: string;
}
export interface QueryGetCctxByNonceRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetCctxByNonceRequest";
  value: QueryGetCctxByNonceRequestAmino;
}
export interface QueryGetCctxByNonceRequestSDKType {
  chainID: bigint;
  nonce: bigint;
}
export interface QueryGetCctxResponse {
  CrossChainTx?: CrossChainTx;
}
export interface QueryGetCctxResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetCctxResponse";
  value: Uint8Array;
}
export interface QueryGetCctxResponseAmino {
  CrossChainTx?: CrossChainTxAmino;
}
export interface QueryGetCctxResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryGetCctxResponse";
  value: QueryGetCctxResponseAmino;
}
export interface QueryGetCctxResponseSDKType {
  CrossChainTx?: CrossChainTxSDKType;
}
export interface QueryAllCctxRequest {
  pagination?: PageRequest;
}
export interface QueryAllCctxRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllCctxRequest";
  value: Uint8Array;
}
export interface QueryAllCctxRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllCctxRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllCctxRequest";
  value: QueryAllCctxRequestAmino;
}
export interface QueryAllCctxRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllCctxResponse {
  CrossChainTx: CrossChainTx[];
  pagination?: PageResponse;
}
export interface QueryAllCctxResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllCctxResponse";
  value: Uint8Array;
}
export interface QueryAllCctxResponseAmino {
  CrossChainTx?: CrossChainTxAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllCctxResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryAllCctxResponse";
  value: QueryAllCctxResponseAmino;
}
export interface QueryAllCctxResponseSDKType {
  CrossChainTx: CrossChainTxSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryListCctxPendingRequest {
  chainId: bigint;
  limit: number;
}
export interface QueryListCctxPendingRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryListCctxPendingRequest";
  value: Uint8Array;
}
export interface QueryListCctxPendingRequestAmino {
  chain_id?: string;
  limit?: number;
}
export interface QueryListCctxPendingRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryListCctxPendingRequest";
  value: QueryListCctxPendingRequestAmino;
}
export interface QueryListCctxPendingRequestSDKType {
  chain_id: bigint;
  limit: number;
}
export interface QueryListCctxPendingResponse {
  CrossChainTx: CrossChainTx[];
  totalPending: bigint;
}
export interface QueryListCctxPendingResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryListCctxPendingResponse";
  value: Uint8Array;
}
export interface QueryListCctxPendingResponseAmino {
  CrossChainTx?: CrossChainTxAmino[];
  totalPending?: string;
}
export interface QueryListCctxPendingResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryListCctxPendingResponse";
  value: QueryListCctxPendingResponseAmino;
}
export interface QueryListCctxPendingResponseSDKType {
  CrossChainTx: CrossChainTxSDKType[];
  totalPending: bigint;
}
export interface QueryLastZetaHeightRequest {}
export interface QueryLastZetaHeightRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryLastZetaHeightRequest";
  value: Uint8Array;
}
export interface QueryLastZetaHeightRequestAmino {}
export interface QueryLastZetaHeightRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryLastZetaHeightRequest";
  value: QueryLastZetaHeightRequestAmino;
}
export interface QueryLastZetaHeightRequestSDKType {}
export interface QueryLastZetaHeightResponse {
  Height: bigint;
}
export interface QueryLastZetaHeightResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryLastZetaHeightResponse";
  value: Uint8Array;
}
export interface QueryLastZetaHeightResponseAmino {
  Height?: string;
}
export interface QueryLastZetaHeightResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryLastZetaHeightResponse";
  value: QueryLastZetaHeightResponseAmino;
}
export interface QueryLastZetaHeightResponseSDKType {
  Height: bigint;
}
export interface QueryConvertGasToZetaRequest {
  chainId: bigint;
  gasLimit: string;
}
export interface QueryConvertGasToZetaRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryConvertGasToZetaRequest";
  value: Uint8Array;
}
export interface QueryConvertGasToZetaRequestAmino {
  chainId?: string;
  gasLimit?: string;
}
export interface QueryConvertGasToZetaRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryConvertGasToZetaRequest";
  value: QueryConvertGasToZetaRequestAmino;
}
export interface QueryConvertGasToZetaRequestSDKType {
  chainId: bigint;
  gasLimit: string;
}
export interface QueryConvertGasToZetaResponse {
  outboundGasInZeta: string;
  protocolFeeInZeta: string;
  ZetaBlockHeight: bigint;
}
export interface QueryConvertGasToZetaResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryConvertGasToZetaResponse";
  value: Uint8Array;
}
export interface QueryConvertGasToZetaResponseAmino {
  outboundGasInZeta?: string;
  protocolFeeInZeta?: string;
  ZetaBlockHeight?: string;
}
export interface QueryConvertGasToZetaResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryConvertGasToZetaResponse";
  value: QueryConvertGasToZetaResponseAmino;
}
export interface QueryConvertGasToZetaResponseSDKType {
  outboundGasInZeta: string;
  protocolFeeInZeta: string;
  ZetaBlockHeight: bigint;
}
export interface QueryMessagePassingProtocolFeeRequest {}
export interface QueryMessagePassingProtocolFeeRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryMessagePassingProtocolFeeRequest";
  value: Uint8Array;
}
export interface QueryMessagePassingProtocolFeeRequestAmino {}
export interface QueryMessagePassingProtocolFeeRequestAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryMessagePassingProtocolFeeRequest";
  value: QueryMessagePassingProtocolFeeRequestAmino;
}
export interface QueryMessagePassingProtocolFeeRequestSDKType {}
export interface QueryMessagePassingProtocolFeeResponse {
  feeInZeta: string;
}
export interface QueryMessagePassingProtocolFeeResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.crosschain.QueryMessagePassingProtocolFeeResponse";
  value: Uint8Array;
}
export interface QueryMessagePassingProtocolFeeResponseAmino {
  feeInZeta?: string;
}
export interface QueryMessagePassingProtocolFeeResponseAminoMsg {
  type: "/zetachain.zetacore.crosschain.QueryMessagePassingProtocolFeeResponse";
  value: QueryMessagePassingProtocolFeeResponseAmino;
}
export interface QueryMessagePassingProtocolFeeResponseSDKType {
  feeInZeta: string;
}
function createBaseQueryGetTssAddressRequest(): QueryGetTssAddressRequest {
  return {};
}
export const QueryGetTssAddressRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetTssAddressRequest",
  encode(_: QueryGetTssAddressRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetTssAddressRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetTssAddressRequest();
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
  fromPartial(_: Partial<QueryGetTssAddressRequest>): QueryGetTssAddressRequest {
    const message = createBaseQueryGetTssAddressRequest();
    return message;
  },
  fromAmino(_: QueryGetTssAddressRequestAmino): QueryGetTssAddressRequest {
    const message = createBaseQueryGetTssAddressRequest();
    return message;
  },
  toAmino(_: QueryGetTssAddressRequest): QueryGetTssAddressRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryGetTssAddressRequestAminoMsg): QueryGetTssAddressRequest {
    return QueryGetTssAddressRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetTssAddressRequestProtoMsg): QueryGetTssAddressRequest {
    return QueryGetTssAddressRequest.decode(message.value);
  },
  toProto(message: QueryGetTssAddressRequest): Uint8Array {
    return QueryGetTssAddressRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetTssAddressRequest): QueryGetTssAddressRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetTssAddressRequest",
      value: QueryGetTssAddressRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetTssAddressResponse(): QueryGetTssAddressResponse {
  return {
    eth: "",
    btc: ""
  };
}
export const QueryGetTssAddressResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetTssAddressResponse",
  encode(message: QueryGetTssAddressResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.eth !== "") {
      writer.uint32(10).string(message.eth);
    }
    if (message.btc !== "") {
      writer.uint32(18).string(message.btc);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetTssAddressResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetTssAddressResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.eth = reader.string();
          break;
        case 2:
          message.btc = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetTssAddressResponse>): QueryGetTssAddressResponse {
    const message = createBaseQueryGetTssAddressResponse();
    message.eth = object.eth ?? "";
    message.btc = object.btc ?? "";
    return message;
  },
  fromAmino(object: QueryGetTssAddressResponseAmino): QueryGetTssAddressResponse {
    const message = createBaseQueryGetTssAddressResponse();
    if (object.eth !== undefined && object.eth !== null) {
      message.eth = object.eth;
    }
    if (object.btc !== undefined && object.btc !== null) {
      message.btc = object.btc;
    }
    return message;
  },
  toAmino(message: QueryGetTssAddressResponse): QueryGetTssAddressResponseAmino {
    const obj: any = {};
    obj.eth = message.eth;
    obj.btc = message.btc;
    return obj;
  },
  fromAminoMsg(object: QueryGetTssAddressResponseAminoMsg): QueryGetTssAddressResponse {
    return QueryGetTssAddressResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetTssAddressResponseProtoMsg): QueryGetTssAddressResponse {
    return QueryGetTssAddressResponse.decode(message.value);
  },
  toProto(message: QueryGetTssAddressResponse): Uint8Array {
    return QueryGetTssAddressResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetTssAddressResponse): QueryGetTssAddressResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetTssAddressResponse",
      value: QueryGetTssAddressResponse.encode(message).finish()
    };
  }
};
function createBaseQueryZetaAccountingRequest(): QueryZetaAccountingRequest {
  return {};
}
export const QueryZetaAccountingRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryZetaAccountingRequest",
  encode(_: QueryZetaAccountingRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryZetaAccountingRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryZetaAccountingRequest();
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
  fromPartial(_: Partial<QueryZetaAccountingRequest>): QueryZetaAccountingRequest {
    const message = createBaseQueryZetaAccountingRequest();
    return message;
  },
  fromAmino(_: QueryZetaAccountingRequestAmino): QueryZetaAccountingRequest {
    const message = createBaseQueryZetaAccountingRequest();
    return message;
  },
  toAmino(_: QueryZetaAccountingRequest): QueryZetaAccountingRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryZetaAccountingRequestAminoMsg): QueryZetaAccountingRequest {
    return QueryZetaAccountingRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryZetaAccountingRequestProtoMsg): QueryZetaAccountingRequest {
    return QueryZetaAccountingRequest.decode(message.value);
  },
  toProto(message: QueryZetaAccountingRequest): Uint8Array {
    return QueryZetaAccountingRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryZetaAccountingRequest): QueryZetaAccountingRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryZetaAccountingRequest",
      value: QueryZetaAccountingRequest.encode(message).finish()
    };
  }
};
function createBaseQueryZetaAccountingResponse(): QueryZetaAccountingResponse {
  return {
    abortedZetaAmount: ""
  };
}
export const QueryZetaAccountingResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryZetaAccountingResponse",
  encode(message: QueryZetaAccountingResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.abortedZetaAmount !== "") {
      writer.uint32(10).string(message.abortedZetaAmount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryZetaAccountingResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryZetaAccountingResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.abortedZetaAmount = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryZetaAccountingResponse>): QueryZetaAccountingResponse {
    const message = createBaseQueryZetaAccountingResponse();
    message.abortedZetaAmount = object.abortedZetaAmount ?? "";
    return message;
  },
  fromAmino(object: QueryZetaAccountingResponseAmino): QueryZetaAccountingResponse {
    const message = createBaseQueryZetaAccountingResponse();
    if (object.aborted_zeta_amount !== undefined && object.aborted_zeta_amount !== null) {
      message.abortedZetaAmount = object.aborted_zeta_amount;
    }
    return message;
  },
  toAmino(message: QueryZetaAccountingResponse): QueryZetaAccountingResponseAmino {
    const obj: any = {};
    obj.aborted_zeta_amount = message.abortedZetaAmount;
    return obj;
  },
  fromAminoMsg(object: QueryZetaAccountingResponseAminoMsg): QueryZetaAccountingResponse {
    return QueryZetaAccountingResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryZetaAccountingResponseProtoMsg): QueryZetaAccountingResponse {
    return QueryZetaAccountingResponse.decode(message.value);
  },
  toProto(message: QueryZetaAccountingResponse): Uint8Array {
    return QueryZetaAccountingResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryZetaAccountingResponse): QueryZetaAccountingResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryZetaAccountingResponse",
      value: QueryZetaAccountingResponse.encode(message).finish()
    };
  }
};
function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}
export const QueryParamsRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryParamsRequest",
  encode(_: QueryParamsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
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
  fromPartial(_: Partial<QueryParamsRequest>): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
  fromAmino(_: QueryParamsRequestAmino): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
  toAmino(_: QueryParamsRequest): QueryParamsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryParamsRequestAminoMsg): QueryParamsRequest {
    return QueryParamsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryParamsRequestProtoMsg): QueryParamsRequest {
    return QueryParamsRequest.decode(message.value);
  },
  toProto(message: QueryParamsRequest): Uint8Array {
    return QueryParamsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryParamsRequest): QueryParamsRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryParamsRequest",
      value: QueryParamsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryParamsResponse(): QueryParamsResponse {
  return {
    params: Params.fromPartial({})
  };
}
export const QueryParamsResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryParamsResponse",
  encode(message: QueryParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryParamsResponse>): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: QueryParamsResponseAmino): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: QueryParamsResponse): QueryParamsResponseAmino {
    const obj: any = {};
    obj.params = message.params ? Params.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryParamsResponseAminoMsg): QueryParamsResponse {
    return QueryParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryParamsResponseProtoMsg): QueryParamsResponse {
    return QueryParamsResponse.decode(message.value);
  },
  toProto(message: QueryParamsResponse): Uint8Array {
    return QueryParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryParamsResponse): QueryParamsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryParamsResponse",
      value: QueryParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetOutTxTrackerRequest(): QueryGetOutTxTrackerRequest {
  return {
    chainID: BigInt(0),
    nonce: BigInt(0)
  };
}
export const QueryGetOutTxTrackerRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetOutTxTrackerRequest",
  encode(message: QueryGetOutTxTrackerRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainID !== BigInt(0)) {
      writer.uint32(8).int64(message.chainID);
    }
    if (message.nonce !== BigInt(0)) {
      writer.uint32(16).uint64(message.nonce);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetOutTxTrackerRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetOutTxTrackerRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainID = reader.int64();
          break;
        case 2:
          message.nonce = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetOutTxTrackerRequest>): QueryGetOutTxTrackerRequest {
    const message = createBaseQueryGetOutTxTrackerRequest();
    message.chainID = object.chainID !== undefined && object.chainID !== null ? BigInt(object.chainID.toString()) : BigInt(0);
    message.nonce = object.nonce !== undefined && object.nonce !== null ? BigInt(object.nonce.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryGetOutTxTrackerRequestAmino): QueryGetOutTxTrackerRequest {
    const message = createBaseQueryGetOutTxTrackerRequest();
    if (object.chainID !== undefined && object.chainID !== null) {
      message.chainID = BigInt(object.chainID);
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = BigInt(object.nonce);
    }
    return message;
  },
  toAmino(message: QueryGetOutTxTrackerRequest): QueryGetOutTxTrackerRequestAmino {
    const obj: any = {};
    obj.chainID = message.chainID ? message.chainID.toString() : undefined;
    obj.nonce = message.nonce ? message.nonce.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetOutTxTrackerRequestAminoMsg): QueryGetOutTxTrackerRequest {
    return QueryGetOutTxTrackerRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetOutTxTrackerRequestProtoMsg): QueryGetOutTxTrackerRequest {
    return QueryGetOutTxTrackerRequest.decode(message.value);
  },
  toProto(message: QueryGetOutTxTrackerRequest): Uint8Array {
    return QueryGetOutTxTrackerRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetOutTxTrackerRequest): QueryGetOutTxTrackerRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetOutTxTrackerRequest",
      value: QueryGetOutTxTrackerRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetOutTxTrackerResponse(): QueryGetOutTxTrackerResponse {
  return {
    outTxTracker: OutTxTracker.fromPartial({})
  };
}
export const QueryGetOutTxTrackerResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetOutTxTrackerResponse",
  encode(message: QueryGetOutTxTrackerResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.outTxTracker !== undefined) {
      OutTxTracker.encode(message.outTxTracker, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetOutTxTrackerResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetOutTxTrackerResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.outTxTracker = OutTxTracker.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetOutTxTrackerResponse>): QueryGetOutTxTrackerResponse {
    const message = createBaseQueryGetOutTxTrackerResponse();
    message.outTxTracker = object.outTxTracker !== undefined && object.outTxTracker !== null ? OutTxTracker.fromPartial(object.outTxTracker) : undefined;
    return message;
  },
  fromAmino(object: QueryGetOutTxTrackerResponseAmino): QueryGetOutTxTrackerResponse {
    const message = createBaseQueryGetOutTxTrackerResponse();
    if (object.outTxTracker !== undefined && object.outTxTracker !== null) {
      message.outTxTracker = OutTxTracker.fromAmino(object.outTxTracker);
    }
    return message;
  },
  toAmino(message: QueryGetOutTxTrackerResponse): QueryGetOutTxTrackerResponseAmino {
    const obj: any = {};
    obj.outTxTracker = message.outTxTracker ? OutTxTracker.toAmino(message.outTxTracker) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetOutTxTrackerResponseAminoMsg): QueryGetOutTxTrackerResponse {
    return QueryGetOutTxTrackerResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetOutTxTrackerResponseProtoMsg): QueryGetOutTxTrackerResponse {
    return QueryGetOutTxTrackerResponse.decode(message.value);
  },
  toProto(message: QueryGetOutTxTrackerResponse): Uint8Array {
    return QueryGetOutTxTrackerResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetOutTxTrackerResponse): QueryGetOutTxTrackerResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetOutTxTrackerResponse",
      value: QueryGetOutTxTrackerResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllOutTxTrackerRequest(): QueryAllOutTxTrackerRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllOutTxTrackerRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerRequest",
  encode(message: QueryAllOutTxTrackerRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllOutTxTrackerRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllOutTxTrackerRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllOutTxTrackerRequest>): QueryAllOutTxTrackerRequest {
    const message = createBaseQueryAllOutTxTrackerRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllOutTxTrackerRequestAmino): QueryAllOutTxTrackerRequest {
    const message = createBaseQueryAllOutTxTrackerRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllOutTxTrackerRequest): QueryAllOutTxTrackerRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllOutTxTrackerRequestAminoMsg): QueryAllOutTxTrackerRequest {
    return QueryAllOutTxTrackerRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllOutTxTrackerRequestProtoMsg): QueryAllOutTxTrackerRequest {
    return QueryAllOutTxTrackerRequest.decode(message.value);
  },
  toProto(message: QueryAllOutTxTrackerRequest): Uint8Array {
    return QueryAllOutTxTrackerRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllOutTxTrackerRequest): QueryAllOutTxTrackerRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerRequest",
      value: QueryAllOutTxTrackerRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllOutTxTrackerResponse(): QueryAllOutTxTrackerResponse {
  return {
    outTxTracker: [],
    pagination: undefined
  };
}
export const QueryAllOutTxTrackerResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerResponse",
  encode(message: QueryAllOutTxTrackerResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.outTxTracker) {
      OutTxTracker.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllOutTxTrackerResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllOutTxTrackerResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.outTxTracker.push(OutTxTracker.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllOutTxTrackerResponse>): QueryAllOutTxTrackerResponse {
    const message = createBaseQueryAllOutTxTrackerResponse();
    message.outTxTracker = object.outTxTracker?.map(e => OutTxTracker.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllOutTxTrackerResponseAmino): QueryAllOutTxTrackerResponse {
    const message = createBaseQueryAllOutTxTrackerResponse();
    message.outTxTracker = object.outTxTracker?.map(e => OutTxTracker.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllOutTxTrackerResponse): QueryAllOutTxTrackerResponseAmino {
    const obj: any = {};
    if (message.outTxTracker) {
      obj.outTxTracker = message.outTxTracker.map(e => e ? OutTxTracker.toAmino(e) : undefined);
    } else {
      obj.outTxTracker = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllOutTxTrackerResponseAminoMsg): QueryAllOutTxTrackerResponse {
    return QueryAllOutTxTrackerResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllOutTxTrackerResponseProtoMsg): QueryAllOutTxTrackerResponse {
    return QueryAllOutTxTrackerResponse.decode(message.value);
  },
  toProto(message: QueryAllOutTxTrackerResponse): Uint8Array {
    return QueryAllOutTxTrackerResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllOutTxTrackerResponse): QueryAllOutTxTrackerResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerResponse",
      value: QueryAllOutTxTrackerResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllOutTxTrackerByChainRequest(): QueryAllOutTxTrackerByChainRequest {
  return {
    chain: BigInt(0),
    pagination: undefined
  };
}
export const QueryAllOutTxTrackerByChainRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerByChainRequest",
  encode(message: QueryAllOutTxTrackerByChainRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chain !== BigInt(0)) {
      writer.uint32(8).int64(message.chain);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllOutTxTrackerByChainRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllOutTxTrackerByChainRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chain = reader.int64();
          break;
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllOutTxTrackerByChainRequest>): QueryAllOutTxTrackerByChainRequest {
    const message = createBaseQueryAllOutTxTrackerByChainRequest();
    message.chain = object.chain !== undefined && object.chain !== null ? BigInt(object.chain.toString()) : BigInt(0);
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllOutTxTrackerByChainRequestAmino): QueryAllOutTxTrackerByChainRequest {
    const message = createBaseQueryAllOutTxTrackerByChainRequest();
    if (object.chain !== undefined && object.chain !== null) {
      message.chain = BigInt(object.chain);
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllOutTxTrackerByChainRequest): QueryAllOutTxTrackerByChainRequestAmino {
    const obj: any = {};
    obj.chain = message.chain ? message.chain.toString() : undefined;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllOutTxTrackerByChainRequestAminoMsg): QueryAllOutTxTrackerByChainRequest {
    return QueryAllOutTxTrackerByChainRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllOutTxTrackerByChainRequestProtoMsg): QueryAllOutTxTrackerByChainRequest {
    return QueryAllOutTxTrackerByChainRequest.decode(message.value);
  },
  toProto(message: QueryAllOutTxTrackerByChainRequest): Uint8Array {
    return QueryAllOutTxTrackerByChainRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllOutTxTrackerByChainRequest): QueryAllOutTxTrackerByChainRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerByChainRequest",
      value: QueryAllOutTxTrackerByChainRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllOutTxTrackerByChainResponse(): QueryAllOutTxTrackerByChainResponse {
  return {
    outTxTracker: [],
    pagination: undefined
  };
}
export const QueryAllOutTxTrackerByChainResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerByChainResponse",
  encode(message: QueryAllOutTxTrackerByChainResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.outTxTracker) {
      OutTxTracker.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllOutTxTrackerByChainResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllOutTxTrackerByChainResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.outTxTracker.push(OutTxTracker.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllOutTxTrackerByChainResponse>): QueryAllOutTxTrackerByChainResponse {
    const message = createBaseQueryAllOutTxTrackerByChainResponse();
    message.outTxTracker = object.outTxTracker?.map(e => OutTxTracker.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllOutTxTrackerByChainResponseAmino): QueryAllOutTxTrackerByChainResponse {
    const message = createBaseQueryAllOutTxTrackerByChainResponse();
    message.outTxTracker = object.outTxTracker?.map(e => OutTxTracker.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllOutTxTrackerByChainResponse): QueryAllOutTxTrackerByChainResponseAmino {
    const obj: any = {};
    if (message.outTxTracker) {
      obj.outTxTracker = message.outTxTracker.map(e => e ? OutTxTracker.toAmino(e) : undefined);
    } else {
      obj.outTxTracker = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllOutTxTrackerByChainResponseAminoMsg): QueryAllOutTxTrackerByChainResponse {
    return QueryAllOutTxTrackerByChainResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllOutTxTrackerByChainResponseProtoMsg): QueryAllOutTxTrackerByChainResponse {
    return QueryAllOutTxTrackerByChainResponse.decode(message.value);
  },
  toProto(message: QueryAllOutTxTrackerByChainResponse): Uint8Array {
    return QueryAllOutTxTrackerByChainResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllOutTxTrackerByChainResponse): QueryAllOutTxTrackerByChainResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllOutTxTrackerByChainResponse",
      value: QueryAllOutTxTrackerByChainResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllInTxTrackerByChainRequest(): QueryAllInTxTrackerByChainRequest {
  return {
    chainId: BigInt(0),
    pagination: undefined
  };
}
export const QueryAllInTxTrackerByChainRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackerByChainRequest",
  encode(message: QueryAllInTxTrackerByChainRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllInTxTrackerByChainRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllInTxTrackerByChainRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllInTxTrackerByChainRequest>): QueryAllInTxTrackerByChainRequest {
    const message = createBaseQueryAllInTxTrackerByChainRequest();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllInTxTrackerByChainRequestAmino): QueryAllInTxTrackerByChainRequest {
    const message = createBaseQueryAllInTxTrackerByChainRequest();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllInTxTrackerByChainRequest): QueryAllInTxTrackerByChainRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllInTxTrackerByChainRequestAminoMsg): QueryAllInTxTrackerByChainRequest {
    return QueryAllInTxTrackerByChainRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllInTxTrackerByChainRequestProtoMsg): QueryAllInTxTrackerByChainRequest {
    return QueryAllInTxTrackerByChainRequest.decode(message.value);
  },
  toProto(message: QueryAllInTxTrackerByChainRequest): Uint8Array {
    return QueryAllInTxTrackerByChainRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllInTxTrackerByChainRequest): QueryAllInTxTrackerByChainRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackerByChainRequest",
      value: QueryAllInTxTrackerByChainRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllInTxTrackerByChainResponse(): QueryAllInTxTrackerByChainResponse {
  return {
    inTxTracker: [],
    pagination: undefined
  };
}
export const QueryAllInTxTrackerByChainResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackerByChainResponse",
  encode(message: QueryAllInTxTrackerByChainResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.inTxTracker) {
      InTxTracker.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllInTxTrackerByChainResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllInTxTrackerByChainResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.inTxTracker.push(InTxTracker.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllInTxTrackerByChainResponse>): QueryAllInTxTrackerByChainResponse {
    const message = createBaseQueryAllInTxTrackerByChainResponse();
    message.inTxTracker = object.inTxTracker?.map(e => InTxTracker.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllInTxTrackerByChainResponseAmino): QueryAllInTxTrackerByChainResponse {
    const message = createBaseQueryAllInTxTrackerByChainResponse();
    message.inTxTracker = object.inTxTracker?.map(e => InTxTracker.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllInTxTrackerByChainResponse): QueryAllInTxTrackerByChainResponseAmino {
    const obj: any = {};
    if (message.inTxTracker) {
      obj.inTxTracker = message.inTxTracker.map(e => e ? InTxTracker.toAmino(e) : undefined);
    } else {
      obj.inTxTracker = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllInTxTrackerByChainResponseAminoMsg): QueryAllInTxTrackerByChainResponse {
    return QueryAllInTxTrackerByChainResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllInTxTrackerByChainResponseProtoMsg): QueryAllInTxTrackerByChainResponse {
    return QueryAllInTxTrackerByChainResponse.decode(message.value);
  },
  toProto(message: QueryAllInTxTrackerByChainResponse): Uint8Array {
    return QueryAllInTxTrackerByChainResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllInTxTrackerByChainResponse): QueryAllInTxTrackerByChainResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackerByChainResponse",
      value: QueryAllInTxTrackerByChainResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllInTxTrackersRequest(): QueryAllInTxTrackersRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllInTxTrackersRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackersRequest",
  encode(message: QueryAllInTxTrackersRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllInTxTrackersRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllInTxTrackersRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllInTxTrackersRequest>): QueryAllInTxTrackersRequest {
    const message = createBaseQueryAllInTxTrackersRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllInTxTrackersRequestAmino): QueryAllInTxTrackersRequest {
    const message = createBaseQueryAllInTxTrackersRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllInTxTrackersRequest): QueryAllInTxTrackersRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllInTxTrackersRequestAminoMsg): QueryAllInTxTrackersRequest {
    return QueryAllInTxTrackersRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllInTxTrackersRequestProtoMsg): QueryAllInTxTrackersRequest {
    return QueryAllInTxTrackersRequest.decode(message.value);
  },
  toProto(message: QueryAllInTxTrackersRequest): Uint8Array {
    return QueryAllInTxTrackersRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllInTxTrackersRequest): QueryAllInTxTrackersRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackersRequest",
      value: QueryAllInTxTrackersRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllInTxTrackersResponse(): QueryAllInTxTrackersResponse {
  return {
    inTxTracker: [],
    pagination: undefined
  };
}
export const QueryAllInTxTrackersResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackersResponse",
  encode(message: QueryAllInTxTrackersResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.inTxTracker) {
      InTxTracker.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllInTxTrackersResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllInTxTrackersResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.inTxTracker.push(InTxTracker.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllInTxTrackersResponse>): QueryAllInTxTrackersResponse {
    const message = createBaseQueryAllInTxTrackersResponse();
    message.inTxTracker = object.inTxTracker?.map(e => InTxTracker.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllInTxTrackersResponseAmino): QueryAllInTxTrackersResponse {
    const message = createBaseQueryAllInTxTrackersResponse();
    message.inTxTracker = object.inTxTracker?.map(e => InTxTracker.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllInTxTrackersResponse): QueryAllInTxTrackersResponseAmino {
    const obj: any = {};
    if (message.inTxTracker) {
      obj.inTxTracker = message.inTxTracker.map(e => e ? InTxTracker.toAmino(e) : undefined);
    } else {
      obj.inTxTracker = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllInTxTrackersResponseAminoMsg): QueryAllInTxTrackersResponse {
    return QueryAllInTxTrackersResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllInTxTrackersResponseProtoMsg): QueryAllInTxTrackersResponse {
    return QueryAllInTxTrackersResponse.decode(message.value);
  },
  toProto(message: QueryAllInTxTrackersResponse): Uint8Array {
    return QueryAllInTxTrackersResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllInTxTrackersResponse): QueryAllInTxTrackersResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxTrackersResponse",
      value: QueryAllInTxTrackersResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetInTxHashToCctxRequest(): QueryGetInTxHashToCctxRequest {
  return {
    inTxHash: ""
  };
}
export const QueryGetInTxHashToCctxRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetInTxHashToCctxRequest",
  encode(message: QueryGetInTxHashToCctxRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.inTxHash !== "") {
      writer.uint32(10).string(message.inTxHash);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetInTxHashToCctxRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetInTxHashToCctxRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.inTxHash = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetInTxHashToCctxRequest>): QueryGetInTxHashToCctxRequest {
    const message = createBaseQueryGetInTxHashToCctxRequest();
    message.inTxHash = object.inTxHash ?? "";
    return message;
  },
  fromAmino(object: QueryGetInTxHashToCctxRequestAmino): QueryGetInTxHashToCctxRequest {
    const message = createBaseQueryGetInTxHashToCctxRequest();
    if (object.inTxHash !== undefined && object.inTxHash !== null) {
      message.inTxHash = object.inTxHash;
    }
    return message;
  },
  toAmino(message: QueryGetInTxHashToCctxRequest): QueryGetInTxHashToCctxRequestAmino {
    const obj: any = {};
    obj.inTxHash = message.inTxHash;
    return obj;
  },
  fromAminoMsg(object: QueryGetInTxHashToCctxRequestAminoMsg): QueryGetInTxHashToCctxRequest {
    return QueryGetInTxHashToCctxRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetInTxHashToCctxRequestProtoMsg): QueryGetInTxHashToCctxRequest {
    return QueryGetInTxHashToCctxRequest.decode(message.value);
  },
  toProto(message: QueryGetInTxHashToCctxRequest): Uint8Array {
    return QueryGetInTxHashToCctxRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetInTxHashToCctxRequest): QueryGetInTxHashToCctxRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetInTxHashToCctxRequest",
      value: QueryGetInTxHashToCctxRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetInTxHashToCctxResponse(): QueryGetInTxHashToCctxResponse {
  return {
    inTxHashToCctx: InTxHashToCctx.fromPartial({})
  };
}
export const QueryGetInTxHashToCctxResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetInTxHashToCctxResponse",
  encode(message: QueryGetInTxHashToCctxResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.inTxHashToCctx !== undefined) {
      InTxHashToCctx.encode(message.inTxHashToCctx, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetInTxHashToCctxResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetInTxHashToCctxResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.inTxHashToCctx = InTxHashToCctx.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetInTxHashToCctxResponse>): QueryGetInTxHashToCctxResponse {
    const message = createBaseQueryGetInTxHashToCctxResponse();
    message.inTxHashToCctx = object.inTxHashToCctx !== undefined && object.inTxHashToCctx !== null ? InTxHashToCctx.fromPartial(object.inTxHashToCctx) : undefined;
    return message;
  },
  fromAmino(object: QueryGetInTxHashToCctxResponseAmino): QueryGetInTxHashToCctxResponse {
    const message = createBaseQueryGetInTxHashToCctxResponse();
    if (object.inTxHashToCctx !== undefined && object.inTxHashToCctx !== null) {
      message.inTxHashToCctx = InTxHashToCctx.fromAmino(object.inTxHashToCctx);
    }
    return message;
  },
  toAmino(message: QueryGetInTxHashToCctxResponse): QueryGetInTxHashToCctxResponseAmino {
    const obj: any = {};
    obj.inTxHashToCctx = message.inTxHashToCctx ? InTxHashToCctx.toAmino(message.inTxHashToCctx) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetInTxHashToCctxResponseAminoMsg): QueryGetInTxHashToCctxResponse {
    return QueryGetInTxHashToCctxResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetInTxHashToCctxResponseProtoMsg): QueryGetInTxHashToCctxResponse {
    return QueryGetInTxHashToCctxResponse.decode(message.value);
  },
  toProto(message: QueryGetInTxHashToCctxResponse): Uint8Array {
    return QueryGetInTxHashToCctxResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetInTxHashToCctxResponse): QueryGetInTxHashToCctxResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetInTxHashToCctxResponse",
      value: QueryGetInTxHashToCctxResponse.encode(message).finish()
    };
  }
};
function createBaseQueryInTxHashToCctxDataRequest(): QueryInTxHashToCctxDataRequest {
  return {
    inTxHash: ""
  };
}
export const QueryInTxHashToCctxDataRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryInTxHashToCctxDataRequest",
  encode(message: QueryInTxHashToCctxDataRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.inTxHash !== "") {
      writer.uint32(10).string(message.inTxHash);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryInTxHashToCctxDataRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryInTxHashToCctxDataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.inTxHash = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryInTxHashToCctxDataRequest>): QueryInTxHashToCctxDataRequest {
    const message = createBaseQueryInTxHashToCctxDataRequest();
    message.inTxHash = object.inTxHash ?? "";
    return message;
  },
  fromAmino(object: QueryInTxHashToCctxDataRequestAmino): QueryInTxHashToCctxDataRequest {
    const message = createBaseQueryInTxHashToCctxDataRequest();
    if (object.inTxHash !== undefined && object.inTxHash !== null) {
      message.inTxHash = object.inTxHash;
    }
    return message;
  },
  toAmino(message: QueryInTxHashToCctxDataRequest): QueryInTxHashToCctxDataRequestAmino {
    const obj: any = {};
    obj.inTxHash = message.inTxHash;
    return obj;
  },
  fromAminoMsg(object: QueryInTxHashToCctxDataRequestAminoMsg): QueryInTxHashToCctxDataRequest {
    return QueryInTxHashToCctxDataRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryInTxHashToCctxDataRequestProtoMsg): QueryInTxHashToCctxDataRequest {
    return QueryInTxHashToCctxDataRequest.decode(message.value);
  },
  toProto(message: QueryInTxHashToCctxDataRequest): Uint8Array {
    return QueryInTxHashToCctxDataRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryInTxHashToCctxDataRequest): QueryInTxHashToCctxDataRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryInTxHashToCctxDataRequest",
      value: QueryInTxHashToCctxDataRequest.encode(message).finish()
    };
  }
};
function createBaseQueryInTxHashToCctxDataResponse(): QueryInTxHashToCctxDataResponse {
  return {
    CrossChainTxs: []
  };
}
export const QueryInTxHashToCctxDataResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryInTxHashToCctxDataResponse",
  encode(message: QueryInTxHashToCctxDataResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.CrossChainTxs) {
      CrossChainTx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryInTxHashToCctxDataResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryInTxHashToCctxDataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.CrossChainTxs.push(CrossChainTx.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryInTxHashToCctxDataResponse>): QueryInTxHashToCctxDataResponse {
    const message = createBaseQueryInTxHashToCctxDataResponse();
    message.CrossChainTxs = object.CrossChainTxs?.map(e => CrossChainTx.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: QueryInTxHashToCctxDataResponseAmino): QueryInTxHashToCctxDataResponse {
    const message = createBaseQueryInTxHashToCctxDataResponse();
    message.CrossChainTxs = object.CrossChainTxs?.map(e => CrossChainTx.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: QueryInTxHashToCctxDataResponse): QueryInTxHashToCctxDataResponseAmino {
    const obj: any = {};
    if (message.CrossChainTxs) {
      obj.CrossChainTxs = message.CrossChainTxs.map(e => e ? CrossChainTx.toAmino(e) : undefined);
    } else {
      obj.CrossChainTxs = [];
    }
    return obj;
  },
  fromAminoMsg(object: QueryInTxHashToCctxDataResponseAminoMsg): QueryInTxHashToCctxDataResponse {
    return QueryInTxHashToCctxDataResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryInTxHashToCctxDataResponseProtoMsg): QueryInTxHashToCctxDataResponse {
    return QueryInTxHashToCctxDataResponse.decode(message.value);
  },
  toProto(message: QueryInTxHashToCctxDataResponse): Uint8Array {
    return QueryInTxHashToCctxDataResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryInTxHashToCctxDataResponse): QueryInTxHashToCctxDataResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryInTxHashToCctxDataResponse",
      value: QueryInTxHashToCctxDataResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllInTxHashToCctxRequest(): QueryAllInTxHashToCctxRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllInTxHashToCctxRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxHashToCctxRequest",
  encode(message: QueryAllInTxHashToCctxRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllInTxHashToCctxRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllInTxHashToCctxRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllInTxHashToCctxRequest>): QueryAllInTxHashToCctxRequest {
    const message = createBaseQueryAllInTxHashToCctxRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllInTxHashToCctxRequestAmino): QueryAllInTxHashToCctxRequest {
    const message = createBaseQueryAllInTxHashToCctxRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllInTxHashToCctxRequest): QueryAllInTxHashToCctxRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllInTxHashToCctxRequestAminoMsg): QueryAllInTxHashToCctxRequest {
    return QueryAllInTxHashToCctxRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllInTxHashToCctxRequestProtoMsg): QueryAllInTxHashToCctxRequest {
    return QueryAllInTxHashToCctxRequest.decode(message.value);
  },
  toProto(message: QueryAllInTxHashToCctxRequest): Uint8Array {
    return QueryAllInTxHashToCctxRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllInTxHashToCctxRequest): QueryAllInTxHashToCctxRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxHashToCctxRequest",
      value: QueryAllInTxHashToCctxRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllInTxHashToCctxResponse(): QueryAllInTxHashToCctxResponse {
  return {
    inTxHashToCctx: [],
    pagination: undefined
  };
}
export const QueryAllInTxHashToCctxResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxHashToCctxResponse",
  encode(message: QueryAllInTxHashToCctxResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.inTxHashToCctx) {
      InTxHashToCctx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllInTxHashToCctxResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllInTxHashToCctxResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.inTxHashToCctx.push(InTxHashToCctx.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllInTxHashToCctxResponse>): QueryAllInTxHashToCctxResponse {
    const message = createBaseQueryAllInTxHashToCctxResponse();
    message.inTxHashToCctx = object.inTxHashToCctx?.map(e => InTxHashToCctx.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllInTxHashToCctxResponseAmino): QueryAllInTxHashToCctxResponse {
    const message = createBaseQueryAllInTxHashToCctxResponse();
    message.inTxHashToCctx = object.inTxHashToCctx?.map(e => InTxHashToCctx.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllInTxHashToCctxResponse): QueryAllInTxHashToCctxResponseAmino {
    const obj: any = {};
    if (message.inTxHashToCctx) {
      obj.inTxHashToCctx = message.inTxHashToCctx.map(e => e ? InTxHashToCctx.toAmino(e) : undefined);
    } else {
      obj.inTxHashToCctx = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllInTxHashToCctxResponseAminoMsg): QueryAllInTxHashToCctxResponse {
    return QueryAllInTxHashToCctxResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllInTxHashToCctxResponseProtoMsg): QueryAllInTxHashToCctxResponse {
    return QueryAllInTxHashToCctxResponse.decode(message.value);
  },
  toProto(message: QueryAllInTxHashToCctxResponse): Uint8Array {
    return QueryAllInTxHashToCctxResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllInTxHashToCctxResponse): QueryAllInTxHashToCctxResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllInTxHashToCctxResponse",
      value: QueryAllInTxHashToCctxResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetGasPriceRequest(): QueryGetGasPriceRequest {
  return {
    index: ""
  };
}
export const QueryGetGasPriceRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetGasPriceRequest",
  encode(message: QueryGetGasPriceRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetGasPriceRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGasPriceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetGasPriceRequest>): QueryGetGasPriceRequest {
    const message = createBaseQueryGetGasPriceRequest();
    message.index = object.index ?? "";
    return message;
  },
  fromAmino(object: QueryGetGasPriceRequestAmino): QueryGetGasPriceRequest {
    const message = createBaseQueryGetGasPriceRequest();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    return message;
  },
  toAmino(message: QueryGetGasPriceRequest): QueryGetGasPriceRequestAmino {
    const obj: any = {};
    obj.index = message.index;
    return obj;
  },
  fromAminoMsg(object: QueryGetGasPriceRequestAminoMsg): QueryGetGasPriceRequest {
    return QueryGetGasPriceRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetGasPriceRequestProtoMsg): QueryGetGasPriceRequest {
    return QueryGetGasPriceRequest.decode(message.value);
  },
  toProto(message: QueryGetGasPriceRequest): Uint8Array {
    return QueryGetGasPriceRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetGasPriceRequest): QueryGetGasPriceRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetGasPriceRequest",
      value: QueryGetGasPriceRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetGasPriceResponse(): QueryGetGasPriceResponse {
  return {
    GasPrice: undefined
  };
}
export const QueryGetGasPriceResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetGasPriceResponse",
  encode(message: QueryGetGasPriceResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.GasPrice !== undefined) {
      GasPrice.encode(message.GasPrice, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetGasPriceResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGasPriceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.GasPrice = GasPrice.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetGasPriceResponse>): QueryGetGasPriceResponse {
    const message = createBaseQueryGetGasPriceResponse();
    message.GasPrice = object.GasPrice !== undefined && object.GasPrice !== null ? GasPrice.fromPartial(object.GasPrice) : undefined;
    return message;
  },
  fromAmino(object: QueryGetGasPriceResponseAmino): QueryGetGasPriceResponse {
    const message = createBaseQueryGetGasPriceResponse();
    if (object.GasPrice !== undefined && object.GasPrice !== null) {
      message.GasPrice = GasPrice.fromAmino(object.GasPrice);
    }
    return message;
  },
  toAmino(message: QueryGetGasPriceResponse): QueryGetGasPriceResponseAmino {
    const obj: any = {};
    obj.GasPrice = message.GasPrice ? GasPrice.toAmino(message.GasPrice) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetGasPriceResponseAminoMsg): QueryGetGasPriceResponse {
    return QueryGetGasPriceResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetGasPriceResponseProtoMsg): QueryGetGasPriceResponse {
    return QueryGetGasPriceResponse.decode(message.value);
  },
  toProto(message: QueryGetGasPriceResponse): Uint8Array {
    return QueryGetGasPriceResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetGasPriceResponse): QueryGetGasPriceResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetGasPriceResponse",
      value: QueryGetGasPriceResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllGasPriceRequest(): QueryAllGasPriceRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllGasPriceRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllGasPriceRequest",
  encode(message: QueryAllGasPriceRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllGasPriceRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllGasPriceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllGasPriceRequest>): QueryAllGasPriceRequest {
    const message = createBaseQueryAllGasPriceRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllGasPriceRequestAmino): QueryAllGasPriceRequest {
    const message = createBaseQueryAllGasPriceRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllGasPriceRequest): QueryAllGasPriceRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllGasPriceRequestAminoMsg): QueryAllGasPriceRequest {
    return QueryAllGasPriceRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllGasPriceRequestProtoMsg): QueryAllGasPriceRequest {
    return QueryAllGasPriceRequest.decode(message.value);
  },
  toProto(message: QueryAllGasPriceRequest): Uint8Array {
    return QueryAllGasPriceRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllGasPriceRequest): QueryAllGasPriceRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllGasPriceRequest",
      value: QueryAllGasPriceRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllGasPriceResponse(): QueryAllGasPriceResponse {
  return {
    GasPrice: [],
    pagination: undefined
  };
}
export const QueryAllGasPriceResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllGasPriceResponse",
  encode(message: QueryAllGasPriceResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.GasPrice) {
      GasPrice.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllGasPriceResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllGasPriceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.GasPrice.push(GasPrice.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllGasPriceResponse>): QueryAllGasPriceResponse {
    const message = createBaseQueryAllGasPriceResponse();
    message.GasPrice = object.GasPrice?.map(e => GasPrice.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllGasPriceResponseAmino): QueryAllGasPriceResponse {
    const message = createBaseQueryAllGasPriceResponse();
    message.GasPrice = object.GasPrice?.map(e => GasPrice.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllGasPriceResponse): QueryAllGasPriceResponseAmino {
    const obj: any = {};
    if (message.GasPrice) {
      obj.GasPrice = message.GasPrice.map(e => e ? GasPrice.toAmino(e) : undefined);
    } else {
      obj.GasPrice = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllGasPriceResponseAminoMsg): QueryAllGasPriceResponse {
    return QueryAllGasPriceResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllGasPriceResponseProtoMsg): QueryAllGasPriceResponse {
    return QueryAllGasPriceResponse.decode(message.value);
  },
  toProto(message: QueryAllGasPriceResponse): Uint8Array {
    return QueryAllGasPriceResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllGasPriceResponse): QueryAllGasPriceResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllGasPriceResponse",
      value: QueryAllGasPriceResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetLastBlockHeightRequest(): QueryGetLastBlockHeightRequest {
  return {
    index: ""
  };
}
export const QueryGetLastBlockHeightRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetLastBlockHeightRequest",
  encode(message: QueryGetLastBlockHeightRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetLastBlockHeightRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetLastBlockHeightRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetLastBlockHeightRequest>): QueryGetLastBlockHeightRequest {
    const message = createBaseQueryGetLastBlockHeightRequest();
    message.index = object.index ?? "";
    return message;
  },
  fromAmino(object: QueryGetLastBlockHeightRequestAmino): QueryGetLastBlockHeightRequest {
    const message = createBaseQueryGetLastBlockHeightRequest();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    return message;
  },
  toAmino(message: QueryGetLastBlockHeightRequest): QueryGetLastBlockHeightRequestAmino {
    const obj: any = {};
    obj.index = message.index;
    return obj;
  },
  fromAminoMsg(object: QueryGetLastBlockHeightRequestAminoMsg): QueryGetLastBlockHeightRequest {
    return QueryGetLastBlockHeightRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetLastBlockHeightRequestProtoMsg): QueryGetLastBlockHeightRequest {
    return QueryGetLastBlockHeightRequest.decode(message.value);
  },
  toProto(message: QueryGetLastBlockHeightRequest): Uint8Array {
    return QueryGetLastBlockHeightRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetLastBlockHeightRequest): QueryGetLastBlockHeightRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetLastBlockHeightRequest",
      value: QueryGetLastBlockHeightRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetLastBlockHeightResponse(): QueryGetLastBlockHeightResponse {
  return {
    LastBlockHeight: undefined
  };
}
export const QueryGetLastBlockHeightResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetLastBlockHeightResponse",
  encode(message: QueryGetLastBlockHeightResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.LastBlockHeight !== undefined) {
      LastBlockHeight.encode(message.LastBlockHeight, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetLastBlockHeightResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetLastBlockHeightResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.LastBlockHeight = LastBlockHeight.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetLastBlockHeightResponse>): QueryGetLastBlockHeightResponse {
    const message = createBaseQueryGetLastBlockHeightResponse();
    message.LastBlockHeight = object.LastBlockHeight !== undefined && object.LastBlockHeight !== null ? LastBlockHeight.fromPartial(object.LastBlockHeight) : undefined;
    return message;
  },
  fromAmino(object: QueryGetLastBlockHeightResponseAmino): QueryGetLastBlockHeightResponse {
    const message = createBaseQueryGetLastBlockHeightResponse();
    if (object.LastBlockHeight !== undefined && object.LastBlockHeight !== null) {
      message.LastBlockHeight = LastBlockHeight.fromAmino(object.LastBlockHeight);
    }
    return message;
  },
  toAmino(message: QueryGetLastBlockHeightResponse): QueryGetLastBlockHeightResponseAmino {
    const obj: any = {};
    obj.LastBlockHeight = message.LastBlockHeight ? LastBlockHeight.toAmino(message.LastBlockHeight) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetLastBlockHeightResponseAminoMsg): QueryGetLastBlockHeightResponse {
    return QueryGetLastBlockHeightResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetLastBlockHeightResponseProtoMsg): QueryGetLastBlockHeightResponse {
    return QueryGetLastBlockHeightResponse.decode(message.value);
  },
  toProto(message: QueryGetLastBlockHeightResponse): Uint8Array {
    return QueryGetLastBlockHeightResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetLastBlockHeightResponse): QueryGetLastBlockHeightResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetLastBlockHeightResponse",
      value: QueryGetLastBlockHeightResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllLastBlockHeightRequest(): QueryAllLastBlockHeightRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllLastBlockHeightRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllLastBlockHeightRequest",
  encode(message: QueryAllLastBlockHeightRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllLastBlockHeightRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllLastBlockHeightRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllLastBlockHeightRequest>): QueryAllLastBlockHeightRequest {
    const message = createBaseQueryAllLastBlockHeightRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllLastBlockHeightRequestAmino): QueryAllLastBlockHeightRequest {
    const message = createBaseQueryAllLastBlockHeightRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllLastBlockHeightRequest): QueryAllLastBlockHeightRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllLastBlockHeightRequestAminoMsg): QueryAllLastBlockHeightRequest {
    return QueryAllLastBlockHeightRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllLastBlockHeightRequestProtoMsg): QueryAllLastBlockHeightRequest {
    return QueryAllLastBlockHeightRequest.decode(message.value);
  },
  toProto(message: QueryAllLastBlockHeightRequest): Uint8Array {
    return QueryAllLastBlockHeightRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllLastBlockHeightRequest): QueryAllLastBlockHeightRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllLastBlockHeightRequest",
      value: QueryAllLastBlockHeightRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllLastBlockHeightResponse(): QueryAllLastBlockHeightResponse {
  return {
    LastBlockHeight: [],
    pagination: undefined
  };
}
export const QueryAllLastBlockHeightResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllLastBlockHeightResponse",
  encode(message: QueryAllLastBlockHeightResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.LastBlockHeight) {
      LastBlockHeight.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllLastBlockHeightResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllLastBlockHeightResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.LastBlockHeight.push(LastBlockHeight.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllLastBlockHeightResponse>): QueryAllLastBlockHeightResponse {
    const message = createBaseQueryAllLastBlockHeightResponse();
    message.LastBlockHeight = object.LastBlockHeight?.map(e => LastBlockHeight.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllLastBlockHeightResponseAmino): QueryAllLastBlockHeightResponse {
    const message = createBaseQueryAllLastBlockHeightResponse();
    message.LastBlockHeight = object.LastBlockHeight?.map(e => LastBlockHeight.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllLastBlockHeightResponse): QueryAllLastBlockHeightResponseAmino {
    const obj: any = {};
    if (message.LastBlockHeight) {
      obj.LastBlockHeight = message.LastBlockHeight.map(e => e ? LastBlockHeight.toAmino(e) : undefined);
    } else {
      obj.LastBlockHeight = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllLastBlockHeightResponseAminoMsg): QueryAllLastBlockHeightResponse {
    return QueryAllLastBlockHeightResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllLastBlockHeightResponseProtoMsg): QueryAllLastBlockHeightResponse {
    return QueryAllLastBlockHeightResponse.decode(message.value);
  },
  toProto(message: QueryAllLastBlockHeightResponse): Uint8Array {
    return QueryAllLastBlockHeightResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllLastBlockHeightResponse): QueryAllLastBlockHeightResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllLastBlockHeightResponse",
      value: QueryAllLastBlockHeightResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetCctxRequest(): QueryGetCctxRequest {
  return {
    index: ""
  };
}
export const QueryGetCctxRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetCctxRequest",
  encode(message: QueryGetCctxRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetCctxRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetCctxRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetCctxRequest>): QueryGetCctxRequest {
    const message = createBaseQueryGetCctxRequest();
    message.index = object.index ?? "";
    return message;
  },
  fromAmino(object: QueryGetCctxRequestAmino): QueryGetCctxRequest {
    const message = createBaseQueryGetCctxRequest();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    return message;
  },
  toAmino(message: QueryGetCctxRequest): QueryGetCctxRequestAmino {
    const obj: any = {};
    obj.index = message.index;
    return obj;
  },
  fromAminoMsg(object: QueryGetCctxRequestAminoMsg): QueryGetCctxRequest {
    return QueryGetCctxRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetCctxRequestProtoMsg): QueryGetCctxRequest {
    return QueryGetCctxRequest.decode(message.value);
  },
  toProto(message: QueryGetCctxRequest): Uint8Array {
    return QueryGetCctxRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetCctxRequest): QueryGetCctxRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetCctxRequest",
      value: QueryGetCctxRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetCctxByNonceRequest(): QueryGetCctxByNonceRequest {
  return {
    chainID: BigInt(0),
    nonce: BigInt(0)
  };
}
export const QueryGetCctxByNonceRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetCctxByNonceRequest",
  encode(message: QueryGetCctxByNonceRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainID !== BigInt(0)) {
      writer.uint32(8).int64(message.chainID);
    }
    if (message.nonce !== BigInt(0)) {
      writer.uint32(16).uint64(message.nonce);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetCctxByNonceRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetCctxByNonceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainID = reader.int64();
          break;
        case 2:
          message.nonce = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetCctxByNonceRequest>): QueryGetCctxByNonceRequest {
    const message = createBaseQueryGetCctxByNonceRequest();
    message.chainID = object.chainID !== undefined && object.chainID !== null ? BigInt(object.chainID.toString()) : BigInt(0);
    message.nonce = object.nonce !== undefined && object.nonce !== null ? BigInt(object.nonce.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryGetCctxByNonceRequestAmino): QueryGetCctxByNonceRequest {
    const message = createBaseQueryGetCctxByNonceRequest();
    if (object.chainID !== undefined && object.chainID !== null) {
      message.chainID = BigInt(object.chainID);
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = BigInt(object.nonce);
    }
    return message;
  },
  toAmino(message: QueryGetCctxByNonceRequest): QueryGetCctxByNonceRequestAmino {
    const obj: any = {};
    obj.chainID = message.chainID ? message.chainID.toString() : undefined;
    obj.nonce = message.nonce ? message.nonce.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetCctxByNonceRequestAminoMsg): QueryGetCctxByNonceRequest {
    return QueryGetCctxByNonceRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetCctxByNonceRequestProtoMsg): QueryGetCctxByNonceRequest {
    return QueryGetCctxByNonceRequest.decode(message.value);
  },
  toProto(message: QueryGetCctxByNonceRequest): Uint8Array {
    return QueryGetCctxByNonceRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetCctxByNonceRequest): QueryGetCctxByNonceRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetCctxByNonceRequest",
      value: QueryGetCctxByNonceRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetCctxResponse(): QueryGetCctxResponse {
  return {
    CrossChainTx: undefined
  };
}
export const QueryGetCctxResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryGetCctxResponse",
  encode(message: QueryGetCctxResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.CrossChainTx !== undefined) {
      CrossChainTx.encode(message.CrossChainTx, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetCctxResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetCctxResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.CrossChainTx = CrossChainTx.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetCctxResponse>): QueryGetCctxResponse {
    const message = createBaseQueryGetCctxResponse();
    message.CrossChainTx = object.CrossChainTx !== undefined && object.CrossChainTx !== null ? CrossChainTx.fromPartial(object.CrossChainTx) : undefined;
    return message;
  },
  fromAmino(object: QueryGetCctxResponseAmino): QueryGetCctxResponse {
    const message = createBaseQueryGetCctxResponse();
    if (object.CrossChainTx !== undefined && object.CrossChainTx !== null) {
      message.CrossChainTx = CrossChainTx.fromAmino(object.CrossChainTx);
    }
    return message;
  },
  toAmino(message: QueryGetCctxResponse): QueryGetCctxResponseAmino {
    const obj: any = {};
    obj.CrossChainTx = message.CrossChainTx ? CrossChainTx.toAmino(message.CrossChainTx) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetCctxResponseAminoMsg): QueryGetCctxResponse {
    return QueryGetCctxResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetCctxResponseProtoMsg): QueryGetCctxResponse {
    return QueryGetCctxResponse.decode(message.value);
  },
  toProto(message: QueryGetCctxResponse): Uint8Array {
    return QueryGetCctxResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetCctxResponse): QueryGetCctxResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryGetCctxResponse",
      value: QueryGetCctxResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllCctxRequest(): QueryAllCctxRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllCctxRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllCctxRequest",
  encode(message: QueryAllCctxRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllCctxRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllCctxRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllCctxRequest>): QueryAllCctxRequest {
    const message = createBaseQueryAllCctxRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllCctxRequestAmino): QueryAllCctxRequest {
    const message = createBaseQueryAllCctxRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllCctxRequest): QueryAllCctxRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllCctxRequestAminoMsg): QueryAllCctxRequest {
    return QueryAllCctxRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllCctxRequestProtoMsg): QueryAllCctxRequest {
    return QueryAllCctxRequest.decode(message.value);
  },
  toProto(message: QueryAllCctxRequest): Uint8Array {
    return QueryAllCctxRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllCctxRequest): QueryAllCctxRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllCctxRequest",
      value: QueryAllCctxRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllCctxResponse(): QueryAllCctxResponse {
  return {
    CrossChainTx: [],
    pagination: undefined
  };
}
export const QueryAllCctxResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryAllCctxResponse",
  encode(message: QueryAllCctxResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.CrossChainTx) {
      CrossChainTx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllCctxResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllCctxResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.CrossChainTx.push(CrossChainTx.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllCctxResponse>): QueryAllCctxResponse {
    const message = createBaseQueryAllCctxResponse();
    message.CrossChainTx = object.CrossChainTx?.map(e => CrossChainTx.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllCctxResponseAmino): QueryAllCctxResponse {
    const message = createBaseQueryAllCctxResponse();
    message.CrossChainTx = object.CrossChainTx?.map(e => CrossChainTx.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllCctxResponse): QueryAllCctxResponseAmino {
    const obj: any = {};
    if (message.CrossChainTx) {
      obj.CrossChainTx = message.CrossChainTx.map(e => e ? CrossChainTx.toAmino(e) : undefined);
    } else {
      obj.CrossChainTx = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllCctxResponseAminoMsg): QueryAllCctxResponse {
    return QueryAllCctxResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllCctxResponseProtoMsg): QueryAllCctxResponse {
    return QueryAllCctxResponse.decode(message.value);
  },
  toProto(message: QueryAllCctxResponse): Uint8Array {
    return QueryAllCctxResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllCctxResponse): QueryAllCctxResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryAllCctxResponse",
      value: QueryAllCctxResponse.encode(message).finish()
    };
  }
};
function createBaseQueryListCctxPendingRequest(): QueryListCctxPendingRequest {
  return {
    chainId: BigInt(0),
    limit: 0
  };
}
export const QueryListCctxPendingRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryListCctxPendingRequest",
  encode(message: QueryListCctxPendingRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.limit !== 0) {
      writer.uint32(16).uint32(message.limit);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryListCctxPendingRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryListCctxPendingRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.limit = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryListCctxPendingRequest>): QueryListCctxPendingRequest {
    const message = createBaseQueryListCctxPendingRequest();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.limit = object.limit ?? 0;
    return message;
  },
  fromAmino(object: QueryListCctxPendingRequestAmino): QueryListCctxPendingRequest {
    const message = createBaseQueryListCctxPendingRequest();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.limit !== undefined && object.limit !== null) {
      message.limit = object.limit;
    }
    return message;
  },
  toAmino(message: QueryListCctxPendingRequest): QueryListCctxPendingRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.limit = message.limit;
    return obj;
  },
  fromAminoMsg(object: QueryListCctxPendingRequestAminoMsg): QueryListCctxPendingRequest {
    return QueryListCctxPendingRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryListCctxPendingRequestProtoMsg): QueryListCctxPendingRequest {
    return QueryListCctxPendingRequest.decode(message.value);
  },
  toProto(message: QueryListCctxPendingRequest): Uint8Array {
    return QueryListCctxPendingRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryListCctxPendingRequest): QueryListCctxPendingRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryListCctxPendingRequest",
      value: QueryListCctxPendingRequest.encode(message).finish()
    };
  }
};
function createBaseQueryListCctxPendingResponse(): QueryListCctxPendingResponse {
  return {
    CrossChainTx: [],
    totalPending: BigInt(0)
  };
}
export const QueryListCctxPendingResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryListCctxPendingResponse",
  encode(message: QueryListCctxPendingResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.CrossChainTx) {
      CrossChainTx.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.totalPending !== BigInt(0)) {
      writer.uint32(16).uint64(message.totalPending);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryListCctxPendingResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryListCctxPendingResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.CrossChainTx.push(CrossChainTx.decode(reader, reader.uint32()));
          break;
        case 2:
          message.totalPending = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryListCctxPendingResponse>): QueryListCctxPendingResponse {
    const message = createBaseQueryListCctxPendingResponse();
    message.CrossChainTx = object.CrossChainTx?.map(e => CrossChainTx.fromPartial(e)) || [];
    message.totalPending = object.totalPending !== undefined && object.totalPending !== null ? BigInt(object.totalPending.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryListCctxPendingResponseAmino): QueryListCctxPendingResponse {
    const message = createBaseQueryListCctxPendingResponse();
    message.CrossChainTx = object.CrossChainTx?.map(e => CrossChainTx.fromAmino(e)) || [];
    if (object.totalPending !== undefined && object.totalPending !== null) {
      message.totalPending = BigInt(object.totalPending);
    }
    return message;
  },
  toAmino(message: QueryListCctxPendingResponse): QueryListCctxPendingResponseAmino {
    const obj: any = {};
    if (message.CrossChainTx) {
      obj.CrossChainTx = message.CrossChainTx.map(e => e ? CrossChainTx.toAmino(e) : undefined);
    } else {
      obj.CrossChainTx = [];
    }
    obj.totalPending = message.totalPending ? message.totalPending.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryListCctxPendingResponseAminoMsg): QueryListCctxPendingResponse {
    return QueryListCctxPendingResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryListCctxPendingResponseProtoMsg): QueryListCctxPendingResponse {
    return QueryListCctxPendingResponse.decode(message.value);
  },
  toProto(message: QueryListCctxPendingResponse): Uint8Array {
    return QueryListCctxPendingResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryListCctxPendingResponse): QueryListCctxPendingResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryListCctxPendingResponse",
      value: QueryListCctxPendingResponse.encode(message).finish()
    };
  }
};
function createBaseQueryLastZetaHeightRequest(): QueryLastZetaHeightRequest {
  return {};
}
export const QueryLastZetaHeightRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryLastZetaHeightRequest",
  encode(_: QueryLastZetaHeightRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryLastZetaHeightRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryLastZetaHeightRequest();
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
  fromPartial(_: Partial<QueryLastZetaHeightRequest>): QueryLastZetaHeightRequest {
    const message = createBaseQueryLastZetaHeightRequest();
    return message;
  },
  fromAmino(_: QueryLastZetaHeightRequestAmino): QueryLastZetaHeightRequest {
    const message = createBaseQueryLastZetaHeightRequest();
    return message;
  },
  toAmino(_: QueryLastZetaHeightRequest): QueryLastZetaHeightRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryLastZetaHeightRequestAminoMsg): QueryLastZetaHeightRequest {
    return QueryLastZetaHeightRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryLastZetaHeightRequestProtoMsg): QueryLastZetaHeightRequest {
    return QueryLastZetaHeightRequest.decode(message.value);
  },
  toProto(message: QueryLastZetaHeightRequest): Uint8Array {
    return QueryLastZetaHeightRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryLastZetaHeightRequest): QueryLastZetaHeightRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryLastZetaHeightRequest",
      value: QueryLastZetaHeightRequest.encode(message).finish()
    };
  }
};
function createBaseQueryLastZetaHeightResponse(): QueryLastZetaHeightResponse {
  return {
    Height: BigInt(0)
  };
}
export const QueryLastZetaHeightResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryLastZetaHeightResponse",
  encode(message: QueryLastZetaHeightResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.Height !== BigInt(0)) {
      writer.uint32(8).int64(message.Height);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryLastZetaHeightResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryLastZetaHeightResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Height = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryLastZetaHeightResponse>): QueryLastZetaHeightResponse {
    const message = createBaseQueryLastZetaHeightResponse();
    message.Height = object.Height !== undefined && object.Height !== null ? BigInt(object.Height.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryLastZetaHeightResponseAmino): QueryLastZetaHeightResponse {
    const message = createBaseQueryLastZetaHeightResponse();
    if (object.Height !== undefined && object.Height !== null) {
      message.Height = BigInt(object.Height);
    }
    return message;
  },
  toAmino(message: QueryLastZetaHeightResponse): QueryLastZetaHeightResponseAmino {
    const obj: any = {};
    obj.Height = message.Height ? message.Height.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryLastZetaHeightResponseAminoMsg): QueryLastZetaHeightResponse {
    return QueryLastZetaHeightResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryLastZetaHeightResponseProtoMsg): QueryLastZetaHeightResponse {
    return QueryLastZetaHeightResponse.decode(message.value);
  },
  toProto(message: QueryLastZetaHeightResponse): Uint8Array {
    return QueryLastZetaHeightResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryLastZetaHeightResponse): QueryLastZetaHeightResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryLastZetaHeightResponse",
      value: QueryLastZetaHeightResponse.encode(message).finish()
    };
  }
};
function createBaseQueryConvertGasToZetaRequest(): QueryConvertGasToZetaRequest {
  return {
    chainId: BigInt(0),
    gasLimit: ""
  };
}
export const QueryConvertGasToZetaRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryConvertGasToZetaRequest",
  encode(message: QueryConvertGasToZetaRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.gasLimit !== "") {
      writer.uint32(18).string(message.gasLimit);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryConvertGasToZetaRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryConvertGasToZetaRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.gasLimit = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryConvertGasToZetaRequest>): QueryConvertGasToZetaRequest {
    const message = createBaseQueryConvertGasToZetaRequest();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.gasLimit = object.gasLimit ?? "";
    return message;
  },
  fromAmino(object: QueryConvertGasToZetaRequestAmino): QueryConvertGasToZetaRequest {
    const message = createBaseQueryConvertGasToZetaRequest();
    if (object.chainId !== undefined && object.chainId !== null) {
      message.chainId = BigInt(object.chainId);
    }
    if (object.gasLimit !== undefined && object.gasLimit !== null) {
      message.gasLimit = object.gasLimit;
    }
    return message;
  },
  toAmino(message: QueryConvertGasToZetaRequest): QueryConvertGasToZetaRequestAmino {
    const obj: any = {};
    obj.chainId = message.chainId ? message.chainId.toString() : undefined;
    obj.gasLimit = message.gasLimit;
    return obj;
  },
  fromAminoMsg(object: QueryConvertGasToZetaRequestAminoMsg): QueryConvertGasToZetaRequest {
    return QueryConvertGasToZetaRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryConvertGasToZetaRequestProtoMsg): QueryConvertGasToZetaRequest {
    return QueryConvertGasToZetaRequest.decode(message.value);
  },
  toProto(message: QueryConvertGasToZetaRequest): Uint8Array {
    return QueryConvertGasToZetaRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryConvertGasToZetaRequest): QueryConvertGasToZetaRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryConvertGasToZetaRequest",
      value: QueryConvertGasToZetaRequest.encode(message).finish()
    };
  }
};
function createBaseQueryConvertGasToZetaResponse(): QueryConvertGasToZetaResponse {
  return {
    outboundGasInZeta: "",
    protocolFeeInZeta: "",
    ZetaBlockHeight: BigInt(0)
  };
}
export const QueryConvertGasToZetaResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryConvertGasToZetaResponse",
  encode(message: QueryConvertGasToZetaResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.outboundGasInZeta !== "") {
      writer.uint32(10).string(message.outboundGasInZeta);
    }
    if (message.protocolFeeInZeta !== "") {
      writer.uint32(18).string(message.protocolFeeInZeta);
    }
    if (message.ZetaBlockHeight !== BigInt(0)) {
      writer.uint32(24).uint64(message.ZetaBlockHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryConvertGasToZetaResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryConvertGasToZetaResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.outboundGasInZeta = reader.string();
          break;
        case 2:
          message.protocolFeeInZeta = reader.string();
          break;
        case 3:
          message.ZetaBlockHeight = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryConvertGasToZetaResponse>): QueryConvertGasToZetaResponse {
    const message = createBaseQueryConvertGasToZetaResponse();
    message.outboundGasInZeta = object.outboundGasInZeta ?? "";
    message.protocolFeeInZeta = object.protocolFeeInZeta ?? "";
    message.ZetaBlockHeight = object.ZetaBlockHeight !== undefined && object.ZetaBlockHeight !== null ? BigInt(object.ZetaBlockHeight.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryConvertGasToZetaResponseAmino): QueryConvertGasToZetaResponse {
    const message = createBaseQueryConvertGasToZetaResponse();
    if (object.outboundGasInZeta !== undefined && object.outboundGasInZeta !== null) {
      message.outboundGasInZeta = object.outboundGasInZeta;
    }
    if (object.protocolFeeInZeta !== undefined && object.protocolFeeInZeta !== null) {
      message.protocolFeeInZeta = object.protocolFeeInZeta;
    }
    if (object.ZetaBlockHeight !== undefined && object.ZetaBlockHeight !== null) {
      message.ZetaBlockHeight = BigInt(object.ZetaBlockHeight);
    }
    return message;
  },
  toAmino(message: QueryConvertGasToZetaResponse): QueryConvertGasToZetaResponseAmino {
    const obj: any = {};
    obj.outboundGasInZeta = message.outboundGasInZeta;
    obj.protocolFeeInZeta = message.protocolFeeInZeta;
    obj.ZetaBlockHeight = message.ZetaBlockHeight ? message.ZetaBlockHeight.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryConvertGasToZetaResponseAminoMsg): QueryConvertGasToZetaResponse {
    return QueryConvertGasToZetaResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryConvertGasToZetaResponseProtoMsg): QueryConvertGasToZetaResponse {
    return QueryConvertGasToZetaResponse.decode(message.value);
  },
  toProto(message: QueryConvertGasToZetaResponse): Uint8Array {
    return QueryConvertGasToZetaResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryConvertGasToZetaResponse): QueryConvertGasToZetaResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryConvertGasToZetaResponse",
      value: QueryConvertGasToZetaResponse.encode(message).finish()
    };
  }
};
function createBaseQueryMessagePassingProtocolFeeRequest(): QueryMessagePassingProtocolFeeRequest {
  return {};
}
export const QueryMessagePassingProtocolFeeRequest = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryMessagePassingProtocolFeeRequest",
  encode(_: QueryMessagePassingProtocolFeeRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryMessagePassingProtocolFeeRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMessagePassingProtocolFeeRequest();
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
  fromPartial(_: Partial<QueryMessagePassingProtocolFeeRequest>): QueryMessagePassingProtocolFeeRequest {
    const message = createBaseQueryMessagePassingProtocolFeeRequest();
    return message;
  },
  fromAmino(_: QueryMessagePassingProtocolFeeRequestAmino): QueryMessagePassingProtocolFeeRequest {
    const message = createBaseQueryMessagePassingProtocolFeeRequest();
    return message;
  },
  toAmino(_: QueryMessagePassingProtocolFeeRequest): QueryMessagePassingProtocolFeeRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryMessagePassingProtocolFeeRequestAminoMsg): QueryMessagePassingProtocolFeeRequest {
    return QueryMessagePassingProtocolFeeRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryMessagePassingProtocolFeeRequestProtoMsg): QueryMessagePassingProtocolFeeRequest {
    return QueryMessagePassingProtocolFeeRequest.decode(message.value);
  },
  toProto(message: QueryMessagePassingProtocolFeeRequest): Uint8Array {
    return QueryMessagePassingProtocolFeeRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryMessagePassingProtocolFeeRequest): QueryMessagePassingProtocolFeeRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryMessagePassingProtocolFeeRequest",
      value: QueryMessagePassingProtocolFeeRequest.encode(message).finish()
    };
  }
};
function createBaseQueryMessagePassingProtocolFeeResponse(): QueryMessagePassingProtocolFeeResponse {
  return {
    feeInZeta: ""
  };
}
export const QueryMessagePassingProtocolFeeResponse = {
  typeUrl: "/zetachain.zetacore.crosschain.QueryMessagePassingProtocolFeeResponse",
  encode(message: QueryMessagePassingProtocolFeeResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.feeInZeta !== "") {
      writer.uint32(10).string(message.feeInZeta);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryMessagePassingProtocolFeeResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMessagePassingProtocolFeeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.feeInZeta = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryMessagePassingProtocolFeeResponse>): QueryMessagePassingProtocolFeeResponse {
    const message = createBaseQueryMessagePassingProtocolFeeResponse();
    message.feeInZeta = object.feeInZeta ?? "";
    return message;
  },
  fromAmino(object: QueryMessagePassingProtocolFeeResponseAmino): QueryMessagePassingProtocolFeeResponse {
    const message = createBaseQueryMessagePassingProtocolFeeResponse();
    if (object.feeInZeta !== undefined && object.feeInZeta !== null) {
      message.feeInZeta = object.feeInZeta;
    }
    return message;
  },
  toAmino(message: QueryMessagePassingProtocolFeeResponse): QueryMessagePassingProtocolFeeResponseAmino {
    const obj: any = {};
    obj.feeInZeta = message.feeInZeta;
    return obj;
  },
  fromAminoMsg(object: QueryMessagePassingProtocolFeeResponseAminoMsg): QueryMessagePassingProtocolFeeResponse {
    return QueryMessagePassingProtocolFeeResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryMessagePassingProtocolFeeResponseProtoMsg): QueryMessagePassingProtocolFeeResponse {
    return QueryMessagePassingProtocolFeeResponse.decode(message.value);
  },
  toProto(message: QueryMessagePassingProtocolFeeResponse): Uint8Array {
    return QueryMessagePassingProtocolFeeResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryMessagePassingProtocolFeeResponse): QueryMessagePassingProtocolFeeResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.crosschain.QueryMessagePassingProtocolFeeResponse",
      value: QueryMessagePassingProtocolFeeResponse.encode(message).finish()
    };
  }
};