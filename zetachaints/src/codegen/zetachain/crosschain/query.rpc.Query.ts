import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryGetTssAddressRequest, QueryGetTssAddressResponse, QueryParamsRequest, QueryParamsResponse, QueryGetOutTxTrackerRequest, QueryGetOutTxTrackerResponse, QueryAllOutTxTrackerRequest, QueryAllOutTxTrackerResponse, QueryAllOutTxTrackerByChainRequest, QueryAllOutTxTrackerByChainResponse, QueryAllInTxTrackerByChainRequest, QueryAllInTxTrackerByChainResponse, QueryAllInTxTrackersRequest, QueryAllInTxTrackersResponse, QueryGetInTxHashToCctxRequest, QueryGetInTxHashToCctxResponse, QueryInTxHashToCctxDataRequest, QueryInTxHashToCctxDataResponse, QueryAllInTxHashToCctxRequest, QueryAllInTxHashToCctxResponse, QueryGetGasPriceRequest, QueryGetGasPriceResponse, QueryAllGasPriceRequest, QueryAllGasPriceResponse, QueryConvertGasToZetaRequest, QueryConvertGasToZetaResponse, QueryMessagePassingProtocolFeeRequest, QueryMessagePassingProtocolFeeResponse, QueryGetLastBlockHeightRequest, QueryGetLastBlockHeightResponse, QueryAllLastBlockHeightRequest, QueryAllLastBlockHeightResponse, QueryGetCctxRequest, QueryGetCctxResponse, QueryGetCctxByNonceRequest, QueryAllCctxRequest, QueryAllCctxResponse, QueryListCctxPendingRequest, QueryListCctxPendingResponse, QueryZetaAccountingRequest, QueryZetaAccountingResponse, QueryLastZetaHeightRequest, QueryLastZetaHeightResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
  /**
   * GetTssAddress queries the tss address of the module.
   * Deprecated: Moved to observer
   * TODO: remove after v12 once upgrade testing is no longer needed with v11
   * https://github.com/zeta-chain/node/issues/1547
   */
  getTssAddress(request?: QueryGetTssAddressRequest): Promise<QueryGetTssAddressResponse>;
  /** Parameters queries the parameters of the module. */
  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a OutTxTracker by index. */
  outTxTracker(request: QueryGetOutTxTrackerRequest): Promise<QueryGetOutTxTrackerResponse>;
  /** Queries a list of OutTxTracker items. */
  outTxTrackerAll(request?: QueryAllOutTxTrackerRequest): Promise<QueryAllOutTxTrackerResponse>;
  outTxTrackerAllByChain(request: QueryAllOutTxTrackerByChainRequest): Promise<QueryAllOutTxTrackerByChainResponse>;
  inTxTrackerAllByChain(request: QueryAllInTxTrackerByChainRequest): Promise<QueryAllInTxTrackerByChainResponse>;
  inTxTrackerAll(request?: QueryAllInTxTrackersRequest): Promise<QueryAllInTxTrackersResponse>;
  /** Queries a InTxHashToCctx by index. */
  inTxHashToCctx(request: QueryGetInTxHashToCctxRequest): Promise<QueryGetInTxHashToCctxResponse>;
  /** Queries a InTxHashToCctx data by index. */
  inTxHashToCctxData(request: QueryInTxHashToCctxDataRequest): Promise<QueryInTxHashToCctxDataResponse>;
  /** Queries a list of InTxHashToCctx items. */
  inTxHashToCctxAll(request?: QueryAllInTxHashToCctxRequest): Promise<QueryAllInTxHashToCctxResponse>;
  /** Queries a gasPrice by index. */
  gasPrice(request: QueryGetGasPriceRequest): Promise<QueryGetGasPriceResponse>;
  /** Queries a list of gasPrice items. */
  gasPriceAll(request?: QueryAllGasPriceRequest): Promise<QueryAllGasPriceResponse>;
  convertGasToZeta(request: QueryConvertGasToZetaRequest): Promise<QueryConvertGasToZetaResponse>;
  protocolFee(request?: QueryMessagePassingProtocolFeeRequest): Promise<QueryMessagePassingProtocolFeeResponse>;
  /** Queries a lastBlockHeight by index. */
  lastBlockHeight(request: QueryGetLastBlockHeightRequest): Promise<QueryGetLastBlockHeightResponse>;
  /** Queries a list of lastBlockHeight items. */
  lastBlockHeightAll(request?: QueryAllLastBlockHeightRequest): Promise<QueryAllLastBlockHeightResponse>;
  /** Queries a send by index. */
  cctx(request: QueryGetCctxRequest): Promise<QueryGetCctxResponse>;
  /** Queries a send by nonce. */
  cctxByNonce(request: QueryGetCctxByNonceRequest): Promise<QueryGetCctxResponse>;
  /** Queries a list of send items. */
  cctxAll(request?: QueryAllCctxRequest): Promise<QueryAllCctxResponse>;
  /** Queries a list of pending cctxs. */
  cctxListPending(request: QueryListCctxPendingRequest): Promise<QueryListCctxPendingResponse>;
  zetaAccounting(request?: QueryZetaAccountingRequest): Promise<QueryZetaAccountingResponse>;
  /** Queries a list of lastMetaHeight items. */
  lastZetaHeight(request?: QueryLastZetaHeightRequest): Promise<QueryLastZetaHeightResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.getTssAddress = this.getTssAddress.bind(this);
    this.params = this.params.bind(this);
    this.outTxTracker = this.outTxTracker.bind(this);
    this.outTxTrackerAll = this.outTxTrackerAll.bind(this);
    this.outTxTrackerAllByChain = this.outTxTrackerAllByChain.bind(this);
    this.inTxTrackerAllByChain = this.inTxTrackerAllByChain.bind(this);
    this.inTxTrackerAll = this.inTxTrackerAll.bind(this);
    this.inTxHashToCctx = this.inTxHashToCctx.bind(this);
    this.inTxHashToCctxData = this.inTxHashToCctxData.bind(this);
    this.inTxHashToCctxAll = this.inTxHashToCctxAll.bind(this);
    this.gasPrice = this.gasPrice.bind(this);
    this.gasPriceAll = this.gasPriceAll.bind(this);
    this.convertGasToZeta = this.convertGasToZeta.bind(this);
    this.protocolFee = this.protocolFee.bind(this);
    this.lastBlockHeight = this.lastBlockHeight.bind(this);
    this.lastBlockHeightAll = this.lastBlockHeightAll.bind(this);
    this.cctx = this.cctx.bind(this);
    this.cctxByNonce = this.cctxByNonce.bind(this);
    this.cctxAll = this.cctxAll.bind(this);
    this.cctxListPending = this.cctxListPending.bind(this);
    this.zetaAccounting = this.zetaAccounting.bind(this);
    this.lastZetaHeight = this.lastZetaHeight.bind(this);
  }
  getTssAddress(request: QueryGetTssAddressRequest = {}): Promise<QueryGetTssAddressResponse> {
    const data = QueryGetTssAddressRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "GetTssAddress", data);
    return promise.then(data => QueryGetTssAddressResponse.decode(new BinaryReader(data)));
  }
  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new BinaryReader(data)));
  }
  outTxTracker(request: QueryGetOutTxTrackerRequest): Promise<QueryGetOutTxTrackerResponse> {
    const data = QueryGetOutTxTrackerRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "OutTxTracker", data);
    return promise.then(data => QueryGetOutTxTrackerResponse.decode(new BinaryReader(data)));
  }
  outTxTrackerAll(request: QueryAllOutTxTrackerRequest = {
    pagination: undefined
  }): Promise<QueryAllOutTxTrackerResponse> {
    const data = QueryAllOutTxTrackerRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "OutTxTrackerAll", data);
    return promise.then(data => QueryAllOutTxTrackerResponse.decode(new BinaryReader(data)));
  }
  outTxTrackerAllByChain(request: QueryAllOutTxTrackerByChainRequest): Promise<QueryAllOutTxTrackerByChainResponse> {
    const data = QueryAllOutTxTrackerByChainRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "OutTxTrackerAllByChain", data);
    return promise.then(data => QueryAllOutTxTrackerByChainResponse.decode(new BinaryReader(data)));
  }
  inTxTrackerAllByChain(request: QueryAllInTxTrackerByChainRequest): Promise<QueryAllInTxTrackerByChainResponse> {
    const data = QueryAllInTxTrackerByChainRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "InTxTrackerAllByChain", data);
    return promise.then(data => QueryAllInTxTrackerByChainResponse.decode(new BinaryReader(data)));
  }
  inTxTrackerAll(request: QueryAllInTxTrackersRequest = {
    pagination: undefined
  }): Promise<QueryAllInTxTrackersResponse> {
    const data = QueryAllInTxTrackersRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "InTxTrackerAll", data);
    return promise.then(data => QueryAllInTxTrackersResponse.decode(new BinaryReader(data)));
  }
  inTxHashToCctx(request: QueryGetInTxHashToCctxRequest): Promise<QueryGetInTxHashToCctxResponse> {
    const data = QueryGetInTxHashToCctxRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "InTxHashToCctx", data);
    return promise.then(data => QueryGetInTxHashToCctxResponse.decode(new BinaryReader(data)));
  }
  inTxHashToCctxData(request: QueryInTxHashToCctxDataRequest): Promise<QueryInTxHashToCctxDataResponse> {
    const data = QueryInTxHashToCctxDataRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "InTxHashToCctxData", data);
    return promise.then(data => QueryInTxHashToCctxDataResponse.decode(new BinaryReader(data)));
  }
  inTxHashToCctxAll(request: QueryAllInTxHashToCctxRequest = {
    pagination: undefined
  }): Promise<QueryAllInTxHashToCctxResponse> {
    const data = QueryAllInTxHashToCctxRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "InTxHashToCctxAll", data);
    return promise.then(data => QueryAllInTxHashToCctxResponse.decode(new BinaryReader(data)));
  }
  gasPrice(request: QueryGetGasPriceRequest): Promise<QueryGetGasPriceResponse> {
    const data = QueryGetGasPriceRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "GasPrice", data);
    return promise.then(data => QueryGetGasPriceResponse.decode(new BinaryReader(data)));
  }
  gasPriceAll(request: QueryAllGasPriceRequest = {
    pagination: undefined
  }): Promise<QueryAllGasPriceResponse> {
    const data = QueryAllGasPriceRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "GasPriceAll", data);
    return promise.then(data => QueryAllGasPriceResponse.decode(new BinaryReader(data)));
  }
  convertGasToZeta(request: QueryConvertGasToZetaRequest): Promise<QueryConvertGasToZetaResponse> {
    const data = QueryConvertGasToZetaRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "ConvertGasToZeta", data);
    return promise.then(data => QueryConvertGasToZetaResponse.decode(new BinaryReader(data)));
  }
  protocolFee(request: QueryMessagePassingProtocolFeeRequest = {}): Promise<QueryMessagePassingProtocolFeeResponse> {
    const data = QueryMessagePassingProtocolFeeRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "ProtocolFee", data);
    return promise.then(data => QueryMessagePassingProtocolFeeResponse.decode(new BinaryReader(data)));
  }
  lastBlockHeight(request: QueryGetLastBlockHeightRequest): Promise<QueryGetLastBlockHeightResponse> {
    const data = QueryGetLastBlockHeightRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "LastBlockHeight", data);
    return promise.then(data => QueryGetLastBlockHeightResponse.decode(new BinaryReader(data)));
  }
  lastBlockHeightAll(request: QueryAllLastBlockHeightRequest = {
    pagination: undefined
  }): Promise<QueryAllLastBlockHeightResponse> {
    const data = QueryAllLastBlockHeightRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "LastBlockHeightAll", data);
    return promise.then(data => QueryAllLastBlockHeightResponse.decode(new BinaryReader(data)));
  }
  cctx(request: QueryGetCctxRequest): Promise<QueryGetCctxResponse> {
    const data = QueryGetCctxRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "Cctx", data);
    return promise.then(data => QueryGetCctxResponse.decode(new BinaryReader(data)));
  }
  cctxByNonce(request: QueryGetCctxByNonceRequest): Promise<QueryGetCctxResponse> {
    const data = QueryGetCctxByNonceRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "CctxByNonce", data);
    return promise.then(data => QueryGetCctxResponse.decode(new BinaryReader(data)));
  }
  cctxAll(request: QueryAllCctxRequest = {
    pagination: undefined
  }): Promise<QueryAllCctxResponse> {
    const data = QueryAllCctxRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "CctxAll", data);
    return promise.then(data => QueryAllCctxResponse.decode(new BinaryReader(data)));
  }
  cctxListPending(request: QueryListCctxPendingRequest): Promise<QueryListCctxPendingResponse> {
    const data = QueryListCctxPendingRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "CctxListPending", data);
    return promise.then(data => QueryListCctxPendingResponse.decode(new BinaryReader(data)));
  }
  zetaAccounting(request: QueryZetaAccountingRequest = {}): Promise<QueryZetaAccountingResponse> {
    const data = QueryZetaAccountingRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "ZetaAccounting", data);
    return promise.then(data => QueryZetaAccountingResponse.decode(new BinaryReader(data)));
  }
  lastZetaHeight(request: QueryLastZetaHeightRequest = {}): Promise<QueryLastZetaHeightResponse> {
    const data = QueryLastZetaHeightRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.crosschain.Query", "LastZetaHeight", data);
    return promise.then(data => QueryLastZetaHeightResponse.decode(new BinaryReader(data)));
  }
}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    getTssAddress(request?: QueryGetTssAddressRequest): Promise<QueryGetTssAddressResponse> {
      return queryService.getTssAddress(request);
    },
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },
    outTxTracker(request: QueryGetOutTxTrackerRequest): Promise<QueryGetOutTxTrackerResponse> {
      return queryService.outTxTracker(request);
    },
    outTxTrackerAll(request?: QueryAllOutTxTrackerRequest): Promise<QueryAllOutTxTrackerResponse> {
      return queryService.outTxTrackerAll(request);
    },
    outTxTrackerAllByChain(request: QueryAllOutTxTrackerByChainRequest): Promise<QueryAllOutTxTrackerByChainResponse> {
      return queryService.outTxTrackerAllByChain(request);
    },
    inTxTrackerAllByChain(request: QueryAllInTxTrackerByChainRequest): Promise<QueryAllInTxTrackerByChainResponse> {
      return queryService.inTxTrackerAllByChain(request);
    },
    inTxTrackerAll(request?: QueryAllInTxTrackersRequest): Promise<QueryAllInTxTrackersResponse> {
      return queryService.inTxTrackerAll(request);
    },
    inTxHashToCctx(request: QueryGetInTxHashToCctxRequest): Promise<QueryGetInTxHashToCctxResponse> {
      return queryService.inTxHashToCctx(request);
    },
    inTxHashToCctxData(request: QueryInTxHashToCctxDataRequest): Promise<QueryInTxHashToCctxDataResponse> {
      return queryService.inTxHashToCctxData(request);
    },
    inTxHashToCctxAll(request?: QueryAllInTxHashToCctxRequest): Promise<QueryAllInTxHashToCctxResponse> {
      return queryService.inTxHashToCctxAll(request);
    },
    gasPrice(request: QueryGetGasPriceRequest): Promise<QueryGetGasPriceResponse> {
      return queryService.gasPrice(request);
    },
    gasPriceAll(request?: QueryAllGasPriceRequest): Promise<QueryAllGasPriceResponse> {
      return queryService.gasPriceAll(request);
    },
    convertGasToZeta(request: QueryConvertGasToZetaRequest): Promise<QueryConvertGasToZetaResponse> {
      return queryService.convertGasToZeta(request);
    },
    protocolFee(request?: QueryMessagePassingProtocolFeeRequest): Promise<QueryMessagePassingProtocolFeeResponse> {
      return queryService.protocolFee(request);
    },
    lastBlockHeight(request: QueryGetLastBlockHeightRequest): Promise<QueryGetLastBlockHeightResponse> {
      return queryService.lastBlockHeight(request);
    },
    lastBlockHeightAll(request?: QueryAllLastBlockHeightRequest): Promise<QueryAllLastBlockHeightResponse> {
      return queryService.lastBlockHeightAll(request);
    },
    cctx(request: QueryGetCctxRequest): Promise<QueryGetCctxResponse> {
      return queryService.cctx(request);
    },
    cctxByNonce(request: QueryGetCctxByNonceRequest): Promise<QueryGetCctxResponse> {
      return queryService.cctxByNonce(request);
    },
    cctxAll(request?: QueryAllCctxRequest): Promise<QueryAllCctxResponse> {
      return queryService.cctxAll(request);
    },
    cctxListPending(request: QueryListCctxPendingRequest): Promise<QueryListCctxPendingResponse> {
      return queryService.cctxListPending(request);
    },
    zetaAccounting(request?: QueryZetaAccountingRequest): Promise<QueryZetaAccountingResponse> {
      return queryService.zetaAccounting(request);
    },
    lastZetaHeight(request?: QueryLastZetaHeightRequest): Promise<QueryLastZetaHeightResponse> {
      return queryService.lastZetaHeight(request);
    }
  };
};