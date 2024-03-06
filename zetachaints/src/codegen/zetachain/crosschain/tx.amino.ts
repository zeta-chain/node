//@ts-nocheck
import { MsgAddToOutTxTracker, MsgAddToInTxTracker, MsgRemoveFromOutTxTracker, MsgGasPriceVoter, MsgVoteOnObservedOutboundTx, MsgVoteOnObservedInboundTx, MsgWhitelistERC20, MsgUpdateTssAddress, MsgMigrateTssFunds, MsgCreateTSSVoter, MsgAbortStuckCCTX, MsgRefundAbortedCCTX } from "./tx";
export const AminoConverter = {
  "/zetachain.zetacore.crosschain.MsgAddToOutTxTracker": {
    aminoType: "/zetachain.zetacore.crosschain.MsgAddToOutTxTracker",
    toAmino: MsgAddToOutTxTracker.toAmino,
    fromAmino: MsgAddToOutTxTracker.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgAddToInTxTracker": {
    aminoType: "/zetachain.zetacore.crosschain.MsgAddToInTxTracker",
    toAmino: MsgAddToInTxTracker.toAmino,
    fromAmino: MsgAddToInTxTracker.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker": {
    aminoType: "/zetachain.zetacore.crosschain.MsgRemoveFromOutTxTracker",
    toAmino: MsgRemoveFromOutTxTracker.toAmino,
    fromAmino: MsgRemoveFromOutTxTracker.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgGasPriceVoter": {
    aminoType: "/zetachain.zetacore.crosschain.MsgGasPriceVoter",
    toAmino: MsgGasPriceVoter.toAmino,
    fromAmino: MsgGasPriceVoter.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx": {
    aminoType: "/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx",
    toAmino: MsgVoteOnObservedOutboundTx.toAmino,
    fromAmino: MsgVoteOnObservedOutboundTx.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx": {
    aminoType: "/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx",
    toAmino: MsgVoteOnObservedInboundTx.toAmino,
    fromAmino: MsgVoteOnObservedInboundTx.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgWhitelistERC20": {
    aminoType: "/zetachain.zetacore.crosschain.MsgWhitelistERC20",
    toAmino: MsgWhitelistERC20.toAmino,
    fromAmino: MsgWhitelistERC20.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgUpdateTssAddress": {
    aminoType: "/zetachain.zetacore.crosschain.MsgUpdateTssAddress",
    toAmino: MsgUpdateTssAddress.toAmino,
    fromAmino: MsgUpdateTssAddress.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgMigrateTssFunds": {
    aminoType: "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
    toAmino: MsgMigrateTssFunds.toAmino,
    fromAmino: MsgMigrateTssFunds.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgCreateTSSVoter": {
    aminoType: "/zetachain.zetacore.crosschain.MsgCreateTSSVoter",
    toAmino: MsgCreateTSSVoter.toAmino,
    fromAmino: MsgCreateTSSVoter.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX": {
    aminoType: "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX",
    toAmino: MsgAbortStuckCCTX.toAmino,
    fromAmino: MsgAbortStuckCCTX.fromAmino
  },
  "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX": {
    aminoType: "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX",
    toAmino: MsgRefundAbortedCCTX.toAmino,
    fromAmino: MsgRefundAbortedCCTX.fromAmino
  }
};