//@ts-nocheck
import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgUpdatePolicies } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/zetachain.zetacore.authority.MsgUpdatePolicies", MsgUpdatePolicies]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    updatePolicies(value: MsgUpdatePolicies) {
      return {
        typeUrl: "/zetachain.zetacore.authority.MsgUpdatePolicies",
        value: MsgUpdatePolicies.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    updatePolicies(value: MsgUpdatePolicies) {
      return {
        typeUrl: "/zetachain.zetacore.authority.MsgUpdatePolicies",
        value
      };
    }
  },
  fromPartial: {
    updatePolicies(value: MsgUpdatePolicies) {
      return {
        typeUrl: "/zetachain.zetacore.authority.MsgUpdatePolicies",
        value: MsgUpdatePolicies.fromPartial(value)
      };
    }
  }
};