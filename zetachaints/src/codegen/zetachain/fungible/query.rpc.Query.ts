import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryGetForeignCoinsRequest, QueryGetForeignCoinsResponse, QueryAllForeignCoinsRequest, QueryAllForeignCoinsResponse, QueryGetSystemContractRequest, QueryGetSystemContractResponse, QueryGetGasStabilityPoolAddress, QueryGetGasStabilityPoolAddressResponse, QueryGetGasStabilityPoolBalance, QueryGetGasStabilityPoolBalanceResponse, QueryAllGasStabilityPoolBalance, QueryAllGasStabilityPoolBalanceResponse, QueryCodeHashRequest, QueryCodeHashResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a ForeignCoins by index. */
  foreignCoins(request: QueryGetForeignCoinsRequest): Promise<QueryGetForeignCoinsResponse>;
  /** Queries a list of ForeignCoins items. */
  foreignCoinsAll(request?: QueryAllForeignCoinsRequest): Promise<QueryAllForeignCoinsResponse>;
  /** Queries SystemContract */
  systemContract(request?: QueryGetSystemContractRequest): Promise<QueryGetSystemContractResponse>;
  /** Queries the address of a gas stability pool on a given chain. */
  gasStabilityPoolAddress(request?: QueryGetGasStabilityPoolAddress): Promise<QueryGetGasStabilityPoolAddressResponse>;
  /** Queries the balance of a gas stability pool on a given chain. */
  gasStabilityPoolBalance(request: QueryGetGasStabilityPoolBalance): Promise<QueryGetGasStabilityPoolBalanceResponse>;
  /** Queries all gas stability pool balances. */
  gasStabilityPoolBalanceAll(request?: QueryAllGasStabilityPoolBalance): Promise<QueryAllGasStabilityPoolBalanceResponse>;
  /** Code hash query the code hash of a contract. */
  codeHash(request: QueryCodeHashRequest): Promise<QueryCodeHashResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.params = this.params.bind(this);
    this.foreignCoins = this.foreignCoins.bind(this);
    this.foreignCoinsAll = this.foreignCoinsAll.bind(this);
    this.systemContract = this.systemContract.bind(this);
    this.gasStabilityPoolAddress = this.gasStabilityPoolAddress.bind(this);
    this.gasStabilityPoolBalance = this.gasStabilityPoolBalance.bind(this);
    this.gasStabilityPoolBalanceAll = this.gasStabilityPoolBalanceAll.bind(this);
    this.codeHash = this.codeHash.bind(this);
  }
  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new BinaryReader(data)));
  }
  foreignCoins(request: QueryGetForeignCoinsRequest): Promise<QueryGetForeignCoinsResponse> {
    const data = QueryGetForeignCoinsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Query", "ForeignCoins", data);
    return promise.then(data => QueryGetForeignCoinsResponse.decode(new BinaryReader(data)));
  }
  foreignCoinsAll(request: QueryAllForeignCoinsRequest = {
    pagination: undefined
  }): Promise<QueryAllForeignCoinsResponse> {
    const data = QueryAllForeignCoinsRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Query", "ForeignCoinsAll", data);
    return promise.then(data => QueryAllForeignCoinsResponse.decode(new BinaryReader(data)));
  }
  systemContract(request: QueryGetSystemContractRequest = {}): Promise<QueryGetSystemContractResponse> {
    const data = QueryGetSystemContractRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Query", "SystemContract", data);
    return promise.then(data => QueryGetSystemContractResponse.decode(new BinaryReader(data)));
  }
  gasStabilityPoolAddress(request: QueryGetGasStabilityPoolAddress = {}): Promise<QueryGetGasStabilityPoolAddressResponse> {
    const data = QueryGetGasStabilityPoolAddress.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Query", "GasStabilityPoolAddress", data);
    return promise.then(data => QueryGetGasStabilityPoolAddressResponse.decode(new BinaryReader(data)));
  }
  gasStabilityPoolBalance(request: QueryGetGasStabilityPoolBalance): Promise<QueryGetGasStabilityPoolBalanceResponse> {
    const data = QueryGetGasStabilityPoolBalance.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Query", "GasStabilityPoolBalance", data);
    return promise.then(data => QueryGetGasStabilityPoolBalanceResponse.decode(new BinaryReader(data)));
  }
  gasStabilityPoolBalanceAll(request: QueryAllGasStabilityPoolBalance = {}): Promise<QueryAllGasStabilityPoolBalanceResponse> {
    const data = QueryAllGasStabilityPoolBalance.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Query", "GasStabilityPoolBalanceAll", data);
    return promise.then(data => QueryAllGasStabilityPoolBalanceResponse.decode(new BinaryReader(data)));
  }
  codeHash(request: QueryCodeHashRequest): Promise<QueryCodeHashResponse> {
    const data = QueryCodeHashRequest.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.fungible.Query", "CodeHash", data);
    return promise.then(data => QueryCodeHashResponse.decode(new BinaryReader(data)));
  }
}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },
    foreignCoins(request: QueryGetForeignCoinsRequest): Promise<QueryGetForeignCoinsResponse> {
      return queryService.foreignCoins(request);
    },
    foreignCoinsAll(request?: QueryAllForeignCoinsRequest): Promise<QueryAllForeignCoinsResponse> {
      return queryService.foreignCoinsAll(request);
    },
    systemContract(request?: QueryGetSystemContractRequest): Promise<QueryGetSystemContractResponse> {
      return queryService.systemContract(request);
    },
    gasStabilityPoolAddress(request?: QueryGetGasStabilityPoolAddress): Promise<QueryGetGasStabilityPoolAddressResponse> {
      return queryService.gasStabilityPoolAddress(request);
    },
    gasStabilityPoolBalance(request: QueryGetGasStabilityPoolBalance): Promise<QueryGetGasStabilityPoolBalanceResponse> {
      return queryService.gasStabilityPoolBalance(request);
    },
    gasStabilityPoolBalanceAll(request?: QueryAllGasStabilityPoolBalance): Promise<QueryAllGasStabilityPoolBalanceResponse> {
      return queryService.gasStabilityPoolBalanceAll(request);
    },
    codeHash(request: QueryCodeHashRequest): Promise<QueryCodeHashResponse> {
      return queryService.codeHash(request);
    }
  };
};