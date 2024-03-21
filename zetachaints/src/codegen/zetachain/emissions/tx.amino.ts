//@ts-nocheck
import { MsgWithdrawEmission } from "./tx";
export const AminoConverter = {
  "/zetachain.zetacore.emissions.MsgWithdrawEmission": {
    aminoType: "/zetachain.zetacore.emissions.MsgWithdrawEmission",
    toAmino: MsgWithdrawEmission.toAmino,
    fromAmino: MsgWithdrawEmission.fromAmino
  }
};