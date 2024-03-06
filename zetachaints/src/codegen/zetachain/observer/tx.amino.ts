//@ts-nocheck
import { MsgAddObserver, MsgUpdateObserver, MsgUpdateChainParams, MsgRemoveChainParams, MsgAddBlameVote, MsgUpdateCrosschainFlags, MsgUpdateKeygen, MsgAddBlockHeader } from "./tx";
export const AminoConverter = {
  "/zetachain.zetacore.observer.MsgAddObserver": {
    aminoType: "/zetachain.zetacore.observer.MsgAddObserver",
    toAmino: MsgAddObserver.toAmino,
    fromAmino: MsgAddObserver.fromAmino
  },
  "/zetachain.zetacore.observer.MsgUpdateObserver": {
    aminoType: "/zetachain.zetacore.observer.MsgUpdateObserver",
    toAmino: MsgUpdateObserver.toAmino,
    fromAmino: MsgUpdateObserver.fromAmino
  },
  "/zetachain.zetacore.observer.MsgUpdateChainParams": {
    aminoType: "/zetachain.zetacore.observer.MsgUpdateChainParams",
    toAmino: MsgUpdateChainParams.toAmino,
    fromAmino: MsgUpdateChainParams.fromAmino
  },
  "/zetachain.zetacore.observer.MsgRemoveChainParams": {
    aminoType: "/zetachain.zetacore.observer.MsgRemoveChainParams",
    toAmino: MsgRemoveChainParams.toAmino,
    fromAmino: MsgRemoveChainParams.fromAmino
  },
  "/zetachain.zetacore.observer.MsgAddBlameVote": {
    aminoType: "/zetachain.zetacore.observer.MsgAddBlameVote",
    toAmino: MsgAddBlameVote.toAmino,
    fromAmino: MsgAddBlameVote.fromAmino
  },
  "/zetachain.zetacore.observer.MsgUpdateCrosschainFlags": {
    aminoType: "/zetachain.zetacore.observer.MsgUpdateCrosschainFlags",
    toAmino: MsgUpdateCrosschainFlags.toAmino,
    fromAmino: MsgUpdateCrosschainFlags.fromAmino
  },
  "/zetachain.zetacore.observer.MsgUpdateKeygen": {
    aminoType: "/zetachain.zetacore.observer.MsgUpdateKeygen",
    toAmino: MsgUpdateKeygen.toAmino,
    fromAmino: MsgUpdateKeygen.fromAmino
  },
  "/zetachain.zetacore.observer.MsgAddBlockHeader": {
    aminoType: "/zetachain.zetacore.observer.MsgAddBlockHeader",
    toAmino: MsgAddBlockHeader.toAmino,
    fromAmino: MsgAddBlockHeader.fromAmino
  }
};