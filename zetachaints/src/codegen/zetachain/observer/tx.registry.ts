//@ts-nocheck
import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgAddObserver, MsgUpdateObserver, MsgUpdateChainParams, MsgRemoveChainParams, MsgAddBlameVote, MsgUpdateCrosschainFlags, MsgUpdateKeygen, MsgAddBlockHeader } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/zetachain.zetacore.observer.MsgAddObserver", MsgAddObserver], ["/zetachain.zetacore.observer.MsgUpdateObserver", MsgUpdateObserver], ["/zetachain.zetacore.observer.MsgUpdateChainParams", MsgUpdateChainParams], ["/zetachain.zetacore.observer.MsgRemoveChainParams", MsgRemoveChainParams], ["/zetachain.zetacore.observer.MsgAddBlameVote", MsgAddBlameVote], ["/zetachain.zetacore.observer.MsgUpdateCrosschainFlags", MsgUpdateCrosschainFlags], ["/zetachain.zetacore.observer.MsgUpdateKeygen", MsgUpdateKeygen], ["/zetachain.zetacore.observer.MsgAddBlockHeader", MsgAddBlockHeader]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    addObserver(value: MsgAddObserver) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgAddObserver",
        value: MsgAddObserver.encode(value).finish()
      };
    },
    updateObserver(value: MsgUpdateObserver) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateObserver",
        value: MsgUpdateObserver.encode(value).finish()
      };
    },
    updateChainParams(value: MsgUpdateChainParams) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateChainParams",
        value: MsgUpdateChainParams.encode(value).finish()
      };
    },
    removeChainParams(value: MsgRemoveChainParams) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgRemoveChainParams",
        value: MsgRemoveChainParams.encode(value).finish()
      };
    },
    addBlameVote(value: MsgAddBlameVote) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgAddBlameVote",
        value: MsgAddBlameVote.encode(value).finish()
      };
    },
    updateCrosschainFlags(value: MsgUpdateCrosschainFlags) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlags",
        value: MsgUpdateCrosschainFlags.encode(value).finish()
      };
    },
    updateKeygen(value: MsgUpdateKeygen) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateKeygen",
        value: MsgUpdateKeygen.encode(value).finish()
      };
    },
    addBlockHeader(value: MsgAddBlockHeader) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgAddBlockHeader",
        value: MsgAddBlockHeader.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    addObserver(value: MsgAddObserver) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgAddObserver",
        value
      };
    },
    updateObserver(value: MsgUpdateObserver) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateObserver",
        value
      };
    },
    updateChainParams(value: MsgUpdateChainParams) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateChainParams",
        value
      };
    },
    removeChainParams(value: MsgRemoveChainParams) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgRemoveChainParams",
        value
      };
    },
    addBlameVote(value: MsgAddBlameVote) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgAddBlameVote",
        value
      };
    },
    updateCrosschainFlags(value: MsgUpdateCrosschainFlags) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlags",
        value
      };
    },
    updateKeygen(value: MsgUpdateKeygen) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateKeygen",
        value
      };
    },
    addBlockHeader(value: MsgAddBlockHeader) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgAddBlockHeader",
        value
      };
    }
  },
  fromPartial: {
    addObserver(value: MsgAddObserver) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgAddObserver",
        value: MsgAddObserver.fromPartial(value)
      };
    },
    updateObserver(value: MsgUpdateObserver) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateObserver",
        value: MsgUpdateObserver.fromPartial(value)
      };
    },
    updateChainParams(value: MsgUpdateChainParams) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateChainParams",
        value: MsgUpdateChainParams.fromPartial(value)
      };
    },
    removeChainParams(value: MsgRemoveChainParams) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgRemoveChainParams",
        value: MsgRemoveChainParams.fromPartial(value)
      };
    },
    addBlameVote(value: MsgAddBlameVote) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgAddBlameVote",
        value: MsgAddBlameVote.fromPartial(value)
      };
    },
    updateCrosschainFlags(value: MsgUpdateCrosschainFlags) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlags",
        value: MsgUpdateCrosschainFlags.fromPartial(value)
      };
    },
    updateKeygen(value: MsgUpdateKeygen) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgUpdateKeygen",
        value: MsgUpdateKeygen.fromPartial(value)
      };
    },
    addBlockHeader(value: MsgAddBlockHeader) {
      return {
        typeUrl: "/zetachain.zetacore.observer.MsgAddBlockHeader",
        value: MsgAddBlockHeader.fromPartial(value)
      };
    }
  }
};