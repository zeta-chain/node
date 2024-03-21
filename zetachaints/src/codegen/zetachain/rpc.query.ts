import { Tendermint34Client, HttpEndpoint } from "@cosmjs/tendermint-rpc";
import { QueryClient } from "@cosmjs/stargate";
export const createRPCQueryClient = async ({
  rpcEndpoint
}: {
  rpcEndpoint: string | HttpEndpoint;
}) => {
  const tmClient = await Tendermint34Client.connect(rpcEndpoint);
  const client = new QueryClient(tmClient);
  return {
    zetachain: {
      zetacore: {
        authority: (await import("./authority/query.rpc.Query")).createRpcQueryExtension(client),
        crosschain: (await import("./crosschain/query.rpc.Query")).createRpcQueryExtension(client),
        emissions: (await import("./emissions/query.rpc.Query")).createRpcQueryExtension(client),
        fungible: (await import("./fungible/query.rpc.Query")).createRpcQueryExtension(client),
        observer: (await import("./observer/query.rpc.Query")).createRpcQueryExtension(client)
      }
    }
  };
};