import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryGetPoliciesRequest, QueryGetPoliciesResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
  /** Queries Policies */
  policies(request?: QueryGetPoliciesRequest): Promise<QueryGetPoliciesResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.policies = this.policies.bind(this);
  }
  policies(request: QueryGetPoliciesRequest = {}): Promise<QueryGetPoliciesResponse> {
    const data = QueryGetPoliciesRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.authority.Query", "Policies", data);
    return promise.then(data => QueryGetPoliciesResponse.decode(new BinaryReader(data)));
  }
}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    policies(request?: QueryGetPoliciesRequest): Promise<QueryGetPoliciesResponse> {
      return queryService.policies(request);
    }
  };
};