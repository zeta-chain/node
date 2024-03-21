//@ts-nocheck
import { MsgDeploySystemContracts, MsgDeployFungibleCoinZRC20, MsgRemoveForeignCoin, MsgUpdateSystemContract, MsgUpdateContractBytecode, MsgUpdateZRC20WithdrawFee, MsgUpdateZRC20PausedStatus, MsgUpdateZRC20LiquidityCap } from "./tx";
export const AminoConverter = {
  "/zetachain.zetacore.fungible.MsgDeploySystemContracts": {
    aminoType: "/zetachain.zetacore.fungible.MsgDeploySystemContracts",
    toAmino: MsgDeploySystemContracts.toAmino,
    fromAmino: MsgDeploySystemContracts.fromAmino
  },
  "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20": {
    aminoType: "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20",
    toAmino: MsgDeployFungibleCoinZRC20.toAmino,
    fromAmino: MsgDeployFungibleCoinZRC20.fromAmino
  },
  "/zetachain.zetacore.fungible.MsgRemoveForeignCoin": {
    aminoType: "/zetachain.zetacore.fungible.MsgRemoveForeignCoin",
    toAmino: MsgRemoveForeignCoin.toAmino,
    fromAmino: MsgRemoveForeignCoin.fromAmino
  },
  "/zetachain.zetacore.fungible.MsgUpdateSystemContract": {
    aminoType: "/zetachain.zetacore.fungible.MsgUpdateSystemContract",
    toAmino: MsgUpdateSystemContract.toAmino,
    fromAmino: MsgUpdateSystemContract.fromAmino
  },
  "/zetachain.zetacore.fungible.MsgUpdateContractBytecode": {
    aminoType: "/zetachain.zetacore.fungible.MsgUpdateContractBytecode",
    toAmino: MsgUpdateContractBytecode.toAmino,
    fromAmino: MsgUpdateContractBytecode.fromAmino
  },
  "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee": {
    aminoType: "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee",
    toAmino: MsgUpdateZRC20WithdrawFee.toAmino,
    fromAmino: MsgUpdateZRC20WithdrawFee.fromAmino
  },
  "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus": {
    aminoType: "/zetachain.zetacore.fungible.MsgUpdateZRC20PausedStatus",
    toAmino: MsgUpdateZRC20PausedStatus.toAmino,
    fromAmino: MsgUpdateZRC20PausedStatus.fromAmino
  },
  "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap": {
    aminoType: "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap",
    toAmino: MsgUpdateZRC20LiquidityCap.toAmino,
    fromAmino: MsgUpdateZRC20LiquidityCap.fromAmino
  }
};