//@ts-nocheck
import { MsgUpdatePolicies } from "./tx";
export const AminoConverter = {
  "/zetachain.zetacore.authority.MsgUpdatePolicies": {
    aminoType: "/zetachain.zetacore.authority.MsgUpdatePolicies",
    toAmino: MsgUpdatePolicies.toAmino,
    fromAmino: MsgUpdatePolicies.fromAmino
  }
};