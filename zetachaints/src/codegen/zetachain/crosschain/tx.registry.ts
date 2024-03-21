//@ts-nocheck
import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgAddToOutTxTracker, MsgAddToInTxTracker, MsgRemoveFromOutTxTracker, MsgGasPriceVoter, MsgVoteOnObservedOutboundTx, MsgVoteOnObservedInboundTx, MsgWhitelistERC20, MsgUpdateTssAddress, MsgMigrateTssFunds, MsgCreateTSSVoter, MsgAbortStuckCCTX, MsgRefundAbortedCCTX } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/zetachain.zetacore.crosschain.MsgAddToOutTxTracker", MsgAddToOutTxTracker], ["/zetachain.zetacore.crosschain.MsgAddToInTxTracker", MsgAddToInTxTracker], ["/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker", MsgRemoveFromOutTxTracker], ["/zetachain.zetacore.crosschain.MsgGasPriceVoter", MsgGasPriceVoter], ["/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx", MsgVoteOnObservedOutboundTx], ["/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx", MsgVoteOnObservedInboundTx], ["/zetachain.zetacore.crosschain.MsgWhitelistERC20", MsgWhitelistERC20], ["/zetachain.zetacore.crosschain.MsgUpdateTssAddress", MsgUpdateTssAddress], ["/zetachain.zetacore.crosschain.MsgMigrateTssFunds", MsgMigrateTssFunds], ["/zetachain.zetacore.crosschain.MsgCreateTSSVoter", MsgCreateTSSVoter], ["/zetachain.zetacore.crosschain.MsgAbortStuckCCTX", MsgAbortStuckCCTX], ["/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX", MsgRefundAbortedCCTX]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    addToOutTxTracker(value: MsgAddToOutTxTracker) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgAddToOutTxTracker",
        value: MsgAddToOutTxTracker.encode(value).finish()
      };
    },
    addToInTxTracker(value: MsgAddToInTxTracker) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgAddToInTxTracker",
        value: MsgAddToInTxTracker.encode(value).finish()
      };
    },
    removeFromOutTxTracker(value: MsgRemoveFromOutTxTracker) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker",
        value: MsgRemoveFromOutTxTracker.encode(value).finish()
      };
    },
    gasPriceVoter(value: MsgGasPriceVoter) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgGasPriceVoter",
        value: MsgGasPriceVoter.encode(value).finish()
      };
    },
    voteOnObservedOutboundTx(value: MsgVoteOnObservedOutboundTx) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx",
        value: MsgVoteOnObservedOutboundTx.encode(value).finish()
      };
    },
    voteOnObservedInboundTx(value: MsgVoteOnObservedInboundTx) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx",
        value: MsgVoteOnObservedInboundTx.encode(value).finish()
      };
    },
    whitelistERC20(value: MsgWhitelistERC20) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgWhitelistERC20",
        value: MsgWhitelistERC20.encode(value).finish()
      };
    },
    updateTssAddress(value: MsgUpdateTssAddress) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgUpdateTssAddress",
        value: MsgUpdateTssAddress.encode(value).finish()
      };
    },
    migrateTssFunds(value: MsgMigrateTssFunds) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
        value: MsgMigrateTssFunds.encode(value).finish()
      };
    },
    createTSSVoter(value: MsgCreateTSSVoter) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgCreateTSSVoter",
        value: MsgCreateTSSVoter.encode(value).finish()
      };
    },
    abortStuckCCTX(value: MsgAbortStuckCCTX) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX",
        value: MsgAbortStuckCCTX.encode(value).finish()
      };
    },
    refundAbortedCCTX(value: MsgRefundAbortedCCTX) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX",
        value: MsgRefundAbortedCCTX.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    addToOutTxTracker(value: MsgAddToOutTxTracker) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgAddToOutTxTracker",
        value
      };
    },
    addToInTxTracker(value: MsgAddToInTxTracker) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgAddToInTxTracker",
        value
      };
    },
    removeFromOutTxTracker(value: MsgRemoveFromOutTxTracker) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker",
        value
      };
    },
    gasPriceVoter(value: MsgGasPriceVoter) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgGasPriceVoter",
        value
      };
    },
    voteOnObservedOutboundTx(value: MsgVoteOnObservedOutboundTx) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx",
        value
      };
    },
    voteOnObservedInboundTx(value: MsgVoteOnObservedInboundTx) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx",
        value
      };
    },
    whitelistERC20(value: MsgWhitelistERC20) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgWhitelistERC20",
        value
      };
    },
    updateTssAddress(value: MsgUpdateTssAddress) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgUpdateTssAddress",
        value
      };
    },
    migrateTssFunds(value: MsgMigrateTssFunds) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
        value
      };
    },
    createTSSVoter(value: MsgCreateTSSVoter) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgCreateTSSVoter",
        value
      };
    },
    abortStuckCCTX(value: MsgAbortStuckCCTX) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX",
        value
      };
    },
    refundAbortedCCTX(value: MsgRefundAbortedCCTX) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX",
        value
      };
    }
  },
  fromPartial: {
    addToOutTxTracker(value: MsgAddToOutTxTracker) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgAddToOutTxTracker",
        value: MsgAddToOutTxTracker.fromPartial(value)
      };
    },
    addToInTxTracker(value: MsgAddToInTxTracker) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgAddToInTxTracker",
        value: MsgAddToInTxTracker.fromPartial(value)
      };
    },
    removeFromOutTxTracker(value: MsgRemoveFromOutTxTracker) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker",
        value: MsgRemoveFromOutTxTracker.fromPartial(value)
      };
    },
    gasPriceVoter(value: MsgGasPriceVoter) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgGasPriceVoter",
        value: MsgGasPriceVoter.fromPartial(value)
      };
    },
    voteOnObservedOutboundTx(value: MsgVoteOnObservedOutboundTx) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx",
        value: MsgVoteOnObservedOutboundTx.fromPartial(value)
      };
    },
    voteOnObservedInboundTx(value: MsgVoteOnObservedInboundTx) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx",
        value: MsgVoteOnObservedInboundTx.fromPartial(value)
      };
    },
    whitelistERC20(value: MsgWhitelistERC20) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgWhitelistERC20",
        value: MsgWhitelistERC20.fromPartial(value)
      };
    },
    updateTssAddress(value: MsgUpdateTssAddress) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgUpdateTssAddress",
        value: MsgUpdateTssAddress.fromPartial(value)
      };
    },
    migrateTssFunds(value: MsgMigrateTssFunds) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
        value: MsgMigrateTssFunds.fromPartial(value)
      };
    },
    createTSSVoter(value: MsgCreateTSSVoter) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgCreateTSSVoter",
        value: MsgCreateTSSVoter.fromPartial(value)
      };
    },
    abortStuckCCTX(value: MsgAbortStuckCCTX) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX",
        value: MsgAbortStuckCCTX.fromPartial(value)
      };
    },
    refundAbortedCCTX(value: MsgRefundAbortedCCTX) {
      return {
        typeUrl: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX",
        value: MsgRefundAbortedCCTX.fromPartial(value)
      };
    }
  }
};