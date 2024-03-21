import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryListPoolAddressesRequest, QueryListPoolAddressesResponse, QueryGetEmissionsFactorsRequest, QueryGetEmissionsFactorsResponse, QueryShowAvailableEmissionsRequest, QueryShowAvailableEmissionsResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a list of ListBalances items. */
  listPoolAddresses(request?: QueryListPoolAddressesRequest): Promise<QueryListPoolAddressesResponse>;
  /** Queries a list of GetEmmisonsFactors items. */
  getEmissionsFactors(request?: QueryGetEmissionsFactorsRequest): Promise<QueryGetEmissionsFactorsResponse>;
  /** Queries a list of ShowAvailableEmissions items. */
  showAvailableEmissions(request: QueryShowAvailableEmissionsRequest): Promise<QueryShowAvailableEmissionsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.params = this.params.bind(this);
    this.listPoolAddresses = this.listPoolAddresses.bind(this);
    this.getEmissionsFactors = this.getEmissionsFactors.bind(this);
    this.showAvailableEmissions = this.showAvailableEmissions.bind(this);
  }
  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.emissions.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new BinaryReader(data)));
  }
  listPoolAddresses(request: QueryListPoolAddressesRequest = {}): Promise<QueryListPoolAddressesResponse> {
    const data = QueryListPoolAddressesRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.emissions.Query", "ListPoolAddresses", data);
    return promise.then(data => QueryListPoolAddressesResponse.decode(new BinaryReader(data)));
  }
  getEmissionsFactors(request: QueryGetEmissionsFactorsRequest = {}): Promise<QueryGetEmissionsFactorsResponse> {
    const data = QueryGetEmissionsFactorsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.emissions.Query", "GetEmissionsFactors", data);
    return promise.then(data => QueryGetEmissionsFactorsResponse.decode(new BinaryReader(data)));
  }
  showAvailableEmissions(request: QueryShowAvailableEmissionsRequest): Promise<QueryShowAvailableEmissionsResponse> {
    const data = QueryShowAvailableEmissionsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.emissions.Query", "ShowAvailableEmissions", data);
    return promise.then(data => QueryShowAvailableEmissionsResponse.decode(new BinaryReader(data)));
  }
}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },
    listPoolAddresses(request?: QueryListPoolAddressesRequest): Promise<QueryListPoolAddressesResponse> {
      return queryService.listPoolAddresses(request);
    },
    getEmissionsFactors(request?: QueryGetEmissionsFactorsRequest): Promise<QueryGetEmissionsFactorsResponse> {
      return queryService.getEmissionsFactors(request);
    },
    showAvailableEmissions(request: QueryShowAvailableEmissionsRequest): Promise<QueryShowAvailableEmissionsResponse> {
      return queryService.showAvailableEmissions(request);
    }
  };
};