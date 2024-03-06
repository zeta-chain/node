import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { MsgWithdrawEmission, MsgWithdrawEmissionResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
  withdrawEmission(request: MsgWithdrawEmission): Promise<MsgWithdrawEmissionResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.withdrawEmission = this.withdrawEmission.bind(this);
  }
  withdrawEmission(request: MsgWithdrawEmission): Promise<MsgWithdrawEmissionResponse> {
    const data = MsgWithdrawEmission.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.emissions.Msg", "WithdrawEmission", data);
    return promise.then(data => MsgWithdrawEmissionResponse.decode(new BinaryReader(data)));
  }
}