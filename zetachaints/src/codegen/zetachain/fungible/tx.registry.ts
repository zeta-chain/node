//@ts-nocheck
import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgDeploySystemContracts, MsgDeployFungibleCoinZRC20, MsgRemoveForeignCoin, MsgUpdateSystemContract, MsgUpdateContractBytecode, MsgUpdateZRC20WithdrawFee, MsgUpdateZRC20PausedStatus, MsgUpdateZRC20LiquidityCap } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/zetachain.zetacore.fungible.MsgDeploySystemContracts", MsgDeploySystemContracts], ["/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20", MsgDeployFungibleCoinZRC20], ["/zetachain.zetacore.fungible.MsgRemoveForeignCoin", MsgRemoveForeignCoin], ["/zetachain.zetacore.fungible.MsgUpdateSystemContract", MsgUpdateSystemContract], ["/zetachain.zetacore.fungible.MsgUpdateContractBytecode", MsgUpdateContractBytecode], ["/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee", MsgUpdateZRC20WithdrawFee], ["/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus", MsgUpdateZRC20PausedStatus], ["/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap", MsgUpdateZRC20LiquidityCap]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    deploySystemContracts(value: MsgDeploySystemContracts) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgDeploySystemContracts",
        value: MsgDeploySystemContracts.encode(value).finish()
      };
    },
    deployFungibleCoinZRC20(value: MsgDeployFungibleCoinZRC20) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20",
        value: MsgDeployFungibleCoinZRC20.encode(value).finish()
      };
    },
    removeForeignCoin(value: MsgRemoveForeignCoin) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgRemoveForeignCoin",
        value: MsgRemoveForeignCoin.encode(value).finish()
      };
    },
    updateSystemContract(value: MsgUpdateSystemContract) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateSystemContract",
        value: MsgUpdateSystemContract.encode(value).finish()
      };
    },
    updateContractBytecode(value: MsgUpdateContractBytecode) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateContractBytecode",
        value: MsgUpdateContractBytecode.encode(value).finish()
      };
    },
    updateZRC20WithdrawFee(value: MsgUpdateZRC20WithdrawFee) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee",
        value: MsgUpdateZRC20WithdrawFee.encode(value).finish()
      };
    },
    updateZRC20PausedStatus(value: MsgUpdateZRC20PausedStatus) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus",
        value: MsgUpdateZRC20PausedStatus.encode(value).finish()
      };
    },
    updateZRC20LiquidityCap(value: MsgUpdateZRC20LiquidityCap) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap",
        value: MsgUpdateZRC20LiquidityCap.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    deploySystemContracts(value: MsgDeploySystemContracts) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgDeploySystemContracts",
        value
      };
    },
    deployFungibleCoinZRC20(value: MsgDeployFungibleCoinZRC20) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20",
        value
      };
    },
    removeForeignCoin(value: MsgRemoveForeignCoin) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgRemoveForeignCoin",
        value
      };
    },
    updateSystemContract(value: MsgUpdateSystemContract) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateSystemContract",
        value
      };
    },
    updateContractBytecode(value: MsgUpdateContractBytecode) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateContractBytecode",
        value
      };
    },
    updateZRC20WithdrawFee(value: MsgUpdateZRC20WithdrawFee) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee",
        value
      };
    },
    updateZRC20PausedStatus(value: MsgUpdateZRC20PausedStatus) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus",
        value
      };
    },
    updateZRC20LiquidityCap(value: MsgUpdateZRC20LiquidityCap) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap",
        value
      };
    }
  },
  fromPartial: {
    deploySystemContracts(value: MsgDeploySystemContracts) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgDeploySystemContracts",
        value: MsgDeploySystemContracts.fromPartial(value)
      };
    },
    deployFungibleCoinZRC20(value: MsgDeployFungibleCoinZRC20) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20",
        value: MsgDeployFungibleCoinZRC20.fromPartial(value)
      };
    },
    removeForeignCoin(value: MsgRemoveForeignCoin) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgRemoveForeignCoin",
        value: MsgRemoveForeignCoin.fromPartial(value)
      };
    },
    updateSystemContract(value: MsgUpdateSystemContract) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateSystemContract",
        value: MsgUpdateSystemContract.fromPartial(value)
      };
    },
    updateContractBytecode(value: MsgUpdateContractBytecode) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateContractBytecode",
        value: MsgUpdateContractBytecode.fromPartial(value)
      };
    },
    updateZRC20WithdrawFee(value: MsgUpdateZRC20WithdrawFee) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee",
        value: MsgUpdateZRC20WithdrawFee.fromPartial(value)
      };
    },
    updateZRC20PausedStatus(value: MsgUpdateZRC20PausedStatus) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus",
        value: MsgUpdateZRC20PausedStatus.fromPartial(value)
      };
    },
    updateZRC20LiquidityCap(value: MsgUpdateZRC20LiquidityCap) {
      return {
        typeUrl: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap",
        value: MsgUpdateZRC20LiquidityCap.fromPartial(value)
      };
    }
  }
};