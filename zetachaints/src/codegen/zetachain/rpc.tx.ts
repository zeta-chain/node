import { Rpc } from "../helpers";
export const createRPCMsgClient = async ({
  rpc
}: {
  rpc: Rpc;
}) => ({
  zetachain: {
    zetacore: {
      authority: new (await import("./authority/tx.rpc.msg")).MsgClientImpl(rpc),
      crosschain: new (await import("./crosschain/tx.rpc.msg")).MsgClientImpl(rpc),
      emissions: new (await import("./emissions/tx.rpc.msg")).MsgClientImpl(rpc),
      fungible: new (await import("./fungible/tx.rpc.msg")).MsgClientImpl(rpc),
      observer: new (await import("./observer/tx.rpc.msg")).MsgClientImpl(rpc)
    }
  }
});