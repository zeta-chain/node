import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Proof, ProofAmino, ProofSDKType, Chain, ChainAmino, ChainSDKType, BlockHeader, BlockHeaderAmino, BlockHeaderSDKType } from "../common/common";
import { ChainNonces, ChainNoncesAmino, ChainNoncesSDKType } from "./chain_nonces";
import { PendingNonces, PendingNoncesAmino, PendingNoncesSDKType } from "./pending_nonces";
import { TSS, TSSAmino, TSSSDKType } from "./tss";
import { Params, ParamsAmino, ParamsSDKType, ChainParams, ChainParamsAmino, ChainParamsSDKType, ChainParamsList, ChainParamsListAmino, ChainParamsListSDKType } from "./params";
import { VoteType, BallotStatus, voteTypeFromJSON, ballotStatusFromJSON } from "./ballot";
import { ObservationType, LastObserverCount, LastObserverCountAmino, LastObserverCountSDKType, observationTypeFromJSON } from "./observer";
import { NodeAccount, NodeAccountAmino, NodeAccountSDKType } from "./node_account";
import { CrosschainFlags, CrosschainFlagsAmino, CrosschainFlagsSDKType } from "./crosschain_flags";
import { Keygen, KeygenAmino, KeygenSDKType } from "./keygen";
import { Blame, BlameAmino, BlameSDKType } from "./blame";
import { BlockHeaderState, BlockHeaderStateAmino, BlockHeaderStateSDKType } from "./block_header";
import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
export interface QueryGetChainNoncesRequest {
  index: string;
}
export interface QueryGetChainNoncesRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainNoncesRequest";
  value: Uint8Array;
}
export interface QueryGetChainNoncesRequestAmino {
  index?: string;
}
export interface QueryGetChainNoncesRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetChainNoncesRequest";
  value: QueryGetChainNoncesRequestAmino;
}
export interface QueryGetChainNoncesRequestSDKType {
  index: string;
}
export interface QueryGetChainNoncesResponse {
  ChainNonces: ChainNonces;
}
export interface QueryGetChainNoncesResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainNoncesResponse";
  value: Uint8Array;
}
export interface QueryGetChainNoncesResponseAmino {
  ChainNonces?: ChainNoncesAmino;
}
export interface QueryGetChainNoncesResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetChainNoncesResponse";
  value: QueryGetChainNoncesResponseAmino;
}
export interface QueryGetChainNoncesResponseSDKType {
  ChainNonces: ChainNoncesSDKType;
}
export interface QueryAllChainNoncesRequest {
  pagination?: PageRequest;
}
export interface QueryAllChainNoncesRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllChainNoncesRequest";
  value: Uint8Array;
}
export interface QueryAllChainNoncesRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllChainNoncesRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllChainNoncesRequest";
  value: QueryAllChainNoncesRequestAmino;
}
export interface QueryAllChainNoncesRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllChainNoncesResponse {
  ChainNonces: ChainNonces[];
  pagination?: PageResponse;
}
export interface QueryAllChainNoncesResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllChainNoncesResponse";
  value: Uint8Array;
}
export interface QueryAllChainNoncesResponseAmino {
  ChainNonces?: ChainNoncesAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllChainNoncesResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllChainNoncesResponse";
  value: QueryAllChainNoncesResponseAmino;
}
export interface QueryAllChainNoncesResponseSDKType {
  ChainNonces: ChainNoncesSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryAllPendingNoncesRequest {
  pagination?: PageRequest;
}
export interface QueryAllPendingNoncesRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllPendingNoncesRequest";
  value: Uint8Array;
}
export interface QueryAllPendingNoncesRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllPendingNoncesRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllPendingNoncesRequest";
  value: QueryAllPendingNoncesRequestAmino;
}
export interface QueryAllPendingNoncesRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllPendingNoncesResponse {
  pendingNonces: PendingNonces[];
  pagination?: PageResponse;
}
export interface QueryAllPendingNoncesResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllPendingNoncesResponse";
  value: Uint8Array;
}
export interface QueryAllPendingNoncesResponseAmino {
  pending_nonces?: PendingNoncesAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllPendingNoncesResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllPendingNoncesResponse";
  value: QueryAllPendingNoncesResponseAmino;
}
export interface QueryAllPendingNoncesResponseSDKType {
  pending_nonces: PendingNoncesSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryPendingNoncesByChainRequest {
  chainId: bigint;
}
export interface QueryPendingNoncesByChainRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryPendingNoncesByChainRequest";
  value: Uint8Array;
}
export interface QueryPendingNoncesByChainRequestAmino {
  chain_id?: string;
}
export interface QueryPendingNoncesByChainRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryPendingNoncesByChainRequest";
  value: QueryPendingNoncesByChainRequestAmino;
}
export interface QueryPendingNoncesByChainRequestSDKType {
  chain_id: bigint;
}
export interface QueryPendingNoncesByChainResponse {
  pendingNonces: PendingNonces;
}
export interface QueryPendingNoncesByChainResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryPendingNoncesByChainResponse";
  value: Uint8Array;
}
export interface QueryPendingNoncesByChainResponseAmino {
  pending_nonces?: PendingNoncesAmino;
}
export interface QueryPendingNoncesByChainResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryPendingNoncesByChainResponse";
  value: QueryPendingNoncesByChainResponseAmino;
}
export interface QueryPendingNoncesByChainResponseSDKType {
  pending_nonces: PendingNoncesSDKType;
}
export interface QueryGetTSSRequest {}
export interface QueryGetTSSRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTSSRequest";
  value: Uint8Array;
}
export interface QueryGetTSSRequestAmino {}
export interface QueryGetTSSRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetTSSRequest";
  value: QueryGetTSSRequestAmino;
}
export interface QueryGetTSSRequestSDKType {}
export interface QueryGetTSSResponse {
  TSS: TSS;
}
export interface QueryGetTSSResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTSSResponse";
  value: Uint8Array;
}
export interface QueryGetTSSResponseAmino {
  TSS?: TSSAmino;
}
export interface QueryGetTSSResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetTSSResponse";
  value: QueryGetTSSResponseAmino;
}
export interface QueryGetTSSResponseSDKType {
  TSS: TSSSDKType;
}
export interface QueryGetTssAddressRequest {
  bitcoinChainId: bigint;
}
export interface QueryGetTssAddressRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressRequest";
  value: Uint8Array;
}
export interface QueryGetTssAddressRequestAmino {
  bitcoin_chain_id?: string;
}
export interface QueryGetTssAddressRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetTssAddressRequest";
  value: QueryGetTssAddressRequestAmino;
}
export interface QueryGetTssAddressRequestSDKType {
  bitcoin_chain_id: bigint;
}
export interface QueryGetTssAddressResponse {
  eth: string;
  btc: string;
}
export interface QueryGetTssAddressResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressResponse";
  value: Uint8Array;
}
export interface QueryGetTssAddressResponseAmino {
  eth?: string;
  btc?: string;
}
export interface QueryGetTssAddressResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetTssAddressResponse";
  value: QueryGetTssAddressResponseAmino;
}
export interface QueryGetTssAddressResponseSDKType {
  eth: string;
  btc: string;
}
export interface QueryGetTssAddressByFinalizedHeightRequest {
  finalizedZetaHeight: bigint;
  bitcoinChainId: bigint;
}
export interface QueryGetTssAddressByFinalizedHeightRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressByFinalizedHeightRequest";
  value: Uint8Array;
}
export interface QueryGetTssAddressByFinalizedHeightRequestAmino {
  finalized_zeta_height?: string;
  bitcoin_chain_id?: string;
}
export interface QueryGetTssAddressByFinalizedHeightRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetTssAddressByFinalizedHeightRequest";
  value: QueryGetTssAddressByFinalizedHeightRequestAmino;
}
export interface QueryGetTssAddressByFinalizedHeightRequestSDKType {
  finalized_zeta_height: bigint;
  bitcoin_chain_id: bigint;
}
export interface QueryGetTssAddressByFinalizedHeightResponse {
  eth: string;
  btc: string;
}
export interface QueryGetTssAddressByFinalizedHeightResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressByFinalizedHeightResponse";
  value: Uint8Array;
}
export interface QueryGetTssAddressByFinalizedHeightResponseAmino {
  eth?: string;
  btc?: string;
}
export interface QueryGetTssAddressByFinalizedHeightResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetTssAddressByFinalizedHeightResponse";
  value: QueryGetTssAddressByFinalizedHeightResponseAmino;
}
export interface QueryGetTssAddressByFinalizedHeightResponseSDKType {
  eth: string;
  btc: string;
}
export interface QueryTssHistoryRequest {
  pagination?: PageRequest;
}
export interface QueryTssHistoryRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryTssHistoryRequest";
  value: Uint8Array;
}
export interface QueryTssHistoryRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryTssHistoryRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryTssHistoryRequest";
  value: QueryTssHistoryRequestAmino;
}
export interface QueryTssHistoryRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryTssHistoryResponse {
  tssList: TSS[];
  pagination?: PageResponse;
}
export interface QueryTssHistoryResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryTssHistoryResponse";
  value: Uint8Array;
}
export interface QueryTssHistoryResponseAmino {
  tss_list?: TSSAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryTssHistoryResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryTssHistoryResponse";
  value: QueryTssHistoryResponseAmino;
}
export interface QueryTssHistoryResponseSDKType {
  tss_list: TSSSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryProveRequest {
  chainId: bigint;
  txHash: string;
  proof?: Proof;
  blockHash: string;
  txIndex: bigint;
}
export interface QueryProveRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryProveRequest";
  value: Uint8Array;
}
export interface QueryProveRequestAmino {
  chain_id?: string;
  tx_hash?: string;
  proof?: ProofAmino;
  block_hash?: string;
  tx_index?: string;
}
export interface QueryProveRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryProveRequest";
  value: QueryProveRequestAmino;
}
export interface QueryProveRequestSDKType {
  chain_id: bigint;
  tx_hash: string;
  proof?: ProofSDKType;
  block_hash: string;
  tx_index: bigint;
}
export interface QueryProveResponse {
  valid: boolean;
}
export interface QueryProveResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryProveResponse";
  value: Uint8Array;
}
export interface QueryProveResponseAmino {
  valid?: boolean;
}
export interface QueryProveResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryProveResponse";
  value: QueryProveResponseAmino;
}
export interface QueryProveResponseSDKType {
  valid: boolean;
}
export interface QueryParamsRequest {}
export interface QueryParamsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryParamsRequest";
  value: Uint8Array;
}
export interface QueryParamsRequestAmino {}
export interface QueryParamsRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryParamsRequest";
  value: QueryParamsRequestAmino;
}
export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params;
}
export interface QueryParamsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryParamsResponse";
  value: Uint8Array;
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponseAmino {
  /** params holds all the parameters of this module. */
  params?: ParamsAmino;
}
export interface QueryParamsResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryParamsResponse";
  value: QueryParamsResponseAmino;
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponseSDKType {
  params: ParamsSDKType;
}
export interface QueryHasVotedRequest {
  ballotIdentifier: string;
  voterAddress: string;
}
export interface QueryHasVotedRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryHasVotedRequest";
  value: Uint8Array;
}
export interface QueryHasVotedRequestAmino {
  ballot_identifier?: string;
  voter_address?: string;
}
export interface QueryHasVotedRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryHasVotedRequest";
  value: QueryHasVotedRequestAmino;
}
export interface QueryHasVotedRequestSDKType {
  ballot_identifier: string;
  voter_address: string;
}
export interface QueryHasVotedResponse {
  hasVoted: boolean;
}
export interface QueryHasVotedResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryHasVotedResponse";
  value: Uint8Array;
}
export interface QueryHasVotedResponseAmino {
  has_voted?: boolean;
}
export interface QueryHasVotedResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryHasVotedResponse";
  value: QueryHasVotedResponseAmino;
}
export interface QueryHasVotedResponseSDKType {
  has_voted: boolean;
}
export interface QueryBallotByIdentifierRequest {
  ballotIdentifier: string;
}
export interface QueryBallotByIdentifierRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryBallotByIdentifierRequest";
  value: Uint8Array;
}
export interface QueryBallotByIdentifierRequestAmino {
  ballot_identifier?: string;
}
export interface QueryBallotByIdentifierRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryBallotByIdentifierRequest";
  value: QueryBallotByIdentifierRequestAmino;
}
export interface QueryBallotByIdentifierRequestSDKType {
  ballot_identifier: string;
}
export interface VoterList {
  voterAddress: string;
  voteType: VoteType;
}
export interface VoterListProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.VoterList";
  value: Uint8Array;
}
export interface VoterListAmino {
  voter_address?: string;
  vote_type?: VoteType;
}
export interface VoterListAminoMsg {
  type: "/zetachain.zetacore.observer.VoterList";
  value: VoterListAmino;
}
export interface VoterListSDKType {
  voter_address: string;
  vote_type: VoteType;
}
export interface QueryBallotByIdentifierResponse {
  ballotIdentifier: string;
  voters: VoterList[];
  observationType: ObservationType;
  ballotStatus: BallotStatus;
}
export interface QueryBallotByIdentifierResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryBallotByIdentifierResponse";
  value: Uint8Array;
}
export interface QueryBallotByIdentifierResponseAmino {
  ballot_identifier?: string;
  voters?: VoterListAmino[];
  observation_type?: ObservationType;
  ballot_status?: BallotStatus;
}
export interface QueryBallotByIdentifierResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryBallotByIdentifierResponse";
  value: QueryBallotByIdentifierResponseAmino;
}
export interface QueryBallotByIdentifierResponseSDKType {
  ballot_identifier: string;
  voters: VoterListSDKType[];
  observation_type: ObservationType;
  ballot_status: BallotStatus;
}
export interface QueryObserverSet {}
export interface QueryObserverSetProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryObserverSet";
  value: Uint8Array;
}
export interface QueryObserverSetAmino {}
export interface QueryObserverSetAminoMsg {
  type: "/zetachain.zetacore.observer.QueryObserverSet";
  value: QueryObserverSetAmino;
}
export interface QueryObserverSetSDKType {}
export interface QueryObserverSetResponse {
  observers: string[];
}
export interface QueryObserverSetResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryObserverSetResponse";
  value: Uint8Array;
}
export interface QueryObserverSetResponseAmino {
  observers?: string[];
}
export interface QueryObserverSetResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryObserverSetResponse";
  value: QueryObserverSetResponseAmino;
}
export interface QueryObserverSetResponseSDKType {
  observers: string[];
}
export interface QuerySupportedChains {}
export interface QuerySupportedChainsProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QuerySupportedChains";
  value: Uint8Array;
}
export interface QuerySupportedChainsAmino {}
export interface QuerySupportedChainsAminoMsg {
  type: "/zetachain.zetacore.observer.QuerySupportedChains";
  value: QuerySupportedChainsAmino;
}
export interface QuerySupportedChainsSDKType {}
export interface QuerySupportedChainsResponse {
  chains: Chain[];
}
export interface QuerySupportedChainsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QuerySupportedChainsResponse";
  value: Uint8Array;
}
export interface QuerySupportedChainsResponseAmino {
  chains?: ChainAmino[];
}
export interface QuerySupportedChainsResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QuerySupportedChainsResponse";
  value: QuerySupportedChainsResponseAmino;
}
export interface QuerySupportedChainsResponseSDKType {
  chains: ChainSDKType[];
}
export interface QueryGetChainParamsForChainRequest {
  chainId: bigint;
}
export interface QueryGetChainParamsForChainRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsForChainRequest";
  value: Uint8Array;
}
export interface QueryGetChainParamsForChainRequestAmino {
  chain_id?: string;
}
export interface QueryGetChainParamsForChainRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetChainParamsForChainRequest";
  value: QueryGetChainParamsForChainRequestAmino;
}
export interface QueryGetChainParamsForChainRequestSDKType {
  chain_id: bigint;
}
export interface QueryGetChainParamsForChainResponse {
  chainParams?: ChainParams;
}
export interface QueryGetChainParamsForChainResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsForChainResponse";
  value: Uint8Array;
}
export interface QueryGetChainParamsForChainResponseAmino {
  chain_params?: ChainParamsAmino;
}
export interface QueryGetChainParamsForChainResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetChainParamsForChainResponse";
  value: QueryGetChainParamsForChainResponseAmino;
}
export interface QueryGetChainParamsForChainResponseSDKType {
  chain_params?: ChainParamsSDKType;
}
export interface QueryGetChainParamsRequest {}
export interface QueryGetChainParamsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsRequest";
  value: Uint8Array;
}
export interface QueryGetChainParamsRequestAmino {}
export interface QueryGetChainParamsRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetChainParamsRequest";
  value: QueryGetChainParamsRequestAmino;
}
export interface QueryGetChainParamsRequestSDKType {}
export interface QueryGetChainParamsResponse {
  chainParams?: ChainParamsList;
}
export interface QueryGetChainParamsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsResponse";
  value: Uint8Array;
}
export interface QueryGetChainParamsResponseAmino {
  chain_params?: ChainParamsListAmino;
}
export interface QueryGetChainParamsResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetChainParamsResponse";
  value: QueryGetChainParamsResponseAmino;
}
export interface QueryGetChainParamsResponseSDKType {
  chain_params?: ChainParamsListSDKType;
}
export interface QueryGetNodeAccountRequest {
  index: string;
}
export interface QueryGetNodeAccountRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetNodeAccountRequest";
  value: Uint8Array;
}
export interface QueryGetNodeAccountRequestAmino {
  index?: string;
}
export interface QueryGetNodeAccountRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetNodeAccountRequest";
  value: QueryGetNodeAccountRequestAmino;
}
export interface QueryGetNodeAccountRequestSDKType {
  index: string;
}
export interface QueryGetNodeAccountResponse {
  nodeAccount?: NodeAccount;
}
export interface QueryGetNodeAccountResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetNodeAccountResponse";
  value: Uint8Array;
}
export interface QueryGetNodeAccountResponseAmino {
  node_account?: NodeAccountAmino;
}
export interface QueryGetNodeAccountResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetNodeAccountResponse";
  value: QueryGetNodeAccountResponseAmino;
}
export interface QueryGetNodeAccountResponseSDKType {
  node_account?: NodeAccountSDKType;
}
export interface QueryAllNodeAccountRequest {
  pagination?: PageRequest;
}
export interface QueryAllNodeAccountRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllNodeAccountRequest";
  value: Uint8Array;
}
export interface QueryAllNodeAccountRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllNodeAccountRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllNodeAccountRequest";
  value: QueryAllNodeAccountRequestAmino;
}
export interface QueryAllNodeAccountRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllNodeAccountResponse {
  NodeAccount: NodeAccount[];
  pagination?: PageResponse;
}
export interface QueryAllNodeAccountResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllNodeAccountResponse";
  value: Uint8Array;
}
export interface QueryAllNodeAccountResponseAmino {
  NodeAccount?: NodeAccountAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllNodeAccountResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllNodeAccountResponse";
  value: QueryAllNodeAccountResponseAmino;
}
export interface QueryAllNodeAccountResponseSDKType {
  NodeAccount: NodeAccountSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryGetCrosschainFlagsRequest {}
export interface QueryGetCrosschainFlagsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetCrosschainFlagsRequest";
  value: Uint8Array;
}
export interface QueryGetCrosschainFlagsRequestAmino {}
export interface QueryGetCrosschainFlagsRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetCrosschainFlagsRequest";
  value: QueryGetCrosschainFlagsRequestAmino;
}
export interface QueryGetCrosschainFlagsRequestSDKType {}
export interface QueryGetCrosschainFlagsResponse {
  crosschainFlags: CrosschainFlags;
}
export interface QueryGetCrosschainFlagsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetCrosschainFlagsResponse";
  value: Uint8Array;
}
export interface QueryGetCrosschainFlagsResponseAmino {
  crosschain_flags?: CrosschainFlagsAmino;
}
export interface QueryGetCrosschainFlagsResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetCrosschainFlagsResponse";
  value: QueryGetCrosschainFlagsResponseAmino;
}
export interface QueryGetCrosschainFlagsResponseSDKType {
  crosschain_flags: CrosschainFlagsSDKType;
}
export interface QueryGetKeygenRequest {}
export interface QueryGetKeygenRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetKeygenRequest";
  value: Uint8Array;
}
export interface QueryGetKeygenRequestAmino {}
export interface QueryGetKeygenRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetKeygenRequest";
  value: QueryGetKeygenRequestAmino;
}
export interface QueryGetKeygenRequestSDKType {}
export interface QueryGetKeygenResponse {
  keygen?: Keygen;
}
export interface QueryGetKeygenResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetKeygenResponse";
  value: Uint8Array;
}
export interface QueryGetKeygenResponseAmino {
  keygen?: KeygenAmino;
}
export interface QueryGetKeygenResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetKeygenResponse";
  value: QueryGetKeygenResponseAmino;
}
export interface QueryGetKeygenResponseSDKType {
  keygen?: KeygenSDKType;
}
export interface QueryShowObserverCountRequest {}
export interface QueryShowObserverCountRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryShowObserverCountRequest";
  value: Uint8Array;
}
export interface QueryShowObserverCountRequestAmino {}
export interface QueryShowObserverCountRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryShowObserverCountRequest";
  value: QueryShowObserverCountRequestAmino;
}
export interface QueryShowObserverCountRequestSDKType {}
export interface QueryShowObserverCountResponse {
  lastObserverCount?: LastObserverCount;
}
export interface QueryShowObserverCountResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryShowObserverCountResponse";
  value: Uint8Array;
}
export interface QueryShowObserverCountResponseAmino {
  last_observer_count?: LastObserverCountAmino;
}
export interface QueryShowObserverCountResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryShowObserverCountResponse";
  value: QueryShowObserverCountResponseAmino;
}
export interface QueryShowObserverCountResponseSDKType {
  last_observer_count?: LastObserverCountSDKType;
}
export interface QueryBlameByIdentifierRequest {
  blameIdentifier: string;
}
export interface QueryBlameByIdentifierRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryBlameByIdentifierRequest";
  value: Uint8Array;
}
export interface QueryBlameByIdentifierRequestAmino {
  blame_identifier?: string;
}
export interface QueryBlameByIdentifierRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryBlameByIdentifierRequest";
  value: QueryBlameByIdentifierRequestAmino;
}
export interface QueryBlameByIdentifierRequestSDKType {
  blame_identifier: string;
}
export interface QueryBlameByIdentifierResponse {
  blameInfo?: Blame;
}
export interface QueryBlameByIdentifierResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryBlameByIdentifierResponse";
  value: Uint8Array;
}
export interface QueryBlameByIdentifierResponseAmino {
  blame_info?: BlameAmino;
}
export interface QueryBlameByIdentifierResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryBlameByIdentifierResponse";
  value: QueryBlameByIdentifierResponseAmino;
}
export interface QueryBlameByIdentifierResponseSDKType {
  blame_info?: BlameSDKType;
}
export interface QueryAllBlameRecordsRequest {
  pagination?: PageRequest;
}
export interface QueryAllBlameRecordsRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllBlameRecordsRequest";
  value: Uint8Array;
}
export interface QueryAllBlameRecordsRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllBlameRecordsRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllBlameRecordsRequest";
  value: QueryAllBlameRecordsRequestAmino;
}
export interface QueryAllBlameRecordsRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllBlameRecordsResponse {
  blameInfo: Blame[];
  pagination?: PageResponse;
}
export interface QueryAllBlameRecordsResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllBlameRecordsResponse";
  value: Uint8Array;
}
export interface QueryAllBlameRecordsResponseAmino {
  blame_info?: BlameAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllBlameRecordsResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllBlameRecordsResponse";
  value: QueryAllBlameRecordsResponseAmino;
}
export interface QueryAllBlameRecordsResponseSDKType {
  blame_info: BlameSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryBlameByChainAndNonceRequest {
  chainId: bigint;
  nonce: bigint;
}
export interface QueryBlameByChainAndNonceRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryBlameByChainAndNonceRequest";
  value: Uint8Array;
}
export interface QueryBlameByChainAndNonceRequestAmino {
  chain_id?: string;
  nonce?: string;
}
export interface QueryBlameByChainAndNonceRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryBlameByChainAndNonceRequest";
  value: QueryBlameByChainAndNonceRequestAmino;
}
export interface QueryBlameByChainAndNonceRequestSDKType {
  chain_id: bigint;
  nonce: bigint;
}
export interface QueryBlameByChainAndNonceResponse {
  blameInfo: Blame[];
}
export interface QueryBlameByChainAndNonceResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryBlameByChainAndNonceResponse";
  value: Uint8Array;
}
export interface QueryBlameByChainAndNonceResponseAmino {
  blame_info?: BlameAmino[];
}
export interface QueryBlameByChainAndNonceResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryBlameByChainAndNonceResponse";
  value: QueryBlameByChainAndNonceResponseAmino;
}
export interface QueryBlameByChainAndNonceResponseSDKType {
  blame_info: BlameSDKType[];
}
export interface QueryAllBlockHeaderRequest {
  pagination?: PageRequest;
}
export interface QueryAllBlockHeaderRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllBlockHeaderRequest";
  value: Uint8Array;
}
export interface QueryAllBlockHeaderRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllBlockHeaderRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllBlockHeaderRequest";
  value: QueryAllBlockHeaderRequestAmino;
}
export interface QueryAllBlockHeaderRequestSDKType {
  pagination?: PageRequestSDKType;
}
export interface QueryAllBlockHeaderResponse {
  blockHeaders: BlockHeader[];
  pagination?: PageResponse;
}
export interface QueryAllBlockHeaderResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryAllBlockHeaderResponse";
  value: Uint8Array;
}
export interface QueryAllBlockHeaderResponseAmino {
  block_headers?: BlockHeaderAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllBlockHeaderResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryAllBlockHeaderResponse";
  value: QueryAllBlockHeaderResponseAmino;
}
export interface QueryAllBlockHeaderResponseSDKType {
  block_headers: BlockHeaderSDKType[];
  pagination?: PageResponseSDKType;
}
export interface QueryGetBlockHeaderByHashRequest {
  blockHash: Uint8Array;
}
export interface QueryGetBlockHeaderByHashRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderByHashRequest";
  value: Uint8Array;
}
export interface QueryGetBlockHeaderByHashRequestAmino {
  block_hash?: string;
}
export interface QueryGetBlockHeaderByHashRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetBlockHeaderByHashRequest";
  value: QueryGetBlockHeaderByHashRequestAmino;
}
export interface QueryGetBlockHeaderByHashRequestSDKType {
  block_hash: Uint8Array;
}
export interface QueryGetBlockHeaderByHashResponse {
  blockHeader?: BlockHeader;
}
export interface QueryGetBlockHeaderByHashResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderByHashResponse";
  value: Uint8Array;
}
export interface QueryGetBlockHeaderByHashResponseAmino {
  block_header?: BlockHeaderAmino;
}
export interface QueryGetBlockHeaderByHashResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetBlockHeaderByHashResponse";
  value: QueryGetBlockHeaderByHashResponseAmino;
}
export interface QueryGetBlockHeaderByHashResponseSDKType {
  block_header?: BlockHeaderSDKType;
}
export interface QueryGetBlockHeaderStateRequest {
  chainId: bigint;
}
export interface QueryGetBlockHeaderStateRequestProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderStateRequest";
  value: Uint8Array;
}
export interface QueryGetBlockHeaderStateRequestAmino {
  chain_id?: string;
}
export interface QueryGetBlockHeaderStateRequestAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetBlockHeaderStateRequest";
  value: QueryGetBlockHeaderStateRequestAmino;
}
export interface QueryGetBlockHeaderStateRequestSDKType {
  chain_id: bigint;
}
export interface QueryGetBlockHeaderStateResponse {
  blockHeaderState?: BlockHeaderState;
}
export interface QueryGetBlockHeaderStateResponseProtoMsg {
  typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderStateResponse";
  value: Uint8Array;
}
export interface QueryGetBlockHeaderStateResponseAmino {
  block_header_state?: BlockHeaderStateAmino;
}
export interface QueryGetBlockHeaderStateResponseAminoMsg {
  type: "/zetachain.zetacore.observer.QueryGetBlockHeaderStateResponse";
  value: QueryGetBlockHeaderStateResponseAmino;
}
export interface QueryGetBlockHeaderStateResponseSDKType {
  block_header_state?: BlockHeaderStateSDKType;
}
function createBaseQueryGetChainNoncesRequest(): QueryGetChainNoncesRequest {
  return {
    index: ""
  };
}
export const QueryGetChainNoncesRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainNoncesRequest",
  encode(message: QueryGetChainNoncesRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetChainNoncesRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetChainNoncesRequest();
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
  fromPartial(object: Partial<QueryGetChainNoncesRequest>): QueryGetChainNoncesRequest {
    const message = createBaseQueryGetChainNoncesRequest();
    message.index = object.index ?? "";
    return message;
  },
  fromAmino(object: QueryGetChainNoncesRequestAmino): QueryGetChainNoncesRequest {
    const message = createBaseQueryGetChainNoncesRequest();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    return message;
  },
  toAmino(message: QueryGetChainNoncesRequest): QueryGetChainNoncesRequestAmino {
    const obj: any = {};
    obj.index = message.index;
    return obj;
  },
  fromAminoMsg(object: QueryGetChainNoncesRequestAminoMsg): QueryGetChainNoncesRequest {
    return QueryGetChainNoncesRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetChainNoncesRequestProtoMsg): QueryGetChainNoncesRequest {
    return QueryGetChainNoncesRequest.decode(message.value);
  },
  toProto(message: QueryGetChainNoncesRequest): Uint8Array {
    return QueryGetChainNoncesRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetChainNoncesRequest): QueryGetChainNoncesRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetChainNoncesRequest",
      value: QueryGetChainNoncesRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetChainNoncesResponse(): QueryGetChainNoncesResponse {
  return {
    ChainNonces: ChainNonces.fromPartial({})
  };
}
export const QueryGetChainNoncesResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainNoncesResponse",
  encode(message: QueryGetChainNoncesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.ChainNonces !== undefined) {
      ChainNonces.encode(message.ChainNonces, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetChainNoncesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetChainNoncesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ChainNonces = ChainNonces.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetChainNoncesResponse>): QueryGetChainNoncesResponse {
    const message = createBaseQueryGetChainNoncesResponse();
    message.ChainNonces = object.ChainNonces !== undefined && object.ChainNonces !== null ? ChainNonces.fromPartial(object.ChainNonces) : undefined;
    return message;
  },
  fromAmino(object: QueryGetChainNoncesResponseAmino): QueryGetChainNoncesResponse {
    const message = createBaseQueryGetChainNoncesResponse();
    if (object.ChainNonces !== undefined && object.ChainNonces !== null) {
      message.ChainNonces = ChainNonces.fromAmino(object.ChainNonces);
    }
    return message;
  },
  toAmino(message: QueryGetChainNoncesResponse): QueryGetChainNoncesResponseAmino {
    const obj: any = {};
    obj.ChainNonces = message.ChainNonces ? ChainNonces.toAmino(message.ChainNonces) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetChainNoncesResponseAminoMsg): QueryGetChainNoncesResponse {
    return QueryGetChainNoncesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetChainNoncesResponseProtoMsg): QueryGetChainNoncesResponse {
    return QueryGetChainNoncesResponse.decode(message.value);
  },
  toProto(message: QueryGetChainNoncesResponse): Uint8Array {
    return QueryGetChainNoncesResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetChainNoncesResponse): QueryGetChainNoncesResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetChainNoncesResponse",
      value: QueryGetChainNoncesResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllChainNoncesRequest(): QueryAllChainNoncesRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllChainNoncesRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllChainNoncesRequest",
  encode(message: QueryAllChainNoncesRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllChainNoncesRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllChainNoncesRequest();
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
  fromPartial(object: Partial<QueryAllChainNoncesRequest>): QueryAllChainNoncesRequest {
    const message = createBaseQueryAllChainNoncesRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllChainNoncesRequestAmino): QueryAllChainNoncesRequest {
    const message = createBaseQueryAllChainNoncesRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllChainNoncesRequest): QueryAllChainNoncesRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllChainNoncesRequestAminoMsg): QueryAllChainNoncesRequest {
    return QueryAllChainNoncesRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllChainNoncesRequestProtoMsg): QueryAllChainNoncesRequest {
    return QueryAllChainNoncesRequest.decode(message.value);
  },
  toProto(message: QueryAllChainNoncesRequest): Uint8Array {
    return QueryAllChainNoncesRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllChainNoncesRequest): QueryAllChainNoncesRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllChainNoncesRequest",
      value: QueryAllChainNoncesRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllChainNoncesResponse(): QueryAllChainNoncesResponse {
  return {
    ChainNonces: [],
    pagination: undefined
  };
}
export const QueryAllChainNoncesResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllChainNoncesResponse",
  encode(message: QueryAllChainNoncesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.ChainNonces) {
      ChainNonces.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllChainNoncesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllChainNoncesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ChainNonces.push(ChainNonces.decode(reader, reader.uint32()));
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
  fromPartial(object: Partial<QueryAllChainNoncesResponse>): QueryAllChainNoncesResponse {
    const message = createBaseQueryAllChainNoncesResponse();
    message.ChainNonces = object.ChainNonces?.map(e => ChainNonces.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllChainNoncesResponseAmino): QueryAllChainNoncesResponse {
    const message = createBaseQueryAllChainNoncesResponse();
    message.ChainNonces = object.ChainNonces?.map(e => ChainNonces.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllChainNoncesResponse): QueryAllChainNoncesResponseAmino {
    const obj: any = {};
    if (message.ChainNonces) {
      obj.ChainNonces = message.ChainNonces.map(e => e ? ChainNonces.toAmino(e) : undefined);
    } else {
      obj.ChainNonces = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllChainNoncesResponseAminoMsg): QueryAllChainNoncesResponse {
    return QueryAllChainNoncesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllChainNoncesResponseProtoMsg): QueryAllChainNoncesResponse {
    return QueryAllChainNoncesResponse.decode(message.value);
  },
  toProto(message: QueryAllChainNoncesResponse): Uint8Array {
    return QueryAllChainNoncesResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllChainNoncesResponse): QueryAllChainNoncesResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllChainNoncesResponse",
      value: QueryAllChainNoncesResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllPendingNoncesRequest(): QueryAllPendingNoncesRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllPendingNoncesRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllPendingNoncesRequest",
  encode(message: QueryAllPendingNoncesRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllPendingNoncesRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllPendingNoncesRequest();
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
  fromPartial(object: Partial<QueryAllPendingNoncesRequest>): QueryAllPendingNoncesRequest {
    const message = createBaseQueryAllPendingNoncesRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllPendingNoncesRequestAmino): QueryAllPendingNoncesRequest {
    const message = createBaseQueryAllPendingNoncesRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllPendingNoncesRequest): QueryAllPendingNoncesRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllPendingNoncesRequestAminoMsg): QueryAllPendingNoncesRequest {
    return QueryAllPendingNoncesRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllPendingNoncesRequestProtoMsg): QueryAllPendingNoncesRequest {
    return QueryAllPendingNoncesRequest.decode(message.value);
  },
  toProto(message: QueryAllPendingNoncesRequest): Uint8Array {
    return QueryAllPendingNoncesRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllPendingNoncesRequest): QueryAllPendingNoncesRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllPendingNoncesRequest",
      value: QueryAllPendingNoncesRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllPendingNoncesResponse(): QueryAllPendingNoncesResponse {
  return {
    pendingNonces: [],
    pagination: undefined
  };
}
export const QueryAllPendingNoncesResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllPendingNoncesResponse",
  encode(message: QueryAllPendingNoncesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.pendingNonces) {
      PendingNonces.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllPendingNoncesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllPendingNoncesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pendingNonces.push(PendingNonces.decode(reader, reader.uint32()));
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
  fromPartial(object: Partial<QueryAllPendingNoncesResponse>): QueryAllPendingNoncesResponse {
    const message = createBaseQueryAllPendingNoncesResponse();
    message.pendingNonces = object.pendingNonces?.map(e => PendingNonces.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllPendingNoncesResponseAmino): QueryAllPendingNoncesResponse {
    const message = createBaseQueryAllPendingNoncesResponse();
    message.pendingNonces = object.pending_nonces?.map(e => PendingNonces.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllPendingNoncesResponse): QueryAllPendingNoncesResponseAmino {
    const obj: any = {};
    if (message.pendingNonces) {
      obj.pending_nonces = message.pendingNonces.map(e => e ? PendingNonces.toAmino(e) : undefined);
    } else {
      obj.pending_nonces = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllPendingNoncesResponseAminoMsg): QueryAllPendingNoncesResponse {
    return QueryAllPendingNoncesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllPendingNoncesResponseProtoMsg): QueryAllPendingNoncesResponse {
    return QueryAllPendingNoncesResponse.decode(message.value);
  },
  toProto(message: QueryAllPendingNoncesResponse): Uint8Array {
    return QueryAllPendingNoncesResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllPendingNoncesResponse): QueryAllPendingNoncesResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllPendingNoncesResponse",
      value: QueryAllPendingNoncesResponse.encode(message).finish()
    };
  }
};
function createBaseQueryPendingNoncesByChainRequest(): QueryPendingNoncesByChainRequest {
  return {
    chainId: BigInt(0)
  };
}
export const QueryPendingNoncesByChainRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryPendingNoncesByChainRequest",
  encode(message: QueryPendingNoncesByChainRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryPendingNoncesByChainRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPendingNoncesByChainRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryPendingNoncesByChainRequest>): QueryPendingNoncesByChainRequest {
    const message = createBaseQueryPendingNoncesByChainRequest();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryPendingNoncesByChainRequestAmino): QueryPendingNoncesByChainRequest {
    const message = createBaseQueryPendingNoncesByChainRequest();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    return message;
  },
  toAmino(message: QueryPendingNoncesByChainRequest): QueryPendingNoncesByChainRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryPendingNoncesByChainRequestAminoMsg): QueryPendingNoncesByChainRequest {
    return QueryPendingNoncesByChainRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryPendingNoncesByChainRequestProtoMsg): QueryPendingNoncesByChainRequest {
    return QueryPendingNoncesByChainRequest.decode(message.value);
  },
  toProto(message: QueryPendingNoncesByChainRequest): Uint8Array {
    return QueryPendingNoncesByChainRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryPendingNoncesByChainRequest): QueryPendingNoncesByChainRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryPendingNoncesByChainRequest",
      value: QueryPendingNoncesByChainRequest.encode(message).finish()
    };
  }
};
function createBaseQueryPendingNoncesByChainResponse(): QueryPendingNoncesByChainResponse {
  return {
    pendingNonces: PendingNonces.fromPartial({})
  };
}
export const QueryPendingNoncesByChainResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryPendingNoncesByChainResponse",
  encode(message: QueryPendingNoncesByChainResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pendingNonces !== undefined) {
      PendingNonces.encode(message.pendingNonces, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryPendingNoncesByChainResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPendingNoncesByChainResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pendingNonces = PendingNonces.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryPendingNoncesByChainResponse>): QueryPendingNoncesByChainResponse {
    const message = createBaseQueryPendingNoncesByChainResponse();
    message.pendingNonces = object.pendingNonces !== undefined && object.pendingNonces !== null ? PendingNonces.fromPartial(object.pendingNonces) : undefined;
    return message;
  },
  fromAmino(object: QueryPendingNoncesByChainResponseAmino): QueryPendingNoncesByChainResponse {
    const message = createBaseQueryPendingNoncesByChainResponse();
    if (object.pending_nonces !== undefined && object.pending_nonces !== null) {
      message.pendingNonces = PendingNonces.fromAmino(object.pending_nonces);
    }
    return message;
  },
  toAmino(message: QueryPendingNoncesByChainResponse): QueryPendingNoncesByChainResponseAmino {
    const obj: any = {};
    obj.pending_nonces = message.pendingNonces ? PendingNonces.toAmino(message.pendingNonces) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryPendingNoncesByChainResponseAminoMsg): QueryPendingNoncesByChainResponse {
    return QueryPendingNoncesByChainResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryPendingNoncesByChainResponseProtoMsg): QueryPendingNoncesByChainResponse {
    return QueryPendingNoncesByChainResponse.decode(message.value);
  },
  toProto(message: QueryPendingNoncesByChainResponse): Uint8Array {
    return QueryPendingNoncesByChainResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryPendingNoncesByChainResponse): QueryPendingNoncesByChainResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryPendingNoncesByChainResponse",
      value: QueryPendingNoncesByChainResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetTSSRequest(): QueryGetTSSRequest {
  return {};
}
export const QueryGetTSSRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTSSRequest",
  encode(_: QueryGetTSSRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetTSSRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetTSSRequest();
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
  fromPartial(_: Partial<QueryGetTSSRequest>): QueryGetTSSRequest {
    const message = createBaseQueryGetTSSRequest();
    return message;
  },
  fromAmino(_: QueryGetTSSRequestAmino): QueryGetTSSRequest {
    const message = createBaseQueryGetTSSRequest();
    return message;
  },
  toAmino(_: QueryGetTSSRequest): QueryGetTSSRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryGetTSSRequestAminoMsg): QueryGetTSSRequest {
    return QueryGetTSSRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetTSSRequestProtoMsg): QueryGetTSSRequest {
    return QueryGetTSSRequest.decode(message.value);
  },
  toProto(message: QueryGetTSSRequest): Uint8Array {
    return QueryGetTSSRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetTSSRequest): QueryGetTSSRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetTSSRequest",
      value: QueryGetTSSRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetTSSResponse(): QueryGetTSSResponse {
  return {
    TSS: TSS.fromPartial({})
  };
}
export const QueryGetTSSResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTSSResponse",
  encode(message: QueryGetTSSResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.TSS !== undefined) {
      TSS.encode(message.TSS, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetTSSResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetTSSResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.TSS = TSS.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetTSSResponse>): QueryGetTSSResponse {
    const message = createBaseQueryGetTSSResponse();
    message.TSS = object.TSS !== undefined && object.TSS !== null ? TSS.fromPartial(object.TSS) : undefined;
    return message;
  },
  fromAmino(object: QueryGetTSSResponseAmino): QueryGetTSSResponse {
    const message = createBaseQueryGetTSSResponse();
    if (object.TSS !== undefined && object.TSS !== null) {
      message.TSS = TSS.fromAmino(object.TSS);
    }
    return message;
  },
  toAmino(message: QueryGetTSSResponse): QueryGetTSSResponseAmino {
    const obj: any = {};
    obj.TSS = message.TSS ? TSS.toAmino(message.TSS) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetTSSResponseAminoMsg): QueryGetTSSResponse {
    return QueryGetTSSResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetTSSResponseProtoMsg): QueryGetTSSResponse {
    return QueryGetTSSResponse.decode(message.value);
  },
  toProto(message: QueryGetTSSResponse): Uint8Array {
    return QueryGetTSSResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetTSSResponse): QueryGetTSSResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetTSSResponse",
      value: QueryGetTSSResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetTssAddressRequest(): QueryGetTssAddressRequest {
  return {
    bitcoinChainId: BigInt(0)
  };
}
export const QueryGetTssAddressRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressRequest",
  encode(message: QueryGetTssAddressRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.bitcoinChainId !== BigInt(0)) {
      writer.uint32(16).int64(message.bitcoinChainId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetTssAddressRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetTssAddressRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.bitcoinChainId = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetTssAddressRequest>): QueryGetTssAddressRequest {
    const message = createBaseQueryGetTssAddressRequest();
    message.bitcoinChainId = object.bitcoinChainId !== undefined && object.bitcoinChainId !== null ? BigInt(object.bitcoinChainId.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryGetTssAddressRequestAmino): QueryGetTssAddressRequest {
    const message = createBaseQueryGetTssAddressRequest();
    if (object.bitcoin_chain_id !== undefined && object.bitcoin_chain_id !== null) {
      message.bitcoinChainId = BigInt(object.bitcoin_chain_id);
    }
    return message;
  },
  toAmino(message: QueryGetTssAddressRequest): QueryGetTssAddressRequestAmino {
    const obj: any = {};
    obj.bitcoin_chain_id = message.bitcoinChainId ? message.bitcoinChainId.toString() : undefined;
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
      typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressRequest",
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
  typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressResponse",
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
      typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressResponse",
      value: QueryGetTssAddressResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetTssAddressByFinalizedHeightRequest(): QueryGetTssAddressByFinalizedHeightRequest {
  return {
    finalizedZetaHeight: BigInt(0),
    bitcoinChainId: BigInt(0)
  };
}
export const QueryGetTssAddressByFinalizedHeightRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressByFinalizedHeightRequest",
  encode(message: QueryGetTssAddressByFinalizedHeightRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.finalizedZetaHeight !== BigInt(0)) {
      writer.uint32(8).int64(message.finalizedZetaHeight);
    }
    if (message.bitcoinChainId !== BigInt(0)) {
      writer.uint32(16).int64(message.bitcoinChainId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetTssAddressByFinalizedHeightRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetTssAddressByFinalizedHeightRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.finalizedZetaHeight = reader.int64();
          break;
        case 2:
          message.bitcoinChainId = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetTssAddressByFinalizedHeightRequest>): QueryGetTssAddressByFinalizedHeightRequest {
    const message = createBaseQueryGetTssAddressByFinalizedHeightRequest();
    message.finalizedZetaHeight = object.finalizedZetaHeight !== undefined && object.finalizedZetaHeight !== null ? BigInt(object.finalizedZetaHeight.toString()) : BigInt(0);
    message.bitcoinChainId = object.bitcoinChainId !== undefined && object.bitcoinChainId !== null ? BigInt(object.bitcoinChainId.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryGetTssAddressByFinalizedHeightRequestAmino): QueryGetTssAddressByFinalizedHeightRequest {
    const message = createBaseQueryGetTssAddressByFinalizedHeightRequest();
    if (object.finalized_zeta_height !== undefined && object.finalized_zeta_height !== null) {
      message.finalizedZetaHeight = BigInt(object.finalized_zeta_height);
    }
    if (object.bitcoin_chain_id !== undefined && object.bitcoin_chain_id !== null) {
      message.bitcoinChainId = BigInt(object.bitcoin_chain_id);
    }
    return message;
  },
  toAmino(message: QueryGetTssAddressByFinalizedHeightRequest): QueryGetTssAddressByFinalizedHeightRequestAmino {
    const obj: any = {};
    obj.finalized_zeta_height = message.finalizedZetaHeight ? message.finalizedZetaHeight.toString() : undefined;
    obj.bitcoin_chain_id = message.bitcoinChainId ? message.bitcoinChainId.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetTssAddressByFinalizedHeightRequestAminoMsg): QueryGetTssAddressByFinalizedHeightRequest {
    return QueryGetTssAddressByFinalizedHeightRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetTssAddressByFinalizedHeightRequestProtoMsg): QueryGetTssAddressByFinalizedHeightRequest {
    return QueryGetTssAddressByFinalizedHeightRequest.decode(message.value);
  },
  toProto(message: QueryGetTssAddressByFinalizedHeightRequest): Uint8Array {
    return QueryGetTssAddressByFinalizedHeightRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetTssAddressByFinalizedHeightRequest): QueryGetTssAddressByFinalizedHeightRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressByFinalizedHeightRequest",
      value: QueryGetTssAddressByFinalizedHeightRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetTssAddressByFinalizedHeightResponse(): QueryGetTssAddressByFinalizedHeightResponse {
  return {
    eth: "",
    btc: ""
  };
}
export const QueryGetTssAddressByFinalizedHeightResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressByFinalizedHeightResponse",
  encode(message: QueryGetTssAddressByFinalizedHeightResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.eth !== "") {
      writer.uint32(10).string(message.eth);
    }
    if (message.btc !== "") {
      writer.uint32(18).string(message.btc);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetTssAddressByFinalizedHeightResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetTssAddressByFinalizedHeightResponse();
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
  fromPartial(object: Partial<QueryGetTssAddressByFinalizedHeightResponse>): QueryGetTssAddressByFinalizedHeightResponse {
    const message = createBaseQueryGetTssAddressByFinalizedHeightResponse();
    message.eth = object.eth ?? "";
    message.btc = object.btc ?? "";
    return message;
  },
  fromAmino(object: QueryGetTssAddressByFinalizedHeightResponseAmino): QueryGetTssAddressByFinalizedHeightResponse {
    const message = createBaseQueryGetTssAddressByFinalizedHeightResponse();
    if (object.eth !== undefined && object.eth !== null) {
      message.eth = object.eth;
    }
    if (object.btc !== undefined && object.btc !== null) {
      message.btc = object.btc;
    }
    return message;
  },
  toAmino(message: QueryGetTssAddressByFinalizedHeightResponse): QueryGetTssAddressByFinalizedHeightResponseAmino {
    const obj: any = {};
    obj.eth = message.eth;
    obj.btc = message.btc;
    return obj;
  },
  fromAminoMsg(object: QueryGetTssAddressByFinalizedHeightResponseAminoMsg): QueryGetTssAddressByFinalizedHeightResponse {
    return QueryGetTssAddressByFinalizedHeightResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetTssAddressByFinalizedHeightResponseProtoMsg): QueryGetTssAddressByFinalizedHeightResponse {
    return QueryGetTssAddressByFinalizedHeightResponse.decode(message.value);
  },
  toProto(message: QueryGetTssAddressByFinalizedHeightResponse): Uint8Array {
    return QueryGetTssAddressByFinalizedHeightResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetTssAddressByFinalizedHeightResponse): QueryGetTssAddressByFinalizedHeightResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetTssAddressByFinalizedHeightResponse",
      value: QueryGetTssAddressByFinalizedHeightResponse.encode(message).finish()
    };
  }
};
function createBaseQueryTssHistoryRequest(): QueryTssHistoryRequest {
  return {
    pagination: undefined
  };
}
export const QueryTssHistoryRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryTssHistoryRequest",
  encode(message: QueryTssHistoryRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryTssHistoryRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryTssHistoryRequest();
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
  fromPartial(object: Partial<QueryTssHistoryRequest>): QueryTssHistoryRequest {
    const message = createBaseQueryTssHistoryRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryTssHistoryRequestAmino): QueryTssHistoryRequest {
    const message = createBaseQueryTssHistoryRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryTssHistoryRequest): QueryTssHistoryRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryTssHistoryRequestAminoMsg): QueryTssHistoryRequest {
    return QueryTssHistoryRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryTssHistoryRequestProtoMsg): QueryTssHistoryRequest {
    return QueryTssHistoryRequest.decode(message.value);
  },
  toProto(message: QueryTssHistoryRequest): Uint8Array {
    return QueryTssHistoryRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryTssHistoryRequest): QueryTssHistoryRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryTssHistoryRequest",
      value: QueryTssHistoryRequest.encode(message).finish()
    };
  }
};
function createBaseQueryTssHistoryResponse(): QueryTssHistoryResponse {
  return {
    tssList: [],
    pagination: undefined
  };
}
export const QueryTssHistoryResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryTssHistoryResponse",
  encode(message: QueryTssHistoryResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.tssList) {
      TSS.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryTssHistoryResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryTssHistoryResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tssList.push(TSS.decode(reader, reader.uint32()));
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
  fromPartial(object: Partial<QueryTssHistoryResponse>): QueryTssHistoryResponse {
    const message = createBaseQueryTssHistoryResponse();
    message.tssList = object.tssList?.map(e => TSS.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryTssHistoryResponseAmino): QueryTssHistoryResponse {
    const message = createBaseQueryTssHistoryResponse();
    message.tssList = object.tss_list?.map(e => TSS.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryTssHistoryResponse): QueryTssHistoryResponseAmino {
    const obj: any = {};
    if (message.tssList) {
      obj.tss_list = message.tssList.map(e => e ? TSS.toAmino(e) : undefined);
    } else {
      obj.tss_list = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryTssHistoryResponseAminoMsg): QueryTssHistoryResponse {
    return QueryTssHistoryResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryTssHistoryResponseProtoMsg): QueryTssHistoryResponse {
    return QueryTssHistoryResponse.decode(message.value);
  },
  toProto(message: QueryTssHistoryResponse): Uint8Array {
    return QueryTssHistoryResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryTssHistoryResponse): QueryTssHistoryResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryTssHistoryResponse",
      value: QueryTssHistoryResponse.encode(message).finish()
    };
  }
};
function createBaseQueryProveRequest(): QueryProveRequest {
  return {
    chainId: BigInt(0),
    txHash: "",
    proof: undefined,
    blockHash: "",
    txIndex: BigInt(0)
  };
}
export const QueryProveRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryProveRequest",
  encode(message: QueryProveRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.txHash !== "") {
      writer.uint32(18).string(message.txHash);
    }
    if (message.proof !== undefined) {
      Proof.encode(message.proof, writer.uint32(26).fork()).ldelim();
    }
    if (message.blockHash !== "") {
      writer.uint32(34).string(message.blockHash);
    }
    if (message.txIndex !== BigInt(0)) {
      writer.uint32(40).int64(message.txIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryProveRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryProveRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.txHash = reader.string();
          break;
        case 3:
          message.proof = Proof.decode(reader, reader.uint32());
          break;
        case 4:
          message.blockHash = reader.string();
          break;
        case 5:
          message.txIndex = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryProveRequest>): QueryProveRequest {
    const message = createBaseQueryProveRequest();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.txHash = object.txHash ?? "";
    message.proof = object.proof !== undefined && object.proof !== null ? Proof.fromPartial(object.proof) : undefined;
    message.blockHash = object.blockHash ?? "";
    message.txIndex = object.txIndex !== undefined && object.txIndex !== null ? BigInt(object.txIndex.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryProveRequestAmino): QueryProveRequest {
    const message = createBaseQueryProveRequest();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
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
  toAmino(message: QueryProveRequest): QueryProveRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.tx_hash = message.txHash;
    obj.proof = message.proof ? Proof.toAmino(message.proof) : undefined;
    obj.block_hash = message.blockHash;
    obj.tx_index = message.txIndex ? message.txIndex.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryProveRequestAminoMsg): QueryProveRequest {
    return QueryProveRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryProveRequestProtoMsg): QueryProveRequest {
    return QueryProveRequest.decode(message.value);
  },
  toProto(message: QueryProveRequest): Uint8Array {
    return QueryProveRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryProveRequest): QueryProveRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryProveRequest",
      value: QueryProveRequest.encode(message).finish()
    };
  }
};
function createBaseQueryProveResponse(): QueryProveResponse {
  return {
    valid: false
  };
}
export const QueryProveResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryProveResponse",
  encode(message: QueryProveResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.valid === true) {
      writer.uint32(8).bool(message.valid);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryProveResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryProveResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.valid = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryProveResponse>): QueryProveResponse {
    const message = createBaseQueryProveResponse();
    message.valid = object.valid ?? false;
    return message;
  },
  fromAmino(object: QueryProveResponseAmino): QueryProveResponse {
    const message = createBaseQueryProveResponse();
    if (object.valid !== undefined && object.valid !== null) {
      message.valid = object.valid;
    }
    return message;
  },
  toAmino(message: QueryProveResponse): QueryProveResponseAmino {
    const obj: any = {};
    obj.valid = message.valid;
    return obj;
  },
  fromAminoMsg(object: QueryProveResponseAminoMsg): QueryProveResponse {
    return QueryProveResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryProveResponseProtoMsg): QueryProveResponse {
    return QueryProveResponse.decode(message.value);
  },
  toProto(message: QueryProveResponse): Uint8Array {
    return QueryProveResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryProveResponse): QueryProveResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryProveResponse",
      value: QueryProveResponse.encode(message).finish()
    };
  }
};
function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}
export const QueryParamsRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryParamsRequest",
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
      typeUrl: "/zetachain.zetacore.observer.QueryParamsRequest",
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
  typeUrl: "/zetachain.zetacore.observer.QueryParamsResponse",
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
      typeUrl: "/zetachain.zetacore.observer.QueryParamsResponse",
      value: QueryParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryHasVotedRequest(): QueryHasVotedRequest {
  return {
    ballotIdentifier: "",
    voterAddress: ""
  };
}
export const QueryHasVotedRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryHasVotedRequest",
  encode(message: QueryHasVotedRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.ballotIdentifier !== "") {
      writer.uint32(10).string(message.ballotIdentifier);
    }
    if (message.voterAddress !== "") {
      writer.uint32(18).string(message.voterAddress);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryHasVotedRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryHasVotedRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ballotIdentifier = reader.string();
          break;
        case 2:
          message.voterAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryHasVotedRequest>): QueryHasVotedRequest {
    const message = createBaseQueryHasVotedRequest();
    message.ballotIdentifier = object.ballotIdentifier ?? "";
    message.voterAddress = object.voterAddress ?? "";
    return message;
  },
  fromAmino(object: QueryHasVotedRequestAmino): QueryHasVotedRequest {
    const message = createBaseQueryHasVotedRequest();
    if (object.ballot_identifier !== undefined && object.ballot_identifier !== null) {
      message.ballotIdentifier = object.ballot_identifier;
    }
    if (object.voter_address !== undefined && object.voter_address !== null) {
      message.voterAddress = object.voter_address;
    }
    return message;
  },
  toAmino(message: QueryHasVotedRequest): QueryHasVotedRequestAmino {
    const obj: any = {};
    obj.ballot_identifier = message.ballotIdentifier;
    obj.voter_address = message.voterAddress;
    return obj;
  },
  fromAminoMsg(object: QueryHasVotedRequestAminoMsg): QueryHasVotedRequest {
    return QueryHasVotedRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryHasVotedRequestProtoMsg): QueryHasVotedRequest {
    return QueryHasVotedRequest.decode(message.value);
  },
  toProto(message: QueryHasVotedRequest): Uint8Array {
    return QueryHasVotedRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryHasVotedRequest): QueryHasVotedRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryHasVotedRequest",
      value: QueryHasVotedRequest.encode(message).finish()
    };
  }
};
function createBaseQueryHasVotedResponse(): QueryHasVotedResponse {
  return {
    hasVoted: false
  };
}
export const QueryHasVotedResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryHasVotedResponse",
  encode(message: QueryHasVotedResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.hasVoted === true) {
      writer.uint32(8).bool(message.hasVoted);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryHasVotedResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryHasVotedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.hasVoted = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryHasVotedResponse>): QueryHasVotedResponse {
    const message = createBaseQueryHasVotedResponse();
    message.hasVoted = object.hasVoted ?? false;
    return message;
  },
  fromAmino(object: QueryHasVotedResponseAmino): QueryHasVotedResponse {
    const message = createBaseQueryHasVotedResponse();
    if (object.has_voted !== undefined && object.has_voted !== null) {
      message.hasVoted = object.has_voted;
    }
    return message;
  },
  toAmino(message: QueryHasVotedResponse): QueryHasVotedResponseAmino {
    const obj: any = {};
    obj.has_voted = message.hasVoted;
    return obj;
  },
  fromAminoMsg(object: QueryHasVotedResponseAminoMsg): QueryHasVotedResponse {
    return QueryHasVotedResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryHasVotedResponseProtoMsg): QueryHasVotedResponse {
    return QueryHasVotedResponse.decode(message.value);
  },
  toProto(message: QueryHasVotedResponse): Uint8Array {
    return QueryHasVotedResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryHasVotedResponse): QueryHasVotedResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryHasVotedResponse",
      value: QueryHasVotedResponse.encode(message).finish()
    };
  }
};
function createBaseQueryBallotByIdentifierRequest(): QueryBallotByIdentifierRequest {
  return {
    ballotIdentifier: ""
  };
}
export const QueryBallotByIdentifierRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryBallotByIdentifierRequest",
  encode(message: QueryBallotByIdentifierRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.ballotIdentifier !== "") {
      writer.uint32(10).string(message.ballotIdentifier);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryBallotByIdentifierRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBallotByIdentifierRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ballotIdentifier = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryBallotByIdentifierRequest>): QueryBallotByIdentifierRequest {
    const message = createBaseQueryBallotByIdentifierRequest();
    message.ballotIdentifier = object.ballotIdentifier ?? "";
    return message;
  },
  fromAmino(object: QueryBallotByIdentifierRequestAmino): QueryBallotByIdentifierRequest {
    const message = createBaseQueryBallotByIdentifierRequest();
    if (object.ballot_identifier !== undefined && object.ballot_identifier !== null) {
      message.ballotIdentifier = object.ballot_identifier;
    }
    return message;
  },
  toAmino(message: QueryBallotByIdentifierRequest): QueryBallotByIdentifierRequestAmino {
    const obj: any = {};
    obj.ballot_identifier = message.ballotIdentifier;
    return obj;
  },
  fromAminoMsg(object: QueryBallotByIdentifierRequestAminoMsg): QueryBallotByIdentifierRequest {
    return QueryBallotByIdentifierRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryBallotByIdentifierRequestProtoMsg): QueryBallotByIdentifierRequest {
    return QueryBallotByIdentifierRequest.decode(message.value);
  },
  toProto(message: QueryBallotByIdentifierRequest): Uint8Array {
    return QueryBallotByIdentifierRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryBallotByIdentifierRequest): QueryBallotByIdentifierRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryBallotByIdentifierRequest",
      value: QueryBallotByIdentifierRequest.encode(message).finish()
    };
  }
};
function createBaseVoterList(): VoterList {
  return {
    voterAddress: "",
    voteType: 0
  };
}
export const VoterList = {
  typeUrl: "/zetachain.zetacore.observer.VoterList",
  encode(message: VoterList, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.voterAddress !== "") {
      writer.uint32(10).string(message.voterAddress);
    }
    if (message.voteType !== 0) {
      writer.uint32(16).int32(message.voteType);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): VoterList {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVoterList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.voterAddress = reader.string();
          break;
        case 2:
          message.voteType = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<VoterList>): VoterList {
    const message = createBaseVoterList();
    message.voterAddress = object.voterAddress ?? "";
    message.voteType = object.voteType ?? 0;
    return message;
  },
  fromAmino(object: VoterListAmino): VoterList {
    const message = createBaseVoterList();
    if (object.voter_address !== undefined && object.voter_address !== null) {
      message.voterAddress = object.voter_address;
    }
    if (object.vote_type !== undefined && object.vote_type !== null) {
      message.voteType = voteTypeFromJSON(object.vote_type);
    }
    return message;
  },
  toAmino(message: VoterList): VoterListAmino {
    const obj: any = {};
    obj.voter_address = message.voterAddress;
    obj.vote_type = message.voteType;
    return obj;
  },
  fromAminoMsg(object: VoterListAminoMsg): VoterList {
    return VoterList.fromAmino(object.value);
  },
  fromProtoMsg(message: VoterListProtoMsg): VoterList {
    return VoterList.decode(message.value);
  },
  toProto(message: VoterList): Uint8Array {
    return VoterList.encode(message).finish();
  },
  toProtoMsg(message: VoterList): VoterListProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.VoterList",
      value: VoterList.encode(message).finish()
    };
  }
};
function createBaseQueryBallotByIdentifierResponse(): QueryBallotByIdentifierResponse {
  return {
    ballotIdentifier: "",
    voters: [],
    observationType: 0,
    ballotStatus: 0
  };
}
export const QueryBallotByIdentifierResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryBallotByIdentifierResponse",
  encode(message: QueryBallotByIdentifierResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.ballotIdentifier !== "") {
      writer.uint32(10).string(message.ballotIdentifier);
    }
    for (const v of message.voters) {
      VoterList.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.observationType !== 0) {
      writer.uint32(24).int32(message.observationType);
    }
    if (message.ballotStatus !== 0) {
      writer.uint32(32).int32(message.ballotStatus);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryBallotByIdentifierResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBallotByIdentifierResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.ballotIdentifier = reader.string();
          break;
        case 2:
          message.voters.push(VoterList.decode(reader, reader.uint32()));
          break;
        case 3:
          message.observationType = (reader.int32() as any);
          break;
        case 4:
          message.ballotStatus = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryBallotByIdentifierResponse>): QueryBallotByIdentifierResponse {
    const message = createBaseQueryBallotByIdentifierResponse();
    message.ballotIdentifier = object.ballotIdentifier ?? "";
    message.voters = object.voters?.map(e => VoterList.fromPartial(e)) || [];
    message.observationType = object.observationType ?? 0;
    message.ballotStatus = object.ballotStatus ?? 0;
    return message;
  },
  fromAmino(object: QueryBallotByIdentifierResponseAmino): QueryBallotByIdentifierResponse {
    const message = createBaseQueryBallotByIdentifierResponse();
    if (object.ballot_identifier !== undefined && object.ballot_identifier !== null) {
      message.ballotIdentifier = object.ballot_identifier;
    }
    message.voters = object.voters?.map(e => VoterList.fromAmino(e)) || [];
    if (object.observation_type !== undefined && object.observation_type !== null) {
      message.observationType = observationTypeFromJSON(object.observation_type);
    }
    if (object.ballot_status !== undefined && object.ballot_status !== null) {
      message.ballotStatus = ballotStatusFromJSON(object.ballot_status);
    }
    return message;
  },
  toAmino(message: QueryBallotByIdentifierResponse): QueryBallotByIdentifierResponseAmino {
    const obj: any = {};
    obj.ballot_identifier = message.ballotIdentifier;
    if (message.voters) {
      obj.voters = message.voters.map(e => e ? VoterList.toAmino(e) : undefined);
    } else {
      obj.voters = [];
    }
    obj.observation_type = message.observationType;
    obj.ballot_status = message.ballotStatus;
    return obj;
  },
  fromAminoMsg(object: QueryBallotByIdentifierResponseAminoMsg): QueryBallotByIdentifierResponse {
    return QueryBallotByIdentifierResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryBallotByIdentifierResponseProtoMsg): QueryBallotByIdentifierResponse {
    return QueryBallotByIdentifierResponse.decode(message.value);
  },
  toProto(message: QueryBallotByIdentifierResponse): Uint8Array {
    return QueryBallotByIdentifierResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryBallotByIdentifierResponse): QueryBallotByIdentifierResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryBallotByIdentifierResponse",
      value: QueryBallotByIdentifierResponse.encode(message).finish()
    };
  }
};
function createBaseQueryObserverSet(): QueryObserverSet {
  return {};
}
export const QueryObserverSet = {
  typeUrl: "/zetachain.zetacore.observer.QueryObserverSet",
  encode(_: QueryObserverSet, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryObserverSet {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryObserverSet();
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
  fromPartial(_: Partial<QueryObserverSet>): QueryObserverSet {
    const message = createBaseQueryObserverSet();
    return message;
  },
  fromAmino(_: QueryObserverSetAmino): QueryObserverSet {
    const message = createBaseQueryObserverSet();
    return message;
  },
  toAmino(_: QueryObserverSet): QueryObserverSetAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryObserverSetAminoMsg): QueryObserverSet {
    return QueryObserverSet.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryObserverSetProtoMsg): QueryObserverSet {
    return QueryObserverSet.decode(message.value);
  },
  toProto(message: QueryObserverSet): Uint8Array {
    return QueryObserverSet.encode(message).finish();
  },
  toProtoMsg(message: QueryObserverSet): QueryObserverSetProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryObserverSet",
      value: QueryObserverSet.encode(message).finish()
    };
  }
};
function createBaseQueryObserverSetResponse(): QueryObserverSetResponse {
  return {
    observers: []
  };
}
export const QueryObserverSetResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryObserverSetResponse",
  encode(message: QueryObserverSetResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.observers) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryObserverSetResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryObserverSetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.observers.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryObserverSetResponse>): QueryObserverSetResponse {
    const message = createBaseQueryObserverSetResponse();
    message.observers = object.observers?.map(e => e) || [];
    return message;
  },
  fromAmino(object: QueryObserverSetResponseAmino): QueryObserverSetResponse {
    const message = createBaseQueryObserverSetResponse();
    message.observers = object.observers?.map(e => e) || [];
    return message;
  },
  toAmino(message: QueryObserverSetResponse): QueryObserverSetResponseAmino {
    const obj: any = {};
    if (message.observers) {
      obj.observers = message.observers.map(e => e);
    } else {
      obj.observers = [];
    }
    return obj;
  },
  fromAminoMsg(object: QueryObserverSetResponseAminoMsg): QueryObserverSetResponse {
    return QueryObserverSetResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryObserverSetResponseProtoMsg): QueryObserverSetResponse {
    return QueryObserverSetResponse.decode(message.value);
  },
  toProto(message: QueryObserverSetResponse): Uint8Array {
    return QueryObserverSetResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryObserverSetResponse): QueryObserverSetResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryObserverSetResponse",
      value: QueryObserverSetResponse.encode(message).finish()
    };
  }
};
function createBaseQuerySupportedChains(): QuerySupportedChains {
  return {};
}
export const QuerySupportedChains = {
  typeUrl: "/zetachain.zetacore.observer.QuerySupportedChains",
  encode(_: QuerySupportedChains, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QuerySupportedChains {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySupportedChains();
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
  fromPartial(_: Partial<QuerySupportedChains>): QuerySupportedChains {
    const message = createBaseQuerySupportedChains();
    return message;
  },
  fromAmino(_: QuerySupportedChainsAmino): QuerySupportedChains {
    const message = createBaseQuerySupportedChains();
    return message;
  },
  toAmino(_: QuerySupportedChains): QuerySupportedChainsAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QuerySupportedChainsAminoMsg): QuerySupportedChains {
    return QuerySupportedChains.fromAmino(object.value);
  },
  fromProtoMsg(message: QuerySupportedChainsProtoMsg): QuerySupportedChains {
    return QuerySupportedChains.decode(message.value);
  },
  toProto(message: QuerySupportedChains): Uint8Array {
    return QuerySupportedChains.encode(message).finish();
  },
  toProtoMsg(message: QuerySupportedChains): QuerySupportedChainsProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QuerySupportedChains",
      value: QuerySupportedChains.encode(message).finish()
    };
  }
};
function createBaseQuerySupportedChainsResponse(): QuerySupportedChainsResponse {
  return {
    chains: []
  };
}
export const QuerySupportedChainsResponse = {
  typeUrl: "/zetachain.zetacore.observer.QuerySupportedChainsResponse",
  encode(message: QuerySupportedChainsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.chains) {
      Chain.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QuerySupportedChainsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySupportedChainsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chains.push(Chain.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QuerySupportedChainsResponse>): QuerySupportedChainsResponse {
    const message = createBaseQuerySupportedChainsResponse();
    message.chains = object.chains?.map(e => Chain.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: QuerySupportedChainsResponseAmino): QuerySupportedChainsResponse {
    const message = createBaseQuerySupportedChainsResponse();
    message.chains = object.chains?.map(e => Chain.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: QuerySupportedChainsResponse): QuerySupportedChainsResponseAmino {
    const obj: any = {};
    if (message.chains) {
      obj.chains = message.chains.map(e => e ? Chain.toAmino(e) : undefined);
    } else {
      obj.chains = [];
    }
    return obj;
  },
  fromAminoMsg(object: QuerySupportedChainsResponseAminoMsg): QuerySupportedChainsResponse {
    return QuerySupportedChainsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QuerySupportedChainsResponseProtoMsg): QuerySupportedChainsResponse {
    return QuerySupportedChainsResponse.decode(message.value);
  },
  toProto(message: QuerySupportedChainsResponse): Uint8Array {
    return QuerySupportedChainsResponse.encode(message).finish();
  },
  toProtoMsg(message: QuerySupportedChainsResponse): QuerySupportedChainsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QuerySupportedChainsResponse",
      value: QuerySupportedChainsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetChainParamsForChainRequest(): QueryGetChainParamsForChainRequest {
  return {
    chainId: BigInt(0)
  };
}
export const QueryGetChainParamsForChainRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsForChainRequest",
  encode(message: QueryGetChainParamsForChainRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetChainParamsForChainRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetChainParamsForChainRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetChainParamsForChainRequest>): QueryGetChainParamsForChainRequest {
    const message = createBaseQueryGetChainParamsForChainRequest();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryGetChainParamsForChainRequestAmino): QueryGetChainParamsForChainRequest {
    const message = createBaseQueryGetChainParamsForChainRequest();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    return message;
  },
  toAmino(message: QueryGetChainParamsForChainRequest): QueryGetChainParamsForChainRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetChainParamsForChainRequestAminoMsg): QueryGetChainParamsForChainRequest {
    return QueryGetChainParamsForChainRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetChainParamsForChainRequestProtoMsg): QueryGetChainParamsForChainRequest {
    return QueryGetChainParamsForChainRequest.decode(message.value);
  },
  toProto(message: QueryGetChainParamsForChainRequest): Uint8Array {
    return QueryGetChainParamsForChainRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetChainParamsForChainRequest): QueryGetChainParamsForChainRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsForChainRequest",
      value: QueryGetChainParamsForChainRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetChainParamsForChainResponse(): QueryGetChainParamsForChainResponse {
  return {
    chainParams: undefined
  };
}
export const QueryGetChainParamsForChainResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsForChainResponse",
  encode(message: QueryGetChainParamsForChainResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainParams !== undefined) {
      ChainParams.encode(message.chainParams, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetChainParamsForChainResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetChainParamsForChainResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainParams = ChainParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetChainParamsForChainResponse>): QueryGetChainParamsForChainResponse {
    const message = createBaseQueryGetChainParamsForChainResponse();
    message.chainParams = object.chainParams !== undefined && object.chainParams !== null ? ChainParams.fromPartial(object.chainParams) : undefined;
    return message;
  },
  fromAmino(object: QueryGetChainParamsForChainResponseAmino): QueryGetChainParamsForChainResponse {
    const message = createBaseQueryGetChainParamsForChainResponse();
    if (object.chain_params !== undefined && object.chain_params !== null) {
      message.chainParams = ChainParams.fromAmino(object.chain_params);
    }
    return message;
  },
  toAmino(message: QueryGetChainParamsForChainResponse): QueryGetChainParamsForChainResponseAmino {
    const obj: any = {};
    obj.chain_params = message.chainParams ? ChainParams.toAmino(message.chainParams) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetChainParamsForChainResponseAminoMsg): QueryGetChainParamsForChainResponse {
    return QueryGetChainParamsForChainResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetChainParamsForChainResponseProtoMsg): QueryGetChainParamsForChainResponse {
    return QueryGetChainParamsForChainResponse.decode(message.value);
  },
  toProto(message: QueryGetChainParamsForChainResponse): Uint8Array {
    return QueryGetChainParamsForChainResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetChainParamsForChainResponse): QueryGetChainParamsForChainResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsForChainResponse",
      value: QueryGetChainParamsForChainResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetChainParamsRequest(): QueryGetChainParamsRequest {
  return {};
}
export const QueryGetChainParamsRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsRequest",
  encode(_: QueryGetChainParamsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetChainParamsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetChainParamsRequest();
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
  fromPartial(_: Partial<QueryGetChainParamsRequest>): QueryGetChainParamsRequest {
    const message = createBaseQueryGetChainParamsRequest();
    return message;
  },
  fromAmino(_: QueryGetChainParamsRequestAmino): QueryGetChainParamsRequest {
    const message = createBaseQueryGetChainParamsRequest();
    return message;
  },
  toAmino(_: QueryGetChainParamsRequest): QueryGetChainParamsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryGetChainParamsRequestAminoMsg): QueryGetChainParamsRequest {
    return QueryGetChainParamsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetChainParamsRequestProtoMsg): QueryGetChainParamsRequest {
    return QueryGetChainParamsRequest.decode(message.value);
  },
  toProto(message: QueryGetChainParamsRequest): Uint8Array {
    return QueryGetChainParamsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetChainParamsRequest): QueryGetChainParamsRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsRequest",
      value: QueryGetChainParamsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetChainParamsResponse(): QueryGetChainParamsResponse {
  return {
    chainParams: undefined
  };
}
export const QueryGetChainParamsResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsResponse",
  encode(message: QueryGetChainParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainParams !== undefined) {
      ChainParamsList.encode(message.chainParams, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetChainParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetChainParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainParams = ChainParamsList.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetChainParamsResponse>): QueryGetChainParamsResponse {
    const message = createBaseQueryGetChainParamsResponse();
    message.chainParams = object.chainParams !== undefined && object.chainParams !== null ? ChainParamsList.fromPartial(object.chainParams) : undefined;
    return message;
  },
  fromAmino(object: QueryGetChainParamsResponseAmino): QueryGetChainParamsResponse {
    const message = createBaseQueryGetChainParamsResponse();
    if (object.chain_params !== undefined && object.chain_params !== null) {
      message.chainParams = ChainParamsList.fromAmino(object.chain_params);
    }
    return message;
  },
  toAmino(message: QueryGetChainParamsResponse): QueryGetChainParamsResponseAmino {
    const obj: any = {};
    obj.chain_params = message.chainParams ? ChainParamsList.toAmino(message.chainParams) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetChainParamsResponseAminoMsg): QueryGetChainParamsResponse {
    return QueryGetChainParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetChainParamsResponseProtoMsg): QueryGetChainParamsResponse {
    return QueryGetChainParamsResponse.decode(message.value);
  },
  toProto(message: QueryGetChainParamsResponse): Uint8Array {
    return QueryGetChainParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetChainParamsResponse): QueryGetChainParamsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetChainParamsResponse",
      value: QueryGetChainParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetNodeAccountRequest(): QueryGetNodeAccountRequest {
  return {
    index: ""
  };
}
export const QueryGetNodeAccountRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetNodeAccountRequest",
  encode(message: QueryGetNodeAccountRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetNodeAccountRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetNodeAccountRequest();
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
  fromPartial(object: Partial<QueryGetNodeAccountRequest>): QueryGetNodeAccountRequest {
    const message = createBaseQueryGetNodeAccountRequest();
    message.index = object.index ?? "";
    return message;
  },
  fromAmino(object: QueryGetNodeAccountRequestAmino): QueryGetNodeAccountRequest {
    const message = createBaseQueryGetNodeAccountRequest();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    return message;
  },
  toAmino(message: QueryGetNodeAccountRequest): QueryGetNodeAccountRequestAmino {
    const obj: any = {};
    obj.index = message.index;
    return obj;
  },
  fromAminoMsg(object: QueryGetNodeAccountRequestAminoMsg): QueryGetNodeAccountRequest {
    return QueryGetNodeAccountRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetNodeAccountRequestProtoMsg): QueryGetNodeAccountRequest {
    return QueryGetNodeAccountRequest.decode(message.value);
  },
  toProto(message: QueryGetNodeAccountRequest): Uint8Array {
    return QueryGetNodeAccountRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetNodeAccountRequest): QueryGetNodeAccountRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetNodeAccountRequest",
      value: QueryGetNodeAccountRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetNodeAccountResponse(): QueryGetNodeAccountResponse {
  return {
    nodeAccount: undefined
  };
}
export const QueryGetNodeAccountResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetNodeAccountResponse",
  encode(message: QueryGetNodeAccountResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.nodeAccount !== undefined) {
      NodeAccount.encode(message.nodeAccount, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetNodeAccountResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetNodeAccountResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.nodeAccount = NodeAccount.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetNodeAccountResponse>): QueryGetNodeAccountResponse {
    const message = createBaseQueryGetNodeAccountResponse();
    message.nodeAccount = object.nodeAccount !== undefined && object.nodeAccount !== null ? NodeAccount.fromPartial(object.nodeAccount) : undefined;
    return message;
  },
  fromAmino(object: QueryGetNodeAccountResponseAmino): QueryGetNodeAccountResponse {
    const message = createBaseQueryGetNodeAccountResponse();
    if (object.node_account !== undefined && object.node_account !== null) {
      message.nodeAccount = NodeAccount.fromAmino(object.node_account);
    }
    return message;
  },
  toAmino(message: QueryGetNodeAccountResponse): QueryGetNodeAccountResponseAmino {
    const obj: any = {};
    obj.node_account = message.nodeAccount ? NodeAccount.toAmino(message.nodeAccount) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetNodeAccountResponseAminoMsg): QueryGetNodeAccountResponse {
    return QueryGetNodeAccountResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetNodeAccountResponseProtoMsg): QueryGetNodeAccountResponse {
    return QueryGetNodeAccountResponse.decode(message.value);
  },
  toProto(message: QueryGetNodeAccountResponse): Uint8Array {
    return QueryGetNodeAccountResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetNodeAccountResponse): QueryGetNodeAccountResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetNodeAccountResponse",
      value: QueryGetNodeAccountResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllNodeAccountRequest(): QueryAllNodeAccountRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllNodeAccountRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllNodeAccountRequest",
  encode(message: QueryAllNodeAccountRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllNodeAccountRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllNodeAccountRequest();
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
  fromPartial(object: Partial<QueryAllNodeAccountRequest>): QueryAllNodeAccountRequest {
    const message = createBaseQueryAllNodeAccountRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllNodeAccountRequestAmino): QueryAllNodeAccountRequest {
    const message = createBaseQueryAllNodeAccountRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllNodeAccountRequest): QueryAllNodeAccountRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllNodeAccountRequestAminoMsg): QueryAllNodeAccountRequest {
    return QueryAllNodeAccountRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllNodeAccountRequestProtoMsg): QueryAllNodeAccountRequest {
    return QueryAllNodeAccountRequest.decode(message.value);
  },
  toProto(message: QueryAllNodeAccountRequest): Uint8Array {
    return QueryAllNodeAccountRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllNodeAccountRequest): QueryAllNodeAccountRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllNodeAccountRequest",
      value: QueryAllNodeAccountRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllNodeAccountResponse(): QueryAllNodeAccountResponse {
  return {
    NodeAccount: [],
    pagination: undefined
  };
}
export const QueryAllNodeAccountResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllNodeAccountResponse",
  encode(message: QueryAllNodeAccountResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.NodeAccount) {
      NodeAccount.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllNodeAccountResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllNodeAccountResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.NodeAccount.push(NodeAccount.decode(reader, reader.uint32()));
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
  fromPartial(object: Partial<QueryAllNodeAccountResponse>): QueryAllNodeAccountResponse {
    const message = createBaseQueryAllNodeAccountResponse();
    message.NodeAccount = object.NodeAccount?.map(e => NodeAccount.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllNodeAccountResponseAmino): QueryAllNodeAccountResponse {
    const message = createBaseQueryAllNodeAccountResponse();
    message.NodeAccount = object.NodeAccount?.map(e => NodeAccount.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllNodeAccountResponse): QueryAllNodeAccountResponseAmino {
    const obj: any = {};
    if (message.NodeAccount) {
      obj.NodeAccount = message.NodeAccount.map(e => e ? NodeAccount.toAmino(e) : undefined);
    } else {
      obj.NodeAccount = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllNodeAccountResponseAminoMsg): QueryAllNodeAccountResponse {
    return QueryAllNodeAccountResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllNodeAccountResponseProtoMsg): QueryAllNodeAccountResponse {
    return QueryAllNodeAccountResponse.decode(message.value);
  },
  toProto(message: QueryAllNodeAccountResponse): Uint8Array {
    return QueryAllNodeAccountResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllNodeAccountResponse): QueryAllNodeAccountResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllNodeAccountResponse",
      value: QueryAllNodeAccountResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetCrosschainFlagsRequest(): QueryGetCrosschainFlagsRequest {
  return {};
}
export const QueryGetCrosschainFlagsRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetCrosschainFlagsRequest",
  encode(_: QueryGetCrosschainFlagsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetCrosschainFlagsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetCrosschainFlagsRequest();
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
  fromPartial(_: Partial<QueryGetCrosschainFlagsRequest>): QueryGetCrosschainFlagsRequest {
    const message = createBaseQueryGetCrosschainFlagsRequest();
    return message;
  },
  fromAmino(_: QueryGetCrosschainFlagsRequestAmino): QueryGetCrosschainFlagsRequest {
    const message = createBaseQueryGetCrosschainFlagsRequest();
    return message;
  },
  toAmino(_: QueryGetCrosschainFlagsRequest): QueryGetCrosschainFlagsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryGetCrosschainFlagsRequestAminoMsg): QueryGetCrosschainFlagsRequest {
    return QueryGetCrosschainFlagsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetCrosschainFlagsRequestProtoMsg): QueryGetCrosschainFlagsRequest {
    return QueryGetCrosschainFlagsRequest.decode(message.value);
  },
  toProto(message: QueryGetCrosschainFlagsRequest): Uint8Array {
    return QueryGetCrosschainFlagsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetCrosschainFlagsRequest): QueryGetCrosschainFlagsRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetCrosschainFlagsRequest",
      value: QueryGetCrosschainFlagsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetCrosschainFlagsResponse(): QueryGetCrosschainFlagsResponse {
  return {
    crosschainFlags: CrosschainFlags.fromPartial({})
  };
}
export const QueryGetCrosschainFlagsResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetCrosschainFlagsResponse",
  encode(message: QueryGetCrosschainFlagsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.crosschainFlags !== undefined) {
      CrosschainFlags.encode(message.crosschainFlags, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetCrosschainFlagsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetCrosschainFlagsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.crosschainFlags = CrosschainFlags.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetCrosschainFlagsResponse>): QueryGetCrosschainFlagsResponse {
    const message = createBaseQueryGetCrosschainFlagsResponse();
    message.crosschainFlags = object.crosschainFlags !== undefined && object.crosschainFlags !== null ? CrosschainFlags.fromPartial(object.crosschainFlags) : undefined;
    return message;
  },
  fromAmino(object: QueryGetCrosschainFlagsResponseAmino): QueryGetCrosschainFlagsResponse {
    const message = createBaseQueryGetCrosschainFlagsResponse();
    if (object.crosschain_flags !== undefined && object.crosschain_flags !== null) {
      message.crosschainFlags = CrosschainFlags.fromAmino(object.crosschain_flags);
    }
    return message;
  },
  toAmino(message: QueryGetCrosschainFlagsResponse): QueryGetCrosschainFlagsResponseAmino {
    const obj: any = {};
    obj.crosschain_flags = message.crosschainFlags ? CrosschainFlags.toAmino(message.crosschainFlags) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetCrosschainFlagsResponseAminoMsg): QueryGetCrosschainFlagsResponse {
    return QueryGetCrosschainFlagsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetCrosschainFlagsResponseProtoMsg): QueryGetCrosschainFlagsResponse {
    return QueryGetCrosschainFlagsResponse.decode(message.value);
  },
  toProto(message: QueryGetCrosschainFlagsResponse): Uint8Array {
    return QueryGetCrosschainFlagsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetCrosschainFlagsResponse): QueryGetCrosschainFlagsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetCrosschainFlagsResponse",
      value: QueryGetCrosschainFlagsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetKeygenRequest(): QueryGetKeygenRequest {
  return {};
}
export const QueryGetKeygenRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetKeygenRequest",
  encode(_: QueryGetKeygenRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetKeygenRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetKeygenRequest();
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
  fromPartial(_: Partial<QueryGetKeygenRequest>): QueryGetKeygenRequest {
    const message = createBaseQueryGetKeygenRequest();
    return message;
  },
  fromAmino(_: QueryGetKeygenRequestAmino): QueryGetKeygenRequest {
    const message = createBaseQueryGetKeygenRequest();
    return message;
  },
  toAmino(_: QueryGetKeygenRequest): QueryGetKeygenRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryGetKeygenRequestAminoMsg): QueryGetKeygenRequest {
    return QueryGetKeygenRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetKeygenRequestProtoMsg): QueryGetKeygenRequest {
    return QueryGetKeygenRequest.decode(message.value);
  },
  toProto(message: QueryGetKeygenRequest): Uint8Array {
    return QueryGetKeygenRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetKeygenRequest): QueryGetKeygenRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetKeygenRequest",
      value: QueryGetKeygenRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetKeygenResponse(): QueryGetKeygenResponse {
  return {
    keygen: undefined
  };
}
export const QueryGetKeygenResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetKeygenResponse",
  encode(message: QueryGetKeygenResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.keygen !== undefined) {
      Keygen.encode(message.keygen, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetKeygenResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetKeygenResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.keygen = Keygen.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetKeygenResponse>): QueryGetKeygenResponse {
    const message = createBaseQueryGetKeygenResponse();
    message.keygen = object.keygen !== undefined && object.keygen !== null ? Keygen.fromPartial(object.keygen) : undefined;
    return message;
  },
  fromAmino(object: QueryGetKeygenResponseAmino): QueryGetKeygenResponse {
    const message = createBaseQueryGetKeygenResponse();
    if (object.keygen !== undefined && object.keygen !== null) {
      message.keygen = Keygen.fromAmino(object.keygen);
    }
    return message;
  },
  toAmino(message: QueryGetKeygenResponse): QueryGetKeygenResponseAmino {
    const obj: any = {};
    obj.keygen = message.keygen ? Keygen.toAmino(message.keygen) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetKeygenResponseAminoMsg): QueryGetKeygenResponse {
    return QueryGetKeygenResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetKeygenResponseProtoMsg): QueryGetKeygenResponse {
    return QueryGetKeygenResponse.decode(message.value);
  },
  toProto(message: QueryGetKeygenResponse): Uint8Array {
    return QueryGetKeygenResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetKeygenResponse): QueryGetKeygenResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetKeygenResponse",
      value: QueryGetKeygenResponse.encode(message).finish()
    };
  }
};
function createBaseQueryShowObserverCountRequest(): QueryShowObserverCountRequest {
  return {};
}
export const QueryShowObserverCountRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryShowObserverCountRequest",
  encode(_: QueryShowObserverCountRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryShowObserverCountRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryShowObserverCountRequest();
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
  fromPartial(_: Partial<QueryShowObserverCountRequest>): QueryShowObserverCountRequest {
    const message = createBaseQueryShowObserverCountRequest();
    return message;
  },
  fromAmino(_: QueryShowObserverCountRequestAmino): QueryShowObserverCountRequest {
    const message = createBaseQueryShowObserverCountRequest();
    return message;
  },
  toAmino(_: QueryShowObserverCountRequest): QueryShowObserverCountRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryShowObserverCountRequestAminoMsg): QueryShowObserverCountRequest {
    return QueryShowObserverCountRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryShowObserverCountRequestProtoMsg): QueryShowObserverCountRequest {
    return QueryShowObserverCountRequest.decode(message.value);
  },
  toProto(message: QueryShowObserverCountRequest): Uint8Array {
    return QueryShowObserverCountRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryShowObserverCountRequest): QueryShowObserverCountRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryShowObserverCountRequest",
      value: QueryShowObserverCountRequest.encode(message).finish()
    };
  }
};
function createBaseQueryShowObserverCountResponse(): QueryShowObserverCountResponse {
  return {
    lastObserverCount: undefined
  };
}
export const QueryShowObserverCountResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryShowObserverCountResponse",
  encode(message: QueryShowObserverCountResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.lastObserverCount !== undefined) {
      LastObserverCount.encode(message.lastObserverCount, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryShowObserverCountResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryShowObserverCountResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.lastObserverCount = LastObserverCount.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryShowObserverCountResponse>): QueryShowObserverCountResponse {
    const message = createBaseQueryShowObserverCountResponse();
    message.lastObserverCount = object.lastObserverCount !== undefined && object.lastObserverCount !== null ? LastObserverCount.fromPartial(object.lastObserverCount) : undefined;
    return message;
  },
  fromAmino(object: QueryShowObserverCountResponseAmino): QueryShowObserverCountResponse {
    const message = createBaseQueryShowObserverCountResponse();
    if (object.last_observer_count !== undefined && object.last_observer_count !== null) {
      message.lastObserverCount = LastObserverCount.fromAmino(object.last_observer_count);
    }
    return message;
  },
  toAmino(message: QueryShowObserverCountResponse): QueryShowObserverCountResponseAmino {
    const obj: any = {};
    obj.last_observer_count = message.lastObserverCount ? LastObserverCount.toAmino(message.lastObserverCount) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryShowObserverCountResponseAminoMsg): QueryShowObserverCountResponse {
    return QueryShowObserverCountResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryShowObserverCountResponseProtoMsg): QueryShowObserverCountResponse {
    return QueryShowObserverCountResponse.decode(message.value);
  },
  toProto(message: QueryShowObserverCountResponse): Uint8Array {
    return QueryShowObserverCountResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryShowObserverCountResponse): QueryShowObserverCountResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryShowObserverCountResponse",
      value: QueryShowObserverCountResponse.encode(message).finish()
    };
  }
};
function createBaseQueryBlameByIdentifierRequest(): QueryBlameByIdentifierRequest {
  return {
    blameIdentifier: ""
  };
}
export const QueryBlameByIdentifierRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryBlameByIdentifierRequest",
  encode(message: QueryBlameByIdentifierRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.blameIdentifier !== "") {
      writer.uint32(10).string(message.blameIdentifier);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryBlameByIdentifierRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBlameByIdentifierRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blameIdentifier = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryBlameByIdentifierRequest>): QueryBlameByIdentifierRequest {
    const message = createBaseQueryBlameByIdentifierRequest();
    message.blameIdentifier = object.blameIdentifier ?? "";
    return message;
  },
  fromAmino(object: QueryBlameByIdentifierRequestAmino): QueryBlameByIdentifierRequest {
    const message = createBaseQueryBlameByIdentifierRequest();
    if (object.blame_identifier !== undefined && object.blame_identifier !== null) {
      message.blameIdentifier = object.blame_identifier;
    }
    return message;
  },
  toAmino(message: QueryBlameByIdentifierRequest): QueryBlameByIdentifierRequestAmino {
    const obj: any = {};
    obj.blame_identifier = message.blameIdentifier;
    return obj;
  },
  fromAminoMsg(object: QueryBlameByIdentifierRequestAminoMsg): QueryBlameByIdentifierRequest {
    return QueryBlameByIdentifierRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryBlameByIdentifierRequestProtoMsg): QueryBlameByIdentifierRequest {
    return QueryBlameByIdentifierRequest.decode(message.value);
  },
  toProto(message: QueryBlameByIdentifierRequest): Uint8Array {
    return QueryBlameByIdentifierRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryBlameByIdentifierRequest): QueryBlameByIdentifierRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryBlameByIdentifierRequest",
      value: QueryBlameByIdentifierRequest.encode(message).finish()
    };
  }
};
function createBaseQueryBlameByIdentifierResponse(): QueryBlameByIdentifierResponse {
  return {
    blameInfo: undefined
  };
}
export const QueryBlameByIdentifierResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryBlameByIdentifierResponse",
  encode(message: QueryBlameByIdentifierResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.blameInfo !== undefined) {
      Blame.encode(message.blameInfo, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryBlameByIdentifierResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBlameByIdentifierResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blameInfo = Blame.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryBlameByIdentifierResponse>): QueryBlameByIdentifierResponse {
    const message = createBaseQueryBlameByIdentifierResponse();
    message.blameInfo = object.blameInfo !== undefined && object.blameInfo !== null ? Blame.fromPartial(object.blameInfo) : undefined;
    return message;
  },
  fromAmino(object: QueryBlameByIdentifierResponseAmino): QueryBlameByIdentifierResponse {
    const message = createBaseQueryBlameByIdentifierResponse();
    if (object.blame_info !== undefined && object.blame_info !== null) {
      message.blameInfo = Blame.fromAmino(object.blame_info);
    }
    return message;
  },
  toAmino(message: QueryBlameByIdentifierResponse): QueryBlameByIdentifierResponseAmino {
    const obj: any = {};
    obj.blame_info = message.blameInfo ? Blame.toAmino(message.blameInfo) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryBlameByIdentifierResponseAminoMsg): QueryBlameByIdentifierResponse {
    return QueryBlameByIdentifierResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryBlameByIdentifierResponseProtoMsg): QueryBlameByIdentifierResponse {
    return QueryBlameByIdentifierResponse.decode(message.value);
  },
  toProto(message: QueryBlameByIdentifierResponse): Uint8Array {
    return QueryBlameByIdentifierResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryBlameByIdentifierResponse): QueryBlameByIdentifierResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryBlameByIdentifierResponse",
      value: QueryBlameByIdentifierResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllBlameRecordsRequest(): QueryAllBlameRecordsRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllBlameRecordsRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllBlameRecordsRequest",
  encode(message: QueryAllBlameRecordsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllBlameRecordsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllBlameRecordsRequest();
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
  fromPartial(object: Partial<QueryAllBlameRecordsRequest>): QueryAllBlameRecordsRequest {
    const message = createBaseQueryAllBlameRecordsRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllBlameRecordsRequestAmino): QueryAllBlameRecordsRequest {
    const message = createBaseQueryAllBlameRecordsRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllBlameRecordsRequest): QueryAllBlameRecordsRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllBlameRecordsRequestAminoMsg): QueryAllBlameRecordsRequest {
    return QueryAllBlameRecordsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllBlameRecordsRequestProtoMsg): QueryAllBlameRecordsRequest {
    return QueryAllBlameRecordsRequest.decode(message.value);
  },
  toProto(message: QueryAllBlameRecordsRequest): Uint8Array {
    return QueryAllBlameRecordsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllBlameRecordsRequest): QueryAllBlameRecordsRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllBlameRecordsRequest",
      value: QueryAllBlameRecordsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllBlameRecordsResponse(): QueryAllBlameRecordsResponse {
  return {
    blameInfo: [],
    pagination: undefined
  };
}
export const QueryAllBlameRecordsResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllBlameRecordsResponse",
  encode(message: QueryAllBlameRecordsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.blameInfo) {
      Blame.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllBlameRecordsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllBlameRecordsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blameInfo.push(Blame.decode(reader, reader.uint32()));
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
  fromPartial(object: Partial<QueryAllBlameRecordsResponse>): QueryAllBlameRecordsResponse {
    const message = createBaseQueryAllBlameRecordsResponse();
    message.blameInfo = object.blameInfo?.map(e => Blame.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllBlameRecordsResponseAmino): QueryAllBlameRecordsResponse {
    const message = createBaseQueryAllBlameRecordsResponse();
    message.blameInfo = object.blame_info?.map(e => Blame.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllBlameRecordsResponse): QueryAllBlameRecordsResponseAmino {
    const obj: any = {};
    if (message.blameInfo) {
      obj.blame_info = message.blameInfo.map(e => e ? Blame.toAmino(e) : undefined);
    } else {
      obj.blame_info = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllBlameRecordsResponseAminoMsg): QueryAllBlameRecordsResponse {
    return QueryAllBlameRecordsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllBlameRecordsResponseProtoMsg): QueryAllBlameRecordsResponse {
    return QueryAllBlameRecordsResponse.decode(message.value);
  },
  toProto(message: QueryAllBlameRecordsResponse): Uint8Array {
    return QueryAllBlameRecordsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllBlameRecordsResponse): QueryAllBlameRecordsResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllBlameRecordsResponse",
      value: QueryAllBlameRecordsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryBlameByChainAndNonceRequest(): QueryBlameByChainAndNonceRequest {
  return {
    chainId: BigInt(0),
    nonce: BigInt(0)
  };
}
export const QueryBlameByChainAndNonceRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryBlameByChainAndNonceRequest",
  encode(message: QueryBlameByChainAndNonceRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    if (message.nonce !== BigInt(0)) {
      writer.uint32(16).int64(message.nonce);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryBlameByChainAndNonceRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBlameByChainAndNonceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        case 2:
          message.nonce = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryBlameByChainAndNonceRequest>): QueryBlameByChainAndNonceRequest {
    const message = createBaseQueryBlameByChainAndNonceRequest();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    message.nonce = object.nonce !== undefined && object.nonce !== null ? BigInt(object.nonce.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryBlameByChainAndNonceRequestAmino): QueryBlameByChainAndNonceRequest {
    const message = createBaseQueryBlameByChainAndNonceRequest();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = BigInt(object.nonce);
    }
    return message;
  },
  toAmino(message: QueryBlameByChainAndNonceRequest): QueryBlameByChainAndNonceRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    obj.nonce = message.nonce ? message.nonce.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryBlameByChainAndNonceRequestAminoMsg): QueryBlameByChainAndNonceRequest {
    return QueryBlameByChainAndNonceRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryBlameByChainAndNonceRequestProtoMsg): QueryBlameByChainAndNonceRequest {
    return QueryBlameByChainAndNonceRequest.decode(message.value);
  },
  toProto(message: QueryBlameByChainAndNonceRequest): Uint8Array {
    return QueryBlameByChainAndNonceRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryBlameByChainAndNonceRequest): QueryBlameByChainAndNonceRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryBlameByChainAndNonceRequest",
      value: QueryBlameByChainAndNonceRequest.encode(message).finish()
    };
  }
};
function createBaseQueryBlameByChainAndNonceResponse(): QueryBlameByChainAndNonceResponse {
  return {
    blameInfo: []
  };
}
export const QueryBlameByChainAndNonceResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryBlameByChainAndNonceResponse",
  encode(message: QueryBlameByChainAndNonceResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.blameInfo) {
      Blame.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryBlameByChainAndNonceResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBlameByChainAndNonceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blameInfo.push(Blame.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryBlameByChainAndNonceResponse>): QueryBlameByChainAndNonceResponse {
    const message = createBaseQueryBlameByChainAndNonceResponse();
    message.blameInfo = object.blameInfo?.map(e => Blame.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: QueryBlameByChainAndNonceResponseAmino): QueryBlameByChainAndNonceResponse {
    const message = createBaseQueryBlameByChainAndNonceResponse();
    message.blameInfo = object.blame_info?.map(e => Blame.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: QueryBlameByChainAndNonceResponse): QueryBlameByChainAndNonceResponseAmino {
    const obj: any = {};
    if (message.blameInfo) {
      obj.blame_info = message.blameInfo.map(e => e ? Blame.toAmino(e) : undefined);
    } else {
      obj.blame_info = [];
    }
    return obj;
  },
  fromAminoMsg(object: QueryBlameByChainAndNonceResponseAminoMsg): QueryBlameByChainAndNonceResponse {
    return QueryBlameByChainAndNonceResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryBlameByChainAndNonceResponseProtoMsg): QueryBlameByChainAndNonceResponse {
    return QueryBlameByChainAndNonceResponse.decode(message.value);
  },
  toProto(message: QueryBlameByChainAndNonceResponse): Uint8Array {
    return QueryBlameByChainAndNonceResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryBlameByChainAndNonceResponse): QueryBlameByChainAndNonceResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryBlameByChainAndNonceResponse",
      value: QueryBlameByChainAndNonceResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllBlockHeaderRequest(): QueryAllBlockHeaderRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllBlockHeaderRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllBlockHeaderRequest",
  encode(message: QueryAllBlockHeaderRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllBlockHeaderRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllBlockHeaderRequest();
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
  fromPartial(object: Partial<QueryAllBlockHeaderRequest>): QueryAllBlockHeaderRequest {
    const message = createBaseQueryAllBlockHeaderRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllBlockHeaderRequestAmino): QueryAllBlockHeaderRequest {
    const message = createBaseQueryAllBlockHeaderRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllBlockHeaderRequest): QueryAllBlockHeaderRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllBlockHeaderRequestAminoMsg): QueryAllBlockHeaderRequest {
    return QueryAllBlockHeaderRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllBlockHeaderRequestProtoMsg): QueryAllBlockHeaderRequest {
    return QueryAllBlockHeaderRequest.decode(message.value);
  },
  toProto(message: QueryAllBlockHeaderRequest): Uint8Array {
    return QueryAllBlockHeaderRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllBlockHeaderRequest): QueryAllBlockHeaderRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllBlockHeaderRequest",
      value: QueryAllBlockHeaderRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllBlockHeaderResponse(): QueryAllBlockHeaderResponse {
  return {
    blockHeaders: [],
    pagination: undefined
  };
}
export const QueryAllBlockHeaderResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryAllBlockHeaderResponse",
  encode(message: QueryAllBlockHeaderResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.blockHeaders) {
      BlockHeader.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllBlockHeaderResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllBlockHeaderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blockHeaders.push(BlockHeader.decode(reader, reader.uint32()));
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
  fromPartial(object: Partial<QueryAllBlockHeaderResponse>): QueryAllBlockHeaderResponse {
    const message = createBaseQueryAllBlockHeaderResponse();
    message.blockHeaders = object.blockHeaders?.map(e => BlockHeader.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllBlockHeaderResponseAmino): QueryAllBlockHeaderResponse {
    const message = createBaseQueryAllBlockHeaderResponse();
    message.blockHeaders = object.block_headers?.map(e => BlockHeader.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllBlockHeaderResponse): QueryAllBlockHeaderResponseAmino {
    const obj: any = {};
    if (message.blockHeaders) {
      obj.block_headers = message.blockHeaders.map(e => e ? BlockHeader.toAmino(e) : undefined);
    } else {
      obj.block_headers = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllBlockHeaderResponseAminoMsg): QueryAllBlockHeaderResponse {
    return QueryAllBlockHeaderResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllBlockHeaderResponseProtoMsg): QueryAllBlockHeaderResponse {
    return QueryAllBlockHeaderResponse.decode(message.value);
  },
  toProto(message: QueryAllBlockHeaderResponse): Uint8Array {
    return QueryAllBlockHeaderResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllBlockHeaderResponse): QueryAllBlockHeaderResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryAllBlockHeaderResponse",
      value: QueryAllBlockHeaderResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetBlockHeaderByHashRequest(): QueryGetBlockHeaderByHashRequest {
  return {
    blockHash: new Uint8Array()
  };
}
export const QueryGetBlockHeaderByHashRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderByHashRequest",
  encode(message: QueryGetBlockHeaderByHashRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.blockHash.length !== 0) {
      writer.uint32(10).bytes(message.blockHash);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetBlockHeaderByHashRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetBlockHeaderByHashRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blockHash = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetBlockHeaderByHashRequest>): QueryGetBlockHeaderByHashRequest {
    const message = createBaseQueryGetBlockHeaderByHashRequest();
    message.blockHash = object.blockHash ?? new Uint8Array();
    return message;
  },
  fromAmino(object: QueryGetBlockHeaderByHashRequestAmino): QueryGetBlockHeaderByHashRequest {
    const message = createBaseQueryGetBlockHeaderByHashRequest();
    if (object.block_hash !== undefined && object.block_hash !== null) {
      message.blockHash = bytesFromBase64(object.block_hash);
    }
    return message;
  },
  toAmino(message: QueryGetBlockHeaderByHashRequest): QueryGetBlockHeaderByHashRequestAmino {
    const obj: any = {};
    obj.block_hash = message.blockHash ? base64FromBytes(message.blockHash) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetBlockHeaderByHashRequestAminoMsg): QueryGetBlockHeaderByHashRequest {
    return QueryGetBlockHeaderByHashRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetBlockHeaderByHashRequestProtoMsg): QueryGetBlockHeaderByHashRequest {
    return QueryGetBlockHeaderByHashRequest.decode(message.value);
  },
  toProto(message: QueryGetBlockHeaderByHashRequest): Uint8Array {
    return QueryGetBlockHeaderByHashRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetBlockHeaderByHashRequest): QueryGetBlockHeaderByHashRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderByHashRequest",
      value: QueryGetBlockHeaderByHashRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetBlockHeaderByHashResponse(): QueryGetBlockHeaderByHashResponse {
  return {
    blockHeader: undefined
  };
}
export const QueryGetBlockHeaderByHashResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderByHashResponse",
  encode(message: QueryGetBlockHeaderByHashResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.blockHeader !== undefined) {
      BlockHeader.encode(message.blockHeader, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetBlockHeaderByHashResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetBlockHeaderByHashResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blockHeader = BlockHeader.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetBlockHeaderByHashResponse>): QueryGetBlockHeaderByHashResponse {
    const message = createBaseQueryGetBlockHeaderByHashResponse();
    message.blockHeader = object.blockHeader !== undefined && object.blockHeader !== null ? BlockHeader.fromPartial(object.blockHeader) : undefined;
    return message;
  },
  fromAmino(object: QueryGetBlockHeaderByHashResponseAmino): QueryGetBlockHeaderByHashResponse {
    const message = createBaseQueryGetBlockHeaderByHashResponse();
    if (object.block_header !== undefined && object.block_header !== null) {
      message.blockHeader = BlockHeader.fromAmino(object.block_header);
    }
    return message;
  },
  toAmino(message: QueryGetBlockHeaderByHashResponse): QueryGetBlockHeaderByHashResponseAmino {
    const obj: any = {};
    obj.block_header = message.blockHeader ? BlockHeader.toAmino(message.blockHeader) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetBlockHeaderByHashResponseAminoMsg): QueryGetBlockHeaderByHashResponse {
    return QueryGetBlockHeaderByHashResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetBlockHeaderByHashResponseProtoMsg): QueryGetBlockHeaderByHashResponse {
    return QueryGetBlockHeaderByHashResponse.decode(message.value);
  },
  toProto(message: QueryGetBlockHeaderByHashResponse): Uint8Array {
    return QueryGetBlockHeaderByHashResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetBlockHeaderByHashResponse): QueryGetBlockHeaderByHashResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderByHashResponse",
      value: QueryGetBlockHeaderByHashResponse.encode(message).finish()
    };
  }
};
function createBaseQueryGetBlockHeaderStateRequest(): QueryGetBlockHeaderStateRequest {
  return {
    chainId: BigInt(0)
  };
}
export const QueryGetBlockHeaderStateRequest = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderStateRequest",
  encode(message: QueryGetBlockHeaderStateRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.chainId !== BigInt(0)) {
      writer.uint32(8).int64(message.chainId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetBlockHeaderStateRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetBlockHeaderStateRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chainId = reader.int64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetBlockHeaderStateRequest>): QueryGetBlockHeaderStateRequest {
    const message = createBaseQueryGetBlockHeaderStateRequest();
    message.chainId = object.chainId !== undefined && object.chainId !== null ? BigInt(object.chainId.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: QueryGetBlockHeaderStateRequestAmino): QueryGetBlockHeaderStateRequest {
    const message = createBaseQueryGetBlockHeaderStateRequest();
    if (object.chain_id !== undefined && object.chain_id !== null) {
      message.chainId = BigInt(object.chain_id);
    }
    return message;
  },
  toAmino(message: QueryGetBlockHeaderStateRequest): QueryGetBlockHeaderStateRequestAmino {
    const obj: any = {};
    obj.chain_id = message.chainId ? message.chainId.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetBlockHeaderStateRequestAminoMsg): QueryGetBlockHeaderStateRequest {
    return QueryGetBlockHeaderStateRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetBlockHeaderStateRequestProtoMsg): QueryGetBlockHeaderStateRequest {
    return QueryGetBlockHeaderStateRequest.decode(message.value);
  },
  toProto(message: QueryGetBlockHeaderStateRequest): Uint8Array {
    return QueryGetBlockHeaderStateRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetBlockHeaderStateRequest): QueryGetBlockHeaderStateRequestProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderStateRequest",
      value: QueryGetBlockHeaderStateRequest.encode(message).finish()
    };
  }
};
function createBaseQueryGetBlockHeaderStateResponse(): QueryGetBlockHeaderStateResponse {
  return {
    blockHeaderState: undefined
  };
}
export const QueryGetBlockHeaderStateResponse = {
  typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderStateResponse",
  encode(message: QueryGetBlockHeaderStateResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.blockHeaderState !== undefined) {
      BlockHeaderState.encode(message.blockHeaderState, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetBlockHeaderStateResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetBlockHeaderStateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blockHeaderState = BlockHeaderState.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetBlockHeaderStateResponse>): QueryGetBlockHeaderStateResponse {
    const message = createBaseQueryGetBlockHeaderStateResponse();
    message.blockHeaderState = object.blockHeaderState !== undefined && object.blockHeaderState !== null ? BlockHeaderState.fromPartial(object.blockHeaderState) : undefined;
    return message;
  },
  fromAmino(object: QueryGetBlockHeaderStateResponseAmino): QueryGetBlockHeaderStateResponse {
    const message = createBaseQueryGetBlockHeaderStateResponse();
    if (object.block_header_state !== undefined && object.block_header_state !== null) {
      message.blockHeaderState = BlockHeaderState.fromAmino(object.block_header_state);
    }
    return message;
  },
  toAmino(message: QueryGetBlockHeaderStateResponse): QueryGetBlockHeaderStateResponseAmino {
    const obj: any = {};
    obj.block_header_state = message.blockHeaderState ? BlockHeaderState.toAmino(message.blockHeaderState) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryGetBlockHeaderStateResponseAminoMsg): QueryGetBlockHeaderStateResponse {
    return QueryGetBlockHeaderStateResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetBlockHeaderStateResponseProtoMsg): QueryGetBlockHeaderStateResponse {
    return QueryGetBlockHeaderStateResponse.decode(message.value);
  },
  toProto(message: QueryGetBlockHeaderStateResponse): Uint8Array {
    return QueryGetBlockHeaderStateResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryGetBlockHeaderStateResponse): QueryGetBlockHeaderStateResponseProtoMsg {
    return {
      typeUrl: "/zetachain.zetacore.observer.QueryGetBlockHeaderStateResponse",
      value: QueryGetBlockHeaderStateResponse.encode(message).finish()
    };
  }
};