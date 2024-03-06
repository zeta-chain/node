import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryHasVotedRequest, QueryHasVotedResponse, QueryBallotByIdentifierRequest, QueryBallotByIdentifierResponse, QueryObserverSet, QueryObserverSetResponse, QuerySupportedChains, QuerySupportedChainsResponse, QueryGetChainParamsForChainRequest, QueryGetChainParamsForChainResponse, QueryGetChainParamsRequest, QueryGetChainParamsResponse, QueryGetNodeAccountRequest, QueryGetNodeAccountResponse, QueryAllNodeAccountRequest, QueryAllNodeAccountResponse, QueryGetCrosschainFlagsRequest, QueryGetCrosschainFlagsResponse, QueryGetKeygenRequest, QueryGetKeygenResponse, QueryShowObserverCountRequest, QueryShowObserverCountResponse, QueryBlameByIdentifierRequest, QueryBlameByIdentifierResponse, QueryAllBlameRecordsRequest, QueryAllBlameRecordsResponse, QueryBlameByChainAndNonceRequest, QueryBlameByChainAndNonceResponse, QueryAllBlockHeaderRequest, QueryAllBlockHeaderResponse, QueryGetBlockHeaderByHashRequest, QueryGetBlockHeaderByHashResponse, QueryGetBlockHeaderStateRequest, QueryGetBlockHeaderStateResponse, QueryProveRequest, QueryProveResponse, QueryGetTssAddressRequest, QueryGetTssAddressResponse, QueryGetTssAddressByFinalizedHeightRequest, QueryGetTssAddressByFinalizedHeightResponse, QueryGetTSSRequest, QueryGetTSSResponse, QueryTssHistoryRequest, QueryTssHistoryResponse, QueryAllPendingNoncesRequest, QueryAllPendingNoncesResponse, QueryPendingNoncesByChainRequest, QueryPendingNoncesByChainResponse, QueryGetChainNoncesRequest, QueryGetChainNoncesResponse, QueryAllChainNoncesRequest, QueryAllChainNoncesResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Query if a voter has voted for a ballot */
  hasVoted(request: QueryHasVotedRequest): Promise<QueryHasVotedResponse>;
  /** Queries a list of VoterByIdentifier items. */
  ballotByIdentifier(request: QueryBallotByIdentifierRequest): Promise<QueryBallotByIdentifierResponse>;
  /** Queries a list of ObserversByChainAndType items. */
  observerSet(request?: QueryObserverSet): Promise<QueryObserverSetResponse>;
  supportedChains(request?: QuerySupportedChains): Promise<QuerySupportedChainsResponse>;
  /** Queries a list of GetChainParamsForChain items. */
  getChainParamsForChain(request: QueryGetChainParamsForChainRequest): Promise<QueryGetChainParamsForChainResponse>;
  /** Queries a list of GetChainParams items. */
  getChainParams(request?: QueryGetChainParamsRequest): Promise<QueryGetChainParamsResponse>;
  /** Queries a nodeAccount by index. */
  nodeAccount(request: QueryGetNodeAccountRequest): Promise<QueryGetNodeAccountResponse>;
  /** Queries a list of nodeAccount items. */
  nodeAccountAll(request?: QueryAllNodeAccountRequest): Promise<QueryAllNodeAccountResponse>;
  crosschainFlags(request?: QueryGetCrosschainFlagsRequest): Promise<QueryGetCrosschainFlagsResponse>;
  /** Queries a keygen by index. */
  keygen(request?: QueryGetKeygenRequest): Promise<QueryGetKeygenResponse>;
  /** Queries a list of ShowObserverCount items. */
  showObserverCount(request?: QueryShowObserverCountRequest): Promise<QueryShowObserverCountResponse>;
  /** Queries a list of VoterByIdentifier items. */
  blameByIdentifier(request: QueryBlameByIdentifierRequest): Promise<QueryBlameByIdentifierResponse>;
  /** Queries a list of VoterByIdentifier items. */
  getAllBlameRecords(request?: QueryAllBlameRecordsRequest): Promise<QueryAllBlameRecordsResponse>;
  /** Queries a list of VoterByIdentifier items. */
  blamesByChainAndNonce(request: QueryBlameByChainAndNonceRequest): Promise<QueryBlameByChainAndNonceResponse>;
  getAllBlockHeaders(request?: QueryAllBlockHeaderRequest): Promise<QueryAllBlockHeaderResponse>;
  getBlockHeaderByHash(request: QueryGetBlockHeaderByHashRequest): Promise<QueryGetBlockHeaderByHashResponse>;
  getBlockHeaderStateByChain(request: QueryGetBlockHeaderStateRequest): Promise<QueryGetBlockHeaderStateResponse>;
  /** merkle proof verification */
  prove(request: QueryProveRequest): Promise<QueryProveResponse>;
  /** Queries a list of GetTssAddress items. */
  getTssAddress(request: QueryGetTssAddressRequest): Promise<QueryGetTssAddressResponse>;
  getTssAddressByFinalizedHeight(request: QueryGetTssAddressByFinalizedHeightRequest): Promise<QueryGetTssAddressByFinalizedHeightResponse>;
  /** Queries a tSS by index. */
  tSS(request?: QueryGetTSSRequest): Promise<QueryGetTSSResponse>;
  tssHistory(request?: QueryTssHistoryRequest): Promise<QueryTssHistoryResponse>;
  pendingNoncesAll(request?: QueryAllPendingNoncesRequest): Promise<QueryAllPendingNoncesResponse>;
  pendingNoncesByChain(request: QueryPendingNoncesByChainRequest): Promise<QueryPendingNoncesByChainResponse>;
  /** Queries a chainNonces by index. */
  chainNonces(request: QueryGetChainNoncesRequest): Promise<QueryGetChainNoncesResponse>;
  /** Queries a list of chainNonces items. */
  chainNoncesAll(request?: QueryAllChainNoncesRequest): Promise<QueryAllChainNoncesResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.params = this.params.bind(this);
    this.hasVoted = this.hasVoted.bind(this);
    this.ballotByIdentifier = this.ballotByIdentifier.bind(this);
    this.observerSet = this.observerSet.bind(this);
    this.supportedChains = this.supportedChains.bind(this);
    this.getChainParamsForChain = this.getChainParamsForChain.bind(this);
    this.getChainParams = this.getChainParams.bind(this);
    this.nodeAccount = this.nodeAccount.bind(this);
    this.nodeAccountAll = this.nodeAccountAll.bind(this);
    this.crosschainFlags = this.crosschainFlags.bind(this);
    this.keygen = this.keygen.bind(this);
    this.showObserverCount = this.showObserverCount.bind(this);
    this.blameByIdentifier = this.blameByIdentifier.bind(this);
    this.getAllBlameRecords = this.getAllBlameRecords.bind(this);
    this.blamesByChainAndNonce = this.blamesByChainAndNonce.bind(this);
    this.getAllBlockHeaders = this.getAllBlockHeaders.bind(this);
    this.getBlockHeaderByHash = this.getBlockHeaderByHash.bind(this);
    this.getBlockHeaderStateByChain = this.getBlockHeaderStateByChain.bind(this);
    this.prove = this.prove.bind(this);
    this.getTssAddress = this.getTssAddress.bind(this);
    this.getTssAddressByFinalizedHeight = this.getTssAddressByFinalizedHeight.bind(this);
    this.tSS = this.tSS.bind(this);
    this.tssHistory = this.tssHistory.bind(this);
    this.pendingNoncesAll = this.pendingNoncesAll.bind(this);
    this.pendingNoncesByChain = this.pendingNoncesByChain.bind(this);
    this.chainNonces = this.chainNonces.bind(this);
    this.chainNoncesAll = this.chainNoncesAll.bind(this);
  }
  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new BinaryReader(data)));
  }
  hasVoted(request: QueryHasVotedRequest): Promise<QueryHasVotedResponse> {
    const data = QueryHasVotedRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "HasVoted", data);
    return promise.then(data => QueryHasVotedResponse.decode(new BinaryReader(data)));
  }
  ballotByIdentifier(request: QueryBallotByIdentifierRequest): Promise<QueryBallotByIdentifierResponse> {
    const data = QueryBallotByIdentifierRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "BallotByIdentifier", data);
    return promise.then(data => QueryBallotByIdentifierResponse.decode(new BinaryReader(data)));
  }
  observerSet(request: QueryObserverSet = {}): Promise<QueryObserverSetResponse> {
    const data = QueryObserverSet.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "ObserverSet", data);
    return promise.then(data => QueryObserverSetResponse.decode(new BinaryReader(data)));
  }
  supportedChains(request: QuerySupportedChains = {}): Promise<QuerySupportedChainsResponse> {
    const data = QuerySupportedChains.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "SupportedChains", data);
    return promise.then(data => QuerySupportedChainsResponse.decode(new BinaryReader(data)));
  }
  getChainParamsForChain(request: QueryGetChainParamsForChainRequest): Promise<QueryGetChainParamsForChainResponse> {
    const data = QueryGetChainParamsForChainRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "GetChainParamsForChain", data);
    return promise.then(data => QueryGetChainParamsForChainResponse.decode(new BinaryReader(data)));
  }
  getChainParams(request: QueryGetChainParamsRequest = {}): Promise<QueryGetChainParamsResponse> {
    const data = QueryGetChainParamsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "GetChainParams", data);
    return promise.then(data => QueryGetChainParamsResponse.decode(new BinaryReader(data)));
  }
  nodeAccount(request: QueryGetNodeAccountRequest): Promise<QueryGetNodeAccountResponse> {
    const data = QueryGetNodeAccountRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "NodeAccount", data);
    return promise.then(data => QueryGetNodeAccountResponse.decode(new BinaryReader(data)));
  }
  nodeAccountAll(request: QueryAllNodeAccountRequest = {
    pagination: undefined
  }): Promise<QueryAllNodeAccountResponse> {
    const data = QueryAllNodeAccountRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "NodeAccountAll", data);
    return promise.then(data => QueryAllNodeAccountResponse.decode(new BinaryReader(data)));
  }
  crosschainFlags(request: QueryGetCrosschainFlagsRequest = {}): Promise<QueryGetCrosschainFlagsResponse> {
    const data = QueryGetCrosschainFlagsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "CrosschainFlags", data);
    return promise.then(data => QueryGetCrosschainFlagsResponse.decode(new BinaryReader(data)));
  }
  keygen(request: QueryGetKeygenRequest = {}): Promise<QueryGetKeygenResponse> {
    const data = QueryGetKeygenRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "Keygen", data);
    return promise.then(data => QueryGetKeygenResponse.decode(new BinaryReader(data)));
  }
  showObserverCount(request: QueryShowObserverCountRequest = {}): Promise<QueryShowObserverCountResponse> {
    const data = QueryShowObserverCountRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "ShowObserverCount", data);
    return promise.then(data => QueryShowObserverCountResponse.decode(new BinaryReader(data)));
  }
  blameByIdentifier(request: QueryBlameByIdentifierRequest): Promise<QueryBlameByIdentifierResponse> {
    const data = QueryBlameByIdentifierRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "BlameByIdentifier", data);
    return promise.then(data => QueryBlameByIdentifierResponse.decode(new BinaryReader(data)));
  }
  getAllBlameRecords(request: QueryAllBlameRecordsRequest = {
    pagination: undefined
  }): Promise<QueryAllBlameRecordsResponse> {
    const data = QueryAllBlameRecordsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "GetAllBlameRecords", data);
    return promise.then(data => QueryAllBlameRecordsResponse.decode(new BinaryReader(data)));
  }
  blamesByChainAndNonce(request: QueryBlameByChainAndNonceRequest): Promise<QueryBlameByChainAndNonceResponse> {
    const data = QueryBlameByChainAndNonceRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "BlamesByChainAndNonce", data);
    return promise.then(data => QueryBlameByChainAndNonceResponse.decode(new BinaryReader(data)));
  }
  getAllBlockHeaders(request: QueryAllBlockHeaderRequest = {
    pagination: undefined
  }): Promise<QueryAllBlockHeaderResponse> {
    const data = QueryAllBlockHeaderRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "GetAllBlockHeaders", data);
    return promise.then(data => QueryAllBlockHeaderResponse.decode(new BinaryReader(data)));
  }
  getBlockHeaderByHash(request: QueryGetBlockHeaderByHashRequest): Promise<QueryGetBlockHeaderByHashResponse> {
    const data = QueryGetBlockHeaderByHashRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "GetBlockHeaderByHash", data);
    return promise.then(data => QueryGetBlockHeaderByHashResponse.decode(new BinaryReader(data)));
  }
  getBlockHeaderStateByChain(request: QueryGetBlockHeaderStateRequest): Promise<QueryGetBlockHeaderStateResponse> {
    const data = QueryGetBlockHeaderStateRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "GetBlockHeaderStateByChain", data);
    return promise.then(data => QueryGetBlockHeaderStateResponse.decode(new BinaryReader(data)));
  }
  prove(request: QueryProveRequest): Promise<QueryProveResponse> {
    const data = QueryProveRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "Prove", data);
    return promise.then(data => QueryProveResponse.decode(new BinaryReader(data)));
  }
  getTssAddress(request: QueryGetTssAddressRequest): Promise<QueryGetTssAddressResponse> {
    const data = QueryGetTssAddressRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "GetTssAddress", data);
    return promise.then(data => QueryGetTssAddressResponse.decode(new BinaryReader(data)));
  }
  getTssAddressByFinalizedHeight(request: QueryGetTssAddressByFinalizedHeightRequest): Promise<QueryGetTssAddressByFinalizedHeightResponse> {
    const data = QueryGetTssAddressByFinalizedHeightRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "GetTssAddressByFinalizedHeight", data);
    return promise.then(data => QueryGetTssAddressByFinalizedHeightResponse.decode(new BinaryReader(data)));
  }
  tSS(request: QueryGetTSSRequest = {}): Promise<QueryGetTSSResponse> {
    const data = QueryGetTSSRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "TSS", data);
    return promise.then(data => QueryGetTSSResponse.decode(new BinaryReader(data)));
  }
  tssHistory(request: QueryTssHistoryRequest = {
    pagination: undefined
  }): Promise<QueryTssHistoryResponse> {
    const data = QueryTssHistoryRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "TssHistory", data);
    return promise.then(data => QueryTssHistoryResponse.decode(new BinaryReader(data)));
  }
  pendingNoncesAll(request: QueryAllPendingNoncesRequest = {
    pagination: undefined
  }): Promise<QueryAllPendingNoncesResponse> {
    const data = QueryAllPendingNoncesRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "PendingNoncesAll", data);
    return promise.then(data => QueryAllPendingNoncesResponse.decode(new BinaryReader(data)));
  }
  pendingNoncesByChain(request: QueryPendingNoncesByChainRequest): Promise<QueryPendingNoncesByChainResponse> {
    const data = QueryPendingNoncesByChainRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "PendingNoncesByChain", data);
    return promise.then(data => QueryPendingNoncesByChainResponse.decode(new BinaryReader(data)));
  }
  chainNonces(request: QueryGetChainNoncesRequest): Promise<QueryGetChainNoncesResponse> {
    const data = QueryGetChainNoncesRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "ChainNonces", data);
    return promise.then(data => QueryGetChainNoncesResponse.decode(new BinaryReader(data)));
  }
  chainNoncesAll(request: QueryAllChainNoncesRequest = {
    pagination: undefined
  }): Promise<QueryAllChainNoncesResponse> {
    const data = QueryAllChainNoncesRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.observer.Query", "ChainNoncesAll", data);
    return promise.then(data => QueryAllChainNoncesResponse.decode(new BinaryReader(data)));
  }
}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },
    hasVoted(request: QueryHasVotedRequest): Promise<QueryHasVotedResponse> {
      return queryService.hasVoted(request);
    },
    ballotByIdentifier(request: QueryBallotByIdentifierRequest): Promise<QueryBallotByIdentifierResponse> {
      return queryService.ballotByIdentifier(request);
    },
    observerSet(request?: QueryObserverSet): Promise<QueryObserverSetResponse> {
      return queryService.observerSet(request);
    },
    supportedChains(request?: QuerySupportedChains): Promise<QuerySupportedChainsResponse> {
      return queryService.supportedChains(request);
    },
    getChainParamsForChain(request: QueryGetChainParamsForChainRequest): Promise<QueryGetChainParamsForChainResponse> {
      return queryService.getChainParamsForChain(request);
    },
    getChainParams(request?: QueryGetChainParamsRequest): Promise<QueryGetChainParamsResponse> {
      return queryService.getChainParams(request);
    },
    nodeAccount(request: QueryGetNodeAccountRequest): Promise<QueryGetNodeAccountResponse> {
      return queryService.nodeAccount(request);
    },
    nodeAccountAll(request?: QueryAllNodeAccountRequest): Promise<QueryAllNodeAccountResponse> {
      return queryService.nodeAccountAll(request);
    },
    crosschainFlags(request?: QueryGetCrosschainFlagsRequest): Promise<QueryGetCrosschainFlagsResponse> {
      return queryService.crosschainFlags(request);
    },
    keygen(request?: QueryGetKeygenRequest): Promise<QueryGetKeygenResponse> {
      return queryService.keygen(request);
    },
    showObserverCount(request?: QueryShowObserverCountRequest): Promise<QueryShowObserverCountResponse> {
      return queryService.showObserverCount(request);
    },
    blameByIdentifier(request: QueryBlameByIdentifierRequest): Promise<QueryBlameByIdentifierResponse> {
      return queryService.blameByIdentifier(request);
    },
    getAllBlameRecords(request?: QueryAllBlameRecordsRequest): Promise<QueryAllBlameRecordsResponse> {
      return queryService.getAllBlameRecords(request);
    },
    blamesByChainAndNonce(request: QueryBlameByChainAndNonceRequest): Promise<QueryBlameByChainAndNonceResponse> {
      return queryService.blamesByChainAndNonce(request);
    },
    getAllBlockHeaders(request?: QueryAllBlockHeaderRequest): Promise<QueryAllBlockHeaderResponse> {
      return queryService.getAllBlockHeaders(request);
    },
    getBlockHeaderByHash(request: QueryGetBlockHeaderByHashRequest): Promise<QueryGetBlockHeaderByHashResponse> {
      return queryService.getBlockHeaderByHash(request);
    },
    getBlockHeaderStateByChain(request: QueryGetBlockHeaderStateRequest): Promise<QueryGetBlockHeaderStateResponse> {
      return queryService.getBlockHeaderStateByChain(request);
    },
    prove(request: QueryProveRequest): Promise<QueryProveResponse> {
      return queryService.prove(request);
    },
    getTssAddress(request: QueryGetTssAddressRequest): Promise<QueryGetTssAddressResponse> {
      return queryService.getTssAddress(request);
    },
    getTssAddressByFinalizedHeight(request: QueryGetTssAddressByFinalizedHeightRequest): Promise<QueryGetTssAddressByFinalizedHeightResponse> {
      return queryService.getTssAddressByFinalizedHeight(request);
    },
    tSS(request?: QueryGetTSSRequest): Promise<QueryGetTSSResponse> {
      return queryService.tSS(request);
    },
    tssHistory(request?: QueryTssHistoryRequest): Promise<QueryTssHistoryResponse> {
      return queryService.tssHistory(request);
    },
    pendingNoncesAll(request?: QueryAllPendingNoncesRequest): Promise<QueryAllPendingNoncesResponse> {
      return queryService.pendingNoncesAll(request);
    },
    pendingNoncesByChain(request: QueryPendingNoncesByChainRequest): Promise<QueryPendingNoncesByChainResponse> {
      return queryService.pendingNoncesByChain(request);
    },
    chainNonces(request: QueryGetChainNoncesRequest): Promise<QueryGetChainNoncesResponse> {
      return queryService.chainNonces(request);
    },
    chainNoncesAll(request?: QueryAllChainNoncesRequest): Promise<QueryAllChainNoncesResponse> {
      return queryService.chainNoncesAll(request);
    }
  };
};