//@ts-nocheck
import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgWithdrawEmission } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/zetachain.zetacore.emissions.MsgWithdrawEmission", MsgWithdrawEmission]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    withdrawEmission(value: MsgWithdrawEmission) {
      return {
        typeUrl: "/zetachain.zetacore.emissions.MsgWithdrawEmission",
        value: MsgWithdrawEmission.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    withdrawEmission(value: MsgWithdrawEmission) {
      return {
        typeUrl: "/zetachain.zetacore.emissions.MsgWithdrawEmission",
        value
      };
    }
  },
  fromPartial: {
    withdrawEmission(value: MsgWithdrawEmission) {
      return {
        typeUrl: "/zetachain.zetacore.emissions.MsgWithdrawEmission",
        value: MsgWithdrawEmission.fromPartial(value)
      };
    }
  }
};