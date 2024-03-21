import { GeneratedType, Registry, OfflineSigner } from "@cosmjs/proto-signing";
import { defaultRegistryTypes, AminoTypes, SigningStargateClient } from "@cosmjs/stargate";
import { HttpEndpoint } from "@cosmjs/tendermint-rpc";
import * as zetachainAuthorityTxRegistry from "./authority/tx.registry";
import * as zetachainCrosschainTxRegistry from "./crosschain/tx.registry";
import * as zetachainEmissionsTxRegistry from "./emissions/tx.registry";
import * as zetachainFungibleTxRegistry from "./fungible/tx.registry";
import * as zetachainObserverTxRegistry from "./observer/tx.registry";
import * as zetachainAuthorityTxAmino from "./authority/tx.amino";
import * as zetachainCrosschainTxAmino from "./crosschain/tx.amino";
import * as zetachainEmissionsTxAmino from "./emissions/tx.amino";
import * as zetachainFungibleTxAmino from "./fungible/tx.amino";
import * as zetachainObserverTxAmino from "./observer/tx.amino";
export const zetachainAminoConverters = {
  ...zetachainAuthorityTxAmino.AminoConverter,
  ...zetachainCrosschainTxAmino.AminoConverter,
  ...zetachainEmissionsTxAmino.AminoConverter,
  ...zetachainFungibleTxAmino.AminoConverter,
  ...zetachainObserverTxAmino.AminoConverter
};
export const zetachainProtoRegistry: ReadonlyArray<[string, GeneratedType]> = [...zetachainAuthorityTxRegistry.registry, ...zetachainCrosschainTxRegistry.registry, ...zetachainEmissionsTxRegistry.registry, ...zetachainFungibleTxRegistry.registry, ...zetachainObserverTxRegistry.registry];
export const getSigningZetachainClientOptions = ({
  defaultTypes = defaultRegistryTypes
}: {
  defaultTypes?: ReadonlyArray<[string, GeneratedType]>;
} = {}): {
  registry: Registry;
  aminoTypes: AminoTypes;
} => {
  const registry = new Registry([...defaultTypes, ...zetachainProtoRegistry]);
  const aminoTypes = new AminoTypes({
    ...zetachainAminoConverters
  });
  return {
    registry,
    aminoTypes
  };
};
export const getSigningZetachainClient = async ({
  rpcEndpoint,
  signer,
  defaultTypes = defaultRegistryTypes
}: {
  rpcEndpoint: string | HttpEndpoint;
  signer: OfflineSigner;
  defaultTypes?: ReadonlyArray<[string, GeneratedType]>;
}) => {
  const {
    registry,
    aminoTypes
  } = getSigningZetachainClientOptions({
    defaultTypes
  });
  const client = await SigningStargateClient.connectWithSigner(rpcEndpoint, signer, {
    registry: (registry as any),
    aminoTypes
  });
  return client;
};