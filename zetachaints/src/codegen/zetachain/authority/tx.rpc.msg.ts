import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { MsgUpdatePolicies, MsgUpdatePoliciesResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
  updatePolicies(request: MsgUpdatePolicies): Promise<MsgUpdatePoliciesResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.updatePolicies = this.updatePolicies.bind(this);
  }
  updatePolicies(request: MsgUpdatePolicies): Promise<MsgUpdatePoliciesResponse> {
    const data = MsgUpdatePolicies.encode(request).finish();
    const promise = this.rpc.request("zetachain.zetacore.authority.Msg", "UpdatePolicies", data);
    return promise.then(data => MsgUpdatePoliciesResponse.decode(new BinaryReader(data)));
  }
}